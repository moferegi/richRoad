<template>
  <view class="learning-home">
    <view class="home-nav">
      <view class="nav-brand">
        <image v-if="appLogo" class="brand-logo" :src="appLogo" mode="aspectFit" />
        <view v-else class="brand-logo fallback">{{ (appTitle || 'R').slice(0, 1).toUpperCase() }}</view>
        <text class="brand-title">{{ appTitle || 'RichRoad' }}</text>
      </view>
    </view>

    <view class="banner-wrap" v-if="bannerList.length > 0">
      <BannerSwiper :lists="bannerList" />
    </view>

    <view class="video-section">
      <view class="section-head">
        <text class="section-title">{{ t('home.video_zone') }}</text>
      </view>

      <scroll-view class="category-tabs" scroll-x show-scrollbar="false">
        <view
          class="tab-chip"
          :class="{ active: currentCategoryId === 0 }"
          @tap="switchCategory(0)"
        >
          {{ t('home.all_category') }}
        </view>
        <view
          v-for="cat in videoCategories"
          :key="cat.id"
          class="tab-chip"
          :class="{ active: currentCategoryId === cat.id }"
          @tap="switchCategory(cat.id)"
        >
          {{ localText(cat.name) }}
        </view>
      </scroll-view>

      <view class="video-grid">
        <view class="video-card" v-for="item in videoList" :key="item.id" @tap="goDetail(item.id)">
          <image class="video-cover" :src="item.coverUrl || ''" mode="aspectFit" />
          <view class="video-content">
            <text class="video-title">{{ localText(item.name) }}</text>
            <view class="video-meta">
              <text>{{ t('video.views') }} {{ item.viewCount || 0 }}</text>
              <text>{{ t('video.users') }} {{ item.userCount || 0 }}</text>
            </view>
          </view>
        </view>
      </view>

      <view v-if="loading" class="load-text">{{ t('common.loading') }}</view>
      <view v-else-if="finished && videoList.length > 0" class="load-text">{{ t('common.no_more') }}</view>
      <view v-else-if="!loading && videoList.length === 0" class="load-text">{{ t('common.empty') }}</view>
    </view>

    <!-- 打卡区域暂时隐藏
    <view class="checkin-board"></view>
    -->
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

const normalizeId = (item) => Number(item?.id || item?.ID || 0)

const loadAppInfo = async () => {
  await appConfigStore.loadConfig({ force: true, localeOnly: true })
}

const loadBanner = async () => {
  const res = await getBannerList()
  if (res.code !== 0) {
    bannerList.value = []
    return
  }
  const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
  bannerList.value = list
}

const loadCategories = async () => {
  const res = await getVideoCategoryList({ page: 1, pageSize: 100, showHome: true })
  if (res.code !== 0 || !res.data) {
    videoCategories.value = []
    return
  }
  const list = Array.isArray(res.data.list) ? res.data.list : []
  videoCategories.value = list
    .map((item) => ({ ...item, id: normalizeId(item) }))
    .filter((item) => item.id > 0)
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

  const res = await getVideoSeriesList(params)
  loading.value = false

  if (res.code !== 0 || !res.data) return

  const list = Array.isArray(res.data.list) ? res.data.list : []
  const normalized = list
    .map((item) => ({ ...item, id: normalizeId(item) }))
    .filter((item) => item.id > 0)

  videoList.value = append ? [...videoList.value, ...normalized] : normalized

  const total = Number(res.data.total || 0)
  finished.value = videoList.value.length >= total || normalized.length < pageSize
}

const switchCategory = async (categoryId) => {
  if (currentCategoryId.value === categoryId && videoList.value.length > 0) return
  currentCategoryId.value = categoryId
  page.value = 1
  finished.value = false
  videoList.value = []
  await loadSeries(false)
}

const goDetail = (seriesId) => {
  uni.navigateTo({ url: `/pages/learning/video-detail?seriesId=${seriesId}` })
}

const reloadHome = async () => {
  await Promise.all([loadAppInfo(), loadBanner(), loadCategories()])
  page.value = 1
  finished.value = false
  videoList.value = []
  await loadSeries(false)
}

onLoad(() => {
  homeInited.value = true
  reloadHome()
})

onShow(() => {
  if (!homeInited.value) {
    homeInited.value = true
    return
  }
  if (!homeShownOnce.value) {
    homeShownOnce.value = true
    return
  }
  reloadHome()
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
  background: linear-gradient(180deg, #f8f4e7 0%, #f4efe1 100%);
}

.home-nav {
  padding: calc(var(--status-bar-height, 0px) + 24rpx) 22rpx 16rpx;
}

.nav-brand {
  height: 86rpx;
  border-radius: 22rpx;
  background: rgba(255, 253, 248, 0.96);
  border: 1rpx solid rgba(20, 184, 166, 0.2);
  box-shadow: 0 12rpx 24rpx rgba(120, 53, 15, 0.1);
  display: flex;
  align-items: center;
  padding: 0 20rpx;
  gap: 14rpx;
}

.brand-logo {
  width: 52rpx;
  height: 52rpx;
  border-radius: 14rpx;
  background: #fff;
}

.brand-logo.fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 26rpx;
  font-weight: 700;
  background: linear-gradient(135deg, #0f766e, #f97316);
}

.brand-title {
  color: #7c2d12;
  font-size: 30rpx;
  font-weight: 700;
  letter-spacing: 1rpx;
}

.banner-wrap {
  margin-top: 4rpx;
}

.video-section {
  margin: 12rpx 20rpx 0;
  padding: 24rpx 20rpx 30rpx;
  border-radius: 24rpx;
  background: rgba(255, 253, 248, 0.96);
  border: 1rpx solid rgba(20, 184, 166, 0.18);
  box-shadow: 0 16rpx 32rpx rgba(120, 53, 15, 0.1);
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 18rpx;
}

.section-title {
  font-size: 32rpx;
  color: #7c2d12;
  font-weight: 700;
}

.category-tabs {
  white-space: nowrap;
  margin-bottom: 20rpx;
}

.tab-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  height: 56rpx;
  padding: 0 24rpx;
  border-radius: 999rpx;
  margin-right: 14rpx;
  font-size: 24rpx;
  color: #64748b;
  background: #fff;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
}

.tab-chip.active {
  color: #ffffff;
  border-color: #0f766e;
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16rpx;
}

.video-card {
  border-radius: 18rpx;
  overflow: hidden;
  background: #fffdf8;
  border: 1rpx solid rgba(20, 184, 166, 0.16);
  display: flex;
  flex-direction: column;
}

.video-cover {
  width: 100%;
  height: 220rpx;
  display: block;
  background: #e5e7eb;
  flex-shrink: 0;
}

.video-content {
  padding: 14rpx;
}

.video-title {
  font-size: 24rpx;
  color: #1e293b;
  font-weight: 600;
  line-height: 1.35;
  min-height: 66rpx;
}

.video-meta {
  margin-top: 8rpx;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  font-size: 20rpx;
  color: #64748b;
}

.load-text {
  text-align: center;
  color: #94a3b8;
  font-size: 22rpx;
  padding-top: 20rpx;
}
</style>
