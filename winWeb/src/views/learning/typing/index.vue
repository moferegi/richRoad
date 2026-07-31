<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  getCategoryList,
  getChapterList,
  getWordList,
  findWord,
  reportWordError,
  getWordProgress,
  saveWordProgress,
  getCollectionList,
  collect,
  uncollect
} from '@/api/learning'
import { localText } from '@/utils/i18n'
import { playWordAudio as playWordTTS, stopLocalTTS, destroyTTS } from '@/utils/learning-tts'

const { t } = useI18n()

// 分类/章节
const categories = ref([])
const chapters = ref([])
const currentCategoryId = ref(null)
const currentChapterId = ref(null)
const showCategoryPicker = ref(false)
const showChapterPicker = ref(false)
const showSettings = ref(false)

// 单词列表
const wordList = ref([])
const currentWordIndex = ref(0)
const loading = ref(false)

// 设置
const settings = reactive({
  showWord: true,
  showPhonetic: true,
  showExplanation: true,
  accent: 'US',
  ttsVoice: 'female'
})

// 跟打模式
const isTypingMode = ref(false)
const typedChars = ref('')
const errorCount = ref(0)
const errorChars = ref({}) // { index: true } 标记错误位置

// 收藏
const collectedMap = ref({})

// 搜索
const searchQuery = ref('')
const showWordList = ref(true)

// 当前单词
const currentWord = computed(() => wordList.value[currentWordIndex.value] || null)

// 是否收藏当前单词
const isCollected = computed(() => {
  if (!currentWord.value) return false
  const id = currentWord.value.ID || currentWord.value.id
  return !!collectedMap.value[id]
})

// 加载分类
async function loadCategories() {
  try {
    const res = await getCategoryList({ page: 1, pageSize: 200 })
    if (res.code === 0) {
      categories.value = res.data?.list || res.data || []
      if (categories.value.length > 0 && !currentCategoryId.value) {
        currentCategoryId.value = categories.value[0].ID || categories.value[0].id
      }
    }
  } catch (e) {
    console.warn('load categories failed', e)
  }
}

// 加载章节
async function loadChapters() {
  if (!currentCategoryId.value) return
  try {
    const res = await getChapterList({ categoryId: currentCategoryId.value, page: 1, pageSize: 200 })
    if (res.code === 0) {
      chapters.value = res.data?.list || res.data || []
      if (chapters.value.length > 0) {
        currentChapterId.value = chapters.value[0].ID || chapters.value[0].id
      }
    }
  } catch (e) {
    console.warn('load chapters failed', e)
  }
}

// 加载单词
async function loadWords() {
  if (!currentChapterId.value) return
  loading.value = true
  try {
    const res = await getWordList({
      chapterId: currentChapterId.value,
      categoryId: currentCategoryId.value,
      page: 1,
      pageSize: 500
    })
    if (res.code === 0) {
      wordList.value = res.data?.list || res.data || []
      currentWordIndex.value = 0
      resetTyping()
    }
  } catch (e) {
    console.warn('load words failed', e)
  } finally {
    loading.value = false
  }
}

// 加载学习进度
async function loadProgress() {
  try {
    const res = await getWordProgress()
    if (res.code === 0 && res.data) {
      const data = res.data
      if (data.chapterId && data.wordIndex >= 0) {
        // 如果当前章节匹配，直接定位
        if (data.chapterId === currentChapterId.value) {
          const idx = Number(data.wordIndex) || 0
          if (idx >= 0 && idx < wordList.value.length) {
            currentWordIndex.value = idx
          }
        }
      }
    }
  } catch (e) {
    console.warn('load progress failed', e)
  }
}

// 加载收藏状态
async function loadCollections() {
  try {
    const res = await getCollectionList({ targetType: 1, pageSize: 500 })
    if (res.code === 0) {
      const list = res.data?.list || res.data || []
      const map = {}
      list.forEach(item => {
        map[item.targetId] = true
      })
      collectedMap.value = map
    }
  } catch (e) {
    console.warn('load collections failed', e)
  }
}

// 切换分类
function selectCategory(cat) {
  currentCategoryId.value = cat.ID || cat.id
  showCategoryPicker.value = false
  loadChapters()
}

// 切换章节
function selectChapter(ch) {
  currentChapterId.value = ch.ID || ch.id
  showChapterPicker.value = false
  loadWords()
}

