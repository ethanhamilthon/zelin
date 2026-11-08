package host

import (
	"context"
	"strings"
	"testing"
)

func TestBashRejectsUnknownLevel(t *testing.T) {
	result := BashIn(context.Background(), "", "true", 1, "huge")
	if result.ExitCode != -1 || !strings.Contains(result.Stderr, "unknown level") {
		t.Fatalf("Unexpected result: %+v", result)
	}
}

func TestBashLevelsRun(t *testing.T) {
	for level := range levels {
		if result := BashIn(context.Background(), "", "printf ok", 5, level); result.Stdout != "ok" || result.ExitCode != 0 {
			t.Fatalf("%s: %+v", level, result)
		}
	}
}

func TestBashMemoryLimit(t *testing.T) {
	if !scopes() {
		t.Skip("systemd user scopes are unavailable")
	}
	if cgroup := BashIn(context.Background(), "", "cat /proc/self/cgroup", 5, "low").Stdout; !strings.Contains(cgroup, "run-") {
		t.Fatalf("Command did not run in a scope: %s", cgroup)
	}
	result := BashIn(context.Background(), "", "python3 -c 'x = bytearray(1500 * 1024 * 1024); x[::4096] = b\"1\" * len(x[::4096])'", 60, "low")
	if result.ExitCode == 0 || !strings.Contains(result.Stderr, "memory limit (1024M)") {
		t.Fatalf("Memory limit was not enforced: %+v", result)
	}
}
