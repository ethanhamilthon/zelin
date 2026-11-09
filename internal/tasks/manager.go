package tasks

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Manager struct {
	Dir     string
	Delay   time.Duration
	Deliver func(session string)
	Changed func(session string)
	mu      sync.Mutex
	items   map[string]*Session
}

func (m *Manager) Session(id string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.items == nil {
		m.items = map[string]*Session{}
	}
	if m.items[id] == nil {
		m.items[id] = &Session{m: m, id: id}
	}
	return m.items[id]
}

func (m *Manager) Find(id string) *Session {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.items[id]
}

func (m *Manager) KillAll() {
	m.mu.Lock()
	all := make([]*Session, 0, len(m.items))
	for _, s := range m.items {
		all = append(all, s)
	}
	m.mu.Unlock()
	for _, s := range all {
		s.killAll()
	}
}

func (m *Manager) Clean() error {
	if m.Dir == "" {
		return nil
	}
	return os.RemoveAll(m.Dir)
}

func (m *Manager) logPath(session string, id int) (string, error) {
	dir := filepath.Join(m.Dir, session)
	if m.Dir == "" {
		dir = filepath.Join(os.TempDir(), "zelin-tasks", session)
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, filepathID(id)+".log"), nil
}

func (m *Manager) changed(session string) {
	if m.Changed != nil {
		m.Changed(session)
	}
}

func (m *Manager) delay() time.Duration {
	if m.Delay > 0 {
		return m.Delay
	}
	return 3 * time.Second
}
