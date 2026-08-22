<template>
  <view class="game-play">
    <!-- 背景装饰 -->
    <view class="bg-glow bg-glow-1"></view>
    <view class="bg-glow bg-glow-2"></view>

    <!-- 顶部栏 -->
    <view class="top-bar">
      <view class="back-btn" @click="playClick(); goBack()">
        <text class="back-icon">‹</text>
      </view>
      <view class="level-info">
        <text class="level-title">{{ difficultyName }}</text>
        <text class="level-sub">第 {{ currentLevelNumber }} 关</text>
      </view>
      <view class="placeholder" />
    </view>

    <!-- 目标卡片 -->
    <view class="target-card">
      <text class="target-label">{{ t('game.target') }}</text>
      <text class="target-num">{{ targetResult }}</text>
    </view>

    <!-- 数字卡片 -->
    <view class="number-grid">
      <view
        v-for="(num, idx) in numbers"
        :key="idx"
        class="number-card"
        :class="{
          selected: num.selected,
          used: num.used
        }"
        @click="onNumberClick(idx)"
      >
        <text class="num-text">{{ num.value }}</text>
      </view>
    </view>

    <!-- 运算符 -->
    <view class="operator-row">
      <view
        v-for="op in operators"
        :key="op"
        class="operator-btn"
        :class="{ active: selectedOp === op }"
        @click="onOperatorClick(op)"
      >
        <text class="op-text">{{ op }}</text>
      </view>
    </view>

    <!-- 操作按钮 -->
    <view class="action-row">
      <view class="action-btn undo-btn" @click="playClick(); undoLast()">
        <text class="action-icon">↩</text>
        <text class="action-text">{{ t('common.back') }}</text>
      </view>
      <view class="action-btn reset-btn" @click="playClick(); resetGame()">
        <text class="action-icon">↻</text>
        <text class="action-text">{{ t('game.reset') }}</text>
      </view>
    </view>

    <!-- 计算历史 -->
    <view class="history-card">
      <text class="history-title">{{ t('game.history') }}</text>
      <view class="history-list">
        <view v-if="history.length === 0" class="history-empty">{{ t('game.noHistory') }}</view>
        <view v-for="(h, idx) in history" :key="idx" class="history-item">
          <text class="history-step">{{ idx + 1 }}</text>
          <text class="history-text">{{ h.text }}</text>
        </view>
      </view>
    </view>

    <!-- 成功弹窗 -->
    <view v-if="showSuccess" class="modal-overlay">
      <view class="modal-card success">
        <view class="modal-icon-wrap success-icon">
          <text class="modal-icon">🎉</text>
        </view>
        <text class="modal-title">{{ t('game.success') }}</text>
        <text class="modal-desc">恭喜你完成了这一关！</text>
        <view class="modal-actions">
          <view class="modal-btn primary" @click="playClick(); goNextLevel()">
            <text>{{ t('game.nextLevel') }}</text>
          </view>
          <view class="modal-btn ghost" @click="playClick(); goBackToHome()">
            <text>{{ t('game.selectCategory') }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 失败弹窗 -->
    <view v-if="showFail" class="modal-overlay">
      <view class="modal-card fail">
        <view class="modal-icon-wrap fail-icon">
          <text class="modal-icon">😅</text>
        </view>
        <text class="modal-title">{{ t('game.fail') }}</text>
        <text class="modal-desc">再试一次，你可以的！</text>
        <view class="modal-actions">
          <view class="modal-btn primary" @click="playClick(); resetGame()">
            <text>{{ t('game.retry') }}</text>
          </view>
          <view class="modal-btn ghost" @click="playClick(); goBackToHome()">
            <text>{{ t('game.selectCategory') }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText } from '@/utils/i18n.js'
import { getLevelDetail, getDifficultyCategory, getDifficultyCategories, getLevelList, submitLevelResult } from '@/api/game.js'

const langStore = useLangStore()
const locale = ref(langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const t = (key) => {
  const text = i18nT(key, locale.value)
  return text !== key ? text : key
}

const levelID = ref(0)
const catID = ref(0)
const currentLevelNumber = ref(0)
const targetResult = ref(24)
const difficultyName = ref('')

const numbers = reactive([])
const operators = ['+', '-', '×', '÷']
const selectedOp = ref('')
const selectedNumIdx = ref(-1)
const history = ref([])
const showSuccess = ref(false)
const showFail = ref(false)
const submitted = ref(false)

const clickAudio = ref(null)

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length <= 1) {
    uni.reLaunch({ url: '/pages/game/index' })
  } else {
    uni.navigateBack()
  }
}

