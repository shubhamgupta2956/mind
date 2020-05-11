#!/bin/sh

set -e

source $PWD/info

PKG_NAME="github.com/shubhamgupta2956/mind"

LD_FLAGS="
  -X '${PKG_NAME}.SlackClientID=${SLACK_CLIENT_ID}'
  -X '${PKG_NAME}.SlackClientSecret=${SLACK_CLIENT_SECRET}'
  -X '${PKG_NAME}.TodoistClientID=${TODOIST_CLIENT_ID}'
  -X '${PKG_NAME}.TodoistClientSecret=${TODOIST_CLIENT_SECRET}'
"

echo "[*] Building as ./target/mind"
test -d target || mkdir target
go build -o "./target/mind" -ldflags "${LD_FLAGS}" "${PKG_NAME}/cmd/mind"
echo "[+] Build complete!"
