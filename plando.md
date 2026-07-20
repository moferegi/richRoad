# richRoad 全局改造施工图纸 (plando.md)

> 使用方法：按阶段/步骤顺序执行，每完成一步勾选 `[x]`。遇到问题记录在对应步骤下。

---

## 阶段一：全局基础设施改造 (Server 端)

### 1.1 双端鉴权物理隔离检查加固

- [ ] **1.1.1** 检查 `server/middleware/jwt.go`，确认 Uni 端和 Web 端使用独立的 JWT 校验路径
  - 文件：`server/middleware/jwt.go`
  - 操作：确认已有 `X-Client-Platform` header 判断逻辑；若无则新增分支
- [ ] **1.1.2** 检查 `server/initialize/router_biz.go`，确认 Uni 路由组统一挂载 `UniJWTAuth()` 中间件
  - 文件：`server/initialize/router_biz.go`
  - 操作：确保所有 `/api/` 前缀的 Uni 路由使用独立鉴权中间件
- [ ] **1.1.3** 检查 Web 路由组统一挂载 `JWTAuth() + CasbinHandler()`
  - 文件：`server/initialize/router.go`（system 路由）
  - 操作：确认 Web 后台路由均使用 Casbin RBAC 中间件
- [ ] **1.1.4** 梳理现有 Uni 路由清单，确保无遗漏
  - 文件：`server/router/client/` 下所有 `*_router.go`
  - 验证：每个路由组都正确挂载了 Uni 鉴权

### 1.2 Uni 端多语言后端自动拦截

- [ ] **1.2.1** 审查 `server/middleware/uni_response_protect.go` 第 88 行的 `LocalizeI18nPayloadByContext` 函数
  - 文件：`server/middleware/uni_response_protect.go`
  - 检查：`server/utils/` 中对应的实现，确认已覆盖的多语言字段
- [ ] **1.2.2** 扩展 `LocalizeI18nPayloadByContext` 的多语言字段列表
  - 新增字段：`seriesName`, `episodeName`, `categoryName`, `chapterName`, `word`, `phoneticUs`, `phoneticUk`, `sentenceSource`
  - 文件：`server/utils/i18n_util.go`（或对应文件）
  - 逻辑：遍历 JSON 数据的每个字段，若值为 `{"zh":"...","en":"..."}` 格式则裁剪为单语字符串
- [ ] **1.2.3** 确认 Web 端响应不走此裁剪逻辑
  - 验证：`shouldProtectForUni()` 函数中 `X-Client-Platform` 不为 `uni` 时直接跳过
- [ ] **1.2.4** 端到端验证：用 Postman 分别模拟 Uni header 和 Web header，检查响应差异

### 1.3 Uni 端接口加密与签名全覆盖

- [ ] **1.3.1** 将 `uniProtectPathPrefixes` 白名单改为前缀匹配
  - 文件：`server/middleware/uni_response_protect.go`
  - 当前：硬编码 14 个路径前缀（第 26-40 行）
  - 改为：`strings.HasPrefix(path, "/api/")` 自动纳入所有 `/api/` 请求
- [ ] **1.3.2** 新增签名校验中间件 `server/middleware/uni_sign_verify.go`
  - 校验 `X-Req-Ts`（时间戳）、`X-Req-Nonce`（随机数）、`X-Req-Sign`（HMAC-SHA256）
  - 时间窗口：±300 秒（5 分钟）
  - Nonce 去重：使用 Redis `SET NX EX` 或内存缓存
  - 签名算法：`HMAC-SHA256(method|path|query|ts|nonce, SHA256(token+"|uni-api-sign-v1"))`
- [ ] **1.3.3** 读取 `learning_api_sign_enabled` 开关
  - 函数：复用 `utils.GetBoolSetting("learning_api_sign_enabled", "LEARNING_API_SIGN_ENABLED", true)`
  - 当开关关闭时，签名校验中间件放行所有请求
