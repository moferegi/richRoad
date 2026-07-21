<template>
  <view class="typing-page" @click="closeWordSearchResult">
    <!-- 渐变顶栏 -->
    <view class="top-bar">
      <view class="top-bar-bg"></view>
      <view class="top-bar-content">
        <view class="bar-selectors">
          <view class="selector" @click="openPickerSheet('category')">
            <text class="selector-text">{{ categoryName || t('typing.category') }}</text>
            <text class="selector-arrow">›</text>
          </view>
          <view v-if="chapters.length > 0" class="selector" @click="openPickerSheet('chapter')">
            <text class="selector-text">{{ chapterName || t('typing.chapter') }}</text>
            <text class="selector-arrow">›</text>
          </view>
        </view>
        <view class="bar-action" @click="showSettingsSheet">
          <text class="action-icon">⚙</text>
        </view>
      </view>
    </view>

    <!-- 单词主卡片 -->
    <view v-if="isPageLoading" class="word-hero-card skeleton-mode">
      <view class="sk-line sk-w1"></view>
      <view class="sk-line sk-w2"></view>
      <view class="sk-line sk-w3"></view>
    </view>
    <view v-else class="word-hero-card" :class="{ 'shake-animation': isWrong }">
      <view class="word-display">
        <text
          v-for="(char, index) in currentWord.word"
          :key="index"
          class="word-char"
          :class="getCharClass(index)"
        >{{ (!settings.showWord && index >= typedChars.length) ? '_' : char }}</text>
      </view>

      <view class="word-meta">
        <view class="audio-trigger" :class="{ playing: isAudioPlaying }" @click.stop="playWordAudio">
          <text class="audio-icon">{{ isAudioPlaying ? '◉' : '♪' }}</text>
        </view>
        <text class="word-phonetic" v-if="settings.showPhonetic">{{ settings.accent === 'US' ? currentWord.phoneticUs : currentWord.phoneticUk }}</text>
      </view>

      <text class="word-meaning" v-if="settings.showExplanation">{{ localText(currentWord.explanation) }}</text>

      <view class="word-footer">
        <view class="collect-btn" @click.stop="collectWord">
          <text class="collect-icon">{{ isCollected ? '★' : '☆' }}</text>
          <text class="collect-text">{{ isCollected ? t('typing.uncollect') : t('typing.collect') }}</text>
        </view>
      </view>
    </view>

    <!-- 工具卡片 -->
    <view class="tools-card" @click.stop>
      <template v-if="isPageLoading">
        <view class="tools-row">
          <view class="sk-line sk-btn"></view>
          <view class="sk-line sk-input"></view>
        </view>
      </template>
      <template v-else>
        <view class="tools-row">
          <view class="tool-btn" @click="openWordDrawer">
            <text class="tool-icon">☰</text>
            <text>{{ t('typing.word_list') }}</text>
          </view>
          <view class="search-field">
            <input
              v-model="wordSearchKeyword"
              class="search-input"
              :placeholder="t('typing.search_word_placeholder')"
              confirm-type="search"
              @confirm="runWordSearch"
            />
            <view v-if="wordSearchKeyword" class="search-clear" @click="clearWordSearchKeyword">
              <text>×</text>
            </view>
            <view class="search-go" @click="runWordSearch">
              <text>{{ t('typing.search_action') }}</text>
            </view>
          </view>
        </view>

        <view v-if="showWordSearchResult" class="search-results">
          <view v-if="wordSearchResults.length === 0" class="search-empty">
            {{ t('typing.search_empty') }}
          </view>
          <view
            v-for="item in wordSearchResults"
            :key="item.id"
            class="search-item"
            @click="selectSearchWord(item)"
          >
            <text class="search-item-word">{{ item.word }}</text>
            <text class="search-item-def">{{ localText(item.explanation) }}</text>
          </view>
        </view>
      </template>
    </view>

    <!-- 状态栏：错误计数 + 开始/停止 -->
    <view v-if="isPageLoading" class="control-bar skeleton-mode">
      <view class="sk-line sk-w2"></view>
      <view class="sk-line sk-btn"></view>
    </view>
    <view v-else class="control-bar">
      <view class="error-badge">
        <text class="error-count">{{ errorCount }}</text>
        <text class="error-label">{{ t('typing.error_count') }}</text>
      </view>
      <button
        class="mode-btn"
        :class="{ typing: isTypingMode }"
        @click="toggleTyping"
      >
        {{ isTypingMode ? t('typing.stop') : t('typing.start') }}
      </button>
    </view>

    <!-- 例句区 -->
    <view class="sentences" :style="{ paddingBottom: isTypingMode ? '450rpx' : '0' }">
      <view class="sentence-card" v-for="(sentence, index) in currentWord.sentences" :key="index">
        <view class="sentence-top">
          <text class="sentence-en">{{ sentence.source }}</text>
          <view class="sentence-audio" @click="playSentenceAudio(sentence)">
            <text class="sentence-audio-icon">♪</text>
          </view>
        </view>
        <text class="sentence-zh" v-if="settings.showExplanation">{{ localText(sentence.translate) }}</text>
      </view>
    </view>

    <!-- 左右切换悬浮按钮 -->
    <view class="float-nav prev" @click="handlePrevNav"><text>{{ prevNavLabel }}</text></view>
    <view class="float-nav next" @click="handleNextNav"><text>{{ nextNavLabel }}</text></view>

    <!-- 自定义 26 键键盘 -->
    <view class="keyboard" :class="{ show: isTypingMode }">
      <view class="kb-row" v-for="(row, rIndex) in keyboardLayout" :key="rIndex">
        <view
          class="kb-key"
          v-for="key in row"
          :key="key"
          @click="onKeyPress(key)"
          hover-class="kb-key-hover"
        >{{ key }}</view>
      </view>
    </view>

    <!-- 底部弹窗：分类/章节选择 -->
    <uni-popup ref="pickerPopup" type="bottom">
      <view class="sheet">
        <view class="sheet-handle"></view>
        <view class="sheet-title">{{ pickerType === 'category' ? tt('typing.select_category', 'typing.category') : tt('typing.select_chapter', 'typing.chapter') }}</view>
        <scroll-view class="sheet-list" scroll-y :scroll-into-view="pickerScrollTarget" :scroll-with-animation="true">
          <view
            v-for="item in pickerOptions"
            :key="item.id"
            :id="'picker-item-' + item.id"
            class="sheet-item"
            :class="{ active: item.active }"
            @click="selectPickerItem(item)"
          >
            <text class="sheet-item-label">{{ item.label }}</text>
            <text v-if="item.active" class="sheet-item-check">✓</text>
          </view>
        </scroll-view>
        <button class="sheet-cancel" @click="closePickerSheet">{{ t('cancel') }}</button>
      </view>
    </uni-popup>

    <!-- 底部弹窗：设置 -->
    <uni-popup ref="settingsPopup" type="bottom">
      <view class="sheet">
        <view class="sheet-handle"></view>
        <view class="sheet-title">⚙</view>
        <view class="set-row"><text>{{ t('typing.show_word') }}</text><switch :checked="settings.showWord" @change="settings.showWord = $event.detail.value" /></view>
        <view class="set-row"><text>{{ t('typing.show_phonetic') }}</text><switch :checked="settings.showPhonetic" @change="settings.showPhonetic = $event.detail.value" /></view>
        <view class="set-row"><text>{{ t('typing.show_explanation') }}</text><switch :checked="settings.showExplanation" @change="settings.showExplanation = $event.detail.value" /></view>
        <view class="set-row"><text>{{ t('typing.accent') }}</text>
          <radio-group @change="settings.accent = $event.detail.value">
            <label class="accent-label"><radio value="US" :checked="settings.accent==='US'"/> US</label>
            <label class="accent-label"><radio value="UK" :checked="settings.accent==='UK'"/> UK</label>
          </radio-group>
        </view>
        <button class="sheet-cancel" @click="closeSettingsSheet">{{ t('cancel') }}</button>
      </view>
    </uni-popup>

    <!-- 左侧抽屉：词表 -->
    <uni-popup ref="wordDrawerPopup" type="left">
      <view class="drawer">
        <view class="drawer-head">
          <text class="drawer-title">{{ t('typing.word_list') }}</text>
          <text class="drawer-close" @click="closeWordDrawer">×</text>
        </view>
        <scroll-view class="drawer-body" scroll-y :scroll-into-view="drawerScrollToId" :scroll-with-animation="true">
          <view
            v-for="(item, idx) in wordList"
            :key="item.id || idx"
            :id="'drawer-word-' + (item.id || idx)"
            class="drawer-item"
            :class="{ active: idx === currentWordIndex }"
            @click="selectWordFromDrawer(idx)"
          >
            <text class="drawer-word">{{ item.word }}</text>
            <text class="drawer-def">{{ localText(item.explanation) }}</text>
          </view>
        </scroll-view>
      </view>
    </uni-popup>
    <custom-tab-bar />
  </view>
