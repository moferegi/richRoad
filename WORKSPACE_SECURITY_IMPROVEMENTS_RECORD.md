# 工作区安全改进完整记录（面向小白）

更新时间：2026-05-20

适用范围：本工作区当前安全改造主线（后端、前端配套、初始化种子、运维文档）。

阅读方式：
- 每一条都按固定格式写明：改了什么、具体有什么用、能防止什么、怎么操作。
- 你可以直接跳到「第 12 章 一键检查清单」先做验收，再回看每条原理。

---

## 1. 核心运行时配置统一（总开关能力）

### 1.1 统一读取安全配置优先级
- 改了什么：新增统一配置读取工具，优先级变为「数据库系统参数 > 环境变量 > 默认值」。
- 文件：server/utils/runtime_settings.go
- 有什么用：安全策略可以不改代码直接调（后台改参数即可生效）。
- 能防止什么：环境变量、代码默认值、后台参数三套配置冲突导致的误开放。
- 小白怎么操作：
  1. 在后台系统参数里改 security 组配置。
  2. 没配的项再用环境变量兜底。
  3. 都没配时才走默认值。

### 1.2 新增完整 security 参数默认值
- 改了什么：补齐并初始化大量 security 配置键（公开限流、WS、初始化接口、自动封禁、错误上报限流、访客心跳、试衣限流等）。
- 文件：server/initialize/other.go
- 有什么用：新环境启动后不再“裸奔”，每个安全模块都有默认保护。
- 能防止什么：上线忘记配参数导致开关失效。
- 小白怎么操作：后台「系统参数 -> security 分组」按业务量微调阈值。

---

## 2. 权限与越权防护

### 2.1 系统敏感接口 API 层白名单
- 改了什么：在 API 层再次校验角色。
- 文件：server/api/v1/system/sys_system.go
- 有什么用：即使 Casbin 配错，敏感接口仍会被代码层拒绝（双保险）。
- 能防止什么：低权限角色误调用系统配置/系统重载/服务器信息接口。
- 小白怎么操作：
  1. 仅超级管理员执行系统重载与服务器信息查询。
  2. 配置读取/修改仅允许指定管理角色。

### 2.2 API 同步/刷新权限收口
- 改了什么：刷新 Casbin 缓存接口增加角色限制（仅管理角色）。
- 文件：server/api/v1/system/sys_api.go, server/router/system/sys_api.go
- 有什么用：避免普通角色随意触发权限缓存刷新。
- 能防止什么：权限抖动、绕审计式操作。
- 小白怎么操作：只让 888/8881 角色执行刷新操作。

### 2.3 客户端公开配置键白名单
- 改了什么：客户端按 key 读取系统参数时，增加公开 allowlist。
- 文件：server/api/v1/client/sys_config.go
- 有什么用：前端只能取允许公开的参数。
- 能防止什么：把后台敏感配置（如 token、内部策略）暴露给客户端。
- 小白怎么操作：新增前端配置项时，先判断是否可公开，再决定是否加入白名单。

### 2.4 订单管理敏感操作角色校验
- 改了什么：退款、删除、确认收款、部分状态变更增加后台角色判断。
- 文件：server/api/v1/shop/order.go
- 有什么用：客户端用户只能做自己的订单流程，不能做后台管理动作。
- 能防止什么：越权退款、越权删单、越权改状态。
- 小白怎么操作：管理动作统一走后台账号（888/8881）。

### 2.5 浏览历史清空接口从公开组收口到私有组
- 改了什么：清空历史接口移入鉴权路由组。
- 文件：server/router/shop/good.go
- 有什么用：需要登录+权限后才能调用。
- 能防止什么：未登录用户清空他人历史或恶意刷接口。
- 小白怎么操作：确认前端调用时必须携带登录态。

### 2.6 个人资料接口权限自修复
- 改了什么：启动时自动补齐 /clientUser/setClientUserInfo 的 API 与 Casbin 权限。
- 文件：server/initialize/security_init.go, server/source/system/api.go, server/source/system/casbin.go, server/router/client/user.go
- 有什么用：历史环境缺权限时自动修复，不用手工补表。
- 能防止什么：客户端用户无法更新个人资料，或权限错配。
- 小白怎么操作：重启后检查角色 8080 是否具备该接口权限。

### 2.7 系统重载权限自修复
- 改了什么：启动时自动补齐 /system/reloadSystem API 与管理员权限，并支持清理历史过授权。
- 文件：server/initialize/security_init.go, server/source/system/api.go, server/source/system/casbin.go
- 有什么用：系统重载权限可持续保持正确。
- 能防止什么：非管理员拥有重载权限。
- 小白怎么操作：
  1. 维护窗口开启过授权清理开关。
  2. 启动后再关闭开关。

### 2.8 历史 Casbin 过授权“一次性清理”机制
- 改了什么：新增可开关清理器，按接口+允许角色删除多余授权。
- 文件：server/initialize/casbin_cleanup.go
- 有什么用：可安全、可灰度地治理历史脏权限。
- 能防止什么：老环境“权限漂移”长期累积。
- 小白怎么操作：
  1. 先备份 casbin_rule。
  2. 打开 security_cleanup_legacy_casbin_overgrant。
  3. 启动一次后关闭。

### 2.9 新模块权限种子范围收敛 + 历史授权清理
- 改了什么：新模块初始化时，仅对既定角色（888/8881/8080/9528）写入种子权限，并执行受管规则清理；sysConfig 读接口中高风险项收口为管理员权限。
- 文件：server/initialize/new_modules_init.go
- 有什么用：避免“全角色批量赋权”造成权限外溢。
- 能防止什么：历史角色遗留获得新模块管理权限。
- 小白怎么操作：初始化后抽查 casbin_rule，确认非目标角色没有对应新模块管理接口授权。

### 2.10 商城权限种子范围收敛 + 管理端接口回收
- 改了什么：商城模块与新模块同样改为仅对既定角色播种权限，并对订单管理、sysConfig、clientUser、外链域名等管理接口执行非管理员授权清理。
- 文件：server/initialize/shop_init.go
- 有什么用：防止业务扩展后“角色越多、误授权越多”。
- 能防止什么：普通运营角色误拿到退款、删单、配置写入等高风险权限。
- 小白怎么操作：上线后验证只有 888/8881 具备上述管理接口权限。

### 2.11 AutoCode 路由去公开化
- 改了什么：自动代码模块原本挂在公开路由组的 llmAuto/initMenu/initAPI/initDictionary 改为私有鉴权路由组。
- 文件：server/router/system/sys_auto_code.go
- 有什么用：自动代码与插件初始化能力不再裸露到公开入口。
- 能防止什么：未登录或低权限绕过后台流程触发高危操作。
- 小白怎么操作：确认这些接口调用都必须携带后台登录态并通过权限校验。

---

## 3. 初始化接口与审计盲区治理

