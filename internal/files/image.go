package files

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
)

const maxImageBytes = 20 << 20

type Image struct {
	Path     string
	MimeType string
	DataURL  string
}

func ReadImage(path string) (Image, error) {
	resolved, err := Resolve(path)
	if err != nil {
		return Image{}, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return Image{}, err
	}
	if info.IsDir() {
		return Image{}, fmt.Errorf("%s is a directory", resolved)
	}
	if info.Size() > maxImageBytes {
		return Image{}, fmt.Errorf("Image exceeds 20 MB: %s", resolved)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return Image{}, err
	}
	mime := http.DetectContentType(data)
	switch mime {
	case "image/png", "image/jpeg", "image/gif", "image/webp":
	default:
		return Image{}, fmt.Errorf("Unsupported image format %s; expected PNG, JPEG, GIF, or WebP", mime)
	}
	return Image{
		Path:     resolved,
		MimeType: mime,
		DataURL:  "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data),
	}, nil
}
