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
		TemplatePath: "examples/salary-slip/template.typ",
		CSVPath:      "examples/salary-slip/data.csv",
		MetadataPath: "examples/salary-slip/metadata.txt",
		OutputDir:    "out",
		Formats:      []typst.Format{typst.PDF},
		TypstBinary:  "typst",
	}
}
