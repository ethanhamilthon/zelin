package files

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

const grepLimit = 500

var skipped = map[string]bool{".git": true, "node_modules": true}

func Grep(pattern, root, glob string) ([]string, error) {
	expression, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("Invalid pattern: %w", err)
	}
	if root == "" {
		root = "."
	}
	if root, err = Resolve(root); err != nil {
		return nil, err
	}
	results := []string{}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if entry.IsDir() {
			if path != root && skipped[entry.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if glob != "" {
			if ok, _ := filepath.Match(glob, entry.Name()); !ok {
				return nil
			}
		}
		results = search(expression, path, results)
		if len(results) >= grepLimit {
			return fs.SkipAll
		}
		return nil
	})
	return results, err
}

func search(expression *regexp.Regexp, path string, results []string) []string {
	info, err := os.Stat(path)
	if err != nil || info.Size() > 2<<20 {
		return results
	}
	data, err := os.ReadFile(path)
	if err != nil || bytes.IndexByte(data[:min(len(data), 8000)], 0) >= 0 {
		return results
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(nil, 2<<20)
	for line := 1; scanner.Scan() && len(results) < grepLimit; line++ {
		if text := scanner.Text(); expression.MatchString(text) {
			if len(text) > 300 {
				text = text[:300] + "…"
			}
			results = append(results, fmt.Sprintf("%s:%d:%s", path, line, text))
		}
	}
	return results
}
