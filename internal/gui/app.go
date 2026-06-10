package gui

import (
	fyneapp "fyne.io/fyne/v2/app"

	"typemerge/internal/app"
)

func Run(service *app.Service) {
	application := fyneapp.NewWithID("dev.typemerge.app")
	NewWindow(application, service).ShowAndRun()
}
