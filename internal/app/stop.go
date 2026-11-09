package app

import "zelin/internal/provider"

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
		if s.paused == nil {
			s.paused = map[string]bool{}
		}
		s.paused[id] = true
		cancel()
		return true
	}
	return false
}
