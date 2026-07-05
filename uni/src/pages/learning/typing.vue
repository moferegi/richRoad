<template>
  <view class="typing-container">
    <!-- Navbar / 设置区域 -->
    <view class="top-nav">
      <view class="nav-left">
        <view class="nav-pill" @click="openPickerSheet('category')">
          <text class="nav-text">{{ categoryName || t('typing.category') }}</text>
          <text class="arrow">▼</text>
        </view>
        <view v-if="chapters.length > 0" class="nav-pill" @click="openPickerSheet('chapter')">
          <text class="nav-text">{{ chapterName || t('typing.chapter') }}</text>
          <text class="arrow">▼</text>
        </view>
      </view>
      <view class="nav-right" @click="showSettingsSheet">
        <text class="icon-settings">⚙️</text>
      </view>
    </view>

    <!-- 单词主要展示区 -->
    <view class="word-card" :class="{ 'shake-animation': isWrong }">
      <view class="word-main">
        <!-- 重点：默写开关开启时，隐藏未打出的字母，改为下划线。关闭时全灰展示 -->
        <text 
          v-for="(char, index) in currentWord.word" 
          :key="index"
          class="char-item"
          :class="getCharClass(index)"
        >
          <!-- 如果开启默写（关闭显示），且还没打到这个字，显示下划线 -->
          {{ (!settings.showWord && index >= typedChars.length) ? '_' : char }}
        </text>
        
        <view class="word-actions">
          <view class="audio-btn" @click="playWordAudio">
            <text>🔊</text>
          </view>
        </view>
      </view>

      <view class="word-meta" v-if="settings.showPhonetic">
        <text>{{ settings.accent === 'US' ? currentWord.phoneticUs : currentWord.phoneticUk }}</text>
      </view>

      <view class="word-explanation" v-if="settings.showExplanation">
        <text>{{ localText(currentWord.explanation) }}</text>
      </view>

      <view class="word-collect-row">
        <view class="inline-collect-btn" @click="collectWord">
          <text>{{ isCollected ? t('typing.uncollect') : t('typing.collect') }}</text>
        </view>
      </view>
    </view>

    <!-- 跟打控制与数据行 -->
    <view class="stats-row">
      <text class="count-stat">❌ {{ t('typing.error_count') }}: {{ errorCount }}</text>
      <button 
        class="toggle-typing-btn" 
        :class="{ active: isTypingMode }" 
        @click="toggleTyping"
      >
        {{ isTypingMode ? t('typing.stop') : t('typing.start') }}
      </button>
    </view>

    <view v-if="audioDebugLogs.length > 0" class="audio-debug-box">
      <view class="audio-debug-head">
        <text class="audio-debug-title">{{ t('typing.audio_debug') }}</text>
        <text class="audio-debug-clear" @click="clearAudioDebugLogs">{{ t('typing.audio_debug_clear') }}</text>
      </view>
      <view class="audio-debug-list">
        <view class="audio-debug-item" v-for="item in audioDebugLogs" :key="item.id">
          <text class="audio-debug-status" :class="item.status">{{ item.statusLabel }}</text>
          <text class="audio-debug-source">{{ item.source }}</text>
          <text class="audio-debug-url">{{ item.url }}</text>
          <text v-if="item.message" class="audio-debug-msg">{{ item.message }}</text>
        </view>
      </view>
    </view>

    <!-- 造句列表区 (需要自适应高度并在键盘弹出时滚动) -->
    <scroll-view class="sentence-list" scroll-y :style="{ paddingBottom: isTypingMode ? '450rpx' : '0' }">
      <view class="sentence-item" v-for="(sentence, index) in currentWord.sentences" :key="index">
        <view class="sentence-en">
          <text>{{ sentence.source }}</text>
          <text class="sentence-audio" @click="playSentenceAudio(sentence)">🔊</text>
        </view>
        <view class="sentence-zh" v-if="settings.showExplanation">
          <text>{{ localText(sentence.translate) }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 左右切换悬浮按钮 -->
    <view class="nav-btn prev-btn" @click="prevWord"> <text>❮</text> </view>
    <view class="nav-btn next-btn" @click="nextWord"> <text>❯</text> </view>

    <!-- 自定义防串扰 26 键英语键盘 (原生键盘容易激活输入法联想和中文，坚决摒弃) -->
    <view class="custom-keyboard" :class="{ 'keyboard-show': isTypingMode }">
      <view class="keyboard-row" v-for="(row, rIndex) in keyboardLayout" :key="rIndex">
        <view 
          class="key-btn" 
          v-for="key in row" 
          :key="key" 
          @click="onKeyPress(key)"
          hover-class="key-hover"
        >
          {{ key }}
        </view>
      </view>
    </view>

    <!-- 功能设置弹窗 (ActionSheet 等价表现) -->
    <uni-popup ref="pickerPopup" type="bottom">
      <view class="picker-sheet">
        <view class="picker-title">{{ pickerType === 'category' ? t('typing.select_category') : t('typing.select_chapter') }}</view>
        <scroll-view class="picker-list" scroll-y>
          <view
            v-for="item in pickerOptions"
            :key="item.id"
            class="picker-item"
            :class="{ active: item.active }"
            @click="selectPickerItem(item)"
          >
            {{ item.label }}
          </view>
        </scroll-view>
        <button class="picker-cancel" @click="closePickerSheet">{{ t('common.cancel') }}</button>
      </view>
    </uni-popup>

    <uni-popup ref="settingsPopup" type="bottom">
      <view class="settings-sheet">
        <view class="set-item"><text>{{ t('typing.show_word') }}</text><switch :checked="settings.showWord" @change="settings.showWord = $event.detail.value" /></view>
        <view class="set-item"><text>{{ t('typing.show_phonetic') }}</text><switch :checked="settings.showPhonetic" @change="settings.showPhonetic = $event.detail.value" /></view>
        <view class="set-item"><text>{{ t('typing.show_explanation') }}</text><switch :checked="settings.showExplanation" @change="settings.showExplanation = $event.detail.value" /></view>
        <view class="set-item"><text>{{ t('typing.accent') }}</text>
          <radio-group @change="settings.accent = $event.detail.value">
            <label><radio value="US" :checked="settings.accent==='US'"/> US</label>
            <label><radio value="UK" :checked="settings.accent==='UK'"/> UK</label>
          </radio-group>
        </view>
        <view class="set-item highlight" @click="closeSettingsSheet">{{ t('common.cancel') }}</view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { collect, uncollect, getCollectionList, getCategoryList, getChapterList, getWordList, getWordProgress, reportWordError, saveWordProgress } from '@/api/learning.js'

const langStore = useLangStore()

const t = (key) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return key
}
const localText = (val) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  return i18nLocalText(val, locale)
}

