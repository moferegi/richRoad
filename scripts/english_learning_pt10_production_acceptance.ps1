param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8888",
  [string]$Token = "",
  [string]$ProbeWord = "destiny",
  [int]$WordId = 0
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($Token)) {
  throw "Token is required. Use an admin x-token."
}

function Build-RequestUri {
  param(
    [string]$BaseUrl,
    [string]$Path,
    [hashtable]$Query
  )

  $baseUri = [System.Uri]$BaseUrl
  $targetUri = [System.Uri]::new($baseUri, $Path.TrimStart('/'))
  if (-not $Query -or $Query.Count -eq 0) {
    return $targetUri.AbsoluteUri
  }

  $queryParts = foreach ($k in $Query.Keys) {
    "{0}={1}" -f [System.Uri]::EscapeDataString([string]$k), [System.Uri]::EscapeDataString([string]$Query[$k])
  }

  $builder = [System.UriBuilder]::new($targetUri)
  $builder.Query = [string]::Join("&", $queryParts)
  return $builder.Uri.AbsoluteUri
}

function Invoke-JsonApi {
  param(
    [ValidateSet("GET", "POST", "PUT", "DELETE")][string]$Method,
    [string]$Path,
    [hashtable]$Query,
    [object]$Body
  )

  $requestUri = Build-RequestUri -BaseUrl $ApiBaseUrl -Path $Path -Query $Query
  $invokeArgs = @{
    Method  = $Method
    Uri     = $requestUri
    Headers = @{ "x-token" = $Token }
  }

  if ($Body -ne $null) {
    $invokeArgs["ContentType"] = "application/json"
    $invokeArgs["Body"] = ($Body | ConvertTo-Json -Depth 10)
  }

  $resp = Invoke-RestMethod @invokeArgs
  if ($null -eq $resp) {
    throw "Empty response for $Method $Path"
  }

  if ($null -ne $resp.code -and [int]$resp.code -ne 0) {
    throw "API failed ($Method $Path): code=$($resp.code), msg=$($resp.msg)"
  }

  return $resp
}

$script:results = @()
function Add-Result {
  param(
    [string]$Name,
    [bool]$Passed,
    [string]$Detail
  )

  $script:results += [PSCustomObject]@{
    Name   = $Name
    Passed = $Passed
    Detail = $Detail
  }
}

function Run-Step {
  param(
    [string]$Name,
    [scriptblock]$Script
  )

  try {
    & $Script
    Add-Result -Name $Name -Passed $true -Detail "OK"
  }
  catch {
    Add-Result -Name $Name -Passed $false -Detail $_.Exception.Message
  }
}

Run-Step -Name "Check config group" -Script {
  $resp = Invoke-JsonApi -Method "GET" -Path "/sysConfig/getSysConfigByGroup" -Query @{ configGroup = "english_learning" }
  if (-not $resp.data) {
    throw "english_learning config group is empty"
  }
}

Run-Step -Name "Preflight US+UK" -Script {
  $resp = Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/preflightTTS" -Body @{ word = $ProbeWord }
  if ([string]::IsNullOrWhiteSpace([string]$resp.data.audioUs) -or [string]::IsNullOrWhiteSpace([string]$resp.data.audioUk)) {
    throw "audioUs/audioUk missing in preflight response"
  }
}

Run-Step -Name "Preflight US only" -Script {
  $resp = Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/preflightTTS" -Body @{ word = $ProbeWord; checkUs = $true; checkUk = $false }
  if ([string]::IsNullOrWhiteSpace([string]$resp.data.audioUs)) {
    throw "audioUs missing when checkUs=true"
  }
}

Run-Step -Name "Preflight UK only" -Script {
  $resp = Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/preflightTTS" -Body @{ word = $ProbeWord; checkUs = $false; checkUk = $true }
  if ([string]::IsNullOrWhiteSpace([string]$resp.data.audioUk)) {
    throw "audioUk missing when checkUk=true"
  }
}

if ($WordId -gt 0) {
  Run-Step -Name "Regenerate audio for WordId" -Script {
    [void](Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/regenerateAudio" -Body @{ ID = $WordId })
  }
}

Write-Host "`n========== PT-10 Production Acceptance =========="
$script:results | Format-Table -AutoSize

$failed = @($script:results | Where-Object { -not $_.Passed })
if ($failed.Count -gt 0) {
  Write-Host "`nAcceptance: FAIL ($($failed.Count) failed steps)" -ForegroundColor Red
  exit 1
}

Write-Host "`nAcceptance: PASS" -ForegroundColor Green
exit 0
