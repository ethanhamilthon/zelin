package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func (s *Store) Create() (Session, error) { return s.CreateIn(DefaultWorkspace) }

func (s *Store) CreateIn(workspaceID string) (Session, error) {
	if workspaceID == "" {
		workspaceID = DefaultWorkspace
	}
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return Session{}, err
	}
	session := Session{ID: hex.EncodeToString(bytes), Title: DefaultTitle, Messages: []Message{}, Updated: time.Now().UTC().Format(time.RFC3339Nano), WorkspaceID: workspaceID}
	_, err := s.db.Exec("INSERT INTO sessions(id,title,updated,workspace_id) VALUES(?,?,?,?)", session.ID, session.Title, session.Updated, workspaceID)
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
	rows, err := s.db.Query("SELECT id,title,updated,workspace_id FROM sessions ORDER BY updated DESC,id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	list := []Session{}
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.ID, &session.Title, &session.Updated, &session.WorkspaceID); err != nil {
			return nil, err
		}
		list = append(list, session)
	}
	return list, rows.Err()
}