const goBackToHome = () => {
  uni.reLaunch({ url: '/pages/game/index' })
}

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

const onNumberClick = (idx) => {
  if (numbers[idx].used) return
  playClick()

  if (selectedNumIdx.value === -1) {
    selectedNumIdx.value = idx
    numbers[idx].selected = true
  } else if (!selectedOp.value) {
    // 未选运算符，切换选中数字
    numbers[selectedNumIdx.value].selected = false
    selectedNumIdx.value = idx
    numbers[idx].selected = true
  } else if (selectedOp.value) {
    const a = numbers[selectedNumIdx.value]
    const b = numbers[idx]
    if (a === b) return

    let result = 0
    switch (selectedOp.value) {
      case '+': result = a.value + b.value; break
      case '-': result = a.value - b.value; break
      case '×': result = a.value * b.value; break
      case '÷':
        if (b.value === 0) {
          uni.showToast({ title: t('game.divideByZero'), icon: 'none' })
          return
        }
        result = a.value / b.value
        if (!Number.isInteger(result)) {
          uni.showToast({ title: t('game.notInteger'), icon: 'none' })
          return
        }
        break
    }

    history.value.push({
      text: `${a.value} ${selectedOp.value} ${b.value} = ${result}`
    })

    a.value = result
    a.selected = true
    b.used = true

    selectedOp.value = ''

    checkResult()
  }
}

const onOperatorClick = (op) => {
  if (selectedNumIdx.value === -1) return
  playClick()
  selectedOp.value = op
}

const checkResult = () => {
  const remaining = numbers.filter(n => !n.used)
  if (remaining.length === 1) {
    if (remaining[0].value === targetResult.value) {
      setTimeout(() => {
        showSuccess.value = true
        submitResult()
      }, 500)
    } else {
      setTimeout(() => {
        showFail.value = true
      }, 500)
    }
  }
}

const submitResult = async () => {
  if (submitted.value) return
  submitted.value = true
  try {
    await submitLevelResult(levelID.value)
  } catch (e) { /* ignore */ }
}

const resetGame = () => {
  showSuccess.value = false
  showFail.value = false
  submitted.value = false
  selectedNumIdx.value = -1
  selectedOp.value = ''
  history.value = []
  for (let i = 0; i < numbers.length; i++) {
    numbers[i].value = numbers[i].originalValue
    numbers[i].selected = false
    numbers[i].used = false
  }
}

const undoLast = () => {
  if (history.value.length === 0) return
  resetGame()
}

const goNextLevel = async () => {
  showSuccess.value = false
  try {
    const res = await getLevelList(catID.value)
    if (res.code === 0) {
      const levels = (res.data || []).sort((a, b) => a.levelNumber - b.levelNumber)
      // 用 levelID 精确定位当前关卡
      const curIdx = levels.findIndex(l => l.ID === levelID.value)
      if (curIdx >= 0 && curIdx < levels.length - 1) {
        const nextLevel = levels[curIdx + 1]
        uni.redirectTo({ url: `/pages/game/play?levelID=${nextLevel.ID}&catID=${catID.value}` })
        return
      }
    }
    // 当前分类最后一关，尝试跳到下一个难度分类的第一关
    const catRes = await getDifficultyCategory(catID.value)
    if (catRes.code === 0 && catRes.data) {
      const gameID = catRes.data.gameID
      const catsRes = await getDifficultyCategories(gameID)
      if (catsRes.code === 0) {
        const cats = (catsRes.data || []).sort((a, b) => a.sort - b.sort)
        const curCatIdx = cats.findIndex(c => c.ID === catID.value)
        if (curCatIdx >= 0 && curCatIdx < cats.length - 1) {
          const nextCat = cats[curCatIdx + 1]
          const nextLevelsRes = await getLevelList(nextCat.ID)
          if (nextLevelsRes.code === 0) {
            const nextLevels = (nextLevelsRes.data || []).sort((a, b) => a.levelNumber - b.levelNumber)
            if (nextLevels.length > 0) {
              uni.redirectTo({ url: `/pages/game/play?levelID=${nextLevels[0].ID}&catID=${nextCat.ID}` })
              return
            }
          }
        }
      }
    }
  } catch (e) { /* ignore */ }
  uni.navigateBack()
}

