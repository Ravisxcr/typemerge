package input

import (
	"fmt"
	"strings"
)

func ValidateColumns(records []Record, required []string) error {
	if len(records) == 0 {
		return fmt.Errorf("no records to validate")
	}
	for _, column := range required {
		column = strings.TrimSpace(column)
		if column == "" {
			continue
		}
		if _, exists := records[0][column]; !exists {
			return fmt.Errorf("CSV is missing required column %q", column)
		}
	}
	return nil
}

func ValidateMetadata(meta map[string]string, required []string) error {
	for _, key := range required {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, exists := meta[key]; !exists {
			return fmt.Errorf("metadata is missing required key %q", key)
		}
	}
	return nil
}
