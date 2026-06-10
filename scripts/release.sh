#!/usr/bin/env sh
set -eu

version="${1:?usage: scripts/release.sh VERSION}"
mkdir -p dist

for target in linux/amd64 linux/arm64 windows/amd64 darwin/amd64 darwin/arm64; do
  os="${target%/*}"
  arch="${target#*/}"
  suffix=""
  if [ "$os" = "windows" ]; then
    suffix=".exe"
  fi
  output="dist/typemerge-${version}-${os}-${arch}${suffix}"
  GOOS="$os" GOARCH="$arch" go build -trimpath -o "$output" ./cmd/typemerge
done
