<template>
  <view class="nf-player">

    <!-- ===== 全屏播放器弹层 ===== -->
    <view class="nf-theater" :class="{ 'nf-theater-show': showTheater }" @tap.self="closeTheater">
      <view class="nf-theater-inner" :class="{ 'nf-theater-fade-in': theaterReady }">
        <!-- 关闭按钮 -->
        <view class="nf-theater-close" @tap="closeTheater">
          <uni-icons type="closeempty" size="22" color="#fff" />
        </view>
        <!-- 正在播放标签 -->
        <view class="nf-theater-badge">
          <view class="nf-theater-dot"></view>
          <text class="nf-theater-badge-text">{{ $t('playerNowPlaying') }}</text>
        </view>
        <!-- 视频播放器 -->
        <view class="nf-theater-video-wrap">
          <!-- 自定义视频播放器 -->
          <div v-if="videoUrl && showTheater" ref="playerRef" class="nf-vp">
            <video ref="videoRef" class="nf-vp-video"
              :src="videoUrl" :poster="getUrl(data.imageUrl)"
              playsinline webkit-playsinline x5-video-player-type="h5" preload="auto"
              :style="brtStyle"
              @timeupdate="vpTimeUpdate" @play="isPlaying = true" @pause="isPlaying = false"
              @ended="vpEnded"
            ></video>
            <!-- 手势层 -->
            <div class="nf-vp-touch" @touchstart.prevent="vpTouchStart" @touchmove.prevent="vpTouchMove" @touchend.prevent="vpTouchEnd"></div>
            <!-- 提示 -->
            <div class="nf-vp-tip" v-if="vpTip">{{ vpTip }}</div>
            <!-- 加载中 -->
            <div class="nf-vp-loading" v-if="isBuffering"><div class="nf-vp-spinner"></div></div>
            <!-- 控制栏 -->
            <div class="nf-vp-ctrl" :class="{ 'nf-vp-ctrl-show': vpCtrlVis }">
              <div class="nf-vp-prog" ref="vpProgRef" @click.stop="vpProgClick">
                <div class="nf-vp-prog-bar"><div class="nf-vp-prog-fill" :style="{ width: vpProg + '%' }"></div></div>
                <div class="nf-vp-prog-dot" :style="{ left: vpProg + '%' }"></div>
              </div>
              <div class="nf-vp-btns">
                <div class="nf-vp-btn" @click.stop="vpTogglePlay">{{ isPlaying ? '⏸' : '▶' }}</div>
                <span class="nf-vp-time">{{ fmtTime(currentTime) }} / {{ fmtTime(duration) }}</span>
                <div class="nf-vp-spacer"></div>
                <div class="nf-vp-btn nf-vp-spd" @click.stop="vpCycleSpeed">{{ playbackSpeed === 1 ? $t('playerSpeed') : playbackSpeed + 'x' }}</div>
                <div class="nf-vp-btn nf-vp-fs" @click.stop="vpToggleFs">{{ vpFs ? $t('playerExitFs') : $t('playerFullscreen') }}</div>
                <div class="nf-vp-btn nf-vp-ori" v-if="vpFs" @click.stop="vpToggleOri">{{ vpLand ? $t('playerPortrait') : $t('playerLandscape') }}</div>
              </div>
            </div>
          </div>
          <!-- 无视频时占位 -->
          <view v-else-if="showTheater" class="nf-theater-poster-wrap">
            <image class="nf-theater-poster" :src="getUrl(data.imageUrl)" mode="aspectFill" />
            <view class="nf-theater-no-video">
              <text class="nf-theater-no-video-text">{{ $t('playerNoVideo') }}</text>
            </view>
          </view>
        </view>
        <!-- 当前集信息 -->
        <view class="nf-theater-info">
          <text class="nf-theater-title">{{ currentEpisode.name || $lt(data.title) }}</text>
          <text class="nf-theater-sub" v-if="currentEpisode.description">{{ currentEpisode.description }}</text>
        </view>
        <!-- 集数快选 -->
        <view class="nf-theater-eps" v-if="episodes.length > 1">
          <scroll-view scroll-x :show-scrollbar="false" class="nf-theater-eps-scroll">
            <view class="nf-theater-eps-list">
              <view
                v-for="(ep, idx) in episodes"
                :key="ep.ID"
                class="nf-theater-ep-btn"
                :class="{ 'nf-theater-ep-active': currentIndex === idx }"
                @tap="switchEpisode(idx)"
              >
                <text class="nf-theater-ep-text">{{ idx + 1 }}</text>
              </view>
            </view>
          </scroll-view>
        </view>
      </view>
    </view>

    <!-- ===== 主页面内容 ===== -->
    <!-- 封面 Hero -->
    <view class="nf-hero">
      <view class="nf-hero-status"></view>
      <image class="nf-hero-img" :src="getUrl(data.imageUrl)" mode="aspectFill" @tap="playFirst" />
      <view class="nf-hero-gradient"></view>
      <!-- 返回按钮 -->
      <view class="nf-hero-back" @tap="goBack">
        <uni-icons type="left" size="20" color="#fff" />
      </view>
    </view>

    <!-- 内容区域 -->
    <scroll-view scroll-y :show-scrollbar="false" class="nf-content-scroll">
      <!-- 标题信息 -->
      <view class="nf-info-section">
        <text class="nf-title">{{ $lt(data.title) }}</text>
        <view class="nf-meta-row">
          <view class="nf-meta-item" v-if="data.rating">
            <text class="nf-meta-star">★</text>
            <text class="nf-meta-value">{{ data.rating }}</text>
            <text class="nf-meta-label">{{ $t('playerRating') }}</text>
          </view>
          <view class="nf-meta-divider" v-if="data.rating && data.view_num"></view>
          <view class="nf-meta-item" v-if="data.view_num">
            <text class="nf-meta-value">{{ formatNum(data.view_num) }}</text>
            <text class="nf-meta-label">{{ $t('playerViews') }}</text>
          </view>
          <view class="nf-meta-divider" v-if="episodes.length"></view>
          <view class="nf-meta-item" v-if="episodes.length">
            <text class="nf-meta-value">{{ episodes.length }}</text>
            <text class="nf-meta-label">{{ $t('playerEpisodes') }}</text>
          </view>
        </view>

        <!-- 操作栏 -->
        <view class="nf-actions">
          <view class="nf-action-btn" @tap="toggleCollect">
            <uni-icons
              :type="collectionFlag ? 'heart-filled' : 'heart'"
              size="22"
              :color="collectionFlag ? '#e50914' : 'rgba(255,255,255,0.6)'"
            />
            <text class="nf-action-text">{{ collectionFlag ? $t('playerCollected') : $t('playerCollect') }}</text>
          </view>
          <view class="nf-action-btn" @tap="shareVideo">
            <uni-icons type="redo" size="22" color="rgba(255,255,255,0.6)" />
            <text class="nf-action-text">{{ $t('playerShare') }}</text>
          </view>
        </view>
      </view>

      <!-- 简介 -->
      <view class="nf-desc-section" v-if="$lt(data.description)">
        <text class="nf-section-title">{{ $t('playerDesc') }}</text>
        <text class="nf-desc-text" :class="{ 'nf-desc-expand': descExpand }">{{ $lt(data.description) }}</text>
        <text class="nf-desc-toggle" @tap="descExpand = !descExpand">
          {{ descExpand ? $t('playerCollapse') : $t('playerExpand') }}
        </text>
      </view>

      <!-- 选集 -->
      <view class="nf-episodes-section" v-if="episodes.length">
        <view class="nf-section-header">
          <text class="nf-section-title">{{ $t('playerEpisodes') }}</text>
          <text class="nf-episode-count">{{ epCountLabel }}</text>
        </view>
        <scroll-view scroll-x :show-scrollbar="false" class="nf-episodes-scroll">
          <view class="nf-episodes-list">
            <view
              v-for="(ep, idx) in episodes"
              :key="ep.ID"
              class="nf-episode-card"
              :class="{ 'nf-episode-active': currentIndex === idx }"
              @tap="playEpisode(idx)"
            >
              <view class="nf-ep-thumb-wrap">
                <image
                  class="nf-ep-thumb"
                  :src="ep.picture ? getUrl(ep.picture) : getUrl(data.imageUrl)"
                  mode="aspectFill"
                />
                <view class="nf-ep-play-icon" v-if="currentIndex !== idx">
                  <text class="nf-ep-play-text">▶</text>
                </view>
                <view class="nf-ep-playing" v-else>
                  <view class="nf-eq-bar"></view>
                  <view class="nf-eq-bar"></view>
                  <view class="nf-eq-bar"></view>
                </view>
              </view>
              <view class="nf-ep-info">
                <text class="nf-ep-name">{{ ep.name || epLabel(idx + 1) }}</text>
                <text class="nf-ep-desc" v-if="ep.description">{{ ep.description }}</text>
                <text class="nf-ep-price" v-if="ep.price > 0">¥{{ (ep.price / 100).toFixed(2) }}</text>
                <text class="nf-ep-free" v-else>{{ $t('playerFree') }}</text>
              </view>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 无剧集时的占位 -->
      <view class="nf-no-episodes" v-if="!episodes.length && !isLoading">
        <text class="nf-no-episodes-text">{{ $t('playerNoEpisodes') }}</text>
      </view>

      <!-- 底部安全区 -->
      <view style="height: 60rpx;"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed, nextTick, onUnmounted, watch } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { findGood } from '@/api/product.js'
