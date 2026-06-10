#!/usr/bin/env sh
set -eu

mkdir -p dist
go build -o dist/typemerge ./cmd/typemerge
go build -o dist/typemerge-gui ./cmd/typemerge-gui