</template>

<script setup>
import { computed, ref, onUnmounted } from 'vue'
import { onShow, onHide } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { collect, uncollect, findWord, getCollectionList, getCategoryList, getChapterList, getWordList, getWordProgress, reportWordError, saveWordProgress } from '@/api/learning.js'
import CustomTabBar from '@/components/custom-tab-bar/custom-tab-bar.vue'

const langStore = useLangStore()

const t = (key) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return key
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

const hasChapters = computed(() => chapters.value.length > 0 && selectedChapterId.value > 0)
const currentChapterIndex = computed(() => chapters.value.findIndex((item) => item.id === selectedChapterId.value))
const isAtFirstWord = computed(() => wordList.value.length > 0 && currentWordIndex.value <= 0)
const isAtLastWord = computed(() => wordList.value.length > 0 && currentWordIndex.value >= wordList.value.length - 1)
const prevNavLabel = computed(() => {
  if (hasChapters.value && isAtFirstWord.value) {
    return t('typing.prev_chapter')
  }
  return t('typing.prev_word')
})
const nextNavLabel = computed(() => {
  if (hasChapters.value && isAtLastWord.value) {
    return t('typing.next_chapter')
  }
  return t('typing.next_word')
})

// 系统配置
const settings = ref({
  showWord: true,        // 显示单词原文(不开就是下划线默写模式)
  showPhonetic: true,    // 显示音标
  showExplanation: true, // 显示释义
  accent: 'US',          // 美音/英音
})

