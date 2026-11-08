package provider

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Delta struct {
	Content  string
	Thinking string
	Calls    []ToolCall
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content          string          `json:"content"`
			Reasoning        string          `json:"reasoning"`
			ReasoningContent string          `json:"reasoning_content"`
			ReasoningDetails json.RawMessage `json:"reasoning_details"`
			ToolCalls        []struct {
				Index    int    `json:"index"`
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"delta"`
	} `json:"choices"`
	Usage *rawUsage `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func readStream(body io.Reader, onDelta func(Delta)) (Message, *rawUsage, error) {
	var message Message
	var usage *rawUsage
	var content, thinking strings.Builder
	var details []json.RawMessage
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 64<<10), 16<<20)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			break
		}
		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		if chunk.Error != nil {
			return Message{}, nil, retryable{fmt.Errorf("Provider stream error: %s", chunk.Error.Message)}
		}
		if chunk.Usage != nil {
			usage = chunk.Usage
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		d := chunk.Choices[0].Delta
		content.WriteString(d.Content)
		think := d.ReasoningContent + d.Reasoning
		if d.ReasoningContent != "" && d.Reasoning != "" {
			think = d.ReasoningContent
		}
		thinking.WriteString(think)
		if len(d.ReasoningDetails) > 0 && string(d.ReasoningDetails) != "null" {
			details = append(details, d.ReasoningDetails)
		}
		for _, tc := range d.ToolCalls {
			for len(message.ToolCalls) <= tc.Index {
				message.ToolCalls = append(message.ToolCalls, ToolCall{Type: "function"})
			}
			call := &message.ToolCalls[tc.Index]
			if tc.ID != "" {
				call.ID = tc.ID
			}
			call.Function.Name += tc.Function.Name
			call.Function.Arguments += tc.Function.Arguments
		}
		message.Content, message.ReasoningContent = content.String(), thinking.String()
		if onDelta != nil && (d.Content != "" || think != "" || len(d.ToolCalls) > 0) {
			onDelta(Delta{Content: message.Content, Thinking: message.ReasoningContent, Calls: append([]ToolCall{}, message.ToolCalls...)})
		}
	}
	if err := scanner.Err(); err != nil {
		return Message{}, nil, retryable{fmt.Errorf("Provider stream failed: %w", err)}
	}
	message.Role = "assistant"
	message.ReasoningDetails = mergeDetails(details)
	return message, usage, nil
}