// 上一个单词
function prevWord() {
  if (currentWordIndex.value > 0) {
    currentWordIndex.value--
    resetTyping()
    playAudio()
    saveCurrentProgress()
  } else {
    // 跳转到上一章最后一个单词
    const idx = chapters.value.findIndex(
      c => (c.ID || c.id) === currentChapterId.value
    )
    if (idx > 0) {
      selectChapter(chapters.value[idx - 1])
    }
  }
}

// 下一个单词
function nextWord() {
  if (currentWordIndex.value < wordList.value.length - 1) {
    currentWordIndex.value++
    resetTyping()
    playAudio()
    saveCurrentProgress()
  } else {
    // 跳转到下一章第一个单词
    const idx = chapters.value.findIndex(
      c => (c.ID || c.id) === currentChapterId.value
    )
    if (idx < chapters.value.length - 1) {
      selectChapter(chapters.value[idx + 1])
    }
  }
}

// 跳转到指定单词
function jumpToWord(index) {
  currentWordIndex.value = index
  resetTyping()
  playAudio()
  saveCurrentProgress()
}

// 播放发音
function playAudio() {
  const w = currentWord.value
  if (!w) return
  playWordTTS(w.word, w.audioUs || '', w.audioUk || '', settings.accent, {
    ttsVoice: settings.ttsVoice,
    onStart: () => {},
    onEnd: () => {},
    onError: () => {}
  })
}

// 播放例句发音
function playSentenceAudio(s) {
  const text = s?.en || s?.english || ''
  if (!text) return
  playWordTTS(text, '', '', settings.accent, {
    ttsVoice: settings.ttsVoice,
    onStart: () => {},
    onEnd: () => {},
    onError: () => {}
  })
}

// 收藏/取消收藏
async function toggleCollect() {
  const w = currentWord.value
  if (!w) return
  const id = w.ID || w.id
  try {
    if (isCollected.value) {
      await uncollect(1, id)
      delete collectedMap.value[id]
    } else {
      await collect(1, id)
      collectedMap.value[id] = true
    }
  } catch (e) {
    console.warn('toggle collect failed', e)
  }
}

// 保存进度
function saveCurrentProgress() {
  const w = currentWord.value
  if (!w) return
  saveWordProgress({
    categoryId: currentCategoryId.value,
    chapterId: currentChapterId.value,
    wordIndex: currentWordIndex.value
  }).catch(() => {})
}

// ========== 跟打相关 ==========
function toggleTypingMode() {
  isTypingMode.value = !isTypingMode.value
  if (isTypingMode.value) {
    resetTyping()
    nextTick(() => {
      document.getElementById('typing-input')?.focus()
    })
  }
}

function resetTyping() {
  typedChars.value = ''
  errorChars.value = {}
}

// 键盘输入
function handleKeyDown(e) {
  if (!isTypingMode.value || !currentWord.value) return
  const word = currentWord.value.word || ''

  // 忽略功能键
  if (e.key === 'Backspace') {
    e.preventDefault()
    if (typedChars.value.length > 0) {
      const idx = typedChars.value.length - 1
      typedChars.value = typedChars.value.slice(0, -1)
      delete errorChars.value[idx]
    }
    return
  }

  if (e.key.length !== 1) return
  if (e.ctrlKey || e.metaKey || e.altKey) return

  e.preventDefault()
  const char = e.key.toLowerCase()
  const expectedChar = word[typedChars.value.length]?.toLowerCase()

  if (!expectedChar) return

  if (char === expectedChar) {
    typedChars.value += char
    // 完成整个单词
    if (typedChars.value.length === word.length) {
      setTimeout(() => {
        nextWord()
      }, 300)
    }
  } else {
    // 错误
    errorChars.value[typedChars.value.length] = true
    errorCount.value++
    // 上报错误
    const w = currentWord.value
    if (w) {
      reportWordError({
        wordId: w.ID || w.id,
        categoryId: currentCategoryId.value,
        chapterId: currentChapterId.value,
        wrongIndex: typedChars.value.length,
        expectedChar: expectedChar,
        inputChar: char
      }).catch(() => {})
    }
  }
}

