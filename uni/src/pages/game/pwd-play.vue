<template>
  <view class="pwd-game">
    <!-- 背景装饰 -->
    <view class="bg-glow bg-glow-1"></view>
    <view class="bg-glow bg-glow-2"></view>

    <!-- 顶部栏 -->
    <view class="top-bar">
      <view class="back-btn" @click="playClick(); goBack()">
        <text class="back-icon">&#8249;</text>
      </view>
      <view class="title-wrap">
        <text class="title">{{ difficultyName }}</text>
        <text class="subtitle">{{ t('game.level') }} {{ currentLevelNumber }}</text>
      </view>
      <view class="placeholder" />
    </view>

    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">{{ t('game.loading') }}</text>
    </view>

    <template v-else>
      <!-- ========== 输入答案显示（最上方） ========== -->
      <view class="input-section">
        <text class="section-title">{{ t('pwdGame.inputAnswer') }}</text>
        <view class="answer-boxes">
          <view
            v-for="(d, di) in answerDigits"
            :key="di"
            class="answer-box"
            :class="{ filled: d !== '' }"
          >
            <text class="answer-digit">{{ d }}</text>
          </view>
        </view>
      </view>

      <!-- ========== 4个提示卡片：每行4个数字 + 文字一一对应 ========== -->
      <view class="hints-section">
        <text class="section-title">{{ t('pwdGame.hints') }}</text>
        <view class="hint-cards">
          <view v-for="card in 4" :key="card" class="hint-card">
            <view class="hint-digits-row">
              <view
                v-for="pos in 4"
                :key="pos"
                class="hint-digit-box"
                :style="hintColors[(card-1)*4 + (pos-1)] ? { backgroundColor: hintColors[(card-1)*4 + (pos-1)], borderColor: hintColors[(card-1)*4 + (pos-1)] } : {}"
                @click="onHintDigitClick(card-1, pos-1)"
              >
                <text class="hint-digit-text">{{ hintDigits[(card-1)*4 + (pos-1)] }}</text>
              </view>
            </view>
            <view class="hint-text-box">
              <text class="hint-text">{{ localText(hintTexts[card-1]) }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- ========== 颜色标记（hint-cards 下方，一行） ========== -->
      <view class="color-section">
        <text class="color-hint">{{ t('pwdGame.colorHint') }}</text>
        <view class="color-row">
          <view
            v-for="(c, ci) in colors"
            :key="ci"
            class="color-block"
            :class="{ active: selectedColorIdx === ci }"
            :style="{ backgroundColor: c.color }"
            @click="onColorClick(ci)"
          >
            <text v-if="selectedColorIdx === ci" class="color-check">&#10003;</text>
          </view>
        </view>
        <view class="reset-colors-btn" @click="resetColors">
          <text class="reset-colors-text">{{ t('pwdGame.reset') }}</text>
        </view>
      </view>
    </template>

    <!-- ========== 数字键盘（固定底部，两行） ========== -->
    <view class="keypad">
      <!-- 第一行 1-5 + 删除 -->
      <view class="keypad-row">
        <view v-for="n in 5" :key="'r1' + n" class="keypad-btn" @click="onKeypadClick(n)">
          <text class="keypad-text">{{ n }}</text>
        </view>
        <view class="keypad-btn keypad-delete" @click="onDelete">
          <text class="keypad-text">{{ t('pwdGame.delete') }}</text>
        </view>
      </view>
      <!-- 第二行 6-0 + 确认 -->
      <view class="keypad-row">
        <view v-for="n in 5" :key="'r2' + (n + 5)" class="keypad-btn" @click="onKeypadClick(n === 5 ? 0 : n + 5)">
          <text class="keypad-text">{{ n === 5 ? 0 : n + 5 }}</text>
        </view>
        <view class="keypad-btn keypad-confirm" @click="onConfirm">
          <text class="keypad-text">&#10003;</text>
        </view>
      </view>
    </view>

    <!-- 成功弹窗 -->
    <view v-if="showSuccess" class="modal-overlay">
      <view class="modal-card success">
        <view class="modal-icon-wrap success-icon">
          <text class="modal-icon">&#127881;</text>
        </view>
        <text class="modal-title">{{ t('pwdGame.correct') }}</text>
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
          <text class="modal-icon">&#128549;</text>
        </view>
        <text class="modal-title">{{ t('pwdGame.wrong') }}</text>
        <view class="modal-actions">
          <view class="modal-btn primary" @click="playClick(); retryAnswer()">
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
import { t, localText } from '@/utils/i18n.js'
import { getPwdLevelDetail, getDifficultyCategory, getDifficultyCategories, getPwdLevelList, submitPwdLevelResult } from '@/api/game.js'

const levelID = ref(0)
const catID = ref(0)
const loading = ref(true)
const showSuccess = ref(false)
const showFail = ref(false)
const submitted = ref(false)
const difficultyName = ref('')
const currentLevelNumber = ref(0)
const correctAnswer = ref('')

// 提示数字（4组×4位=16位，按 cardIdx*4+digitIdx 索引）
const hintDigits = reactive(Array(16).fill(''))
const hintTexts = reactive(['', '', '', ''])
// 每个卡片每个数字的颜色状态（4卡片×4数字=16位，索引 = cardIdx*4 + digitIdx）
const hintColors = reactive(Array(16).fill(null))
const colors = [
  { color: '#e94560' },
  { color: '#f0a500' },
  { color: '#00b894' }
]
const selectedColorIdx = ref(-1)

// 输入答案
const answerDigits = ref(['', '', '', ''])
const answerIndex = ref(0)

// 音效
const clickAudio = ref(null)

const playClick = () => {
  const soundEnabled = uni.getStorageSync('game-sound-enabled') !== '0'
  if (!soundEnabled || !clickAudio.value) return
  clickAudio.value.stop()
  clickAudio.value.play()
}

const initAudio = () => {
  const soundEnabled = uni.getStorageSync('game-sound-enabled') !== '0'
  if (soundEnabled) {
    clickAudio.value = uni.createInnerAudioContext()
    clickAudio.value.src = '/static/click-click.mp3'
  }
}

const goBack = () => {
  uni.navigateBack()
}

const goBackToHome = () => {
  uni.reLaunch({ url: '/pages/game/index' })
}

// 解析 hintTexts：兼容 i18n 裁剪后的数组 和 原始 JSON 字符串
const parseHintTexts = (raw) => {
  if (!raw) return []
  if (Array.isArray(raw)) return raw
  if (typeof raw === 'string') {
    try {
      const parsed = JSON.parse(raw)
      return Array.isArray(parsed) ? parsed : []
    } catch { return [] }
  }
  return []
}

const loadLevel = async () => {
  try {
    const res = await getPwdLevelDetail(levelID.value)
    if (res.code === 0 && res.data) {
      currentLevelNumber.value = res.data.levelNumber || 0
      correctAnswer.value = res.data.answer || ''

      const rawDigits = (res.data.hintDigits || '').trim()
      let groups = []
      if (rawDigits.includes(',')) {
        groups = rawDigits.split(',').map(d => d.trim()).filter(Boolean)
      } else {
        groups = [rawDigits]
      }
      const texts = parseHintTexts(res.data.hintTexts)

      // 展平为 16 元素数组：每组 4 个数字，不足补第一组
      const flat = Array(16).fill('')
      for (let g = 0; g < 4; g++) {
        const group = groups[g] || groups[0] || ''
        const chars = group.split('')
        for (let c = 0; c < 4; c++) {
          flat[g * 4 + c] = chars[c] || ''
        }
      }
      hintDigits.splice(0, 16, ...flat)

      const t = ['', '', '', '']
      for (let i = 0; i < 4; i++) {
        t[i] = texts[i] || ''
      }
      hintTexts.splice(0, 4, t[0], t[1], t[2], t[3])
    }
  } catch (e) { /* ignore */ }
  loading.value = false
}

const loadDifficultyName = async () => {
  try {
    const res = await getDifficultyCategory(catID.value)
    if (res.code === 0 && res.data) {
      difficultyName.value = localText(res.data.name) || ''
    }
  } catch (e) { /* ignore */ }
}

// 颜色标记点击
const onColorClick = (idx) => {
  playClick()
  selectedColorIdx.value = selectedColorIdx.value === idx ? -1 : idx
}

// 提示数字点击：同色取消，异色直接切换（cardIdx: 0-3, digitIdx: 0-3）
const onHintDigitClick = (cardIdx, digitIdx) => {
  playClick()
  if (selectedColorIdx.value === -1) return
  const idx = cardIdx * 4 + digitIdx
  const newColor = colors[selectedColorIdx.value].color
  hintColors[idx] = hintColors[idx] === newColor ? null : newColor
}

// 数字键盘
const onKeypadClick = (num) => {
  playClick()
  if (answerIndex.value >= 4) return
  answerDigits.value[answerIndex.value] = String(num)
  answerIndex.value++
}

// 删除
const onDelete = () => {
  playClick()
  if (answerIndex.value === 0) return
  answerIndex.value--
  answerDigits.value[answerIndex.value] = ''
}

// 确认提交
const onConfirm = () => {
  playClick()
  if (answerIndex.value < 4) return
  const input = answerDigits.value.join('')
  if (input === correctAnswer.value) {
    showSuccess.value = true
    submitResult()
  } else {
    showFail.value = true
  }
}

const submitResult = async () => {
  if (submitted.value) return
  submitted.value = true
  try {
    await submitPwdLevelResult(levelID.value)
  } catch (e) { /* ignore */ }
}

const resetGame = () => {
  showSuccess.value = false
  showFail.value = false
  submitted.value = false
  answerDigits.value = ['', '', '', '']
  answerIndex.value = 0
  selectedColorIdx.value = -1
  hintColors.splice(0, 16, ...Array(16).fill(null))
}

// 仅重置答案输入（失败后重试，保留颜色标记）
const retryAnswer = () => {
  showFail.value = false
  answerDigits.value = ['', '', '', '']
  answerIndex.value = 0
}

// 重置所有提示数字颜色
const resetColors = () => {
  playClick()
  hintColors.splice(0, 16, ...Array(16).fill(null))
}

const goNextLevel = async () => {
  showSuccess.value = false
  try {
    const res = await getPwdLevelList(catID.value)
    if (res.code === 0) {
      const levels = (res.data || []).sort((a, b) => a.levelNumber - b.levelNumber)
      const curIdx = levels.findIndex(l => l.ID === levelID.value)
      if (curIdx >= 0 && curIdx < levels.length - 1) {
        const nextLevel = levels[curIdx + 1]
        uni.redirectTo({ url: `/pages/game/pwd-play?levelID=${nextLevel.ID}&catID=${catID.value}` })
        return
      }
    }
    const catRes = await getDifficultyCategory(catID.value)
    if (catRes.code === 0 && catRes.data) {
      const gameID = catRes.data.gameID
      const catsRes = await getDifficultyCategories(gameID)
      if (catsRes.code === 0) {
        const cats = (catsRes.data || []).sort((a, b) => a.sort - b.sort)
        const curCatIdx = cats.findIndex(c => c.ID === catID.value)
        if (curCatIdx >= 0 && curCatIdx < cats.length - 1) {
          const nextCat = cats[curCatIdx + 1]
          const nextLevelsRes = await getPwdLevelList(nextCat.ID)
          if (nextLevelsRes.code === 0) {
            const nextLevels = (nextLevelsRes.data || []).sort((a, b) => a.levelNumber - b.levelNumber)
            if (nextLevels.length > 0) {
              uni.redirectTo({ url: `/pages/game/pwd-play?levelID=${nextLevels[0].ID}&catID=${nextCat.ID}` })
              return
            }
          }
        }
      }
    }
  } catch (e) { /* ignore */ }
  uni.navigateBack()
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
.pwd-game {
  min-height: 100vh;
  background: #0a0a1a;
  padding: 0 32rpx 180rpx;
  position: relative;
  overflow: hidden;
}

.bg-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(120rpx);
  opacity: 0.15;
  pointer-events: none;
}
.bg-glow-1 { width: 500rpx; height: 500rpx; background: #e94560; top: -100rpx; right: -100rpx; }
.bg-glow-2 { width: 400rpx; height: 400rpx; background: #6c5ce7; bottom: 200rpx; left: -80rpx; }

/* 顶部栏 */
.top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 0 20rpx;
  position: relative;
  z-index: 1;
}
.back-btn {
  width: 64rpx; height: 64rpx;
  border-radius: 50%;
  background: rgba(255,255,255,0.08);
  display: flex;
  align-items: center;
  justify-content: center;
}
.back-icon { font-size: 48rpx; color: #fff; font-weight: 300; line-height: 1; }
.title-wrap { text-align: center; }
.title { font-size: 32rpx; font-weight: 700; color: #fff; }
.subtitle { font-size: 22rpx; color: rgba(255,255,255,0.5); margin-top: 4rpx; display: block; }
.placeholder { width: 64rpx; }

.loading-wrap { display: flex; justify-content: center; padding: 120rpx 0; }
.loading-text { color: rgba(255,255,255,0.5); font-size: 28rpx; }

.section-title {
  font-size: 22rpx;
  color: rgba(255,255,255,0.4);
  margin-bottom: 14rpx;
  display: block;
  text-transform: uppercase;
  letter-spacing: 2rpx;
}

/* ========== 输入答案显示（最上方） ========== */
.input-section {
  position: relative;
  z-index: 1;
  margin-top: 12rpx;
}
.answer-boxes {
  display: flex;
  gap: 20rpx;
  justify-content: center;
}
.answer-box {
  width: 100rpx;
  height: 110rpx;
  border-radius: 16rpx;
  background: rgba(255,255,255,0.05);
  border: 2rpx solid rgba(255,255,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.answer-box.filled {
  border-color: #e94560;
  background: rgba(233,69,96,0.12);
  box-shadow: 0 0 16rpx rgba(233,69,96,0.2);
}
.answer-digit {
  font-size: 48rpx;
  font-weight: 700;
  color: #fff;
  font-style: italic;
}

/* ========== 提示卡片 ========== */
.hints-section {
  position: relative;
  z-index: 1;
  margin-top: 24rpx;
}
.hint-cards {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}
.hint-card {
  display: flex;
  flex-direction: column;
  gap: 14rpx;
  background: rgba(255,255,255,0.04);
  border: 1rpx solid rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 16rpx 20rpx;
}
.hint-digits-row {
  display: flex;
  gap: 16rpx;
  justify-content: center;
}
.hint-digit-box {
  width: 100rpx;
  height: 80rpx;
  border-radius: 14rpx;
  background: rgba(255,255,255,0.08);
  border: 2rpx solid rgba(255,255,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.hint-digit-text {
  font-size: 40rpx;
  font-weight: 800;
  color: #fff;
  font-style: italic;
}
.hint-text-box {
  flex: 1;
  display: flex;
  align-items: flex-start;
  gap: 12rpx;
}
.hint-text {
  font-size: 26rpx;
  color: rgba(255,255,255,0.75);
  line-height: 1.5;
  flex: 1;
}

/* ========== 颜色标记（hint-cards 下方，一行） ========== */
.color-section {
  position: relative;
  z-index: 1;
  margin-top: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16rpx;
}
.color-hint {
  font-size: 22rpx;
  color: rgba(255,255,255,0.4);
  flex-shrink: 0;
}
.color-row {
  display: flex;
  gap: 20rpx;
}
.color-block {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  opacity: 0.45;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3rpx solid transparent;
}
.color-block.active {
  opacity: 1;
  transform: scale(1.2);
  border-color: #fff;
  box-shadow: 0 0 20rpx currentColor;
}
.color-check {
  font-size: 30rpx;
  color: #fff;
  font-weight: 700;
}
.reset-colors-btn {
  padding: 10rpx 24rpx;
  border-radius: 24rpx;
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.12);
  flex-shrink: 0;
}
.reset-colors-btn:active {
  background: rgba(255,255,255,0.15);
}
.reset-colors-text {
  font-size: 22rpx;
  color: rgba(255,255,255,0.55);
}

/* ========== 数字键盘（固定底部，两行） ========== */
.keypad {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 10;
  background: linear-gradient(to top, rgba(10,10,26,0.98), rgba(10,10,26,0.9));
  backdrop-filter: blur(20rpx);
  padding: 16rpx 32rpx 32rpx;
  display: flex;
  flex-direction: column;
  gap: 14rpx;
}
.keypad-row {
  display: flex;
  gap: 14rpx;
}
.keypad-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 16rpx;
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}
.keypad-btn:active {
  background: rgba(233,69,96,0.3);
  transform: scale(0.95);
}
.keypad-text {
  font-size: 36rpx;
  font-weight: 600;
  color: #fff;
}
.keypad-delete {
  background: rgba(255,255,255,0.05);
  flex: 1.2;
}
.keypad-delete .keypad-text {
  font-size: 24rpx;
  color: rgba(255,255,255,0.55);
}
.keypad-confirm {
  background: linear-gradient(135deg, #00b894, #00cec9);
  border: none;
  flex: 1.2;
  box-shadow: 0 4rpx 16rpx rgba(0,184,148,0.3);
}
.keypad-confirm:active {
  background: linear-gradient(135deg, #00a381, #00b894);
}
.keypad-confirm .keypad-text {
  font-size: 42rpx;
  font-weight: 700;
}

/* ========== 弹窗 ========== */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8rpx);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal-card {
  width: 560rpx;
  background: rgba(30,30,60,0.96);
  border: 1rpx solid rgba(255,255,255,0.1);
  border-radius: 32rpx;
  padding: 48rpx 40rpx 36rpx;
  text-align: center;
  backdrop-filter: blur(20rpx);
}
.modal-icon-wrap {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  margin: 0 auto 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.success-icon { background: rgba(0,184,148,0.15); }
.fail-icon { background: rgba(233,69,96,0.15); }
.modal-icon { font-size: 52rpx; }
.modal-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12rpx;
}
.modal-desc {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
  margin-bottom: 32rpx;
}
.modal-actions { display: flex; flex-direction: column; gap: 16rpx; }
.modal-btn {
  padding: 22rpx;
  border-radius: 16rpx;
  font-size: 28rpx;
  font-weight: 600;
}
.modal-btn.primary {
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  color: #fff;
}
.modal-btn.ghost {
  background: rgba(255,255,255,0.08);
  color: rgba(255,255,255,0.7);
}
</style>