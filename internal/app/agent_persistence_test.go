package app

import (
	"context"
	"path/filepath"
	"testing"

	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
)

func TestAnswerIsDurableBeforeReturningToUI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sessions.db")
	store, err := sessions.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	session, err := store.Create()
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(session.ID, []sessions.Message{{Role: "user", Content: "Calculate"}}, "", nil); err != nil {
		t.Fatal(err)
	}
	service := &AgentService{store: store, agents: map[string]*agent.Agent{session.ID: agent.New(&serviceClient{success: true})}}
	result, err := service.Run(context.Background(), session.ID, []provider.Message{{Role: "user", Content: "Calculate"}}, nil)
	if err != nil || result.Content != "Result is 2" {
		t.Fatalf("%+v %v", result, err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := sessions.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	saved, err := reopened.Get(session.ID)
	if err != nil || len(saved.Messages) != 2 || saved.Messages[1].Content != "Result is 2" || saved.Messages[1].Executions[0].Output != "2" || saved.Messages[1].Turn == nil || result.Turn == nil {
		t.Fatalf("Response was not durable: %+v %v", saved, err)
	}
	if len(saved.State.Transcript) != 5 {
		t.Fatalf("Provider context was not retained: %+v", saved.State)
	}
}
