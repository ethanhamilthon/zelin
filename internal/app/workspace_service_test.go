package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"zelin/internal/sessions"
	"zelin/internal/workspace"
)

func newService(t *testing.T) (*WorkspaceService, *AgentService, *sessions.Store) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	if err := store.EnsureDefault(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	return &WorkspaceService{store: store}, &AgentService{store: store}, store
}

func TestNewSessionUsesWorkspaceAndEventsFire(t *testing.T) {
	w, agents, _ := newService(t)
	fired := 0
	w.notify = func(string, any) { fired++ }
	dir := t.TempDir()
	item, err := w.Create("Proj", "", dir)
	if err != nil || fired != 1 {
		t.Fatal(err, fired)
	}
	session, err := w.NewSession(item.ID)
	if err != nil || session.WorkspaceID != item.ID {
		t.Fatalf("%+v %v", session, err)
	}
	if got := agents.workspaceDir(session.WorkspaceID); got != dir {
		t.Fatal(got)
	}
	if _, err := w.NewSession("missing"); err == nil {
		t.Fatal("Created a session in a missing workspace")
	}
	if _, err := w.Update(item.ID, "", "", dir); err == nil || fired != 1 {
		t.Fatal("Invalid update must not notify")
	}
}

func TestChangingDirectoryKeepsExistingSessionState(t *testing.T) {
	w, agents, store := newService(t)
	first, second := t.TempDir(), t.TempDir()
	item, _ := w.Create("Proj", "", first)
	session, _ := w.NewSession(item.ID)
	if _, err := w.Update(item.ID, "Proj", "", second); err != nil {
		t.Fatal(err)
	}
	if got := agents.workspaceDir(item.ID); got != second {
		t.Fatal(got)
	}
	started := agents.workspaceDir(session.WorkspaceID)
	if started != second {
		t.Fatal(started)
	}
	record, _ := store.Get(session.ID)
	if record.State.Cwd != "" {
		t.Fatalf("Saved session cwd must stay explicit: %q", record.State.Cwd)
	}
}

func TestInstructionsComeFromTheSessionsWorkspace(t *testing.T) {
	w, agents, _ := newService(t)
	a, b := t.TempDir(), t.TempDir()
	os.WriteFile(workspace.AgentsPath(a), []byte("Rules A"), 0644)
	os.WriteFile(workspace.AgentsPath(b), []byte("Rules B"), 0644)
	wa, _ := w.Create("A", "", a)
	wb, _ := w.Create("B", "", b)
	sa, _ := w.NewSession(wa.ID)
	sb, _ := w.NewSession(wb.ID)
	ta, _ := agents.workspaceInstructions(sa.ID)
	tb, _ := agents.workspaceInstructions(sb.ID)
	if !strings.Contains(ta, "Rules A") || strings.Contains(ta, "Rules B") || !strings.Contains(tb, "Rules B") || strings.Contains(tb, "Rules A") {
		t.Fatalf("%q %q", ta, tb)
	}
	os.WriteFile(workspace.AgentsPath(a), []byte("Rules A2"), 0644)
	if ta, _ = agents.workspaceInstructions(sa.ID); !strings.Contains(ta, "Rules A2") {
		t.Fatal("Edit not picked up before the next turn")
	}
}

func TestAgentsFileReport(t *testing.T) {
	w, _, _ := newService(t)
	dir := t.TempDir()
	item, _ := w.Create("P", "", dir)
	file, err := w.Agents(item.ID)
	if err != nil || file.Exists || file.Path != filepath.Join(dir, "AGENTS.md") {
		t.Fatalf("%+v %v", file, err)
	}
	os.WriteFile(file.Path, []byte("x"), 0644)
	if file, _ = w.Agents(item.ID); !file.Exists {
		t.Fatal("Existing file not reported")
	}
	if _, err := w.Agents("missing"); err == nil {
		t.Fatal("Reported a missing workspace")
	}
}
