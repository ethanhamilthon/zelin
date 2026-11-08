package tagged

import (
	"fmt"
	"path/filepath"
	"strings"

	"zelin/internal/files"
	"zelin/internal/workspace"
)

type Result struct {
	Content string
	Images  []string
}

func resolve(root string, ref Ref) (string, string, error) {
	if ref.Path == "" {
		return "", "", fmt.Errorf("Reference %s has no path", ref.Label)
	}
	path, err := workspace.Within(root, ref.Path)
	if err != nil {
		return "", "", err
	}
	real, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", "", fmt.Errorf("Tagged item %s is missing", path)
	}
	base, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", "", err
	}
	if !workspace.Inside(base, real) {
		return "", "", fmt.Errorf("%s links outside the workspace to %s", path, real)
	}
	return path, real, nil
}

func image(root string, ref Ref) (string, error) {
	if ref.Data != "" {
		return ref.Data, nil
	}
	_, real, err := resolve(root, ref)
	if err != nil {
		return "", err
	}
	img, err := files.ReadImage(real)
	return img.DataURL, err
}

func Expand(root, workspaceID, content string, images []string, refs []Ref) (Result, error) {
	result := Result{Content: content, Images: append([]string{}, images...)}
	known := map[string]Ref{}
	for _, ref := range Used(content, refs) {
		known[ref.ID] = ref
	}
	var blocks []string
	done, size := map[string]bool{}, 0
	replacement := map[string]string{}
	for _, ref := range known {
		if ref.Workspace != "" && workspaceID != "" && ref.Workspace != workspaceID {
			return result, fmt.Errorf("Reference %s belongs to another workspace", ref.Label)
		}
	}
	var failure error
	result.Content = Token.ReplaceAllStringFunc(content, func(token string) string {
		match := Token.FindStringSubmatch(token)
		ref, ok := known[match[1]]
		if !ok || failure != nil {
			return Plain(token)
		}
		if text, seen := replacement[ref.ID]; seen {
			return text
		}
		text := "@" + ref.Label
		switch ref.Kind {
		case "image":
			url, err := image(root, ref)
			if err != nil {
				failure = fmt.Errorf("Tagged image %s: %w", ref.Label, err)
				return token
			}
			result.Images = append(result.Images, url)
			text = fmt.Sprintf("[Image #%d]", len(result.Images))
		case "file", "dir":
			path, real, err := resolve(root, ref)
			if err != nil {
				failure = err
				return token
			}
			if !done[real] {
				done[real] = true
				var block string
				if ref.Kind == "dir" {
					block, err = folderBlock(Ref{Path: path}, real)
				} else {
					block, err = fileBlock(Ref{Path: path}, real, &size)
				}
				if err != nil {
					failure = err
					return token
				}
				blocks = append(blocks, block)
			}
		default:
			failure = fmt.Errorf("Unknown reference kind %q", ref.Kind)
			return token
		}
		replacement[ref.ID] = text
		return text
	})
	if failure != nil {
		return Result{Content: content, Images: images}, failure
	}
	if len(blocks) > 0 {
		result.Content += "\n\n" + strings.Join(blocks, "\n")
	}
	return result, nil
}
