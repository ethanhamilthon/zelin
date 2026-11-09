package tasks

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MaxRunning = 8

type Session struct {
	m       *Manager
	id      string
	mu      sync.Mutex
	tasks   []*task
	pending []string
	timer   *time.Timer
}

func filepathID(id int) string { return strconv.Itoa(id) }

func (s *Session) Run(dir, cmd, name, level string, timeout float64) (Info, error) {
	if strings.TrimSpace(cmd) == "" {
		return Info{}, fmt.Errorf("cmd is required")
	}
	s.mu.Lock()
	if s.running() >= MaxRunning {
		s.mu.Unlock()
		return Info{}, fmt.Errorf("Too many running tasks (%d); wait or kill one", MaxRunning)
	}
	if name == "" {
		name = cmd
	}
	t := &task{info: Info{ID: len(s.tasks) + 1, Name: name, Cmd: cmd, Level: level, State: "running"}}
	s.tasks = append(s.tasks, t)
	s.mu.Unlock()
	path, err := s.m.logPath(s.id, t.info.ID)
	if err == nil {
		t.log = path
		err = s.start(t, dir, timeout)
	}
	if err != nil {
		s.mu.Lock()
		t.info.State, t.ended = "failed", time.Now()
		s.mu.Unlock()
		return Info{}, err
	}
	s.m.changed(s.id)
	return s.info(t), nil
}

func (s *Session) running() int {
	count := 0
	for _, t := range s.tasks {
		if t.info.State == "running" {
			count++
		}
	}
	return count
}

func (s *Session) info(t *task) Info {
	s.mu.Lock()
	defer s.mu.Unlock()
	return t.snapshot()
}

func (s *Session) List(active bool) []Info {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Info{}
	for _, t := range s.tasks {
		if !active || t.info.State == "running" {
			out = append(out, t.snapshot())
		}
	}
	return out
}

func (s *Session) find(id int) (*task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id < 1 || id > len(s.tasks) {
		return nil, fmt.Errorf("Task %d not found", id)
	}
	return s.tasks[id-1], nil
}

func (s *Session) Kill(id int, quiet bool) error {
	t, err := s.find(id)
	if err != nil {
		return err
	}
	s.mu.Lock()
	if t.info.State != "running" || t.cancel == nil {
		s.mu.Unlock()
		return fmt.Errorf("Task %d is not running", id)
	}
	t.killed, t.silent = true, quiet
	cancel := t.cancel
	s.mu.Unlock()
	cancel()
	return nil
}

func (s *Session) killAll() {
	for _, info := range s.List(true) {
		s.Kill(info.ID, true)
	}
}
