package host

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
)

type limit struct {
	memoryMB int
	cpuShare float64
}

var levels = map[string]limit{
	"low":    {memoryMB: 1024, cpuShare: 1.0 / 8},
	"medium": {memoryMB: 2048, cpuShare: 1.0 / 4},
	"high":   {memoryMB: 3072, cpuShare: 1.0 / 2},
}

var scopes = sync.OnceValue(func() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	return exec.Command("systemd-run", "--user", "--scope", "--quiet", "true").Run() == nil
})

func resolveLevel(level string) (string, error) {
	if level == "" {
		return "low", nil
	}
	if _, ok := levels[level]; !ok {
		return "", fmt.Errorf("unknown level %q, use low, medium or high", level)
	}
	return level, nil
}

func shell(level, command string) (string, []string) {
	if !scopes() {
		return "bash", []string{"-c", command}
	}
	l := levels[level]
	quota := max(int(l.cpuShare*float64(runtime.NumCPU())*100), 10)
	return "systemd-run", []string{"--user", "--scope", "--quiet",
		"-p", fmt.Sprintf("MemoryMax=%dM", l.memoryMB),
		"-p", "MemorySwapMax=0",
		"-p", fmt.Sprintf("CPUQuota=%d%%", quota),
		"-p", "TasksMax=512",
		"bash", "-c", command}
}

func killedByLimit(err error, code int) bool {
	if !scopes() {
		return false
	}
	if code == 137 {
		return true
	}
	exit, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	status, ok := exit.Sys().(syscall.WaitStatus)
	return ok && status.Signaled() && status.Signal() == syscall.SIGKILL
}

func limitNote(level string) string {
	next := map[string]string{"low": ` Retry with level "medium" or "high".`, "medium": ` Retry with level "high".`}[level]
	return fmt.Sprintf("\nKilled, likely by the %s level memory limit (%dM).%s", level, levels[level].memoryMB, next)
}