### 3.1 初始化接口加开关与内网限制
- 改了什么：/init/checkdb 与 /init/initdb 增加统一访问校验：开关控制 + 可选仅内网/回环。
- 文件：server/api/v1/system/sys_initdb.go
- 有什么用：生产环境默认可关闭初始化接口。
- 能防止什么：公网误暴露初始化入口导致的高危操作。
- 小白怎么操作：
  1. 生产建议 security_init_api_enabled=false。
  2. 若临时开启，也保持 private_network_only=true。

### 3.2 路由注册层按开关决定是否挂 init 路由
- 改了什么：路由启动时根据 security_init_api_enabled 决定是否注册初始化路由。
- 文件：server/initialize/router.go
- 有什么用：接口从源头不注册，攻击面更小。
- 能防止什么：即使 API 层判断出错，路由仍暴露。
- 小白怎么操作：生产常态不挂 init 路由。

### 3.3 清理敏感 ignore API（审计盲区）
- 改了什么：启动时自动删除敏感接口的 ignore 规则。
- 文件：server/initialize/security_init.go
- 有什么用：敏感接口重新进入鉴权审计链路。
- 能防止什么：敏感接口被“忽略”后绕过审计。
- 小白怎么操作：启动后检查 sys_ignore_apis，不应包含 initdb/checkdb/reloadSystem。

### 3.4 默认 ignore 种子移除敏感接口
- 改了什么：初始化忽略接口种子仅保留真正公开且低风险接口。
- 文件：server/source/system/api_ignore.go
- 有什么用：新环境默认不会把敏感接口放进忽略名单。
- 能防止什么：新部署即带安全隐患。
- 小白怎么操作：如果你手工加了忽略，请定期复核。

---

## 4. 防刷、限流、自动封禁

### 4.1 公共路由统一 IP 限流
- 改了什么：PublicGroup 增加统一限流中间件（Redis 计数）。
- 文件：server/middleware/public_rate_limit.go, server/initialize/router.go
- 有什么用：公开接口全局防洪峰与恶意刷请求。
- 能防止什么：脚本批量打公开接口拖垮服务。
- 小白怎么操作：调 security_public_rate_limit_* 三个参数即可。

### 4.2 登录 IP 频率限制
- 改了什么：新增登录按 IP 频率限制（Redis，失败回退本地内存计数）。
- 文件：server/service/client/security.go
- 有什么用：同一 IP 爆破时快速限速。
- 能防止什么：高频密码爆破。
- 小白怎么操作：调 security_login_ip_rate_limit_per_minute 与 window_seconds。

### 4.3 登录失败锁定仍保留并与新限流联动
- 改了什么：登录失败计数机制保留，并在用户名/手机号登录都接入。
- 文件：server/service/client/security.go, server/api/v1/client/user.go
- 有什么用：既控“频率”，又控“连续失败次数”。
- 能防止什么：低频慢速爆破绕过单一限流。
- 小白怎么操作：调 login_fail_max 与 login_fail_wait_seconds。

### 4.4 图形验证码接口频率限制
- 改了什么：获取验证码前先做 IP 频率限制。
- 文件：server/api/v1/system/sys_captcha.go, server/service/client/security.go
- 有什么用：减轻验证码接口被刷造成的压力。
- 能防止什么：验证码接口被当放大器攻击。
- 小白怎么操作：调 captcha_rate_limit 与 captcha_rate_limit_window_seconds。

### 4.5 访客心跳双重防刷（限流 + 去重）
- 改了什么：访客心跳新增按 IP 限流、按 visitorID 去重窗口。
- 文件：server/api/v1/client/visitor.go, server/service/client/visitor.go
- 有什么用：减少脏数据与无意义写库。
- 能防止什么：伪造心跳刷统计、数据库写放大。
- 小白怎么操作：调 security_visitor_heartbeat_rate_limit_per_minute 与 dedupe_seconds。

### 4.6 试衣创建限流 + 并发上限
- 改了什么：创建试衣任务前校验“每分钟次数”和“processing 并发数”。
- 文件：server/service/client/tryon_task.go
- 有什么用：控制高成本模型调用与队列堆积。
- 能防止什么：单账号恶意消耗算力。
- 小白怎么操作：调 security_tryon_create_rate_limit_per_minute 与 security_tryon_create_concurrency_limit。

### 4.7 错误上报接口专项限流
- 改了什么：/sysError/createSysError 增加限流（Redis+本地兜底），并返回等待秒数。
- 文件：server/api/v1/system/sys_error.go
- 有什么用：防止错误上报接口被刷爆。
- 能防止什么：伪造报错淹没真实日志。
- 小白怎么操作：调 security_sys_error_create_rate_limit_*。

### 4.8 错误上报请求体收紧
- 改了什么：新增 CreateSysErrorRequest，限制字段与长度。
- 文件：server/model/system/request/sys_error.go
- 有什么用：避免客户端随意写状态/大字段污染。
- 能防止什么：异常大 payload、字段滥用。
- 小白怎么操作：客户端只传 form/info/level。

### 4.9 攻击记录按类型拆分 + 自动封禁阈值分离
- 改了什么：攻击计数键从“按 IP 单桶”改为“attackType + IP”；sys_error_rate_limit 有独立阈值。
- 文件：server/service/system/sys_banned_ip.go
- 有什么用：不同攻击场景可单独调阈值，不互相污染。
- 能防止什么：错误上报刷子误触发登录类封禁，或反向绕过。
- 小白怎么操作：调 security_attack_auto_ban_threshold_default 与 security_attack_auto_ban_threshold_sys_error_rate_limit。

### 4.10 攻击统计新增 sysErrorRate 指标
- 改了什么：攻击统计返回项增加 sysErrorRate。
- 文件：server/service/system/sys_banned_ip.go
- 有什么用：运维可以快速识别“错误上报被刷”场景。
- 能防止什么：误判攻击来源，处置不精准。
- 小白怎么操作：后台看 /sysBannedIP/getAttackStats 的新字段。

### 4.11 商城订单风控加固（建单限频 + 待支付上限 + 改价前状态锁）
- 改了什么：
  - 下单前新增“每分钟建单限频”（用户维度）；
  - 新增“待支付/待确认订单上限”（用户维度）；
  - 改券与改积分改为事务内加锁读取，并强制仅允许“待支付且未超时”订单修改。
- 文件：server/service/shop/order.go, server/initialize/other.go
- 有什么用：把“刷单占库存”和“已支付后再改金额”的风险压住，且阈值可后台调参。
- 能防止什么：
  - 高频建单占用库存；
  - 待支付单无限累积；
  - 已支付或已超时订单继续改券/改积分导致金额不一致。
- 小白怎么操作：
  1. 在后台系统参数 security 分组设置：security_order_create_rate_limit_per_minute。
  2. 在后台系统参数 security 分组设置：security_order_pending_limit_per_user。
  3. 若大促期间误拦截，可临时上调阈值，不需要改代码。

### 4.12 试衣币充值风控加固（建单限频 + 待支付上限 + 超时校验）
- 改了什么：
  - 充值建单新增“每分钟建单限频”（用户维度）；
  - 新增“待支付/待确认充值单上限”（用户维度）；
  - 提交付款确认与管理端确认时增加 close_time 超时校验；
  - 补充新增错误码的多语言文案。
