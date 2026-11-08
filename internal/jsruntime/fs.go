package jsruntime

import (
	"github.com/dop251/goja"
	"zash/internal/files"
)

type options struct {
	Offset, Limit int
	Glob          string
}

func (r *Runtime) options(value goja.Value) options {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return options{}
	}
	object := value.ToObject(r.vm)
	field := func(name string) goja.Value {
		if item := object.Get(name); item != nil && !goja.IsUndefined(item) {
			return item
		}
		return r.vm.ToValue("")
	}
	return options{Offset: int(field("offset").ToInteger()), Limit: int(field("limit").ToInteger()), Glob: field("glob").String()}
}

func failure(err error) map[string]any { return map[string]any{"error": err.Error()} }

func done(err error) map[string]any {
	if err != nil {
		return failure(err)
	}
	return map[string]any{"ok": true}
}

func (r *Runtime) fs() map[string]any {
	return map[string]any{
		"read": func(path string, value goja.Value) any {
			opts := r.options(value)
			text, err := files.Read(path, opts.Offset, opts.Limit)
			if err != nil {
				return failure(err)
			}
			return text
		},
		"write": func(path, text string) any { return done(files.Write(path, text)) },
		"edit":  func(path, old, text string) any { return done(files.Edit(path, old, text)) },
		"patch": func(path, diff string) any { return done(files.Patch(path, diff)) },
		"grep": func(pattern, path string, value goja.Value) any {
			results, err := files.Grep(pattern, path, r.options(value).Glob)
			if err != nil {
				return failure(err)
			}
			return results
		},
	}
}
