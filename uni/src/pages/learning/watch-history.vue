<template>
  <view class="watch-history-container">
    <!-- 渐变 Hero 头部 -->
    <view class="page-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-content">
        <view class="hero-back-btn" @click="goBack">
          <text class="hero-back-icon">‹</text>
        </view>
        <text class="hero-title">{{ t('learningWatchHistoryTitle') }}</text>
        <view class="hero-gap"></view>
      </view>
    </view>

    <scroll-view class="list-scroll" scroll-y>
      <view v-if="!loading && historyList.length === 0" class="empty-wrap">
        <view class="empty-icon">
          <text class="empty-icon-text">▶</text>
        </view>
        <text class="empty-text">{{ t('learningWatchHistoryEmpty') }}</text>
      </view>

      <view v-for="(item, index) in historyList" :key="`${item.episodeId}-${index}`" class="history-card">
        <!-- 顶部：封面 + 标题区 -->
        <view class="card-top">
          <view class="cover-wrap">
            <image class="cover" :src="getExternalUrl(item.coverUrl || '')" mode="aspectFill" />
            <view class="cover-duration">
              <text class="duration-text">{{ formatSeconds(item.progressSecs) }}</text>
            </view>
          </view>
          <view class="card-info">
            <text class="series-name">{{ localText(item.seriesName) || t('learningWatchHistoryUnknownSeries') }}</text>
            <text class="episode-name">{{ localText(item.episodeName) || t('learningWatchHistoryUnknownEpisode') }}</text>
          </view>
        </view>

        <!-- 底部：进度信息 + 按钮 -->
        <view class="card-bottom">
          <view class="meta-line">
            <text class="meta-item">{{ t('learningWatchHistoryProgress') }} {{ formatSeconds(item.progressSecs) }}</text>
            <text class="meta-dot">·</text>
            <text class="meta-item">{{ formatTime(item.updatedAt) }}</text>
          </view>
          <button class="continue-btn" @click.stop="continueWatch(item)">
            <text class="btn-text">{{ t('learningWatchHistoryContinue') }}</text>
            <text class="btn-arrow">→</text>
          </button>
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
import { getExternalUrl } from '@/utils/url.js'

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

/* === 列表 === */
.list-scroll {
  height: calc(100vh - 220rpx);
  padding: 16rpx 24rpx;
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
  font-size: 36rpx;
  color: #6D5BFF;
}

.empty-text {
  color: #A9AECB;
  font-size: 28rpx;
}

.history-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  margin-bottom: 20rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.08);
  overflow: hidden;
}

.history-card:active {
  transform: scale(0.98);
}

/* === 卡片顶部：封面 + 标题 === */
.card-top {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 20rpx 0;
}

.cover-wrap {
  position: relative;
  flex-shrink: 0;
}

.cover {
  width: 200rpx;
  height: 120rpx;
  border-radius: 16rpx;
  background: #EDE9FE;
}

.cover-duration {
  position: absolute;
  bottom: 8rpx;
  right: 8rpx;
  padding: 4rpx 12rpx;
  border-radius: 999rpx;
  background: rgba(26, 27, 58, 0.7);
}

.duration-text {
  font-size: 20rpx;
  color: #fff;
}

.card-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8rpx;
}

.series-name {
  font-size: 28rpx;
  color: #1A1B3A;
  font-weight: 600;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.episode-name {
  font-size: 24rpx;
  color: #6B6F8D;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* === 卡片底部：进度信息 + 按钮 === */
.card-bottom {
  padding: 20rpx;
  display: flex;
  flex-direction: column;
}

.meta-line {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding-bottom: 16rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.08);
}

.meta-item {
  font-size: 22rpx;
  color: #6B6F8D;
}

.meta-dot {
  font-size: 22rpx;
  color: #A9AECB;
}

.continue-btn {
  width: auto;
  align-self: flex-end;
  height: 76rpx;
  margin-top: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  padding: 0 36rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  border: none;
  border-radius: 999rpx;
  font-size: 26rpx;
  font-weight: 700;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.3);
}

.continue-btn::after {
  border: none;
}

.continue-btn:active {
  transform: scale(0.98);
  opacity: 0.9;
}

.btn-arrow {
  font-size: 24rpx;
}

.load-more {
  text-align: center;
  color: #A9AECB;
  font-size: 24rpx;
  padding: 20rpx 0 40rpx;
}

.load-more-btn {
  border: 1rpx solid rgba(108, 91, 255, 0.3);
  color: #6D5BFF;
  background: #FFFFFF;
  border-radius: 999rpx;
  padding: 0 28rpx;
  font-weight: 600;
}
</style>