#!/bin/bash

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
GO_DIR="$ROOT_DIR/src"

if [ ! -f "$GO_DIR/go.mod" ]; then
    echo "go.mod not found. Please run 'go mod init' and 'go mod tidy' before building."
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Go not found. Please install Go and ensure it's in your PATH."
    exit 1
fi

BIN_NAME="exam-transcriber"

(cd "$GO_DIR" && CGO_ENABLED=0 go build -ldflags="-s -w" -o "$ROOT_DIR/$BIN_NAME" ./cmd/exam-transcriber)

if [ $? -ne 0 ]; then
    echo "Build failed."
    exit 1
fi

if command -v upx &> /dev/null
then
    echo "UPX found, will compress the binary."
    (cd "$ROOT_DIR" && upx --best "$BIN_NAME")
else
    echo "UPX not found, skipping compression."
fi
