package workspace

import (
	"fmt"
	"os"
	"strings"

	"zelin/internal/files"
)

func Resolve(dir string) (string, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return "", fmt.Errorf("Directory is required")
	}
	path, err := files.Resolve(dir)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("Directory is not accessible: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}
	handle, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("Directory is not accessible: %w", err)
	}
	return path, handle.Close()
}
