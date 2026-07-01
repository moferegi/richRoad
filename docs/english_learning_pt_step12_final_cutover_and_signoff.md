# PT完善 Step12 最终切换与结项验收

更新时间: 2026-07-01

## 1. 本步目标
- 将 PT-6、PT-11、PT-7、PT-10 串联为一个总控流程，形成最终切换与结项验收入口。
- 明确生产口径与演练口径的通过标准，防止 mock 配置误判上线。

## 2. 本步交付

### 2.1 新增总控脚本
- 文件：`scripts/english_learning_pt12_final_cutover.ps1`
- 编排顺序：
  1. PT-6 rollout（必要时更新配置并预检）
  2. PT-11 provider gate（真实供应商门禁）
  3. PT-7 readiness（就绪检查）
  4. PT-10 acceptance（生产验收）
- 默认策略：某一步失败后，后续步骤标记 `SKIPPED`，并以失败退出。
- 可选参数：
  - `-ProviderUrl/-ApiKey/-TimeoutMs/-VoiceUS/-VoiceUK` 用于一次性切换配置
  - `-UseEnvSecrets/-ProviderUrlEnvVar/-ApiKeyEnvVar` 支持从环境变量读取 provider 与 key
  - `-PromptApiKey` 支持在终端内保密输入 key（不回显）
  - `-SkipRollout` 跳过 PT-6
  - `-AllowMockProvider` 放开门禁用于演练环境
  - `-ContinueOnFailure` 失败后继续执行以收集完整诊断

### 2.2 安全密钥注入模式（推荐）
- 适用场景：需要执行真实切换，但不希望把 `ApiKey` 明文写入命令行历史。

```powershell
$env:LEARNING_TTS_PROVIDER_URL = "https://your-real-tts-provider.example.com/tts"
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt12_final_cutover.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny" -UseEnvSecrets -PromptApiKey
```

- 上述命令会：
  - 从环境变量读取 provider URL
  - 在终端里提示输入真实 API key（不回显）
  - 自动串联 PT-6/11/7/10 完成最终校验

## 3. 实跑结果

### 3.1 严格生产口径（默认）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt12_final_cutover.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny"
```
- 结果：`PT-12: FAIL`
- 步骤摘要：
  - PT-6 rollout: PASS
  - PT-11 provider gate: FAIL
  - PT-7 readiness: SKIPPED
  - PT-10 acceptance: SKIPPED
- 失败原因：
  - provider 仍为本地 mock（`127.0.0.1`）
  - `learning_tts_api_key` 为空或占位

### 3.2 演练口径（放开 mock）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt12_final_cutover.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny" -AllowMockProvider
```
- 结果：`PT-12: PASS`
- 步骤摘要：
  - PT-6 rollout: PASS
  - PT-11 provider gate: PASS
  - PT-7 readiness: PASS
  - PT-10 acceptance: PASS

### 3.3 脚本增强回归（2026-07-01）
- 在新增 `-UseEnvSecrets/-PromptApiKey` 后，复跑演练口径：`PT-12: PASS`。
- 说明：增强仅影响密钥注入方式，不改变既有门禁判定逻辑。

## 4. 生产放行门槛
- 必须满足：
  - `learning_tts_provider_url` 指向真实供应商域名（非 `127.0.0.1/localhost`）
  - `learning_tts_api_key` 为有效真实密钥
  - PT-12 默认口径（不加 `-AllowMockProvider`）执行结果为 PASS

## 5. 结论
- Step12 自动化编排与验收闭环已完成。
- 当前环境仍是 mock 基线，因此生产口径尚未放行。
- 下一动作仅剩“注入真实供应商配置并重跑 PT-12 默认口径”。
