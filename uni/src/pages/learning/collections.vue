<template>
  <view class="collection-container">
    <view class="header-row">
      <view class="back-btn" @click="goBack">‹</view>
      <text class="header-title">{{ t('collection.title', '我的收藏') }}</text>
      <view class="header-gap"></view>
    </view>

    <view class="tab-row">
      <view
        v-for="tab in tabs"
        :key="tab.type"
        class="tab-item"
        :class="{ active: targetType === tab.type }"
        @click="switchType(tab.type)"
      >
        {{ t(tab.key, tab.label) }}
      </view>
    </view>

    <scroll-view class="list-scroll" scroll-y @scrolltolower="loadMore">
      <view v-if="!loading && collectionList.length === 0" class="empty-wrap">
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
            :src="item.video?.coverUrl || ''"
            mode="aspectFill"
          />
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
        <text v-else>{{ t('collection.pull_more') }}</text>
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

const langStore = useLangStore()
const fallbackTexts = {
  'collection.tab_all': '全部',
  'collection.tab_word': '单词',
  'collection.tab_sentence': '句子',
  'collection.tab_video': '视频',
  'collection.title': '我的收藏',
  'collection.empty': '暂无收藏内容',
  'collection.remove': '取消收藏',
  'collection.loading': '加载中...',
  'collection.no_more': '没有更多了',
  'collection.pull_more': '上拉加载更多',
  'collection.word': '单词',
  'collection.sentence': '句子',
  'collection.video': '视频',
  'collection.open_failed': '当前内容不可打开',
  'collection.remove_success': '已取消收藏',
  'collection.remove_failed': '取消失败'
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

const tabs = [
  { type: 0, key: 'collection.tab_all', label: '全部' },
  { type: 1, key: 'collection.tab_word', label: '单词' },
  { type: 2, key: 'collection.tab_sentence', label: '句子' },
  { type: 3, key: 'collection.tab_video', label: '视频' }
]

const targetType = ref(0)
const page = ref(1)
const pageSize = 20
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

.tab-row {
  display: flex;
  align-items: center;
  background: #ffffff;
  border-radius: 16rpx;
  padding: 10rpx;
  margin-bottom: 20rpx;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 14rpx 8rpx;
  color: #6b7280;
  font-size: 26rpx;
  border-radius: 12rpx;
}

.tab-item.active {
  background: #2563eb;
  color: #ffffff;
  font-weight: 600;
}

.list-scroll {
  height: calc(100vh - 160rpx);
}

.empty-wrap {
  padding: 120rpx 0;
  text-align: center;
}

.empty-text {
  color: #9ca3af;
  font-size: 28rpx;
}

.collection-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 20rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 4rpx 16rpx rgba(15, 23, 42, 0.05);
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14rpx;
}

.type-tag {
  background: #dbeafe;
  color: #1d4ed8;
  font-size: 22rpx;
  padding: 6rpx 12rpx;
  border-radius: 10rpx;
}

.time-text {
  color: #9ca3af;
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
  border-radius: 12rpx;
  background: #e5e7eb;
  flex-shrink: 0;
}

.text-wrap {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.title-text {
  font-size: 30rpx;
  color: #111827;
  font-weight: 600;
  line-height: 1.4;
}

.desc-text {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #6b7280;
  line-height: 1.5;
}

.card-footer {
  margin-top: 16rpx;
  display: flex;
  justify-content: flex-end;
}

.remove-btn {
  background: #fef2f2;
  color: #dc2626;
  border: 1rpx solid #fecaca;
  font-size: 24rpx;
  border-radius: 10rpx;
}

.load-more {
  text-align: center;
  color: #9ca3af;
  font-size: 24rpx;
  padding: 20rpx 0 40rpx;
}
</style>
