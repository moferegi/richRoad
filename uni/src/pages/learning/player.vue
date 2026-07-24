<template>
  <view class="player-container">
    <!-- 顶部视频区 (固定不动) -->
    <view class="video-section">
      <video 
        id="englishVideo" 
        class="main-video" 
        :src="videoSrc" 
        :autoplay="true"
        :show-center-play-btn="!showVideoLoading"
        :loop="controls.loop"
        @timeupdate="onTimeUpdate"
        @loadedmetadata="onVideoLoadedMeta"
        @canplay="onVideoCanPlay"
        @seeking="onVideoSeeking"
        @seeked="onVideoSeeked"
        @waiting="onVideoWaiting"
        @playing="onVideoPlaying"
        @play="onVideoPlay"
        @pause="onVideoPause"
        @error="onVideoError"
        @ended="onVideoEnded"
      ></video>
      <view v-if="showVideoLoading" class="video-loading-mask">
        <image v-if="loadingLogo" class="video-loading-logo spin" :src="loadingLogo" mode="aspectFit" />
        <view v-else class="video-loading-logo-fallback spin">R</view>
        <text class="video-loading-text">{{ t('player.video_loading') }}</text>
      </view>
      <!-- 手势滑动暂停提示文字 -->
      <view v-if="showPauseHint" class="pause-hint-overlay" @click="hidePauseHint">
        <text class="pause-hint-text">{{ t('player.tap_subtitle_to_play') }}</text>
      </view>
    </view>

    <!-- 视频进度条 -->
    <view class="progress-section" v-if="videoInfo.duration > 0">
      <view class="progress-time-row">
        <text class="progress-current-time">{{ formatSeconds(currentTime) }}</text>
        <text class="progress-total-time">{{ formatSeconds(videoInfo.duration) }}</text>
      </view>
      <view class="progress-row">
        <slider 
          class="video-progress-slider" 
          :min="0" 
          :max="videoInfo.duration || 100" 
          :value="currentTime"
          :step="0.5"
          activeColor="#6D5BFF"
          backgroundColor="rgba(148,163,184,0.24)"
          block-size="20"
          @change="onProgressChange"
          @changing="onProgressChanging"
        />
      </view>
    </view>

    <!-- 中央外挂字幕滚动区 -->
    <scroll-view 
      class="subtitle-section" 
      scroll-y 
      scroll-with-animation
      :scroll-into-view="activeSentenceId"
      v-if="controls.showSubtitles"
      @touchstart="onSubtitleTouchStart"
      @touchmove="onSubtitleTouchMove"
    >
      <view class="subtitle-padding-top"></view>
      
      <view 
        v-for="(item, index) in subtitleList" 
        :key="index"
        :id="'sentence_' + index"
        class="sentence-row"
        :class="{ 'is-active': activeSentenceIndex === index, 'is-focused': focusedSentenceId === Number(item.id) }"
        @click="jumpBySubtitle(item, index)"
      >
        <view class="sentence-meta">
          <text class="time-badge">{{ formatSeconds(item.startTime) }} - {{ formatSeconds(item.endTime) }}</text>
          <text class="collect-btn" @click.stop="collectSentence(item.id)">⭐ {{ sentenceCollectedMap[Number(item.id)] ? t('typing.uncollect') : t('player.collect') }}</text>
        </view>
        
        <!-- 英文句子 (带高亮重点词可点击分析) -->
        <view class="english-text">
          <text 
            v-for="(seg, sIdx) in parseHighlight(item.english)" 
            :key="sIdx"
            :class="{ 'highlight-word': seg.isHighlight }"
            @click="handleSegClick(seg, $event)"
          >
            {{ seg.text }}
          </text>
        </view>

        <!-- 对应多语翻译 -->
        <view class="translated-text" v-if="controls.showBilingual">
          {{ localText(item.translate) }}
        </view>
      </view>
      
      <view class="subtitle-padding-bottom"></view>
    </scroll-view>

    <!-- 底部控制栏 -->
    <view class="bottom-controls">
      <view class="control-item" @click="toggleBilingual">
        <view class="c-icon-circle" :class="{ 'c-circle-on': controls.showBilingual }">
          <text class="c-icon">🌐</text>
        </view>
        <text class="c-text" :class="{ 'c-text-active': controls.showBilingual }">{{ t('player.bilingual') }}</text>
      </view>
      
      <view class="control-item" @click="toggleRepeat">
        <view class="c-icon-circle" :class="{ 'c-circle-on': repeatMode }">
          <text class="c-icon">🔄</text>
        </view>
        <text class="c-text" :class="{ 'c-text-active': repeatMode }">{{ t('player.repeat') }}</text>
      </view>
      
      <view class="control-item play-btn" @click="togglePlay">
        <view class="c-play-circle">
          <text class="c-icon c-play-icon">{{ showVideoLoading ? '⏳' : (isPlaying ? '⏸' : '▶') }}</text>
        </view>
      </view>
      
      <view class="control-item" @click="toggleSubtitles">
        <view class="c-icon-circle" :class="{ 'c-circle-on': controls.showSubtitles }">
          <text class="c-icon">💬</text>
        </view>
        <text class="c-text" :class="{ 'c-text-active': controls.showSubtitles }">{{ t('player.subtitle') }}</text>
      </view>

      <view class="control-item" @click="openMoreSheet">
        <view class="c-icon-circle" :class="{ 'c-circle-on': controls.loop }">
          <text class="c-icon">⚙</text>
        </view>
        <text class="c-text">{{ tt('player.more', 'common.more') }}</text>
      </view>
    </view>

    <!-- 重点词汇查词弹窗 -->
    <uni-popup ref="wordPopup" type="bottom" background-color="#fff">
      <view class="word-detail-box" v-if="currentWord">
        <view class="popup-handle"></view>
        <view class="w-header">
          <text class="w-title">{{ currentWord.word }}</text>
          <view class="w-audio-circle" @click="playWordAudio(currentWord.audioUs)">
            <text class="w-audio-icon">🔊</text>
          </view>
        </view>
        <text class="w-phonetic">{{ currentWord.phoneticUs }}</text>
        <view class="w-divider"></view>
        <view class="w-exp">{{ localText(currentWord.explanation) }}</view>
        <button class="w-collect" @click="collectWord(currentWord.id)">⭐ {{ t('typing.collect') }}</button>
      </view>
    </uni-popup>

    <uni-popup ref="morePopup" type="bottom" background-color="#fff">
      <view class="more-sheet">
        <view class="popup-handle"></view>
        <view class="more-sheet-head">
          <text class="more-sheet-title">{{ tt('player.more_settings', 'player.more') }}</text>
          <text class="more-sheet-subtitle">{{ tt('player.more_settings_desc', 'player.loop_desc') }}</text>
        </view>

        <view class="more-card" @click="toggleLoop">
          <view class="more-card-main">
            <text class="more-card-title">{{ t('player.loop') }}</text>
            <text class="more-card-desc">{{ tt('player.loop_desc', 'player.loop') }}</text>
          </view>
          <view class="more-state-pill" :class="{ on: controls.loop }">
            {{ controls.loop ? tt('player.state_on') : tt('player.state_off') }}
          </view>
        </view>

        <!-- 播放速率选择 -->
        <view class="more-card speed-card">
          <view class="more-card-main">
            <text class="more-card-title">{{ t('player.speed') }}</text>
            <text class="more-card-desc">{{ tt('player.speed_desc', 'player.speed') }}</text>
          </view>
        </view>
        <view class="speed-pills">
          <text 
            v-for="rate in speedRates" 
            :key="rate"
            class="speed-pill"
            :class="{ 'speed-pill-active': controls.speed === rate }"
            @click="selectSpeed(rate)"
          >{{ rate }}x</text>
        </view>

        <view class="more-actions">
          <button class="more-action-btn ghost" @click="handleBack">{{ tt('player.exit_page', 'player.sheet_back', 'common.back') }}</button>
          <button class="more-action-btn" @click="closeMoreSheet">{{ tt('player.sheet_cancel', 'common.cancel', 'cancel') }}</button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { getExternalUrl } from '@/utils/url.js'
