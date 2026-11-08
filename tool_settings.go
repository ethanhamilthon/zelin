package main

import "encoding/json"

func (s *AgentService) toolOutputLimit() int {
	settings := struct {
		Limit int `json:"toolOutputLimit"`
	}{4096}
	if s.store != nil {
		raw, _ := s.store.Setting("ui")
		_ = json.Unmarshal([]byte(raw), &settings)
	}
	switch settings.Limit {
	case 0, 1024, 2048, 4096, 8192, 16384:
		return settings.Limit
	default:
		return 4096
	}
}
