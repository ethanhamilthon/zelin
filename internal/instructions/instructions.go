package instructions

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Document struct {
	File    string
	Default string
}

type Store struct {
	Dir       string
	Documents map[string]Document
}

func (s *Store) document(name string) (Document, string, error) {
	document, ok := s.Documents[name]
	if !ok {
		return Document{}, "", fmt.Errorf("Unknown document %s", name)
	}
	return document, filepath.Join(s.Dir, document.File), nil
}

func (s *Store) Read(name string) (string, error) {
	document, path, err := s.document(name)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return document.Default, nil
	}
	return string(data), err
}

func (s *Store) Write(name, text string) error {
	_, path, err := s.document(name)
	if err != nil {
		return err
	}
	temp := path + ".tmp"
	if err := os.WriteFile(temp, []byte(text), 0600); err != nil {
		return err
	}
	return os.Rename(temp, path)
}
