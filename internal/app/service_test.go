package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"typemerge/internal/typst"
)

func TestGenerateTypFiles(t *testing.T) {
	root := t.TempDir()
	templatePath := filepath.Join(root, "template.typ")
	csvPath := filepath.Join(root, "data.csv")
	metadataPath := filepath.Join(root, "metadata.txt")

	writeTestFile(t, templatePath, `Hello {{ get .Record "name" }} from {{ get .Meta "company" }}`)
	writeTestFile(t, csvPath, "employee_id,name\nEMP 1,Ada\nEMP 1,Grace\n")
	writeTestFile(t, metadataPath, "company=Acme\n")

	result, err := NewService().Generate(context.Background(), Options{
		TemplatePath: templatePath,
		CSVPath:      csvPath,
		MetadataPath: metadataPath,
		OutputDir:    filepath.Join(root, "out"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Files) != 2 {
		t.Fatalf("generated %d files, want 2", len(result.Files))
	}

	first, err := os.ReadFile(result.Files[0].Path)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(first)); got != "Hello Ada from Acme" {
		t.Fatalf("first output = %q", got)
	}
	if filepath.Base(result.Files[1].Path) != "emp-1-2.typ" {
		t.Fatalf("duplicate output path = %q", result.Files[1].Path)
	}
}

func TestGenerateCompiledFiles(t *testing.T) {
	root := t.TempDir()
	templatePath := filepath.Join(root, "template.typ")
	csvPath := filepath.Join(root, "data.csv")
	metadataPath := filepath.Join(root, "metadata.txt")

	writeTestFile(t, templatePath, `Hello {{ get .Record "name" }}`)
	writeTestFile(t, csvPath, "id,name\n1,Ada\n")
	writeTestFile(t, metadataPath, "company=Acme\n")

	service := NewService()
	service.NewCompiler = func(binary string) (typst.Compiler, error) {
		if binary != "fake-typst" {
			t.Fatalf("binary = %q", binary)
		}
		return fakeCompiler{}, nil
	}
	result, err := service.Generate(context.Background(), Options{
		TemplatePath: templatePath,
		CSVPath:      csvPath,
		MetadataPath: metadataPath,
		OutputDir:    filepath.Join(root, "out"),
		Formats:      []typst.Format{typst.PDF},
		TypstBinary:  "fake-typst",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Files) != 2 || result.Files[1].Format != "pdf" {
		t.Fatalf("unexpected result: %#v", result)
	}
	compiled, err := os.ReadFile(result.Files[1].Path)
	if err != nil {
		t.Fatal(err)
	}
	if string(compiled) != "compiled" {
		t.Fatalf("compiled output = %q", compiled)
	}
}

func TestGenerateWritesAllTypFilesBeforeFindingCompiler(t *testing.T) {
	root := t.TempDir()
	templatePath := filepath.Join(root, "template.typ")
	csvPath := filepath.Join(root, "data.csv")
	metadataPath := filepath.Join(root, "metadata.txt")
	outputPath := filepath.Join(root, "out")

	writeTestFile(t, templatePath, `Hello {{ get .Record "name" }}`)
	writeTestFile(t, csvPath, "id,name\n1,Ada\n2,Grace\n")
	writeTestFile(t, metadataPath, "company=Acme\n")

	service := NewService()
	service.NewCompiler = func(string) (typst.Compiler, error) {
		files, err := filepath.Glob(filepath.Join(outputPath, "typ", "*.typ"))
		if err != nil {
			t.Fatal(err)
		}
		if len(files) != 2 {
			t.Fatalf("found %d typ files before compiler lookup, want 2", len(files))
		}
		return nil, errors.New("compiler unavailable")
	}

	result, err := service.Generate(context.Background(), Options{
		TemplatePath: templatePath,
		CSVPath:      csvPath,
		MetadataPath: metadataPath,
		OutputDir:    outputPath,
		Formats:      []typst.Format{typst.PDF},
		TypstBinary:  "missing-typst",
	})
	if err == nil {
		t.Fatal("Generate() error = nil")
	}
	if len(result.Files) != 2 {
		t.Fatalf("result contains %d files, want 2 typ files", len(result.Files))
	}
	for _, file := range result.Files {
		if file.Format != "typ" {
			t.Fatalf("unexpected generated file: %#v", file)
		}
		if _, err := os.Stat(file.Path); err != nil {
			t.Fatalf("stat generated typ file: %v", err)
		}
	}
}

type fakeCompiler struct{}

func (fakeCompiler) Compile(_ context.Context, _, destinationPath string) error {
	return os.WriteFile(destinationPath, []byte("compiled"), 0o644)
}

func writeTestFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
