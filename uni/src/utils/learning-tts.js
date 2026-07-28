/**
 * 学习模块 TTS 语音播放公共工具
 * 优先级流程：本地离线TTS -> 数据库音频 -> HTTPS重试 -> 有道兜底 -> 本地TTS最终兜底
 *
 * 使用方式：
 * import { playWordAudio, playSentenceAudio, stopLocalTTS, destroyTTS } from '@/utils/learning-tts'
 *
 * const callbacks = {
 *   onStart: () => { isPlaying.value = true },
 *   onEnd: () => { isPlaying.value = false },
 *   onError: () => { isPlaying.value = false }
 * }
 * playWordAudio(wordText, audioUs, audioUk, 'US', callbacks)
 */

// ========== 内部状态 ==========
let audioCtx = null
let audioCandidates = []
let currentAudioAttempt = null
let currentCallbacks = {}
let currentAccent = 'US'
let currentFallbackText = ''
let currentTtsVoice = 'female'

// ========== 辅助函数 ==========

const normalizeAudioSource = (raw) => {
  const src = String(raw || '').trim()
  if (!src) {
    return ''
  }
  if (src.startsWith('//')) {
    return `https:${src}`
  }
  return src.replace(/ /g, '%20')
}

const isLocalAudioSource = (src) => {
  return /^https?:\/\/(127\.0\.0\.1|localhost)/i.test(src)
}

const isPrivateNetworkAudioSource = (src) => {
  return /^https?:\/\/(10\.|192\.168\.|172\.(1[6-9]|2\d|3[0-1])\.)/i.test(src)
}

const buildFallbackAudioSource = (text, accent) => {
  const source = String(text || '').trim()
  if (!source) {
    return ''
  }
  const type = accent === 'UK' ? '2' : '1'
  return `https://dict.youdao.com/dictvoice?audio=${encodeURIComponent(source)}&type=${type}&le=en`
}

// ========== 本地离线 TTS ==========

/**
 * 停止所有本地 TTS 播放
 */
export const stopLocalTTS = () => {
  // #ifdef H5
  try {
    if (typeof window !== 'undefined' && window.speechSynthesis) {
      window.speechSynthesis.cancel()
    }
  } catch (e) { /* ignore */ }
  // #endif
  // #ifdef APP-PLUS
  try {
    if (typeof plus !== 'undefined' && plus.speech) {
      plus.speech.stop()
    }
  } catch (e) { /* ignore */ }
  // #endif
}

/**
 * 尝试本地离线 TTS 播放（优先级最高）
 * @param {string} text - 要播放的文本
 * @param {string} accent - 口音 'US' 或 'UK'
 * @param {object} options - { onStart, onEnd, onError, ttsVoice }
 * @returns {boolean} 是否成功启动本地 TTS
 */
export const tryLocalTTS = (text, accent = 'US', options = {}) => {
  const content = String(text || '').trim()
  if (!content) return false

  const { onStart, onEnd, onError, ttsVoice } = options
  const voiceType = ttsVoice || currentTtsVoice

  // #ifdef APP-PLUS
  try {
    if (typeof plus !== 'undefined' && plus.speech) {
      plus.speech.stop()
      const lang = accent === 'UK' ? 'en-GB' : 'en-US'
      plus.speech.start({
        text: content,
        lang: lang,
        rate: 90,
        onstart: () => { if (onStart) onStart() },
        onend: () => { if (onEnd) onEnd() },
        onerror: () => { if (onError) onError() }
      })
      return true
    }
  } catch (e) { /* plus.speech 不可用 */ }
  // #endif

  // #ifdef H5
  try {
    if (typeof window !== 'undefined' && window.speechSynthesis) {
      window.speechSynthesis.cancel()
      const utterance = new window.SpeechSynthesisUtterance(content)

      const langPrefix = accent === 'UK' ? 'en-GB' : 'en-US'

      // 尝试匹配对应音色的语音
      const voices = window.speechSynthesis.getVoices()
      let bestVoice = null

      if (voices.length > 0) {
        // 优先级逐级匹配：语言+性别 -> 语言 -> 性别兜底
        for (const v of voices) {
          const vLang = (v.lang || '').toLowerCase()
          const vName = (v.name || '').toLowerCase()
          const targetLang = langPrefix.toLowerCase()

          if (voiceType === 'female') {
            if (vLang.startsWith(targetLang) &&
                (vName.includes('female') || vName.includes('samantha') || vName.includes('zira') ||
                 vName.includes('susan') || vName.includes('karen') || vName.includes('monica'))) {
              bestVoice = v
              break
            }
          } else {
            if (vLang.startsWith(targetLang) &&
                (vName.includes('male') || vName.includes('david') || vName.includes('mark') ||
                 vName.includes('daniel') || vName.includes('tom'))) {
              bestVoice = v
              break
            }
          }
        }
        // 降级：匹配对应语言即可
        if (!bestVoice) {
          for (const v of voices) {
            if (v.lang.toLowerCase().startsWith(langPrefix.toLowerCase())) {
              bestVoice = v
              break
            }
          }
        }
        // 最终降级：任意英文语音
        if (!bestVoice) {
          for (const v of voices) {
            if (v.lang.toLowerCase().startsWith('en-')) {
              bestVoice = v
              break
            }
          }
        }
      }

      if (bestVoice) {
        utterance.voice = bestVoice
      }
      utterance.lang = langPrefix
      utterance.rate = 0.9
      utterance.onstart = () => { if (onStart) onStart() }
      utterance.onend = () => { if (onEnd) onEnd() }
      utterance.onerror = () => { if (onError) onError() }
      window.speechSynthesis.speak(utterance)
      return true
    }
  } catch (e) {
    // SpeechSynthesis 不可用
  }
  // #endif

  return false
}

