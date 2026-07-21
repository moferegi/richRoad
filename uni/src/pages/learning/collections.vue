<template>
  <view class="collection-container">
    <!-- 渐变 Hero 头部 -->
    <view class="page-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-content">
        <view class="hero-back-btn" @click="goBack">
          <text class="hero-back-icon">‹</text>
        </view>
        <text class="hero-title">{{ t('collection.title') }}</text>
        <view class="hero-gap"></view>
      </view>
    </view>

    <!-- Tab 胶囊选择器 -->
    <scroll-view class="tab-scroll" scroll-x show-scrollbar="false">
      <view class="tab-track">
        <view
          v-for="tab in tabs"
          :key="tab.type"
          class="tab-pill"
          :class="{ active: targetType === tab.type }"
          @click="switchType(tab.type)"
        >
          <text class="tab-pill-label">{{ t(tab.key) }}</text>
        </view>
      </view>
    </scroll-view>

    <scroll-view class="list-scroll" scroll-y>
      <view v-if="!loading && collectionList.length === 0" class="empty-wrap">
        <view class="empty-icon">
          <text class="empty-icon-text">★</text>
        </view>
        <text class="empty-text">{{ t('collection.empty') }}</text>
      </view>

      <view
        v-for="item in collectionList"
        :key="`${item.id}-${item.targetType}-${item.targetId}`"
        class="collection-card"
      >
        <view class="card-head">
          <text class="type-tag">{{ getTypeLabel(item.targetType) }}</text>
          <text class="time-text">{{ formatTime(item.createdAt) }}</text>
        </view>

        <view class="card-body" @click="openDetail(item)">
          <image
            v-if="item.targetType === 3"
            class="video-cover"
            :src="getExternalUrl(item.video?.coverUrl || '')"
            mode="aspectFill"
          />
          <view v-else class="type-icon-circle">
            <text class="type-icon-char">{{ item.targetType === 1 ? 'W' : 'S' }}</text>
          </view>
          <view class="text-wrap">
            <text class="title-text">{{ getTitle(item) }}</text>
            <text class="desc-text">{{ getDesc(item) }}</text>
          </view>
        </view>

        <view class="card-footer">
          <button size="mini" class="remove-btn" @click.stop="removeCollection(item)">
            {{ t('collection.remove') }}
          </button>
        </view>
      </view>

      <view v-if="collectionList.length > 0" class="load-more">
        <text v-if="loading">{{ t('collection.loading') }}</text>
        <text v-else-if="noMore">{{ t('collection.no_more') }}</text>
        <button v-else class="load-more-btn" size="mini" @click="loadMore">{{ t('common.load_more') }}</button>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText as i18nLocalText, t as i18nT } from '@/utils/i18n.js'
import { getCollectionDetailList, uncollect } from '@/api/learning.js'
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

const tabs = [
  { type: 0, key: 'collection.tab_all' },
  { type: 1, key: 'collection.tab_word' },
  { type: 2, key: 'collection.tab_sentence' },
  { type: 3, key: 'collection.tab_video' }
]

const targetType = ref(0)
const page = ref(1)
const pageSize = 10
const total = ref(0)
const loading = ref(false)
const noMore = ref(false)
const collectionList = ref([])

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const stripWordTags = (text) => {
  const source = String(text || '')
  return source
    .replace(/<w id="\d+">(.*?)<\/w>/g, '$1')
    .replace(/<[^>]+>/g, '')
    .trim()
}

const normalizeCollectionItem = (item) => ({
  id: getEntityId(item),
  targetType: Number(item?.targetType || 0),
  targetId: Number(item?.targetId || 0),
  createdAt: item?.createdAt || '',
  word: item?.word || null,
  sentence: item?.sentence || null,
  video: item?.video || null
})

const getTypeLabel = (type) => {
  if (type === 1) return t('collection.word')
  if (type === 2) return t('collection.sentence')
  if (type === 3) return t('collection.video')
  return '-'
}

const getTitle = (item) => {
  if (item.targetType === 1) {
    return String(item.word?.word || `#${item.targetId}`)
  }
  if (item.targetType === 2) {
    return stripWordTags(item.sentence?.english || `#${item.targetId}`)
  }
  if (item.targetType === 3) {
    const episodeName = localText(item.video?.episodeName || '')
    const seriesName = localText(item.video?.seriesName || '')
    return episodeName || seriesName || `#${item.targetId}`
  }
  return `#${item.targetId}`
}

const getDesc = (item) => {
  if (item.targetType === 1) {
    return localText(item.word?.explanation || '')
  }
  if (item.targetType === 2) {
    return localText(item.sentence?.translate || '')
  }
  if (item.targetType === 3) {
    const seriesName = localText(item.video?.seriesName || '')
    const episodeName = localText(item.video?.episodeName || '')
    if (seriesName && episodeName) {
      return `${seriesName} · ${episodeName}`
    }
    return seriesName || episodeName || ''
  }
  return ''
}

