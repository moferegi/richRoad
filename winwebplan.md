# WinWeb 电脑网页版开发计划

> 将 uni 移动端学习模块的 17 个页面迁移为 PC 电脑网页版
> 技术栈：Vite + Vue 3 + Pinia + Vue Router
> 数据、逻辑、接口保持不变，仅做 PC 端布局适配
> 设计主题风格与 uni 移动端保持一致（紫色渐变主题）

---

## 一、项目总览

### 1.1 目录结构

```
winWeb/                    # 与 uni、web 同一层
├── public/
│   └── favicon.ico
├── src/
│   ├── api/               # 接口层（直接复用 uni 中的接口逻辑）
│   │   ├── learning.js    # 学习模块接口
│   │   ├── base.js        # 基础接口（登录、验证码等）
│   │   ├── kefu.js        # 客服接口
│   │   ├── collect.js     # 收藏接口
│   │   └── ...
│   ├── assets/            # 静态资源
│   │   ├── images/
│   │   └── styles/
│   │       ├── variables.scss   # 全局样式变量（紫色主题）
│   │       └── global.scss      # 全局样式
│   ├── components/        # 公共组件
│   │   ├── AppHeader.vue        # 顶部导航栏
│   │   ├── AppSidebar.vue       # 左侧侧边栏（学习模块导航）
│   │   ├── PlayerVideo.vue      # 视频播放器组件
│   │   ├── KbdKeyboard.vue      # 跟打键盘组件
│   │   ├── WordCard.vue         # 单词卡片组件
│   │   ├── SubtitleScroll.vue   # 字幕滚动组件
│   │   ├── Pagination.vue       # 分页组件
│   │   ├── EmptyState.vue       # 空状态组件
│   │   ├── ModalPopup.vue       # 弹窗组件
│   │   └── LoadingSpinner.vue   # 加载动画
│   ├── layouts/           # 布局组件
│   │   └── MainLayout.vue       # 主布局（Header + Sidebar + Content）
│   ├── router/            # 路由配置
│   │   └── index.js
│   ├── pinia/             # 状态管理
│   │   ├── index.js
│   │   └── modules/
│   │       ├── user.js
│   │       └── lang.js
│   ├── utils/             # 工具函数
│   │   ├── request.js     # 请求封装（axios）
│   │   ├── i18n.js        # 多语言
│   │   └── storage.js     # 本地存储
│   ├── views/             # 页面组件
│   │   ├── login/
│   │   │   └── index.vue
│   │   ├── register/
│   │   │   └── index.vue
│   │   ├── home/
│   │   │   └── index.vue
│   │   ├── learning/
│   │   │   ├── diary/
│   │   │   │   └── index.vue
│   │   │   ├── diary-detail/
│   │   │   │   └── index.vue
│   │   │   ├── typing/
│   │   │   │   └── index.vue
│   │   │   ├── profile/
│   │   │   │   └── index.vue
│   │   │   ├── player/
│   │   │   │   └── index.vue
│   │   │   ├── video-detail/
│   │   │   │   └── index.vue
│   │   │   ├── collections/
│   │   │   │   └── index.vue
│   │   │   ├── watch-history/
│   │   │   │   └── index.vue
│   │   │   ├── checkin-record/
│   │   │   │   └── index.vue
│   │   │   ├── error-log/
│   │   │   │   └── index.vue
│   │   │   ├── point-history/
│   │   │   │   └── index.vue
│   │   │   ├── free-time-history/
│   │   │   │   └── index.vue
│   │   │   └── tag-filter/
│   │   │       └── index.vue
│   │   └── kefu/
│   │       └── index.vue
│   ├── App.vue
│   └── main.js
├── .env.development
├── .env.production
├── index.html
├── package.json
├── vite.config.js
└── jsconfig.json
```

### 1.2 技术选型

