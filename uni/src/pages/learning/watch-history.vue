<template>
  <view class="watch-history-container">
    <view class="header-row">
      <view class="back-btn" @click="goBack">‹</view>
      <text class="header-title">{{ t('learningWatchHistoryTitle') }}</text>
      <view class="header-gap"></view>
    </view>

    <scroll-view class="list-scroll" scroll-y>
      <view v-if="!loading && historyList.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('learningWatchHistoryEmpty') }}</text>
      </view>

      <view v-for="(item, index) in historyList" :key="`${item.episodeId}-${index}`" class="history-card">
        <image class="cover" :src="item.coverUrl || ''" mode="aspectFit" />

        <view class="content-wrap">
          <text class="series-name">{{ localText(item.seriesName) || t('learningWatchHistoryUnknownSeries') }}</text>
          <text class="episode-name">{{ localText(item.episodeName) || t('learningWatchHistoryUnknownEpisode') }}</text>

          <view class="meta-row">
            <text class="meta-label">{{ t('learningWatchHistoryProgress') }}</text>
            <text class="meta-value">{{ formatSeconds(item.progressSecs) }}</text>
          </view>
          <view class="meta-row">
            <text class="meta-label">{{ t('learningWatchHistoryUpdatedAt') }}</text>
            <text class="meta-value">{{ formatTime(item.updatedAt) }}</text>
          </view>

          <view class="action-row">
            <button size="mini" class="continue-btn" @click.stop="continueWatch(item)">{{ t('learningWatchHistoryContinue') }}</button>
          </view>
        </view>
      </view>

      <view v-if="historyList.length > 0" class="load-more">
        <text v-if="loading">{{ t('learningWatchHistoryLoading') }}</text>
        <text v-else-if="noMore">{{ t('learningWatchHistoryNoMore') }}</text>
        <button v-else class="load-more-btn" size="mini" @click="loadMore">{{ t('common.load_more') }}</button>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText as i18nLocalText, resolveApiMessage, t as i18nT } from '@/utils/i18n.js'
import { getWatchHistoryList } from '@/api/learning.js'

const langStore = useLangStore()
const page = ref(1)
const pageSize = 10
const total = ref(0)
const loading = ref(false)
const noMore = ref(false)
const historyList = ref([])

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

const normalizeHistoryItem = (item) => ({
  episodeId: Number(item?.episodeId || item?.EpisodeID || item?.episodeID || 0),
  episodeName: item?.episodeName || item?.EpisodeName || '',
  seriesId: Number(item?.seriesId || item?.SeriesID || item?.seriesID || 0),
  seriesName: item?.seriesName || item?.SeriesName || '',
  coverUrl: String(item?.coverUrl || item?.CoverURL || ''),
  progressSecs: Number(item?.progressSecs || item?.ProgressSecs || 0),
  updatedAt: item?.updatedAt || item?.UpdatedAt || ''
})

const formatSeconds = (value) => {
  const totalSecs = Math.max(0, Math.floor(Number(value || 0)))
  const h = Math.floor(totalSecs / 3600)
  const m = Math.floor((totalSecs % 3600) / 60)
  const s = totalSecs % 60
  if (h > 0) {
    return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  }
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
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

const fetchWatchHistory = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) {
    return
  }

  loading.value = true
  const nextPage = isLoadMore ? page.value + 1 : 1
  try {
    const res = await getWatchHistoryList({ page: nextPage, pageSize })
    if (res.code !== 0) {
      uni.showToast({ title: resolveApiMessage(res.msg, 'operationFailed'), icon: 'none' })
      return
    }

    const data = res.data || {}
    const list = Array.isArray(data.list) ? data.list.map(normalizeHistoryItem) : []
    total.value = Number(data.total || 0)
    page.value = nextPage

    if (isLoadMore) {
      historyList.value = [...historyList.value, ...list]
    } else {
      historyList.value = list
    }

    noMore.value = historyList.value.length >= total.value || list.length < pageSize
  } catch (error) {
    uni.showToast({ title: t('networkError'), icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  fetchWatchHistory(true)
}

const continueWatch = (item) => {
  if (!item.episodeId) {
    uni.showToast({ title: t('learningWatchHistoryUnknownEpisode'), icon: 'none' })
    return
  }
  uni.navigateTo({ url: `/pages/learning/player?id=${item.episodeId}` })
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
  fetchWatchHistory(false)
})
</script>

<style scoped>
.watch-history-container {
  min-height: 100vh;
  background: #f5f7fb;
  padding: 20rpx;
  box-sizing: border-box;
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
  color: #111827;
  font-weight: 700;
}

.header-gap {
  width: 64rpx;
  height: 64rpx;
}

.list-scroll {
  height: calc(100vh - 120rpx);
}

.empty-wrap {
  padding: 120rpx 0;
  text-align: center;
}

.empty-text {
  color: #9ca3af;
  font-size: 28rpx;
}

.history-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 16rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(15, 23, 42, 0.05);
  display: flex;
  gap: 16rpx;
}

.cover {
  width: 180rpx;
  height: 108rpx;
  border-radius: 10rpx;
  background: #e5e7eb;
  flex-shrink: 0;
}

.content-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.series-name {
  font-size: 28rpx;
  color: #111827;
  font-weight: 600;
  line-height: 1.4;
}

.episode-name {
  margin-top: 4rpx;
  font-size: 24rpx;
  color: #4b5563;
  line-height: 1.4;
}

.meta-row {
  margin-top: 8rpx;
  display: flex;
  justify-content: space-between;
}

.meta-label {
  font-size: 22rpx;
  color: #6b7280;
}

.meta-value {
  font-size: 22rpx;
  color: #111827;
}

.action-row {
  margin-top: 12rpx;
  display: flex;
  justify-content: flex-end;
}

.continue-btn {
  background: #eff6ff;
  color: #2563eb;
  border: 1rpx solid #bfdbfe;
  border-radius: 10rpx;
  font-size: 24rpx;
}

.load-more {
  text-align: center;
  color: #9ca3af;
  font-size: 24rpx;
  padding: 20rpx 0 40rpx;
}

.load-more-btn {
  border: 1rpx solid rgba(20, 184, 166, 0.28);
  color: #0f766e;
  background: #fffdf8;
  border-radius: 999rpx;
  padding: 0 28rpx;
}

.watch-history-container {
  background: linear-gradient(180deg, #f8f4e7 0%, #f4efe1 100%);
}

.back-btn {
  border-radius: 18rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
  background: #fffdf8;
  color: #7c2d12;
}

.header-title { color: #7c2d12; }

.history-card {
  background: rgba(255, 253, 248, 0.96);
  border: 1rpx solid rgba(20, 184, 166, 0.18);
  box-shadow: 0 10rpx 24rpx rgba(120, 53, 15, 0.1);
}

.series-name { color: #1e293b; }
.episode-name, .meta-label { color: #64748b; }

.continue-btn {
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
  color: #fff;
  border: none;
  border-radius: 999rpx;
}

.load-more, .empty-text { color: #94a3b8; }
</style>