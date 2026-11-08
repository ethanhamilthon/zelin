package sessions

import (
	"fmt"
	"strings"
)

func lockKey(id string) string { return "title-locked:" + id }

func (s *Store) SetTitle(id, title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return fmt.Errorf("Title cannot be empty")
	}
	if err := s.Rename(id, title); err != nil {
		return err
	}
	return s.SetSetting(lockKey(id), "1")
}

func (s *Store) AutoTitle(id, title string) (bool, error) {
	locked, err := s.Setting(lockKey(id))
	if err != nil || locked != "" {
		return false, err
	}
	return true, s.Rename(id, title)
}