| 技术 | 选型 | 说明 |
|------|------|------|
| 构建工具 | Vite 5.x | 与 uni 端保持一致 |
| 框架 | Vue 3 + Composition API | 与 uni 端一致 |
| 路由 | Vue Router 4 | PC 端 SPA 路由 |
| 状态管理 | Pinia | 与 uni 端一致 |
| HTTP | axios | 替代 uni.request |
| 样式 | SCSS | 与 uni 端一致 |
| 视频播放 | hls.js | 与 uni H5 端一致 |
| 多语言 | vue-i18n | 与 uni 端一致 |
| UI 组件库 | 自定义组件 | 不引入 Element Plus 等，保持与 uni 设计风格一致 |

### 1.3 设计规范

**主色调**: `#6D5BFF → #9B8FFF` 紫色渐变
**背景色**: `#F5F3FF` 浅紫灰
**卡片**: 白色圆角 12px + 柔和紫色阴影
**按钮**: 胶囊形 + 紫色渐变
**文字色**: `#1A1B3A`（主）/ `#6B6F8D`（次）/ `#A9AECB`（弱）
**布局**: 1200px 最大宽度居中 + 左侧 240px 侧边栏 + 右侧内容区

---

## 二、开发步骤（按顺序）

### 阶段一：项目脚手架搭建 ✅ 已完成

- [x] **2.1 初始化 Vite + Vue 3 项目**
  - 创建 winWeb 目录及基础文件
  - 配置 package.json（vue、pinia、vue-router、axios、hls.js、sass、vue-i18n）
  - 配置 vite.config.js（路径别名 @、代理）
  - 配置 jsconfig.json
  - 配置 .env.development / .env.production

- [x] **2.2 基础架构搭建**
  - 创建 src/main.js 入口文件
  - 创建 src/App.vue 根组件
  - 配置 src/router/index.js 路由框架（17个路由 + 路由守卫）
  - 配置 src/pinia/index.js 及 user/lang 模块
  - 创建 src/utils/request.js 请求封装（axios + 401拦截）
  - 创建 src/utils/storage.js 存储封装
  - 创建 src/utils/i18n.js 多语言配置（zh/en）

- [x] **2.3 全局样式与主题**
  - 创建 src/assets/styles/variables.scss（紫色主题设计 token）
  - 创建 src/assets/styles/global.scss（全局重置、基础样式、通用类）
  - 定义设计 token：紫色主题、圆角、阴影、渐变

- [x] **2.4 布局组件**
  - 创建 layouts/MainLayout.vue（Header + Sidebar + Content）
  - 创建 components/AppHeader.vue（顶部导航：Logo、搜索、用户头像下拉）
  - 创建 components/AppSidebar.vue（左侧学习模块导航：首页、日记、单词跟打、我的+子菜单）

### 阶段二：接口层迁移 ✅ 已完成

- [x] **2.5 API 接口迁移**
  - **learning.js**（学习模块 37 个接口）
    - 打卡/积分：getCheckinStats、getLearningCheckinRecordList、doCheckin、exchangeTime、getLearningPointRecordList
    - 资产/时长：getLearningFreeTimeRecordList、getAsset、heartbeat
    - 内容分类：getCategoryList、getChapterList、getVideoCategoryList
    - 视频系列/集：getVideoSeriesList、findVideoSeries、findVideoEpisode、getVideoEpisodeList、getVideoEpisodeListByTag
    - 视频标签：getVideoTagList
    - 单词：getWordList、findWord、reportWordError、getWordErrorLogList、deleteWordErrorLog
    - 句子/字幕：getSentenceList
    - 用户数据/进度：saveWordProgress、getWordProgress、getWatchProgress、getSeriesWatchProgressList
    - 观看历史：getWatchHistoryList
    - 收藏：collect、uncollect、getCollectionList、getCollectionDetailList
    - 日记：getDiaryCategoryList、getDiaryTagPublic、getDiaryList、getDiaryListByTag、findDiary、getDiarySentenceList
  - **base.js**（9 个接口：getCaptcha、register、login、getUserInfo、phoneLogin、phoneRegister、setClientUserInfo、changePassword、getMyInviteInfo）
  - **kefu.js**（3 个接口：getKefuList、getCsConfig、getSysConfigByKey）
  - **homePage.js**（1 个接口：getBannerList）
  - **sysConfig.js**（4 个接口：getSysConfigByGroup、getSysConfigByKey、getLoginConfig、getPointsExchangeRate）
  - **phoneAreaCode.js**（1 个接口：getEnabledPhoneAreaCodes）
  - uni.request → axios 替换，请求/响应拦截器保持一致逻辑（401 跳转登录）

