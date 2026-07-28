# 英语学习模块改造计划

> 范围：uni 端 6 个页面 + web 后台 2 个页面 + 后端插件层。原则：最小改动、不误动其它功能。

---

## 一、需求总览

| 编号 | 页面 | 需求 |
|------|------|------|
| A1 | diary.vue | 轮播图高度增加 1/2 |
| A2 | diary.vue | 移除 page-hero，筛选区参照 typing.vue 顶栏风格 |
| A3 | diary.vue + tag-filter.vue | 标签选择弹窗换新风格（选中置顶+椭圆+X取消+激活样式） |
| A4 | diary.vue | 第一张轮播字幕应用试看比例限制 |
| B1 | diary-detail.vue | 字幕激活滚动被固定顶部遮挡修复 + 离底部控制栏间距 5-8px |
| B2 | diary-detail.vue + player.vue | 单词弹窗加取消+收藏按钮，语音播放改用 typing.vue 优先级流程+动画 |
| B3 | diary-detail.vue + collections.vue | 日记句子收藏与日记整体收藏区分显示 |
| B4 | diary-detail.vue | "更多"弹窗：移除口音切换，加速率选择（参照 player.vue） |
| B5 | diary-detail.vue | US/UK 双字幕（两套时间区间） |
| B6 | diary-detail.vue + player.vue | 试看比例过滤（锁头显示） |
| C1 | collections.vue | 所有类型加"查看"按钮（唯一跳转入口），"取消收藏"移到左边 |
| D1 | web diary.vue | 编辑日记标签自动填充 |
| D2 | web diary.vue | 试看比例生效（加会员模式开关） |
| D3 | web diary.vue | 字幕文件自动填充 + 查看字幕功能 |
| D4 | web diary.vue | US/UK 双字幕上传与解析 |
| E1 | web video.vue | 字幕文件自动填充 |
| F1 | uni 全局 | 底部弹窗防 tab/控制栏遮挡 |

---

## 二、后端改动（english_learning 插件）

### 2.1 新增收藏类型 5 = 日记句子

**问题**：日记句子收藏当前用 `targetType=2`（句子），但后端 `GetCollectionDetailList` 对 type 2 只查 `VideoSentence` 表，导致日记句子收藏无法解析详情。

**方案**：新增 `targetType=5`（日记句子），收藏时 targetId 存 `DiarySentence.ID`。

| 文件 | 改动 |
|------|------|
| `model/user_learning_model.go` | `UserCollection.TargetType` 注释更新为 `1单词 2视频句子 3视频 4日记 5日记句子` |
| `service/user_data.go` `GetCollectionDetailList` | 新增 case 5：查 `DiarySentence` 表获取句子，再关联查 `Diary` 获取日记信息，返回时同时填充 `Sentence`（复用结构）和 `Diary` 字段 |
| `model/request/user_data.go` | `CollectionReq.TargetType` 无需改（已是 int） |

### 2.2 DiarySentence 模型增加 UK 时间字段

**问题**：美式/英式音频时长不同，单套时间区间导致切换口音后字幕激活不准。

**方案**：在 `DiarySentence` 增加 `StartTimeUk`、`EndTimeUk` 字段，US 时间保持原 `StartTime`/`EndTime`。文字内容（English/Translate）两套共用。

| 文件 | 改动 |
|------|------|
| `model/diary_model.go` | `DiarySentence` 增加 `StartTimeUk float64` 和 `EndTimeUk float64`，默认 0 表示无 UK 时间（回退用 US 时间） |
| `initialize/gorm.go` | AutoMigrate 已包含 DiarySentence，无需额外操作 |
| `service/diary.go` `ParseAndSaveDiarySentencesFromFiles` | 增加可选的 UK 字幕文件参数，解析后分别填充 US/UK 时间 |
| `model/request/diary.go` | `ParseDiarySubtitleFilesReq` 增加 `EnglishSubtitleUrlUk` 字段；`DiarySubtitleFileItem` 保持不变 |

### 2.3 字幕列表 API 增加试看过滤

**问题**：试看比例限制下，超出范围的字幕内容不应返回，但时间区间要显示（内容区显示锁头）。

**方案**：在客户端获取字幕列表的 API/Service 中，根据用户权限和试看比例，对超出范围的句子做内容脱敏（清空 English/Translate，加 `locked: true` 标记），时间区间保留。

