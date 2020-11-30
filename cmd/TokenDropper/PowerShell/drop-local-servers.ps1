<#
.SYNOPSIS
    Drops a list of tokens to local machine

.NOTES
    Last Edit: 2020-10-27
    Version 1.0 - initial release
#>

# you can optionally set the full path to the dropper executable here
# if left empty, it will look for it in the same dir as the script
param (
    [string]$Dropper = ""
)

# CSV of (filename), (where to drop), (token type)
# script will iterte over those entries creating those tokens
# the second field (where to drop) is relative to "C:"
# e.g: setting the second field to "Backup" will drop to "c:\Backup"
# another e.g: setting the second field to "windows\SYSVOL" will drop to "c:\windows\SYSVOL"
$full_path_entries = @"
Network_Inventory,c:\Users\Admnistrator\Desktop\Full_Backup,windows-dir
Archive_NYT,C:\ProgramData,windows-dir
logon_configuration.txt,c:\Users\Admnistrator\Desktop,aws-id
Web Credentials,c:\Users\Admnistrator\Desktop,aws-id
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

#####################
# Full Path entries #
#####################
foreach ($entry in $($full_path_entries -split "`n")) {
    $fields = $entry.Trim() -split ","

    $filename = $fields[0]
    $path = $fields[1]
    $kind = $fields[2]

    Write-Host -ForegroundColor Green "[*] FULL_PATH_ENTRY: dropping $filename of $kind token to $path" 
    Write-Host -ForegroundColor Green "[*] Executing: ```Dropper -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`"``` "
    & "$Dropper" -kind `"$kind`" -where `"$path`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`" 2>&1 | ForEach-Object { "$_" }
}