<template>
  <view class="game-category">
    <!-- 背景装饰 -->
    <view class="bg-glow bg-glow-1"></view>
    <view class="bg-glow bg-glow-2"></view>

    <!-- 顶部栏 -->
    <view class="top-bar">
      <view class="back-btn" @click="playClick(); goBack()">
        <text class="back-icon">‹</text>
      </view>
      <text class="title">{{ gameName }}</text>
      <view class="placeholder" />
    </view>

    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('game.loading') }}</text>
    </view>

    <template v-else-if="difficultyList.length > 0">
      <!-- 难度分类 Tab -->
      <scroll-view class="cat-tabs" scroll-x :show-scrollbar="false">
        <view class="cat-tabs-inner">
          <view
            v-for="(cat, ci) in difficultyList"
            :key="cat.ID"
            class="cat-tab"
            :class="{ active: ci === currentIndex }"
            @click="playClick(); currentIndex = ci"
          >
            <text class="cat-tab-text">{{ localText(cat.name) }}</text>
          </view>
        </view>
      </scroll-view>

      <!-- 关卡网格 -->
      <view class="levels-wrap">
        <view class="level-grid">
          <view
            v-for="level in getLevelsByCategory(difficultyList[currentIndex].ID)"
            :key="level.ID"
            class="level-cell"
            :class="{
              locked: !isLevelUnlocked(difficultyList[currentIndex].ID, level),
              cleared: isLevelCleared(difficultyList[currentIndex].ID, level)
            }"
            @click="playClick(); onLevelClick(difficultyList[currentIndex], level)"
          >
            <view v-if="!isLevelUnlocked(difficultyList[currentIndex].ID, level)" class="level-locked">
              <text class="lock-icon">🔒</text>
            </view>
            <template v-else>
              <view class="level-inner">
                <text class="level-num">{{ level.levelNumber }}</text>
                <view v-if="isLevelCleared(difficultyList[currentIndex].ID, level)" class="cleared-badge">
                  <text class="check-icon">✓</text>
                </view>
              </view>
            </template>
          </view>
        </view>
      </view>
    </template>

    <view v-else class="empty-wrap">
      <text class="empty-icon">📭</text>
      <text class="empty-text">{{ t('game.noData') }}</text>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { getGameCategory, getDifficultyCategories, getUserProgress, getLevelList } from '@/api/game.js'

const langStore = useLangStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const t = (key) => {
  const text = i18nT(key, locale.value)
  return text !== key ? text : key
}
const localText = (value) => i18nLocalText(value, locale.value)

const gameID = ref(0)
const gameKey = ref('')
const playPage = ref('play')
const gameName = ref('')
const loading = ref(true)
const difficultyList = ref([])
const currentIndex = ref(0)
const clickAudio = ref(null)
const levelCache = reactive({})

const initAudio = () => {
  const soundEnabled = uni.getStorageSync('game-sound-enabled') !== '0'
  if (soundEnabled) {
    clickAudio.value = uni.createInnerAudioContext()
    clickAudio.value.src = '/static/click-click.mp3'
  }
}

const playClick = () => {
  if (clickAudio.value) {
    clickAudio.value.stop()
    clickAudio.value.play()
  }
}

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length <= 1) {
    uni.reLaunch({ url: '/pages/game/index' })
  } else {
    uni.navigateBack()
  }
}

const getLevelsByCategory = (catID) => {
  if (!levelCache[catID]) return []
  return levelCache[catID].levels || []
}

const isLevelCleared = (catID, level) => {
  if (!levelCache[catID] || !levelCache[catID].progress) return false
  return levelCache[catID].progress.some(p => p.levelID === level.ID && p.status === 1)
}

const isLevelUnlocked = (catID, level) => {
  const catIndex = difficultyList.value.findIndex(c => c.ID === catID)
  if (catIndex < 0) return false

  // 第一个难度分类的第一关默认解锁
  if (catIndex === 0 && level.levelNumber === 1) return true

  // 同一分类内，前一关通过则解锁
  if (level.levelNumber > 1) {
    const prevLevel = (levelCache[catID]?.levels || []).find(l => l.levelNumber === level.levelNumber - 1)
    if (prevLevel && isLevelCleared(catID, prevLevel)) return true
  }

  // 非第一个分类的第一关：需要前一个分类全部通关
  if (level.levelNumber === 1 && catIndex > 0) {
    const prevCat = difficultyList.value[catIndex - 1]
    if (prevCat && levelCache[prevCat.ID]) {
      const prevLevels = levelCache[prevCat.ID].levels || []
      if (prevLevels.length > 0 && prevLevels.every(l => isLevelCleared(prevCat.ID, l))) return true
    }
  }

  return false
}

const onLevelClick = (cat, level) => {
  if (!isLevelUnlocked(cat.ID, level)) return
  uni.navigateTo({ url: `/pages/game/${playPage.value}?levelID=${level.ID}&catID=${cat.ID}` })
}

