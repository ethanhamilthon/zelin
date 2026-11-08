package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const readLimit = 2000

func Resolve(path string) (string, error) { return ResolveIn("", path) }

func ResolveIn(base, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("Path is required")
	}
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[1:])
	}
	if base != "" && !filepath.IsAbs(path) {
		path = filepath.Join(base, path)
	}
	return filepath.Abs(path)
}

func Read(path string, offset, limit int) (string, error) {
	path, err := Resolve(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if offset < 1 {
		offset = 1
	}
	if limit < 1 {
		limit = readLimit
	}
	if offset > len(lines) {
		return "", fmt.Errorf("Offset %d is past the end of the file (%d lines)", offset, len(lines))
	}
	end := min(offset-1+limit, len(lines))
	var out strings.Builder
	for i := offset - 1; i < end; i++ {
		fmt.Fprintf(&out, "%d\t%s\n", i+1, lines[i])
	}
	if offset > 1 || end < len(lines) {
		fmt.Fprintf(&out, "[lines %d-%d of %d]\n", offset, end, len(lines))
	}
	return out.String(), nil
}

func Write(path, text string) error {
	path, err := Resolve(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return replace(path, text)
}

func Edit(path, old, text string) error {
	path, err := Resolve(path)
	if err != nil {
		return err
	}
	if old == "" {
		return fmt.Errorf("Text to replace is empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	switch count := strings.Count(string(data), old); count {
	case 0:
		return fmt.Errorf("Text to replace was not found in %s", path)
	case 1:
		return replace(path, strings.Replace(string(data), old, text, 1))
	default:
		return fmt.Errorf("Text to replace matches %d times in %s; include more context", count, path)
	}
}

func replace(path, text string) error {
	mode := os.FileMode(0644)
	if info, err := os.Stat(path); err == nil {
		mode = info.Mode().Perm()
	}
	temp := path + ".zash-tmp"
	if err := os.WriteFile(temp, []byte(text), mode); err != nil {
		return err
	}
	return os.Rename(temp, path)
}
