package agent

import (
	"fmt"
	"slices"
	"strings"

	"zelin/internal/provider"
	"zelin/internal/tagged"
)

func (a *Agent) begin(history []provider.Message) ([]provider.Message, error) {
	extra := ""
	if a.Extra != nil {
		text, err := a.Extra()
		if err != nil {
			return nil, err
		}
		extra = text
	}
	pin := a.compose(extra)
	a.pinned = &pin
	if len(history) == 0 || history[len(history)-1].Role != "user" {
		a.pinned = nil
		return nil, fmt.Errorf("A user message is required")
	}
	messages := []provider.Message{pin}
	for _, message := range history {
		if (message.Role != "user" && message.Role != "assistant") || strings.TrimSpace(message.Content) == "" {
			a.pinned = nil
			return nil, fmt.Errorf("Invalid conversation history")
		}
		messages = append(messages, provider.Message{Role: message.Role, Content: tagged.Plain(message.Content), Images: message.Images})
	}
	if a.Expand != nil {
		expanded, err := a.Expand(history[len(history)-1])
		if err != nil {
			a.pinned = nil
			return nil, err
		}
		messages[len(messages)-1] = expanded
	}
	return messages, nil
}

func (a *Agent) continues(history []provider.Message) bool {
	return len(history) == len(a.history)+1 && slices.EqualFunc(history[:len(history)-1], a.history, func(x, y provider.Message) bool {
		return x.Role == y.Role && x.Content == y.Content
	})
}
