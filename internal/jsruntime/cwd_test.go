package jsruntime

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func run(r *Runtime, code string) string { return r.Run(context.Background(), code) }

func quoted(path string) string { return "\"" + path + "\"" }

func TestCdAffectsEveryOperation(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "with space")
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "a.txt"), []byte("alpha\n"), 0644)
	os.WriteFile(filepath.Join(dir, "pic.png"), []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"), 0644)
	r := NewIn(root)
	if got := run(r, `cd("with space").path`); got != quoted(dir) {
		t.Fatal(got)
	}
	if got := run(r, `fs.read("a.txt")`); !strings.Contains(got, "alpha") {
		t.Fatal(got)
	}
	if got := run(r, `fs.write("sub/b.txt", "beta").ok`); got != "true" {
		t.Fatal(got)
	}
	if _, err := os.Stat(filepath.Join(dir, "sub", "b.txt")); err != nil {
		t.Fatal(err)
	}
	if got := run(r, `fs.edit("a.txt", "alpha", "gamma").ok`); got != "true" {
		t.Fatal(got)
	}
	if got := run(r, `fs.grep("gamma", ".")`); !strings.Contains(got, "a.txt:1:gamma") {
		t.Fatal(got)
	}
	want := quoted(dir + "\\n" + dir + "\\n")
	if got := run(r, `bash("pwd; echo $PWD", 5).stdout`); got != want {
		t.Fatal(got, want)
	}
	if got := run(r, `view_image("pic.png").ok`); got != "true" {
		t.Fatal(got)
	}
	if got := run(r, `cd("sub").path`); got != quoted(filepath.Join(dir, "sub")) {
		t.Fatal(got)
	}
	if got := run(r, `cd("..").path`); got != quoted(dir) {
		t.Fatal(got)
	}
	if got := run(r, `cd().path`); got != quoted(dir) {
		t.Fatal(got)
	}
}

func TestCdValidationAndShellIsolation(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0644)
	r := NewIn(root)
	for _, code := range []string{`cd("missing").error`, `cd("file").error`} {
		if got := run(r, code); got == "undefined" || got == "null" {
			t.Fatal(code, got)
		}
	}
	if r.Cwd() != root {
		t.Fatal(r.Cwd())
	}
	run(r, `bash("cd /tmp", 5)`)
	if got := run(r, `bash("pwd", 5).stdout`); got != quoted(root+"\\n") {
		t.Fatal(got)
	}
	home, _ := os.UserHomeDir()
	if got := run(r, `cd("~").path`); got != quoted(home) {
		t.Fatal(got)
	}
	if got := run(r, "cd("+quoted(root)+").path"); got != quoted(root) {
		t.Fatal(got)
	}
}

func TestRuntimesAreIsolated(t *testing.T) {
	a, b := NewIn(t.TempDir()), NewIn(t.TempDir())
	run(a, `cd("/tmp")`)
	if b.Cwd() == "/tmp" || a.Cwd() != "/tmp" {
		t.Fatal(a.Cwd(), b.Cwd())
	}
}

func TestDeletedDirectoryReportsError(t *testing.T) {
	root := t.TempDir()
	gone := filepath.Join(root, "gone")
	os.Mkdir(gone, 0755)
	r := NewIn(gone)
	os.Remove(gone)
	if got := run(r, `fs.read("x")`); !strings.Contains(got, "unavailable") {
		t.Fatal(got)
	}
	if got := run(r, `bash("pwd", 5).stderr`); !strings.Contains(got, "unavailable") {
		t.Fatal(got)
	}
	if got := run(r, "cd("+quoted(root)+").ok"); got != "true" {
		t.Fatal(got)
	}
}