| 文件 | 改动 |
|------|------|
| `service/diary.go` | 新增 `GetDiarySentenceListWithAuth(diaryID, hasFullAuth, trialPercent)` 方法：查全部句子，对 `startTime > maxTrialTime` 的句子清空 English/Translate 并标记 `locked` |
| `model/diary_model.go` | `DiarySentence` 增加 `Locked bool` 字段 `gorm:"-"`（不入库，仅传输用） |
| `api/content.go` 或 `api/diary.go` | 客户端获取日记句子的接口调用新方法 |
| `service/video_subtitle.go` | 视频句子列表同理增加试看过滤逻辑 |

**边界规则**：`maxTrialTime = duration * trialPercent / 100`。若某句 `startTime <= maxTrialTime` 但 `endTime > maxTrialTime`（跨边界），该句可显示（不锁）。即仅锁 `startTime > maxTrialTime` 的句子。

### 2.4 日记列表接口返回试看信息

**问题**：diary.vue 第一张轮播字幕也需要试看限制。

**方案**：日记列表/详情接口已返回 `trialPercent` 和 `hasFullAuth`（gorm:"-"），前端据此过滤预览字幕即可，无需后端改动。但需确认客户端日记详情接口会填充 `HasFullAuth`。

| 文件 | 改动 |
|------|------|
| `api/content.go` | 确认获取日记详情时调用 `LearningAuthzService` 填充 `HasFullAuth`（若已有则无需改动） |

---

## 三、uni 端改动

### 3.1 pages/learning/diary.vue

#### A1. 轮播图高度增加 1/2
- **当前**：`.diary-swiper { height: 400rpx; }`（第 521-524 行）
- **改为**：`height: 600rpx;`

#### A2. 移除 page-hero，筛选区改 typing.vue 顶栏风格
- **删除**：`<view class="page-hero">...</view>`（第 4-10 行）及其全部样式（第 373-416 行）
- **新增**：参照 typing.vue 的 `.top-bar` 结构（第 4-21 行），在页面顶部放渐变顶栏，内含分类+标签选择器（胶囊形、半透明白底、白字、旋转箭头）+ 右侧重置按钮
- 选择器样式参照 typing.vue 的 `.selector`（`border-radius: 999rpx; background: rgba(255,255,255,0.22); backdrop-filter: blur(8rpx); color: #fff;`）
- 分类弹窗也参照 typing.vue 的 `.sheet` 样式（渐变把手、圆角卡片项、激活边框+背景）

#### A3. 标签选择弹窗换新风格
- **新交互**：
  - 已选标签置顶显示，每个为椭圆形 chip，右上角有 X 取消按钮
  - 未选标签在下方列表，点击选中后移动到置顶区
  - 有明确的激活/选中样式（当前缺失）
- **diary.vue 改动**：重写标签弹窗（第 116-152 行），分两个区域：
  1. 已选区（横向滚动）：椭圆形 chip `标签名 ✕`，点 X 移除
  2. 未选区（列表）：点击添加到已选区，选中项高亮
- **tag-filter.vue 改动**：同样替换标签弹窗（第 57-85 行）为新风格
- **新样式**：
  ```css
  .selected-chip { display: inline-flex; align-items: center; padding: 10rpx 20rpx; border-radius: 999rpx; background: linear-gradient(135deg,#6D5BFF,#9B8FFF); color: #fff; }
  .chip-close { margin-left: 8rpx; font-size: 24rpx; }
  .tag-option.active { color: #6D5BFF; font-weight: 700; background: rgba(108,91,255,0.12); }
  ```

#### A4. 第一张轮播字幕试看限制
- **当前**：`loadPreviewSentences`（第 267-290 行）取前 10 条，无试看过滤
- **改为**：根据日记的 `trialPercent` 和 `hasFullAuth`，过滤掉 `startTime > maxTrialTime` 的句子（与后端 2.3 规则一致）

### 3.2 pages/learning/diary-detail.vue

#### B1. 字幕区间距修复
- **问题**：字幕激活行被固定顶部信息区遮挡；字幕区离底部控制栏太远
- **修复**：
  - `.subtitle-section` 高度计算中，顶部避让值从 `110rpx` 调整为实际顶部高度 + `margin-top` 间距（加 `margin-top: 16rpx` 或调大 padding-top）
  - `.sentence-row` 的 `scroll-margin-top` 从 `120rpx` 增大至能完全露出激活行的高度（约 `160rpx`）
  - 底部间距：`.subtitle-section` 高度计算中底部避让从 `130rpx` 减小，使字幕区离底部控制栏 5-8px（约 `10-16rpx`）

