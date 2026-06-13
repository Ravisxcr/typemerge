package typemerge

import "testing"

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()
	if options.TemplatePath != "" || options.CSVPath != "" || options.MetadataPath != "" {
		t.Fatalf("input paths should require explicit configuration: %#v", options)
	}
	if options.OutputDir != "out" {
		t.Fatalf("output directory = %q, want %q", options.OutputDir, "out")
	}
	if len(options.Formats) != 1 || options.Formats[0] != FormatPDF {
		t.Fatalf("unexpected formats: %#v", options.Formats)
	}
}
