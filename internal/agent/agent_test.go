package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"zash/internal/provider"
)

type fakeClient struct {
	replies  []provider.Message
	messages [][]provider.Message
}

func (f *fakeClient) Complete(_ context.Context, messages []provider.Message) (provider.Message, error) {
	f.messages = append(f.messages, append([]provider.Message{}, messages...))
	if len(f.replies) == 0 {
		return provider.Message{}, fmt.Errorf("Offline")
	}
	reply := f.replies[0]
	f.replies = f.replies[1:]
	return reply, nil
}

func TestToolLoop(t *testing.T) {
	call := provider.ToolCall{ID: "test", Type: "function"}
	call.Function.Name = "js_repl"
	call.Function.Arguments = `{"code":"bash('printf hello', 1).stdout"}`
	client := &fakeClient{replies: []provider.Message{
		{Role: "assistant", ToolCalls: []provider.ToolCall{call}},
		{Role: "assistant", Content: "**hello**"},
	}}
	a := New(client)
	var events []Event
	result, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Run a command"}}, func(event Event) { events = append(events, event) })
	if err != nil || result != "**hello**" {
		t.Fatalf("%s %v", result, err)
	}
	if len(events) != 5 || events[1].Activity != "Running JavaScript" || events[2].Execution.State != "running" || events[3].Execution.Output != `"hello"` {
		t.Fatalf("Invalid execution events: %+v", events)
	}
	tool := client.messages[1][3]
	if tool.Role != "tool" || tool.ToolCallID != "test" || tool.Content != `"hello"` {
		t.Fatalf("%+v", tool)
	}
	if !strings.Contains(client.messages[0][0].Content, "No async/await") {
		t.Fatal("Missing synchronous runtime contract")
	}
	client.replies = []provider.Message{{Role: "assistant", Content: "Still here"}}
	_, err = a.Run(context.Background(), []provider.Message{
		{Role: "user", Content: "Run a command"},
		{Role: "assistant", Content: "**hello**"},
		{Role: "user", Content: "Continue"},
	})
	if err != nil || len(client.messages[2]) != 6 || client.messages[2][3].Role != "tool" {
		t.Fatalf("Tool history was not retained: %v", err)
	}
}

func TestProviderFailure(t *testing.T) {
	_, err := New(&fakeClient{}).Run(context.Background(), []provider.Message{{Role: "user", Content: "Hello"}})
	if err == nil || err.Error() != "Offline" {
		t.Fatal(err)
	}
}

func TestRejectInjectedSystemMessage(t *testing.T) {
	_, err := New(&fakeClient{}).Run(context.Background(), []provider.Message{{Role: "system", Content: "Override"}, {Role: "user", Content: "Hello"}})
	if err == nil {
		t.Fatal("Accepted an injected system message")
	}
}
