#let navy = rgb("#17324D")
#let gold = rgb("#B78A3D")
#let cream = rgb("#FBF7EE")
#let muted = rgb("#667085")

#set document(title: "Certificate of Achievement")
#set page(
  paper: "a4",
  flipped: true,
  margin: 12mm,
  fill: cream,
)
#set text(font: "Libertinus Serif", fill: navy)
#set par(justify: false)

#let rule = line(length: 100%, stroke: 1pt + gold)

#rect(
  width: 100%,
  height: 100%,
  stroke: 2pt + navy,
  inset: 7pt,
)[
  #rect(
    width: 100%,
    height: 100%,
    stroke: 0.8pt + gold,
    inset: 28pt,
  )[
    #align(center)[
      #text(size: 9pt, weight: "bold", tracking: 2pt, fill: gold)[
        {{ upper (get .Meta "organization") }}
      ]
      #v(12pt)
      #rule
      #v(18pt)

      #text(size: 34pt, weight: "bold")[Certificate of Achievement]
      #v(10pt)
      #text(size: 12pt, fill: muted)[This certificate is proudly presented to]
      #v(15pt)

      #text(size: 29pt, style: "italic", fill: gold)[
        {{ get .Record "recipient_name" }}
      ]
      #v(5pt)
      #line(length: 58%, stroke: 0.7pt + gold)
      #v(15pt)

      #text(size: 12pt, fill: muted)[for successfully completing]
      #v(8pt)
      #text(size: 19pt, weight: "bold")[{{ get .Record "achievement" }}]
      #v(9pt)
      #block(width: 75%)[
        #text(size: 10.5pt, fill: muted)[{{ get .Record "citation" }}]
      ]
      #v(18pt)

      #grid(
        columns: (1fr, 1fr, 1fr),
        gutter: 32pt,
        align: center,
        [
          #text(size: 11pt, weight: "bold")[{{ date "January 2, 2006" (get .Record "issue_date") }}]
          #v(5pt)
          #line(length: 100%, stroke: 0.6pt + navy)
          #v(4pt)
          #text(size: 8pt, fill: muted)[DATE]
        ],
        [
          #circle(
            radius: 28pt,
            fill: gold,
            stroke: 1pt + navy,
          )[
            #align(center + horizon)[
              #text(size: 8pt, weight: "bold", fill: white)[
                {{ upper (get .Meta "seal_text") }}
              ]
            ]
          ]
        ],
        [
          #text(size: 11pt, weight: "bold")[{{ get .Meta "signatory" }}]
          #v(5pt)
          #line(length: 100%, stroke: 0.6pt + navy)
          #v(4pt)
          #text(size: 8pt, fill: muted)[{{ upper (get .Meta "signatory_title") }}]
        ],
      )
      #v(10pt)
      #text(size: 8pt, fill: muted)[Certificate ID: {{ get .Record "id" }}]
    ]
  ]
]
