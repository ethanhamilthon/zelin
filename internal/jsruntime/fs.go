package jsruntime

import (
	"github.com/dop251/goja"
	"zelin/internal/files"
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
			path, err := r.abs(path)
			if err != nil {
				return failure(err)
			}
			text, err := files.Read(path, opts.Offset, opts.Limit)
			if err != nil {
				return failure(err)
			}
			return text
		},
		"write": func(path, text string) any {
			return r.mutate(path, func(p string) error { return files.Write(p, text) })
		},
		"edit": func(path, old, text string) any {
			return r.mutate(path, func(p string) error { return files.Edit(p, old, text) })
		},
		"patch": func(path, diff string) any {
			return r.mutate(path, func(p string) error { return files.Patch(p, diff) })
		},
		"grep": func(pattern, path string, value goja.Value) any {
			if path == "" {
				path = "."
			}
			path, err := r.abs(path)
			if err != nil {
				return failure(err)
			}
			results, err := files.Grep(pattern, path, r.options(value).Glob)
			if err != nil {
				return failure(err)
			}
			return results
		},
	}
}

func (r *Runtime) mutate(path string, change func(string) error) any {
	resolved, err := r.abs(path)
	if err != nil {
		return failure(err)
	}
	return done(change(resolved))
}
