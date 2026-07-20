# uni-shop-plus 项目指南

## 项目启动与构建指令说明

本项目基于uni-app框架开发，支持多平台编译和部署。以下是可用的启动和构建指令说明：

### 开发环境指令 (dev)

开发环境指令用于本地开发和调试，会启动热更新服务器。

| 指令 | 说明 | 命令 |
| --- | --- | --- |
| `npm run dev:h5` | 启动H5开发环境 | `uni` |
| `npm run dev:h5:ssr` | 启动H5开发环境（SSR模式） | `uni --ssr` |
| `npm run dev:mp-weixin` | 启动微信小程序开发环境 | `uni -p mp-weixin` |
| `npm run dev:mp-alipay` | 启动支付宝小程序开发环境 | `uni -p mp-alipay` |
| `npm run dev:mp-baidu` | 启动百度小程序开发环境 | `uni -p mp-baidu` |
| `npm run dev:mp-toutiao` | 启动字节跳动小程序开发环境 | `uni -p mp-toutiao` |
| `npm run dev:mp-qq` | 启动QQ小程序开发环境 | `uni -p mp-qq` |
| `npm run dev:mp-kuaishou` | 启动快手小程序开发环境 | `uni -p mp-kuaishou` |
| `npm run dev:mp-lark` | 启动飞书小程序开发环境 | `uni -p mp-lark` |
| `npm run dev:mp-jd` | 启动京东小程序开发环境 | `uni -p mp-jd` |
| `npm run dev:mp-xhs` | 启动小红书小程序开发环境 | `uni -p mp-xhs` |
| `npm run dev:mp-harmony` | 启动鸿蒙小程序开发环境 | `uni -p mp-harmony` |
| `npm run dev:quickapp-webview` | 启动快应用开发环境 | `uni -p quickapp-webview` |
| `npm run dev:quickapp-webview-huawei` | 启动华为快应用开发环境 | `uni -p quickapp-webview-huawei` |
| `npm run dev:quickapp-webview-union` | 启动快应用联盟开发环境 | `uni -p quickapp-webview-union` |
| `npm run dev:custom` | 自定义平台开发环境 | `uni -p` |

### 生产环境构建指令 (build)

生产环境构建指令用于生成可部署的生产环境代码。

| 指令 | 说明 | 命令 |
| --- | --- | --- |
| `npm run build:h5` | 构建H5生产环境 | `uni build` |
| `npm run build:h5:ssr` | 构建H5生产环境（SSR模式） | `uni build --ssr` |
| `npm run build:mp-weixin` | 构建微信小程序生产环境 | `uni build -p mp-weixin` |
| `npm run build:mp-alipay` | 构建支付宝小程序生产环境 | `uni build -p mp-alipay` |
| `npm run build:mp-baidu` | 构建百度小程序生产环境 | `uni build -p mp-baidu` |
| `npm run build:mp-toutiao` | 构建字节跳动小程序生产环境 | `uni build -p mp-toutiao` |
| `npm run build:mp-qq` | 构建QQ小程序生产环境 | `uni build -p mp-qq` |
| `npm run build:mp-kuaishou` | 构建快手小程序生产环境 | `uni build -p mp-kuaishou` |
| `npm run build:mp-lark` | 构建飞书小程序生产环境 | `uni build -p mp-lark` |
| `npm run build:mp-jd` | 构建京东小程序生产环境 | `uni build -p mp-jd` |
| `npm run build:mp-xhs` | 构建小红书小程序生产环境 | `uni build -p mp-xhs` |
| `npm run build:mp-harmony` | 构建鸿蒙小程序生产环境 | `uni build -p mp-harmony` |
| `npm run build:quickapp-webview` | 构建快应用生产环境 | `uni build -p quickapp-webview` |
| `npm run build:quickapp-webview-huawei` | 构建华为快应用生产环境 | `uni build -p quickapp-webview-huawei` |
| `npm run build:quickapp-webview-union` | 构建快应用联盟生产环境 | `uni build -p quickapp-webview-union` |
| `npm run build:custom` | 自定义平台生产环境构建 | `uni build -p` |

## 使用说明

1. 开发时，根据目标平台选择对应的开发指令，例如开发微信小程序使用 `npm run dev:mp-weixin`
2. 构建生产版本时，使用对应平台的构建指令，例如构建H5版本使用 `npm run build:h5`
3. 自定义平台开发和构建时，需要在命令后面指定平台标识

视频防盗链说明
Fix 2 — 视频 URL 联动 ExternalLinkDomain
根因：SignLearningVideoURL 原逻辑是 有防盗链 CDN → 签名，没有 → 原样返回。当 hotlink 未配置时，DB 里的相对路径 /file/.../ep1.mp4 或残留的 https://video.mnmovie.icu/... 直接透传给前端，前端无法播放。修复：video_security.go 新增二级回退链：
SignLearningVideoURL(rawURL)
  ├── Hotlink.CdnDomain 有值？
  │     YES → normalizePath → upload.SignURL()    (防盗链签名)
  │     NO  → resolveWithDefaultDomain(rawURL)     (ExternalLinkDomain回退)
  │              ├── rawURL 是完整HTTP URL？→ 提取Path
  │              ├── 确保Path以 / 开头
  │              ├── getCachedDefaultDomain()        (查DB or 缓存)
  │              └── domain + filePath → https://默认云域名/file/.../ep1.mp4
  └── 都没配 → rawURL 原样返回
getCachedDefaultDomain() 带 sync.RWMutex 缓存，只查一次 DB，后续调用零开销
查询逻辑与前端 getExternalUrl 一致：优先 is_default=true → 回退第一个 is_enabled → Domain 优先、空则 BaseURL