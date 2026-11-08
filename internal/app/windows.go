package app

import (
	"github.com/wailsapp/wails/v3/pkg/application"
	"zelin/internal/version"
)

func openWindow(app *application.App) {
	app.Window.NewWithOptions(application.WebviewWindowOptions{Title: version.Name, Width: 1100, Height: 760, MinWidth: 360, MinHeight: 240, BackgroundColour: application.NewRGB(17, 17, 17), URL: "/"})
}
