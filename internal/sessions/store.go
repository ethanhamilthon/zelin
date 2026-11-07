package sessions

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"zash/internal/agent"
	"zash/internal/provider"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	Role       string            `json:"role"`
	Content    string            `json:"content"`
	Images     []string          `json:"images,omitempty"`
	Executions []agent.Execution `json:"executions,omitempty"`
	Turn       *Turn             `json:"turn,omitempty"`
}

type Session struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Messages    []Message       `json:"messages"`
	Draft       string          `json:"draft"`
	Attachments []string        `json:"attachments"`
	Updated     string          `json:"updated"`
	Usage       provider.Totals `json:"usage"`
	State       agent.State     `json:"-"`
}

const DefaultTitle = "New session"

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	file.Close()
	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
 id TEXT PRIMARY KEY, title TEXT NOT NULL, messages TEXT NOT NULL DEFAULT '[]',
 draft TEXT NOT NULL DEFAULT '', state TEXT NOT NULL DEFAULT '{}', updated TEXT NOT NULL);
 CREATE TABLE IF NOT EXISTS settings (key TEXT PRIMARY KEY, value TEXT NOT NULL);`)
	if err == nil {
		_, err = db.Exec(`ALTER TABLE sessions ADD COLUMN attachments TEXT NOT NULL DEFAULT '[]'`)
		if err != nil && strings.Contains(err.Error(), "duplicate column") {
			err = nil
		}
	}
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Get(id string) (Session, error) {
	var session Session
	var messages, attachments, state string
	err := s.db.QueryRow("SELECT id,title,messages,draft,attachments,state,updated FROM sessions WHERE id=?", id).
		Scan(&session.ID, &session.Title, &messages, &session.Draft, &attachments, &state, &session.Updated)
	if err != nil {
		return session, err
	}
	if err = json.Unmarshal([]byte(messages), &session.Messages); err != nil {
		return session, err
	}
	if err = json.Unmarshal([]byte(attachments), &session.Attachments); err != nil {
		return session, err
	}
	err = json.Unmarshal([]byte(state), &session.State)
	session.Usage = session.State.Usage
	return session, err
}
