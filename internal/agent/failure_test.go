package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"zelin/internal/provider"
)

func TestFailureKeepsToolHistoryAndRuntime(t *testing.T) {
	call := func(id, code string) provider.ToolCall {
		c := provider.ToolCall{ID: id, Type: "function"}
		c.Function.Name = "js_repl"
		c.Function.Arguments = `{"code":"` + code + `"}`
		return c
	}
	client := &fakeClient{replies: []provider.Message{{Role: "assistant", ToolCalls: []provider.ToolCall{call("one", "var kept = 41; kept")}}}}
	a := New(client)
	if _, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Start"}}); err == nil {
		t.Fatal("Expected provider failure")
	}
	client.replies = []provider.Message{{Role: "assistant", ToolCalls: []provider.ToolCall{call("two", "kept + 1")}}, {Role: "assistant", Content: "Done"}}
	var outputs []string
	_, err := a.Run(context.Background(), []provider.Message{
		{Role: "user", Content: "Start"},
		{Role: "assistant", Content: StoppedMessage},
		{Role: "user", Content: "Continue"},
	}, func(event Event) {
		if event.Execution != nil && event.Execution.State == "complete" {
			outputs = append(outputs, event.Execution.Output)
		}
	})
	if err != nil || len(outputs) != 1 || outputs[0] != "42" {
		t.Fatalf("JavaScript state was reset: %v %v", outputs, err)
	}
	sent := client.messages[len(client.messages)-2]
	if len(sent) != 7 || sent[2].ToolCalls[0].ID != "one" || sent[3].Role != "tool" || sent[3].Content != "41" || sent[4].Role != "system" || !strings.Contains(sent[4].Content, "Offline") || sent[5].Content != StoppedMessage {
		t.Fatalf("Tool history was dropped: %+v", sent)
	}
}

func TestWorkspaceInstructionsAreAppliedPerTurn(t *testing.T) {
	client := &fakeClient{replies: []provider.Message{{Role: "assistant", Content: "one"}, {Role: "assistant", Content: "two"}}}
	a := New(client)
	text := "Rule A"
	a.Extra = func() (string, error) { return text, nil }
	user := []provider.Message{{Role: "user", Content: "Hi"}}
	if _, err := a.Run(context.Background(), user); err != nil {
		t.Fatal(err)
	}
	text = "Rule B"
	if _, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Hi"}, {Role: "assistant", Content: "one"}, {Role: "user", Content: "Again"}}); err != nil {
		t.Fatal(err)
	}
	first, second := client.messages[0][0].Content, client.messages[1][0].Content
	if !strings.HasSuffix(first, "Rule A") || !strings.HasSuffix(second, "Rule B") || strings.Contains(second, "Rule A") {
		t.Fatalf("%q %q", first, second)
	}
}

func TestUnreadableWorkspaceInstructionsStopTheTurn(t *testing.T) {
	client := &fakeClient{replies: []provider.Message{{Role: "assistant", Content: "x"}}}
	a := New(client)
	a.Extra = func() (string, error) { return "", fmt.Errorf("AGENTS.md is larger than 64 KB") }
	_, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Hi"}})
	if err == nil || !strings.Contains(err.Error(), "AGENTS.md") || len(client.messages) != 0 {
		t.Fatalf("%v %d", err, len(client.messages))
	}
}

func TestCwdPersistsInState(t *testing.T) {
	dir := t.TempDir()
	a := NewIn(&fakeClient{}, dir)
	if a.State().Cwd != dir {
		t.Fatal(a.State().Cwd)
	}
	restored := New(&fakeClient{})
	restored.Restore(State{Cwd: dir})
	if restored.State().Cwd != dir {
		t.Fatal(restored.State().Cwd)
	}
}

func TestFailureWithoutStepsKeepsTranscriptClean(t *testing.T) {
	client := &fakeClient{}
	a := New(client)
	if _, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Hi"}}); err == nil {
		t.Fatal("Expected provider failure")
	}
	client.replies = []provider.Message{{Role: "assistant", Content: "Hello"}}
	if _, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "Hi"}}); err != nil {
		t.Fatal(err)
	}
	if sent := client.messages[1]; len(sent) != 2 {
		t.Fatalf("Failed user message was duplicated: %+v", sent)
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
