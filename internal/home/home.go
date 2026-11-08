package home

import (
	"os"
	"path/filepath"
)

var Name = "zelin"

const Legacy = "zash"

func Dir() (string, error) {
	if dir := os.Getenv("ZELIN_HOME"); dir != "" {
		return dir, nil
	}
	if dir := os.Getenv("ZASH_HOME"); dir != "" {
		return dir, nil
	}
	return base(Name)
}

func base(name string) (string, error) {
	root := os.Getenv("XDG_DATA_HOME")
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		root = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(root, name), nil
}

func Path(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, name), nil
}
