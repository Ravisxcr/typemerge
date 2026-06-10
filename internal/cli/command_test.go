package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunGeneratesTypOnly(t *testing.T) {
	root := t.TempDir()
	templatePath := filepath.Join(root, "template.typ")
	csvPath := filepath.Join(root, "data.csv")
	metadataPath := filepath.Join(root, "metadata.txt")
	outputPath := filepath.Join(root, "out")

	writeFile(t, templatePath, `Hello {{ get .Employee "name" }}`)
	writeFile(t, csvPath, "id,name\n1,Ada\n")
	writeFile(t, metadataPath, "company=Acme\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err := Run(context.Background(), []string{
		"-template", templatePath,
		"-csv", csvPath,
		"-metadata", metadataPath,
		"-out", outputPath,
		"-formats", "",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Run() error = %v, stderr = %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), filepath.Join("typ", "1.typ")) {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
