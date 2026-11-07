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

func (f *fakeClient) Complete(_ context.Context, request provider.Request) (provider.Message, provider.Usage, error) {
	messages := request.Messages
	if len(request.Tools) != 1 || !strings.Contains(request.Tools[0].Description, "No async/await") {
		return provider.Message{}, provider.Usage{}, fmt.Errorf("Missing synchronous runtime contract")
	}
	f.messages = append(f.messages, append([]provider.Message{}, messages...))
	if len(f.replies) == 0 {
		return provider.Message{}, provider.Usage{}, fmt.Errorf("Offline")
	}
	reply := f.replies[0]
	f.replies = f.replies[1:]
	input, output := int64(10), int64(2)
	return reply, provider.Usage{Input: &input, Output: &output}, nil
}

func TestToolLoop(t *testing.T) {
	call := provider.ToolCall{ID: "test", Type: "function"}
	call.Function.Name = "js_repl"
	call.Function.Arguments = `{"code":"bash('printf hello', 1).stdout"}`
	client := &fakeClient{replies: []provider.Message{
		{Role: "assistant", ToolCalls: []provider.ToolCall{call}},
		{Role: "assistant", Content: "**hello**", ReasoningContent: "Done"},
	}}
	a := New(client)
	var events []Event
	result, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Run a command"}}, func(event Event) { events = append(events, event) })
	if err != nil || result != "**hello**" {
		t.Fatalf("%s %v", result, err)
	}
	if len(events) != 8 || events[7].Execution.Kind != "thinking" || events[7].Execution.Code != "Done" || events[1].Usage.Requests != 1 || events[2].Activity != "Running JavaScript" || events[3].Execution.State != "running" || events[4].Execution.Output != `"hello"` {
		t.Fatalf("Invalid execution events: %+v", events)
	}
	if usage := a.State().Usage; usage.Requests != 2 || *usage.Input != 20 || *usage.Context != 12 {
		t.Fatalf("%+v", usage)
	}
	tool := client.messages[1][3]
	if tool.Role != "tool" || tool.ToolCallID != "test" || tool.Content != `"hello"` {
		t.Fatalf("%+v", tool)
	}
	if client.messages[0][0].Content != DefaultPrompt {
		t.Fatal("Missing system prompt")
	}
	client.replies = []provider.Message{{Role: "assistant", Content: "Still here"}}
	a.Prompt = func() string { return "Edited" }
	_, err = a.Run(context.Background(), []provider.Message{
		{Role: "user", Content: "Run a command"},
		{Role: "assistant", Content: "**hello**"},
		{Role: "user", Content: "Continue"},
	})
	if err != nil || len(client.messages[2]) != 6 || client.messages[2][3].Role != "tool" || client.messages[2][0].Content != "Edited" || client.messages[2][4].ReasoningContent != "Done" {
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
