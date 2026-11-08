package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func file(t *testing.T, text string) string {
	path := filepath.Join(t.TempDir(), "a.txt")
	if err := os.WriteFile(path, []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func content(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestReadNumbersAndPages(t *testing.T) {
	path := file(t, "one\ntwo\nthree\n")
	text, err := Read(path, 0, 0)
	if err != nil || text != "1\tone\n2\ttwo\n3\tthree\n" {
		t.Fatalf("%q %v", text, err)
	}
	text, _ = Read(path, 2, 1)
	if text != "2\ttwo\n[lines 2-2 of 3]\n" {
		t.Fatalf("%q", text)
	}
	if _, err := Read(path, 9, 0); err == nil {
		t.Fatal("expected offset error")
	}
}

func TestWriteCreatesDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x", "y", "z.txt")
	if err := Write(path, "hi"); err != nil || content(t, path) != "hi" {
		t.Fatal(err)
	}
}

func TestEditRequiresOneMatch(t *testing.T) {
	path := file(t, "a b a")
	if err := Edit(path, "c", "d"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatal(err)
	}
	if err := Edit(path, "a", "d"); err == nil || !strings.Contains(err.Error(), "2 times") {
		t.Fatal(err)
	}
	if err := Edit(path, "b", "c"); err != nil || content(t, path) != "a c a" {
		t.Fatal(err, content(t, path))
	}
}

func TestPatchAppliesHunks(t *testing.T) {
	path := file(t, "a\nb\nc\nd\ne\n")
	diff := "--- a.txt\n+++ a.txt\n@@ -1,3 +1,3 @@\n a\n-b\n+B\n c\n@@ -4,2 +4,3 @@\n d\n+x\n e\n"
	if err := Patch(path, diff); err != nil || content(t, path) != "a\nB\nc\nd\nx\ne\n" {
		t.Fatal(err, content(t, path))
	}
}

func TestPatchFailsWithoutPartialWrite(t *testing.T) {
	path := file(t, "a\nb\nc\n")
	diff := "@@ -1,2 +1,2 @@\n a\n-b\n+B\n@@ -3,1 +3,1 @@\n-missing\n+x\n"
	if err := Patch(path, diff); err == nil || content(t, path) != "a\nb\nc\n" {
		t.Fatal(err, content(t, path))
	}
}

func TestGrepSkipsIgnoredDirectories(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "node_modules"), 0755)
	os.WriteFile(filepath.Join(dir, "a.go"), []byte("x\nneedle here\n"), 0600)
	os.WriteFile(filepath.Join(dir, "b.md"), []byte("needle"), 0600)
	os.WriteFile(filepath.Join(dir, "node_modules", "c.go"), []byte("needle"), 0600)
	results, err := Grep("needle", dir, "*.go")
	if err != nil || len(results) != 1 || results[0] != filepath.Join(dir, "a.go")+":2:needle here" {
		t.Fatal(results, err)
	}
}

func TestLocateFindsUniqueNestedRelativeFile(t *testing.T) {
	dir := t.TempDir()
	nested := filepath.Join(dir, "internal", "prompts", "default_system.md")
	_ = os.MkdirAll(filepath.Dir(nested), 0755)
	_ = os.WriteFile(nested, []byte("prompt"), 0600)
	prev, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer os.Chdir(prev)
	got, err := Locate("default_system.md")
	if err != nil || got != nested {
		t.Fatalf("got=%q err=%v", got, err)
	}
}
