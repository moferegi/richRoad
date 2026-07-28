<template>
  <view class="profile-container">
    <!-- 渐变 Hero 头部 -->
    <view class="profile-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-user">
        <view class="avatar-ring">
          <image class="avatar" :src="avatarSrc" mode="aspectFill"></image>
        </view>
        <text class="nickname">{{ userInfo.nickname }}</text>
        <view class="hero-badge">
          <text class="badge-icon">★</text>
          <text class="badge-text">{{ t('profile.points') }} {{ userAsset.totalPoints }}</text>
        </view>
      </view>
    </view>

    <!-- 资产网格卡片 -->
    <view class="asset-section">
      <view class="asset-card" @click="goTo('point_history')">
        <text class="asset-val">{{ userAsset.totalPoints }}</text>
        <text class="asset-label">{{ t('profile.points') }}</text>
      </view>
      <view class="asset-card" @click="goTo('free_time_history')">
        <text class="asset-val">{{ userAsset.freeMinutes }}</text>
        <text class="asset-label">{{ t('profile.free_time') }}</text>
      </view>
    </view>

    <!-- 兑换按钮卡片 -->
    <view class="exchange-card" @click="showExchangePopup">
      <view class="exchange-left">
        <text class="exchange-title">{{ t('profile.exchange_time') }}</text>
      </view>
      <view class="exchange-arrow">
        <text class="arrow-icon">›</text>
      </view>
    </view>

    <!-- 列表菜单分组 -->
    <view class="menu-group">
      <view class="menu-list">
        <view class="menu-item" @click="goTo('watch_history')">
          <view class="menu-icon-wrap icon-video">
            <image class="menu-icon-img" src="/static/images/learning/icon-play-purple.svg" mode="aspectFit" />
          </view>
          <text class="menu-label">{{ t('profile.watch_history') }}</text>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="goTo('checkin_calendar')">
          <view class="menu-icon-wrap icon-calendar">
            <image class="menu-icon-img" src="/static/images/learning/icon-calendar.svg" mode="aspectFit" />
          </view>
          <text class="menu-label">{{ t('profile.checkin_record') }}</text>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="goTo('collections')">
          <view class="menu-icon-wrap icon-star">
            <image class="menu-icon-img" src="/static/images/learning/icon-bookmark-gold.svg" mode="aspectFit" />
          </view>
          <text class="menu-label">{{ t('profile.my_collections') }}</text>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="goTo('error_words')">
          <view class="menu-icon-wrap icon-book">
            <image class="menu-icon-img" src="/static/images/learning/icon-book.svg" mode="aspectFit" />
          </view>
          <text class="menu-label">{{ t('profile.error_book') }}</text>
          <text class="arrow">›</text>
        </view>
      </view>
    </view>

    <view class="menu-group">
      <view class="menu-list">
        <view class="menu-item" @click="contactService">
          <view class="menu-icon-wrap icon-service">
            <image class="menu-icon-img" src="/static/images/learning/icon-headset.svg" mode="aspectFit" />
          </view>
          <text class="menu-label">{{ t('profile.contact_service') }}</text>
          <text class="arrow">›</text>
        </view>
        <view class="menu-item" @click="showLanguageSwitcher">
          <view class="menu-icon-wrap icon-lang">
            <image class="menu-icon-img" src="/static/images/learning/icon-language.svg" mode="aspectFit" />
          </view>
          <text class="menu-label">{{ t('profile.switch_language') }}</text>
          <text class="arrow">›</text>
        </view>
      </view>
    </view>

    <button class="logout-btn" @click="logout">{{ t('profile.logout') }}</button>

    <lang-switch v-model="showLangPicker" />

    <!-- 积分兑换弹窗 -->
    <uni-popup ref="exchangePopup" type="center">
      <view class="popup-box">
        <view class="popup-handle"></view>
        <view class="popup-title">{{ t('profile.exchange_popup_title') }}</view>
        <view class="popup-content">
          <view class="popup-points-row">
            <text class="popup-points-label">{{ t('profile.current_points') }}</text>
            <text class="popup-points-val">{{ userAsset.totalPoints }}</text>
          </view>
          <input type="number" v-model="exchangePoints" :placeholder="t('profile.input_points_to_exchange')" class="exchange-input"/>
          <view class="exchange-tips-box">
            <text class="exchange-tips">{{ exchangePoints }} {{ t('profile.points_equals') }} {{ exchangeMinutes }} {{ t('profile.minutes') }}</text>
          </view>
        </view>
        <view class="popup-actions">
          <button class="popup-btn-cancel" @click="closeExchangePopup">{{ t('common.cancel') }}</button>
          <button class="popup-btn-confirm" @click="submitExchange">{{ t('common.confirm') }}</button>
        </view>
      </view>
    </uni-popup>
    <custom-tab-bar />
  </view>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { onShow, onHide } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText as i18nLocalText, t as i18nT } from '@/utils/i18n.js'
