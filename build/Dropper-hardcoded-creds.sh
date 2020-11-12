#!/bin/bash
set -ex

BUILDTIME=$(date '+%Y-%m-%dT%H:%M:%S')
SHA1VER=$(git rev-parse HEAD)

if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters $#"
    exit 2
fi


# create the work dir
# mkdir -p ./$3/bin/windows
# mkdir -p ./$3/bin/macos
# mkdir -p ./$3/bin/linux


# build the binaries
pushd ../cmd/TokenDropper
GOOS=darwin go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/macos/TokenDropper
GOOS=linux go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2  -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/linux/TokenDropper
GOOS=windows go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2  -X main.BUILDTIME=$BUILDTIME -X main.SHA1VER=$SHA1VER -w -s -linkmode=internal" -o ../../build/$3/windows/TokenDropper.exe
popd

# prep the jamf setup
pushd ./$3/macos
zip -9 TokenDropper.zip TokenDropper
popd
cp ../cmd/TokenDropper/Jamf/drop-with-filenames.sh ./$3/macos/
tar cvfj ./$3/macos/dropper_jamf.tar.gz ./$3/macos/*.sh ./$3/macos/TokenDropper.zip
mv ./$3/macos/dropper_jamf.tar.gz ./$3/

# prep the windows files
cp ../cmd/TokenDropper/PowerShell/*.ps1 ./$3/windows/
zip -9 ./$3/windows/dropper_windows.zip ./$3/windows/TokenDropper.exe
mv ./$3/windows/dropper_windows.zip ./$3/
