package app

import (
	"fmt"
	"path/filepath"

	"zelin/internal/home"
	"zelin/internal/provider"
	"zelin/internal/tagged"
)

func (s *AgentService) expander(id string, refs []tagged.Ref) func(provider.Message) (provider.Message, error) {
	return func(message provider.Message) (provider.Message, error) {
		return s.expand(id, message, refs)
	}
}

func (s *AgentService) expand(id string, message provider.Message, refs []tagged.Ref) (provider.Message, error) {
	if len(refs) == 0 {
		message.Content = tagged.Plain(message.Content)
		return message, nil
	}
	record, err := s.store.Get(id)
	if err != nil {
		return message, err
	}
	item, err := s.store.Workspace(record.WorkspaceID)
	if err != nil {
		return message, fmt.Errorf("Workspace is unavailable: %w", err)
	}
	result, err := tagged.Expand(item.Dir, s.assets(id), record.WorkspaceID, message.Content, message.Images, refs)
	if err != nil {
		return message, err
	}
	message.Content, message.Images = result.Content, result.Images
	return message, nil
}

func (s *AgentService) assets(id string) string {
	dir, err := home.Dir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "images", id)
}
