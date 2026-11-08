package sessions

import (
	"path/filepath"
	"testing"

	"zelin/internal/agent"
)

func TestPendingStepsSurviveRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sessions.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	session, _ := store.Current()
	if err := store.Save(session.ID, []Message{{Role: "user", Content: "go"}}, "", nil); err != nil {
		t.Fatal(err)
	}
	if err := store.AddStep(session.ID, agent.Execution{ID: "1", Code: "1", Output: "1", State: "complete"}); err != nil {
		t.Fatal(err)
	}
	store.Close()
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	loaded, _ := store.Get(session.ID)
	if len(loaded.Messages) != 2 || loaded.Messages[1].Role != "assistant" || len(loaded.Messages[1].Executions) != 1 {
		t.Fatalf("Steps lost: %+v", loaded.Messages)
	}
	if err := store.Finish(session.ID, Message{Role: "assistant", Content: "ok"}, agent.State{}); err != nil {
		t.Fatal(err)
	}
	store.Close()
	store, _ = Open(path)
	defer store.Close()
	again, _ := store.Get(session.ID)
	if len(again.Messages) != 3 {
		t.Fatalf("Unexpected recovery: %+v", again.Messages)
	}
}
