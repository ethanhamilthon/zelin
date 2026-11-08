package main

import (
	"path/filepath"
	"testing"
	"zelin/internal/sessions"
)

func TestToolOutputSettings(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	service := &AgentService{store: store}
	for _, item := range []struct {
		raw  string
		want int
	}{
		{"", 4096}, {"{}", 4096}, {"{\"toolOutputLimit\":0}", 0},
		{"{\"toolOutputLimit\":1024}", 1024}, {"{\"toolOutputLimit\":16384}", 16384},
		{"{\"toolOutputLimit\":-1}", 4096},
	} {
		if err := store.SetSetting("ui", item.raw); err != nil {
			t.Fatal(err)
		}
		if got := service.toolOutputLimit(); got != item.want {
			t.Fatalf("%s: %d", item.raw, got)
		}
	}
}
