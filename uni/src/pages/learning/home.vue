<template>
  <view class="learning-home">
    <!-- 渐变 Hero 头部 -->
    <view class="home-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-content">
        <view class="hero-logo-wrap">
          <image v-if="appLogo" class="hero-logo" :src="appLogo" mode="aspectFit" />
          <view v-else class="hero-logo fallback">{{ (appTitle || 'R').slice(0, 1).toUpperCase() }}</view>
        </view>
        <view class="hero-text">
          <text class="hero-title">{{ appTitle || 'RichRoad' }}</text>
          <text class="hero-subtitle">{{ t('home.video_zone') }}</text>
        </view>
      </view>
    </view>

    <!-- Banner 卡片 -->
    <view class="banner-wrap" v-if="bannerList.length > 0">
      <BannerSwiper :lists="bannerList" />
    </view>

    <!-- 视频区 -->
    <view class="video-zone">
      <view class="zone-header">
        <view class="zone-accent"></view>
        <text class="zone-title">{{ t('home.video_zone') }}</text>
        <view class="zone-header-right">
          <view class="tag-entry-btn" @tap="goTagFilter">
            <text class="tag-entry-label">{{ t('home.more_tags') }}</text>
            <text class="tag-entry-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- 胶囊式分类选择器 -->
      <scroll-view class="cat-scroll" scroll-x show-scrollbar="false" :scroll-into-view="catScrollIntoView">
        <view class="cat-track">
          <view
            id="cat-pill-0"
            class="cat-pill"
            :class="{ active: currentCategoryId === 0 }"
            @tap="switchCategory(0)"
          >
            <text class="cat-pill-label">{{ t('home.all_category') }}</text>
          </view>
          <view
            v-for="cat in videoCategories"
            :key="cat.id"
            :id="'cat-pill-' + cat.id"
            class="cat-pill"
            :class="{ active: currentCategoryId === cat.id }"
            @tap="switchCategory(cat.id)"
          >
            <text class="cat-pill-label">{{ localText(cat.name) }}</text>
          </view>
        </view>
      </scroll-view>

      <!-- 视频卡片网格 -->
      <view class="video-grid">
        <view class="video-card" v-for="item in videoList" :key="item.id" @tap="goDetail(item.id)">
          <image class="video-cover" :src="getExternalUrl(item.coverUrl || '')" mode="aspectFill" />
          <view class="video-play-badge">
            <text class="play-icon">▶</text>
          </view>
          <view class="video-overlay">
            <text class="video-title">{{ localText(item.name) }}</text>
            <view class="video-stats">
              <text class="stat-item">{{ t('video.views') }} {{ item.viewCount || 0 }}</text>
              <text class="stat-dot">·</text>
              <text class="stat-item">{{ t('video.users') }} {{ item.userCount || 0 }}</text>
            </view>
          </view>
        </view>
      </view>

      <view v-if="loading" class="load-status">{{ t('common.loading') }}</view>
      <view v-else-if="finished && videoList.length > 0" class="load-status">{{ t('common.no_more') }}</view>
      <view v-else-if="!loading && videoList.length === 0" class="load-status">{{ t('common.empty') }}</view>
    </view>
    <custom-tab-bar />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { getBannerList } from '@/api/homePage.js'
import { getVideoCategoryList, getVideoSeriesList } from '@/api/learning.js'
import { getExternalUrl } from '@/utils/url.js'
import BannerSwiper from '@/pages/tabBar/components/swiper.vue'
import CustomTabBar from '@/components/custom-tab-bar/custom-tab-bar.vue'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const t = (key) => {
  const text = i18nT(key, locale.value)
  if (text && text !== key) return text
  return key
}

const localText = (value) => i18nLocalText(value, locale.value)

const appTitle = computed(() => appConfigStore.appName || '')
const appLogo = computed(() => getExternalUrl(appConfigStore.appLogo || ''))
const bannerList = ref([])

const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

const videoCategories = ref([])
const currentCategoryId = ref(0)
const videoList = ref([])
const homeInited = ref(false)
const homeShownOnce = ref(false)
const homeLocaleLoaded = ref('')
const reloading = ref(false)
const lastLoadTime = ref(0)
const catScrollIntoView = ref('')
const CACHE_TTL = 5 * 60 * 1000 // 5 分钟缓存窗口

const normalizeId = (item) => Number(item?.id || item?.ID || 0)

const loadAppInfo = async () => {
  await appConfigStore.loadConfig({ force: true, localeOnly: true })
}

