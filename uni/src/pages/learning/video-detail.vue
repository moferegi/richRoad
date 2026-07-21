<template>
  <view class="detail-container">
    <!-- 渐变 Hero 头部 -->
    <view class="page-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-bar">
        <view class="hero-back-btn" @click="goBack">‹</view>
        <text class="hero-title">{{ t('videoDetail.title') }}</text>
        <view class="hero-placeholder"></view>
      </view>
    </view>

    <!-- 系列信息浮卡 -->
    <view class="series-card">
      <image class="series-cover" :src="getExternalUrl(seriesInfo.coverUrl || '')" mode="aspectFit" />
      <view class="series-main">
        <text class="series-name">{{ localText(seriesInfo.name) || `#${seriesInfo.id || '-'}` }}</text>
        <text class="series-sub">
          {{ t('video.views') }} {{ Number(seriesInfo.viewCount || 0) }}
          ·
          {{ t('video.users') }} {{ Number(seriesInfo.userCount || 0) }}
        </text>
        <text class="series-sub" v-if="latestWatchAtText">
          {{ t('videoDetail.last_watch') }} {{ latestWatchAtText }}
        </text>
        <button class="continue-btn" size="mini" @click="continueWatch">{{ t('videoDetail.continue_watch') }}</button>
      </view>
    </view>

    <!-- 集数列表标题 -->
    <view class="episode-head">
      <view class="episode-accent"></view>
      <text class="episode-title">{{ t('videoDetail.episode_list') }}</text>
      <text class="episode-count">{{ t('videoDetail.total_episode') }} {{ episodeList.length }}</text>
    </view>

    <scroll-view class="episode-scroll" scroll-y>
      <view v-if="!loading && episodeList.length === 0" class="empty-wrap">
        <view class="empty-icon-wrap">
          <text class="empty-icon">▶</text>
        </view>
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
import { getExternalUrl } from '@/utils/url.js'

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
  background: #F5F3FF;
  padding-bottom: 40rpx;
  box-sizing: border-box;
  overflow-x: hidden;
}

/* === 渐变 Hero 头部 === */
.page-hero {
  position: relative;
  overflow: hidden;
  padding: calc(var(--status-bar-height, 0px) + 32rpx) 32rpx 48rpx;
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

.hero-bar {
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
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 44rpx;
  line-height: 1;
}

.hero-title {
  font-size: 36rpx;
  color: #fff;
  font-weight: 800;
  letter-spacing: 0.5rpx;
}

.hero-placeholder {
  width: 64rpx;
  height: 64rpx;
}

/* === 系列信息浮卡 === */
.series-card {
  margin: -28rpx 24rpx 24rpx;
  position: relative;
  z-index: 5;
  background: #FFFFFF;
  border-radius: 24rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  padding: 24rpx;
  display: flex;
  gap: 20rpx;
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.12);
}

.series-cover {
  width: 220rpx;
  height: 132rpx;
  border-radius: 16rpx;
  background: #EDE9FE;
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
  color: #1A1B3A;
  font-weight: 700;
  line-height: 1.4;
}

.series-sub {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #6B6F8D;
}

.continue-btn {
  margin-top: 12rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #ffffff;
  border-radius: 999rpx;
  font-size: 24rpx;
  font-weight: 700;
  border: none;
}

/* === 集数列表 === */
.episode-head {
  display: flex;
  align-items: center;
  gap: 12rpx;
  padding: 0 24rpx;
  margin-bottom: 16rpx;
}

.episode-accent {
  width: 8rpx;
  height: 32rpx;
  border-radius: 4rpx;
  background: linear-gradient(180deg, #6D5BFF, #9B8FFF);
}

.episode-title {
  flex: 1;
  font-size: 34rpx;
  color: #1A1B3A;
  font-weight: 800;
}

.episode-count {
  font-size: 24rpx;
  color: #6B6F8D;
}

.episode-scroll {
  height: calc(100vh - 480rpx);
  padding: 0 24rpx;
  box-sizing: border-box;
  width: 100%;
}

/* === 空状态 === */
.empty-wrap {
  padding: 100rpx 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
}

.empty-icon-wrap {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-icon {
  font-size: 36rpx;
  color: #6D5BFF;
}

.empty-text {
  color: #A9AECB;
  font-size: 28rpx;
}

/* === 集数卡片 === */
.episode-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.episode-card:active {
  transform: scale(0.98);
}

.episode-card.active {
  border-color: rgba(108, 91, 255, 0.5);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.2);
}

.episode-card.watched {
  border-color: rgba(108, 91, 255, 0.3);
}

.episode-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.episode-name {
  font-size: 28rpx;
  color: #1A1B3A;
  font-weight: 600;
  flex: 1;
  line-height: 1.4;
}

.play-tip {
  margin-left: 16rpx;
  font-size: 24rpx;
  color: #6D5BFF;
  font-weight: 600;
}

.progress-row {
  margin-top: 10rpx;
}

.progress-text {
  font-size: 23rpx;
  color: #6B6F8D;
}

.progress-text.time {
  margin-left: 16rpx;
}
</style>
