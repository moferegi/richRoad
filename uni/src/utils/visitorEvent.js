import { visitorHeartbeat } from '@/api/visitor.js'
import { generateFingerprint, getSessionId, getPlatform } from '@/utils/fingerprint.js'

const DEFAULT_CATEGORY = 'payment_kefu_guide'

const safeString = (val, maxLen) => {
  const text = String(val || '').trim()
  if (!text) return ''
  return text.slice(0, maxLen)
}

const safeNumber = (val) => {
  const num = Number(val)
  return Number.isFinite(num) ? num : 0
}

const toExtraString = (extra, maxLen = 4000) => {
  if (!extra) return ''
  try {
    return JSON.stringify(extra).slice(0, maxLen)
  } catch (e) {
    return ''
  }
}

const getCurrentPagePath = () => {
  try {
    const pages = getCurrentPages()
    if (!pages || pages.length === 0) return '/'
    return '/' + String(pages[pages.length - 1].route || '').replace(/^\//, '')
  } catch (e) {
    return '/'
  }
}

export const trackVisitorEvent = (payload = {}) => {
  const action = safeString(payload.action, 64)
  if (!action) return Promise.resolve()

  let info = {}
  try {
    info = uni.getSystemInfoSync() || {}
  } catch (e) {
    info = {}
  }

  return visitorHeartbeat({
    visitorId: generateFingerprint(),
    sessionId: getSessionId(),
    platform: getPlatform(),
    pagePath: safeString(payload.pagePath, 255) || getCurrentPagePath(),
    referer: safeString(payload.referer, 512),
    screenWidth: safeNumber(info.screenWidth),
    screenHeight: safeNumber(info.screenHeight),
    language: safeString(uni.getStorageSync('app-lang') || info.language || 'zh', 16),
    eventCategory: safeString(payload.category || DEFAULT_CATEGORY, 64),
    eventAction: action,
    eventLabel: safeString(payload.label, 128),
    eventValue: safeNumber(payload.value),
    eventExtra: toExtraString(payload.extra)
  }).catch(() => {})
}
