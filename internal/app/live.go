package app

import (
	"time"
	"zelin/internal/agent"
)

type RunSnapshot struct {
	Started int64         `json:"started"`
	Events  []agent.Event `json:"events"`
}

func (s *AgentService) Active(id string) *RunSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	if live := s.live[id]; live != nil {
		return &RunSnapshot{Started: live.Started, Events: append([]agent.Event{}, live.Events...)}
	}
	return nil
}

func (s *AgentService) beginLive(id string) {
	s.mu.Lock()
	if s.live == nil {
		s.live = map[string]*RunSnapshot{}
	}
	s.activityLocked(id, true)
	s.live[id] = &RunSnapshot{Started: time.Now().UnixMilli(), Events: []agent.Event{}}
	s.mu.Unlock()
	if s.notify != nil {
		s.notify("session:changed", id)
	}
}

func (s *AgentService) recordLive(id string, event agent.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if live := s.live[id]; live != nil {
		if event.Kind == "steer" || event.Kind == "notice" {
			kept := live.Events[:0]
			for _, previous := range live.Events {
				if previous.Kind != "execution" {
					kept = append(kept, previous)
				}
			}
			live.Events = kept
			return
		}
		for i, previous := range live.Events {
			if previous.Kind == event.Kind && (event.Execution == nil || previous.Execution != nil && previous.Execution.ID == event.Execution.ID) {
				live.Events[i] = event
				return
			}
		}
		live.Events = append(live.Events, event)
	}
}

func (s *AgentService) endLive(id string) {
	s.mu.Lock()
	delete(s.live, id)
	s.activityLocked(id, false)
	s.mu.Unlock()
	if s.notify != nil {
		s.notify("session:changed", id)
	}
}