#### B2. 单词弹窗改造
- **当前**（第 119-134 行）：只有收藏按钮，无取消按钮；音频图标是 🔊 emoji
- **改为**：
  1. 加"取消"按钮（底部，与收藏按钮并排或上方关闭按钮）
  2. 音频 logo 改为 typing.vue 风格的 `♪`/`◉` 圆形图标（64rpx 渐变紫圆）
  3. 播放逻辑改用 typing.vue 优先级流程：
     - 第一优先：`tryLocalTTS`（JS 生成 TTS，APP 用 plus.speech，H5 用 speechSynthesis）
     - 第二优先：在线音频（db-primary -> https 重试 -> 有道兜底）
     - 第三优先：最终本地 TTS 兜底
  4. 播放时 `pulse` 动画（0.6s 缩放+透明度循环）
- **player.vue 同步改动**：单词弹窗（第 140-155 行）同样加取消按钮 + 收藏按钮 + typing.vue 风格音频播放
- **需提取公共逻辑**：将 typing.vue 的 `tryLocalTTS`、`playOnlineAudio`、`buildFallbackAudioSource` 提取为 `utils/learning-tts.js` 公共模块，三页面共用

#### B3. 收藏逻辑调整
- **当前**：日记句子收藏用 `collect(2, sentenceId)`（type 2 = 视频句子），导致 collections 无法解析
- **改为**：日记句子收藏用 `collect(5, sentenceId)`（type 5 = 日记句子）
- `collectSentence` 函数（第 454-467 行）：`collect(2, ...)` → `collect(5, ...)`，`uncollect(2, ...)` → `uncollect(5, ...)`
- `sentenceCollectedMap` 查询逻辑同步调整

#### B4. "更多"弹窗改造
- **当前**（第 136-160 行）：只有口音切换卡片
- **改为**：
  1. 移除口音切换卡片（`more-card` + `toggleAccentInMore`）
  2. 新增速率选择卡片 + pill 按钮组，样式完全参照 player.vue 的 `speed-pills`（第 157-197 行 + 第 1445-1545 行样式）
  3. 速率数据：`[0.75, 1.0, 1.25, 1.5, 2.0]`
  4. 速率应用：`audioCtx.playbackRate` —— 注意 `uni.createInnerAudioContext` 不支持 playbackRate，需确认平台支持情况；H5 可用 `audioCtx.playbackRate`，APP 需验证
- 口音切换保留在顶部信息区的 `accent-switch` 快捷按钮（第 17-21 行），不从更多弹窗操作

#### B5. US/UK 双字幕
- **当前**：单套字幕，切换口音仅切换音频源
- **改为**：
  1. 加载字幕时同时获取 US/UK 两套时间区间
  2. `updateActiveSentence`（第 361-375 行）根据当前 `accent` 选择对应时间区间匹配
  3. 若 UK 时间为 0（未上传 UK 字幕），回退用 US 时间
  4. 切换口音时重新计算激活行
- **数据结构**：每条句子增加 `startTimeUk`、`endTimeUk`（后端 2.2 提供）
- **性能优化**：不维护两份数组，在 `updateActiveSentence` 中动态选择字段即可，零额外内存

#### B6. 试看比例过滤
- **当前**：无试看限制
- **改为**：
  1. 获取日记详情时拿到 `trialPercent` 和 `hasFullAuth`
  2. 字幕列表中 `locked: true` 的句子：时间区间正常显示，内容区居中显示锁头图标
  3. 锁定行不可点击（不触发 jumpBySubtitle、不显示收藏按钮）
  4. 播放进度超过试看时间时暂停并提示（参照 player.vue 的 `ensureTrialLimit`）

### 3.3 pages/learning/player.vue

#### B2-player. 单词弹窗改造
- 同 diary-detail.vue B2，加取消+收藏按钮，音频播放改用公共 TTS 模块

#### B6-player. 试看比例过滤（字幕锁头）
- **当前**：试看限制只回弹视频时间，字幕完整显示
- **改为**：字幕列表中超出试看范围的句子显示锁头（时间区间保留，内容区锁头图标）
- 后端 2.3 会返回 `locked` 标记，前端据此渲染

### 3.4 pages/learning/collections.vue

#### C1. 按钮重构
- **当前**：
  - 点击 card-body 跳转（`openDetail`），无独立"查看"按钮
  - "取消收藏"在右侧，日记类型有播放按钮在左侧
- **改为**：
  1. **移除** card-body 的 `@click="openDetail"`，点击卡片不再跳转
  2. **所有类型**增加"查看"按钮（唯一跳转入口）
  3. "取消收藏"按钮移到**左侧**
  4. "查看"按钮在**右侧**
- **按钮布局**：`card-footer { justify-content: space-between; }`，取消收藏在左，查看在右

