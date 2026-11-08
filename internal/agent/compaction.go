package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"zelin/internal/prompts"
	"zelin/internal/provider"
)

var CompactPrompt = prompts.Compact

const DefaultCompactionTimeout = 90 * time.Second

func (a *Agent) compactTimeout() time.Duration {
	if a.CompactTimeout > 0 {
		return a.CompactTimeout
	}
	return DefaultCompactionTimeout
}

func (a *Agent) compact(ctx context.Context, messages []provider.Message, reason string, threshold, mainWin, utilWin int, emit func(Event)) ([]provider.Message, error) {
	emit(Event{Kind: "activity", Activity: "Compacting context"})
	started := time.Now()
	model := ""
	if a.UtilityModel != nil {
		model = a.UtilityModel()
	}
	if utilWin <= 0 && a.UtilityWindow != nil {
		utilWin = a.UtilityWindow()
	}
	estimate := a.ContextTokens(messages)
	keepBudget := int64(8000)
	if mainWin > 0 {
		keepBudget = max(int64(200), int64(mainWin)*20/100)
	}
	older, recent := splitWorkingContext(messages[1:], keepBudget)
	chunks := chunkOlder(older, utilWin)
	diag := CompactionDiag{Reason: reason, UtilityModel: model, Threshold: threshold, MainWindow: mainWin, UtilityWindow: utilWin, ContextEstimate: estimate, Chunks: len(chunks)}
	summary := ""
	for _, chunk := range chunks {
		userContent := chunk
		if summary != "" {
			userContent = "Existing summary of earlier turns:\n" + summary + "\n\nNext transcript segment to fold into the summary:\n" + chunk
		}
		reqCtx, cancel := context.WithTimeout(ctx, a.compactTimeout())
		reply, usage, err := a.client.Complete(reqCtx, provider.Request{
			Model:    model,
			Messages: []provider.Message{{Role: "system", Content: CompactPrompt}, {Role: "user", Content: userContent}},
		})
		cancel()
		if err != nil {
			diag.DurationMs = time.Since(started).Milliseconds()
			diag.Error = err.Error()
			a.LastCompaction = diag
			logCompaction(diag)
			return nil, fmt.Errorf("Compaction failed: %w", err)
		}
		next := strings.TrimSpace(reply.Content)
		if next == "" {
			diag.DurationMs = time.Since(started).Milliseconds()
			diag.Error = "empty summary"
			a.LastCompaction = diag
			logCompaction(diag)
			return nil, fmt.Errorf("Compaction returned an empty summary")
		}
		prevCtx := a.usage.Context
		a.usage = a.usage.Add(usage)
		a.usage.Context = prevCtx
		summary = next
	}
	diag.DurationMs = time.Since(started).Milliseconds()
	a.LastCompaction = diag
	logCompaction(diag)
	result := append([]provider.Message{a.system(), {Role: "user", Content: "Conversation summary (continue from here; JavaScript bindings are unchanged):\n" + summary}}, recent...)
	size := Estimate(result)
	a.usage.Context = &size
	a.reportedLen = len(result)
	a.compactFailed = false
	totals := a.usage
	emit(Event{Kind: "usage", Usage: &totals})
	return result, nil
}

func (a *Agent) Compact(ctx context.Context, fallback []provider.Message, emit func(Event)) error {
	messages := a.transcript
	if len(messages) < 2 {
		messages = append([]provider.Message{a.system()}, fallback...)
	}
	if len(messages) < 2 {
		return fmt.Errorf("No conversation to compact")
	}
	threshold, mainWin := 0, 0
	if a.AutoCompact != nil {
		_, threshold, mainWin = a.AutoCompact()
	}
	compacted, err := a.compact(ctx, messages, "manual", threshold, mainWin, 0, emit)
	if err != nil {
		return err
	}
	if len(a.history) == 0 {
		a.history = append([]provider.Message{}, fallback...)
	}
	a.transcript = compacted
	return nil
}

func (a *Agent) autoCompact(ctx context.Context, messages []provider.Message, emit func(Event)) ([]provider.Message, error) {
	if a.compactFailed || a.AutoCompact == nil {
		return messages, nil
	}
	enabled, threshold, mainWin := a.AutoCompact()
	if !enabled || mainWin <= 0 {
		return messages, nil
	}
	if a.ContextTokens(messages)*100 < int64(mainWin*threshold) {
		return messages, nil
	}
	compacted, err := a.compact(ctx, messages, "auto", threshold, mainWin, 0, emit)
	if err != nil {
		a.compactFailed = true
		emit(Event{Kind: "compact_error", Activity: err.Error()})
		return messages, nil
	}
	return compacted, nil
}
