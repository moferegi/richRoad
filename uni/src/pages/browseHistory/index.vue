<template>
  <view class="nf-history">
    <view class="nf-history-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('browseHistory') }}</text>
        <view class="nf-navbar-action" @tap="clearAll" v-if="historyList.length">
          <text class="nf-navbar-action-text">{{ $t('clearAll') }}</text>
        </view>
      </view>
    </view>

    <!-- 商品列表 -->
    <scroll-view
      scroll-y
      :show-scrollbar="false"
      class="nf-history-scroll"
      @scrolltolower="debouncedLower"
    >
      <!-- 按日期分组 -->
      <view v-for="(group, gIdx) in groupedList" :key="gIdx" class="nf-history-group">
        <view class="nf-history-date">
          <text class="nf-history-date-text">{{ group.date }}</text>
        </view>
        <view
          class="nf-goods-card"
          v-for="(item, index) in group.items"
          :key="index"
          @tap="goDetail(item)"
        >
          <view class="nf-goods-img-wrap">
            <image class="nf-goods-img" :src="getUrl(item.imageUrl)" mode="aspectFill" />
            <view class="nf-goods-img-overlay"></view>
          </view>
          <view class="nf-goods-info">
            <text class="nf-goods-title">{{ $lt(item.title) }}</text>
            <view class="nf-goods-bottom">
              <text class="nf-price">¥{{ formatPrice(item.price) }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 底部加载状态 -->
      <view class="nf-history-footer" v-if="historyList.length > 0">
        <text class="nf-history-footer-text">{{ isBottom ? $t('reachedBottom') : $t('loading') }}</text>
      </view>

      <!-- 空状态 -->
      <view class="nf-history-empty" v-if="!historyList.length && !loading">
        <view class="nf-history-empty-icon">
          <uni-icons type="eye" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-history-empty-text">{{ $t('noHistoryData') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getGoodHistory } from '@/api/order'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { usePlayHistoryStore } from '@/pinia/modules/playHistory.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const playHistoryStore = usePlayHistoryStore()
const historyList = ref([])
const loading = ref(false)
const isBottom = ref(true)

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  const cents = parseInt(priceInCents)
  if (isNaN(cents)) return '0.00'
  return (cents / 100).toFixed(2)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const target = new Date(d.getFullYear(), d.getMonth(), d.getDate())
  const diff = (today - target) / 86400000

  if (diff === 0) return $t.value('today')
  if (diff === 1) return $t.value('yesterday')
  return `${d.getMonth() + 1}/${d.getDate()}`
}

const groupedList = computed(() => {
  const groups = {}
  historyList.value.forEach(item => {
    const date = formatDate(item.CreatedAt || item.createdAt || item.viewTime)
    if (!groups[date]) groups[date] = []
    groups[date].push(item)
  })
  return Object.keys(groups).map(date => ({ date, items: groups[date] }))
})

const loadHistory = async () => {
  loading.value = true
  try {
    const token = uni.getStorageSync('x-token')
    if (token) {
      const res = await getGoodHistory()
      if (res.code === 0 && res.data && res.data.length) {
        historyList.value = res.data
        return
      }
    }
    // 未登录或无远程数据，使用本地历史
    const localList = playHistoryStore.getRecentList(50)
    historyList.value = localList || []
  } finally {
    loading.value = false
  }
}

onShow(() => {
  loadHistory()
})

const debounce = (func, delay) => {
  let timer
  return function (...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => func.apply(this, args), delay)
  }
}
const debouncedLower = debounce(() => {}, 300)

const clearAll = () => {
  uni.showModal({
    title: $t.value('confirmClearHistory'),
    confirmColor: '#e50914',
    success: (res) => {
      if (res.confirm) {
        playHistoryStore.clearAll()
        historyList.value = []
        uni.showToast({ title: $t.value('cleared'), icon: 'none' })
      }
    }
  })
}

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const goDetail = (item) => {
  if (item && item.ID) {
    uni.navigateTo({ url: `/pages/player/index?id=${item.ID}` })
  }
}
</script>

<style lang="scss" scoped>
.nf-history {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-history-bg {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(180deg, #1a1a2e 0%, #141414 100%);
  z-index: 0;
}
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(20rpx);
}
.nf-navbar-status {
  height: var(--status-bar-height, 44rpx);
}
.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 24rpx;
}
.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-navbar-title {
  font-size: 34rpx;
  font-weight: bold;
  color: #fff;
}
.nf-navbar-action {
  width: 64rpx;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}
.nf-navbar-action-text {
  font-size: 24rpx;
  color: rgba(255,255,255,0.6);
}
.nf-history-scroll {
  position: relative;
  z-index: 1;
  height: 100vh;
  padding-top: calc(var(--status-bar-height, 44rpx) + 88rpx);
}
.nf-history-group {
  padding: 0 24rpx;
}
.nf-history-date {
  padding: 20rpx 0 12rpx;
}
.nf-history-date-text {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
  font-weight: bold;
}
.nf-goods-card {
  display: flex;
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 16rpx;
}
.nf-goods-img-wrap {
  position: relative;
  width: 200rpx;
  min-height: 200rpx;
  flex-shrink: 0;
}
.nf-goods-img {
  width: 100%;
  height: 100%;
}
.nf-goods-img-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(90deg, transparent 60%, rgba(0,0,0,0.3) 100%);
}
.nf-goods-info {
  flex: 1;
  padding: 16rpx 20rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.nf-goods-title {
  font-size: 28rpx;
  font-weight: bold;
  color: #fff;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.nf-goods-bottom {
  margin-top: 12rpx;
}
.nf-price {
  font-size: 30rpx;
  font-weight: bold;
  color: #e50914;
}
.nf-history-footer {
  text-align: center;
  padding: 30rpx 0 60rpx;
}
.nf-history-footer-text {
  font-size: 24rpx;
  color: rgba(255,255,255,0.3);
}
.nf-history-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 200rpx;
}
.nf-history-empty-icon {
  margin-bottom: 20rpx;
}
.nf-history-empty-text {
  font-size: 28rpx;
  color: rgba(255,255,255,0.4);
}
</style>
