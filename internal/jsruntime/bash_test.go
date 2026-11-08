package jsruntime

import (
	"context"
	"strings"
	"testing"
)

func TestBashObjectArguments(t *testing.T) {
	r := New()
	if got := r.Run(context.Background(), `bash({cmd: "printf hi", timeout: 5, level: "medium"}).stdout`); got != `"hi"` {
		t.Fatalf("Unexpected output: %s", got)
	}
	if got := r.Run(context.Background(), `bash({cmd: "true", level: "huge"}).stderr`); !strings.Contains(got, "unknown level") {
		t.Fatalf("Unknown level was accepted: %s", got)
	}
}
