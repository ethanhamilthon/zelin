package agent

import (
	"context"
	"testing"

	"zelin/internal/provider"
)

func TestInterimAssistantTextBecomesStep(t *testing.T) {
	call := provider.ToolCall{ID: "c", Type: "function"}
	call.Function.Name = "js_repl"
	call.Function.Arguments = `{"code":"1"}`
	client := &fakeClient{replies: []provider.Message{
		{Role: "assistant", Content: "Checking files", ToolCalls: []provider.ToolCall{call}},
		{Role: "assistant", Content: "Done"},
	}}
	var notes []string
	_, err := New(client).Run(context.Background(), []provider.Message{{Role: "user", Content: "Go"}}, func(event Event) {
		if event.Execution != nil && event.Execution.Kind == "message" {
			notes = append(notes, event.Execution.Code)
		}
	})
	if err != nil || len(notes) != 1 || notes[0] != "Checking files" {
		t.Fatalf("%v %v", notes, err)
	}
}