> **注意**：uni 端另有 collect.js（商城收藏）、signIn.js（商城签到）、language.js、visitor.js、loginApi.js（小程序）等接口文件，均为商城/小程序专用，学习模块 PC 端无需迁移。

### 阶段三：登录注册 ✅ 已完成

- [x] **2.6 登录页** (`/login`)
  - PC 双栏布局：左侧紫色渐变品牌展示区（Logo+名称+Slogan+特性列表），右侧表单区
  - 双模式切换胶囊：账号登录 / 手机号登录
  - 表单字段：用户名/手机号、密码、验证码（可点击刷新）
  - 区号选择弹窗（列表形式）
  - 语言切换按钮（中文/英文）
  - 「去注册」入口
  - 登录成功跳转到首页（支持 redirect 参数）

- [x] **2.7 注册页** (`/register`)
  - 与登录页相同的双栏布局
  - 双模式：用户名注册 / 手机号注册
  - 表单项：账号/手机、密码、确认密码、验证码、邀请码（选填）
  - 邀请码从 URL 参数自动带入
  - 「去登录」入口
  - 注册成功跳转登录页

### 阶段四：首页与视频页面 ✅ 已完成

- [x] **2.8 首页** (`/home`)
  - 紫色渐变 Hero 区（装饰圆形 + 标题 + Slogan）
  - Banner 轮播区（自动轮播 + 鼠标悬停暂停 + 圆点指示器）
  - 视频分类胶囊选择器（横向排列，"全部" + 各分类）
  - 视频卡片网格（4列自适应，PC 端 4~5 列，hover 上浮效果）
  - 更多标签入口（跳转 tag-filter）
  - 分页加载 + 加载更多按钮
  - 空状态

- [x] **2.9 视频系列详情页** (`/learning/video-detail/:seriesId`)
  - 顶部 Hero 区（背景模糊封面 + 紫色遮罩 + 系列封面 + 标题 + 统计信息）
  - 继续观看按钮（带播放图标）
  - 集数列表（3列网格，进度条 + 百分比显示）
  - 返回按钮
  - 观看进度同步显示

- [x] **2.10 视频播放器页** (`/learning/player/:episodeId`)
  - PC 左右分栏布局：左视频区 + 右字幕区
  - 自定义控制栏（进度条、播放/暂停、倍速、字幕切换、复读、循环）
  - 字幕滚动高亮（当前句左侧紫条高亮 + 放大 + 自动滚动居中）
  - 试看锁定机制（进度条标记 + 试看结束遮罩）
  - 重点词汇查词弹窗（预留接口）
  - 心跳防盗机制（每 30 秒上报）
  - hls.js 支持 m3u8 播放
  - 收藏/取消收藏
  - 支持 sentenceId 参数定位句子

### 阶段五：单词跟打模块 ✅ 已完成

- [x] **2.11 单词跟打页** (`/learning/typing`)
  - PC 三栏布局：左侧词表面板 + 中间主卡片 + 右侧例句
  - 顶部栏：分类/章节选择器 + 错误计数徽章 + 搜索框 + 设置按钮
  - 主卡片：大字显示单词（64px）+ 音标 + 发音按钮 + 释义 + 收藏 + 进度条
  - 跟打模式：物理键盘输入，字符颜色状态（正确紫色，错误红色删除线，当前光标闪烁）
  - 错误上报（位置 + 期望字符 + 输入字符）
  - 上一个/下一个单词 + 自动跳章节
  - TTS 发音（audio URL 优先，Web Speech API 兜底）
  - 设置弹窗：显示单词/音标/释义、US/UK 口音切换
  - 收藏/取消收藏
  - 学习进度保存与恢复
  - 例句区域（前 3 句）

