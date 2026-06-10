package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func BuildForm(window fyne.Window, state *State, generate *widget.Button) fyne.CanvasObject {
	return container.NewVBox(
		fileRow("Template", state.Template, window),
		fileRow("CSV data", state.CSV, window),
		fileRow("Metadata", state.Metadata, window),
		folderRow("Output", state.Output, window),
		widget.NewForm(
			widget.NewFormItem("Typst binary", state.Typst),
			widget.NewFormItem("Formats", container.NewHBox(state.PDF, state.PNG, state.SVG)),
		),
		generate,
	)
}

func fileRow(label string, entry *widget.Entry, window fyne.Window) fyne.CanvasObject {
	browse := widget.NewButton("Browse", func() {
		picker := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			entry.SetText(reader.URI().Path())
		}, window)
		picker.Show()
	})
	return container.NewBorder(nil, nil, widget.NewLabel(label), browse, entry)
}

func folderRow(label string, entry *widget.Entry, window fyne.Window) fyne.CanvasObject {
	browse := widget.NewButton("Browse", func() {
		picker := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			entry.SetText(uri.Path())
		}, window)
		picker.Show()
	})
	return container.NewBorder(nil, nil, widget.NewLabel(label), browse, entry)
}
