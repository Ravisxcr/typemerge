package typst

import (
	"fmt"
	"strings"
)

type Format string

const (
	PDF Format = "pdf"
	PNG Format = "png"
	SVG Format = "svg"
)

func (f Format) Valid() bool {
	switch f {
	case PDF, PNG, SVG:
		return true
	default:
		return false
	}
}

func ParseFormats(value string) ([]Format, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	seen := make(map[Format]bool)
	var formats []Format
	for _, raw := range strings.Split(value, ",") {
		format := Format(strings.ToLower(strings.TrimSpace(raw)))
		if !format.Valid() {
			return nil, fmt.Errorf("unsupported output format %q", raw)
		}
		if !seen[format] {
			seen[format] = true
			formats = append(formats, format)
		}
	}
	return formats, nil
}