// ========== 在线音频播放（兜底流程） ==========

const playNextAudioCandidate = () => {
  while (audioCandidates.length > 0) {
    const next = audioCandidates.shift()
    if (!next || !next.url) {
      continue
    }
    currentAudioAttempt = next
    audioCtx.src = next.url
    audioCtx.play()
    return
  }
  currentAudioAttempt = null
  if (currentCallbacks.onError) {
    currentCallbacks.onError()
  }
}

const ensureAudioCtx = () => {
  if (audioCtx) return audioCtx
  audioCtx = uni.createInnerAudioContext()
  audioCtx.onPlay(() => {
    if (currentCallbacks.onStart) currentCallbacks.onStart()
  })
  audioCtx.onEnded(() => {
    if (currentCallbacks.onEnd) currentCallbacks.onEnd()
  })
  audioCtx.onError(() => {
    // 所有在线源失败后尝试本地TTS作为最后兜底
    if (audioCandidates.length === 0 && currentFallbackText) {
      const spoke = tryLocalTTS(currentFallbackText, currentAccent, currentCallbacks)
      if (spoke) {
        currentFallbackText = ''
        return
      }
    }
    playNextAudioCandidate()
  })
  return audioCtx
}

/**
 * 在线音频播放（兜底流程）
 * @param {string} src - 主音频 URL
 * @param {string} fallbackText - 兜底文本（用于有道TTS和本地TTS）
 * @param {string} accent - 口音 'US' 或 'UK'
 * @param {object} callbacks - { onStart, onEnd, onError }
 */
export const playOnlineAudio = (src, fallbackText = '', accent = 'US', callbacks = {}) => {
  const audioSrc = normalizeAudioSource(src)
  const fallbackAudio = buildFallbackAudioSource(fallbackText, accent)
  currentFallbackText = fallbackText || ''
  currentCallbacks = callbacks
  currentAccent = accent

  const isUnreachable = audioSrc && (isLocalAudioSource(audioSrc) || isPrivateNetworkAudioSource(audioSrc))
  const primary = audioSrc && !isUnreachable ? audioSrc : ''

  audioCandidates = []
  if (primary) {
    audioCandidates.push({ url: primary, source: 'db-primary' })
    if (primary.startsWith('http://')) {
      audioCandidates.push({ url: primary.replace(/^http:\/\//, 'https://'), source: 'db-https-retry' })
    }
  }
  if (fallbackAudio && !audioCandidates.some((item) => item.url === fallbackAudio)) {
    audioCandidates.push({ url: fallbackAudio, source: 'free-fallback' })
  }

  if (audioCandidates.length === 0) {
    if (callbacks.onError) callbacks.onError()
    return
  }

  ensureAudioCtx()
  playNextAudioCandidate()
}

// ========== 公共 API ==========

/**
 * 播放单词/句子音频（完整优先级流程）
 * 优先本地TTS，失败后兜底在线音频
 *
 * @param {string} text - 要播放的文本
 * @param {string} audioUs - 美式音频 URL
 * @param {string} audioUk - 英式音频 URL
 * @param {string} accent - 口音 'US' 或 'UK'
 * @param {object} callbacks - { onStart, onEnd, onError, ttsVoice }
 */
export const playWordAudio = (text, audioUs, audioUk, accent = 'US', callbacks = {}) => {
  const primary = accent === 'US' ? audioUs : audioUk
  const fallback = accent === 'US' ? audioUk : audioUs

  // 默认使用本地离线TTS
  const localOk = tryLocalTTS(text, accent, callbacks)
  if (localOk) return

  // 兜底：在线音频
  playOnlineAudio(primary || fallback, text, accent, callbacks)
}

/**
 * 播放句子音频（与 playWordAudio 相同的优先级流程）
 * 兼容 playWordAudio 的接口
 */
export const playSentenceAudio = playWordAudio

/**
 * 销毁 TTS 资源（页面卸载时调用）
 */
export const destroyTTS = () => {
  stopLocalTTS()
  if (audioCtx) {
    audioCtx.destroy()
    audioCtx = null
  }
  audioCandidates = []
  currentAudioAttempt = null
  currentFallbackText = ''
  currentCallbacks = {}
}
