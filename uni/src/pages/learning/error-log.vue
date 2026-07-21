<template>
  <view class="error-log-container">
    <!-- 渐变 Hero 头部 -->
    <view class="page-hero">
      <view class="hero-decor-circle decor-1"></view>
      <view class="hero-decor-circle decor-2"></view>
      <view class="hero-content">
        <view class="hero-back-btn" @click="goBack">
          <text class="hero-back-icon">‹</text>
        </view>
        <text class="hero-title">{{ t('learningErrorLogTitle') }}</text>
        <view class="hero-gap"></view>
      </view>
    </view>

    <!-- 摘要浮卡 -->
    <view class="summary-card">
      <view class="summary-left">
        <text class="summary-label">{{ t('learningErrorLogTotal') }}</text>
      </view>
      <text class="summary-value">{{ totalWrongCount }}</text>
    </view>

    <scroll-view class="list-scroll" scroll-y>
      <view v-if="!loading && errorLogList.length === 0" class="empty-wrap">
        <view class="empty-icon">
          <text class="empty-icon-text">✎</text>
        </view>
        <text class="empty-text">{{ t('learningErrorLogEmpty') }}</text>
      </view>

      <view
        v-for="(item, index) in errorLogList"
        :key="`${item.id}-${item.wordId}-${index}`"
        class="log-card"
      >
        <!-- 头部：图标 + 单词 + 错误次数胶囊 -->
        <view class="card-head">
          <view class="card-head-left">
            <view class="word-icon-wrap">
              <text class="word-icon">✎</text>
            </view>
            <view class="word-info">
              <text class="word-text">{{ getWordLabel(item) }}</text>
              <text class="word-desc">{{ localText(item.word?.explanation) }}</text>
            </view>
          </view>
          <view class="count-pill">
            <text class="count-pill-text">{{ item.wrongCount }}</text>
          </view>
        </view>

        <!-- Meta 网格 -->
        <view class="meta-grid">
          <view class="meta-cell">
            <text class="meta-label">{{ t('learningErrorLogWrongIndex') }}</text>
            <text class="meta-value">{{ item.lastWrongIndex }}</text>
          </view>
          <view class="meta-cell">
            <text class="meta-label">{{ t('learningErrorLogExpectedChar') }}</text>
            <text class="meta-value">{{ getSafeChar(item.lastExpectedChar) }}</text>
          </view>
          <view class="meta-cell">
            <text class="meta-label">{{ t('learningErrorLogInputChar') }}</text>
            <text class="meta-value">{{ getSafeChar(item.lastInputChar) }}</text>
          </view>
          <view class="meta-cell">
            <text class="meta-label">{{ t('learningErrorLogUpdatedAt') }}</text>
            <text class="meta-value">{{ formatTime(item.updatedAt) }}</text>
          </view>
        </view>

        <!-- 底部按钮 -->
        <view class="card-footer">
          <button class="practice-btn" size="mini" @click.stop="goTyping(item)">{{ t('learningErrorLogGoPractice') }} →</button>
          <button class="delete-btn" size="mini" @click.stop="removeLog(item)">{{ tt('learningErrorLogDelete', 'delete') }}</button>
        </view>
      </view>

      <view v-if="errorLogList.length > 0" class="load-more">
        <text v-if="loading">{{ t('learningErrorLogLoading') }}</text>
        <text v-else-if="noMore">{{ t('learningErrorLogNoMore') }}</text>
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
import { deleteWordErrorLog, getWordErrorLogList } from '@/api/learning.js'

const langStore = useLangStore()
const page = ref(1)
const pageSize = 10
const totalWrongCount = ref(0)
const loading = ref(false)
const noMore = ref(false)
const errorLogList = ref([])

const t = (key) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  const text = i18nT(key, locale)
  if (text && text !== key) return text
  return key
}

const tt = (...keys) => {
  for (const key of keys) {
    const text = t(key)
    if (text && text !== key) {
      return text
    }
  }
  const last = keys[keys.length - 1]
  return typeof last === 'string' ? last : ''
}

const localText = (value) => {
  const locale = langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  return i18nLocalText(value, locale)
}

const normalizeLogItem = (item) => {
  const word = item?.word || item?.Word || {}
  return {
    id: Number(item?.id || item?.ID || 0),
    wordId: Number(item?.wordId || item?.WordID || item?.wordID || 0),
    wrongCount: Number(item?.wrongCount || item?.WrongCount || 0),
    lastWrongIndex: Number(item?.lastWrongIndex ?? item?.LastWrongIndex ?? -1),
    lastExpectedChar: String(item?.lastExpectedChar || item?.LastExpectedChar || ''),
    lastInputChar: String(item?.lastInputChar || item?.LastInputChar || ''),
    updatedAt: item?.updatedAt || item?.UpdatedAt || item?.createdAt || item?.CreatedAt || '',
    word
  }
}

const getWordLabel = (item) => {
  const text = String(item?.word?.word || item?.word?.Word || '').trim()
  return text || t('learningErrorLogUnknownWord')
}

const getSafeChar = (value) => {
  const text = String(value || '').trim()
  return text || t('learningErrorLogUnknownChar')
}

