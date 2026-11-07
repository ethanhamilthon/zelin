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
			Model     string    `json:"model"`
			Messages  []Message `json:"messages"`
			Reasoning string    `json:"reasoning_effort"`
			Tools     []struct {
				Function struct {
					Name string `json:"name"`
				} `json:"function"`
			} `json:"tools"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Model != "test" || body.Reasoning != "high" || len(body.Messages) != 1 || len(body.Tools) != 1 || body.Tools[0].Function.Name != "js_repl" {
			t.Errorf("Invalid body: %+v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Hello","reasoning_content":"Hmm","reasoning_details":[{"type":"reasoning.text","signature":"sig"}]}}],"usage":{"prompt_tokens":100,"completion_tokens":20,"cost":0.5,"prompt_tokens_details":{"cached_tokens":40}}}`))
	}))
	defer server.Close()
	client := &Client{URL: server.URL, Key: "secret", Model: "test", Effort: "high", HTTP: server.Client()}
	reply, usage, err := client.Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "Hi"}}, Tools: []Tool{{Name: "js_repl"}}})
	if err != nil || reply.Content != "Hello" || reply.Thinking() != "Hmm" || !strings.Contains(string(reply.ReasoningDetails), "sig") {
		t.Fatalf("%+v %v", reply, err)
	}
	if *usage.Input != 100 || *usage.Output != 20 || *usage.CacheRead != 40 || usage.CacheWrite != nil || *usage.Cost != 0.5 || usage.Estimated {
		t.Fatalf("%+v", usage)
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
			_, _, err := (&Client{URL: server.URL, Key: "secret", Model: "test", HTTP: server.Client()}).Complete(context.Background(), Request{})
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatal(err)
			}
		})
	}
}

func TestImagesBecomeContentParts(t *testing.T) {
	message := Message{Role: "user", Content: "See [Image #2] and [Image #1] or [Image #9]", Images: []string{"data:a", "data:b"}}
	data, _ := json.Marshal(wire([]Message{message, {Role: "user", Content: "plain"}}))
	want := `[{"role":"user","content":[{"text":"See [Image #2]","type":"text"},{"image_url":{"url":"data:b"},"type":"image_url"},{"text":" and [Image #1]","type":"text"},{"image_url":{"url":"data:a"},"type":"image_url"},{"text":" or [Image #9]","type":"text"}]},{"role":"user","content":"plain"}]`
	if string(data) != want {
		t.Fatalf("%s", data)
	}
}
