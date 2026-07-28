<template>
  <view class="diary-detail-container">
    <!-- 固定顶部信息区 -->
    <view class="top-bar">
      <view class="top-bar-bg"></view>
      <view class="top-bar-content">
        <view class="back-btn" @click="goBack">
          <text class="back-icon">‹</text>
        </view>
        <view class="top-info">
          <text class="diary-title">{{ localText(diaryInfo.name) }}</text>
          <view class="audio-duration" :class="{ 'duration-playing': isAudioPlaying }">
            <text class="duration-text">{{ formatSeconds(currentDuration) }}</text>
          </view>
        </view>
        <view class="top-actions">
          <view class="accent-switch" @click="toggleAccent">
            <template v-if="isAccentSwitching">
              <text class="accent-switching-text">切换中</text>
            </template>
            <template v-else>
              <text class="accent-text" :class="{ active: accent === 'US' }">US</text>
              <text class="accent-divider">/</text>
              <text class="accent-text" :class="{ active: accent === 'UK' }">UK</text>
            </template>
          </view>
          <view class="image-btn" @click="viewImage">
            <image class="image-btn-icon" src="/static/images/learning/icon-image-white.svg" mode="aspectFit" />
          </view>
        </view>
      </view>
    </view>

    <!-- 音频播放器（隐藏视频） -->
    <view class="audio-section" style="display:none;">
    </view>

    <!-- 字幕滚动区 -->
    <scroll-view 
      class="subtitle-section" 
      scroll-y 
      scroll-with-animation
      :scroll-top="subtitleScrollTop"
      v-if="controls.showSentences"
    >
      <view class="subtitle-padding-top"></view>
      
      <view 
        v-for="(item, index) in subtitleList" 
        :key="index"
        :id="'sentence_' + index"
        class="sentence-row"
        :class="{ 'is-active': activeSentenceIndex === index }"
        @click="!item.locked && jumpBySubtitle(item, index)"
      >
        <!-- 音频图标 -->
        <view class="sentence-audio-icon" :class="{ 'audio-playing': activeSentenceIndex === index && isAudioPlaying }">
          <text class="audio-icon-char">{{ activeSentenceIndex === index && isAudioPlaying ? '◉' : '♪' }}</text>
        </view>

        <!-- 锁定句子显示锁头和时间 -->
        <view v-if="item.locked" class="sentence-locked">
          <text class="lock-time">{{ formatTimeRange(item) }}</text>
          <text class="lock-icon">🔒</text>
        </view>

        <!-- 英文句子 -->
        <view v-else class="sentence-text-wrap">
          <view class="english-text" v-if="diaryInfo.showEnglish !== false">
            <text 
              v-for="(seg, sIdx) in parseHighlight(item.english)" 
              :key="sIdx"
              :class="{ 'highlight-word': seg.isHighlight }"
              @click.stop="handleSegClick(seg)"
            >{{ seg.text }}</text>
          </view>

          <!-- 对应多语翻译 -->
          <view class="translated-text" v-if="controls.showTranslate">
            {{ localText(item.translate) }}
          </view>
        </view>

        <!-- 收藏句子按钮（锁定句子不显示） -->
        <view v-if="!item.locked" class="sentence-collect-btn" @click.stop="collectSentence(item, index)">
          <text class="collect-icon">{{ sentenceCollectedMap[getSentenceId(item)] ? '★' : '☆' }}</text>
        </view>
      </view>
      
      <view class="subtitle-padding-bottom"></view>
    </scroll-view>

    <!-- 底部控制栏 -->
    <view class="bottom-controls">
      <view class="control-item" @click="toggleSentences">
        <view class="c-icon-circle" :class="{ 'c-circle-on': controls.showSentences }">
          <image class="c-icon-img" :src="controls.showSentences ? '/static/images/learning/icon-sentence-active.svg' : '/static/images/learning/icon-sentence.svg'" mode="aspectFit" />
        </view>
        <text class="c-text" :class="{ 'c-text-active': controls.showSentences }">{{ t('diary.sentence') }}</text>
      </view>
      
      <view class="control-item" @click="toggleTranslate">
        <view class="c-icon-circle" :class="{ 'c-circle-on': controls.showTranslate }">
          <image class="c-icon-img" :src="controls.showTranslate ? '/static/images/learning/icon-translate-active.svg' : '/static/images/learning/icon-translate.svg'" mode="aspectFit" />
        </view>
        <text class="c-text" :class="{ 'c-text-active': controls.showTranslate }">{{ t('diary.translate') }}</text>
      </view>
      
      <view class="control-item play-btn" @click="togglePlay">
        <view class="c-play-circle">
          <image v-if="isAudioLoading" class="c-play-img" src="/static/images/learning/icon-play-white.svg" mode="aspectFit" />
          <image v-else class="c-play-img" :src="isAudioPlaying ? '/static/images/learning/icon-pause-white.svg' : '/static/images/learning/icon-play-white.svg'" mode="aspectFit" />
        </view>
      </view>
      
      <view class="control-item" @click="collectDiary">
        <view class="c-icon-circle" :class="{ 'c-circle-on': isCollected }">
          <image class="c-icon-img" :src="isCollected ? '/static/images/learning/icon-collect-active.svg' : '/static/images/learning/icon-collect.svg'" mode="aspectFit" />
        </view>
        <text class="c-text" :class="{ 'c-text-active': isCollected }">{{ isCollected ? t('diary.uncollect') : t('diary.collect') }}</text>
      </view>

      <view class="control-item" @click="openMoreSheet">
        <view class="c-icon-circle">
          <image class="c-icon-img" src="/static/images/learning/icon-more.svg" mode="aspectFit" />
        </view>
        <text class="c-text">{{ t('diary.more') }}</text>
      </view>
    </view>

    <!-- 重点词汇查词弹窗 -->
    <uni-popup ref="wordPopup" type="bottom" background-color="#fff">
      <view class="word-detail-box" v-if="currentWord">
        <view class="popup-handle"></view>
        <view class="w-header">
          <text class="w-title">{{ currentWord.word }}</text>
          <view class="w-audio-circle" :class="{ playing: isWordAudioPlaying }" @click="playWordAudio">
            <text class="w-audio-icon">{{ isWordAudioPlaying ? '◉' : '♪' }}</text>
          </view>
        </view>
        <text class="w-phonetic">{{ currentWord.phoneticUs }}</text>
        <view class="w-divider"></view>
        <view class="w-exp">{{ localText(currentWord.explanation) }}</view>
        <view class="w-actions">
          <button class="w-cancel-btn" @click="closeWordPopup">{{ t('common.cancel') || '取消' }}</button>
          <button class="w-collect" @click="collectWord(currentWord.id)">{{ isWordCollected ? '★ ' + t('diary.uncollect') : '☆ ' + t('diary.collect') }}</button>
        </view>
      </view>
    </uni-popup>

    <!-- 更多设置弹窗 -->
    <uni-popup ref="morePopup" type="bottom" background-color="#fff">
      <view class="more-sheet">
        <view class="popup-handle"></view>
        <view class="more-sheet-head">
          <text class="more-sheet-title">{{ t('diary.more_settings') }}</text>
        </view>

        <!-- 播放速率选择 -->
        <view class="speed-hint">点击字幕空白处可跳到该处播放</view>
        <view class="more-card speed-card">
          <view class="more-card-main">
            <text class="more-card-title speed-disabled-title">{{ t('player.speed') || '倍速' }}</text>
            <text class="speed-unavailable-text">暂不可用</text>
          </view>
        </view>
        <view class="speed-pills">
          <text
            v-for="rate in speedRates"
            :key="rate"
            class="speed-pill speed-pill-disabled"
          >{{ rate }}x</text>
        </view>

        <view class="more-actions">
          <button class="more-action-btn ghost" @click="handleBack">{{ t('diary.exit_page') }}</button>
          <button class="more-action-btn" @click="closeMoreSheet">{{ t('common.cancel') }}</button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { getExternalUrl, initExternalDomain } from '@/utils/url.js'
