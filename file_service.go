package main

import (
	"fmt"
	"os"
	"strings"

	"zelin/internal/files"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type FileService struct{}

func (s *FileService) Read(path string) (string, error) {
	path, err := files.Resolve(path)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	return string(data), err
}

func (s *FileService) Write(path, text string) error { return files.Write(path, text) }

func (s *FileService) Resolve(path string) (string, error) {
	return files.Locate(path)
}

func (s *FileService) Open(path string) error {
	path, err := s.Resolve(path)
	if err != nil {
		return err
	}
	return application.Get().Browser.OpenFile(path)
}

func (s *FileService) OpenURL(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "mailto:") {
		return fmt.Errorf("Unsupported link: %s", url)
	}
	return application.Get().Browser.OpenURL(url)
}
