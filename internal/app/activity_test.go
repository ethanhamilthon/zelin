package app

import "testing"

func TestSharedActivityLifecycle(t *testing.T) {
	service := &AgentService{}
	service.beginLive("one")
	service.beginLive("two")
	service.endLive("one")
	snapshot := service.Activity()
	if !snapshot["two"].Running || !snapshot["one"].Unread {
		t.Fatal(snapshot)
	}
	delete(snapshot, "two")
	if !service.Activity()["two"].Running {
		t.Fatal("Shared mutable snapshot")
	}
	service.Read("one")
	if service.Activity()["one"].Unread {
		t.Fatal("Read did not clear unread")
	}
	service.Read("missing")
	service.beginLive("one")
	if service.Activity()["one"].Unread || !service.Activity()["one"].Running {
		t.Fatal("Restart did not clear unread")
	}
}
