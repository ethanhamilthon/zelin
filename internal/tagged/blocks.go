package tagged

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	FileLimit  = 128 << 10
	TotalLimit = 384 << 10
	listLimit  = 200
)

func attr(text string) string {
	return strings.NewReplacer("&", "&amp;", "\"", "&quot;", "<", "&lt;", ">", "&gt;").Replace(text)
}

func cdata(text string) string {
	return "<![CDATA[" + strings.ReplaceAll(text, "]]>", "]]]]><![CDATA[>") + "]]>"
}

func where(path, real string) string {
	out := `path="` + attr(path) + `"`
	if real != "" && real != path {
		out += ` target="` + attr(real) + `"`
	}
	return out
}

func isText(data []byte) bool {
	probe := data[:min(len(data), 8000)]
	return !strings.ContainsRune(string(probe), 0) && utf8.Valid(data)
}

func fileBlock(ref Ref, real string, used *int) (string, error) {
	info, err := os.Stat(real)
	if err != nil {
		return "", fmt.Errorf("Tagged file %s is not readable: %w", ref.Path, err)
	}
	if info.IsDir() {
		return "", fmt.Errorf("Tagged file %s is a directory", ref.Path)
	}
	if info.Size() > FileLimit {
		return "", fmt.Errorf("Tagged file %s is too large (%d KB; limit %d KB)", ref.Path, info.Size()>>10, FileLimit>>10)
	}
	data, err := os.ReadFile(real)
	if err != nil {
		return "", fmt.Errorf("Tagged file %s is not readable: %w", ref.Path, err)
	}
	if !isText(data) {
		return fmt.Sprintf(`<tagged-file %s type="binary" size="%d"></tagged-file>`, where(ref.Path, real), len(data)), nil
	}
	if *used += len(data); *used > TotalLimit {
		return "", fmt.Errorf("Tagged files exceed the combined limit of %d KB; remove some references", TotalLimit>>10)
	}
	return `<tagged-file ` + where(ref.Path, real) + `>` + "\n" + cdata(string(data)) + "\n</tagged-file>", nil
}

func folderBlock(ref Ref, real string) (string, error) {
	items, err := os.ReadDir(real)
	if err != nil {
		return "", fmt.Errorf("Tagged folder %s is not readable: %w", ref.Path, err)
	}
	var list strings.Builder
	for i, item := range items {
		if i == listLimit {
			fmt.Fprintf(&list, "... %d more\n", len(items)-listLimit)
			break
		}
		list.WriteString(item.Name())
		if item.IsDir() {
			list.WriteString("/")
		}
		list.WriteString("\n")
	}
	return `<tagged-folder ` + where(ref.Path, real) + `>` + "\n" + cdata(list.String()) + "</tagged-folder>", nil
}
