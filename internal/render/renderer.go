package render

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"typemerge/internal/input"
)

type Data struct {
	Record input.Record
	Meta   map[string]string
}

type Renderer struct {
	template *template.Template
}

func Load(path string) (*Renderer, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read template: %w", err)
	}
	return Parse(filepath.Base(path), string(raw))
}

func Parse(name, source string) (*Renderer, error) {
	tmpl, err := template.New(name).
		Funcs(Funcs()).
		Option("missingkey=zero").
		Parse(source)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	return &Renderer{template: tmpl}, nil
}

func (r *Renderer) Render(data Data) ([]byte, error) {
	var output bytes.Buffer
	if err := r.template.Execute(&output, data); err != nil {
		return nil, fmt.Errorf("render template: %w", err)
	}
	return output.Bytes(), nil
}
