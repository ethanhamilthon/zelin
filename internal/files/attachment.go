package files

import (
	"fmt"
	"os"
	"strings"
)

type Attachment struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func NewAttachment(name, path string) (Attachment, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Attachment{}, fmt.Errorf("Attachment name is required")
	}
	info, err := os.Stat(path)
	if err != nil {
		return Attachment{}, err
	}
	if !info.Mode().IsRegular() {
		return Attachment{}, fmt.Errorf("Attachment must be a regular file: %s", path)
	}
	return Attachment{Name: name, Path: path}, nil
}
