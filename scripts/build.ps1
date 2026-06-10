$ErrorActionPreference = "Stop"

New-Item -ItemType Directory -Force -Path dist | Out-Null
go build -o dist/typemerge.exe ./cmd/typemerge
go build -o dist/typemerge-gui.exe ./cmd/typemerge-gui
