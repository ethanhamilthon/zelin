package agent

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"zelin/internal/provider"
)

const tinyPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func TestViewImageQueuesImageForNextStep(t *testing.T) {
	dir := t.TempDir()
	imgPath := filepath.Join(dir, "pixel.png")
	txtPath := filepath.Join(dir, "note.txt")
	data, _ := base64.StdEncoding.DecodeString(tinyPNG)
	if err := os.WriteFile(imgPath, data, 0600); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(txtPath, []byte("not an image"), 0600)
	rawArgs, _ := json.Marshal(map[string]string{"code": fmt.Sprintf("var a = view_image(%q); var b = view_image(%q); [a, b]", imgPath, txtPath)})
	call := provider.ToolCall{ID: "img-call", Type: "function"}
	call.Function.Name = "js_repl"
	call.Function.Arguments = string(rawArgs)
	client := &fakeClient{replies: []provider.Message{
		{Role: "assistant", ToolCalls: []provider.ToolCall{call}},
		{Role: "assistant", Content: "Saw the pixel"},
	}}
	a := New(client)
	got, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Inspect image"}})
	if err != nil || got != "Saw the pixel" {
		t.Fatalf("%q %v", got, err)
	}
	if len(client.messages) != 2 || len(client.messages[1]) != 5 {
		t.Fatalf("unexpected messages: %+v", client.messages)
	}
	toolMsg := client.messages[1][3]
	if !strings.Contains(toolMsg.Content, "image/png") || !strings.Contains(toolMsg.Content, "Unsupported image format") {
		t.Fatalf("unexpected tool output: %s", toolMsg.Content)
	}
	injected := client.messages[1][4]
	if injected.Role != "user" || !strings.Contains(injected.Content, "[Image #1]") || len(injected.Images) != 1 || !strings.HasPrefix(injected.Images[0], "data:image/png;base64,") {
		t.Fatalf("unexpected injected image message: %+v", injected)
	}
	if len(a.State().History) != 2 {
		t.Fatalf("view_image leaked into visible chat history: %+v", a.State().History)
	}
}