#### C1-细节. 各类型"查看"按钮行为
| 类型 | 查看按钮行为 |
|------|-------------|
| 单词(1) | 跳转 typing 页（传 wordId） |
| 视频句子(2) | 跳转 player 页（传 episodeId + sentenceId） |
| 视频(3) | 跳转 video-detail 或 player 页 |
| 日记(4) | 跳转 diary-detail 页（传 diaryId），**无播放按钮** |
| 日记句子(5) | **查看按钮 = 播放按钮**，点击播放该句子音频（typing.vue TTS 优先级流程+动画），不跳转 |

#### C1-tab. 日记 tab 显示 type 4 + type 5
- tabs 定义（第 109-115 行）：日记 tab `{ type: 4 }` 改为同时查 4 和 5
- 后端 `GetCollectionList` 当前 `if info.TargetType > 0` 只查单一类型，需改为支持查询多个类型（`target_type IN ?`）
- 或前端请求两次合并（不推荐，分页不准）
- **方案**：后端 `CollectionSearch` 增加 `TargetTypes []int` 字段，`GetCollectionList` 和 `GetCollectionDetailList` 支持 `IN` 查询

### 3.5 pages/learning/tag-filter.vue

#### A3-tag-filter. 标签弹窗换新风格
- 替换第 57-85 行的标签弹窗为 diary.vue A3 的新风格（已选置顶+椭圆 chip+X+激活样式）
- 复用同一套样式类名

### 3.6 公共 TTS 模块提取

新建 `uni/src/utils/learning-tts.js`：
- 从 typing.vue 提取 `tryLocalTTS`、`playOnlineAudio`、`buildFallbackAudioSource`、`playNextAudioCandidate`
- 导出 `playWordTTS(text, audioUs, audioUk, accent, callbacks)` 和 `playSentenceTTS(text, audioUs, audioUk, accent, callbacks)`
- typing.vue、diary-detail.vue、player.vue、collections.vue 共用

---

## 四、web 后台改动

### 4.1 plugin/english_learning/view/diary.vue

#### D1. 编辑日记标签自动填充
- **当前**（第 730-756 行 `openDiaryDialog`）：直接用 `row.tags`，不调详情接口；若列表未返回 tags 则标签为空
- **改为**：参照 video.vue 的 `openEpisodeDialog`（第 1319-1329 行），编辑时调用 `findDiary` 获取详情，从详情中加载 tags
  ```js
  const openDiaryDialog = async (row) => {
    if (row) {
      diaryDialogMode.value = 'update'
      diaryForm.value = { /* ...现有字段... */ }
      // 新增：调用详情接口加载标签
      try {
        const detailRes = await findDiary({ ID: row.ID })
        if (detailRes.code === 0 && detailRes.data?.tags) {
          diaryForm.value.tagIds = detailRes.data.tags.map(t => t.ID)
        }
      } catch (e) {}
    }
    // ...
  }
  ```
- **后端确认**：`findDiary` 接口（`GetDiary` service）需返回 tags —— 当前 `GetDiary` 只返回 Diary 结构体不含 tags。需在后端 `GetDiary` 或对应 API 中关联查 tags 并返回

#### D2. 试看比例生效
- **当前**：日记有 `trialPercent` 字段但客户端 diary-detail 未使用
- **改动**：
  1. 后端 2.3/2.4 已让试看比例在字幕过滤中生效
  2. 用户提到"如果需要 video 的需要会员模式，就放在日记来设置多一层类似的"——在日记表单增加"需要会员"开关（`needVip`），控制是否启用试看限制
- **后端**：
  - `model.Diary` 增加 `NeedVip *bool` 字段
  - `CreateDiaryReq`/`UpdateDiaryReq` 增加 `NeedVip` 字段
  - 客户端获取日记时，若 `needVip=true` 且用户非 VIP，则应用 `trialPercent` 限制
- **前端 diary.vue**：表单增加"需要会员"开关（参照 video.vue 第 392-394 行的 el-switch）

#### D3. 字幕文件自动填充 + 查看字幕功能
- **字幕自动填充**：
  - **当前**：`openDiarySubtitleParser`（第 819-858 行）重置 `englishSubtitleUrl` 为空
  - **改为**：Diary 模型增加 `englishSubtitleUrl` 字段存储已上传的字幕文件 URL；打开解析时回填
  - 或：从已有句子反推不可行（句子不含源文件 URL），需在 Diary 模型增加字幕文件 URL 字段
