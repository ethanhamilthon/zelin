package main

import (
	"embed"
	"log"

	"zash/internal/agent"
	"zash/internal/home"
	"zash/internal/instructions"
	"zash/internal/metadata"
	"zash/internal/provider"
	"zash/internal/sessions"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
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
	app := application.New(application.Options{
		Name:        home.Name,
		Description: "Zen AI agent",
		Services: []application.Service{
			application.NewService(service),
			application.NewService(&SessionService{store: store}),
			application.NewService(&SettingsService{store: store, instructions: docs}),
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
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            home.Name,
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