- 文件：server/service/client/tryon_recharge_order.go, server/utils/i18n/i18n.go, server/initialize/other.go
- 有什么用：防止充值链路被刷单与历史超时单继续推进状态。
- 能防止什么：
  - 高频创建充值单压垮审核链路；
  - 待支付充值单堆积；
  - 超时订单继续提交/确认付款。
- 小白怎么操作：
  1. 在后台系统参数 security 分组设置：security_tryon_recharge_create_rate_limit_per_minute。
  2. 在后台系统参数 security 分组设置：security_tryon_recharge_pending_limit_per_user。
  3. 如果提示“订单已超时”，引导用户重新创建订单即可。

### 4.13 待后台确认订单自动关闭（可配置）
- 改了什么：
  - 商品订单从“待支付(0)”提交到“待后台确认(8)”时，刷新 close_time；
  - 试衣币充值订单提交付款确认进入“待后台确认(8)”时，刷新 close_time；
  - 定时清理任务由仅处理状态0扩展为处理状态0和8。
- 文件：server/service/shop/order.go, server/service/client/tryon_recharge_order.go, server/task/closeOrder.go, server/initialize/other.go
- 有什么用：后台审核单不会无限堆积，可按业务配置自动过期释放。
- 能防止什么：
  - 待确认订单长期滞留造成库存/订单占用；
  - 老旧待确认订单被迟到处理带来的资金与履约风险。
- 小白怎么操作：
  1. 在后台系统参数 order 分组设置：order_pending_confirm_close_minutes。
  2. 在后台系统参数 tryon 分组设置：tryon_recharge_pending_confirm_close_minutes。
  3. 审核链路较长时可临时调大分钟数，不需要改代码。

### 4.14 订单状态机收口（合法流转 + 幂等 + 行锁）
- 改了什么：
  - UpdateOrderStatus 增加状态流转校验（仅允许合法前置状态进入目标状态）；
  - 同状态重复写入直接幂等返回，避免重复副作用；
  - 订单行加 `FOR UPDATE` 锁，降低并发下重复加销量/重复恢复库存风险；
  - 批量改状态从“直接改字段”改为逐单走安全状态流转函数。
- 文件：server/service/shop/order.go
- 有什么用：把订单状态变更统一收敛到一个安全入口，副作用（库存、积分、销量）不会被批量接口绕过。
- 能防止什么：
  - 管理端误操作把已支付单直接取消导致库存异常；
  - 重复回调/重复提交导致重复加销量或重复恢复库存；
  - 批量改状态跳过业务规则造成数据不一致。
- 小白怎么操作：
  1. 批量改状态前确认目标状态与当前状态是否匹配（如取消仅对待支付/待确认）。
  2. 如果返回“当前订单状态不支持…”，先在后台查看该订单当前状态再操作。

### 4.15 订单错误多语言统一（服务层 key 化 + API 层翻译）
- 改了什么：
  - 订单服务层将关键业务错误统一改为 i18n key（如 `orderCreateTooFrequent`、`orderExpiredRecreate`、`orderRefundNotPaid`）；
  - 订单 API 层统一通过 `failWithErr/failWithKey` 返回翻译后的消息，移除主要中文硬编码错误；
  - i18n 词典补齐订单状态流转、改券改积分、库存、退款申请等场景的多语言词条；
  - 相关单元测试断言同步改为 key，避免文案变更导致测试误报。
- 文件：server/service/shop/order.go, server/api/v1/shop/order.go, server/utils/i18n/i18n.go, server/service/shop/order_status_pending_confirm_test.go
- 有什么用：Uni/Web 在切语种时，订单错误提示能稳定展示当前语言，不再“时中时英时中文”。
- 能防止什么：
  - 服务层新增错误后 API 层直接透传英文/中文，造成多语言失效；
  - 文案修改影响测试稳定性，导致回归噪音。
- 小白怎么操作：
  1. 新增订单错误时，只在 service 返回 error key（不要写中文句子）。
  2. API 层统一用 i18n 翻译后返回给前端。
  3. 如果前端看到 key 原文，优先检查 i18n 词典是否缺该 key。

### 4.16 充值订单 API 错误返回收口（与订单 API 风格对齐）
- 改了什么：
  - 试衣币充值订单 API 的错误返回已统一收口到 client 包通用 helper（`failClientWithErr/failClientWithKey`）；
  - 将各接口重复的 `FailWithMessage(i18n.T(...))` 调用统一替换，避免漏翻译或写法漂移；
  - 权限不足、参数非法、业务错误三类返回路径统一到同一处理方式。
- 文件：server/api/v1/client/tryon_recharge_order.go
- 有什么用：充值订单链路与商城订单链路在错误返回层面保持同一规范，后续维护和审计更稳定。
- 能防止什么：
  - 新接口复制旧代码时遗漏 i18n 翻译；
  - 同一业务在不同接口返回格式不一致，导致前端兜底逻辑复杂化。
- 小白怎么操作：
  1. 后续在充值订单 API 新增错误返回时，优先调用 `failClientWithErr/failClientWithKey`，不要手写 `FailWithMessage`。
  2. 服务层继续返回 error key，由 API 层统一翻译后再输出给前端。

### 4.17 Client 侧交易资源 API 错误 helper 通用化
- 改了什么：
  - 在 client API 包新增通用错误返回 helper（`failClientWithErr/failClientWithKey`）；
  - 将“我的模特”“我的衣橱”“充值订单”三条链路统一切换到通用 helper；
  - 移除充值订单文件内的专用 helper，避免同类工具函数重复实现。
- 文件：server/api/v1/client/error_helper.go, server/api/v1/client/tryon_model.go, server/api/v1/client/tryon_cloth.go, server/api/v1/client/tryon_recharge_order.go
- 有什么用：client 侧核心交易/资源 API 错误返回风格统一，后续新增接口不容易出现“有的翻译、有的直返”的不一致。
- 能防止什么：
  - 同一模块内 helper 命名和行为不一致导致维护成本上升；
  - 改一处漏一处，出现多语言回退不稳定。
- 小白怎么操作：
  1. client API 层遇到错误返回时，优先使用 `failClientWithErr/failClientWithKey`。
  2. 服务层继续返回 i18n key，不要在 service 中写死中文文案。

### 4.18 Client 用户与试衣任务接口收口补齐
- 改了什么：
  - 试衣任务 API（`tryon_task`）中固定 key 错误返回全面替换为 `failClientWithKey`；
  - 用户 API（`user`）中直接透传 `err.Error()` 的关键分支统一改为 `failClientWithErr`；
  - 保持动态拼接错误（如 `setFail:xxx`）不改动，避免改变既有前端提示结构。
