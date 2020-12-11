<#
.SYNOPSIS
    Triggers various actions to generate Canary alerts,
    It will read the list of canaries to be poked from a text file (canaries.txt by default),
    Each line in that text file should be a canary IP or hostname

.NOTES
    Last Edit: 2020-12-07
    Version 1.0 - initial release
    Version 1.1 - Add support for SharePoint Skin
#>

param (
    [string]$CanariesFile = "canaries.txt"
)

# Invoke-PokeCanary function
function Invoke-PokeCanary {
    param (
        [string]$Canary
    )
    Write-Host -ForegroundColor Green "[+] Poking $Canary ..."

    Invoke-PortScanAlert -Canary $Canary
    Invoke-SMBShareAlert -Canary $Canary
    Invoke-HTTPLoginAlert -Canary $Canary
}

function Invoke-PortScanAlert {
    param (
        [string]$Canary
    )
    # port scanning ... we only need five ports to trigger 
    Write-Host -ForegroundColor Yellow "[!] Poke: Port scanning $Canary."
    $ports = @(80, 8080, 22, 21, 1433)
    $ports | ForEach-Object {
        (New-Object Net.Sockets.TcpClient).Connect($Canary, $_)
    }
}

function Invoke-SMBShareAlert {
    param (
        [string]$Canary
    )
    # Opening a share
    # by default, Canary will trigger an alert only if a file is accessed,
    # not merely opening the share.
    # so we'll have to list shares, list files, then copy one of them.
    Write-Host -ForegroundColor Yellow "[!] Poke: Openning a share $Canary."
    $shares = &net.exe view \\$Canary /all | Select-Object -Skip 7 | Where-Object { $_ -match 'disk*' } | ForEach-Object { $_ -match '^(.+?)\s+Disk*' | out-null; $matches[1] } | Where-Object { $_ -notmatch '.+\$' }
    # now the $shares variable should have the shares enabled on that canary,
    # simple sanity check
    if (!$shares) {
        Write-Error "[x] The canary doesn't seem to have shares enabled `'$Canary`'!" 
        return
    }
    foreach ($share in $shares) {
        # get folders under each share
        $doc_files = $(&cmd.exe /c dir /s /b "\\$Canary\$share\*.docx") -split '`n'
        
        # check...
        if (!$doc_files) {
            # no .docx in this share
            continue
        }

        # pick first .docx entry
        $doc_file = $doc_files[0]

        Write-Host -ForegroundColor Green "[+] Reading file off a share '$doc_file'"

        # trigger the alert...
        Get-Content $doc_file > $null

        # one is enough, let's break here.
        break
    }
}

function Invoke-HTTPLoginAlert {
    param (
        [string]$Canary
    )

    $user = 'PokeCanary'
    $pass = 'PokeCanary'

    $pair = "$($user):$($pass)"

    $encodedCreds = [System.Convert]::ToBase64String([System.Text.Encoding]::ASCII.GetBytes($pair))

    $basicAuthValue = "Basic $encodedCreds"

    $Headers = @{
        Authorization = $basicAuthValue
    }

    # Get (cisco)
    Write-Host -ForegroundColor Yellow "[!] Poke: HTTP Login (GET Request) $Canary."
    Invoke-WebRequest -Uri "http://$Canary" -Headers $Headers

    # Post (synoligy)
    $postParams = @{username = $user; password = $pass }
    Write-Host -ForegroundColor Yellow "[!] Poke: HTTP Login (POST Request - Synoligy) $Canary."
    Invoke-WebRequest -Uri "http://$Canary/index.html" -Method POST -Body $postParams

    # Post (SharePoint)
    $postParams = @{"ctl00`$PlaceHolderMain`$signInControl`$UserName" = $user; "ctl00`$PlaceHolderMain`$signInControl`$password" = $pass }
    Write-Host -ForegroundColor Yellow "[!] Poke: HTTP Login (POST Request - SharePoint) $Canary."
    Invoke-WebRequest -Uri "http://$Canary/_forms/default.aspx?ReturnUrl=%2f_layouts%2fAuthenticate.aspx%3fSource%3d%252F&Source=%2F" -Headers $Headers

}


# Check if mandatory param has been provided.
if (-not $CanariesFile) { Write-Error -ErrorAction Stop "[x] You must provide a value for -CanariesFile" }

# Does the file exist?
if (-not (Test-Path $CanariesFile)) {
    Write-Error -ErrorAction Stop "[x] The file `'$CanariesFile`' does not exist" 
}

# Getting content of the Canaries TXT file
# this should have a list of Canary device IPs or host name, each on its own line.
Write-Host -ForegroundColor Green "[+] Reading Canaries' IPs/Hostnames from '$CanariesFile'"
$CanariesText = Get-Content $CanariesFile -ErrorAction Stop

# convert the file content to an array, skipping empty lines
$Canaries = $($CanariesText -split "`n").Where( { $_.Trim() -ne "" })

# id the file empty?
if (!$Canaries) {
    Write-Error -ErrorAction Stop "[x] The file `'$CanariesFile`' is empty!" 
}

# iterate over canaries, poking them one by one
foreach ($Canary in $Canaries) {
    Invoke-PokeCanary -Canary $Canary
}