import { findDiary, getDiarySentenceList, collect, getCollectionList, uncollect, findWord } from '@/api/learning.js'
import { playWordAudio as playWordTTS, stopLocalTTS, destroyTTS } from '@/utils/learning-tts'

const langStore = useLangStore()

const t = (key) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return key
}

const localText = (value) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  return i18nLocalText(value, locale)
}

const formatSeconds = (secs) => {
  const m = Math.floor(secs / 60).toString().padStart(2, '0')
  const s = Math.floor(secs % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}

const formatTimeRange = (item) => {
  const start = formatSeconds(item.startTime || 0)
  const end = formatSeconds(item.endTime || 0)
  return `${start} - ${end}`
}

let audioCtx = null
const diaryId = ref(0)
const isAudioPlaying = ref(false)
const isAudioLoading = ref(false)
const currentTime = ref(0)
const accent = ref('US')
const isAccentSwitching = ref(false)

const diaryInfo = ref({
  id: 0,
  name: '',
  audioUs: '',
  audioUk: '',
  duration: 0,
  imageI18n: '',
  showTranslate: true,
  showEnglish: true
})

// 各口音音频时长（秒），以实际音频为准
const audioDurationUs = ref(0)
const audioDurationUk = ref(0)

const controls = ref({
  showSentences: true,
  showTranslate: true
})

const subtitleList = ref([])
const activeSentenceIndex = ref(-1)
const activeSentenceId = ref('')
const subtitleScrollTop = ref(0)
const isCollected = ref(false)
let progressTimer = null

const morePopup = ref(null)

// 重点单词弹窗
const wordPopup = ref(null)
const currentWord = ref(null)
const isWordAudioPlaying = ref(false)
const isWordCollected = ref(false)

// 播放速率
const speedRates = [0.75, 1.0, 1.25, 1.5, 2.0]
const playbackSpeed = ref(1.0)

// 句子收藏状态
const sentenceCollectedMap = ref({})

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const getSentenceId = (item) => Number(item?.id || item?.ID || 0)

// 当前显示的音频时长（以实际音频为准）
const currentDuration = computed(() => {
  return accent.value === 'US' ? audioDurationUs.value : audioDurationUk.value
})

const getAudioUrl = () => {
  const raw = accent.value === 'US' ? diaryInfo.value.audioUs : diaryInfo.value.audioUk
  return getExternalUrl(raw)
}

const toggleAccent = () => {
  const wasPlaying = isAudioPlaying.value
  if (wasPlaying) {
    stopAudio()
  }
  accent.value = accent.value === 'US' ? 'UK' : 'US'
  isAccentSwitching.value = true
  // 切换口音时重新初始化音频
  nextTick(() => {
    initAudio()
    if (wasPlaying) {
      playAudio()
    }
  })
}

const initAudio = () => {
  if (audioCtx) {
    audioCtx.destroy()
  }
  const url = getAudioUrl()
  if (!url) {
    uni.showToast({ title: t('diary.no_audio'), icon: 'none' })
    return
  }
  audioCtx = uni.createInnerAudioContext()
  audioCtx.src = url
  audioCtx.onPlay(() => {
    isAudioPlaying.value = true
    isAudioLoading.value = false
  })
  audioCtx.onPause(() => {
    isAudioPlaying.value = false
  })
  audioCtx.onStop(() => {
    isAudioPlaying.value = false
  })
  audioCtx.onEnded(() => {
    isAudioPlaying.value = false
    currentTime.value = 0
    activeSentenceIndex.value = -1
    activeSentenceId.value = ''
  })
  audioCtx.onTimeUpdate(() => {
    currentTime.value = audioCtx.currentTime || 0
    // 记录实际音频时长
    const dur = audioCtx.duration || 0
    if (dur && isFinite(dur) && dur > 0) {
      if (accent.value === 'US') {
        audioDurationUs.value = Math.round(dur * 10) / 10
      } else {
        audioDurationUk.value = Math.round(dur * 10) / 10
      }
    }
    // 试看限制检查
    if (ensureTrialLimit(currentTime.value)) {
      return
    }
    updateActiveSentence()
  })
  audioCtx.onError((err) => {
    console.error('Audio error:', err)
    isAudioPlaying.value = false
    isAudioLoading.value = false
    isAccentSwitching.value = false
    uni.showToast({ title: t('diary.audio_play_failed'), icon: 'none' })
  })
  audioCtx.onCanplay(() => {
    isAudioLoading.value = false
    isAccentSwitching.value = false
    // 获取时长
    const dur = audioCtx.duration || 0
    if (dur && isFinite(dur) && dur > 0) {
      if (accent.value === 'US') {
        audioDurationUs.value = Math.round(dur * 10) / 10
      } else {
        audioDurationUk.value = Math.round(dur * 10) / 10
      }
    }
  })
  audioCtx.onWaiting(() => {
    isAudioLoading.value = true
  })
}

const playAudio = () => {
  if (!audioCtx) {
    initAudio()
  }
  if (audioCtx) {
    isAudioLoading.value = true
    audioCtx.play()
  }
}

const pauseAudio = () => {
  if (audioCtx) {
    audioCtx.pause()
  }
}

const stopAudio = () => {
  if (audioCtx) {
    audioCtx.stop()
  }
}

const togglePlay = () => {
  if (isAudioPlaying.value) {
    pauseAudio()
  } else {
    playAudio()
  }
}

const updateActiveSentence = () => {
  const time = currentTime.value
  const useUk = accent.value === 'UK'
  let found = -1
  for (let i = subtitleList.value.length - 1; i >= 0; i--) {
    const item = subtitleList.value[i]
    // UK 口音使用 UK 时间轴（若有），否则回退用 US 时间
    const startTime = useUk && item.startTimeUk > 0 ? item.startTimeUk : item.startTime
    if (time >= startTime) {
      found = i
      break
    }
  }
  if (found !== activeSentenceIndex.value) {
    activeSentenceIndex.value = found
    if (found >= 0) {
      const rowHeight = 100
      subtitleScrollTop.value = Math.max(0, found * rowHeight)
    }
  }
}

// 试看限制：非全权限用户播放超过试看比例时暂停
const lastTrialModalAt = ref(0)
const isTrialLocked = ref(false)

const maxTrialAllowedTime = () => {
  const duration = currentDuration.value || Number(diaryInfo.value.duration || 0)
  if (diaryInfo.value.hasFullAuth || duration <= 0) {
    return Number.POSITIVE_INFINITY
  }
  const safePercent = Math.max(1, Number(diaryInfo.value.trialPercent || 8))
  // 与后端字幕过滤逻辑一致：使用 startTime >= maxTrialTime 判断
  return (duration * safePercent) / 100
}

const ensureTrialLimit = (time) => {
  if (diaryInfo.value.hasFullAuth) {
    return false
  }
  const limitTime = maxTrialAllowedTime()
  if (!isFinite(limitTime) || limitTime <= 0) {
    return false
  }
  // 与后端字幕过滤一致：currentTime >= limitTime 时触发限制
  if (time < limitTime) {
    if (isTrialLocked.value) {
      isTrialLocked.value = false
    }
    return false
  }
  // 超过试看时间，暂停并回弹到限制时间前0.5秒
  if (audioCtx) {
    pauseAudio()
    audioCtx.seek(Math.max(0, limitTime - 0.5))
  }
  isTrialLocked.value = true
  const now = Date.now()
  if (now - lastTrialModalAt.value < 4000) {
    return true
  }
  lastTrialModalAt.value = now
  uni.showModal({
    title: '试看结束',
    content: '试看时间已结束，开通会员可完整收听',
    confirmText: '去开通',
    success: (res) => {
      if (res.confirm) {
        uni.navigateTo({ url: '/pages/learning/profile' })
      }
    }
  })
  return true
}

const jumpBySubtitle = (item, index) => {
  if (!audioCtx) return
  // If clicking the same sentence that's already active and audio is playing, pause
  if (activeSentenceIndex.value === index && isAudioPlaying.value) {
    pauseAudio()
    return
  }
  // If clicking the same sentence that's paused, resume
  if (activeSentenceIndex.value === index && !isAudioPlaying.value) {
    playAudio()
    return
  }
  // Different sentence: seek and play
  const useUk = accent.value === 'UK'
  const startTime = useUk && item.startTimeUk > 0 ? item.startTimeUk : item.startTime
  audioCtx.seek(startTime)
  activeSentenceIndex.value = index
  activeSentenceId.value = 'sentence_' + index
  const rowHeight = 100
  subtitleScrollTop.value = Math.max(0, index * rowHeight)
  playAudio()
}

const parseHighlight = (text) => {
  if (!text) return [{ text: '', isHighlight: false }]
  const parts = []
  const regex = /<w id="(\d+)">(.*?)<\/w>/g
  let lastIdx = 0
  let match
  while ((match = regex.exec(text)) !== null) {
    if (match.index > lastIdx) {
      parts.push({ text: text.substring(lastIdx, match.index), isHighlight: false })
    }
    parts.push({ text: match[2], isHighlight: true, wordId: match[1] })
    lastIdx = match.index + match[0].length
  }
  if (lastIdx < text.length) {
    parts.push({ text: text.substring(lastIdx), isHighlight: false })
  }
  if (parts.length === 0) {
    parts.push({ text, isHighlight: false })
  }
  return parts
}

// 重点单词点击
const handleSegClick = (seg) => {
  if (seg && seg.isHighlight) {
    showWordDetail(seg.wordId)
  }
}

const showWordDetail = async (wordId) => {
  // 暂停音频
  if (isAudioPlaying.value) {
    pauseAudio()
  }
  const res = await findWord(Number(wordId))
  if (res.code !== 0 || !res.data) {
    return
  }
  currentWord.value = {
    ...res.data,
    id: getEntityId(res.data)
  }
  isWordCollected.value = false
  wordPopup.value?.open()
}

const playWordAudio = () => {
  if (!currentWord.value) return
  if (isWordAudioPlaying.value) {
    stopLocalTTS()
    isWordAudioPlaying.value = false
    return
  }
  const wordText = currentWord.value.word
  const audioUs = currentWord.value.audioUs || ''
  const audioUk = currentWord.value.audioUk || ''
  playWordTTS(wordText, audioUs, audioUk, accent.value, {
    onStart: () => { isWordAudioPlaying.value = true },
    onEnd: () => { isWordAudioPlaying.value = false },
    onError: () => { isWordAudioPlaying.value = false }
  })
}

const closeWordPopup = () => {
  stopLocalTTS()
  isWordAudioPlaying.value = false
  wordPopup.value?.close()
}

const selectSpeed = (rate) => {
  playbackSpeed.value = rate
  if (audioCtx) {
    audioCtx.playbackRate(rate)
  }
}

const collectWord = async (wordId) => {
  if (!wordId) return
  const collected = isWordCollected.value
  const res = collected ? await uncollect(1, wordId) : await collect(1, wordId)
  if (res.code === 0) {
    isWordCollected.value = !collected
    uni.showToast({ title: !collected ? t('diary.collect_success') : t('diary.uncollect_success'), icon: 'success' })
  }
}

// 句子收藏（使用 targetType=5 日记句子）
const collectSentence = async (item, index) => {
  const sentenceId = getSentenceId(item)
  if (!sentenceId) return
  const collected = !!sentenceCollectedMap.value[sentenceId]
  const res = collected ? await uncollect(5, sentenceId) : await collect(5, sentenceId)
  if (res.code === 0) {
    sentenceCollectedMap.value = {
      ...sentenceCollectedMap.value,
      [sentenceId]: !collected
    }
    uni.showToast({ title: !collected ? t('diary.collect_success') : t('diary.uncollect_success'), icon: 'success' })
  }
}

const loadSentenceCollectionState = async () => {
  const sentenceIds = subtitleList.value
    .map((item) => getSentenceId(item))
    .filter((id) => id > 0)
  if (sentenceIds.length === 0) {
    sentenceCollectedMap.value = {}
    return
  }
  const idSet = new Set(sentenceIds)
  const map = {}
  const maxPages = 20
  let page = 1
  try {
    while (page <= maxPages) {
      const res = await getCollectionList({ targetType: 5, page, pageSize: 100 })
      if (res.code !== 0 || !res.data) {
        sentenceCollectedMap.value = {}
        return
      }
      const list = Array.isArray(res.data.list) ? res.data.list : []
      for (const item of list) {
        const targetId = Number(item?.targetId || 0)
        if (targetId > 0 && idSet.has(targetId)) {
          map[targetId] = true
        }
      }
      if (Object.keys(map).length >= idSet.size) break
      if (list.length < 100) break
      page += 1
    }
    sentenceCollectedMap.value = map
  } catch (e) {
    sentenceCollectedMap.value = {}
  }
}

const loadDiary = async () => {
  const res = await findDiary(diaryId.value)
  if (res.code !== 0 || !res.data) {
    uni.showToast({ title: t('diary.load_failed'), icon: 'none' })
    return
  }
  const data = res.data
  diaryInfo.value = {
    id: getEntityId(data),
    name: data?.name || '',
    audioUs: data?.audioUs || '',
    audioUk: data?.audioUk || '',
    duration: Number(data?.duration || 0),
    imageI18n: data?.imageI18n || '',
    showTranslate: data?.showTranslate !== false,
    showEnglish: data?.showEnglish !== false,
    hasFullAuth: data?.hasFullAuth !== false,
    trialPercent: Number(data?.trialPercent || 8)
  }
  // 初始时长用数据库的（后续音频加载后以实际为准）
  audioDurationUs.value = diaryInfo.value.duration
  audioDurationUk.value = diaryInfo.value.duration
}

const loadSentences = async () => {
  const res = await getDiarySentenceList(diaryId.value)
  if (res.code === 0 && res.data) {
    subtitleList.value = res.data || []
    await loadSentenceCollectionState()
  }
}

const loadCollectionState = async () => {
  const res = await getCollectionList({ targetType: 4, page: 1, pageSize: 999 })
  if (res.code === 0 && res.data) {
    const list = Array.isArray(res.data.list) ? res.data.list : []
    isCollected.value = list.some(item => Number(item.targetId) === diaryId.value)
  }
}

const collectDiary = async () => {
  if (isCollected.value) {
    const res = await uncollect(4, diaryId.value)
    if (res.code === 0) {
      isCollected.value = false
      uni.showToast({ title: t('diary.uncollect_success'), icon: 'success' })
    }
  } else {
    const res = await collect(4, diaryId.value)
    if (res.code === 0) {
      isCollected.value = true
      uni.showToast({ title: t('diary.collect_success'), icon: 'success' })
    }
  }
}

const viewImage = () => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  let images = {}
  try {
    images = typeof diaryInfo.value.imageI18n === 'object' ? diaryInfo.value.imageI18n : JSON.parse(diaryInfo.value.imageI18n || '{}')
  } catch (e) {}
  const url = getExternalUrl(images[locale] || images['en'] || images['zh'] || '')
  if (!url) {
    uni.showToast({ title: t('diary.no_image'), icon: 'none' })
    return
  }
  uni.previewImage({
    urls: [url],
    current: url
  })
}

