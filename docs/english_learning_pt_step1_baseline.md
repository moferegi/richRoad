# PT完善 Step1 基线冻结与验收矩阵

更新时间: 2026-07-01

## 1. 基线信息
- 代码基线: 5f63daf9
- 工作区状态: dirty（已存在多处并行改动，本文仅冻结英语学习相关能力现状）
- 冻结范围:
  - 后端: `server/plugin/english_learning/**`
  - 初始化与菜单: `server/initialize/new_modules_init.go`, `server/initialize/other.go`, `server/plugin/register.go`
  - Web: `web/src/plugin/english_learning/**`, `web/src/utils/asyncRouter.js`
  - Uni: `uni/src/pages/learning/**`, `uni/src/api/learning.js`, `uni/src/pages.json`

## 2. 模块盘点
- 后端英语插件: 已具备 API/Service/Router/Model/Init/Gorm/Test 分层结构（共 42 个文件）。
- Uni 学习端: 已有 11 个学习页面（首页、跟打、播放器、视频详情、我的、错词本、收藏、观看记录、签到记录、积分流水、免费时长流水）。
- Web 英语管理端: 仅 2 个页面 + 1 个 API 文件，仍是轻量原型态。

## 3. 18项需求映射矩阵（PT-1 冻结）
说明: 以下 R01-R18 按本次会话中用户 18 条需求语义归并，作为后续 PT-2 及以后步骤的验收基准。

