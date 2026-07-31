<template>
  <div class="player-page">
    <!-- 顶部返回 + 标题 -->
    <div class="player-header">
      <button class="back-btn" @click="handleBack">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
        {{ t('common.back') }}
      </button>
      <h2 class="video-title">{{ videoInfo?.name || '视频播放' }}</h2>
      <div class="header-actions">
        <button class="action-btn" :class="{ active: isEpisodeCollected }" @click="toggleEpisodeCollect" :title="isEpisodeCollected ? t('common.uncollect') : t('common.collect')">
          <svg viewBox="0 0 24 24" :fill="isEpisodeCollected ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
          </svg>
        </button>
      </div>
    </div>

    <div class="player-body">
      <!-- 左侧视频区 -->
      <div class="video-section">
        <div class="video-wrapper" @mousemove="onMouseMove" @mouseleave="showControls=true">
          <video
            ref="videoRef"
            class="video-player"
            playsinline
            :loop="controls.loop"
            @timeupdate="onTimeUpdate"
            @loadedmetadata="onLoadedMetadata"
            @canplay="onCanPlay"
            @seeking="onSeeking"
            @seeked="onSeeked"
            @waiting="onWaiting"
            @playing="onPlaying"
            @play="onPlay"
            @pause="onPause"
            @error="onVideoError"
            @ended="onEnded"
          ></video>

          <!-- 加载遮罩 -->
          <div v-if="showVideoLoading" class="video-loading-mask">
            <div class="loading-logo spin">R</div>
            <span class="loading-text">{{ t('player.video_loading') }}</span>
          </div>

          <!-- 暂停提示 -->
          <div v-if="showPauseHint" class="pause-hint-overlay" @click="hidePauseHint">
            <span class="pause-hint-text">{{ t('player.tap_subtitle_to_play') }}</span>
          </div>

          <!-- 试看遮罩 -->
          <div v-if="isTrialLocked" class="trial-mask" @click="onTrialClick">
            <div class="trial-card">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              <h3>{{ t('player.trial_end_title') }}</h3>
              <p>{{ t('player.trial_end_desc') }}</p>
              <button class="btn btn-primary" @click.stop="goVip">{{ t('player.btn_vip') }}</button>
            </div>
          </div>

          <!-- 控制栏 -->
          <div class="video-controls" :class="{ hidden: !showControls && isPlaying }">
            <!-- 进度条 -->
            <div class="progress-bar" @mousedown="onProgressMouseDown">
              <div class="progress-played" :style="{ width: progressPercent + '%' }">
                <div class="progress-thumb" :class="{ dragging: isDragging }"></div>
              </div>
              <div v-if="!authData.hasFullAuth && authData.trialPercent < 100" class="trial-marker" :style="{ left: authData.trialPercent + '%' }">
                <span class="trial-label">{{ authData.trialPercent }}%</span>
              </div>
            </div>

            <div class="controls-row">
              <div class="controls-left">
                <button class="ctrl-btn" @click="togglePlay">
                  <svg v-if="!isPlaying" viewBox="0 0 24 24" fill="currentColor">
                    <polygon points="5 3 19 12 5 21 5 3"></polygon>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="currentColor">
                    <rect x="6" y="4" width="4" height="16"></rect>
                    <rect x="14" y="4" width="4" height="16"></rect>
                  </svg>
                </button>
                <span class="time-text">{{ formatTime(currentTime) }} / {{ formatTime(duration) }}</span>
              </div>

              <div class="controls-center">
                <button class="ctrl-btn" :class="{ active: controls.showBilingual }" @click="toggleBilingual" :title="t('player.bilingual')">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M5 8l6 6"></path>
                    <path d="M4 14l6-6 2-3"></path>
                    <path d="M2 5h12"></path>
                    <path d="M7 2h1"></path>
                    <path d="M22 22l-5-10-5 10"></path>
                    <path d="M14 18h6"></path>
                  </svg>
                </button>

                <button class="ctrl-btn" :class="{ active: isRepeatMode }" @click="toggleRepeat" :title="t('player.repeat')">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="1 4 1 10 7 10"></polyline>
                    <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path>
                  </svg>
                </button>

                <button class="ctrl-btn" :class="{ active: controls.showSubtitles }" @click="toggleSubtitles" :title="t('player.subtitle')">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="2" y="4" width="20" height="16" rx="2"></rect>
                    <path d="M6 12h4M14 12h4M6 16h4M14 16h2"></path>
                  </svg>
                </button>

                <div class="speed-selector">
                  <button class="ctrl-btn speed-btn" @click="showSettings = !showSettings">
                    {{ controls.speed }}x
                  </button>
                  <div v-if="showSettings" class="speed-menu">
                    <div class="menu-title">{{ t('player.speed') }}</div>
                    <div
                      v-for="rate in speedRates"
                      :key="rate"
                      class="speed-item"
                      :class="{ active: controls.speed === rate }"
                      @click="selectSpeed(rate)"
                    >
                      {{ rate }}x
                    </div>
                    <div class="speed-divider"></div>
                    <div class="speed-item" :class="{ active: controls.loop }" @click="toggleLoop">
                      <span>{{ t('player.loop') }}</span>
                      <span v-if="controls.loop" class="check">✓</span>
                    </div>
                    <div class="speed-divider"></div>
                    <div class="speed-item" @click="handleBack">
                      <span>{{ t('player.exit_page') }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧字幕区 -->
      <div v-show="controls.showSubtitles" class="subtitle-section">
        <div class="subtitle-header">
          <h3>{{ t('player.subtitle') }}</h3>
        </div>

        <div ref="subtitleScrollRef" class="subtitle-scroll" @mousedown="onSubtitleMouseDown">
          <div class="subtitle-padding-top"></div>
          <div
            v-for="(item, index) in subtitleList"
            :key="index"
            :ref="el => { if (index === activeSentenceIndex) activeEl = el }"
            class="sentence-row"
            :class="{
              'is-active': activeSentenceIndex === index,
              'is-focused': focusedSentenceId === getEntityId(item)
            }"
            @click="!item._locked && jumpBySubtitle(item, index)"
          >
            <div class="sentence-meta">
              <span class="time-badge">{{ formatTime(item.startTime) }} - {{ formatTime(item.endTime) }}</span>
              <span class="collect-btn" @click.stop="collectSentence(item)">
                {{ sentenceCollectedMap[getEntityId(item)] ? '★' : '☆' }}
              </span>
            </div>

            <div v-if="item._locked" class="sentence-locked">
              <span class="lock-time">{{ formatTime(item.startTime) }} - {{ formatTime(item.endTime) }}</span>
              <span class="lock-icon">🔒</span>
            </div>
            <div v-else class="sentence-text-wrap">
              <div class="english-text">
                <template v-for="(seg, sIdx) in parseHighlight(item.english)" :key="sIdx">
                  <span
                    v-if="seg.isHighlight"
                    class="highlight-word"
                    @click="handleSegClick(seg, $event)"
                  >{{ seg.text }}</span>
                  <span v-else>{{ seg.text }}</span>
                </template>
              </div>
              <div v-if="controls.showBilingual" class="translated-text">
                {{ localText(item.translate) }}
              </div>
            </div>
          </div>
          <div class="subtitle-padding-bottom"></div>
        </div>
      </div>
    </div>

    <!-- 单词弹窗 -->
    <div v-if="wordPopup.visible && wordPopup.word" class="word-popup-overlay" @click="closeWordPopup">
      <div class="word-detail-box" @click.stop>
        <div class="w-header">
          <span class="w-title">{{ wordPopup.word.word }}</span>
          <button class="w-audio-circle" :class="{ playing: isWordAudioPlaying }" @click="playWordAudioAction">
            <span class="w-audio-icon">{{ isWordAudioPlaying ? '◉' : '♪' }}</span>
          </button>
        </div>
        <span class="w-phonetic">{{ wordPopup.word.phoneticUs }}</span>
        <div class="w-divider"></div>
        <div class="w-exp">{{ localText(wordPopup.word.explanation) }}</div>
        <div class="w-actions">
          <button class="w-cancel-btn" @click="closeWordPopup">{{ t('common.cancel') }}</button>
          <button class="w-collect" @click="collectCurrentWord">
            {{ isWordCollected ? '★ ' + t('common.uncollect') : '☆ ' + t('common.collect') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Hls from 'hls.js'
import {
  findVideoEpisode,
  getSentenceList,
  getWatchProgress,
  collect,
  uncollect,
  getCollectionList,
  heartbeat,
  findWord
} from '@/api/learning'
import { baseUrl } from '@/utils/request'
import { localText, resolveApiMessage } from '@/utils/i18n'
import { playWordAudio as playWordTTS, stopLocalTTS, destroyTTS } from '@/utils/learning-tts'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const episodeId = ref(Number(route.params.episodeId) || 0)
const sentenceId = ref(Number(route.query.sentenceId) || 0)

// Video state
const videoRef = ref(null)
let hlsInstance = null
let nativeVideoEl = null
let currentBlobUrl = null
let hlsFatalErrorCount = 0
const HLS_MAX_FATAL_ERRORS = 5

const isPlaying = ref(false)
const isVideoReady = ref(false)
const isVideoBuffering = ref(true)
const isRefreshingSignedUrl = ref(false)
const showVideoLoading = computed(() => isVideoBuffering.value || isRefreshingSignedUrl.value)
const autoPlayPending = ref(true)
const currentTime = ref(0)
const duration = ref(0)
let lastPlaybackTime = 0
let lastSignedRefreshAt = 0

const showControls = ref(true)
let controlsTimer = null
const isDragging = ref(false)
let wasPlayingBeforeDrag = false

// Video info & auth
const videoInfo = ref({ id: 0, name: '', videoUrl: '', duration: 0 })
const authData = ref({ hasFullAuth: false, trialPercent: 8 })
const isTrialLocked = ref(false)
let lastTrialModalAt = 0

// Controls
const controls = ref({
  showBilingual: true,
  showSubtitles: true,
  speed: 1.0,
  loop: false
})
const speedRates = [0.75, 1.0, 1.25, 1.5, 2.0]
const showSettings = ref(false)

// Subtitles
const subtitleList = ref([])
const activeSentenceIndex = ref(-1)
let activeEl = null
const focusedSentenceId = ref(0)
let focusSentenceTimer = null
const subtitleScrollRef = ref(null)
let isSubtitleDragging = false

// Repeat
const isRepeatMode = ref(false)
let lastRepeatSeekTime = 0

// Pause hint
const showPauseHint = ref(false)

// Episode collect
const isEpisodeCollected = ref(false)

// Sentence collect map
const sentenceCollectedMap = ref({})

// Word popup
const wordPopup = ref({ visible: false, word: null })
const currentWord = ref(null)
const isWordAudioPlaying = ref(false)
const isWordCollected = ref(false)

// Heartbeat
let lastHeartbeatTime = 0
let accumulatedViewTime = 0
let queuedHeartbeatSecs = 0
let heartbeatSending = false
const HEARTBEAT_FLUSH_THRESHOLD = 30
const HEARTBEAT_CHUNK_MAX = 120
const HEARTBEAT_QUEUE_MAX = 600

// Loading timeout
let loadingTimeoutTimer = null

const clearLoadingTimeout = () => {
  if (loadingTimeoutTimer) { clearTimeout(loadingTimeoutTimer); loadingTimeoutTimer = null }
}
const startLoadingTimeout = () => {
  clearLoadingTimeout()
  loadingTimeoutTimer = setTimeout(() => {
    if (isVideoBuffering.value && !isVideoReady.value) {
      console.warn('[player] 加载超时（8秒）')
      isVideoBuffering.value = false
    }
  }, 8000)
}

// Helpers
const getEntityId = (item) => Number(item?.id || item?.ID || 0)
const formatTime = (secs) => {
  if (!secs || isNaN(secs)) return '00:00'
  const m = Math.floor(secs / 60).toString().padStart(2, '0')
  const s = Math.floor(secs % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}
const isHlsUrl = (url) => /\.m3u8(\?|$)/i.test(url || '')

// Seek
const seekVideo = (secs) => {
  const target = Math.max(0, Number(secs || 0))
  const v = videoRef.value
  if (v) { v.currentTime = target } else { pendingSeekTime = target }
}
const safePlay = () => {
  const v = videoRef.value
  if (!v) return
  try { const p = v.play(); if (p && typeof p.catch === 'function') p.catch(() => {}) } catch (_) {}
}

// HLS
const rewriteM3u8 = async (m3u8Url) => {
  const resp = await fetch(m3u8Url)
  if (!resp.ok) throw new Error(`Fetch m3u8 failed: ${resp.status}`)
  const content = await resp.text()
  let backendBase = baseUrl
  if (backendBase.startsWith('/')) backendBase = window.location.origin + backendBase
  const m3u8Base = m3u8Url.substring(0, m3u8Url.lastIndexOf('/') + 1)
  const lines = content.split('\n')
  const rewritten = lines.map(line => {
    if (line.startsWith('#EXT-X-KEY')) {
      return line.replace(/URI="([^"]*)"/, (_m, uri) => {
        if (uri.startsWith('http://') || uri.startsWith('https://')) return _m
        return `URI="${backendBase}${uri.startsWith('/') ? uri : '/' + uri}"`
      })
    }
    if (line && !line.startsWith('#') && /\.ts(\?|$)/i.test(line)) {
      if (!line.startsWith('http://') && !line.startsWith('https://')) return m3u8Base + line
    }
    return line
  })
  const blob = new Blob([rewritten.join('\n')], { type: 'application/vnd.apple.mpegurl' })
  return URL.createObjectURL(blob)
}

