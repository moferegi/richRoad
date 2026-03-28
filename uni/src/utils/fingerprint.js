/**
 * 生成浏览器指纹ID
 * 基于屏幕信息、时区、平台等不变特征生成哈希
 */
function generateFingerprint() {
  const cached = uni.getStorageSync('visitor_id')
  if (cached) return cached

  const info = uni.getSystemInfoSync()
  const raw = [
    info.screenWidth,
    info.screenHeight,
    info.platform,
    info.system,
    info.language || 'unknown',
    info.model || 'unknown',
    new Date().getTimezoneOffset()
  ].join('|')

  // 简单哈希
  let hash = 0
  for (let i = 0; i < raw.length; i++) {
    const char = raw.charCodeAt(i)
    hash = ((hash << 5) - hash) + char
    hash = hash & hash // 转32位整数
  }
  const visitorId = 'v_' + Math.abs(hash).toString(36) + '_' + Date.now().toString(36)

  uni.setStorageSync('visitor_id', visitorId)
  return visitorId
}

/**
 * 获取会话ID（每次打开APP/页面重新生成）
 */
function getSessionId() {
  let sessionId = uni.getStorageSync('visitor_session_id')
  if (!sessionId) {
    sessionId = 's_' + Date.now().toString(36) + '_' + Math.random().toString(36).slice(2, 8)
    uni.setStorageSync('visitor_session_id', sessionId)
  }
  return sessionId
}

/**
 * 获取平台信息
 */
function getPlatform() {
  let platform = 'h5'
  // #ifdef APP-PLUS
  platform = uni.getSystemInfoSync().platform // ios / android
  // #endif
  // #ifdef MP-WEIXIN
  platform = 'wxmp'
  // #endif
  return platform
}

export { generateFingerprint, getSessionId, getPlatform }
