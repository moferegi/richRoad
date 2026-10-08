# 广告盈利系统设计规划 (AdPlan)

> 版本: v1.0 | 日期: 2026-08-23 | 状态: 规划阶段，未开始编码

---

## 一、概述

### 1.1 目标

在现有项目中新增广告盈利系统，通过展示广告视频/图片实现流量变现。系统分为三部分：

- **Uni 端（移动端）**：在指定页面位置展示"看广告"入口，弹出广告视频/图片，有效时间未到不可关闭。
- **Web 管理后台**：在"客户端"菜单下新增广告管理，运营人员配置广告内容和投放策略。
- **Server 端**：提供广告数据 API、视频切片上传、多云存储分发。

### 1.2 参考模块

| 参考模块 | 借鉴内容 |
|----------|----------|
| `layout/englishApp/englishLearningVideo` | 视频切片上传模式（HLS）、FFmpeg 切片流程 |
| `layout/client/externalLinkDomain` | 多云存储策略（默认上传 + 同步上传）、域名管理 |
| `server/service/example/exa_file_upload_download.go` | 文件上传集成多云存储的通用模式 |

---

## 二、数据库设计

### 2.1 广告位置表 `ad_positions`

定义广告展示的"点位"，一个位置对应一个页面上的广告入口（如 `pages/game/category` 的"看广告"按钮）。

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | uint | 主键 |
| `name` | varchar(100) | 位置名称（如"游戏分类页-看广告"） |
| `position_key` | varchar(200) | 位置标记（UNIQUE，如 `pages/game/category` 中的 `ad_watch_button`） |
| `description` | varchar(500) | 位置描述 |
| `is_enabled` | tinyint(1) | 是否启用该位置（0=关闭 1=开启），默认 1 |
| `sort` | int | 排序 |
| `created_at` | datetime | 创建时间 |
| `updated_at` | datetime | 更新时间 |

### 2.2 广告视频表 `ad_videos`

一条广告记录，可以是视频或图片。多个广告视频可绑定到同一个位置。

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | uint | 主键 |
| `position_id` | uint | 关联广告位置 ID（外键 `ad_positions.id`） |
| `title` | varchar(200) | 广告标题 |
| `media_type` | varchar(20) | 媒体类型：`video`（视频）、`image`（图片） |
| `media_url` | varchar(500) | 媒体资源相对路径（不含域名，运行时拼接） |
| `thumbnail_url` | varchar(500) | 封面缩略图（可选，用于视频的封面） |
| `duration` | int | 总时长（秒），视频用 ffprobe 自动获取，图片手动设置 |
| `min_watch_seconds` | int | 用户必须看完的秒数（最少观看时长），>=1 |
| `is_enabled` | tinyint(1) | 是否启用该视频（0=关闭 1=开启），默认 1 |
| `sort` | int | 排序（同一位置下多个广告按此排序播放） |
| `storage_key` | varchar(200) | 云存储目录标识（用于切片上传时的目录前缀） |
| `upload_folder` | varchar(200) | 上传目录（可修改，如 `ad/video/`） |
| `created_at` | datetime | 创建时间 |
| `updated_at` | datetime | 更新时间 |

### 2.3 广告观看记录表 `ad_watch_records`

记录用户观看广告的行为，用于统计和防刷。

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | uint | 主键 |
| `user_id` | uint | 用户 ID |
| `ad_video_id` | uint | 广告视频 ID |
| `position_id` | uint | 广告位置 ID |
| `watched_seconds` | int | 实际观看秒数 |
| `is_completed` | tinyint(1) | 是否达到最低观看时长 |
| `ip` | varchar(50) | 用户 IP |
| `created_at` | datetime | 观看时间 |

### 2.4 表关系

```
ad_positions (1) ──< (N) ad_videos ──< (N) ad_watch_records
```

---

## 三、Uni 端设计

### 3.1 "看广告"按钮

**位置**：`pages/game/category` 页面右上角 top-bar 区域。

现有布局：
```
[‹ 返回]    [游戏名称标题]    [占位]
```

