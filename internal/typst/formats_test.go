package typst

import "testing"

func TestParseFormats(t *testing.T) {
	formats, err := ParseFormats("pdf, png,pdf")
	if err != nil {
		t.Fatal(err)
	}
	if len(formats) != 2 || formats[0] != PDF || formats[1] != PNG {
		t.Fatalf("unexpected formats: %#v", formats)
	}
}

func TestParseFormatsRejectsUnknownFormat(t *testing.T) {
	if _, err := ParseFormats("docx"); err == nil {
		t.Fatal("expected unsupported format error")
	}
}
