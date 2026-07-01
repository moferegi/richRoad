<template>
  <view class="player-container">
    <!-- 顶部视频区 (固定不动) -->
    <view class="video-section">
      <video 
        id="englishVideo" 
        class="main-video" 
        :src="videoInfo.videoUrl" 
        :autoplay="true"
        :loop="controls.loop"
        @timeupdate="onTimeUpdate"
        @play="isPlaying = true"
        @pause="isPlaying = false"
        @error="onVideoError"
        @ended="onVideoEnded"
      ></video>
    </view>

    <!-- 中央外挂字幕滚动区 -->
    <scroll-view 
      class="subtitle-section" 
      scroll-y 
      scroll-with-animation
      :scroll-into-view="activeSentenceId"
      v-if="controls.showSubtitles"
    >
      <view class="subtitle-padding-top"></view>
      
      <view 
        v-for="(item, index) in subtitleList" 
        :key="index"
        :id="'sentence_' + index"
        class="sentence-row"
        :class="{ 'is-active': activeSentenceIndex === index, 'is-focused': focusedSentenceId === Number(item.id) }"
      >
        <view class="sentence-meta">
          <text class="time-badge">{{ formatSeconds(item.startTime) }} - {{ formatSeconds(item.endTime) }}</text>
          <text class="collect-btn" @click="collectSentence(item.id)">⭐ {{ t('player.collect') }}</text>
        </view>
        
        <!-- 英文句子 (带高亮重点词可点击分析) -->
        <view class="english-text">
          <text 
            v-for="(seg, sIdx) in parseHighlight(item.english)" 
            :key="sIdx"
            :class="{ 'highlight-word': seg.isHighlight }"
            @click="seg.isHighlight ? showWordDetail(seg.wordId) : null"
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
        <text class="c-icon">{{ controls.showBilingual ? '🇺🇳' : '🇬🇧' }}</text>
        <text class="c-text">{{ t('player.bilingual') }}</text>
      </view>
      
      <view class="control-item" @click="changeSpeed">
        <text class="c-icon">⚡</text>
        <text class="c-text">{{ controls.speed }}x</text>
      </view>
      
      <view class="control-item play-btn" @click="togglePlay">
        <text class="c-icon">{{ isPlaying ? '⏸' : '▶' }}</text>
      </view>
      
      <view class="control-item" @click="toggleLoop">
        <text class="c-icon" :class="{ 'c-active': controls.loop }">🔁</text>
        <text class="c-text">{{ t('player.loop') }}</text>
      </view>
      
      <view class="control-item" @click="toggleSubtitles">
        <text class="c-icon" :class="{ 'c-active': controls.showSubtitles }">💬</text>
        <text class="c-text">{{ t('player.subtitle') }}</text>
      </view>
    </view>

    <!-- 重点词汇查词弹窗 -->
    <uni-popup ref="wordPopup" type="bottom" background-color="#fff">
      <view class="word-detail-box" v-if="currentWord">
        <view class="w-header">
          <text class="w-title">{{ currentWord.word }}</text>
          <text class="w-audio" @click="playWordAudio(currentWord.audioUs)">🔊</text>
        </view>
        <text class="w-phonetic">{{ currentWord.phoneticUs }}</text>
        <view class="w-exp">{{ localText(currentWord.explanation) }}</view>
        <button class="w-collect" @click="collectWord(currentWord.id)">⭐ {{ t('typing.collect') }}</button>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { collect, findVideoEpisode, findWord, getSentenceList, getWatchProgress, heartbeat } from '@/api/learning.js'

const langStore = useLangStore()
const fallbackTexts = {
  'player.collect': '收藏句子',
  'player.bilingual': '双语',
  'player.loop': '循环',
  'player.subtitle': '字幕',
  'typing.collect': '收藏单词',
  'player.trial_end_title': '试看结束',
  'player.trial_end_desc': '免费额度已用完，积分兑换观看时长或升级会员解锁全集。',
  'player.btn_vip': '去开通',
  'player.load_failed': '视频加载失败',
  'player.no_subtitle': '暂无字幕',
  'player.collect_success': '收藏成功',
  'player.sentence_focused': '已定位到目标句子',
  'player.sentence_not_found': '未找到目标句子',
}