import { exchangeTime, getAsset } from '@/api/learning.js'
import { getPointsExchangeRate } from '@/api/sysConfig.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { getExternalUrl } from '@/utils/url.js'
import langSwitch from '@/components/lang-switch/lang-switch.vue'
import CustomTabBar from '@/components/custom-tab-bar/custom-tab-bar.vue'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
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

const userInfo = ref({ avatar: '', nickname: '' })
const userAsset = ref({ totalPoints: 0, freeMinutes: 0 })
const exchangeRate = ref(100)
const exchangePoints = ref(100)

const exchangePopup = ref(null)
const showLangPicker = ref(false)
const localeRefreshing = ref(false)
const avatarSrc = computed(() => {
  const logo = getExternalUrl(appConfigStore.appLogo || '')
  if (logo) {
    return logo
  }
  return userInfo.value.avatar || ''
})

const exchangeMinutes = computed(() => Math.floor(Number(exchangePoints.value || 0) / exchangeRate.value))

const loadUserInfo = () => {
  const raw = uni.getStorageSync('userInfo')
  if (raw && typeof raw === 'object') {
    const localizedNickname = localText(raw.nickName || raw.nickname)
    userInfo.value = {
      avatar: raw.headerImg || raw.avatar || '',
      nickname: String(localizedNickname || raw.nickName || raw.nickname || t('profile.default_nickname'))
    }
    return
  }
  userInfo.value = { avatar: '', nickname: t('profile.default_nickname') }
}

const loadAsset = async () => {
  const res = await getAsset()
  if (res.code !== 0 || !res.data) {
    return
  }
  userAsset.value = {
    totalPoints: Number(res.data.totalPoints || 0),
    freeMinutes: Number(res.data.freeMinutes || 0)
  }
}

const loadExchangeRate = async () => {
  const res = await getPointsExchangeRate()
  const rate = Number(res?.data?.configValue || 100)
  exchangeRate.value = rate > 0 ? rate : 100
}

const goTo = async (route) => {
  if (route === 'watch_history') {
    uni.navigateTo({ url: '/pages/learning/watch-history' })
    return
  }
  if (route === 'checkin_calendar') {
    uni.navigateTo({ url: '/pages/learning/checkin-record' })
    return
  }
  if (route === 'collections') {
    uni.navigateTo({ url: '/pages/learning/collections' })
    return
  }
  if (route === 'error_words') {
    uni.navigateTo({ url: '/pages/learning/error-log' })
    return
  }
  if (route === 'point_history') {
    uni.navigateTo({ url: '/pages/learning/point-history' })
    return
  }
  if (route === 'free_time_history') {
    uni.navigateTo({ url: '/pages/learning/free-time-history' })
    return
  }
  uni.showToast({ title: t('profile.coming_soon'), icon: 'none' })
}

const showExchangePopup = () => {
  exchangePoints.value = exchangeRate.value
  exchangePopup.value.open()
}

const closeExchangePopup = () => {
  exchangePopup.value?.close()
}

const submitExchange = () => {
  const points = Number(exchangePoints.value)
  if (!points || points <= 0) {
    uni.showToast({ title: t('profile.exchange_invalid'), icon: 'none' })
    return
  }
  exchangeTime(points).then((res) => {
    if (res.code === 0) {
      uni.showToast({ title: t('profile.exchange_success'), icon: 'success' })
      exchangePopup.value?.close()
      loadAsset()
    }
  })
}

const contactService = () => {
  uni.navigateTo({ url: '/pages/kefu/index?from=learning_profile' })
}

const showLanguageSwitcher = async () => {
  showLangPicker.value = true
}

const logout = () => {
  uni.showModal({
    title: t('logout'),
    content: t('confirmLogout'),
    cancelText: t('cancel'),
    confirmText: t('confirm'),
    success: (res) => {
      if (!res.confirm) {
        return
      }
      uni.removeStorageSync('x-token')
      uni.removeStorageSync('userInfo')
      uni.showToast({ title: t('logoutSuccess'), icon: 'none' })
      uni.reLaunch({ url: '/pages/user/login' })
    }
  })
}

