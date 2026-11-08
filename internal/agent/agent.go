package agent

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"zash/internal/jsruntime"
	"zash/internal/provider"
)

type Completer interface {
	Complete(context.Context, provider.Request) (provider.Message, provider.Usage, error)
}

type Agent struct {
	Prompt      func() string
	client      Completer
	runtime     *jsruntime.Runtime
	history     []provider.Message
	transcript  []provider.Message
	usage       provider.Totals
	steering    steering
	AutoCompact func() (bool, int, int)
}

func New(client Completer) *Agent { return &Agent{client: client, runtime: jsruntime.New()} }

func (a *Agent) system() provider.Message {
	prompt := DefaultPrompt
	if a.Prompt != nil {
		prompt = a.Prompt()
	}
	return provider.Message{Role: "system", Content: prompt}
}

func (a *Agent) Run(ctx context.Context, history []provider.Message, observers ...func(Event)) (string, error) {
	emit := func(event Event) {
		for _, observer := range observers {
			observer(event)
		}
	}
	messages := []provider.Message{a.system()}
	if len(history) == 0 || history[len(history)-1].Role != "user" {
		return "", fmt.Errorf("A user message is required")
	}
	for _, message := range history {
		if (message.Role != "user" && message.Role != "assistant") || strings.TrimSpace(message.Content) == "" {
			return "", fmt.Errorf("Invalid conversation history")
		}
		messages = append(messages, provider.Message{Role: message.Role, Content: message.Content, Images: message.Images})
	}
	matches := len(history) == len(a.history)+1 && slices.EqualFunc(history[:len(history)-1], a.history, func(x, y provider.Message) bool {
		return x.Role == y.Role && x.Content == y.Content
	})
	if matches && len(a.transcript) > 0 {
		messages = append(append([]provider.Message{a.system()}, a.transcript[1:]...), messages[len(messages)-1])
	} else if len(a.history) > 0 {
		a.runtime = jsruntime.New()
	}
	history = append([]provider.Message{}, history...)
	defer func() { a.transcript = messages }()
	for step := 0; ; step++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		if message := a.takeSteer(emit); message != nil {
			messages = append(messages, *message)
			history = append(history, *message)
		}
		compacted, compactErr := a.autoCompact(ctx, messages, emit)
		if compactErr != nil {
			return "", compactErr
		}
		messages = compacted
		emit(Event{Kind: "activity", Activity: "Thinking"})
		reply, usage, err := a.client.Complete(ctx, provider.Request{Messages: messages, Tools: []provider.Tool{Tool}})
		if err != nil {
			return "", err
		}
		a.usage = a.usage.Add(usage)
		totals := a.usage
		emit(Event{Kind: "usage", Usage: &totals})
		if thinking := strings.TrimSpace(reply.Thinking()); thinking != "" {
			id := fmt.Sprintf("thinking-%d-%d", time.Now().UnixNano(), step)
			emit(Event{Kind: "execution", Execution: &Execution{ID: id, Kind: "thinking", Code: thinking, State: "complete"}})
		}
		reply.Role = "assistant"
		if len(reply.ToolCalls) == 0 {
			if strings.TrimSpace(reply.Content) == "" {
				return "", fmt.Errorf("Agent returned an empty response")
			}
			if message := a.takeSteer(emit); message != nil {
				messages = append(messages, reply, *message)
				history = append(history, *message)
				continue
			}
			messages = append(messages, reply)
			a.history = append(append([]provider.Message{}, history...), provider.Message{Role: "assistant", Content: reply.Content})
			return reply.Content, nil
		}
		messages = append(messages, reply)
		for _, call := range reply.ToolCalls {
			messages = append(messages, provider.Message{Role: "tool", Content: a.call(ctx, call, emit), ToolCallID: call.ID})
		}
	}
}
