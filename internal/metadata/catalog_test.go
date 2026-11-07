package metadata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestKeyNormalizesProviderIDs(t *testing.T) {
	for id, want := range map[string]string{
		"anthropic/claude-sonnet-4.6":        "claude-sonnet-4-6",
		"claude-opus-4-5-20251101":           "claude-opus-4-5",
		"gemini-3.8-flash-high":              "gemini-3-8-flash",
		"openai/gpt-6.1-terra:free":          "gpt-6-1-terra",
		"bedrock/claude-3-7-sonnet-20250219": "claude-3-7-sonnet",
	} {
		if got := Key(id); got != want {
			t.Errorf("Key(%q) = %q, want %q", id, got, want)
		}
	}
}

func TestCatalogMergesSourcesAndCaches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/or" {
			w.Write([]byte(`{"data":[{"id":"x/model-a:batch","context_length":1,"pricing":{"prompt":"9"}},{"id":"x/model-a","context_length":1000,"pricing":{"prompt":"0.000001","completion":"0.000002"}}]}`))
			return
		}
		w.Write([]byte(`{"model-a":{"max_input_tokens":5},"vertex/model-b":{"max_input_tokens":7},"model-b":{"max_input_tokens":2000,"input_cost_per_token":0.5}}`))
	}))
	defer server.Close()
	catalog := New(filepath.Join(t.TempDir(), "metadata.json"))
	catalog.Sources = [2]string{server.URL + "/or", server.URL + "/ll"}
	a, ok := catalog.Lookup(context.Background(), "model-a")
	if !ok || a.ContextWindow != 1000 || *a.Pricing.Completion != 0.000002 {
		t.Fatalf("model-a = %+v", a)
	}
	b, _ := catalog.Lookup(context.Background(), "model-b-20250101")
	if b.ContextWindow != 2000 || *b.Pricing.Prompt != 0.5 {
		t.Fatalf("model-b = %+v", b)
	}
	server.Close()
	cached := New(catalog.Path)
	if entry, ok := cached.Lookup(context.Background(), "model-a"); !ok || entry.ContextWindow != 1000 {
		t.Fatalf("cache = %+v", entry)
	}
}
