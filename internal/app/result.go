package app

type GeneratedFile struct {
	Record int    `json:"record"`
	Format string `json:"format"`
	Path   string `json:"path"`
}

type Result struct {
	Files []GeneratedFile `json:"files"`
}
