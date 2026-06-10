package main

import (
	"typemerge/internal/app"
	"typemerge/internal/gui"
)

func main() {
	gui.Run(app.NewService())
}
