<template>
  <view class="detail-container">
    <view class="header-row">
      <view class="back-btn" @click="goBack">‹</view>
      <text class="header-title">{{ t('videoDetail.title', '视频详情') }}</text>
      <view class="header-gap"></view>
    </view>

    <view class="series-card">
      <image class="series-cover" :src="seriesInfo.coverUrl || ''" mode="aspectFill" />
      <view class="series-main">
        <text class="series-name">{{ localText(seriesInfo.name) || `#${seriesInfo.id || '-'}` }}</text>
        <text class="series-sub">
          {{ t('video.views', '浏览量') }} {{ Number(seriesInfo.viewCount || 0) }}
          ·
          {{ t('video.users', '观看人数') }} {{ Number(seriesInfo.userCount || 0) }}
        </text>
        <text class="series-sub" v-if="latestWatchAtText">
          {{ t('videoDetail.last_watch') }} {{ latestWatchAtText }}
        </text>
        <button class="continue-btn" size="mini" @click="continueWatch">{{ t('videoDetail.continue_watch') }}</button>
      </view>
    </view>

    <view class="episode-head">
      <text class="episode-title">{{ t('videoDetail.episode_list') }}</text>
      <text class="episode-count">{{ t('videoDetail.total_episode') }} {{ episodeList.length }}</text>
    </view>

    <scroll-view class="episode-scroll" scroll-y>
      <view v-if="!loading && episodeList.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('videoDetail.no_episode') }}</text>
      </view>

      <view
        v-for="episode in episodeList"
        :key="episode.id"
        class="episode-card"
        :class="{ active: resumeEpisodeId === episode.id, watched: !!progressMap[episode.id] }"
        @click="openEpisode(episode.id)"
      >
        <view class="episode-line">
          <text class="episode-name">{{ localText(episode.name) || `${t('videoDetail.episode')} ${episode.sort || episode.id}` }}</text>
          <text class="play-tip">{{ t('videoDetail.play_now') }}</text>
        </view>

        <view class="progress-row" v-if="progressMap[episode.id]">
          <text class="progress-text">{{ t('videoDetail.watched_to') }} {{ formatSeconds(progressMap[episode.id]) }}</text>
          <text class="progress-text time" v-if="progressUpdatedMap[episode.id]">{{ t('videoDetail.last_watch') }} {{ progressUpdatedMap[episode.id] }}</text>
        </view>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText as i18nLocalText, t as i18nT } from '@/utils/i18n.js'
import { findVideoEpisode, findVideoSeries, getSeriesWatchProgressList, getVideoEpisodeList } from '@/api/learning.js'

const langStore = useLangStore()
const fallbackTexts = {
  'videoDetail.title': '视频详情',
  'videoDetail.continue_watch': '继续观看',
  'videoDetail.episode_list': '分集列表',
  'videoDetail.total_episode': '共',
  'videoDetail.no_episode': '暂无可播放分集',
  'videoDetail.episode': '第',
  'videoDetail.play_now': '立即播放',
  'videoDetail.watched_to': '已观看到',
  'videoDetail.last_watch': '最近观看',
  'videoDetail.load_failed': '加载失败',
  'videoDetail.invalid_series': '剧集参数异常',
  'video.views': '浏览量',
  'video.users': '观看人数'
}

const t = (key, defaultText = '') => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return defaultText || fallbackTexts[key] || key
}

const localText = (value) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  return i18nLocalText(value, locale)
}

const seriesId = ref(0)
const preferredEpisodeId = ref(0)
const loading = ref(false)
const seriesInfo = ref({
  id: 0,
  name: '',
  coverUrl: '',
  viewCount: 0,
  userCount: 0
})
const episodeList = ref([])
const progressMap = ref({})
const progressUpdatedMap = ref({})
const latestProgressEpisodeId = ref(0)
const latestWatchAtText = ref('')
const initialized = ref(false)

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const resumeEpisodeId = computed(() => {
  if (preferredEpisodeId.value && episodeList.value.some((item) => item.id === preferredEpisodeId.value)) {
    return preferredEpisodeId.value
  }
  if (latestProgressEpisodeId.value && episodeList.value.some((item) => item.id === latestProgressEpisodeId.value)) {
    return latestProgressEpisodeId.value
  }
  return episodeList.value[0]?.id || 0
})

const formatSeconds = (secs) => {
  const total = Math.max(0, Math.floor(Number(secs || 0)))
  const m = Math.floor(total / 60)
  const s = total % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

const formatDateTime = (value) => {
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

const goBack = () => {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/learning/home' })
    }
  })
}

const openEpisode = (id) => {
  if (!id) {
    return
  }
  uni.navigateTo({ url: `/pages/learning/player?id=${id}` })
}

const continueWatch = () => {
  if (!resumeEpisodeId.value) {
    uni.showToast({ title: t('videoDetail.no_episode'), icon: 'none' })
    return
  }
  openEpisode(resumeEpisodeId.value)
}

