package tasks

import (
	"fmt"
	"strings"
	"time"
)

func (s *Session) done(t *task, code int, timedOut bool) {
	s.mu.Lock()
	t.ended = time.Now()
	t.info.ExitCode = &code
	t.info.State = "exited"
	switch {
	case timedOut:
		t.info.State = "timed out"
	case t.killed:
		t.info.State = "killed"
	}
	if !t.silent {
		s.pending = append(s.pending, summary(t.snapshot()))
		if s.timer == nil {
			s.timer = time.AfterFunc(s.m.delay(), s.flush)
		}
	}
	s.mu.Unlock()
	s.m.changed(s.id)
}

func summary(info Info) string {
	result := fmt.Sprintf("exit %d", *info.ExitCode)
	if info.State != "exited" {
		result = info.State
	}
	return fmt.Sprintf("- #%d %q: %s after %s", info.ID, info.Name, result, time.Duration(info.Elapsed)*time.Second)
}

func (s *Session) Take(force bool) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timer != nil && !force {
		return ""
	}
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}
	if len(s.pending) == 0 {
		return ""
	}
	text := "Background tasks finished:\n" + strings.Join(s.pending, "\n") + "\nRead the logs with tasks.output(id)."
	s.pending = nil
	return text
}

func (s *Session) Ready() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.timer == nil && len(s.pending) > 0
}

func (s *Session) flush() {
	s.mu.Lock()
	s.timer = nil
	ready := len(s.pending) > 0
	s.mu.Unlock()
	if ready && s.m.Deliver != nil {
		s.m.Deliver(s.id)
	}
}