import { findCollect, createCollect } from '@/api/collect.js'
import { getUrl } from '@/utils/url.js'
import { useUserStore } from '@/pinia/modules/user'
import { useLangStore } from '@/pinia/modules/lang.js'
import { usePlayHistoryStore } from '@/pinia/modules/playHistory.js'


const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const userStore = useUserStore()
const token = userStore.token || ''
const playHistoryStore = usePlayHistoryStore()

const data = ref({})
const episodes = ref([])
const currentIndex = ref(-1)
const currentEpisode = ref({})
const videoUrl = ref('')
const collectionFlag = ref(false)
const isLoading = ref(true)
const descExpand = ref(false)
const goodID = ref(0)
const showTheater = ref(false)
const theaterReady = ref(false)
const isPlaying = ref(false)
const currentTime = ref(0)
const duration = ref(0)
const playbackSpeed = ref(1.0)

const epLabel = (n) => {
  return $t.value('playerEp').replace('{{n}}', n)
}

const epCountLabel = computed(() => {
  return $t.value('playerEpCount').replace('{{n}}', episodes.value.length)
})

const formatNum = (num) => {
  if (!num) return '0'
  if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
  return String(num)
}

onLoad(async (options) => {
  if (options.id) {
    goodID.value = options.id
    await init()
    if (options.autoplay) {
      playFirst()
    }
  }
})

const init = async () => {
  isLoading.value = true
  try {
    const res = await findGood(goodID.value)
    if (res.code === 0 && res.data.regood) {
      data.value = res.data.regood
      episodes.value = res.data.regood.skus || []
    }
    if (token) {
      const status = await findCollect({ goodID: goodID.value })
      if (status.code === 0) collectionFlag.value = status.data
    }
  } catch (e) {
    console.error($t.value('playerLoadFail'), e)
  } finally {
    isLoading.value = false
  }
}

