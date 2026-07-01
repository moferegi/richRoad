param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8888",
  [string]$Token = "",
  [string]$ProbeWord = "destiny",
  [switch]$AllowMockProvider
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

$failures = @()
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
if ($missing.Count -gt 0) {
  $failures += "Missing required keys: $([string]::Join(', ', $missing))"
}

$enabled = [string]$keyMap["learning_tts_enabled"]
$providerUrl = [string]$keyMap["learning_tts_provider_url"]
$apiKey = [string]$keyMap["learning_tts_api_key"]

if ($enabled.ToLowerInvariant() -ne "true") {
  $failures += "learning_tts_enabled is not true"
}

if ([string]::IsNullOrWhiteSpace($providerUrl)) {
  $failures += "learning_tts_provider_url is empty"
}

$providerHost = ""
if (-not [string]::IsNullOrWhiteSpace($providerUrl)) {
  try {
    $providerHost = ([System.Uri]$providerUrl).Host
  }
  catch {
    $failures += "learning_tts_provider_url is not a valid URI"
  }
}

if (-not $AllowMockProvider) {
  if ($providerHost -eq "127.0.0.1" -or $providerHost -eq "localhost") {
    $failures += "Provider host is local mock ($providerHost), real provider is required"
  }

  if ([string]::IsNullOrWhiteSpace($apiKey) -or $apiKey -eq "******") {
    $failures += "learning_tts_api_key is empty or redacted placeholder"
  }
}

$preflightCode = -1
$preflightMsg = ""
$preflightData = $null
try {
  $preflightResp = Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/preflightTTS" -Body @{ word = $ProbeWord }
  $preflightCode = [int]$preflightResp.code
  $preflightMsg = [string]$preflightResp.msg
  $preflightData = $preflightResp.data
  if ($preflightCode -ne 0) {
    $failures += "preflightTTS failed: code=$preflightCode, msg=$preflightMsg"
  }
}
catch {
  $preflightMsg = $_.Exception.Message
  $failures += "preflightTTS HTTP error: $preflightMsg"
}

Write-Host "`n========== PT-11 Provider Cutover Gate =========="
Write-Host "AllowMockProvider: $AllowMockProvider"
Write-Host "ProviderUrl: $providerUrl"
Write-Host "ProviderHost: $providerHost"
Write-Host "TTS enabled: $enabled"
Write-Host "Preflight: code=$preflightCode, msg=$preflightMsg"
if ($preflightCode -eq 0 -and $preflightData) {
  Write-Host "US audio: $($preflightData.audioUs)"
  Write-Host "UK audio: $($preflightData.audioUk)"
}

if ($failures.Count -gt 0) {
  Write-Host "`nGate: FAIL" -ForegroundColor Red
  foreach ($f in $failures) {
    Write-Host " - $f"
  }
  exit 1
}

Write-Host "`nGate: PASS" -ForegroundColor Green
exit 0
