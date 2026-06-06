# receipt-generator

A small Go CLI that renders employee salary-slip `.typ` files from:

- a Typst template
- an employee CSV
- a `metadata.txt` file for month, company name, payment date, etc.

It can also call `typst compile` for every rendered `.typ` file to produce PDFs.

## Usage

Install Typst first if you want PDF compilation:

```bash
typst.exe --version
```

Generate `.typ` files and PDFs:

```bash
go run ./cmd/receipt-generator \
  -template examples/salary-slip-template.typ \
  -csv examples/employees.csv \
  -metadata examples/metadata.txt \
  -out out \
  -typst typst.exe
```

Generate only `.typ` files:

```bash
go run ./cmd/receipt-generator -compile=false
```

Outputs are written to:

- `out/typ/*.typ`
- `out/pdf/*.pdf` when `-compile=true`

## CSV format

The first row must contain column names. Each later row creates one salary slip.

Example:

```csv
employee_id,name,designation,basic,hra,allowances,deductions,net_pay
EMP001,Asha Rao,Software Engineer,70000,28000,12000,5000,105000
```

## metadata.txt format

Use `key=value` lines:

```text
company_name=Acme Software Pvt Ltd
month=May
year=2026
payment_date=2026-05-31
currency=INR
```

Blank lines and lines beginning with `#` are ignored.

## Template placeholders

The template is a normal Typst file with Go template placeholders.

Use employee CSV fields:

```gotemplate
{{ get .Employee "name" }}
{{ get .Employee "net_pay" }}
```

Use metadata fields:

```gotemplate
{{ get .Meta "company_name" }}
{{ get .Meta "month" }}
```