const settingsPopup = ref(null)
const pickerPopup = ref(null)
const wordDrawerPopup = ref(null)
const pickerType = ref('category')
const pickerScrollTarget = ref('')
const pickerOptions = ref([])
const wordSearchKeyword = ref('')
const showWordSearchResult = ref(false)
const wordSearchResults = ref([])
const isCollected = ref(false)
const collectedWordIdSet = ref(new Set())
const collectionsLoadedAt = ref(0)
const typingInitialized = ref(false)
// 3.3.1 骨架屏 + 3.3.2 语音动画 + 3.3.3 离线发音 + 3.3.4 抽屉聚焦
const isPageLoading = ref(true)
const isAudioPlaying = ref(false)
let offlineFallbackText = ''
const drawerScrollToId = ref('')
const showSettingsSheet = () => settingsPopup.value.open()
const closeSettingsSheet = () => settingsPopup.value?.close()
const openWordDrawer = () => {
  // 3.3.4 自动聚焦当前单词
  wordDrawerPopup.value?.open()
  const targetId = currentWord.value.id || currentWordIndex.value
  drawerScrollToId.value = ''
  setTimeout(() => {
    drawerScrollToId.value = 'drawer-word-' + targetId
  }, 300)
}
const closeWordDrawer = () => wordDrawerPopup.value?.close()
const closeWordSearchResult = () => {
  showWordSearchResult.value = false
}
const clearWordSearchKeyword = () => {
  wordSearchKeyword.value = ''
  wordSearchResults.value = []
  showWordSearchResult.value = false
}

const loadWordCollections = async (force = false) => {
  const now = Date.now()
  if (!force && now - Number(collectionsLoadedAt.value || 0) < 20000) {
    return
  }
  try {
    const res = await getCollectionList({ targetType: 1, page: 1, pageSize: 200 })
    if (res.code !== 0 || !res.data) {
      return
    }
    const list = Array.isArray(res.data.list) ? res.data.list : []
    const nextSet = new Set()
    list.forEach((item) => {
      const targetId = Number(item?.targetId || 0)
      if (targetId > 0) {
        nextSet.add(targetId)
      }
    })
    collectedWordIdSet.value = nextSet
    collectionsLoadedAt.value = now
  } catch (error) {
  }
}

const refreshWordCollectedState = () => {
  const wordId = Number(currentWord.value.id || 0)
  if (!wordId) {
    isCollected.value = false
    return
  }
  isCollected.value = collectedWordIdSet.value.has(wordId)
}