import { collect, findVideoEpisode, findWord, getCollectionList, getSentenceList, getWatchProgress, heartbeat, uncollect } from '@/api/learning.js'
import { baseUrl } from '@/utils/request.js'
// #ifdef H5
import Hls from 'hls.js'
// #endif

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()

const t = (k) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(k, locale)
  if (text && text !== k) return text
  return k
}

const tt = (...keys) => {
  for (const key of keys) {
    const text = t(key)
    if (text && text !== key) {
      return text
    }
  }
  const last = keys[keys.length - 1]
  return typeof last === 'string' ? last : ''
}

const localText = (val) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  return i18nLocalText(val, locale)
}
const formatSeconds = (secs) => {
  const m = Math.floor(secs / 60).toString().padStart(2, '0')
  const s = Math.floor(secs % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}

let videoCtx = null

// hls.js 实例（仅 H5 端，用于 .m3u8 播放；mp4 走原生 video）
// #ifdef H5
let hlsInstance = null
let nativeVideoEl = null
let hlsFatalErrorCount = 0
const HLS_MAX_FATAL_ERRORS = 5
// #endif

const isHlsUrl = (url) => /\.m3u8(\?|$)/i.test(url || '')

const videoSrc = computed(() => {
  const url = videoInfo.value.videoUrl || ''
  // m3u8 也设 src：Safari（iOS/macOS）原生 HLS 需要它；Chrome 端 hls.js 会通过 MediaSource 接管
  return url
})
let audioCtx = null
const isPlaying = ref(false)
const isVideoReady = ref(false)
const isTrialLocked = ref(false)
const lastTrialModalAt = ref(0)
const episodeId = ref(0)
const sentenceId = ref(0)

// 鉴权限制体系
const authData = ref({
  hasFullAuth: false,
  trialPercent: 8,    // 只能试看百分之几
})

const videoInfo = ref({
  id: 0,
  name: '',
  videoUrl: '',
  duration: 0,
})

const controls = ref({
  showBilingual: true,
  showSubtitles: true,
  speed: 1.0,
  loop: false
})

const subtitleList = ref([])

const activeSentenceIndex = ref(-1)
const activeSentenceId = ref('')
const focusedSentenceId = ref(0)
const pendingSeekTime = ref(0)
const isRefreshingSignedUrl = ref(false)
const autoPlayPending = ref(true)
const isVideoBuffering = ref(true)
const lastPlaybackTime = ref(0)
const loadingLogo = computed(() => getExternalUrl(appConfigStore.appLogo || ''))
const showVideoLoading = computed(() => isVideoBuffering.value || isRefreshingSignedUrl.value)
const sentenceCollectedMap = ref({})
let lastSignedRefreshAt = 0
let focusSentenceTimer = null
const morePopup = ref(null)
const currentTime = ref(0)
const showPauseHint = ref(false)
const repeatMode = ref(false)
const speedRates = [0.75, 1.0, 1.25, 1.5, 2.0]
let isProgressDragging = false
const showSpeedPicker = ref(false)
let lastRepeatSeekTime = 0
let loadingTimeoutTimer = null

const clearLoadingTimeout = () => {
  if (loadingTimeoutTimer) {
    clearTimeout(loadingTimeoutTimer)
    loadingTimeoutTimer = null
  }
}

const startLoadingTimeout = () => {
  clearLoadingTimeout()
  loadingTimeoutTimer = setTimeout(() => {
    if (isVideoBuffering.value && !isVideoReady.value) {
      console.warn('[player] 加载超时（8秒），清除等待状态')
      isVideoBuffering.value = false
    }
  }, 8000)
}

onMounted(() => {
  videoCtx = uni.createVideoContext('englishVideo')
  // #ifdef H5
  nativeVideoEl = document.getElementById('englishVideo')
  // #endif
  appConfigStore.loadConfig({ force: true, localeOnly: true })
  // isVideoBuffering 由 initPage() 设置，此处不重复置 true（避免覆盖 canplay 事件的结果）
  isVideoReady.value = false
  if (pendingSeekTime.value > 0) {
    videoCtx.seek(pendingSeekTime.value)
    pendingSeekTime.value = 0
  }
  // #ifdef H5
  initHlsIfNeeded()
  // #endif
})

onUnmounted(() => {
  flushHeartbeat(lastHeartbeatTime, true)
  clearLoadingTimeout()
  // #ifdef H5
  destroyHls()
  // #endif
  if (focusSentenceTimer) {
    clearTimeout(focusSentenceTimer)
    focusSentenceTimer = null
  }
  if (audioCtx) {
    audioCtx.destroy()
    audioCtx = null
  }
})

onUnload(() => {
  flushHeartbeat(lastHeartbeatTime, true)
})

// #ifdef H5
// 监听 videoUrl 变化：m3u8 → hls.js 接管；mp4 → 原生播放
watch(() => videoInfo.value.videoUrl, (newUrl) => {
  nextTick(() => {
    initHlsIfNeeded()
  })
})
// #endif

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const seekVideo = (secs) => {
  const target = Math.max(0, Number(secs || 0))
  if (videoCtx) {
    videoCtx.seek(target)
  } else {
    pendingSeekTime.value = target
  }
}

// H5 浏览器自动播放策略限制：play() 需要在用户手势内调用，否则抛出 NotAllowedError
// 包装 safePlay 静默吞掉该异常，避免未捕获的 Promise rejection
const safePlay = () => {
  if (!videoCtx) return
  try {
    const p = videoCtx.play()
    if (p && typeof p.catch === 'function') {
      p.catch(() => {})
    }
  } catch (_) {}
}

// ========== hls.js 集成（H5 端 .m3u8 播放，mp4 不受影响） ==========
// #ifdef H5
let currentBlobUrl = null

// 拉取 m3u8 并重写 URI：密钥 → 后端地址，.ts 分段 → 绝对路径
const rewriteM3u8 = async (m3u8Url) => {
  const resp = await fetch(m3u8Url)
  if (!resp.ok) throw new Error(`Fetch m3u8 failed: ${resp.status}`)
  const content = await resp.text()

  // baseUrl 可能是相对路径（如 /api）或绝对 URL（如 https://enapi.235235.vip）
  // blob m3u8 中不能有相对 URI，必须转为绝对地址，否则 hls.js 解析出错
  let backendBase = baseUrl
  if (backendBase.startsWith('/')) {
    backendBase = window.location.origin + backendBase
  }
  const m3u8Base = m3u8Url.substring(0, m3u8Url.lastIndexOf('/') + 1)

  const lines = content.split('\n')
  const rewritten = lines.map(line => {
    // 替换密钥 URI：相对路径 → 后端绝对地址
    if (line.startsWith('#EXT-X-KEY')) {
      return line.replace(/URI="([^"]*)"/, (_m, uri) => {
        if (uri.startsWith('http://') || uri.startsWith('https://')) return _m
        return `URI="${backendBase}${uri.startsWith('/') ? uri : '/' + uri}"`
      })
    }
    // .ts 分段改用绝对路径，确保 blob URL 场景下也能正确加载
    if (line && !line.startsWith('#') && /\.ts(\?|$)/i.test(line)) {
      if (!line.startsWith('http://') && !line.startsWith('https://')) {
        return m3u8Base + line
      }
    }
    return line
  })

  const blob = new Blob([rewritten.join('\n')], { type: 'application/vnd.apple.mpegurl' })
  return URL.createObjectURL(blob)
}

