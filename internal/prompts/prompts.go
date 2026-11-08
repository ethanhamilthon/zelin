package prompts

import (
	_ "embed"
	"strings"
)

//go:embed default_system.md
var DefaultSystem string

//go:embed js_docs.md
var rawJSDocs string

//go:embed js_tool.md
var rawJSTool string

//go:embed title.md
var rawTitle string

//go:embed compact.md
var rawCompact string

//go:embed handoff.md
var rawHandoff string

//go:embed restart.md
var rawRestart string

var (
	JSDocs  = strings.TrimSpace(rawJSDocs)
	JSTool  = strings.TrimSpace(rawJSTool)
	Title   = strings.TrimSpace(rawTitle)
	Compact = strings.TrimSpace(rawCompact)
	Handoff = strings.TrimSpace(rawHandoff)
	Restart = strings.TrimSpace(rawRestart)
)
