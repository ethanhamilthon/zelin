package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"testing"

	"zash/internal/provider"
	"zash/internal/sessions"
)

func TestConfigurePersistsSelection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"id":"fast"},{"id":"deep"}]}`))
	}))
	defer server.Close()
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	client := &provider.Client{URL: server.URL + "/chat/completions", Model: "fast", HTTP: server.Client()}
	service := &AgentService{store: store, client: client}
	info, err := service.Configure(context.Background(), "deep", "xhigh")
	if err != nil || info.Model != "deep" || info.Effort != "xhigh" || !slices.Equal(info.Efforts, provider.Efforts) {
		t.Fatalf("%+v %v", info, err)
	}
	if info, _ = service.Configure(context.Background(), "fast", "turbo"); info.Effort != provider.DefaultEffort {
		t.Fatalf("Unknown effort was kept: %+v", info)
	}
	if _, err := service.Configure(context.Background(), "missing", ""); err == nil {
		t.Fatal("Accepted an unknown model")
	}
	if info, err = service.SetTitleModel(context.Background(), "fast"); err != nil || info.TitleModel != "fast" {
		t.Fatalf("%+v %v", info, err)
	}
	service.Configure(context.Background(), "deep", "none")
	restored := &provider.Client{Model: "fast"}
	if err := restoreSelection(store, restored); err != nil || restored.Model != "deep" || restored.Effort != "none" {
		t.Fatalf("%+v %v", restored, err)
	}
}