const toggleSentences = () => {
  controls.value.showSentences = !controls.value.showSentences
}

const toggleTranslate = () => {
  controls.value.showTranslate = !controls.value.showTranslate
}

const openMoreSheet = () => {
  morePopup.value?.open()
}

const closeMoreSheet = () => {
  morePopup.value?.close()
}

const goBack = () => {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/learning/diary' })
    }
  })
}

const handleBack = () => {
  closeMoreSheet()
  goBack()
}

onLoad((options) => {
  diaryId.value = Number(options?.id || 0)
  if (diaryId.value > 0) {
    initExternalDomain().then(() => {
      loadDiary().then(() => {
        loadSentences()
        loadCollectionState()
        initAudio()
      })
    })
  }
})

onUnmounted(() => {
  if (progressTimer) {
    clearInterval(progressTimer)
  }
  if (audioCtx) {
    audioCtx.destroy()
    audioCtx = null
  }
})

onUnload(() => {
  destroyTTS()
  if (audioCtx) {
    audioCtx.destroy()
    audioCtx = null
  }
})
</script>

<style scoped>
.diary-detail-container {
  min-height: 100vh;
  height: 100vh;
  background: #FFFFFF;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
  padding-top: calc(var(--status-bar-height, 0px) + 110rpx);
  box-sizing: border-box;
}