onMounted(() => {
  appConfigStore.loadConfig({ force: true, localeOnly: true })
  loadUserInfo()
  loadAsset()
  loadExchangeRate()
})

onShow(() => {
  uni.hideTabBar()
  appConfigStore.loadConfig({ force: true, localeOnly: true })
  loadUserInfo()
  loadAsset()
  loadExchangeRate()
})

onHide(() => {
  // 切换页面时关闭语言选择弹窗，防止遮罩阻塞其他页面滚动
  showLangPicker.value = false
})

const forceRefreshByLocale = () => {
  if (localeRefreshing.value) return
  localeRefreshing.value = true
  uni.reLaunch({ url: '/pages/learning/profile' })
}

onMounted(() => {
  if (typeof uni.$on === 'function') {
    uni.$on('app:locale-force-refresh', forceRefreshByLocale)
  }
})

onUnmounted(() => {
  if (typeof uni.$off === 'function') {
    uni.$off('app:locale-force-refresh', forceRefreshByLocale)
  }
})
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  padding: 0 26rpx calc(96rpx + env(safe-area-inset-bottom));
  box-sizing: border-box;
  background: linear-gradient(180deg, #F5F3FF 0%, #EDE9FE 100%);
}

/* === 渐变 Hero 头部 === */
.profile-hero {
  position: relative;
  overflow: hidden;
  margin: 0 -26rpx 24rpx;
  padding: calc(var(--status-bar-height, 0px) + 48rpx) 32rpx 56rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-radius: 0 0 36rpx 36rpx;
}

.hero-decor-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.12);
}

.decor-1 {
  width: 280rpx;
  height: 280rpx;
  top: -100rpx;
  right: -60rpx;
}

.decor-2 {
  width: 180rpx;
  height: 180rpx;
  bottom: -80rpx;
  left: 160rpx;
}

.hero-user {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14rpx;
}

.avatar-ring {
  width: 140rpx;
  height: 140rpx;
  border-radius: 50%;
  padding: 6rpx;
  background: rgba(255, 255, 255, 0.4);
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.18);
}

.avatar {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: #EDE9FE;
}

.nickname {
  font-size: 38rpx;
  font-weight: 800;
  color: #fff;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.16);
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 8rpx 24rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.22);
  backdrop-filter: blur(8rpx);
}

.badge-icon {
  font-size: 24rpx;
  color: #FFD66B;
}

.badge-text {
  font-size: 24rpx;
  color: #fff;
  font-weight: 600;
}

/* === 资产网格卡片 === */
.asset-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.asset-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32rpx 20rpx;
  border-radius: 24rpx;
  background: #FFFFFF;
  border: 1rpx solid rgba(108, 91, 255, 0.12);
  box-shadow: 0 8rpx 32rpx rgba(108, 91, 255, 0.1);
}

.asset-card:active {
  transform: scale(0.97);
  background: rgba(108, 91, 255, 0.04);
}

.asset-val {
  font-size: 40rpx;
  font-weight: 800;
  color: #6D5BFF;
}

.asset-label {
  font-size: 24rpx;
  color: #6B6F8D;
  margin-top: 8rpx;
}

/* === 兑换按钮卡片 === */
.exchange-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-radius: 24rpx;
  background: linear-gradient(135deg, rgba(108, 91, 255, 0.1) 0%, rgba(155, 143, 255, 0.1) 100%);
  border: 1rpx solid rgba(108, 91, 255, 0.2);
  box-shadow: 0 8rpx 32rpx rgba(108, 91, 255, 0.08);
  margin-bottom: 28rpx;
}

.exchange-card:active {
  transform: scale(0.98);
}

.exchange-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #6D5BFF;
}

.exchange-arrow {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(108, 91, 255, 0.36);
}

.arrow-icon {
  color: #fff;
  font-size: 32rpx;
  font-weight: 700;
}

/* === 菜单分组 === */
.menu-group {
  margin-bottom: 24rpx;
}

.menu-group-title {
  font-size: 24rpx;
  color: #6B6F8D;
  font-weight: 600;
  padding: 0 12rpx 12rpx;
  letter-spacing: 1rpx;
}

