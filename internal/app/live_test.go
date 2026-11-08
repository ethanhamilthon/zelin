package app

import (
	"testing"
	"zelin/internal/agent"
)

func TestLiveSnapshotLifecycle(t *testing.T) {
	var names []string
	service := &AgentService{notify: func(name string, _ any) { names = append(names, name) }}
	service.beginLive("one")
	service.recordLive("one", agent.Event{Kind: "activity", Activity: "Thinking"})
	service.recordLive("one", agent.Event{Kind: "activity", Activity: "Working"})
	service.recordLive("one", agent.Event{Kind: "execution", Execution: &agent.Execution{ID: "tool", State: "running"}})
	service.recordLive("one", agent.Event{Kind: "execution", Execution: &agent.Execution{ID: "tool", State: "complete"}})
	snapshot := service.Active("one")
	if snapshot == nil || snapshot.Started == 0 || len(snapshot.Events) != 2 || snapshot.Events[1].Execution.State != "complete" {
		t.Fatalf("Invalid snapshot: %+v", snapshot)
	}
	snapshot.Events[0].Activity = "Changed"
	if service.Active("one").Events[0].Activity != "Working" {
		t.Fatal("Snapshot shares event slice")
	}
	service.endLive("one")
	if service.Active("one") != nil || len(names) != 2 {
		t.Fatal("Run did not finish")
	}
}
