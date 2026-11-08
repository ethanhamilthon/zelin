package agent

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"zelin/internal/provider"
)

type streamer struct {
	emit     func(Event)
	thinking string
	last     time.Time
	pending  *provider.Delta
}

func (a *Agent) stream(ctx context.Context, messages []provider.Message, emit func(Event), thinkingID string) (provider.Message, provider.Usage, error) {
	s := newStreamer(emit, thinkingID)
	reply, usage, err := a.client.Complete(ctx, provider.Request{Messages: messages, Tools: []provider.Tool{Tool}, OnDelta: s.delta})
	s.flush()
	if len(reply.ToolCalls) > 0 || err != nil {
		emit(Event{Kind: "stream", Message: &provider.Message{Role: "assistant"}})
	}
	return reply, usage, err
}

func newStreamer(emit func(Event), thinkingID string) *streamer {
	emit(Event{Kind: "stream", Message: &provider.Message{Role: "assistant"}})
	return &streamer{emit: emit, thinking: thinkingID}
}

func (s *streamer) delta(d provider.Delta) {
	s.pending = &d
	if time.Since(s.last) >= 50*time.Millisecond {
		s.flush()
	}
}

func (s *streamer) flush() {
	if s.pending == nil {
		return
	}
	d := *s.pending
	s.pending, s.last = nil, time.Now()
	if thinking := strings.TrimSpace(d.Thinking); thinking != "" {
		s.emit(Event{Kind: "execution", Execution: &Execution{ID: s.thinking, Kind: "thinking", Code: thinking, State: "running"}})
	}
	for _, call := range d.Calls {
		if call.ID != "" {
			s.emit(Event{Kind: "execution", Execution: &Execution{ID: call.ID, Code: partialCode(call.Function.Arguments), State: "running"}})
		}
	}
	if d.Content != "" {
		s.emit(Event{Kind: "stream", Message: &provider.Message{Role: "assistant", Content: d.Content}})
	}
}

func partialCode(args string) string {
	var full struct {
		Code string `json:"code"`
	}
	if json.Unmarshal([]byte(args), &full) == nil {
		return full.Code
	}
	start := strings.Index(args, `"code"`)
	if start < 0 {
		return ""
	}
	rest := args[start+6:]
	quote := strings.Index(rest, `"`)
	if quote < 0 {
		return ""
	}
	raw := strings.TrimSuffix(rest[quote+1:], `\`)
	if end := strings.LastIndex(raw, `"`); end >= 0 && (end == 0 || raw[end-1] != '\\') && strings.TrimSpace(raw[end+1:]) == "}" {
		raw = raw[:end]
	}
	for cut := 0; cut <= 6 && cut <= len(raw); cut++ {
		var text string
		if json.Unmarshal([]byte(`"`+raw[:len(raw)-cut]+`"`), &text) == nil {
			return text
		}
	}
	return ""
}

func emitSteps(emit func(Event), reply provider.Message, id string) {
	if thinking := strings.TrimSpace(reply.Thinking()); thinking != "" {
		emit(Event{Kind: "execution", Execution: &Execution{ID: id, Kind: "thinking", Code: thinking, State: "complete"}})
	}
	if note := strings.TrimSpace(reply.Content); note != "" && len(reply.ToolCalls) > 0 {
		emit(Event{Kind: "execution", Execution: &Execution{ID: id + "-note", Kind: "message", Code: note, State: "complete"}})
	}
}
