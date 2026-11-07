package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"zash/internal/agent"
	"zash/internal/instructions"
	"zash/internal/provider"
	"zash/internal/sessions"
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
}

type AgentResult struct {
	Content    string            `json:"content"`
	Executions []agent.Execution `json:"executions"`
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
		current = agent.New(s.client)
		current.Restore(record.State)
		s.agents[id] = current
	}
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

func (s *AgentService) Run(ctx context.Context, id string, history []provider.Message) (AgentResult, error) {
	current, ctx, err := s.start(ctx, id)
	if err != nil {
		return AgentResult{}, err
	}
	defer s.finish(id)
	if len(history) == 1 && s.client != nil {
		go s.title(id, history[0].Content)
	}
	result := AgentResult{Executions: []agent.Execution{}}
	content, err := current.Run(ctx, history, func(event agent.Event) {
		if event.Execution != nil && event.Execution.State == "complete" {
			result.Executions = append(result.Executions, *event.Execution)
		}
		if s.emit != nil {
			event.Session = id
			s.emit(event)
		}
	})
	result.Content = content
	if errors.Is(err, context.Canceled) {
		result.Error = "Stopped"
	} else if err != nil {
		result.Error = err.Error()
	}
	var saveErr error
	if content != "" || len(result.Executions) > 0 {
		if content == "" {
			content = "The agent stopped before completing this turn."
		}
		saveErr = s.store.Finish(id, sessions.Message{Role: "assistant", Content: content, Executions: result.Executions}, current.State())
	} else {
		saveErr = s.store.SaveState(id, current.State())
	}
	if saveErr != nil {
		return result, fmt.Errorf("Could not save session: %w", saveErr)
	}
	return result, nil
}
