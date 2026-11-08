package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zelin/internal/jsruntime"
	"zelin/internal/provider"
)

type Completer interface {
	Complete(context.Context, provider.Request) (provider.Message, provider.Usage, error)
}

type Agent struct {
	OutputLimit    func() int
	Prompt         func() string
	Extra          func() (string, error)
	Expand         func(provider.Message) (provider.Message, error)
	UtilityModel   func() string
	UtilityWindow  func() int
	AutoCompact    func() (bool, int, int)
	CompactTimeout time.Duration
	LastCompaction CompactionDiag
	client         Completer
	runtime        *jsruntime.Runtime
	pinned         *provider.Message
	history        []provider.Message
	transcript     []provider.Message
	usage          provider.Totals
	reportedLen    int
	compactFailed  bool
	steering       steering
}

func New(client Completer) *Agent { return NewIn(client, "") }

func NewIn(client Completer, dir string) *Agent {
	return &Agent{client: client, runtime: jsruntime.NewIn(dir)}
}

func (a *Agent) system() provider.Message {
	if a.pinned != nil {
		return *a.pinned
	}
	return a.compose("")
}

func (a *Agent) compose(extra string) provider.Message {
	prompt := DefaultPrompt
	if a.Prompt != nil {
		prompt = a.Prompt()
	}
	content := SystemPrompt(prompt)
	if extra != "" {
		content += "\n\n" + extra
	}
	return provider.Message{Role: "system", Content: content}
}

func (a *Agent) Run(ctx context.Context, history []provider.Message, observers ...func(Event)) (_ string, runErr error) {
	progress := false
	emit := func(event Event) {
		progress = progress || event.Kind == "execution"
		for _, observer := range observers {
			observer(event)
		}
	}
	messages, err := a.begin(history)
	if err != nil {
		return "", err
	}
	if a.continues(history) && len(a.transcript) > 0 {
		messages = append(append([]provider.Message{a.system()}, a.transcript[1:]...), messages[len(messages)-1])
	} else {
		a.reportedLen = 0
		if len(a.history) > 0 {
			a.runtime = jsruntime.NewIn(a.runtime.Cwd())
		}
	}
	history = append([]provider.Message{}, history...)
	prior := a.transcript
	defer func() { a.settle(runErr, progress, prior, messages, history) }()
	for step := 0; ; step++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		if imgMessage := a.takeImages(); imgMessage != nil {
			messages = append(messages, *imgMessage)
		}
		if message, model := a.takeSteer(emit); message != nil {
			messages = append(messages, *model)
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
		prevCtx := a.usage.Context
		a.usage = a.usage.Add(usage)
		if usage.Input != nil {
			a.reportedLen = len(messages) + 1
		} else if prevCtx != nil {
			nextCtx := *prevCtx + EstimateMessage(reply)
			a.usage.Context = &nextCtx
			a.reportedLen = len(messages) + 1
		}
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
			if message, model := a.takeSteer(emit); message != nil {
				messages = append(messages, reply, *model)
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