const t = (k, d = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(k, locale)
  if (text && text !== k) return text
  return d || fallbackTexts[k] || k
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
let audioCtx = null
const isPlaying = ref(true)
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
let lastSignedRefreshAt = 0
let focusSentenceTimer = null

onMounted(() => {
  videoCtx = uni.createVideoContext('englishVideo')
  if (pendingSeekTime.value > 0) {
    videoCtx.seek(pendingSeekTime.value)
    pendingSeekTime.value = 0
  }
})

onUnmounted(() => {
  flushHeartbeat(lastHeartbeatTime, true)
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

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const seekVideo = (secs) => {
  const target = Math.max(0, Number(secs || 0))
  if (videoCtx) {
    videoCtx.seek(target)
  } else {
    pendingSeekTime.value = target
  }
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
  activeSentenceId.value = 'sentence_' + Math.max(0, targetIndex - 1)
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
  const previousUrl = videoInfo.value.videoUrl
  const resumeTarget = Math.max(0, Number(resumeSecs || 0))

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
    if (urlChanged) {
      setTimeout(() => {
        if (resumeTarget > 0) {
          seekVideo(Math.max(0, resumeTarget - 1))
        }
        if (isPlaying.value && videoCtx) {
          videoCtx.play()
        }
      }, 350)
    }

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
  await loadEpisode()
  await loadSubtitleList()
  const focused = focusSentence()
  if (!focused) {
    if (sentenceId.value > 0) {
      uni.showToast({ title: t('player.sentence_not_found'), icon: 'none' })
    }
    await loadWatchHistory()
  }
}

onLoad((options) => {
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
  const currentTime = e.detail.currentTime
  const duration = e.detail.duration || videoInfo.value.duration
  if (duration > 0) {
    videoInfo.value.duration = duration
  }

  // 1. 防盗验证：如果没买且过了 8%，直接暂停+弹窗斩断
  if (!authData.value.hasFullAuth && duration > 0) {
    const percent = (currentTime / duration) * 100
    if (percent >= authData.value.trialPercent) {
      videoCtx.pause()
      uni.showModal({
        title: t('player.trial_end_title', '试看结束'),
        content: t('player.trial_end_desc', '免费额度已用完，积分兑换观看时长或升级会员解锁全集。'),
        confirmText: t('player.btn_vip', '去开通'),
        success: (res) => { if(res.confirm) { uni.navigateTo({url: '/pages/learning/profile'}) } }
      })
      // 可以防止用户拖进度条强看
      videoCtx.seek(currentTime - 2) 
      return
    }
  }

  // 2. 匹配外挂字幕轴并滚动
  const currentIdx = subtitleList.value.findIndex(sub => currentTime >= sub.startTime && currentTime <= sub.endTime)
  if (currentIdx !== -1 && currentIdx !== activeSentenceIndex.value) {
    activeSentenceIndex.value = currentIdx
    activeSentenceId.value = 'sentence_' + Math.max(0, currentIdx - 1) // 滚动锚点上移一行体验更好
  } else if (currentIdx === -1) {
    activeSentenceIndex.value = -1 // 空白期
  }

  // 3. 签名链接即将到期时，提前刷新，避免播放中途 403。
  refreshSignedUrlIfNeeded(currentTime)

  // 4. 心跳累加与上报逻辑 (Phase 2 的前端呼应)
  // 如果正常播放，累加观看差值，每满 30 秒 Post 给后端一次
  if (lastHeartbeatTime > 0 && currentTime > lastHeartbeatTime) {
    accumulatedViewTime += (currentTime - lastHeartbeatTime)
    if (accumulatedViewTime >= HEARTBEAT_FLUSH_THRESHOLD) {
      flushHeartbeat(currentTime)
    }
  }
  lastHeartbeatTime = currentTime
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

const onVideoError = async () => {
  const recovered = await refreshSignedEpisodeUrl(lastHeartbeatTime)
  if (!recovered) {
    uni.showToast({ title: t('player.load_failed'), icon: 'none' })
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

// 查词查词交互
const wordPopup = ref(null)
const currentWord = ref(null)
const showWordDetail = async (wordId) => {
  if (videoCtx) {
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
const togglePlay = () => { isPlaying.value ? videoCtx.pause() : videoCtx.play() }
const toggleBilingual = () => { controls.value.showBilingual = !controls.value.showBilingual }
const toggleSubtitles = () => { controls.value.showSubtitles = !controls.value.showSubtitles }
const toggleLoop = () => { controls.value.loop = !controls.value.loop }
const changeSpeed = () => {
  let next = controls.value.speed + 0.25
  if (next > 2.0) next = 0.75
  controls.value.speed = next
  videoCtx.playbackRate(next)
}
const collectSentence = async (id) => {
  const res = await collect(2, Number(id || 0))
  if (res.code === 0) {
    uni.showToast({ title: t('player.collect_success'), icon: 'success' })
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
  audioCtx.src = src
  audioCtx.play()
}
</script>

<style scoped>
.player-container { display: flex; flex-direction: column; height: 100vh; background: #1a1a1a; color: #fff;}

/* 视频区域 */
.video-section { width: 100%; height: 420rpx; background: #000; flex-shrink: 0; }
.main-video { width: 100%; height: 100%; }

/* 外挂字幕区域 */
.subtitle-section { flex: 1; overflow: hidden; position: relative; }
.subtitle-padding-top, .subtitle-padding-bottom { height: 40%; }
.sentence-row { padding: 20rpx 40rpx; margin-bottom: 20rpx; transition: all 0.3s; opacity: 0.6; }
.sentence-row.is-active { opacity: 1; transform: scale(1.05); background: rgba(255,255,255,0.05); border-radius: 12rpx; border-left: 6rpx solid #409eff; }
.sentence-row.is-focused { opacity: 1; border-left: 6rpx solid #facc15; box-shadow: 0 0 0 2rpx rgba(250, 204, 21, 0.35) inset; }

.sentence-meta { display: flex; justify-content: space-between; margin-bottom: 10rpx; }
.time-badge { font-size: 22rpx; color: #999; background: rgba(0,0,0,0.5); padding: 4rpx 10rpx; border-radius: 6rpx; }
.collect-btn { font-size: 22rpx; color: #ff9800; cursor: pointer; }

/* 自主控制高亮词汇，绕过 rich-text 不能点击内部标签的缺陷 */
.english-text { font-size: 34rpx; font-weight: 500; margin-bottom: 15rpx; line-height: 1.5; color: #ececec; }
.highlight-word { color: #fed700; font-weight: bold; border-bottom: 1px dashed #fed700; display: inline-block; padding: 0 4rpx; margin: 0 4rpx; }
.translated-text { font-size: 28rpx; color: #bbb; line-height: 1.4; }

/* 底部操作区 */
.bottom-controls { height: 120rpx; background: #000; flex-shrink: 0; display: flex; justify-content: space-around; align-items: center; padding-bottom: env(safe-area-inset-bottom); border-top: 1px solid #333; }
.control-item { display: flex; flex-direction: column; align-items: center; justify-content: center; }
.c-icon { font-size: 40rpx; margin-bottom: 6rpx; filter: grayscale(1); }
.c-active { filter: grayscale(0); text-shadow: 0 0 10rpx rgba(255,255,255,0.8); }
.c-text { font-size: 20rpx; color: #888; }
.play-btn .c-icon { font-size: 60rpx; filter: grayscale(0); color: #fff;}

/* 弹窗内容 */
.word-detail-box { padding: 40rpx; display: flex; flex-direction: column; }
.w-header { display: flex; align-items: center; justify-content: space-between; }
.w-title { font-size: 48rpx; font-weight: bold; color: #333; }
.w-audio { font-size: 40rpx; }
.w-phonetic { font-size: 28rpx; color: #888; margin: 10rpx 0 30rpx 0; }
.w-exp { font-size: 32rpx; color: #444; line-height: 1.6; margin-bottom: 40rpx; }
.w-collect { background: #409eff; color: #fff; border-radius: 40rpx; }
</style>
