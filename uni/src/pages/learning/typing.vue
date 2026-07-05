<template>
  <view class="typing-container">
    <!-- Navbar / 设置区域 -->
    <view class="top-nav">
      <view class="nav-left">
        <view class="nav-pill" @click="showCategorySheet">
          <text class="nav-text">{{ categoryName || t('typing.category', '分类') }}</text>
          <text class="arrow">▼</text>
        </view>
        <view v-if="chapters.length > 0" class="nav-pill" @click="showChapterSheet">
          <text class="nav-text">{{ chapterName || t('typing.chapter', '章节') }}</text>
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
        
        <view class="audio-btn" @click="playWordAudio">
          <text>🔊</text>
        </view>
      </view>

      <view class="word-meta" v-if="settings.showPhonetic">
        <text>{{ settings.accent === 'US' ? currentWord.phoneticUs : currentWord.phoneticUk }}</text>
      </view>

      <view class="word-explanation" v-if="settings.showExplanation">
        <text>{{ localText(currentWord.explanation) }}</text>
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
        <view class="set-item highlight" @click="collectWord">{{ isCollected ? t('typing.uncollect') : t('typing.collect') }}</view>
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
const fallbackTexts = {
  'typing.error_count': '错误次数',
  'typing.stop': '停止跟打',
  'typing.start': '开始跟打',
  'typing.show_word': '显示单词',
  'typing.show_phonetic': '显示音标',
  'typing.show_explanation': '显示释义',
  'typing.accent': '发音口音',
  'typing.collect': '收藏单词',
  'typing.uncollect': '取消收藏',
  'typing.category': '分类',
  'typing.chapter': '章节',
  'typing.load_empty': '当前章节暂无单词',
  'typing.collect_success': '收藏成功',
  'typing.uncollect_success': '已取消收藏',
  'typing.word_done': '完成',
  'typing.no_chapter': '暂无章节',
  'typing.audio_debug': '语音播放日志',
  'typing.audio_debug_clear': '清空',
  'typing.audio_ok': '成功',
  'typing.audio_fail': '失败',
  'typing.audio_skip': '跳过',
}

const t = (key, dt = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return dt || fallbackTexts[key] || key
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
const isCollected = ref(false)
const showSettingsSheet = () => settingsPopup.value.open()

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
  uni.showToast({ title: t('player.load_failed', '音频播放失败'), icon: 'none' })
}

