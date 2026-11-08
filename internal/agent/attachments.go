package agent

import "zelin/internal/files"

func (a *Agent) TakeAttachments() []files.Attachment { return a.runtime.TakeAttachments() }