- 文件：server/api/v1/client/tryon_task.go, server/api/v1/client/user.go
- 有什么用：把 client 侧高频接口的错误处理进一步统一，减少“同一个模块不同文件写法不同”的维护噪音。
- 能防止什么：
  - 新人复制旧代码时继续扩散 `i18n.T(c, err.Error())` 直返写法；
  - 代码审计时因风格不一致漏掉关键错误出口。
- 小白怎么操作：
  1. 在 client API 层需要返回错误 key 时，用 `failClientWithKey`。
  2. 需要返回 service 透传错误时，用 `failClientWithErr`。

### 4.19 Client 账户与签到收藏接口风格对齐
- 改了什么：
  - 用户 API（`user`）中固定 key 错误返回进一步统一替换为 `failClientWithKey`，与前面已替换的 `failClientWithErr` 一起形成完整收口；
  - 签到 API（`sign_in`）与收藏 API（`collect`）同步切换到 `failClientWithKey/failClientWithErr`；
  - 用户接口中带上下文拼接的动态提示（`setFail:具体错误`）保持原样，避免影响前端现有展示逻辑。
- 文件：server/api/v1/client/user.go, server/api/v1/client/sign_in.go, server/api/v1/client/collect.go
- 有什么用：client 端账户与签到收藏等高频入口的错误返回风格统一，后续排查和审计更直接。
- 能防止什么：
  - 多个 API 文件并行演进时，错误返回方式再次分叉；
  - 前端因历史接口返回格式差异产生额外兼容分支。
- 小白怎么操作：
  1. 写 client API 新接口时优先使用 `failClientWithKey/failClientWithErr`。
  2. 只有确实需要拼接上下文时，才保留 `key + 详细信息` 形式。

### 4.20 Client 地址与语言区号接口收口补齐
- 改了什么：
  - 地址 API（`address`）中的固定 key 错误返回统一替换为 `failClientWithKey`；
  - 语言 API（`language`）与国际区号 API（`phone_area_code`）同步完成同样替换；
  - 保持成功返回与本地化数据结构不变，仅收口错误出口。
- 文件：server/api/v1/client/address.go, server/api/v1/client/language.go, server/api/v1/client/phone_area_code.go
- 有什么用：client 侧基础配置与资料相关接口的错误处理风格与交易链路一致，减少维护时的认知切换。
- 能防止什么：
  - 同一目录下部分文件已收口、部分未收口导致代码规范断层；
  - 后续批量审计错误返回时遗漏基础模块。
- 小白怎么操作：
  1. 这三类接口新增错误返回时直接复用 `failClientWithKey`。
  2. 若后续出现需要透传服务错误，再按需使用 `failClientWithErr`。

### 4.21 Client 配置与统计接口收口补齐
- 改了什么：
  - 系统配置 API（`sys_config`）中的固定 key 错误返回统一替换为 `failClientWithKey`；
  - 访客统计 API（`visitor`）、积分记录 API（`point_record`）与外链域名 API（`external_link_domain`）同步完成同类替换；
  - 仅收口错误出口，原有成功返回与 `i18n.TWithSuffix(...)` 动态错误拼接逻辑保持不变。
- 文件：server/api/v1/client/sys_config.go, server/api/v1/client/visitor.go, server/api/v1/client/point_record.go, server/api/v1/client/external_link_domain.go
- 有什么用：client 目录下配置、统计、管理相关接口的错误返回风格基本统一，后续维护可按同一模式扩展。
- 能防止什么：
  - 不同模块沿用历史写法导致错误出口风格碎片化；
  - 代码审计时因返回方式不统一遗漏关键分支。
- 小白怎么操作：
  1. 在这些模块新增固定 key 错误提示时，直接用 `failClientWithKey`。
  2. 需要携带底层错误上下文时，继续使用 `i18n.TWithSuffix(...)` 等动态方式。

### 4.22 Client 目录错误出口扫尾核查
- 改了什么：
  - 对 `server/api/v1/client` 全目录执行旧写法扫描，确认固定 key 与 `err.Error()` 直返已收口完成；
  - 扫描结果仅保留 `error_helper.go` 内部翻译出口（这是预期的统一出口，不属于遗留问题）。
- 文件：server/api/v1/client/error_helper.go
- 有什么用：后续做多语言与安全审计时，能够快速确认 client 层错误返回的统一性。
- 能防止什么：
  - 局部收口后误以为“已全量完成”，但仍留有分散旧写法；
  - 新改动回归时缺少基线，不易快速判断是否风格回退。
- 小白怎么操作：
  1. 例行巡检时用同一 grep 规则扫描 `server/api/v1/client/**/*.go`。
  2. 如果命中不在 `error_helper.go`，优先改为 `failClientWithKey/failClientWithErr`。

---

## 5. 上传链路与文件库安全

### 5.1 文件库管理接口 API 层角色白名单
- 改了什么：delete/edit/list/import/listFolders 等敏感接口先做角色白名单判断。
- 文件：server/api/v1/example/exa_file_upload_download.go
- 有什么用：Casbin 之外再加一层硬校验。
- 能防止什么：误授权角色越权操作文件库。
- 小白怎么操作：仅让 888/8881/9528 管理文件库。

### 5.2 SignURL 路径规范化校验
- 改了什么：签名接口先校验 filePath（禁止空、路径穿越、query/fragment、超长）。
- 文件：server/api/v1/example/exa_file_upload_download.go, server/model/example/request/exa_file_upload_and_downloads.go
- 有什么用：只对规范路径签名。
- 能防止什么：构造恶意路径骗签名。
- 小白怎么操作：前端只传标准文件路径，不传整段危险 URL。

### 5.3 上传模型增加创建者字段
- 改了什么：exa_file_upload_and_downloads 增加 created_by。
- 文件：server/model/example/exa_file_upload_download.go
- 有什么用：后续可按人隔离文件操作范围。
- 能防止什么：同权限角色互删他人文件。
- 小白怎么操作：如历史数据 created_by=0，先用超级管理员治理。

### 5.4 上传输入严格校验（大小、扩展名、MIME）
- 改了什么：上传前校验空文件、超大小、危险扩展名、危险内容类型。
- 文件：server/service/example/exa_file_upload_download.go
- 有什么用：在入口阻断高风险文件。
- 能防止什么：脚本木马、可执行文件伪装上传。
- 小白怎么操作：
  1. security_upload_max_size_mb 控制体积。
  2. security_upload_strict_validation_enabled 控制严格模式。

### 5.5 删除对象前安全清洗 key
- 改了什么：删除文件时对对象存储 key 做清洗；非法 key 不执行对象删除，仅删库记录。
- 文件：server/service/example/exa_file_upload_download.go
- 有什么用：避免把异常 key 传给 OSS 删除接口。
- 能防止什么：误删前缀、路径穿越式删除。
- 小白怎么操作：发现历史脏 key，先清理数据再做批量删。

### 5.6 编辑/列表/导入边界保护
- 改了什么：
  - 编辑文件名校验 ID 与长度；
  - 列表限制分页上限、关键词长度；
  - 导入 URL 限批次、限长度、校验协议与路径。
