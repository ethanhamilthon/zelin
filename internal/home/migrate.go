package home

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Migration struct {
	From, To string
	Copied   bool
}

func legacyName() string { return strings.Replace(Name, "zelin", Legacy, 1) }

func MigrateLegacy(snapshot func(from, to string) error) (Migration, error) {
	if os.Getenv("ZELIN_HOME") != "" || os.Getenv("ZASH_HOME") != "" {
		return Migration{}, nil
	}
	to, err := Dir()
	if err != nil {
		return Migration{}, err
	}
	from, err := base(legacyName())
	if err != nil {
		return Migration{}, err
	}
	result := Migration{From: from, To: to}
	if _, err := os.Stat(from); err != nil {
		return result, nil
	}
	if _, err := os.Stat(filepath.Join(to, "sessions.db")); err == nil {
		return result, nil
	}
	if entries, err := os.ReadDir(to); err == nil && len(entries) > 0 {
		return result, fmt.Errorf("%s already has data; not overwriting it with %s", to, from)
	}
	if err := copyTree(from, to, snapshot); err != nil {
		os.RemoveAll(to)
		return result, fmt.Errorf("Could not copy %s to %s: %w", from, to, err)
	}
	result.Copied = true
	return result, nil
}

func copyTree(from, to string, snapshot func(from, to string) error) error {
	return filepath.WalkDir(from, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(from, path)
		target := filepath.Join(to, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o700)
		}
		if strings.HasSuffix(entry.Name(), "-wal") || strings.HasSuffix(entry.Name(), "-shm") || !entry.Type().IsRegular() {
			return nil
		}
		if entry.Name() == "sessions.db" {
			return snapshot(path, target)
		}
		return copyFile(path, target)
	})
}

func copyFile(from, to string) error {
	in, err := os.Open(from)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(to, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}
