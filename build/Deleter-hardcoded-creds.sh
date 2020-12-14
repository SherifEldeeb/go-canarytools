#!/bin/bash
set -ex

if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters $#"
    exit 2
fi

BUILDTIME=$(date '+%Y-%m-%dT%H:%M:%S')
SHA1VER=$(git rev-parse HEAD)
TOOL="CanaryDeleter"

# build the binaries
pushd ../cmd/$TOOL
GOOS=darwin go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/$TOOL-macos/$TOOL
GOOS=linux go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2  -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/$TOOL-linux/$TOOL
GOOS=windows go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2  -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/$TOOL-windows/$TOOL.exe
popd
