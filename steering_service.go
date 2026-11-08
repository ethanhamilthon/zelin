package main

import (
	"strings"
	"zash/internal/provider"
)

func (s *AgentService) Steer(id string, message provider.Message) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runs[id] == nil || s.agents[id] == nil {
		return false
	}
	message.Role = "user"
	if strings.TrimSpace(message.Content) == "" {
		message.Content = ""
	}
	s.agents[id].Steer(message)
	return true
}
