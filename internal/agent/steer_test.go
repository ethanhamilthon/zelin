package agent

import (
	"context"
	"testing"

	"zelin/internal/provider"
)

type steerClient struct {
	replies  []provider.Message
	requests [][]provider.Message
	onCall   map[int]func()
}

func (c *steerClient) Complete(_ context.Context, r provider.Request) (provider.Message, provider.Usage, error) {
	c.requests = append(c.requests, append([]provider.Message{}, r.Messages...))
	if f := c.onCall[len(c.requests)]; f != nil {
		f()
	}
	reply := c.replies[0]
	c.replies = c.replies[1:]
	return reply, provider.Usage{}, nil
}

func toolCall(id, code string) provider.Message {
	c := provider.ToolCall{ID: id, Type: "function"}
	c.Function.Name = "js_repl"
	c.Function.Arguments = `{"code":"` + code + `"}`
	return provider.Message{Role: "assistant", ToolCalls: []provider.ToolCall{c}}
}

func TestSteerKeepsPriorActionsAcrossTurns(t *testing.T) {
	client := &steerClient{replies: []provider.Message{toolCall("1", "var x = 5; x"), {Role: "assistant", Content: "done"}, {Role: "assistant", Content: "after steer"}, {Role: "assistant", Content: "final"}}, onCall: map[int]func(){}}
	a := New(client)
	client.onCall[2] = func() { a.Steer(provider.Message{Role: "user", Content: "steer"}) }
	history := []provider.Message{{Role: "user", Content: "go"}}
	if _, err := a.Run(context.Background(), history); err != nil {
		t.Fatal(err)
	}
	history = append(history, provider.Message{Role: "user", Content: "steer"}, provider.Message{Role: "assistant", Content: "after steer"}, provider.Message{Role: "user", Content: "again"})
	if _, err := a.Run(context.Background(), history); err != nil {
		t.Fatal(err)
	}
	last := client.requests[len(client.requests)-1]
	if last[2].ToolCalls[0].ID != "1" || last[3].Role != "tool" || last[3].Content != "5" {
		t.Fatalf("Tool context lost after steering: %+v", last)
	}
}

func TestTranscriptOnlyGrowsAndPersistsEveryStep(t *testing.T) {
	client := &steerClient{replies: []provider.Message{toolCall("1", "1"), toolCall("2", "2"), {Role: "assistant", Content: "done"}}}
	a := New(client)
	var sizes []int
	a.Persist = func(state State) error {
		if n := len(sizes); n > 0 && len(state.Transcript) < sizes[n-1] {
			t.Fatalf("Transcript shrank: %d after %d", len(state.Transcript), sizes[n-1])
		}
		sizes = append(sizes, len(state.Transcript))
		return nil
	}
	if _, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "go"}}); err != nil {
		t.Fatal(err)
	}
	if len(sizes) < 4 {
		t.Fatalf("Steps were not persisted incrementally: %v", sizes)
	}
}

func TestRejectedTurnIsRecordedNotWiped(t *testing.T) {
	client := &steerClient{replies: []provider.Message{{Role: "assistant", Content: "one"}}}
	a := New(client)
	history := []provider.Message{{Role: "user", Content: "a"}}
	if _, err := a.Run(context.Background(), history); err != nil {
		t.Fatal(err)
	}
	a.Reject(append(append([]provider.Message{}, history...), provider.Message{Role: "assistant", Content: "one"}, provider.Message{Role: "user", Content: "b"}), context.Canceled)
	state := a.State()
	if len(state.History) != 4 || state.History[3].Content != StoppedMessage || state.Transcript[len(state.Transcript)-2].Role != "system" {
		t.Fatalf("Rejected turn not recorded: %+v", state)
	}
}
