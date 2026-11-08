package app

import (
	"strings"
	"testing"

	"zelin/internal/version"
)

func TestAboutReportsOnlySafeFacts(t *testing.T) {
	t.Setenv("ZELIN_HOME", "/data/zelin-test")
	t.Setenv("OPENAI_API_KEY", "sk-secret-value")
	info := (&AboutService{}).Info()
	if info.Name != "Zelin" || info.Version != version.Number || info.Profile != "/data/zelin-test" || info.Platform == "" || info.Runtime == "" {
		t.Fatalf("%+v", info)
	}
	if strings.Contains(info.Description+info.Daemon+info.Profile+info.Commit, "sk-secret-value") {
		t.Fatal("Leaked a secret")
	}
}
