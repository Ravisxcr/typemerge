package output

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"typemerge/internal/input"
)

var unsafeName = regexp.MustCompile(`[^a-z0-9]+`)

func Stem(record input.Record, rowNumber int) string {
	keys := make([]string, 0, len(record))
	for key := range record {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if value := strings.TrimSpace(record[key]); value != "" {
			return Slug(value)
		}
	}
	return fmt.Sprintf("record-%03d", rowNumber)
}

func Slug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = unsafeName.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if value == "" {
		return "record"
	}
	return value
}