- [ ] **1.3.4** 在 Uni 路由组中挂载签名中间件
  - 文件：`server/initialize/router_biz.go`
  - 在 `UniJWTAuth()` 之后、业务 handler 之前挂载 `UniSignVerify()`

### 1.4 去域名化存储规范落地

- [ ] **1.4.1** 创建统一文件上传路径剥离工具函数
  - 文件：`server/utils/upload.go`（新建或扩展）
  - 函数：`StripDomain(rawUrl string) string`
  - 逻辑：解析 URL，若为已知云存储域名则剥离，返回相对路径
- [ ] **1.4.2** 在文件上传入库的 Service 层统一调用 `StripDomain()`
  - 涉及：`server/service/` 下所有处理文件上传的 Service 函数
  - 每次保存文件路径到数据库前，调用 `StripDomain()` 确保只存相对路径
- [ ] **1.4.3** 审计现有数据库记录
  - 编写 SQL 查询脚本，找出 `english_learning_*` 相关表中包含 `https://` 的记录
  - 编写 Go 迁移脚本批量清理（可选，视数据量决定）

### 1.5 Uni 请求签名前端配合

- [ ] **1.5.1** 确认 `uni/src/utils/request.js` 签名头已正确生成
  - 当前已有 `buildSignHeaders()` 函数（第 366-379 行）
  - 验证：nonce 格式、ts 精度、sign 算法与后端一致
- [ ] **1.5.2** 确认 `learning_api_sign_enabled` 开关前端读取
  - 存储 key：`learning-api-sign-enabled`
  - 当关闭时不发送签名头

---

## 阶段二：Uni 移动端 - 视频播放器升级

> 主文件：`uni/src/pages/learning/player.vue`

### 2.1 字幕初始定位

- [ ] **2.1.1** 在 `initPage()` 完成数据加载后，强制滚动到激活字幕位置
  - 当前逻辑（第 527-533 行）：`focusSentence()` 已尝试聚焦，但 `activeSentenceId` 计算可能不够准确
  - 改进：在 `loadSubtitleList()` 完成后，用 `$nextTick` + 延迟 300ms 确保 DOM 渲染后再设置 `activeSentenceId`
- [ ] **2.1.2** 修改 `activeSentenceId` 的计算逻辑
  - 当前：`'sentence_' + Math.max(0, targetIndex - 1)`（向上偏移一行）
  - 目标：初始定位时使用 `'sentence_0'`（第一行紧贴视频播放器下方）

### 2.2 切换双语/字幕时滚动位置稳定

- [ ] **2.2.1** 在 `toggleBilingual()` 和 `toggleSubtitles()` 切换后重新锚定
  - 在切换函数末尾调用 `reAnchorScroll()`，根据当前 `activeSentenceIndex` 重新设置 `activeSentenceId`
- [ ] **2.2.2** 新增 `reAnchorScroll()` 函数
  - 先用空字符串重置 `activeSentenceId`（触发 scroll-view 重置）
  - 下一帧再设置正确的 `activeSentenceId`
- [ ] **2.2.3** 针对滚动慢半拍问题，引入提前偏置
  - 在 `activeSentenceId` 设置前，加 `setTimeout(..., 50)` 微延迟等待渲染

### 2.3 底部播放时字幕跟踪

- [ ] **2.3.1** 确保 `onTimeUpdate` 在非播放器内播放时也能触发
  - 检查 uni-app `video` 组件在页面不可见时是否仍触发 `@timeupdate`
  - 若不触发，改为使用 `setInterval` 轮询 `videoCtx` 获取当前时间
- [ ] **2.3.2** 统一字幕点亮逻辑
  - `onTimeUpdate`（第 602-608 行）已有基础匹配逻辑
  - 优化：当 `activeSentenceIndex` 变化时，同步更新 `activeSentenceId` 触发滚动

### 2.4 手势交互与弹窗暂停

- [ ] **2.4.1** 在字幕 `scroll-view` 上绑定 `@touchstart` / `@touchmove` 事件
  - 滑动开始时调用 `videoCtx.pause()`
  - 显示中央提示文字："点击字幕可播放当前位置"
