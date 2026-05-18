---
name: uni-i18n-display-unification
description: "Uni 多语言展示统一工作流。用于修复 SKU 规格在切语时显示中文/key、统一 nameI18n/labelI18n/valueI18n 解析方式，并标准化 includeI18n 与 useI18nDisplay 的落地流程。"
---

# Uni 多语言展示统一工作流

## 适用场景
- SKU 规格在英文环境仍显示中文
- 页面显示 `labelI18n: "color"`、`valueI18n: "red"` 但未翻译
- 同类页面各自写 `resolveDisplayText`，逻辑分叉
- 页面切换语言后，配置文案不更新或显示旧语种

## 目标
- 保持 Uni 端默认单语返回，按需全量回退
- 所有展示文案统一走 `useI18nDisplay`
- 统一兼容 4 种输入：对象 / JSON 字符串 / i18n key / 纯文本

## 配套资产
- 手工回归清单：`REGRESSION-CHECKLIST.md`
- 建议在每次批量改造后按清单执行中英文切换验证，并记录通过截图或日志。

## 标准调用链（必须理解）
1. Uni 请求统一设置请求头
- `Accept-Language`: 当前语言
- `X-Client-App: uni`
- 代码位置：`uni/src/utils/request.js`

2. 后端中间件解析语言
- `middleware.Locale()` 读取 `Accept-Language` 写入上下文
- 代码位置：`server/middleware/locale.go`

3. API 返回前统一过 `LocalizeResponseData`
- 代码位置：`server/utils/i18n/payload_localizer.go`
- 规则：
  - Uni 默认单语（`X-Client-App=uni`）
  - `includeI18n=1` 或 `i18nMode=full` 返回全量多语

4. 前端展示统一解析
- composable: `uni/src/composables/useI18nDisplay.js`
- 核心函数: `resolveI18nDisplayText` in `uni/src/utils/i18n.js`

## 决策规则（默认单语 + 按需全量）
1. 默认不传 `includeI18n`
- 列表/详情页只需当前语言展示即可
- 减少 payload，响应更轻

2. 必须传 `includeI18n: true`
- 页面内不离页切语，需要即时重算文案
- 配置类接口需保留原始多语对象用于动态切语
- 典型接口：
  - `/sysConfig/getPaymentConfig`
  - `/sysConfig/getUniPreferredPayConfig`
  - `/sysConfig/getAnnouncementConfig`
  - `/sysConfig/getTryonConfig`

## 统一模板

### A. API 层模板（支持按需回退）
```js
const withI18nFallback = (params = {}, options = {}) => {
  if (options && options.includeI18n) {
    return { ...params, includeI18n: 1 }
  }
  return params
}

export const getPaymentConfig = (options = {}) => {
  return request({
    url: '/sysConfig/getPaymentConfig',
    method: 'get',
    params: withI18nFallback({}, options)
  })
}
```

### B. 页面层模板（统一展示解析）
```js
import { computed } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useI18nDisplay } from '@/composables/useI18nDisplay.js'

const langStore = useLangStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')
const { resolveDisplayText } = useI18nDisplay(locale)

const label = resolveDisplayText(spec?.labelI18n || spec?.nameI18n || spec?.label || spec?.name, spec?.label || spec?.name || '')
const value = resolveDisplayText(spec?.valueI18n || spec?.value, spec?.value || '')
```

### C. SKU 组件模板（保持命名语义）
```js
const { resolveDisplayText: resolveSpecText } = useI18nDisplay(locale)
```

## 改造步骤（可直接执行）
1. 扫描风险点
- `resolveI18nDisplayText` 直调
- `$lt(nameI18n/labelI18n/valueI18n)` 直出
- 页面内自定义 `resolveDisplayText` 包装

2. 改造策略
- 优先替换为 `useI18nDisplay(locale)`
- 仅替换解析入口，不改业务流程
- 组件内保持原函数命名（如 `resolveSpecText`）可读性更高

3. 校验
- 静态检查：目标文件无报错
- 构建校验：`npm -C uni run build:h5`
- 手工回归：中文 -> 英文 -> 中文，页面不离开时文案应即时更新

## 常见错误（禁止）
- 用 `String(value)` 渲染多语言对象/JSON
- 在页面里直接 `localText(...)` 后写死到长期状态，切语不更新
- 在需要页面内动态切语的场景不传 `includeI18n`
- 同一类字段在不同页面混用不同解析策略

## 验收清单
- `nameI18n/labelI18n/valueI18n` 展示统一走 `resolveDisplayText`
- 页面中无 `resolveI18nDisplayText(...)` 直调（除 `useI18nDisplay` 与 `utils/i18n.js`）
- 关键页面切语不离页可即时更新文案
- Uni 构建通过（允许既有 warning，但不能有新 error）