const destroyHls = () => {
  clearLoadingTimeout()
  hlsFatalErrorCount = 0
  if (hlsInstance) { hlsInstance.destroy(); hlsInstance = null }
  if (currentBlobUrl) { URL.revokeObjectURL(currentBlobUrl); currentBlobUrl = null }
}

const initHls = async (url) => {
  destroyHls()
  hlsFatalErrorCount = 0
  nativeVideoEl = videoRef.value
  if (!nativeVideoEl || !url) return
  isVideoBuffering.value = true; isVideoReady.value = false; startLoadingTimeout()
  let blobUrl = url
  try { blobUrl = await rewriteM3u8(url); currentBlobUrl = blobUrl } catch (e) { console.error('[hls] rewrite failed', e) }
  if (Hls.isSupported()) {
    nativeVideoEl.removeAttribute('src')
    hlsInstance = new Hls({ fragLoadingMaxRetry: 0, manifestLoadingMaxRetry: 0, levelLoadingMaxRetry: 0 })
    hlsInstance.attachMedia(nativeVideoEl)
    hlsInstance.on(Hls.Events.MEDIA_ATTACHED, () => hlsInstance.loadSource(blobUrl))
    hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => { clearLoadingTimeout(); isVideoBuffering.value = false; isVideoReady.value = true })
    hlsInstance.on(Hls.Events.ERROR, (_event, data) => {
      console.error('[hls error]', data.type, data.details, data.fatal)
      if (data.details === 'keyLoadError' || data.details === 'fragDecryptError') {
        hlsInstance.stopLoad(); destroyHls(); isVideoBuffering.value = false; showToast(t('player.decrypt_failed')); return
      }
      if (data.fatal) {
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
            hlsFatalErrorCount++
            if (hlsFatalErrorCount > HLS_MAX_FATAL_ERRORS) { hlsInstance.stopLoad(); destroyHls(); isVideoBuffering.value = false; return }
            hlsInstance.startLoad(); break
          case Hls.ErrorTypes.MEDIA_ERROR:
            hlsFatalErrorCount++
            if (hlsFatalErrorCount > HLS_MAX_FATAL_ERRORS) { hlsInstance.stopLoad(); destroyHls(); isVideoBuffering.value = false; return }
            hlsInstance.recoverMediaError(); break
          default: destroyHls(); break
        }
      }
    })
  } else if (nativeVideoEl.canPlayType('application/vnd.apple.mpegurl')) {
    nativeVideoEl.src = blobUrl
  }
}

