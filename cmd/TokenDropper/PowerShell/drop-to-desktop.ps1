<#
.SYNOPSIS
    Drops Canarytokens to all users' directories

.NOTES
    Last Edit: 2020-10-20
    Version 1.0 - initial release
#>

# you can optionally set the full path to the dropper executable here
# if left empty, it will look for it in the same dir as the script
param (
    [string]$Dropper = ""
)

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

# CSV of (filename), (where to drop), (token type)
# script will iterte over all users under 'c:\users\*'
$entries = @"
private_key.docx,Desktop,doc-msword
Website Production Password.txt,Desktop,aws-id
"@

# get list of users, sans "Public"
$users = $(Get-ChildItem c:/Users/ | Where-Object { ($_.PSIsContainer) -and (($_.Name -ne "Public") -and ($_.Name -ne "sqladm") -and ($_.Name -ne "svc_account")) })

# for each entry, we will create 
foreach ($user in $users) {
    foreach ($entry in $($entries -split "`n")) {
        $fields = $entry.Trim() -split ","

        $filename = $fields[0]
        $user_dir = $fields[1]
        $kind = $fields[2]

        $where = Join-Path -Path "c:" -ChildPath "/Users/" | Join-Path -ChildPath "$user" | Join-Path -ChildPath "$user_dir"
        Write-Host "[*] dropping $filename of $kind token to $user_dir" -ForegroundColor Green
        Write-Host "[*] Executing: ```Dropper -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false``` "
        & "$Dropper" -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false 2>&1 | ForEach-Object{ "$_" }
    }  
}
