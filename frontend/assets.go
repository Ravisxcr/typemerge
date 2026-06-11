package frontend

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var files embed.FS

var Assets fs.FS

func init() {
	Assets, _ = fs.Sub(files, "dist")
}
