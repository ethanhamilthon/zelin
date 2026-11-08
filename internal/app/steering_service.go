package app

import (
	"strings"

	"zelin/internal/provider"
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
