package output

import "path/filepath"

type Paths struct {
	Root string
	Typ  string
}

func NewPaths(root string) Paths {
	return Paths{
		Root: root,
		Typ:  filepath.Join(root, "typ"),
	}
}

func (p Paths) FormatDir(format string) string {
	return filepath.Join(p.Root, format)
}