const initHlsIfNeeded = () => {
  const url = videoInfo.value.videoUrl || ''
  if (isHlsUrl(url)) initHls(url)
  else { destroyHls(); const v = videoRef.value; if (v && url) v.src = url }
}

// Trial
const maxTrialAllowedTime = () => {
  const d = Number(videoInfo.value.duration || 0)
  if (authData.value.hasFullAuth || d <= 0) return Number.POSITIVE_INFINITY
  return (d * Math.max(1, Number(authData.value.trialPercent || 8))) / 100
}

const ensureTrialLimit = (curTime) => {
  if (authData.value.hasFullAuth || !videoInfo.value.duration) return false
  const limit = maxTrialAllowedTime()
  if (curTime < limit) { if (isTrialLocked.value) isTrialLocked.value = false; return false }
  const v = videoRef.value
  if (v) v.pause()
  seekVideo(Math.max(0, limit - 0.5))
  isTrialLocked.value = true; isPlaying.value = false
  const now = Date.now()
  if (now - lastTrialModalAt < 4000) return true
  lastTrialModalAt = now
  showTrialModal()
  return true
}

const showTrialModal = () => {
  if (confirm(t('player.trial_end_title') + '\n' + t('player.trial_end_desc'))) router.push('/learning/profile')
}
const onTrialClick = () => showTrialModal()
const goVip = () => router.push('/learning/profile')

// Signed URL refresh
const signedUrlDeadline = (rawUrl) => {
  if (!rawUrl) return 0
  try { const p = new URL(rawUrl, window.location.origin); const tv = p.searchParams.get('t'); if (!tv) return 0; const d = parseInt(tv, 16); return isFinite(d) && d > 0 ? d : 0 } catch (e) { return 0 }
}

