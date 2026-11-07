package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"zash/internal/jsruntime"
	"zash/internal/provider"
)

type Completer interface {
	Complete(context.Context, []provider.Message) (provider.Message, provider.Usage, error)
}

type Agent struct {
	client     Completer
	runtime    *jsruntime.Runtime
	history    []provider.Message
	transcript []provider.Message
	usage      provider.Totals
}

func New(client Completer) *Agent { return &Agent{client: client, runtime: jsruntime.New()} }

func (a *Agent) Run(ctx context.Context, history []provider.Message, observers ...func(Event)) (string, error) {
	emit := func(event Event) {
		for _, observer := range observers {
			observer(event)
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	messages := []provider.Message{{Role: "system", Content: "You are Zash, a helpful coding agent. Reply in Markdown. Your only tool is js_repl, synchronous JavaScript in Goja. Its only host function is bash(cmd, timeout): executes a shell command and returns {stdout, stderr, exitCode, timedOut}. Timeout is in seconds, default 30, maximum 120. No async/await, promises, console, require, or other host APIs. Use the last expression for tool output. Variables persist across turns; reuse variables or use var to avoid redeclaration errors. Commands run with the user's permissions in the app's working directory. Never claim success without checking command results. Do not run destructive commands unless the user explicitly requests them."}}
	if len(history) == 0 || history[len(history)-1].Role != "user" {
		return "", fmt.Errorf("A user message is required")
	}
	for _, message := range history {
		if (message.Role != "user" && message.Role != "assistant") || strings.TrimSpace(message.Content) == "" {
			return "", fmt.Errorf("Invalid conversation history")
		}
		messages = append(messages, provider.Message{Role: message.Role, Content: message.Content})
	}
	matches := len(history) == len(a.history)+1 && slices.EqualFunc(history[:len(history)-1], a.history, func(x, y provider.Message) bool {
		return x.Role == y.Role && x.Content == y.Content
	})
	if matches && len(a.transcript) > 0 {
		messages = append(append([]provider.Message{}, a.transcript...), messages[len(messages)-1])
	} else if len(a.history) > 0 {
		a.runtime = jsruntime.New()
	}
	for step := 0; step < 12; step++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		emit(Event{Kind: "activity", Activity: "Thinking"})
		reply, usage, err := a.client.Complete(ctx, messages)
		if err != nil {
			return "", err
		}
		a.usage = a.usage.Add(usage)
		totals := a.usage
		emit(Event{Kind: "usage", Usage: &totals})
		if len(reply.ToolCalls) == 0 {
			if strings.TrimSpace(reply.Content) == "" {
				return "", fmt.Errorf("Agent returned an empty response")
			}
			reply.Role = "assistant"
			a.transcript = append(messages, reply)
			a.history = append(append([]provider.Message{}, history...), provider.Message{Role: "assistant", Content: reply.Content})
			return reply.Content, nil
		}
		reply.Role = "assistant"
		messages = append(messages, reply)
		for _, call := range reply.ToolCalls {
			var args struct {
				Code string `json:"code"`
			}
			result := "Unknown tool; use js_repl only"
			if call.Function.Name == "js_repl" {
				if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
					result = "Invalid arguments; expected a JSON object with a code string"
				} else {
					execution := Execution{ID: call.ID, Code: args.Code, State: "running"}
					emit(Event{Kind: "activity", Activity: "Running JavaScript"})
					emit(Event{Kind: "execution", Execution: &execution})
					result = a.runtime.Run(ctx, args.Code)
					completed := Execution{ID: call.ID, Code: args.Code, Output: result, State: "complete"}
					emit(Event{Kind: "execution", Execution: &completed})
				}
			}
			messages = append(messages, provider.Message{Role: "tool", Content: result, ToolCallID: call.ID})
		}
	}
	return "", fmt.Errorf("Agent reached the tool step limit")
}