const currentWord = ref({
  id: 0,
  word: '-',
  phoneticUs: '',
  phoneticUk: '',
  audioUs: '',
  audioUk: '',
  explanation: '',
  sentences: []
})

const categoryName = ref('')
const chapterName = ref('')
const categories = ref([])
const chapters = ref([])
const wordList = ref([])
const selectedCategoryId = ref(0)
const selectedChapterId = ref(0)
const currentWordIndex = ref(0)

// 系统配置
const settings = ref({
  showWord: true,        // 显示单词原文(不开就是下划线默写模式)
  showPhonetic: true,    // 显示音标
  showExplanation: true, // 显示释义
  accent: 'US',          // 美音/英音
})

const settingsPopup = ref(null)
const pickerPopup = ref(null)
const pickerType = ref('category')
const pickerOptions = ref([])
const isCollected = ref(false)
const showSettingsSheet = () => settingsPopup.value.open()
const closeSettingsSheet = () => settingsPopup.value?.close()

const refreshWordCollectedState = async () => {
  const wordId = Number(currentWord.value.id || 0)
  if (!wordId) {
    isCollected.value = false
    return
  }
  try {
    const res = await getCollectionList({ targetType: 1, page: 1, pageSize: 200 })
    if (res.code !== 0 || !res.data) {
      isCollected.value = false
      return
    }
    const list = Array.isArray(res.data.list) ? res.data.list : []
    isCollected.value = list.some((item) => Number(item?.targetId || 0) === wordId)
  } catch (error) {
    isCollected.value = false
  }
}

const collectWord = async () => {
  const wordId = Number(currentWord.value.id || 0)
  if (!wordId) {
    return
  }
  const res = isCollected.value ? await uncollect(1, wordId) : await collect(1, wordId)
  if (res.code === 0) {
    isCollected.value = !isCollected.value
    uni.showToast({ title: isCollected.value ? t('typing.collect_success') : t('typing.uncollect_success'), icon: 'success' })
  }
}