- [ ] **2.4.2** 新增提示文字 UI 元素
  - 在 `video-section` 内添加遮罩层，滑动时显示、点击后消失
- [ ] **2.4.3** 点击字幕或播放器恢复播放时，提示文字消失
  - 在 `jumpBySubtitle()` 和 `togglePlay()` 中隐藏提示
- [ ] **2.4.4** 确认 `showWordDetail()` 中已有 `videoCtx.pause()`（第 741 行已实现）

### 2.5 控制功能重构

- [ ] **2.5.1** 将"速率"功能从底部控制栏移除
  - 删除底部 `changeSpeed` 按钮（第 82-85 行的 `control-item`）
- [ ] **2.5.2** 在"更多设置"面板中新增速率选择
  - 在 `more-sheet` 内增加速率选择行（0.75x / 1.0x / 1.25x / 1.5x / 2.0x）
- [ ] **2.5.3** 原"速率"位置替换为"复读"按钮
  - 新增按钮：图标 `🔄`，文字 `t('player.repeat')`
  - 点击后开启单句循环模式：当前激活字幕播放完毕后自动 seek 回 `startTime`
- [ ] **2.5.4** 实现单句循环逻辑
  - 在 `onTimeUpdate` 中检测循环模式 + 当前时间超过 `endTime` 时，自动 `seek(startTime)`
- [ ] **2.5.5** "更多设置"左下角"返回"改为"退出页面"
  - 修改 `handleBack` 函数及相关文案 key 从 `player.sheet_back` 改为 `player.exit_page`

### 2.6 进度条开发

- [ ] **2.6.1** 在视频播放器下方新增进度条组件
  - 位置：`video-section` 与 `subtitle-section` 之间
  - 结构：`<slider>` + 右侧 `<text>` 显示总时长
- [ ] **2.6.2** 进度条绑定视频播放进度
  - `value` 绑定 `currentTime`（从 `onTimeUpdate` 获取）
  - `max` 绑定 `videoInfo.duration`
- [ ] **2.6.3** 支持拖拽跳转
  - `@change` 事件调用 `seekVideo(newValue)`
- [ ] **2.6.4** 右侧显示总时长
  - 格式：`mm:ss`，使用已有 `formatSeconds()` 函数

---

## 阶段三：Uni 移动端 - 首页与打字页

### 3.1 Uni 首页缓存无感渲染

> 文件：`uni/src/pages/learning/home.vue`

- [ ] **3.1.1** 分析当前 `onShow` 逻辑
  - 问题：第 236 行无条件的 `reloadHome()` 导致每次返回首页都刷新
- [ ] **3.1.2** 新增缓存时间戳变量 `lastLoadTime`
  - 单位：毫秒，初始值 0
- [ ] **3.1.3** 修改 `onShow` 逻辑
  - 若 `lastLoadTime > 0` 且 `Date.now() - lastLoadTime < 5 * 60 * 1000`（5 分钟）
  - 且有缓存数据（`videoList.length > 0`），则直接渲染不刷新
  - 同时检查语言是否变更（`homeLocaleLoaded !== locale`），变更时才刷新
- [ ] **3.1.4** 新增静默后台刷新机制
  - 超过缓存时间后，数据静默更新（不显示 loading），更新后替换列表
- [ ] **3.1.5** `switchCategory` 切换分类时正常刷新，不走缓存

### 3.2 单词跟打页面交互升级

> 文件：`uni/src/pages/learning/typing.vue`

- [ ] **3.2.1** Skeleton 骨架屏
  - 在 `top-nav` 下方、`word-card` 区域添加骨架屏组件
  - 当 `wordList.value.length === 0` 且正在加载时，显示骨架占位
  - 数据加载完成后切换为真实内容
- [ ] **3.2.2** 语音播放动画
  - 在 `audio-btn` 点击时添加 CSS class `is-playing`，触发脉冲/旋转动画
  - 利用 `audioCtx.onPlay` 和 `audioCtx.onEnded` / `onError` 控制动画状态
  - 播放期间 `pointer-events: none` 防止重复点击
