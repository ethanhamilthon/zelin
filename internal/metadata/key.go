package metadata

import (
	"regexp"
	"strings"
)

var (
	dated   = regexp.MustCompile(`-(\d{8}|\d{4}-\d{2}-\d{2})$`)
	variant = regexp.MustCompile(`-(none|minimal|low|medium|high|xhigh|max|thinking|latest)$`)
)

func Key(id string) string {
	id = strings.ToLower(id)
	id = id[strings.LastIndex(id, "/")+1:]
	id, _, _ = strings.Cut(id, ":")
	id = strings.ReplaceAll(id, ".", "-")
	id = dated.ReplaceAllString(id, "")
	return variant.ReplaceAllString(id, "")
}
