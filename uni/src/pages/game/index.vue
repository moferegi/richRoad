<template>
  <view class="game-index">
    <!-- 背景装饰 -->
    <view class="bg-glow bg-glow-1"></view>
    <view class="bg-glow bg-glow-2"></view>

    <!-- 顶部栏 -->
    <view class="top-bar">
      <view class="brand">
        <view class="brand-logo">
          <text class="logo-text">24</text>
        </view>
        <text class="app-name">{{ appName }}</text>
      </view>
      <view class="top-actions">
        <view class="user-chip">
          <text class="user-avatar">{{ username.charAt(0).toUpperCase() }}</text>
          <text class="user-name">{{ username }}</text>
        </view>
        <view class="icon-btn" @click="playClick(); showSettings = true">
          <text class="settings-icon">⚙</text>
        </view>
      </view>
    </view>

    <!-- 游戏选择 -->
    <view class="hero-section">
      <text class="hero-label">{{ t('game.gameList') }}</text>
      <view class="game-grid">
        <view
          v-for="game in gameList"
          :key="game.ID"
          class="game-card"
          :class="{ disabled: game.status !== 1 }"
          @click="game.status === 1 ? (playClick(), goToCategory(game)) : null"
        >
          <view class="game-card-icon">
            <text class="game-icon-text">{{ game.gameKey === '24point' ? '24' : game.gameKey === '36point' ? '36' : '🔐' }}</text>
          </view>
          <text class="game-card-name">{{ localText(game.name) }}</text>
          <text v-if="game.status !== 1" class="game-card-badge">{{ t('game.comingSoon') }}</text>
          <text v-else class="game-card-enter">进入 →</text>
        </view>
      </view>
    </view>

    <!-- 排行榜 -->
    <view class="leaderboard-section">
      <view class="section-header">
        <text class="section-title">{{ t('game.leaderboardAll') }}</text>
        <view class="tab-switch">
          <text
            class="tab-btn"
            :class="{ active: leaderboardType === 'all' }"
            @click="playClick(); switchLeaderboard('all')"
          >{{ t('game.leaderboardAll') }}</text>
          <text
            class="tab-btn"
            :class="{ active: leaderboardType === 'daily' }"
            @click="playClick(); switchLeaderboard('daily')"
          >{{ t('game.leaderboardDaily') }}</text>
        </view>
      </view>
      <view class="leaderboard-card">
        <view v-if="leaderboard.length === 0" class="empty-state">
          <text class="empty-icon">🏆</text>
          <text class="empty-text">{{ t('game.noData') }}</text>
        </view>
        <view
          v-for="(item, idx) in leaderboard.slice(0, 10)"
          :key="item.userID"
          class="rank-row"
          :class="'rank-row-' + (idx + 1)"
        >
          <view class="rank-badge" :class="'rank-badge-' + (idx + 1)">
            <text v-if="idx < 3" class="rank-crown">{{ ['🥇','🥈','🥉'][idx] }}</text>
            <text v-else class="rank-num">{{ idx + 1 }}</text>
          </view>
          <view class="rank-user">
            <text class="rank-name">{{ item.nickname || item.username }}</text>
          </view>
          <view class="rank-score">
            <text class="score-num">{{ item.total }}</text>
            <text class="score-label">{{ t('game.passed') }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 设置弹窗 -->
    <view v-if="showSettings" class="modal-overlay" @click="playClick(); showSettings = false">
      <view class="modal-panel" @click.stop>
        <view class="modal-handle"></view>
        <text class="modal-title">{{ t('game.settings') }}</text>
        <view class="setting-list">
          <view class="setting-row">
            <view class="setting-left">
              <text class="setting-icon">🔊</text>
              <text class="setting-label">{{ t('game.sound') }}</text>
            </view>
            <switch :checked="soundEnabled" @change="toggleSound" color="#e94560" />
          </view>
          <view class="setting-row">
            <view class="setting-left">
              <text class="setting-icon">🎵</text>
              <text class="setting-label">{{ t('game.music') }}</text>
            </view>
            <switch :checked="musicEnabled" @change="toggleMusic" color="#e94560" />
          </view>
          <view class="setting-row">
            <view class="setting-left">
              <text class="setting-icon">👤</text>
              <text class="setting-label">{{ t('game.username') }}</text>
            </view>
            <text class="setting-value">{{ username }}</text>
          </view>
          <view class="setting-row" @click="playClick(); openLanguageSwitch()">
            <view class="setting-left">
              <text class="setting-icon">🌐</text>
              <text class="setting-label">{{ t('game.language') }}</text>
            </view>
            <view class="setting-right">
              <text class="setting-value">{{ locale === 'zh' ? '中文' : 'English' }}</text>
              <text class="chevron">›</text>
            </view>
          </view>
        </view>
        <view class="modal-footer">
          <view class="btn-primary" @click="playClick(); showSettings = false">
            <text>{{ t('common.confirm') }}</text>
          </view>
        </view>
      </view>
    </view>

    <lang-switch v-model="showLangPicker" />
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useUserStore } from '@/pinia/modules/user.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { getGameList, getLeaderboard } from '@/api/game.js'
import { autoRegister } from '@/api/base.js'
import { getAppName } from '@/api/sysConfig.js'

