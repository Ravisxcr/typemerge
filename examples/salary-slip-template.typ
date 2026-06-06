// ------------------------------------------------------------
//  Salary Slip Template
//  Usage: typst compile salary_slip.typ
// ------------------------------------------------------------

// Brand colors. Change these once to rebrand the document.
#let ink = rgb("#111827")
#let muted = rgb("#6B7280")
#let soft = rgb("#F8FAFC")
#let border = rgb("#E5E7EB")
#let accent = rgb("#1D4ED8")
#let accent-soft = rgb("#EFF6FF")
#let success = rgb("#047857")
#let danger = rgb("#B91C1C")

#set document(title: "Salary Slip")
#set page(
  paper: "a4",
  margin: (x: 20mm, y: 18mm),
)
#set text(font: "Libertinus Serif", size: 10pt, fill: ink)
#set par(leading: 5.5pt)

#let label(content) = text(size: 7pt, weight: "bold", fill: muted, tracking: 0.8pt)[#content]
#let value(content) = text(size: 9pt, weight: "bold", fill: ink)[#content]

#let kv(name, val) = grid(
  columns: (1fr, auto),
  gutter: 8pt,
  align: (left, right),
  label(upper(name)),
  value(val),
)

#let section-title(content) = [
  #text(size: 8pt, fill: accent, weight: "bold", tracking: 1pt)[#upper(content)]
  #v(5pt)
]

#let detail-card(title, rows) = block(
  fill: white,
  stroke: 0.6pt + border,
  radius: 6pt,
  inset: 12pt,
  width: 100%,
)[
  #stack(
    dir: ttb,
    spacing: 7pt,
    label(title),
    line(length: 100%, stroke: 0.5pt + border),
    rows,
  )
]

#let amount-style(content, fill: ink, strong: false) = text(
  size: if strong { 10pt } else { 9pt },
  weight: if strong { "bold" } else { "regular" },
  fill: fill,
)[#content]

#let slip-row(component, kind, amount, kind-fill: muted, is-dark: false) = block(
  fill: if is-dark { soft } else { white },
  width: 100%,
)[
  #pad(x: 14pt, y: 6pt)[
    #grid(
      columns: (2fr, 1fr, 1fr),
      gutter: 8pt,
      align: (left, center, right),
      text(size: 9pt, fill: ink)[#component],
      text(size: 8pt, fill: kind-fill)[#kind],
      amount-style(amount),
    )
  ]
  #line(length: 100%, stroke: 0.4pt + border)
]

// ------------------------------------------------------------
// Header
// ------------------------------------------------------------
#grid(
  columns: (1fr, auto),
  gutter: 18pt,
  align: (left, right),
  [
    #stack(
      dir: ltr,
      spacing: 12pt,
      block(
        width: 38pt,
        height: 38pt,
        fill: accent,
        radius: 7pt,
        align(center + horizon, text(size: 18pt, weight: "bold", fill: white)[
          {{ initial (get .Meta "company_name") }}
        ]),
      ),
      stack(
        dir: ttb,
        spacing: 3pt,
        text(size: 16pt, weight: "bold")[{{ get .Meta "company_name" }}],
        text(size: 8.5pt, fill: muted)[{{ get .Meta "address" }}],
      ),
    )
  ],
  [
    #block(
      fill: accent-soft,
      stroke: 0.6pt + rgb("#BFDBFE"),
      radius: 6pt,
      inset: (x: 14pt, y: 8pt),
    )[
      #stack(
        dir: ttb,
        spacing: 2pt,
        align(right, label("SALARY SLIP")),
        align(right, text(size: 13pt, weight: "bold", fill: accent)[{{ get .Meta "month" }} {{ get .Meta "year" }}]),
      )
    ]
  ],
)

#v(10pt)
#line(length: 100%, stroke: 0.7pt + border)
#v(12pt)

// ------------------------------------------------------------
// Employee and payment details
// ------------------------------------------------------------
#grid(
  columns: (1fr, 1fr),
  gutter: 12pt,
  detail-card("Employee Details", [
    #kv("Employee ID", "{{ get .Employee "employee_id" }}")
    #kv("Full Name", "{{ get .Employee "name" }}")
    #kv("Designation", "{{ get .Employee "designation" }}")
  ]),
  detail-card("Payment Details", [
    #kv("Payment Date", "{{ get .Meta "payment_date" }}")
    #kv("Pay Period", "{{ get .Meta "month" }} {{ get .Meta "year" }}")
    #kv("Currency", "{{ get .Meta "currency" }}")
  ]),
)

#v(14pt)

// ------------------------------------------------------------
// Earnings and deductions
// ------------------------------------------------------------
#section-title("Earnings and Deductions")

#block(
  fill: white,
  stroke: 0.6pt + border,
  radius: 7pt,
  clip: true,
  width: 100%,
)[
  #block(fill: rgb("#F1F5F9"), width: 100%)[
    #pad(x: 14pt, y: 7pt)[
      #grid(
        columns: (2fr, 1fr, 1fr),
        gutter: 8pt,
        align: (left, center, right),
        label("COMPONENT"),
        label("TYPE"),
        label("AMOUNT"),
      )
    ]
  ]
  #line(length: 100%, stroke: 0.5pt + border)

  #slip-row("Basic Salary", "Earning", "{{ money (get .Employee "basic") }}", kind-fill: success)
  #slip-row("House Rent Allowance", "Earning", "{{ money (get .Employee "hra") }}", kind-fill: success, is-dark: true)
  #slip-row("Other Allowances", "Earning", "{{ money (get .Employee "allowances") }}", kind-fill: success)
  #slip-row("Deductions", "Deduction", "{{ money (get .Employee "deductions") }}", kind-fill: danger, is-dark: true)

  #block(fill: accent, width: 100%)[
    #pad(x: 14pt, y: 8pt)[
      #grid(
        columns: (2fr, 1fr, 1fr),
        gutter: 8pt,
        align: (left, center, right),
        text(size: 10pt, weight: "bold", fill: white)[Net Pay],
        [],
        amount-style("{{ money (get .Employee "net_pay") }}", fill: white, strong: true),
      )
    ]
  ]
]

#v(12pt)

#block(
  fill: accent-soft,
  stroke: 0.6pt + rgb("#BFDBFE"),
  radius: 6pt,
  inset: 12pt,
  width: 100%,
)[
  #grid(
    columns: (auto, 1fr),
    gutter: 10pt,
    align: (left, left),
    text(size: 18pt, weight: "bold", fill: accent)[{{ get .Meta "currency" }}],
    stack(
      dir: ttb,
      spacing: 2pt,
      text(size: 7pt, fill: muted, weight: "bold", tracking: 0.8pt)[TOTAL NET PAY],
      text(size: 15pt, fill: ink, weight: "bold")[{{ money (get .Employee "net_pay") }}],
    ),
  )
]

#v(1fr)

// ------------------------------------------------------------
// Footer
// ------------------------------------------------------------
#line(length: 100%, stroke: 0.5pt + border)
#v(7pt)
#grid(
  columns: (1fr, auto),
  gutter: 18pt,
  align: (left + bottom, right + bottom),
  text(size: 7.5pt, fill: muted)[
    This is a computer-generated salary slip issued by {{ get .Meta "company_name" }}.
    It does not require a physical signature.
  ],
  stack(
    dir: ttb,
    spacing: 4pt,
    align(right, line(length: 90pt, stroke: 0.6pt + rgb("#9CA3AF"))),
    align(right, text(size: 8pt, fill: muted)[Authorized Signatory]),
  ),
)
