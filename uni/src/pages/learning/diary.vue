<template>
  <view class="diary-container">
    <!-- 渐变顶栏（参照 typing.vue 风格） -->
    <view class="top-bar">
      <view class="top-bar-bg"></view>
      <view class="top-bar-content">
        <view class="bar-selectors">
          <view class="selector" @click="openCategoryPicker">
            <text class="selector-text">{{ selectedCategoryName || t('diary.all_category') }}</text>
            <text class="selector-arrow">›</text>
          </view>
          <view class="selector" @click="openTagPicker">
            <text class="selector-text">{{ tagDisplayText }}</text>
            <text class="selector-arrow">›</text>
          </view>
        </view>
        <view class="bar-action" @click="resetFilters">
          <text class="action-icon">↻</text>
        </view>
      </view>
    </view>

    <!-- 日记列表 -->
    <scroll-view class="list-scroll" scroll-y @scrolltolower="loadMore">
      <view v-if="!loading && diaryList.length === 0" class="empty-wrap">
        <view class="empty-icon">
          <text class="empty-icon-text">📖</text>
        </view>
        <text class="empty-text">{{ t('diary.empty') }}</text>
      </view>

      <view v-for="item in diaryList" :key="item.ID" class="diary-card">
        <!-- 轮播区域 -->
        <swiper class="diary-swiper" :indicator-dots="true" indicator-color="rgba(255,255,255,0.4)" indicator-active-color="#6D5BFF" :current="0" :autoplay="true" :interval="4000" circular>
          <!-- 第一张：字幕文本预览 -->
          <swiper-item>
            <view class="swiper-content text-preview">
              <text class="diary-name">{{ localText(item.name) }}</text>
              <view class="text-lines">
                <text 
                  v-for="(line, idx) in getPreviewLines(item)" 
                  :key="idx"
                  class="text-line"
                >{{ line }}</text>
              </view>
            </view>
          </swiper-item>
          <!-- 第二张：日记图片 -->
          <swiper-item>
            <view class="swiper-content image-preview">
              <text class="diary-name">{{ localText(item.name) }}</text>
              <image 
                v-if="getDiaryImage(item)" 
                class="diary-image" 
                :src="getDiaryImage(item)" 
                mode="aspectFit"
              />
              <view v-else class="no-image-hint">
                <text class="no-image-icon">🖼</text>
                <text class="no-image-text">{{ t('diary.no_image') }}</text>
              </view>
            </view>
          </swiper-item>
        </swiper>

        <!-- 轮播底部操作 -->
        <view class="card-footer">
          <view class="dots-placeholder"></view>
          <view class="view-btn" @click="goToDetail(item)">
            <text class="view-btn-text">{{ t('diary.view') }}</text>
          </view>
        </view>
      </view>

      <view v-if="diaryList.length > 0" class="load-more">
        <text v-if="loading">{{ t('diary.loading') }}</text>
        <text v-else-if="noMore">{{ t('diary.no_more') }}</text>
      </view>
    </scroll-view>

    <!-- 分类选择器弹窗 -->
    <uni-popup ref="categoryPopup" type="bottom">
      <view class="sheet">
        <view class="sheet-handle"></view>
        <view class="sheet-title">{{ t('diary.select_category') }}</view>
        <scroll-view class="sheet-list" scroll-y>
          <view
            class="sheet-item"
            :class="{ active: selectedCategoryId === 0 }"
            @click="selectCategory(0, t('diary.all_category'))"
          >
            <text class="sheet-item-label">{{ t('diary.all_category') }}</text>
            <text v-if="selectedCategoryId === 0" class="sheet-item-check">✓</text>
          </view>
          <view
            v-for="cat in categories"
            :key="cat.ID"
            class="sheet-item"
            :class="{ active: selectedCategoryId === cat.ID }"
            @click="selectCategory(cat.ID, localText(cat.name))"
          >
            <text class="sheet-item-label">{{ localText(cat.name) }}</text>
            <text v-if="selectedCategoryId === cat.ID" class="sheet-item-check">✓</text>
          </view>
        </scroll-view>
        <button class="sheet-cancel" @click="closeCategoryPicker">{{ t('common.cancel') }}</button>
      </view>
    </uni-popup>

    <!-- 标签选择器弹窗 -->
    <uni-popup ref="tagPopup" type="bottom" background-color="#fff">
      <view class="tag-sheet">
        <view class="sheet-handle"></view>
        <view class="tag-sheet-head">
          <text class="tag-sheet-title">{{ t('diary.select_tags') }}</text>
          <view class="tag-sheet-reset" @click="resetTags">
            <text class="reset-text">{{ t('diary.reset') }}</text>
          </view>
        </view>

        <!-- Selected tags area (chips at top) -->
        <view class="selected-tags-area" v-if="tempTagIds.length > 0">
          <scroll-view scroll-x class="selected-tags-scroll" :show-scrollbar="false">
            <view class="selected-tags-list">
              <view
                v-for="tagId in tempTagIds"
                :key="'selected-' + tagId"
                class="selected-chip"
                @click="toggleTag(tagId)"
              >
                <text class="selected-chip-text">{{ getTagName(tagId) }}</text>
                <text class="selected-chip-close">✕</text>
              </view>
            </view>
          </scroll-view>
        </view>

        <!-- Unselected tags list -->
        <scroll-view class="tag-sheet-list" scroll-y>
          <view
            v-for="tag in unselectedTags"
            :key="tag.ID"
            class="tag-option"
            @click="toggleTag(tag.ID)"
          >
            <text class="tag-option-text">{{ localText(tag.nameI18n) || tag.name }}</text>
            <text class="tag-option-add">+</text>
          </view>
          <view v-if="unselectedTags.length === 0" class="tag-empty">
            <text class="tag-empty-text">{{ t('common.all_selected') }}</text>
          </view>
        </scroll-view>

        <view class="tag-sheet-actions">
          <button class="tag-btn-cancel" @click="closeTagPicker">{{ t('common.cancel') }}</button>
          <button class="tag-btn-confirm" @click="confirmTags">{{ t('diary.confirm') }}</button>
        </view>
      </view>
    </uni-popup>

    <custom-tab-bar />
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import CustomTabBar from '@/components/custom-tab-bar/custom-tab-bar.vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { localText as i18nLocalText, t as i18nT } from '@/utils/i18n.js'
import { getExternalUrl, initExternalDomain } from '@/utils/url.js'
import { getDiaryCategoryList, getDiaryListByTag, getDiaryTagPublic, getDiarySentenceList } from '@/api/learning.js'

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

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

