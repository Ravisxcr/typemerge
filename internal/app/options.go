package app

import "typemerge/internal/typst"

type Options struct {
	TemplatePath     string
	CSVPath          string
	MetadataPath     string
	OutputDir        string
	Formats          []typst.Format
	TypstBinary      string
	RequiredColumns  []string
	RequiredMetadata []string
}

func DefaultOptions() Options {
	return Options{
		OutputDir:   "out",
		Formats:     []typst.Format{typst.PDF},
		TypstBinary: "typst",
	}
}
