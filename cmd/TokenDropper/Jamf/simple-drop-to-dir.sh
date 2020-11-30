#!/bin/bash
set -eux
set -o pipefail

########################################################################
########################################################################
## This JamfPRO bash script will use TokenDropper to drop one aws-id  ##
## token to /opt/classified, and place it in the "TokenDropper" flock.##
## If the dir does not exist, it will be created.                     ##
########################################################################
########################################################################

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
./TokenDropper -count 1 -kind "$KIND" -where "$WHERE" -flock "$FLOCK" 2> tokendropper.log

# Logs printed to stdout are logged to the JAMF console (useful for debugging)
cat "$WORK_DIR/tokendropper.log"

# Return from work dir
popd

# Self-destruct work dir (avoids an avenue of detection)
rm -rf "$WORK_DIR"

# Self-destruct package (avoids an avenue of detection)
rm -f "$DROPPER"