- 文件：server/service/example/exa_file_upload_download.go
- 有什么用：防资源放大与参数滥用。
- 能防止什么：超大分页拖库、批量脏 URL 导入。
- 小白怎么操作：导入 URL 单次别超过限制，失败看具体条目序号。

### 5.7 文件归属隔离（非超级只能操作自己）
- 改了什么：delete/edit/list 操作都按 created_by 做 owner scope；超级管理员保留全局范围。
- 文件：server/service/example/exa_file_upload_download.go
- 有什么用：同角色不同用户之间实现资源隔离。
- 能防止什么：横向越权（删除/修改/查看别人文件）。
- 小白怎么操作：普通管理员只会看到自己上传的数据，这是预期行为。

### 5.8 客服图片上传适配新上传签名 + 强校验
- 改了什么：客服图片上传改为传 operatorUserID，并加入大小/扩展名/MIME 严格校验。
- 文件：server/plugin/customer_service/api/message.go
- 有什么用：客服通道上传也统一进入安全规则。
- 能防止什么：聊天通道上传恶意文件。
- 小白怎么操作：在客服配置里限制允许扩展名和大小。

### 5.9 Hotlink 权限收口
- 改了什么：listFolders 改为后台管理角色；启动自动删除 8080 的历史授权。
- 文件：server/initialize/hotlink_init.go
- 有什么用：客户端角色不再可枚举存储目录。
- 能防止什么：目录探测与资源枚举。
- 小白怎么操作：若老环境还查得到目录，重启后复查 Casbin。

---

## 6. WebSocket 安全加固

### 6.1 Origin 校验升级（同源/同主域/本地联调）
- 改了什么：Origin 校验支持同 host、同根域、loopback 互通；IP 直连不走根域比较。
- 文件：server/plugin/customer_service/api/ws.go
- 有什么用：兼顾安全与本地联调可用性。
- 能防止什么：跨站恶意 Origin 建连。
- 小白怎么操作：线上使用正式域名；本地允许 localhost/127.0.0.1/::1 互通。

### 6.2 WS token 传递策略收口
- 改了什么：优先 Authorization/子协议取 token，query token 受开关控制，默认禁用。
- 文件：server/plugin/customer_service/api/ws.go
- 有什么用：减少 token 出现在 URL 中。
- 能防止什么：URL token 泄露（日志、代理、分享链接）。
- 小白怎么操作：
  1. 客户端升级为 Authorization 或 Sec-WebSocket-Protocol。
  2. 保持 security_ws_allow_query_token=false。

### 6.3 单 IP 连接数上限
- 改了什么：用户通道和坐席通道都加入每 IP 最大连接数限制。
- 文件：server/plugin/customer_service/api/ws.go
- 有什么用：防止单 IP 建大量连接耗尽资源。
- 能防止什么：WS 连接洪泛。
- 小白怎么操作：调 security_ws_max_conns_per_ip。

---

## 7. 导出下载链路安全

### 7.1 一次性导出 token + 过期清理
- 改了什么：导出改为先拿一次性 token，再用 token 下载；token 用后删除。
- 文件：server/api/v1/system/sys_export_template.go
- 有什么用：下载链路更可控，不可长期复用。
- 能防止什么：历史 token 被二次利用。
- 小白怎么操作：前端先调 export 接口获取 payload，再执行下载。

### 7.2 Header 传 token，query token 变可选兼容
- 改了什么：支持 Authorization 或 X-Export-Token；query token 由开关控制。
- 文件：server/api/v1/system/sys_export_template.go
- 有什么用：默认更安全，兼容期仍可临时放开。
- 能防止什么：URL token 泄露。
- 小白怎么操作：生产建议 security_export_allow_query_token=false。

### 7.3 Web 前端下载实现改造
- 改了什么：前端导出组件改为 fetch + Authorization: ExportToken 下载，并处理文件名解析。
- 文件：web/src/components/exportExcel/exportExcel.vue, web/src/components/exportExcel/exportTemplate.vue
- 有什么用：与后端新安全策略一致。
- 能防止什么：前端继续用 query token 导致后端拒绝或泄露。
- 小白怎么操作：如果老前端还打不开下载，先更新前端组件再上线。

---

## 8. 路由与网络入口层防护

### 8.1 Trusted Proxies 可配置
- 改了什么：Gin Trusted Proxies 改为可从系统参数/环境变量读取。
- 文件：server/initialize/router.go
- 有什么用：正确识别真实客户端 IP。
- 能防止什么：错误代理链下 IP 伪造或封禁误伤。
- 小白怎么操作：按你的网关实际地址设置 security_trusted_proxies。

### 8.2 PublicGroup 默认叠加封禁检查与统一限流
- 改了什么：公开路由链路统一挂 BanIPCheck + PublicRateLimit。
- 文件：server/initialize/router.go
- 有什么用：不用逐接口重复写防护。
- 能防止什么：某些公开接口漏加防护。
- 小白怎么操作：看日志里是否有 rate:public:ip 命中记录。

---

## 9. 种子数据与权限自愈补齐

### 9.1 API 种子补齐关键敏感接口
- 改了什么：种子中明确包含 /system/reloadSystem、/clientUser/setClientUserInfo 等关键接口。
- 文件：server/source/system/api.go
- 有什么用：新环境初始化后权限体系更完整。
- 能防止什么：接口存在但 API 表缺记录，导致权限流程异常。
- 小白怎么操作：新库初始化后检查 sys_apis 是否包含这些接口。

### 9.2 Casbin 种子补齐客户端个人资料权限
- 改了什么：为 8080 角色补齐 /clientUser/setClientUserInfo 权限种子。
- 文件：server/source/system/casbin.go
- 有什么用：避免客户端改个人资料时报无权限。
- 能防止什么：上线后出现“功能可见但永远失败”。
- 小白怎么操作：权限异常时先核对 casbin_rule 种子是否落库。

---

## 10. 文案与运维配套

### 10.1 i18n 增加安全相关提示词
- 改了什么：新增 requestTooFrequent、visitorHeartbeatTooFrequent、loginLocked、tryonCreateTooFrequent 等多语言文案。
- 文件：server/utils/i18n/i18n.go
- 有什么用：限流与封控场景返回用户可读提示。
- 能防止什么：前端看到生硬英文或 key 泄露。
- 小白怎么操作：新增安全错误码时同步补 i18n。

### 10.2 上线手册文档
- 改了什么：新增安全加固运行手册，覆盖上线、验证、回滚、SQL 检查。
- 文件：docs/security-hardening-runbook.md
- 有什么用：团队可按步骤执行，不靠口头传达。
- 能防止什么：误操作、漏验证、回滚无章法。
- 小白怎么操作：上线前先照手册做备份，再按步骤启开关。

### 10.3 README 入口补齐
- 改了什么：中英文 README 增加 runbook 链接。
- 文件：README.md, README-en.md
- 有什么用：新同学能快速找到安全上线流程。
- 能防止什么：文档存在但没人知道在哪。
- 小白怎么操作：把 runbook 当作发布前必读文档。