// 筛选
const selectedCategoryId = ref(0)
const selectedCategoryName = ref('')
const selectedTagIds = ref([])
const tempTagIds = ref([])

const tagDisplayText = computed(() => {
  if (selectedTagIds.value.length === 0) return t('diary.all_tags')
  return t('diary.selected_tags').replace('{0}', selectedTagIds.value.length)
})

// 分类列表
const categories = ref([])
const loadCategories = async () => {
  const res = await getDiaryCategoryList({ page: 1, pageSize: 100 })
  if (res.code === 0 && res.data) {
    categories.value = res.data.list || []
  }
}

// 标签列表
const tags = ref([])
const loadTags = async () => {
  const res = await getDiaryTagPublic()
  if (res.code === 0) {
    tags.value = res.data || []
  }
}

const unselectedTags = computed(() => {
  return tags.value.filter(tag => !tempTagIds.value.includes(tag.ID))
})

const getTagName = (tagId) => {
  const tag = tags.value.find(t => t.ID === tagId)
  if (!tag) return ''
  return localText(tag.nameI18n) || tag.name
}

// 日记列表
const diaryList = ref([])
const page = ref(1)
const pageSize = 10
const loading = ref(false)
const noMore = ref(false)

const fetchDiaryList = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) return
  loading.value = true

  const nextPage = isLoadMore ? page.value + 1 : 1
  const params = {
    page: nextPage,
    pageSize
  }
  if (selectedCategoryId.value > 0) {
    params.categoryId = selectedCategoryId.value
  }
  if (selectedTagIds.value.length > 0) {
    params.tagIds = selectedTagIds.value
  }

  const res = await getDiaryListByTag(params)
  if (res.code === 0 && res.data) {
    const list = Array.isArray(res.data.list) ? res.data.list : []
    const total = Number(res.data.total || 0)
    page.value = nextPage
    if (isLoadMore) {
      diaryList.value = [...diaryList.value, ...list]
    } else {
      diaryList.value = list
    }
    noMore.value = diaryList.value.length >= total || list.length < pageSize
  }
  loading.value = false
}