const formatTime = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const fetchCollectionList = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) {
    return
  }
  loading.value = true

  const nextPage = isLoadMore ? page.value + 1 : 1
  const params = {
    page: nextPage,
    pageSize
  }
  if (targetType.value > 0) {
    params.targetType = targetType.value
  }

  const res = await getCollectionDetailList(params)
  if (res.code === 0 && res.data) {
    const list = Array.isArray(res.data.list) ? res.data.list.map(normalizeCollectionItem) : []
    total.value = Number(res.data.total || 0)
    page.value = nextPage
    if (isLoadMore) {
      collectionList.value = [...collectionList.value, ...list]
    } else {
      collectionList.value = list
    }
    noMore.value = collectionList.value.length >= total.value || list.length < pageSize
  }

  loading.value = false
}

const switchType = (type) => {
  if (targetType.value === type) {
    return
  }
  targetType.value = type
  page.value = 1
  noMore.value = false
  fetchCollectionList(false)
}

const loadMore = () => {
  fetchCollectionList(true)
}

const goBack = () => {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/learning/profile' })
    }
  })
}

const openDetail = (item) => {
  if (item.targetType === 1) {
    // switchTab 不支持传参，通过 storage 传递 wordId 给 typing 页
    uni.setStorageSync('typing-pending-word-id', Number(item.targetId || 0))
    uni.switchTab({ url: '/pages/learning/typing' })
    return
  }

  if (item.targetType === 2) {
    const episodeId = Number(item.sentence?.episodeId || 0)
    const sentenceId = Number(item.sentence?.id || item.sentence?.ID || 0)
    if (episodeId > 0) {
      const query = sentenceId > 0 ? `?id=${episodeId}&sentenceId=${sentenceId}` : `?id=${episodeId}`
      uni.navigateTo({ url: `/pages/learning/player${query}` })
      return
    }
  }

  if (item.targetType === 3) {
    const seriesId = Number(item.video?.seriesId || 0)
    const episodeId = Number(item.video?.episodeId || 0)
    if (seriesId > 0) {
      const query = episodeId > 0 ? `?seriesId=${seriesId}&episodeId=${episodeId}` : `?seriesId=${seriesId}`
      uni.navigateTo({ url: `/pages/learning/video-detail${query}` })
      return
    }
    if (episodeId > 0) {
      uni.navigateTo({ url: `/pages/learning/player?id=${episodeId}` })
      return
    }
  }

  uni.showToast({ title: t('collection.open_failed'), icon: 'none' })
}

const removeCollection = async (item) => {
  const res = await uncollect(item.targetType, item.targetId)
  if (res.code !== 0) {
    uni.showToast({ title: t('collection.remove_failed'), icon: 'none' })
    return
  }
  uni.showToast({ title: t('collection.remove_success'), icon: 'success' })
  collectionList.value = collectionList.value.filter((row) => !(row.targetType === item.targetType && row.targetId === item.targetId))
  total.value = Math.max(0, total.value - 1)
  if (collectionList.value.length === 0) {
    noMore.value = true
  }
}

onShow(() => {
  fetchCollectionList(false)
})
</script>

<style scoped>
.collection-container {
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

/* === Tab 胶囊选择器 === */
.tab-scroll {
  white-space: nowrap;
  padding: 16rpx 24rpx;
}

.tab-track {
  display: inline-flex;
  gap: 16rpx;
  padding: 4rpx 0;
}

.tab-pill {
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

.tab-pill-label {
  font-size: 26rpx;
  color: #6B6F8D;
  font-weight: 500;
  white-space: nowrap;
}

.tab-pill.active {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-color: transparent;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
}

.tab-pill.active .tab-pill-label {
  color: #fff;
  font-weight: 700;
}

.tab-pill:active {
  transform: scale(0.96);
}

/* === 列表 === */
.list-scroll {
  height: calc(100vh - 280rpx);
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

.collection-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.08);
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 14rpx;
}

.type-tag {
  background: rgba(108, 91, 255, 0.12);
  color: #6D5BFF;
  font-size: 22rpx;
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  font-weight: 600;
}

.time-text {
  color: #A9AECB;
  font-size: 22rpx;
}

.card-body {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.video-cover {
  width: 160rpx;
  height: 96rpx;
  border-radius: 16rpx;
  background: #EDE9FE;
  flex-shrink: 0;
}

.type-icon-circle {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.1);
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.type-icon-char {
  font-size: 32rpx;
  font-weight: 800;
  color: #6D5BFF;
}

.text-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.title-text {
  font-size: 30rpx;
  color: #1A1B3A;
  font-weight: 600;
  line-height: 1.4;
}

.desc-text {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #6B6F8D;
  line-height: 1.5;
}

.card-footer {
  margin-top: 16rpx;
  display: flex;
  justify-content: flex-end;
}

.remove-btn {
  background: rgba(239, 68, 68, 0.1);
  color: #EF4444;
  border: 1rpx solid rgba(239, 68, 68, 0.22);
  font-size: 24rpx;
  border-radius: 999rpx;
  font-weight: 600;
}

.remove-btn:active {
  transform: scale(0.95);
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