const collectWord = async () => {
  const wordId = Number(currentWord.value.id || 0)
  if (!wordId) {
    return
  }
  const res = isCollected.value ? await uncollect(1, wordId) : await collect(1, wordId)
  if (res.code === 0) {
    const nextSet = new Set(collectedWordIdSet.value)
    if (isCollected.value) {
      nextSet.delete(wordId)
    } else {
      nextSet.add(wordId)
    }
    collectedWordIdSet.value = nextSet
    refreshWordCollectedState()
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
        currentWordIndex.value = (currentWordIndex.value + 1) % wordList.value.length
        syncCurrentWord()
        persistProgress()
        resetStats()
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

const jumpToChapter = async (offset, targetWordAtEnd = false) => {
  if (chapters.value.length === 0) {
    uni.showToast({ title: t('typing.no_chapter'), icon: 'none' })
    return
  }
  const currentIdx = chapters.value.findIndex((item) => item.id === selectedChapterId.value)
  if (currentIdx < 0) {
    return
  }
  const targetIdx = currentIdx + offset
  if (targetIdx < 0) {
    uni.showToast({ title: t('typing.chapter_first'), icon: 'none' })
    return
  }
  if (targetIdx >= chapters.value.length) {
    uni.showToast({ title: t('typing.chapter_last'), icon: 'none' })
    return
  }

  selectedChapterId.value = chapters.value[targetIdx].id
  await loadWords(selectedChapterId.value, selectedCategoryId.value)
  // 3.3.5 切换章节后重聚焦抽屉
  reFocusDrawerIfOpen()
  if (wordList.value.length > 0) {
    currentWordIndex.value = targetWordAtEnd ? wordList.value.length - 1 : 0
    syncCurrentWord()
  } else {
    currentWordIndex.value = 0
  }
  updateHeaderNames()
  persistProgress()
  resetStats()
}

const prevWord = () => {
  if (wordList.value.length === 0 || currentWordIndex.value <= 0) {
    return
  }
  currentWordIndex.value -= 1
  syncCurrentWord()
  persistProgress()
  resetStats()
}

const nextWord = () => {
  if (wordList.value.length === 0 || currentWordIndex.value >= wordList.value.length - 1) {
    return
  }
  currentWordIndex.value += 1
  syncCurrentWord()
  persistProgress()
  resetStats()
}

const handlePrevNav = () => {
  if (hasChapters.value && isAtFirstWord.value) {
    if (currentChapterIndex.value <= 0) {
      uni.showToast({ title: t('typing.chapter_first'), icon: 'none' })
      return
    }
    jumpToChapter(-1, true)
    return
  }
  prevWord()
}

const handleNextNav = () => {
  if (hasChapters.value && isAtLastWord.value) {
    if (currentChapterIndex.value >= chapters.value.length - 1) {
      uni.showToast({ title: t('typing.chapter_last'), icon: 'none' })
      return
    }
    jumpToChapter(1, false)
    return
  }
  nextWord()
}
const resetStats = () => {
  errorCount.value = 0
  typedChars.value = ''
  isWrong.value = false
}

const selectWordFromDrawer = (index) => {
  if (index < 0 || index >= wordList.value.length) {
    return
  }
  currentWordIndex.value = index
  syncCurrentWord()
  persistProgress()
  resetStats()
  closeWordDrawer()
}

const normalizeWordBindingIds = (source) => {
  if (!Array.isArray(source)) {
    return []
  }
  return source.map((id) => Number(id || 0)).filter((id) => id > 0)
}

const jumpToWordBySearchResult = async (wordItem) => {
  const wordID = Number(wordItem?.id || wordItem?.ID || 0)
  if (!wordID) {
    return
  }

  const detailRes = await findWord(wordID)
  if (detailRes.code !== 0 || !detailRes.data) {
    return
  }
  const detail = detailRes.data
  const categoryIds = normalizeWordBindingIds(detail.categoryIds)
  const chapterIds = normalizeWordBindingIds(detail.chapterIds)

  let targetCategoryId = categoryIds[0] || selectedCategoryId.value || 0
  if (!targetCategoryId) {
    await loadCategories()
    targetCategoryId = categories.value[0]?.id || 0
  }
  if (!targetCategoryId) {
    return
  }

  selectedCategoryId.value = targetCategoryId
  await loadChapters(targetCategoryId)

  let targetChapterId = 0
  if (chapters.value.length > 0) {
    targetChapterId = chapterIds.find((id) => chapters.value.some((item) => item.id === id)) || chapters.value[0].id
  }
  selectedChapterId.value = targetChapterId
  await loadWords(targetChapterId, targetCategoryId)

  let targetWordIndex = wordList.value.findIndex((item) => item.id === wordID)
  if (targetWordIndex < 0 && chapterIds.length > 1) {
    for (const chapterId of chapterIds) {
      if (!chapters.value.some((item) => item.id === chapterId) || chapterId === targetChapterId) {
        continue
      }
      selectedChapterId.value = chapterId
      await loadWords(chapterId, targetCategoryId)
      targetWordIndex = wordList.value.findIndex((item) => item.id === wordID)
      if (targetWordIndex >= 0) {
        break
      }
    }
  }

  if (targetWordIndex < 0) {
    selectedChapterId.value = 0
    await loadWords(0, targetCategoryId)
    targetWordIndex = wordList.value.findIndex((item) => item.id === wordID)
  }

  if (targetWordIndex < 0) {
    uni.showToast({ title: t('typing.search_empty'), icon: 'none' })
    return
  }

  currentWordIndex.value = targetWordIndex
  syncCurrentWord()
  updateHeaderNames()
  persistProgress()
  resetStats()
}

const selectSearchWord = async (item) => {
  if (!item || typeof item !== 'object') {
    return
  }
  showWordSearchResult.value = false
  wordSearchKeyword.value = item.word || ''
  await jumpToWordBySearchResult(item)
}

const runWordSearch = async () => {
  const keyword = String(wordSearchKeyword.value || '').trim()
  if (!keyword) {
    uni.showToast({ title: t('typing.search_no_keyword'), icon: 'none' })
    showWordSearchResult.value = false
    wordSearchResults.value = []
    return
  }

  const res = await getWordList({ page: 1, pageSize: 50, keyword })
  if (res.code !== 0 || !res.data) {
    showWordSearchResult.value = false
    wordSearchResults.value = []
    return
  }

  const list = Array.isArray(res.data.list) ? res.data.list : []
  wordSearchResults.value = list.map((item) => normalizeWord(item || {}))
  showWordSearchResult.value = true
}

let audioCtx = null
let audioCandidates = []
let currentAudioAttempt = null
let currentAudioLoggedSuccess = false
onShow(() => {
  uni.hideTabBar()
  if (!typingInitialized.value) {
    typingInitialized.value = true
    loadTypingData()
  }
  // 检测从收藏页跳转过来的待处理 wordId
  const pendingWordId = uni.getStorageSync('typing-pending-word-id')
  if (pendingWordId) {
    uni.removeStorageSync('typing-pending-word-id')
    setTimeout(() => {
      jumpToWordBySearchResult({ id: Number(pendingWordId) })
    }, 400)
  }
})

onHide(() => {
  // 切换页面时关闭所有弹窗，防止遮罩阻塞其他页面滚动
  pickerPopup.value?.close()
  settingsPopup.value?.close()
  wordDrawerPopup.value?.close()
  isTypingMode.value = false
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
  // 3.3.3 存储离线发音降级文本
  offlineFallbackText = fallbackText || ''
  const isUnreachable = audioSrc && (isLocalAudioSource(audioSrc) || isPrivateNetworkAudioSource(audioSrc))
  const primary = audioSrc && !isUnreachable ? audioSrc : ''

  audioCandidates = []
  if (audioSrc && isUnreachable) {
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
      // 3.3.2 语音播放动画开始
      isAudioPlaying.value = true
      if (currentAudioAttempt && !currentAudioLoggedSuccess) {
        currentAudioLoggedSuccess = true
      }
    })
    audioCtx.onEnded(() => {
      // 3.3.2 语音播放动画结束
      isAudioPlaying.value = false
    })
    audioCtx.onError((err) => {
      // 3.3.2 语音播放动画结束
      isAudioPlaying.value = false
      // 3.3.3 所有在线源失败后尝试离线语音合成
      if (audioCandidates.length === 0 && offlineFallbackText) {
        const spoke = speakOffline(offlineFallbackText)
        if (spoke) {
          offlineFallbackText = ''
          return
        }
      }
      playNextAudioCandidate()
    })
  }
  playNextAudioCandidate()
}

// 3.3.3 离线发音降级方案：Web Speech API
const speakOffline = (text) => {
  // #ifdef H5
  try {
    if (typeof window !== 'undefined' && window.speechSynthesis) {
      window.speechSynthesis.cancel()
      const utterance = new window.SpeechSynthesisUtterance(String(text || '').trim())
      utterance.lang = 'en-US'
      utterance.rate = 0.9
      utterance.onstart = () => { isAudioPlaying.value = true }
      utterance.onend = () => { isAudioPlaying.value = false }
      utterance.onerror = () => { isAudioPlaying.value = false }
      window.speechSynthesis.speak(utterance)
      return true
    }
  } catch (e) {
    // SpeechSynthesis 不可用
  }
  // #endif
  return false
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

const buildEmptyWord = () => ({
  id: 0,
  word: '-',
  phoneticUs: '',
  phoneticUk: '',
  audioUs: '',
  audioUk: '',
  explanation: '',
  sentences: []
})

const updateHeaderNames = () => {
  const category = categories.value.find((item) => item.id === selectedCategoryId.value)
  const chapter = chapters.value.find((item) => item.id === selectedChapterId.value)
  categoryName.value = category ? localText(category.name) : ''
  chapterName.value = chapter ? localText(chapter.name) : ''
}

const syncCurrentWord = () => {
  if (wordList.value.length === 0) {
    currentWord.value = buildEmptyWord()
    return
  }
  if (currentWordIndex.value < 0) {
    currentWordIndex.value = 0
  }
  if (currentWordIndex.value >= wordList.value.length) {
    currentWordIndex.value = wordList.value.length - 1
  }
  const target = wordList.value[currentWordIndex.value]
  currentWord.value = target ? normalizeWord(target) : buildEmptyWord()
  refreshWordCollectedState()
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

// 3.3.5 切换章节后若抽屉已打开则重新聚焦
const reFocusDrawerIfOpen = () => {
  if (!wordDrawerPopup.value) return
  const targetId = currentWord.value.id || currentWordIndex.value
  drawerScrollToId.value = ''
  setTimeout(() => {
    drawerScrollToId.value = 'drawer-word-' + targetId
  }, 200)
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
    isPageLoading.value = false
    return
  }

  const progressCategoryId = Number(progress?.categoryId || 0)
  const progressChapterId = Number(progress?.chapterId || 0)
  const progressWordIndex = Number(progress?.wordIndex || 0)
  const hasProgressCategory = categories.value.some((item) => item.id === progressCategoryId)
  const categoryId = hasProgressCategory ? progressCategoryId : categories.value[0].id

  await selectCategory(categoryId, progressChapterId, progressWordIndex)
  await loadWordCollections(true)
  refreshWordCollectedState()
  isPageLoading.value = false
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
  const activeId = type === 'category' ? selectedCategoryId.value : selectedChapterId.value
  // 原生 DOM scrollIntoView，不依赖 uni-app 框架层（scroll-into-view/scroll-top 在 popup 内均不可靠）
  setTimeout(() => {
    const el = document.getElementById('picker-item-' + activeId)
    if (el) {
      el.scrollIntoView({ block: 'start', behavior: 'smooth' })
    }
  }, 500)
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
    // 3.3.5 切换章节后重聚焦抽屉
    reFocusDrawerIfOpen()
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
/* === uni-popup z-index 覆盖：高于custom-tab-bar的9999 === */
:deep(.uni-popup) {
  z-index: 10001 !important;
}
:deep(.uni-mask) {
  z-index: 10001 !important;
}
:deep(.uni-popup__wrapperbox) {
  z-index: 10002 !important;
}

/* === 页面基底 === */
.typing-page {
  position: relative;
  height: calc(100vh - 96rpx - env(safe-area-inset-bottom));
  display: flex;
  flex-direction: column;
  background: #F5F3FF;
}

/* === 渐变顶栏 === */
.top-bar {
  position: relative;
  overflow: hidden;
}

.top-bar-bg {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 100%;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
}

.top-bar-content {
  position: relative;
  z-index: 2;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: calc(var(--status-bar-height, 0px) + 20rpx) 32rpx 20rpx;
}

.bar-selectors {
  display: flex;
  gap: 16rpx;
  align-items: center;
}

.selector {
  display: flex;
  align-items: center;
  gap: 6rpx;
  padding: 10rpx 24rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.22);
  backdrop-filter: blur(8rpx);
}

.selector:active {
  background: rgba(255, 255, 255, 0.32);
}

.selector-text {
  font-size: 27rpx;
  font-weight: 600;
  color: #fff;
}

.selector-arrow {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  transform: rotate(90deg);
}

.bar-action {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.22);
}

.bar-action:active {
  background: rgba(255, 255, 255, 0.34);
}

.action-icon {
  font-size: 32rpx;
  color: #fff;
}

/* === 单词主卡片 === */
.word-hero-card {
  margin: 20rpx 24rpx 16rpx;
  padding: 40rpx 32rpx 28rpx;
  border-radius: 28rpx;
  background: #FFFFFF;
  border: 1rpx solid rgba(108, 91, 255, 0.12);
  box-shadow: 0 8rpx 32rpx rgba(108, 91, 255, 0.12);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.word-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4rpx;
}

.word-char {
  font-size: 76rpx;
  font-weight: 800;
  font-family: 'SF Mono', 'Menlo', monospace;
  letter-spacing: 2rpx;
}

.char-initial { color: #1A1B3A; }
.char-pending { color: #C7CBE5; }
.char-correct { color: #1A1B3A; text-decoration: underline; text-decoration-color: #6D5BFF; text-underline-offset: 8rpx; }
.char-wrong { color: #EF4444; }

.word-meta {
  margin-top: 20rpx;
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.audio-trigger {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
  transition: transform 0.15s;
}

.audio-icon {
  font-size: 30rpx;
  color: #fff;
}

.audio-trigger:active {
  transform: scale(0.92);
}

.audio-trigger.playing {
  pointer-events: none;
  animation: pulse 0.6s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); opacity: 1; }
  50% { transform: scale(1.12); opacity: 0.75; }
}

.word-phonetic {
  font-size: 28rpx;
  color: #6B6F8D;
}

.word-meaning {
  margin-top: 16rpx;
  font-size: 30rpx;
  color: #1A1B3A;
  text-align: center;
  line-height: 1.5;
}

.word-footer {
  margin-top: 24rpx;
}

.collect-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 28rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.36);
  border-radius: 999rpx;
  background: rgba(108, 91, 255, 0.04);
}

.collect-icon {
  font-size: 26rpx;
  color: #F59E0B;
}

.collect-text {
  font-size: 22rpx;
  color: #6D5BFF;
  font-weight: 600;
}

.collect-btn:active {
  background: rgba(108, 91, 255, 0.12);
}

/* === 工具卡片 === */
.tools-card {
  margin: 0 24rpx 16rpx;
  padding: 20rpx;
  border-radius: 24rpx;
  background: #FFFFFF;
  border: 1rpx solid rgba(108, 91, 255, 0.1);
  box-shadow: 0 4rpx 20rpx rgba(108, 91, 255, 0.08);
  position: relative;
  z-index: 20;
}

.tools-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.tool-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  height: 64rpx;
  padding: 0 24rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.36);
  border-radius: 16rpx;
  font-size: 24rpx;
  font-weight: 600;
  color: #6D5BFF;
  background: rgba(108, 91, 255, 0.04);
}

