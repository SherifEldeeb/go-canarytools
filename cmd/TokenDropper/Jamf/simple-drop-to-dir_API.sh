#!/bin/bash
set -eux
set -o pipefail

###############################################################################
###############################################################################
## This JamfPRO bash script will use TokenDropper to drop one aws-id         ##
## token to /opt/classified, and place it in the "TokenDropper" flock.       ##
## If the dir does not exist, it will be created.                            ##
## Since JamfPRO reserves variables from $0 to $3 for its own execution, the ##
## first varaible available to JamfPRO operators is $4:                      ##
##   - $4 is the console domain hash (e.g. aabbccddee), NOT the full domain  ##
##   - $5 is the console API key                                             ##
###############################################################################
###############################################################################

# Console domain hash
# if your domain is aabbccddee.canary.tools, then ONLY provide the aabbccddee part
DOMAIN="$4"

# API Auth Token
APIKEY="$5"

### Options ###
# "Kind" of tokens to be dropped; popular ones are 'aws-id' & 'doc-msword'.
KIND="aws-id"
# This is where the token will be droppped; if dir does not exist it will 
# be created.
WHERE="/opt/classified/"
# Flock that will be used to host the dropped tokens.
FLOCK="TokenDropper"
### End of Options ###

### Activity Starts here ###
# This is the JAMF caching location.
# Jamf Package will be stored locally here...
DROPPER="/Library/Application Support/JAMF/Waiting Room/TokenDropper.zip"

# create a temp workdir for package extraction
WORK_DIR=$(mktemp -d)

# Copy the Package file to temp work dir
cp "$DROPPER" "$WORK_DIR"

# Change directory to work dir
pushd "$WORK_DIR"

# Extract the tool from the package
unzip TokenDropper.zip

# Set the execution bit to the tool
chmod +x TokenDropper

# run the tool
./TokenDropper -count 1 -kind "$KIND" -where "$WHERE" -flock "$FLOCK" -domain $DOMAIN -apikey $APIKEY 2> tokendropper.log

# Logs printed to stdout are logged to the JAMF console (useful for debugging)
cat "$WORK_DIR/tokendropper.log"

# Return from work dir
popd

# Self-destruct work dir (avoids an avenue of detection)
rm -rf "$WORK_DIR"

# Self-destruct package (avoids an avenue of detection)
rm -f "$DROPPER"
