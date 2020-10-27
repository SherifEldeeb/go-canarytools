<#
.SYNOPSIS
    Drops a doc-msword token to c:\windows\SYSVOL, 
    and creates a "windows dir" token on every user's Desktop
    This could be suitable for Domain Controllers

.NOTES
    Last Edit: 2020-10-23
    Version 1.0 - initial release
#>

# you can optionally set the full path to the dropper executable here
# if left empty, it will look for it in the same dir as the script
param (
    [string]$Dropper = ""
)

# list of servers, replace with yours.
$servers = @"
testdomain-dc01
testdomain-dc02
"@

# CSV of (filename), (where to drop), (token type)
# script will iterte over all users under 'c:\users\*'
# and create the full set of tokens per user.
# the second field (where to drop) is relative to each user's home directory
# e.g setting the second field to "Downloads" will drop to "c:\Users\user\Downloads"
# for full path entries that are irrelevant to users
# please check the next variable "$full_path_entries"
#
# Tokens will be sequentially created ... so you can chain them:
# e.g. You can first to create a windows folder canarytoken, then place a word token inside that folder.
#      in the following example: we will create a Canarytoken folder "Full Backup\Latest", then place 
#      an AWS-ID & a Work document inside it.
#
# Note: it is NOT recommended to place the Windows dir Canarytoken folder directly on the desktop, 
#       becuase this will cause false positives, instead, place the token within another folder before
#       placing it on the desktop ... tool will automatically create the full directory tree.
#       e.g. set first field to "Latest", second fields to "Desktop\Full Backup" ... this will create
#       Desktop\Full Backup\Latest
$user_relative_entries = @"
Latest,Desktop\Full Backup,windows-dir
Website Production Password.txt,Desktop\Full Backup\Latest,aws-id
private_key.docx,Desktop\Full Backup\Latest,doc-msword
"@

# CSV of (filename), (where to drop), (token type)
# script will iterte over those entries creating those tokens
# the second field (where to drop) is relative to "C:"
# e.g: setting the second field to "Backup" will drop to "c:\Backup"
# another e.g: setting the second field to "windows\SYSVOL" will drop to "c:\windows\SYSVOL"
$full_path_entries = @"
internal_credentials.docx,windows\SYSVOL,doc-msword
"@

###############################
# if $PATH_TO_DROPPER is empty,
# set the full path for the dropper
# in this example we'll assume the exe is in the same dir as the ps1 script
if ($Dropper -eq "") {
    Write-Host "[!] Path to dropper is not set, will assume it's same as ps1 script" -ForegroundColor Yellow

    $scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Definition
    if ($scriptPath -eq "") {
        $scriptPath = "."
    }
    $Dropper = Join-Path $scriptPath "TokenDropper.exe"
}

if (-not (Test-Path $Dropper)) {
    Write-Error "[x] The file `'$Dropper`' does not exist" -ErrorAction Stop
}

Write-Host "[*] Using: $Dropper" -ForegroundColor Green

foreach ($server in $($servers -split "`n")) {
    $RootPath = "\\" + $server.Trim() + "\C`$"

    #####################
    # Full Path entries #
    #####################
    foreach ($entry in $($full_path_entries -split "`n")) {
        $fields = $entry.Trim() -split ","

        $filename = $fields[0]
        $path = $fields[1]
        $kind = $fields[2]

        $where = Join-Path -Path "$RootPath" -ChildPath "$path" 
        Write-Host -ForegroundColor Green "[*] FULL_PATH_ENTRY: dropping $filename of $kind token to $where" 
        Write-Host -ForegroundColor Green "[*] Executing: ```Dropper -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`"``` "
        & "$Dropper" -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`" 2>&1 | ForEach-Object { "$_" }
    }


    #########################
    # User Relative entries #
    #########################
    # get list of users, sans "Public" & "sqladm"
    $users = $(Get-ChildItem "$RootPath\Users\" | Where-Object { ($_.PSIsContainer) -and (($_.Name -ne "Public") -and ($_.Name -ne "sqladm")) })

    foreach ($user in $users) {
        foreach ($entry in $($user_relative_entries -split "`n")) {
            $fields = $entry.Trim() -split ","

            $filename = $fields[0]
            $user_dir = $fields[1]
            $kind = $fields[2]

            $where = Join-Path -Path "$RootPath" -ChildPath "/Users/" | Join-Path -ChildPath "$user" | Join-Path -ChildPath "$user_dir"
            Write-Host -ForegroundColor Green "[*] USER_RELATIVE_ENTRY: dropping $filename of $kind token to $user_dir" 
            Write-Host -ForegroundColor Green "[*] Executing: ```Dropper -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`"``` "
            & "$Dropper" -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`" 2>&1 | ForEach-Object { "$_" }
        }  
    }
}