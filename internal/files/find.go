package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ignoredDirs = map[string]bool{"node_modules": true, "dist": true, "vendor": true, "target": true}

func Locate(path string) (string, error) {
	resolved, err := Resolve(path)
	if err != nil {
		return "", err
	}
	if info, statErr := os.Stat(resolved); statErr == nil && !info.IsDir() {
		return resolved, nil
	}
	if filepath.IsAbs(path) || path == "~" || strings.HasPrefix(path, "~/") {
		return "", fmt.Errorf("File not found: %s", resolved)
	}
	cwd, err := os.Getwd()
	home, _ := os.UserHomeDir()
	if err != nil || cwd == "/" || cwd == home {
		return "", fmt.Errorf("File not found: %s", resolved)
	}
	want := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "./")
	var matches []string
	visited := 0
	_ = filepath.WalkDir(cwd, func(current string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		visited++
		if visited > 15000 {
			return fs.SkipAll
		}
		if entry.IsDir() {
			if current != cwd && (strings.HasPrefix(entry.Name(), ".") || ignoredDirs[entry.Name()]) {
				return filepath.SkipDir
			}
			return nil
		}
		rel, relErr := filepath.Rel(cwd, current)
		if relErr != nil {
			return nil
		}
		slashed := filepath.ToSlash(rel)
		if slashed == want || strings.HasSuffix(slashed, "/"+want) {
			matches = append(matches, current)
			if len(matches) > 1 {
				return fs.SkipAll
			}
		}
		return nil
	})
	if len(matches) == 1 {
		return matches[0], nil
	}
	return "", fmt.Errorf("File not found: %s", resolved)
}
