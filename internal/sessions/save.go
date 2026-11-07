package sessions

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"zash/internal/agent"
)

func (s *Store) Save(id string, messages []Message, draft string, attachments []string) error {
	if messages == nil {
		messages = []Message{}
	}
	if attachments == nil {
		attachments = []string{}
	}
	title := DefaultTitle
	for _, message := range messages {
		if message.Role == "user" {
			runes := []rune(message.Content)
			if len(runes) > 60 {
				runes = runes[:60]
			}
			title = string(runes)
			break
		}
	}
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	files, err := json.Marshal(attachments)
	if err != nil {
		return err
	}
	result, err := s.db.Exec("UPDATE sessions SET title=CASE WHEN title=? THEN ? ELSE title END,messages=?,draft=?,attachments=?,updated=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=?",
		DefaultTitle, title, string(data), draft, string(files), id)
	return changed(result, err)
}

func (s *Store) Rename(id, title string) error {
	result, err := s.db.Exec("UPDATE sessions SET title=? WHERE id=?", title, id)
	return changed(result, err)
}

func (s *Store) SaveState(id string, state agent.State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	result, err := s.db.Exec("UPDATE sessions SET state=? WHERE id=?", string(data), id)
	return changed(result, err)
}

func changed(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return fmt.Errorf("Session not found")
	}
	return err
}