const openTheater = () => {
  showTheater.value = true
  nextTick(() => {
    setTimeout(() => {
      theaterReady.value = true
      initPlayer()
    }, 80)
  })
}

const closeTheater = () => {
  // 关闭前立即保存进度
  if (goodID.value && currentTime.value > 5) {
    playHistoryStore.saveProgress(goodID.value, {
      episodeIndex: currentIndex.value,
      currentTime: currentTime.value,
      duration: duration.value,
      episodeName: currentEpisode.value?.name || epLabel(currentIndex.value + 1),
      imageUrl: data.value.imageUrl,
      title: data.value.title
    })
  }
  theaterReady.value = false
  isPlaying.value = false
  destroyPlayer()
  setTimeout(() => {
    showTheater.value = false
    videoUrl.value = ''
    currentTime.value = 0
    duration.value = 0
  }, 300)
}

// === 自定义播放器 ===
const videoRef = ref(null)
const playerRef = ref(null)
const vpProgRef = ref(null)
const isBuffering = ref(false)
const vpFs = ref(false)
const vpLand = ref(true)
const vpCtrlVis = ref(true)
const vpTip = ref('')
const brightnessVal = ref(1)
const usingCssRotation = ref(false)
const brtStyle = computed(() => ({ filter: `brightness(${brightnessVal.value})` }))
const vpProg = computed(() => duration.value > 0 ? (currentTime.value / duration.value * 100) : 0)

// 检测是否支持原生屏幕方向锁定（Android 支持，iOS 不支持）
const canNativeLock = () => {
  try {
    return typeof screen !== 'undefined' && screen.orientation && typeof screen.orientation.lock === 'function'
  } catch (e) {
    return false
  }
}

// 获取原生 video DOM 元素（uni-app H5 ref 可能是组件实例）
const getVideo = () => {
  const v = videoRef.value
  if (!v) return null
  if (v instanceof HTMLVideoElement) return v
  if (v.$el instanceof HTMLVideoElement) return v.$el
  return v.$el?.querySelector?.('video') || null
}

let ctrlTimer = null
let longTimer = null
let tipTimer = null
let tStartX = 0, tStartY = 0, tStartTime = 0
let tDir = '', tSeeking = false, seekTarget = 0, seekStart = 0
let speedBoosting = false, prevRate = 1
const SPEEDS = [1, 1.25, 1.5, 2, 0.5, 0.75]

