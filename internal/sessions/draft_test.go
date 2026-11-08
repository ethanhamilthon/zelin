package sessions

import (
	"path/filepath"
	"testing"
)

func TestDraftSavesKeepMessages(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sessions.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	session, err := store.Create()
	if err != nil {
		t.Fatal(err)
	}
	messages := []Message{{Role: "user", Content: "Keep me"}, {Role: "assistant", Content: "Done"}}
	if err := store.Save(session.ID, messages, "", []string{"data:image/png;base64,AA"}); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveDraft(session.ID, "Typing"); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveAttachments(session.ID, []string{"data:image/png;base64,BB"}); err != nil {
		t.Fatal(err)
	}
	store.Close()
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	loaded, err := store.Get(session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Draft != "Typing" || len(loaded.Messages) != 2 || loaded.Title != "Keep me" || loaded.Attachments[0] != "data:image/png;base64,BB" {
		t.Fatalf("Draft save changed history: %+v", loaded)
	}
	if err := store.SaveAttachments(session.ID, nil); err != nil {
		t.Fatal(err)
	}
	if cleared, _ := store.Get(session.ID); len(cleared.Attachments) != 0 {
		t.Fatalf("Attachments were not cleared: %+v", cleared)
	}
	if store.SaveDraft("missing", "x") == nil || store.SaveAttachments("missing", nil) == nil {
		t.Fatal("Accepted missing session")
	}
}
