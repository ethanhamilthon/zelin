package app

import (
	"zelin/internal/instructions"
	"zelin/internal/sessions"
)

type SettingsService struct {
	store        *sessions.Store
	instructions *instructions.Store
	notify       func(string, any)
}

func (s *SettingsService) Load() (string, error) { return s.store.Setting("ui") }
func (s *SettingsService) Save(value string) error {
	if err := s.store.SetSetting("ui", value); err != nil {
		return err
	}
	if s.notify != nil {
		s.notify("settings:changed", value)
	}
	return nil
}
func (s *SettingsService) Read(name string) (string, error) { return s.instructions.Read(name) }
func (s *SettingsService) Write(name, text string) error    { return s.instructions.Write(name, text) }
