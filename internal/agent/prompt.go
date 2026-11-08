package agent

import (
	"strings"

	"zelin/internal/prompts"
)

var DefaultPrompt = prompts.DefaultSystem

func SystemPrompt(base string) string {
	trimmed := strings.TrimSpace(base)
	if trimmed == "" {
		return prompts.JSDocs
	}
	return trimmed + "\n\n" + prompts.JSDocs
}
