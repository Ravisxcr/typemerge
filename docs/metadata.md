# Metadata

Metadata files contain one `key=value` pair per line:

```text
company_name=Acme Software Pvt Ltd
payment_date=2026-05-31
currency=INR
```

Whitespace around keys and values is trimmed. Blank lines and lines beginning
with `#` are ignored. Values may contain additional `=` characters. Duplicate
keys are rejected.
