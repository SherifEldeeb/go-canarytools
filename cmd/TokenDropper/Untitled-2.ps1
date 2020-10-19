# VARIABLES - Dynamic/ChangeMe
$CanaryHost = "353ce306.canary.tools"
$CanaryApiKey = "47859316b488477a720d070d115c6aac"
$FlockId = "flock:3bb3005767a64d9eda0d50eae985c128"
$Memo = "Example Memo"
$Kind = "doc-msword"
$DocPath = "C:\Users\jpsaroud\Documents\CERT-EU\001-Meli_Project\Automation\Scripts"
$DocName = "RTReport_ECA_ver_0.docx"
$DocType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
$OutFilePath = "$DocPath"
$OutFileName = "Tokenised-$DocName"
$OutFile = "$OutFilePath\$OutFileName"
# VARIABLES - Static/LeaveMe
$Uri = "https://$CanaryHost/api/v1"
$ApiPing = "$Uri/ping?auth_token=$CanaryApiKey"
$ApiCreateToken = "$Uri/canarytoken/create"
$ApiDownloadToken = "$Uri/canarytoken/download"
$Doc = "$DocPath\$DocName"
# Add the System.Net.Http Assembly
add-type -AssemblyName System.Net.Http
$error.Clear()
# Continue past errors
$ErrorActionPreference = 'SilentlyContinue'
# Confirm path of document
Write-Host "Confirming path = $Doc"
# Check that the API is reachable (both host and API key are valid)
Invoke-WebRequest -Uri $ApiPing | out-null
if ($error.count -ne 0) {
    Write-Host -ForegroundColor Red "Could not connect to the Canary Console at $CanaryHost."
    Write-Host -ForegroundColor White -BackgroundColor Black "The error was:"
    Write-Host $error[0]
    Write-Host -ForegroundColor White -BackgroundColor Black "Please ensure the Canary Console hostname and API key are present in the script."
    return
}
### CREATE TOKEN SECTION ###
# Initiate multipartContent
$multipartContent = [System.Net.Http.MultipartFormDataContent]::new()
# Flock content
$stringHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$stringHeader.Name = "flock_id"
$StringContent = [System.Net.Http.StringContent]::new($FlockId)
$StringContent.Headers.ContentDisposition = $stringHeader
$multipartContent.Add($stringContent)
# Memo content
$stringHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$stringHeader.Name = "memo"
$StringContent = [System.Net.Http.StringContent]::new($Memo)
$StringContent.Headers.ContentDisposition = $stringHeader
$multipartContent.Add($stringContent)
# Kind content
$stringHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$stringHeader.Name = "kind"
$StringContent = [System.Net.Http.StringContent]::new($Kind)
$StringContent.Headers.ContentDisposition = $stringHeader
$multipartContent.Add($stringContent)
# File content
$multipartFile = $Doc
$FileStream = [System.IO.FileStream]::new($multipartFile, [System.IO.FileMode]::Open)
$fileHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$fileHeader.Name = "doc"
$fileHeader.FileName = $DocName
$fileContent = [System.Net.Http.StreamContent]::new($FileStream)
$fileContent.Headers.ContentType = [System.Net.Http.Headers.MediaTypeHeaderValue]::Parse($DocType)
$fileContent.Headers.ContentDisposition = $fileHeader
$multipartContent.Add($fileContent)
# Auth Token content
$stringHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$stringHeader.Name = "auth_token"
$StringContent = [System.Net.Http.StringContent]::new($CanaryApiKey)
$StringContent.Headers.ContentDisposition = $stringHeader
$multipartContent.Add($stringContent)
# "Bodybuilding"
$Body = $multipartContent
# Tokenise the document
Write-Host "Tokenising Document"
$response = Invoke-RestMethod $ApiCreateToken -Method 'POST' -Body $Body
if ($error.count -ne 0) {
    Write-Host -ForegroundColor Red "Could not create Canarytoken."
    Write-Host -ForegroundColor White -BackgroundColor Black "The error was:"
    Write-Host $error
    return}
$response | ConvertTo-Json
### DOWNLOAD TOKEN SECTION ###
# Initiate multipartContent
$multipartContent = [System.Net.Http.MultipartFormDataContent]::new()
# Auth Token content
$stringHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$stringHeader.Name = "auth_token"
$StringContent = [System.Net.Http.StringContent]::new($CanaryApiKey)
$StringContent.Headers.ContentDisposition = $stringHeader
$multipartContent.Add($stringContent)
# Canarytoken ID content
$stringHeader = [System.Net.Http.Headers.ContentDispositionHeaderValue]::new("form-data")
$stringHeader.Name = "canarytoken"
$StringContent = [System.Net.Http.StringContent]::new($response.canarytoken.canarytoken)
$StringContent.Headers.ContentDisposition = $stringHeader
$multipartContent.Add($stringContent)
# "Bodybuilding"
$Body = $multipartContent
# Download the tokenised document
Write-Host "Downloading Tokenised Document"
Invoke-RestMethod $ApiDownloadToken -Method 'GET' -Body $Body -OutFile $OutFile
if ($error.count -ne 0) {
    Write-Host -ForegroundColor Red "Could not download Canarytoken."
    Write-Host -ForegroundColor White -BackgroundColor Black "The error was:"
    Write-Host $error
    return}