#!/bin/bash
set -eux
set -o pipefail
IFS=$'\n'

# this is a csv of the provided info
# (filename), (where to drop "relative to user's home dir"), (kind of canarytoken)
read -d '' entries << EOF || true
private_key.docx,.,doc-msword
passwords.docx,.,doc-msword
internal_credentials.docx,Documents,doc-msword
confidential_invoice.docx,Documents,doc-msword
classified_biden_campaign_source_transcript.pdf,Documents,pdf-acrobat-reader
top_secret_times_internal.pdf,Documents,pdf-acrobat-reader
secret_access_keys.docx,Downloads,doc-msword
confidential_salary_payroll_data.docx,Downloads,doc-msword
confidential_trump_taxes_do_not_publish.docx,Downloads,doc-msword
top_secret_putin_emails.docx,Downloads,doc-msword
confidential_sources.pdf,Downloads,pdf-acrobat-reader
Website Production Password.txt,.,aws-id
Story Publishing Credentials.txt,Documents,aws-id
Private Keys.txt,Downloads,aws-id
EOF

# This is the JAMF caching location.
DROPPER="/Library/Application Support/JAMF/Waiting Room/TokenDropper.zip"

# Listing all users
TARGET_USERS=$(ls -ld /Users/* | grep -v Shared | awk '{print $3}')

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
            -randomize-filenames=false \
            -overwrite-files
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