const playAudio = (src, fallbackText = '') => {
  const audioSrc = normalizeAudioSource(src)
  const fallbackAudio = buildFallbackAudioSource(fallbackText)
  const isUnreachable = audioSrc && (isLocalAudioSource(audioSrc) || isPrivateNetworkAudioSource(audioSrc))
  const primary = audioSrc && !isUnreachable ? audioSrc : ''

  audioCandidates = []
  if (audioSrc && isUnreachable) {
    pushAudioDebugLog({ status: 'skip', source: 'db-local-skip', url: audioSrc, message: '本地/内网地址，直接跳过' })
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
    uni.showToast({ title: t('player.load_failed', '音频地址无效'), icon: 'none' })
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
          message: err?.errMsg || '未知错误'
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

const showChapterSheet = () => {
  if (chapters.value.length === 0) {
    uni.showToast({ title: t('typing.no_chapter'), icon: 'none' })
    return
  }
  uni.showActionSheet({
    itemList: chapters.value.map((item) => localText(item.name) || `#${item.id}`),
    success: async (res) => {
      const target = chapters.value[res.tapIndex]
      if (!target || target.id === selectedChapterId.value) {
        return
      }
      selectedChapterId.value = target.id
      currentWordIndex.value = 0
      await loadWords(selectedChapterId.value, selectedCategoryId.value)
      updateHeaderNames()
      persistProgress()
      resetStats()
    }
  })
}

const showCategorySheet = () => {
  if (categories.value.length === 0) {
    return
  }
  uni.showActionSheet({
    itemList: categories.value.map((item) => localText(item.name) || `#${item.id}`),
    success: async (res) => {
      const target = categories.value[res.tapIndex]
      if (!target || target.id === selectedCategoryId.value) {
        return
      }
      currentWordIndex.value = 0
      await selectCategory(target.id)
      persistProgress()
      resetStats()
    }
  })
}
</script>

<style scoped>
.typing-container { position: relative; height: 100vh; display: flex; flex-direction: column; background-color: #f7f9fc; }

.top-nav { display: flex; justify-content: space-between; padding: 30rpx; background: #fff; box-shadow: 0 2rpx 10rpx rgba(0,0,0,0.02); z-index: 10;}
.nav-left { display: flex; gap: 14rpx; align-items: center; }
.nav-pill { display: flex; align-items: center; gap: 8rpx; padding: 10rpx 18rpx; background: #f3f4f6; border-radius: 999rpx; }
.nav-text { font-size: 28rpx; font-weight: bold; }
.icon-settings { font-size: 40rpx; }

/* 单词卡片核心区域 */
.word-card { background: #fff; margin: 30rpx; padding: 60rpx 40rpx; border-radius: 20rpx; align-items: center; display: flex; flex-direction: column; box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.05); }
.word-main { display: flex; align-items: center; margin-bottom: 20rpx; }
.char-item { font-size: 80rpx; font-weight: bold; font-family: monospace; letter-spacing: 4rpx; }
.char-initial { color: #111; }
.char-pending { color: #ccc; }
.char-correct { color: #409eff; }
.char-wrong { color: #e53935; }
.audio-btn { margin-left: 20rpx; font-size: 40rpx; }
.word-meta { color: #888; margin-bottom: 10rpx; font-size: 28rpx; }
.word-explanation { color: #555; text-align: center; font-size: 30rpx; margin-top: 20rpx; line-height: 1.5; }

/* 防止作弊：打错时严重抖动与警示红屏 */
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-15rpx); }
  50% { transform: translateX(15rpx); }
  75% { transform: translateX(-15rpx); }
}
.shake-animation { animation: shake 0.4s ease-in-out; }

/* 统计与按钮区 */
.stats-row { display: flex; justify-content: space-between; align-items: center; padding: 0 40rpx; margin-bottom: 20rpx; }
.count-stat { font-size: 28rpx; color: #666; font-family: monospace; }
.toggle-typing-btn { background: #409eff; color: #fff; font-size: 28rpx; border-radius: 40rpx; margin: 0; padding: 0 40rpx; height: 60rpx; line-height: 60rpx; }
.toggle-typing-btn.active { background: #e53935; }

/* 句子区域 */
.sentence-list { flex: 1; padding: 0 30rpx; box-sizing: border-box; }
.sentence-item { background: #fff; padding: 30rpx; border-radius: 16rpx; margin-bottom: 20rpx; }
.sentence-en { font-size: 32rpx; font-weight: bold; margin-bottom: 15rpx; color: #333; display: flex; justify-content: space-between; }
.sentence-zh { font-size: 28rpx; color: #777; }

.audio-debug-box { margin: 0 30rpx 16rpx; background: #fff; border-radius: 14rpx; padding: 16rpx 18rpx; box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04); }
.audio-debug-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10rpx; }
.audio-debug-title { font-size: 24rpx; color: #334155; font-weight: 600; }
.audio-debug-clear { font-size: 22rpx; color: #2563eb; }
.audio-debug-list { display: flex; flex-direction: column; gap: 8rpx; }
.audio-debug-item { background: #f8fafc; border-radius: 10rpx; padding: 10rpx 12rpx; display: flex; flex-direction: column; gap: 4rpx; }
.audio-debug-status { font-size: 22rpx; font-weight: 600; }
.audio-debug-status.success { color: #15803d; }
.audio-debug-status.fail { color: #dc2626; }
.audio-debug-status.skip { color: #92400e; }
.audio-debug-source { font-size: 22rpx; color: #475569; }
.audio-debug-url { font-size: 20rpx; color: #64748b; line-height: 1.3; word-break: break-all; }
.audio-debug-msg { font-size: 20rpx; color: #b91c1c; line-height: 1.3; word-break: break-all; }

/* 浮动切换按钮 */
.nav-btn { position: fixed; top: 40%; width: 80rpx; height: 100rpx; background: rgba(0,0,0,0.3); color: #fff; display: flex; align-items: center; justify-content: center; font-size: 40rpx; z-index: 100; border-radius: 10rpx; }
.nav-btn { top: auto; bottom: 320rpx; }
.prev-btn { left: 0; border-top-left-radius: 0; border-bottom-left-radius: 0; }
.next-btn { right: 0; border-top-right-radius: 0; border-bottom-right-radius: 0; }

/* 前端绘制纯安全虚拟键盘 (核心防原声键盘中文作弊机制) */
.custom-keyboard { position: fixed; bottom: -450rpx; left: 0; width: 100%; height: 450rpx; background: #e0e4e8; transition: bottom 0.3s; padding: 20rpx 10rpx 40rpx 10rpx; box-sizing: border-box; display: flex; flex-direction: column; justify-content: space-between; z-index: 999; }
.keyboard-show { bottom: 0; }
.keyboard-row { display: flex; justify-content: center; gap: 10rpx; }
.key-btn { background: #fff; height: 90rpx; min-width: 60rpx; flex: 1; border-radius: 10rpx; display: flex; justify-content: center; align-items: center; font-size: 40rpx; font-family: monospace; font-weight: bold; color: #333; box-shadow: 0 2rpx 0 #999; text-transform: uppercase; }
.key-hover { background: #ccc; }

.settings-sheet { background: #fff; border-top-left-radius: 20rpx; border-top-right-radius: 20rpx; padding: 40rpx; padding-bottom: calc(40rpx + env(safe-area-inset-bottom) + 120rpx); max-height: 68vh; overflow-y: auto; box-sizing: border-box; }
.set-item { display: flex; justify-content: space-between; padding: 30rpx 0; border-bottom: 1px solid #eee; font-size: 32rpx; }
.highlight { color: #409eff; font-weight: bold; justify-content: center; }
</style>
