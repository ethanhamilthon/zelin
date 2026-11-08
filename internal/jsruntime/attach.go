package jsruntime

import "zelin/internal/files"

func (r *Runtime) attach(name, path string) any {
	resolved, err := r.abs(path)
	if err != nil {
		return failure(err)
	}
	item, err := files.NewAttachment(name, resolved)
	if err != nil {
		return failure(err)
	}
	for i, previous := range r.attachments {
		if previous.Path == item.Path {
			r.attachments[i] = item
			return map[string]any{"ok": true, "name": item.Name, "path": item.Path}
		}
	}
	r.attachments = append(r.attachments, item)
	return map[string]any{"ok": true, "name": item.Name, "path": item.Path}
}

func (r *Runtime) TakeAttachments() []files.Attachment {
	result := r.attachments
	r.attachments = nil
	return result
}