const langStore = useLangStore()
const userStore = useUserStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const t = (key) => {
  const text = i18nT(key, locale.value)
  return text !== key ? text : key
}
const localText = (value) => i18nLocalText(value, locale.value)

const appName = ref('')
const username = ref('')
const gameList = ref([])
const leaderboard = ref([])
const leaderboardType = ref('all')
const showSettings = ref(false)
const showLangPicker = ref(false)
const soundEnabled = ref(true)
const musicEnabled = ref(true)

const clickAudio = ref(null)
const bgMusic = ref(null)

const loadAppName = async () => {
  try {
    const res = await getAppName()
    if (res.code === 0 && res.data) {
      appName.value = localText(res.data)
    }
  } catch (e) { /* ignore */ }
}

const loadGameList = async () => {
  try {
    const res = await getGameList()
    if (res.code === 0) gameList.value = res.data || []
  } catch (e) { /* ignore */ }
}

const loadLeaderboard = async () => {
  try {
    const res = await getLeaderboard(leaderboardType.value, 50)
    if (res.code === 0) leaderboard.value = res.data || []
  } catch (e) { /* ignore */ }
}

const switchLeaderboard = (type) => {
  leaderboardType.value = type
  loadLeaderboard()
}

const goToCategory = (game) => {
  if (game.status !== 1) return
  uni.navigateTo({ url: `/pages/game/category?gameID=${game.ID}&gameKey=${game.gameKey}` })
}

const initAudio = () => {
  clickAudio.value = uni.createInnerAudioContext()
  clickAudio.value.src = '/static/click-click.mp3'
  bgMusic.value = uni.createInnerAudioContext()
  bgMusic.value.src = '/static/back-music.mp3'
  bgMusic.value.loop = true
  soundEnabled.value = uni.getStorageSync('game-sound-enabled') !== '0'
  musicEnabled.value = uni.getStorageSync('game-music-enabled') !== '0'
  if (musicEnabled.value) {
    bgMusic.value.play()
  }
}

const toggleSound = (e) => {
  soundEnabled.value = e.detail.value
  uni.setStorageSync('game-sound-enabled', soundEnabled.value ? '1' : '0')
}

const toggleMusic = (e) => {
  musicEnabled.value = e.detail.value
  uni.setStorageSync('game-music-enabled', musicEnabled.value ? '1' : '0')
  if (musicEnabled.value) {
    bgMusic.value?.play()
  } else {
    bgMusic.value?.stop()
  }
}

const playClick = () => {
  if (soundEnabled.value && clickAudio.value) {
    clickAudio.value.stop()
    clickAudio.value.play()
  }
}

const openLanguageSwitch = () => {
  showSettings.value = false
  showLangPicker.value = true
}

