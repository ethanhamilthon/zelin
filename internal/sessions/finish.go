package sessions

import (
	"encoding/json"
	"zelin/internal/agent"
)

func (s *Store) Finish(id string, message Message, state agent.State) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var data string
	if err := tx.QueryRow("SELECT messages FROM sessions WHERE id=?", id).Scan(&data); err != nil {
		return err
	}
	var messages []Message
	if err := json.Unmarshal([]byte(data), &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	history, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	snapshot, err := json.Marshal(state)
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE sessions SET messages=?,state=?,updated=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=?", string(history), string(snapshot), id); err != nil {
		return err
	}
	return tx.Commit()
}
