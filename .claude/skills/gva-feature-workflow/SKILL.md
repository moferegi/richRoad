---
name: gva-feature-workflow
description: "GVA 功能开发标准工作流。用于新建/扩展插件或业务模块时，按 Model->Service->API->Router->Initialize->Plugin->Web/Uni 顺序生成可上线代码，并强制检查 Swagger、Casbin、菜单、数据类型一致性。"
---

# GVA 功能开发标准工作流

## 适用场景
- 新建业务模块、插件、后台页面、uni 页面
- 迁移老功能到可复用架构（试衣、电影等）
- 修复“接口已写但权限/菜单/初始化遗漏”类问题

## 关联技能
- 遇到 Uni 多语言展示一致性问题（尤其 SKU 规格切语、nameI18n/labelI18n/valueI18n 显示 key 或错误回退）时，必须同时应用 `uni-i18n-display-unification`。
- 推荐执行顺序：先按本工作流完成模块分层与接口，再按 `uni-i18n-display-unification` 完成展示层统一解析和回归验证。

## 必做顺序
1. 先设计数据模型与请求模型
- 后端先写 model 与 model/request
- 统一字段类型，避免 model/request/response 类型不一致
- 指针与非指针转换在 service 层显式处理

2. 再写 Service
- 只做业务与数据库，不接触 gin.Context
- 统一返回 (data, error) 或 error

3. 再写 API
- 参数绑定与校验在 API 层
- 用统一 response 包返回
- 每个对外接口补全 Swagger 注释

4. 再写 Router
- 按 public/private 分组
- 鉴权路由统一挂 JWTAuth + CasbinHandler
- 需要审计的写接口挂 OperationRecord

5. 初始化与注册
- initialize/api.go: 注册 sys_apis
- initialize/menu.go: 注册菜单并分配角色
- initialize/casbin.go 或 init 文件: 注册并收敛规则（必要时清理历史脏权限）
- initialize/gorm.go: AutoMigrate
- plugin.go: init 注册插件 + Register 挂载路由
- plugin/register.go: 匿名导入激活插件

6. 前端落地
- web/src/api 按模块封装接口
- web/src/plugin/<name>/view 增加页面
- uni 端文案必须使用 i18n key，禁止新增硬编码
- uni 新增/修改后端接口字段若会给用户展示，必须定义 i18n 方案（i18n key 或多语言对象）
- 后端返回多语言对象或 JSON 字符串（如 {mn:"",en:""}）时，必须使用 localText(value, langStore.locale) 解析后展示
- uni 的 toast/modal/弹窗标题与内容/placeholder/JS 拼接文案必须纳入 i18n
- 试衣间公告 announcement_content 必须按多语言解析，禁止 String(...) 直接渲染
- 新增词条时同步更新 layout/client/language（或项目实际 i18n 文件）
- 上传规则（格式/大小）优先走后端配置源

## 管理端多语言字段编辑规范（强制执行）
1. 后端多语言字段统一入口
- 只要字段最终会以多语言形式下发（对象或 JSON 字符串），管理端必须统一使用 `web/src/components/multilingual/multi-lang-editor.vue` 编辑
- 禁止在页面中手写分散的 `v-for + input/textarea` 多语言输入实现，避免交互与序列化策略不一致

2. 一键翻译能力必须保留
- 多语言编辑必须保留“一键翻译空白 / 覆盖翻译”能力，禁止去掉翻译工具栏后改为手动逐语填写
- 对翻译结果必须提供人工复核流程（尤其运营可见文案、支付文案、教程文案）

3. 富文本字段必须走多语言组件
- 对教程、公告、协议、活动说明等富文本多语言字段，仍需以 `MultiLangEditor` 为外层容器，并通过 richtext 模式或 slot 接入 `RichEdit`
- 禁止使用“单语富文本 + 其他语种纯文本”的混搭编辑方式
- 富文本一键翻译后必须检查标签结构、换行和链接是否正确
- 若页面通过 slot 自定义 `RichEdit`（如 `#editor` 插槽）渲染多语富文本，提交前必须显式执行结构校验钩子（推荐 `validateRichTextI18nStructure`），校验失败时阻断保存并提示用户修正
- 对 slot 富文本场景，需保留统一注释标识“富文本多语校验钩子”，便于代码巡检与批量审计

4. 保存与兼容约束
- 保存格式统一为多语言 JSON 对象字符串（如 `{ "zh": "...", "en": "..." }`），不得混存为单语纯字符串
- 编辑器加载历史值时必须兼容对象/JSON 字符串/旧单语字符串，并在保存时收敛为统一多语言结构
- 不得按启用语种裁剪已填写键，避免未启用语种内容在保存时被误删

