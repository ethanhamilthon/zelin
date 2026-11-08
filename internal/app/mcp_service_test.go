package app

import (
	"path/filepath"
	"testing"

	"zelin/internal/mcp"
	"zelin/internal/sessions"
)

func TestMCPServiceAddRemove(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	s := &MCPService{store: store}
	if len(s.List()) != 0 {
		t.Fatal("Expected no servers")
	}
	if _, err := s.Add(mcp.Server{Name: " ", Command: "x"}); err == nil {
		t.Fatal("Accepted empty name")
	}
	if list, err := s.Add(mcp.Server{Name: "fs", Command: "npx", Args: []string{"-y", "server"}}); err != nil || len(list) != 1 {
		t.Fatalf("%v %v", list, err)
	}
	if _, err := s.Add(mcp.Server{Name: "fs", Command: "other"}); err == nil {
		t.Fatal("Accepted duplicate")
	}
	if got := s.List(); len(got) != 1 || got[0].Args[1] != "server" {
		t.Fatalf("%+v", got)
	}
	if list, err := s.Remove("fs"); err != nil || len(list) != 0 {
		t.Fatalf("%v %v", list, err)
	}
}
