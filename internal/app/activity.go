package app

type SessionActivity struct {
	Running bool `json:"running"`
	Unread  bool `json:"unread"`
}

func (s *AgentService) Activity() map[string]SessionActivity {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := map[string]SessionActivity{}
	for id, item := range s.activity {
		result[id] = item
	}
	return result
}

func (s *AgentService) Read(id string) {
	s.mu.Lock()
	item := s.activity[id]
	item.Unread = false
	if s.activity != nil {
		s.activity[id] = item
	}
	s.mu.Unlock()
}

func (s *AgentService) activityLocked(id string, running bool) {
	if s.activity == nil {
		s.activity = map[string]SessionActivity{}
	}
	s.activity[id] = SessionActivity{Running: running, Unread: !running}
}
