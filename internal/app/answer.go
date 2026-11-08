package app

import (
	"zelin/internal/agent"
	"zelin/internal/sessions"
)

func (s *AgentService) saveAnswer(id string, content string, result AgentResult, current *agent.Agent) error {
	var saveErr error
	if content != "" || len(result.Executions) > 0 || len(result.Files) > 0 {
		if content == "" {
			content = agent.StoppedMessage
		}
		saveErr = s.store.Finish(id, sessions.Message{Role: "assistant", Content: content, Executions: result.Executions, Turn: result.Turn, Files: result.Files}, current.State())
	} else {
		saveErr = s.store.SaveState(id, current.State())
	}
	return saveErr
}