// 跟打逻辑与互动参数
const isTypingMode = ref(false)
const typedChars = ref('') // 当前按正确的拼接缓存
const isWrong = ref(false) // 触发错误动画
const errorCount = ref(0)  // 错误次数

// 安全键盘排版 QWERTY
const keyboardLayout = [
  ['q','w','e','r','t','y','u','i','o','p'],
  ['a','s','d','f','g','h','j','k','l'],
  ['z','x','c','v','b','n','m','-'] // '-' could be space or generic symbol
]

const toggleTyping = () => {
  isTypingMode.value = !isTypingMode.value
  typedChars.value = '' // 重新开始清空输入
}

// 核心键盘敲击校验
const onKeyPress = (key) => {
  if (!isTypingMode.value || isWrong.value) return

  const targetWord = currentWord.value.word.toLowerCase()
  const currentIndex = typedChars.value.length
  
  if (key === targetWord[currentIndex]) {
    // 按对了
    typedChars.value += key
    // 判断是否全打了
    if (typedChars.value.length === targetWord.length) {
      uni.showToast({ title: t('typing.word_done'), icon: 'success' })
      setTimeout(() => {
        nextWord()
      }, 800)
    }
  } else {
    const wordId = Number(currentWord.value.id || 0)
    if (wordId > 0) {
      reportWordError({
        wordId,
        categoryId: selectedCategoryId.value,
        chapterId: selectedChapterId.value,
        wrongIndex: currentIndex,
        expectedChar: targetWord[currentIndex] || '',
        inputChar: key
      }).catch(() => {})
    }

    // 扣错报错逻辑
    isWrong.value = true
    errorCount.value++
    typedChars.value = '' // 清空重来

    setTimeout(() => {
      isWrong.value = false
    }, 600) // 动画抖动 0.6s 后恢复
  }
}

// 获取各字母显示颜色
const getCharClass = (index) => {
  if (typedChars.value.length === 0) return 'char-initial'
  if (isWrong.value) return 'char-wrong' // 飘红
  if (index < typedChars.value.length) return 'char-correct' // 按对的是蓝色
  return 'char-pending' // 尚未按到的是灰色
}

// 上下首切换逻辑
const nextWord = () => {
  if (wordList.value.length === 0) {
    return
  }
  currentWordIndex.value = (currentWordIndex.value + 1) % wordList.value.length
  syncCurrentWord()
  persistProgress()
  resetStats()
}
const prevWord = () => {
  if (wordList.value.length === 0) {
    return
  }
  currentWordIndex.value = currentWordIndex.value <= 0 ? wordList.value.length - 1 : currentWordIndex.value - 1
  syncCurrentWord()
  persistProgress()
  resetStats()
}
const resetStats = () => {
  errorCount.value = 0
  typedChars.value = ''
  isWrong.value = false
}

let audioCtx = null
let audioCandidates = []
let currentAudioAttempt = null
let currentAudioLoggedSuccess = false
let audioDebugSeed = 1
const audioDebugLogs = ref([])

const pushAudioDebugLog = ({ status, source, url, message = '' }) => {
  const statusKey = status === 'success' ? 'typing.audio_ok' : status === 'skip' ? 'typing.audio_skip' : 'typing.audio_fail'
  audioDebugLogs.value = [
    {
      id: `${Date.now()}_${audioDebugSeed++}`,
      status,
      statusLabel: t(statusKey),
      source: String(source || ''),
      url: String(url || ''),
      message: String(message || '')
    },
    ...audioDebugLogs.value
  ].slice(0, 8)
}

const clearAudioDebugLogs = () => {
  audioDebugLogs.value = []
}
onMounted(() => {
  loadTypingData()
})
onShow(() => {
  loadTypingData()
})
onUnmounted(() => {
  persistProgress()
  if (audioCtx) {
    audioCtx.destroy()
    audioCtx = null
  }
})

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

const buildFallbackAudioSource = (text) => {
  const source = String(text || '').trim()
  if (!source) {
    return ''
  }
  const type = settings.value.accent === 'UK' ? '2' : '1'
  return `https://dict.youdao.com/dictvoice?audio=${encodeURIComponent(source)}&type=${type}&le=en`
}

