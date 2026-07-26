<template>
  <view class="tag-filter-page">
    <!-- 渐变 Hero 头部（与 video-detail 一致） -->
    <view class="page-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-bar">
        <view class="hero-back-btn" @click="goBack">‹</view>
        <text class="hero-title">{{ t('home.more_tags') }}</text>
        <view class="hero-placeholder"></view>
      </view>
    </view>

    <!-- 筛选区 -->
    <view class="filter-section">
      <view class="filter-row">
        <!-- 分类下拉 -->
        <view class="dropdown-trigger" @click="showCatPicker = true">
          <text class="dropdown-label">{{ currentCategoryLabel }}</text>
          <text class="dropdown-arrow">▼</text>
        </view>

        <!-- 标签下拉（多选） -->
        <view class="dropdown-trigger tag-trigger" @click="showTagPicker = true">
          <text class="dropdown-label">{{ tagTriggerLabel }}</text>
          <text class="dropdown-arrow">▼</text>
        </view>
      </view>
    </view>

    <!-- 分类选择器弹层 -->
    <view class="picker-overlay" v-if="showCatPicker" @click="showCatPicker = false">
      <view class="picker-panel" @click.stop>
        <view class="picker-list">
          <view
            class="picker-item"
            :class="{ active: currentCategoryId === 0 }"
            @click="selectCategory(0)"
          >
            <text>{{ t('home.all_category') }}</text>
            <text v-if="currentCategoryId === 0" class="picker-check">✓</text>
          </view>
          <view
            v-for="cat in categories"
            :key="cat.id"
            class="picker-item"
            :class="{ active: currentCategoryId === cat.id }"
            @click="selectCategory(cat.id)"
          >
            <text>{{ localText(cat.name) }}</text>
            <text v-if="currentCategoryId === cat.id" class="picker-check">✓</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 标签多选弹层 -->
    <view class="picker-overlay" v-if="showTagPicker" @click="cancelTagPicker">
      <view class="picker-panel tag-panel" @click.stop>
        <view class="tag-panel-header">
          <text class="tag-panel-title">{{ t('home.more_tags') }}</text>
          <view class="tag-panel-reset" @click="resetTags">
            <text class="reset-text">{{ t('common.reset') }}</text>
          </view>
        </view>
        <scroll-view class="tag-panel-list" scroll-y>
          <view
            v-for="tag in tags"
            :key="tag.id"
            class="tag-panel-item"
            :class="{ checked: pendingTagIds.includes(tag.id) }"
            @click="togglePendingTag(tag.id)"
          >
            <view class="tag-checkbox" :class="{ on: pendingTagIds.includes(tag.id) }">
              <text v-if="pendingTagIds.includes(tag.id)" class="checkbox-icon">✓</text>
            </view>
            <text class="tag-panel-label">{{ localText(tag.nameI18n || tag.name_i18n || tag.name) }}</text>
          </view>
        </scroll-view>
        <view class="tag-panel-actions">
          <button class="tag-btn-cancel" @click="cancelTagPicker">{{ t('common.cancel') }}</button>
          <button class="tag-btn-confirm" @click="confirmTagPicker">{{ t('common.confirm') }}</button>
        </view>
      </view>
    </view>

    <!-- 视频单集列表 -->
    <scroll-view class="episode-scroll" scroll-y @scrolltolower="onLoadMore">
      <view v-if="!loading && episodeList.length === 0" class="empty-wrap">
        <view class="empty-icon-wrap">
          <text class="empty-icon">▶</text>
        </view>
        <text class="empty-text">{{ t('common.empty') }}</text>
      </view>

      <view
        v-for="item in episodeList"
        :key="item.id"
        class="episode-card"
      >
        <view class="episode-info">
          <text class="episode-name">{{ localText(item.name) }}</text>
          <view class="episode-tags" v-if="item.tags && item.tags.length > 0">
            <text
              v-for="tag in item.tags"
              :key="tag.id"
              class="episode-tag"
            >{{ localText(tag.nameI18n || tag.name_i18n || tag.name) }}</text>
          </view>
          <text class="episode-duration" v-if="item.duration">
            {{ formatDuration(item.duration) }}
          </text>
        </view>
        <button class="play-btn" size="mini" @click="goPlay(item.id)">{{ t('videoDetail.play_now') }}</button>
      </view>

      <view v-if="loading" class="load-status">{{ t('common.loading') }}</view>
      <view v-else-if="finished && episodeList.length > 0" class="load-status">{{ t('common.no_more') }}</view>

      <view class="bottom-safe"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { getVideoCategoryList, getVideoEpisodeListByTag, getVideoTagList } from '@/api/learning.js'

