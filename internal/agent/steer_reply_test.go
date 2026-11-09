package agent

import (
	"context"
	"testing"

	"zelin/internal/provider"
)

func TestReplyBeforeSteerIsShown(t *testing.T) {
	client := &steerClient{replies: []provider.Message{{Role: "assistant", Content: "first answer"}, {Role: "assistant", Content: "second answer"}}, onCall: map[int]func(){}}
	a := New(client)
	client.onCall[1] = func() { a.Steer(provider.Message{Role: "user", Content: "steer"}) }
	var order []string
	_, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: "go"}}, func(event Event) {
		if event.Execution != nil && event.Execution.Kind == "message" {
			order = append(order, "note:"+event.Execution.Code)
		}
		if event.Kind == "steer" {
			order = append(order, "steer")
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 || order[0] != "note:first answer" || order[1] != "steer" {
		t.Fatalf("Reply before steer was lost: %v", order)
	}
}
