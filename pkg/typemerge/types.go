package typemerge

type Format string

const (
	FormatPDF Format = "pdf"
	FormatPNG Format = "png"
	FormatSVG Format = "svg"
)

type Options struct {
	TemplatePath     string
	CSVPath          string
	MetadataPath     string
	OutputDir        string
	Formats          []Format
	TypstBinary      string
	RequiredColumns  []string
	RequiredMetadata []string
}

type GeneratedFile struct {
	Record int
	Format string
	Path   string
}

type Result struct {
	Files []GeneratedFile
}