const langStore = useLangStore()
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const t = (key, params) => {
  const text = i18nT(key, locale.value)
  if (params && text && text !== key) {
    return text.replace(/\{(\w+)\}/g, (_, k) => (params[k] != null ? params[k] : ''))
  }
  return text !== key ? text : key
}
const localText = (obj) => i18nLocalText(obj, locale.value)

// 分类
const categories = ref([])
const currentCategoryId = ref(0)
const showCatPicker = ref(false)

const currentCategoryLabel = computed(() => {
  if (currentCategoryId.value === 0) return t('home.all_category')
  const cat = categories.value.find(c => c.id === currentCategoryId.value)
  return cat ? localText(cat.name) : t('home.all_category')
})

// 标签
const tags = ref([])
const selectedTagIds = ref([])
const pendingTagIds = ref([])
const showTagPicker = ref(false)

const tagTriggerLabel = computed(() => {
  if (selectedTagIds.value.length === 0) return t('home.all_category')
  const lang = locale.value
  const selectedTags = tags.value.filter(t => selectedTagIds.value.includes(t.id))
  if (selectedTags.length === 1) {
    return localText(selectedTags[0].nameI18n || selectedTags[0].name_i18n || selectedTags[0].name)
  }
  return `${t('tag.selected_count', { count: selectedTagIds.value.length })}`
})

// 列表
const episodeList = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const finished = ref(false)

const goBack = () => {
  uni.navigateBack({
    fail: () => {
      uni.switchTab({ url: '/pages/learning/home' })
    }
  })
}

const loadCategories = async () => {
  try {
    const res = await getVideoCategoryList({ page: 1, pageSize: 100 })
    const list = (res && res.data && res.data.list) ? res.data.list : []
    categories.value = Array.isArray(list) ? list.map(item => ({ ...item, id: Number(item.id || item.ID || 0) })).filter(item => item.id > 0) : []
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

const loadTags = async () => {
  try {
    const res = await getVideoTagList()
    const list = Array.isArray(res.data) ? res.data : (res && res.data && res.data.list) ? res.data.list : []
    tags.value = Array.isArray(list) ? list.map(item => ({ ...item, id: Number(item.id || item.ID || 0) })).filter(item => item.id > 0) : []
  } catch (e) {
    console.error('加载标签失败', e)
  }
}

const loadEpisodes = async (append = false) => {
  if (loading.value) return
  loading.value = true
  try {
    const params = { page: page.value, pageSize }
    if (currentCategoryId.value > 0) {
      params.categoryId = currentCategoryId.value
    }
    if (selectedTagIds.value.length > 0) {
      params.tagIds = selectedTagIds.value.slice()
    }
    const res = await getVideoEpisodeListByTag(params)
    const list = (res && res.data && res.data.list) ? res.data.list : (res && res.list) ? res.list : []
    const total = (res && res.data && res.data.total != null) ? res.data.total : (res && res.total != null) ? res.total : list.length

    const normalized = Array.isArray(list) ? list.map(item => ({
      ...item,
      id: Number(item.id || item.ID || 0),
      tags: Array.isArray(item.tags) ? item.tags.map(t => ({ ...t, id: Number(t.id || t.ID || 0) })) : []
    })) : []

    episodeList.value = append ? [...episodeList.value, ...normalized] : normalized
    if (episodeList.value.length >= total) finished.value = true
  } catch (e) {
    console.error('加载单集列表失败', e)
    if (!append) episodeList.value = []
  } finally {
    loading.value = false
  }
}

// 分类选择
const selectCategory = (categoryId) => {
  currentCategoryId.value = categoryId
  showCatPicker.value = false
  resetAndLoad()
}

// 标签选择
const togglePendingTag = (tagId) => {
  const idx = pendingTagIds.value.indexOf(tagId)
  if (idx >= 0) {
    pendingTagIds.value.splice(idx, 1)
  } else {
    pendingTagIds.value.push(tagId)
  }
}

const confirmTagPicker = () => {
  selectedTagIds.value = [...pendingTagIds.value]
  showTagPicker.value = false
  resetAndLoad()
}

const cancelTagPicker = () => {
  pendingTagIds.value = [...selectedTagIds.value]
  showTagPicker.value = false
}

const resetTags = () => {
  pendingTagIds.value = []
}

const resetAndLoad = () => {
  page.value = 1
  finished.value = false
  episodeList.value = []
  loadEpisodes(false)
}

const onLoadMore = () => {
  if (loading.value || finished.value) return
  page.value++
  loadEpisodes(true)
}

const goPlay = (episodeId) => {
  uni.navigateTo({ url: `/pages/learning/player?id=${episodeId}` })
}

const formatDuration = (seconds) => {
  if (!seconds || seconds <= 0) return ''
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

onLoad(() => {
  loadCategories()
  loadTags()
  loadEpisodes(false)
})
</script>

<style scoped>
.tag-filter-page {
  min-height: 100vh;
  background: #F5F3FF;
  padding-bottom: 40rpx;
  box-sizing: border-box;
  overflow-x: hidden;
}

/* === 渐变 Hero 头部（与 video-detail 一致） === */
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

/* === 筛选区 === */
.filter-section {
  padding: 20rpx 24rpx 0;
  box-sizing: border-box;
  width: 100%;
}

.filter-row {
  display: flex;
  gap: 16rpx;
}

.dropdown-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18rpx 24rpx;
  border-radius: 16rpx;
  background: #FFFFFF;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 2rpx 12rpx rgba(108, 91, 255, 0.06);
  flex: 1;
  min-width: 0;
}

.dropdown-label {
  font-size: 26rpx;
  color: #1A1B3A;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.dropdown-arrow {
  font-size: 20rpx;
  color: #9B8FFF;
  margin-left: 12rpx;
  flex-shrink: 0;
}

/* === 弹层 === */
.picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(15, 23, 42, 0.36);
  z-index: 1000;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.picker-panel {
  width: 100%;
  max-width: 750rpx;
  max-height: 70vh;
  background: #FFFFFF;
  border-radius: 32rpx 32rpx 0 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.picker-list {
  padding: 16rpx 0;
  max-height: 60vh;
  overflow-y: auto;
}

.picker-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 40rpx;
  font-size: 30rpx;
  color: #1A1B3A;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.06);
}

.picker-item:active {
  background: rgba(108, 91, 255, 0.04);
}

.picker-item.active {
  color: #6D5BFF;
  font-weight: 700;
}

.picker-check {
  color: #6D5BFF;
  font-size: 32rpx;
  font-weight: 700;
}

/* === 标签弹层 === */
.tag-panel {
  max-height: 75vh;
}

.tag-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx 16rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.08);
  flex-shrink: 0;
}

