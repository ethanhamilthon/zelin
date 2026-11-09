package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
	"zelin/internal/agent"
	"zelin/internal/files"
	"zelin/internal/instructions"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/tagged"
	"zelin/internal/tasks"
)

type AgentService struct {
	mu           sync.Mutex
	agents       map[string]*agent.Agent
	runs         map[string]context.CancelFunc
	store        *sessions.Store
	client       *provider.Client
	instructions *instructions.Store
	emit         func(agent.Event)
	notify       func(name string, data any)
	steerRefs    map[string][]tagged.Ref
	live         map[string]*RunSnapshot
	activity     map[string]SessionActivity
	tasks        *tasks.Manager
	paused       map[string]bool
}

type AgentResult struct {
	Files      []files.Attachment `json:"files,omitempty"`
	Content    string             `json:"content"`
	Executions []agent.Execution  `json:"executions"`
	Turn       *sessions.Turn     `json:"turn,omitempty"`
	Error      string             `json:"error,omitempty"`
}

func (s *AgentService) start(ctx context.Context, id string) (*agent.Agent, context.Context, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runs == nil {
		s.runs = map[string]context.CancelFunc{}
	}
	if s.runs[id] != nil {
		return nil, nil, fmt.Errorf("Agent is already running in this session")
	}
	current := s.agents[id]
	if current == nil {
		record, err := s.store.Get(id)
		if err != nil {
			return nil, nil, err
		}
		current = agent.NewIn(sessionClient{s, id}, s.workspaceDir(record.WorkspaceID))
		current.Restore(record.State)
		s.agents[id] = current
	}
	current.Extra = func() (string, error) { return s.workspaceInstructions(id) }
	current.OutputLimit = s.toolOutputLimit
	current.UtilityModel = s.utilityModel
	current.UtilityWindow = s.utilityWindow
	if s.instructions != nil {
		current.Prompt = func() string { text, _ := s.instructions.Read("system"); return text }
	}
	ctx, cancel := context.WithCancel(context.WithoutCancel(ctx))
	s.runs[id] = cancel
	return current, ctx, nil
}

func (s *AgentService) run(ctx context.Context, id string, history []provider.Message, refs []tagged.Ref, claim func() (string, error)) (AgentResult, error) {
	current, ctx, err := s.start(ctx, id)
	if err != nil {
		return AgentResult{}, err
	}
	defer s.finish(id)
	if history, err = claimed(history, claim); err != nil {
		return AgentResult{}, err
	}
	s.beginLive(id)
	defer s.endLive(id)
	if s.client != nil {
		if err := s.pin(id); err != nil {
			current.Reject(history, err)
			return AgentResult{}, err
		}
	}
	if s.client != nil && shouldTitle(history) {
		go s.title(id, plainHistory(history))
	}
	current.TakeAttachments()
	current.Persist = func(state agent.State) error { return s.store.SaveState(id, state) }
	current.Expand = s.expander(id, refs)
	current.AutoCompact = func() (bool, int, int) { return s.compactionSettings(id) }
	s.bindTasks(id, current)
	started, before, result := time.Now(), current.State().Usage, AgentResult{Executions: []agent.Execution{}}
	content, err := current.Run(ctx, history, func(event agent.Event) {
		if event.Kind == "steer" && event.Message != nil {
			s.recordSteer(id, event, &result, current)
		}
		if event.Kind == "notice" && event.Message != nil {
			s.recordNotice(id, event, &result, current)
		}
		if event.Execution != nil && event.Execution.State == "complete" {
			result.Executions = append(result.Executions, *event.Execution)
			if err := s.store.AddStep(id, *event.Execution); err != nil {
				log.Printf("persist step: %v", err)
			}
		}
		s.recordLive(id, event)
		if s.emit != nil {
			event.Session = id
			s.emit(event)
		}
	})
	result.Files = current.TakeAttachments()
	result.Content = content
	result.Turn = s.turn(id, started, before, current.State().Usage)
	if errors.Is(err, context.Canceled) {
		result.Error = "Stopped"
	} else if err != nil {
		result.Error = err.Error()
	}
	if saveErr := s.saveAnswer(id, content, result, current); saveErr != nil {
		return result, fmt.Errorf("Could not save session: %w", saveErr)
	}
	return result, nil
}
