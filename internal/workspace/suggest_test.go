package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSuggestIsBoundedFilteredAndOrdered(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"alpha.go", "alpine", "beta.go", "gamma-alpha.txt", "alp1", "alp2", "alp3", ".alpha-hidden"} {
		path := filepath.Join(root, name)
		if name == "alpine" {
			os.Mkdir(path, 0755)
		} else {
			os.WriteFile(path, nil, 0644)
		}
	}
	got, err := Suggest(root, "", "alp")
	if err != nil || len(got) != SearchLimit {
		t.Fatalf("%d %v", len(got), err)
	}
	if got[0].Name != "alpine" {
		t.Fatalf("Directories first within a rank: %s", names(got))
	}
	for _, e := range got {
		if strings.HasPrefix(e.Name, ".") {
			t.Fatal("Hidden files must need an explicit dot")
		}
	}
	again, _ := Suggest(root, "", "alp")
	if names(got) != names(again) {
		t.Fatal("Ordering is not stable")
	}
	contains, _ := Suggest(root, "", "lpha")
	if len(contains) != 2 || contains[0].Name != "alpha.go" {
		t.Fatalf("%s", names(contains))
	}
	hidden, _ := Suggest(root, "", ".al")
	if len(hidden) != 1 || hidden[0].Name != ".alpha-hidden" {
		t.Fatalf("%s", names(hidden))
	}
}

func TestSuggestTraversesDirectoriesAndRejectsEscapes(t *testing.T) {
	root := tree(t)
	got, err := Suggest(root, "", "src/")
	if err != nil || names(got) != "deep:dir,main.go:file" {
		t.Fatalf("%s %v", names(got), err)
	}
	got, _ = Suggest(root, "", "src/ma")
	if names(got) != "main.go:file" {
		t.Fatal(names(got))
	}
	got, _ = Suggest(root, "", "my d")
	if names(got) != "my docs:dir" {
		t.Fatalf("Names with spaces: %s", names(got))
	}
	if _, err := Suggest(root, "", "../"); err == nil {
		t.Fatal("Traversed outside the workspace")
	}
	if _, err := Suggest(root, "/etc", ""); err == nil {
		t.Fatal("Accepted a base outside the workspace")
	}
}
