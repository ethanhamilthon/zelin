package main

import "github.com/wailsapp/wails/v3/pkg/application"

func (w *WorkspaceService) ChooseDirectory() (string, error) {
	return application.Get().Dialog.OpenFile().CanChooseDirectories(true).CanChooseFiles(false).PromptForSingleSelection()
}
