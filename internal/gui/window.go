package gui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"typemerge/internal/app"
)

type MainWindow struct {
	window   fyne.Window
	service  *app.Service
	state    *State
	logs     *LogView
	preview  *Preview
	generate *widget.Button
}

func NewWindow(application fyne.App, service *app.Service) fyne.Window {
	window := application.NewWindow("typemerge")
	main := &MainWindow{
		window:  window,
		service: service,
		state:   NewState(),
		logs:    NewLogView(),
		preview: NewPreview(),
	}
	main.generate = widget.NewButton("Generate", main.startGeneration)

	form := BuildForm(window, main.state, main.generate)
	output := container.NewGridWithColumns(2,
		container.NewBorder(widget.NewLabel("Summary"), nil, nil, nil, main.preview.Object()),
		container.NewBorder(widget.NewLabel("Log"), nil, nil, nil, main.logs.Object()),
	)
	window.SetContent(container.NewBorder(nil, nil, nil, nil, container.NewVSplit(form, output)))
	window.Resize(fyne.NewSize(900, 620))
	return window
}

func (w *MainWindow) startGeneration() {
	options := w.state.Options()
	w.generate.Disable()
	w.logs.Reset("Generating...")

	go func() {
		result, err := w.service.Generate(context.Background(), options)
		fyne.Do(func() {
			defer w.generate.Enable()
			if err != nil {
				w.logs.Reset(fmt.Sprintf("Error: %v", err))
				return
			}

			lines := make([]string, len(result.Files))
			for index, file := range result.Files {
				lines[index] = "wrote " + file.Path
			}
			w.logs.ShowLines(lines)
			w.preview.Show(result)
		})
	}()
}
