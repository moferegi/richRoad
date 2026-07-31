<template>
  <div class="diary-detail-page">
    <!-- 顶部信息栏 -->
    <div class="detail-header">
      <button class="back-btn" @click="goBack">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15 18 9 12 15 6"></polyline>
        </svg>
        {{ t('common.back') }}
      </button>

      <h1 class="diary-title">{{ localText(diaryInfo?.name) }}</h1>

      <div class="header-actions">
        <span class="duration-info">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <polyline points="12 6 12 12 16 14"></polyline>
          </svg>
          {{ formatTime(currentTime) }}
        </span>

        <button class="accent-switch-btn" @click="toggleAccent" :disabled="isAccentSwitching">
          <template v-if="isAccentSwitching">{{ t('diaryDetail.accent') }}...</template>
          <template v-else>
            <span :class="{ active: accent === 'US' }">US</span>
            <span class="divider">/</span>
            <span :class="{ active: accent === 'UK' }">UK</span>
          </template>
        </button>

        <button class="image-btn" @click="viewImage" :title="t('diaryDetail.bilingual')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
            <circle cx="8.5" cy="8.5" r="1.5"></circle>
            <polyline points="21 15 16 10 5 21"></polyline>
          </svg>
        </button>

        <button
          class="action-btn"
          :class="{ active: isCollected }"
          @click="collectDiary"
          :title="isCollected ? t('common.uncollect') : t('common.collect')"
        >
          <svg viewBox="0 0 24 24" :fill="isCollected ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
          </svg>
        </button>
      </div>
    </div>

    <!-- 主体区域 -->
    <div class="detail-body">
      <!-- 字幕区 -->
      <div class="subtitle-section" v-if="controls.showSentences">
        <div ref="subtitleScrollRef" class="subtitle-scroll">
          <div class="subtitle-padding-top"></div>
          <div
            v-for="(item, index) in subtitleList"
            :key="index"
            :ref="el => { if (index === activeSentenceIndex) activeEl = el }"
            class="sentence-row"
            :class="{ 'is-active': activeSentenceIndex === index }"
            @click="!item._locked && jumpBySubtitle(item, index)"
          >
            <div class="sentence-audio-icon" :class="{ 'audio-playing': activeSentenceIndex === index && isAudioPlaying }">
              <span class="audio-icon-char">{{ activeSentenceIndex === index && isAudioPlaying ? '◉' : '♪' }}</span>
            </div>

            <div v-if="item._locked" class="sentence-locked">
              <span class="lock-time">{{ formatTime(item.startTime) }} - {{ formatTime(item.endTime) }}</span>
              <span class="lock-icon">🔒</span>
            </div>
            <div v-else class="sentence-text-wrap">
              <div class="english-text" v-if="diaryInfo?.showEnglish !== false">
                <template v-for="(seg, sIdx) in parseHighlight(item.english)" :key="sIdx">
                  <span
                    v-if="seg.isHighlight"
                    class="highlight-word"
                    @click.stop="handleSegClick(seg)"
                  >{{ seg.text }}</span>
                  <span v-else>{{ seg.text }}</span>
                </template>
              </div>
              <div class="translated-text" v-if="controls.showTranslate">
                {{ localText(item.translate) }}
              </div>
            </div>

            <div v-if="!item._locked" class="sentence-collect-btn" @click.stop="collectSentence(item)">
              <span class="collect-icon">{{ sentenceCollectedMap[getSentenceId(item)] ? '★' : '☆' }}</span>
            </div>
          </div>
          <div class="subtitle-padding-bottom"></div>
        </div>
      </div>

      <!-- 右侧图片区 -->
      <div class="image-section">
        <div class="image-wrapper">
          <img
            v-if="currentImageUrl"
            :src="currentImageUrl"
            :alt="localText(diaryInfo?.name)"
            class="diary-image"
            @click="viewImage"
          />
          <div v-else class="image-placeholder">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <circle cx="8.5" cy="8.5" r="1.5"></circle>
              <polyline points="21 15 16 10 5 21"></polyline>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部控制栏 -->
    <div class="bottom-controls">
      <div class="control-item" @click="toggleSentences">
        <div class="c-icon-circle" :class="{ 'c-circle-on': controls.showSentences }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
            <path d="M4 6h16M4 12h16M4 18h10"></path>
          </svg>
        </div>
        <span class="c-text" :class="{ 'c-text-active': controls.showSentences }">{{ t('diary.sentence') }}</span>
      </div>

      <div class="control-item" @click="toggleTranslate">
        <div class="c-icon-circle" :class="{ 'c-circle-on': controls.showTranslate }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
            <path d="M5 8l6 6"></path>
            <path d="M4 14l6-6 2-3"></path>
            <path d="M2 5h12"></path>
            <path d="M22 22l-5-10-5 10"></path>
          </svg>
        </div>
        <span class="c-text" :class="{ 'c-text-active': controls.showTranslate }">{{ t('diary.translate') }}</span>
      </div>

      <div class="control-item play-btn" @click="togglePlay">
        <div class="c-play-circle">
          <svg v-if="!isAudioPlaying" viewBox="0 0 24 24" fill="currentColor" width="24" height="24">
            <polygon points="5 3 19 12 5 21 5 3"></polygon>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="currentColor" width="24" height="24">
            <rect x="6" y="4" width="4" height="16"></rect>
            <rect x="14" y="4" width="4" height="16"></rect>
          </svg>
        </div>
      </div>

      <div class="control-item" @click="collectDiary">
        <div class="c-icon-circle" :class="{ 'c-circle-on': isCollected }">
          <svg viewBox="0 0 24 24" :fill="isCollected ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2" width="20" height="20">
            <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
          </svg>
        </div>
        <span class="c-text" :class="{ 'c-text-active': isCollected }">{{ isCollected ? t('diary.uncollect') : t('diary.collect') }}</span>
      </div>

      <div class="control-item" @click="showMoreSheet = true">
        <div class="c-icon-circle">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="20" height="20">
            <circle cx="12" cy="12" r="1"></circle>
            <circle cx="19" cy="12" r="1"></circle>
            <circle cx="5" cy="12" r="1"></circle>
          </svg>
        </div>
        <span class="c-text">{{ t('common.settings') }}</span>
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

    <!-- 更多弹窗 -->
    <div v-if="showMoreSheet" class="more-overlay" @click="showMoreSheet = false">
      <div class="more-sheet" @click.stop>
        <div class="popup-handle"></div>
        <div class="more-sheet-head">
          <span class="more-sheet-title">{{ t('common.settings') }}</span>
        </div>
        <div class="speed-hint">{{ t('player.tap_subtitle_hint') }}</div>
        <div class="more-card">
          <div class="more-card-main">
            <span class="more-card-title">{{ t('player.speed') }}</span>
          </div>
        </div>
        <div class="speed-pills">
          <span v-for="rate in speedRates" :key="rate" class="speed-pill speed-pill-disabled">{{ rate }}x</span>
        </div>
        <div class="more-actions">
          <button class="more-action-btn ghost" @click="handleBack">{{ t('common.back') }}</button>
          <button class="more-action-btn" @click="showMoreSheet = false">{{ t('common.cancel') }}</button>
        </div>
      </div>
    </div>

    <!-- 隐藏音频 -->
    <audio ref="audioRef" style="display:none;"></audio>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  findDiary,
  getDiarySentenceList,
  collect,
  uncollect,
  getCollectionList,
  findWord
} from '@/api/learning'
import { localText, resolveApiMessage } from '@/utils/i18n'
import { getExternalUrl, initExternalDomain } from '@/utils/url'
import { playWordAudio as playWordTTS, stopLocalTTS, destroyTTS } from '@/utils/learning-tts'
import storage from '@/utils/storage'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const diaryId = ref(Number(route.params.id) || 0)

