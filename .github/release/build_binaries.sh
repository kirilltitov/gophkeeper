#!/bin/bash

set -e

rm -rf ./build
mkdir ./build
export CGO_ENABLED=0

NOW=$(date +'%Y-%m-%d %H:%M:%S')

echo "Build linux"
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.buildVersion=${TAG_NAME}' -X 'main.buildDate=${NOW}' -X 'main.buildCommit=${COMMIT_ID}'" -o ./build/client-linux-amd64 -v ./cmd/client

echo "Build Mac"
GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.buildVersion=${TAG_NAME}' -X 'main.buildDate=${NOW}' -X 'main.buildCommit=${COMMIT_ID}'" -o ./build/client-darwin-amd64 -v ./cmd/client

echo "Build Windows"
GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.buildVersion=${TAG_NAME}' -X 'main.buildDate=${NOW}' -X 'main.buildCommit=${COMMIT_ID}'" -o ./build/client-windows-amd64.exe -v ./cmd/client

echo "Gzip binaries"
tar -czf ./build/client-linux-amd64.${TAG_NAME}.tar.gz ./build/client-linux-amd64
tar -czf ./build/client-darwin-amd64.${TAG_NAME}.tar.gz ./build/client-darwin-amd64
tar -czf ./build/client-windows-amd64.${TAG_NAME}.tar.gz ./build/client-windows-amd64.exe
