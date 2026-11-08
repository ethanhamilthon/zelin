package agent

import "zelin/internal/provider"

const ImageTokens int64 = 1000

func EstimateMessage(m provider.Message) int64 {
	chars := len(m.Content) + len(m.Reasoning) + len(m.ReasoningContent)
	for _, call := range m.ToolCalls {
		chars += len(call.Function.Name) + len(call.Function.Arguments)
	}
	tokens := int64((chars+3)/4) + 4
	for _, img := range m.Images {
		if img != "" {
			tokens += ImageTokens
		}
	}
	return tokens
}

func Estimate(messages []provider.Message) int64 {
	var total int64
	for _, m := range messages {
		total += EstimateMessage(m)
	}
	return total
}

func (a *Agent) ContextTokens(messages []provider.Message) int64 {
	if a.usage.Context != nil && a.reportedLen > 0 && a.reportedLen <= len(messages) {
		return *a.usage.Context + Estimate(messages[a.reportedLen:])
	}
	if a.usage.Context != nil && len(messages) == 0 {
		return *a.usage.Context
	}
	return Estimate(messages)
}
