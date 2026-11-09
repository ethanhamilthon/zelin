package tasks

import (
	"context"
	"errors"
	"os"
	"time"

	"zelin/internal/host"
)

type Info struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Cmd      string  `json:"cmd"`
	Level    string  `json:"level"`
	State    string  `json:"state"`
	ExitCode *int    `json:"exitCode,omitempty"`
	Elapsed  float64 `json:"elapsed"`
}

type task struct {
	info    Info
	log     string
	started time.Time
	ended   time.Time
	cancel  context.CancelFunc
	killed  bool
	silent  bool
}

func (t *task) snapshot() Info {
	info := t.info
	end := t.ended
	if end.IsZero() {
		end = time.Now()
	}
	info.Elapsed = end.Sub(t.started).Round(time.Second).Seconds()
	return info
}

func (s *Session) start(t *task, dir string, timeout float64) error {
	ctx, cancel := context.WithCancel(context.Background())
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout*float64(time.Second)))
	}
	cmd, level, err := host.Command(ctx, dir, t.info.Cmd, t.info.Level)
	if err != nil {
		cancel()
		return err
	}
	file, err := os.Create(t.log)
	if err != nil {
		cancel()
		return err
	}
	cmd.Stdout, cmd.Stderr = file, file
	if err := cmd.Start(); err != nil {
		cancel()
		file.Close()
		return err
	}
	t.info.Level, t.cancel, t.started = level, cancel, time.Now()
	go func() {
		err := cmd.Wait()
		code, _ := host.ExitCode(err)
		if err != nil && ctx.Err() == nil {
			file.WriteString(host.LimitNote(err, code, level))
		}
		file.Close()
		cancel()
		s.done(t, code, errors.Is(ctx.Err(), context.DeadlineExceeded))
	}()
	return nil
}
