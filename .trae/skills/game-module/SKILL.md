---
name: "game-module"
description: "游戏模块开发指南：新增游戏分类、关卡管理、多语言设计、API 加密签名。当需要新增游戏玩法或理解游戏模块架构时使用。"
---

# 游戏模块开发指南

## 一、架构总览

游戏模块采用三层级联模型，基于 `gameKey` 区分不同玩法：

```
GameCategory (game_categories)          ← gameKey 区分玩法（24point/pwd-guess/...）
    └── GameDifficultyCategory          ← gameID 关联，难度层级
            └── GameLevel (game_levels) ← 统一关卡表，game_data JSON 存玩法特定数据
                    └── GameUserProgress (game_user_progress) ← 通用进度表
```

**关键文件清单：**

| 层级 | 文件路径 | 说明 |
|------|----------|------|
| 数据模型 | `server/model/client/game.go` | GameCategory / GameDifficultyCategory / GameLevel / GameUserProgress / PwdGameLevel |
| 业务逻辑 | `server/service/client/game.go` | CRUD + 闯关 + 排行榜 |
| API 接口 | `server/api/v1/client/game.go` | Uni 端 + 管理端 Handler |
| 路由注册 | `server/router/client/game.go` | 路由分组 |
| 数据库迁移 | `server/initialize/gorm_biz.go` | 表注册 |
| 种子数据 | `server/initialize/new_modules_init.go` | `initGameSeedData` 函数 |
| 菜单定义 | `server/initialize/new_modules_init.go` | englishApp 下的游戏菜单 |
| Uni API 封装 | `uni/src/api/game.js` | Uni 端 API 请求函数 |
| Uni 游戏首页 | `uni/src/pages/game/index.vue` | 游戏列表 + 排行榜 |
| Uni 分类页 | `uni/src/pages/game/category.vue` | 难度选择 + 关卡网格 |
| Uni 游戏页 | `uni/src/pages/game/play.vue` | 24点计算游戏 |
| Uni 密码页 | `uni/src/pages/game/pwd-play.vue` | 密码推理游戏 |
| Web 管理端 API | `web/src/api/client/game.js` | 管理端 API 封装 |
| Web 分类管理 | `web/src/view/client/game/category.vue` | 游戏分类管理页 |
| Web 难度管理 | `web/src/view/client/game/difficulty.vue` | 难度分类管理页 |
| Web 关卡管理 | `web/src/view/client/game/level.vue` | 普通关卡管理页 |
| Web 密码关卡 | `web/src/view/client/game/pwdLevel.vue` | 密码关卡管理页 |

---

## 二、新增游戏分类的标准流程

### 步骤 1：种子数据（`server/initialize/new_modules_init.go`）

在 `initGameSeedData` 函数的 `cats` 切片中添加新游戏：

```go
{GameKey: "sudoku", Name: `{"zh":"数独挑战","en":"Sudoku",...}`, Sort: 4, Status: 1},
```

必须包含 12 种语言：`zh`, `en`, `mn`, `th`, `hi`, `id`, `vi`, `ar`, `ja`, `ko`, `ms`。

### 步骤 2：数据库模型（`server/model/client/game.go`）

**推荐方案：统一关卡表**。在 `GameLevel` 中加 `GameData` JSON 字段，避免为每个玩法建独立表：

```go
type GameLevel struct {
    global.GVA_MODEL
    CategoryID   uint   `json:"categoryID" gorm:"column:category_id;index;"`
    LevelNumber  int    `json:"levelNumber" gorm:"column:level_number;"`
    Sort         int    `json:"sort" gorm:"column:sort;default:0;"`
    GameData     string `json:"gameData" gorm:"column:game_data;type:json;comment:玩法数据JSON;"` // 新增
}
```

`game_data` 示例：
- 24点：`{"type":"calc","numbers":"3,8,3,8","targetResult":24}`
- 密码推理：`{"type":"pwd","answer":"1234","hintDigits":"1,2,3,4","hintTexts":[{...}]}`
- 数独：`{"type":"sudoku","puzzle":"530070000...","solution":"534678912..."}`

### 步骤 3：Uni 前端游戏页面

#### 3.1 API 封装（`uni/src/api/game.js`）

