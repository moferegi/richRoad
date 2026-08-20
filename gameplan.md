# 游戏模块开发方案（GamePlan）

> 版本：v1.0 | 日期：2026-08-12 | 阶段：仅先做"24闯关模式"

---

## 一、项目概述

在现有 `richRoad` 项目中新增游戏模块，首期实现"24点闯关模式"。用户进入首页自动注册（随机7位数字用户名，角色888），可游玩24点计算游戏。后续扩展"36点"和"密码破解"。

---

## 二、关键设计决策

### 2.1 计算策略：纯本地计算，闯关成功才请求接口

| 场景 | 策略 | 理由 |
|------|------|------|
| 生成题目（4个数字 + 目标结果） | **后端下发** | 保证题目可控、难度可配置、防作弊 |
| 用户操作过程（点击数字、符号、逐步计算） | **纯前端本地计算** | 零延迟体验，减少服务端压力 |
| 闯关成功 | **调用后端接口提交** | 记录成绩、更新排行榜 |
| 闯关失败/放弃 | **不请求接口** | 减少无效请求 |

### 2.2 Uni端接口限流

| 接口 | 限流策略 | 理由 |
|------|----------|------|
| 获取题目 | 每关每天最多请求 N 次（如20次） | 防止恶意刷题 |
| 提交闯关结果 | 同关同用户去重（已通过则拒绝重复提交） | 防重放攻击 |
| 排行榜查询 | 不做特殊限制 | 读操作无副作用 |

### 2.3 加密与签名联动

遵循现有 `learning_api_encrypt_enabled` 和 `learning_api_sign_enabled` 配置，**所有新增 Uni 端游戏接口均联动这两个开关**：

- 请求签名：通过 `X-Req-Ts` / `X-Req-Nonce` / `X-Req-Sign` 头
- 响应加密：通过 `X-Resp-Encrypt: 1` + `X-Client-Platform: uni` 头
- 后端中间件：复用 `UniSignVerify` + `UniResponseProtect`（或游戏模块独立中间件，但共享同一开关）

### 2.4 多语言裁剪

所有 Uni 端接口返回的多语言字段（如 `{"zh":"简单","en":"Easy"}`）在 `request.js` 的 `i18nNormalizeEndpoints` 列表中注册后自动裁剪。前端参照 `learning/profile.vue` 模式使用 `t()` 和 `localText()`。

---

## 三、数据库设计

### 3.1 游戏大分类表 `game_categories`

```sql
CREATE TABLE game_categories (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  game_key      VARCHAR(64)  NOT NULL UNIQUE COMMENT '游戏标识: 24point / 36point / password_crack',
  name          JSON         NOT NULL COMMENT '多语言名称 {"zh":"24闯关模式","en":"24 Game"}',
  description   JSON         NULL     COMMENT '多语言描述',
  sort          INT          DEFAULT 0 COMMENT '排序',
  status        TINYINT      DEFAULT 1 COMMENT '1=启用 0=禁用',
  created_at    DATETIME(3)  NULL,
  updated_at    DATETIME(3)  NULL,
  deleted_at    DATETIME(3)  NULL
);
```

### 3.2 游戏难度分类表 `game_difficulty_categories`

```sql
CREATE TABLE game_difficulty_categories (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  game_id       BIGINT UNSIGNED NOT NULL COMMENT '关联 game_categories.id',
  name          JSON         NOT NULL COMMENT '多语言名称 {"zh":"简单","en":"Easy"}',
  sort          INT          DEFAULT 0 COMMENT '排序（越高越难）',
  status        TINYINT      DEFAULT 1 COMMENT '1=启用 0=禁用',
  created_at    DATETIME(3)  NULL,
  updated_at    DATETIME(3)  NULL,
  deleted_at    DATETIME(3)  NULL,
  FOREIGN KEY (game_id) REFERENCES game_categories(id)
);
```

### 3.3 关卡表 `game_levels`

```sql
CREATE TABLE game_levels (
  id                  BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  category_id         BIGINT UNSIGNED NOT NULL COMMENT '关联 game_difficulty_categories.id',
  level_number        INT          NOT NULL COMMENT '该分类下第几关（从1开始）',
  numbers             VARCHAR(128) NOT NULL COMMENT '4个数字，逗号分隔 "3,8,3,8"',
  target_result       INT          NOT NULL DEFAULT 24 COMMENT '目标结果',
  sort                INT          DEFAULT 0,
  created_at          DATETIME(3)  NULL,
  updated_at          DATETIME(3)  NULL,
  deleted_at          DATETIME(3)  NULL,
  FOREIGN KEY (category_id) REFERENCES game_difficulty_categories(id),
  UNIQUE KEY uk_cat_level (category_id, level_number)
);
```

