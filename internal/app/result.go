package app

type GeneratedFile struct {
	Record int
	Format string
	Path   string
}

type Result struct {
	Files []GeneratedFile
}
