package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"zelin/internal/sessions"
	"zelin/internal/tagged"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type SessionService struct {
	store  *sessions.Store
	notify func(name string, data any)
}

func (s *SessionService) Current() (sessions.Session, error)       { return s.store.Current() }
func (s *SessionService) Create() (sessions.Session, error)        { return s.store.Create() }
func (s *SessionService) List() ([]sessions.Session, error)        { return s.store.List() }
func (s *SessionService) Open(id string) (sessions.Session, error) { return s.store.Activate(id) }
func (s *SessionService) Save(id string, messages []sessions.Message, draft string, attachments []string) error {
	return s.store.Save(id, messages, draft, attachments)
}
func (s *SessionService) SaveRefs(id string, refs []tagged.Ref) error {
	return s.store.SaveRefs(id, refs)
}
func (s *SessionService) SaveDraft(id, draft string) error { return s.store.SaveDraft(id, draft) }
func (s *SessionService) SaveAttachments(id string, attachments []string) error {
	return s.store.SaveAttachments(id, attachments)
}

func (s *SessionService) Rename(id, title string) (string, error) {
	if err := s.store.SetTitle(id, title); err != nil {
		return "", err
	}
	title = strings.TrimSpace(title)
	if s.notify != nil {
		s.notify("session:title", map[string]string{"id": id, "title": title})
	}
	return title, nil
}

var unsafeName = regexp.MustCompile(`[^\p{L}\p{N}._ -]+`)

func exportName(title string) string {
	name := strings.TrimSpace(strings.TrimLeft(unsafeName.ReplaceAllString(title, ""), "."))
	if name == "" {
		name = "session"
	}
	return name + ".md"
}

func (s *SessionService) Export(id string) (string, error) {
	session, err := s.store.Get(id)
	if err != nil {
		return "", err
	}
	path, err := application.Get().Dialog.SaveFile().
		SetFilename(exportName(session.Title)).
		AddFilter("Markdown", "*.md").
		CanCreateDirectories(true).
		PromptForSingleSelection()
	if err != nil || path == "" {
		return "", err
	}
	if filepath.Ext(path) == "" {
		path += ".md"
	}
	if err := os.WriteFile(path, []byte(session.Markdown()), 0o600); err != nil {
		return "", fmt.Errorf("Could not export session: %w", err)
	}
	return path, nil
}
