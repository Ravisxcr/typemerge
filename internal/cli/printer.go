package cli

import (
	"fmt"
	"io"

	"typemerge/internal/app"
)

func PrintResult(writer io.Writer, result app.Result) {
	for _, file := range result.Files {
		fmt.Fprintf(writer, "wrote %s\n", file.Path)
	}
}
