package home

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveEnvReplacesKeys(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ZASH_HOME", dir)
	os.WriteFile(filepath.Join(dir, ".env"), []byte("# keys\nOTHER=1\nexport OPENAI_API_KEY=old\n"), 0600)
	if err := SaveEnv(map[string]string{"OPENAI_API_KEY": "new", "OPENAI_BASE_URL": "https://x"}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(filepath.Join(dir, ".env"))
	want := "# keys\nOTHER=1\nOPENAI_API_KEY=\"new\"\nOPENAI_BASE_URL=\"https://x\"\n"
	if string(data) != want {
		t.Fatalf("%q", data)
	}
}