- [ ] **3.2.3** 离线发音方案（Web Speech API）
  - 检测 `window.speechSynthesis` 可用性
  - 新增 `speakOffline(text, lang='en-US')` 函数
  - 当在线音频源加载失败或网络不可用时，降级到 `SpeechSynthesis`
  - 在 `playAudio` 的错误回调中增加此降级路径
- [ ] **3.2.4** 单词抽屉自动聚焦
  - 在 `openWordDrawer()` 中，给当前单词的 `view` 元素设置 `id="drawer-word-{currentWordIndex}"`
  - 使用 `scroll-into-view` 属性定位
  - 或使用 `uni.createSelectorQuery()` 计算滚动位置后手动 `scrollTo()`
- [ ] **3.2.5** 切换章节后抽屉聚焦
  - 在 `selectPickerItem()` 中，如果当前是章节切换，操作完成后若抽屉打开则重新聚焦
  - 在 `jumpToChapter()` 函数末尾调用 `reFocusDrawer()`

---

## 阶段四：Web 后台管理系统升级

### 4.1 英语视频标签管理

#### 4.1.1 后端：数据模型

- [ ] 创建 `server/model/client/video_tag.go`
  - 参考文件：`server/model/shop/tag.go`
  - 独立表名：`video_tags`（不可与 `Tags` 混用）
  - 字段：`Name`, `NameI18n`, `Description`, `Color` + `global.GVA_MODEL`
- [ ] 创建 `server/model/client/request/video_tag.go`
  - 参考文件：`server/model/shop/request/tag.go`
  - 含 `VideoTagSearch` 分页查询结构体

#### 4.1.2 后端：Service 层

- [ ] 创建 `server/service/client/video_tag.go`
  - `CreateVideoTag`, `DeleteVideoTag`, `DeleteVideoTagByIds`, `UpdateVideoTag`, `FindVideoTag`, `GetVideoTagList`
  - 复制 `server/service/shop/tag.go` 的同名函数，替换模型引用
- [ ] 在 `server/service/client/enter.go` 中注册 `VideoTagService`

#### 4.1.3 后端：API 层

- [ ] 创建 `server/api/v1/client/video_tag.go`
  - 复制 `server/api/v1/shop/tag.go` 的 API 函数，替换 Service 引用
- [ ] 在 `server/api/v1/client/enter.go` 中注册 `VideoTagApi`

#### 4.1.4 后端：Router 层

- [ ] 创建 `server/router/client/video_tag.go`
  - 路由前缀：`/videoTag`
  - 挂载 Web 管理后台鉴权中间件（非 Uni 端）
- [ ] 在 `server/router/client/enter.go` 中注册 `VideoTagRouter`

#### 4.1.5 后端：数据库迁移

- [ ] 在 `server/initialize/gorm_biz.go` 中注册 `client.VideoTag` 的 `AutoMigrate`

#### 4.1.6 前端：API 封装

- [ ] 创建 `web/src/api/client/videoTag.js`
  - 复制 `web/src/api/shop/tag.js`，调整 URL 前缀

#### 4.1.7 前端：页面

- [ ] 创建 `web/src/view/client/videoTag/videoTag.vue`
  - 复制 `web/src/view/shop/tag/tag.vue` 页面逻辑
  - 修改页面标题、API 引用

#### 4.1.8 前端：菜单注册

- [ ] 在菜单管理中添加"视频标签"菜单项
  - 父级菜单：`layout/client`
  - 路由路径：`/client/videoTag`

### 4.2 重点单词多选与字幕解析

> 涉及：英语学习视频的字幕解析相关接口（具体文件需根据实际代码路径确认）

- [ ] **4.2.1** 改造字幕解析接口：拆分"扫描提取"和"确认入库"两步
  - 第一步（`POST /englishLearning/video/scanKeywords`）：解析字幕文件，提取其中包含的重点单词列表
  - 返回单词列表（含 id, word, explanation），供管理员多选