// 搜索单词
async function handleSearch() {
  const keyword = searchQuery.value.trim()
  if (!keyword) return
  try {
    const res = await getWordList({ keyword, page: 1, pageSize: 20 })
    if (res.code === 0 && res.data) {
      const list = res.data?.list || res.data || []
      if (list.length > 0) {
        const word = list[0]
        const idx = wordList.value.findIndex(
          w => (w.ID || w.id) === (word.ID || word.id)
        )
        if (idx >= 0) {
          jumpToWord(idx)
        } else {
          wordList.value.unshift(word)
          jumpToWord(0)
        }
      }
    }
  } catch (e) {
    console.warn('search word failed', e)
  }
}

// 当前单词的显示字符数组
const wordChars = computed(() => {
  const word = currentWord.value?.word || ''
  return word.split('').map((char, i) => ({
    char,
    typed: i < typedChars.value.length,
    error: errorChars.value[i],
    current: i === typedChars.value.length
  }))
})

// 过滤后的词表
const filteredWordList = computed(() => {
  if (!searchQuery.value.trim()) return wordList.value
  const q = searchQuery.value.toLowerCase()
  return wordList.value.filter(w =>
    (w.word || '').toLowerCase().includes(q)
  )
})

// 分类名称
const currentCategoryName = computed(() => {
  const cat = categories.value.find(c => (c.ID || c.id) === currentCategoryId.value)
  return cat ? localText(cat.name) : t('typing.category')
})

// 章节名称
const currentChapterName = computed(() => {
  const ch = chapters.value.find(c => (c.ID || c.id) === currentChapterId.value)
  return ch ? localText(ch.name) : t('typing.chapter')
})

onMounted(async () => {
  await loadCategories()
  await loadChapters()
  await loadWords()
  loadProgress()
  loadCollections()

  // 监听键盘
  window.addEventListener('keydown', handleKeyDown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown)
  saveCurrentProgress()
  destroyTTS()
})
</script>

