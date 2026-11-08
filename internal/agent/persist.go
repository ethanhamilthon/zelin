package agent

import (
	"context"
	"errors"
	"log"

	"zelin/internal/provider"
	"zelin/internal/tagged"
)

func (a *Agent) persist() {
	if a.Persist == nil {
		return
	}
	if err := a.Persist(a.State()); err != nil {
		log.Printf("persist state: %v", err)
	}
}

func (a *Agent) checkpoint(messages, history []provider.Message) {
	stopped := provider.Message{Role: "assistant", Content: StoppedMessage}
	a.transcript = messages
	a.history = append(history[:len(history):len(history)], stopped)
	a.persist()
}

func (a *Agent) complete(messages, history []provider.Message) {
	a.transcript, a.history = messages, history
	a.persist()
}

func (a *Agent) settle(err error, messages, history []provider.Message) {
	a.pinned = nil
	if err == nil {
		return
	}
	note := provider.Message{Role: "system", Content: failureNote(err)}
	a.checkpoint(append(append([]provider.Message{}, messages...), note), history)
	a.transcript = append(a.transcript, a.history[len(a.history)-1])
	a.persist()
}

func (a *Agent) Reject(history []provider.Message, err error) {
	last := len(history) - 1
	if last < 0 || history[last].Role != "user" {
		return
	}
	stopped := provider.Message{Role: "assistant", Content: StoppedMessage}
	a.history = append(append([]provider.Message{}, history...), stopped)
	if len(a.transcript) > 0 {
		user := provider.Message{Role: "user", Content: tagged.Plain(history[last].Content)}
		note := provider.Message{Role: "system", Content: failureNote(err)}
		a.transcript = append(append([]provider.Message{}, a.transcript...), user, note, stopped)
	}
	a.persist()
}

func failureNote(err error) string {
	if errors.Is(err, context.Canceled) {
		return "The user stopped the agent before this turn was completed."
	}
	return "The turn failed before completion: " + err.Error()
}
