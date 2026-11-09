package tasks

import (
	"os"
	"strings"
)

func (s *Session) Output(id, tail int) (string, error) {
	t, err := s.find(id)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(t.log)
	if err != nil {
		return "", err
	}
	if tail <= 0 {
		tail = 50
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) > tail {
		lines = lines[len(lines)-tail:]
	}
	text := strings.Join(lines, "\n")
	if len(text) > 64000 {
		text = text[len(text)-64000:]
	}
	return text, nil
}
