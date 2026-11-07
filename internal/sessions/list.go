package sessions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"
)

func (s *Store) Create() (Session, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return Session{}, err
	}
	session := Session{ID: hex.EncodeToString(bytes), Title: "New session", Messages: []Message{}, Updated: time.Now().UTC().Format(time.RFC3339Nano)}
	_, err := s.db.Exec("INSERT INTO sessions(id,title,updated) VALUES(?,?,?)", session.ID, session.Title, session.Updated)
	if err != nil {
		return Session{}, err
	}
	return s.Activate(session.ID)
}

func (s *Store) Activate(id string) (Session, error) {
	session, err := s.Get(id)
	if err != nil {
		return session, err
	}
	_, err = s.db.Exec("INSERT INTO settings(key,value) VALUES('active',?) ON CONFLICT(key) DO UPDATE SET value=excluded.value", id)
	return session, err
}

func (s *Store) Current() (Session, error) {
	var id string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key='active'").Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return s.Create()
	}
	if err != nil {
		return Session{}, err
	}
	return s.Get(id)
}

func (s *Store) List() ([]Session, error) {
	rows, err := s.db.Query("SELECT id,title,updated FROM sessions ORDER BY updated DESC,id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Session{}
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.ID, &session.Title, &session.Updated); err != nil {
			return nil, err
		}
		list = append(list, session)
	}
	return list, rows.Err()
}
