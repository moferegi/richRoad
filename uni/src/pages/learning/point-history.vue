<template>
  <view class="point-history-container">
    <!-- 渐变 Hero 头部 -->
    <view class="page-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-content">
        <view class="hero-back-btn" @click="goBack">
          <text class="hero-back-icon">‹</text>
        </view>
        <text class="hero-title">{{ t('learningPointHistoryTitle') }}</text>
        <view class="hero-gap"></view>
      </view>
    </view>

    <!-- 摘要浮卡 -->
    <view class="summary-card">
      <view class="summary-left">
        <text class="summary-label">{{ t('learningPointHistoryCurrentPoints') }}</text>
      </view>
      <text class="summary-value">{{ currentPoints }}</text>
    </view>

    <scroll-view class="list-scroll" scroll-y @scrolltolower="loadMore">
      <view v-if="!loading && recordList.length === 0" class="empty-wrap">
        <view class="empty-icon">
          <text class="empty-icon-text">★</text>
        </view>
        <text class="empty-text">{{ t('learningPointHistoryEmpty') }}</text>
      </view>

      <view v-for="(item, index) in recordList" :key="`${item.id}-${index}`" class="record-card">
        <view class="record-icon-wrap" :class="item.pointChange >= 0 ? 'is-plus' : 'is-minus'">
          <text class="record-icon">{{ item.pointChange >= 0 ? '+' : '−' }}</text>
        </view>
        <view class="record-content">
          <text class="record-reason">{{ getReasonText(item.reason) }}</text>
          <text class="record-time">{{ formatTime(item.createdAt) }}</text>
        </view>
        <view class="record-value-pill" :class="item.pointChange >= 0 ? 'is-plus' : 'is-minus'">
          <text class="record-value-text">{{ item.pointChange > 0 ? '+' : '' }}{{ item.pointChange }}</text>
        </view>
      </view>

      <view v-if="recordList.length > 0" class="load-more">
        <text v-if="loading">{{ t('learningPointHistoryLoading') }}</text>
        <text v-else-if="noMore">{{ t('learningPointHistoryNoMore') }}</text>
        <text v-else>{{ t('learningPointHistoryPullMore') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { resolveApiMessage, t as i18nT } from '@/utils/i18n.js'
import { getAsset, getLearningPointRecordList } from '@/api/learning.js'

const langStore = useLangStore()
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const noMore = ref(false)
const recordList = ref([])
const currentPoints = ref(0)

const t = (key, defaultText = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return defaultText || key
}

const normalizeRecordItem = (item) => ({
  id: Number(item?.id || item?.ID || 0),
  pointChange: Number(item?.pointChange || item?.PointChange || 0),
  reason: String(item?.reason || item?.Reason || ''),
  createdAt: item?.createdAt || item?.CreatedAt || ''
})

const getReasonText = (reason) => {
  const key = String(reason || '').trim()
  if (!key) return t('pointsChange')
  const translated = t(key)
  if (translated && translated !== key) return translated
  return key
}

const formatTime = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''

  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${hh}:${mm}`
}

const fetchAsset = async () => {
  try {
    const res = await getAsset()
    if (res.code === 0 && res.data) {
      currentPoints.value = Number(res.data.totalPoints || 0)
    }
  } catch (error) {
    currentPoints.value = 0
  }
}

const fetchPointRecordList = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) {
    return
  }

  loading.value = true
  const nextPage = isLoadMore ? page.value + 1 : 1

  try {
    const res = await getLearningPointRecordList({ page: nextPage, pageSize })
    if (res.code !== 0) {
      uni.showToast({ title: resolveApiMessage(res.msg, 'operationFailed'), icon: 'none' })
      return
    }

    const data = res.data || {}
    const list = Array.isArray(data.list) ? data.list.map(normalizeRecordItem) : []
    total.value = Number(data.total || 0)
    page.value = nextPage

    if (isLoadMore) {
      recordList.value = [...recordList.value, ...list]
    } else {
      recordList.value = list
    }

    noMore.value = recordList.value.length >= total.value || list.length < pageSize
  } catch (error) {
    uni.showToast({ title: t('networkError'), icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  fetchPointRecordList(true)
}

const goBack = () => {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/learning/profile' })
    }
  })
}

onShow(() => {
  page.value = 1
  noMore.value = false
  fetchAsset()
  fetchPointRecordList(false)
})
</script>

<style scoped>
.point-history-container {
  min-height: 100vh;
  background: #F5F3FF;
  overflow-x: hidden;
}

/* === 渐变 Hero 头部 === */
.page-hero {
  position: relative;
  overflow: hidden;
  padding: calc(var(--status-bar-height, 0px) + 36rpx) 32rpx 48rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-radius: 0 0 36rpx 36rpx;
  box-sizing: border-box;
  width: 100%;
}

.hero-decor-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.12);
}

.decor-1 {
  width: 240rpx;
  height: 240rpx;
  top: -80rpx;
  right: -40rpx;
}

.decor-2 {
  width: 160rpx;
  height: 160rpx;
  bottom: -60rpx;
  left: 200rpx;
}

.hero-content {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.hero-back-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hero-back-btn:active {
  transform: scale(0.92);
  background: rgba(255, 255, 255, 0.35);
}

.hero-back-icon {
  color: #fff;
  font-size: 44rpx;
  line-height: 1;
}

.hero-title {
  font-size: 36rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 0.5rpx;
}

.hero-gap {
  width: 64rpx;
  height: 64rpx;
}

/* === 摘要浮卡 === */
.summary-card {
  margin: -28rpx 24rpx 16rpx;
  position: relative;
  z-index: 5;
  padding: 28rpx 32rpx;
  border-radius: 24rpx;
  background: #FFFFFF;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.12);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-left {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.summary-label {
  font-size: 26rpx;
  color: #6B6F8D;
  font-weight: 500;
}

.summary-value {
  font-size: 40rpx;
  color: #6D5BFF;
  font-weight: 800;
}

/* === 列表 === */
.list-scroll {
  height: calc(100vh - 340rpx);
  padding: 0 24rpx;
  box-sizing: border-box;
  width: 100%;
}

.empty-wrap {
  padding: 120rpx 0;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
}

.empty-icon {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.1);
  border: 1rpx solid rgba(108, 91, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-icon-text {
  font-size: 40rpx;
  color: #6D5BFF;
}

.empty-text {
  color: #A9AECB;
  font-size: 28rpx;
}

.record-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.06);
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.record-card:active {
  transform: scale(0.98);
}

.record-icon-wrap {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.record-icon-wrap.is-plus {
  background: rgba(108, 91, 255, 0.1);
}

.record-icon-wrap.is-minus {
  background: rgba(239, 68, 68, 0.1);
}

.record-icon {
  font-size: 36rpx;
  font-weight: 800;
}

.record-icon-wrap.is-plus .record-icon {
  color: #6D5BFF;
}

.record-icon-wrap.is-minus .record-icon {
  color: #EF4444;
}

.record-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  min-width: 0;
}

.record-reason {
  font-size: 28rpx;
  color: #1A1B3A;
  font-weight: 600;
}

.record-time {
  font-size: 24rpx;
  color: #6B6F8D;
}

.record-value-pill {
  min-width: 100rpx;
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.record-value-pill.is-plus {
  background: rgba(108, 91, 255, 0.12);
}

.record-value-pill.is-minus {
  background: rgba(239, 68, 68, 0.12);
}

.record-value-text {
  font-size: 28rpx;
  font-weight: 800;
}

.record-value-pill.is-plus .record-value-text {
  color: #6D5BFF;
}

.record-value-pill.is-minus .record-value-text {
  color: #EF4444;
}

.load-more {
  text-align: center;
  padding: 20rpx 0 40rpx;
  color: #A9AECB;
  font-size: 24rpx;
}
</style>
