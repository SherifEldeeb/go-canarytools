#!/bin/bash
set -ex

# create the work dir
pushd ..
mkdir -p ./$3/bin/windows
mkdir -p ./$3/bin/macos
mkdir -p ./$3/bin/linux


# build the binaries
pushd ../cmd/TokenDropper
GOOS=darwin go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -w -s -linkmode=internal" -o ../../build/$3/macos/TokenDropper
GOOS=linux go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -w -s -linkmode=internal" -o ../../build/$3/linux/TokenDropper
GOOS=windows go build -v -ldflags "-X main.DOMAIN=$1  -X main.APIKEY=$2 -w -s -linkmode=internal" -o ../../build/$3/windows/TokenDropper.exe
popd

# prep the jamf setup
pushd ./$3/macos
zip -9 TokenDropper.zip TokenDropper
popd
cp ../cmd/TokenDropper/Jamf/drop-with-filenames.sh ./$3/macos/
tar xvfj ./$3/macos/dropper_jamf.tar.gz ./$3/macos/*.sh ./$3/macos/TokenDropper.zip
mv ./$3/macos/dropper_jamf.tar.gz ./$3/

# prep the windows files
cp ../cmd/TokenDropper/PowerShell/*.ps1 ./$3/windows/
zip -9 ./$3/windows/dropper_windows.zip ./$3/windows/drop-with-filenames.ps1 ./$3/windows/TokenDropper.exe
mv ./$3/windows/dropper_windows.zip ./$3/

popd