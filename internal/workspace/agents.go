package workspace

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const AgentsFile = "AGENTS.md"
const maxAgents = 64 << 10

func AgentsPath(dir string) string { return filepath.Join(dir, AgentsFile) }

func Exists(dir string) bool {
	info, err := os.Stat(AgentsPath(dir))
	return err == nil && !info.IsDir()
}

func ReadAgents(dir string) (string, error) {
	path := AgentsPath(dir)
	info, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("Cannot read %s: %w", path, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s is a directory; remove or rename it", path)
	}
	if info.Size() > maxAgents {
		return "", fmt.Errorf("%s is larger than %d KB; shorten it", path, maxAgents>>10)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("Cannot read %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func Instructions(dir string) (string, error) {
	text, err := ReadAgents(dir)
	if err != nil || text == "" {
		return "", err
	}
	return "<workspace-instructions path=\"" + escape(AgentsPath(dir)) + "\">\n" + text + "\n</workspace-instructions>", nil
}

func escape(text string) string {
	return strings.NewReplacer("&", "&amp;", "\"", "&quot;", "<", "&lt;", ">", "&gt;").Replace(text)
}