// Audio
const audioRef = ref(null)
let audioCtx = null
const isAudioPlaying = ref(false)
const isAudioLoading = ref(false)
const currentTime = ref(0)
const accent = ref('US')
const isAccentSwitching = ref(false)
const audioDurationUs = ref(0)
const audioDurationUk = ref(0)

const diaryInfo = ref({
  id: 0, name: '', audioUs: '', audioUk: '', duration: 0,
  imageI18n: '', showTranslate: true, showEnglish: true,
  hasFullAuth: true, trialPercent: 100
})

const controls = ref({ showSentences: true, showTranslate: true })
const subtitleList = ref([])
const activeSentenceIndex = ref(-1)
let activeEl = null
const subtitleScrollRef = ref(null)
const isCollected = ref(false)

const showMoreSheet = ref(false)
const speedRates = [0.75, 1.0, 1.25, 1.5, 2.0]

const isTrialLocked = ref(false)
let lastTrialModalAt = 0

// Sentence collect
const sentenceCollectedMap = ref({})

// Word popup
const wordPopup = ref({ visible: false, word: null })
const currentWord = ref(null)
const isWordAudioPlaying = ref(false)
const isWordCollected = ref(false)

const getEntityId = (item) => Number(item?.id || item?.ID || 0)
const getSentenceId = (item) => Number(item?.id || item?.ID || 0)