const initHls = async (url) => {
  destroyHls()
  hlsFatalErrorCount = 0
  const wrapper = document.getElementById('englishVideo')
  const realVideo = wrapper ? (wrapper.querySelector('video') || wrapper.getElementsByTagName('video')[0]) : null
  nativeVideoEl = realVideo

  if (!nativeVideoEl || !url) {
    if (!realVideo && wrapper) {
      console.warn('[hls.js] 未找到 <uni-video> 内部的 <video> 元素，稍后重试')
      setTimeout(() => initHls(url), 200)
    }
    return
  }

  // 开始加载流程，启动超时保护
  isVideoBuffering.value = true
  isVideoReady.value = false
  startLoadingTimeout()

  // 拉取并重写 m3u8，动态替换密钥 URI 为当前环境后端地址
  let blobUrl = url
  try {
    blobUrl = await rewriteM3u8(url)
    currentBlobUrl = blobUrl
  } catch (e) {
    console.error('[hls.js] 重写 m3u8 失败，回退原始 URL', e)
  }

  if (Hls.isSupported()) {
    // Chrome / Edge / Firefox / Android：用 hls.js 软解
    nativeVideoEl.removeAttribute('src')
    hlsInstance = new Hls({
      // 禁用 hls.js 内部重试，由外部统一控制（避免非 fatal 错误绕开我们的计数器）
      fragLoadingMaxRetry: 0,
      manifestLoadingMaxRetry: 0,
      levelLoadingMaxRetry: 0
    })
    hlsInstance.attachMedia(nativeVideoEl)
    hlsInstance.on(Hls.Events.MEDIA_ATTACHED, () => {
      hlsInstance.loadSource(blobUrl)
    })
    // m3u8 解析完成即说明视频源可达，提前结束加载状态
    hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
      clearLoadingTimeout()
      isVideoBuffering.value = false
      isVideoReady.value = true
      // 注意：不在此处重置 hlsFatalErrorCount
      // MANIFEST_PARSED 会在 startLoad() 恢复时再次触发，重置会导致死循环
    })
    hlsInstance.on(Hls.Events.ERROR, (_event, data) => {
      console.error('[hls.js] 播放错误', data.type, data.details, data.fatal)

      // keyLoadError / fragDecryptError：密钥或解密问题，不管 fatal 与否都立即停止
      if (data.details === 'keyLoadError' || data.details === 'fragDecryptError') {
        console.error('[hls.js] 密钥/解密失败，停止播放')
        hlsInstance.stopLoad()
        destroyHls()
        isVideoBuffering.value = false
        uni.showToast({ title: '视频解密失败，请刷新重试', icon: 'none' })
        return
      }

      if (data.fatal) {
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
            hlsFatalErrorCount++
            if (hlsFatalErrorCount > HLS_MAX_FATAL_ERRORS) {
              console.error('[hls.js] 致命错误次数超限，停止播放')
              hlsInstance.stopLoad()
              destroyHls()
              isVideoBuffering.value = false
              return
            }
            console.error(`[hls.js] 网络错误，尝试恢复（${hlsFatalErrorCount}/${HLS_MAX_FATAL_ERRORS}）...`)
            hlsInstance.startLoad()
            break
          case Hls.ErrorTypes.MEDIA_ERROR:
            hlsFatalErrorCount++
            if (hlsFatalErrorCount > HLS_MAX_FATAL_ERRORS) {
              console.error('[hls.js] 致命错误次数超限，停止播放')
              hlsInstance.stopLoad()
              destroyHls()
              isVideoBuffering.value = false
              return
            }
            console.error(`[hls.js] 媒体错误，尝试恢复（${hlsFatalErrorCount}/${HLS_MAX_FATAL_ERRORS}）...`)
            hlsInstance.recoverMediaError()
            break
          default:
            destroyHls()
            break
        }
      }
    })
  } else if (nativeVideoEl.canPlayType('application/vnd.apple.mpegurl')) {
    // Safari：原生支持 HLS，直接设 src
    nativeVideoEl.src = blobUrl
  }
}

