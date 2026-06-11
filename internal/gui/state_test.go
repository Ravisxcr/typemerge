package gui

import (
	"os"
	"path/filepath"
	"testing"

	"typemerge/internal/typst"
)

func TestNewStateUsesUserOutputDirectory(t *testing.T) {
	state := NewState()

	if state.TemplatePath != "" || state.CSVPath != "" || state.MetadataPath != "" {
		t.Fatalf("new GUI state should not contain repository sample paths: %#v", state)
	}
	if state.OutputDir == "" || filepath.Base(state.OutputDir) != "output" {
		t.Fatalf("unexpected default output directory: %q", state.OutputDir)
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		relative, err := filepath.Rel(home, state.OutputDir)
		if err != nil || relative == ".." || filepath.IsAbs(relative) {
			t.Fatalf("default output directory %q is not below home %q", state.OutputDir, home)
		}
	}
	if !state.PDF {
		t.Fatal("PDF should be selected by default")
	}
}

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

func TestDiscoverProjectFindsConventionalFiles(t *testing.T) {
	directory := t.TempDir()
	for _, name := range []string{"template.typ", "data.csv", "metadata.txt"} {
		if err := os.WriteFile(filepath.Join(directory, name), []byte("test"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	setup, err := discoverProject(directory)
	if err != nil {
		t.Fatal(err)
	}
	if setup.Directory != directory || setup.Name != filepath.Base(directory) {
		t.Fatalf("unexpected project identity: %#v", setup)
	}
	if filepath.Base(setup.TemplatePath) != "template.typ" {
		t.Fatalf("unexpected template: %q", setup.TemplatePath)
	}
	if filepath.Base(setup.CSVPath) != "data.csv" {
		t.Fatalf("unexpected CSV: %q", setup.CSVPath)
	}
	if filepath.Base(setup.MetadataPath) != "metadata.txt" {
		t.Fatalf("unexpected metadata: %q", setup.MetadataPath)
	}
	if setup.OutputDir != filepath.Join(directory, "output") {
		t.Fatalf("unexpected output directory: %q", setup.OutputDir)
	}
}

func TestUseTemplatePresetWritesProjectTemplate(t *testing.T) {
	app := NewApp(nil)
	path, err := app.UseTemplatePreset(t.TempDir(), "invoice")
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(contents) == 0 || filepath.Base(path) != "invoice.typ" {
		t.Fatalf("unexpected preset template %q", path)
	}
}