const fmtTime = (s) => {
  if (!s || isNaN(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return m + ':' + (sec < 10 ? '0' : '') + sec
}

const showTip = (text, ms = 800) => {
  vpTip.value = text
  clearTimeout(tipTimer)
  tipTimer = setTimeout(() => { vpTip.value = '' }, ms)
}

const showCtrlBriefly = () => {
  vpCtrlVis.value = true
  clearTimeout(ctrlTimer)
  if (isPlaying.value) {
    ctrlTimer = setTimeout(() => { vpCtrlVis.value = false }, 4000)
  }
}

watch(isPlaying, (val) => {
  if (!val) {
    vpCtrlVis.value = true
    clearTimeout(ctrlTimer)
  } else {
    showCtrlBriefly()
  }
})

const onWaiting = () => { isBuffering.value = true }
const onCanPlay = () => { isBuffering.value = false }

const initPlayer = () => {
  destroyPlayer()
  if (!videoUrl.value) return
  nextTick(() => {
    const v = getVideo()
    if (!v) return
    // 原生监听 buffering 事件（Vue 模板绑定在 uni-app 中可能不可靠）
    v.addEventListener('waiting', onWaiting)
    v.addEventListener('playing', onCanPlay)
    v.addEventListener('canplay', onCanPlay)
    v.addEventListener('loadeddata', onCanPlay)
    v.playbackRate = playbackSpeed.value
    v.play().catch(() => {})
    // 恢复历史进度
    if (pendingResumeTime > 0) {
      const targetTime = pendingResumeTime
      const epNum = currentIndex.value + 1
      pendingResumeTime = 0
      const doSeek = () => {
        v.removeEventListener('loadedmetadata', doSeek)
        v.removeEventListener('canplay', doSeek)
        v.currentTime = targetTime
        showTip(
          $t.value('playerContinueFrom')
            .replace('{{ep}}', epNum)
            .replace('{{time}}', fmtTime(targetTime)),
          3000
        )
      }
      if (v.readyState >= 1) {
        doSeek()
      } else {
        v.addEventListener('loadedmetadata', doSeek)
        v.addEventListener('canplay', doSeek)
      }
    }
    showCtrlBriefly()
    document.addEventListener('fullscreenchange', onFsChange)
    document.addEventListener('webkitfullscreenchange', onFsChange)
  })
}

const destroyPlayer = () => {
  clearTimeout(ctrlTimer)
  clearTimeout(longTimer)
  clearTimeout(tipTimer)
  clearTimeout(switchTimer)
  switchingEpisode = false
  seekingByClick = false
  // 移除原生 video 事件
  const v = getVideo()
  if (v) {
    v.removeEventListener('waiting', onWaiting)
    v.removeEventListener('playing', onCanPlay)
    v.removeEventListener('canplay', onCanPlay)
    v.removeEventListener('loadeddata', onCanPlay)
  }
  isBuffering.value = false
  document.removeEventListener('fullscreenchange', onFsChange)
  document.removeEventListener('webkitfullscreenchange', onFsChange)
  if (document.fullscreenElement || document.webkitFullscreenElement) {
    try { (document.exitFullscreen || document.webkitExitFullscreen)?.call(document) } catch(e) {}
  }
  try { screen.orientation?.unlock?.() } catch(e) {}
  // 清理 CSS 全屏
  removeCssOrientation(playerRef.value)
  playerRef.value?.classList.remove('nf-vp-css-fs', 'nf-vp-css-fs-portrait')
  vpFs.value = false
  vpLand.value = true
}

const onFsChange = () => {
  const isNativeFs = !!(document.fullscreenElement || document.webkitFullscreenElement)
  if (!isNativeFs && vpFs.value && !playerRef.value?.classList.contains('nf-vp-css-fs')) {
    vpFs.value = false
    removeCssOrientation(playerRef.value)
    try { screen.orientation?.unlock?.() } catch(e) {}
  }
  if (isNativeFs) vpFs.value = true
  if (!vpFs.value) {
    vpLand.value = true
    showCtrlBriefly()
  }
}

const vpTogglePlay = () => {
  const v = getVideo()
  if (!v) return
  v.paused ? v.play().catch(() => {}) : v.pause()
  showCtrlBriefly()
}

const vpCycleSpeed = () => {
  const idx = SPEEDS.indexOf(playbackSpeed.value)
  playbackSpeed.value = SPEEDS[(idx + 1) % SPEEDS.length]
  if (getVideo()) getVideo().playbackRate = playbackSpeed.value
  showTip(playbackSpeed.value + 'x')
  showCtrlBriefly()
}

const vpToggleFs = () => {
  const el = playerRef.value
  if (!el) return
  if (vpFs.value) {
    if (document.fullscreenElement || document.webkitFullscreenElement) {
      (document.exitFullscreen || document.webkitExitFullscreen)?.call(document)
      try { screen.orientation?.unlock?.() } catch(e) {}
    } else {
      // CSS 全屏回退
      vpFs.value = false
      removeCssOrientation(el)
      el.classList.remove('nf-vp-css-fs', 'nf-vp-css-fs-portrait')
    }
  } else {
    const fn = el.requestFullscreen || el.webkitRequestFullscreen
    if (fn) {
      fn.call(el).then(() => {
        vpLand.value = true
        if (canNativeLock()) {
          // Android：使用原生方向锁定
          screen.orientation.lock('landscape').catch(() => {
            applyCssLandscape(el, true)
          })
        } else {
          // iOS：使用 CSS transform 模拟横屏
          applyCssLandscape(el, true)
        }
      }).catch(() => {
        enterCssFs(el)
      })
    } else {
      enterCssFs(el)
    }
  }
}

const enterCssFs = (el) => {
  el.classList.add('nf-vp-css-fs')
  vpFs.value = true
  vpLand.value = true
  applyCssLandscape(el, true)
}

// 用 CSS transform 模拟横屏（iOS 不支持 screen.orientation.lock）
const applyCssLandscape = (el, landscape) => {
  if (!el) return
  if (landscape) {
    el.classList.add('nf-vp-landscape')
    el.classList.remove('nf-vp-portrait')
    usingCssRotation.value = true
  } else {
    el.classList.remove('nf-vp-landscape')
    el.classList.add('nf-vp-portrait')
    usingCssRotation.value = false
  }
}

const removeCssOrientation = (el) => {
  if (!el) return
  el.classList.remove('nf-vp-landscape', 'nf-vp-portrait')
  usingCssRotation.value = false
}

const vpToggleOri = () => {
  if (!vpFs.value) return
  vpLand.value = !vpLand.value
  const el = playerRef.value
  if (canNativeLock() && (document.fullscreenElement || document.webkitFullscreenElement)) {
    // Android：使用原生方向锁定
    screen.orientation.lock(vpLand.value ? 'landscape' : 'portrait').catch(() => {
      applyCssLandscape(el, vpLand.value)
    })
  } else {
    // iOS / CSS 全屏：使用 CSS transform
    applyCssLandscape(el, vpLand.value)
  }
}

let seekingByClick = false
let seekPrevTime = 0
let switchingEpisode = false
let switchTimer = null
let pendingResumeTime = 0   // 待恢复的播放时间（由 playFirst 设置）
let lastSaveTime = 0        // 上次保存进度的时间戳（节流用）

const vpProgClick = (e) => {
  const bar = vpProgRef.value
  const v = getVideo()
  if (!bar || !v) return
  const rect = bar.getBoundingClientRect()
  let pct
  if (usingCssRotation.value) {
    pct = Math.max(0, Math.min(1, (e.clientY - rect.top) / rect.height))
  } else {
    pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
  }
  const targetTime = pct * duration.value
  seekPrevTime = currentTime.value
  // 立即更新进度条视觉位置
  currentTime.value = targetTime
  isBuffering.value = true
  seekingByClick = true
  v.currentTime = targetTime
  // 监听 seek 结果
  const cleanupSeek = () => {
    v.removeEventListener('seeked', onSeeked)
    v.removeEventListener('error', onSeekErr)
    clearTimeout(seekTimeout)
  }
  const onSeeked = () => {
    cleanupSeek()
    seekingByClick = false
  }
  const onSeekErr = () => {
    cleanupSeek()
    // seek 失败，回退到之前的位置
    currentTime.value = seekPrevTime
    v.currentTime = seekPrevTime
    seekingByClick = false
    isBuffering.value = false
  }
  // iOS 可能不触发 seeked，5秒超时兜底
  const seekTimeout = setTimeout(() => {
    cleanupSeek()
    seekingByClick = false
    // 用 video 实际位置同步
    currentTime.value = v.currentTime
  }, 5000)
  v.addEventListener('seeked', onSeeked)
  v.addEventListener('error', onSeekErr)
  showCtrlBriefly()
}

const vpTimeUpdate = () => {
  const v = getVideo()
  if (!v) return
  // seek 过程中不让 timeupdate 覆盖视觉进度和 buffering 状态
  if (seekingByClick) return
  currentTime.value = v.currentTime
  duration.value = v.duration || 0
  // timeupdate 说明在播放，清除加载状态
  if (isBuffering.value && !switchingEpisode) isBuffering.value = false
  // 每10秒节流保存播放进度
  const now = Date.now()
  if (goodID.value && currentTime.value > 5 && now - lastSaveTime > 10000) {
    lastSaveTime = now
    playHistoryStore.saveProgress(goodID.value, {
      episodeIndex: currentIndex.value,
      currentTime: currentTime.value,
      duration: duration.value,
      episodeName: currentEpisode.value?.name || epLabel(currentIndex.value + 1),
      imageUrl: data.value.imageUrl,
      title: data.value.title
    })
  }
}

const vpEnded = () => {
  isPlaying.value = false
  vpCtrlVis.value = true
  if (currentIndex.value < episodes.value.length - 1) {
    switchEpisode(currentIndex.value + 1)
  }
}

// === 手势控制 ===
const vpTouchStart = (e) => {
  const t = e.touches[0]
  if (!t) return
  tStartX = t.clientX
  tStartY = t.clientY
  tStartTime = Date.now()
  tDir = ''
  tSeeking = false
  seekStart = currentTime.value
  longTimer = setTimeout(() => {
    speedBoosting = true
    prevRate = getVideo()?.playbackRate || 1
    if (getVideo()) getVideo().playbackRate = 2
    showTip($t.value('playerFastFwd'), 60000)
  }, 500)
}

const vpTouchMove = (e) => {
  const t = e.touches[0]
  if (!t) return
  clearTimeout(longTimer)
  longTimer = null
  if (speedBoosting) return

  let dx = t.clientX - tStartX
  let dy = t.clientY - tStartY

  // CSS rotate(90deg) 后坐标系需要旋转映射：
  // 视觉水平(快进) = 屏幕竖直方向，视觉竖直(亮度/音量) = 屏幕水平方向(反向)
  if (usingCssRotation.value) {
    const rawDx = dx
    dx = dy
    dy = -rawDx
  }

  if (!tDir) {
    if (Math.abs(dx) > 10 || Math.abs(dy) > 10) {
      tDir = Math.abs(dx) > Math.abs(dy) ? 'h' : 'v'
    } else return
  }

  const rect = playerRef.value?.getBoundingClientRect()
  if (!rect) return

  // CSS 旋转后 getBoundingClientRect 的宽高互换
  const visualWidth = usingCssRotation.value ? rect.height : rect.width
  const visualHeight = usingCssRotation.value ? rect.width : rect.height

  if (tDir === 'h') {
    tSeeking = true
    const seekDelta = (dx / visualWidth) * duration.value * 0.5
    seekTarget = Math.max(0, Math.min(duration.value, seekStart + seekDelta))
    showTip(`${seekDelta >= 0 ? '+' : ''}${Math.round(seekDelta)}s → ${fmtTime(seekTarget)}`, 60000)
  } else {
    const delta = -dy / (visualHeight * 0.7)
    // CSS 旋转后：视频左半边 = 屏幕上半边
    const isLeftHalf = usingCssRotation.value
      ? (tStartY < rect.top + rect.height / 2)
      : (tStartX < rect.left + rect.width / 2)
    if (isLeftHalf) {
      brightnessVal.value = Math.max(0.2, Math.min(1.5, brightnessVal.value + delta * 0.01))
      showTip(`☀ ${Math.round(brightnessVal.value * 100)}%`, 60000)
    } else {
      if (getVideo()) {
        const v = Math.max(0, Math.min(1, getVideo().volume + delta * 0.01))
        getVideo().volume = v
        showTip(`🔊 ${Math.round(v * 100)}%`, 60000)
      }
    }
  }
}

const vpTouchEnd = () => {
  clearTimeout(longTimer)
  longTimer = null
  if (speedBoosting) {
    speedBoosting = false
    if (getVideo()) getVideo().playbackRate = prevRate
    vpTip.value = ''
    return
  }
  if (tSeeking) {
    if (getVideo()) getVideo().currentTime = seekTarget
    tSeeking = false
    setTimeout(() => { vpTip.value = '' }, 400)
    return
  }
  vpTip.value = ''
  if (tDir) return
  if (Date.now() - tStartTime < 300) {
    if (vpCtrlVis.value) {
      vpCtrlVis.value = false
      clearTimeout(ctrlTimer)
    } else {
      showCtrlBriefly()
    }
  }
}

const resolveVideoUrl = (ep) => {
  const sources = [ep.attrs, ep.specs]
  for (const raw of sources) {
    if (!raw) continue
    try {
      const d = typeof raw === 'string' ? JSON.parse(raw) : raw
      if (Array.isArray(d)) {
        const found = d.find(s =>
          s.label && /^(videoUrl|video_url|url|视频链接|视频地址|视频)$/i.test(s.label)
        )
        if (found && found.value) return found.value
      } else if (d && typeof d === 'object') {
        const url = d.videoUrl || d.video_url || d.url || ''
        if (url) return url
      }
    } catch (e) {}
  }
  return ''
}

const playFirst = () => {
  if (episodes.value.length > 0) {
    const progress = playHistoryStore.getProgress(goodID.value)
    const epIdx = progress ? Math.min(progress.episodeIndex, episodes.value.length - 1) : 0
    playEpisode(epIdx)
  } else {
    pendingResumeTime = 0
    openTheater()
  }
}

const playEpisode = (idx) => {
  const ep = episodes.value[idx]
  if (!ep) return
  currentIndex.value = idx
  currentEpisode.value = ep
  videoUrl.value = resolveVideoUrl(ep)
  // 查询该集的历史进度
  const epProg = playHistoryStore.getEpisodeProgress(goodID.value, idx)
  pendingResumeTime = epProg ? epProg.currentTime : 0
  lastSaveTime = 0  // 新集数，重置节流计时
  openTheater()
}

const switchEpisode = (idx) => {
  const ep = episodes.value[idx]
  if (!ep) return
  // 切换前先保存当前集进度
  if (goodID.value && currentTime.value > 5) {
    playHistoryStore.saveProgress(goodID.value, {
      episodeIndex: currentIndex.value,
      currentTime: currentTime.value,
      duration: duration.value,
      episodeName: currentEpisode.value?.name || epLabel(currentIndex.value + 1),
      imageUrl: data.value.imageUrl,
      title: data.value.title
    })
  }
  currentIndex.value = idx
  currentEpisode.value = ep
  videoUrl.value = resolveVideoUrl(ep)
  currentTime.value = 0
  duration.value = 0
  isBuffering.value = true
  switchingEpisode = true
  // 查询该集的历史进度，有则恢复
  const epProg = playHistoryStore.getEpisodeProgress(goodID.value, idx)
  pendingResumeTime = epProg ? epProg.currentTime : 0
  lastSaveTime = 0       // 新集数重置节流
  clearTimeout(switchTimer)
  showTip($t.value('playerSwitchEp').replace('{{n}}', idx + 1), 2000)
  const v = getVideo()
  if (v) {
    v.src = videoUrl.value
    v.playbackRate = playbackSpeed.value
    const cleanup = () => {
      v.removeEventListener('canplay', onReady)
      v.removeEventListener('error', onSwitchErr)
      clearTimeout(switchTimer)
    }
    const onReady = () => {
      cleanup()
      switchingEpisode = false
      isBuffering.value = false
      // 如果有历史进度要恢复
      if (pendingResumeTime > 0) {
        const targetTime = pendingResumeTime
        const epNum = currentIndex.value + 1
        pendingResumeTime = 0
        v.currentTime = targetTime
        showTip(
          $t.value('playerContinueFrom')
            .replace('{{ep}}', epNum)
            .replace('{{time}}', fmtTime(targetTime)),
          3000
        )
      }
      v.play().catch(() => {})
    }
    const onSwitchErr = () => {
      cleanup()
      switchingEpisode = false
      isBuffering.value = false
      showTip($t.value('playerLoadFail') + '，' + $t.value('playerTapRetry'), 5000)
    }
    v.addEventListener('canplay', onReady)
    v.addEventListener('error', onSwitchErr)
    // 15秒超时：可能网络极差，canplay 不触发
    switchTimer = setTimeout(() => {
      cleanup()
      switchingEpisode = false
      // 如果 duration 有了说明其实已在播放
      if (v.duration > 0 && !v.paused) {
        isBuffering.value = false
        return
      }
      isBuffering.value = false
      showTip($t.value('playerLoadSlow'), 5000)
    }, 15000)
    v.load()
  }
}

const toggleCollect = async () => {
  if (!token) {
    uni.showToast({ title: $t.value('loginFirst') || '请登录', icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  const res = await createCollect({ goodID: Number(goodID.value) })
  if (res.code === 0) {
    collectionFlag.value = !collectionFlag.value
    uni.showToast({
      title: collectionFlag.value ? $t.value('collected') : $t.value('uncollected'),
      icon: 'none'
    })
  }
}

const shareVideo = () => {
  uni.showToast({ title: $t.value('playerShare'), icon: 'none' })
}

const goBack = () => {
  uni.navigateBack()
}

onUnmounted(() => {
  destroyPlayer()
})
</script>

<style lang="scss">
page {
  background-color: #000;
}

.nf-player {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #000;
  position: relative;
}

/* 压制浏览器原生 video 控件条（部分移动端会强制显示加载进度条） */
.nf-vp-video::-webkit-media-controls,
.nf-vp-video::-webkit-media-controls-panel,
.nf-vp-video::-webkit-media-controls-enclosure,
.nf-vp-video::-webkit-media-controls-overlay-play-button {
  display: none !important;
  -webkit-appearance: none;
}

/* ========== 全屏影院弹层 ========== */
.nf-theater {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 9999;
  background: rgba(0, 0, 0, 0.96);
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.nf-theater-show {
  pointer-events: auto;
  opacity: 1;
}

.nf-theater-inner {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.35s cubic-bezier(0.16, 1, 0.3, 1);
}

.nf-theater-fade-in {
  opacity: 1;
}

.nf-theater-close {
  position: absolute;
  top: 0;
  right: 0;
  width: 88rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  margin-top: 16rpx;
  margin-right: 16rpx;

  &:active {
    opacity: 0.6;
    transform: scale(0.9);
  }
}

.nf-theater-badge {
  display: flex;
  align-items: center;
  gap: 10rpx;
  margin-bottom: 24rpx;
}

.nf-theater-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
  background: #e50914;
  box-shadow: 0 0 10rpx rgba(229, 9, 20, 0.8);
  animation: pulseDot 1.5s ease-in-out infinite;
}

@keyframes pulseDot {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.8); }
}

.nf-theater-badge-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 3rpx;
  text-transform: uppercase;
}

.nf-theater-video-wrap {
  width: 100%;
  background: #000;
  border-radius: 0;
  box-shadow: 0 20rpx 80rpx rgba(0, 0, 0, 0.6);
  position: relative;
}

.nf-theater-video {
  width: 100%;
  height: 100%;
}

.nf-theater-poster-wrap {
  width: 100%;
  height: 56.25vw;
  position: relative;
}

.nf-theater-poster {
  width: 100%;
  height: 100%;
}

.nf-theater-no-video {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.55);
}

.nf-theater-no-video-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.5);
  letter-spacing: 2rpx;
}