---

## 11. 配置总表（小白常用）

以下是本次重点安全开关，建议先用默认值上线，再按流量调整：

- security_init_api_enabled：生产建议 false
- security_init_api_private_network_only：建议 true
- security_public_rate_limit_enabled：建议 true
- security_public_rate_limit_window_seconds：建议 60
- security_public_rate_limit_max_requests：建议 300
- security_login_ip_rate_limit_per_minute：建议 30
- security_login_ip_rate_limit_window_seconds：建议 60
- captcha_rate_limit：建议 10
- captcha_rate_limit_window_seconds：建议 60
- security_sys_error_create_rate_limit_per_minute：建议 30
- security_sys_error_create_rate_limit_window_seconds：建议 60
- security_visitor_heartbeat_rate_limit_per_minute：建议 120
- security_visitor_heartbeat_dedupe_seconds：建议 3
- security_tryon_create_rate_limit_per_minute：建议 20
- security_tryon_create_concurrency_limit：建议 2
- security_order_create_rate_limit_per_minute：建议 30
- security_order_pending_limit_per_user：建议 10
- security_tryon_recharge_create_rate_limit_per_minute：建议 20
- security_tryon_recharge_pending_limit_per_user：建议 8
- order_pending_confirm_close_minutes：建议 180
- tryon_recharge_pending_confirm_close_minutes：建议 180
- security_attack_auto_ban_enabled：建议 true
- security_attack_auto_ban_window_seconds：建议 3600
- security_attack_auto_ban_duration_minutes：建议 60
- security_attack_auto_ban_threshold_default：建议 10
- security_attack_auto_ban_threshold_sys_error_rate_limit：建议 30
- security_export_allow_query_token：建议 false
- security_ws_allow_query_token：建议 false
- security_ws_max_conns_per_ip：建议 8
- security_upload_max_size_mb：建议 20
- security_upload_strict_validation_enabled：建议 true
- security_trusted_proxies：按网关实际设置
- security_cleanup_legacy_casbin_overgrant：默认 false，仅维护窗口临时 true

---

## 12. 一键检查清单（上线后）

### 12.1 错误信息直透巡检（本轮新增）
- 改了什么：对 client/shop/system 三条 API 线的 `err.Error()` 对外返回做了统一收口，未知错误不再直接下发原始文本。
- 文件：server/api/v1/client/error_helper.go, server/api/v1/client/tryon_task.go, server/api/v1/client/point_record.go, server/api/v1/client/user.go, server/api/v1/shop/order.go, server/api/v1/system/sys_error.go
- 有什么用：把“内部错误文本/数据库报错/未定义 key”直接暴露给前端的风险降到最低。
- 能防止什么：
  1. 前端拿到内部实现细节（表名、字段、第三方返回原文）后被用于探测。
  2. i18n 词条缺失时把 key 原样展示给用户。
- 小白怎么操作：
  1. 新增 API 错误返回时，优先走 `failClientWithErr` 或 `failClientWithKey`。
  2. service 层如需返回可翻译错误，请返回 i18n key；不要返回中文句子或底层原始错误。
  3. API 层需要记录细节时写日志（zap.Error），不要拼接到 response 文案里。

### 12.2 i18n key 存在性检查能力（本轮新增）
- 改了什么：i18n 工具新增 `HasKey` 判断函数，供 API 错误出口在翻译前确认 key 是否存在。
- 文件：server/utils/i18n/i18n.go
- 有什么用：遇到未知 key 时可以稳定回落到通用失败文案，而不是把 key 原样返回。
- 能防止什么：新功能尚未补词条时，用户界面出现“裸 key”或内部命名信息。
- 小白怎么操作：
  1. 写错误返回 helper 时先用 `i18n.HasKey(key)` 判定。
  2. 判定失败统一回落 `fail` 或业务默认 key。

### 12.3 system/shop 第二批错误出口收口（本轮新增）
- 改了什么：继续移除管理端与商城接口中参数绑定错误、业务错误的 `err.Error()` 直返与拼接式返回。
- 文件：server/api/v1/system/sys_api.go, server/api/v1/system/sys_api_token.go, server/api/v1/system/db_inspector.go, server/api/v1/system/sys_system.go, server/api/v1/system/sys_params.go, server/api/v1/shop/cart.go, server/api/v1/shop/coupon.go, server/api/v1/shop/banner.go, server/api/v1/shop/category.go, server/api/v1/shop/comment.go
- 有什么用：扩大“错误不泄露内部细节”的覆盖面，减少后台接口被探测时的信息回显。
- 能防止什么：
  1. 参数校验失败直接回显内部结构体错误。
  2. 服务层异常文本拼接后直接暴露给前端。
- 小白怎么操作：
  1. 这批文件后续新增错误返回时继续沿用固定文案或统一 helper，不要拼接 `err.Error()`。
  2. 详细异常只写日志，不回给客户端。

### 12.4 高密度 CRUD 接口第三批收口（本轮新增）
- 改了什么：继续清理 system/shop 历史高密度接口中的 `err.Error()` 直返，覆盖用户、菜单、版本与多组商城 CRUD 接口。
- 文件：server/api/v1/system/sys_user.go, server/api/v1/system/sys_menu.go, server/api/v1/system/sys_version.go, server/api/v1/shop/tag.go, server/api/v1/shop/promotion.go, server/api/v1/shop/sku_spec.go, server/api/v1/shop/kefu_service.go
- 有什么用：把“批量生成类 API 文件”中的历史错误回显统一到安全文案，降低后续复制粘贴扩散风险。
- 能防止什么：
  1. 后台和商城管理接口在异常场景回显内部错误细节。
  2. 攻击者利用差异化报错信息推断表结构和业务分支。
- 小白怎么操作：
  1. 新增同类 CRUD 接口时，参数异常统一返回“参数错误”。
  2. 业务异常统一返回“创建/更新/删除/查询/获取失败”等固定文案。
  3. 详细异常继续保留在 `zap.Error(err)` 日志里排障。

### 12.5 持续扫尾（order + 通用系统列表页）（本轮新增）
- 改了什么：
  1. 清理 `order` 中残留 `i18n.TWithSuffix(..., err.Error())` 错误回显。
  2. 清理系统通用列表页接口中的参数错误直透（操作日志、登录日志、按钮权限）。
- 文件：server/api/v1/shop/order.go, server/api/v1/system/sys_operation_record.go, server/api/v1/system/sys_login_log.go, server/api/v1/system/sys_authority_btn.go
- 有什么用：把“高频查询/列表页”接口的错误返回也统一到安全文案，减少遗漏面。
- 能防止什么：
  1. 物流查询等场景在失败时把下游错误细节透给前端。
  2. 管理后台列表接口在参数异常时返回内部校验细节。
- 小白怎么操作：
  1. 对接前端时只依赖标准 `msg` 与 `code`，不要依赖底层错误文本。
  2. 遇到排障需求，统一在后端日志检索 `zap.Error(err)`。

