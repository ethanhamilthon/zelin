package agent

import (
	"sync"
	"zelin/internal/provider"
)

type steering struct {
	mu      sync.Mutex
	pending *provider.Message
	model   *provider.Message
}

func (a *Agent) Steer(message provider.Message) { a.SteerWith(message, message) }

func (a *Agent) SteerWith(message, model provider.Message) {
	a.steering.mu.Lock()
	defer a.steering.mu.Unlock()
	if message.Content == "" {
		a.steering.pending, a.steering.model = nil, nil
		return
	}
	a.steering.pending, a.steering.model = &message, &model
}

func (a *Agent) takeSteer(emit func(Event)) (*provider.Message, *provider.Message) {
	a.steering.mu.Lock()
	message, model := a.steering.pending, a.steering.model
	a.steering.pending, a.steering.model = nil, nil
	a.steering.mu.Unlock()
	if message != nil {
		emit(Event{Kind: "steer", Message: message})
	}
	return message, model
}