const destroyHls = () => {
  clearLoadingTimeout()
  hlsFatalErrorCount = 0
  if (hlsInstance) {
    hlsInstance.destroy()
    hlsInstance = null
  }
  if (currentBlobUrl) {
    URL.revokeObjectURL(currentBlobUrl)
    currentBlobUrl = null
  }
}

const initHlsIfNeeded = () => {
  const url = videoInfo.value.videoUrl || ''
  if (isHlsUrl(url)) {
    initHls(url)
  } else {
    destroyHls()
  }
}
// #endif
// ========== hls.js 集成结束 ==========

const maxTrialAllowedTime = () => {
  const duration = Number(videoInfo.value.duration || 0)
  if (authData.value.hasFullAuth || duration <= 0) {
    return Number.POSITIVE_INFINITY
  }
  const safePercent = Math.max(1, Number(authData.value.trialPercent || 8))
  return Math.max(0, Math.floor((duration * safePercent) / 100) - 1)
}

const ensureTrialLimit = (currentTime) => {
  if (authData.value.hasFullAuth || !videoInfo.value.duration) {
    return false
  }

  const limitTime = maxTrialAllowedTime()
  if (currentTime <= limitTime) {
    if (isTrialLocked.value) {
      isTrialLocked.value = false
    }
    return false
  }

  if (videoCtx) {
    videoCtx.pause()
  }
  seekVideo(limitTime)
  isTrialLocked.value = true

  const now = Date.now()
  if (now - lastTrialModalAt.value < 4000) {
    return true
  }
  lastTrialModalAt.value = now

  uni.showModal({
    title: t('player.trial_end_title'),
    content: t('player.trial_end_desc'),
    confirmText: t('player.btn_vip'),
    success: (res) => {
      if (res.confirm) {
        uni.navigateTo({ url: '/pages/learning/profile' })
      }
    }
  })
  return true
}

const focusSentence = () => {
  if (!sentenceId.value || subtitleList.value.length === 0) {
    return false
  }
  const targetIndex = subtitleList.value.findIndex((item) => Number(item.id) === sentenceId.value)
  if (targetIndex < 0) {
    return false
  }

  activeSentenceIndex.value = targetIndex
  activeSentenceId.value = 'sentence_' + targetIndex
  focusedSentenceId.value = sentenceId.value
  if (focusSentenceTimer) {
    clearTimeout(focusSentenceTimer)
  }
  focusSentenceTimer = setTimeout(() => {
    focusedSentenceId.value = 0
  }, 2500)
  seekVideo(subtitleList.value[targetIndex]?.startTime || 0)
  uni.showToast({ title: t('player.sentence_focused'), icon: 'none' })
  return true
}

const loadEpisode = async () => {
  const res = await findVideoEpisode(episodeId.value)
  if (res.code !== 0 || !res.data) {
    uni.showToast({ title: t('player.load_failed'), icon: 'none' })
    return
  }
  const data = res.data
  videoInfo.value = {
    id: getEntityId(data),
    name: data?.name || '',
    videoUrl: data?.videoUrl || '',
    duration: Number(data?.duration || 0)
  }
  authData.value.hasFullAuth = !!data?.hasFullAuth
  const percent = Number(data?.trialPercent || 8)
  authData.value.trialPercent = percent > 0 ? percent : 8
}

const signedUrlDeadline = (rawUrl) => {
  if (!rawUrl) {
    return 0
  }
  try {
    const parsed = new URL(rawUrl)
    const t = parsed.searchParams.get('t')
    if (!t) {
      return 0
    }
    const deadline = Number.parseInt(t, 16)
    return Number.isFinite(deadline) && deadline > 0 ? deadline : 0
  } catch (e) {
    return 0
  }
}

const refreshSignedEpisodeUrl = async (resumeSecs = 0) => {
  if (!episodeId.value || isRefreshingSignedUrl.value) {
    return false
  }

  isRefreshingSignedUrl.value = true
  isVideoBuffering.value = true
  isVideoReady.value = false
  const wasPlayingBeforeRefresh = isPlaying.value
  // 刷新过程中重置播放状态，等新源就绪后再恢复
  isPlaying.value = false
  const previousUrl = videoInfo.value.videoUrl
  const resumeTarget = Math.max(0, Number(resumeSecs || 0))
  const shouldResumeAfterRefresh = wasPlayingBeforeRefresh || autoPlayPending.value || resumeTarget <= 0
  autoPlayPending.value = shouldResumeAfterRefresh

  try {
    const res = await findVideoEpisode(episodeId.value)
    if (res.code !== 0 || !res.data) {
      return false
    }

    const data = res.data
    const nextUrl = String(data?.videoUrl || '')
    if (!nextUrl) {
      return false
    }

    videoInfo.value.videoUrl = nextUrl
    authData.value.hasFullAuth = !!data?.hasFullAuth
    const percent = Number(data?.trialPercent || 8)
    authData.value.trialPercent = percent > 0 ? percent : 8

    const urlChanged = nextUrl !== previousUrl
    setTimeout(() => {
      if (urlChanged && resumeTarget > 0) {
        seekVideo(Math.max(0, resumeTarget - 1))
      }
      if (shouldResumeAfterRefresh && videoCtx && !isTrialLocked.value) {
        isPlaying.value = true
        safePlay()
      }
    }, 350)

    return true
  } finally {
    isRefreshingSignedUrl.value = false
  }
}

const refreshSignedUrlIfNeeded = (currentSecs) => {
  const deadline = signedUrlDeadline(videoInfo.value.videoUrl)
  if (!deadline) {
    return
  }

  const now = Math.floor(Date.now() / 1000)
  if (deadline-now > 45) {
    return
  }

  if (isRefreshingSignedUrl.value || now-lastSignedRefreshAt < 10) {
    return
  }

  lastSignedRefreshAt = now
  refreshSignedEpisodeUrl(currentSecs)
}

