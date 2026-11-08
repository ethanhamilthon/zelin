package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tree(t *testing.T) string {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, "src", "deep"), 0755)
	os.MkdirAll(filepath.Join(root, "my docs"), 0755)
	os.WriteFile(filepath.Join(root, "README.md"), []byte("hi"), 0644)
	os.WriteFile(filepath.Join(root, "src", "main.go"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(root, "src", "deep", "hidden.txt"), []byte("x"), 0644)
	os.WriteFile(filepath.Join(root, ".env"), []byte("x"), 0644)
	return root
}

func names(entries []Entry) string {
	out := []string{}
	for _, e := range entries {
		out = append(out, e.Name+":"+e.Kind)
	}
	return strings.Join(out, ",")
}

func TestListIsShallowSortedAndNavigable(t *testing.T) {
	root := tree(t)
	top, err := List(root, root)
	if err != nil || top.Parent != "" {
		t.Fatalf("%+v %v", top, err)
	}
	if got := names(top.Entries); got != "my docs:dir,src:dir,.env:file,README.md:file" {
		t.Fatal(got)
	}
	sub, err := List(root, filepath.Join(root, "src"))
	if err != nil || sub.Parent != root || names(sub.Entries) != "deep:dir,main.go:file" {
		t.Fatalf("%+v %v", sub, err)
	}
}

func TestListRejectsEscapes(t *testing.T) {
	root := tree(t)
	for _, dir := range []string{filepath.Dir(root), filepath.Join(root, ".."), "/etc"} {
		if _, err := List(root, dir); err == nil {
			t.Fatalf("Listed %s", dir)
		}
	}
	if _, err := List(root, root+"-sibling"); err == nil {
		t.Fatal("A sibling prefix must not count as inside")
	}
}

func TestSymlinksShowTargetsAndDoNotChangeRoot(t *testing.T) {
	root := tree(t)
	outside := t.TempDir()
	os.WriteFile(filepath.Join(outside, "o.txt"), []byte("x"), 0644)
	os.Symlink(outside, filepath.Join(root, "out"))
	os.Symlink(filepath.Join(root, "README.md"), filepath.Join(root, "link.md"))
	os.Symlink(filepath.Join(root, "missing"), filepath.Join(root, "broken"))
	os.Symlink(root, filepath.Join(root, "src", "loop"))
	top, _ := List(root, root)
	found := map[string]Entry{}
	for _, e := range top.Entries {
		found[e.Name] = e
	}
	resolvedOut, _ := filepath.EvalSymlinks(outside)
	if found["out"].Kind != "dir" || found["out"].Link != resolvedOut {
		t.Fatalf("%+v", found["out"])
	}
	if found["link.md"].Kind != "file" || !strings.HasSuffix(found["link.md"].Link, "README.md") {
		t.Fatalf("%+v", found["link.md"])
	}
	if !found["broken"].Broken {
		t.Fatalf("%+v", found["broken"])
	}
	sub, err := List(root, filepath.Join(root, "src"))
	if err != nil || sub.Parent != root {
		t.Fatalf("A cyclic link must not break listing: %v", err)
	}
}

func TestLargeDirectoriesAreBounded(t *testing.T) {
	root := t.TempDir()
	for i := 0; i < ListLimit+50; i++ {
		os.WriteFile(filepath.Join(root, fmt.Sprintf("f%04d", i)), nil, 0644)
	}
	listing, err := List(root, root)
	if err != nil || len(listing.Entries) != ListLimit || !listing.Truncated {
		t.Fatalf("%d %v %v", len(listing.Entries), listing.Truncated, err)
	}
}

func TestPermissionErrorIsReported(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("root ignores permissions")
	}
	root := t.TempDir()
	locked := filepath.Join(root, "locked")
	os.Mkdir(locked, 0000)
	defer os.Chmod(locked, 0755)
	if _, err := List(root, locked); err == nil {
		t.Fatal("Listed an unreadable directory")
	}
}
