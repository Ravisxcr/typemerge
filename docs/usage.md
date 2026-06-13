# CLI Usage

Pass all input paths explicitly so the command does not depend on repository
sample files.

```sh
go run ./cmd/typemerge \
  -template path/to/template.typ \
  -csv path/to/data.csv \
  -metadata path/to/metadata.txt \
  -out path/to/output \
  -formats pdf
```

PowerShell users can place the command on one line or replace each trailing
`\` with a backtick.

To render editable `.typ` files without installing or invoking Typst:

```sh
go run ./cmd/typemerge \
  -template path/to/template.typ \
  -csv path/to/data.csv \
  -metadata path/to/metadata.txt \
  -out path/to/output \
  -formats ""
```

## Options

| Option | Default | Description |
| --- | --- | --- |
| `-template` | required | Typst file parsed as a Go template |
| `-csv` | required | CSV file; each data row creates one document |
| `-metadata` | required | Shared `key=value` metadata file |
| `-out` | `out` | Root directory for generated files |
| `-formats` | `pdf` | Comma-separated `pdf`, `png`, and/or `svg`; empty means `.typ` only |
| `-typst` | `typst` | Typst executable name or full path |
| `-require-columns` | empty | Comma-separated CSV columns that must exist |
| `-require-metadata` | empty | Comma-separated metadata keys that must exist |

Display the built-in help:

```sh
go run ./cmd/typemerge -h
```

The `-template`, `-csv`, and `-metadata` paths must be provided explicitly.

Format names are case-insensitive, surrounding whitespace is ignored, and
duplicates are removed. Unsupported formats cause an error.

## Multiple Formats

```sh
go run ./cmd/typemerge -formats pdf,png,svg
```

Every document is rendered to `.typ` first and then compiled into each selected
format. If compilation fails, the `.typ` files already written remain available
for inspection or manual compilation.

To render without invoking Typst:

```sh
go run ./cmd/typemerge -formats ""
```

## Input Validation

Require fields used by a template:

```sh
go run ./cmd/typemerge \
  -template template.typ \
  -csv data.csv \
  -metadata metadata.txt \
  -require-columns id,title,amount \
  -require-metadata document_type,generated_on \
  -formats ""
```

Required CSV validation checks the header. Required metadata validation checks
that each key exists. It does not require values to be non-empty.

CSV input follows standard CSV quoting rules. In addition:

- the first row must contain unique, non-empty headers
- whitespace around headers and values is trimmed
- every data row must have the same number of fields as the header
- at least one data row is required


## Output Files

For each CSV row, `typemerge` writes:

```text
<out>/typ/<name>.typ
<out>/pdf/<name>.pdf
<out>/png/<name>.png
<out>/svg/<name>.svg
```

Only requested compiled-format directories are created.

The filename uses the first non-empty value after sorting column names. When
the row has no non-empty values, it uses `record-NNN`.

Names are lowercased and converted to URL-like slugs. For example,
`Record 001` becomes `record-001`. Repeated names receive `-2`, `-3`, and so
on.

## Typst Executable

`-typst` accepts an executable name available on `PATH` or a full path:

```sh
go run ./cmd/typemerge -typst /path/to/typst -formats pdf
```

For example, WSL can call a Windows executable:

```sh
go run ./cmd/typemerge \
  -typst /mnt/c/path/to/typst.exe \
  -formats pdf,png,svg
```

## Programmatic Use

The `pkg/typemerge` package exposes the same generation workflow:

```go
package main

import (
	"context"
	"log"

	"typemerge/pkg/typemerge"
)

func main() {
	options := typemerge.DefaultOptions()
	options.TemplatePath = "template.typ"
	options.CSVPath = "data.csv"
	options.MetadataPath = "metadata.txt"
	options.OutputDir = "out"
	options.Formats = []typemerge.Format{typemerge.FormatPDF}

	result, err := typemerge.Generate(context.Background(), options)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range result.Files {
		log.Printf("wrote %s", file.Path)
	}
}
```

The import path above matches the current module name and is directly usable
inside this module. A published module should replace it with its full
repository module path.

## Troubleshooting

**`find Typst executable "typst"`**

Typst is not installed or is not on `PATH`. Install Typst, pass its full path
with `-typst`, or use `-formats ""`.

**A template value breaks Typst compilation**

Template values are inserted as raw text. Escape or constrain data that may
contain Typst syntax such as `#`, `[`, `]`, or quotes.

**The command cannot find an input file**

Check the current working directory and pass explicit paths with `-template`,
`-csv`, and `-metadata`.

**Some `.typ` files exist but compiled output is incomplete**

Rendering happens before compilation. Fix the Typst error shown by the command,
then rerun generation.
