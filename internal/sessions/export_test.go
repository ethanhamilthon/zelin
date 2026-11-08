package sessions

import (
	"path/filepath"
	"strings"
	"testing"

	"zelin/internal/agent"
)

func TestMarkdownExport(t *testing.T) {
	session := Session{Title: " Plan ", Messages: []Message{
		{Role: "user", Content: "Run it", Images: []string{"data:image/png;base64,AA"}},
		{Role: "assistant", Content: "Done", Executions: []agent.Execution{
			{Kind: "thinking", Code: "Think"},
			{Code: "1+1", Output: "2 ``` inside"},
		}},
	}}
	text := session.Markdown()
	for _, want := range []string{"# Plan\n", "## User", "## Assistant", "```js\n1+1\n```", "````text\n2 ``` inside\n````", "Think", "Done", "1 image attachment(s) not included"} {
		if !strings.Contains(text, want) {
			t.Fatalf("Missing %q in:\n%s", want, text)
		}
	}
	if strings.Contains(text, "base64") {
		t.Fatal("Exported image data")
	}
}

func TestManualTitleBeatsAutomatic(t *testing.T) {
	store, err := Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	if ok, err := store.AutoTitle(session.ID, "Auto"); !ok || err != nil {
		t.Fatal(ok, err)
	}
	if err := store.SetTitle(session.ID, "   "); err == nil {
		t.Fatal("Accepted an empty title")
	}
	if err := store.SetTitle(session.ID, "  Mine  "); err != nil {
		t.Fatal(err)
	}
	if ok, _ := store.AutoTitle(session.ID, "Late auto"); ok {
		t.Fatal("Automatic title overwrote the manual one")
	}
	if got, _ := store.Get(session.ID); got.Title != "Mine" {
		t.Fatal(got.Title)
	}
	if err := store.SetTitle("missing", "x"); err == nil {
		t.Fatal("Renamed a missing session")
	}
}
