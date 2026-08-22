# 密码推理游戏 (pwd-guess) 开发计划

## 一、游戏玩法设计

### 核心规则
- 正确答案是 **4 位数字**（如 `3847`），数字和位置必须完全匹配才过关
- 每关提供 **4 个提示**，每个提示由两部分组成：
  - **固定数字**：4 个提示数字（始终显示在界面上方，如 `1 2 3 4`）
  - **文字提示**：每条提示对应一个数字，给出推理线索（如"其中两个数字正确但位置不对"）
- 玩家需要用 **3 种颜色标记** 辅助推理，标记提示数字以区分不同类别

### 界面交互
- **提示区**：4 个固定提示数字 + 4 条文字提示
- **颜色标记**：3 个颜色块（红/黄/绿），点击选中颜色 → 点击提示数字应用颜色背景 → 再次点击取消
- **数字键盘**：0-9 数字键 + 删除键，输入 4 位数字
- **确认按钮**：输入完 4 位后点击确认，判断是否正确

---

## 二、数据库设计

### 新增模型 `PwdGameLevel`（表名 `pwd_game_levels`）

```go
type PwdGameLevel struct {
    global.GVA_MODEL
    CategoryID  uint   `json:"categoryID" gorm:"column:category_id;index;comment:关联 game_difficulty_categories.id"`
    LevelNumber int    `json:"levelNumber" gorm:"column:level_number;comment:该分类下第几关"`
    Answer      string `json:"answer" gorm:"column:answer;size:4;comment:正确答案4位数字"`
    HintDigits  string `json:"hintDigits" gorm:"column:hint_digits;size:128;comment:4个提示数字逗号分隔"`
    HintTexts   string `json:"hintTexts" gorm:"column:hint_texts;type:json;comment:4条文字提示，多语言JSON数组"`
    Sort        int    `json:"sort" gorm:"column:sort;default:0;comment:排序"`
}
```

- `Answer`：`"3847"` 4 位数字字符串
- `HintDigits`：`"1,2,3,4"` 逗号分隔，始终显示
- `HintTexts`：JSON 多语言数组，如 `["hint1_zh|hint1_en", "hint2_zh|hint2_en", ...]` 或 `[{"zh":"...","en":"..."}, ...]`

### 复用现有模型（无需修改）
- `GameCategory` — 游戏分类，新增 `gameKey = "pwd-guess"` 一条
- `GameDifficultyCategory` — 难度分类，关联 `gameID`
- `GameUserProgress` — 用户进度，复用 `LevelID` 指向 `PwdGameLevel.ID`

---

## 三、后端改动

### 3.1 模型层 `server/model/client/game.go`
- 新增 `PwdGameLevel` 结构体
- 自动迁移加入 `PwdGameLevel`

### 3.2 服务层 `server/service/client/game.go`
新增方法：
- `CreatePwdLevel(level *PwdGameLevel) error`
- `UpdatePwdLevel(id uint, data map[string]interface{}) error`
- `DeletePwdLevel(id uint) error` — 级联删除进度
- `GetPwdLevelList(categoryID uint) ([]PwdGameLevel, error)`
- `GetPwdLevelListAdmin(categoryID uint, page, pageSize int) ([]PwdGameLevel, int64, error)`
- `GetPwdLevelByID(id uint) (*PwdGameLevel, error)`
- `SubmitPwdLevelResult(userID, levelID uint, success bool) error` — 复用 SubmitLevelResult 逻辑，但 LevelID 指向 PwdGameLevel

### 3.3 API 层 `server/api/v1/client/game.go`
新增 Uni 端接口：
- `GET /game/getPwdLevelDetail?levelID=` — 获取密码关卡详情
- `GET /game/getPwdLevelList?categoryID=` — 获取密码关卡列表
- `POST /game/submitPwdLevelResult` — 提交密码关卡结果

新增管理端接口：
- `POST /game/admin/createPwdLevel`
- `PUT /game/admin/updatePwdLevel`
- `DELETE /game/admin/deletePwdLevel`
- `GET /game/admin/getPwdLevelList`

### 3.4 路由 `server/router/client/game.go`
- 注册上述新路由到对应的 `gameRouter` / `gameAdminRouter`

### 3.5 菜单与权限 `server/initialize/new_modules_init.go`
- API 声明新增 7 条
- 菜单新增 "密码关卡管理"（`view/client/game/pwdLevel.vue`），排序 14.5
- Casbin 策略新增对应 API 路径

### 3.6 种子数据
- `initGameSeedData` 中 `password_crack`（status=0）改为 `pwd-guess`（status=1 启用）

---

## 四、Uni 前端

