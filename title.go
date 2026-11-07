package main

import (
	"context"
	"strings"

	"zash/internal/provider"
)

const titlePrompt = "Write a short title (2 to 6 words) for a conversation that starts with the user's message below. Reply with the title only, without quotes or punctuation at the end."

func (s *AgentService) title(id, message string) {
	ctx := context.Background()
	model, _ := s.store.Setting("titleModel")
	if model == "" {
		model, _ = s.client.Selection()
	}
	reply, _, err := s.client.Complete(ctx, provider.Request{Model: model, Effort: "low", Messages: []provider.Message{
		{Role: "system", Content: titlePrompt},
		{Role: "user", Content: message},
	}})
	if err != nil {
		return
	}
	title := cleanTitle(reply.Content)
	if title == "" || s.store.Rename(id, title) != nil || s.notify == nil {
		return
	}
	s.notify("session:title", map[string]string{"id": id, "title": title})
}

func cleanTitle(text string) string {
	text, _, _ = strings.Cut(strings.TrimSpace(text), "\n")
	text = strings.Trim(strings.TrimSpace(text), "\"'`*#.")
	runes := []rune(strings.TrimSpace(text))
	if len(runes) > 60 {
		runes = runes[:60]
	}
	return string(runes)
}
