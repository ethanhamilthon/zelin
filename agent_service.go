package main

import (
	"context"
	"fmt"
	"sync"

	"zash/internal/agent"
	"zash/internal/provider"
	"zash/internal/sessions"
)

type AgentService struct {
	mu     sync.Mutex
	agents map[string]*agent.Agent
	store  *sessions.Store
	client *provider.Client
	emit   func(agent.Event)
}

type AgentResult struct {
	Content    string            `json:"content"`
	Executions []agent.Execution `json:"executions"`
	Error      string            `json:"error,omitempty"`
}

func (s *AgentService) Info(ctx context.Context) provider.Metadata {
	return s.client.Metadata(ctx)
}

func (s *AgentService) Run(ctx context.Context, id string, history []provider.Message) (AgentResult, error) {
	if !s.mu.TryLock() {
		return AgentResult{}, fmt.Errorf("Agent is already running")
	}
	defer s.mu.Unlock()
	record, err := s.store.Get(id)
	if err != nil {
		return AgentResult{}, err
	}
	current := s.agents[id]
	if current == nil {
		current = agent.New(s.client)
		current.Restore(record.State)
		s.agents[id] = current
	}
	result := AgentResult{Executions: []agent.Execution{}}
	content, err := current.Run(ctx, history, func(event agent.Event) {
		if event.Execution != nil && event.Execution.State == "complete" {
			result.Executions = append(result.Executions, *event.Execution)
		}
		if s.emit != nil {
			s.emit(event)
		}
	})
	result.Content = content
	if err != nil {
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
