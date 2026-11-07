package sessions

import (
	"crypto/rand"
	"encoding/hex"
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
	return session, s.SetSetting("active", id)
}

func (s *Store) Current() (Session, error) {
	id, err := s.Setting("active")
	if err != nil {
		return Session{}, err
	}
	if id == "" {
		return s.Create()
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