const loadSubtitleList = async () => {
  const res = await getSentenceList(episodeId.value)
  if (res.code !== 0) {
    subtitleList.value = []
    return
  }
  const list = Array.isArray(res.data) ? res.data : []
  subtitleList.value = list.map((item, index) => ({
    ...item,
    id: getEntityId(item) || index + 1,
    startTime: Number(item?.startTime || 0),
    endTime: Number(item?.endTime || 0),
    english: String(item?.english || ''),
    translate: item?.translate || ''
  }))
  if (subtitleList.value.length === 0) {
    uni.showToast({ title: t('player.no_subtitle'), icon: 'none' })
  }
}

const loadSentenceCollectionState = async () => {
  const subtitleIds = subtitleList.value
    .map((item) => Number(item?.id || 0))
    .filter((id) => id > 0)

  if (subtitleIds.length === 0) {
    sentenceCollectedMap.value = {}
    return
  }

  const subtitleIdSet = new Set(subtitleIds)
  const map = {}
  const maxPages = 20
  let page = 1

  try {
    while (page <= maxPages) {
      const res = await getCollectionList({ targetType: 2, page, pageSize: 100 })
      if (res.code !== 0 || !res.data) {
        sentenceCollectedMap.value = {}
        return
      }

      const list = Array.isArray(res.data.list) ? res.data.list : []
      for (const item of list) {
        const targetId = Number(item?.targetId || 0)
        if (targetId > 0 && subtitleIdSet.has(targetId)) {
          map[targetId] = true
        }
      }

      if (Object.keys(map).length >= subtitleIdSet.size) {
        break
      }
      if (list.length < 100) {
        break
      }
      page += 1
    }

    sentenceCollectedMap.value = map
  } catch (error) {
    sentenceCollectedMap.value = {}
  }
}

const loadWatchHistory = async () => {
  const res = await getWatchProgress(episodeId.value)
  if (res.code !== 0 || !res.data) {
    return
  }
  const progress = Number(res.data.progressSecs || 0)
  if (progress <= 0) {
    return
  }
  seekVideo(progress)
}

const initPage = async () => {
  if (!episodeId.value) {
    uni.showToast({ title: t('player.load_failed'), icon: 'none' })
    return
  }
  autoPlayPending.value = true
  isVideoBuffering.value = true
  isVideoReady.value = false
  startLoadingTimeout()
  lastPlaybackTime.value = 0
  await loadEpisode()
  await loadSubtitleList()
  await loadSentenceCollectionState()
  const focused = focusSentence()
  if (!focused) {
    if (sentenceId.value > 0) {
      uni.showToast({ title: t('player.sentence_not_found'), icon: 'none' })
    }
    await loadWatchHistory()
  } else {
    // 3.1.1 字幕初始定位：等待 scroll-view DOM 渲染后强制定位
    await new Promise(resolve => setTimeout(resolve, 350))
    reAnchorScroll()
  }
}

onLoad((options) => {
  // H5 浏览器自动播放限制：:autoplay 属性触发原生 video.play() 可能被浏览器拒绝，
  // 该异常不在 safePlay 管辖范围，需全局捕获静默处理。
  // #ifdef H5
  window.addEventListener('unhandledrejection', (e) => {
    if (e.reason && e.reason.name === 'NotAllowedError') {
      e.preventDefault()
    }
  })
  // #endif
  episodeId.value = Number(options?.id || 0)
  sentenceId.value = Number(options?.sentenceId || 0)
  initPage()
})

// 心跳包计时期防伪
let lastHeartbeatTime = 0
let accumulatedViewTime = 0
let queuedHeartbeatSecs = 0
let heartbeatSending = false

const HEARTBEAT_FLUSH_THRESHOLD = 30
const HEARTBEAT_CHUNK_MAX = 120
const HEARTBEAT_QUEUE_MAX = 600

const enqueueHeartbeatSecs = (secs) => {
  const delta = Math.max(0, Math.floor(Number(secs || 0)))
  if (delta <= 0) {
    return
  }
  queuedHeartbeatSecs = Math.min(HEARTBEAT_QUEUE_MAX, queuedHeartbeatSecs + delta)
}

const sendHeartbeatQueue = async (progressSecs) => {
  if (!episodeId.value || heartbeatSending || queuedHeartbeatSecs <= 0) {
    return
  }
  heartbeatSending = true
  try {
    while (queuedHeartbeatSecs > 0) {
      const usingTimeSecs = Math.min(HEARTBEAT_CHUNK_MAX, queuedHeartbeatSecs)
      const res = await heartbeat({
        episodeId: episodeId.value,
        progressSecs,
        usingTimeSecs
      })
      if (!res || res.code !== 0) {
        break
      }
      queuedHeartbeatSecs -= usingTimeSecs
    }
  } finally {
    heartbeatSending = false
  }
}

// 核心：视频时间推进事件
const onTimeUpdate = (e) => {
  const curTime = e.detail.currentTime
  const duration = e.detail.duration || videoInfo.value.duration
  // 3.1.8 进度条：同步当前播放时间（拖拽中不更新，避免抖动）
  if (!isProgressDragging) {
    currentTime.value = curTime
  }
  if (curTime > lastPlaybackTime.value + 0.05) {
    isVideoBuffering.value = false
    isVideoReady.value = true
  }
  lastPlaybackTime.value = curTime
  if (duration > 0) {
    videoInfo.value.duration = duration
  }

  // 1. 防盗验证：超出试看比例时回弹到阈值，不重复弹窗轰炸。
  if (ensureTrialLimit(curTime)) {
    return
  }

  // 2. 匹配外挂字幕轴并滚动（复读模式下冻结，避免短句闪烁和相邻句跳跃）
  if (!repeatMode.value) {
    const currentIdx = subtitleList.value.findIndex(sub => curTime >= sub.startTime && curTime <= sub.endTime)
    if (currentIdx !== -1 && currentIdx !== activeSentenceIndex.value) {
      activeSentenceIndex.value = currentIdx
      activeSentenceId.value = 'sentence_' + currentIdx
    } else if (currentIdx === -1) {
      activeSentenceIndex.value = -1 // 空白期
    }
  }

  // 2b. 3.1.6 复读模式：当前句播放完毕后自动 seek 回句首
  if (repeatMode.value) {
    const repeatIdx = subtitleList.value.findIndex(sub => curTime >= sub.startTime && curTime <= sub.endTime)
    if (repeatIdx !== -1) {
      const curSentence = subtitleList.value[repeatIdx]
      if (curSentence && curTime >= curSentence.endTime - 0.3) {
        const now = Date.now()
        if (now - lastRepeatSeekTime > 600) {
          lastRepeatSeekTime = now
          seekVideo(curSentence.startTime)
        }
      }
    }
  }

  // 3. 签名链接即将到期时，提前刷新，避免播放中途 403。
  refreshSignedUrlIfNeeded(curTime)

  // 4. 心跳累加与上报逻辑 (Phase 2 的前端呼应)
  // 如果正常播放，累加观看差值，每满 30 秒 Post 给后端一次
  if (lastHeartbeatTime > 0 && curTime > lastHeartbeatTime) {
    accumulatedViewTime += (curTime - lastHeartbeatTime)
    if (accumulatedViewTime >= HEARTBEAT_FLUSH_THRESHOLD) {
      flushHeartbeat(curTime)
    }
  }
  lastHeartbeatTime = curTime
}

