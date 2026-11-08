package home

import (
	"os"
	"path/filepath"
	"testing"
)

func fake(from, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	return os.WriteFile(to, append(data, '!'), 0o600)
}

func legacy(t *testing.T) (root, from string) {
	root = t.TempDir()
	t.Setenv("XDG_DATA_HOME", root)
	t.Setenv("ZELIN_HOME", "")
	t.Setenv("ZASH_HOME", "")
	from = filepath.Join(root, "zash")
	os.MkdirAll(filepath.Join(from, "storage"), 0o700)
	os.WriteFile(filepath.Join(from, "sessions.db"), []byte("db"), 0o600)
	os.WriteFile(filepath.Join(from, "sessions.db-wal"), []byte("wal"), 0o600)
	os.WriteFile(filepath.Join(from, ".env"), []byte("KEY=1"), 0o600)
	os.WriteFile(filepath.Join(from, "SYSTEM.md"), []byte("prompt"), 0o600)
	os.WriteFile(filepath.Join(from, "storage", "x"), []byte("x"), 0o600)
	return root, from
}

func TestMigrationCopiesAndKeepsTheLegacyProfile(t *testing.T) {
	root, from := legacy(t)
	result, err := MigrateLegacy(fake)
	if err != nil || !result.Copied {
		t.Fatalf("%+v %v", result, err)
	}
	to := filepath.Join(root, "zelin")
	for file, want := range map[string]string{"sessions.db": "db!", ".env": "KEY=1", "SYSTEM.md": "prompt", "storage/x": "x"} {
		got, err := os.ReadFile(filepath.Join(to, file))
		if err != nil || string(got) != want {
			t.Fatalf("%s: %q %v", file, got, err)
		}
	}
	if _, err := os.Stat(filepath.Join(to, "sessions.db-wal")); err == nil {
		t.Fatal("Journal files must be rebuilt, not copied")
	}
	if got, _ := os.ReadFile(filepath.Join(from, "sessions.db")); string(got) != "db" {
		t.Fatal("The legacy profile was modified")
	}
	if info, _ := os.Stat(filepath.Join(to, ".env")); info.Mode().Perm() != 0o600 {
		t.Fatalf("Secrets must stay private: %v", info.Mode())
	}
}

func TestMigrationNeverOverwrites(t *testing.T) {
	root, _ := legacy(t)
	to := filepath.Join(root, "zelin")
	os.MkdirAll(to, 0o700)
	os.WriteFile(filepath.Join(to, "sessions.db"), []byte("new"), 0o600)
	if result, err := MigrateLegacy(fake); err != nil || result.Copied {
		t.Fatalf("%+v %v", result, err)
	}
	if got, _ := os.ReadFile(filepath.Join(to, "sessions.db")); string(got) != "new" {
		t.Fatal("Existing data was overwritten")
	}
	os.Remove(filepath.Join(to, "sessions.db"))
	os.WriteFile(filepath.Join(to, "stray"), nil, 0o600)
	if _, err := MigrateLegacy(fake); err == nil {
		t.Fatal("A non-empty destination without a database must be reported")
	}
	if _, err := os.Stat(filepath.Join(to, "stray")); err != nil {
		t.Fatal("A failed migration removed existing files")
	}
}

func TestMigrationSkipsWithoutLegacyOrWithOverrides(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_DATA_HOME", root)
	t.Setenv("ZELIN_HOME", "")
	t.Setenv("ZASH_HOME", "")
	if result, err := MigrateLegacy(fake); err != nil || result.Copied {
		t.Fatalf("%+v %v", result, err)
	}
	legacyRoot, _ := legacy(t)
	t.Setenv("ZELIN_HOME", filepath.Join(legacyRoot, "custom"))
	if result, _ := MigrateLegacy(fake); result.Copied {
		t.Fatal("An explicit home must not trigger migration")
	}
}

func TestHomePrecedenceAndProfileIsolation(t *testing.T) {
	t.Setenv("ZASH_HOME", "/legacy")
	t.Setenv("ZELIN_HOME", "")
	if dir, _ := Dir(); dir != "/legacy" {
		t.Fatal(dir)
	}
	t.Setenv("ZELIN_HOME", "/new")
	if dir, _ := Dir(); dir != "/new" {
		t.Fatal(dir)
	}
	t.Setenv("ZELIN_HOME", "")
	t.Setenv("ZASH_HOME", "")
	t.Setenv("XDG_DATA_HOME", "/data")
	old := Name
	defer func() { Name = old }()
	Name = "zelin-prod"
	if dir, _ := Dir(); dir != "/data/zelin-prod" || legacyName() != "zash-prod" {
		t.Fatalf("%s %s", dir, legacyName())
	}
	Name = "zelin"
	if dir, _ := Dir(); dir != "/data/zelin" || legacyName() != "zash" {
		t.Fatalf("%s %s", dir, legacyName())
	}
}