const refreshSignedEpisodeUrl = async (resumeSecs = 0) => {
  if (!episodeId.value || isRefreshingSignedUrl.value) return false
  isRefreshingSignedUrl.value = true; isVideoBuffering.value = true; isVideoReady.value = false
  const wasPlaying = isPlaying.value; isPlaying.value = false
  const prevUrl = videoInfo.value.videoUrl
  const resumeTarget = Math.max(0, Number(resumeSecs))
  const shouldResume = wasPlaying || autoPlayPending.value
  autoPlayPending.value = shouldResume
  try {
    const res = await findVideoEpisode(episodeId.value)
    if (res.code !== 0 || !res.data) return false
    const d = res.data
    const nextUrl = String(d?.videoUrl || '')
    if (!nextUrl) return false
    videoInfo.value.videoUrl = nextUrl
    videoInfo.value.name = d?.name || videoInfo.value.name
    if (d?.duration) videoInfo.value.duration = Number(d.duration)
    authData.value.hasFullAuth = !!d?.hasFullAuth
    const pct = Number(d?.trialPercent || 8); authData.value.trialPercent = pct > 0 ? pct : 8
    const urlChanged = nextUrl !== prevUrl
    setTimeout(() => {
      if (urlChanged) initHlsIfNeeded()
      setTimeout(() => {
        if (urlChanged && resumeTarget > 0) seekVideo(Math.max(0, resumeTarget - 1))
        if (shouldResume && !isTrialLocked.value) { isPlaying.value = true; safePlay() }
      }, 350)
    }, 100)
    return true
  } finally { isRefreshingSignedUrl.value = false }
}

const refreshSignedUrlIfNeeded = (curSecs) => {
  const dl = signedUrlDeadline(videoInfo.value.videoUrl); if (!dl) return
  const now = Math.floor(Date.now() / 1000); if (dl - now > 45) return
  if (isRefreshingSignedUrl.value || now - lastSignedRefreshAt < 10) return
  lastSignedRefreshAt = now; refreshSignedEpisodeUrl(curSecs)
}

// Data loading
const loadEpisode = async () => {
  const res = await findVideoEpisode(episodeId.value)
  if (res.code !== 0 || !res.data) { showToast(t('player.load_failed')); return }
  const d = res.data
  videoInfo.value = { id: getEntityId(d), name: d?.name || '', videoUrl: d?.videoUrl || '', duration: Number(d?.duration || 0) }
  authData.value.hasFullAuth = !!d?.hasFullAuth
  const pct = Number(d?.trialPercent || 8); authData.value.trialPercent = pct > 0 ? pct : 8
}

const loadSubtitleList = async () => {
  const res = await getSentenceList(episodeId.value)
  if (res.code !== 0) { subtitleList.value = []; return }
  const list = Array.isArray(res.data) ? res.data : []
  const limit = maxTrialAllowedTime()
  subtitleList.value = list.map((item, index) => {
    const st = Number(item?.startTime || 0)
    return { ...item, id: getEntityId(item) || index + 1, startTime: st, endTime: Number(item?.endTime || 0), english: String(item?.english || ''), translate: item?.translate || '', _locked: !authData.value.hasFullAuth && st >= limit }
  })
  if (subtitleList.value.length === 0) showToast(t('player.no_subtitle'))
}

const loadSentenceCollectionState = async () => {
  const ids = subtitleList.value.map(i => Number(i?.id || 0)).filter(id => id > 0)
  if (ids.length === 0) { sentenceCollectedMap.value = {}; return }
  const idSet = new Set(ids); const map = {}
  try {
    for (let page = 1; page <= 20; page++) {
      const res = await getCollectionList({ targetType: 2, page, pageSize: 100 })
      if (res.code !== 0 || !res.data) break
      const list = Array.isArray(res.data.list) ? res.data.list : []
      for (const item of list) { const tid = Number(item?.targetId || 0); if (tid > 0 && idSet.has(tid)) map[tid] = true }
      if (Object.keys(map).length >= idSet.size) break
      if (list.length < 100) break
    }
    sentenceCollectedMap.value = map
  } catch (e) { sentenceCollectedMap.value = {} }
}

const loadEpisodeCollectionState = async () => {
  try {
    const res = await getCollectionList({ targetType: 3, pageSize: 100 })
    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data.list) ? res.data.list : []
      isEpisodeCollected.value = list.some(item => Number(item.targetId) === Number(episodeId.value))
    }
  } catch (e) {}
}

const loadWatchHistory = async () => {
  const res = await getWatchProgress(episodeId.value)
  if (res.code !== 0 || !res.data) return
  const progress = Number(res.data.progressSecs || 0)
  if (progress > 0) seekVideo(progress)
}

const focusSentence = () => {
  if (!sentenceId.value || subtitleList.value.length === 0) return false
  const idx = subtitleList.value.findIndex(i => getEntityId(i) === sentenceId.value)
  if (idx < 0) return false
  activeSentenceIndex.value = idx; focusedSentenceId.value = sentenceId.value
  if (focusSentenceTimer) clearTimeout(focusSentenceTimer)
  focusSentenceTimer = setTimeout(() => { focusedSentenceId.value = 0 }, 2500)
  seekVideo(subtitleList.value[idx]?.startTime || 0)
  showToast(t('player.sentence_focused'))
  return true
}

const reAnchorScroll = () => {
  nextTick(() => { if (activeEl) activeEl.scrollIntoView({ behavior: 'smooth', block: 'center' }) })
}

