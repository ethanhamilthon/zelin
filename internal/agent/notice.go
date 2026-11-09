package agent

import "zelin/internal/provider"

func (a *Agent) inbox(emit func(Event)) (model, shown []provider.Message) {
	if message, expanded := a.takeSteer(emit); message != nil {
		model, shown = append(model, *expanded), append(shown, *message)
	}
	if notice := a.takeNotice(emit); notice != nil {
		model, shown = append(model, *notice), append(shown, *notice)
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
