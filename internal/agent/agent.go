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
	Persist        func(State) error
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
	emit := func(event Event) {
		for _, observer := range observers {
			observer(event)
		}
	}
	messages, err := a.begin(history)
	if err != nil {
		a.Reject(history, err)
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
	defer func() { a.settle(runErr, messages, history) }()
	sync := func() { a.checkpoint(messages, history) }
	sync()
	for step := 0; ; step++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		if imgMessage := a.takeImages(); imgMessage != nil {
			messages = append(messages, *imgMessage)
			sync()
		}
		if message, model := a.takeSteer(emit); message != nil {
			messages = append(messages, *model)
			history = append(history, *message)
			sync()
		}
		compacted, compactErr := a.autoCompact(ctx, messages, emit)
		if compactErr != nil {
			return "", compactErr
		}
		messages = compacted
		sync()
		emit(Event{Kind: "activity", Activity: "Thinking"})
		thinkingID := fmt.Sprintf("thinking-%d-%d", time.Now().UnixNano(), step)
		reply, usage, err := a.stream(ctx, messages, emit, thinkingID)
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
		emitSteps(emit, reply, thinkingID)
		reply.Role = "assistant"
		if len(reply.ToolCalls) == 0 {
			if strings.TrimSpace(reply.Content) == "" {
				return "", fmt.Errorf("Agent returned an empty response")
			}
			if message, model := a.takeSteer(emit); message != nil {
				messages = append(messages, reply, *model)
				history = append(history, *message)
				sync()
				continue
			}
			messages = append(messages, reply)
			history = append(history, provider.Message{Role: "assistant", Content: reply.Content})
			a.complete(messages, history)
			return reply.Content, nil
		}
		messages = append(messages, reply)
		for _, call := range reply.ToolCalls {
			messages = append(messages, provider.Message{Role: "tool", Content: a.call(ctx, call, emit), ToolCallID: call.ID})
		}
		sync()
	}
}
