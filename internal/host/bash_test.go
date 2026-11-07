package host

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBashOutputAndExit(t *testing.T) {
	result := Bash(context.Background(), "printf hello; printf problem >&2; exit 7", 1)
	if result.Stdout != "hello" || result.Stderr != "problem" || result.ExitCode != 7 || result.TimedOut {
		t.Fatalf("Unexpected result: %+v", result)
	}
}

func TestBashTimeout(t *testing.T) {
	start := time.Now()
	result := Bash(context.Background(), "sleep 20 & wait", 0.05)
	if !result.TimedOut || result.ExitCode == 0 || time.Since(start) > 2*time.Second {
		t.Fatalf("Timeout did not terminate the process group: %+v", result)
	}
}

func TestBashLimitsOutput(t *testing.T) {
	result := Bash(context.Background(), "yes x | head -c 100000", 1)
	if len(result.Stdout) != 64000 || !strings.HasPrefix(result.Stdout, "x\n") {
		t.Fatalf("Output was not bounded: %d", len(result.Stdout))
	}
}
