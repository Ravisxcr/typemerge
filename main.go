package main

import (
	"log"

	"typemerge/internal/app"
	"typemerge/internal/gui"
)

func main() {
	if err := gui.Run(app.NewService()); err != nil {
		log.Fatal(err)
	}
}