.tool-icon {
  font-size: 26rpx;
}

.tool-btn:active {
  background: rgba(108, 91, 255, 0.12);
}

.search-field {
  flex: 1;
  height: 64rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.24);
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  overflow: hidden;
  background: #fff;
}

.search-input {
  flex: 1;
  height: 64rpx;
  padding: 0 20rpx;
  font-size: 24rpx;
  color: #1A1B3A;
}

.search-clear {
  width: 44rpx;
  height: 44rpx;
  margin-right: 8rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6B6F8D;
  font-size: 28rpx;
}

.search-go {
  height: 64rpx;
  padding: 0 24rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 24rpx;
  font-weight: 600;
}

.search-go:active {
  opacity: 0.88;
}

.search-results {
  margin-top: 12rpx;
  max-height: 320rpx;
  overflow-y: auto;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  border-radius: 16rpx;
  background: rgba(108, 91, 255, 0.03);
}

.search-empty {
  padding: 24rpx;
  color: #A9AECB;
  font-size: 24rpx;
  text-align: center;
}

.search-item {
  padding: 18rpx 24rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.1);
}

.search-item:last-child { border-bottom: none; }
.search-item:active { background: rgba(108, 91, 255, 0.08); }

.search-item-word {
  font-size: 28rpx;
  font-weight: 700;
  color: #1A1B3A;
  display: block;
}