const currentDuration = computed(() => accent.value === 'US' ? audioDurationUs.value : audioDurationUk.value)

const currentImageUrl = computed(() => {
  if (!diaryInfo.value?.imageI18n) return ''
  const locale = storage.get('app-lang') || 'zh'
  let images = {}
  try { images = typeof diaryInfo.value.imageI18n === 'object' ? diaryInfo.value.imageI18n : JSON.parse(diaryInfo.value.imageI18n || '{}') } catch (e) {}
  const raw = images[locale] || images['en'] || images['zh'] || ''
  return getExternalUrl(raw)
})

const formatTime = (secs) => {
  if (!secs || isNaN(secs)) return '00:00'
  const m = Math.floor(secs / 60).toString().padStart(2, '0')
  const s = Math.floor(secs % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}

const getAudioUrl = () => {
  const raw = accent.value === 'US' ? diaryInfo.value.audioUs : diaryInfo.value.audioUk
  return getExternalUrl(raw)
}

// Audio
const initAudio = () => {
  destroyAudio()
  const url = getAudioUrl()
  if (!url) { showToast('No audio available'); return }
  audioCtx = new Audio()
  audioCtx.src = url
  audioCtx.addEventListener('play', () => { isAudioPlaying.value = true; isAudioLoading.value = false })
  audioCtx.addEventListener('pause', () => { isAudioPlaying.value = false })
  audioCtx.addEventListener('ended', () => { isAudioPlaying.value = false; currentTime.value = 0; activeSentenceIndex.value = -1 })
  audioCtx.addEventListener('timeupdate', () => {
    currentTime.value = audioCtx.currentTime || 0
    const dur = audioCtx.duration || 0
    if (dur && isFinite(dur) && dur > 0) {
      if (accent.value === 'US') audioDurationUs.value = Math.round(dur * 10) / 10
      else audioDurationUk.value = Math.round(dur * 10) / 10
    }
    if (ensureTrialLimit(currentTime.value)) return
    updateActiveSentence()
  })
  audioCtx.addEventListener('error', () => { isAudioPlaying.value = false; isAudioLoading.value = false; isAccentSwitching.value = false; showToast('Audio play failed') })
  audioCtx.addEventListener('canplay', () => {
    isAudioLoading.value = false; isAccentSwitching.value = false
    const dur = audioCtx.duration || 0
    if (dur && isFinite(dur) && dur > 0) {
      if (accent.value === 'US') audioDurationUs.value = Math.round(dur * 10) / 10
      else audioDurationUk.value = Math.round(dur * 10) / 10
    }
  })
  audioCtx.addEventListener('waiting', () => { isAudioLoading.value = true })
}

const destroyAudio = () => {
  if (audioCtx) { audioCtx.pause(); audioCtx.src = ''; audioCtx = null }
}

const playAudio = () => {
  if (!audioCtx) initAudio()
  if (audioCtx) { isAudioLoading.value = true; audioCtx.play().catch(() => { isAudioLoading.value = false }) }
}
const pauseAudio = () => { if (audioCtx) audioCtx.pause() }

const togglePlay = () => {
  if (isAudioPlaying.value) pauseAudio()
  else playAudio()
}

const toggleAccent = () => {
  const was = isAudioPlaying.value
  if (was) pauseAudio()
  accent.value = accent.value === 'US' ? 'UK' : 'US'
  isAccentSwitching.value = true
  nextTick(() => { initAudio(); if (was) playAudio() })
}

const updateActiveSentence = () => {
  const time = currentTime.value
  const useUk = accent.value === 'UK'
  let found = -1
  for (let i = subtitleList.value.length - 1; i >= 0; i--) {
    const item = subtitleList.value[i]
    const st = useUk && item.startTimeUk > 0 ? item.startTimeUk : item.startTime
    if (time >= st) { found = i; break }
  }
  if (found !== activeSentenceIndex.value) {
    activeSentenceIndex.value = found
    nextTick(() => scrollToActiveSentence())
  }
}

// Trial
const maxTrialAllowedTime = () => {
  const dur = currentDuration.value || Number(diaryInfo.value.duration || 0)
  if (diaryInfo.value.hasFullAuth || dur <= 0) return Number.POSITIVE_INFINITY
  const pct = Math.max(1, Number(diaryInfo.value.trialPercent || 8))
  return (dur * pct) / 100
}

const ensureTrialLimit = (time) => {
  if (diaryInfo.value.hasFullAuth) return false
  const limit = maxTrialAllowedTime()
  if (!isFinite(limit) || limit <= 0) return false
  if (time < limit) { if (isTrialLocked.value) isTrialLocked.value = false; return false }
  pauseAudio()
  if (audioCtx) audioCtx.currentTime = Math.max(0, limit - 0.5)
  isTrialLocked.value = true
  const now = Date.now()
  if (now - lastTrialModalAt < 4000) return true
  lastTrialModalAt = now
  if (confirm(t('player.trial_end_title') + '\n' + t('player.trial_end_desc'))) router.push('/learning/profile')
  return true
}

const jumpBySubtitle = (item, index) => {
  if (!audioCtx) return
  if (activeSentenceIndex.value === index && isAudioPlaying.value) { pauseAudio(); return }
  if (activeSentenceIndex.value === index && !isAudioPlaying.value) { playAudio(); return }
  const useUk = accent.value === 'UK'
  const st = useUk && item.startTimeUk > 0 ? item.startTimeUk : item.startTime
  audioCtx.currentTime = st
  activeSentenceIndex.value = index
  nextTick(() => scrollToActiveSentence())
  playAudio()
}

// Highlight parsing
const parseHighlight = (text) => {
  if (!text) return [{ text: '', isHighlight: false }]
  const parts = []; const regex = /<w id="(\d+)">(.*?)<\/w>/g; let lastIdx = 0; let match
  while ((match = regex.exec(text)) !== null) {
    if (match.index > lastIdx) parts.push({ text: text.substring(lastIdx, match.index), isHighlight: false })
    parts.push({ text: match[2], isHighlight: true, wordId: match[1] })
    lastIdx = match.index + match[0].length
  }
  if (lastIdx < text.length) parts.push({ text: text.substring(lastIdx), isHighlight: false })
  return parts.length ? parts : [{ text, isHighlight: false }]
}

// Word popup
const handleSegClick = (seg) => { if (seg?.isHighlight) showWordDetail(seg.wordId) }

const showWordDetail = async (wordId) => {
  if (isAudioPlaying.value) pauseAudio()
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

const closeWordPopup = () => { stopLocalTTS(); isWordAudioPlaying.value = false; wordPopup.value.visible = false }

const playWordAudioAction = () => {
  if (!currentWord.value) return
  if (isWordAudioPlaying.value) { stopLocalTTS(); isWordAudioPlaying.value = false; return }
  playWordTTS(currentWord.value.word, currentWord.value.audioUs || '', currentWord.value.audioUk || '', accent.value, {
    onStart: () => { isWordAudioPlaying.value = true },
    onEnd: () => { isWordAudioPlaying.value = false },
    onError: () => { isWordAudioPlaying.value = false }
  })
}

const collectCurrentWord = async () => {
  const id = getEntityId(currentWord.value); if (!id) return
  const col = isWordCollected.value
  const res = col ? await uncollect(1, id) : await collect(1, id)
  if (res.code === 0) { isWordCollected.value = !col; showToast(!col ? t('player.collect_success') : t('player.uncollect_success')) }
  else showToast(resolveApiMessage(res.msg, 'player.operationFailed'))
}

// Sentence collect (targetType=5)
const collectSentence = async (item) => {
  const sid = getSentenceId(item); if (!sid) return
  const col = !!sentenceCollectedMap.value[sid]
  const res = col ? await uncollect(5, sid) : await collect(5, sid)
  if (res.code === 0) { sentenceCollectedMap.value = { ...sentenceCollectedMap.value, [sid]: !col }; showToast(!col ? t('player.collect_success') : t('player.uncollect_success')) }
  else showToast(resolveApiMessage(res.msg, 'player.operationFailed'))
}

const loadSentenceCollectionState = async () => {
  const ids = subtitleList.value.map(i => getSentenceId(i)).filter(id => id > 0)
  if (ids.length === 0) { sentenceCollectedMap.value = {}; return }
  const idSet = new Set(ids); const map = {}
  try {
    for (let page = 1; page <= 20; page++) {
      const res = await getCollectionList({ targetType: 5, page, pageSize: 100 })
      if (res.code !== 0 || !res.data) break
      const list = Array.isArray(res.data.list) ? res.data.list : []
      for (const item of list) { const tid = Number(item?.targetId || 0); if (tid > 0 && idSet.has(tid)) map[tid] = true }
      if (Object.keys(map).length >= idSet.size) break
      if (list.length < 100) break
    }
    sentenceCollectedMap.value = map
  } catch (e) { sentenceCollectedMap.value = {} }
}

// Data
const loadDiary = async () => {
  const res = await findDiary(diaryId.value)
  if (res.code !== 0 || !res.data) { showToast(t('player.load_failed')); return }
  const d = res.data
  diaryInfo.value = {
    id: getEntityId(d), name: d?.name || '', audioUs: d?.audioUs || '', audioUk: d?.audioUk || '',
    duration: Number(d?.duration || 0), imageI18n: d?.imageI18n || '',
    showTranslate: d?.showTranslate !== false, showEnglish: d?.showEnglish !== false,
    hasFullAuth: d?.hasFullAuth !== false, trialPercent: Number(d?.trialPercent || 8)
  }
  audioDurationUs.value = diaryInfo.value.duration
  audioDurationUk.value = diaryInfo.value.duration
}

const loadSentences = async () => {
  const res = await getDiarySentenceList(diaryId.value)
  if (res.code === 0 && res.data) {
    const limit = maxTrialAllowedTime()
    subtitleList.value = (res.data || []).map((item, index) => {
      const st = Number(item?.startTime || 0)
      return { ...item, _locked: !diaryInfo.value.hasFullAuth && st >= limit }
    })
    await loadSentenceCollectionState()
  }
}

const loadCollectionState = async () => {
  try {
    const res = await getCollectionList({ targetType: 4, page: 1, pageSize: 999 })
    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data.list) ? res.data.list : []
      isCollected.value = list.some(item => Number(item.targetId) === diaryId.value)
    }
  } catch (e) {}
}

