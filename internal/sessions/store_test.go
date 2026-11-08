package sessions

import (
	"path/filepath"
	"testing"

	"zelin/internal/agent"
	"zelin/internal/provider"
)

func TestSessionsSurviveRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "zash", "sessions.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	first, err := store.Current()
	if err != nil {
		t.Fatal(err)
	}
	messages := []Message{
		{Role: "user", Content: "Hello SQLite"},
		{Role: "assistant", Content: "Done", Executions: []agent.Execution{{ID: "call", Code: "2+2", Output: "4", State: "complete"}}},
	}
	if err := store.Save(first.ID, messages, "Next draft", []string{"data:image/png;base64,AA"}); err != nil {
		t.Fatal(err)
	}
	state := agent.State{History: []provider.Message{{Role: "user", Content: "Hello SQLite"}}, Transcript: []provider.Message{{Role: "tool", Content: "4", ToolCallID: "call"}}, Usage: provider.Totals{Requests: 2}}
	if err := store.SaveState(first.ID, state); err != nil {
		t.Fatal(err)
	}
	second, err := store.Create()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Activate(first.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.SetSetting("model", "test"); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	loaded, err := store.Current()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.ID != first.ID || loaded.Title != "Hello SQLite" || loaded.Draft != "Next draft" || len(loaded.Messages) != 2 || loaded.Attachments[0] != "data:image/png;base64,AA" {
		t.Fatalf("Session was not restored: %+v", loaded)
	}
	if loaded.Messages[1].Executions[0].Output != "4" || loaded.State.Transcript[0].ToolCallID != "call" || loaded.Usage.Requests != 2 {
		t.Fatalf("History missing: %+v", loaded)
	}
	other, err := store.Get(second.ID)
	if err != nil || len(other.Messages) != 0 || len(other.State.History) != 0 {
		t.Fatalf("Session data leaked: %+v %v", other, err)
	}
	list, err := store.List()
	if err != nil || len(list) != 2 {
		t.Fatalf("%+v %v", list, err)
	}
	if model, err := store.Setting("model"); err != nil || model != "test" {
		t.Fatalf("%s %v", model, err)
	}
	if err := store.Save("missing", nil, "", nil); err == nil {
		t.Fatal("Accepted missing session")
	}
	if err := store.Rename(first.ID, "Generated"); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(first.ID, messages, "", nil); err != nil {
		t.Fatal(err)
	}
	if renamed, _ := store.Get(first.ID); renamed.Title != "Generated" || len(renamed.Attachments) != 0 {
		t.Fatalf("Generated title was overwritten: %+v", renamed)
	}
	if _, err := store.Activate("missing"); err == nil {
		t.Fatal("Activated missing session")
	}
}
