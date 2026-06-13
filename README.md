# typemerge

`typemerge` creates personalized [Typst](https://typst.app/) documents from a
CSV file. Each CSV row produces one document by combining:

- a Typst file written as a Go `text/template`
- row-specific values from CSV
- shared values from a `key=value` metadata file

The rendered `.typ` sources are always kept. Optionally, `typemerge` can call
the Typst CLI to compile them to PDF, PNG, or SVG.

## Requirements

- Go 1.22 or newer
- Typst on `PATH` when generating PDF, PNG, or SVG files

Typst is not required when using `-formats ""` to render `.typ` files only.

## Quick Start

Prepare these three files in your working directory:

```text
template.typ
data.csv
metadata.txt
```

Render one editable `.typ` document per CSV row:

```sh
typemerge \
  -template template.typ \
  -csv data.csv \
  -metadata metadata.txt \
  -out out \
  -formats ""
```

To compile PDFs as well, install Typst and change the final option to
`-formats pdf`.

When running from source, replace `typemerge` with:

```sh
go run ./cmd/typemerge
```

For example:

```sh
go run ./cmd/typemerge \
  -template template.typ \
  -csv data.csv \
  -metadata metadata.txt \
  -out out \
  -formats ""
```

## Build

Build the command-line executable:

```sh
go build -o dist/typemerge ./cmd/typemerge
```

On Windows, use:

```powershell
go build -o dist/typemerge.exe ./cmd/typemerge
```

Example CSV:

```csv
id,title,amount
record-001,First document,1250
record-002,Second document,2400
```

Example metadata:

```text
document_type=Statement
generated_on=2026-06-12
```

Inside the template, read CSV and metadata values with:

```gotemplate
{{ get .Record "title" }}
{{ get .Meta "document_type" }}
```

## Output

For an output root of `out`, files are grouped by format:

```text
out/
  typ/
    record-001.typ
    record-002.typ
  pdf/
    record-001.pdf
    record-002.pdf
```

Output names prefer the CSV columns `employee_id`, `Employee ID`, `id`, `ID`,
`name`, or `Name`, in that order. See [CLI usage](docs/usage.md#output-files)
for the full naming behavior.

## Documentation

- [CLI usage](docs/usage.md): flags, validation, output, and troubleshooting

## Development

Run the test suite:

```sh
go test ./...
```

The command-line application delegates to `internal/app.Service`. Parsing,
rendering, output naming, filesystem writing, and Typst compilation are kept in
separate packages. The Go entrypoint intended for programmatic use is
`pkg/typemerge`.
