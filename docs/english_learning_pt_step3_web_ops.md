# PT完善 Step3 Web运营端补齐与CRUD闭环

更新时间: 2026-07-01

## 1. 本步目标
- 完成英语学习 Web 运营端可用化，解决“有菜单但不可运营”的问题。
- 打通单词侧与视频侧后台完整增删改查，不再停留在原型表单。
- 完成本步可验收闭环：代码改造 + 测试/lint 验证 + 记录沉淀。

## 2. 本步交付

### 2.1 单词运营台（Web）
- 重构页面：`web/src/plugin/english_learning/view/word.vue`
- 完整能力：
  - 分类管理：新增、编辑、删除、分页、筛选（needVip）
  - 章节管理：新增、编辑、删除、分页、按分类筛选
  - 单词管理：新增、编辑、删除、分页、按章节筛选
- 表单与数据约束：
  - 名称/释义 JSON 格式校验
  - 创建单词时强制至少绑定一个章节
  - 编辑单词时支持可选重绑章节

### 2.2 视频运营台（Web）
- 重构页面：`web/src/plugin/english_learning/view/video.vue`
- 完整能力：
  - 视频分类管理：新增、编辑、删除、分页
  - 剧集管理：新增、编辑、删除、分页、按分类筛选
  - 单集管理：新增、编辑、删除、分页、按剧集筛选
  - 字幕解析：按单集提交字幕 JSON 到解析接口
  - 用户授权管理：按用户 + 资源类型查询，支持授权与撤销
- 资源授权类型覆盖：
  - `english_category`
  - `video_series`
  - `video_episode`

### 2.3 Web API 封装补齐
- 扩展接口文件：`web/src/plugin/english_learning/api/english.js`
- 已覆盖运营端所需全部 API：
  - 单词/章节/分类 CRUD
  - 视频分类/剧集/单集 CRUD
  - 字幕解析
  - 用户授权 grant/revoke/list

### 2.4 后端联动补齐（本步前置已合入）
- 单词 update/delete 接口在 PT-3 中落地并已接入 Web 页面。
- 测试初始化补齐删除路径依赖表结构，避免单测误失败：
  - `server/plugin/english_learning/api/english_word_api_test.go`

### 2.5 质量修正
- 清理 Web 插件入口占位导致的 lint 报错：
  - `web/src/plugin/english_learning/index.js`

## 3. 本步验收结果
- 后端测试通过：
  - `cd server && go test ./plugin/english_learning/...`
- 前端 lint 通过：
  - `cd web && npm run lint -- src/plugin/english_learning/api/english.js src/plugin/english_learning/view/word.vue src/plugin/english_learning/view/video.vue`
- VS Code 诊断：本步相关文件无错误。

## 4. 结论
- PT-3 已完成：英语学习 Web 运营端从原型页升级为可运营后台。
- 当前状态达到“菜单可见 + 页面可用 + 接口可用 + 校验可过”的一步闭环标准。
- 可进入 PT-4（生产化能力增强，如 TTS 真实链路、运营批处理、数据看板等）。