修改后：
```
[‹ 返回]    [游戏名称标题]    [看广告]
```

**按钮行为**：
1. 进入页面时调用 API 查询该位置（`position_key = "pages/game/category"`）是否启用且有可用广告。
2. 如果有可用广告，显示"看广告"按钮；否则隐藏。
3. 点击按钮后弹出广告播放弹窗。

### 3.2 广告播放弹窗

**核心规则**：如果 `min_watch_seconds` 时间未到，弹窗**不可关闭**（无关闭按钮、无返回手势、无物理返回键）。

**弹窗设计**：
- 全屏或半屏模态弹窗，背景半透明遮罩。
- 视频类型：使用 `<video>` 组件播放，显示倒计时进度条。
- 图片类型：显示图片 + 倒计时进度条。
- 倒计时到 `min_watch_seconds` 后：
  - 显示"关闭"按钮或"✕"图标。
  - 如果还有剩余时长（视频总时长 > min_watch_seconds），用户可以手动关闭，也可以继续看完。
- 如果该位置有多个广告视频：
  - 按 `sort` 排序依次播放。
  - 当前广告看完（或达到最小观看时长后用户关闭）后，自动播放下一个。
  - 用户也可以点击"下一个"跳过（仅当当前广告已满足最小观看时长时）。
  - 全部播放完毕后弹窗关闭。

**API 调用**：
- `GET /ad/getAdByPosition?positionKey=pages/game/category` — 获取该位置的所有启用广告列表
- `POST /ad/reportWatch` — 上报观看记录（达到最小观看时长时调用）

### 3.3 通用化设计

"看广告"按钮和弹窗应封装为可复用的组件，以便后续在其他页面（如首页、结果页等）快速接入：

```
uni/src/components/ad-watch/
├── AdButton.vue        # "看广告"按钮组件
├── AdPlayer.vue        # 广告播放弹窗组件
└── adApi.js            # 广告 API 调用
```

---

## 四、Web 管理后台设计

### 4.1 菜单位置

在"客户端"（`client`）父菜单下新增子菜单，注册在 `server/initialize/new_modules_init.go`：

```go
menuDef{"adManage", "adManage", "view/client/ad/adManage.vue", "广告管理", "advertising", clientParent.ID, 22},
```

菜单路径：`客户端 > 广告管理`

### 4.2 页面功能设计

#### Tab 1：广告位置管理

| 功能 | 说明 |
|------|------|
| 列表展示 | 表格：位置名称、位置标记、启用状态、广告数量、排序、操作 |
| 新增/编辑位置 | 弹窗表单：名称、位置标记（唯一）、描述、启用开关、排序 |
| 删除 | 仅当该位置下无广告视频时可删除 |
| 启用/禁用 | 一键切换，禁用后 Uni 端不再显示该位置的"看广告"按钮 |

#### Tab 2：广告视频管理

| 功能 | 说明 |
|------|------|
| 筛选 | 按广告位置筛选 |
| 列表展示 | 表格：标题、媒体类型、预览缩略图、总时长、最低观看秒数、启用状态、排序、操作 |
| 新增/编辑 | 弹窗表单（详见下方） |
| 删除 | 删除广告视频（同时清理云存储文件，需确认） |
| 启用/禁用 | 单个视频可独立开关，关闭后 Uni 端跳过该视频 |
| 排序 | 拖拽或数字排序，决定 Uni 端播放顺序 |

#### 广告视频新增/编辑弹窗

| 字段 | 组件 | 说明 |
|------|------|------|
| 所属位置 | 下拉选择 | 选择 `ad_positions` 中的启用位置 |
| 标题 | 文本输入 | 广告标题 |
| 媒体类型 | 单选 | 视频 / 图片 |
| 上传媒体 | 文件上传 | 视频文件走切片上传（参照 englishLearningVideo），图片直接上传 |
| 媒体 URL | 文本输入 | 也支持手动输入外部 URL |
| 封面缩略图 | 图片上传 | 可选 |
| 切片上传目录 | 文本输入 | 可修改，默认 `ad/video/` |
| 总时长 | 数字输入 | 视频自动获取（ffprobe），图片手动设置 |
| 最低观看秒数 | 数字输入 | 必须 >=1，且 <= 总时长 |
| 启用状态 | 开关 | 控制该视频是否生效 |
| 排序 | 数字输入 | 同一位置下的播放顺序 |

