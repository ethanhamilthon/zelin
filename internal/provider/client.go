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

	"zelin/internal/metadata"
)

var Efforts = []string{"none", "low", "medium", "high", "xhigh", "max"}

const DefaultEffort = "medium"

type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type Request struct {
	Model    string
	Effort   string
	Messages []Message
	Tools    []Tool
	OnDelta  func(Delta)
}

type Client struct {
	URL     string
	Key     string
	Model   string
	Effort  string
	HTTP    *http.Client
	Meta    *metadata.Catalog
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
	return &Client{URL: strings.TrimRight(base, "/") + "/chat/completions", Key: os.Getenv("OPENAI_API_KEY"), Model: model, Effort: DefaultEffort, HTTP: &http.Client{}}
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

func (c *Client) complete(ctx context.Context, request Request) (Message, Usage, error) {
	url, key := c.endpoint()
	if key == "" {
		return Message{}, Usage{}, fmt.Errorf("No API key configured")
	}
	model, effort := c.Selection()
	if request.Model != "" {
		model, effort = request.Model, request.Effort
	}
	meta := c.metadata(ctx, model)
	payload := map[string]any{"model": model, "messages": wire(request.Messages)}
	if len(request.Tools) > 0 {
		tools := make([]map[string]any, len(request.Tools))
		for i, tool := range request.Tools {
			tools[i] = map[string]any{"type": "function", "function": tool}
		}
		payload["tools"] = tools
	}
	if effort != "" && effort != "none" && (meta.Reasoning == nil || *meta.Reasoning) {
		payload["reasoning_effort"] = effort
	}
	if request.OnDelta != nil {
		payload["stream"] = true
		payload["stream_options"] = map[string]any{"include_usage": true}
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return Message{}, Usage{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return Message{}, Usage{}, err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return Message{}, Usage{}, retryable{fmt.Errorf("Provider request failed: %w", err)}
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err := httpError(res.StatusCode, res.Body, key)
		if res.StatusCode == 408 || res.StatusCode == 429 || res.StatusCode >= 500 {
			return Message{}, Usage{}, retryable{err}
		}
		return Message{}, Usage{}, err
	}
	if request.OnDelta != nil {
		message, raw, err := readStream(res.Body, request.OnDelta)
		return message, raw.usage(meta.Pricing), err
	}
	var data struct {
		Choices []struct {
			Message Message `json:"message"`
		} `json:"choices"`
		Usage *rawUsage `json:"usage"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 16<<20)).Decode(&data); err != nil {
		return Message{}, Usage{}, retryable{fmt.Errorf("Invalid provider response: %w", err)}
	}
	if len(data.Choices) == 0 {
		return Message{}, Usage{}, fmt.Errorf("Provider returned no choices")
	}
	return data.Choices[0].Message, data.Usage.usage(meta.Pricing), nil
}
