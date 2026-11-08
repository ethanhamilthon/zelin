package app

import (
	"path/filepath"
	"testing"

	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
)

func TestSteerTakesPriorStepsWithIt(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	if err := store.Save(session.ID, []sessions.Message{{Role: "user", Content: "go"}}, "", nil); err != nil {
		t.Fatal(err)
	}
	service := &AgentService{store: store}
	result := AgentResult{Executions: []agent.Execution{{ID: "a", State: "complete"}, {ID: "b", State: "complete"}}}
	event := agent.Event{Kind: "steer", Message: &provider.Message{Role: "user", Content: "change course"}}
	service.recordSteer(session.ID, event, &result, agent.New(nil))
	record, err := store.Get(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	last := record.Messages[len(record.Messages)-1]
	if last.Role != "user" || len(last.Executions) != 2 || len(result.Executions) != 0 {
		t.Fatalf("Steps not attached to steer: %+v %+v", last, result.Executions)
	}
}
