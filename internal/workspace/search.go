package workspace

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Suggest(root, dir, query string) ([]Entry, error) {
	base := root
	if dir != "" {
		var err error
		if base, err = Within(root, dir); err != nil {
			return nil, err
		}
	}
	if slash := strings.LastIndex(query, "/"); slash >= 0 {
		sub, err := Within(root, filepath.Join(base, query[:slash]))
		if err != nil {
			return nil, err
		}
		base, query = sub, query[slash+1:]
	}
	items, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	needle := strings.ToLower(query)
	type scored struct {
		entry Entry
		rank  int
	}
	found := []scored{}
	for _, item := range items {
		name := strings.ToLower(item.Name())
		if strings.HasPrefix(item.Name(), ".") && !strings.HasPrefix(needle, ".") {
			continue
		}
		rank := 2
		switch {
		case needle == "" || strings.HasPrefix(name, needle):
			rank = 0
		case strings.Contains(name, needle):
			rank = 1
		default:
			continue
		}
		found = append(found, scored{entry(base, item), rank})
	}
	sort.SliceStable(found, func(i, j int) bool {
		a, b := found[i], found[j]
		if a.rank != b.rank {
			return a.rank < b.rank
		}
		if (a.entry.Kind == "dir") != (b.entry.Kind == "dir") {
			return a.entry.Kind == "dir"
		}
		return strings.ToLower(a.entry.Name) < strings.ToLower(b.entry.Name)
	})
	out := []Entry{}
	for _, item := range found {
		out = append(out, item.entry)
	}
	return out, nil
}
