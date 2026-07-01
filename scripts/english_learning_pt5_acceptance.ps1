param(
  [string]$RepoRoot = "",
  [string]$ApiBaseUrl = "http://127.0.0.1:8888",
  [string]$Token = "",
  [int]$WordId = 0,
  [int]$EpisodeId = 0
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($RepoRoot)) {
  $RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
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
  } catch {
    Add-Result -Name $Name -Passed $false -Detail $_.Exception.Message
  }
}

function Invoke-EnglishApi {
  param(
    [ValidateSet("GET", "POST", "PUT", "DELETE")][string]$Method,
    [string]$Path,
    [hashtable]$Query,
    [object]$Body
  )

  if ([string]::IsNullOrWhiteSpace($Token)) {
    throw "Token is empty. Use -Token for API smoke checks."
  }

  $baseUri = [System.Uri]$ApiBaseUrl
  $targetUri = [System.Uri]::new($baseUri, $Path.TrimStart('/'))
  $requestUri = $targetUri.AbsoluteUri

  if ($Query -and $Query.Count -gt 0) {
    $queryParts = foreach ($k in $Query.Keys) {
      "{0}={1}" -f [System.Uri]::EscapeDataString([string]$k), [System.Uri]::EscapeDataString([string]$Query[$k])
    }
    $builder = [System.UriBuilder]::new($targetUri)
    $builder.Query = [string]::Join("&", $queryParts)
    $requestUri = $builder.Uri.AbsoluteUri
  }

  $invokeArgs = @{
    Method  = $Method
    Uri     = $requestUri
    Headers = @{ "x-token" = $Token }
  }

  if ($Body -ne $null) {
    $invokeArgs["ContentType"] = "application/json"
    $invokeArgs["Body"] = ($Body | ConvertTo-Json -Depth 8)
  }

  return Invoke-RestMethod @invokeArgs
}

Run-Step -Name "Go tests: english_learning" -Script {
  Push-Location (Join-Path $RepoRoot "server")
  try {
    go test ./plugin/english_learning/...
    if ($LASTEXITCODE -ne 0) {
      throw "go test failed with exit code $LASTEXITCODE"
    }
  } finally {
    Pop-Location
  }
}

Run-Step -Name "Web lint: english plugin files" -Script {
  Push-Location (Join-Path $RepoRoot "web")
  try {
    npm run lint -- src/plugin/english_learning/api/english.js src/plugin/english_learning/view/word.vue src/plugin/english_learning/view/video.vue
    if ($LASTEXITCODE -ne 0) {
      throw "npm lint failed with exit code $LASTEXITCODE"
    }
  } finally {
    Pop-Location
  }
}

if (-not [string]::IsNullOrWhiteSpace($Token)) {
  Run-Step -Name "API smoke: category list" -Script {
    $resp = Invoke-EnglishApi -Method "GET" -Path "/englishLearning/content/getCategoryList" -Query @{ page = 1; pageSize = 5 }
    if ($resp.code -ne 0) {
      throw "code=$($resp.code), msg=$($resp.msg)"
    }
  }

  Run-Step -Name "API smoke: series list" -Script {
    $resp = Invoke-EnglishApi -Method "GET" -Path "/englishLearning/content/getVideoSeriesList" -Query @{ page = 1; pageSize = 5 }
    if ($resp.code -ne 0) {
      throw "code=$($resp.code), msg=$($resp.msg)"
    }
  }

  Run-Step -Name "API smoke: word list" -Script {
    $resp = Invoke-EnglishApi -Method "GET" -Path "/englishLearning/word/getWordList" -Query @{ page = 1; pageSize = 5 }
    if ($resp.code -ne 0) {
      throw "code=$($resp.code), msg=$($resp.msg)"
    }
  }

  Run-Step -Name "API smoke: entitlement list" -Script {
    $resp = Invoke-EnglishApi -Method "GET" -Path "/englishLearning/asset/getEntitlementList" -Query @{ page = 1; pageSize = 5 }
    if ($resp.code -ne 0) {
      throw "code=$($resp.code), msg=$($resp.msg)"
    }
  }

  if ($WordId -gt 0) {
    Run-Step -Name "API smoke: regenerate word audio" -Script {
      $resp = Invoke-EnglishApi -Method "POST" -Path "/englishLearning/word/regenerateAudio" -Body @{ ID = $WordId }
      if ($resp.code -ne 0) {
        throw "code=$($resp.code), msg=$($resp.msg)"
      }
    }
  }

  if ($EpisodeId -gt 0) {
    Run-Step -Name "API smoke: sentence list" -Script {
      $resp = Invoke-EnglishApi -Method "GET" -Path "/englishLearning/video/getSentenceList" -Query @{ episodeId = $EpisodeId }
      if ($resp.code -ne 0) {
        throw "code=$($resp.code), msg=$($resp.msg)"
      }
    }
  }
}

Write-Host "`n========== PT-5 Acceptance Results =========="
$script:results | Format-Table -AutoSize

$failed = @($script:results | Where-Object { -not $_.Passed })
if ($failed.Count -gt 0) {
  Write-Host "`nFailed steps: $($failed.Count)" -ForegroundColor Red
  exit 1
}

Write-Host "`nAll steps passed." -ForegroundColor Green
exit 0