.search-item-def {
  margin-top: 4rpx;
  font-size: 22rpx;
  color: #6B6F8D;
  display: block;
}

/* === 控制栏 === */
.control-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24rpx 20rpx;
}

.error-badge {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  background: rgba(239, 68, 68, 0.1);
}

.error-count {
  font-size: 32rpx;
  color: #EF4444;
  font-family: monospace;
  font-weight: 800;
}

.error-label {
  font-size: 22rpx;
  color: #6B6F8D;
}

.mode-btn {
  height: 72rpx;
  line-height: 72rpx;
  margin: 0;
  padding: 0 48rpx;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border: none;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
}

.mode-btn:active { opacity: 0.88; transform: scale(0.97); }

.mode-btn.typing {
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  box-shadow: 0 4rpx 16rpx rgba(239, 68, 68, 0.36);
}

/* === 例句区 === */
.sentences {
  flex: 1;
  padding: 0 24rpx;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.sentence-card {
  padding: 24rpx;
  margin-bottom: 16rpx;
  border-radius: 20rpx;
  background: #fff;
  border: 1rpx solid rgba(108, 91, 255, 0.1);
  box-shadow: 0 4rpx 20rpx rgba(108, 91, 255, 0.08);
}

.sentence-card:last-child { margin-bottom: 0; }

.sentence-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.sentence-en {
  font-size: 29rpx;
  font-weight: 600;
  color: #1A1B3A;
  line-height: 1.45;
  flex: 1;
}

.sentence-audio {
  margin-left: 16rpx;
  width: 52rpx;
  height: 52rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(108, 91, 255, 0.12);
  flex-shrink: 0;
}

.sentence-audio-icon {
  font-size: 26rpx;
  color: #6D5BFF;
}

.sentence-audio:active {
  background: rgba(108, 91, 255, 0.22);
}

.sentence-zh {
  margin-top: 10rpx;
  font-size: 25rpx;
  color: #6B6F8D;
  line-height: 1.5;
}

/* === 悬浮切换按钮 === */
.float-nav {
  position: fixed;
  bottom: 340rpx;
  padding: 16rpx 28rpx;
  background: linear-gradient(135deg, rgba(108, 91, 255, 0.92) 0%, rgba(155, 143, 255, 0.92) 100%);
  color: #fff;
  font-size: 22rpx;
  font-weight: 600;
  border-radius: 999rpx;
  z-index: 100;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
}

.float-nav:active { opacity: 0.85; transform: scale(0.95); }
.float-nav.prev { left: 20rpx; }
.float-nav.next { right: 20rpx; }

/* === 键盘 === */
.keyboard {
  position: fixed;
  bottom: -450rpx;
  left: 0;
  width: 100%;
  height: 450rpx;
  background: #EDE9FE;
  transition: bottom 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  padding: 16rpx 12rpx calc(20rpx + env(safe-area-inset-bottom));
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  z-index: 10000;
  border-top: 1rpx solid rgba(108, 91, 255, 0.16);
  border-radius: 28rpx 28rpx 0 0;
}

.keyboard.show { bottom: calc(88rpx + env(safe-area-inset-bottom)); }

.kb-row {
  display: flex;
  justify-content: center;
  gap: 10rpx;
}

.kb-key {
  flex: 1;
  height: 88rpx;
  min-width: 58rpx;
  background: #fff;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 38rpx;
  font-family: monospace;
  font-weight: 600;
  color: #1A1B3A;
  box-shadow: 0 2rpx 8rpx rgba(108, 91, 255, 0.14);
  text-transform: uppercase;
}

.kb-key-hover {
  background: rgba(108, 91, 255, 0.16);
  transform: translateY(1rpx);
  color: #6D5BFF;
}

/* === 底部弹窗 === */
.sheet {
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  box-shadow: 0 -12rpx 36rpx rgba(108, 91, 255, 0.14);
  max-height: 75vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sheet-handle {
  width: 60rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #6D5BFF, #9B8FFF);
  margin: 16rpx auto 0;
  flex-shrink: 0;
}

.sheet-title {
  text-align: center;
  font-size: 32rpx;
  font-weight: 800;
  color: #1A1B3A;
  margin: 20rpx 0;
  padding: 0 32rpx;
  flex-shrink: 0;
}

.sheet-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 0 32rpx;
  box-sizing: border-box;
  width: 100%;
}

.sheet-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 92rpx;
  padding: 0 24rpx;
  border-radius: 16rpx;
  margin-bottom: 8rpx;
  font-size: 28rpx;
  color: #1A1B3A;
  background: rgba(108, 91, 255, 0.06);
  box-sizing: border-box;
  width: 100%;
}

