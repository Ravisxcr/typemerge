package input

import (
	"strings"
	"testing"
)

func TestParseMetadata(t *testing.T) {
	meta, err := ParseMetadata(strings.NewReader("# payroll\nmonth = May\ncompany=Acme=Labs\n"))
	if err != nil {
		t.Fatal(err)
	}
	if meta["month"] != "May" || meta["company"] != "Acme=Labs" {
		t.Fatalf("unexpected metadata: %#v", meta)
	}
}

func TestParseMetadataRejectsDuplicateKeys(t *testing.T) {
	_, err := ParseMetadata(strings.NewReader("month=May\nmonth=June\n"))
	if err == nil {
		t.Fatal("expected duplicate key error")
	}
}
