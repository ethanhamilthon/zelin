package home

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvKeepsExistingVariables(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ZASH_HOME", dir)
	t.Setenv("ZASH_TEST_SET", "shell")
	os.Unsetenv("ZASH_TEST_KEY")
	defer os.Unsetenv("ZASH_TEST_KEY")
	os.WriteFile(filepath.Join(dir, ".env"), []byte("# keys\nexport ZASH_TEST_KEY=\"secret\"\nZASH_TEST_SET=file\n"), 0600)
	if err := LoadEnv(); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("ZASH_TEST_KEY") != "secret" || os.Getenv("ZASH_TEST_SET") != "shell" {
		t.Fatalf("%q %q", os.Getenv("ZASH_TEST_KEY"), os.Getenv("ZASH_TEST_SET"))
	}
}
