package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"zelin/internal/provider"
)

type contextClient struct {
	requests []provider.Request
	fail     bool
}

func (c *contextClient) Complete(_ context.Context, r provider.Request) (provider.Message, provider.Usage, error) {
	c.requests = append(c.requests, r)
	if c.fail {
		return provider.Message{}, provider.Usage{}, fmt.Errorf("offline")
	}
	text := "done"
	if len(r.Tools) == 0 {
		text = "Preserved goal and exact file paths"
	}
	return provider.Message{Role: "assistant", Content: text}, provider.Usage{}, nil
}

func TestManualCompactionSurvivesRestore(t *testing.T) {
	c := &contextClient{}
	a := New(c)
	history := []provider.Message{{Role: "user", Content: strings.Repeat("large history ", 100)}}
	if err := a.Compact(context.Background(), history, func(Event) {}); err != nil {
		t.Fatal(err)
	}
	state := a.State()
	if len(state.Transcript) != 2 || !strings.Contains(state.Transcript[1].Content, "Preserved goal") {
		t.Fatal(state)
	}
	b := New(c)
	b.Restore(state)
	_, err := b.Run(context.Background(), append(history, provider.Message{Role: "user", Content: "continue"}))
	if err != nil {
		t.Fatal(err)
	}
	request := c.requests[1]
	if len(request.Messages) != 4 || !strings.Contains(request.Messages[1].Content, "Preserved goal") {
		t.Fatal(request)
	}
}

func TestAutoCompactionAndToggle(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		c := &contextClient{}
		a := New(c)
		a.AutoCompact = func() (bool, int, int) { return enabled, 60, 1000 }
		_, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: strings.Repeat("x", 3000)}})
		if err != nil {
			t.Fatal(err)
		}
		expected := 1
		if enabled {
			expected = 2
		}
		if len(c.requests) != expected {
			t.Fatalf("enabled=%v requests=%d", enabled, len(c.requests))
		}
	}
}

func TestFailedCompactionLeavesStateUntouched(t *testing.T) {
	c := &contextClient{fail: true}
	a := New(c)
	a.transcript = []provider.Message{a.system(), {Role: "user", Content: "original"}}
	if err := a.Compact(context.Background(), nil, func(Event) {}); err == nil {
		t.Fatal("expected failure")
	}
	if a.transcript[1].Content != "original" {
		t.Fatal("lost original transcript")
	}
}

func TestSteerAfterWholeToolBatch(t *testing.T) {
	call := provider.ToolCall{ID: "one", Type: "function"}
	call.Function.Name, call.Function.Arguments = "js_repl", `{"code":"1+1"}`
	second := call
	second.ID = "two"
	c := &fakeClient{replies: []provider.Message{{ToolCalls: []provider.ToolCall{call, second}}, {Content: "done"}}}
	a := New(c)
	seen := 0
	_, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "start"}}, func(e Event) {
		if e.Execution != nil && e.Execution.State == "running" {
			a.Steer(provider.Message{Role: "user", Content: "steer"})
		}
		if e.Kind == "steer" {
			seen++
		}
	})
	if err != nil || seen != 1 {
		t.Fatalf("%v seen=%d", err, seen)
	}
	messages := c.messages[1]
	if len(messages) != 6 || messages[3].Role != "tool" || messages[4].Role != "tool" || messages[5].Content != "steer" {
		t.Fatal(messages)
	}
	if len(a.State().History) != 3 {
		t.Fatal(a.State().History)
	}
}

func TestSteerCancellation(t *testing.T) {
	a := New(&contextClient{})
	a.Steer(provider.Message{Role: "user", Content: "queued"})
	a.Steer(provider.Message{})
	if message, _ := a.takeSteer(func(Event) { t.Fatal("unexpected event") }); message != nil {
		t.Fatal("steer was not removed")
	}
}
