package app

import "testing"

func TestExportName(t *testing.T) {
	for title, want := range map[string]string{"Plan: a/b": "Plan ab.md", "   ": "session.md", "Привет мир": "Привет мир.md", "../../x": "x.md"} {
		if got := exportName(title); got != want {
			t.Fatalf("%q: %q != %q", title, got, want)
		}
	}
}