const loadBanner = async () => {
  try {
    const res = await getBannerList()
    if (res.code !== 0) {
      bannerList.value = []
      return
    }
    const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
    bannerList.value = list
  } catch (error) {
    bannerList.value = []
  }
}

const loadCategories = async () => {
  try {
    const res = await getVideoCategoryList({ page: 1, pageSize: 100, showHome: true })
    if (res.code !== 0 || !res.data) {
      videoCategories.value = []
      return
    }
    const list = Array.isArray(res.data.list) ? res.data.list : []
    videoCategories.value = list
      .map((item) => ({ ...item, id: normalizeId(item) }))
      .filter((item) => item.id > 0)
  } catch (error) {
    videoCategories.value = []
  }
}

const loadSeries = async (append = false) => {
  if (loading.value || (finished.value && append)) return
  loading.value = true

  const params = {
    page: page.value,
    pageSize,
    showHome: true
  }
  if (currentCategoryId.value > 0) {
    params.categoryId = currentCategoryId.value
  }

  try {
    const res = await getVideoSeriesList(params)
    if (res.code !== 0 || !res.data) {
      if (!append) {
        videoList.value = []
      }
      return
    }

    const list = Array.isArray(res.data.list) ? res.data.list : []
    const normalized = list
      .map((item) => ({ ...item, id: normalizeId(item) }))
      .filter((item) => item.id > 0)

    videoList.value = append ? [...videoList.value, ...normalized] : normalized

    const total = Number(res.data.total || 0)
    finished.value = videoList.value.length >= total || normalized.length < pageSize
  } catch (error) {
    if (!append) {
      videoList.value = []
      finished.value = false
    }
  } finally {
    loading.value = false
  }
}

const switchCategory = async (categoryId) => {
  if (currentCategoryId.value === categoryId && videoList.value.length > 0) return
  currentCategoryId.value = categoryId
  catScrollIntoView.value = 'cat-pill-' + categoryId
  page.value = 1
  finished.value = false
  videoList.value = []
  await loadSeries(false)
}

const goTagFilter = () => {
  uni.navigateTo({ url: '/pages/learning/tag-filter' })
}

const goDetail = (seriesId) => {
  uni.navigateTo({ url: `/pages/learning/video-detail?seriesId=${seriesId}` })
}

const reloadHome = async () => {
  if (reloading.value) {
    return
  }
  reloading.value = true
  try {
    await Promise.allSettled([loadAppInfo(), loadBanner(), loadCategories()])
    page.value = 1
    finished.value = false
    videoList.value = []
    await loadSeries(false)
    homeLocaleLoaded.value = locale.value
    lastLoadTime.value = Date.now()
  } finally {
    reloading.value = false
  }
}

// 3.2 静默后台刷新：缓存过期后无感更新数据，不显示 loading
const silentRefresh = async () => {
  if (reloading.value) return
  reloading.value = true
  try {
    await Promise.allSettled([loadAppInfo(), loadBanner(), loadCategories()])
    const params = { page: 1, pageSize, showHome: true }
    if (currentCategoryId.value > 0) {
      params.categoryId = currentCategoryId.value
    }
    const res = await getVideoSeriesList(params)
    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data.list) ? res.data.list : []
      const normalized = list
        .map((item) => ({ ...item, id: normalizeId(item) }))
        .filter((item) => item.id > 0)
      if (normalized.length > 0) {
        videoList.value = normalized
      }
    }
    page.value = 1
    finished.value = false
    homeLocaleLoaded.value = locale.value
    lastLoadTime.value = Date.now()
  } finally {
    reloading.value = false
  }
}

onLoad(() => {
  homeInited.value = true
  reloadHome()
})

onShow(() => {
  uni.hideTabBar()
  langStore.updateTabBar(locale.value)
  if (!homeInited.value) {
    homeInited.value = true
    return
  }
  if (!homeShownOnce.value) {
    homeShownOnce.value = true
    if (!videoList.value.length) {
      reloadHome()
    }
    return
  }
  // 语言变更 → 全量重载
  if (homeLocaleLoaded.value !== locale.value) {
    reloadHome()
    return
  }
  // 无缓存 → 全量重载
  if (!videoList.value.length) {
    reloadHome()
    return
  }
  // 3.2 缓存命中（TTL 内）→ 直接渲染，跳过请求
  const now = Date.now()
  if (lastLoadTime.value > 0 && (now - lastLoadTime.value) < CACHE_TTL) {
    return
  }
  // 缓存过期 → 静默后台更新（保留旧数据不闪烁）
  silentRefresh()
})