const initPage = async () => {
  if (!episodeId.value) { showToast(t('player.load_failed')); return }
  autoPlayPending.value = true; isVideoBuffering.value = true; isVideoReady.value = false; startLoadingTimeout(); lastPlaybackTime = 0
  await loadEpisode()
  await loadSubtitleList()
  await loadSentenceCollectionState()
  await loadEpisodeCollectionState()
  await nextTick()
  initHlsIfNeeded()
  const focused = focusSentence()
  if (!focused) {
    if (sentenceId.value > 0) showToast(t('player.sentence_not_found'))
    await loadWatchHistory()
  } else {
    setTimeout(() => reAnchorScroll(), 350)
  }
}

// Highlight parsing
const parseHighlight = (htmlStr) => {
  if (!htmlStr) return [{ text: '', isHighlight: false }]
  const parts = []; const regex = /<w id="(\d+)">(.*?)<\/w>/g; let lastIdx = 0; let match
  while ((match = regex.exec(htmlStr)) !== null) {
    if (match.index > lastIdx) parts.push({ text: htmlStr.slice(lastIdx, match.index), isHighlight: false })
    parts.push({ text: match[2], isHighlight: true, wordId: match[1] })
    lastIdx = regex.lastIndex
  }
  if (lastIdx < htmlStr.length) parts.push({ text: htmlStr.slice(lastIdx), isHighlight: false })
  return parts.length ? parts : [{ text: htmlStr, isHighlight: false }]
}

// Word popup
const showWordDetail = async (wordId) => {
  if (isPlaying.value && videoRef.value) videoRef.value.pause()
  try {
    const res = await findWord(Number(wordId))
    if (res.code !== 0 || !res.data) return
    currentWord.value = { ...res.data, id: getEntityId(res.data) }
    isWordCollected.value = false
    try {
      const cr = await getCollectionList({ targetType: 1, pageSize: 100 })
      if (cr.code === 0 && cr.data) {
        const list = Array.isArray(cr.data.list) ? cr.data.list : []
        isWordCollected.value = list.some(item => Number(item.targetId) === getEntityId(res.data))
      }
    } catch (e) {}
    wordPopup.value = { visible: true, word: currentWord.value }
  } catch (e) {}
}

const handleSegClick = (seg, event) => {
  if (seg && seg.isHighlight) {
    if (event?.stopPropagation) event.stopPropagation()
    showWordDetail(seg.wordId)
  }
}

const closeWordPopup = () => { stopLocalTTS(); isWordAudioPlaying.value = false; wordPopup.value.visible = false }

const playWordAudioAction = () => {
  if (!currentWord.value) return
  if (isWordAudioPlaying.value) { stopLocalTTS(); isWordAudioPlaying.value = false; return }
  playWordTTS(currentWord.value.word, currentWord.value.audioUs || '', currentWord.value.audioUk || '', 'US', {
    onStart: () => { isWordAudioPlaying.value = true },
    onEnd: () => { isWordAudioPlaying.value = false },
    onError: () => { isWordAudioPlaying.value = false }
  })
}

const collectCurrentWord = async () => {
  const id = getEntityId(currentWord.value); if (!id) return
  const collected = isWordCollected.value
  const res = collected ? await uncollect(1, id) : await collect(1, id)
  if (res.code === 0) {
    isWordCollected.value = !collected
    showToast(!collected ? t('player.collect_success') : t('player.uncollect_success'))
  } else { showToast(resolveApiMessage(res.msg, 'player.operationFailed')) }
}

// Heartbeat
const enqueueHeartbeatSecs = (secs) => {
  const d = Math.max(0, Math.floor(Number(secs || 0)))
  if (d > 0) queuedHeartbeatSecs = Math.min(HEARTBEAT_QUEUE_MAX, queuedHeartbeatSecs + d)
}

const sendHeartbeatQueue = async (progressSecs) => {
  if (!episodeId.value || heartbeatSending || queuedHeartbeatSecs <= 0) return
  heartbeatSending = true
  try {
    while (queuedHeartbeatSecs > 0) {
      const usingTimeSecs = Math.min(HEARTBEAT_CHUNK_MAX, queuedHeartbeatSecs)
      const res = await heartbeat({ episodeId: episodeId.value, progressSecs, usingTimeSecs })
      if (!res || res.code !== 0) break
      queuedHeartbeatSecs -= usingTimeSecs
    }
  } finally { heartbeatSending = false }
}

const flushHeartbeat = (progressSecs, force = false) => {
  const usingTimeSecs = Math.floor(accumulatedViewTime)
  if (usingTimeSecs > 0) { enqueueHeartbeatSecs(usingTimeSecs); accumulatedViewTime = 0 }
  if (!force && queuedHeartbeatSecs < HEARTBEAT_FLUSH_THRESHOLD) return
  void sendHeartbeatQueue(progressSecs)
}

// Video events
const onTimeUpdate = () => {
  const v = videoRef.value; if (!v) return
  const cur = v.currentTime; const dur = v.duration || videoInfo.value.duration
  if (!isDragging.value) currentTime.value = cur
  if (cur > lastPlaybackTime + 0.05) { isVideoBuffering.value = false; isVideoReady.value = true }
  lastPlaybackTime = cur
  if (dur > 0) duration.value = dur
  if (ensureTrialLimit(cur)) return
  if (!isRepeatMode.value) {
    const idx = subtitleList.value.findIndex(s => cur >= s.startTime && cur <= s.endTime)
    if (idx !== -1 && idx !== activeSentenceIndex.value) { activeSentenceIndex.value = idx; nextTick(() => scrollToActiveSentence()) }
  }
  if (isRepeatMode.value) {
    const ri = subtitleList.value.findIndex(s => cur >= s.startTime && cur <= s.endTime)
    if (ri !== -1) {
      const cs = subtitleList.value[ri]
      if (cs && cur >= cs.endTime - 0.3) {
        const now = Date.now()
        if (now - lastRepeatSeekTime > 600) { lastRepeatSeekTime = now; seekVideo(cs.startTime) }
      }
    }
  }
  refreshSignedUrlIfNeeded(cur)
  if (lastHeartbeatTime > 0 && cur > lastHeartbeatTime) {
    accumulatedViewTime += (cur - lastHeartbeatTime)
    if (accumulatedViewTime >= HEARTBEAT_FLUSH_THRESHOLD) flushHeartbeat(cur)
  }
  lastHeartbeatTime = cur
}

