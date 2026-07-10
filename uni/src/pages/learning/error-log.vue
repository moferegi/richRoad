<template>
  <view class="error-log-container">
    <view class="header-row">
      <view class="back-btn" @click="goBack">‹</view>
      <text class="header-title">{{ t('learningErrorLogTitle') }}</text>
      <view class="header-gap"></view>
    </view>

    <view class="summary-card">
      <text class="summary-label">{{ t('learningErrorLogTotal') }}</text>
      <text class="summary-value">{{ totalWrongCount }}</text>
    </view>

    <scroll-view class="list-scroll" scroll-y>
      <view v-if="!loading && errorLogList.length === 0" class="empty-wrap">
        <text class="empty-text">{{ t('learningErrorLogEmpty') }}</text>
      </view>

      <view
        v-for="(item, index) in errorLogList"
        :key="`${item.id}-${item.wordId}-${index}`"
        class="log-card"
      >
        <view class="card-head">
          <text class="word-text">{{ getWordLabel(item) }}</text>
          <text class="count-pill">{{ t('learningErrorLogWrongCount') }} {{ item.wrongCount }}</text>
        </view>

        <text class="word-desc">{{ localText(item.word?.explanation) }}</text>

        <view class="meta-row">
          <text class="meta-label">{{ t('learningErrorLogWrongIndex') }}</text>
          <text class="meta-value">{{ item.lastWrongIndex }}</text>
        </view>
        <view class="meta-row">
          <text class="meta-label">{{ t('learningErrorLogExpectedChar') }}</text>
          <text class="meta-value">{{ getSafeChar(item.lastExpectedChar) }}</text>
        </view>
        <view class="meta-row">
          <text class="meta-label">{{ t('learningErrorLogInputChar') }}</text>
          <text class="meta-value">{{ getSafeChar(item.lastInputChar) }}</text>
        </view>
        <view class="meta-row">
          <text class="meta-label">{{ t('learningErrorLogUpdatedAt') }}</text>
          <text class="meta-value">{{ formatTime(item.updatedAt) }}</text>
        </view>

        <view class="card-footer">
          <button class="practice-btn" size="mini" @click.stop="goTyping">{{ t('learningErrorLogGoPractice') }}</button>
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

const goTyping = () => {
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
  padding: 20rpx;
  box-sizing: border-box;
  background: #f5f7fb;
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
  font-weight: 700;
  color: #111827;
}

.header-gap {
  width: 64rpx;
  height: 64rpx;
}

.summary-card {
  margin-bottom: 16rpx;
  padding: 24rpx;
  border-radius: 16rpx;
  background: #ffffff;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-label {
  font-size: 26rpx;
  color: #6b7280;
}

.summary-value {
  font-size: 36rpx;
  color: #dc2626;
  font-weight: 700;
}

.list-scroll {
  height: calc(100vh - 200rpx);
}

.empty-wrap {
  padding: 120rpx 0;
  text-align: center;
}

.empty-text {
  color: #9ca3af;
  font-size: 28rpx;
}

.log-card {
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
  margin-bottom: 8rpx;
}

.word-text {
  font-size: 34rpx;
  font-weight: 700;
  color: #111827;
}

.count-pill {
  font-size: 22rpx;
  color: #ffffff;
  background: #dc2626;
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
}

.word-desc {
  display: block;
  margin-bottom: 12rpx;
  font-size: 24rpx;
  color: #6b7280;
  line-height: 1.5;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  margin-top: 8rpx;
}

.meta-label {
  color: #6b7280;
  font-size: 24rpx;
}

.meta-value {
  color: #111827;
  font-size: 24rpx;
}

.card-footer {
  margin-top: 18rpx;
  display: flex;
  justify-content: flex-end;
  gap: 10rpx;
}

.practice-btn {
  background: #eff6ff;
  color: #2563eb;
  border: 1rpx solid #bfdbfe;
  border-radius: 10rpx;
  font-size: 24rpx;
}

.load-more {
  text-align: center;
  color: #9ca3af;
  font-size: 24rpx;
  padding: 20rpx 0 40rpx;
}

.error-log-container {
  background: linear-gradient(180deg, #f8f4e7 0%, #f4efe1 100%);
}

.back-btn {
  border-radius: 18rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
  background: #fffdf8;
  color: #7c2d12;
}

.header-title { color: #7c2d12; }

.summary-card,
.log-card {
  background: rgba(255, 253, 248, 0.96);
  border: 1rpx solid rgba(20, 184, 166, 0.18);
}

.summary-label,
.meta-label,
.word-desc { color: #64748b; }

.word-text,
.meta-value { color: #1e293b; }

.count-pill {
  background: linear-gradient(120deg, #dc2626 0%, #f97316 100%);
}

.practice-btn {
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
  color: #fff;
  border: none;
  border-radius: 999rpx;
}

.delete-btn {
  background: rgba(220, 38, 38, 0.08);
  color: #dc2626;
  border: 1rpx solid rgba(220, 38, 38, 0.26);
  border-radius: 999rpx;
}

.load-more-btn {
  border: 1rpx solid rgba(20, 184, 166, 0.28);
  color: #0f766e;
  background: #fffdf8;
  border-radius: 999rpx;
  padding: 0 28rpx;
}

.load-more, .empty-text { color: #94a3b8; }
</style>