const loadMore = () => {
  fetchDiaryList(true)
}

// 字幕预览 - 正常显示语句，去除 <W> 标签
const sentenceCache = ref({})
const stripHighlightTags = (text) => {
  if (!text) return ''
  return text.replace(/<\/?w[^>]*>/g, '')
}

const formatSecondsPreview = (secs) => {
  const m = Math.floor(secs / 60).toString().padStart(2, '0')
  const s = Math.floor(secs % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}

const getPreviewLines = (item) => {
  const id = getEntityId(item)
  if (sentenceCache.value[id]) return sentenceCache.value[id]
  return []
}

const loadPreviewSentences = async () => {
  for (const item of diaryList.value) {
    const id = getEntityId(item)
    if (sentenceCache.value[id]) continue
    try {
      const res = await getDiarySentenceList(id)
      if (res.code === 0 && res.data) {
        const lines = (res.data || [])
          .slice(0, 10)
          .map(s => {
            if (s.locked) {
              // 锁定句子显示时间区间 + 锁标识
              const start = formatSecondsPreview(s.startTime || 0)
              const end = formatSecondsPreview(s.endTime || 0)
              return `🔒 ${start} - ${end}`
            }
            const showTranslate = item.showTranslate !== false
            const showEnglish = item.showEnglish !== false
            const parts = []
            if (showEnglish) parts.push(stripHighlightTags(s.english || ''))
            if (showTranslate) parts.push(localText(s.translate))
            return parts.join('\n')
          }).filter(Boolean)
        sentenceCache.value = { ...sentenceCache.value, [id]: lines }
      }
    } catch (e) {
      sentenceCache.value = { ...sentenceCache.value, [id]: [] }
    }
  }
}

// 日记图片 - 使用 externalLinkDomain 拼接完整URL
const getDiaryImage = (item) => {
  if (!item.imageI18n) return ''
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  let images = {}
  try {
    images = typeof item.imageI18n === 'object' ? item.imageI18n : JSON.parse(item.imageI18n || '{}')
  } catch (e) {}
  const url = images[locale] || images['en'] || images['zh'] || ''
  return getExternalUrl(url)
}

// 分类选择
const categoryPopup = ref(null)
const openCategoryPicker = () => categoryPopup.value?.open()
const closeCategoryPicker = () => categoryPopup.value?.close()
const selectCategory = (id, name) => {
  selectedCategoryId.value = id
  selectedCategoryName.value = name
  closeCategoryPicker()
  fetchDiaryList(false)
}

// 标签选择
const tagPopup = ref(null)
const openTagPicker = () => {
  tempTagIds.value = [...selectedTagIds.value]
  tagPopup.value?.open()
}
const closeTagPicker = () => tagPopup.value?.close()
const toggleTag = (id) => {
  const idx = tempTagIds.value.indexOf(id)
  if (idx >= 0) {
    tempTagIds.value.splice(idx, 1)
  } else {
    tempTagIds.value.push(id)
  }
}
const selectAllTags = () => {
  tempTagIds.value = []
}
const resetTags = () => {
  tempTagIds.value = []
}
const confirmTags = () => {
  selectedTagIds.value = [...tempTagIds.value]
  closeTagPicker()
  fetchDiaryList(false)
}

// 重置
const resetFilters = () => {
  selectedCategoryId.value = 0
  selectedCategoryName.value = ''
  selectedTagIds.value = []
  fetchDiaryList(false)
}

// 跳转详情
const goToDetail = (item) => {
  const id = getEntityId(item)
  uni.navigateTo({ url: `/pages/learning/diary-detail?id=${id}` })
}

onShow(() => {
  initExternalDomain()
  loadCategories()
  loadTags()
  fetchDiaryList(false).then(() => loadPreviewSentences())
})
</script>

<style scoped>
.diary-container {
  min-height: 100vh;
  background: #F5F3FF;
  overflow-x: hidden;
  padding-bottom: env(safe-area-inset-bottom);
}

/* === 渐变顶栏 === */
.top-bar {
  position: relative;
  overflow: hidden;
}

.top-bar-bg {
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 100%;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
}

.top-bar-content {
  position: relative;
  z-index: 2;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: calc(var(--status-bar-height, 0px) + 20rpx) 32rpx 20rpx;
}

.bar-selectors {
  display: flex;
  gap: 16rpx;
  align-items: center;
  flex: 1;
}

.selector {
  display: flex;
  align-items: center;
  gap: 6rpx;
  padding: 10rpx 24rpx;
  border-radius: 999rpx;
  background: rgba(255, 255, 255, 0.22);
  backdrop-filter: blur(8rpx);
}

.selector:active {
  background: rgba(255, 255, 255, 0.32);
}

.selector-text {
  font-size: 27rpx;
  font-weight: 600;
  color: #fff;
  max-width: 200rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.selector-arrow {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  transform: rotate(90deg);
}

.bar-action {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.22);
  flex-shrink: 0;
  margin-left: 16rpx;
}

.bar-action:active {
  background: rgba(255, 255, 255, 0.34);
}

.action-icon {
  font-size: 32rpx;
  color: #fff;
}

/* === 列表 === */
.list-scroll {
  height: calc(100vh - 280rpx - 88rpx - env(safe-area-inset-bottom));
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
}

.empty-text {
  color: #A9AECB;
  font-size: 28rpx;
}

/* === 日记卡片 === */
.diary-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  margin-bottom: 20rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.08);
  overflow: hidden;
}

