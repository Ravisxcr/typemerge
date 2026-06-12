package output

import (
	"testing"

	"typemerge/internal/input"
)

func TestStemPrefersEmployeeID(t *testing.T) {
	record := input.Record{"employee_id": "EMP 001", "name": "Ada Lovelace"}
	if got := Stem(record, 1); got != "emp-001" {
		t.Fatalf("Stem() = %q", got)
	}
}

func TestStemFallsBackToRowNumber(t *testing.T) {
	if got := Stem(input.Record{}, 7); got != "record-007" {
		t.Fatalf("Stem() = %q", got)
	}
}
