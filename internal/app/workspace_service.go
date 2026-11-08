package app

import (
	"os"

	"zelin/internal/sessions"
	"zelin/internal/workspace"
)

func (s *AgentService) workspaceDir(id string) string {
	if item, err := s.store.Workspace(id); err == nil {
		return item.Dir
	}
	dir, _ := os.Getwd()
	return dir
}

func (s *AgentService) workspaceInstructions(session string) (string, error) {
	record, err := s.store.Get(session)
	if err != nil {
		return "", err
	}
	item, err := s.store.Workspace(record.WorkspaceID)
	if err != nil {
		return "", nil
	}
	return workspace.Instructions(item.Dir)
}

type WorkspaceService struct {
	store  *sessions.Store
	notify func(name string, data any)
}

type AgentsFile struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

func (w *WorkspaceService) changed() {
	if w.notify != nil {
		w.notify("workspace:changed", nil)
	}
}

func (w *WorkspaceService) List() ([]sessions.Workspace, error) { return w.store.Workspaces() }

func (w *WorkspaceService) Create(name, icon, dir string) (sessions.Workspace, error) {
	item, err := w.store.CreateWorkspace(name, icon, dir)
	if err == nil {
		w.changed()
	}
	return item, err
}

func (w *WorkspaceService) Update(id, name, icon, dir string) (sessions.Workspace, error) {
	item, err := w.store.UpdateWorkspace(id, name, icon, dir)
	if err == nil {
		w.changed()
	}
	return item, err
}

func (w *WorkspaceService) NewSession(id string) (sessions.Session, error) {
	if _, err := w.store.Workspace(id); err != nil {
		return sessions.Session{}, err
	}
	session, err := w.store.CreateIn(id)
	if err != nil {
		return session, err
	}
	return w.store.Activate(session.ID)
}

func (w *WorkspaceService) Agents(id string) (AgentsFile, error) {
	item, err := w.store.Workspace(id)
	if err != nil {
		return AgentsFile{}, err
	}
	return AgentsFile{Path: workspace.AgentsPath(item.Dir), Exists: workspace.Exists(item.Dir)}, nil
}
