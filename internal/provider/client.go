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
	"sync"
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
	URL     string
	Key     string
	Model   string
	Effort  string
	HTTP    *http.Client
	mu      sync.Mutex
	catalog sync.Mutex
	models  []Model
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
	return &Client{URL: strings.TrimRight(base, "/") + "/chat/completions", Key: os.Getenv("OPENAI_API_KEY"), Model: model, HTTP: &http.Client{Timeout: 5 * time.Minute}}
}

func (c *Client) Selection() (string, string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Model, c.Effort
}

func (c *Client) Select(model, effort string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Model, c.Effort = model, effort
}

func (c *Client) pricing(id string) Pricing {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, model := range c.models {
		if model.ID == id {
			return model.Pricing
		}
	}
	return Pricing{}
}

func (c *Client) Complete(ctx context.Context, messages []Message) (Message, Usage, error) {
	if c.Key == "" {
		return Message{}, Usage{}, fmt.Errorf("Set OPENAI_API_KEY before starting Zash")
	}
	model, effort := c.Selection()
	tool := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name": "js_repl", "description": "Execute synchronous JavaScript in a persistent Goja REPL. The only host function is bash(cmd, timeout). Return the final expression as output. No async/await.",
			"parameters": map[string]any{"type": "object", "properties": map[string]any{"code": map[string]string{"type": "string"}}, "required": []string{"code"}, "additionalProperties": false},
		},
	}
	payload := map[string]any{"model": model, "messages": messages, "tools": []any{tool}}
	if effort != "" {
		payload["reasoning"] = map[string]string{"effort": effort}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return Message{}, Usage{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.URL, bytes.NewReader(body))
	if err != nil {
		return Message{}, Usage{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return Message{}, Usage{}, fmt.Errorf("Provider request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return Message{}, Usage{}, fmt.Errorf("Provider returned HTTP %d", res.StatusCode)
	}
	var data struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage *rawUsage `json:"usage"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 4<<20)).Decode(&data); err != nil {
		return Message{}, Usage{}, fmt.Errorf("Invalid provider response: %w", err)
	}
	if len(data.Choices) == 0 {
		return Message{}, Usage{}, fmt.Errorf("Provider returned no choices")
	}
	return data.Choices[0].Message, data.Usage.usage(c.pricing(model)), nil
}