- **查看字幕功能**：
  - 参照 video.vue 的 `previewEpisodeSubtitles`（第 1478-1509 行）+ 字幕预览弹窗（第 504-594 行）
  - 在日记列表操作列增加"查看字幕"按钮
  - 弹窗展示 DiarySentence 列表（时间、英文、翻译），支持编辑

#### D4. US/UK 双字幕上传与解析
- **当前**：字幕解析只上传一个英文字幕文件
- **改为**：
  1. 增加可选的"英式英文字幕文件"上传（用于 UK 时间轴）
  2. `ParseDiarySubtitleFilesReq` 增加 `EnglishSubtitleUrlUk`（后端 2.2 已规划）
  3. 解析时：US 字幕填充 `StartTime`/`EndTime`，UK 字幕填充 `StartTimeUk`/`EndTimeUk`
  4. 文字内容（English/Translate）以 US 字幕为准，UK 仅取时间轴

### 4.2 plugin/english_learning/view/video.vue

#### E1. 字幕文件自动填充
- **当前**：`openSubtitleParser`（第 1463-1476 行）重置为空
- **改为**：VideoEpisode 模型增加 `englishSubtitleUrl` 字段存储已上传字幕 URL；打开解析时回填
- 后端 `VideoEpisode` 模型增加字幕文件 URL 字段

---

## 五、底部弹窗防遮挡（uni 全局）

### F1. 弹窗 safe-area 适配检查
- 检查所有 `uni-popup type="bottom"` 和自定义底部弹窗：
  - diary.vue：分类弹窗、标签弹窗
  - diary-detail.vue：单词弹窗、更多弹窗
  - player.vue：单词弹窗、更多弹窗
  - tag-filter.vue：分类弹窗、标签弹窗
  - typing.vue：选择器弹窗、设置弹窗
- **规则**：弹窗底部 padding 必须包含 `env(safe-area-inset-bottom)` + tab 栏高度（若在 tab 页）
- diary.vue 和 typing.vue 是 tab 页，底部弹窗需额外避开 tab 栏（约 100rpx + safe-area）
- diary-detail.vue 和 player.vue 非 tab 页，只需 safe-area
- **改动**：为缺少 safe-area 适配的弹窗补上 `padding-bottom: calc(基础值 + env(safe-area-inset-bottom))`
- tab 页弹窗额外加 `+ 100rpx`（tab 栏高度）或使用 `uni.getSystemInfoSync()` 动态计算

---

## 六、实施顺序

### 阶段 1：后端基础（不破坏现有功能）
1. DiarySentence 模型加 UK 时间字段 + Locked 标记
2. Diary 模型加 NeedVip + 字幕文件 URL 字段
3. 新增 targetType 5 收藏详情解析
4. 字幕列表试看过滤方法
5. 收藏列表支持多类型查询
6. findDiary 返回 tags

### 阶段 2：web 后台
7. diary.vue 标签自动填充
8. diary.vue 会员模式开关
9. diary.vue 字幕文件回填 + 查看字幕 + UK 字幕上传
10. video.vue 字幕文件回填

### 阶段 3：uni 公共模块
11. 提取 learning-tts.js

### 阶段 4：uni 端页面
12. diary.vue（轮播高度 + 顶栏 + 标签弹窗 + 试看过滤）
13. diary-detail.vue（间距 + 单词弹窗 + 收藏类型 + 更多弹窗 + 双字幕 + 试看锁头）
14. player.vue（单词弹窗 + 试看锁头）
15. collections.vue（按钮重构 + 日记 tab）
16. tag-filter.vue（标签弹窗）
17. 底部弹窗 safe-area 检查修复

---

## 七、风险与注意事项

1. **playbackRate 兼容性**：`uni.createInnerAudioContext` 的 `playbackRate` 在部分平台不支持，需验证；若不支持则速率选择仅对视频（videoCtx）有效，音频日记可能需要降级处理
2. **TTS 平台差异**：`plus.speech` 仅 APP 端可用，`speechSynthesis` 仅 H5 可用，小程序端无本地 TTS —— 需确保在线音频兜底覆盖所有平台
3. **数据迁移**：DiarySentence 新增 UK 时间字段后，历史数据 UK 时间为 0，需回退用 US 时间（已在逻辑中处理）
4. **收藏类型兼容**：新增 type 5 后，历史 type 2 的日记句子收藏仍无法解析 —— 可选择写迁移脚本将历史日记句子收藏从 type 2 迁移到 type 5（通过判断 targetId 是否在 DiarySentence 表中）
5. **不误动原则**：所有改动仅涉及上述文件，不修改 i18n 配置（除非新增 key）、不修改路由配置、不修改其他业务模块