const flushHeartbeat = (progressSecs, force = false) => {
  const usingTimeSecs = Math.floor(accumulatedViewTime)
  if (usingTimeSecs > 0) {
    enqueueHeartbeatSecs(usingTimeSecs)
    accumulatedViewTime = 0
  }
  if (!force && queuedHeartbeatSecs < HEARTBEAT_FLUSH_THRESHOLD) {
    return
  }
  void sendHeartbeatQueue(progressSecs)
}

const onVideoEnded = () => {
  flushHeartbeat(videoInfo.value.duration || lastHeartbeatTime, true)
}

const onVideoLoadedMeta = (e) => {
  const duration = Number(e?.detail?.duration || 0)
  if (duration > 0) {
    videoInfo.value.duration = duration
  }
  // 新视频源加载中，重置播放状态避免僵死
  isPlaying.value = false
  isVideoBuffering.value = true
  isVideoReady.value = false
  startLoadingTimeout()
}

const onVideoCanPlay = () => {
  clearLoadingTimeout()
  isVideoBuffering.value = false
  isVideoReady.value = true
  if (autoPlayPending.value && videoCtx && !isTrialLocked.value) {
    safePlay()
  }
}

const onVideoSeeking = () => {
  isVideoBuffering.value = true
  startLoadingTimeout()
}

const onVideoSeeked = () => {
  clearLoadingTimeout()
  if (videoCtx && !isTrialLocked.value) {
    safePlay()
  }
}

const onVideoWaiting = () => {
  isVideoBuffering.value = true
  isVideoReady.value = false
  startLoadingTimeout()
}

const onVideoPlaying = () => {
  clearLoadingTimeout()
  isVideoBuffering.value = false
  isVideoReady.value = true
  isPlaying.value = true
  autoPlayPending.value = false
}

const onVideoPlay = () => {
  clearLoadingTimeout()
  isVideoBuffering.value = false
  isPlaying.value = true
  isVideoReady.value = true
  autoPlayPending.value = false
}

const onVideoPause = () => {
  isPlaying.value = false
}

const onVideoError = async () => {
  // 出错时重置播放状态，避免按钮显示 ⏸ 但实际已停止
  isPlaying.value = false
  isVideoBuffering.value = true
  isVideoReady.value = false
  autoPlayPending.value = true
  clearLoadingTimeout()
  const recovered = await refreshSignedEpisodeUrl(lastHeartbeatTime)
  if (!recovered) {
    uni.showToast({ title: t('player.load_failed'), icon: 'none' })
  }
}

const jumpBySubtitle = (item, index) => {
  showPauseHint.value = false
  const startTime = Number(item?.startTime || 0)
  activeSentenceIndex.value = index
  activeSentenceId.value = 'sentence_' + index
  seekVideo(startTime)
  heartbeat({
    episodeId: episodeId.value,
    progressSecs: startTime,
    usingTimeSecs: 0
  }).catch(() => {})
  if (videoCtx) {
    safePlay()
  }
}

// 重点词渲染解析器 (<w id="x">) - 将其拆分成 Vue 安全遍历结构
const parseHighlight = (htmlStr) => {
  // 简易正则拆分出 <w id="x">文本</w>
  const parts = []
  const regex = /<w id="(\d+)">(.*?)<\/w>/g
  let lastIndex = 0
  let match

  while ((match = regex.exec(htmlStr)) !== null) {
    if (match.index > lastIndex) {
      parts.push({ text: htmlStr.slice(lastIndex, match.index), isHighlight: false })
    }
    parts.push({ text: match[2], isHighlight: true, wordId: match[1] })
    lastIndex = regex.lastIndex
  }
  if (lastIndex < htmlStr.length) {
    parts.push({ text: htmlStr.slice(lastIndex), isHighlight: false })
  }
  return parts.length ? parts : [{ text: htmlStr, isHighlight: false }]
}

// 字幕区文字点击：重点词阻止冒泡弹窗查词，非重点词透传给父级jumpBySubtitle
const handleSegClick = (seg, event) => {
  if (seg && seg.isHighlight) {
    if (event && event.stopPropagation) event.stopPropagation()
    showWordDetail(seg.wordId)
  }
  // 非重点词不做任何处理，事件自然冒泡到父级 .sentence-row 触发 jumpBySubtitle
}

// 查词查词交互
const wordPopup = ref(null)
const currentWord = ref(null)
const showWordDetail = async (wordId) => {
  // 暂停时点击重点单词只弹出单词窗，不改变播放状态
  if (isPlaying.value && videoCtx) {
    videoCtx.pause()
  }
  const res = await findWord(Number(wordId))
  if (res.code !== 0 || !res.data) {
    return
  }
  currentWord.value = {
    ...res.data,
    id: getEntityId(res.data)
  }
  wordPopup.value.open()
}