const isLocalAudioSource = (src) => {
  return /^https?:\/\/(127\.0\.0\.1|localhost)/i.test(src)
}

const isPrivateNetworkAudioSource = (src) => {
  return /^https?:\/\/(10\.|192\.168\.|172\.(1[6-9]|2\d|3[0-1])\.)/i.test(src)
}

const playNextAudioCandidate = () => {
  while (audioCandidates.length > 0) {
    const next = audioCandidates.shift()
    if (!next || !next.url) {
      continue
    }
    currentAudioAttempt = next
    currentAudioLoggedSuccess = false
    audioCtx.src = next.url
    audioCtx.play()
    return
  }
  currentAudioAttempt = null
  currentAudioLoggedSuccess = false
  uni.showToast({ title: t('typing.audio_play_failed'), icon: 'none' })
}

const playAudio = (src, fallbackText = '') => {
  const audioSrc = normalizeAudioSource(src)
  const fallbackAudio = buildFallbackAudioSource(fallbackText)
  const isUnreachable = audioSrc && (isLocalAudioSource(audioSrc) || isPrivateNetworkAudioSource(audioSrc))
  const primary = audioSrc && !isUnreachable ? audioSrc : ''

  audioCandidates = []
  if (audioSrc && isUnreachable) {
    pushAudioDebugLog({ status: 'skip', source: 'db-local-skip', url: audioSrc, message: t('typing.audio_local_skip') })
  }
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
    uni.showToast({ title: t('typing.audio_invalid_url'), icon: 'none' })
    return
  }

  if (!audioCtx) {
    audioCtx = uni.createInnerAudioContext()
    audioCtx.onPlay(() => {
      if (currentAudioAttempt && !currentAudioLoggedSuccess) {
        currentAudioLoggedSuccess = true
        pushAudioDebugLog({ status: 'success', source: currentAudioAttempt.source, url: currentAudioAttempt.url })
      }
    })
    audioCtx.onError((err) => {
      if (currentAudioAttempt) {
        pushAudioDebugLog({
          status: 'fail',
          source: currentAudioAttempt.source,
          url: currentAudioAttempt.url,
          message: err?.errMsg || t('typing.audio_unknown_error')
        })
      }
      playNextAudioCandidate()
    })
  }
  playNextAudioCandidate()
}

const playWordAudio = () => {
  const primary = settings.value.accent === 'US' ? currentWord.value.audioUs : currentWord.value.audioUk
  const fallback = settings.value.accent === 'US' ? currentWord.value.audioUk : currentWord.value.audioUs
  playAudio(primary || fallback, currentWord.value.word)
}

const playSentenceAudio = (sentence) => {
  const primary = settings.value.accent === 'US'
    ? (sentence.audioUs || currentWord.value.audioUs)
    : (sentence.audioUk || currentWord.value.audioUk)
  const fallback = settings.value.accent === 'US'
    ? (sentence.audioUk || currentWord.value.audioUk)
    : (sentence.audioUs || currentWord.value.audioUs)
  playAudio(primary || fallback, sentence?.source || currentWord.value.word)
}

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const normalizeWord = (item) => ({
  ...item,
  id: getEntityId(item),
  word: String(item?.word || ''),
  phoneticUs: String(item?.phoneticUs || ''),
  phoneticUk: String(item?.phoneticUk || ''),
  audioUs: String(item?.audioUs || ''),
  audioUk: String(item?.audioUk || ''),
  explanation: item?.explanation || '',
  sentences: Array.isArray(item?.sentences) ? item.sentences : []
})

const updateHeaderNames = () => {
  const category = categories.value.find((item) => item.id === selectedCategoryId.value)
  const chapter = chapters.value.find((item) => item.id === selectedChapterId.value)
  categoryName.value = category ? localText(category.name) : ''
  chapterName.value = chapter ? localText(chapter.name) : ''
}

const syncCurrentWord = () => {
  if (wordList.value.length === 0) {
    currentWord.value = {
      id: 0,
      word: '-',
      phoneticUs: '',
      phoneticUk: '',
      audioUs: '',
      audioUk: '',
      explanation: '',
      sentences: []
    }
    return
  }
  currentWord.value = wordList.value[currentWordIndex.value]
  void refreshWordCollectedState()
}

