package main

import (
	"fmt"
	"path/filepath"

	"zelin/internal/workspace"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (w *WorkspaceService) root(id string) (string, error) {
	item, err := w.store.Workspace(id)
	if err != nil {
		return "", err
	}
	return workspace.Resolve(item.Dir)
}

func (w *WorkspaceService) Browse(id, dir string) (workspace.Listing, error) {
	root, err := w.root(id)
	if err != nil {
		return workspace.Listing{}, err
	}
	if dir == "" {
		dir = root
	}
	return workspace.List(root, dir)
}

func (w *WorkspaceService) Suggest(id, dir, query string) ([]workspace.Entry, error) {
	root, err := w.root(id)
	if err != nil {
		return nil, err
	}
	return workspace.Suggest(root, dir, query)
}

func (w *WorkspaceService) Open(id, path string) error {
	root, err := w.root(id)
	if err != nil {
		return err
	}
	abs, err := workspace.Within(root, path)
	if err != nil {
		return err
	}
	real, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return fmt.Errorf("Cannot open %s: %w", abs, err)
	}
	return application.Get().Browser.OpenFile(real)
}
