package jsruntime

import (
	"fmt"
	"os"

	"zelin/internal/files"
)

var launch, _ = os.Getwd()

func (r *Runtime) Cwd() string { return r.cwd }

func (r *Runtime) usable() error {
	info, err := os.Stat(r.cwd)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("Working directory %s is unavailable; call cd() with an existing directory", r.cwd)
	}
	return nil
}

func (r *Runtime) abs(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("Path is required")
	}
	if err := r.usable(); err != nil {
		return "", err
	}
	return files.ResolveIn(r.cwd, path)
}

func (r *Runtime) chdir(path string) any {
	if path == "" {
		return map[string]any{"ok": true, "path": r.cwd}
	}
	target, err := files.ResolveIn(r.cwd, path)
	if err != nil {
		return failure(err)
	}
	info, err := os.Stat(target)
	if err != nil {
		return failure(err)
	}
	if !info.IsDir() {
		return failure(fmt.Errorf("%s is not a directory", target))
	}
	r.cwd = target
	return map[string]any{"ok": true, "path": target}
}
