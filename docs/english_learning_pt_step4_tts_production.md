# PT完善 Step4 TTS生产链路定稿与运营闭环

更新时间: 2026-07-01

## 1. 本步目标
- 解决 R18：将英语单词发音从占位URL升级为可配置、可落地的真实TTS调用链路。
- 提供运营可执行的“重生成发音”能力，支持对历史词条补音频。
- 完成本步可验收闭环：后端实现、Web入口、自动化测试与lint/测试通过。

## 2. 本步交付

### 2.1 后端TTS链路生产化
- 文件：`server/plugin/english_learning/service/english_word.go`
- 关键改造：
  - 移除原先占位的 `https://oss.example.com/...` 假URL生成逻辑。
  - 新增可配置TTS调用流程：
    - 配置开关：`learning_tts_enabled`
    - 服务地址：`learning_tts_provider_url`
    - 服务密钥：`learning_tts_api_key`
    - 超时配置：`learning_tts_timeout_ms`
    - 音色配置：`learning_tts_voice_us` / `learning_tts_voice_uk`
  - 支持从多种响应结构提取音频URL：
    - `url`
    - `audioUrl`
    - `data.url`
    - `data.audioUrl`
- 运行行为：
  - 创建单词时，若 `audioUs/audioUk` 为空，按配置尝试自动生成。
  - 未开启或未配置服务时，不再写入假链接。

### 2.2 单词发音手动重生成能力
- 新增请求模型：`server/plugin/english_learning/model/request/english.go`
  - `RegenerateWordAudioReq`
- 新增服务方法：
  - `RegenerateWordAudio`
- 新增API：`POST /englishLearning/word/regenerateAudio`
  - 文件：`server/plugin/english_learning/api/english_word.go`
- 新增路由：
  - 文件：`server/plugin/english_learning/router/english_word.go`

### 2.3 Web运营端入口补齐
- 文件：`web/src/plugin/english_learning/view/word.vue`
  - 单词操作列新增“重生发音”按钮，支持运营一键重生成。
- 文件：`web/src/plugin/english_learning/api/english.js`
  - 新增 `regenerateWordAudio` API 封装。

### 2.4 默认配置初始化
- 文件：`server/initialize/other.go`
- english_learning 分组新增6个配置项：
  - `learning_tts_enabled`
  - `learning_tts_provider_url`
  - `learning_tts_api_key`
  - `learning_tts_timeout_ms`
  - `learning_tts_voice_us`
  - `learning_tts_voice_uk`

## 3. 测试与验证

### 3.1 新增测试
- `server/plugin/english_learning/service/english_word_tts_test.go`
  - 覆盖音频URL提取与TTS调用成功路径。
- `server/plugin/english_learning/api/english_word_api_test.go`
  - 补充 `regenerateAudio` 参数错误/成功用例。

### 3.2 回归结果
- 后端：`cd server && go test ./plugin/english_learning/...` 通过。
- 前端：`cd web && npm run lint -- src/plugin/english_learning/api/english.js src/plugin/english_learning/view/word.vue src/plugin/english_learning/view/video.vue` 通过。

## 4. 结论
- PT-4 已完成：TTS链路从“假数据占位”升级为“可配置真实调用 + 运营可手工补发音”。
- 英语学习运营后台已具备内容维护与发音生成双闭环能力。
- 可进入 PT-5（测试矩阵扩展与端到端验收脚本）。
