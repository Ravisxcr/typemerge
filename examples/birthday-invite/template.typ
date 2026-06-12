#let purple = rgb("#6C3BAA")
#let pink = rgb("#E66AA2")
#let yellow = rgb("#F7C948")
#let cream = rgb("#FFF9EE")
#let ink = rgb("#382C45")
#let muted = rgb("#776A80")

#set document(title: "Birthday Invitation")
#set page(
  width: 148mm,
  height: 210mm,
  margin: 12mm,
  fill: cream,
)
#set text(font: "Libertinus Serif", fill: ink)
#set par(justify: false)

#rect(
  width: 100%,
  height: 100%,
  stroke: 1.2pt + pink,
  radius: 7pt,
  inset: 8pt,
)[
  #rect(
    width: 100%,
    height: 100%,
    stroke: 0.5pt + yellow,
    radius: 4pt,
    inset: 22pt,
  )[
    #align(center)[
      #circle(radius: 3pt, fill: yellow)
      #v(7pt)
      #text(size: 8pt, weight: "bold", tracking: 1.6pt, fill: muted)[
        {{ upper (get .Meta "host_line") }}
      ]
      #v(10pt)
      #text(size: 11pt, style: "italic", fill: muted)[
        Come celebrate
      ]
      #v(7pt)

      #text(size: 32pt, style: "italic", fill: purple)[
        {{ get .Employee "celebrant_name" }}
      ]
      #v(5pt)
      #text(size: 13pt, fill: pink)[is turning]
      #v(3pt)
      #text(size: 34pt, weight: "bold", fill: purple)[
        {{ get .Employee "age" }}
      ]

      #v(12pt)
      #line(length: 55%, stroke: 0.7pt + yellow)
      #v(12pt)

      #text(size: 17pt, weight: "bold", fill: purple)[
        {{ date "Monday, January 2, 2006" (get .Employee "birthday_date") }}
      ]
      #v(5pt)
      #text(size: 12pt)[at {{ get .Employee "party_time" }}]
      #v(11pt)

      #text(size: 12pt, weight: "bold")[{{ get .Employee "venue_name" }}]
      #v(3pt)
      #block(width: 78%)[
        #text(size: 9pt, fill: muted)[{{ get .Employee "venue_address" }}]
      ]

      #v(12pt)
      #rect(
        width: 78%,
        fill: white,
        stroke: 0.5pt + pink,
        radius: 4pt,
        inset: 9pt,
      )[
        #text(size: 8pt, fill: muted)[PARTY PLANS]
        #v(3pt)
        #text(size: 10pt, weight: "bold")[
          {{ get .Employee "party_details" }}
        ]
      ]

      #v(1fr)
      #text(size: 9pt, style: "italic", fill: muted)[{{ get .Employee "note" }}]
      #v(9pt)
      #text(size: 8pt, weight: "bold", tracking: 1pt, fill: purple)[
        RSVP BY {{ upper (date "January 2" (get .Employee "rsvp_date")) }}
      ]
      #v(3pt)
      #text(size: 8pt, fill: muted)[{{ get .Meta "rsvp_contact" }}]
      #v(8pt)
      #circle(radius: 3pt, fill: yellow)
    ]
  ]
]
