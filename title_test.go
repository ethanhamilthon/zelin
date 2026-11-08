package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"zelin/internal/provider"
	"zelin/internal/sessions"
)

func TestCleanTitleCapsAtSixWordsAndStripsFormatting(t *testing.T) {
	got := cleanTitle("  \"**Скорая помощь в Казахстане сегодня вечером быстро**\"\nextra line  ")
	if got != "Скорая помощь в Казахстане сегодня вечером" {
		t.Fatalf("unexpected title: %q", got)
	}
}

func TestShouldTitleTwice(t *testing.T) {
	u := provider.Message{Role: "user", Content: "hi"}
	a := provider.Message{Role: "assistant", Content: "hello"}
	if !shouldTitle([]provider.Message{u}) {
		t.Fatal("expected first turn to title")
	}
	if !shouldTitle([]provider.Message{u, a, u}) {
		t.Fatal("expected follow-up after first answer to title")
	}
	if shouldTitle([]provider.Message{u, a, u, a, u}) {
		t.Fatal("expected third turn not to retitle")
	}
}

func TestTitleWrapsConversationDataAndUsesUtilityModel(t *testing.T) {
	var mu sync.Mutex
	var seenModel, seenUser string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Model    string
			Messages []provider.Message
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		seenModel = body.Model
		if len(body.Messages) > 1 {
			seenUser = body.Messages[1].Content
		}
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]string{"content": "**Аудит старого туду в zash.**"}}}})
	}))
	defer server.Close()
	store, err := sessions.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	_ = store.SetSetting("titleModel", "mini-title")
	notified := ""
	service := &AgentService{
		store:  store,
		client: &provider.Client{URL: server.URL + "/chat/completions", Key: "k", Model: "main", HTTP: server.Client()},
		notify: func(name string, data any) {
			if name == "session:title" {
				notified = data.(map[string]string)["title"]
			}
		},
	}
	service.title(session.ID, []provider.Message{
		{Role: "user", Content: "Привет! Открой ~/Projects/zash", Images: []string{"data:image/png;base64,AAA"}},
		{Role: "assistant", Content: "Сейчас открою"},
		{Role: "user", Content: "Что там по туду?"},
	})
	mu.Lock()
	defer mu.Unlock()
	if seenModel != "mini-title" || !strings.Contains(seenUser, "<conversation>") || strings.Contains(seenUser, "base64") {
		t.Fatalf("model=%q user=%q", seenModel, seenUser)
	}
	if notified != "Аудит старого туду в zash" {
		t.Fatalf("notified=%q", notified)
	}
}