<template>
  <div class="typing-page">
    <!-- 顶部栏 -->
    <div class="typing-header">
      <div class="header-left">
        <button class="picker-btn" @click="showCategoryPicker = true">
          <span class="picker-label">{{ t('typing.category') }}</span>
          <span class="picker-value">{{ currentCategoryName }}</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
        <span class="header-sep">/</span>
        <button class="picker-btn" @click="showChapterPicker = true">
          <span class="picker-label">{{ t('typing.chapter') }}</span>
          <span class="picker-value">{{ currentChapterName }}</span>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6 9 12 15 18 9"></polyline>
          </svg>
        </button>
      </div>

      <div class="header-center">
        <div class="error-badge" :title="t('typing.errorCount')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="15" y1="9" x2="9" y2="15"></line>
            <line x1="9" y1="9" x2="15" y2="15"></line>
          </svg>
          {{ errorCount }}
        </div>
      </div>

      <div class="header-right">
        <div class="search-box">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t('typing.search')"
            @keyup.enter="handleSearch"
          />
          <button @click="handleSearch">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
          </button>
        </div>
        <button class="icon-btn" @click="showSettings = true" :title="t('typing.settings')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"></circle>
            <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"></path>
          </svg>
        </button>
      </div>
    </div>

    <div class="typing-body">
      <!-- 左侧词表 -->
      <div class="wordlist-panel">
        <div class="panel-header">
          <h3>{{ t('typing.wordList') }}</h3>
          <span class="word-count">{{ wordList.length }}</span>
        </div>
        <div class="wordlist-scroll">
          <div
            v-for="(w, i) in filteredWordList"
            :key="w.ID || w.id || i"
            class="wordlist-item"
            :class="{ active: currentWordIndex === i, collected: collectedMap[w.ID || w.id] }"
            @click="jumpToWord(i)"
          >
            <span class="item-index">{{ i + 1 }}</span>
            <span class="item-word">{{ w.word }}</span>
            <svg v-if="collectedMap[w.ID || w.id]" class="item-star" viewBox="0 0 24 24" fill="currentColor">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
            </svg>
          </div>
        </div>
      </div>

      <!-- 中间主卡片 -->
      <div class="main-panel">
        <div v-if="loading" class="loading-wrap">
          <div class="loading-spinner"></div>
          <span>{{ t('common.loading') }}</span>
        </div>

        <div v-else-if="currentWord" class="word-card-wrap">
          <div class="word-card card">
            <!-- 收藏按钮 -->
            <button class="collect-btn" :class="{ active: isCollected }" @click="toggleCollect">
              <svg viewBox="0 0 24 24" :fill="isCollected ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
                <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
              </svg>
            </button>

            <!-- 单词 -->
            <div class="word-display" @click="playAudio">
              <div v-if="settings.showWord" class="word-text">
                <span
                  v-for="(c, i) in wordChars"
                  :key="i"
                  class="word-char"
                  :class="{ typed: c.typed, error: c.error, current: c.current && isTypingMode }"
                >{{ c.char }}</span>
              </div>
              <div v-else class="word-hidden">***</div>

              <button class="play-btn" @click.stop="playAudio">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <polygon points="5 3 19 12 5 21 5 3"></polygon>
                </svg>
              </button>
            </div>

            <!-- 音标 -->
            <div v-if="settings.showPhonetic" class="phonetic-row">
              <span class="phonetic">US: {{ currentWord.phoneticUs }}</span>
              <span class="phonetic">UK: {{ currentWord.phoneticUk }}</span>
            </div>

            <!-- 释义 -->
            <div v-if="settings.showExplanation" class="explanation">
              {{ currentWord.explanation || currentWord.meaning }}
            </div>

            <!-- 进度 -->
            <div class="word-progress">
              <div class="progress-info">
                <span>{{ currentWordIndex + 1 }} / {{ wordList.length }}</span>
                <span>{{ ((currentWordIndex + 1) / wordList.length * 100).toFixed(1) }}%</span>
              </div>
              <div class="progress-bar">
                <div
                  class="progress-fill"
                  :style="{ width: ((currentWordIndex + 1) / wordList.length * 100) + '%' }"
                ></div>
              </div>
            </div>
          </div>

          <!-- 控制按钮 -->
          <div class="control-row">
            <button class="ctrl-btn prev-btn" @click="prevWord">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15 18 9 12 15 6"></polyline>
              </svg>
              {{ t('typing.prevWord') }}
            </button>

            <button
              class="typing-toggle-btn"
              :class="{ active: isTypingMode }"
              @click="toggleTypingMode"
            >
              {{ isTypingMode ? t('typing.stopTyping') : t('typing.startTyping') }}
            </button>

            <button class="ctrl-btn next-btn" @click="nextWord">
              {{ t('typing.nextWord') }}
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9 18 15 12 9 6"></polyline>
              </svg>
            </button>
          </div>

          <!-- 隐藏输入框用于捕获物理键盘 -->
          <input
            id="typing-input"
            type="text"
            style="position: absolute; opacity: 0; pointer-events: none;"
            @blur="isTypingMode && document.getElementById('typing-input')?.focus()"
          />

          <!-- 例句区 -->
          <div v-if="currentWord.sentences?.length > 0" class="sentences-section">
            <h4 class="section-title">例句</h4>
            <div
              v-for="(s, i) in currentWord.sentences"
              :key="i"
              class="sentence-card"
            >
              <div class="sentence-row-header">
                <p class="sentence-en">{{ s.english || s.en }}</p>
                <button class="sentence-audio-btn" @click="playSentenceAudio(s)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                    <polygon points="5 3 19 12 5 21 5 3"></polygon>
                  </svg>
                </button>
              </div>
              <p v-if="s.translate || s.zh" class="sentence-zh">{{ s.translate || s.zh }}</p>
            </div>
          </div>
        </div>

        <div v-else class="empty-state">
          <p>{{ t('common.empty') }}</p>
        </div>
      </div>
    </div>

    <!-- 分类选择弹窗 -->
    <div v-if="showCategoryPicker" class="modal-mask" @click.self="showCategoryPicker = false">
      <div class="modal-box">
        <div class="modal-header">
          <h3>{{ t('typing.category') }}</h3>
          <button class="modal-close" @click="showCategoryPicker = false">×</button>
        </div>
        <div class="modal-body">
          <div
            v-for="cat in categories"
            :key="cat.ID || cat.id"
            class="list-item"
            :class="{ active: currentCategoryId === (cat.ID || cat.id) }"
            @click="selectCategory(cat)"
          >
            {{ cat.nameI18n || cat.name }}
          </div>
        </div>
      </div>
    </div>

    <!-- 章节选择弹窗 -->
    <div v-if="showChapterPicker" class="modal-mask" @click.self="showChapterPicker = false">
      <div class="modal-box">
        <div class="modal-header">
          <h3>{{ t('typing.chapter') }}</h3>
          <button class="modal-close" @click="showChapterPicker = false">×</button>
        </div>
        <div class="modal-body">
          <div
            v-for="ch in chapters"
            :key="ch.ID || ch.id"
            class="list-item"
            :class="{ active: currentChapterId === (ch.ID || ch.id) }"
            @click="selectChapter(ch)"
          >
            {{ ch.nameI18n || ch.name }}
          </div>
        </div>
      </div>
    </div>

    <!-- 设置弹窗 -->
    <div v-if="showSettings" class="modal-mask" @click.self="showSettings = false">
      <div class="modal-box settings-modal">
        <div class="modal-header">
          <h3>{{ t('typing.settings') }}</h3>
          <button class="modal-close" @click="showSettings = false">×</button>
        </div>
        <div class="modal-body">
          <div class="setting-row">
            <label>{{ t('typing.showWord') }}</label>
            <div class="toggle-switch" :class="{ on: settings.showWord }" @click="settings.showWord = !settings.showWord">
              <div class="toggle-dot"></div>
            </div>
          </div>
          <div class="setting-row">
            <label>{{ t('typing.showPhonetic') }}</label>
            <div class="toggle-switch" :class="{ on: settings.showPhonetic }" @click="settings.showPhonetic = !settings.showPhonetic">
              <div class="toggle-dot"></div>
            </div>
          </div>
          <div class="setting-row">
            <label>{{ t('typing.showExplanation') }}</label>
            <div class="toggle-switch" :class="{ on: settings.showExplanation }" @click="settings.showExplanation = !settings.showExplanation">
              <div class="toggle-dot"></div>
            </div>
          </div>
          <div class="setting-row">
            <label>{{ t('typing.accent') }}</label>
            <div class="accent-selector">
              <button :class="{ active: settings.accent === 'US' }" @click="settings.accent = 'US'">{{ t('typing.usAccent') }}</button>
              <button :class="{ active: settings.accent === 'UK' }" @click="settings.accent = 'UK'">{{ t('typing.ukAccent') }}</button>
            </div>
          </div>
          <div class="setting-row">
            <label>{{ t('typing.ttsVoice') }}</label>
            <div class="accent-selector">
              <button :class="{ active: settings.ttsVoice === 'female' }" @click="settings.ttsVoice = 'female'">{{ t('typing.femaleVoice') }}</button>
              <button :class="{ active: settings.ttsVoice === 'male' }" @click="settings.ttsVoice = 'male'">{{ t('typing.maleVoice') }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.typing-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - #{$header-height} - 48px);
  max-width: $content-max-width;
  margin: 0 auto;
}