/* === 顶部固定信息区 === */
.top-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: calc(var(--status-bar-height, 0px) + 16rpx) 24rpx 16rpx;
  overflow: hidden;
}

.top-bar-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-radius: 0 0 24rpx 24rpx;
}

.top-bar-content {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.back-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255,255,255,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.back-btn:active {
  transform: scale(0.92);
}

.back-icon {
  font-size: 44rpx;
  color: #fff;
  line-height: 1;
}

.top-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.diary-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
}

.audio-duration {
  display: flex;
  align-items: center;
}

.duration-text {
  font-size: 22rpx;
  color: rgba(255,255,255,0.7);
  font-weight: 500;
  transition: all 0.3s ease;
}

.duration-playing .duration-text {
  color: #fff;
  animation: pulse-duration 0.8s ease-in-out infinite alternate;
}

@keyframes pulse-duration {
  0% { transform: scale(1); }
  100% { transform: scale(1.15); }
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 12rpx;
  flex-shrink: 0;
}

.accent-switch {
  display: flex;
  align-items: center;
  padding: 8rpx 16rpx;
  background: rgba(255,255,255,0.2);
  border-radius: 999rpx;
}

.accent-switching-text {
  font-size: 20rpx;
  color: #fff;
  font-weight: 600;
  animation: switching-pulse 1s ease-in-out infinite;
}

