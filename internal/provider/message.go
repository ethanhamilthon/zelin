package provider

import (
	"encoding/json"
	"regexp"
	"strconv"
)

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type Message struct {
	Role             string          `json:"role"`
	Content          string          `json:"content"`
	Images           []string        `json:"images,omitempty"`
	ToolCalls        []ToolCall      `json:"tool_calls,omitempty"`
	ToolCallID       string          `json:"tool_call_id,omitempty"`
	Reasoning        string          `json:"reasoning,omitempty"`
	ReasoningContent string          `json:"reasoning_content,omitempty"`
	ReasoningDetails json.RawMessage `json:"reasoning_details,omitempty"`
}

func (m Message) Thinking() string {
	if m.ReasoningContent != "" {
		return m.ReasoningContent
	}
	return m.Reasoning
}

var imageToken = regexp.MustCompile(`\[Image #(\d+)\]`)

type wireMessage struct {
	Message
	Content any `json:"content"`
	Images  any `json:"images,omitempty"`
}

func wire(messages []Message) []wireMessage {
	result := make([]wireMessage, len(messages))
	for i, message := range messages {
		result[i] = wireMessage{Message: message, Content: content(message)}
	}
	return result
}

func content(message Message) any {
	if len(message.Images) == 0 {
		return message.Content
	}
	parts, last := []map[string]any{}, 0
	text := func(value string) {
		if value != "" {
			parts = append(parts, map[string]any{"type": "text", "text": value})
		}
	}
	for _, match := range imageToken.FindAllStringSubmatchIndex(message.Content, -1) {
		number, _ := strconv.Atoi(message.Content[match[2]:match[3]])
		if number < 1 || number > len(message.Images) || message.Images[number-1] == "" {
			continue
		}
		text(message.Content[last:match[1]])
		parts = append(parts, map[string]any{"type": "image_url", "image_url": map[string]string{"url": message.Images[number-1]}})
		last = match[1]
	}
	text(message.Content[last:])
	return parts
}