const onEnded = () => {
  flushHeartbeat(duration.value || lastHeartbeatTime, true)
  if (controls.value.loop) { seekVideo(0); safePlay() } else { isPlaying.value = false }
}
const onLoadedMetadata = () => { const d = videoRef.value?.duration || 0; if (d > 0) duration.value = d; isPlaying.value = false; isVideoBuffering.value = true; isVideoReady.value = false; startLoadingTimeout() }
const onCanPlay = () => { clearLoadingTimeout(); isVideoBuffering.value = false; isVideoReady.value = true; if (autoPlayPending.value && !isTrialLocked.value) safePlay() }
const onSeeking = () => { isVideoBuffering.value = true; startLoadingTimeout() }
const onSeeked = () => { clearLoadingTimeout(); if (!isTrialLocked.value) safePlay() }
const onWaiting = () => { isVideoBuffering.value = true; isVideoReady.value = false; startLoadingTimeout() }
const onPlaying = () => { clearLoadingTimeout(); isVideoBuffering.value = false; isVideoReady.value = true; isPlaying.value = true; autoPlayPending.value = false }
const onPlay = () => { clearLoadingTimeout(); isVideoBuffering.value = false; isPlaying.value = true; isVideoReady.value = true; autoPlayPending.value = false }
const onPause = () => { isPlaying.value = false }
const onVideoError = async () => {
  isPlaying.value = false; isVideoBuffering.value = true; isVideoReady.value = false; autoPlayPending.value = true; clearLoadingTimeout()
  const ok = await refreshSignedEpisodeUrl(lastPlaybackTime); if (!ok) showToast(t('player.load_failed'))
}

// UI actions
const togglePlay = () => {
  showPauseHint.value = false; const v = videoRef.value; if (!v) return
  if (!authData.value.hasFullAuth && isTrialLocked.value) return
  if (showVideoLoading.value) { autoPlayPending.value = true; safePlay(); return }
  isPlaying.value ? v.pause() : safePlay()
}
const toggleBilingual = () => { controls.value.showBilingual = !controls.value.showBilingual; reAnchorScroll() }
const toggleSubtitles = () => { controls.value.showSubtitles = !controls.value.showSubtitles; reAnchorScroll() }
const toggleLoop = () => { controls.value.loop = !controls.value.loop }
const toggleRepeat = () => { isRepeatMode.value = !isRepeatMode.value; if (isRepeatMode.value) showToast(t('player.repeat_on')) }
const selectSpeed = (rate) => { controls.value.speed = rate; const v = videoRef.value; if (v) v.playbackRate = rate }

const jumpBySubtitle = (item, index) => {
  showPauseHint.value = false
  const st = Number(item?.startTime || 0)
  activeSentenceIndex.value = index; seekVideo(st)
  try { heartbeat({ episodeId: episodeId.value, progressSecs: st, usingTimeSecs: 0 }).catch(() => {}) } catch (e) {}
  safePlay()
}

const collectSentence = async (item) => {
  const id = getEntityId(item); if (!id) return
  const col = !!sentenceCollectedMap.value[id]
  const res = col ? await uncollect(2, id) : await collect(2, id)
  if (res.code === 0) { sentenceCollectedMap.value = { ...sentenceCollectedMap.value, [id]: !col }; showToast(!col ? t('player.collect_success') : t('player.uncollect_success')) }
  else showToast(resolveApiMessage(res.msg, 'player.operationFailed'))
}

const toggleEpisodeCollect = async () => {
  const id = episodeId.value; const col = isEpisodeCollected.value
  const res = col ? await uncollect(3, id) : await collect(3, id)
  if (res.code === 0) { isEpisodeCollected.value = !col; showToast(!col ? t('player.collect_success') : t('player.uncollect_success')) }
  else showToast(resolveApiMessage(res.msg, 'player.operationFailed'))
}

const handleBack = () => { showSettings.value = false; if (window.history.length > 1) router.back(); else router.push('/home') }

// Progress drag
const updateProgressFromEvent = (e) => {
  const bar = document.querySelector('.progress-bar'); if (!bar) return
  const rect = bar.getBoundingClientRect()
  const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
  currentTime.value = pct * duration.value
}

const onProgressMouseDown = (e) => {
  isDragging.value = true; wasPlayingBeforeDrag = isPlaying.value; const v = videoRef.value; if (v) v.pause()
  document.addEventListener('mousemove', onProgressDrag); document.addEventListener('mouseup', onProgressDragEnd)
  updateProgressFromEvent(e)
}
const onProgressDrag = (e) => { if (isDragging.value) updateProgressFromEvent(e) }
const onProgressDragEnd = (e) => {
  if (!isDragging.value) return
  isDragging.value = false; document.removeEventListener('mousemove', onProgressDrag); document.removeEventListener('mouseup', onProgressDragEnd)
  const bar = document.querySelector('.progress-bar')
  if (bar) {
    const rect = bar.getBoundingClientRect()
    const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
    seekVideo(pct * duration.value); if (wasPlayingBeforeDrag && !isTrialLocked.value) safePlay()
  }
}

const progressPercent = computed(() => duration.value <= 0 ? 0 : (currentTime.value / duration.value) * 100)

const onSubtitleMouseDown = () => { const v = videoRef.value; if (v && !v.paused) { v.pause(); showPauseHint.value = true } }
const hidePauseHint = () => { showPauseHint.value = false }

const onMouseMove = () => {
  showControls.value = true; if (controlsTimer) clearTimeout(controlsTimer)
  controlsTimer = setTimeout(() => { if (isPlaying.value) showControls.value = false }, 3000)
}

const scrollToActiveSentence = () => {
  const c = subtitleScrollRef.value; if (!c || !activeEl) return
  const cr = c.getBoundingClientRect(); const ar = activeEl.getBoundingClientRect()
  c.scrollBy({ top: ar.top - cr.top - c.clientHeight / 2 + ar.height / 2, behavior: 'smooth' })
}