@keyframes switching-pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

.accent-text {
  font-size: 22rpx;
  color: rgba(255,255,255,0.6);
  font-weight: 600;
}

.accent-text.active {
  color: #fff;
}

.accent-divider {
  font-size: 22rpx;
  color: rgba(255,255,255,0.4);
  margin: 0 4rpx;
}

.image-btn {
  width: 60rpx;
  height: 60rpx;
  border-radius: 50%;
  background: rgba(255,255,255,0.22);
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-btn:active {
  transform: scale(0.92);
  background: rgba(255,255,255,0.3);
}

.image-btn-icon {
  width: 32rpx;
  height: 32rpx;
}

/* === 字幕滚动区 === */
.subtitle-section {
  height: calc(100vh - var(--status-bar-height, 0px) - 110rpx - 150rpx - env(safe-area-inset-bottom));
  margin-top: 16rpx;
  padding: 0 24rpx;
  box-sizing: border-box;
}

.subtitle-padding-top {
  height: 40rpx;
}

.subtitle-padding-bottom {
  height: 40rpx;
}

.sentence-row {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  padding: 20rpx 16rpx;
  margin-bottom: 8rpx;
  border-radius: 16rpx;
  transition: all 0.25s ease;
  scroll-margin-top: 200rpx;
}

.sentence-row.is-active {
  background: rgba(108, 91, 255, 0.06);
}

.sentence-audio-icon {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #F0EEFF;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 4rpx;
}

.audio-icon-char {
  font-size: 24rpx;
  color: #A9AECB;
}

.sentence-audio-icon.audio-playing {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  animation: pulse-icon 0.8s ease-in-out infinite alternate;
}

.sentence-audio-icon.audio-playing .audio-icon-char {
  color: #fff;
}

@keyframes pulse-icon {
  0% { transform: scale(1); }
  100% { transform: scale(1.12); }
}

.sentence-text-wrap {
  flex: 1;
}

.english-text {
  font-size: 28rpx;
  color: #2D2E4A;
  line-height: 1.8;
  font-weight: 500;
}

.highlight-word {
  color: #6D5BFF;
  font-weight: 700;
  border-bottom: 2rpx solid rgba(108, 91, 255, 0.3);
}

.translated-text {
  font-size: 24rpx;
  color: #6B6F8D;
  line-height: 1.6;
  margin-top: 8rpx;
}

/* === 底部控制栏 === */
.bottom-controls {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 18rpx 24rpx 14rpx;
  padding-bottom: calc(14rpx + env(safe-area-inset-bottom));
  background: #FFFFFF;
  border-top: 1rpx solid #F0EEFF;
  box-shadow: 0 -4rpx 16rpx rgba(108, 91, 255, 0.04);
}

.control-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6rpx;
}