- [ ] **4.2.2** 第二步（`POST /englishLearning/video/parseSubtitles`）：接收管理员选中的单词 ID 列表，执行字幕解析入库
- [ ] **4.2.3** "查看字幕"预览弹窗增加重点单词列表展示
  - 前端组件：在现有字幕预览弹窗中增加一个 Tab 或折叠面板
  - 显示该集关联的重点单词列表
- [ ] **4.2.4** 支持在预览中重新多选单词并更新入库
  - 增加"编辑重点单词"按钮 → 打开多选弹窗
  - 提交后调用更新接口

### 4.3 外链域名与多云上传配置

> 文件：`web/src/view/client/externalLinkDomain/externalLinkDomain.vue` + `server/service/client/external_link_domain.go`

#### 4.3.1 扩展数据模型

- [ ] 扩展 `server/model/client/external_link_domain.go`
  - 新增字段：`Provider`（qiniu/r2/b2）、`AccessKey`、`SecretKey`、`Bucket`、`Region`、`Endpoint`
  - 新增字段：`IsDefaultUpload`（默认上传云）、`SyncToAll`（同步上传至所有云）
- [ ] 更新数据库迁移（`AutoMigrate`）

#### 4.3.2 检测连接功能

- [ ] 新增 API：`POST /extDomain/testConnection`
  - 参数：云存储配置（Provider + 密钥等）
  - 返回：连接成功/失败状态 + 延迟时间
  - 实现：调用各云 SDK 的 list buckets 或 head bucket 接口

#### 4.3.3 一键目录比对

- [ ] 新增 API：`POST /extDomain/compareDirectories`
  - 参数：源云和目标云标识
  - 遍历两边的目录和文件名，生成差异报告（源有目标无、目标有源无）
  - 只比对目录名和文件名，不比对文件内容
- [ ] 前端展示差异报告表格

#### 4.3.4 在线文件浏览器

- [ ] 新增 API：`POST /extDomain/listFiles`
  - 参数：云存储标识 + 当前目录路径
  - 返回：该目录下的子目录列表 + 文件列表
- [ ] 前端实现目录树 + 文件列表组件

#### 4.3.5 默认上传云设置

- [ ] 在域名编辑弹窗中增加"设为默认上传云"开关
- [ ] 同时只能有一个云为默认上传云（设置时取消其他的默认状态）

### 4.4 后台统一文件上传组件重构

#### 4.4.1 创建通用上传组件

- [ ] 创建 `web/src/components/FileUploadWithDir/index.vue`
  - 包含：上传按钮 + 上传目录输入框（只读显示）+ "修改/创建目录"按钮
  - 上传完成后返回相对路径
  - 支持配置上传目录前缀（如 `english-learn/word/`）
- [ ] 创建"修改/创建目录"弹窗组件
  - 输入新目录名，自动拼接在配置的前缀下

#### 4.4.2 englishLearningWord 上传组件替换

- [ ] "单词管理"美式发音 URL 字段：替换为 `FileUploadWithDir`
  - 目录前缀：`english-learn/word/audio/us/`
- [ ] "单词管理"英式发音 URL 字段：替换为 `FileUploadWithDir`
  - 目录前缀：`english-learn/word/audio/uk/`
- [ ] "分类管理"logo 地址字段：替换为 `FileUploadWithDir`
  - 目录前缀：`english-learn/category/logo/`

#### 4.4.3 englishLearningVideo 上传组件替换

- [ ] "剧集管理"封面地址：替换为 `FileUploadWithDir`
  - 目录前缀：`english-learn/video/cover/`
- [ ] "单集与字幕"视频地址：替换为 `FileUploadWithDir`
  - 目录前缀：`english-learn/video/episode/`
- [ ] "字幕解析入库"字幕文件上传：替换为 `FileUploadWithDir`
  - 目录前缀：`english-learn/video/subtitle/`

#### 4.4.4 shop/popup 上传组件替换

