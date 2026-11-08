package agent

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestLimitOutput(t *testing.T) {
	for _, limit := range []int{1024, 2048, 4096, 8192, 16384} {
		text := strings.Repeat("界", 6000)
		got := LimitOutput(text, limit)
		prefix := strings.Split(got, "\n")[0]
		if !utf8.ValidString(got) || len(prefix) > limit || !strings.Contains(got, "output truncated") {
			t.Fatalf("Invalid output for %d", limit)
		}
	}
	text := strings.Repeat("x", 200000)
	if LimitOutput(text, 0) != text {
		t.Fatal("Disable must preserve full output")
	}
	if LimitOutput("small", 1024) != "small" {
		t.Fatal("Small output changed")
	}
}
