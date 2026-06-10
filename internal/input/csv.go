package input

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Record map[string]string

func ReadCSV(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open CSV: %w", err)
	}
	defer file.Close()

	return ParseCSV(file)
}

func ParseCSV(source io.Reader) ([]Record, error) {
	reader := csv.NewReader(source)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read CSV header: %w", err)
	}

	seen := make(map[string]struct{}, len(headers))
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
		if headers[i] == "" {
			return nil, fmt.Errorf("CSV header %d is empty", i+1)
		}
		if _, exists := seen[headers[i]]; exists {
			return nil, fmt.Errorf("CSV header %q is duplicated", headers[i])
		}
		seen[headers[i]] = struct{}{}
	}

	var records []Record
	for line := 2; ; line++ {
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

		record := make(Record, len(headers))
		for i, header := range headers {
			record[header] = strings.TrimSpace(row[i])
		}
		records = append(records, record)
	}

	if len(records) == 0 {
		return nil, errors.New("CSV has no data rows")
	}
	return records, nil
}
