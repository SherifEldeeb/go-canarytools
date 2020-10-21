<#
.SYNOPSIS
    Drops Canarytokens to all users' directories on a set of remote machines

.NOTES
    Last Edit: 2020-10-20
    Version 1.0 - initial release
#>

# you can optionally set the full path to the dropper executable here
# if left empty, it will look for it in the same dir as the script
# RootPath can be set to drop tokens to a remote machine
# e.g. you can drop tokens to 192.168.0.10 frmo the current machine like the following: 
# .\drop-tokens.ps1 -RootPath "\\192.168.0.10\C$" 
param (
    [string]$Dropper = "",
    [string]$RootPath = ""
)

# CSV of (filename), (where to drop), (token type)
# script will iterte over all users under 'c:\users\*'
$entries = @"
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
"@

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

# set the $RootPath
if ($RootPath -eq "") {
    $RootPath = $env:HOMEDRIVE
} 

    
# get list of users, sans "Public"
$users = $(Get-ChildItem $RootPath/Users/ | Where-Object { ($_.PSIsContainer) -and ($_.Name -ne "Public") })

foreach ($user in $users) {
    foreach ($entry in $($entries -split "`n")) {
        $fields = $entry.Trim() -split ","

        $filename = $fields[0]
        $user_dir = $fields[1]
        $kind = $fields[2]

        $where = Join-Path -Path "$RootPath" -ChildPath "/Users/" | Join-Path -ChildPath "$user" | Join-Path -ChildPath "$user_dir"
        Write-Host -ForegroundColor Green "[*] dropping $filename of $kind token to $user_dir" 
        Write-Host -ForegroundColor Green "[*] Executing: ```Dropper -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`"``` "
        & "$Dropper" -kind `"$kind`" -where `"$where`" -filename `"$filename`" -flock `"TokenDropper`" -randomize-filenames=false -memo `"RootPath:$RootPath`" 2>&1 | ForEach-Object { "$_" }
    }  
}