### 12.6 角色管理接口收口（本轮新增）
- 改了什么：清理 `sys_authority` 中参数校验与业务失败场景的 `err.Error()` 拼接回显。
- 文件：server/api/v1/system/sys_authority.go
- 有什么用：角色创建/复制/删除/更新等核心权限接口的错误输出口径统一，避免暴露内部约束细节。
- 能防止什么：攻击者利用角色管理接口错误差异探测权限模型与数据关系。
- 小白怎么操作：角色管理相关 API 出错时，只看统一失败文案；具体原因在后端日志中排查。

### 12.7 system/shop 错误直透清零（本轮新增）
- 改了什么：对 `server/api/v1/system` 与 `server/api/v1/shop` 继续分批清理 `FailWithMessage(err.Error())`、`"..."+err.Error()`、`TWithSuffix(..., err.Error())`，将参数异常统一为“参数错误”，业务异常统一为固定失败文案。
- 文件：
  1. shop：server/api/v1/shop/good_purchase.go, server/api/v1/shop/good.go, server/api/v1/shop/couponorderuser.go
  2. system：server/api/v1/system/sys_banned_ip.go, server/api/v1/system/sys_casbin.go, server/api/v1/system/sys_dictionary.go, server/api/v1/system/sys_dictionary_detail.go, server/api/v1/system/auto_code_history.go, server/api/v1/system/auto_code_template.go, server/api/v1/system/sys_auto_code.go, server/api/v1/system/auto_code_package.go, server/api/v1/system/auto_code_mcp.go, server/api/v1/system/auto_code_plugin.go, server/api/v1/system/sys_export_template.go
- 有什么用：system/shop 两条后端主链路的历史错误回显已统一到安全口径，排障走日志，前端不再收到底层原始错误文本。
- 能防止什么：
  1. 通过错误文案探测内部实现细节（字段名、依赖报错、流程分支）。
  2. 历史拼接式错误在新增接口中继续扩散。
- 小白怎么操作：
  1. 新增 API 错误返回时继续遵循“参数错误 + 固定业务失败文案 + zap 日志细节”三原则。
  2. 验证时执行全局扫描：`grep_search` 模式 `FailWithMessage(...err.Error...)|TWithSuffix(..., err.Error())`，当前 system/shop 结果应为 0。

### 12.8 api/v1 全量清零（本轮新增）
- 改了什么：继续清理 example 模块剩余点位后，`server/api/v1/**` 在同一口径扫描下已无命中。
- 文件：server/api/v1/example/exa_customer.go, server/api/v1/example/exa_file_upload_download.go, server/api/v1/example/exa_attachment_category.go, server/api/v1/example/exa_breakpoint_continue.go
- 有什么用：后端 API 主入口整体统一到“固定失败文案 + 日志留痕”模式，跨模块行为一致。
- 能防止什么：
  1. 示例/工具接口泄露底层错误文本，被用来侧信道探测。
  2. 老代码风格在新接口中继续复制扩散。
- 小白怎么操作：
  1. 新增 API 时不要写 `FailWithMessage(err.Error(), c)` 或 `"xx"+err.Error()`。
  2. 参数失败返回“参数错误”，业务失败返回固定文案，细节写 `zap.Error(err)`。

### 12.9 插件 API 层清零（本轮新增）
- 改了什么：将清理范围从 `api/v1` 扩展到 `server/plugin/**/api`，完成 wxpay、geo、announcement、email、customer_service 的 `err.Error()` 直返与拼接回显收口。
- 文件：
  1. server/plugin/wxpay/api/api.go
  2. server/plugin/geo/api/api.go
  3. server/plugin/announcement/api/info.go
  4. server/plugin/email/api/sys_email.go
  5. server/plugin/customer_service/api/agent.go
  6. server/plugin/customer_service/api/blacklist.go
  7. server/plugin/customer_service/api/config.go
  8. server/plugin/customer_service/api/conversation.go
  9. server/plugin/customer_service/api/quick_reply.go
  10. server/plugin/customer_service/api/message.go
- 有什么用：插件侧接口与主线 API 的错误返回口径统一，避免“主系统安全、插件直透”造成短板。
- 能防止什么：
  1. 插件管理接口泄露内部报错细节（权限判断、会话状态、数据库异常等）。
  2. 插件复制模板继续扩散 `err.Error()` 直返模式。
- 小白怎么操作：
  1. 插件 API 参数错误统一返回“参数错误”。
  2. 权限检查失败与业务失败返回固定语义文案（如“无权限访问该会话”“发送失败”），详细错误只写日志。
  3. 复检口径：`server/**/*.go` 下 `FailWithMessage(...err.Error...)|TWithSuffix(..., err.Error())|"..."+err.Error()` 当前应为 0。

### 12.10 中间件层错误出口收口（本轮新增）
- 改了什么：修复鉴权和限流中间件中的原始错误文本回传，避免在认证失败和限频场景直接下发底层错误。
- 文件：server/middleware/jwt.go, server/middleware/limit_ip.go
- 有什么用：把中间件层也纳入统一错误输出策略，避免“API 已收口但中间件还直出”的旁路泄露。
- 能防止什么：
  1. JWT 解析失败时暴露签名/格式等内部错误细节。
  2. 限流检查异常时暴露底层 Redis/中间件错误文本。
- 小白怎么操作：
  1. 中间件返回统一文案（如 `tokenInvalid`、`请求太过频繁，请稍后再试`）。
  2. 详细错误只进日志，禁止通过 `response` 或 `c.JSON` 直接返回 `err.Error()`。
  3. 全仓复检口径：`server/**/*.go` 下相关模式当前应为 0。

### 12.11 Go 依赖漏洞清零（本轮新增）
- 改了什么：
  1. 升级可修复三方依赖：`golang.org/x/net` `v0.50.0 -> v0.53.0`、`github.com/ulikunitz/xz` `v0.5.12 -> v0.5.15`。
  2. 同步升级相关 `x/*` 依赖链（`x/crypto`、`x/sync`、`x/sys`、`x/text`、`x/mod`、`x/tools`）。
  3. 在 `go.mod` 固化补丁级工具链：`toolchain go1.26.3`，用于覆盖标准库漏洞修复版本。
- 文件：server/go.mod, server/go.sum
- 有什么用：把 govulncheck 报告中的“可达漏洞”从 13 个降为 0。
- 能防止什么：
  1. `x/net` HTTP/2 DoS 相关漏洞。
  2. `xz` 解压内存/资源消耗类漏洞。
  3. 标准库补丁级安全缺陷在构建环境中长期滞留。
- 小白怎么操作：
  1. 后端构建机执行 `go version`，确认实际使用 `go1.26.3` 或更高补丁版本。
  2. 发布前执行 `govulncheck ./...`，当前预期：`Your code is affected by 0 vulnerabilities`。

