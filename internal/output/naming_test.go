package output

import (
	"testing"

	"typemerge/internal/input"
)

func TestStemUsesFirstNonEmptyValueBySortedColumnName(t *testing.T) {
	record := input.Record{"id": "42", "account": "Primary", "blank": " "}
	if got := Stem(record, 1); got != "primary" {
		t.Fatalf("Stem() = %q", got)
	}
}

func TestStemFallsBackToRowNumber(t *testing.T) {
	if got := Stem(input.Record{}, 7); got != "record-007" {
		t.Fatalf("Stem() = %q", got)
	}
}
