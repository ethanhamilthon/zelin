package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	root := t.TempDir()
	spaced := filepath.Join(root, "a b")
	os.Mkdir(spaced, 0755)
	if got, err := Resolve(" " + spaced + " "); err != nil || got != spaced {
		t.Fatal(got, err)
	}
	if _, err := Resolve(""); err == nil {
		t.Fatal("Accepted empty")
	}
	if _, err := Resolve(filepath.Join(root, "nope")); err == nil {
		t.Fatal("Accepted missing")
	}
}

func TestAgentsStates(t *testing.T) {
	dir := t.TempDir()
	if text, err := Instructions(dir); text != "" || err != nil || Exists(dir) {
		t.Fatalf("missing file must be no instructions: %q %v", text, err)
	}
	os.WriteFile(AgentsPath(dir), []byte("  Be kind <b>\n"), 0644)
	text, err := Instructions(dir)
	if err != nil || !strings.Contains(text, "Be kind <b>") || !strings.Contains(text, "<workspace-instructions path=") || !Exists(dir) {
		t.Fatalf("%q %v", text, err)
	}
	os.WriteFile(AgentsPath(dir), []byte("changed"), 0644)
	if text, _ := Instructions(dir); !strings.Contains(text, "changed") {
		t.Fatal("Edits were not picked up")
	}
	os.WriteFile(AgentsPath(dir), []byte(strings.Repeat("x", 70<<10)), 0644)
	if _, err := Instructions(dir); err == nil || !strings.Contains(err.Error(), "larger") {
		t.Fatalf("Oversized file not reported: %v", err)
	}
	os.Remove(AgentsPath(dir))
	os.Mkdir(AgentsPath(dir), 0755)
	if _, err := Instructions(dir); err == nil || !strings.Contains(err.Error(), "directory") {
		t.Fatalf("Directory not reported: %v", err)
	}
	if Exists(dir) {
		t.Fatal("A directory is not an AGENTS.md")
	}
}

func TestUnreadableAgentsIsReported(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root ignores permissions")
	}
	dir := t.TempDir()
	os.WriteFile(AgentsPath(dir), []byte("secret"), 0000)
	if _, err := Instructions(dir); err == nil || !strings.Contains(err.Error(), "Cannot read") {
		t.Fatalf("%v", err)
	}
}

func TestInstructionsEscapeThePath(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "a\"<b>")
	os.Mkdir(dir, 0755)
	os.WriteFile(AgentsPath(dir), []byte("x"), 0644)
	text, _ := Instructions(dir)
	if strings.Contains(text, "a\"<b>") {
		t.Fatalf("Path was not escaped: %s", text)
	}
}