const loadDifficultyName = async () => {
  try {
    const res = await getDifficultyCategory(catID.value)
    if (res.code === 0 && res.data) {
      difficultyName.value = localText(res.data.name, locale.value)
    }
  } catch (e) { /* ignore */ }
}

const loadLevel = async () => {
  try {
    const res = await getLevelDetail(levelID.value)
    if (res.code === 0 && res.data) {
      let numsStr = res.data.numbers || ''
      let target = res.data.targetResult || 24

      // 兼容 gameData 统一字段：从 gameData JSON 中解析
      if (!numsStr && res.data.gameData) {
        try {
          const gd = typeof res.data.gameData === 'string' ? JSON.parse(res.data.gameData) : res.data.gameData
          if (gd.numbers) numsStr = gd.numbers
          if (gd.targetResult) target = gd.targetResult
        } catch (e) { /* ignore */ }
      }

      const nums = numsStr.split(',').map(Number)
      targetResult.value = target
      currentLevelNumber.value = res.data.levelNumber || 0
      numbers.length = 0
      nums.forEach((n, i) => {
        numbers.push({
          value: n,
          originalValue: n,
          selected: false,
          used: false,
          computed: false,
          originalIdx: i
        })
      })
    }
  } catch (e) {
    uni.showToast({ title: t('game.loading'), icon: 'none' })
  }
}

onMounted(() => {
  initAudio()
  const pages = getCurrentPages()
  const current = pages[pages.length - 1]
  if (current.$page) {
    const opts = current.$page.options || {}
    levelID.value = parseInt(opts.levelID) || 0
    catID.value = parseInt(opts.catID) || 0
  }
  loadDifficultyName()
  loadLevel()
})
</script>

<style scoped>
.game-play {
  min-height: 100vh;
  background: #0a0a1a;
  padding-bottom: 40rpx;
  position: relative;
  overflow: hidden;
}

