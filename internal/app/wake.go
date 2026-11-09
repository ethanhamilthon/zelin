package app

import (
	"context"
	"errors"
	"log"

	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/tagged"
	"zelin/internal/tasks"
)

func (s *AgentService) taskSession(id string) *tasks.Session {
	if s.tasks == nil {
		return nil
	}
	return s.tasks.Session(id)
}

func (s *AgentService) bindTasks(id string, current *agent.Agent) {
	session := s.taskSession(id)
	current.Tasks = session
	if session == nil {
		current.Notices = nil
		return
	}
	forced := session.Take(true)
	current.Notices = func() string {
		if text := forced; text != "" {
			forced = ""
			return text
		}
		return session.Take(false)
	}
}

func (s *AgentService) running(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.runs[id] != nil || s.paused[id]
}

func (s *AgentService) deliver(id string) {
	if !s.running(id) {
		go s.wake(id)
	}
}

func (s *AgentService) wake(id string) {
	session := s.taskSession(id)
	if session == nil || s.running(id) || !session.Ready() {
		return
	}
	record, err := s.store.Get(id)
	if err != nil {
		log.Printf("task wake: %v", err)
		return
	}
	history := []provider.Message{}
	for _, message := range record.Messages {
		history = append(history, provider.Message{Role: message.Role, Content: message.Content, Images: message.Images})
	}
	result, err := s.run(context.Background(), id, history, nil, func() (string, error) {
		text := session.Take(true)
		if text == "" {
			return "", errNothing
		}
		return text, s.store.Append(id, sessions.Message{Role: "user", Kind: "task", Content: text})
	})
	if err == nil {
		s.afterRun(id, result)
	} else if !errors.Is(err, errNothing) {
		log.Printf("task wake: %v", err)
	}
}

func (s *AgentService) recordNotice(id string, event agent.Event, result *AgentResult, current *agent.Agent) {
	message := sessions.Message{Role: "user", Kind: "task", Content: event.Message.Content, Executions: result.Executions}
	if err := s.store.Finish(id, message, current.State()); err != nil {
		result.Error = err.Error()
		s.Stop(id)
		return
	}
	result.Executions = []agent.Execution{}
}

func (s *AgentService) Run(ctx context.Context, id string, history []provider.Message, refs []tagged.Ref) (AgentResult, error) {
	s.mu.Lock()
	delete(s.paused, id)
	s.mu.Unlock()
	result, err := s.run(ctx, id, history, refs, nil)
	if err == nil {
		s.afterRun(id, result)
	}
	return result, err
}

func (s *AgentService) afterRun(id string, result AgentResult) {
	if session := s.taskSession(id); session != nil && session.Ready() {
		go s.wake(id)
	}
}

func claimed(history []provider.Message, claim func() (string, error)) ([]provider.Message, error) {
	if claim == nil {
		return history, nil
	}
	text, err := claim()
	if err != nil {
		return nil, err
	}
	return append(history, provider.Message{Role: "user", Content: text}), nil
}

var errNothing = errors.New("no task notices")
