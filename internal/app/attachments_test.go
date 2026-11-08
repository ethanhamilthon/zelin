package app

import (
	"context"
	"path/filepath"
	"testing"
	"zelin/internal/agent"
	"zelin/internal/provider"
	"zelin/internal/sessions"
)

type attachmentClient struct{ calls int }

func (c *attachmentClient) Complete(context.Context, provider.Request) (provider.Message, provider.Usage, error) {
	c.calls++
	if c.calls == 1 {
		return provider.Message{Role: "assistant", ToolCalls: []provider.ToolCall{{ID: "file", Type: "function", Function: struct {
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		}{Name: "js_repl", Arguments: `{"code":"fs.write('report.md', '# Report'); attach('Report', 'report.md')"}`}}}}, provider.Usage{}, nil
	}
	return provider.Message{Role: "assistant", Content: "Your report"}, provider.Usage{}, nil
}
func TestAttachedFilesPersistWithAnswer(t *testing.T) {
	dir := t.TempDir()
	store, err := sessions.Open(filepath.Join(dir, "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	session, _ := store.Create()
	service := &AgentService{store: store, agents: map[string]*agent.Agent{session.ID: agent.NewIn(&attachmentClient{}, dir)}}
	result, err := service.Run(context.Background(), session.ID, []provider.Message{{Role: "user", Content: "Create a report"}}, nil)
	if err != nil || len(result.Files) != 1 {
		t.Fatalf("%+v %v", result, err)
	}
	stored, err := store.Get(session.ID)
	if err != nil || len(stored.Messages) != 1 || len(stored.Messages[0].Files) != 1 {
		t.Fatalf("%+v %v", stored, err)
	}
	if stored.Messages[0].Files[0] != result.Files[0] {
		t.Fatal("Attachment changed during persistence")
	}
}
