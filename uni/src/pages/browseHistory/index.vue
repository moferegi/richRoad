<template>
  <view class="history-page">
    <view class="nav">
      <view class="nav-status"></view>
      <view class="nav-content">
        <view class="nav-btn" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="nav-title">{{ $t('browseHistory') }}</text>
        <view class="nav-btn nav-action" @tap="clearAll" v-if="visibleList.length">
          <text class="nav-action-text">{{ $t('clearAll') }}</text>
        </view>
        <view class="nav-btn" v-else></view>
      </view>
    </view>

    <scroll-view scroll-y :show-scrollbar="false" class="history-scroll" @scrolltolower="onScrollToLower">
      <view class="group" v-for="(group, gIdx) in groupedList" :key="gIdx">
        <view class="group-date">{{ group.date }}</view>
        <view class="card" v-for="(item, index) in group.items" :key="index" @tap="goDetail(item)">
          <LazyImage
            class="card-image"
            :src="item.externalImagePath ? getExternalUrl(item.externalImagePath) : getUrl(item.imageUrl || item.picture || item.image || '')"
            mode="aspectFill"
          />
          <view class="card-main">
            <text class="card-title">{{ $lt(item.title || item.name || item.label || '') }}</text>
            <view class="card-bottom">
              <text class="card-price">{{ cs }}{{ formatPrice(item.price) }}</text>
              <text class="card-time">{{ formatDateTime(item.CreatedAt || item.createdAt || item.viewTime) }}</text>
            </view>
          </view>
        </view>
      </view>

      <view class="load-more-row" v-if="hasMore && visibleList.length && reachedBottom">
        <view class="load-more-btn" :class="{ disabled: loadingMore }" @tap="loadMore">
          {{ loadingMore ? $t('loading') : $t('loadMore') }}
        </view>
      </view>

      <view class="footer" v-else-if="visibleList.length && !hasMore">
        {{ $t('reachedBottom') }}
      </view>

      <view class="empty" v-if="!visibleList.length && !loading">
        <view class="empty-icon">
          <uni-icons type="eye" size="44" color="rgba(15,23,42,0.28)" />
        </view>
        <text class="empty-text">{{ $t('noHistoryData') }}</text>
      </view>

      <view style="height: 40rpx"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { clearGoodHistory, getGoodHistory } from '@/api/order'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { usePlayHistoryStore } from '@/pinia/modules/playHistory.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'

const PAGE_SIZE = 10

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const playHistoryStore = usePlayHistoryStore()

const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const sourceList = ref([])
const visibleList = ref([])
const loading = ref(false)
const loadingMore = ref(false)
const page = ref(1)
const hasMore = ref(false)
const reachedBottom = ref(false)

const sortTimestamp = (item) => {
  const raw = item?.CreatedAt || item?.createdAt || item?.viewTime || item?.updatedAt || 0
  const time = new Date(raw).getTime()
  return Number.isNaN(time) ? 0 : time
}

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  const cents = Number(priceInCents)
  if (Number.isNaN(cents)) return '0.00'
  return (cents / 100).toFixed(2)
}

const toDate = (raw) => {
  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) return null
  return date
}

