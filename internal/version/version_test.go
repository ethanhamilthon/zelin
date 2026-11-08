package version

import (
	"encoding/json"
	"os"
	"regexp"
	"testing"
)

func TestVersionIsSemanticAndMatchesFrontend(t *testing.T) {
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(Number) {
		t.Fatalf("Not semantic: %s", Number)
	}
	data, err := os.ReadFile("../../frontend/package.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Version != Number || manifest.Name != "zelin" {
		t.Fatalf("Frontend metadata disagrees: %+v vs %s", manifest, Number)
	}
	if got := Get(); got.Name != "Zelin" || got.Version != Number {
		t.Fatalf("%+v", got)
	}
}
