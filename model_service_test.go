package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"zash/internal/provider"
	"zash/internal/sessions"
)

func TestConfigurePersistsSelection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"id":"fast","context_length":1000},{"id":"deep","context_length":2000,"reasoning":{"supported_efforts":["high","low"]}}]}`))
	}))
	defer server.Close()
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	client := &provider.Client{URL: server.URL + "/chat/completions", Model: "fast", HTTP: server.Client()}
	service := &AgentService{store: store, client: client}
	info, err := service.Configure(context.Background(), "deep", "high")
	if err != nil || info.Model != "deep" || info.Effort != "high" || info.ContextWindow != 2000 || len(info.Efforts) != 2 {
		t.Fatalf("%+v %v", info, err)
	}
	if info, _ = service.Configure(context.Background(), "fast", "high"); info.Effort != "" {
		t.Fatalf("Unsupported effort was kept: %+v", info)
	}
	if _, err := service.Configure(context.Background(), "missing", ""); err == nil {
		t.Fatal("Accepted an unknown model")
	}
	service.Configure(context.Background(), "deep", "low")
	restored := &provider.Client{Model: "fast"}
	if err := restoreSelection(store, restored); err != nil || restored.Model != "deep" || restored.Effort != "low" {
		t.Fatalf("%+v %v", restored, err)
	}
}
