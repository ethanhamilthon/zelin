package sessions

import "testing"

func TestDeleteWorkspaceMovesSessionsToDefault(t *testing.T) {
	store, _ := open(t)
	if err := store.EnsureDefault(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	item, err := store.CreateWorkspace("Temp", "", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	session, err := store.CreateIn(item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.DeleteWorkspace(DefaultWorkspace); err == nil {
		t.Fatal("Deleted the default workspace")
	}
	if err := store.DeleteWorkspace(item.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Workspace(item.ID); err == nil {
		t.Fatal("Workspace still exists")
	}
	if saved, err := store.Get(session.ID); err != nil || saved.WorkspaceID != DefaultWorkspace {
		t.Fatalf("%+v %v", saved.WorkspaceID, err)
	}
	if err := store.DeleteWorkspace(item.ID); err == nil {
		t.Fatal("Deleted a missing workspace")
	}
}
