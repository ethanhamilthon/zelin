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

func modelService(t *testing.T) (*AgentService, *sessions.Store) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"id":"fast"},{"id":"deep"},{"id":"third"}]}`))
	}))
	t.Cleanup(server.Close)
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	client := &provider.Client{URL: server.URL + "/chat/completions", Model: "fast", Effort: "medium", HTTP: server.Client()}
	return &AgentService{store: store, client: client}, store
}

func TestConfigurePersistsSelection(t *testing.T) {
	service, store := modelService(t)
	session, _ := store.Create()
	ctx := context.Background()
	info, err := service.Configure(ctx, session.ID, "deep", "xhigh")
	if err != nil || info.Model != "deep" || info.Effort != "xhigh" || !slices.Equal(info.Efforts, provider.Efforts) {
		t.Fatalf("%+v %v", info, err)
	}
	if info, _ = service.Configure(ctx, session.ID, "fast", "turbo"); info.Effort != provider.DefaultEffort {
		t.Fatalf("Unknown effort was kept: %+v", info)
	}
	if _, err := service.Configure(ctx, session.ID, "missing", ""); err == nil {
		t.Fatal("Accepted an unknown model")
	}
	if info, err = service.SetTitleModel(ctx, "fast"); err != nil || info.TitleModel != "fast" {
		t.Fatalf("%+v %v", info, err)
	}
	service.Configure(ctx, session.ID, "deep", "none")
	restored := &provider.Client{Model: "fast"}
	if err := restoreSelection(store, restored); err != nil || restored.Model != "deep" || restored.Effort != "none" {
		t.Fatalf("%+v %v", restored, err)
	}
}

func TestModelIsPerSession(t *testing.T) {
	service, store := modelService(t)
	ctx := context.Background()
	first, _ := store.Create()
	store.Save(first.ID, []sessions.Message{{Role: "user", Content: "one"}}, "", nil)
	second, _ := store.Create()
	store.Save(second.ID, []sessions.Message{{Role: "user", Content: "two"}}, "", nil)
	service.pin(second.ID)
	if _, err := service.Configure(ctx, second.ID, "deep", "high"); err != nil {
		t.Fatal(err)
	}
	if model, _ := service.selection(first.ID); model != "fast" {
		t.Fatalf("first session changed to %s", model)
	}
	if model, effort := service.selection(second.ID); model != "deep" || effort != "high" {
		t.Fatal(model, effort)
	}
	if model, _ := service.client.Selection(); model != "fast" {
		t.Fatalf("existing session changed the default to %s", model)
	}
	blank, _ := store.Create()
	if _, err := service.Configure(ctx, blank.ID, "third", "low"); err != nil {
		t.Fatal(err)
	}
	if model, _ := service.client.Selection(); model != "third" {
		t.Fatalf("blank session did not change the default: %s", model)
	}
	if model, _ := service.selection(first.ID); model != "fast" {
		t.Fatalf("legacy session followed the new default: %s", model)
	}
	if model, _ := service.selection(second.ID); model != "deep" {
		t.Fatalf("second session changed: %s", model)
	}
	next, _ := store.Create()
	if model, effort := service.selection(next.ID); model != "third" || effort != "low" || service.Info(ctx, next.ID).Model != "third" {
		t.Fatal(model, effort)
	}
}