### 阶段六：日记模块 ✅ 已完成

- [x] **2.12 日记列表页** (`/learning/diary`)
  - 顶部分类筛选 + 标签筛选栏（chip 展示已选标签 + 重置按钮）
  - 日记卡片网格布局（PC 端 3 列）
  - 卡片：左图片预览 + 右标题/文本预览/查看按钮
  - 分类/标签选择弹窗（分类单选，标签多选）
  - 分页加载 + 加载更多
  - 空状态
  - 骨架屏加载动画

- [x] **2.13 日记详情页** (`/learning/diary-detail/:id`)
  - 顶部信息栏：返回 + 标题 + 时长 + US/UK 口音切换 + 收藏按钮
  - PC 左右分栏：左侧字幕滚动区 + 右侧图片预览区
  - 底部控制栏：播放/暂停 + 双语切换 + 倍速
  - 字幕点击跳转播放位置
  - TTS 播放（audio URL 优先，Web Speech API 兜底）
  - 试看锁定机制（trialPercent）
  - 收藏/取消收藏（targetType=4）

### 阶段七：个人中心模块 ✅ 已完成

- [x] **2.14 个人中心页** (`/learning/profile`)
  - 顶部用户信息卡（头像 + 昵称 + 积分）
  - 资产概览卡（总积分 + 免费时长）
  - 积分兑换时长功能（弹窗输入积分，调用 exchangeTime 接口）
  - 功能菜单网格（观看历史、打卡记录、我的收藏、错误本、积分记录、时长记录、联系客服、切换语言）
  - 退出登录按钮

- [x] **2.15 我的收藏页** (`/learning/collections`)
  - Tab 切换：全部 / 单词 / 句子 / 视频 / 日记
  - 收藏列表卡片（PC 端列表式布局，左缩略图 + 右信息）
  - 分页加载
  - 移除收藏 / 查看/播放 操作

- [x] **2.16 观看历史页** (`/learning/watch-history`)
  - 历史记录列表（左封面 + 右信息：系列名、集名、进度、时间）
  - 继续观看按钮
  - 分页加载

- [x] **2.17 打卡记录页** (`/learning/checkin-record`)
  - 顶部统计卡（连续打卡、累计打卡、今日状态）
  - 打卡按钮
  - 打卡记录列表（日期 + 积分奖励）
  - 分页加载

- [x] **2.18 错误本页** (`/learning/error-log`)
  - 错误总数摘要卡
  - 错误记录列表（单词 + 释义 + 错误次数 + 错误详情）
  - 「去练习」按钮（跳转到跟打页）
  - 删除操作
  - 分页加载

- [x] **2.19 积分记录页** (`/learning/point-history`)
  - 积分概览
  - 记录列表（时间 + 类型 + 积分变动）
  - 分页加载

- [x] **2.20 免费时长记录页** (`/learning/free-time-history`)
  - 时长概览
  - 记录列表（时间 + 类型 + 时长变动）
  - 分页加载

### 阶段八：其他页面 ✅ 已完成

- [x] **2.21 客服页** (`/kefu?from=learning_profile`)
  - 客服列表卡片
  - 外部客服：头像 + 名称 + 状态 + 联系按钮
  - 平台客服入口（二维码弹窗）
  - 支付上下文提示卡（从订单跳入时）
  - 空状态
  - 复制功能 + 复制成功提示

- [x] **2.22 标签筛选页** (`/learning/tag-filter`)
  - 所有标签分类展示（按 group 分组）
  - 标签选择交互（多选 chip 样式）
  - 已选标签顶部展示
  - 确认跳转视频列表（带 tags 查询参数）
  - 重置按钮

---

## 三、页面路由表