.nf-theater-info {
  padding: 28rpx 40rpx 16rpx;
  width: 100%;
  box-sizing: border-box;
}

.nf-theater-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1rpx;
  margin-bottom: 8rpx;
}

.nf-theater-sub {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.4);
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

/* 弹层内集数快选 */
.nf-theater-eps {
  width: 100%;
  padding: 16rpx 0 0;
}

.nf-theater-eps-scroll {
  width: 100%;
  white-space: nowrap;
}

.nf-theater-eps-list {
  display: flex;
  padding: 0 40rpx 24rpx;
  gap: 16rpx;
}

.nf-theater-ep-btn {
  width: 72rpx;
  height: 72rpx;
  border-radius: 16rpx;
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.25s;

  &:active {
    transform: scale(0.92);
  }
}

.nf-theater-ep-active {
  background: linear-gradient(135deg, #e50914, #b20710);
  border-color: #e50914;
  box-shadow: 0 0 20rpx rgba(229, 9, 20, 0.4);
}

.nf-theater-ep-text {
  font-size: 26rpx;
  font-weight: 700;
  color: #fff;
}

/* ========== 封面 Hero ========== */
.nf-hero {
  width: 100%;
  height: 560rpx;
  position: relative;
  flex-shrink: 0;
  overflow: hidden;
}

.nf-hero-status {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: var(--status-bar-height, 0px);
  z-index: 5;
}

.nf-hero-img {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0; left: 0;
}

.nf-hero-gradient {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  height: 320rpx;
  background: linear-gradient(to top, #000 0%, rgba(0,0,0,0.7) 50%, transparent 100%);
  pointer-events: none;
}

.nf-hero-back {
  position: absolute;
  top: var(--status-bar-height, 0px);
  left: 20rpx;
  margin-top: 16rpx;
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;

  &:active {
    background: rgba(0, 0, 0, 0.7);
    transform: scale(0.92);
  }
}

/* ===== 内容滚动区 ===== */
.nf-content-scroll {
  flex: 1;
  width: 100%;
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
  }
}

/* ===== 标题信息区 ===== */
.nf-info-section {
  padding: 28rpx 32rpx 0;
}

.nf-title {
  font-size: 40rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 1rpx;
  line-height: 1.3;
  margin-bottom: 16rpx;
}

.nf-meta-row {
  display: flex;
  align-items: center;
  margin-bottom: 24rpx;
}

.nf-meta-item {
  display: flex;
  align-items: baseline;
  gap: 6rpx;
}

.nf-meta-star {
  color: #ffd700;
  font-size: 26rpx;
}

.nf-meta-value {
  color: #fff;
  font-size: 28rpx;
  font-weight: 700;
}

.nf-meta-label {
  color: rgba(255, 255, 255, 0.4);
  font-size: 22rpx;
  margin-left: 4rpx;
}

.nf-meta-divider {
  width: 1rpx;
  height: 24rpx;
  background: rgba(255, 255, 255, 0.15);
  margin: 0 24rpx;
}

/* ===== 操作栏 ===== */
.nf-actions {
  display: flex;
  gap: 48rpx;
  padding: 16rpx 0 24rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;

  &:active {
    opacity: 0.7;
    transform: scale(0.95);
  }
}

.nf-action-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== 简介 ===== */
.nf-desc-section {
  padding: 24rpx 32rpx;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-section-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 1rpx;
  margin-bottom: 12rpx;
}

.nf-desc-text {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.55);
  line-height: 1.65;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}

