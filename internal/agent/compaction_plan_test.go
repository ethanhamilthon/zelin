package agent

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"zelin/internal/provider"
)

type planClient struct {
	requests  []provider.Request
	compactFn func(req provider.Request, n int) (provider.Message, error)
	replies   []provider.Message
	usages    []provider.Usage
}

func (c *planClient) Complete(ctx context.Context, req provider.Request) (provider.Message, provider.Usage, error) {
	c.requests = append(c.requests, req)
	if len(req.Tools) == 0 {
		if c.compactFn != nil {
			msg, err := c.compactFn(req, len(c.requests))
			return msg, provider.Usage{}, err
		}
		return provider.Message{Role: "assistant", Content: "Summary: preserved goal, decisions, files, tests, next steps"}, provider.Usage{}, nil
	}
	if err := ctx.Err(); err != nil {
		return provider.Message{}, provider.Usage{}, err
	}
	reply := provider.Message{Role: "assistant", Content: "ok"}
	if len(c.replies) > 0 {
		reply = c.replies[0]
		c.replies = c.replies[1:]
	}
	var usage provider.Usage
	if len(c.usages) > 0 {
		usage = c.usages[0]
		c.usages = c.usages[1:]
	}
	return reply, usage, nil
}

func Test217KIn1MWindowWithLargeImageDoesNotTrigger60PercentThreshold(t *testing.T) {
	inTok, outTok := int64(215000), int64(2000)
	c := &planClient{usages: []provider.Usage{{Input: &inTok, Output: &outTok}, {Input: &inTok, Output: &outTok}}}
	a := New(c)
	a.AutoCompact = func() (bool, int, int) { return true, 60, 1000000 }
	h1 := []provider.Message{{Role: "user", Content: "Start task"}}
	if _, err := a.Run(context.Background(), h1); err != nil {
		t.Fatal(err)
	}
	hugeBase64 := "data:image/png;base64," + strings.Repeat("A", 2000000)
	h2 := []provider.Message{
		{Role: "user", Content: "Start task"},
		{Role: "assistant", Content: "ok"},
		{Role: "user", Content: "Check [Image #1]", Images: []string{hugeBase64}},
	}
	if _, err := a.Run(context.Background(), h2); err != nil {
		t.Fatal(err)
	}
	if len(c.requests) != 2 {
		t.Fatalf("expected 0 compaction calls at 217K/1M, got %d total requests", len(c.requests))
	}
}

func TestSmallerUtilityWindowChunksAndPreservesRecentMessagesAndOmitsEffort(t *testing.T) {
	inTok, outTok := int64(610000), int64(5000)
	compactCalls := 0
	c := &planClient{
		usages: []provider.Usage{{Input: &inTok, Output: &outTok}},
		compactFn: func(req provider.Request, _ int) (provider.Message, error) {
			compactCalls++
			if req.Effort != "" || req.Model != "haiku-small" {
				t.Fatalf("unexpected compact request: model=%q effort=%q", req.Model, req.Effort)
			}
			return provider.Message{Role: "assistant", Content: fmt.Sprintf("chunk-summary-%d", compactCalls)}, nil
		},
	}
	a := New(c)
	a.UtilityModel = func() string { return "haiku-small" }
	a.UtilityWindow = func() int { return 600 }
	a.AutoCompact = func() (bool, int, int) { return true, 60, 1000000 }
	h1 := []provider.Message{
		{Role: "user", Content: strings.Repeat("old-one ", 200)},
		{Role: "assistant", Content: strings.Repeat("ans-one ", 200)},
		{Role: "user", Content: strings.Repeat("old-two ", 200)},
	}
	if _, err := a.Run(context.Background(), h1); err != nil {
		t.Fatal(err)
	}
	h2 := append(append([]provider.Message{}, h1...),
		provider.Message{Role: "assistant", Content: "ok"},
		provider.Message{Role: "user", Content: "keep this recent message intact"},
	)
	if _, err := a.Run(context.Background(), h2); err != nil {
		t.Fatal(err)
	}
	if compactCalls < 2 {
		t.Fatalf("expected chunked compaction for small utility window, got %d calls", compactCalls)
	}
	lastReq := c.requests[len(c.requests)-1]
	if lastReq.Messages[len(lastReq.Messages)-1].Content != "keep this recent message intact" || !strings.Contains(lastReq.Messages[1].Content, "chunk-summary-") {
		t.Fatalf("recent message or summary missing: %+v", lastReq.Messages)
	}
}

func TestAutoCompactionFailureDoesNotRetryEveryStepAndTimeoutApplies(t *testing.T) {
	call := provider.ToolCall{ID: "t1", Type: "function"}
	call.Function.Name, call.Function.Arguments = "js_repl", `{"code":"1+1"}`
	compactAttempts := 0
	c := &planClient{
		replies: []provider.Message{{ToolCalls: []provider.ToolCall{call}}, {Content: "completed despite 400"}},
		compactFn: func(req provider.Request, _ int) (provider.Message, error) {
			compactAttempts++
			return provider.Message{}, fmt.Errorf("Provider returned HTTP 400: unsupported reasoning_effort")
		},
	}
	a := New(c)
	a.CompactTimeout = 50 * time.Millisecond
	a.AutoCompact = func() (bool, int, int) { return true, 60, 500 }
	out, err := a.Run(context.Background(), []provider.Message{{Role: "user", Content: strings.Repeat("z", 2000)}})
	if err != nil || out != "completed despite 400" || compactAttempts != 1 {
		t.Fatalf("out=%q err=%v attempts=%d", out, err, compactAttempts)
	}
	if !strings.Contains(a.LastCompaction.Error, "HTTP 400") {
		t.Fatalf("missing diagnostic error: %+v", a.LastCompaction)
	}
}
