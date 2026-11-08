package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("FAKE_MCP") == "1" {
		fake()
		return
	}
	os.Exit(m.Run())
}

func fake() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req struct {
			ID     *int64 `json:"id"`
			Method string `json:"method"`
			Params struct {
				Name      string         `json:"name"`
				Arguments map[string]any `json:"arguments"`
			} `json:"params"`
		}
		json.Unmarshal(scanner.Bytes(), &req)
		if req.ID == nil {
			continue
		}
		result := map[string]any{}
		switch req.Method {
		case "tools/list":
			result["tools"] = []map[string]any{{"name": "echo", "description": "Echo text"}}
		case "tools/call":
			if req.Params.Name == "fail" {
				result = map[string]any{"isError": true, "content": []map[string]string{{"type": "text", "text": "boom"}}}
			} else {
				result["content"] = []map[string]string{{"type": "text", "text": fmt.Sprint("echo:", req.Params.Arguments["text"], ":", os.Getenv("TOKEN"))}}
			}
		}
		data, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": *req.ID, "result": result})
		fmt.Println(string(data))
	}
}

func TestManagerListsAndCallsTools(t *testing.T) {
	exe, _ := os.Executable()
	servers := []Server{{Name: "fake", Command: exe, Env: map[string]string{"FAKE_MCP": "1", "TOKEN": "t1"}}}
	m := &Manager{Servers: func() []Server { return servers }}
	defer m.Reset()
	ctx := context.Background()
	tools, err := m.Tools(ctx, "fake", t.TempDir())
	if err != nil || len(tools) != 1 || tools[0].Name != "echo" {
		t.Fatalf("%+v %v", tools, err)
	}
	out, err := m.Call(ctx, "fake", "", "echo", map[string]any{"text": "hi"})
	if err != nil || out != "echo:hi:t1" {
		t.Fatalf("%q %v", out, err)
	}
	if _, err := m.Call(ctx, "fake", "", "fail", nil); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatal(err)
	}
	if _, err := m.Tools(ctx, "missing", ""); err == nil {
		t.Fatal("Missing server accepted")
	}
}
