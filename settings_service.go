package main

import (
	"zash/internal/instructions"
	"zash/internal/sessions"
)

type SettingsService struct {
	store        *sessions.Store
	instructions *instructions.Store
}

func (s *SettingsService) Load() (string, error)            { return s.store.Setting("ui") }
func (s *SettingsService) Save(value string) error          { return s.store.SetSetting("ui", value) }
func (s *SettingsService) Read(name string) (string, error) { return s.instructions.Read(name) }
func (s *SettingsService) Write(name, text string) error    { return s.instructions.Write(name, text) }
