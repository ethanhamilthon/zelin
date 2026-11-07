package jsruntime

import (
	"context"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"zash/internal/host"
)

type Runtime struct {
	vm        *goja.Runtime
	stringify goja.Callable
}

func New() *Runtime {
	vm := goja.New()
	stringify, _ := goja.AssertFunction(vm.Get("JSON").ToObject(vm).Get("stringify"))
	return &Runtime{vm: vm, stringify: stringify}
}

func (r *Runtime) Run(ctx context.Context, code string) string {
	ctx, cancel := context.WithTimeout(ctx, 125*time.Second)
	defer cancel()
	r.vm.Set("bash", func(command string, timeout float64) map[string]any {
		result := host.Bash(ctx, command, timeout)
		return map[string]any{"stdout": result.Stdout, "stderr": result.Stderr, "exitCode": result.ExitCode, "timedOut": result.TimedOut}
	})
	done := make(chan struct{})
	stop := context.AfterFunc(ctx, func() {
		r.vm.Interrupt("JavaScript execution timed out or was cancelled")
		close(done)
	})
	defer func() {
		if !stop() {
			<-done
		}
		r.vm.ClearInterrupt()
	}()
	value, err := r.vm.RunString(code)
	if err != nil {
		return fmt.Sprintf("JavaScript error: %v", err)
	}
	if value == nil || goja.IsUndefined(value) {
		return "undefined"
	}
	serialized, err := r.stringify(goja.Undefined(), value)
	if err != nil {
		return fmt.Sprintf("Cannot serialize JavaScript result: %v", err)
	}
	result := serialized.String()
	if len(result) > 128000 {
		return result[:128000] + "\n[output truncated]"
	}
	return result
}