### 4.1 新增页面 `uni/src/pages/game/pwd-play.vue`
- 路由参数：`levelID`、`catID`
- 顶部栏：返回按钮、难度名称、关卡号
- **提示区**：4 个固定数字 + 4 条对应文字提示卡片
- **颜色标记区**：3 个颜色块（红 `#e94560`、黄 `#f0a500`、绿 `#00b894`），点击选中
- **数字键盘**：0-9 网格布局 + 删除键 + 确认按钮
- **输入显示区**：4 个方框显示已输入的数字
- 成功/失败弹窗
- 多语言：通过 `t()` 和 `localText()`
- 音效：复用 `initAudio` + `playClick`
- 页面配置自动加入 `pages.json`

### 4.2 新增 API `uni/src/api/game.js`
```js
export const getPwdLevelDetail = (levelID) => request({ url: '/game/getPwdLevelDetail', method: 'GET', params: { levelID } })
export const getPwdLevelList = (categoryID) => request({ url: '/game/getPwdLevelList', method: 'GET', params: { categoryID } })
export const submitPwdLevelResult = (levelID, answer) => request({ url: '/game/submitPwdLevelResult', method: 'POST', data: { levelID, answer } })
```

### 4.3 路由跳转
- `category.vue` 中，根据 `gameKey` 判断跳转 `play.vue` 还是 `pwd-play.vue`
- `index.vue` 中，游戏卡片点击时根据 `gameKey` 传递参数

### 4.4 多语言 `uni/src/utils/i18n-locales/zh.js` + `en.js`
新增键：
```
pwdGame.title: '密码推理'
pwdGame.hints: '提示'
pwdGame.inputAnswer: '输入答案'
pwdGame.confirm: '确认'
pwdGame.clear: '清除'
pwdGame.colorHint: '点击颜色标记提示数字'
pwdGame.delete: '删除'
```

---

## 五、Web 管理端

### 5.1 新增页面 `web/src/view/client/game/pwdLevel.vue`
- 两级联动筛选：游戏 → 难度分类
- 表格列：ID、关卡序号、答案（4位数字）、提示数字、提示文字（多语言）、排序
- 弹窗表单：关卡序号、答案（4位数字输入）、提示数字（4个逗号分隔）、提示文字（MultiLangEditor × 4）
- 参考 `level.vue` 结构

### 5.2 新增 API `web/src/api/client/game.js`
```js
export const createPwdLevel = (data) => service.post('/game/admin/createPwdLevel', data)
export const updatePwdLevel = (data) => service.put('/game/admin/updatePwdLevel', data)
export const deletePwdLevel = (data) => service.delete('/game/admin/deletePwdLevel', { data })
export const getPwdLevelList = (params) => service.get('/game/admin/getPwdLevelList', { params })
```

---

## 六、多语言后端处理

### 6.1 `HintTexts` 字段的 i18n 裁剪
- `HintTexts` 字段名加入 `i18n_payload.go` 的字段白名单
- `UniResponseProtect` 中间件自动裁剪 `hintTexts` 为当前语言

### 6.2 字段白名单 `server/utils/i18n_payload.go`
- 新增 `"hintTexts"` 到白名单

---

## 七、API 加密

### 7.1 Uni 端请求
- 请求自动携带 `X-Client-Platform: uni` 头
- 若 `learning_api_sign_enabled` 开启，请求自动签名
- 若 `learning_api_encrypt_enabled` 开启，响应自动 AES 加密

### 7.2 中间件
- `UniSignVerify` — 校验 Uni 请求签名（已有，无需修改）
- `UniResponseProtect` — i18n 裁剪 + AES 加密响应（已有，无需修改）

---

## 八、开发步骤

| 步骤 | 文件 | 说明 |
|------|------|------|
| 1 | `server/model/client/game.go` | 新增 `PwdGameLevel` 模型 |
| 2 | `server/service/client/game.go` | 新增密码关卡 CRUD 服务方法 |
| 3 | `server/api/v1/client/game.go` | 新增 Uni + Admin API 处理器 |
| 4 | `server/router/client/game.go` | 注册新路由 |
| 5 | `server/initialize/new_modules_init.go` | API 声明 + 菜单 + Casbin + 种子数据 |
| 6 | `server/utils/i18n_payload.go` | 白名单加 `hintTexts` |
| 7 | `uni/src/api/game.js` | 新增 3 个 Uni API |
| 8 | `uni/src/utils/i18n-locales/zh.js` + `en.js` | 新增 i18n 键 |
| 9 | `uni/src/pages/game/pwd-play.vue` | 密码推理游戏页面 |
| 10 | `uni/src/pages/game/category.vue` | gameKey 判断跳转页面 |
| 11 | `uni/src/pages.json` | 注册 pwd-play 页面 |
| 12 | `web/src/api/client/game.js` | 新增 4 个管理端 API |
| 13 | `web/src/view/client/game/pwdLevel.vue` | 密码关卡管理页面 |
| 14 | 编译测试 | 后端 + Uni + Web 全链路验证 |