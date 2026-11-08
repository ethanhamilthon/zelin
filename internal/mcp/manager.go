package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	InputSchema json.RawMessage `json:"inputSchema,omitempty"`
}

type Manager struct {
	Servers func() []Server
	mu      sync.Mutex
	clients map[string]*Client
}

func (m *Manager) find(name string) (Server, error) {
	for _, server := range m.Servers() {
		if server.Name == name {
			return server, nil
		}
	}
	return Server{}, fmt.Errorf("Unknown MCP server %q", name)
}

func (m *Manager) client(ctx context.Context, name, dir string) (*Client, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c := m.clients[name]; c != nil {
		select {
		case <-c.done:
			delete(m.clients, name)
		default:
			return c, nil
		}
	}
	server, err := m.find(name)
	if err != nil {
		return nil, err
	}
	c, err := Start(server, dir)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	_, err = c.Call(ctx, "initialize", map[string]any{
		"protocolVersion": "2025-06-18",
		"capabilities":    map[string]any{},
		"clientInfo":      map[string]string{"name": "zelin", "version": "1"},
	})
	if err == nil {
		err = c.Notify("notifications/initialized")
	}
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("MCP server %s failed to initialize: %w", name, err)
	}
	if m.clients == nil {
		m.clients = map[string]*Client{}
	}
	m.clients[name] = c
	return c, nil
}

func (m *Manager) Tools(ctx context.Context, name, dir string) ([]Tool, error) {
	c, err := m.client(ctx, name, dir)
	if err != nil {
		return nil, err
	}
	raw, err := c.Call(ctx, "tools/list", map[string]any{})
	if err != nil {
		return nil, err
	}
	var result struct {
		Tools []Tool `json:"tools"`
	}
	return result.Tools, json.Unmarshal(raw, &result)
}

func (m *Manager) Call(ctx context.Context, name, dir, tool string, args map[string]any) (string, error) {
	c, err := m.client(ctx, name, dir)
	if err != nil {
		return "", err
	}
	if args == nil {
		args = map[string]any{}
	}
	raw, err := c.Call(ctx, "tools/call", map[string]any{"name": tool, "arguments": args})
	if err != nil {
		return "", err
	}
	return render(raw)
}

func render(raw json.RawMessage) (string, error) {
	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Structured json.RawMessage `json:"structuredContent"`
		IsError    bool            `json:"isError"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", err
	}
	var parts []string
	for _, item := range result.Content {
		if item.Type == "text" {
			parts = append(parts, item.Text)
		} else {
			parts = append(parts, "["+item.Type+" content]")
		}
	}
	if len(parts) == 0 && len(result.Structured) > 0 {
		parts = append(parts, string(result.Structured))
	}
	text := strings.Join(parts, "\n")
	if result.IsError {
		return "", fmt.Errorf("%s", text)
	}
	return text, nil
}

func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for name, c := range m.clients {
		c.Close()
		delete(m.clients, name)
	}
}
