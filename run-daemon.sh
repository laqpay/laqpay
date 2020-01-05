#!/usr/bin/env bash

# Runs laqpay in daemon mode configuration

set -x

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo "laqpay binary dir:" "$DIR"
pushd "$DIR" >/dev/null

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
GOLDFLAGS="-X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

GORUNFLAGS=${GORUNFLAGS:-}

go run -ldflags "${GOLDFLAGS}" $GORUNFLAGS cmd/laqpay/laqpay.go \
    -enable-gui=false \
    -launch-browser=false \
    -log-level=debug \
    $@

popd >/dev/null
