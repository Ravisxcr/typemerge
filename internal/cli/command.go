package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strings"

	"typemerge/internal/app"
	"typemerge/internal/typst"
	"typemerge/internal/version"
)

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	options := app.DefaultOptions()
	formatsValue := formatsString(options.Formats)
	var requiredColumns string
	var requiredMetadata string
	var showVersion bool

	flags := flag.NewFlagSet("typemerge", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.BoolVar(&showVersion, "version", false, "display version information")
	flags.StringVar(&options.TemplatePath, "template", options.TemplatePath, "Typst template file")
	flags.StringVar(&options.CSVPath, "csv", options.CSVPath, "CSV data file")
	flags.StringVar(&options.MetadataPath, "metadata", options.MetadataPath, "key=value metadata file")
	flags.StringVar(&options.OutputDir, "out", options.OutputDir, "output directory")
	flags.StringVar(&formatsValue, "formats", formatsValue, "compiled formats: pdf,png,svg; empty writes only .typ")
	flags.StringVar(&options.TypstBinary, "typst", options.TypstBinary, "Typst executable name or path")
	flags.StringVar(&requiredColumns, "require-columns", "", "comma-separated required CSV columns")
	flags.StringVar(&requiredMetadata, "require-metadata", "", "comma-separated required metadata keys")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if showVersion {
		fmt.Fprintf(stdout, "typemerge %s\n", version.Value)
		return nil
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(flags.Args(), " "))
	}

	formats, err := typst.ParseFormats(formatsValue)
	if err != nil {
		return err
	}
	options.Formats = formats
	options.RequiredColumns = splitList(requiredColumns)
	options.RequiredMetadata = splitList(requiredMetadata)

	result, err := app.NewService().Generate(ctx, options)
	if err != nil {
		return err
	}
	PrintResult(stdout, result)
	return nil
}

func splitList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func formatsString(formats []typst.Format) string {
	values := make([]string, len(formats))
	for index, format := range formats {
		values[index] = string(format)
	}
	return strings.Join(values, ",")
}
