package agent

import "zash/internal/provider"

type State struct {
	History    []provider.Message `json:"history"`
	Transcript []provider.Message `json:"transcript"`
}

func (a *Agent) State() State { return State{History: a.history, Transcript: a.transcript} }

func (a *Agent) Restore(state State) {
	a.history, a.transcript = state.History, state.Transcript
	if len(a.transcript) > 0 {
		a.transcript = append(a.transcript, provider.Message{Role: "system", Content: "The application restarted. JavaScript bindings have been reset; previous tool results remain in the conversation. Recreate any variables you need."})
	}
}
