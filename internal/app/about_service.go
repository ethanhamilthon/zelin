package app

import (
	"runtime"

	"zelin/internal/home"
	"zelin/internal/version"
)

type AboutService struct{ provider string }

type About struct {
	version.Info
	Description string `json:"description"`
	Runtime     string `json:"runtime"`
	Platform    string `json:"platform"`
	Profile     string `json:"profile"`
	Daemon      string `json:"daemon"`
}

func (a *AboutService) Info() About {
	dir, _ := home.Dir()
	return About{
		Info:        version.Get(),
		Description: "Zen AI agent: a desktop coding agent that works through a JavaScript runtime.",
		Runtime:     runtime.Version(),
		Platform:    runtime.GOOS + "/" + runtime.GOARCH,
		Profile:     dir,
		Daemon:      "Not used: the backend runs inside the application",
	}
}