// ========== 顶部栏 ==========
.typing-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md 0;
  margin-bottom: $spacing-md;

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .header-sep {
    color: $text-tertiary;
  }

  .picker-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 14px;
    background: $bg-card;
    border-radius: $radius-md;
    transition: all 0.2s ease;

    &:hover {
      background: $bg-hover;
    }

    .picker-label {
      font-size: 12px;
      color: $text-tertiary;
    }

    .picker-value {
      font-size: 14px;
      font-weight: $font-weight-medium;
      color: $text-primary;
      max-width: 120px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }

    svg {
      width: 14px;
      height: 14px;
      color: $text-tertiary;
    }
  }

  .error-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 14px;
    background: rgba(255, 59, 48, 0.1);
    color: $error-color;
    border-radius: $radius-pill;
    font-size: 13px;
    font-weight: $font-weight-medium;

    svg {
      width: 16px;
      height: 16px;
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .search-box {
    display: flex;
    align-items: center;
    background: $bg-card;
    border-radius: $radius-pill;
    padding: 4px 4px 4px 14px;
    border: 1px solid $border-light;
    transition: all 0.2s ease;

    &:focus-within {
      border-color: $primary-color;
      box-shadow: 0 0 0 3px rgba(109, 91, 255, 0.1);
    }

    input {
      border: none;
      background: transparent;
      font-size: 13px;
      padding: 6px 0;
      width: 160px;
    }

    button {
      width: 28px;
      height: 28px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: $gradient-primary;
      color: $text-white;
      border-radius: 50%;

      svg {
        width: 14px;
        height: 14px;
      }
    }
  }

  .icon-btn {
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: $bg-card;
    border-radius: 50%;
    color: $text-secondary;
    transition: all 0.2s ease;

    &:hover {
      color: $primary-color;
      background: $bg-active;
    }

    svg {
      width: 20px;
      height: 20px;
    }
  }
}