.c-icon-circle {
  width: 76rpx;
  height: 76rpx;
  border-radius: 50%;
  background: #F5F3FF;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.25s ease;
}

.c-icon-circle:active {
  transform: scale(0.92);
}

.c-circle-on {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 4rpx 14rpx rgba(108, 91, 255, 0.35);
}

.c-icon-img {
  width: 36rpx;
  height: 36rpx;
}

.c-text {
  font-size: 22rpx;
  color: #A9AECB;
  font-weight: 500;
}

.c-text-active {
  color: #6D5BFF;
  font-weight: 600;
}

.play-btn .c-play-circle {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6rpx 20rpx rgba(108, 91, 255, 0.4);
}

.play-btn .c-play-circle:active {
  transform: scale(0.92);
}

.c-play-img {
  width: 44rpx;
  height: 44rpx;
  margin-left: 4rpx;
}

/* === 更多设置弹窗 === */
.more-sheet {
  border-radius: 32rpx 32rpx 0 0;
  padding-bottom: calc(180rpx + env(safe-area-inset-bottom));
}

.popup-handle {
  width: 60rpx;
  height: 6rpx;
  background: #D1D5DB;
  border-radius: 6rpx;
  margin: 16rpx auto;
}

.more-sheet-head {
  padding: 16rpx 32rpx;
  border-bottom: 1rpx solid #F0EEFF;
}

