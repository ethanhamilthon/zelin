package sessions

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func open(t *testing.T) (*Store, string) {
	path := filepath.Join(t.TempDir(), "s.db")
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store, path
}

func TestWorkspaceValidationAndPaths(t *testing.T) {
	store, _ := open(t)
	root := t.TempDir()
	spaced := filepath.Join(root, "my project")
	os.Mkdir(spaced, 0755)
	file := filepath.Join(root, "file")
	os.WriteFile(file, nil, 0644)
	for _, bad := range [][3]string{{"", "", spaced}, {"  ", "", spaced}, {"x", "", ""}, {"x", "", filepath.Join(root, "missing")}, {"x", "", file}} {
		if _, err := store.CreateWorkspace(bad[0], bad[1], bad[2]); err == nil {
			t.Fatalf("Accepted %v", bad)
		}
	}
	item, err := store.CreateWorkspace("  Mine ", "★", spaced)
	if err != nil || item.Name != "Mine" || item.Dir != spaced || item.ID == "" {
		t.Fatalf("%+v %v", item, err)
	}
	home, _ := os.UserHomeDir()
	if item, err := store.CreateWorkspace("Home", "", "~"); err != nil || item.Dir != home {
		t.Fatalf("%+v %v", item, err)
	}
	updated, err := store.UpdateWorkspace(item.ID, "Renamed", "", root)
	if err != nil || updated.Dir != root {
		t.Fatalf("%+v %v", updated, err)
	}
	if _, err := store.UpdateWorkspace(item.ID, "Renamed", "", filepath.Join(root, "missing")); err == nil {
		t.Fatal("Accepted a missing directory")
	}
	if got, _ := store.Workspace(item.ID); got.Dir != root {
		t.Fatalf("Failed update changed state: %+v", got)
	}
	if _, err := store.UpdateWorkspace("missing", "x", "", root); err == nil {
		t.Fatal("Updated a missing workspace")
	}
}

func TestLegacySessionsMigrateToDefaultWorkspace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE sessions (id TEXT PRIMARY KEY, title TEXT NOT NULL, messages TEXT NOT NULL DEFAULT '[]', draft TEXT NOT NULL DEFAULT '', state TEXT NOT NULL DEFAULT '{}', updated TEXT NOT NULL);
 CREATE TABLE settings (key TEXT PRIMARY KEY, value TEXT NOT NULL);
 INSERT INTO sessions(id,title,messages,draft,updated) VALUES('old','Old chat','[{"role":"user","content":"hi"}]','keep me','2024-01-01T00:00:00Z');`)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	dir := t.TempDir()
	for i := 0; i < 2; i++ {
		if err := store.EnsureDefault(dir); err != nil {
			t.Fatal(err)
		}
	}
	old, err := store.Get("old")
	if err != nil || old.WorkspaceID != DefaultWorkspace || old.Draft != "keep me" || len(old.Messages) != 1 {
		t.Fatalf("%+v %v", old, err)
	}
	list, _ := store.Workspaces()
	if len(list) != 1 || list[0].Dir != dir {
		t.Fatalf("%+v", list)
	}
}

func TestSessionsBelongToWorkspaces(t *testing.T) {
	store, _ := open(t)
	if err := store.EnsureDefault(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	other, err := store.CreateWorkspace("Other", "", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := store.Create()
	b, _ := store.CreateIn(other.ID)
	if a.WorkspaceID != DefaultWorkspace || b.WorkspaceID != other.ID {
		t.Fatalf("%+v %+v", a, b)
	}
	list, _ := store.List()
	got := map[string]string{}
	for _, item := range list {
		got[item.ID] = item.WorkspaceID
	}
	if got[a.ID] != DefaultWorkspace || got[b.ID] != other.ID {
		t.Fatalf("%v", got)
	}
}

func TestWorkspaceSurvivesRestart(t *testing.T) {
	store, path := open(t)
	store.EnsureDefault(t.TempDir())
	item, _ := store.CreateWorkspace("Kept", "", t.TempDir())
	store.Close()
	again, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer again.Close()
	if got, err := again.Workspace(item.ID); err != nil || !strings.Contains(got.Name, "Kept") {
		t.Fatalf("%+v %v", got, err)
	}
}
