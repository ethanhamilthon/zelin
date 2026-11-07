package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetadata(t *testing.T) {
	t.Setenv("OPENAI_CONTEXT_WINDOW", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Error("Missing authentication")
		}
		w.Write([]byte(`{"data":[{"id":"other","context_length":4096},{"id":"test","context_length":128000}]}`))
	}))
	defer server.Close()
	client := &Client{URL: server.URL + "/v1/chat/completions", Key: "secret", Model: "test", HTTP: server.Client()}
	info := client.Metadata(context.Background())
	if info.Model != "test" || info.ContextWindow != 128000 {
		t.Fatalf("%+v", info)
	}
	t.Setenv("OPENAI_CONTEXT_WINDOW", "256000")
	if info := client.Metadata(context.Background()); info.ContextWindow != 256000 {
		t.Fatalf("%+v", info)
	}
}

func TestUnknownContextIsNotInvented(t *testing.T) {
	t.Setenv("OPENAI_CONTEXT_WINDOW", "")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"data":[{"id":"test"}]}`))
	}))
	defer server.Close()
	client := &Client{URL: server.URL + "/chat/completions", Model: "test", HTTP: server.Client()}
	if info := client.Metadata(context.Background()); info.ContextWindow != 0 {
		t.Fatalf("%+v", info)
	}
}
