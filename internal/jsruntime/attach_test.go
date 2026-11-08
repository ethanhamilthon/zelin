package jsruntime

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAttachResolvesValidatesAndDrains(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "report.md"), []byte("Report"), 0600)
	runtime := NewIn(dir)
	result := runtime.Run(context.Background(), "attach('Report', 'report.md')")
	if !strings.Contains(result, "true") {
		t.Fatal(result)
	}
	runtime.Run(context.Background(), "attach('Updated report', 'report.md')")
	items := runtime.TakeAttachments()
	if len(items) != 1 || items[0].Path != filepath.Join(dir, "report.md") || items[0].Name != "Updated report" {
		t.Fatal(items)
	}
	if len(runtime.TakeAttachments()) != 0 {
		t.Fatal("Attachments leaked into next turn")
	}
	for _, code := range []string{"attach('', 'report.md')", "attach('Missing', 'absent')", "attach('Directory', '.')"} {
		if result := runtime.Run(context.Background(), code); !strings.Contains(result, "error") {
			t.Fatal(result)
		}
	}
	if len(runtime.TakeAttachments()) != 0 {
		t.Fatal("Invalid file attached")
	}
}