onReachBottom(async () => {
  if (loading.value || finished.value) return
  page.value += 1
  await loadSeries(true)
})
</script>

<style scoped>
.learning-home {
  min-height: 100vh;
  background: #F5F3FF;
  padding-bottom: calc(96rpx + env(safe-area-inset-bottom));
}

/* === 渐变 Hero 头部 === */
.home-hero {
  position: relative;
  overflow: hidden;
  padding: calc(var(--status-bar-height, 0px) + 36rpx) 32rpx 48rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-radius: 0 0 36rpx 36rpx;
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
  gap: 20rpx;
}

.hero-logo-wrap {
  flex-shrink: 0;
}

.hero-logo {
  width: 80rpx;
  height: 80rpx;
  border-radius: 20rpx;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.18);
}

.hero-logo.fallback {
  width: 80rpx;
  height: 80rpx;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6D5BFF;
  font-size: 36rpx;
  font-weight: 800;
  background: #fff;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.18);
}

.hero-text {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.hero-title {
  font-size: 38rpx;
  font-weight: 800;
  color: #fff;
  letter-spacing: 0.5rpx;
  text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.12);
}

.hero-subtitle {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.82);
  font-weight: 500;
}

/* === Banner === */
.banner-wrap {
  margin: -28rpx 24rpx 8rpx;
  position: relative;
  z-index: 5;
}

/* === 视频区 === */
.video-zone {
  padding: 20rpx 32rpx calc(40rpx + 96rpx + env(safe-area-inset-bottom));
}

.zone-header {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 24rpx;
}

.zone-accent {
  width: 8rpx;
  height: 32rpx;
  border-radius: 4rpx;
  background: linear-gradient(180deg, #6D5BFF, #9B8FFF);
}

.zone-title {
  font-size: 34rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.zone-header-right {
  margin-left: auto;
  display: flex;
  align-items: center;
}

.tag-entry-btn {
  display: flex;
  align-items: center;
  gap: 4rpx;
  padding: 8rpx 16rpx 8rpx 20rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 2rpx 12rpx rgba(108, 91, 255, 0.28);
}

.tag-entry-label {
  font-size: 22rpx;
  color: #fff;
  font-weight: 600;
}

.tag-entry-arrow {
  font-size: 28rpx;
  color: #fff;
  font-weight: 700;
  line-height: 1;
}

/* === 胶囊式分类选择器 === */
.cat-scroll {
  white-space: nowrap;
  margin-bottom: 28rpx;
}

.cat-track {
  display: inline-flex;
  gap: 16rpx;
  padding: 4rpx 0;
}

.cat-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 12rpx 32rpx;
  border-radius: 999rpx;
  background: #fff;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 2rpx 8rpx rgba(108, 91, 255, 0.06);
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.cat-pill-label {
  font-size: 26rpx;
  color: #6B6F8D;
  font-weight: 500;
  white-space: nowrap;
}

.cat-pill.active {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-color: transparent;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
}

.cat-pill.active .cat-pill-label {
  color: #fff;
  font-weight: 700;
}

.cat-pill:active {
  transform: scale(0.96);
}

/* === 视频卡片 === */
.video-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 24rpx;
}

.video-card {
  position: relative;
  border-radius: 24rpx;
  overflow: hidden;
  aspect-ratio: 4 / 3;
  background: #EDE9FE;
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.12);
}

.video-card:active {
  opacity: 0.92;
  transform: scale(0.97);
}

.video-cover {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.video-play-badge {
  position: absolute;
  top: 16rpx;
  right: 16rpx;
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.92);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3;
  box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.2);
}

.play-icon {
  font-size: 20rpx;
  color: #fff;
  margin-left: 4rpx;
}

.video-overlay {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 56rpx 20rpx 18rpx;
  background: linear-gradient(to top, rgba(26, 27, 58, 0.82) 0%, rgba(26, 27, 58, 0) 100%);
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.video-title {
  font-size: 25rpx;
  color: #fff;
  font-weight: 600;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.video-stats {
  margin-top: 8rpx;
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.stat-item {
  font-size: 19rpx;
  color: rgba(255, 255, 255, 0.82);
}

.stat-dot {
  font-size: 19rpx;
  color: rgba(255, 255, 255, 0.5);
}

/* === 加载状态 === */
.load-status {
  text-align: center;
  color: #A9AECB;
  font-size: 22rpx;
  padding: 40rpx 0 16rpx;
  letter-spacing: 1rpx;
}
</style>