const persistProgress = () => {
  if (!selectedCategoryId.value) {
    return
  }
  saveWordProgress({
    categoryId: selectedCategoryId.value,
    chapterId: selectedChapterId.value || 0,
    wordIndex: currentWordIndex.value
  }).catch(() => {})
}

const loadCategories = async () => {
  const res = await getCategoryList({ page: 1, pageSize: 200 })
  if (res.code !== 0 || !res.data) {
    categories.value = []
    return
  }
  const list = Array.isArray(res.data.list) ? res.data.list : []
  categories.value = list.map((item) => ({ ...item, id: getEntityId(item) })).filter((item) => item.id > 0)
}

const loadChapters = async (categoryId) => {
  const res = await getChapterList({ page: 1, pageSize: 200, categoryId })
  if (res.code !== 0 || !res.data) {
    chapters.value = []
    return
  }
  const list = Array.isArray(res.data.list) ? res.data.list : []
  chapters.value = list.map((item) => ({ ...item, id: getEntityId(item) })).filter((item) => item.id > 0)
}

const loadWords = async (chapterId, categoryId = 0) => {
  const params = { page: 1, pageSize: 500 }
  if (chapterId > 0) {
    params.chapterId = chapterId
  } else if (categoryId > 0) {
    params.categoryId = categoryId
  }

  const res = await getWordList(params)
  if (res.code !== 0 || !res.data) {
    wordList.value = []
    syncCurrentWord()
    return
  }
  const list = Array.isArray(res.data.list) ? res.data.list : []
  wordList.value = list.map(normalizeWord)
  if (wordList.value.length === 0) {
    uni.showToast({ title: t('typing.load_empty'), icon: 'none' })
  }
  syncCurrentWord()
}

const selectCategory = async (categoryId, preferredChapterId = 0, preferredWordIndex = 0) => {
  if (!categoryId) {
    return
  }
  selectedCategoryId.value = categoryId
  await loadChapters(categoryId)
  if (chapters.value.length === 0) {
    selectedChapterId.value = 0
    await loadWords(0, categoryId)
    currentWordIndex.value = wordList.value.length > 0
      ? Math.min(Math.max(0, Number(preferredWordIndex || 0)), wordList.value.length - 1)
      : 0
    syncCurrentWord()
    updateHeaderNames()
    return
  }

  const chapterExists = chapters.value.some((item) => item.id === preferredChapterId)
  selectedChapterId.value = chapterExists ? preferredChapterId : chapters.value[0].id
  await loadWords(selectedChapterId.value, categoryId)

  if (wordList.value.length > 0) {
    const maxIndex = wordList.value.length - 1
    currentWordIndex.value = Math.min(Math.max(0, Number(preferredWordIndex || 0)), maxIndex)
    syncCurrentWord()
  } else {
    currentWordIndex.value = 0
  }

  updateHeaderNames()
}

const loadTypingData = async () => {
  const progressRes = await getWordProgress()
  const progress = progressRes.code === 0 && progressRes.data ? progressRes.data : {}

  await loadCategories()
  if (categories.value.length === 0) {
    selectedCategoryId.value = 0
    selectedChapterId.value = 0
    wordList.value = []
    currentWordIndex.value = 0
    syncCurrentWord()
    categoryName.value = ''
    chapterName.value = ''
    return
  }

  const progressCategoryId = Number(progress?.categoryId || 0)
  const progressChapterId = Number(progress?.chapterId || 0)
  const progressWordIndex = Number(progress?.wordIndex || 0)
  const hasProgressCategory = categories.value.some((item) => item.id === progressCategoryId)
  const categoryId = hasProgressCategory ? progressCategoryId : categories.value[0].id

  await selectCategory(categoryId, progressChapterId, progressWordIndex)
}

const openPickerSheet = (type) => {
  isTypingMode.value = false
  if (type === 'chapter' && chapters.value.length === 0) {
    uni.showToast({ title: t('typing.no_chapter'), icon: 'none' })
    return
  }
  if (type === 'category' && categories.value.length === 0) {
    return
  }

  pickerType.value = type
  const source = type === 'category' ? categories.value : chapters.value
  pickerOptions.value = source.map((item) => ({
    id: item.id,
    label: localText(item.name) || `#${item.id}`,
    active: type === 'category' ? item.id === selectedCategoryId.value : item.id === selectedChapterId.value
  }))
  pickerPopup.value?.open()
}