.menu-list {
  background: #FFFFFF;
  border-radius: 24rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.12);
  overflow: hidden;
  box-shadow: 0 8rpx 32rpx rgba(108, 91, 255, 0.1);
  box-sizing: border-box;
  width: 100%;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 26rpx 28rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.08);
  box-sizing: border-box;
  width: 100%;
}

.menu-item:active {
  background: rgba(108, 91, 255, 0.06);
}

.menu-item:last-child { border-bottom: none; }

.menu-icon-wrap {
  width: 64rpx;
  height: 64rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.icon-video { background: linear-gradient(135deg, rgba(108, 91, 255, 0.18), rgba(155, 143, 255, 0.18)); }
.icon-calendar { background: linear-gradient(135deg, rgba(34, 197, 94, 0.18), rgba(74, 222, 128, 0.18)); }
.icon-star { background: linear-gradient(135deg, rgba(255, 214, 107, 0.22), rgba(251, 191, 36, 0.22)); }
.icon-book { background: linear-gradient(135deg, rgba(239, 68, 68, 0.16), rgba(248, 113, 113, 0.16)); }
.icon-service { background: linear-gradient(135deg, rgba(14, 165, 233, 0.18), rgba(56, 189, 248, 0.18)); }
.icon-lang { background: linear-gradient(135deg, rgba(168, 85, 247, 0.18), rgba(192, 132, 252, 0.18)); }

.menu-icon {
  font-size: 30rpx;
  font-weight: 700;
}

.icon-video .menu-icon { color: #6D5BFF; }
.icon-calendar .menu-icon { color: #22C55E; }
.icon-star .menu-icon { color: #F59E0B; }
.icon-book .menu-icon { color: #EF4444; }
.icon-service .menu-icon { color: #0EA5E9; }
.icon-lang .menu-icon { color: #A855F7; }

.menu-label {
  flex: 1;
  font-size: 30rpx;
  color: #1A1B3A;
}

.arrow { color: #9B8FFF; font-weight: 700; font-size: 36rpx; }

.logout-btn {
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 999rpx;
  font-size: 30rpx;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(135deg, #EF4444 0%, #F87171 100%);
  box-shadow: 0 8rpx 24rpx rgba(239, 68, 68, 0.28);
  margin-top: 16rpx;
  margin-bottom: calc(32rpx + env(safe-area-inset-bottom));
}

.logout-btn:active {
  opacity: 0.88;
  transform: scale(0.97);
}

/* === 弹窗 === */
.popup-box {
  width: 620rpx;
  background: #FFFFFF;
  border-radius: 32rpx;
  padding: 0 34rpx 34rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.12);
  box-shadow: 0 24rpx 64rpx rgba(108, 91, 255, 0.22);
  overflow: hidden;
}

.popup-handle {
  width: 60rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #6D5BFF, #9B8FFF);
  margin: 24rpx auto 0;
}

.popup-title {
  font-size: 34rpx;
  font-weight: 800;
  text-align: center;
  color: #1A1B3A;
  margin: 20rpx 0 28rpx;
}

.popup-content {
  margin-bottom: 28rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.popup-points-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 24rpx;
  border-radius: 16rpx;
  background: rgba(108, 91, 255, 0.06);
}

.popup-points-label {
  font-size: 26rpx;
  color: #6B6F8D;
}

.popup-points-val {
  font-size: 34rpx;
  font-weight: 800;
  color: #6D5BFF;
}

.exchange-input {
  border: 2rpx solid rgba(108, 91, 255, 0.24);
  background: rgba(108, 91, 255, 0.04);
  padding: 22rpx 24rpx;
  border-radius: 16rpx;
  font-size: 30rpx;
  color: #1A1B3A;
  font-weight: 600;
}

.exchange-tips-box {
  padding: 16rpx 24rpx;
  border-radius: 12rpx;
  background: rgba(108, 91, 255, 0.06);
}

.exchange-tips { font-size: 24rpx; color: #6B6F8D; }

.popup-actions {
  display: flex;
  gap: 16rpx;
}

.popup-btn-cancel,
.popup-btn-confirm {
  flex: 1;
  margin: 0;
  border-radius: 999rpx;
  font-size: 30rpx;
  font-weight: 700;
  height: 84rpx;
  line-height: 84rpx;
}

.popup-btn-cancel {
  background: #fff;
  color: #6D5BFF;
  border: 2rpx solid rgba(108, 91, 255, 0.32);
}

.popup-btn-confirm {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  border: none;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
}
</style>