const initUser = async () => {
  if (userStore.token && userStore.userInfo) {
    username.value = userStore.userInfo.nickname || userStore.userInfo.username || 'User'
    return
  }
  try {
    const res = await autoRegister()
    if (res.code === 0 && res.data && res.data.user) {
      const user = res.data.user
      const token = res.data.token
      userStore.setToken(token)
      userStore.setUserInfo(user)
      username.value = user.nickname || user.username || 'User'
    } else {
      username.value = 'User'
    }
  } catch (e) {
    username.value = 'User'
  }
}

onMounted(() => {
  initAudio()
  initUser()
  loadAppName()
  loadGameList()
  loadLeaderboard()
})
</script>

<style scoped>
.game-index {
  min-height: 100vh;
  background: #0a0a1a;
  padding: 0 32rpx;
  padding-bottom: 60rpx;
  position: relative;
  overflow: hidden;
}

/* 背景光晕 */
.bg-glow {
  position: fixed;
  border-radius: 50%;
  filter: blur(120rpx);
  opacity: 0.4;
  pointer-events: none;
  z-index: 0;
}
.bg-glow-1 {
  width: 500rpx;
  height: 500rpx;
  background: #e94560;
  top: -150rpx;
  right: -100rpx;
}
.bg-glow-2 {
  width: 400rpx;
  height: 400rpx;
  background: #533483;
  bottom: 100rpx;
  left: -100rpx;
}