const closePickerSheet = () => {
  pickerPopup.value?.close()
}

const selectPickerItem = async (item) => {
  if (!item?.id) {
    return
  }
  if (pickerType.value === 'chapter') {
    if (item.id === selectedChapterId.value) {
      closePickerSheet()
      return
    }
    selectedChapterId.value = item.id
    currentWordIndex.value = 0
    await loadWords(selectedChapterId.value, selectedCategoryId.value)
    updateHeaderNames()
    persistProgress()
    resetStats()
    closePickerSheet()
    return
  }

  if (item.id === selectedCategoryId.value) {
    closePickerSheet()
    return
  }
  currentWordIndex.value = 0
  await selectCategory(item.id)
  persistProgress()
  resetStats()
  closePickerSheet()
}
</script>

<style scoped>
.typing-container {
  position: relative;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, #f8f4e7 0%, #f4efe1 100%);
}

.top-nav {
  display: flex;
  justify-content: space-between;
  padding: 24rpx;
  background: rgba(255, 253, 248, 0.95);
  border-bottom: 1rpx solid rgba(146, 64, 14, 0.12);
}

.nav-left { display: flex; gap: 12rpx; align-items: center; }

.nav-pill {
  display: flex;
  align-items: center;
  gap: 8rpx;
  height: 60rpx;
  padding: 0 20rpx;
  border-radius: 999rpx;
  background: #fff;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
}

