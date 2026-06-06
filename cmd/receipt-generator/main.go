package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

type employeeRecord map[string]string

type templateData struct {
	Employee employeeRecord
	Meta     map[string]string
}

type options struct {
	templatePath string
	csvPath      string
	metadataPath string
	outDir       string
	compilePDF   bool
	typstBin     string
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var opts options
	fs := flag.NewFlagSet("receipt-generator", flag.ExitOnError)
	fs.StringVar(&opts.templatePath, "template", "examples/salary-slip-template.typ", "Typst template file")
	fs.StringVar(&opts.csvPath, "csv", "examples/employees.csv", "CSV file containing employee salary data")
	fs.StringVar(&opts.metadataPath, "metadata", "examples/metadata.txt", "metadata.txt file containing key=value lines")
	fs.StringVar(&opts.outDir, "out", "out", "output directory")
	fs.BoolVar(&opts.compilePDF, "compile", true, "compile rendered Typst files to PDF using typst")
	fs.StringVar(&opts.typstBin, "typst", "typst.exe", "typst executable path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	records, err := readCSV(opts.csvPath)
	if err != nil {
		return err
	}
	meta, err := readMetadata(opts.metadataPath)
	if err != nil {
		return err
	}
	tmpl, err := readTemplate(opts.templatePath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(opts.outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	typDir := filepath.Join(opts.outDir, "typ")
	pdfDir := filepath.Join(opts.outDir, "pdf")
	if err := os.MkdirAll(typDir, 0o755); err != nil {
		return fmt.Errorf("create typ output directory: %w", err)
	}
	if opts.compilePDF {
		if err := os.MkdirAll(pdfDir, 0o755); err != nil {
			return fmt.Errorf("create pdf output directory: %w", err)
		}
	}

	for i, record := range records {
		stem := outputStem(record, i+1)
		typPath := filepath.Join(typDir, stem+".typ")
		pdfPath := filepath.Join(pdfDir, stem+".pdf")

		if err := renderTypFile(tmpl, typPath, templateData{Employee: record, Meta: meta}); err != nil {
			return err
		}
		fmt.Println("wrote", typPath)

		if opts.compilePDF {
			if err := compileTypst(opts.typstBin, typPath, pdfPath); err != nil {
				return err
			}
			fmt.Println("wrote", pdfPath)
		}
	}

	return nil
}

func readTemplate(path string) (*template.Template, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}

	funcs := template.FuncMap{
		"get": func(values map[string]string, key string) string {
			return values[key]
		},
		"initial": firstInitial,
		"lower":   strings.ToLower,
		"money":   formatMoney,
	}

	tmpl, err := template.New(filepath.Base(path)).Funcs(funcs).Option("missingkey=zero").Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	return tmpl, nil
}

func readCSV(path string) ([]employeeRecord, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
		if headers[i] == "" {
			return nil, fmt.Errorf("CSV header %d is empty", i+1)
		}
	}

	var records []employeeRecord
	line := 1
	for {
		line++
		row, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read CSV line %d: %w", line, err)
		}
		if len(row) != len(headers) {
			return nil, fmt.Errorf("CSV line %d has %d fields, expected %d", line, len(row), len(headers))
		}

		record := make(employeeRecord, len(headers))
		for i, header := range headers {
			record[header] = strings.TrimSpace(row[i])
		}
		records = append(records, record)
	}
	if len(records) == 0 {
		return nil, errors.New("CSV has no employee rows")
	}
	return records, nil
}

func readMetadata(path string) (map[string]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read metadata: %w", err)
	}

	meta := make(map[string]string)
	lines := strings.Split(string(raw), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("metadata line %d must be key=value", i+1)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("metadata line %d has an empty key", i+1)
		}
		meta[key] = strings.TrimSpace(value)
	}
	return meta, nil
}

func renderTypFile(tmpl *template.Template, path string, data templateData) error {
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("render %s: %w", path, err)
	}
	if err := os.WriteFile(path, rendered.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func compileTypst(typstBin, typPath, pdfPath string) error {
	cmd := exec.Command(typstBin, "compile", typPath, pdfPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compile %s: %w\n%s", typPath, err, strings.TrimSpace(string(output)))
	}
	return nil
}

func outputStem(record employeeRecord, rowNumber int) string {
	for _, key := range []string{"employee_id", "Employee ID", "id", "ID", "name", "Name"} {
		if value := strings.TrimSpace(record[key]); value != "" {
			return slug(value)
		}
	}

	keys := make([]string, 0, len(record))
	for key := range record {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if value := strings.TrimSpace(record[key]); value != "" {
			return slug(value)
		}
	}
	return fmt.Sprintf("employee-%03d", rowNumber)
}

func slug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return "employee"
	}
	return value
}

func formatMoney(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "0"
	}
	return value
}

func firstInitial(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	runes := []rune(value)
	return strings.ToUpper(string(runes[0]))
}