/* 背景光晕 */
.bg-glow {
  position: fixed;
  border-radius: 50%;
  filter: blur(120rpx);
  opacity: 0.3;
  pointer-events: none;
  z-index: 0;
}
.bg-glow-1 {
  width: 400rpx;
  height: 400rpx;
  background: #e94560;
  top: -80rpx;
  right: -80rpx;
}
.bg-glow-2 {
  width: 350rpx;
  height: 350rpx;
  background: #533483;
  bottom: 100rpx;
  left: -80rpx;
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
.level-info {
  flex: 1;
  text-align: center;
}
.level-title {
  display: block;
  font-size: 32rpx;
  color: #fff;
  font-weight: 700;
}
.level-sub {
  display: block;
  font-size: 22rpx;
  color: rgba(255,255,255,0.5);
  margin-top: 4rpx;
}
.placeholder { width: 72rpx; }

/* 目标卡片 */
.target-card {
  position: relative;
  z-index: 1;
  margin: 30rpx 32rpx 0;
  background: linear-gradient(135deg, rgba(233, 69, 96, 0.2), rgba(233, 69, 96, 0.05));
  border: 1rpx solid rgba(233, 69, 96, 0.3);
  border-radius: 28rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  backdrop-filter: blur(10px);
}
.target-label {
  font-size: 28rpx;
  color: rgba(255,255,255,0.7);
  font-weight: 500;
}
.target-num {
  font-size: 64rpx;
  font-weight: 900;
  color: #e94560;
  font-style: italic;
  text-shadow: 0 0 30rpx rgba(233, 69, 96, 0.5);
}

/* 数字卡片 */
.number-grid {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24rpx;
  margin: 50rpx auto 0;
  padding: 0 32rpx;
  max-width: 500rpx;
}
.number-card {
  aspect-ratio: 1;
  background: linear-gradient(145deg, rgba(255,255,255,0.12), rgba(255,255,255,0.04));
  border: 2rpx solid rgba(255,255,255,0.1);
  border-radius: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
}
.number-card.selected {
  border-color: #e94560;
  background: linear-gradient(145deg, rgba(233, 69, 96, 0.3), rgba(233, 69, 96, 0.1));
  transform: translateY(-6rpx) scale(1.05);
  box-shadow: 0 12rpx 30rpx rgba(233, 69, 96, 0.4);
}
.number-card.used {
  opacity: 0.2;
  transform: scale(0.85);
}
.num-text {
  font-size: 80rpx;
  color: #fff;
  font-weight: 900;
  font-style: italic;
}

/* 运算符 */
.operator-row {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: center;
  gap: 24rpx;
  margin-top: 50rpx;
  padding: 0 32rpx;
}
.operator-btn {
  width: 110rpx;
  height: 110rpx;
  background: linear-gradient(145deg, rgba(255,255,255,0.1), rgba(255,255,255,0.04));
  border: 2rpx solid rgba(255,255,255,0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
}
.operator-btn.active {
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  border-color: transparent;
  transform: scale(1.1);
  box-shadow: 0 8rpx 24rpx rgba(233, 69, 96, 0.5);
}
.op-text {
  font-size: 48rpx;
  color: #fff;
  font-weight: 700;
}

/* 操作按钮 */
.action-row {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: center;
  gap: 24rpx;
  margin-top: 40rpx;
  padding: 0 32rpx;
}
.action-btn {
  flex: 1;
  max-width: 260rpx;
  height: 88rpx;
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  transition: all 0.2s;
}
.reset-btn {
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  box-shadow: 0 6rpx 20rpx rgba(233, 69, 96, 0.4);
}
.undo-btn {
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.12);
}
.action-icon {
  font-size: 32rpx;
  color: #fff;
}
.action-text {
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
}

/* 历史记录 */
.history-card {
  position: relative;
  z-index: 1;
  margin: 40rpx 32rpx 0;
  background: linear-gradient(145deg, rgba(255,255,255,0.06), rgba(255,255,255,0.02));
  border: 1rpx solid rgba(255,255,255,0.08);
  border-radius: 24rpx;
  padding: 24rpx;
  backdrop-filter: blur(10px);
}
.history-title {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
  font-weight: 500;
  margin-bottom: 16rpx;
  display: block;
}
.history-list {
  max-height: 240rpx;
  overflow-y: auto;
}
.history-empty {
  text-align: center;
  padding: 30rpx 0;
  font-size: 24rpx;
  color: rgba(255,255,255,0.3);
}
.history-item {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 14rpx 0;
  border-bottom: 1rpx solid rgba(255,255,255,0.04);
}
.history-item:last-child {
  border-bottom: none;
}
.history-step {
  width: 36rpx;
  height: 36rpx;
  line-height: 36rpx;
  text-align: center;
  background: rgba(233, 69, 96, 0.2);
  color: #e94560;
  font-size: 20rpx;
  font-weight: 700;
  border-radius: 50%;
  flex-shrink: 0;
}
.history-text {
  font-size: 26rpx;
  color: rgba(255,255,255,0.7);
  font-family: monospace;
}

/* 弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.75);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
}
.modal-card {
  width: 600rpx;
  background: linear-gradient(180deg, #1a1a2e 0%, #0f0f1e 100%);
  border: 1rpx solid rgba(255,255,255,0.1);
  border-radius: 36rpx;
  padding: 60rpx 40rpx 48rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.modal-icon-wrap {
  width: 140rpx;
  height: 140rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 30rpx;
}
.success-icon {
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.3), rgba(76, 175, 80, 0.1));
  box-shadow: 0 0 60rpx rgba(76, 175, 80, 0.3);
}
.fail-icon {
  background: linear-gradient(135deg, rgba(233, 69, 96, 0.3), rgba(233, 69, 96, 0.1));
  box-shadow: 0 0 60rpx rgba(233, 69, 96, 0.3);
}
.modal-icon {
  font-size: 72rpx;
}
.modal-title {
  font-size: 40rpx;
  color: #fff;
  font-weight: 800;
  margin-bottom: 12rpx;
}
.modal-desc {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
  margin-bottom: 40rpx;
}
.modal-actions {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  width: 100%;
}
.modal-btn {
  width: 100%;
  height: 92rpx;
  border-radius: 46rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30rpx;
  font-weight: 600;
  transition: all 0.2s;
}
.modal-btn.primary {
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  color: #fff;
  box-shadow: 0 8rpx 24rpx rgba(233, 69, 96, 0.4);
}
.modal-btn.ghost {
  background: rgba(255,255,255,0.06);
  border: 1rpx solid rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.7);
}
</style>
