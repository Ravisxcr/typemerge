# Templates

Templates are Typst documents parsed with Go's `text/template` package.

CSV fields are available through `.Employee`:

```gotemplate
{{ get .Employee "name" }}
```

Shared metadata is available through `.Meta`:

```gotemplate
{{ get .Meta "company_name" }}
```

Available helpers:

- `get MAP KEY`
- `initial VALUE`
- `lower VALUE`
- `upper VALUE`
- `money VALUE`
- `date LAYOUT VALUE`

The bundled example is at `examples/salary-slip/template.typ`.
