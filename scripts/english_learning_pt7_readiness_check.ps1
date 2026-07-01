param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8888",
  [string]$Token = "",
  [string]$ProbeWord = "destiny"
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

  return Invoke-RestMethod @invokeArgs
}

$requiredKeys = @(
  "learning_tts_enabled",
  "learning_tts_provider_url",
  "learning_tts_api_key",
  "learning_tts_timeout_ms",
  "learning_tts_voice_us",
  "learning_tts_voice_uk"
)

$configResp = Invoke-JsonApi -Method "GET" -Path "/sysConfig/getSysConfigByGroup" -Query @{ configGroup = "english_learning" }
if ([int]$configResp.code -ne 0) {
  throw "Failed to query english_learning config group: code=$($configResp.code), msg=$($configResp.msg)"
}

$keyMap = @{}
foreach ($item in @($configResp.data)) {
  $keyMap[[string]$item.configKey] = [string]$item.configValue
}

$missing = @($requiredKeys | Where-Object { -not $keyMap.ContainsKey($_) })

$preflightResp = $null
$preflightCode = -1
$preflightMsg = ""
try {
  $preflightResp = Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/preflightTTS" -Body @{ word = $ProbeWord }
  $preflightCode = [int]$preflightResp.code
  $preflightMsg = [string]$preflightResp.msg
}
catch {
  $preflightMsg = "HTTP error while calling /englishLearning/word/preflightTTS: $($_.Exception.Message)"
}

Write-Host "`n========== PT-7 Readiness Check =========="
Write-Host "Required keys count: $($requiredKeys.Count)"
Write-Host "Configured keys count: $($keyMap.Keys.Count)"

if ($missing.Count -gt 0) {
  Write-Host "Missing keys:" -ForegroundColor Yellow
  foreach ($k in $missing) {
    Write-Host " - $k"
  }
} else {
  Write-Host "Missing keys: none"
}

Write-Host "`nPreflight result: code=$preflightCode, msg=$preflightMsg"
if ($preflightCode -eq 0 -and $preflightResp -and $preflightResp.data) {
  Write-Host "US audio: $($preflightResp.data.audioUs)"
  Write-Host "UK audio: $($preflightResp.data.audioUk)"
  Write-Host "Duration: $($preflightResp.data.durationMs) ms"
}

if ($missing.Count -eq 0 -and $preflightCode -eq 0) {
  Write-Host "`nReadiness: PASS" -ForegroundColor Green
  exit 0
}

Write-Host "`nReadiness: FAIL" -ForegroundColor Red
exit 1
