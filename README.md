# typemerge

`typemerge` renders one Typst document per CSV row using a Go template and
shared `key=value` metadata. It can keep the rendered `.typ` files or invoke
Typst to produce PDF, PNG, and SVG output.

## Quick start

Generate the bundled salary slips as `.typ` files:

```bash
go run ./cmd/typemerge -formats ""
```

Generate PDFs (requires `typst` on `PATH`):

```bash
go run ./cmd/typemerge
```

Generated files are written below `out/`.

## Inputs

- `-template`: Typst file containing Go template expressions
- `-csv`: tabular data; each row creates one document
- `-metadata`: shared `key=value` data
- `-formats`: comma-separated `pdf`, `png`, or `svg`; empty means `.typ` only
- `-typst`: Typst executable name or path, including a WSL path to `typst.exe`
- `-out`: output root

See [docs/usage.md](docs/usage.md), [docs/templates.md](docs/templates.md), and
[docs/metadata.md](docs/metadata.md) for details.

## Architecture

The CLI calls `internal/app.Service`. Parsing, rendering, naming, filesystem
output, and Typst compilation are isolated in their own packages. The supported
library entrypoint is `pkg/typemerge`.