// ========== 主体 ==========
.typing-body {
  flex: 1;
  display: flex;
  gap: $spacing-lg;
  min-height: 0;
}

// ========== 词表面板 ==========
.wordlist-panel {
  width: 240px;
  flex-shrink: 0;
  background: $bg-card;
  border-radius: $radius-lg;
  display: flex;
  flex-direction: column;
  box-shadow: $shadow-card;
  border: 1px solid $border-light;
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid $border-light;

  h3 {
    font-size: 15px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
  }

  .word-count {
    font-size: 12px;
    color: $text-tertiary;
    background: $bg-input;
    padding: 2px 10px;
    border-radius: $radius-pill;
  }
}

.wordlist-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.wordlist-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: $radius-sm;
  cursor: pointer;
  transition: all 0.15s ease;

  &:hover {
    background: $bg-hover;
  }

  &.active {
    background: $bg-active;

    .item-word {
      color: $primary-color;
      font-weight: $font-weight-semibold;
    }

    .item-index {
      background: $gradient-primary;
      color: $text-white;
    }
  }

  .item-index {
    width: 22px;
    height: 22px;
    background: $bg-input;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 11px;
    color: $text-tertiary;
    font-weight: $font-weight-medium;
    flex-shrink: 0;
  }

  .item-word {
    flex: 1;
    font-size: 14px;
    color: $text-primary;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .item-star {
    width: 14px;
    height: 14px;
    color: $warning-color;
    flex-shrink: 0;
  }
}

// ========== 主面板 ==========
.main-panel {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
}

.loading-wrap {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  height: 200px;
  color: $text-secondary;

  .loading-spinner {
    width: 20px;
    height: 20px;
    border: 2px solid $border-color;
    border-top-color: $primary-color;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.word-card-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-lg;
}

.word-card {
  position: relative;
  width: 100%;
  max-width: 680px;
  padding: $spacing-xxl $spacing-xl;
  text-align: center;

  .collect-btn {
    position: absolute;
    top: 20px;
    right: 20px;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    color: $text-tertiary;
    transition: all 0.2s ease;

    &:hover,
    &.active {
      color: $warning-color;
      background: rgba(255, 149, 0, 0.1);
    }

    svg {
      width: 22px;
      height: 22px;
    }
  }

  .word-display {
    position: relative;
    margin-bottom: 16px;
    cursor: pointer;

    .word-text {
      font-size: 64px;
      font-weight: $font-weight-bold;
      color: $text-primary;
      letter-spacing: 2px;
      line-height: 1.2;
    }

    .word-hidden {
      font-size: 64px;
      font-weight: $font-weight-bold;
      color: $border-primary;
      letter-spacing: 8px;
    }

    .word-char {
      display: inline-block;
      transition: color 0.15s ease;

      &.typed {
        color: $primary-color;
      }

      &.error {
        color: $error-color;
        text-decoration: line-through;
      }

      &.current {
        position: relative;

        &::after {
          content: '';
          position: absolute;
          bottom: -4px;
          left: 0;
          right: 0;
          height: 3px;
          background: $primary-color;
          border-radius: 2px;
          animation: blink 1s infinite;
        }
      }
    }

    .play-btn {
      position: absolute;
      right: -60px;
      top: 50%;
      transform: translateY(-50%);
      width: 48px;
      height: 48px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: $gradient-primary;
      color: $text-white;
      border-radius: 50%;
      box-shadow: 0 4px 16px rgba(109, 91, 255, 0.3);
      transition: all 0.2s ease;

      &:hover {
        transform: translateY(-50%) scale(1.08);
      }

      svg {
        width: 20px;
        height: 20px;
        margin-left: 2px;
      }
    }
  }

  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }

  .phonetic-row {
    display: flex;
    justify-content: center;
    gap: 24px;
    margin-bottom: 16px;

    .phonetic {
      font-size: 15px;
      color: $text-secondary;
      font-style: italic;
    }
  }

  .explanation {
    font-size: 16px;
    color: $text-primary;
    line-height: 1.6;
    max-width: 480px;
    margin: 0 auto $spacing-xl;
  }

  .word-progress {
    max-width: 360px;
    margin: 0 auto;

    .progress-info {
      display: flex;
      justify-content: space-between;
      font-size: 12px;
      color: $text-tertiary;
      margin-bottom: 6px;
    }

    .progress-bar {
      height: 6px;
      background: $bg-input;
      border-radius: 3px;
      overflow: hidden;

      .progress-fill {
        height: 100%;
        background: $gradient-primary;
        border-radius: 3px;
        transition: width 0.3s ease;
      }
    }
  }
}

