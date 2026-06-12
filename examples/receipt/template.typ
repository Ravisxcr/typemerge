#let ink = rgb("#202124")
#let muted = rgb("#6B7280")
#let border = rgb("#D1D5DB")
#let accent = rgb("#0F766E")
#let soft = rgb("#F0FDFA")

#set document(title: "Payment Receipt")
#set page(
  width: 105mm,
  height: 180mm,
  margin: 10mm,
)
#set text(font: "Libertinus Serif", size: 9pt, fill: ink)
#set par(leading: 4.5pt)

#let row(label, value, strong: false) = grid(
  columns: (1fr, auto),
  gutter: 8pt,
  text(size: 8pt, fill: muted)[#label],
  text(size: if strong { 10pt } else { 8pt }, weight: if strong { "bold" } else { "regular" })[#value],
)

#align(center)[
  #rect(
    width: 42pt,
    height: 42pt,
    fill: accent,
    radius: 8pt,
  )[
    #align(center + horizon)[
      #text(size: 18pt, weight: "bold", fill: white)[
        {{ initial (get .Meta "business_name") }}
      ]
    ]
  ]
  #v(7pt)
  #text(size: 15pt, weight: "bold")[{{ get .Meta "business_name" }}]
  #v(2pt)
  #text(size: 7.5pt, fill: muted)[{{ get .Meta "address" }}]
  #linebreak()
  #text(size: 7.5pt, fill: muted)[{{ get .Meta "contact" }}]
]

#v(10pt)
#line(length: 100%, stroke: 0.7pt + border)
#v(8pt)

#grid(
  columns: (1fr, auto),
  align: (left, right),
  [
    #text(size: 13pt, weight: "bold", fill: accent)[PAYMENT RECEIPT]
    #v(2pt)
    #text(size: 8pt, fill: muted)[Thank you for your payment]
  ],
  [
    #text(size: 8pt, fill: muted)[RECEIPT NO.]
    #linebreak()
    #text(size: 9pt, weight: "bold")[{{ get .Employee "id" }}]
  ],
)

#v(10pt)
#rect(
  width: 100%,
  fill: soft,
  stroke: 0.6pt + rgb("#99F6E4"),
  radius: 5pt,
  inset: 9pt,
)[
  #row("Received from", "{{ get .Employee "customer_name" }}", strong: true)
  #v(5pt)
  #row("Date", "{{ date "January 2, 2006" (get .Employee "payment_date") }}")
  #v(5pt)
  #row("Payment method", "{{ get .Employee "payment_method" }}")
]

#v(11pt)
#text(size: 7pt, weight: "bold", tracking: 1pt, fill: muted)[PAYMENT DETAILS]
#v(5pt)
#line(length: 100%, stroke: 0.6pt + border)
#v(7pt)
#row("Description", "{{ get .Employee "description" }}")
#v(7pt)
#row("Reference", "{{ get .Employee "reference" }}")
#v(7pt)
#row("Subtotal", "{{ money (get .Employee "subtotal") }}")
#v(7pt)
#row("Tax", "{{ money (get .Employee "tax") }}")
#v(7pt)
#line(length: 100%, stroke: 0.6pt + border)
#v(8pt)

#rect(
  width: 100%,
  fill: accent,
  radius: 5pt,
  inset: 10pt,
)[
  #grid(
    columns: (1fr, auto),
    text(size: 10pt, weight: "bold", fill: white)[TOTAL PAID],
    text(size: 15pt, weight: "bold", fill: white)[
      {{ get .Meta "currency" }} {{ money (get .Employee "total") }}
    ],
  )
]

#v(1fr)
#align(center)[
  #line(length: 55%, stroke: 0.6pt + border)
  #v(4pt)
  #text(size: 8pt, fill: muted)[Authorized by {{ get .Meta "authorized_by" }}]
  #v(8pt)
  #text(size: 7pt, fill: muted)[This computer-generated receipt is valid without a signature.]
]
