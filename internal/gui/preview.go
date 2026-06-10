package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	"typemerge/internal/app"
)

type Preview struct {
	label *widget.Label
}

func NewPreview() *Preview {
	label := widget.NewLabel("No files generated yet.")
	label.Wrapping = fyne.TextWrapWord
	return &Preview{label: label}
}

func (p *Preview) Object() fyne.CanvasObject {
	return p.label
}

func (p *Preview) Show(result app.Result) {
	var previewPath string
	for _, file := range result.Files {
		if file.Format == "pdf" || file.Format == "png" || file.Format == "svg" {
			previewPath = file.Path
			break
		}
	}
	if previewPath == "" && len(result.Files) > 0 {
		previewPath = result.Files[0].Path
	}
	p.label.SetText(fmt.Sprintf(
		"Generated %d files.\n\nFirst output:\n%s",
		len(result.Files),
		strings.TrimSpace(previewPath),
	))
}
