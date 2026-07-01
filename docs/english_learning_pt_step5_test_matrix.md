# PT完善 Step5 测试矩阵扩展与端到端验收脚本

更新时间: 2026-07-01

## 1. 本步目标
- 把英语学习 PT-1 ~ PT-4 的能力整理成可执行、可复用的验收矩阵。
- 在后端补齐 TTS 失败路径测试，避免“仅成功路径通过”的假稳定。
- 提供一键验收脚本，降低回归成本并固化交付门槛。

## 2. 本步交付

### 2.1 测试矩阵扩展
- 文件：`server/plugin/english_learning/service/english_word_tts_test.go`
- 新增覆盖：
  - TTS 开关关闭：`TestRequestTTSAudioURL_Disabled`
  - TTS 服务非 2xx：`TestRequestTTSAudioURL_HTTPStatusError`
  - TTS 响应缺失音频 URL：`TestRequestTTSAudioURL_MissingAudioURL`
- 既有成功路径保持：
  - `TestRequestTTSAudioURL_Success`
  - `TestExtractTTSAudioURL`

### 2.2 API 负向回归补齐
- 文件：`server/plugin/english_learning/api/english_word_api_test.go`
- 新增用例：
  - `TestEnglishWordAPI_RegenerateWordAudio_BothDisabled`
- 作用：当 `regenerateUs=false` 且 `regenerateUk=false` 时，接口必须失败，防止无效请求被误判成功。

### 2.3 一键验收脚本
- 文件：`scripts/english_learning_pt5_acceptance.ps1`
- 脚本能力：
  - 执行后端测试：`go test ./plugin/english_learning/...`
  - 执行前端 lint：`npm run lint -- src/plugin/english_learning/api/english.js src/plugin/english_learning/view/word.vue src/plugin/english_learning/view/video.vue`
  - 可选 API smoke（传入 Token 后启用）：
    - 分类列表
    - 剧集列表
    - 单词列表
    - 授权列表
    - 可选单词发音重生成（需 `-WordId`）
    - 可选字幕列表（需 `-EpisodeId`）

## 3. 验收命令

### 3.1 基础验收（无 Token）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt5_acceptance.ps1
```

### 3.2 含 API smoke 验收（有 Token）
```powershell
powershell -ExecutionPolicy Bypass -File scripts/english_learning_pt5_acceptance.ps1 -ApiBaseUrl "http://127.0.0.1:8888" -Token "<x-token>" -WordId 1 -EpisodeId 1
```

## 4. 通过标准
- 后端插件测试全部通过。
- 前端英语运营端目标文件 lint 全部通过。
- （可选）API smoke 全部返回 `code=0`。
- 新增 TTS 失败路径测试全部通过。

## 5. 结论
- PT-5 已完成：测试矩阵从“功能可用”提升到“成功+失败路径可验证”，并落地一键验收脚本。
- 英语学习模块已具备持续回归与发布前快速自检能力。
- 可进入 PT-6（生产联调演练：真实 TTS 服务接入、配置上线与回滚演练）。
