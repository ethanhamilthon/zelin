package agent

import (
	"strings"

	"zelin/internal/provider"
)

func (a *Agent) inbox(emit func(Event), reply provider.Message, id string) (model, shown []provider.Message) {
	var events []Event
	collect := func(event Event) { events = append(events, event) }
	if message, expanded := a.takeSteer(collect); message != nil {
		model, shown = append(model, *expanded), append(shown, *message)
	}
	if notice := a.takeNotice(collect); notice != nil {
		model, shown = append(model, *notice), append(shown, *notice)
	}
	if len(model) == 0 {
		return nil, nil
	}
	if note := strings.TrimSpace(reply.Content); note != "" {
		emit(Event{Kind: "execution", Execution: &Execution{ID: id + "-note", Kind: "message", Code: note, State: "complete"}})
	}
	for _, event := range events {
		emit(event)
	}
	return model, shown
}

func (a *Agent) takeNotice(emit func(Event)) *provider.Message {
	if a.Notices == nil {
		return nil
	}
	text := a.Notices()
	if text == "" {
		return nil
	}
	message := provider.Message{Role: "user", Content: text}
	emit(Event{Kind: "notice", Message: &message})
	return &message
}