.nav-text { font-size: 26rpx; font-weight: 700; color: #7c2d12; }
.arrow { color: rgba(124, 45, 18, 0.6); font-size: 20rpx; }

.nav-right {
  width: 62rpx;
  height: 62rpx;
  border-radius: 18rpx;
  background: #fff;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-settings { font-size: 34rpx; }

.word-card {
  background: rgba(255, 253, 248, 0.96);
  margin: 22rpx;
  padding: 50rpx 36rpx;
  border-radius: 24rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 14rpx 32rpx rgba(120, 53, 15, 0.1);
}

.word-main { display: flex; align-items: center; justify-content: center; margin-bottom: 14rpx; }
.char-item { font-size: 78rpx; font-weight: 700; font-family: monospace; letter-spacing: 3rpx; }
.char-initial { color: #0f172a; }
.char-pending { color: #cbd5e1; }
.char-correct { color: #0f766e; }
.char-wrong { color: #dc2626; }
.word-actions { display: flex; align-items: center; gap: 10rpx; margin-left: 18rpx; }
.audio-btn { font-size: 40rpx; }
.word-collect-row {
  margin-top: 16rpx;
  width: 100%;
  display: flex;
  justify-content: center;
}
.inline-collect-btn {
  height: 56rpx;
  padding: 0 18rpx;
  border-radius: 999rpx;
  background: rgba(15, 118, 110, 0.14);
  border: 1rpx solid rgba(15, 118, 110, 0.22);
  color: #0f766e;
  font-size: 22rpx;
  display: flex;
  align-items: center;
}
.word-meta { color: #78716c; margin-bottom: 10rpx; font-size: 28rpx; }
.word-explanation { color: #334155; text-align: center; font-size: 29rpx; margin-top: 12rpx; line-height: 1.55; }

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-15rpx); }
  50% { transform: translateX(15rpx); }
  75% { transform: translateX(-15rpx); }
}
.shake-animation { animation: shake 0.4s ease-in-out; }

.stats-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24rpx;
  margin-bottom: 16rpx;
}

.count-stat { font-size: 27rpx; color: #78716c; font-family: monospace; }

.toggle-typing-btn {
  height: 66rpx;
  line-height: 66rpx;
  margin: 0;
  padding: 0 36rpx;
  border-radius: 999rpx;
  font-size: 27rpx;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
}

.toggle-typing-btn.active { background: linear-gradient(120deg, #dc2626 0%, #f97316 100%); }

.sentence-list { flex: 1; padding: 0 22rpx; box-sizing: border-box; }

.sentence-item {
  background: rgba(255, 253, 248, 0.96);
  padding: 24rpx;
  border-radius: 18rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.16);
  margin-bottom: 14rpx;
}

.sentence-en {
  font-size: 31rpx;
  font-weight: 700;
  margin-bottom: 10rpx;
  color: #1e293b;
  display: flex;
  justify-content: space-between;
}

.sentence-zh { font-size: 26rpx; color: #64748b; line-height: 1.5; }

.audio-debug-box {
  margin: 0 22rpx 14rpx;
  background: #fffdf8;
  border-radius: 16rpx;
  padding: 16rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.16);
}
.audio-debug-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10rpx; }
.audio-debug-title { font-size: 24rpx; color: #334155; font-weight: 700; }
.audio-debug-clear { font-size: 22rpx; color: #0f766e; }
.audio-debug-list { display: flex; flex-direction: column; gap: 8rpx; }
.audio-debug-item { background: #f8fafc; border-radius: 10rpx; padding: 10rpx 12rpx; display: flex; flex-direction: column; gap: 4rpx; }
.audio-debug-status { font-size: 22rpx; font-weight: 600; }
.audio-debug-status.success { color: #15803d; }
.audio-debug-status.fail { color: #dc2626; }
.audio-debug-status.skip { color: #92400e; }
.audio-debug-source { font-size: 22rpx; color: #475569; }
.audio-debug-url { font-size: 20rpx; color: #64748b; line-height: 1.3; word-break: break-all; }
.audio-debug-msg { font-size: 20rpx; color: #b91c1c; line-height: 1.3; word-break: break-all; }

.nav-btn {
  position: fixed;
  bottom: 340rpx;
  width: 76rpx;
  height: 92rpx;
  background: rgba(15, 118, 110, 0.66);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 38rpx;
  z-index: 100;
}

.prev-btn { left: 0; border-top-left-radius: 0; border-bottom-left-radius: 0; }
.next-btn { right: 0; border-top-right-radius: 0; border-bottom-right-radius: 0; }

.custom-keyboard {
  position: fixed;
  bottom: -450rpx;
  left: 0;
  width: 100%;
  height: 450rpx;
  background: #e7ecee;
  transition: bottom 0.3s;
  padding: 20rpx 10rpx 40rpx;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  z-index: 999;
}
.keyboard-show { bottom: 0; }
.keyboard-row { display: flex; justify-content: center; gap: 10rpx; }

.key-btn {
  background: #fff;
  height: 90rpx;
  min-width: 60rpx;
  flex: 1;
  border-radius: 12rpx;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 40rpx;
  font-family: monospace;
  font-weight: 700;
  color: #334155;
  box-shadow: 0 2rpx 0 #94a3b8;
  text-transform: uppercase;
}

.key-hover { background: #dbe3e6; }

.settings-sheet {
  background: #fffdf8;
  border-top-left-radius: 22rpx;
  border-top-right-radius: 22rpx;
  padding: 34rpx;
  padding-bottom: calc(34rpx + env(safe-area-inset-bottom) + 120rpx);
  max-height: 68vh;
  overflow-y: auto;
  box-sizing: border-box;
}

.picker-sheet {
  background: #fffdf8;
  border-top-left-radius: 22rpx;
  border-top-right-radius: 22rpx;
  padding: 28rpx;
  padding-bottom: calc(28rpx + env(safe-area-inset-bottom) + 120rpx);
  max-height: 68vh;
  overflow-y: auto;
  box-sizing: border-box;
}

.picker-title {
  text-align: center;
  font-size: 30rpx;
  font-weight: 700;
  color: #7c2d12;
  margin-bottom: 14rpx;
}

.picker-list {
  max-height: 46vh;
}

.picker-item {
  height: 84rpx;
  border-radius: 14rpx;
  padding: 0 22rpx;
  display: flex;
  align-items: center;
  margin-bottom: 10rpx;
  color: #334155;
  background: #fff;
  border: 1rpx solid rgba(20, 184, 166, 0.16);
}

.picker-item.active {
  color: #0f766e;
  border-color: rgba(15, 118, 110, 0.34);
  background: rgba(15, 118, 110, 0.08);
}

.picker-cancel {
  margin-top: 10rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(146, 64, 14, 0.14);
  background: #fff;
  color: #7c2d12;
}

.set-item {
  display: flex;
  justify-content: space-between;
  padding: 28rpx 0;
  border-bottom: 1px solid rgba(146, 64, 14, 0.12);
  font-size: 30rpx;
  color: #334155;
}

.highlight { color: #0f766e; font-weight: 700; justify-content: center; }
</style>
