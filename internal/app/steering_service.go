package app

import (
	"strings"

	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/tagged"
)

func (s *AgentService) Steer(id string, message provider.Message, refs []tagged.Ref) (bool, error) {
	s.mu.Lock()
	running := s.runs[id] != nil && s.agents[id] != nil
	current := s.agents[id]
	s.mu.Unlock()
	if !running {
		return false, nil
	}
	message.Role = "user"
	if strings.TrimSpace(message.Content) == "" {
		current.Steer(provider.Message{})
		return true, nil
	}
	model, err := s.expand(id, message, refs)
	if err != nil {
		return false, err
	}
	s.mu.Lock()
	if s.steerRefs == nil {
		s.steerRefs = map[string][]tagged.Ref{}
	}
	s.steerRefs[id] = tagged.Used(message.Content, refs)
	s.mu.Unlock()
	current.SteerWith(message, model)
	return true, nil
}

func (s *AgentService) takeRefs(id string) []tagged.Ref {
	s.mu.Lock()
	defer s.mu.Unlock()
	refs := s.steerRefs[id]
	delete(s.steerRefs, id)
	return refs
}

func (s *AgentService) recordSteer(id string, event agent.Event, result *AgentResult, current *agent.Agent) {
	message := sessions.Message{Role: "user", Content: event.Message.Content, Images: event.Message.Images, Refs: s.takeRefs(id), Executions: result.Executions}
	if err := s.store.Finish(id, message, current.State()); err != nil {
		result.Error = err.Error()
		s.Stop(id)
		return
	}
	result.Executions = []agent.Execution{}
}
