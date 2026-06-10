package render

import (
	"strings"
	"text/template"
	"time"
)

func Funcs() template.FuncMap {
	return template.FuncMap{
		"date":    formatDate,
		"get":     get,
		"initial": initial,
		"lower":   strings.ToLower,
		"money":   money,
		"upper":   strings.ToUpper,
	}
}

func get(values map[string]string, key string) string {
	return values[key]
}

func initial(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.ToUpper(string([]rune(value)[0]))
}

func money(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "0"
	}
	return value
}

func formatDate(layout, value string) string {
	value = strings.TrimSpace(value)
	for _, inputLayout := range []string{"2006-01-02", time.RFC3339, "02/01/2006", "01/02/2006"} {
		parsed, err := time.Parse(inputLayout, value)
		if err == nil {
			return parsed.Format(layout)
		}
	}
	return value
}