.more-sheet-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1A1B3A;
}

.more-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #F5F3FF;
}

.more-card:active {
  background: rgba(108, 91, 255, 0.04);
}

.more-card-main {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.more-card-title {
  font-size: 28rpx;
  color: #1A1B3A;
  font-weight: 600;
}

.more-card-desc {
  font-size: 22rpx;
  color: #A9AECB;
}

.more-state-pill {
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  background: #F5F3FF;
  font-size: 22rpx;
  color: #A9AECB;
  font-weight: 600;
}

.more-state-pill.on {
  background: rgba(108, 91, 255, 0.12);
  color: #6D5BFF;
}

.more-actions {
  display: flex;
  gap: 16rpx;
  padding: 24rpx 32rpx 28rpx;
  margin-top: 12rpx;
}

.more-action-btn {
  flex: 1;
  height: 76rpx;
  line-height: 76rpx;
  padding: 0;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  font-size: 26rpx;
  font-weight: 600;
  border: none;
  text-align: center;
  box-shadow: 0 4rpx 14rpx rgba(108, 91, 255, 0.32);
}

.more-action-btn.ghost {
  background: #F5F3FF;
  color: #6D5BFF;
  border: 1.5rpx solid rgba(108, 91, 255, 0.2);
  box-shadow: none;
}

.more-action-btn:active {
  transform: scale(0.97);
  opacity: 0.92;
}

/* === 句子收藏按钮 === */
.sentence-collect-btn {
  width: 48rpx;
  height: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 4rpx;
}

.sentence-collect-btn:active {
  transform: scale(0.88);
}

.collect-icon {
  font-size: 28rpx;
  color: #A9AECB;
}

/* === 重点单词弹窗 === */
.word-detail-box {
  padding: 24rpx 40rpx calc(140rpx + env(safe-area-inset-bottom)) 40rpx;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  background: #FFFFFF;
  border-top-left-radius: 28rpx;
  border-top-right-radius: 28rpx;
}

.w-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.w-title {
  font-size: 48rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.w-audio-circle {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border: none;
}

.w-audio-icon {
  font-size: 32rpx;
  color: #fff;
}

.w-audio-circle.playing {
  pointer-events: none;
  animation: pulse 0.6s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.12); opacity: 0.75; }
}