const collectDiary = async () => {
  const res = isCollected.value ? await uncollect(4, diaryId.value) : await collect(4, diaryId.value)
  if (res.code === 0) { isCollected.value = !isCollected.value; showToast(!isCollected.value ? t('player.collect_success') : t('player.uncollect_success')) }
  else showToast(resolveApiMessage(res.msg, 'player.operationFailed'))
}

const viewImage = () => {
  const url = currentImageUrl.value
  if (!url) { showToast('No image'); return }
  // Open in new tab / lightbox
  const overlay = document.createElement('div')
  overlay.style.cssText = 'position:fixed;inset:0;background:rgba(0,0,0,0.9);z-index:10000;display:flex;align-items:center;justify-content:center;cursor:pointer;'
  const img = document.createElement('img')
  img.src = url; img.style.cssText = 'max-width:90%;max-height:90%;object-fit:contain;'
  overlay.appendChild(img)
  overlay.onclick = () => overlay.remove()
  document.body.appendChild(overlay)
}

const toggleSentences = () => { controls.value.showSentences = !controls.value.showSentences }
const toggleTranslate = () => { controls.value.showTranslate = !controls.value.showTranslate }

const scrollToActiveSentence = () => {
  const c = subtitleScrollRef.value; if (!c || !activeEl) return
  const cr = c.getBoundingClientRect(); const ar = activeEl.getBoundingClientRect()
  c.scrollBy({ top: ar.top - cr.top - c.clientHeight / 2 + ar.height / 2, behavior: 'smooth' })
}

