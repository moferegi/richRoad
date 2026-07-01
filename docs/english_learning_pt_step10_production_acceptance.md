# PT完善 Step10 生产验收脚本与切换说明

更新时间: 2026-07-01

## 1. 本步目标
- 给英语学习 TTS 链路提供可复用的“生产前验收脚本”。
- 让联调、预检、验收形成标准化命令，降低人工操作误差。

## 2. 本步交付

### 2.1 新增生产验收脚本
- 文件：`scripts/english_learning_pt10_production_acceptance.ps1`
- 验收项：
  - `english_learning` 配置组可读。
  - `preflightTTS` 默认（US+UK）可用且音频地址不为空。
  - `preflightTTS` US-only 可用。
  - `preflightTTS` UK-only 可用。
  - 可选：指定 `WordId` 时执行 `regenerateAudio`。
- 结果输出：
  - 表格化步骤结果
  - 最终 `Acceptance: PASS/FAIL`

### 2.2 执行命令
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt10_production_acceptance.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny"
```

可选（验收重生成能力）：
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt10_production_acceptance.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<admin-token>" -ProbeWord "destiny" -WordId 1
```

## 3. 切换真实供应商说明
- 当前基线 provider 为本地 mock：`http://127.0.0.1:19090/tts`。
- 切换真实供应商时，建议顺序：
  1. 使用 `scripts/english_learning_pt6_tts_rollout.ps1` 更新 `provider_url/api_key/timeout/voice`。
  2. 运行 `scripts/english_learning_pt7_readiness_check.ps1`。
  3. 运行 `scripts/english_learning_pt10_production_acceptance.ps1`。
- 若第 2/3 步失败，使用 PT-6 脚本的回滚能力恢复。

## 4. 结论
- PT-10 已完成：生产前验收动作脚本化，形成可复用的上线前标准流程。
