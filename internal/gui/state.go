package gui

import (
	"os"
	"path/filepath"

	"typemerge/internal/app"
	"typemerge/internal/typst"
)

type State struct {
	TemplatePath string `json:"templatePath"`
	CSVPath      string `json:"csvPath"`
	MetadataPath string `json:"metadataPath"`
	OutputDir    string `json:"outputDir"`
	TypstBinary  string `json:"typstBinary"`
	PDF          bool   `json:"pdf"`
	PNG          bool   `json:"png"`
	SVG          bool   `json:"svg"`
}

type ProjectSetup struct {
	Directory    string `json:"directory"`
	Name         string `json:"name"`
	TemplatePath string `json:"templatePath"`
	CSVPath      string `json:"csvPath"`
	MetadataPath string `json:"metadataPath"`
	OutputDir    string `json:"outputDir"`
}

func NewState() State {
	defaults := app.DefaultOptions()
	return State{
		OutputDir:   defaultOutputDir(),
		TypstBinary: defaults.TypstBinary,
		PDF:         true,
	}
}

func defaultOutputDir() string {
	return filepath.Join(defaultDialogDirectory(), "output")
}

func defaultDialogDirectory() string {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		documents := filepath.Join(home, "Documents")
		if info, statErr := os.Stat(documents); statErr == nil && info.IsDir() {
			return documents
		}
		return home
	}

	workingDir, err := os.Getwd()
	if err == nil && workingDir != "" {
		return workingDir
	}
	return "."
}

func (s State) Options() app.Options {
	formats := make([]typst.Format, 0, 3)
	if s.PDF {
		formats = append(formats, typst.PDF)
	}
	if s.PNG {
		formats = append(formats, typst.PNG)
	}
	if s.SVG {
		formats = append(formats, typst.SVG)
	}
	return app.Options{
		TemplatePath: s.TemplatePath,
		CSVPath:      s.CSVPath,
		MetadataPath: s.MetadataPath,
		OutputDir:    s.OutputDir,
		Formats:      formats,
		TypstBinary:  s.TypstBinary,
	}
}