// 底部控制
const togglePlay = () => {
  showPauseHint.value = false
  if (!videoCtx) {
    return
  }
  if (!authData.value.hasFullAuth && isTrialLocked.value) {
    return
  }
  if (showVideoLoading.value) {
    autoPlayPending.value = true
    safePlay()
    return
  }
  // H5 端：用原生 video.paused 纠正可能僵死的播放状态
  // #ifdef H5
  if (nativeVideoEl && nativeVideoEl.paused !== !isPlaying.value) {
    isPlaying.value = !nativeVideoEl.paused
  }
  // #endif
  isPlaying.value ? videoCtx.pause() : safePlay()
}
const toggleBilingual = () => { 
  controls.value.showBilingual = !controls.value.showBilingual
  reAnchorScroll()
}
const toggleSubtitles = () => { 
  controls.value.showSubtitles = !controls.value.showSubtitles
  reAnchorScroll()
}
const toggleLoop = () => { controls.value.loop = !controls.value.loop }

// 3.1.6 复读按钮：切换单句循环模式
const toggleRepeat = () => {
  repeatMode.value = !repeatMode.value
  if (repeatMode.value) {
    uni.showToast({ title: t('player.repeat_on'), icon: 'none' })
  }
}

// 3.1.6 在更多设置面板中选择播放速率
const selectSpeed = (rate) => {
  controls.value.speed = rate
  if (videoCtx) {
    videoCtx.playbackRate(rate)
  }
}

// 3.1.2 切换双语/字幕时重新锚定滚动位置
const reAnchorScroll = () => {
  if (activeSentenceIndex.value < 0) return
  // 先重置触发 scroll-view 重新计算布局
  activeSentenceId.value = ''
  setTimeout(() => {
    activeSentenceId.value = 'sentence_' + Math.max(0, activeSentenceIndex.value)
  }, 50)
}

// 3.1.4 手势滑动 → 暂停视频并显示提示
const onSubtitleTouchStart = () => {
  if (videoCtx && isPlaying.value) {
    videoCtx.pause()
    showPauseHint.value = true
  }
}
const onSubtitleTouchMove = () => {
  // touchmove 时保持提示可见
  if (!showPauseHint.value && videoCtx && isPlaying.value) {
    videoCtx.pause()
    showPauseHint.value = true
  }
}
const hidePauseHint = () => {
  showPauseHint.value = false
}

// 3.1.8 进度条拖拽跳转
const onProgressChange = (e) => {
  isProgressDragging = false
  const value = Number(e.detail.value || 0)
  currentTime.value = value
  seekVideo(value)
  if (videoCtx && !isTrialLocked.value) {
    safePlay()
  }
}
const onProgressChanging = (e) => {
  isProgressDragging = true
  currentTime.value = Number(e.detail.value || 0)
}

const openMoreSheet = () => {
  morePopup.value?.open()
}

const closeMoreSheet = () => {
  morePopup.value?.close()
}

const handleBack = () => {
  closeMoreSheet()
  const pages = getCurrentPages()
  if (Array.isArray(pages) && pages.length > 1) {
    uni.navigateBack()
    return
  }
  uni.switchTab({ url: '/pages/learning/home' })
}
const collectSentence = async (id) => {
  const sentenceIdNum = Number(id || 0)
  if (!sentenceIdNum) {
    return
  }
  const collected = !!sentenceCollectedMap.value[sentenceIdNum]
  const res = collected ? await uncollect(2, sentenceIdNum) : await collect(2, sentenceIdNum)
  if (res.code === 0) {
    sentenceCollectedMap.value = {
      ...sentenceCollectedMap.value,
      [sentenceIdNum]: !collected,
    }
    uni.showToast({ title: !collected ? t('player.collect_success') : t('typing.uncollect_success'), icon: 'success' })
  }
}

const collectWord = async (id) => {
  const res = await collect(1, Number(id || 0))
  if (res.code === 0) {
    uni.showToast({ title: t('player.collect_success'), icon: 'success' })
  }
}

const playWordAudio = (src) => {
  if (!src) {
    return
  }
  if (!audioCtx) {
    audioCtx = uni.createInnerAudioContext()
  }
  audioCtx.src = getExternalUrl(src)
  audioCtx.play()
}
</script>