const loadGameName = async () => {
  try {
    const res = await getGameCategory(gameID.value)
    if (res.code === 0 && res.data) {
      gameName.value = localText(res.data.name)
      playPage.value = res.data.playPage || (gameKey.value === 'pwd-guess' ? 'pwd-play' : 'play')
    }
  } catch (e) { /* ignore */ }
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await getDifficultyCategories(gameID.value)
    if (res.code === 0) {
      difficultyList.value = res.data || []
      await refreshAllProgress()
    }
  } catch (e) { /* ignore */ }
  loading.value = false
}

const refreshAllProgress = async () => {
  for (const cat of difficultyList.value) {
    await loadCategoryProgress(cat)
  }
}

const loadCategoryProgress = async (cat) => {
  try {
    const [levelRes, progressRes] = await Promise.all([
      getLevelList(cat.ID),
      getUserProgress(cat.ID)
    ])
    levelCache[cat.ID] = {
      levels: levelRes.code === 0 ? (levelRes.data || []) : [],
      progress: progressRes.code === 0 ? (progressRes.data || []) : []
    }
  } catch (e) {
    levelCache[cat.ID] = { levels: [], progress: [] }
  }
}

onMounted(() => {
  initAudio()
  const pages = getCurrentPages()
  const current = pages[pages.length - 1]
  if (current.$page) {
    const opts = current.$page.options || {}
    gameID.value = parseInt(opts.gameID) || 0
    gameKey.value = opts.gameKey || ''
  }
  loadGameName()
  loadData()
})

onShow(() => {
  if (difficultyList.value.length > 0) {
    refreshAllProgress()
  }
})
</script>

<style scoped>
.game-category {
  min-height: 100vh;
  background: #0a0a1a;
  position: relative;
  overflow: hidden;
}

/* 背景光晕 */
.bg-glow {
  position: fixed;
  border-radius: 50%;
  filter: blur(120rpx);
  opacity: 0.35;
  pointer-events: none;
  z-index: 0;
}
.bg-glow-1 {
  width: 450rpx;
  height: 450rpx;
  background: #e94560;
  top: -100rpx;
  left: -100rpx;
}
.bg-glow-2 {
  width: 400rpx;
  height: 400rpx;
  background: #533483;
  bottom: 200rpx;
  right: -100rpx;
}

/* 顶部栏 */
.top-bar {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  padding: 80rpx 32rpx 20rpx;
}
.back-btn {
  width: 72rpx;
  height: 72rpx;
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(10px);
}
.back-icon {
  font-size: 48rpx;
  color: #fff;
  line-height: 1;
  font-weight: 300;
}
.title {
  flex: 1;
  text-align: center;
  font-size: 34rpx;
  color: #fff;
  font-weight: 700;
  letter-spacing: 1rpx;
}
.placeholder { width: 72rpx; }

/* 分类 Tab */
.cat-tabs {
  position: relative;
  z-index: 1;
  white-space: nowrap;
  padding: 20rpx 32rpx 10rpx;
}
.cat-tabs-inner {
  display: inline-flex;
  gap: 16rpx;
  background: rgba(255,255,255,0.06);
  border: 1rpx solid rgba(255,255,255,0.08);
  border-radius: 40rpx;
  padding: 8rpx;
  backdrop-filter: blur(10px);
}
.cat-tab {
  padding: 16rpx 36rpx;
  border-radius: 32rpx;
}
.cat-tab.active {
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  box-shadow: 0 4rpx 16rpx rgba(233, 69, 96, 0.4);
}
.cat-tab-text {
  font-size: 26rpx;
  color: rgba(255,255,255,0.7);
  font-weight: 500;
}
.cat-tab.active .cat-tab-text {
  color: #fff;
  font-weight: 600;
}

/* 关卡网格 */
.levels-wrap {
  position: relative;
  z-index: 1;
  padding: 40rpx 32rpx;
}
.level-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24rpx;
}
.level-cell {
  aspect-ratio: 1;
  position: relative;
  border-radius: 24rpx;
  background: linear-gradient(145deg, rgba(255,255,255,0.1), rgba(255,255,255,0.04));
  border: 1rpx solid rgba(255,255,255,0.1);
  backdrop-filter: blur(10px);
}
.level-cell.locked {
  opacity: 0.35;
}
.level-cell.cleared {
  background: linear-gradient(145deg, rgba(233, 69, 96, 0.25), rgba(233, 69, 96, 0.1));
  border: 1rpx solid rgba(233, 69, 96, 0.5);
  box-shadow: 0 4rpx 16rpx rgba(233, 69, 96, 0.2);
}
.level-locked {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.lock-icon {
  font-size: 40rpx;
}
.level-inner {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}
.level-num {
  font-size: 44rpx;
  color: #fff;
  font-weight: 800;
  font-style: italic;
}
.cleared-badge {
  position: absolute;
  top: -6rpx;
  right: -6rpx;
  width: 36rpx;
  height: 36rpx;
  background: linear-gradient(135deg, #4caf50, #2e7d32);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2rpx 8rpx rgba(76, 175, 80, 0.4);
}
.check-icon {
  font-size: 22rpx;
  color: #fff;
  font-weight: bold;
}

/* 加载/空状态 */
.loading-wrap, .empty-wrap {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}
.loading-text, .empty-text {
  font-size: 28rpx;
  color: rgba(255,255,255,0.4);
}
.empty-icon {
  font-size: 80rpx;
  margin-bottom: 20rpx;
}
</style>