- [ ] "弹窗图片"：替换为 `FileUploadWithDir`
  - 目录前缀：`shop/popup/image/`
- [ ] "多语言弹窗内容中的图片"：替换为 `FileUploadWithDir`
  - 目录前缀：`shop/popup/content/`

#### 4.4.5 shop/banner 上传组件替换

- [ ] Banner"图片（上传）"：替换为 `FileUploadWithDir`
  - 目录前缀：`shop/banner/`

#### 4.4.6 去域名化验证

- [ ] 逐一检查以上所有上传字段，确认数据库保存的是相对路径
- [ ] 前端确认 `getExternalUrl()` 或类似函数正确拼接了域名

---

## 阶段五：全链路验收

### 5.1 多语言验收

- [ ] 用 Uni 端切换 zh/en/mn 语言，验证接口返回字段已裁剪为单语
- [ ] 用 Web 端访问同一接口，验证返回完整多语言 JSON

### 5.2 接口安全验收

- [ ] 关闭 `learning_api_encrypt_enabled`，验证响应不再加密
- [ ] 开启加密，用 Charles 抓包验证 data 字段已加密
- [ ] 关闭 `learning_api_sign_enabled`，验证不校验签名
- [ ] 开启签名，用错误签名请求验证被拦截

### 5.3 去域名化验收

- [ ] 上传文件后查询数据库，确认 URL 字段为相对路径
- [ ] Uni 端和 Web 端均能正常展示图片/视频

### 5.4 播放器验收

- [ ] 进入播放器，首行字幕紧贴视频下方
- [ ] 切换双语/字幕，滚动位置不跳动
- [ ] 底部播放时字幕跟踪正常
- [ ] 滑动字幕暂停 + 提示文字显示/消失
- [ ] 点击单词弹窗后视频暂停
- [ ] 复读按钮单句循环正常
- [ ] 进度条拖拽跳转正常

### 5.5 首页验收

- [ ] 从其他页面返回首页不闪烁
- [ ] 超过 5 分钟后返回首页数据已静默更新

### 5.6 打字页验收

- [ ] 骨架屏在加载时显示、加载完成后消失
- [ ] 语音图标点击后有动画、不可重复点击
- [ ] 离线发音降级方案可用
- [ ] 单词抽屉自动聚焦到当前词

### 5.7 Web 后台验收

- [ ] 视频标签 CRUD 正常
- [ ] 字幕解析两步流程正常
- [ ] 多云连接检测、目录比对、文件浏览正常
- [ ] 各上传组件目录显示正确、路径去域名化

---

## 附录：文件变更清单速查

| 阶段 | 文件路径 | 操作 |
|------|---------|------|
| 1.1 | `server/middleware/jwt.go` | 审查/加固 |
| 1.2 | `server/middleware/uni_response_protect.go` | 扩展字段 |
| 1.2 | `server/utils/i18n_util.go` | 扩展字段列表 |
| 1.3 | `server/middleware/uni_sign_verify.go` | 新建 |
| 1.3 | `server/initialize/router_biz.go` | 挂载中间件 |
| 1.4 | `server/utils/upload.go` | 新建/扩展 |
| 2.x | `uni/src/pages/learning/player.vue` | 多项改造 |
| 3.1 | `uni/src/pages/learning/home.vue` | 缓存改造 |
| 3.2 | `uni/src/pages/learning/typing.vue` | 交互升级 |
| 4.1 | `server/model/client/video_tag.go` | 新建 |
| 4.1 | `server/service/client/video_tag.go` | 新建 |
| 4.1 | `server/api/v1/client/video_tag.go` | 新建 |
| 4.1 | `server/router/client/video_tag.go` | 新建 |
| 4.1 | `web/src/view/client/videoTag/videoTag.vue` | 新建 |
| 4.3 | `server/model/client/external_link_domain.go` | 扩展字段 |
| 4.3 | `web/src/view/client/externalLinkDomain/externalLinkDomain.vue` | 扩展面板 |
| 4.4 | `web/src/components/FileUploadWithDir/index.vue` | 新建 |
