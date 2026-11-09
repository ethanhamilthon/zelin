package host

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func Command(ctx context.Context, dir, command, level string) (*exec.Cmd, string, error) {
	level, err := resolveLevel(level)
	if err != nil {
		return nil, "", err
	}
	name, args := shell(level, command)
	cmd := exec.CommandContext(ctx, name, args...)
	if dir != "" {
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "PWD="+dir)
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = time.Second
	return cmd, level, nil
}

func ExitCode(err error) (int, bool) {
	if err == nil {
		return 0, true
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode(), true
	}
	return -1, false
}

func LimitNote(err error, code int, level string) string {
	if killedByLimit(err, code) {
		return limitNote(level)
	}
	return ""
}
