package tagged

import "regexp"

type Ref struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Path      string `json:"path,omitempty"`
	Label     string `json:"label"`
	Data      string `json:"data,omitempty"`
	Workspace string `json:"workspaceId,omitempty"`
	Session   string `json:"sessionId,omitempty"`
}

var Token = regexp.MustCompile(`\[@([A-Za-z0-9_-]{4,40})(?:\|([^\]\n]*))?\]`)

func Used(content string, refs []Ref) []Ref {
	seen, out := map[string]bool{}, []Ref{}
	for _, match := range Token.FindAllStringSubmatch(content, -1) {
		for _, ref := range refs {
			if ref.ID == match[1] && !seen[ref.ID] {
				seen[ref.ID] = true
				out = append(out, ref)
			}
		}
	}
	return out
}

func Plain(content string) string {
	return Token.ReplaceAllStringFunc(content, func(token string) string {
		match := Token.FindStringSubmatch(token)
		if match[2] == "" {
			return token
		}
		return "@" + match[2]
	})
}