.w-phonetic {
  font-size: 28rpx;
  color: #6B6F8D;
  margin: 4rpx 0 20rpx 0;
}

.w-divider {
  height: 1rpx;
  background: rgba(108, 91, 255, 0.16);
  margin-bottom: 24rpx;
}

.w-exp {
  font-size: 32rpx;
  color: #1A1B3A;
  line-height: 1.6;
  margin-bottom: 36rpx;
}

.w-actions {
  display: flex;
  gap: 16rpx;
}

.w-cancel-btn {
  flex: 1;
  height: 88rpx;
  line-height: 88rpx;
  margin: 0;
  border-radius: 999rpx;
  border: 2rpx solid rgba(108, 91, 255, 0.32);
  background: #fff;
  color: #6D5BFF;
  font-size: 30rpx;
  font-weight: 700;
}

.w-collect {
  flex: 1;
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

/* === 锁定句子样式 === */
.sentence-locked {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  padding: 20rpx 16rpx;
}

.lock-time {
  font-size: 22rpx;
  color: #A9AECB;
  font-weight: 500;
}

.lock-icon {
  font-size: 40rpx;
  opacity: 0.4;
}

/* === 速率选择器 === */
.speed-hint {
  padding: 20rpx 32rpx 0;
  font-size: 24rpx;
  color: #6D5BFF;
  font-weight: 600;
}

.speed-card {
  margin-top: 14rpx;
}

.speed-disabled-title {
  color: #A9AECB;
}

.speed-unavailable-text {
  font-size: 20rpx;
  color: #A9AECB;
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

.speed-pill-disabled {
  color: #A9AECB;
  border-color: #E8E8ED;
  background: #F5F5F8;
}

.speed-pill-active {
  color: #fff;
  border-color: transparent;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 4rpx 12rpx rgba(108, 91, 255, 0.3);
}
</style>