```js
// 统一关卡详情接口（不需要为每种玩法单独定义）
export const getLevelDetail = (levelID) => {
  return request({ url: '/game/getLevelDetail', method: 'GET', params: { levelID } })
}
export const submitLevelResult = (levelID) => {
  return request({ url: '/game/submitLevelResult', method: 'POST', params: { levelID } })
}
```

#### 3.2 新建游戏页面（`uni/src/pages/game/xxx-play.vue`）

从 `game_data` JSON 中解析玩法特定数据渲染 UI。参考 `play.vue` 和 `pwd-play.vue` 的模板结构。

可提取公共逻辑为 composable（`useGamePlay.js`）：

```js
// uni/src/composables/useGamePlay.js
export function useGamePlay() {
  // 公共逻辑：loadLevel, submitResult, goNextLevel, initAudio, playClick, goBack
}
```

#### 3.3 路由分发（`uni/src/pages/game/category.vue`）

**优化方案**：在 `GameCategory` 模型中增加 `playPage` 字段，前端直接读取跳转，不再硬编码 `gameKey` 判断：

```javascript
// 原代码（硬编码）：
const target = gameKey.value === 'pwd-guess' ? 'pwd-play' : 'play'

// 优化后（数据驱动）：
const target = game.playPage || 'play'
uni.navigateTo({ url: `/pages/game/${target}?levelID=${level.ID}&catID=${cat.ID}` })
```

### 步骤 4：管理端菜单（`server/initialize/new_modules_init.go`）

在 `englishAppParent` 菜单组中添加：

```go
menuDef{"gameXxxLevel", "gameXxxLevel", "view/client/game/xxxLevel.vue", "XXX关卡管理", "grid", englishAppParent.ID, 19},
```

### 步骤 5：管理端关卡管理页面（`web/src/view/client/game/xxxLevel.vue`）

参考 `level.vue` 或 `pwdLevel.vue` 创建，使用 `MultiLangEditor` 组件处理多语言字段。

---

## 三、多语言字段设计规范

### 3.1 数据库存储格式

多语言字段在 Go 模型中定义为 `type:json` 或 `type:text`，数据库存储 JSON 字符串：

```go
Name string `json:"name" gorm:"column:name;type:json;comment:多语言名称;"`
```

JSON 格式（12 种语言）：

```json
{"zh":"中文","en":"English","mn":"Монгол","th":"ไทย","hi":"हिन्दी","id":"Bahasa","vi":"Tiếng Việt","ar":"العربية","ja":"日本語","ko":"한국어","ms":"Melayu"}
```

### 3.2 Uni 前端解析

使用 `localText` 函数（`uni/src/utils/i18n.js`）：

```js
import { localText } from '@/utils/i18n.js'

const displayName = localText(item.name) // 自动取当前语言，降级链：当前语言 → en → zh → 任意值
```

### 3.3 Uni 后端自动裁剪

**中间件：** `server/middleware/uni_response_protect.go` 中的 `UniResponseProtect` 中间件挂载在所有路由组上，当检测到 `X-Client-Platform` 为 `uni`/`uniapp`/`uni-app` 时，自动调用 `LocalizeI18nPayloadByContext` 裁剪响应中的多语言字段。

**白名单字段**（`server/utils/i18n_payload.go` 第 12-32 行）：

```go
var i18nPayloadFieldKeys = map[string]struct{}{
    "name": {}, "nameI18n": {}, "title": {}, "titleI18n": {},
    "content": {}, "contentI18n": {}, "countryName": {},
    "seriesName": {}, "episodeName": {}, "categoryName": {},
    "chapterName": {}, "tagName": {}, "description": {},
    "explanation": {}, "translate": {}, "announcementContent": {},
    "popupTitle": {}, "popupContent": {}, "bannerTitle": {},
}
```

**新增多语言字段时**，如果字段名不在白名单中，需要在此 map 中添加。

**降级逻辑**（`pickI18nText`）：请求语言 → 主语言 → 繁体中文 → 英语 → 中文 → 任意非空值。

### 3.4 Web 管理端多语言编辑器

使用 `MultiLangEditor` 组件：

```html
<MultiLangEditor :model="formData.name" title="名称多语言" />
```

提交时序列化为 JSON 字符串：