.nf-desc-expand {
  -webkit-line-clamp: unset;
}

.nf-desc-toggle {
  font-size: 24rpx;
  color: #e50914;
  margin-top: 8rpx;

  &:active {
    opacity: 0.7;
  }
}

/* ===== 选集 ===== */
.nf-episodes-section {
  padding: 24rpx 0 0;
}

.nf-section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 32rpx 16rpx;
}

.nf-episode-count {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.35);
}

.nf-episodes-scroll {
  width: 100%;
  white-space: nowrap;
}

.nf-episodes-list {
  display: flex;
  padding: 0 32rpx 24rpx;
  gap: 20rpx;
}

.nf-episode-card {
  width: 280rpx;
  flex-shrink: 0;
  border-radius: 16rpx;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  transition: all 0.3s;

  &:active {
    transform: scale(0.97);
  }
}

.nf-episode-active {
  border-color: rgba(229, 9, 20, 0.5);
  background: rgba(229, 9, 20, 0.08);
  box-shadow: 0 0 20rpx rgba(229, 9, 20, 0.15);
}

.nf-ep-thumb-wrap {
  width: 100%;
  height: 160rpx;
  position: relative;
  overflow: hidden;
}

.nf-ep-thumb {
  width: 100%;
  height: 100%;
}

.nf-ep-play-icon {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 50rpx;
  height: 50rpx;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-ep-play-text {
  font-size: 22rpx;
  color: #fff;
  margin-left: 3rpx;
}

.nf-ep-playing {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: flex-end;
  gap: 4rpx;
  height: 36rpx;
}

.nf-eq-bar {
  width: 6rpx;
  background: #e50914;
  border-radius: 3rpx;
  animation: eqBounce 0.8s ease-in-out infinite alternate;

  &:nth-child(1) { height: 60%; animation-delay: 0s; }
  &:nth-child(2) { height: 100%; animation-delay: 0.2s; }
  &:nth-child(3) { height: 40%; animation-delay: 0.4s; }
}

@keyframes eqBounce {
  0% { height: 30%; }
  100% { height: 100%; }
}

.nf-ep-info {
  padding: 14rpx 16rpx;
  white-space: normal;
}

.nf-ep-name {
  font-size: 24rpx;
  font-weight: 600;
  color: #fff;
  margin-bottom: 6rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nf-ep-desc {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.4);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 6rpx;
}

.nf-ep-price {
  font-size: 24rpx;
  color: #e50914;
  font-weight: 700;
}

.nf-ep-free {
  font-size: 22rpx;
  color: #22c55e;
  font-weight: 600;
}

/* ===== 无剧集 ===== */
.nf-no-episodes {
  padding: 80rpx 0;
  text-align: center;
}

.nf-no-episodes-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.3);
}