const showToast = (msg) => {
  const t = document.createElement('div')
  t.textContent = msg
  t.style.cssText = 'position:fixed;top:20px;left:50%;transform:translateX(-50%);background:rgba(26,27,58,0.85);color:#fff;padding:10px 24px;border-radius:999px;z-index:9999;font-size:14px;pointer-events:none;'
  document.body.appendChild(t); setTimeout(() => t.remove(), 2500)
}

watch(() => videoInfo.value.videoUrl, () => { nextTick(() => initHlsIfNeeded()) })

onMounted(() => { initPage() })

onBeforeUnmount(() => {
  flushHeartbeat(lastHeartbeatTime, true); clearLoadingTimeout(); destroyHls(); destroyTTS()
  if (focusSentenceTimer) clearTimeout(focusSentenceTimer)
  if (controlsTimer) clearTimeout(controlsTimer)
  document.removeEventListener('mousemove', onProgressDrag); document.removeEventListener('mouseup', onProgressDragEnd)
})
</script>

<style lang="scss" scoped>
.player-page {
  display: flex; flex-direction: column;
  height: calc(100vh - #{$header-height} - 48px);
  max-width: $content-max-width; margin: 0 auto;
}
.player-header {
  display: flex; align-items: center; padding: $spacing-md 0; margin-bottom: $spacing-md; flex-shrink: 0;
  .back-btn {
    display: flex; align-items: center; gap: 6px; padding: 8px 16px;
    background: $bg-card; border-radius: $radius-pill; color: $text-secondary;
    font-size: 14px; transition: all 0.2s ease; border: none; cursor: pointer;
    &:hover { background: $bg-hover; color: $primary-color; }
    svg { width: 16px; height: 16px; }
  }
  .video-title { flex: 1; margin-left: 16px; font-size: 18px; font-weight: $font-weight-semibold; color: $text-primary; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .header-actions { display: flex; gap: 8px;
    .action-btn {
      width: 40px; height: 40px; display: flex; align-items: center; justify-content: center;
      background: $bg-card; border-radius: 50%; color: $text-secondary; transition: all 0.2s ease; border: none; cursor: pointer;
      &:hover, &.active { color: $primary-color; background: $bg-active; }
      svg { width: 20px; height: 20px; }
    }
  }
}
.player-body { flex: 1; display: flex; gap: $spacing-lg; min-height: 0; }
.video-section { flex: 1; min-width: 0; }
.video-wrapper { position: relative; width: 100%; height: 100%; background: #000; border-radius: $radius-lg; overflow: hidden; }
.video-player { width: 100%; height: 100%; object-fit: contain; background: #000; }

.video-loading-mask { position: absolute; inset: 0; background: rgba(26,27,58,0.72); display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 14px; pointer-events: none; z-index: 5; }
.loading-logo { width: 48px; height: 48px; border-radius: 12px; background: $gradient-primary; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 20px; font-weight: 700; }
.loading-text { color: #e2e8f0; font-size: 14px; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
.spin { animation: spin 1.5s linear infinite; }

.pause-hint-overlay { position: absolute; inset: 0; background: rgba(26,27,58,0.55); display: flex; align-items: center; justify-content: center; z-index: 10; cursor: pointer; }
.pause-hint-text { color: #fff; font-size: 16px; background: rgba(109,91,255,0.9); padding: 10px 24px; border-radius: 999px; font-weight: 600; }

.trial-mask { position: absolute; inset: 0; background: rgba(0,0,0,0.8); display: flex; align-items: center; justify-content: center; z-index: 10;
  .trial-card { text-align: center; color: #fff; padding: 40px;
    svg { width: 56px; height: 56px; margin-bottom: 16px; color: $primary-light; }
    h3 { font-size: 22px; font-weight: $font-weight-bold; margin-bottom: 8px; }
    p { font-size: 14px; opacity: 0.7; margin-bottom: 24px; }
  }
}
.btn-primary { padding: 10px 32px; border: none; border-radius: $radius-pill; background: $gradient-primary; color: #fff; font-size: 15px; font-weight: 600; cursor: pointer; box-shadow: 0 4px 14px rgba(109,91,255,0.35); }

.video-controls {
  position: absolute; left: 0; right: 0; bottom: 0; padding: 40px 20px 16px;
  background: linear-gradient(transparent, rgba(0,0,0,0.7)); color: #fff; transition: opacity 0.3s ease; z-index: 4;
  &.hidden { opacity: 0; pointer-events: none; }
  .progress-bar {
    position: relative; height: 4px; background: rgba(255,255,255,0.3); border-radius: 2px; cursor: pointer; margin-bottom: 12px; transition: height 0.15s;
    &:hover { height: 6px; }
    .progress-played { position: absolute; left: 0; top: 0; height: 100%; background: $gradient-primary; border-radius: 2px;
      .progress-thumb { position: absolute; right: -6px; top: 50%; transform: translateY(-50%); width: 12px; height: 12px; background: #fff; border-radius: 50%; box-shadow: 0 2px 6px rgba(0,0,0,0.3); opacity: 0; transition: opacity 0.2s;
        &.dragging, .progress-bar:hover & { opacity: 1; }
      }
    }
    .trial-marker { position: absolute; top: -2px; width: 2px; height: 8px; background: #f59e0b;
      .trial-label { position: absolute; top: -22px; left: 50%; transform: translateX(-50%); font-size: 10px; white-space: nowrap; color: #f59e0b; background: rgba(0,0,0,0.6); padding: 2px 6px; border-radius: 4px; }
    }
  }
  .controls-row { display: flex; align-items: center; justify-content: space-between; }
  .controls-left, .controls-center { display: flex; align-items: center; gap: 8px; }
  .ctrl-btn {
    width: 36px; height: 36px; display: flex; align-items: center; justify-content: center;
    color: #fff; border-radius: 50%; transition: all 0.2s ease; border: none; background: none; cursor: pointer;
    &:hover { background: rgba(255,255,255,0.15); }
    &.active { color: $primary-light; }
    svg { width: 20px; height: 20px; }
  }
  .speed-btn { width: auto; padding: 0 12px; border-radius: $radius-pill; font-size: 13px; font-weight: $font-weight-medium; }
  .time-text { font-size: 13px; font-family: monospace; margin-left: 8px; }
  .speed-selector { position: relative; }
  .speed-menu {
    position: absolute; bottom: calc(100% + 8px); left: 50%; transform: translateX(-50%);
    background: rgba(0,0,0,0.85); border-radius: $radius-md; padding: 8px 0; min-width: 140px; backdrop-filter: blur(10px);
    .menu-title { padding: 4px 16px 8px; font-size: 12px; color: rgba(255,255,255,0.5); }
    .speed-item { padding: 8px 16px; font-size: 13px; cursor: pointer; display: flex; align-items: center; justify-content: space-between; gap: 8px; transition: background 0.15s;
      &:hover { background: rgba(255,255,255,0.1); }
      &.active { color: $primary-light; }
      .check { color: $primary-light; }
    }
    .speed-divider { height: 1px; background: rgba(255,255,255,0.1); margin: 6px 0; }
  }
}

.subtitle-section {
  width: 380px; flex-shrink: 0; background: $bg-card; border-radius: $radius-lg;
  display: flex; flex-direction: column; box-shadow: $shadow-card; border: 1px solid $border-light; overflow: hidden;
}
.subtitle-header { display: flex; align-items: center; justify-content: space-between; padding: 16px 20px; border-bottom: 1px solid $border-light; flex-shrink: 0;
  h3 { font-size: 16px; font-weight: $font-weight-semibold; color: $text-primary; margin: 0; }
}
.subtitle-scroll { flex: 1; overflow-y: auto; background: $bg-page; }
.subtitle-padding-top { height: 12px; }
.subtitle-padding-bottom { height: 16px; }

.sentence-row {
  padding: 14px 18px; margin: 0 12px 8px; transition: all 0.25s ease; opacity: 0.55;
  border-radius: $radius-md; background: #fff; border: 1px solid rgba(109,91,255,0.08); cursor: pointer;
  &:hover { opacity: 0.8; background: rgba(109,91,255,0.03); }
  &.is-active {
    opacity: 1; background: linear-gradient(135deg, rgba(109,91,255,0.06) 0%, rgba(155,143,255,0.04) 100%);
    border-left: 4px solid $primary-color; border: 1.5px solid rgba(109,91,255,0.28); box-shadow: 0 4px 12px rgba(109,91,255,0.12);
  }
  &.is-focused { opacity: 1; border-left: 4px solid $primary-light; box-shadow: 0 4px 12px rgba(155,143,255,0.15); }
  .sentence-meta { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
  .time-badge { font-size: 11px; color: $primary-color; background: rgba(109,91,255,0.1); padding: 2px 10px; border-radius: 999px; font-variant-numeric: tabular-nums; font-weight: 600; }
  .collect-btn { font-size: 14px; color: $primary-color; padding: 2px 8px; border-radius: 999px; cursor: pointer; transition: all 0.2s; &:hover { background: rgba(109,91,255,0.1); } }
  .english-text { font-size: 14px; font-weight: 500; line-height: 1.6; color: $text-primary; }
  .highlight-word {
    color: $primary-color; font-weight: 700; border-bottom: 2px dashed rgba(109,91,255,0.4);
    display: inline; padding: 0 2px; margin: 0 2px; cursor: pointer;
    &:hover { background: rgba(109,91,255,0.1); }
  }
  .translated-text { font-size: 13px; color: $text-secondary; line-height: 1.6; margin-top: 8px; padding-top: 8px; border-top: 1px solid rgba(109,91,255,0.08); }
}
.sentence-locked { display: flex; align-items: center; justify-content: center; gap: 8px; padding: 12px 8px;
  .lock-time { font-size: 12px; color: $text-tertiary; font-variant-numeric: tabular-nums; }
  .lock-icon { font-size: 18px; opacity: 0.5; }
}

.word-popup-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.3); display: flex; align-items: center; justify-content: center; z-index: $z-modal; }
.word-detail-box { width: 400px; max-width: 90vw; padding: 24px 28px; background: #fff; border-radius: $radius-lg; box-shadow: $shadow-xl; }
.w-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 4px; }
.w-title { font-size: 24px; font-weight: 700; color: $text-primary; }
.w-audio-circle {
  width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center;
  background: $gradient-primary; border: none; cursor: pointer; color: #fff; transition: transform 0.15s;
  &:hover { transform: scale(1.05); }
  &.playing { animation: pulse 0.6s ease-in-out infinite; pointer-events: none; }
}
@keyframes pulse { 0%, 100% { transform: scale(1); opacity: 1; } 50% { transform: scale(1.12); opacity: 0.75; } }
.w-audio-icon { font-size: 16px; }
.w-phonetic { font-size: 14px; color: $text-secondary; display: block; margin: 4px 0 12px; }
.w-divider { height: 1px; background: rgba(109,91,255,0.16); margin-bottom: 16px; }
.w-exp { font-size: 15px; color: $text-primary; line-height: 1.6; margin-bottom: 20px; }
.w-actions { display: flex; gap: 12px; }
.w-cancel-btn {
  flex: 1; height: 44px; border-radius: $radius-pill; border: 1.5px solid rgba(109,91,255,0.32);
  background: #fff; color: $primary-color; font-size: 15px; font-weight: 600; cursor: pointer;
}
.w-collect {
  flex: 1; height: 44px; border: none; background: $gradient-primary; color: #fff;
  border-radius: $radius-pill; font-size: 15px; font-weight: 600; cursor: pointer; box-shadow: 0 4px 14px rgba(109,91,255,0.3);
}
</style>