### 12.12 前端依赖审计压降（本轮新增）
- 改了什么：
  1. web：执行 `npm audit fix` 并移除历史遗留 Vue CLI 依赖链（`@vue/cli-*`, `@vue/cli-service`）及 `npm` 运行时依赖。
  2. uni：执行 `npm audit fix` 完成非破坏性升级。
  3. 验证结论：uni 在当前 `@dcloudio/uni-* 3.0.0-405...` 生态下，单独升级 `vite` 到 `5.4.21` 会把漏洞从 51 提高到 54，已回退到 `5.2.8`。
- 文件：web/package.json, web/package-lock.json, uni/package-lock.json
- 有什么用：前端漏洞面明显下降且不影响构建通过。
- 能防止什么：
  1. 旧 CLI 链路引入的高危 transitive 漏洞长期残留。
  2. 运行时无关依赖（如 `npm`）将额外漏洞打包进审计面。
- 小白怎么操作：
  1. web 当前审计从 `31` 降到 `7`（剩余主要是 `wangeditor` 无官方修复与 `vue3-sfc-loader` 链路）。
  2. uni 当前审计从 `70` 降到 `51`（主要集中在 `@dcloudio/uni` 生态，修复多为 `--force` 级别，需单独联调窗口处理）。
  3. 构建验证：`web npm run build` 通过；`uni npm run build:h5` 退出码 `0`。

---

## 13. 下一步建议（安全）

1. 继续做全仓 `err.Error()` 对外返回清点（优先 public/client 路由），把剩余拼接式错误逐步迁移到 key 化返回。
2. 进行 Go 依赖最小升级（jwt、redis、x/net）并跑一次回归测试，降低已知高危依赖风险。
3. 对 web/uni 依赖高危项按“可直接升级/需联调”分组治理，先处理可无破坏升级的包。

### 12.1 权限检查
- 角色 8080 访问 /system/reloadSystem 应失败。
- 角色 8080 访问 /fileUploadAndDownload/listFolders 应失败。
- 角色 888/8881 对管理端接口应正常。

### 12.2 限流检查
- 连续刷 /base/captcha 应出现“请求过于频繁”。
- 连续刷 /sysError/createSysError 应出现限流并记录攻击日志。
- 连续刷访客心跳应触发限流或去重。
- 连续高频调用商城下单/购物车下单，应出现频率限制提示。
- 单用户累计待支付/待确认订单过多时，应被拒绝新建订单。
- 连续高频创建试衣币充值单，应出现频率限制提示。
- 将订单提交到待后台确认后，超过配置分钟数应被任务自动关闭。
- 对同一订单重复执行同一状态更新（例如重复取消）应幂等，不应再次变动库存/销量。

### 12.3 上传检查
- 上传 .php/.exe 应被拒绝。
- 超过 security_upload_max_size_mb 应被拒绝。
- 非超级管理员只能看/改/删自己 created_by 的文件记录。

### 12.4 导出检查
- 新前端下载走 Authorization: ExportToken。
- 当 security_export_allow_query_token=false 时，query token 下载应被拒绝。

### 12.5 WebSocket 检查
- query token 在默认配置下应被拒绝。
- 超过单 IP 连接数应被拒绝。
- 非同源/非同根域 Origin 应被拒绝。

---

## 13. 上线步骤（小白版）

1. 先备份数据库（至少 casbin_rule）。
2. 合并代码并部署。
3. 在后台系统参数里设置 security 组推荐值。
4. 重启服务。
5. 按第 12 章做冒烟检查。
6. 若要清理历史过授权：
   - 临时开 security_cleanup_legacy_casbin_overgrant=true；
   - 启动一次；
   - 立即改回 false。
7. 导出与 WS 老客户端若未升级，可短期打开 query 兼容开关，升级完成后再关回去。

---

## 14. 回滚步骤（小白版）

1. 先把高风险开关恢复为保守值（例如关清理开关）。
2. 如权限异常，恢复 casbin_rule 备份。
3. 重启服务。
4. 再执行第 12 章检查，确认恢复。

---

## 15. 额外说明：工作区中的“非本次安全主线”改动

以下文件在当前 git status 中也有变化，但不是本次安全主线的核心条目，先登记不展开：
- server/source/geo/generate_iso_geo_sql.mjs
- server/source/geo/geo_language_country_templates.sql
- server/source/geo/geo_zh_translate_cache.json

说明：
- server/initialize/new_modules_init.go
- server/initialize/shop_init.go
- server/router/system/sys_auto_code.go
以上三项已并入第 2 章（权限与越权防护），不再归为“非主线”。

已补充独立文档：WORKSPACE_NON_SECURITY_CHANGES_RECORD.md

---

## 16. 2026-05-20 全面检测收尾修复（静态诊断清零）

### 16.1 试衣分割重试循环修复“无条件终止”
- 改了什么：`tryon_task` 在阿里分割重试循环中，把失败分支的直接 `return` 调整为 `continue`，保留循环后统一返回最后错误。
- 文件：server/service/client/tryon_task.go
- 有什么用：真正执行多次重试，而不是首轮失败就提前退出。
- 能防止什么：重试逻辑名存实亡，导致容错能力下降。
- 小白怎么操作：无需额外配置，保持现有重试参数即可。

### 16.2 去除精修分支冗余 nil 判断
- 改了什么：`tryon_task` 中 `modelCfg` 已在上游判空，移除后续冗余 `modelCfg != nil` 条件。
- 文件：server/service/client/tryon_task.go
- 有什么用：消除静态分析“恒真条件”噪音，代码可读性更稳定。
- 能防止什么：后续维护时被误导为“这里可能为 nil”。
- 小白怎么操作：无感知变更。

### 16.3 依赖直连标记修正
- 改了什么：`go.mod` 中 `golang.org/x/net` 从 `// indirect` 调整为直接依赖标记。
- 文件：server/go.mod
- 有什么用：与真实引用关系一致，减少依赖诊断误报。
- 能防止什么：CI 或 IDE 对依赖关系反复提示“应为 direct”。
- 小白怎么操作：`go mod tidy` 时不再反复出现该提示。

### 16.4 封禁中间件时间计算规范化
- 改了什么：封禁剩余分钟计算由 `t.Sub(time.Now())` 改为 `time.Until(t)`。
- 文件：server/middleware/ban_ip.go
- 有什么用：表达更直接，符合 Go 常见时间差写法。
- 能防止什么：静态检查告警堆积，影响真正问题排查。
- 小白怎么操作：无感知变更。

### 16.5 MCP 列表接口资源释放顺序修正
- 改了什么：`MCPList` 中先检查 `client.NewClient(...)` 错误，再 `defer testClient.Close()`。
- 文件：server/api/v1/system/auto_code_mcp.go
- 有什么用：避免在创建失败时对无效对象执行 Close。
- 能防止什么：潜在空指针/无效资源释放路径。
- 小白怎么操作：无感知变更。

### 16.6 本轮验证结果
- `get_errors` 针对上述 4 个目标文件复查：已全部 `No errors found`。
- 后端回归验证通过：
  1. `go test ./... -run TestDoesNotExist -count=1`
  2. `go test ./service/shop ./service/client ./task ./api/v1/client -count=1`