/* ===== 自定义视频播放器 ===== */
.nf-vp {
  position: relative;
  width: 100%;
  background: #000;
  line-height: 0;
  user-select: none;
  -webkit-user-select: none;
  -webkit-touch-callout: none;
}

.nf-vp-video {
  width: 100%;
  display: block;
}

.nf-vp:fullscreen,
.nf-vp:-webkit-full-screen {
  width: 100vw !important;
  height: 100vh !important;
  height: 100dvh !important;
}

.nf-vp:fullscreen .nf-vp-video,
.nf-vp:-webkit-full-screen .nf-vp-video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.nf-vp-touch {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 2;
}

.nf-vp-tip {
  position: absolute;
  top: 50%; left: 50%;
  transform: translate(-50%, -50%);
  z-index: 15;
  padding: 8px 18px;
  border-radius: 8px;
  background: rgba(0, 0, 0, 0.75);
  color: #fff;
  font-size: 14px;
  line-height: 1.4;
  white-space: nowrap;
  pointer-events: none;
}

.nf-vp-loading {
  position: absolute;
  top: 50%; left: 50%;
  transform: translate(-50%, -50%);
  z-index: 15;
  pointer-events: none;
}

.nf-vp-spinner {
  width: 36px; height: 36px;
  border: 3px solid rgba(255,255,255,0.2);
  border-top-color: #e50914;
  border-radius: 50%;
  animation: vpSpin 0.8s linear infinite;
}

