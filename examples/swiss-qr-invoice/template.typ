#import "@preview/payqr-swiss:0.4.1": swiss-qr-bill

#let ink = rgb("#172033")
#let muted = rgb("#64748B")
#let accent = rgb("#1D4ED8")
#let border = rgb("#CBD5E1")
#let soft = rgb("#EFF6FF")

#set document(title: "Swiss QR Invoice {{ get .Record "id" }}")
#set page(paper: "a4", margin: 20mm)
#set text(font: "Libertinus Serif", size: 9pt, fill: ink)
#set par(leading: 5pt)

#let detail(label, value) = grid(
  columns: (30mm, 1fr),
  gutter: 5mm,
  text(size: 8pt, fill: muted)[#label],
  text(weight: "bold")[#value],
)

#grid(
  columns: (1fr, auto),
  gutter: 15mm,
  [
    #text(size: 18pt, weight: "bold", fill: accent)[{{ get .Meta "company_name" }}]
    #v(3pt)
    #text(fill: muted)[
      {{ get .Meta "creditor_street" }} {{ get .Meta "creditor_building" }}
      #linebreak()
      {{ get .Meta "creditor_postal_code" }} {{ get .Meta "creditor_city" }}
      #linebreak()
      {{ get .Meta "company_email" }}
    ]
  ],
  [
    #align(right)[
      #text(size: 20pt, weight: "bold")[INVOICE]
      #v(3pt)
      #text(size: 11pt, fill: accent)[{{ get .Record "id" }}]
    ]
  ],
)

#v(14mm)

#grid(
  columns: (1fr, 1fr),
  gutter: 15mm,
  [
    #text(size: 8pt, weight: "bold", fill: muted)[BILL TO]
    #v(4pt)
    #text(weight: "bold")[{{ get .Record "customer_name" }}]
    #linebreak()
    {{ get .Record "customer_street" }} {{ get .Record "customer_building" }}
    #linebreak()
    {{ get .Record "customer_postal_code" }} {{ get .Record "customer_city" }}
  ],
  [
    #detail("Invoice date", "{{ date "02.01.2006" (get .Record "invoice_date") }}")
    #v(4pt)
    #detail("Due date", "{{ date "02.01.2006" (get .Record "due_date") }}")
    #v(4pt)
    #detail("Currency", "{{ get .Meta "currency" }}")
  ],
)

#v(14mm)

#table(
  columns: (1fr, 35mm),
  inset: (x: 8pt, y: 7pt),
  stroke: (x: none, y: 0.6pt + border),
  align: (left, right),
  table.header(
    table.cell(fill: soft)[#text(weight: "bold", fill: accent)[DESCRIPTION]],
    table.cell(fill: soft)[#text(weight: "bold", fill: accent)[AMOUNT]],
  ),
  [{{ get .Record "description" }}],
  [{{ get .Meta "currency" }} {{ money (get .Record "amount") }}],
  table.cell(colspan: 1, align: right)[#text(weight: "bold")[Total due]],
  table.cell(fill: accent)[
    #text(weight: "bold", fill: white)[
      {{ get .Meta "currency" }} {{ money (get .Record "amount") }}
    ]
  ],
)

#v(8mm)
#rect(
  width: 100%,
  fill: soft,
  stroke: 0.6pt + border,
  radius: 4pt,
  inset: 10pt,
)[
  #text(weight: "bold", fill: accent)[Payment information]
  #v(3pt)
  Please use the Swiss QR payment section below. Payment is due by
  {{ date "January 2, 2006" (get .Record "due_date") }}.
]

#place(
  bottom + center,
  dy: 20mm,
)[
  #swiss-qr-bill(
    account: "{{ get .Meta "account" }}",
    creditor-name: "{{ get .Meta "company_name" }}",
    creditor-street: "{{ get .Meta "creditor_street" }}",
    creditor-building: "{{ get .Meta "creditor_building" }}",
    creditor-postal-code: "{{ get .Meta "creditor_postal_code" }}",
    creditor-city: "{{ get .Meta "creditor_city" }}",
    creditor-country: "{{ get .Meta "creditor_country" }}",
    amount: {{ get .Record "amount" }},
    currency: "{{ get .Meta "currency" }}",
    debtor-name: "{{ get .Record "customer_name" }}",
    debtor-street: "{{ get .Record "customer_street" }}",
    debtor-building: "{{ get .Record "customer_building" }}",
    debtor-postal-code: "{{ get .Record "customer_postal_code" }}",
    debtor-city: "{{ get .Record "customer_city" }}",
    debtor-country: "{{ get .Record "customer_country" }}",
    reference-type: "NON",
    additional-info: "Invoice {{ get .Record "id" }}",
    language: "en",
    font: "page",
  )
]