| ID | 需求项 | 状态 | 证据 | 缺口/风险 |
| --- | --- | --- | --- | --- |
| R01 | 英语学习模块独立，不复用商城分类/SKU | 已完成 | `server/plugin/english_learning/plugin.go`, `server/plugin/english_learning/model/english_model.go:8`, `server/plugin/english_learning/model/video_model.go:8` | 无 |
| R02 | 英语学习核心模型完整（分类/章节/单词/句子/视频分类/系列/单集/字幕/句子） | 已完成 | `server/plugin/english_learning/model/english_model.go:8`, `server/plugin/english_learning/model/video_model.go:8` | 无 |
| R03 | Uni 学习主导航与页面闭环（Home/Typing/My） | 已完成 | `uni/src/pages.json:487`, `uni/src/pages.json:494`, `uni/src/pages.json:500` | 无 |
| R04 | Web 后台可见英语菜单（单词/视频） | 已完成 | `server/initialize/new_modules_init.go:190`, `server/initialize/new_modules_init.go:191`, `web/src/utils/asyncRouter.js:2` | 依赖初始化执行与角色菜单关系写入 |
| R05 | Web 运营端完整内容管理（分类、章节、词汇、视频全链路） | 部分完成 | `web/src/plugin/english_learning/view/word.vue:62`, `web/src/plugin/english_learning/view/video.vue:36`, `web/src/plugin/english_learning/api/english.js:4` | 仅覆盖“建词+字幕解析”原型，缺完整 CRUD/筛选/关联管理 |
| R06 | 单词跟打默认/恢复进度逻辑 | 已完成 | `uni/src/pages/learning/typing.vue`（已接 `getWordProgress`/`saveWordProgress`） | 仍缺按“用户单独授权分类”过滤能力 |
| R07 | 视频试看比例（默认8%）与中断提示 | 已完成 | `server/plugin/english_learning/model/video_model.go:33`, `uni/src/pages/learning/player.vue:153`, `uni/src/pages/learning/player.vue:442` | 当前主要在前端播放层执行，需后续加强服务端强约束 |
| R08 | 后端下发当前用户视频权限态（hasFullAuth） | 已完成 | `server/plugin/english_learning/api/content.go:593`, `server/plugin/english_learning/api/content.go:599` | 仅下发全局态，未到资源粒度 |
| R09 | 权限口径包含免费期/免费分钟/VIP | 部分完成 | `server/plugin/english_learning/middleware/learning_auth.go:49`, `server/plugin/english_learning/middleware/learning_auth.go:54`, `server/plugin/english_learning/middleware/learning_auth.go:60` | 仍是“全局资产态”，未融合资源收费与用户专属授权 |
| R10 | 新用户默认全站免费期（如24小时）自动赋值 | 未完成 | `server/plugin/english_learning/model/user_learning_model.go:15`, `server/plugin/english_learning/middleware/learning_auth.go:49` | 仅有字段与判定，未发现注册时写入 `FreeTimeExpire` 的链路 |
| R11 | 用户粒度授权（分类/系列/单集）能力 | 未完成 | 模块内未发现专属授权模型/API（仅 `Price/NeedVip/TrialPercent`） | 无法满足“给某用户单独开权限”运营诉求 |
| R12 | 心跳可信扣时长、防刷限流与进度校验 | 已完成 | `server/plugin/english_learning/service/user_learning_asset.go:23`, `server/plugin/english_learning/service/user_learning_asset.go:60`, `server/plugin/english_learning/service/user_learning_asset.go:242` | 需补充压测与异常链路观测 |
| R13 | 视频防盗链签名与播放续签 | 已完成 | `server/plugin/english_learning/service/video_security.go:12`, `server/plugin/english_learning/service/video_security.go:14`, `uni/src/pages/learning/player.vue`（续签逻辑已接） | CDN 未配置时会回退原地址 |
| R14 | 学习关键参数配置化（daily target/签到规则/兑换倍率） | 已完成 | `server/initialize/other.go:138`, `server/initialize/other.go:143`, `server/plugin/english_learning/api/checkin.go:45`, `server/plugin/english_learning/api/checkin.go:82` | 需补后台配置说明文档 |
| R15 | Uni 多语言展示统一（t/localText） | 已完成 | `uni/src/pages/learning/**/*.vue`（多处 `t()`/`localText()`） | Web 运营端文案仍以中文硬编码为主 |
| R16 | 外链资源策略（音频/视频/字幕） | 部分完成 | `server/plugin/english_learning/model/english_model.go:31`, `server/plugin/english_learning/model/video_model.go:32`, `server/plugin/english_learning/model/video_model.go:44` | 缺“上传策略、回源、失效、兜底”统一方案文档 |
| R17 | 自动化测试与验收基线 | 部分完成 | `server/plugin/english_learning/api/*_test.go`, `server/plugin/english_learning/service/*_test.go` | 缺 Web/Uni 端到端回归与授权冲突专项测试 |
| R18 | TTS/音频方案（本地/离线/API/成本）定稿 | 未完成 | `server/plugin/english_learning/service/english_word.go:66`, `server/plugin/english_learning/service/english_word.go:73` | 当前 `generateEdgeTTS` 为模拟实现，未形成可上线音频生产链路 |

## 4. 冲突与风险清单（冻结）
- 授权冲突未解: 目前“全局资产态”与“资源价格/试看”并存，但缺统一优先级引擎和用户粒度授权。
- 新用户免费期未自动发放: 字段存在但缺初始化写入流程，导致规则可配不可落地。
- Web 运营能力不对称: 后端能力多、Web 页面少，运营无法完整管理英语内容资产。
- TTS 生产链路未落地: 当前仅模拟 URL，不具备质量、成本、审计与可追踪性。

## 5. PT-1 结论
- PT-1 已完成: 基线已冻结，需求矩阵已形成。
- 当前整体完成度评估:
  - C端学习体验: 86%
  - 后端接口闭环: 89%
  - Web运营工具完备度: 62%
  - 上线前综合完成度: 80%

## 6. 下一步门槛（进入 PT-2 前）
- 先实现统一授权引擎与用户粒度授权模型（解决 R09/R10/R11）。
- 再补 Web 全链路运营页面（解决 R05）。
- 最后定稿音频/TTS生产方案并补测试矩阵（解决 R18/R17）。
