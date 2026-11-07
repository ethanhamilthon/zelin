package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type Client struct {
	URL   string
	Key   string
	Model string
	HTTP  *http.Client
}

func FromEnv() *Client {
	base := os.Getenv("OPENAI_BASE_URL")
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4.1-mini"
	}
	return &Client{strings.TrimRight(base, "/") + "/chat/completions", os.Getenv("OPENAI_API_KEY"), model, &http.Client{Timeout: 90 * time.Second}}
}

func (c *Client) Complete(ctx context.Context, messages []Message) (Message, error) {
	if c.Key == "" {
		return Message{}, fmt.Errorf("Set OPENAI_API_KEY before starting Zash")
	}
	tool := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name": "js_repl", "description": "Execute synchronous JavaScript in a persistent Goja REPL. The only host function is bash(cmd, timeout). Return the final expression as output. No async/await.",
			"parameters": map[string]any{"type": "object", "properties": map[string]any{"code": map[string]string{"type": "string"}}, "required": []string{"code"}, "additionalProperties": false},
		},
	}
	body, err := json.Marshal(map[string]any{"model": c.Model, "messages": messages, "tools": []any{tool}})
	if err != nil {
		return Message{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return Message{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return Message{}, fmt.Errorf("Provider request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Message{}, fmt.Errorf("Provider returned HTTP %d", res.StatusCode)
	}
	var data struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 4<<20)).Decode(&data); err != nil {
		return Message{}, fmt.Errorf("Invalid provider response: %w", err)
	}
	if len(data.Choices) == 0 {
		return Message{}, fmt.Errorf("Provider returned no choices")
	}
	return data.Choices[0].Message, nil
}