#### Tab 3：观看记录/统计

| 功能 | 说明 |
|------|------|
| 观看记录列表 | 按用户、广告、位置筛选，显示观看时长、是否完成 |
| 统计概览 | 总观看次数、完成率、各位置/各广告的数据统计 |

---

## 五、Server 端设计

### 5.1 路由设计

路由前缀：`/ad`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| **管理端** | | | |
| POST | `/ad/createPosition` | 创建广告位置 | 管理端 |
| PUT | `/ad/updatePosition` | 更新广告位置 | 管理端 |
| DELETE | `/ad/deletePosition` | 删除广告位置 | 管理端 |
| GET | `/ad/getPositionList` | 获取广告位置列表 | 管理端 |
| POST | `/ad/createAdVideo` | 创建广告视频 | 管理端 |
| PUT | `/ad/updateAdVideo` | 更新广告视频 | 管理端 |
| DELETE | `/ad/deleteAdVideo` | 删除广告视频 | 管理端 |
| GET | `/ad/getAdVideoList` | 获取广告视频列表 | 管理端 |
| POST | `/ad/sliceAdVideo` | 视频切片上传（参照 englishLearningVideo 的 sliceVideoEpisode） | 管理端 |
| GET | `/ad/checkFfmpeg` | 检查 FFmpeg 是否可用 | 管理端 |
| GET | `/ad/getWatchRecordList` | 获取观看记录列表 | 管理端 |
| GET | `/ad/getWatchStats` | 获取观看统计 | 管理端 |
| **Uni 端** | | | |
| GET | `/ad/getAdByPosition` | 根据位置标记获取可用广告列表 | Uni 端 |
| POST | `/ad/reportWatch` | 上报观看记录 | Uni 端 |

### 5.2 文件结构

```
server/
├── model/client/
│   ├── ad.go                          # 数据模型：AdPosition, AdVideo, AdWatchRecord
│   └── request/
│       └── ad.go                      # 请求/响应结构体
├── api/v1/client/
│   └── ad.go                          # API 控制器
├── service/client/
│   ├── ad.go                          # 广告业务逻辑
│   └── ad_slice.go                    # 广告视频切片上传服务
├── router/client/
│   └── ad.go                          # 路由注册
└── initialize/
    └── new_modules_init.go            # 菜单/API/Casbin 权限注册（新增广告模块）
```

### 5.3 视频切片上传策略（参照 englishLearningVideo）

完全参照 `server/plugin/english_learning/service/video_slice.go` 的 `SliceAndUpload` 方法：

1. 前端上传视频文件到服务端临时目录。
2. 调用 `ffprobe` 获取视频总时长。
3. 调用 `ffmpeg` 进行 HLS 切片（`-c copy -hls_time 10`），生成 `.m3u8` + `.ts` 片段。
4. 调用 `CloudStorageService.GetDefaultUploadDomain()` 获取默认上传云。
5. 并发上传所有 `.ts` 和 `.m3u8` 文件到默认云存储。
6. 异步调用 `syncUploadToClouds()` 同步到其他标记了 `sync_upload=true` 的云。
7. 数据库存储相对路径（不含域名），Uni 端运行时通过 `getDefaultDomain` 接口拼接完整 URL。

**关键差异**：
- 广告视频不需要 AES-128 加密（HLS 加密仅用于保护付费内容，广告不需要）。
- 广告视频切片不需要 `hlsKey` 端点。
- 如果媒体类型是图片，则走普通文件上传（`UploadFileToCloud`），不切片。

### 5.4 多云存储策略（参照 externalLinkDomain）

与现有 `externalLinkDomain` 的 `UploadFile` 方法保持一致：

1. 默认上传云：优先 `defaultUpload=true`，其次 `isDefault=true`。
2. 同步上传云：所有 `syncUpload=true` 的云（排除默认云）。
3. 上传后数据库存储相对路径，运行时拼接域名。

