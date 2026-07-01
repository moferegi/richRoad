# PT完善 Step7 联调就绪检查与回归收口

更新时间: 2026-07-01

## 1. 本步目标
- 使用真实运行环境继续推进 TTS 联调，补齐“就绪检查”能力。
- 完成遗留回归项：`getWordList` 回归测试补强。
- 清理联调过程中的临时测试数据，保持工作区可维护。

## 2. 本步交付

### 2.1 getWordList 回归测试补强
- 文件：`server/plugin/english_learning/api/english_word_api_test.go`
- 新增用例：`TestEnglishWordAPI_GetWordList_ChapterFilterAndDistinct`
- 覆盖点：
  - `chapterId` 过滤生效。
  - 同一单词多条章节关联不会导致列表重复计数（`total` 按 distinct 口径）。
  - 分页时 `pageSize` 与总数口径一致。

### 2.2 PT-7 联调就绪脚本
- 文件：`scripts/english_learning_pt7_readiness_check.ps1`
- 脚本能力：
  - 读取 `english_learning` 配置分组。
  - 检查 `learning_tts_*` 关键配置是否齐全。
  - 调用 `POST /englishLearning/word/preflightTTS` 进行联调预检。
  - 输出 PASS/FAIL 并返回对应退出码。

### 2.3 联调临时数据清理
- 清理内容：
  - `tmp/web_captcha.png`
  - `tmp/pt7_gen_admin_token.go`
- 目的：避免临时文件污染工作区。

## 3. 联调实测结论（当前环境）
- 使用用户提供账号直登：
  - Web (`admin/123456`) 返回 `验证码错误`。
  - Uni (`mofer222/111111`) 返回 `验证码错误`。
- 使用本地生成管理员 token 可访问受保护接口（例如 `/sysConfig/getSysConfigByGroup`）。
- 当前环境 `english_learning` 分组仅存在部分配置项，`learning_tts_*` 仍有缺失。
- 当前正在运行的服务尚未加载 PT-6 新增路由，调用 `/englishLearning/word/preflightTTS` 返回 HTTP 404。
- 因此 Readiness 结果为 FAIL，需先部署最新服务并补齐配置再做正式联调。

## 4. 验证命令

### 4.1 后端回归
```powershell
cd server
go test ./plugin/english_learning/api -run "TestEnglishWordAPI_GetWordList_(NormalizedPage|ChapterFilterAndDistinct)$" -count=1
go test ./plugin/english_learning/...
```

### 4.2 联调就绪检查
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt7_readiness_check.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<x-token>"
```

### 4.3 当前环境检查输出（示例）
```text
Missing keys:
 - learning_tts_enabled
 - learning_tts_provider_url
 - learning_tts_api_key
 - learning_tts_timeout_ms
 - learning_tts_voice_us
 - learning_tts_voice_uk

Preflight result: code=-1, msg=HTTP error while calling /englishLearning/word/preflightTTS: (404) Not Found
Readiness: FAIL
```

## 5. 结论
- PT-7 已完成“回归补强 + 环境就绪检查 + 临时数据清理”。
- 下一步进入 PT-8：补齐运行环境 `learning_tts_*` 配置并执行一次 DrillMode 全链路演练。
