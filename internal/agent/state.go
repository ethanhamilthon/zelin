package agent

import (
	"zelin/internal/jsruntime"
	"zelin/internal/prompts"
	"zelin/internal/provider"
)

type State struct {
	History    []provider.Message `json:"history"`
	Transcript []provider.Message `json:"transcript"`
	Usage      provider.Totals    `json:"usage"`
	Cwd        string             `json:"cwd,omitempty"`
}

func (a *Agent) State() State {
	return State{History: a.history, Transcript: a.transcript, Usage: a.usage, Cwd: a.runtime.Cwd()}
}

func (a *Agent) Restore(state State) {
	a.history, a.transcript, a.usage = state.History, state.Transcript, state.Usage
	a.reportedLen = 0
	a.compactFailed = false
	if state.Cwd != "" {
		a.runtime = jsruntime.NewIn(state.Cwd)
	}
	if len(a.transcript) > 0 {
		if a.usage.Context != nil {
			a.reportedLen = len(a.transcript)
		}
		a.transcript = append(a.transcript, provider.Message{Role: "system", Content: prompts.Restart})
	}
}