### 3.4 用户闯关记录表 `game_user_progress`

```sql
CREATE TABLE game_user_progress (
  id            BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id       BIGINT UNSIGNED NOT NULL COMMENT '关联 client_users.id',
  level_id      BIGINT UNSIGNED NOT NULL COMMENT '关联 game_levels.id',
  status        TINYINT      DEFAULT 1 COMMENT '1=已通过 0=未通过',
  passed_at     DATETIME(3)  NULL COMMENT '通过时间',
  attempt_count INT          DEFAULT 0 COMMENT '尝试次数',
  created_at    DATETIME(3)  NULL,
  updated_at    DATETIME(3)  NULL,
  FOREIGN KEY (user_id) REFERENCES client_users(id),
  FOREIGN KEY (level_id) REFERENCES game_levels(id),
  UNIQUE KEY uk_user_level (user_id, level_id)
);
```

### 3.5 排行榜设计

**不单独建表**，通过聚合查询 `game_user_progress` 实现：

- **累计闯关榜**：`SELECT user_id, COUNT(*) as total FROM game_user_progress WHERE status=1 GROUP BY user_id ORDER BY total DESC`
- **每日闯关榜**：同上 + `WHERE DATE(passed_at) = CURDATE()`

---

## 四、后端 API 设计

### 4.1 Uni 端接口（需加密+签名联动）

| 方法 | 路径 | 说明 | 限流 |
|------|------|------|------|
| GET | `/game/getGameList` | 获取游戏列表（含大分类） | 无 |
| GET | `/game/getDifficultyCategories` | 获取某游戏的难度分类列表 | 无 |
| GET | `/game/getLevelDetail` | 获取关卡详情（含4个数字） | 每关每天N次 |
| POST | `/game/submitLevelResult` | 提交闯关结果 | 去重 |
| GET | `/game/getUserProgress` | 获取用户闯关进度 | 无 |
| GET | `/game/getLeaderboard` | 获取排行榜（累计/每日） | 无 |

### 4.2 Web 管理端接口（需 Casbin 权限）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/game/admin/createCategory` | 创建游戏大分类 |
| PUT | `/game/admin/updateCategory` | 更新游戏大分类 |
| DELETE | `/game/admin/deleteCategory` | 删除游戏大分类 |
| GET | `/game/admin/getCategoryList` | 获取游戏大分类列表 |
| POST | `/game/admin/createDifficultyCategory` | 创建难度分类 |
| PUT | `/game/admin/updateDifficultyCategory` | 更新难度分类 |
| DELETE | `/game/admin/deleteDifficultyCategory` | 删除难度分类 |
| GET | `/game/admin/getDifficultyCategoryList` | 获取难度分类列表 |
| POST | `/game/admin/createLevel` | 创建关卡 |
| PUT | `/game/admin/updateLevel` | 更新关卡 |
| DELETE | `/game/admin/deleteLevel` | 删除关卡 |
| GET | `/game/admin/getLevelList` | 获取关卡列表 |
| GET | `/game/admin/getLeaderboard` | 管理端排行榜 |
| GET | `/game/admin/getUserProgress` | 管理端用户进度 |

### 4.3 提交闯关结果的安全校验

```go
// submitLevelResult 防作弊逻辑
func (s *GameService) SubmitLevelResult(userID, levelID uint, ...) error {
    // 1. 校验关卡是否存在
    // 2. 校验同一用户同一关是否已通过（去重）
    // 3. 校验前置关卡是否已通过（顺序闯关）
    // 4. 记录通过
}
```

---

## 五、Uni 前端设计

### 5.1 页面路由

在 `uni/src/pages.json` 中新增：

```json
{
  "pages": [
    { "path": "pages/game/index", "style": { "navigationStyle": "custom" } },
    { "path": "pages/game/category", "style": { "navigationStyle": "custom" } },
    { "path": "pages/game/play", "style": { "navigationStyle": "custom" } }
  ]
}
```

### 5.2 页面结构

```
pages/game/
├── index.vue          # 游戏首页（游戏列表 + 排行榜）
├── category.vue       # 分类选择页（左右滑动切换分类）
├── play.vue           # 游戏进行页（24点计算）
└── components/
    ├── GameCard.vue       # 游戏卡片组件
    ├── CategorySlider.vue # 分类滑动组件
    ├── LevelGrid.vue      # 关卡网格（锁/解锁状态）
    ├── CalcBoard.vue      # 计算面板（4个数字 + 运算符）
    └── LeaderBoard.vue    # 排行榜组件
```

### 5.3 游戏首页（game/index）

