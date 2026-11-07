package sessions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"zash/internal/agent"
	"zash/internal/provider"
)

type Message struct {
	Role       string            `json:"role"`
	Content    string            `json:"content"`
	Executions []agent.Execution `json:"executions,omitempty"`
}

type Session struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Messages []Message       `json:"messages"`
	Draft    string          `json:"draft"`
	Updated  string          `json:"updated"`
	Usage    provider.Totals `json:"usage"`
	State    agent.State     `json:"-"`
}

type Store struct{ db *sql.DB }

func DefaultPath() (string, error) {
	base := os.Getenv("XDG_DATA_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(base, "zash", "sessions.db"), nil
}

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
	if err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) Get(id string) (Session, error) {
	var session Session
	var messages, state string
	err := s.db.QueryRow("SELECT id,title,messages,draft,state,updated FROM sessions WHERE id=?", id).
		Scan(&session.ID, &session.Title, &messages, &session.Draft, &state, &session.Updated)
	if err != nil {
		return session, err
	}
	if err = json.Unmarshal([]byte(messages), &session.Messages); err != nil {
		return session, err
	}
	err = json.Unmarshal([]byte(state), &session.State)
	session.Usage = session.State.Usage
	return session, err
}

func (s *Store) Save(id string, messages []Message, draft string) error {
	if messages == nil {
		messages = []Message{}
	}
	title := "New session"
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
	result, err := s.db.Exec("UPDATE sessions SET title=?,messages=?,draft=?,updated=strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id=?", title, string(data), draft, id)
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