const goBack = () => { if (window.history.length > 1) router.back(); else router.push('/learning/diary') }
const handleBack = () => { showMoreSheet.value = false; goBack() }

const showToast = (msg) => {
  const el = document.createElement('div')
  el.textContent = msg
  el.style.cssText = 'position:fixed;top:20px;left:50%;transform:translateX(-50%);background:rgba(26,27,58,0.85);color:#fff;padding:10px 24px;border-radius:999px;z-index:9999;font-size:14px;pointer-events:none;'
  document.body.appendChild(el); setTimeout(() => el.remove(), 2500)
}

onMounted(async () => {
  await initExternalDomain()
  await loadDiary()
  await loadSentences()
  await loadCollectionState()
  nextTick(() => initAudio())
})

onBeforeUnmount(() => {
  destroyAudio()
  destroyTTS()
})
</script>

<style lang="scss" scoped>
.diary-detail-page {
  display: flex; flex-direction: column;
  height: calc(100vh - #{$header-height} - 48px);
  max-width: $content-max-width; margin: 0 auto;
}
.detail-header {
  display: flex; align-items: center; padding: $spacing-md 0; margin-bottom: $spacing-md; gap: 16px; flex-shrink: 0;
  .back-btn {
    display: flex; align-items: center; gap: 6px; padding: 8px 16px;
    background: $bg-card; border-radius: $radius-pill; color: $text-secondary; font-size: 14px;
    cursor: pointer; transition: all 0.2s ease; border: none;
    &:hover { background: $bg-hover; color: $primary-color; }
    svg { width: 16px; height: 16px; }
  }
  .diary-title { flex: 1; margin: 0; font-size: 18px; font-weight: $font-weight-semibold; color: $text-primary; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .header-actions { display: flex; align-items: center; gap: 12px; flex-shrink: 0;
    .duration-info { display: flex; align-items: center; gap: 4px; font-size: 13px; color: $text-secondary;
      svg { width: 16px; height: 16px; }
    }
    .accent-switch-btn {
      display: flex; align-items: center; gap: 2px; padding: 5px 14px;
      background: $bg-input; border-radius: $radius-pill; font-size: 12px; font-weight: $font-weight-medium;
      color: $text-secondary; border: none; cursor: pointer; transition: all 0.2s;
      span { color: $text-tertiary; font-weight: 600; &.active { color: $primary-color; } &.divider { margin: 0 2px; opacity: 0.4; } }
    }
    .image-btn, .action-btn {
      width: 36px; height: 36px; display: flex; align-items: center; justify-content: center;
      background: $bg-card; border-radius: 50%; color: $text-secondary; border: none; cursor: pointer; transition: all 0.2s;
      &:hover, &.active { color: $primary-color; background: $bg-active; }
      svg { width: 18px; height: 18px; }
    }
  }
}
.detail-body { flex: 1; display: flex; gap: $spacing-lg; min-height: 0; }

.subtitle-section {
  flex: 1; min-width: 0; background: $bg-card; border-radius: $radius-lg;
  display: flex; flex-direction: column; box-shadow: $shadow-card; border: 1px solid $border-light; overflow: hidden;
}
.subtitle-scroll { flex: 1; overflow-y: auto; background: $bg-page; }
.subtitle-padding-top { height: 12px; }
.subtitle-padding-bottom { height: 16px; }

.sentence-row {
  display: flex; align-items: flex-start; gap: 12px; padding: 14px 18px; margin: 0 12px 8px;
  transition: all 0.25s ease; border-radius: $radius-md; background: #fff; border: 1px solid rgba(109,91,255,0.08); cursor: pointer;
  &:hover { background: rgba(109,91,255,0.03); }
  &.is-active { background: rgba(109,91,255,0.06); border-left: 4px solid $primary-color; box-shadow: 0 4px 12px rgba(109,91,255,0.12); }
  .sentence-audio-icon {
    width: 32px; height: 32px; border-radius: 50%; background: $bg-page; display: flex; align-items: center; justify-content: center;
    flex-shrink: 0; margin-top: 2px;
    &.audio-playing { background: $gradient-primary; animation: pulse-icon 0.8s ease-in-out infinite alternate; }
    .audio-icon-char { font-size: 14px; color: $text-tertiary; }
    &.audio-playing .audio-icon-char { color: #fff; }
  }
  .sentence-text-wrap { flex: 1; min-width: 0; }
  .english-text { font-size: 14px; font-weight: 500; line-height: 1.8; color: $text-primary; }
  .highlight-word { color: $primary-color; font-weight: 700; border-bottom: 2px dashed rgba(109,91,255,0.4); display: inline; padding: 0 2px; margin: 0 2px; cursor: pointer;
    &:hover { background: rgba(109,91,255,0.1); }
  }
  .translated-text { font-size: 13px; color: $text-secondary; line-height: 1.6; margin-top: 6px; padding-top: 6px; border-top: 1px solid rgba(109,91,255,0.08); }
  .sentence-collect-btn { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; cursor: pointer;
    .collect-icon { font-size: 16px; color: $text-tertiary; transition: color 0.2s; }
    &:hover .collect-icon { color: $primary-color; }
  }
}
.sentence-locked { flex: 1; display: flex; align-items: center; justify-content: center; gap: 8px; padding: 12px;
  .lock-time { font-size: 12px; color: $text-tertiary; font-variant-numeric: tabular-nums; }
  .lock-icon { font-size: 18px; opacity: 0.4; }
}
@keyframes pulse-icon { 0% { transform: scale(1); } 100% { transform: scale(1.12); } }

.image-section { width: 340px; flex-shrink: 0; display: flex; flex-direction: column; gap: $spacing-md; }
.image-wrapper { position: relative; width: 100%; background: $bg-card; border-radius: $radius-lg; overflow: hidden; box-shadow: $shadow-card; border: 1px solid $border-light; aspect-ratio: 4/3;
  .diary-image { width: 100%; height: 100%; object-fit: cover; cursor: pointer; }
  .image-placeholder { position: absolute; inset: 0; display: flex; align-items: center; justify-content: center; color: $border-primary; background: $gradient-card;
    svg { width: 48px; height: 48px; }
  }
}

.bottom-controls {
  flex-shrink: 0; display: flex; align-items: center; justify-content: space-around;
  padding: 14px 24px; background: $bg-card; border-radius: $radius-lg; margin: $spacing-md 0;
  box-shadow: $shadow-card; border: 1px solid $border-light;
  .control-item { display: flex; flex-direction: column; align-items: center; gap: 4px; cursor: pointer; }
  .c-icon-circle {
    width: 44px; height: 44px; border-radius: 50%; background: $bg-page; display: flex; align-items: center; justify-content: center;
    transition: all 0.25s ease; color: $text-secondary;
  }
  .c-circle-on { background: $gradient-primary; color: #fff; box-shadow: 0 4px 14px rgba(109,91,255,0.35); }
  .c-text { font-size: 11px; color: $text-tertiary; font-weight: 500; }
  .c-text-active { color: $primary-color; font-weight: 600; }
  .c-play-circle {
    width: 56px; height: 56px; border-radius: 50%; background: $gradient-primary; display: flex; align-items: center; justify-content: center;
    color: #fff; box-shadow: 0 6px 20px rgba(109,91,255,0.4);
  }
}

.word-popup-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.3); display: flex; align-items: center; justify-content: center; z-index: $z-modal; }
.word-detail-box { width: 400px; max-width: 90vw; padding: 24px 28px; background: #fff; border-radius: $radius-lg; box-shadow: $shadow-xl; }
.w-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 4px; }
.w-title { font-size: 24px; font-weight: 700; color: $text-primary; }
.w-audio-circle { width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center; background: $gradient-primary; border: none; cursor: pointer; color: #fff;
  &.playing { animation: pulse 0.6s ease-in-out infinite; pointer-events: none; }
}
@keyframes pulse { 0%, 100% { transform: scale(1); opacity: 1; } 50% { transform: scale(1.12); opacity: 0.75; } }
.w-audio-icon { font-size: 16px; }
.w-phonetic { font-size: 14px; color: $text-secondary; display: block; margin: 4px 0 12px; }
.w-divider { height: 1px; background: rgba(109,91,255,0.16); margin-bottom: 16px; }
.w-exp { font-size: 15px; color: $text-primary; line-height: 1.6; margin-bottom: 20px; }
.w-actions { display: flex; gap: 12px; }
.w-cancel-btn { flex: 1; height: 44px; border-radius: $radius-pill; border: 1.5px solid rgba(109,91,255,0.32); background: #fff; color: $primary-color; font-size: 15px; font-weight: 600; cursor: pointer; }
.w-collect { flex: 1; height: 44px; border: none; background: $gradient-primary; color: #fff; border-radius: $radius-pill; font-size: 15px; font-weight: 600; cursor: pointer; box-shadow: 0 4px 14px rgba(109,91,255,0.3); }

.more-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4); z-index: $z-modal; display: flex; align-items: flex-end; justify-content: center; }
.more-sheet { width: 100%; max-width: 500px; background: #fff; border-radius: 20px 20px 0 0; padding: 8px 0 24px;
  .popup-handle { width: 40px; height: 4px; background: #D1D5DB; border-radius: 4px; margin: 12px auto; }
  .more-sheet-head { padding: 8px 24px 16px; border-bottom: 1px solid #F0EEFF;
    .more-sheet-title { font-size: 18px; font-weight: 700; color: $text-primary; }
  }
  .speed-hint { padding: 16px 24px 0; font-size: 13px; color: $primary-color; font-weight: 600; }
  .more-card { padding: 16px 24px; .more-card-main { .more-card-title { font-size: 15px; color: $text-primary; font-weight: 600; } } }
  .speed-pills { display: flex; flex-wrap: wrap; gap: 8px; padding: 8px 24px 0;
    .speed-pill { flex: 1; min-width: 60px; height: 36px; line-height: 36px; text-align: center; border-radius: 999px; border: 1px solid rgba(109,91,255,0.16); font-size: 13px; font-weight: 600; color: $text-secondary; background: #fff; }
    .speed-pill-disabled { color: $text-tertiary; border-color: #E8E8ED; background: #F5F5F8; }
  }
  .more-actions { display: flex; gap: 12px; padding: 20px 24px 0;
    .more-action-btn { flex: 1; height: 44px; border-radius: $radius-pill; background: $gradient-primary; color: #fff; font-size: 14px; font-weight: 600; border: none; cursor: pointer; box-shadow: 0 4px 14px rgba(109,91,255,0.32);
      &.ghost { background: $bg-page; color: $primary-color; border: 1.5px solid rgba(109,91,255,0.2); box-shadow: none; }
    }
  }
}

@media (max-width: 1000px) {
  .detail-body { flex-direction: column; }
  .image-section { width: 100%; }
  .image-wrapper { aspect-ratio: 16/9; }
}
</style>
