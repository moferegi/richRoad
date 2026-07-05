<template>
  <view class="profile-container">
    <view class="user-header">
      <image class="avatar" :src="userInfo.avatar"></image>
      <text class="nickname">{{ userInfo.nickname }}</text>
    </view>

    <!-- 资产面板 -->
    <view class="asset-panel">
      <view class="asset-item" @click="goTo('point_history')">
        <text class="asset-val">{{ userAsset.totalPoints }}</text>
        <text class="asset-label">{{ t('profile.points') }}</text>
      </view>
      <view class="asset-item" @click="goTo('free_time_history')">
        <text class="asset-val">{{ userAsset.freeMinutes }} {{ t('profile.minutes') }}</text>
        <text class="asset-label">{{ t('profile.free_time') }}</text>
      </view>
      <view class="exchange-btn-wrap">
        <button class="exchange-btn" size="mini" @click="showExchangePopup">{{ t('profile.exchange_time') }}</button>
      </view>
    </view>

    <!-- 列表菜单 -->
    <view class="menu-list">
      <view class="menu-item" @click="goTo('watch_history')">
        <text>{{ t('profile.watch_history') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="goTo('checkin_calendar')">
        <text>{{ t('profile.checkin_record') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="goTo('collections')">
        <text>{{ t('profile.my_collections') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="goTo('error_words')">
        <text>{{ t('profile.error_book') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="contactService">
        <text>{{ t('profile.contact_service') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="showLanguageSwitcher">
        <text>{{ t('profile.switch_language') }}</text>
        <text class="arrow">></text>
      </view>
    </view>

    <button class="logout-btn" @click="logout">{{ t('profile.logout') }}</button>

    <lang-switch v-model="showLangPicker" />

    <!-- 积分兑换弹窗 -->
    <uni-popup ref="exchangePopup" type="center">
      <view class="popup-box">
        <view class="popup-title">{{ t('profile.exchange_popup_title') }}</view>
        <view class="popup-content">
          <text>{{ t('profile.current_points') }}: {{ userAsset.totalPoints }}</text>
          <input type="number" v-model="exchangePoints" :placeholder="t('profile.input_points_to_exchange')" class="exchange-input"/>
          <text class="exchange-tips">{{ exchangePoints }} {{ t('profile.points_equals') }} {{ exchangeMinutes }} {{ t('profile.minutes') }}</text>
        </view>
        <view class="popup-actions">
          <button @click="closeExchangePopup">{{ t('common.cancel') }}</button>
          <button type="primary" @click="submitExchange">{{ t('common.confirm') }}</button>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText as i18nLocalText, t as i18nT } from '@/utils/i18n.js'
import { exchangeTime, getAsset } from '@/api/learning.js'
import { getPointsExchangeRate } from '@/api/sysConfig.js'
import langSwitch from '@/components/lang-switch/lang-switch.vue'

const langStore = useLangStore()
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
    title: t('profile.logout_confirm_title'),
    content: t('profile.logout_confirm_desc'),
    cancelText: t('common.cancel'),
    confirmText: t('common.confirm'),
    success: (res) => {
      if (!res.confirm) {
        return
      }
      uni.removeStorageSync('x-token')
      uni.removeStorageSync('userInfo')
      uni.showToast({ title: t('profile.logout_success'), icon: 'none' })
      uni.reLaunch({ url: '/pages/user/login' })
    }
  })
}

onMounted(() => {
  loadUserInfo()
  loadAsset()
  loadExchangeRate()
})

onShow(() => {
  loadUserInfo()
  loadAsset()
  loadExchangeRate()
})
</script>

<style scoped>
.profile-container {
  min-height: 100vh;
  padding: 26rpx;
  box-sizing: border-box;
  background: linear-gradient(180deg, #f8f4e7 0%, #f4efe1 100%);
}

.user-header {
  display: flex;
  align-items: center;
  margin-bottom: 28rpx;
  padding: 22rpx;
  border-radius: 22rpx;
  background: rgba(255, 253, 248, 0.96);
  border: 1rpx solid rgba(20, 184, 166, 0.2);
  box-shadow: 0 12rpx 26rpx rgba(120, 53, 15, 0.1);
}

.avatar {
  width: 114rpx;
  height: 114rpx;
  border-radius: 50%;
  background: #f1f5f9;
  margin-right: 22rpx;
  border: 3rpx solid rgba(15, 118, 110, 0.28);
}

.nickname {
  font-size: 36rpx;
  font-weight: 700;
  color: #7c2d12;
}

.asset-panel {
  display: flex;
  align-items: center;
  background: rgba(255, 253, 248, 0.96);
  border-radius: 22rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.2);
  box-shadow: 0 12rpx 26rpx rgba(120, 53, 15, 0.1);
  padding: 28rpx 12rpx;
  margin-bottom: 26rpx;
}

.asset-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  border-right: 1rpx solid rgba(20, 184, 166, 0.18);
}

.asset-item:last-of-type { border-right: none; }

.asset-val {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f766e;
}

.asset-label {
  font-size: 24rpx;
  color: #64748b;
  margin-top: 8rpx;
}

.exchange-btn-wrap { flex: 1; display: flex; justify-content: center; }

.exchange-btn {
  height: 66rpx;
  line-height: 66rpx;
  margin: 0;
  padding: 0 26rpx;
  border-radius: 999rpx;
  font-size: 24rpx;
  color: #fff;
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
}

.menu-list {
  background: rgba(255, 253, 248, 0.96);
  border-radius: 22rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.2);
  overflow: hidden;
  margin-bottom: 40rpx;
}

.menu-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28rpx;
  border-bottom: 1rpx solid rgba(20, 184, 166, 0.12);
  font-size: 30rpx;
  color: #334155;
}

.menu-item:last-child { border-bottom: none; }
.arrow { color: #0f766e; font-weight: 700; }

.logout-btn {
  height: 86rpx;
  line-height: 86rpx;
  border-radius: 20rpx;
  font-size: 30rpx;
  font-weight: 700;
  color: #fff;
  background: linear-gradient(120deg, #dc2626 0%, #f97316 100%);
}

.popup-box {
  width: 620rpx;
  background: #fffdf8;
  border-radius: 22rpx;
  padding: 34rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.2);
}

.popup-title {
  font-size: 32rpx;
  font-weight: 700;
  text-align: center;
  color: #7c2d12;
  margin-bottom: 24rpx;
}

.popup-content {
  margin-bottom: 26rpx;
  display: flex;
  flex-direction: column;
  gap: 18rpx;
  color: #334155;
}

.exchange-input {
  border: 1rpx solid rgba(146, 64, 14, 0.16);
  background: #fff;
  padding: 18rpx;
  border-radius: 12rpx;
}

.exchange-tips { font-size: 24rpx; color: #64748b; }

.popup-actions {
  display: flex;
  justify-content: space-between;
  gap: 14rpx;
}

.popup-actions button {
  flex: 1;
  margin: 0;
  border-radius: 14rpx;
}
</style>
