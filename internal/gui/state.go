package gui

import (
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

func NewState() State {
	defaults := app.DefaultOptions()
	return State{
		TemplatePath: defaults.TemplatePath,
		CSVPath:      defaults.CSVPath,
		MetadataPath: defaults.MetadataPath,
		OutputDir:    defaults.OutputDir,
		TypstBinary:  defaults.TypstBinary,
		PDF:          true,
	}
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
