package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"zelin/internal/agent"
	"zelin/internal/instructions"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/tagged"
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
}

type AgentResult struct {
	Content    string            `json:"content"`
	Executions []agent.Execution `json:"executions"`
	Turn       *sessions.Turn    `json:"turn,omitempty"`
	Error      string            `json:"error,omitempty"`
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
	ctx, cancel := context.WithCancel(ctx)
	s.runs[id] = cancel
	return current, ctx, nil
}

func (s *AgentService) finish(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cancel := s.runs[id]; cancel != nil {
		cancel()
	}
	if current := s.agents[id]; current != nil {
		current.Steer(provider.Message{})
	}
	delete(s.runs, id)
}

func (s *AgentService) Stop(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cancel := s.runs[id]; cancel != nil {
		cancel()
		return true
	}
	return false
}

func (s *AgentService) Run(ctx context.Context, id string, history []provider.Message, refs []tagged.Ref) (AgentResult, error) {
	current, ctx, err := s.start(ctx, id)
	if err != nil {
		return AgentResult{}, err
	}
	defer s.finish(id)
	if s.client != nil {
		if err := s.pin(id); err != nil {
			return AgentResult{}, err
		}
	}
	if s.client != nil && shouldTitle(history) {
		titled := append([]provider.Message{}, history...)
		for i := range titled {
			titled[i].Content = tagged.Plain(titled[i].Content)
		}
		go s.title(id, titled)
	}
	current.Expand = s.expander(id, refs)
	current.AutoCompact = func() (bool, int, int) { return s.compactionSettings(id) }
	started, before := time.Now(), current.State().Usage
	result := AgentResult{Executions: []agent.Execution{}}
	content, err := current.Run(ctx, history, func(event agent.Event) {
		if event.Kind == "steer" && event.Message != nil {
			if err := s.store.Finish(id, sessions.Message{Role: "user", Content: event.Message.Content, Images: event.Message.Images, Refs: s.takeRefs(id)}, current.State()); err != nil {
				result.Error = err.Error()
				s.Stop(id)
			}
		}
		if event.Execution != nil && event.Execution.State == "complete" {
			result.Executions = append(result.Executions, *event.Execution)
		}
		if s.emit != nil {
			event.Session = id
			s.emit(event)
		}
	})
	result.Content = content
	result.Turn = s.turn(id, started, before, current.State().Usage)
	if errors.Is(err, context.Canceled) {
		result.Error = "Stopped"
	} else if err != nil {
		result.Error = err.Error()
	}
	var saveErr error
	if content != "" || len(result.Executions) > 0 {
		if content == "" {
			content = agent.StoppedMessage
		}
		saveErr = s.store.Finish(id, sessions.Message{Role: "assistant", Content: content, Executions: result.Executions, Turn: result.Turn}, current.State())
	} else {
		saveErr = s.store.SaveState(id, current.State())
	}
	if saveErr != nil {
		return result, fmt.Errorf("Could not save session: %w", saveErr)
	}
	return result, nil
}
