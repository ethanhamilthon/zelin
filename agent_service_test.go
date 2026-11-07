package main

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"zash/internal/sessions"

	"zash/internal/agent"
	"zash/internal/provider"
)

type serviceClient struct {
	step    int
	success bool
}

func (c *serviceClient) Complete(context.Context, []provider.Message) (provider.Message, provider.Usage, error) {
	c.step++
	if c.step == 1 {
		call := provider.ToolCall{ID: "call", Type: "function"}
		call.Function.Name = "js_repl"
		call.Function.Arguments = `{"code":"1 + 1"}`
		return provider.Message{Role: "assistant", ToolCalls: []provider.ToolCall{call}}, provider.Usage{}, nil
	}
	if c.success {
		return provider.Message{Role: "assistant", Content: "Result is 2"}, provider.Usage{}, nil
	}
	return provider.Message{}, provider.Usage{}, fmt.Errorf("Provider disconnected")
}

func TestServicePreservesExecutionsOnFailure(t *testing.T) {
	var events []agent.Event
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, err := store.Create()
	if err != nil {
		t.Fatal(err)
	}
	service := &AgentService{
		store:  store,
		agents: map[string]*agent.Agent{session.ID: agent.New(&serviceClient{})},
		emit:   func(event agent.Event) { events = append(events, event) },
	}
	result, err := service.Run(context.Background(), session.ID, []provider.Message{{Role: "user", Content: "Calculate"}})
	if err != nil || result.Error != "Provider disconnected" || len(result.Executions) != 1 {
		t.Fatalf("%+v %v", result, err)
	}
	if result.Executions[0].Output != "2" || len(events) != 6 {
		t.Fatalf("%+v", result)
	}
	if saved, err := store.Get(session.ID); err != nil || saved.Usage.Requests != 1 {
		t.Fatalf("Usage was not saved: %+v %v", saved.Usage, err)
	}
}