@keyframes vpSpin { to { transform: rotate(360deg); } }

.nf-vp-ctrl {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  z-index: 10;
  background: linear-gradient(transparent, rgba(0,0,0,0.85));
  padding: 20px 10px calc(6px + env(safe-area-inset-bottom, 0px));
  opacity: 0;
  transition: opacity 0.25s;
  pointer-events: none;
}

.nf-vp-ctrl-show {
  opacity: 1;
  pointer-events: auto;
}

.nf-vp-prog {
  position: relative;
  height: 22px;
  display: flex;
  align-items: center;
  margin-bottom: 2px;
}

.nf-vp-prog-bar {
  width: 100%;
  height: 3px;
  background: rgba(255,255,255,0.25);
  border-radius: 2px;
  overflow: hidden;
}

.nf-vp-prog-fill {
  height: 100%;
  background: #e50914;
  border-radius: 2px;
}

.nf-vp-prog-dot {
  position: absolute;
  width: 12px; height: 12px;
  border-radius: 50%;
  background: #e50914;
  box-shadow: 0 0 4px rgba(229,9,20,0.5);
  top: 50%;
  margin-left: -6px;
  transform: translateY(-50%);
  pointer-events: none;
}

.nf-vp-btns {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nf-vp-btn {
  min-width: 36px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 15px;
  -webkit-tap-highlight-color: transparent;
  user-select: none;
  -webkit-user-select: none;
  -webkit-touch-callout: none;

  &:active { opacity: 0.6; }
}

.nf-vp-time {
  font-size: 12px;
  color: rgba(255,255,255,0.65);
  white-space: nowrap;
  line-height: 1;
}

.nf-vp-spacer { flex: 1; }

.nf-vp-spd,
.nf-vp-fs,
.nf-vp-ori {
  font-size: 12px;
  padding: 0 8px;
  border: 1px solid rgba(255,255,255,0.3);
  border-radius: 4px;
  height: 26px;
}

.nf-vp-ori {
  background: rgba(255,255,255,0.12);
  border-color: rgba(255,255,255,0.2);
}

/* CSS 全屏（iOS Safari 回退） */
.nf-vp-css-fs {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
  width: 100vw !important;
  height: 100vh !important;
  height: 100dvh !important;
  z-index: 999999 !important;
  border-radius: 0 !important;
  background: #000 !important;
}

.nf-vp-css-fs .nf-vp-video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

/* CSS 全屏时控制栏额外底部安全距离 */
.nf-vp-css-fs .nf-vp-ctrl {
  padding-bottom: calc(6px + env(safe-area-inset-bottom, 0px) + 10px);
}

/* CSS 全屏竖屏方向 */
.nf-vp-css-fs-portrait {
  width: 100vw !important;
  height: 100vh !important;
  height: 100dvh !important;
}

/* CSS transform 模拟横屏（兼容所有浏览器） */
.nf-vp-landscape {
  transform: rotate(90deg) !important;
  transform-origin: center center !important;
  width: 100vh !important;
  width: 100dvh !important;
  height: 100vw !important;
  /* 居中修正：旋转后需要平移 */
  position: fixed !important;
  top: 50% !important;
  left: 50% !important;
  margin: 0 !important;
  translate: -50% -50% !important;
}

.nf-vp-landscape .nf-vp-video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.nf-vp-landscape .nf-vp-ctrl {
  padding-bottom: calc(6px + env(safe-area-inset-bottom, 0px) + 10px);
}

/* 竖屏模式（全屏下切换到竖屏） */
.nf-vp-portrait {
  transform: rotate(0deg) !important;
  transform-origin: center center !important;
  width: 100vw !important;
  height: 100vh !important;
  height: 100dvh !important;
}
</style>
