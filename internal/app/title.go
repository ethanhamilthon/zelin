package app

import (
	"context"
	"encoding/json"
	"strings"

	"zelin/internal/prompts"
	"zelin/internal/provider"
)

type titleEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func shouldTitle(history []provider.Message) bool {
	if len(history) == 1 {
		return true
	}
	assistants := 0
	for _, msg := range history {
		if msg.Role == "assistant" {
			assistants++
		}
	}
	return assistants == 1
}

func (s *AgentService) utilityModel() string {
	if s.store != nil {
		if model, _ := s.store.Setting("titleModel"); model != "" {
			return model
		}
	}
	if s.client != nil {
		model, _ := s.client.Selection()
		return model
	}
	return ""
}

func (s *AgentService) title(id string, history []provider.Message) {
	items := make([]titleEntry, 0, len(history))
	for _, msg := range history {
		if msg.Role == "user" || msg.Role == "assistant" {
			items = append(items, titleEntry{Role: msg.Role, Content: msg.Content})
		}
	}
	data, err := json.Marshal(items)
	if err != nil || len(items) == 0 {
		return
	}
	payload := "Generate a concise topic title (target 4 words, maximum 6 words) for the conversation transcript below. Do not reply to the messages; output only the title.\n\n<conversation>\n" + string(data) + "\n</conversation>"
	reply, _, err := s.client.Complete(context.Background(), provider.Request{
		Model: s.utilityModel(),
		Messages: []provider.Message{
			{Role: "system", Content: prompts.Title},
			{Role: "user", Content: payload},
		},
	})
	if err != nil {
		return
	}
	title := cleanTitle(reply.Content)
	if title == "" || s.notify == nil {
		return
	}
	if ok, err := s.store.AutoTitle(id, title); !ok || err != nil {
		return
	}
	s.notify("session:title", map[string]string{"id": id, "title": title})
}

func cleanTitle(text string) string {
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			text = line
			break
		}
	}
	text = strings.NewReplacer("**", "", "`", "", "\"", "", "«", "", "»", "").Replace(text)
	text = strings.Trim(strings.TrimSpace(text), "'*#.-—:;,!?")
	words := strings.Fields(text)
	if len(words) > 6 {
		words = words[:6]
	}
	text = strings.TrimRight(strings.Join(words, " "), ".,;:!?-—")
	runes := []rune(strings.TrimSpace(text))
	if len(runes) > 60 {
		runes = runes[:60]
	}
	return strings.TrimSpace(string(runes))
}
