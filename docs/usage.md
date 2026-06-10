# CLI usage

Run with the bundled salary-slip example:

```bash
go run ./cmd/typemerge
```

Render without compiling:

```bash
go run ./cmd/typemerge -formats ""
```

Choose multiple compiled formats:

```bash
go run ./cmd/typemerge -formats pdf,png,svg
```

Typemerge always renders every editable `.typ` source into `out/typ` first.
It then compiles those files into the requested formats. If Typst is missing or
compilation fails, the generated `.typ` files remain available for editing and
manual compilation.

Use a Windows Typst executable from WSL:

```bash
go run ./cmd/typemerge \
  -typst /mnt/c/Users/rranj/AppData/Local/typst-x86_64-pc-windows-msvc/typst.exe \
  -formats pdf,png,svg
```

The `-typst` option accepts either an executable name available on `PATH` or
the full path to the executable.

Use custom inputs:

```bash
go run ./cmd/typemerge \
  -template path/to/template.typ \
  -csv path/to/data.csv \
  -metadata path/to/metadata.txt \
  -out out
```

Validate required input fields before rendering:

```bash
go run ./cmd/typemerge \
  -require-columns employee_id,name,net_pay \
  -require-metadata company_name,payment_date
```
