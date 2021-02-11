#!/bin/bash
set -eux
set -o pipefail
IFS=$'\n'

###############################################################################
###############################################################################
## This JamfPRO bash script drops tokens to directories under users' home    ##
## (e.g. /User/john/Downloads, /User/john/secret ... etc.).                  ##
## In this script, you specify three pieces of info (as CSV):                ##
##   - filename,                                                             ##
##   - where to drop token relative to user's home directory, and            ##
##   - token type.                                                           ##
##                                                                           ##
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

# this is a csv of the provided info
# (filename), (where to drop "relative to user's home dir"), (kind of canarytoken)
read -d '' entries << EOF || true
passwords.docx,.,doc-msword
internal_credentials.docx,Documents,doc-msword
Story Publishing Credentials.txt,Documents,aws-id
Private Keys.txt,Downloads,aws-id
EOF

# This is the JAMF caching location.
DROPPER="/Library/Application Support/JAMF/Waiting Room/TokenDropper.zip"

# Listing all users
TARGET_USERS=$(ls -ld /Users/* | grep -v Shared | grep -v sysops |  awk '{print $3}')

# Iterate over all users on this machine
for target_user in $TARGET_USERS;
do
    # create a temp working dir owned by the user
    work_dir=$(sudo -u $target_user mktemp -d)
    # copy the dropper tool archive from the Jamf cache location to the temp work dir
    cp "$DROPPER" "$work_dir"
    # change file ownership to the user
    chown -R $target_user "$work_dir"
    # change dir to temp (push to stack)
    pushd $work_dir
    # extract the droper tool from archive
    sudo -u $target_user unzip TokenDropper.zip
    # set the execution bit to the dropper tool
    sudo -u $target_user chmod +x TokenDropper

    # we are now ready to execute the tool
    # this will read tthe entries, one by one
    # then split each line at ","
    while read entry; 
    do
        # first value is the filename
        filename=$(echo $entry | cut -d, -f1)
        # second value is the location where tokens will be dropped
        user_dir=$(echo $entry | cut -d, -f2)
        # third value is the kind of token to be dropped
        kind=$(echo $entry | cut -d, -f3)
        
        # this will hold the full path to the dir where tokens will be dropped
        # if it didn't exist, it will be created
        where="/Users/$target_user/$user_dir"

        echo "Dropping token: $kind, $filename for $target_user at $where"

        sudo -u $target_user ./TokenDropper -count 1 \
            -kind "$kind" \
            -where "$where" \
            -filename "$filename" \
            -flock "TokenDropper" \
            -randomize-filenames=false  \
            -domain $DOMAIN \
            -apikey $APIKEY
            2> tokendropper.log
        # Logs printed to stdout are logged to the JAMF console (useful for debugging)
        sudo -u $target_user cat "$work_dir/tokendropper.log"
    done <<< "$entries"
    # return to the previous directory
    popd
    # Self-destruct work dir (avoids an avenue of detection)
    rm -rf "$work_dir"
done

# Self-destruct dropper tool archive (avoids an avenue of detection)
rm -f "$DROPPER"
