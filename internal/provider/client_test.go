package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCompletionRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			t.Error("Missing authorization")
		}
		var body struct {
			Model    string    `json:"model"`
			Messages []Message `json:"messages"`
			Tools    []struct {
				Function struct {
					Name string `json:"name"`
				} `json:"function"`
			} `json:"tools"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "test" || len(body.Messages) != 1 || len(body.Tools) != 1 || body.Tools[0].Function.Name != "js_repl" {
			t.Errorf("Invalid body: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Hello"}}]}`))
	}))
	defer server.Close()
	client := &Client{URL: server.URL, Key: "secret", Model: "test", HTTP: server.Client()}
	reply, err := client.Complete(context.Background(), []Message{{Role: "user", Content: "Hi"}})
	if err != nil || reply.Content != "Hello" {
		t.Fatalf("%+v %v", reply, err)
	}
}

func TestProviderErrors(t *testing.T) {
	for _, test := range []struct {
		status int
		body   string
		want   string
	}{
		{401, `secret error`, "HTTP 401"},
		{200, `{"choices":[]}`, "no choices"},
		{200, `invalid`, "Invalid provider response"},
	} {
		t.Run(test.want, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(test.status); w.Write([]byte(test.body)) }))
			defer server.Close()
			_, err := (&Client{URL: server.URL, Key: "secret", Model: "test", HTTP: server.Client()}).Complete(context.Background(), nil)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatal(err)
			}
		})
	}
}