.diary-swiper {
  width: 100%;
  height: 600rpx;
}

.swiper-content {
  width: 100%;
  height: 100%;
  padding: 32rpx;
  box-sizing: border-box;
  position: relative;
}

.text-preview {
  background: linear-gradient(180deg, #FAF9FF 0%, #F5F3FF 100%);
}

.image-preview {
  background: #FAF9FF;
  display: flex;
  flex-direction: column;
}

.diary-name {
  font-size: 30rpx;
  font-weight: 700;
  color: #1A1B3A;
  margin-bottom: 20rpx;
  display: block;
}

.text-lines {
  flex: 1;
  overflow: hidden;
}

.text-line {
  display: block;
  font-size: 24rpx;
  color: #4A4B6A;
  line-height: 1.8;
  margin-bottom: 4rpx;
}

.diary-image {
  flex: 1;
  width: 100%;
  border-radius: 16rpx;
}

.no-image-hint {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
}

.no-image-icon {
  font-size: 64rpx;
}

.no-image-text {
  font-size: 24rpx;
  color: #A9AECB;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 32rpx;
}

.dots-placeholder {
  width: 80rpx;
}

.view-btn {
  padding: 12rpx 32rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  border-radius: 999rpx;
}

.view-btn:active {
  transform: scale(0.95);
}

.view-btn-text {
  font-size: 24rpx;
  color: #fff;
  font-weight: 600;
}

.load-more {
  text-align: center;
  color: #A9AECB;
  font-size: 24rpx;
  padding: 20rpx 0 40rpx;
}

/* === 底部弹窗 === */
.sheet {
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  box-shadow: 0 -12rpx 36rpx rgba(108, 91, 255, 0.14);
  max-height: 72vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sheet-handle {
  width: 60rpx;
  height: 6rpx;
  border-radius: 3rpx;
  background: linear-gradient(90deg, #6D5BFF, #9B8FFF);
  margin: 16rpx auto 0;
  flex-shrink: 0;
}

.sheet-title {
  text-align: center;
  font-size: 32rpx;
  font-weight: 700;
  color: #1A1B3A;
  margin: 20rpx 0;
  padding: 0 32rpx;
  flex-shrink: 0;
}

.sheet-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 0 24rpx;
  box-sizing: border-box;
  width: 100%;
}

.sheet-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 96rpx;
  padding: 0 28rpx;
  border-radius: 20rpx;
  margin-bottom: 10rpx;
  font-size: 28rpx;
  color: #1A1B3A;
  background: #F7F5FF;
  box-sizing: border-box;
  transition: all 0.2s ease;
}

.sheet-item:active {
  background: rgba(108, 91, 255, 0.14);
  transform: scale(0.98);
}

.sheet-item.active {
  background: linear-gradient(135deg, rgba(109, 91, 255, 0.12) 0%, rgba(155, 143, 255, 0.12) 100%);
  border: 1.5rpx solid rgba(108, 91, 255, 0.4);
}

.sheet-item-label { font-weight: 500; }

.sheet-item.active .sheet-item-label {
  color: #6D5BFF;
  font-weight: 700;
}

.sheet-item-check {
  color: #6D5BFF;
  font-size: 32rpx;
  font-weight: 700;
  width: 36rpx;
  height: 36rpx;
  line-height: 36rpx;
  text-align: center;
  border-radius: 50%;
  background: rgba(108, 91, 255, 0.12);
}

.sheet-cancel {
  flex-shrink: 0;
  margin: 20rpx 32rpx;
  margin-bottom: calc(24rpx + 140rpx + env(safe-area-inset-bottom));
  border-radius: 999rpx;
  border: none;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  font-size: 28rpx;
  font-weight: 600;
  height: 80rpx;
  line-height: 80rpx;
  box-shadow: 0 6rpx 18rpx rgba(108, 91, 255, 0.32);
}

.sheet-cancel:active {
  transform: scale(0.97);
  opacity: 0.92;
}

/* === 标签弹窗 === */
.tag-sheet {
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  max-height: 75vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.tag-sheet-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx 16rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.08);
  flex-shrink: 0;
}

