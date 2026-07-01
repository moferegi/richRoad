# PT完善 Step9 非Drill正式联调与就绪PASS

更新时间: 2026-07-01

## 1. 本步目标
- 基于 PT-8 的服务切换结果，执行一次非 DrillMode 的正式联调。
- 保留成功配置并复测 PT-7 就绪检查，目标是 `Readiness: PASS`。

## 2. 本步执行

### 2.1 非 DrillMode 正式联调
- 命令：
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt6_tts_rollout.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -EnableTTS -ProviderUrl "http://127.0.0.1:19090/tts" -TimeoutMs 8000 -VoiceUS "en-US-JennyNeural" -VoiceUK "en-GB-SoniaNeural" -ProbeWord "destiny"
```
- 结果：
  - `Status: PREFLIGHT_OK`
  - US/UK 音频均成功返回
  - 配置保持生效（`RolledBack=False`）

### 2.2 就绪检查复测
- 命令：
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt7_readiness_check.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>"
```
- 结果：
  - `Missing keys: none`
  - `Preflight result: code=0`
  - `Readiness: PASS`

### 2.3 联调后关键配置
- `learning_tts_enabled = true`
- `learning_tts_provider_url = http://127.0.0.1:19090/tts`
- `learning_tts_timeout_ms = 8000`
- `learning_tts_voice_us = en-US-JennyNeural`
- `learning_tts_voice_uk = en-GB-SoniaNeural`

## 3. 说明
- 本步使用本地 mock TTS 服务作为供应商替身，目的是确认配置保留与正式联调链路可用。
- 切到真实供应商时，仅需替换 `provider_url` 与 `api_key` 并复跑 PT-7 就绪检查。

## 4. 结论
- PT-9 已完成：正式联调成功且配置保留，环境已达到就绪 PASS 状态。
- 可进入 PT-10：替换真实供应商配置并执行生产前验收清单。
