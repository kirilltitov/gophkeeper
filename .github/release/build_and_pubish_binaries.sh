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

createRelease() {
  curl -L -sf \
    -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    https://api.github.com/repos/kirilltitov/gophkeeper/releases \
    -d "{\"tag_name\":\"${TAG_NAME}\",\"target_commitish\":\"main\",\"name\":\"${TAG_NAME}\",\"body\":\"\",\"draft\":false,\"prerelease\":false,\"generate_release_notes\":false}"
}

uploadAsset() {
  local ID=$1
  local FILENAME=$2
  local url=https://uploads.github.com/repos/kirilltitov/gophkeeper/releases/${ID}/assets?name=${FILENAME}
  echo "Upload asset ${FILENAME} to release ${ID}"
  echo "POST ${url}"
  curl -sf \
    -X POST \
    -H "Accept: application/vnd.github+json" \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    -H "Content-Type: application/octet-stream" \
    "${url}" \
    --data-binary "@${FILENAME}"
}

echo "Create Release"
RELEASE_ID=$(createRelease | grep -m1 -oP '(?<="id": )([^,]*)')
echo "Created RELEASE_ID: ${RELEASE_ID}"
cd ./build
ls -la
uploadAsset "${RELEASE_ID}" "client-linux-amd64.${TAG_NAME}.tar.gz"
uploadAsset "${RELEASE_ID}" "client-darwin-amd64.${TAG_NAME}.tar.gz"
uploadAsset "${RELEASE_ID}" "client-windows-amd64.${TAG_NAME}.tar.gz"
