#!/bin/bash
set -eux
set -o pipefail

###############################################################################
###############################################################################
## This JamfPRO bash script drops tokens to directories under users' home    ##
## (e.g. /User/john/Downloads, /User/john/secret ... etc.).                  ##
## Since JamfPRO reserves variables from $0 to $3 for its own execution, the ##
## first varaible available to JamfPRO operators is $4:                      ##
##   - $4 is the console domain hash (e.g. aabbccddee), NOT the full domain  ##
##   - $5 is the console API key                                             ##
##   - $6 is the types of tokens to be dropped, separated by comma ","       ##
##     e.g. "aws-id,doc-msword"                                              ##
##   - $7 are the directories relevant to users' home, separated by space    ##
##     e.g. "./ Documents Secret" will drop tokens to the following:         ##
##            - "~/" home directory                                          ##
##            - /Users/john/Docments                                         ##
##            - /Users/john/Secret                                           ##
###############################################################################
###############################################################################

# Console domain hash
# if your domain is aabbccddee.canary.tools, then ONLY provide the aabbccddee part
DOMAIN="$4"

# API Auth Token
APIKEY="$5"

# A list of token-types to be dropped, seperated by space
# e.g "aws-id,doc-msword,msword-macro,msexcel-macro"
TOKEN_TYPES="$6"

# A list of user directories we want to drop to "/Users/$target_user/$directory"
# seperated by space
# e.g ". Documents Downloads Secret" ... if dir does not exist, it will be created
USER_DIRS="$7"


echo "Dropping token types $TOKEN_TYPES to $USER_DIRS"

# This is the JAMF caching location.
DROPPER="/Library/Application Support/JAMF/Waiting Room/TokenDropper.zip"

# Listing all users
TARGET_USERS=$(ls -ld /Users/* | grep -v Shared | awk '{print $3}')


for target_user in $TARGET_USERS;
do
    work_dir=$(sudo -u $target_user mktemp -d)
    cp "$DROPPER" "$work_dir"
    chown -R $target_user "$work_dir"
    pushd $work_dir
    sudo -u $target_user unzip TokenDropper.zip
    sudo -u $target_user chmod +x TokenDropper
    for user_dir in $USER_DIRS;
    do
        where="/Users/$target_user/$user_dir"
        echo "Generating token types: $TOKEN_TYPES for $target_user at $where"

        sudo -u $target_user ./TokenDropper -count 1 \
            -kind "$TOKEN_TYPES" \
            -where "$where" \
            -flock "TokenDropper" \
            -domain $DOMAIN \
            -apikey $APIKEY
            2> tokendropper.log
        # Logs printed to stdout are logged to the JAMF console (useful for debugging)
        sudo -u $target_user cat "$work_dir/tokendropper.log"
    done
    popd
    # Self-destruct (avoids an avenue of detection)
    rm -rf "$work_dir"
done

# Self-destruct (avoids an avenue of detection)
rm -f "$DROPPER"
