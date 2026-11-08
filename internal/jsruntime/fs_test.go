package jsruntime

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileFunctions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "note.txt")
	r := New()
	code := `var p = ` + "\"" + path + "\"" + `;
fs.write(p, "alpha\nbeta\n");
fs.edit(p, "beta", "gamma");
fs.patch(p, "@@ -1,1 +1,1 @@\n-alpha\n+delta\n");
[fs.read(p), fs.read(p, {offset: 2}), fs.grep("gam", p).length, fs.edit(p, "nope", "x").error]`
	result := r.Run(context.Background(), code)
	for _, want := range []string{`1\tdelta\n2\tgamma`, `[lines 2-2 of 2]`, `,1,`, `not found`} {
		if !strings.Contains(result, want) {
			t.Fatalf("missing %q in %s", want, result)
		}
	}
}
