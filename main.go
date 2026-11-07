package main

import (
	"embed"
	"log"

	"zash/internal/agent"
	"zash/internal/provider"
	"zash/internal/sessions"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	path, err := sessions.DefaultPath()
	if err != nil {
		log.Fatal(err)
	}
	store, err := sessions.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	client := provider.FromEnv()
	service := &AgentService{client: client, store: store, agents: map[string]*agent.Agent{}}
	app := application.New(application.Options{
		Name:        "zash",
		Description: "Zen AI agent",
		Services: []application.Service{
			application.NewService(service),
			application.NewService(&SessionService{store: store}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	service.emit = func(event agent.Event) { app.Event.Emit("agent:update", event) }
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "zash",
		Width:            1100,
		Height:           760,
		MinWidth:         360,
		MinHeight:        240,
		BackgroundColour: application.NewRGB(23, 23, 23),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