.tag-sheet-title {
  font-size: 32rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.tag-sheet-reset {
  padding: 8rpx 20rpx;
  border-radius: 999rpx;
  background: rgba(108, 91, 255, 0.08);
}

.reset-text {
  font-size: 26rpx;
  color: #6D5BFF;
  font-weight: 600;
}

.selected-tags-area {
  padding: 16rpx 32rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.08);
  flex-shrink: 0;
}

.selected-tags-scroll { white-space: nowrap; }

.selected-tags-list {
  display: inline-flex;
  gap: 12rpx;
  padding: 4rpx 0;
}

.selected-chip {
  display: inline-flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 20rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  box-shadow: 0 4rpx 12rpx rgba(108, 91, 255, 0.3);
  flex-shrink: 0;
}

.selected-chip-text {
  font-size: 26rpx;
  font-weight: 600;
  color: #fff;
}

.selected-chip-close {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.8);
  width: 32rpx;
  height: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
}

.tag-sheet-list {
  flex: 1;
  min-height: 0;
  max-height: 40vh;
  padding: 8rpx 0;
}

.tag-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 32rpx;
  border-bottom: 1rpx solid rgba(108, 91, 255, 0.06);
}

.tag-option:active { background: rgba(108, 91, 255, 0.04); }

.tag-option-text {
  font-size: 28rpx;
  color: #1A1B3A;
}

.tag-option-add {
  font-size: 32rpx;
  color: #6D5BFF;
  font-weight: 700;
}

.tag-empty {
  padding: 48rpx 32rpx;
  text-align: center;
}

.tag-empty-text {
  font-size: 26rpx;
  color: #A9AECB;
}

.tag-sheet-actions {
  display: flex;
  gap: 16rpx;
  padding: 20rpx 32rpx calc(32rpx + env(safe-area-inset-bottom) + 100rpx);
  border-top: 1rpx solid rgba(108, 91, 255, 0.08);
  flex-shrink: 0;
}

.tag-btn-cancel, .tag-btn-confirm {
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
</style>