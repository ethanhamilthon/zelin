package jsruntime

import (
	"context"

	"github.com/dop251/goja"
	"zelin/internal/host"
)

type bashArgs struct {
	cmd     string
	timeout float64
	level   string
}

func (r *Runtime) bindBash(ctx context.Context) {
	r.vm.Set("bash", func(call goja.FunctionCall) goja.Value {
		args := r.bashArgs(call)
		if err := r.usable(); err != nil {
			return r.vm.ToValue(bashResult(host.Result{Stderr: err.Error(), ExitCode: -1}))
		}
		return r.vm.ToValue(bashResult(host.BashIn(ctx, r.cwd, args.cmd, args.timeout, args.level)))
	})
}

func (r *Runtime) bashArgs(call goja.FunctionCall) bashArgs {
	first := call.Argument(0)
	if _, ok := first.Export().(map[string]any); ok {
		obj := first.ToObject(r.vm)
		return bashArgs{cmd: exportString(obj.Get("cmd")), timeout: number(obj.Get("timeout")), level: exportString(obj.Get("level"))}
	}
	return bashArgs{cmd: exportString(first), timeout: number(call.Argument(1)), level: exportString(call.Argument(2))}
}

func exportString(v goja.Value) string {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return ""
	}
	return v.String()
}

func number(v goja.Value) float64 {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return 0
	}
	return v.ToFloat()
}

func bashResult(result host.Result) map[string]any {
	return map[string]any{"stdout": result.Stdout, "stderr": result.Stderr, "exitCode": result.ExitCode, "timedOut": result.TimedOut}
}