**功能**：
- 进入即自动注册（调用 `/clientUser/autoRegister` 或复用现有注册逻辑，随机7位数字用户名，角色888）
- 顶部显示应用名称（`app_name` 配置，文艺字体）
- 设置按钮 → 弹出设置面板：
  - 音效开关（click-click.mp3，点击时播放）
  - 背景音乐开关（back-music.mp3，循环播放）
  - 用户名显示
  - 语言切换（复用 `<lang-switch>` 组件）
- 游戏列表（卡片式）：
  - "24闯关模式"（本期实现）
  - "36闯关36点"（灰显/锁，后续）
  - "密码破解"（灰显/锁，后续）
- 累计闯关榜 + 每日闯关榜（显示名称 + 总关数）

### 5.4 分类选择页（game/category）

**功能**：
- 左右滑动切换难度分类（Swiper）
- 下方圆点指示器（几个分类几个点）
- 每个分类下显示关卡网格：
  - 第一关：可点击"开始游戏"
  - 后续关卡：未通过时显示锁头图标，通过后解锁
- 点击已解锁关卡 → 进入 `game/play`

### 5.5 游戏进行页（game/play）—— 24点计算

**核心玩法**：
1. 显示4个数字卡片（如 3, 8, 3, 8）
2. 用户点击数字1 → 点击运算符（+ - ×）→ 点击数字2
3. 前端本地计算：数字1 op 数字2 → 结果
4. 动画：数字2消失，数字1位置显示计算结果，背景色区分已计算过的数字
5. 重复直到只剩1个数字
6. 最终结果 = 24 → 弹窗"闯关成功！"，调用 `/game/submitLevelResult` 提交
7. 最终结果 ≠ 24 → 弹窗"再试一次"，可重置

**技术细节**：
- 所有计算纯前端，使用 `reactive` 管理状态
- 动画使用 CSS transition / uni-app 动画API
- 音效：点击时播放 `click-click.mp3`（如果开启）
- 背景音乐：进入页面开始循环播放 `back-music.mp3`（如果开启）

### 5.6 多语言实现

**所有新增 Uni 页面遵循现有模式**：

```js
// 每个页面统一封装
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'

const langStore = useLangStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')
const t = (key) => { const text = i18nT(key, locale.value); return text !== key ? text : key }
const localText = (value) => i18nLocalText(value, locale.value)
```

**翻译文件**：在 `uni/src/utils/i18n-locales/` 各语言文件中新增游戏相关 key。

### 5.7 接口加密签名联动

在 `uni/src/api/` 新建 `game.js`，所有 API 通过 `request.js` 发送，自动联动加密和签名。新增接口路径需注册到 `request.js` 的 `i18nNormalizeEndpoints` 列表中。

---

## 六、Web 管理后台设计

### 6.1 菜单结构

在"英语端"父菜单（`englishApp`）同级新增"游戏设置"父菜单：

```
游戏设置 (gameSettings)
├── 游戏分类 (gameCategory)     → 管理 game_categories
├── 难度分类 (gameDifficulty)  → 管理 game_difficulty_categories
├── 关卡管理 (gameLevel)        → 管理 game_levels
├── 累计闯关榜 (gameLeaderboardAll)
└── 每日闯关榜 (gameLeaderboardDaily)
```

### 6.2 初始化方式

在 `server/initialize/new_modules_init.go` 中新增游戏模块的 API、菜单、Casbin 权限初始化，遵循现有模式。

### 6.3 权限

- 所有游戏管理 API 自动授予 888 角色
- 在 `initNewModulesCasbin` 中注册

---

## 七、实施步骤

### 阶段一：后端基础设施

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 创建数据模型 | `server/model/client/game.go` |
| 2 | 注册数据库迁移 | `server/initialize/gorm_biz.go` |
| 3 | 创建 Service 层 | `server/service/client/game.go` |
| 4 | 创建 API 处理器 | `server/api/v1/client/game.go` |
| 5 | 创建 Router | `server/router/client/game.go` |
| 6 | 注册路由 | `server/initialize/router_biz.go` |
| 7 | 注册 API/菜单/Casbin | `server/initialize/new_modules_init.go` |
| 8 | 种子数据（默认游戏分类） | `server/initialize/new_modules_init.go` |

### 阶段二：Web 管理后台

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 游戏分类管理页 | `web/src/view/client/game/category.vue` |
| 2 | 难度分类管理页 | `web/src/view/client/game/difficulty.vue` |
| 3 | 关卡管理页 | `web/src/view/client/game/level.vue` |
| 4 | 排行榜页 | `web/src/view/client/game/leaderboard.vue` |
| 5 | API 封装 | `web/src/api/client/game.js` |

### 阶段三：Uni 前端

