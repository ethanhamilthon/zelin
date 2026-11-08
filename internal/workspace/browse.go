package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const ListLimit = 500

type Entry struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	Link   string `json:"link,omitempty"`
	Broken bool   `json:"broken,omitempty"`
	Size   int64  `json:"size"`
}

type Listing struct {
	Dir       string  `json:"dir"`
	Parent    string  `json:"parent"`
	Entries   []Entry `json:"entries"`
	Truncated bool    `json:"truncated"`
}

func Inside(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func Within(root, path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if !Inside(root, abs) {
		return "", fmt.Errorf("%s is outside the workspace", abs)
	}
	return abs, nil
}

func entry(dir string, item os.DirEntry) Entry {
	path := filepath.Join(dir, item.Name())
	out := Entry{Name: item.Name(), Path: path, Kind: "file"}
	info, err := item.Info()
	if err == nil {
		out.Size = info.Size()
	}
	if item.Type()&os.ModeSymlink != 0 {
		target, err := filepath.EvalSymlinks(path)
		if err != nil {
			out.Broken, out.Kind = true, "file"
			return out
		}
		out.Link = target
		if stat, err := os.Stat(target); err == nil {
			out.Size = stat.Size()
			if stat.IsDir() {
				out.Kind = "dir"
			}
		}
		return out
	}
	if item.IsDir() {
		out.Kind = "dir"
	}
	return out
}

func List(root, dir string) (Listing, error) {
	abs, err := Within(root, dir)
	if err != nil {
		return Listing{}, err
	}
	if real, err := filepath.EvalSymlinks(abs); err != nil {
		return Listing{}, err
	} else if base, err := filepath.EvalSymlinks(root); err != nil {
		return Listing{}, err
	} else if !Inside(base, real) {
		return Listing{}, fmt.Errorf("%s links outside the workspace to %s", abs, real)
	}
	items, err := os.ReadDir(abs)
	if err != nil {
		return Listing{}, err
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].IsDir() != items[j].IsDir() {
			return items[i].IsDir()
		}
		return strings.ToLower(items[i].Name()) < strings.ToLower(items[j].Name())
	})
	listing := Listing{Dir: abs, Entries: []Entry{}}
	if abs != root {
		listing.Parent = filepath.Dir(abs)
	}
	for _, item := range items {
		if len(listing.Entries) >= ListLimit {
			listing.Truncated = true
			break
		}
		listing.Entries = append(listing.Entries, entry(abs, item))
	}
	return listing, nil
}
