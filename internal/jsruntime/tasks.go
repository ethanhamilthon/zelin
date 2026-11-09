package jsruntime

import (
	"fmt"

	"github.com/dop251/goja"
)

func (r *Runtime) bindTasks() {
	r.vm.Set("tasks", map[string]any{
		"run": func(call goja.FunctionCall) goja.Value { return r.vm.ToValue(r.runTask(call)) },
		"list": func() any {
			if r.Tasks == nil {
				return []any{}
			}
			return plain(r.Tasks.List(false))
		},
		"output": func(id int, options *goja.Object) any {
			if r.Tasks == nil {
				return failure(fmt.Errorf("Background tasks are unavailable"))
			}
			tail := 0.0
			if options != nil {
				tail = number(options.Get("tail"))
			}
			text, err := r.Tasks.Output(id, int(tail))
			if err != nil {
				return failure(err)
			}
			return text
		},
		"kill": func(id int) any {
			if r.Tasks == nil {
				return failure(fmt.Errorf("Background tasks are unavailable"))
			}
			if err := r.Tasks.Kill(id, false); err != nil {
				return failure(err)
			}
			return map[string]any{"ok": true}
		},
	})
}

func (r *Runtime) runTask(call goja.FunctionCall) any {
	if r.Tasks == nil {
		return failure(fmt.Errorf("Background tasks are unavailable"))
	}
	if err := r.usable(); err != nil {
		return failure(err)
	}
	args := r.bashArgs(call)
	name := ""
	if obj, ok := call.Argument(0).(*goja.Object); ok {
		name = exportString(obj.Get("name"))
	}
	info, err := r.Tasks.Run(r.cwd, args.cmd, name, args.level, args.timeout)
	if err != nil {
		return failure(err)
	}
	return plain(info)
}
