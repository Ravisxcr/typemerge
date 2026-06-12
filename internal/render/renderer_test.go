package render

import (
	"strings"
	"testing"

	"typemerge/internal/input"
)

func TestRendererUsesTemplateFunctions(t *testing.T) {
	renderer, err := Parse("test", `{{ initial (get .Meta "company") }} {{ upper (get .Employee "name") }} {{ money "" }}`)
	if err != nil {
		t.Fatal(err)
	}

	rendered, err := renderer.Render(Data{
		Employee: input.Record{"name": "Ada"},
		Meta:     map[string]string{"company": "example"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(string(rendered)); got != "E ADA 0" {
		t.Fatalf("rendered = %q", got)
	}
}
