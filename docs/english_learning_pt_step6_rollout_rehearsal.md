# PT完善 Step6 生产联调演练与配置回滚闭环

更新时间: 2026-07-01

## 1. 本步目标
- 进入真实 TTS 供应商联调阶段，提供上线前可执行预检。
- 将配置上线动作标准化为脚本流程，并支持失败自动回滚。
- 固化演练文档与回归测试，确保 PT-6 可复现、可审计。

## 2. 本步交付

### 2.1 新增 TTS 预检接口
- 文件：`server/plugin/english_learning/api/english_word.go`
- 新增接口：`POST /englishLearning/word/preflightTTS`
- 能力说明：
  - 支持传入预检单词（默认 `destiny`）。
  - 支持按需选择美式/英式发音预检。
  - 返回音色、音频 URL 与总耗时。

### 2.2 新增服务层预检逻辑
- 文件：`server/plugin/english_learning/service/english_word.go`
- 新增能力：`PreflightTTS`
- 行为说明：
  - 复用生产配置读取逻辑（`learning_tts_*`）。
  - 复用真实 HTTP 调用链路，不走占位逻辑。
  - 对 US/UK 分开发音进行逐项校验并输出详细错误。

### 2.3 路由与请求模型补齐
- 文件：
  - `server/plugin/english_learning/router/english_word.go`
  - `server/plugin/english_learning/model/request/english.go`
- 结果：预检接口完成路由注册，且具备标准 DTO。

### 2.4 上线演练脚本（含自动回滚）
- 文件：`scripts/english_learning_pt6_tts_rollout.ps1`
- 核心能力：
  - 拉取 `english_learning` 分组的 TTS 配置快照。
  - 按参数更新 `learning_tts_*` 配置。
  - 触发 `preflightTTS` 接口执行真实联调。
  - 失败时自动回滚（可关闭）。
  - 演练模式（`-DrillMode`）下：即使预检成功也回滚。
- 关键参数：
  - `-Token`（必须，建议超管 token）
  - `-EnableTTS` / `-DisableTTS`
  - `-ProviderUrl` / `-ApiKey` / `-TimeoutMs`
  - `-VoiceUS` / `-VoiceUK`
  - `-ProbeWord`
  - `-DrillMode`

### 2.5 Web API 封装补齐
- 文件：`web/src/plugin/english_learning/api/english.js`
- 新增方法：`preflightWordTTS`

## 3. PT-6 演练命令

### 3.1 只做预检，不改配置
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt6_tts_rollout.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<x-token>"
```

### 3.2 演练模式（成功后自动回滚）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt6_tts_rollout.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<x-token>" -EnableTTS -ProviderUrl "https://tts-provider.example.com/speech" -ApiKey "<api-key>" -TimeoutMs 8000 -VoiceUS "en-US-JennyNeural" -VoiceUK "en-GB-SoniaNeural" -ProbeWord "destiny" -DrillMode
```

### 3.3 正式上线（不自动回滚成功结果）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt6_tts_rollout.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<x-token>" -EnableTTS -ProviderUrl "https://tts-provider.example.com/speech" -ApiKey "<api-key>" -TimeoutMs 8000
```

## 4. 测试与验证
- 新增服务测试：
  - `TestPreflightTTS_Success`
  - `TestPreflightTTS_BothDisabled`
- 新增 API 测试：
  - `TestEnglishWordAPI_PreflightTTS_Success`
  - `TestEnglishWordAPI_PreflightTTS_BothDisabled`
- 回归要求：
  - `go test ./plugin/english_learning/...`
  - `npm run lint -- src/plugin/english_learning/api/english.js src/plugin/english_learning/view/word.vue src/plugin/english_learning/view/video.vue`

## 5. 结论
- PT-6 已把“真实联调 + 上线变更 + 回滚保护”做成标准化流程。
- 英语学习 TTS 链路具备生产前演练能力，可直接进入实网供应商联调窗口。
