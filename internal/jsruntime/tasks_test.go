package jsruntime

import (
	"context"
	"strings"
	"testing"
	"time"

	"zelin/internal/tasks"
)

func TestTasksAPI(t *testing.T) {
	m := &tasks.Manager{Dir: t.TempDir(), Delay: time.Hour}
	defer m.KillAll()
	r := New()
	r.Tasks = m.Session("s")
	if got := r.Run(context.Background(), `JSON.stringify(tasks.run({cmd: "printf hi", name: "greet", level: "medium"}))`); !strings.Contains(got, `\"id\":1`) {
		t.Fatalf("run: %s", got)
	}
	time.Sleep(200 * time.Millisecond)
	if got := r.Run(context.Background(), `tasks.output(1)`); got != `"hi"` {
		t.Fatalf("output: %s", got)
	}
	if got := r.Run(context.Background(), `tasks.list()[0].state + " " + tasks.list()[0].name`); got != `"exited greet"` {
		t.Fatalf("list: %s", got)
	}
	r.Run(context.Background(), `tasks.run("sleep 30")`)
	if got := r.Run(context.Background(), `tasks.kill(2).ok`); got != "true" {
		t.Fatalf("kill: %s", got)
	}
}
