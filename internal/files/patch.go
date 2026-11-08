package files

import (
	"fmt"
	"os"
	"strings"
)

type hunk struct {
	start    int
	old, new []string
}

func Patch(path, diff string) error {
	path, err := Resolve(path)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	hunks, err := parse(diff)
	if err != nil {
		return err
	}
	text := string(data)
	trailing := strings.HasSuffix(text, "\n")
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	shift := 0
	for i, h := range hunks {
		at := locate(lines, h, h.start-1+shift)
		if at < 0 {
			return fmt.Errorf("Hunk %d does not match %s", i+1, path)
		}
		lines = append(lines[:at], append(append([]string{}, h.new...), lines[at+len(h.old):]...)...)
		shift = at - (h.start - 1) + len(h.new) - len(h.old)
	}
	result := strings.Join(lines, "\n")
	if trailing {
		result += "\n"
	}
	return replace(path, result)
}

func parse(diff string) ([]hunk, error) {
	var hunks []hunk
	for _, line := range strings.Split(strings.TrimSuffix(diff, "\n"), "\n") {
		switch {
		case strings.HasPrefix(line, "@@"):
			var start int
			fmt.Sscanf(strings.TrimPrefix(line, "@@ -"), "%d", &start)
			hunks = append(hunks, hunk{start: max(start, 1)})
		case len(hunks) == 0 || strings.HasPrefix(line, "\\"):
		case strings.HasPrefix(line, "+"):
			hunks[len(hunks)-1].new = append(hunks[len(hunks)-1].new, line[1:])
		case strings.HasPrefix(line, "-"):
			hunks[len(hunks)-1].old = append(hunks[len(hunks)-1].old, line[1:])
		default:
			context := strings.TrimPrefix(line, " ")
			hunks[len(hunks)-1].old = append(hunks[len(hunks)-1].old, context)
			hunks[len(hunks)-1].new = append(hunks[len(hunks)-1].new, context)
		}
	}
	if len(hunks) == 0 {
		return nil, fmt.Errorf("Patch has no hunks; use unified diff format with @@ headers")
	}
	return hunks, nil
}

func locate(lines []string, h hunk, hint int) int {
	if matches(lines, h.old, hint) {
		return hint
	}
	found := -1
	for at := 0; at+len(h.old) <= len(lines); at++ {
		if matches(lines, h.old, at) {
			if found >= 0 {
				return -1
			}
			found = at
		}
	}
	return found
}

func matches(lines, old []string, at int) bool {
	if at < 0 || at+len(old) > len(lines) {
		return false
	}
	for i, line := range old {
		if lines[at+i] != line {
			return false
		}
	}
	return true
}
