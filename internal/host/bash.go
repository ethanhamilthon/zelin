package host

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
	TimedOut bool   `json:"timedOut"`
}

type output struct {
	sync.Mutex
	data []byte
}

func (b *output) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	n := len(p)
	if remaining := 64000 - len(b.data); remaining > 0 {
		if len(p) > remaining {
			p = p[:remaining]
		}
		b.data = append(b.data, p...)
	}
	return n, nil
}

func Bash(parent context.Context, command string, timeout float64) Result {
	return BashIn(parent, "", command, timeout, "")
}

func BashIn(parent context.Context, dir, command string, timeout float64, level string) Result {
	level, err := resolveLevel(level)
	if err != nil {
		return Result{Stderr: err.Error(), ExitCode: -1}
	}
	if timeout <= 0 || timeout > 120 {
		timeout = 30
	}
	ctx, cancel := context.WithTimeout(parent, time.Duration(timeout*float64(time.Second)))
	defer cancel()
	name, args := shell(level, command)
	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "PWD="+dir)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = time.Second
	stdout, stderr := &output{}, &output{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err = cmd.Run()
	result := Result{Stdout: string(stdout.data), Stderr: string(stderr.data), TimedOut: errors.Is(ctx.Err(), context.DeadlineExceeded)}
	if err != nil {
		result.ExitCode = -1
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			result.ExitCode = exit.ExitCode()
		} else {
			result.Stderr += err.Error()
		}
		if ctx.Err() == nil && killedByLimit(err, result.ExitCode) {
			result.Stderr += limitNote(level)
		}
	}
	return result
}