const formatTime = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${d} ${hh}:${mm}`
}

const fetchErrorLogList = async (isLoadMore = false) => {
  if (loading.value || (isLoadMore && noMore.value)) {
    return
  }
  loading.value = true

  const nextPage = isLoadMore ? page.value + 1 : 1
  try {
    const res = await getWordErrorLogList({ page: nextPage, pageSize })
    if (res.code !== 0) {
      uni.showToast({ title: resolveApiMessage(res.msg, 'operationFailed'), icon: 'none' })
      return
    }

    const data = res.data || {}
    const list = Array.isArray(data.list) ? data.list.map(normalizeLogItem) : []
    totalWrongCount.value = Number(data.total || 0)
    page.value = nextPage

    if (isLoadMore) {
      errorLogList.value = [...errorLogList.value, ...list]
    } else {
      errorLogList.value = list
    }

    noMore.value = errorLogList.value.length >= totalWrongCount.value || list.length < pageSize
  } catch (error) {
    uni.showToast({ title: t('networkError'), icon: 'none' })
  } finally {
    loading.value = false
  }
}

const loadMore = () => {
  fetchErrorLogList(true)
}

const removeLog = (item) => {
  const wordId = Number(item?.wordId || 0)
  if (!wordId) {
    uni.showToast({ title: t('learningErrorLogDeleteFailed'), icon: 'none' })
    return
  }

  uni.showModal({
    title: t('learningErrorLogTitle'),
    content: tt('learningErrorLogDeleteConfirm', 'confirmDelete'),
    cancelText: tt('common.cancel', 'cancel'),
    confirmText: tt('common.confirm', 'confirm'),
    success: async (res) => {
      if (!res.confirm) {
        return
      }
      const resp = await deleteWordErrorLog(wordId)
      if (resp.code !== 0) {
        uni.showToast({ title: t('learningErrorLogDeleteFailed'), icon: 'none' })
        return
      }
      uni.showToast({ title: t('learningErrorLogDeleteSuccess'), icon: 'success' })
      errorLogList.value = errorLogList.value.filter((row) => Number(row.wordId || 0) !== wordId)
      totalWrongCount.value = Math.max(0, totalWrongCount.value - 1)
      if (errorLogList.value.length === 0) {
        noMore.value = true
      }
    }
  })
}

const goTyping = (item) => {
  const wordId = Number(item?.wordId || 0)
  if (wordId) {
    uni.setStorageSync('typing-pending-word-id', wordId)
  }
  uni.switchTab({ url: '/pages/learning/typing' })
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
  fetchErrorLogList(false)
})
</script>

<style scoped>
.error-log-container {
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

/* === 摘要浮卡 === */
.summary-card {
  margin: -28rpx 24rpx 16rpx;
  position: relative;
  z-index: 5;
  padding: 28rpx 32rpx;
  border-radius: 24rpx;
  background: #FFFFFF;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.12);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-left {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.summary-label {
  font-size: 26rpx;
  color: #6B6F8D;
  font-weight: 500;
}

.summary-value {
  font-size: 40rpx;
  color: #EF4444;
  font-weight: 800;
}

/* === 列表 === */
.list-scroll {
  height: calc(100vh - 340rpx);
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

.log-card {
  background: #FFFFFF;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  border: 1rpx solid rgba(108, 91, 255, 0.16);
  box-shadow: 0 8rpx 24rpx rgba(108, 91, 255, 0.08);
}

/* === 头部：图标 + 单词 + 错误次数胶囊 === */
.card-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16rpx;
}

.card-head-left {
  flex: 1;
  display: flex;
  gap: 16rpx;
  min-width: 0;
}

.word-icon-wrap {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.word-icon {
  font-size: 32rpx;
  color: #EF4444;
  font-weight: 700;
}

.word-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
  min-width: 0;
  padding-top: 4rpx;
}

.word-text {
  font-size: 34rpx;
  font-weight: 800;
  color: #1A1B3A;
}

.word-desc {
  font-size: 24rpx;
  color: #6B6F8D;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.count-pill {
  min-width: 56rpx;
  padding: 8rpx 16rpx;
  border-radius: 999rpx;
  background: linear-gradient(135deg, #EF4444 0%, #F97316 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.count-pill-text {
  font-size: 26rpx;
  color: #fff;
  font-weight: 800;
}

/* === Meta 网格 === */
.meta-grid {
  display: flex;
  flex-wrap: wrap;
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid rgba(108, 91, 255, 0.08);
  gap: 0;
}

.meta-cell {
  width: 50%;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
  padding: 8rpx 0;
}

.meta-label {
  color: #6B6F8D;
  font-size: 22rpx;
}

.meta-value {
  color: #1A1B3A;
  font-size: 26rpx;
  font-weight: 600;
}

/* === 底部按钮 === */
.card-footer {
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid rgba(108, 91, 255, 0.08);
  display: flex;
  justify-content: flex-end;
  gap: 12rpx;
}

.practice-btn {
  background: linear-gradient(135deg, #6D5BFF 0%, #9B8FFF 100%);
  color: #fff;
  border: none;
  border-radius: 999rpx;
  font-size: 24rpx;
  font-weight: 600;
  box-shadow: 0 4rpx 12rpx rgba(108, 91, 255, 0.3);
  padding: 0 28rpx;
}

.practice-btn:active {
  transform: scale(0.95);
}

.delete-btn {
  background: rgba(239, 68, 68, 0.1);
  color: #EF4444;
  border: 1rpx solid rgba(239, 68, 68, 0.22);
  border-radius: 999rpx;
  font-size: 24rpx;
  font-weight: 600;
  padding: 0 28rpx;
}

.delete-btn:active {
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