<style scoped>
.player-container { display: flex; flex-direction: column; height: 100vh; min-height: 100dvh; background: linear-gradient(180deg, #1A1B3A 0%, #2D2E5C 100%); color: #fff;}

/* 视频区域 */
.video-section { width: 100%; height: 420rpx; background: #000; flex-shrink: 0; position: relative; }
.main-video { width: 100%; height: 100%; }

.video-loading-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 420rpx;
  background: rgba(26, 27, 58, 0.66);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14rpx;
  pointer-events: none;
}

.video-loading-logo,
.video-loading-logo-fallback {
  width: 88rpx;
  height: 88rpx;
  border-radius: 22rpx;
}

.video-loading-logo-fallback {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 34rpx;
  font-weight: 700;
}

.video-loading-text {
  color: #e2e8f0;
  font-size: 24rpx;
}

@keyframes logo-spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.spin {
  animation: logo-spin 1.5s linear infinite;
}

/* 外挂字幕区域 */
.subtitle-section { flex: 1; overflow: hidden; position: relative; padding-bottom: calc(132rpx + env(safe-area-inset-bottom)); box-sizing: border-box; }
.subtitle-padding-top, .subtitle-padding-bottom { height: 20rpx; }
.sentence-row { padding: 20rpx 40rpx; margin-bottom: 20rpx; transition: all 0.3s; opacity: 0.6; border-radius: 14rpx; background: rgba(108, 91, 255, 0.06); }
.sentence-row.is-active { opacity: 1; transform: scale(1.05); background: rgba(108, 91, 255, 0.12); border-left: 6rpx solid #6D5BFF; box-shadow: 0 0 0 2rpx rgba(108, 91, 255, 0.28) inset; }
.sentence-row.is-focused { opacity: 1; border-left: 6rpx solid #9B8FFF; box-shadow: 0 0 0 2rpx rgba(155, 143, 255, 0.35) inset; }

.sentence-meta { display: flex; justify-content: space-between; align-items: center; margin-bottom: 12rpx; }
.time-badge { font-size: 20rpx; color: #9B8FFF; background: rgba(108, 91, 255, 0.14); padding: 4rpx 14rpx; border-radius: 999rpx; font-variant-numeric: tabular-nums; }
.collect-btn { font-size: 20rpx; color: #9B8FFF; padding: 4rpx 16rpx; border-radius: 999rpx; border: 1rpx solid rgba(155, 143, 255, 0.3); background: rgba(155, 143, 255, 0.06); }

/* 自主控制高亮词汇，绕过 rich-text 不能点击内部标签的缺陷 */
.english-text { font-size: 34rpx; font-weight: 500; margin-bottom: 0; line-height: 1.5; color: #ececec; }
.highlight-word { color: #9B8FFF; font-weight: bold; border-bottom: 1px dashed #9B8FFF; display: inline-block; padding: 0 4rpx; margin: 0 4rpx; }
.translated-text { font-size: 28rpx; color: #bbb; line-height: 1.5; margin-top: 16rpx; padding-top: 16rpx; border-top: 1rpx solid rgba(108, 91, 255, 0.1); }

/* 底部操作区 */
.bottom-controls {
  height: 132rpx;
  background: rgba(26, 27, 58, 0.96);
  flex-shrink: 0;
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding-bottom: env(safe-area-inset-bottom);
  border-top: 1px solid rgba(108, 91, 255, 0.24);
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 30;
}
.control-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6rpx;
}
.c-icon-circle {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(108, 91, 255, 0.08);
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  transition: all 0.2s;
}
.c-circle-on {
  background: rgba(108, 91, 255, 0.24);
  border-color: rgba(108, 91, 255, 0.5);
  box-shadow: 0 0 12rpx rgba(155, 143, 255, 0.4);
}
.c-icon {
  font-size: 32rpx;
  filter: grayscale(0.6);
  opacity: 0.7;
}
.c-circle-on .c-icon {
  filter: grayscale(0);
  opacity: 1;
}
.c-text { font-size: 20rpx; color: #888; }
.c-text-active { color: #9B8FFF; font-weight: 600; }
.play-btn { gap: 0; }
.c-play-circle {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 4rpx 20rpx rgba(108, 91, 255, 0.5);
}
.c-play-icon {
  font-size: 44rpx;
  filter: grayscale(0);
  opacity: 1;
  color: #fff;
}

/* 弹窗拖动条 */
.popup-handle {
  width: 64rpx;
  height: 8rpx;
  border-radius: 999rpx;
  background: rgba(108, 91, 255, 0.2);
  margin: 0 auto 28rpx auto;
}

/* 弹窗内容 */
.word-detail-box {
  padding: 24rpx 40rpx 40rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  background: #FFFFFF;
  border-top-left-radius: 28rpx;
  border-top-right-radius: 28rpx;
}
.w-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 8rpx; }
.w-title { font-size: 48rpx; font-weight: 800; color: #1A1B3A; }
.w-audio-circle {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(108, 91, 255, 0.12);
  border: 1rpx solid rgba(108, 91, 255, 0.2);
}
.w-audio-icon { font-size: 36rpx; }
.w-phonetic { font-size: 28rpx; color: #6B6F8D; margin: 4rpx 0 20rpx 0; }
.w-divider { height: 1rpx; background: rgba(108, 91, 255, 0.16); margin-bottom: 24rpx; }
.w-exp { font-size: 32rpx; color: #1A1B3A; line-height: 1.6; margin-bottom: 36rpx; }
.w-collect {
  width: 100%;
  height: 88rpx;
  line-height: 88rpx;
  margin: 0;
  border: none;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  border-radius: 999rpx;
  font-size: 30rpx;
  font-weight: 700;
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.3);
}

.more-sheet {
  padding: 24rpx 24rpx;
  padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
  background: #FFFFFF;
  color: #1A1B3A;
  border-top-left-radius: 28rpx;
  border-top-right-radius: 28rpx;
  box-shadow: 0 -12rpx 30rpx rgba(26, 27, 58, 0.12);
}
.more-sheet-head {
  margin-bottom: 16rpx;
}

.more-sheet-title {
  display: block;
  font-size: 32rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.more-sheet-subtitle {
  display: block;
  margin-top: 6rpx;
  font-size: 22rpx;
  color: #6B6F8D;
}

.more-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  border-radius: 24rpx;
  background: #FFFFFF;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.06);
}

.more-card-main {
  display: flex;
  flex-direction: column;
}

.more-card-title {
  font-size: 30rpx;
  color: #1A1B3A;
  font-weight: 700;
}

.more-card-desc {
  margin-top: 6rpx;
  font-size: 22rpx;
  color: #6B6F8D;
}

.more-state-pill {
  min-width: 120rpx;
  height: 52rpx;
  padding: 0 18rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  color: #6B6F8D;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  font-weight: 700;
}

.more-state-pill.on {
  color: #6D5BFF;
  border-color: rgba(108, 91, 255, 0.4);
  background: rgba(108, 91, 255, 0.12);
}

.more-actions {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  margin-top: 24rpx;
}

.more-action-btn {
  width: 100%;
  height: 84rpx;
  line-height: 84rpx;
  margin: 0;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 6rpx 20rpx rgba(108, 91, 255, 0.25);
}

.more-action-btn.ghost {
  color: #6D5BFF;
  background: #fff;
  border: 1rpx solid rgba(108, 91, 255, 0.3);
}

/* 3.1.8 进度条 */
.progress-section {
  flex-shrink: 0;
  padding: 8rpx 32rpx 8rpx 32rpx;
  background: rgba(26, 27, 58, 0.95);
}

.progress-time-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rpx;
}

.progress-current-time {
  font-size: 22rpx;
  color: #9B8FFF;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}

.progress-total-time {
  font-size: 22rpx;
  color: #94a3b8;
  font-variant-numeric: tabular-nums;
}

.progress-row {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.video-progress-slider {
  flex: 1;
  margin: 0;
}

/* 3.1.4 手势滑动暂停提示 */
.pause-hint-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(26, 27, 58, 0.55);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

.pause-hint-text {
  color: #fff;
  font-size: 28rpx;
  background: rgba(108, 91, 255, 0.85);
  padding: 16rpx 36rpx;
  border-radius: 999rpx;
  letter-spacing: 2rpx;
}

/* 3.1.6 速率选择器 */
.speed-card {
  margin-top: 14rpx;
}

.speed-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-top: 14rpx;
}

.speed-pill {
  flex: 1;
  min-width: 80rpx;
  height: 60rpx;
  line-height: 60rpx;
  text-align: center;
  border-radius: 999rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  font-size: 24rpx;
  font-weight: 600;
  color: #6B6F8D;
  background: #fff;
}

.speed-pill-active {
  color: #fff;
  border-color: transparent;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 4rpx 12rpx rgba(108, 91, 255, 0.3);
}
</style>
