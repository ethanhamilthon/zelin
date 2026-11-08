package tagged

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkspaceOwnership(t *testing.T) {
	root := setup(t)
	r := ref(root, "own001", "a.txt", "file")
	r.Workspace = "other"
	if _, err := Expand(root, "mine", "[@own001]", nil, []Ref{r}); err == nil {
		t.Fatal("Accepted a reference from another workspace")
	}
	r.Workspace = "mine"
	if _, err := Expand(root, "mine", "[@own001]", nil, []Ref{r}); err != nil {
		t.Fatal(err)
	}
}

func TestUnusedAndLiteralTokensAreIgnored(t *testing.T) {
	root := setup(t)
	res, err := Expand(root, "", "email me@example.com and [@nope99]", nil, []Ref{ref(root, "id9999", "a.txt", "file")})
	if err != nil || res.Content != "email me@example.com and [@nope99]" {
		t.Fatalf("%v %s", err, res.Content)
	}
}

func TestImagesUseTheProviderAttachmentPath(t *testing.T) {
	root := t.TempDir()
	png := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")
	os.WriteFile(filepath.Join(root, "p.png"), png, 0644)
	file := ref(root, "img001", "p.png", "image")
	pasted := Ref{ID: "img002", Kind: "image", Label: "pasted", Data: "data:image/png;base64,QQ=="}
	res, err := Expand(root, "", "see [@img001] and [@img002]", []string{"data:image/png;base64,OLD"}, []Ref{file, pasted})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Images) != 3 || !strings.HasPrefix(res.Images[1], "data:image/png;base64,") || res.Images[2] != pasted.Data {
		t.Fatalf("%v", res.Images)
	}
	if strings.Contains(res.Content, "tagged-file") || strings.Contains(res.Content, "base64") || !strings.Contains(res.Content, "[Image #2]") || !strings.Contains(res.Content, "[Image #3]") {
		t.Fatalf("%s", res.Content)
	}
}

func TestUnreadableImageIsReported(t *testing.T) {
	root := setup(t)
	if _, err := Expand(root, "", "[@img001]", nil, []Ref{ref(root, "img001", "a.txt", "image")}); err == nil {
		t.Fatal("Accepted a non-image as an image")
	}
}
