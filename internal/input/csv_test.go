package input

import (
	"strings"
	"testing"
)

func TestParseCSV(t *testing.T) {
	records, err := ParseCSV(strings.NewReader("id, name\n1, Ada\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got := records[0]["name"]; got != "Ada" {
		t.Fatalf("name = %q, want Ada", got)
	}
}

func TestParseCSVRejectsDuplicateHeaders(t *testing.T) {
	_, err := ParseCSV(strings.NewReader("id,id\n1,2\n"))
	if err == nil {
		t.Fatal("expected duplicate header error")
	}
}

func TestParseCSVRejectsRowsWithWrongFieldCount(t *testing.T) {
	_, err := ParseCSV(strings.NewReader("id,name\n1\n"))
	if err == nil {
		t.Fatal("expected field count error")
	}
}
