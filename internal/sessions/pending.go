package sessions

import (
	"database/sql"
	"encoding/json"

	"zelin/internal/agent"
)

func (s *Store) AddStep(id string, step agent.Execution) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var data string
	if err := tx.QueryRow("SELECT pending FROM sessions WHERE id=?", id).Scan(&data); err != nil {
		return err
	}
	var steps []agent.Execution
	if err := json.Unmarshal([]byte(data), &steps); err != nil {
		return err
	}
	next, err := json.Marshal(append(steps, step))
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE sessions SET pending=? WHERE id=?", string(next), id); err != nil {
		return err
	}
	return tx.Commit()
}

func recoverPending(db *sql.DB) error {
	rows, err := db.Query("SELECT id FROM sessions WHERE pending != '[]'")
	if err != nil {
		return err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		ids = append(ids, id)
	}
	rows.Close()
	store := &Store{db: db}
	for _, id := range ids {
		if err := store.settlePending(id); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) settlePending(id string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var data, pending string
	if err := tx.QueryRow("SELECT messages,pending FROM sessions WHERE id=?", id).Scan(&data, &pending); err != nil {
		return err
	}
	var messages []Message
	var steps []agent.Execution
	if err := json.Unmarshal([]byte(data), &messages); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(pending), &steps); err != nil {
		return err
	}
	next, err := json.Marshal(append(messages, Message{Role: "assistant", Content: agent.StoppedMessage, Executions: steps}))
	if err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE sessions SET messages=?,pending='[]' WHERE id=?", string(next), id); err != nil {
		return err
	}
	return tx.Commit()
}
