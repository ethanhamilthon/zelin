package sessions

import (
	"path/filepath"
	"testing"
)

func TestSnapshotIncludesUnflushedJournalAndLeavesSourceAlone(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "live.db")
	store, err := Open(source)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	store.Save(session.ID, []Message{{Role: "user", Content: "recent"}}, "draft", nil)
	target := filepath.Join(dir, "copy.db")
	if err := Snapshot(source, target); err != nil {
		t.Fatal(err)
	}
	copied, err := Open(target)
	if err != nil {
		t.Fatal(err)
	}
	defer copied.Close()
	got, err := copied.Get(session.ID)
	if err != nil || len(got.Messages) != 1 || got.Draft != "draft" {
		t.Fatalf("%+v %v", got, err)
	}
	if err := Snapshot(source, target); err == nil {
		t.Fatal("Overwrote an existing snapshot")
	}
}