// 控制按钮
.control-row {
  display: flex;
  align-items: center;
  gap: $spacing-lg;

  .ctrl-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 12px 24px;
    background: $bg-card;
    border-radius: $radius-pill;
    color: $text-secondary;
    font-size: 14px;
    border: 1px solid $border-color;
    transition: all 0.2s ease;

    &:hover {
      border-color: $primary-color;
      color: $primary-color;
    }

    svg {
      width: 18px;
      height: 18px;
    }
  }

  .typing-toggle-btn {
    padding: 14px 36px;
    background: $gradient-primary;
    color: $text-white;
    border-radius: $radius-pill;
    font-size: 15px;
    font-weight: $font-weight-semibold;
    box-shadow: 0 4px 16px rgba(109, 91, 255, 0.3);
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(109, 91, 255, 0.4);
    }

    &.active {
      background: $text-primary;
      box-shadow: 0 4px 16px rgba(26, 27, 58, 0.3);
    }
  }
}

// 例句
.sentences-section {
  width: 100%;
  max-width: 680px;

  .section-title {
    font-size: 16px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
    margin-bottom: 12px;
  }
}

.sentence-card {
  padding: 16px 20px;
  background: $bg-card;
  border-radius: $radius-md;
  margin-bottom: 12px;
  border-left: 3px solid $primary-color;
  box-shadow: $shadow-card;

  .sentence-row-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 6px;
  }

  .sentence-en {
    font-size: 15px;
    color: $text-primary;
    line-height: 1.6;
    flex: 1;
  }

  .sentence-audio-btn {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid $border-primary;
    background: $bg-page;
    color: $primary-color;
    cursor: pointer;
    flex-shrink: 0;
    transition: all 0.2s;

    &:hover {
      background: $gradient-primary;
      color: #fff;
      border-color: transparent;
    }
  }

  .sentence-zh {
    font-size: 13px;
    color: $text-secondary;
    line-height: 1.5;
  }
}

// ========== 弹窗 ==========
.modal-mask {
  position: fixed;
  inset: 0;
  background: $bg-mask;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: $z-modal;
}

.modal-box {
  width: 360px;
  max-height: 70vh;
  background: $bg-card;
  border-radius: $radius-xl;
  display: flex;
  flex-direction: column;
  box-shadow: $shadow-xl;
  animation: modalIn 0.25s ease;

  &.settings-modal {
    width: 400px;
  }
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid $border-light;

  h3 {
    font-size: 17px;
    font-weight: $font-weight-semibold;
    color: $text-primary;
  }

  .modal-close {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 22px;
    color: $text-tertiary;
    border-radius: 50%;

    &:hover {
      background: $bg-hover;
      color: $text-primary;
    }
  }
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px 0;
}

.list-item {
  padding: 12px 24px;
  font-size: 14px;
  color: $text-primary;
  cursor: pointer;
  transition: background 0.15s ease;

  &:hover {
    background: $bg-hover;
  }

  &.active {
    color: $primary-color;
    font-weight: $font-weight-medium;
    background: $bg-active;
  }
}

// 设置项
.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 24px;

  label {
    font-size: 14px;
    color: $text-primary;
  }
}

.toggle-switch {
  position: relative;
  width: 44px;
  height: 24px;
  background: $border-color;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s ease;

  &.on {
    background: $primary-color;

    .toggle-dot {
      left: 22px;
    }
  }

  .toggle-dot {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 20px;
    height: 20px;
    background: $text-white;
    border-radius: 50%;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    transition: left 0.2s ease;
  }
}

.accent-selector {
  display: flex;
  background: $bg-input;
  border-radius: $radius-pill;
  padding: 2px;

  button {
    padding: 6px 16px;
    font-size: 13px;
    color: $text-secondary;
    border-radius: $radius-pill;
    transition: all 0.2s ease;

    &.active {
      background: $bg-card;
      color: $primary-color;
      font-weight: $font-weight-medium;
      box-shadow: $shadow-sm;
    }
  }
}

.empty-state {
  text-align: center;
  padding: 60px 0;
  color: $text-tertiary;
  font-size: 14px;
}
</style>
