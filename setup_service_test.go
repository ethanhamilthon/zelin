package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"zelin/internal/provider"
	"zelin/internal/sessions"
)

func provide(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if strings.HasSuffix(r.URL.Path, "/models") {
			w.Write([]byte(`{"data":[{"id":"main"},{"id":"title"},{"id":"broken"}]}`))
			return
		}
		var body strings.Builder
		buf := make([]byte, 4096)
		n, _ := r.Body.Read(buf)
		body.Write(buf[:n])
		if strings.Contains(body.String(), `"broken"`) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"OK"}}]}`))
	}))
}

func TestSetupFlow(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ZASH_HOME", dir)
	server := provide(t)
	defer server.Close()
	store, err := sessions.Open(filepath.Join(dir, "sessions.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	client := &provider.Client{HTTP: server.Client(), Effort: provider.DefaultEffort}
	setup := &SetupService{client: client, store: store, agent: &AgentService{client: client, store: store}}
	ctx := context.Background()
	if required, _ := setup.Required(); !required {
		t.Fatal("setup should be required on a fresh store")
	}
	if _, err := setup.Connect(ctx, server.URL, "wrong"); err == nil {
		t.Fatal("expected an auth error")
	}
	models, err := setup.Connect(ctx, server.URL+"/", "key")
	if err != nil || len(models) != 3 {
		t.Fatal(models, err)
	}
	if env, _ := os.ReadFile(filepath.Join(dir, ".env")); !strings.Contains(string(env), `OPENAI_API_KEY="key"`) {
		t.Fatalf("%s", env)
	}
	if check := setup.Check(ctx, "main"); check.Error != "" || check.Reply != "OK" {
		t.Fatalf("%+v", check)
	}
	if check := setup.Check(ctx, "broken"); check.Error == "" {
		t.Fatalf("%+v", check)
	}
	info, err := setup.Finish(ctx, "main", "title")
	if err != nil || info.Model != "main" || info.TitleModel != "title" {
		t.Fatal(info, err)
	}
	if required, _ := setup.Required(); required {
		t.Fatal("setup should be done")
	}
}
