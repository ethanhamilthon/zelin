package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

const catalog = `{"data":[
{"id":"plain","name":"Plain","context_length":4096,"pricing":{"prompt":"0.000001","completion":"0.000002"}},
{"id":"open","top_provider":{"context_length":128000},"pricing":{"prompt":"-1"},"reasoning":{"supported_efforts":null}},
{"id":"strict","context_length":400000,"reasoning":{"supported_efforts":["high","low","none"],"default_effort":"low","mandatory":true}}]}`

func TestModels(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/v1/models" || r.Header.Get("Authorization") != "Bearer secret" {
			t.Errorf("Unexpected request: %s", r.URL.Path)
		}
		w.Write([]byte(catalog))
	}))
	defer server.Close()
	client := &Client{URL: server.URL + "/v1/chat/completions", Key: "secret", HTTP: server.Client()}
	plain, err := client.Lookup(context.Background(), "plain")
	if err != nil || plain.Name != "Plain" || plain.ContextWindow != 4096 || *plain.Pricing.Prompt != 0.000001 || plain.Pricing.CacheRead != nil || len(plain.Efforts) != 0 {
		t.Fatalf("%+v %v", plain, err)
	}
	open, _ := client.Lookup(context.Background(), "open")
	if open.Name != "open" || open.ContextWindow != 128000 || open.Pricing.Prompt != nil || !slices.Equal(open.Efforts, gatewayEfforts) {
		t.Fatalf("%+v", open)
	}
	strict, _ := client.Lookup(context.Background(), "strict")
	if !slices.Equal(strict.Efforts, []string{"high", "low"}) || strict.DefaultEffort != "low" {
		t.Fatalf("%+v", strict)
	}
	if _, err := client.Lookup(context.Background(), "missing"); err == nil || calls != 1 {
		t.Fatalf("%v %d", err, calls)
	}
}

func TestModelListFailureIsNotCached(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusBadGateway) }))
	defer server.Close()
	client := &Client{URL: server.URL + "/chat/completions", HTTP: server.Client()}
	if _, err := client.Models(context.Background()); err == nil || client.cached() != nil {
		t.Fatal("Expected an uncached failure")
	}
}
