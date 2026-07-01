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
        <text>{{ t('profileErrorBook') }}</text>
        <text class="arrow">></text>
      </view>
      <view class="menu-item" @click="contactService">
        <text>{{ t('profile.contact_service') }}</text>
        <text class="arrow">></text>
      </view>
    </view>

    <button class="logout-btn" @click="logout">{{ t('profile.logout') }}</button>

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
          <button @click="closeExchangePopup">{{ t('common.cancel', '取消') }}</button>
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
import { t as i18nT } from '@/utils/i18n.js'
import { exchangeTime, getAsset } from '@/api/learning.js'
import { getPointsExchangeRate } from '@/api/sysConfig.js'

const langStore = useLangStore()
const fallbackTexts = {
  'profile.points': '积分',
  'profile.minutes': '分钟',
  'profile.free_time': '免费时长',
  'profile.exchange_time': '兑换时长',
  'profile.watch_history': '观看记录',
  'profile.checkin_record': '签到记录',
  'profile.my_collections': '我的收藏',
  profileErrorBook: 'Error Book',
  'profile.contact_service': '联系客服',
  'profile.logout': '退出登录',
  'profile.exchange_popup_title': '积分兑换时长',
  'profile.current_points': '当前积分',
  'profile.input_points_to_exchange': '请输入要兑换的积分',
  'profile.points_equals': '积分可兑换',
  'profile.exchange_invalid': '请输入有效积分',
  'profile.exchange_success': '兑换成功',
  'profile.coming_soon': '功能开发中',
  'profile.logout_success': '已退出登录',
  'common.confirm': '确定',
}

const t = (key, defaultText = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return defaultText || fallbackTexts[key] || key
}

const userInfo = ref({ avatar: '', nickname: 'English Learner' })
const userAsset = ref({ totalPoints: 0, freeMinutes: 0 })
const exchangeRate = ref(100)
const exchangePoints = ref(100)

const exchangePopup = ref(null)

const exchangeMinutes = computed(() => Math.floor(Number(exchangePoints.value || 0) / exchangeRate.value))

const loadUserInfo = () => {
  const raw = uni.getStorageSync('userInfo')
  if (raw && typeof raw === 'object') {
    userInfo.value = {
      avatar: raw.headerImg || raw.avatar || '',
      nickname: raw.nickName || raw.nickname || 'English Learner'
    }
  }
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

const logout = () => {
  uni.removeStorageSync('x-token')
  uni.removeStorageSync('userInfo')
  uni.showToast({ title: t('profile.logout_success'), icon: 'none' })
  uni.reLaunch({ url: '/pages/user/login' })
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
.profile-container { padding: 30rpx; }
.user-header { display: flex; align-items: center; margin-bottom: 40rpx; }
.avatar { width: 120rpx; height: 120rpx; border-radius: 50%; background: #eee; margin-right: 30rpx; }
.nickname { font-size: 36rpx; font-weight: bold; }

.asset-panel { display: flex; background: #fff; border-radius: 20rpx; padding: 40rpx 20rpx; box-shadow: 0 4rpx 10rpx rgba(0,0,0,0.05); margin-bottom: 40rpx; align-items: center; }
.asset-item { flex: 1; display: flex; flex-direction: column; align-items: center; border-right: 1px solid #f0f0f0; }
.asset-item:last-of-type { border-right: none; }
.asset-val { font-size: 36rpx; font-weight: bold; color: #ff9800; }
.asset-label { font-size: 24rpx; color: #888; margin-top: 10rpx; }
.exchange-btn-wrap { flex: 1; display: flex; justify-content: center; }
.exchange-btn { background: #409eff; color: #fff; font-size: 24rpx; }

.menu-list { background: #fff; border-radius: 20rpx; overflow: hidden; margin-bottom: 60rpx; }
.menu-item { display: flex; justify-content: space-between; padding: 30rpx; border-bottom: 1px solid #f5f5f5; font-size: 30rpx; }
.menu-item:last-child { border-bottom: none; }
.arrow { color: #ccc; }

.logout-btn { background: #f5f5f5; color: #e53935; }

/* Popup styles */
.popup-box { width: 600rpx; background: #fff; border-radius: 20rpx; padding: 40rpx; }
.popup-title { font-size: 32rpx; font-weight: bold; text-align: center; margin-bottom: 30rpx; }
.popup-content { margin-bottom: 30rpx; display:flex; flex-direction:column; gap:20rpx;}
.exchange-input { border: 1px solid #ddd; padding: 20rpx; border-radius: 10rpx; margin-top: 20rpx;}
.exchange-tips { font-size: 24rpx; color: #888; }
.popup-actions { display: flex; justify-content: space-between; }
.popup-actions button { width: 45%; }
</style>
