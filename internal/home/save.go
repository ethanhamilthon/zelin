package home

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func SaveEnv(values map[string]string) error {
	path, err := Path(".env")
	if err != nil {
		return err
	}
	lines := []string{}
	if data, err := os.ReadFile(path); err == nil {
		for _, line := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
			key, _, _ := strings.Cut(strings.TrimPrefix(strings.TrimSpace(line), "export "), "=")
			if _, replaced := values[strings.TrimSpace(key)]; !replaced && line != "" {
				lines = append(lines, line)
			}
		}
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("%s=\"%s\"", key, values[key]))
	}
	temp := filepath.Join(filepath.Dir(path), ".env.tmp")
	if err := os.WriteFile(temp, []byte(strings.Join(lines, "\n")+"\n"), 0600); err != nil {
		return err
	}
	return os.Rename(temp, path)
}
