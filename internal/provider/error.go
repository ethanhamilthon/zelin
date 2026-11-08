package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func httpError(status int, body io.Reader, secret string) error {
	raw, _ := io.ReadAll(io.LimitReader(body, 4096))
	detail := extractErrorText(raw)
	if secret != "" {
		detail = strings.ReplaceAll(detail, secret, "[redacted]")
	}
	detail = strings.Join(strings.Fields(detail), " ")
	if len(detail) > 360 {
		detail = detail[:360] + "..."
	}
	if detail != "" {
		return fmt.Errorf("Provider returned HTTP %d: %s", status, detail)
	}
	return fmt.Errorf("Provider returned HTTP %d", status)
}

func extractErrorText(raw []byte) string {
	var envelope struct {
		Error   any    `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(raw, &envelope) == nil {
		switch val := envelope.Error.(type) {
		case string:
			if val != "" {
				return val
			}
		case map[string]any:
			if msg, ok := val["message"].(string); ok && msg != "" {
				return msg
			}
		}
		if envelope.Message != "" {
			return envelope.Message
		}
	}
	return strings.TrimSpace(string(raw))
}
