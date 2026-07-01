param(
  [string]$ApiBaseUrl = "http://127.0.0.1:8888",
  [string]$Token = "",
  [string]$ProbeWord = "destiny",
  [int]$WordId = 0,
  [string]$ProviderUrl = "",
  [string]$ApiKey = "",
  [int]$TimeoutMs = 0,
  [string]$VoiceUS = "",
  [string]$VoiceUK = "",
  [switch]$UseEnvSecrets,
  [string]$ProviderUrlEnvVar = "LEARNING_TTS_PROVIDER_URL",
  [string]$ApiKeyEnvVar = "LEARNING_TTS_API_KEY",
  [switch]$PromptApiKey,
  [switch]$SkipRollout,
  [switch]$AllowMockProvider,
  [switch]$ContinueOnFailure
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($Token)) {
  throw "Token is required. Use a super admin x-token for PT-12 cutover."
}

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$pt6Script = Join-Path $scriptDir "english_learning_pt6_tts_rollout.ps1"
$pt7Script = Join-Path $scriptDir "english_learning_pt7_readiness_check.ps1"
$pt10Script = Join-Path $scriptDir "english_learning_pt10_production_acceptance.ps1"
$pt11Script = Join-Path $scriptDir "english_learning_pt11_provider_cutover_gate.ps1"

foreach ($path in @($pt6Script, $pt7Script, $pt10Script, $pt11Script)) {
  if (-not (Test-Path $path)) {
    throw "Missing dependent script: $path"
  }
}

$results = @()
$hasFailure = $false

function Add-StepResult {
  param(
    [string]$Name,
    [string]$Status,
    [int]$ExitCode,
    [string]$Detail
  )

  $script:results += [PSCustomObject]@{
    Step     = $Name
    Status   = $Status
    ExitCode = $ExitCode
    Detail   = $Detail
  }
}

function Run-ChildScript {
  param(
    [string]$Name,
    [string]$ScriptPath,
    [hashtable]$NamedArgs
  )

  if ($script:hasFailure -and -not $ContinueOnFailure) {
    Add-StepResult -Name $Name -Status "SKIPPED" -ExitCode 0 -Detail "Skipped because a previous step failed"
    return
  }

  try {
    & $ScriptPath @NamedArgs
    $code = $LASTEXITCODE
    if ($code -ne 0) {
      throw "Child script exited with code $code"
    }

    Add-StepResult -Name $Name -Status "PASS" -ExitCode $code -Detail "OK"
  }
  catch {
    $code = if ($LASTEXITCODE -ne $null) { [int]$LASTEXITCODE } else { 1 }
    Add-StepResult -Name $Name -Status "FAIL" -ExitCode $code -Detail $_.Exception.Message
    $script:hasFailure = $true
  }
}

function Get-EnvValue {
  param(
    [string]$Name
  )

  if ([string]::IsNullOrWhiteSpace($Name)) {
    return ""
  }

  $value = [Environment]::GetEnvironmentVariable($Name)
  return [string]$value
}

function Read-SecretFromPrompt {
  param(
    [string]$Prompt
  )

  $secure = Read-Host -Prompt $Prompt -AsSecureString
  if ($null -eq $secure) {
    return ""
  }

  $ptr = [IntPtr]::Zero
  try {
    $ptr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
    return [Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr)
  }
  finally {
    if ($ptr -ne [IntPtr]::Zero) {
      [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr)
    }
  }
}

if (-not $SkipRollout) {
  $effectiveProviderUrl = $ProviderUrl
  $effectiveApiKey = $ApiKey

  if ($UseEnvSecrets) {
    if ([string]::IsNullOrWhiteSpace($effectiveProviderUrl)) {
      $effectiveProviderUrl = Get-EnvValue -Name $ProviderUrlEnvVar
    }
    if ([string]::IsNullOrWhiteSpace($effectiveApiKey)) {
      $effectiveApiKey = Get-EnvValue -Name $ApiKeyEnvVar
    }
  }

  if ($PromptApiKey -and [string]::IsNullOrWhiteSpace($effectiveApiKey)) {
    $effectiveApiKey = Read-SecretFromPrompt -Prompt "Enter real TTS API key"
    if ([string]::IsNullOrWhiteSpace($effectiveApiKey)) {
      throw "API key input is empty."
    }
  }

  $rolloutNamedArgs = @{
    ApiBaseUrl = $ApiBaseUrl
    Token      = $Token
    ProbeWord  = $ProbeWord
    EnableTTS  = $true
  }

  if (-not [string]::IsNullOrWhiteSpace($effectiveProviderUrl)) {
    $rolloutNamedArgs.ProviderUrl = $effectiveProviderUrl.Trim()
  }
  if (-not [string]::IsNullOrWhiteSpace($effectiveApiKey)) {
    $rolloutNamedArgs.ApiKey = $effectiveApiKey
  }
  if ($TimeoutMs -gt 0) {
    $rolloutNamedArgs.TimeoutMs = $TimeoutMs
  }
  if (-not [string]::IsNullOrWhiteSpace($VoiceUS)) {
    $rolloutNamedArgs.VoiceUS = $VoiceUS.Trim()
  }
  if (-not [string]::IsNullOrWhiteSpace($VoiceUK)) {
    $rolloutNamedArgs.VoiceUK = $VoiceUK.Trim()
  }

  Run-ChildScript -Name "PT-6 rollout" -ScriptPath $pt6Script -NamedArgs $rolloutNamedArgs
}
else {
  Add-StepResult -Name "PT-6 rollout" -Status "SKIPPED" -ExitCode 0 -Detail "SkipRollout was set"
}

$gateNamedArgs = @{
  ApiBaseUrl = $ApiBaseUrl
  Token      = $Token
  ProbeWord  = $ProbeWord
}
if ($AllowMockProvider) {
  $gateNamedArgs.AllowMockProvider = $true
}
Run-ChildScript -Name "PT-11 provider gate" -ScriptPath $pt11Script -NamedArgs $gateNamedArgs

$readinessNamedArgs = @{
  ApiBaseUrl = $ApiBaseUrl
  Token      = $Token
  ProbeWord  = $ProbeWord
}
Run-ChildScript -Name "PT-7 readiness" -ScriptPath $pt7Script -NamedArgs $readinessNamedArgs

$acceptanceNamedArgs = @{
  ApiBaseUrl = $ApiBaseUrl
  Token      = $Token
  ProbeWord  = $ProbeWord
}
if ($WordId -gt 0) {
  $acceptanceNamedArgs.WordId = $WordId
}
Run-ChildScript -Name "PT-10 acceptance" -ScriptPath $pt10Script -NamedArgs $acceptanceNamedArgs

Write-Host "`n========== PT-12 Final Cutover Summary =========="
$results | Format-Table -AutoSize

$failed = @($results | Where-Object { $_.Status -eq "FAIL" })
if ($failed.Count -gt 0) {
  Write-Host "`nPT-12: FAIL ($($failed.Count) step(s) failed)." -ForegroundColor Red
  exit 1
}

Write-Host "`nPT-12: PASS" -ForegroundColor Green
exit 0
