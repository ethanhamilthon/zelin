package jsruntime

import (
	"context"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"zelin/internal/files"
)

type Runtime struct {
	cwd         string
	vm          *goja.Runtime
	stringify   goja.Callable
	images      []files.Image
	attachments []files.Attachment
}

func New() *Runtime { return NewIn(launch) }

func NewIn(dir string) *Runtime {
	if dir == "" {
		dir = launch
	}
	vm := goja.New()
	stringify, _ := goja.AssertFunction(vm.Get("JSON").ToObject(vm).Get("stringify"))
	r := &Runtime{vm: vm, stringify: stringify, cwd: dir}
	vm.Set("fs", r.fs())
	vm.Set("view_image", r.viewImage)
	vm.Set("attach", r.attach)
	vm.Set("cd", r.chdir)
	return r
}

func (r *Runtime) viewImage(path string) any {
	resolved, err := r.abs(path)
	if err != nil {
		return failure(err)
	}
	img, err := files.ReadImage(resolved)
	if err != nil {
		return failure(err)
	}
	r.images = append(r.images, img)
	return map[string]any{"ok": true, "path": img.Path, "mimeType": img.MimeType}
}

func (r *Runtime) TakeImages() []files.Image {
	out := r.images
	r.images = nil
	return out
}

func (r *Runtime) Run(ctx context.Context, code string) string {
	ctx, cancel := context.WithTimeout(ctx, 125*time.Second)
	defer cancel()
	r.bindBash(ctx)
	r.bindMCP(ctx)
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
	return result
}
