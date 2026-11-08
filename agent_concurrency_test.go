package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
)

type blockingClient struct{ started chan struct{} }

func (c *blockingClient) Complete(ctx context.Context, _ provider.Request) (provider.Message, provider.Usage, error) {
	close(c.started)
	<-ctx.Done()
	return provider.Message{}, provider.Usage{}, ctx.Err()
}

type replyClient struct{}

func (replyClient) Complete(context.Context, provider.Request) (provider.Message, provider.Usage, error) {
	return provider.Message{Role: "assistant", Content: "Done", ReasoningContent: "Plan"}, provider.Usage{}, nil
}

func TestSessionsRunConcurrentlyAndStop(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	first, _ := store.Create()
	second, _ := store.Create()
	slow := &blockingClient{started: make(chan struct{})}
	events := make(chan agent.Event, 32)
	service := &AgentService{store: store, emit: func(event agent.Event) { events <- event }, agents: map[string]*agent.Agent{
		first.ID: agent.New(slow), second.ID: agent.New(replyClient{}),
	}}
	history := []provider.Message{{Role: "user", Content: "Go"}}
	done := make(chan AgentResult)
	go func() { result, _ := service.Run(context.Background(), first.ID, history, nil); done <- result }()
	<-slow.started
	if _, err := service.Run(context.Background(), first.ID, history, nil); err == nil {
		t.Fatal("Started the same session twice")
	}
	result, err := service.Run(context.Background(), second.ID, history, nil)
	if err != nil || result.Content != "Done" || len(result.Executions) != 1 || result.Executions[0].Kind != "thinking" {
		t.Fatalf("%+v %v", result, err)
	}
	if !service.Stop(first.ID) {
		t.Fatal("Stop found no run")
	}
	select {
	case stopped := <-done:
		if stopped.Error != "Stopped" {
			t.Fatalf("%+v", stopped)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not stop")
	}
	if service.Stop(first.ID) {
		t.Fatal("Stopped an idle session")
	}
	for len(events) > 0 {
		if event := <-events; event.Session != first.ID && event.Session != second.ID {
			t.Fatalf("Event without session: %+v", event)
		}
	}
}

func TestCleanTitle(t *testing.T) {
	if got := cleanTitle("  \"Fix the build.\"\nextra"); got != "Fix the build" {
		t.Fatalf("%q", got)
	}
}