.sheet-item:active { background: rgba(108, 91, 255, 0.14); }

.sheet-item.active {
  background: rgba(108, 91, 255, 0.16);
  border: 1rpx solid rgba(108, 91, 255, 0.32);
}

.sheet-item-label {
  font-weight: 600;
}

.sheet-item.active .sheet-item-label {
  color: #6D5BFF;
  font-weight: 700;
}

.sheet-item-check {
  color: #6D5BFF;
  font-size: 30rpx;
  font-weight: 700;
}

.sheet-cancel {
  flex-shrink: 0;
  margin: 16rpx 32rpx;
  margin-bottom: calc(16rpx + 88rpx + env(safe-area-inset-bottom));
  border-radius: 999rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.32);
  background: #fff;
  color: #6D5BFF;
  font-size: 30rpx;
  font-weight: 600;
  height: 84rpx;
  line-height: 84rpx;
}

.set-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 32rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.1);
  font-size: 30rpx;
  color: #1A1B3A;
  flex-shrink: 0;
  box-sizing: border-box;
  width: 100%;
}

.set-row:last-of-type {
  border-bottom: none;
}

.accent-label {
  display: inline-flex;
  align-items: center;
  gap: 6rpx;
  margin-right: 24rpx;
  font-size: 28rpx;
  color: #1A1B3A;
}

