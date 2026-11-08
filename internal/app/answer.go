package app

import (
	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/tagged"
)

func (s *AgentService) saveAnswer(id string, content string, result AgentResult, current *agent.Agent) error {
	var saveErr error
	if content != "" || len(result.Executions) > 0 || len(result.Files) > 0 || result.Error != "" {
		if content == "" {
			content = agent.StoppedMessage
		}
		saveErr = s.store.Finish(id, sessions.Message{Role: "assistant", Content: content, Executions: result.Executions, Turn: result.Turn, Files: result.Files}, current.State())
	} else {
		saveErr = s.store.SaveState(id, current.State())
	}
	return saveErr
}

func plainHistory(history []provider.Message) []provider.Message {
	plain := append([]provider.Message{}, history...)
	for i := range plain {
		plain[i].Content = tagged.Plain(plain[i].Content)
	}
	return plain
}