## 试衣模型配置专项规范（强制）
1. 单一配置源
- 统一使用 tryon_models 作为试衣/试鞋/取衣/美肤模型入口
- 禁止新增或恢复 shoe_models 等并行模型列表
- C 端下发配置接口仅返回 tryon_models，不再暴露废弃模型键

2. 模型独立参数
- 每个模型必须自带 provider、mode、url、token、taskQueryUrl 等调用参数
- 禁止依赖全局 tryon_provider_mode / tryon_provider_url / tryon_provider_token 回退
- 管理端编辑器按模型能力展示字段，禁止“全模型统一堆叠无关参数”

3. 协同模型绑定规则
- 协同能力（如精修、美肤、分割）通过 modelKey 或模型内显式字段关联
- 试衣主模型优先通过 beautifyModelKey / refinerModelKey / parsingModelKey 绑定协同模型
- 禁止不同模型共享同一 token 配置作为默认回退
- 精修 token 留空只允许回退到当前模型 token，不允许跨模型回退

4. 参数治理
- 对当前运行链路无效的参数必须删除或隐藏，不得长期保留“占位参数”
- 兼容历史字段时，必须在迁移完成后清理数据库旧键与前端入口
- 配置说明需明确“哪些字段生效于哪些 provider/scene”

5. 统计与运营联动
- 新增或重命名试衣模型后，必须同步检查 tryonPointRecord 与 tryonTaskManage 的统计口径
- 至少验证：模型维度筛选、成功/失败计数、精修计数、扣退币统计在新模型下正确
- 若模型拆分为多子模型，需补充对应管理端可观测指标，避免运营看板失真

## Uni 多语言详细处理（强制执行）
1. 后端接口返回规范
- 参数绑定错误统一返回通用 key（如 invalidParams），禁止把绑定错误文本直接透传到前端
- 业务错误优先返回 i18n key（如 noPermission、querySuccess、createSuccess），避免返回中文句子
- 若历史接口暂时返回 msg 字符串，API 层应先收敛为可翻译 key，再逐步迁移 service 错误来源

2. Uni 页面消费规范
- 所有后端 msg 展示必须经过 resolveApiMessage(msg, fallbackKey)
- 禁止使用 data.msg || $t('xxx')、res.msg || '中文'、e?.message || $t('xxx') 这类直透模式
- 对后端多语言对象或 JSON 字符串字段，必须先 localText(value, langStore.locale) 再渲染
- 禁止直接 String(value) 渲染可能是多语言对象的字段

3. 适用范围（必须全覆盖）
- template 可见文本
- uni.showToast / uni.showModal 的 title、content
- placeholder、空态、按钮文案、校验提示
- JS 中拼接或 reject(new Error(...)) 的错误文本
- 上传/下单/支付/客服聊天等链路的异常提示

4. 词条补齐与语言覆盖
- 新增 key 后，必须同步补齐所有已启用语言，不能只补中英蒙
- 以项目当前语言集合为准逐个补齐（例如 zh-CN、en、mn、zh-TW 等）
- 提交前至少核对新增 key 在每个 locale 中均存在，避免缺词回退成 key 原文

5. 扫描与回归（提交前）
- 扫描高风险模式：data.msg ||、res.msg ||、e?.message ||、String(multilingualValue)
- 扫描 Uni 目录中的 toast/modal 调用，确认全部使用 i18n key 或 resolveApiMessage
- 对涉及多语言的配置项（重点 announcement_content）执行切换语言回归，确认内容随语言变更
- 新增后端 key 时，同步检查 uni/src/utils/i18n.js 与后端 i18n 字典是否均存在该 key

6. 非缓存类常见错误（必须规避）
- fallback 顺序错误：多语言对象/JSON 解析必须优先 locale，再 en，最后 zh，禁止先回退 zh 导致其他语种误显示中文
- 文案固化：禁止在初始化阶段把 localText(...) 结果写死到长期状态；应存 raw 值并在 computed/watch(locale) 中实时解析
- 竞态覆盖：切语言触发重拉时，需防旧请求晚到覆盖新语言数据（可用请求戳/最后一次请求保护）
- 语言状态不一致：所有语言判定统一来源于 langStore.locale（必要时回退 app-lang），避免页面内混用导致局部旧语种
- 系统消息不可回刷：对本地插入的系统提示（如客服聊天提示）应保存 i18nKey，而不是仅保存已翻译文本

7. 表单页面特殊规则（防止切语丢草稿）
- 编辑态页面切语言返回时，禁止无条件重拉并覆盖用户输入
- 推荐策略：查看模式可重拉；编辑模式仅刷新标签/下拉文案/验证码等轻量内容
- 地址省市区、支付方式、评价草稿等场景必须明确“哪些状态允许重建，哪些状态必须保留”

