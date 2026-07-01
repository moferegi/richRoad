# PT完善 Step11 真实供应商切换门禁

更新时间: 2026-07-01

## 1. 本步目标
- 在已联调通过的基础上，增加“真实供应商切换门禁”。
- 明确区分：
  - 演练环境（允许 mock）
  - 生产环境（必须真实 provider + 有效 api_key）

## 2. 本步交付

### 2.1 新增门禁脚本
- 文件：`scripts/english_learning_pt11_provider_cutover_gate.ps1`
- 核心能力：
  - 校验 `learning_tts_*` 关键配置是否齐全。
  - 校验 `learning_tts_enabled=true`。
  - 调用 `preflightTTS` 确认链路可用。
  - 默认禁止本地 mock provider（`127.0.0.1/localhost`）和空 api_key。
  - 支持 `-AllowMockProvider` 放开演练环境门禁。

## 3. 本步验证结果

### 3.1 生产门禁口径（默认）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt11_provider_cutover_gate.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny"
```
- 结果：`Gate: FAIL`
- 失败原因：
  - provider 指向本地 mock（`127.0.0.1`）
  - `learning_tts_api_key` 为空或占位

### 3.2 演练口径（放开 mock）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt11_provider_cutover_gate.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny" -AllowMockProvider
```
- 结果：`Gate: PASS`

## 4. 结论
- PT-11 已完成：门禁脚本已把“能跑”与“可上生产”明确分层。
- 当前环境属于演练通过，但尚未满足真实供应商生产门禁。
- 下一步进入 PT-12：替换真实 provider 与 api_key，并让 PT-11 默认口径（不加 `-AllowMockProvider`）通过。
