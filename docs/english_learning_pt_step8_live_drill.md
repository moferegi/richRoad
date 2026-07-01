# PT完善 Step8 服务切换与实网联调演练

更新时间: 2026-07-01

## 1. 本步目标
- 将运行中的后端服务切换为包含 PT-6/PT-7 新能力的最新代码。
- 验证 `preflightTTS` 路由与 `learning_tts_*` 配置在运行环境中可用。
- 完成一次 DrillMode 联调演练并确认自动回滚生效。

## 2. 本步执行

### 2.1 服务切换
- 停止旧进程（8888）：PID `35616`（临时 `go run` 产物）。
- 启动新进程（8888）：PID `2616`。
- 验证结果：
  - `/englishLearning/word/preflightTTS` 已可访问（不再 404）。
  - 在未开启 TTS 时返回业务错误：`TTS自动生成未开启`（符合预期）。

### 2.2 配置状态确认
- 调用 `GET /sysConfig/getSysConfigByGroup?configGroup=english_learning`。
- 已存在 TTS 关键配置项：
  - `learning_tts_enabled`
  - `learning_tts_provider_url`
  - `learning_tts_api_key`
  - `learning_tts_timeout_ms`
  - `learning_tts_voice_us`
  - `learning_tts_voice_uk`

### 2.3 DrillMode 联调演练
- 本地启动 mock TTS 服务：`http://127.0.0.1:19090/tts`（返回固定 `audioUrl`）。
- 执行命令：
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt6_tts_rollout.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -EnableTTS -ProviderUrl "http://127.0.0.1:19090/tts" -TimeoutMs 8000 -VoiceUS "en-US-JennyNeural" -VoiceUK "en-GB-SoniaNeural" -ProbeWord "destiny" -DrillMode
```
- 实际输出：
  - `Status: PREFLIGHT_OK_ROLLED_BACK`
  - `US audio: http://127.0.0.1:19090/audio/mock.mp3`
  - `UK audio: http://127.0.0.1:19090/audio/mock.mp3`
  - 变更项已回滚：
    - `learning_tts_enabled: false -> true -> false`
    - `learning_tts_provider_url: "" -> mock_url -> ""`

## 3. 风险与说明
- 本次演练使用本地 mock TTS，验证的是链路与回滚机制，不代表真实供应商 SLA。
- Web/Uni 账号密码登录仍受验证码策略约束，自动化联调应继续使用管理员 token 或可控验证码链路。

## 4. 结论
- PT-8 已完成：服务已切换到最新代码，关键配置已补齐，DrillMode 演练通过且回滚生效。
- 可进入 PT-9：接入真实供应商地址与密钥，执行一次“非 DrillMode”正式联调并记录回放。
