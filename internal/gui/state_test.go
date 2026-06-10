package gui

import (
	"testing"

	"typemerge/internal/typst"
)

func TestStateOptions(t *testing.T) {
	state := NewState()
	state.Template.SetText("invoice.typ")
	state.CSV.SetText("data.csv")
	state.Metadata.SetText("meta.txt")
	state.Output.SetText("generated")
	state.PDF.SetChecked(false)
	state.SVG.SetChecked(true)

	options := state.Options()
	if options.TemplatePath != "invoice.typ" || options.OutputDir != "generated" {
		t.Fatalf("unexpected options: %#v", options)
	}
	if len(options.Formats) != 1 || options.Formats[0] != typst.SVG {
		t.Fatalf("unexpected formats: %#v", options.Formats)
	}
}
