/**
 * 多语言国际化工具
 * 支持：中文(zh)、英语(en)、蒙古国语(mn)、繁體中文(zh-TW)、泰语(th)、印地语(hi)、印尼语(id)、越南语(vi)、阿拉伯语(ar)、日语(ja)、韩语(ko)、马来语(ms)
 */

import messages from './i18n-locales/messages.js'

/**
 * 获取翻译文本
 * @param {string} key - 翻译键
 * @param {string} lang - 语言代码
 * @returns {string}
 */
export function t(key, lang) {
  const locale = lang || uni.getStorageSync('app-lang') || 'zh'
  return (messages[locale] && messages[locale][key]) || messages.zh[key] || key
}

/**
 * 解析多语言字段值
 * 支持两种格式：
 * 1. 普通字符串 → 原样返回
 * 2. JSON对象 {"zh":"中文","en":"English","mn":"Монгол"} → 按当前语言返回
 * @param {string|object} value - 字段值（可能是字符串或JSON对象）
 * @param {string} [lang] - 语言代码，不传则自动获取当前语言
 * @returns {string}
 */
export function localText(value, lang) {
  if (!value) return ''
  const locale = lang || uni.getStorageSync('app-lang') || 'zh'
  // 已经是对象
  if (typeof value === 'object') {
    return value[locale] || value['en'] || value['zh'] || Object.values(value)[0] || ''
  }
  // 字符串，尝试解析 JSON
  if (typeof value === 'string') {
    if (value.charAt(0) === '{') {
      try {
        const obj = JSON.parse(value)
        return obj[locale] || obj['en'] || obj['zh'] || Object.values(obj)[0] || value
      } catch (e) {
        return value
      }
    }
    return value
  }
  return String(value)
}

/**
 * 规范化候选文案：兼容对象、JSON字符串、key字符串。
 * @param {string|object|undefined|null} value
 * @returns {string|object|null}
 */
export function normalizeI18nCandidate(value) {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    return value
  }

  if (typeof value === 'string') {
    const text = value.trim()
    if (!text) return null
    if (text.charAt(0) === '{' && text.charAt(text.length - 1) === '}') {
      try {
        const parsed = JSON.parse(text)
        if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
          return parsed
        }
      } catch {
        return text
      }
    }
    return text
  }

  return null
}

/**
 * 统一展示文案解析：
 * 1. 先按 localText 解析对象/JSON
 * 2. 字符串再尝试按 i18n key 翻译
 * 3. key 未命中时，英文保留 key，其他语种优先回退 fallback
 * @param {string|object|undefined|null} candidate
 * @param {string|object|undefined|null} fallback
 * @param {string} [lang]
 * @returns {string}
 */
export function resolveI18nDisplayText(candidate, fallback = '', lang) {
  const locale = lang || uni.getStorageSync('app-lang') || 'zh'

  const localized = String(localText(candidate, locale) || '').trim()
  if (localized && !(typeof candidate === 'string' && localized === candidate.trim())) {
    return localized
  }

  const normalizedCandidate = typeof candidate === 'string' ? candidate.trim() : ''
  const fallbackText = String(localText(fallback, locale) || '').trim()

  if (normalizedCandidate) {
    const byKey = String(t(normalizedCandidate, locale) || '').trim()
    if (byKey && byKey !== normalizedCandidate) {
      return byKey
    }

    if (!String(locale || '').toLowerCase().startsWith('en') && fallbackText) {
      return fallbackText
    }

    if (localized) {
      return localized
    }
  }

  return fallbackText || String(fallback || '')
}

const API_MESSAGE_ALIAS_MAP = {
  '上传文件失败': 'uploadFail',
  '接收文件失败': 'uploadFail',
  '权限不足': 'noPermission',
  '无权限操作': 'noPermission',
  '未登录': 'notLogin',
  '未登录或非法访问，请登录': 'notLogin',
  'token 过期': 'tokenExpired',
  'token已过期': 'tokenExpired',
  '您的帐户异地登陆或令牌失效': 'tokenInvalid',
}

/**
 * 统一解析后端或异常消息：支持 i18n key、多语言对象、JSON 字符串和普通文本。
 * @param {string|object|undefined|null} value - 原始消息
 * @param {string} [fallbackKey='operationFailed'] - 回退词条 key
 * @param {string} [lang] - 语言代码
 * @returns {string}
 */
export function resolveApiMessage(value, fallbackKey = 'operationFailed', lang) {
  const locale = lang || uni.getStorageSync('app-lang') || 'zh'
  const fallback = fallbackKey ? t(fallbackKey, locale) : ''

  if (value === undefined || value === null) {
    return fallback
  }

  if (typeof value === 'object') {
    const text = String(localText(value, locale) || '').trim()
    if (!text) return fallback
    const translated = t(text, locale)
    return translated !== text ? translated : text
  }

  const raw = String(value || '').trim()
  if (!raw) {
    return fallback
  }

  const normalized = raw.charAt(0) === '{'
    ? String(localText(raw, locale) || '').trim()
    : raw
  if (!normalized) {
    return fallback
  }

  const compact = normalized.replace(/\s+/g, '')
  const aliasKey = API_MESSAGE_ALIAS_MAP[normalized] || API_MESSAGE_ALIAS_MAP[compact]
  if (aliasKey) {
    return t(aliasKey, locale)
  }

  const translated = t(normalized, locale)
  return translated !== normalized ? translated : normalized
}

export default messages
