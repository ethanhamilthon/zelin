package tasks

import (
	"strings"
	"sync"
	"testing"
	"time"
)

func manager(t *testing.T) (*Manager, chan string) {
	delivered := make(chan string, 4)
	m := &Manager{Dir: t.TempDir(), Delay: 150 * time.Millisecond}
	m.Deliver = func(session string) { delivered <- m.Session(session).Take(true) }
	return m, delivered
}

func TestTasksBatchFinishedNotices(t *testing.T) {
	m, delivered := manager(t)
	s := m.Session("a")
	s.Run("", "printf one", "first", "", 0)
	s.Run("", "sleep 0.05; exit 3", "second", "", 0)
	select {
	case text := <-delivered:
		if !strings.Contains(text, "#1 \"first\": exit 0") || !strings.Contains(text, "#2 \"second\": exit 3") {
			t.Fatalf("Unexpected notice: %s", text)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("No notice")
	}
	select {
	case text := <-delivered:
		t.Fatalf("Second notice: %s", text)
	case <-time.After(300 * time.Millisecond):
	}
	if out, _ := s.Output(1, 0); out != "one" {
		t.Fatalf("Unexpected output: %q", out)
	}
}

func TestTasksTimeoutAndKill(t *testing.T) {
	m, delivered := manager(t)
	s := m.Session("a")
	s.Run("", "sleep 30", "slow", "", 0.1)
	s.Run("", "sleep 30", "manual", "", 0)
	time.Sleep(50 * time.Millisecond)
	if err := s.Kill(2, false); err != nil {
		t.Fatal(err)
	}
	text := <-delivered
	if !strings.Contains(text, "\"slow\": timed out") || !strings.Contains(text, "\"manual\": killed") {
		t.Fatalf("Unexpected notice: %s", text)
	}
	if len(s.List(true)) != 0 {
		t.Fatalf("Tasks still running: %+v", s.List(true))
	}
}

func TestTasksWithoutTimeoutKeepRunning(t *testing.T) {
	m, _ := manager(t)
	s := m.Session("a")
	s.Run("", "sleep 30", "", "", 0)
	time.Sleep(200 * time.Millisecond)
	if list := s.List(true); len(list) != 1 || list[0].State != "running" {
		t.Fatalf("Unexpected list: %+v", list)
	}
	m.KillAll()
	time.Sleep(200 * time.Millisecond)
	if len(s.List(true)) != 0 || s.Take(true) != "" {
		t.Fatal("KillAll should stop tasks silently")
	}
}

func TestTasksTakeWaitsForBatch(t *testing.T) {
	m := &Manager{Dir: t.TempDir(), Delay: time.Hour}
	var once sync.Once
	m.Deliver = func(string) { once.Do(func() {}) }
	s := m.Session("a")
	s.Run("", "true", "", "", 0)
	time.Sleep(200 * time.Millisecond)
	if s.Take(false) != "" || s.Ready() {
		t.Fatal("Batch should still be collecting")
	}
	if text := s.Take(true); !strings.Contains(text, "#1") {
		t.Fatalf("Forced take: %q", text)
	}
}

func TestTasksLimit(t *testing.T) {
	m, _ := manager(t)
	s := m.Session("a")
	defer m.KillAll()
	for i := 0; i < MaxRunning; i++ {
		if _, err := s.Run("", "sleep 30", "", "", 0); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := s.Run("", "sleep 30", "", "", 0); err == nil {
		t.Fatal("Limit was not enforced")
	}
}
