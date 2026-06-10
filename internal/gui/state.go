package gui

import (
	"fyne.io/fyne/v2/widget"

	"typemerge/internal/app"
	"typemerge/internal/typst"
)

type State struct {
	Template *widget.Entry
	CSV      *widget.Entry
	Metadata *widget.Entry
	Output   *widget.Entry
	Typst    *widget.Entry
	PDF      *widget.Check
	PNG      *widget.Check
	SVG      *widget.Check
}

func NewState() *State {
	defaults := app.DefaultOptions()
	state := &State{
		Template: widget.NewEntry(),
		CSV:      widget.NewEntry(),
		Metadata: widget.NewEntry(),
		Output:   widget.NewEntry(),
		Typst:    widget.NewEntry(),
		PDF:      widget.NewCheck("PDF", nil),
		PNG:      widget.NewCheck("PNG", nil),
		SVG:      widget.NewCheck("SVG", nil),
	}
	state.Template.SetText(defaults.TemplatePath)
	state.CSV.SetText(defaults.CSVPath)
	state.Metadata.SetText(defaults.MetadataPath)
	state.Output.SetText(defaults.OutputDir)
	state.Typst.SetText(defaults.TypstBinary)
	state.PDF.SetChecked(true)
	return state
}

func (s *State) Options() app.Options {
	formats := make([]typst.Format, 0, 3)
	if s.PDF.Checked {
		formats = append(formats, typst.PDF)
	}
	if s.PNG.Checked {
		formats = append(formats, typst.PNG)
	}
	if s.SVG.Checked {
		formats = append(formats, typst.SVG)
	}
	return app.Options{
		TemplatePath: s.Template.Text,
		CSVPath:      s.CSV.Text,
		MetadataPath: s.Metadata.Text,
		OutputDir:    s.Output.Text,
		Formats:      formats,
		TypstBinary:  s.Typst.Text,
	}
}