/* === 左侧抽屉 === */
.drawer {
  width: 560rpx;
  height: 100vh;
  background: #F5F3FF;
  display: flex;
  flex-direction: column;
  box-shadow: 8rpx 0 32rpx rgba(108, 91, 255, 0.18);
}

.drawer-head {
  height: calc(96rpx + env(safe-area-inset-top));
  padding: env(safe-area-inset-top) 24rpx 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.12);
  background: #fff;
  box-shadow: 0 2rpx 12rpx rgba(108, 91, 255, 0.06);
  box-sizing: border-box;
  width: 100%;
}

.drawer-title {
  font-size: 32rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.drawer-close {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
  color: #6D5BFF;
}

.drawer-body {
  flex: 1;
  padding: 16rpx 20rpx calc(40rpx + env(safe-area-inset-bottom));
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  box-sizing: border-box;
  width: 100%;
}

.drawer-item {
  padding: 20rpx 24rpx;
  border-radius: 16rpx;
  margin-bottom: 10rpx;
  background: #fff;
  border: 1rpx solid rgba(108, 91, 255, 0.08);
  box-sizing: border-box;
  width: 100%;
}

.drawer-item:active { background: rgba(108, 91, 255, 0.08); }

.drawer-item.active {
  background: rgba(108, 91, 255, 0.14);
  border-color: rgba(108, 91, 255, 0.32);
}

.drawer-word {
  font-size: 28rpx;
  font-weight: 700;
  color: #1A1B3A;
  display: block;
}

.drawer-item.active .drawer-word {
  color: #6D5BFF;
}

.drawer-def {
  margin-top: 6rpx;
  font-size: 22rpx;
  color: #6B6F8D;
  display: block;
}

/* === 骨架屏 === */
.skeleton-mode { pointer-events: none; }

.sk-line {
  border-radius: 8rpx;
  background: linear-gradient(90deg, #EDE9FE 25%, #F5F3FF 50%, #EDE9FE 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.sk-w1 { width: 300rpx; height: 72rpx; margin-bottom: 16rpx; }
.sk-w2 { width: 180rpx; height: 28rpx; margin-bottom: 12rpx; }
.sk-w3 { width: 400rpx; height: 24rpx; }
.sk-btn { width: 160rpx; height: 64rpx; border-radius: 16rpx; }
.sk-input { flex: 1; height: 64rpx; border-radius: 16rpx; }

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* === 抖动动画 === */
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-12rpx); }
  50% { transform: translateX(12rpx); }
  75% { transform: translateX(-12rpx); }
}

.shake-animation { animation: shake 0.35s ease-in-out; }
</style>
