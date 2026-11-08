package main

import (
	"embed"

	"zelin/internal/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() { app.Run(assets) }
