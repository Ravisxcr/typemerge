package typemerge

import (
	"context"

	"typemerge/internal/app"
	internaltypst "typemerge/internal/typst"
)

func DefaultOptions() Options {
	return fromInternalOptions(app.DefaultOptions())
}

func Generate(ctx context.Context, options Options) (Result, error) {
	result, err := app.NewService().Generate(ctx, toInternalOptions(options))
	if err != nil {
		return Result{}, err
	}

	files := make([]GeneratedFile, len(result.Files))
	for index, file := range result.Files {
		files[index] = GeneratedFile(file)
	}
	return Result{Files: files}, nil
}

func toInternalOptions(options Options) app.Options {
	formats := make([]internaltypst.Format, len(options.Formats))
	for index, format := range options.Formats {
		formats[index] = internaltypst.Format(format)
	}
	return app.Options{
		TemplatePath:     options.TemplatePath,
		CSVPath:          options.CSVPath,
		MetadataPath:     options.MetadataPath,
		OutputDir:        options.OutputDir,
		Formats:          formats,
		TypstBinary:      options.TypstBinary,
		RequiredColumns:  options.RequiredColumns,
		RequiredMetadata: options.RequiredMetadata,
	}
}

func fromInternalOptions(options app.Options) Options {
	formats := make([]Format, len(options.Formats))
	for index, format := range options.Formats {
		formats[index] = Format(format)
	}
	return Options{
		TemplatePath:     options.TemplatePath,
		CSVPath:          options.CSVPath,
		MetadataPath:     options.MetadataPath,
		OutputDir:        options.OutputDir,
		Formats:          formats,
		TypstBinary:      options.TypstBinary,
		RequiredColumns:  options.RequiredColumns,
		RequiredMetadata: options.RequiredMetadata,
	}
}
