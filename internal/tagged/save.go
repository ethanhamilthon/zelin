package tagged

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

func imagePath(path, assets string, ref Ref) string {
	if ref.Data == "" {
		return path
	}
	_, payload, ok := strings.Cut(ref.Data, ",")
	if !ok || assets == "" {
		return ""
	}
	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return ""
	}
	name := unsafeName.ReplaceAllString(filepath.Base(ref.Label), "_")
	target := filepath.Join(assets, ref.ID+"-"+name)
	if os.MkdirAll(assets, 0700) != nil || os.WriteFile(target, data, 0600) != nil {
		return ""
	}
	return target
}
