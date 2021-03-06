#!/bin/bash
set -ex

if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters $#"
    exit 2
fi

BUILDTIME=$(date '+%Y-%m-%dT%H:%M:%S')
SHA1VER=$(git rev-parse HEAD)
TOOL="TokenDropper"

pushd ../cmd/$TOOL
# build the binaries
# macos
GOOS=darwin go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/$TOOL-macos/$TOOL
pushd ../../build/$3/$TOOL-macos/
zip -r -9 $TOOL.zip $TOOL && rm $TOOL
popd

# linux
GOOS=linux go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2  -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/$TOOL-linux/$TOOL
pushd ../../build/$3/$TOOL-linux/
zip -r -9 $TOOL.zip $TOOL && rm $TOOL
popd

# windows
GOOS=windows go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2  -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/$TOOL-windows/$TOOL.exe
pushd ../../build/$3/$TOOL-windows/
zip -r -9 $TOOL.zip $TOOL.exe && rm $TOOL.exe
popd

popd

zip -r -9 "$3-$(date '+%Y-%m-%dT%H_%M_%S').zip" $3