const formatDate = (raw) => {
  const date = toDate(raw)
  if (!date) return '0000-00-00'
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const formatDateTime = (raw) => {
  const date = toDate(raw)
  if (!date) return '--:--'
  const hh = String(date.getHours()).padStart(2, '0')
  const mm = String(date.getMinutes()).padStart(2, '0')
  return `${formatDate(raw)} ${hh}:${mm}`
}

const groupedList = computed(() => {
  const groups = {}
  visibleList.value.forEach((item) => {
    const key = formatDate(item.CreatedAt || item.createdAt || item.viewTime)
    if (!groups[key]) groups[key] = []
    groups[key].push(item)
  })

  return Object.keys(groups).map((date) => ({
    date,
    items: groups[date],
  }))
})

const applyVisible = () => {
  const visibleCount = page.value * PAGE_SIZE
  visibleList.value = sourceList.value.slice(0, visibleCount)
  hasMore.value = sourceList.value.length > visibleList.value.length
}

const loadHistory = async () => {
  if (loading.value) return
  loading.value = true
  try {
    playHistoryStore.reload()
    const token = uni.getStorageSync('x-token')
    let list = []

    if (token) {
      const res = await getGoodHistory()
      if (res.code === 0 && Array.isArray(res.data) && res.data.length) {
        list = res.data
      }
    }

    if (!list.length) {
      list = playHistoryStore.getRecentBrowseList(100) || []
    }

    sourceList.value = [...list].sort((a, b) => sortTimestamp(b) - sortTimestamp(a))
    page.value = 1
    reachedBottom.value = false
    applyVisible()
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = () => {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value += 1
  applyVisible()
  reachedBottom.value = false
  loadingMore.value = false
}

const onScrollToLower = () => {
  if (!hasMore.value) return
  reachedBottom.value = true
}

const clearAll = () => {
  uni.showModal({
    title: $t.value('confirmClearHistory'),
    confirmColor: '#dc2626',
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: async (res) => {
      if (!res.confirm) return

      uni.showLoading({ title: $t.value('loading'), mask: true })
      try {
        const token = uni.getStorageSync('x-token')
        if (token) {
          const remoteRes = await clearGoodHistory()
          if (remoteRes?.code !== 0) {
            uni.showToast({ title: resolveApiMessage(remoteRes?.msg, 'operationFailed'), icon: 'none' })
            await loadHistory()
            return
          }
        }

        if (typeof playHistoryStore.clearBrowse === 'function') {
          playHistoryStore.clearBrowse()
        } else {
          playHistoryStore.clearAll()
        }

        sourceList.value = []
        visibleList.value = []
        hasMore.value = false
        reachedBottom.value = false
        page.value = 1
        uni.showToast({ title: $t.value('cleared'), icon: 'none' })
      } finally {
        uni.hideLoading()
      }
    }
  })
}

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const goDetail = (item) => {
  const goodID = item?.ID || item?.id || item?.goodID || item?.good_id
  if (!goodID) return
  uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${goodID}` })
}

onShow(() => {
  loadHistory()
})
</script>

<style lang="scss" scoped>
.history-page {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
  color: #0f172a;
}

.nav {
  position: sticky;
  top: 0;
  z-index: 50;
  padding: 0 20rpx 12rpx;
  backdrop-filter: blur(16rpx);
  background: rgba(244, 247, 251, 0.74);
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.06);
}

.nav-status {
  height: var(--status-bar-height, 44rpx);
}

.nav-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-btn {
  width: 120rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-action {
  justify-content: flex-end;
}

.nav-title {
  font-size: 32rpx;
  font-weight: 700;
}

.nav-action-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.62);
}

.history-scroll {
  height: calc(100vh - var(--status-bar-height, 44rpx) - 100rpx);
  padding: 16rpx 20rpx 0;
}

.group {
  margin-bottom: 14rpx;
}

.group-date {
  margin-bottom: 10rpx;
  padding-left: 6rpx;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.56);
}

.card {
  display: flex;
  gap: 12rpx;
  padding: 12rpx;
  margin-bottom: 12rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.card-image {
  width: 170rpx;
  height: 170rpx;
  flex-shrink: 0;
  border-radius: 10rpx;
  background: #e2e8f0;
}

.card-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-title {
  font-size: 28rpx;
  line-height: 38rpx;
  font-weight: 600;
  color: #0f172a;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10rpx;
}

.card-price {
  font-size: 30rpx;
  font-weight: 700;
  color: #dc2626;
}

.card-time {
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.5);
}

.load-more-row {
  margin-top: 12rpx;
  display: flex;
  justify-content: center;
}

.load-more-btn {
  width: 280rpx;
  height: 66rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.95);
  color: #0f172a;
  font-size: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.load-more-btn.disabled {
  opacity: 0.65;
}

.footer {
  margin-top: 12rpx;
  text-align: center;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.46);
}

.empty {
  margin-top: 170rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-icon {
  margin-bottom: 14rpx;
}

.empty-text {
  font-size: 26rpx;
  color: rgba(15, 23, 42, 0.5);
}
</style>
