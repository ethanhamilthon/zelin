package agent

import (
	"sync"
	"zash/internal/provider"
)

type steering struct {
	mu      sync.Mutex
	pending *provider.Message
}

func (a *Agent) Steer(message provider.Message) {
	a.steering.mu.Lock()
	defer a.steering.mu.Unlock()
	if message.Content == "" {
		a.steering.pending = nil
	} else {
		a.steering.pending = &message
	}
}

func (a *Agent) takeSteer(emit func(Event)) *provider.Message {
	a.steering.mu.Lock()
	message := a.steering.pending
	a.steering.pending = nil
	a.steering.mu.Unlock()
	if message != nil {
		emit(Event{Kind: "steer", Message: message})
	}
	return message
}
