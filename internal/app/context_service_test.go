package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
)

func TestHandoffUsesTitleModelAndFullVisibleHistory(t *testing.T) {
	var payload struct {
		Model    string
		Messages []provider.Message
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte(`{"data":[]}`))
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Error(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"Continue editing example.go; tests passed"}}]}`))
	}))
	defer server.Close()
	store, err := sessions.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	messages := []sessions.Message{{Role: "user", Content: "Fix example.go"}, {Role: "assistant", Content: "Tested", Executions: []agent.Execution{{Code: "go test", Output: "PASS", State: "complete"}}}}
	if err := store.Save(session.ID, messages, "draft", nil); err != nil {
		t.Fatal(err)
	}
	store.SetSetting("titleModel", "small")
	client := &provider.Client{URL: server.URL + "/chat/completions", Key: "test", Model: "large", HTTP: server.Client()}
	service := &AgentService{store: store, client: client, agents: map[string]*agent.Agent{}}
	next, err := service.Handoff(context.Background(), session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if payload.Model != "small" || !strings.Contains(payload.Messages[1].Content, "PASS") {
		t.Fatal(payload)
	}
	if next.ID == session.ID || len(next.Messages) != 1 || !strings.Contains(next.Messages[0].Content, "example.go") {
		t.Fatal(next)
	}
	original, _ := store.Get(session.ID)
	if len(original.Messages) != 2 || original.Draft != "draft" {
		t.Fatal("original changed")
	}
	active, _ := store.Current()
	if active.ID != next.ID {
		t.Fatal("new session not active")
	}
}

func TestCompactionDefaults(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	t.Setenv("OPENAI_CONTEXT_WINDOW", "10000")
	client := &provider.Client{Model: "test", HTTP: http.DefaultClient}
	service := &AgentService{store: store, client: client}
	enabled, threshold, window := service.compactionSettings("")
	if !enabled || threshold != 60 || window != 10000 {
		t.Fatal(enabled, threshold, window)
	}
	store.SetSetting("ui", `{"autoCompact":false,"compactThreshold":80}`)
	enabled, threshold, _ = service.compactionSettings("")
	if enabled || threshold != 80 {
		t.Fatal(enabled, threshold)
	}
}
