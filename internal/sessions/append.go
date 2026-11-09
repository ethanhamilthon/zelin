package sessions

import "encoding/json"

func (s *Store) Append(id string, message Message) error {
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
	next, err := json.Marshal(append(messages, message))
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE sessions SET messages=?,updated=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=?", string(next), id); err != nil {
		return err
	}
	return tx.Commit()
}
