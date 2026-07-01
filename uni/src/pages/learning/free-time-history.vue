<template>
  <view class="free-time-history-container">
    <view class="header-row">
      <view class="back-btn" @click="goBack">‹</view>
      <text class="header-title">{{ t('learningFreeTimeHistoryTitle') }}</text>
      <view class="header-gap"></view>
    </view>

    <view class="summary-card">
      <text class="summary-label">{{ t('learningFreeTimeHistoryCurrentMinutes') }}</text>
      <text class="summary-value">{{ currentFreeMinutes }} {{ t('profile.minutes') }}</text>
    </view>

    <scroll-view class="list-scroll" scroll-y @scrolltolower="loadMore">
      <view v-if="!loading && recordList.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('learningFreeTimeHistoryEmpty') }}</text>
      </view>

      <view v-for="(item, index) in recordList" :key="`${item.id}-${index}`" class="record-card">
        <view class="record-left">
          <text class="record-reason">{{ getReasonText(item.reason) }}</text>
          <text class="record-time">{{ formatTime(item.createdAt) }}</text>
        </view>

        <view class="record-right" :class="item.minuteChange >= 0 ? 'is-plus' : 'is-minus'">
          {{ item.minuteChange > 0 ? '+' : '' }}{{ item.minuteChange }} {{ t('profile.minutes') }}
        </view>
      </view>

      <view v-if="recordList.length > 0" class="load-more">
        <text v-if="loading">{{ t('learningFreeTimeHistoryLoading') }}</text>
        <text v-else-if="noMore">{{ t('learningFreeTimeHistoryNoMore') }}</text>
        <text v-else>{{ t('learningFreeTimeHistoryPullMore') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { resolveApiMessage, t as i18nT } from '@/utils/i18n.js'
import { getAsset, getLearningFreeTimeRecordList } from '@/api/learning.js'

const langStore = useLangStore()
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const noMore = ref(false)
const recordList = ref([])
const currentFreeMinutes = ref(0)

const t = (key, defaultText = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return defaultText || key
}

const normalizeRecordItem = (item) => ({
  id: Number(item?.id || item?.ID || 0),
  minuteChange: Number(item?.minuteChange || item?.MinuteChange || 0),
  reason: String(item?.reason || item?.Reason || ''),
  createdAt: item?.createdAt || item?.CreatedAt || ''
})

const getReasonText = (reason) => {
  const key = String(reason || '').trim()
  if (!key) return t('learningFreeTimeHistoryTitle')
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
      currentFreeMinutes.value = Number(res.data.freeMinutes || 0)
    }
  } catch (error) {
    currentFreeMinutes.value = 0
  }
}

const fetchFreeTimeRecordList = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) {
    return
  }

  loading.value = true
  const nextPage = isLoadMore ? page.value + 1 : 1

  try {
    const res = await getLearningFreeTimeRecordList({ page: nextPage, pageSize })
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
  fetchFreeTimeRecordList(true)
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
  fetchFreeTimeRecordList(false)
})
</script>

<style scoped>
.free-time-history-container {
  min-height: 100vh;
  padding: 20rpx;
  box-sizing: border-box;
  background: #f5f7fb;
}

.header-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.back-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 32rpx;
  background: #ffffff;
  color: #111827;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 44rpx;
  line-height: 1;
}

.header-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #111827;
}

.header-gap {
  width: 64rpx;
  height: 64rpx;
}

.summary-card {
  margin-bottom: 16rpx;
  padding: 24rpx;
  border-radius: 16rpx;
  background: #ffffff;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-label {
  font-size: 26rpx;
  color: #6b7280;
}

.summary-value {
  font-size: 36rpx;
  color: #0284c7;
  font-weight: 700;
}

.list-scroll {
  height: calc(100vh - 200rpx);
}

.empty-wrap {
  padding: 120rpx 0;
  text-align: center;
}

.empty-text {
  color: #9ca3af;
  font-size: 28rpx;
}

.record-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(15, 23, 42, 0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.record-left {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  flex: 1;
}

.record-reason {
  font-size: 28rpx;
  color: #111827;
  font-weight: 600;
}

.record-time {
  font-size: 24rpx;
  color: #6b7280;
}

.record-right {
  font-size: 30rpx;
  font-weight: 700;
  margin-left: 20rpx;
}

.record-right.is-plus {
  color: #059669;
}

.record-right.is-minus {
  color: #dc2626;
}

.load-more {
  text-align: center;
  padding: 20rpx 0 40rpx;
  color: #9ca3af;
  font-size: 24rpx;
}
</style>
