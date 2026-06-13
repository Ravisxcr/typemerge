package app

import (
	"context"
	"fmt"
	"path/filepath"

	"typemerge/internal/input"
	"typemerge/internal/output"
	"typemerge/internal/render"
	"typemerge/internal/typst"
)

type Service struct {
	Writer      output.Writer
	NewCompiler func(string) (typst.Compiler, error)
}

func NewService() *Service {
	return &Service{
		Writer: output.Writer{},
		NewCompiler: func(binary string) (typst.Compiler, error) {
			path, err := typst.Find(binary)
			if err != nil {
				return nil, err
			}
			return typst.CommandCompiler{Binary: path}, nil
		},
	}
}

func (s *Service) Generate(ctx context.Context, options Options) (Result, error) {
	if err := validateOptions(options); err != nil {
		return Result{}, err
	}

	records, err := input.ReadCSV(options.CSVPath)
	if err != nil {
		return Result{}, err
	}
	if err := input.ValidateColumns(records, options.RequiredColumns); err != nil {
		return Result{}, err
	}

	meta, err := input.ReadMetadata(options.MetadataPath)
	if err != nil {
		return Result{}, err
	}
	if err := input.ValidateMetadata(meta, options.RequiredMetadata); err != nil {
		return Result{}, err
	}

	renderer, err := render.Load(options.TemplatePath)
	if err != nil {
		return Result{}, err
	}

	paths := output.NewPaths(options.OutputDir)
	usedStems := make(map[string]int)
	result := Result{}
	type renderedFile struct {
		record int
		stem   string
		path   string
	}
	rendered := make([]renderedFile, 0, len(records))

	for index, record := range records {
		if err := ctx.Err(); err != nil {
			return result, err
		}

		stem := uniqueStem(output.Stem(record, index+1), usedStems)
		typPath := filepath.Join(paths.Typ, stem+".typ")
		contents, err := renderer.Render(render.Data{Record: record, Meta: meta})
		if err != nil {
			return result, fmt.Errorf("render record %d: %w", index+1, err)
		}
		if err := render.ValidateSource(contents); err != nil {
			return result, fmt.Errorf("validate record %d: %w", index+1, err)
		}
		if err := s.Writer.Write(typPath, contents); err != nil {
			return result, err
		}
		result.Files = append(result.Files, GeneratedFile{Record: index + 1, Format: "typ", Path: typPath})
		rendered = append(rendered, renderedFile{record: index + 1, stem: stem, path: typPath})
	}

	if len(options.Formats) == 0 || len(rendered) == 0 {
		return result, nil
	}

	newCompiler := s.NewCompiler
	if newCompiler == nil {
		newCompiler = NewService().NewCompiler
	}
	compiler, err := newCompiler(options.TypstBinary)
	if err != nil {
		return result, err
	}

	for _, source := range rendered {
		for _, format := range options.Formats {
			destination := filepath.Join(paths.FormatDir(string(format)), source.stem+"."+string(format))
			if err := s.Writer.EnsureDir(filepath.Dir(destination)); err != nil {
				return result, err
			}
			if err := compiler.Compile(ctx, source.path, destination); err != nil {
				return result, err
			}
			result.Files = append(result.Files, GeneratedFile{
				Record: source.record,
				Format: string(format),
				Path:   destination,
			})
		}
	}

	return result, nil
}

func validateOptions(options Options) error {
	switch {
	case options.TemplatePath == "":
		return fmt.Errorf("template path is required")
	case options.CSVPath == "":
		return fmt.Errorf("CSV path is required")
	case options.MetadataPath == "":
		return fmt.Errorf("metadata path is required")
	case options.OutputDir == "":
		return fmt.Errorf("output directory is required")
	case len(options.Formats) > 0 && options.TypstBinary == "":
		return fmt.Errorf("Typst binary is required when compiled formats are enabled")
	}
	for _, format := range options.Formats {
		if !format.Valid() {
			return fmt.Errorf("unsupported output format %q", format)
		}
	}
	return nil
}

func uniqueStem(stem string, used map[string]int) string {
	used[stem]++
	if used[stem] == 1 {
		return stem
	}
	return fmt.Sprintf("%s-%d", stem, used[stem])
}
