package main

import "zash/internal/sessions"

type SessionService struct{ store *sessions.Store }

func (s *SessionService) Current() (sessions.Session, error)       { return s.store.Current() }
func (s *SessionService) Create() (sessions.Session, error)        { return s.store.Create() }
func (s *SessionService) List() ([]sessions.Session, error)        { return s.store.List() }
func (s *SessionService) Open(id string) (sessions.Session, error) { return s.store.Activate(id) }
func (s *SessionService) Save(id string, messages []sessions.Message, draft string) error {
	return s.store.Save(id, messages, draft)
}
