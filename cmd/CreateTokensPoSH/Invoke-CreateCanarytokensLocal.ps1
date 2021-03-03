<#
.SYNOPSIS
    Creates Canarytokens and drops them to local host.
    Uses main console API, which is not a good idea if you want to run it on many hosts.
    If you want to run this on multiple hosts, there is another version that uses Canarytokes
    Factory instead.

.NOTES
    For this tool to work, you must have your Canary Console API enabled, please 
    follow this link to learn how to do so:
    https://help.canary.tools/hc/en-gb/articles/360012727537-How-does-the-API-work-

    ###################
    How does this work?
    ###################
    1. Make ure the host has access to the internet.
    2. Run powershell as a user that has read/write access on the target directory.
    
    Last Edit: 2021-03-03
    Version 1.0 - initial release

.EXAMPLE
    .\Invoke-CreateCanarytokensLocal.ps1
    This will run the tool with the default params

    .\Invoke-CreateCanarytokens.ps1 -TargetDirectory "c:\secret" -TokenType aws-id -TokenFilename aws_secret.txt
    creates an AWS-ID Canarytoken, using aws_secret.txt as the filename, and place it under c:\secret
#>

Param (
    # Full canary domain (e.g. aabbccdd.canary.tools), 
    # if empty, will be asked for interactively 
    [string]$Domain = '',

    # API Auth token, 
    # if empty, will be asked for interactively 
    [string]$ApiAuth = '',

    # Set the target Directory on hosts' root
    # e.g. 'c:\Backup'
    # will be created if not exists
    [string]$TargetDirectory = "c:\Backup",

    # Valid TokenType are as follows:
    #   "aws-id":"Amazon API Key",
    #   "doc-msword":"MS Word .docx Document",
    #   "msexcel-macro":"MS Excel .xlsm Document",
    #   "msword-macro":"MS Word .docm Document",
    #   "pdf-acrobat-reader":"Acrobat Reader PDF Document",
    # if you change $TokenType, make sure to pick an appropriate filename extension in next line
    [string]$TokenType = 'doc-msword' ,
    [string]$TokenFilename = "credentials.docx"
)

# We force TLS1.2 since our API doesn't support lower.
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12;
Set-StrictMode -Version 2.0

# Connect to API
# Get Console Domain
$ApiHost = [string]::Empty
if ($Domain -ne '') {
    $ApiHost = $Domain
} else {
    Do {
        $ApiHost = Read-Host -Prompt "[+] Enter your Full Canary domain (e.g. 'xyz.canary.tools')"
    } Until (($ApiHost.Length -gt 0) -and ([System.Net.Dns]::GetHostEntry($ApiHost).AddressList[0].IPAddressToString))
}

# Get API Auth Token
$ApiToken = [string]::Empty
if ($ApiAuth -ne '') {
    $ApiToken = $ApiAuth
} else {
    $ApiTokenSecure = New-Object System.Security.SecureString
    Do {
        $ApiTokenSecure = Read-Host -AsSecureString -Prompt "[+] Enter your Canary API key"
    } Until ($ApiTokenSecure.Length -gt 0)
    $ApiToken = (New-Object System.Management.Automation.PSCredential "user", $ApiTokenSecure).GetNetworkCredential().Password
}

Write-Host -ForegroundColor Green "[*] Starting Script with the following params:
        Console Domain   = $ApiHost
        Target Directory = $TargetDirectory 
        Token Type       = $TokenType
        Token Filename   = $TokenFilename
"

$ApiBaseURL = '/api/v1'
Write-Host -ForegroundColor Green "[*] Pinging Console..."

$PingResult = Invoke-RestMethod -Method Get -Uri "https://$ApiHost$ApiBaseURL/ping?auth_token=$ApiToken"
$Result = $PingResult.result
If ($Result -ne 'success') {
    Write-Host -ForegroundColor Red "[X] Cannot ping Canary API. Bad token?"
    Exit
}
Else {
    Write-Host -ForegroundColor Green "[*] Canary API available for service!"
}

Write-Host -ForegroundColor Green "[*] Checking if '$TargetDirectory' exists..."

# Create the target Dir if not exist
If (!(Test-Path $TargetDirectory)) {
    Write-Host -ForegroundColor Green "[*] '$TargetDirectory' doesn't exist, creating it ..."
    New-Item -ItemType Directory -Force -Verbose -ErrorAction Stop -Path "$TargetDirectory"
}
# Check whether token already exists
$OutputFileName = "$TargetDirectory\$TokenFilename"
Write-Host -ForegroundColor Green "[*] Dropping '$OutputFileName' ..."

If (Test-Path $OutputFileName) {
    Write-Host Skipping $OutputFileName, file already exists.
    Continue        
}

# Create token
$TokenName = $OutputFileName
$PostData = @{
    auth_token = "$ApiToken"
    kind       = "$TokenType"
    memo       = "$([System.Net.Dns]::GetHostName()) - $TokenName"
}
Write-Host -ForegroundColor Green "[*] Hitting API to create token ..."
$CreateResult = Invoke-RestMethod -Method Post -Uri "https://$ApiHost$ApiBaseURL/canarytoken/create" -Body $PostData
$Result = $CreateResult.result
If ($Result -ne 'success') {
    Write-Host -ForegroundColor Red "[X] Creation of $TokenName failed."
    Exit
}
Else {
    $WordTokenID = $($CreateResult).canarytoken.canarytoken
    Write-Host -ForegroundColor Green "[*] Token Created (ID: $WordTokenID)."
}

# Download token
Write-Host -ForegroundColor Green "[*] Downloading Token from Console..."
Invoke-RestMethod -Method Get -Uri "https://$ApiHost$ApiBaseURL/canarytoken/download?auth_token=$ApiToken&canarytoken=$WordTokenID" -OutFile "$OutputFileName"
Write-Host -ForegroundColor Green "[*] Token Successfully written to destination: '$OutputFileName'."
