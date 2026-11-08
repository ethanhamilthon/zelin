package agent

import (
	"context"
	"testing"

	"zelin/internal/provider"
)

type loopClient struct {
	calls  int
	cancel context.CancelFunc
}

func (c *loopClient) Complete(ctx context.Context, _ provider.Request) (provider.Message, provider.Usage, error) {
	if _, ok := ctx.Deadline(); ok {
		panic("Unexpected turn deadline")
	}
	c.calls++
	if c.cancel != nil {
		c.cancel()
	} else if c.calls == 21 {
		return provider.Message{Content: "Done"}, provider.Usage{}, nil
	}
	call := provider.ToolCall{ID: "loop", Type: "function"}
	call.Function.Name = "js_repl"
	call.Function.Arguments = `{"code":"1 + 1"}`
	return provider.Message{ToolCalls: []provider.ToolCall{call}}, provider.Usage{}, nil
}

func TestRunWithoutStepLimit(t *testing.T) {
	client := &loopClient{}
	text, err := New(client).Run(context.Background(), []provider.Message{{Role: "user", Content: "Go"}})
	if err != nil || text != "Done" || client.calls != 21 {
		t.Fatalf("Unexpected result after %d calls: %q, %v", client.calls, text, err)
	}
}

func TestUnlimitedRunCanBeCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := &loopClient{cancel: cancel}
	_, err := New(client).Run(ctx, []provider.Message{{Role: "user", Content: "Go"}})
	if err != context.Canceled || client.calls != 1 {
		t.Fatalf("Unexpected cancellation: %v, %d calls", err, client.calls)
	}
}