/* 顶部栏 */
.top-bar {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 80rpx 0 40rpx;
}
.brand {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.brand-logo {
  width: 72rpx;
  height: 72rpx;
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(233, 69, 96, 0.4);
}
.logo-text {
  font-size: 32rpx;
  font-weight: 900;
  color: #fff;
  font-style: italic;
}
.app-name {
  font-size: 40rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 1rpx;
}
.top-actions {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.user-chip {
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.1);
  padding: 10rpx 20rpx 10rpx 10rpx;
  border-radius: 40rpx;
  backdrop-filter: blur(10px);
}
.user-avatar {
  width: 48rpx;
  height: 48rpx;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 50%;
  color: #fff;
  font-size: 24rpx;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 48rpx;
  text-align: center;
}
.user-name {
  font-size: 24rpx;
  color: #fff;
  max-width: 120rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.icon-btn {
  width: 72rpx;
  height: 72rpx;
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  backdrop-filter: blur(10px);
}
.settings-icon {
  font-size: 32rpx;
  line-height: 1;
}

/* 游戏选择区 */
.hero-section {
  position: relative;
  z-index: 1;
  margin-top: 20rpx;
}
.hero-label {
  font-size: 28rpx;
  color: rgba(255,255,255,0.6);
  font-weight: 500;
  margin-bottom: 24rpx;
  display: block;
}
.game-grid {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}
.game-card {
  position: relative;
  border-radius: 24rpx;
  padding: 32rpx;
  background: linear-gradient(145deg, rgba(255,255,255,0.1), rgba(255,255,255,0.03));
  border: 1rpx solid rgba(255,255,255,0.1);
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  gap: 24rpx;
}
.game-card.disabled {
  opacity: 0.4;
}
.game-card-icon {
  width: 96rpx;
  height: 96rpx;
  min-width: 96rpx;
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(233, 69, 96, 0.4);
}
.game-icon-text {
  font-size: 40rpx;
  font-weight: 900;
  color: #fff;
  font-style: italic;
}
.game-card-name {
  flex: 1;
  font-size: 30rpx;
  font-weight: 700;
  color: #fff;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  word-break: break-all;
}
.game-card-badge {
  font-size: 22rpx;
  color: rgba(255,255,255,0.5);
  background: rgba(255,255,255,0.1);
  padding: 8rpx 20rpx;
  border-radius: 20rpx;
  flex-shrink: 0;
}
.game-card-enter {
  font-size: 24rpx;
  color: rgba(255,255,255,0.6);
  flex-shrink: 0;
}

/* 排行榜 */
.leaderboard-section {
  position: relative;
  z-index: 1;
  margin-top: 48rpx;
}
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}
.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
}
.tab-switch {
  display: flex;
  background: rgba(255,255,255,0.08);
  border-radius: 30rpx;
  padding: 6rpx;
  border: 1rpx solid rgba(255,255,255,0.08);
}
.tab-btn {
  padding: 10rpx 24rpx;
  font-size: 24rpx;
  color: rgba(255,255,255,0.6);
  border-radius: 24rpx;
  transition: all 0.3s;
}
.tab-btn.active {
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  color: #fff;
  font-weight: 600;
  box-shadow: 0 4rpx 12rpx rgba(233, 69, 96, 0.3);
}
.leaderboard-card {
  background: linear-gradient(145deg, rgba(255,255,255,0.08), rgba(255,255,255,0.03));
  border: 1rpx solid rgba(255,255,255,0.08);
  border-radius: 28rpx;
  padding: 20rpx 24rpx;
  backdrop-filter: blur(20px);
}
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 0;
}
.empty-icon {
  font-size: 80rpx;
  margin-bottom: 16rpx;
}
.empty-text {
  font-size: 26rpx;
  color: rgba(255,255,255,0.4);
}
.rank-row {
  display: flex;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid rgba(255,255,255,0.05);
}
.rank-row:last-child {
  border-bottom: none;
}
.rank-badge {
  width: 56rpx;
  height: 56rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  background: rgba(255,255,255,0.08);
}
.rank-badge-1 {
  background: linear-gradient(135deg, #ffd700, #ff8c00);
  box-shadow: 0 4rpx 12rpx rgba(255, 215, 0, 0.3);
}
.rank-badge-2 {
  background: linear-gradient(135deg, #e8e8e8, #a0a0a0);
  box-shadow: 0 4rpx 12rpx rgba(192, 192, 192, 0.3);
}
.rank-badge-3 {
  background: linear-gradient(135deg, #cd7f32, #8b4513);
  box-shadow: 0 4rpx 12rpx rgba(205, 127, 50, 0.3);
}
.rank-crown {
  font-size: 30rpx;
}
.rank-num {
  font-size: 28rpx;
  font-weight: 700;
  color: rgba(255,255,255,0.6);
}
.rank-user {
  flex: 1;
}
.rank-name {
  font-size: 28rpx;
  color: #fff;
  font-weight: 500;
}
.rank-score {
  text-align: right;
}
.score-num {
  font-size: 32rpx;
  font-weight: 800;
  color: #e94560;
  display: block;
  line-height: 1;
}
.score-label {
  font-size: 20rpx;
  color: rgba(255,255,255,0.4);
  margin-top: 4rpx;
  display: block;
}

/* 弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-end;
  z-index: 999;
}
.modal-panel {
  width: 100%;
  background: linear-gradient(180deg, #1a1a2e 0%, #0f0f1e 100%);
  border-radius: 40rpx 40rpx 0 0;
  padding: 20rpx 32rpx 60rpx;
  border-top: 1rpx solid rgba(255,255,255,0.1);
}
.modal-handle {
  width: 80rpx;
  height: 8rpx;
  background: rgba(255,255,255,0.2);
  border-radius: 4rpx;
  margin: 0 auto 30rpx;
}
.modal-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
  text-align: center;
  margin-bottom: 32rpx;
}
.setting-list {
  background: rgba(255,255,255,0.05);
  border-radius: 24rpx;
  border: 1rpx solid rgba(255,255,255,0.06);
  overflow: hidden;
}
.setting-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28rpx 24rpx;
  border-bottom: 1rpx solid rgba(255,255,255,0.05);
}
.setting-row:last-child {
  border-bottom: none;
}
.setting-left {
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.setting-icon {
  font-size: 36rpx;
}
.setting-label {
  font-size: 28rpx;
  color: #fff;
}
.setting-right {
  display: flex;
  align-items: center;
  gap: 8rpx;
}
.setting-value {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
}
.chevron {
  font-size: 36rpx;
  color: rgba(255,255,255,0.3);
  font-weight: 300;
}
.modal-footer {
  margin-top: 40rpx;
}
.btn-primary {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #e94560, #ff6b6b);
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 30rpx;
  font-weight: 600;
  box-shadow: 0 8rpx 24rpx rgba(233, 69, 96, 0.4);
}
</style>