| 路径 | 页面 | 对应 uni 页面 |
|------|------|--------------|
| `/login` | 登录页 | pages/user/login |
| `/register` | 注册页 | pages/user/register |
| `/home` | 首页 | pages/learning/home |
| `/learning/diary` | 日记列表 | pages/learning/diary |
| `/learning/diary-detail/:id` | 日记详情 | pages/learning/diary-detail |
| `/learning/typing` | 单词跟打 | pages/learning/typing |
| `/learning/profile` | 个人中心 | pages/learning/profile |
| `/learning/player/:episodeId` | 视频播放器 | pages/learning/player |
| `/learning/video-detail/:seriesId` | 视频系列详情 | pages/learning/video-detail |
| `/learning/collections` | 我的收藏 | pages/learning/collections |
| `/learning/watch-history` | 观看历史 | pages/learning/watch-history |
| `/learning/checkin-record` | 打卡记录 | pages/learning/checkin-record |
| `/learning/error-log` | 错误本 | pages/learning/error-log |
| `/learning/point-history` | 积分记录 | pages/learning/point-history |
| `/learning/free-time-history` | 免费时长记录 | pages/learning/free-time-history |
| `/learning/tag-filter` | 标签筛选 | pages/learning/tag-filter |
| `/kefu` | 客服 | pages/kefu/index |

---

## 四、核心迁移策略

### 4.1 接口迁移
- 接口函数签名保持完全一致
- 将 `uni.request` 替换为 `axios`
- 错误处理逻辑保持一致（401 跳转登录、错误 toast 提示）
- Token 存储从 `uni.setStorageSync` 改为 `localStorage`

### 4.2 状态管理迁移
- Pinia store 结构保持一致
- user module：token、userInfo、login/logout 方法
- lang module：多语言切换逻辑

### 4.3 组件迁移
- uni-app 内置组件 → 原生 HTML + 自定义组件
  - `<view>` → `<div>`
  - `<text>` → `<span>`
  - `<image>` → `<img>`
  - `<swiper>` → 自定义轮播 / 左右分栏
  - `<video>` → `<video>` + hls.js
  - `<scroll-view>` → 原生滚动 div
  - `<uni-popup>` → 自定义 Modal 组件
- UI 交互：底部弹窗 → 侧边抽屉 / 居中弹窗

### 4.4 PC 端布局适配原则
1. **侧边导航**：移动端底部 TabBar → PC 端左侧常驻侧边栏（240px）
2. **内容宽度**：移动端全屏 → PC 端 1200px 最大宽度居中
3. **弹窗交互**：移动端底部上滑 → PC 端居中弹窗 / 右侧抽屉
4. **列表布局**：移动端单列 → PC 端多列网格 / 列表式（左图右文）
5. **播放器**：移动端竖屏上下布局 → PC 端横屏左右分栏（左视频右字幕）
6. **键盘输入**：移动端屏幕键盘 → PC 端物理键盘输入为主，屏幕键盘为辅

### 4.5 保持不变的部分
- 所有 API 接口地址、参数、返回值
- 业务逻辑（鉴权、试看、心跳、打卡等）
- 数据结构
- 设计风格（紫色主题、卡片、按钮样式）
- 多语言体系

---

## 五、开发环境配置

### 5.1 环境变量

`.env.development`:
```
VITE_API_BASE_URL=/api
VITE_APP_TITLE=RichRoad Learning
```
> 开发环境使用 `/api` 前缀 + Vite 代理 rewrite，与 uni H5 端一致。

`.env.production`:
```
VITE_API_BASE_URL=https://www.xxxxxxxx.com/api/wechat
VITE_APP_TITLE=RichRoad Learning
```

### 5.2 Vite 代理

```js
server: {
  port: 5222,
  proxy: {
    '/api': {
      target: 'http://localhost:8888',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```

> 端口使用 5222，避免与 uni 端（5173）冲突。代理需 rewrite 去掉 `/api` 前缀，与 uni 端保持一致。

---

---

## 七、已知未完成项（待后续迭代）

以下功能在首版中简化或暂未实现，记录在此供后续迭代：