### 5.5 API 注册（new_modules_init.go）

需要在 `initNewModulesApis` 中新增广告相关 API 记录，在 `initNewModulesMenus` 中新增广告菜单，在 `initNewModulesCasbin` 中新增权限。

---

## 六、改进建议

基于对现有系统的分析，提出以下可完善广告设计的建议：

### 6.1 广告频控策略

- **每日观看上限**：每个用户每天最多观看 N 次广告（可在 `ad_positions` 或全局配置中设置）。
- **冷却时间**：同一位置两次广告之间至少间隔 X 分钟，防止用户刷广告。
- 这些配置可以在管理后台的"广告位置"编辑中设置。

### 6.2 广告奖励机制

- 看完广告后给予用户奖励（如积分、试衣币、游戏道具等），激励用户主动观看广告。
- 在 `ad_positions` 中增加 `reward_type`（积分/试衣币/无）和 `reward_amount` 字段。
- 仅当 `watched_seconds >= min_watch_seconds` 时发放奖励。

### 6.3 广告类型扩展

- **激励视频广告**：用户主动点击观看，看完获得奖励（当前设计）。
- **插屏广告**：在页面切换或游戏结束时自动弹出（可扩展为 `ad_positions` 的 `trigger_type` 字段：`manual` 手动 / `auto` 自动）。
- **Banner 广告**：页面底部固定展示（可扩展为不同位置类型）。

### 6.4 广告数据统计增强

- 在 `ad_watch_records` 基础上增加统计看板：
  - 按天/周/月的广告曝光量、点击率、完播率。
  - 各广告位置的收益对比。
  - 用户观看行为分析（哪个时间段观看最多）。

### 6.5 广告跳转链接

- 广告视频/图片可配置跳转链接（`ad_videos` 表增加 `link_url` 字段）。
- 用户在达到最小观看时长后，可点击广告跳转到目标页面（如商品详情、外部 H5）。

### 6.6 A/B 测试支持

- 同一位置可配置多组广告，随机分组展示，统计各组转化率。
- 可扩展为 `ad_positions` 的 `ab_group` 字段。

### 6.7 广告素材预加载

- Uni 端在进入页面时预加载广告素材（视频第一帧/图片），减少弹窗打开时的等待时间。
- 可在 `AdButton` 组件挂载时触发预加载。

---

## 七、实施步骤（建议顺序）

| 阶段 | 内容 | 涉及模块 | 状态 |
|------|------|----------|------|
| 1 | 数据模型 + 自动迁移 | server/model | ✅ 已完成 |
| 2 | 广告位置 CRUD API | server/api, service, router | ✅ 已完成 |
| 3 | 广告视频 CRUD API + 切片上传 | server/api, service/router, ad_slice | ✅ 已完成 |
| 4 | 菜单注册 + Casbin 权限 | server/initialize | ✅ 已完成 |
| 5 | Web 管理后台页面 | web/src/view/client/ad/ | ✅ 已完成 |
| 6 | Uni 端广告组件 + API 对接 | uni/src/components/ad-watch/ | ✅ 已完成 |
| 7 | 观看记录上报 + 统计 | server + web + uni | ✅ 已完成 |
| 8 | 可选改进（频控、奖励、跳转等） | 全栈 | ⏳ 待实施 |

---

## 八、注意事项

1. **不要影响现有功能**：所有新增代码在独立目录/文件中，不修改现有模块的代码逻辑。
2. **上传策略一致性**：严格参照 `englishLearningVideo` 的切片模式 + `externalLinkDomain` 的多云存储，保持代码风格和架构一致。
3. **图片广告**：图片不走切片，直接上传。`duration` 字段手动设置（如 10 秒），`min_watch_seconds` 必须 <= duration。
4. **不可关闭机制**：Uni 端需同时处理：隐藏关闭按钮、拦截返回手势、拦截物理返回键。
5. **广告观看记录**：需要防刷，同一用户同一广告短时间内不应重复计数（可限制 XX 分钟内同广告不重复记录）。