## Uni 多语言传输策略（单语优先 + 按需全量回退 + 页面内实时重算）（强制执行）
1. 单语优先（默认）
- Uni 请求默认按当前语言返回，避免多语言 payload 全量下发导致响应体积膨胀
- 列表/详情类接口默认不传 includeI18n，优先走单语响应
- 不允许为“图省事”全局开启 includeI18n

2. 按需全量回退（仅必要页面）
- 页面在不离页的情况下切换语言且需保留状态时，接口调用必须显式传 includeI18n
- 推荐在 Uni API 封装层统一透传 options.includeI18n，避免页面直接拼 query
- 典型需要全量回退的配置：payment methods、tryon_models、announcement_content、invite_share_link_tip_text、tryon_recharge_plans

3. 页面内实时重算（避免旧语种缓存）
- 禁止在加载时将 localText(...) 结果直接写死到长期状态（ref/store 字符串）
- 应保存原始值 raw（对象/JSON/字符串），通过 computed 或 watch(locale) 实时解析展示
- 对拼接文案（如订单摘要、地址省市区映射）必须提供 rebuild 方法，并在 locale 变化时重算

4. 页面缓存返回策略（tabBar / 复用页）
- 对 tabBar、navigateBack 可复用页、keep-alive 场景，必须在 onShow 增加 locale 变化守卫（推荐 lastLoadedLocale 对比）
- 禁止“已有数据即 return”短路导致旧语种残留；必须加入 localeChanged 例外分支
- 语言变化时优先做增量刷新（重拉 locale-sensitive 数据或重算派生文案），避免整页重置
- 长连接页面（如客服聊天）切语后应优先回刷本地可翻译文案，避免不必要的重连与状态抖动

5. 页面重载策略
- 切换语言后强制重载页面可以作为兜底，但不能作为唯一方案
- 长流程页面（支付、试衣生成、表单编辑）优先使用“全量回退 + 页面内重算”，避免重载导致状态丢失
- 若某页面明确采用重载方案，需在需求说明中标注“会丢失页面瞬时状态”

6. 交付验收标准
- 关键文案切换语言后无需离页即可正确更新（不出现旧语种残留）
- 接口抓包可验证：默认单语，只有必要请求携带 includeI18n
- 相关页面改造后执行 Uni 构建与核心路径手测（至少支付/试衣/邀请/地址）

## 提交前检查清单
- 新增 API 是否同步注册到 sys_apis
- 新增 API 是否有对应 Casbin 规则
- 是否存在 authority=8080 越权路径
- 写接口是否记录操作日志
- uni 与 web 是否跟随后端单一配置源
- 是否存在中文硬编码未进入 i18n
- uni.showToast / uni.showModal / JS 文案是否仍有硬编码
- Uni 是否仍存在 data.msg ||、res.msg ||、e?.message || 等 raw msg 直透
- 后端 API 是否仍存在 err.Error() 直接返回给 C 端用户
- 后端新增 key 是否在 Uni 所有启用语种中补齐
- 试衣间公告 announcement_content 是否按 localText 解析并随语言切换生效
- Uni 接口是否保持“默认单语”，且仅在必要场景显式 includeI18n
- 页面是否存在“localText 后写死状态”导致切语不更新（应改为 raw + computed/watch 重算）
- tabBar/复用页面是否有 onShow + lastLoadedLocale（或等价）语言守卫，避免返回后旧语种残留
- 多语言 JSON 的 fallback 顺序是否为 locale -> en -> zh（避免缺词时误回中文）
- 管理端后端多语言字段是否统一使用 MultiLangEditor（含一键翻译入口），无手写散落多语输入
- 富文本多语言字段是否通过 MultiLangEditor + RichEdit 接入，避免单语富文本或混搭输入
- 本地系统提示消息是否保存 i18nKey 并支持语言切换后回刷
- 切语重拉是否存在旧请求覆盖新语言数据风险（需有竞态保护）
- 表单编辑页是否避免“切语即重拉覆盖草稿”
- 修改 initialize/gorm*.go（含 gorm_biz.go）后，是否清理未使用 model 导入，且导入模型与 AutoMigrate 列表保持一致
- 修改初始化或迁移文件后，是否至少执行 go test ./initialize 做快速编译验证，防止 unused import 回归
- 新增 Private 路由后，是否同步更新 initialize 下对应的 SysApi 注册与 Casbin 路径清单（例如 invite_init/new_modules_init/shop_init），防止出现“权限不足”