.tag-panel-title {
  font-size: 32rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.tag-panel-reset {
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  background: rgba(108, 91, 255, 0.08);
}

.reset-text {
  font-size: 24rpx;
  color: #6D5BFF;
  font-weight: 600;
}

.tag-panel-list {
  flex: 1;
  max-height: 50vh;
  padding: 8rpx 0;
}

.tag-panel-item {
  display: flex;
  align-items: center;
  padding: 24rpx 32rpx;
  gap: 20rpx;
}

.tag-panel-item:active {
  background: rgba(108, 91, 255, 0.04);
}

.tag-checkbox {
  width: 40rpx;
  height: 40rpx;
  border-radius: 10rpx;
  border: 2rpx solid rgba(108, 91, 255, 0.24);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.tag-checkbox.on {
  background: #6D5BFF;
  border-color: #6D5BFF;
}

.checkbox-icon {
  color: #fff;
  font-size: 24rpx;
  font-weight: 700;
}

.tag-panel-label {
  font-size: 28rpx;
  color: #1A1B3A;
  flex: 1;
}

.tag-panel-actions {
  display: flex;
  gap: 16rpx;
  padding: 20rpx 32rpx calc(32rpx + env(safe-area-inset-bottom));
  border-top: 1rpx solid rgba(108, 91, 255, 0.08);
  flex-shrink: 0;
}

.tag-btn-cancel,
.tag-btn-confirm {
  flex: 1;
  margin: 0;
  border-radius: 999rpx;
  font-size: 28rpx;
  font-weight: 700;
  height: 84rpx;
  line-height: 84rpx;
}

.tag-btn-cancel {
  background: #fff;
  color: #6D5BFF;
  border: 2rpx solid rgba(108, 91, 255, 0.32);
}

.tag-btn-confirm {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  border: none;
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.36);
}

/* === 视频列表 === */
.episode-scroll {
  height: calc(100vh - 320rpx);
  padding: 20rpx 24rpx 0;
  box-sizing: border-box;
  width: 100%;
}

.episode-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 4rpx 16rpx rgba(108, 91, 255, 0.06);
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.episode-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.episode-name {
  font-size: 28rpx;
  color: #1A1B3A;
  font-weight: 600;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.episode-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
}

.episode-tag {
  font-size: 20rpx;
  color: #6D5BFF;
  background: rgba(108, 91, 255, 0.08);
  padding: 4rpx 14rpx;
  border-radius: 999rpx;
  font-weight: 500;
}

.episode-duration {
  font-size: 22rpx;
  color: #A9AECB;
}

.play-btn {
  margin: 0;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  border-radius: 999rpx;
  font-size: 24rpx;
  font-weight: 700;
  border: none;
  flex-shrink: 0;
  padding: 0 28rpx;
  height: 60rpx;
  line-height: 60rpx;
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

/* === 底部 === */
.bottom-safe {
  height: calc(40rpx + env(safe-area-inset-bottom));
}

.load-status {
  text-align: center;
  color: #A9AECB;
  font-size: 22rpx;
  padding: 40rpx 0 16rpx;
}
</style>