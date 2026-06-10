package input

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadMetadata(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open metadata: %w", err)
	}
	defer file.Close()

	return ParseMetadata(file)
}

func ParseMetadata(source io.Reader) (map[string]string, error) {
	meta := make(map[string]string)
	scanner := bufio.NewScanner(source)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("metadata line %d must be key=value", lineNumber)
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("metadata line %d has an empty key", lineNumber)
		}
		if _, exists := meta[key]; exists {
			return nil, fmt.Errorf("metadata key %q is duplicated", key)
		}
		meta[key] = strings.TrimSpace(value)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read metadata: %w", err)
	}
	return meta, nil
}
