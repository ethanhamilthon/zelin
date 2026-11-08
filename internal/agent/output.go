package agent

import (
	"fmt"
	"unicode/utf8"
)

func LimitOutput(text string, limit int) string {
	if limit <= 0 || len(text) <= limit {
		return text
	}
	end := limit
	for end > 0 && !utf8.RuneStart(text[end]) {
		end--
	}
	return text[:end] + fmt.Sprintf("\n[output truncated: showing %d of %d bytes; change JS tool output limit in Settings]", end, len(text))
}