| 序号 | 模块 | 说明 | 优先级 |
|------|------|------|--------|
| 1 | 公共组件 | `Pagination.vue`、`EmptyState.vue`、`ModalPopup.vue`、`LoadingSpinner.vue` 未独立抽离，各页面内联实现 | 低 |
| 2 | 视频播放器 | 重点词汇查词弹窗预留接口，未完整实现 | 中 |
| 3 | 视频播放器 | 复读 / 循环 功能部分简化 | 中 |
| 4 | 单词跟打 | 屏幕键盘（KbdKeyboard）组件未独立，物理键盘输入为主 | 低 |
| 5 | 单词跟打 | WordCard / SubtitleScroll 组件未独立抽离 | 低 |
| 6 | 日记详情 | 图片预览区交互简化，未做放大查看 | 低 |
| 7 | 多语言 | 仅 zh/en 两种语言，未同步 uni 端全部 11 种语言 | 中 |
| 8 | 骨架屏 | 日记列表骨架屏未实现，使用 loading 替代 | 低 |
| 9 | 客服页 | 支付上下文提示卡逻辑简化 | 低 |
| 10 | 接口层 | `kefu.js` 中 `getSysConfigByKey` 与 `sysConfig.js` 重复定义，建议统一到 sysConfig.js | 低 |
| 11 | 接口层 | `getDiaryList`（普通日记列表）导出但未使用，日记页统一使用 `getDiaryListByTag`；需后端确认 tagIds 为空时行为是否一致 | 中 |
| 12 | 接口层 | `getVideoEpisodeListByTag` 导出但未被 tag-filter 页使用，标签筛选页跳转回首页带 query 参数实现 | 低 |
| 13 | 接口层 | `setClientUserInfo`、`changePassword`、`getMyInviteInfo` 导出但页面未调用（个人中心暂未做修改信息/密码/邀请功能） | 低 |

---

## 八、部署上线

### 8.1 构建命令

```bash
# 开发环境
npm run dev

# 生产构建
npm run build

# 预览构建结果
npm run preview
```

### 8.2 构建产物

- 输出目录：`dist/`
- 静态资源：`dist/assets/`
- 入口文件：`dist/index.html`

### 8.3 Nginx 部署配置

```nginx
server {
    listen 80;
    server_name www.example.com;
    root /path/to/winWeb/dist;
    index index.html;

    # SPA 路由 fallback
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # API 代理（如需同域部署）
    location /api/ {
        proxy_pass http://backend:8888/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 8.4 环境变量区分

- 开发环境 `.env.development`：本地接口代理
- 生产环境 `.env.production`：线上正式接口地址

---

## 九、验收标准

- [x] 17 个页面全部完成，样式与设计规范一致
- [x] 所有接口调用正常，数据展示正确
- [x] 核心功能可用（登录注册、视频播放、单词跟打、日记播放、打卡、收藏等）
- [x] 响应式布局（1280 / 1440 / 1920 分辨率自适应）
- [x] 多语言切换正常（zh / en）
- [x] 路由跳转关系与移动端一致
- [ ] 线上部署并通过域名访问
- [ ] 主流浏览器兼容性测试（Chrome / Edge / Firefox / Safari）

---

## 十、项目总结

### 10.1 完成概览

| 阶段 | 内容 | 状态 |
|------|------|------|
| 阶段一 | 项目脚手架搭建 | ✅ |
| 阶段二 | 接口层迁移 | ✅ |
| 阶段三 | 登录注册 | ✅ |
| 阶段四 | 首页与视频页面 | ✅ |
| 阶段五 | 单词跟打模块 | ✅ |
| 阶段六 | 日记模块 | ✅ |
| 阶段七 | 个人中心模块（7 个页面） | ✅ |
| 阶段八 | 其他页面（客服、标签筛选） | ✅ |

**总计：17 个页面全部完成开发。**

### 10.2 后续优化方向

1. **组件抽离**：将各页面中重复的分页、空状态、弹窗等抽离为公共组件
2. **多语言扩展**：同步 uni 端全部 11 种语言
3. **播放器增强**：完善查词弹窗、复读、AB 循环等高级功能
4. **性能优化**：路由懒加载、图片懒加载、虚拟列表
5. **骨架屏**：补充列表页骨架屏加载动画
6. **单元测试**：核心业务逻辑添加单元测试
7. **错误监控**：接入前端错误监控系统
