package main

import (
	"embed"
	"log"
	"os"

	"zelin/internal/agent"
	"zelin/internal/home"
	"zelin/internal/instructions"
	"zelin/internal/metadata"
	"zelin/internal/provider"
	"zelin/internal/sessions"
	"zelin/internal/version"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if migration, err := home.MigrateLegacy(sessions.Snapshot); err != nil {
		log.Fatal(err)
	} else if migration.Copied {
		log.Printf("Copied the legacy profile %s to %s; the original was left untouched", migration.From, migration.To)
	}
	if err := home.LoadEnv(); err != nil {
		log.Fatal(err)
	}
	dir, err := home.Dir()
	if err != nil {
		log.Fatal(err)
	}
	path, err := home.Path("sessions.db")
	if err != nil {
		log.Fatal(err)
	}
	store, err := sessions.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	client := provider.FromEnv()
	client.Meta = metadata.New(dir + "/metadata.json")
	if err := restoreSelection(store, client); err != nil {
		log.Fatal(err)
	}
	docs := &instructions.Store{Dir: dir, Documents: map[string]instructions.Document{"system": {File: "SYSTEM.md", Default: agent.DefaultPrompt}}}
	service := &AgentService{client: client, store: store, instructions: docs, agents: map[string]*agent.Agent{}}
	if wd, err := os.Getwd(); err == nil {
		if err := store.EnsureDefault(wd); err != nil {
			log.Fatal(err)
		}
	}
	sessionService := &SessionService{store: store}
	workspaces := &WorkspaceService{store: store}
	app := application.New(application.Options{
		Name:        version.Name,
		Description: "Zen AI agent",
		Services: []application.Service{
			application.NewService(service),
			application.NewService(sessionService),
			application.NewService(workspaces),
			application.NewService(&AboutService{}),
			application.NewService(&SettingsService{store: store, instructions: docs}),
			application.NewService(&FileService{}),
			application.NewService(&SetupService{client: client, store: store, agent: service}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	service.emit = func(event agent.Event) { app.Event.Emit("agent:update", event) }
	service.notify = func(name string, data any) { app.Event.Emit(name, data) }
	sessionService.notify = service.notify
	workspaces.notify = service.notify
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            version.Name,
		Width:            1100,
		Height:           760,
		MinWidth:         360,
		MinHeight:        240,
		BackgroundColour: application.NewRGB(17, 17, 17),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
