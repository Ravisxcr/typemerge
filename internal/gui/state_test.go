package gui

import (
	"testing"

	"typemerge/internal/typst"
)

func TestStateOptions(t *testing.T) {
	state := NewState()
	state.TemplatePath = "invoice.typ"
	state.CSVPath = "data.csv"
	state.MetadataPath = "meta.txt"
	state.OutputDir = "generated"
	state.PDF = false
	state.SVG = true

	options := state.Options()
	if options.TemplatePath != "invoice.typ" || options.OutputDir != "generated" {
		t.Fatalf("unexpected options: %#v", options)
	}
	if len(options.Formats) != 1 || options.Formats[0] != typst.SVG {
		t.Fatalf("unexpected formats: %#v", options.Formats)
	}
}
