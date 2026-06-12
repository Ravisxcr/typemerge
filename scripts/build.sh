#!/usr/bin/env sh
set -eu

mkdir -p dist
go build -o dist/typemerge ./cmd/typemerge
