<template>
  <view class="checkin-record-container">
    <view class="header-row">
      <view class="back-btn" @click="goBack">‹</view>
      <text class="header-title">{{ t('learningCheckinRecordTitle') }}</text>
      <view class="header-gap"></view>
    </view>

    <view class="stats-card">
      <view class="stats-top">
        <view class="stats-item">
          <text class="stats-label">{{ t('learningCheckinRecordContinuousDays') }}</text>
          <text class="stats-value">{{ continuousDays }}</text>
        </view>
        <view class="stats-item">
          <text class="stats-label">{{ t('learningCheckinRecordTotalDays') }}</text>
          <text class="stats-value">{{ totalDays }}</text>
        </view>
      </view>

      <view class="stats-bottom">
        <text class="checkin-status">{{ checkedToday ? t('learningCheckinRecordCheckedToday') : t('learningCheckinRecordNotCheckedToday') }}</text>
        <button
          v-if="!checkedToday"
          class="checkin-btn"
          size="mini"
          :loading="checkinLoading"
          @click="handleCheckin"
        >
          {{ t('learningCheckinRecordCheckinNow') }}
        </button>
        <text v-else class="checked-tag">{{ t('learningCheckinRecordCheckedTag') }}</text>
      </view>
    </view>

    <scroll-view class="list-scroll" scroll-y @scrolltolower="loadMore">
      <view v-if="!loading && recordList.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('learningCheckinRecordEmpty') }}</text>
      </view>

      <view v-for="(item, index) in recordList" :key="`${item.id}-${index}`" class="record-card">
        <view class="record-left">
          <text class="record-title">{{ t('learningCheckinRecordAward') }}</text>
          <text class="record-date">{{ formatDate(item.checkinDate) }}</text>
        </view>
        <text class="record-right">+{{ item.pointAward }}</text>
      </view>

      <view v-if="recordList.length > 0" class="load-more">
        <text v-if="loading">{{ t('learningCheckinRecordLoading') }}</text>
        <text v-else-if="noMore">{{ t('learningCheckinRecordNoMore') }}</text>
        <text v-else>{{ t('learningCheckinRecordPullMore') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { resolveApiMessage, t as i18nT } from '@/utils/i18n.js'
import { doCheckin, getCheckinStats, getLearningCheckinRecordList } from '@/api/learning.js'

const langStore = useLangStore()
const page = ref(1)
const pageSize = 20
const total = ref(0)
const loading = ref(false)
const checkinLoading = ref(false)
const noMore = ref(false)
const recordList = ref([])
const continuousDays = ref(0)
const totalDays = ref(0)
const checkedToday = ref(false)

const t = (key, defaultText = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return defaultText || key
}

const normalizeRecordItem = (item) => ({
  id: Number(item?.id || item?.ID || 0),
  checkinDate: item?.checkinDate || item?.CheckinDate || '',
  pointAward: Number(item?.pointAward || item?.PointAward || 0)
})

const formatDate = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const fetchStats = async () => {
  try {
    const res = await getCheckinStats()
    if (res.code !== 0 || !res.data) return
    continuousDays.value = Number(res.data.continuousDays || 0)
    totalDays.value = Number(res.data.totalDays || 0)
    checkedToday.value = !!res.data.checkedToday
  } catch (error) {
    continuousDays.value = 0
    totalDays.value = 0
    checkedToday.value = false
  }
}

const fetchCheckinRecords = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) {
    return
  }

  loading.value = true
  const nextPage = isLoadMore ? page.value + 1 : 1

  try {
    const res = await getLearningCheckinRecordList({ page: nextPage, pageSize })
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

const handleCheckin = async () => {
  if (checkinLoading.value || checkedToday.value) {
    return
  }

  checkinLoading.value = true
  try {
    const res = await doCheckin()
    if (res.code !== 0) {
      uni.showToast({ title: resolveApiMessage(res.msg, 'operationFailed'), icon: 'none' })
      return
    }

    const award = Number(res?.data?.award || 0)
    uni.showToast({ title: `${t('learningCheckinRecordSuccess')} +${award}`, icon: 'none' })

    page.value = 1
    noMore.value = false
    await fetchStats()
    await fetchCheckinRecords(false)
  } catch (error) {
    uni.showToast({ title: t('networkError'), icon: 'none' })
  } finally {
    checkinLoading.value = false
  }
}

const loadMore = () => {
  fetchCheckinRecords(true)
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
  fetchStats()
  fetchCheckinRecords(false)
})
</script>

<style scoped>
.checkin-record-container {
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

.stats-card {
  margin-bottom: 16rpx;
  padding: 24rpx;
  border-radius: 16rpx;
  background: #ffffff;
}

.stats-top {
  display: flex;
  justify-content: space-between;
}

.stats-item {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.stats-label {
  font-size: 24rpx;
  color: #6b7280;
}

.stats-value {
  font-size: 36rpx;
  color: #111827;
  font-weight: 700;
}

.stats-bottom {
  margin-top: 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.checkin-status {
  font-size: 24rpx;
  color: #374151;
}

.checkin-btn {
  background: #2563eb;
  color: #ffffff;
  font-size: 24rpx;
  padding: 0 24rpx;
}

.checked-tag {
  font-size: 24rpx;
  color: #059669;
}

.list-scroll {
  height: calc(100vh - 320rpx);
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
}

.record-title {
  font-size: 28rpx;
  color: #111827;
  font-weight: 600;
}

.record-date {
  font-size: 24rpx;
  color: #6b7280;
}

.record-right {
  font-size: 32rpx;
  color: #059669;
  font-weight: 700;
}

.load-more {
  text-align: center;
  padding: 20rpx 0 40rpx;
  color: #9ca3af;
  font-size: 24rpx;
}
</style>
