package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamingCompletion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		for _, line := range []string{
			`{"choices":[{"delta":{"reasoning_content":"Hm","reasoning_details":[{"type":"reasoning.text","index":0,"text":"Hm"}]}}]}`,
			`{"choices":[{"delta":{"reasoning_content":"m","reasoning_details":[{"type":"reasoning.text","index":0,"text":"m","signature":"sig"}]}}]}`,
			`{"choices":[{"delta":{"content":"Hel"}}]}`,
			`{"choices":[{"delta":{"content":"lo","tool_calls":[{"index":0,"id":"c1","function":{"name":"js_repl","arguments":"{\"co"}}]}}]}`,
			`{"choices":[{"delta":{"tool_calls":[{"index":0,"function":{"arguments":"de\":1}"}}]}}]}`,
			`{"choices":[],"usage":{"prompt_tokens":10,"completion_tokens":2}}`,
		} {
			w.Write([]byte("data: " + line + "\n\n"))
		}
		w.Write([]byte("data: [DONE]\n\n"))
	}))
	defer server.Close()
	client := &Client{URL: server.URL, Key: "k", Model: "m", HTTP: server.Client()}
	var deltas []Delta
	reply, usage, err := client.Complete(context.Background(), Request{Messages: []Message{{Role: "user", Content: "Hi"}}, OnDelta: func(d Delta) { deltas = append(deltas, d) }})
	if err != nil || reply.Content != "Hello" || reply.Thinking() != "Hmm" || reply.Role != "assistant" {
		t.Fatalf("%+v %v", reply, err)
	}
	if len(reply.ToolCalls) != 1 || reply.ToolCalls[0].ID != "c1" || reply.ToolCalls[0].Function.Arguments != `{"code":1}` {
		t.Fatalf("%+v", reply.ToolCalls)
	}
	if details := string(reply.ReasoningDetails); !strings.Contains(details, `"text":"Hmm"`) || !strings.Contains(details, "sig") {
		t.Fatal(details)
	}
	if usage.Input == nil || *usage.Input != 10 || len(deltas) != 5 || deltas[2].Content != "Hel" {
		t.Fatalf("%+v %d", usage, len(deltas))
	}
}
