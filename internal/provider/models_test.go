package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"zelin/internal/metadata"
)

func TestModelsUseMetadataCatalog(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/models":
			calls++
			if r.Header.Get("Authorization") != "Bearer secret" {
				t.Error("Missing authorization")
			}
			w.Write([]byte(`{"data":[{"id":"plain","name":"Plain","context_length":1},{"id":"claude-x-20250101"}]}`))
		case "/or":
			w.Write([]byte(`{"data":[{"id":"anthropic/claude-x","context_length":200000,"pricing":{"prompt":"0.000003"}}]}`))
		default:
			w.Write([]byte(`{}`))
		}
	}))
	defer server.Close()
	meta := metadata.New(filepath.Join(t.TempDir(), "meta.json"))
	meta.Sources = [2]string{server.URL + "/or", server.URL + "/ll"}
	client := &Client{URL: server.URL + "/v1/chat/completions", Key: "secret", HTTP: server.Client(), Meta: meta}
	plain, err := client.Lookup(context.Background(), "plain")
	if err != nil || plain.Name != "Plain" || plain.ContextWindow != 0 || plain.Pricing.Prompt != nil {
		t.Fatalf("%+v %v", plain, err)
	}
	claude, _ := client.Lookup(context.Background(), "claude-x-20250101")
	if claude.Name != "claude-x-20250101" || claude.ContextWindow != 200000 || *claude.Pricing.Prompt != 0.000003 {
		t.Fatalf("%+v", claude)
	}
	if _, err := client.Lookup(context.Background(), "missing"); err == nil || calls != 1 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}