const loadSeries = async () => {
  const res = await findVideoSeries(seriesId.value)
  if (res.code !== 0 || !res.data) {
    return false
  }
  seriesInfo.value = {
    id: getEntityId(res.data),
    name: res.data?.name || '',
    coverUrl: res.data?.coverUrl || '',
    viewCount: Number(res.data?.viewCount || 0),
    userCount: Number(res.data?.userCount || 0)
  }
  return true
}

const loadEpisodes = async () => {
  const res = await getVideoEpisodeList({ page: 1, pageSize: 300, seriesId: seriesId.value })
  if (res.code !== 0 || !res.data) {
    episodeList.value = []
    return false
  }

  const list = Array.isArray(res.data.list) ? res.data.list : []
  episodeList.value = list.map((item) => ({
    ...item,
    id: getEntityId(item),
    sort: Number(item?.sort || 0)
  }))
  return true
}

const loadProgress = async () => {
  const res = await getSeriesWatchProgressList(seriesId.value)
  if (res.code !== 0) {
    progressMap.value = {}
    progressUpdatedMap.value = {}
    latestWatchAtText.value = ''
    latestProgressEpisodeId.value = 0
    return
  }

  const list = Array.isArray(res.data) ? res.data : []
  const map = {}
  const updatedMap = {}
  for (const item of list) {
    const episodeId = Number(item?.episodeId || 0)
    if (episodeId <= 0) continue
    map[episodeId] = Number(item?.progressSecs || 0)
    updatedMap[episodeId] = formatDateTime(item?.updatedAt)
  }
  progressMap.value = map
  progressUpdatedMap.value = updatedMap
  latestProgressEpisodeId.value = Number(list[0]?.episodeId || 0)
  latestWatchAtText.value = formatDateTime(list[0]?.updatedAt)
}

const resolveSeriesIdByEpisode = async () => {
  if (!preferredEpisodeId.value) {
    return false
  }
  const res = await findVideoEpisode(preferredEpisodeId.value)
  if (res.code !== 0 || !res.data) {
    return false
  }
  const sid = Number(res.data?.seriesId || 0)
  if (!sid) {
    return false
  }
  seriesId.value = sid
  return true
}

const initPage = async () => {
  if (!seriesId.value) {
    const ok = await resolveSeriesIdByEpisode()
    if (!ok) {
      uni.showToast({ title: t('videoDetail.invalid_series'), icon: 'none' })
      return
    }
  }

  loading.value = true
  const [seriesOk, episodeOk] = await Promise.all([loadSeries(), loadEpisodes(), loadProgress()])
  loading.value = false
  initialized.value = true
  if (!seriesOk || !episodeOk) {
    uni.showToast({ title: t('videoDetail.load_failed'), icon: 'none' })
  }
}

onLoad((options) => {
  seriesId.value = Number(options?.seriesId || 0)
  preferredEpisodeId.value = Number(options?.episodeId || 0)
  initPage()
})

onShow(() => {
  if (!initialized.value || !seriesId.value) {
    return
  }
  loadProgress()
})
</script>

<style scoped>
.detail-container {
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

.series-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 16rpx;
  display: flex;
  gap: 16rpx;
  margin-bottom: 20rpx;
}

.series-cover {
  width: 220rpx;
  height: 132rpx;
  border-radius: 12rpx;
  background: #e5e7eb;
  flex-shrink: 0;
}

.series-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.series-name {
  font-size: 30rpx;
  color: #111827;
  font-weight: 600;
  line-height: 1.4;
}

.series-sub {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #6b7280;
}

.continue-btn {
  margin-top: 12rpx;
  background: #2563eb;
  color: #ffffff;
  border-radius: 10rpx;
  font-size: 24rpx;
}

.episode-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.episode-title {
  font-size: 30rpx;
  color: #111827;
  font-weight: 600;
}

.episode-count {
  font-size: 24rpx;
  color: #6b7280;
}

.episode-scroll {
  height: calc(100vh - 360rpx);
}

.empty-wrap {
  padding: 120rpx 0;
  text-align: center;
}

.empty-text {
  color: #9ca3af;
  font-size: 28rpx;
}

.episode-card {
  background: #ffffff;
  border-radius: 14rpx;
  padding: 20rpx;
  margin-bottom: 14rpx;
  border: 1rpx solid #eef2ff;
}

.episode-card.active {
  border-color: #2563eb;
  box-shadow: 0 6rpx 20rpx rgba(37, 99, 235, 0.12);
}

.episode-card.watched {
  border-color: #bfdbfe;
}

.episode-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.episode-name {
  font-size: 28rpx;
  color: #111827;
  font-weight: 500;
  flex: 1;
  line-height: 1.4;
}

.play-tip {
  margin-left: 16rpx;
  font-size: 24rpx;
  color: #2563eb;
}

.progress-row {
  margin-top: 10rpx;
}

.progress-text {
  font-size: 23rpx;
  color: #6b7280;
}

.progress-text.time {
  margin-left: 16rpx;
}
</style>
