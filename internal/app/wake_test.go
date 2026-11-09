package app

import (
	"context"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/tasks"
)

type taskClient struct {
	mu       sync.Mutex
	requests [][]provider.Message
	code     string
}

func (c *taskClient) Complete(_ context.Context, request provider.Request) (provider.Message, provider.Usage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.requests = append(c.requests, request.Messages)
	if len(c.requests) == 1 {
		call := provider.ToolCall{ID: "call", Type: "function"}
		call.Function.Name = "js_repl"
		call.Function.Arguments = c.code
		return provider.Message{Role: "assistant", ToolCalls: []provider.ToolCall{call}}, provider.Usage{}, nil
	}
	return provider.Message{Role: "assistant", Content: "ok"}, provider.Usage{}, nil
}

func (c *taskClient) count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.requests)
}

func TestTasksWakeIdleAgentWithOneBatch(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	client := &taskClient{code: `{"code":"tasks.run({cmd: 'true', name: 'a'}); tasks.run({cmd: 'sleep 0.1', name: 'b'}).id"}`}
	service := &AgentService{store: store, agents: map[string]*agent.Agent{session.ID: agent.New(client)}}
	service.tasks = &tasks.Manager{Dir: t.TempDir(), Delay: 300 * time.Millisecond, Deliver: service.deliver}
	defer service.tasks.KillAll()
	history := []provider.Message{{Role: "user", Content: "Start"}}
	store.Save(session.ID, []sessions.Message{{Role: "user", Content: "Start"}}, "", nil)
	result, err := service.Run(context.Background(), session.ID, history, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = result
	deadline := time.Now().Add(5 * time.Second)
	for client.count() < 3 && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(500 * time.Millisecond)
	if client.count() != 3 {
		t.Fatalf("Expected one wake, got %d requests", client.count())
	}
	last := client.requests[2][len(client.requests[2])-1]
	if !strings.Contains(last.Content, "\"a\": exit 0") || !strings.Contains(last.Content, "\"b\": exit 0") {
		t.Fatalf("Unexpected notice: %+v", last)
	}
	saved, _ := store.Get(session.ID)
	kinds := []string{}
	for _, message := range saved.Messages {
		kinds = append(kinds, message.Role+":"+message.Kind)
	}
	if strings.Join(kinds, ",") != "user:,assistant:,user:task,assistant:" {
		t.Fatalf("Unexpected history: %v", kinds)
	}
}

func TestTaskNoticeJoinsRunningTurn(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	client := &taskClient{code: `{"code":"tasks.run({cmd: 'true', name: 'quick'}); bash('sleep 0.4'); 1"}`}
	service := &AgentService{store: store, agents: map[string]*agent.Agent{session.ID: agent.New(client)}}
	service.tasks = &tasks.Manager{Dir: t.TempDir(), Delay: 100 * time.Millisecond, Deliver: service.deliver}
	store.Save(session.ID, []sessions.Message{{Role: "user", Content: "Start"}}, "", nil)
	if _, err := service.Run(context.Background(), session.ID, []provider.Message{{Role: "user", Content: "Start"}}, nil); err != nil {
		t.Fatal(err)
	}
	time.Sleep(400 * time.Millisecond)
	if client.count() != 2 {
		t.Fatalf("Expected the notice inside the turn, got %d requests", client.count())
	}
	last := client.requests[1][len(client.requests[1])-1]
	if !strings.Contains(last.Content, "\"quick\": exit 0") {
		t.Fatalf("Unexpected last message: %+v", last)
	}
}

func TestStoppedSessionDoesNotWake(t *testing.T) {
	store, err := sessions.Open(filepath.Join(t.TempDir(), "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	client := &taskClient{code: `{"code":"1"}`}
	service := &AgentService{store: store, agents: map[string]*agent.Agent{session.ID: agent.New(client)}}
	service.tasks = &tasks.Manager{Dir: t.TempDir(), Delay: 50 * time.Millisecond, Deliver: service.deliver}
	service.paused = map[string]bool{session.ID: true}
	service.tasks.Session(session.ID).Run("", "true", "x", "", 0)
	time.Sleep(400 * time.Millisecond)
	if client.count() != 0 {
		t.Fatalf("Paused session woke up: %d", client.count())
	}
	store.Save(session.ID, []sessions.Message{{Role: "user", Content: "Next"}}, "", nil)
	if _, err := service.Run(context.Background(), session.ID, []provider.Message{{Role: "user", Content: "Next"}}, nil); err != nil {
		t.Fatal(err)
	}
	first := client.requests[0]
	if last := first[len(first)-1]; !strings.Contains(last.Content, "\"x\": exit 0") {
		t.Fatalf("Notice was not attached to the next message: %+v", last)
	}
}
