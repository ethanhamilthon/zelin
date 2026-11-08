package agent

import "testing"

func TestPartialCode(t *testing.T) {
	cases := map[string]string{
		``:                    "",
		`{"co`:                "",
		`{"code":"a`:          "a",
		`{"code":"a\nb`:       "a\nb",
		`{"code":"x\`:         "x",
		`{"code":"x\u00`:      "x",
		`{"code":"say \"hi\"`: `say "hi"`,
		`{"code":"done"}`:     "done",
		`{"code":"q\"`:        `q"`,
	}
	for input, want := range cases {
		if got := partialCode(input); got != want {
			t.Errorf("%q: got %q want %q", input, got, want)
		}
	}
}
