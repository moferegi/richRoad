param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8888",
  [string]$Token = "",
  [string]$ProbeWord = "destiny",
  [string]$ProviderUrl = "",
  [string]$ApiKey = "",
  [int]$TimeoutMs = 0,
  [string]$VoiceUS = "",
  [string]$VoiceUK = "",
  [switch]$EnableTTS,
  [switch]$DisableTTS,
  [switch]$DrillMode,
  [bool]$AutoRollback = $true
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($Token)) {
  throw "Token is required. Use a super admin x-token for rollout and rollback."
}

if ($EnableTTS -and $DisableTTS) {
  throw "EnableTTS and DisableTTS cannot be used together."
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
    $msg = [string]$resp.msg
    throw "API failed ($Method $Path): code=$($resp.code), msg=$msg"
  }

  return $resp
}

function Get-EnglishLearningConfigMap {
  $resp = Invoke-JsonApi -Method "GET" -Path "/sysConfig/getSysConfigByGroup" -Query @{ configGroup = "english_learning" }
  $map = @{}
  foreach ($row in @($resp.data)) {
    $map[[string]$row.configKey] = $row
  }
  return $map
}

function Update-SysConfigValue {
  param(
    [object]$ConfigRow,
    [string]$NewValue,
    [string]$Remark
  )

  $body = @{
    id          = [int]$ConfigRow.id
    configValue = $NewValue
    remark      = $Remark
  }
  [void](Invoke-JsonApi -Method "PUT" -Path "/sysConfig/updateSysConfig" -Body $body)
}

$ttsKeys = @(
  "learning_tts_enabled",
  "learning_tts_provider_url",
  "learning_tts_api_key",
  "learning_tts_timeout_ms",
  "learning_tts_voice_us",
  "learning_tts_voice_uk"
)

$configMap = Get-EnglishLearningConfigMap
foreach ($k in $ttsKeys) {
  if (-not $configMap.ContainsKey($k)) {
    throw "Missing sys config key in english_learning group: $k"
  }
}

$desired = @{}
if ($EnableTTS) {
  $desired["learning_tts_enabled"] = "true"
}
if ($DisableTTS) {
  $desired["learning_tts_enabled"] = "false"
}
if (-not [string]::IsNullOrWhiteSpace($ProviderUrl)) {
  $desired["learning_tts_provider_url"] = $ProviderUrl.Trim()
}
if (-not [string]::IsNullOrWhiteSpace($ApiKey)) {
  $desired["learning_tts_api_key"] = $ApiKey
}
if ($TimeoutMs -gt 0) {
  $desired["learning_tts_timeout_ms"] = [string]$TimeoutMs
}
if (-not [string]::IsNullOrWhiteSpace($VoiceUS)) {
  $desired["learning_tts_voice_us"] = $VoiceUS.Trim()
}
if (-not [string]::IsNullOrWhiteSpace($VoiceUK)) {
  $desired["learning_tts_voice_uk"] = $VoiceUK.Trim()
}

if ($desired.ContainsKey("learning_tts_api_key")) {
  $currentApiKeyValue = [string]$configMap["learning_tts_api_key"].configValue
  if ($currentApiKeyValue -eq "******") {
    throw "Current API key is redacted. Use a super admin token to safely snapshot and rollback API key."
  }
}

$snapshot = @{}
$changed = @()
$status = "INIT"
$preflightData = $null

try {
  foreach ($entry in $desired.GetEnumerator()) {
    $key = [string]$entry.Key
    $newValue = [string]$entry.Value
    $row = $configMap[$key]
    $oldValue = [string]$row.configValue

    $snapshot[$key] = $oldValue
    if ($oldValue -eq $newValue) {
      continue
    }

    Update-SysConfigValue -ConfigRow $row -NewValue $newValue -Remark "PT-6 rollout rehearsal"
    $changed += [PSCustomObject]@{
      Key      = $key
      OldValue = $oldValue
      NewValue = $newValue
      RolledBack = $false
    }
  }

  $preflightBody = @{ word = $ProbeWord }
  $preflightResp = Invoke-JsonApi -Method "POST" -Path "/englishLearning/word/preflightTTS" -Body $preflightBody
  $preflightData = $preflightResp.data
  $status = "PREFLIGHT_OK"

  if ($DrillMode -and $changed.Count -gt 0) {
    foreach ($item in $changed) {
      $row = $configMap[[string]$item.Key]
      Update-SysConfigValue -ConfigRow $row -NewValue ([string]$item.OldValue) -Remark "PT-6 drill rollback"
      $item.RolledBack = $true
    }
    $status = "PREFLIGHT_OK_ROLLED_BACK"
  }
}
catch {
  $status = "FAILED"
  $failureMessage = $_.Exception.Message

  if ($AutoRollback -and $changed.Count -gt 0) {
    foreach ($item in $changed) {
      try {
        $row = $configMap[[string]$item.Key]
        Update-SysConfigValue -ConfigRow $row -NewValue ([string]$item.OldValue) -Remark "PT-6 auto rollback"
        $item.RolledBack = $true
      }
      catch {
        Write-Host "Rollback failed for $($item.Key): $($_.Exception.Message)" -ForegroundColor Yellow
      }
    }
  }

  Write-Host "PT-6 rollout failed: $failureMessage" -ForegroundColor Red
  if ($changed.Count -gt 0) {
    Write-Host "Changed keys:" -ForegroundColor Yellow
    $changed | Format-Table -AutoSize
  }
  exit 1
}

Write-Host "`n========== PT-6 Rollout Result =========="
Write-Host "Status: $status"
if ($preflightData) {
  Write-Host "Probe word: $($preflightData.word)"
  Write-Host "US audio:  $($preflightData.audioUs)"
  Write-Host "UK audio:  $($preflightData.audioUk)"
  Write-Host "Duration:  $($preflightData.durationMs) ms"
}

if ($changed.Count -gt 0) {
  Write-Host "`nChanged keys:"
  $changed | Format-Table -AutoSize
} else {
  Write-Host "`nNo config changes were applied."
}

exit 0