```js
formData.name = JSON.stringify(formData.name)
```

---

## 四、API 加密与签名配置

### 4.1 配置项

在 `server/initialize/other.go` 中定义了两个开关：

| 配置 Key | 默认值 | 说明 |
|----------|--------|------|
| `learning_api_encrypt_enabled` | `true` | 响应加密开关，关闭后 Uni 端返回明文 JSON |
| `learning_api_sign_enabled` | `true` | 请求签名校验开关，关闭后跳过签名验证 |

**读取优先级**：sysConfig 数据库值 > 环境变量 > 默认值（`server/utils/runtime_settings.go`）。

### 4.2 签名校验（`server/middleware/uni_sign_verify.go`）

- 通过 `isUniSignEnabled()` 读取 `learning_api_sign_enabled` 开关
- 仅对 `X-Client-Platform=uni` 的请求生效
- 校验三个请求头：`X-Req-Ts`（时间戳，±5分钟窗口）、`X-Req-Nonce`（一次性随机数）、`X-Req-Sign`（HMAC-SHA256 签名）
- 签名算法：`HMAC-SHA256(SHA256(token|salt), METHOD|path|query|ts|nonce)`

### 4.3 响应加密（`server/middleware/uni_response_protect.go`）

- 通过 `isUniEncryptEnabled()` 读取 `learning_api_encrypt_enabled` 开关
- 触发条件：`learning_api_encrypt_enabled=true` + `X-Resp-Encrypt=1` 请求头 + 响应 `code=0`
- 加密算法：**AES-256-CBC**，密钥由 `SHA256(token|uni-api-v1)` 派生
- 加密后响应格式：`{"data":{"alg":"aes-cbc","payload":"base64...","iv":"base64...","ts":...},"_enc":"aes-cbc","_enc_ver":1}`

**重要：i18n 裁剪不受加密开关影响**，只要检测到 Uni 平台请求就会执行裁剪（第 77 行）。

### 4.4 调试时关闭

开发调试时可以通过数据库修改 sysConfig 关闭：

```sql
UPDATE sys_configs SET config_value = 'false' WHERE config_key = 'learning_api_encrypt_enabled';
UPDATE sys_configs SET config_value = 'false' WHERE config_key = 'learning_api_sign_enabled';
```

或设置环境变量：

```bash
LEARNING_API_ENCRYPT_ENABLED=false
LEARNING_API_SIGN_ENABLED=false
```

---

## 五、优化后的通用 Service 模式

### 统一闯关提交

代替为每个玩法写独立的 `SubmitXxxLevelResult`，统一方法：

```go
func (s *GameService) SubmitLevelResult(userID uint, levelID uint) error {
    // 1. 校验关卡存在
    var level client.GameLevel
    if err := global.GVA_DB.First(&level, levelID).Error; err != nil {
        return err
    }
    // 2. 检查前置关卡（同 category_id 内 level_number 更小的关卡）
    // 3. 幂等检查（已通过则直接返回成功）
    // 4. 保存 GameUserProgress
}
```

### 统一关卡 CRUD

新增玩法只需在 `game_data` JSON 中定义数据格式，前端根据 `type` 字段渲染不同 UI。后端 Service 层不需要任何改动。

---

## 六、新增游戏分类 Checklist

- [ ] `server/initialize/new_modules_init.go` → `initGameSeedData` 添加种子数据
- [ ] `server/model/client/game.go` → 如需新字段，在 `GameLevel.GameData` 中定义 JSON 结构
- [ ] `server/utils/i18n_payload.go` → 如有新多语言字段名，加入白名单
- [ ] `uni/src/pages/game/xxx-play.vue` → 新建游戏页面，解析 `game_data` 渲染
- [ ] `uni/src/pages/game/category.vue` → 添加路由分发（或改用 `playPage` 字段）
- [ ] `uni/src/api/game.js` → 如需新 API，添加请求函数
- [ ] `web/src/view/client/game/xxxLevel.vue` → 新建管理端关卡管理页
- [ ] `server/initialize/new_modules_init.go` → 添加管理端菜单项
- [ ] `web/src/api/client/game.js` → 如需新管理端 API，添加请求函数
- [ ] `uni/src/utils/i18n-locales/zh.js` → 添加新游戏相关的中文文案