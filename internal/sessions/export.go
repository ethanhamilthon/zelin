package sessions

import (
	"fmt"
	"strings"

	"zelin/internal/tagged"
)

func fence(text string) string {
	ticks := "```"
	for strings.Contains(text, ticks) {
		ticks += "`"
	}
	return ticks
}

func block(out *strings.Builder, label, text string) {
	mark := fence(text)
	fmt.Fprintf(out, "%s%s\n%s\n%s\n\n", mark, label, text, mark)
}

func (s Session) Markdown() string {
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n\n", strings.TrimSpace(s.Title))
	for _, message := range s.Messages {
		heading := "User"
		if message.Role == "assistant" {
			heading = "Assistant"
		}
		fmt.Fprintf(&out, "## %s\n\n", heading)
		for _, step := range message.Executions {
			if step.Kind == "thinking" {
				block(&out, "text", step.Code)
				continue
			}
			block(&out, "js", step.Code)
			if step.Output != "" {
				block(&out, "text", step.Output)
			}
		}
		out.WriteString(strings.TrimSpace(tagged.Plain(message.Content)) + "\n\n")
		if n := len(message.Images); n > 0 {
			fmt.Fprintf(&out, "_%d image attachment(s) not included_\n\n", n)
		}
	}
	return strings.TrimRight(out.String(), "\n") + "\n"
}
