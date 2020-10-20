#!/bin/bash
set -e

pushd ../cmd/TokenDropper

go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -w -s -linkmode=internal" -o ../../build/TokenDropper

popd

zip -9 TokenDropper.zip TokenDropper