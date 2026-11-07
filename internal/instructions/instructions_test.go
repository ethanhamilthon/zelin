package instructions

import "testing"

func TestReadDefaultsAndWrite(t *testing.T) {
	store := &Store{Dir: t.TempDir(), Documents: map[string]Document{"system": {File: "SYSTEM.md", Default: "Default"}}}
	if text, err := store.Read("system"); err != nil || text != "Default" {
		t.Fatalf("%q %v", text, err)
	}
	if err := store.Write("system", "# Edited"); err != nil {
		t.Fatal(err)
	}
	if text, _ := store.Read("system"); text != "# Edited" {
		t.Fatalf("%q", text)
	}
	if _, err := store.Read("../secret"); err == nil {
		t.Fatal("Read an unknown document")
	}
}