| 步骤 | 内容 | 文件 |
|------|------|------|
| 1 | 多语言翻译文件 | `uni/src/utils/i18n-locales/*.js` |
| 2 | API 封装 | `uni/src/api/game.js` |
| 3 | 游戏首页 | `uni/src/pages/game/index.vue` |
| 4 | 分类选择页 | `uni/src/pages/game/category.vue` |
| 5 | 游戏进行页 | `uni/src/pages/game/play.vue` |
| 6 | 子组件 | `uni/src/pages/game/components/*.vue` |
| 7 | 路由注册 | `uni/src/pages.json` |

### 阶段四：联调与验收

| 步骤 | 内容 |
|------|------|
| 1 | Uni 端加密签名联动测试 |
| 2 | 多语言裁剪验证 |
| 3 | 排行榜数据验证 |
| 4 | 闯关流程端到端测试 |

---

## 八、优化建议

### 8.1 24点题目生成算法

**后端生成**：随机4个数字（1-13），用回溯法验证是否有解（有解才入库），确保每关都有可行解。

### 8.2 前端计算防作弊

虽然是本地计算，但提交时后端应重新校验：
- 后端根据关卡ID查出4个数字
- 不验证具体计算过程（成本高）
- 仅做通过标记和去重

### 8.3 音效管理

封装一个 `useGameAudio` composable：
```js
// 全局音频管理
const useGameAudio = () => {
  const clickAudio = uni.createInnerAudioContext()
  const bgMusic = uni.createInnerAudioContext()
  // 读取 localStorage 中的开关状态
  // 提供 playClick / stopBgMusic 等方法
}
```

### 8.4 自动注册

首页 `onMounted` 时检查是否已登录，未登录则自动注册：
```js
if (!userStore.token) {
  const username = String(Math.floor(1000000 + Math.random() * 9000000)) // 7位随机数字
  await autoRegister({ username, password: '123456', authorityId: 888 })
  // 自动登录
}
```

### 8.5 后端游戏模块不采用插件模式

与英语学习模块不同，游戏模块建议直接放在主项目 `server/` 下的 `api/v1/client/`、`service/client/`、`router/client/` 中，而非 `server/plugin/` 下。因为：
- 游戏模块与主项目紧密耦合（共享用户表、自注册、权限888）
- 不需要独立的热插拔能力
- 遵循现有 client 模块的代码组织方式

---

## 九、不涉及的内容（明确边界）

- 不修改任何现有模块代码（learning、shop、tryon 等）
- 不新增数据库（使用现有数据库）
- 不修改现有中间件逻辑
- 不修改现有用户注册/登录流程（仅新增自动注册入口）
- 36点、密码破解游戏本期不做

---

## 十、文件清单总览

```
server/
├── model/client/game.go                          # [新增] 数据模型
├── service/client/game.go                        # [新增] 业务逻辑
├── service/client/enter.go                       # [修改] 嵌入 GameService
├── api/v1/client/game.go                         # [新增] API 处理器
├── api/v1/client/enter.go                        # [修改] 嵌入 GameApi
├── router/client/game.go                         # [新增] 路由
├── router/client/enter.go                        # [修改] 嵌入 GameRouter
├── initialize/gorm_biz.go                        # [修改] 添加 AutoMigrate
├── initialize/router_biz.go                      # [修改] 注册游戏路由
├── initialize/new_modules_init.go                # [修改] 注册 API/菜单/Casbin

web/
├── src/api/client/game.js                        # [新增] Web API
├── src/view/client/game/category.vue             # [新增] 游戏分类管理
├── src/view/client/game/difficulty.vue           # [新增] 难度分类管理
├── src/view/client/game/level.vue                # [新增] 关卡管理
├── src/view/client/game/leaderboard.vue          # [新增] 排行榜

uni/
├── src/api/game.js                               # [新增] Uni API
├── src/pages/game/index.vue                      # [新增] 游戏首页
├── src/pages/game/category.vue                   # [新增] 分类选择页
├── src/pages/game/play.vue                       # [新增] 游戏进行页
├── src/pages/game/components/GameCard.vue        # [新增] 游戏卡片
├── src/pages/game/components/CategorySlider.vue  # [新增] 分类滑动
├── src/pages/game/components/LevelGrid.vue       # [新增] 关卡网格
├── src/pages/game/components/CalcBoard.vue       # [新增] 计算面板
├── src/pages/game/components/LeaderBoard.vue     # [新增] 排行榜
├── src/pages.json                                # [修改] 添加路由
├── src/utils/i18n-locales/zh.js                  # [修改] 添加翻译
├── src/utils/i18n-locales/en.js                  # [修改] 添加翻译
├── src/utils/i18n-locales/...(其他语言)           # [修改] 添加翻译
└── src/utils/request.js                          # [修改] 注册 i18nNormalizeEndpoints
```