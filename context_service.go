package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"zash/internal/agent"
	"zash/internal/provider"
	"zash/internal/sessions"
)

func (s *AgentService) compactionSettings(id string) (bool, int, int) {
	if s.client == nil {
		return false, 60, 0
	}
	settings := struct {
		AutoCompact      bool `json:"autoCompact"`
		CompactThreshold int  `json:"compactThreshold"`
	}{true, 60}
	raw, _ := s.store.Setting("ui")
	_ = json.Unmarshal([]byte(raw), &settings)
	if settings.CompactThreshold < 1 || settings.CompactThreshold > 100 {
		settings.CompactThreshold = 60
	}
	return settings.AutoCompact, settings.CompactThreshold, s.Info(context.Background(), id).ContextWindow
}

func (s *AgentService) Compact(ctx context.Context, id string) (provider.Totals, error) {
	current, ctx, err := s.start(ctx, id)
	if err != nil {
		return provider.Totals{}, err
	}
	defer s.finish(id)
	record, err := s.store.Get(id)
	if err != nil {
		return provider.Totals{}, err
	}
	history := []provider.Message{}
	for _, message := range record.Messages {
		history = append(history, provider.Message{Role: message.Role, Content: message.Content, Images: message.Images})
	}
	if err := current.Compact(ctx, history, func(event agent.Event) {
		if s.emit != nil {
			event.Session = id
			s.emit(event)
		}
	}); err != nil {
		return provider.Totals{}, err
	}
	return current.State().Usage, s.store.SaveState(id, current.State())
}

func (s *AgentService) Handoff(ctx context.Context, id string) (sessions.Session, error) {
	_, ctx, err := s.start(ctx, id)
	if err != nil {
		return sessions.Session{}, err
	}
	defer s.finish(id)
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	record, err := s.store.Get(id)
	if err != nil {
		return sessions.Session{}, err
	}
	if len(record.Messages) == 0 {
		return sessions.Session{}, fmt.Errorf("No conversation to hand off")
	}
	data, err := json.Marshal(record.Messages)
	if err != nil {
		return sessions.Session{}, err
	}
	model, _ := s.store.Setting("titleModel")
	if model == "" {
		model, _ = s.client.Selection()
	}
	reply, _, err := s.client.Complete(ctx, provider.Request{Model: model, Effort: "low", Messages: []provider.Message{
		{Role: "system", Content: agent.CompactPrompt + " This is a handoff to a NEW session: include all information needed to resume without the original history. JavaScript bindings will be reset."},
		{Role: "user", Content: string(data)},
	}})
	if err != nil {
		return sessions.Session{}, err
	}
	if strings.TrimSpace(reply.Content) == "" {
		return sessions.Session{}, fmt.Errorf("Handoff returned an empty summary")
	}
	session, err := s.store.Create()
	if err != nil {
		return session, err
	}
	messages := []sessions.Message{{Role: "user", Content: "Handoff from the previous session:\n\n" + reply.Content}}
	if model, effort := s.selection(id); model != "" {
		if err := s.store.SetModel(session.ID, model, effort); err != nil {
			return session, err
		}
	}
	if err := s.store.Save(session.ID, messages, "", nil); err != nil {
		return session, err
	}
	if err := s.store.Rename(session.ID, "Handoff: "+record.Title); err != nil {
		return session, err
	}
	return s.store.Get(session.ID)
}
