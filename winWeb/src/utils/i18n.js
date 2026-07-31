import { createI18n } from 'vue-i18n'
import zh from './i18n-locales/zh'
import en from './i18n-locales/en'
import storage from './storage'

const i18n = createI18n({
  legacy: false,
  locale: storage.get('app-lang') || 'zh',
  fallbackLocale: 'zh',
  messages: {
    zh,
    en
  }
})

export default i18n

/**
 * 解析多语言字段值（后端返回的 i18n 对象或 JSON 字符串）
 * @param {string|object} value - 字段值
 * @param {string} [lang] - 语言代码
 * @returns {string}
 */
export function localText(value, lang) {
  if (!value) return ''
  const locale = lang || storage.get('app-lang') || 'zh'
  if (typeof value === 'object') {
    return value[locale] || value['en'] || value['zh'] || Object.values(value)[0] || ''
  }
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
 * 统一解析后端或异常消息
 */
export function resolveApiMessage(value, fallbackKey = 'common.fail', lang) {
  const locale = lang || storage.get('app-lang') || 'zh'
  const messages = i18n.global.messages.value || i18n.global.messages
  const fallback = fallbackKey ? (messages[locale]?.[fallbackKey] || messages.zh?.[fallbackKey] || fallbackKey) : ''

  if (value === undefined || value === null) return fallback

  if (typeof value === 'object') {
    const text = String(localText(value, locale) || '').trim()
    if (!text) return fallback
    const translated = messages[locale]?.[text] || messages.zh?.[text]
    return translated !== text ? translated : text
  }

  const raw = String(value || '').trim()
  if (!raw) return fallback

  const normalized = raw.charAt(0) === '{'
    ? String(localText(raw, locale) || '').trim()
    : raw
  if (!normalized) return fallback

  const compact = normalized.replace(/\s+/g, '')
  const aliasKey = API_MESSAGE_ALIAS_MAP[normalized] || API_MESSAGE_ALIAS_MAP[compact]
  if (aliasKey) {
    return messages[locale]?.[aliasKey] || messages.zh?.[aliasKey] || aliasKey
  }

  const translated = messages[locale]?.[normalized] || messages.zh?.[normalized]
  return translated !== normalized ? translated : normalized
}

/**
 * 简易翻译函数（用于非组件环境）
 */
export function t(key, lang) {
  const locale = lang || storage.get('app-lang') || 'zh'
  const keys = key.split('.')
  let val = i18n.global.messages[locale] || i18n.global.messages.value?.[locale]
  for (const k of keys) {
    if (val && typeof val === 'object') val = val[k]
    else return key
  }
  return val || key
}

/**
 * 多 key 回退翻译
 */
export function tt(...keys) {
  for (const key of keys) {
    const text = t(key)
    if (text && text !== key) return text
  }
  const last = keys[keys.length - 1]
  return typeof last === 'string' ? last : ''
}
