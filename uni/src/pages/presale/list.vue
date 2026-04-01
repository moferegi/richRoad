<template>
  <view class="nf-presale">
    <view class="nf-presale-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('presaleZone') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 商品列表 -->
    <scroll-view
      scroll-y
      :show-scrollbar="false"
      class="nf-presale-scroll"
      @scrolltolower="debouncedLower"
    >
      <view v-if="list.length" class="nf-goods-list">
        <view
          class="nf-goods-card"
          v-for="(item, index) in list"
          :key="index"
          @tap="goDetail(item)"
        >
          <view class="nf-goods-img-wrap">
            <image class="nf-goods-img" :src="getUrl(item.good && item.good.imageUrl)" mode="aspectFill" />
            <view class="nf-goods-img-overlay"></view>
            <!-- 预售标签 -->
            <view class="nf-presale-badge">
              <text class="nf-presale-badge-text">{{ $t('presale') }}</text>
            </view>
          </view>
          <view class="nf-goods-info">
            <text class="nf-goods-title">{{ $lt(item.good && item.good.title) }}</text>

            <!-- 倒计时 -->
            <view class="nf-presale-countdown">
              <text class="nf-countdown-label" v-if="getCountdownType(item) === 'start'">
                {{ $t('presaleStartsIn') }}
              </text>
              <text class="nf-countdown-label" v-else-if="getCountdownType(item) === 'end'">
                {{ $t('presaleEndsIn') }}
              </text>
              <text class="nf-countdown-label nf-countdown-ended" v-else>
                {{ $t('presaleEnded') }}
              </text>
              <text class="nf-countdown-time" v-if="getCountdownType(item) !== 'ended'">
                {{ formatCountdown(item) }}
              </text>
            </view>

            <view class="nf-goods-bottom">
              <view class="nf-price-row">
                <text class="nf-price-label">{{ $t('presalePrice') }}</text>
                <text class="nf-price">¥{{ formatPrice(item.presalePrice || (item.good && item.good.price)) }}</text>
              </view>
              <view class="nf-presale-progress">
                <view class="nf-progress-bar">
                  <view class="nf-progress-fill" :style="{ width: getProgress(item) + '%' }"></view>
                </view>
                <text class="nf-progress-text">{{ $t('sold') }} {{ item.presaleSold || 0 }}/{{ item.presaleTotal || 0 }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <!-- 底部加载状态 -->
      <view class="nf-presale-footer" v-if="list.length > 0">
        <text class="nf-presale-footer-text">{{ isBottom ? $t('reachedBottom') : $t('loading') }}</text>
      </view>

      <!-- 空状态 -->
      <view class="nf-presale-empty" v-if="!list.length && !loading">
        <view class="nf-presale-empty-icon">
          <uni-icons type="shop" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-presale-empty-text">{{ $t('noPresaleGoods') }}</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getPresaleGoodList } from '@/api/presale'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const list = ref([])
const loading = ref(false)
const isBottom = ref(false)
let params = { page: 1, pageSize: 10 }

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  const cents = parseInt(priceInCents)
  if (isNaN(cents)) return '0.00'
  return (cents / 100).toFixed(2)
}

const getCountdownType = (item) => {
  const now = Date.now()
  const start = new Date(item.presaleStartTime).getTime()
  const end = new Date(item.presaleEndTime).getTime()
  if (now < start) return 'start'
  if (now < end) return 'end'
  return 'ended'
}

const formatCountdown = (item) => {
  const now = Date.now()
  const type = getCountdownType(item)
  let target
  if (type === 'start') target = new Date(item.presaleStartTime).getTime()
  else target = new Date(item.presaleEndTime).getTime()

  const diff = Math.max(0, target - now)
  const d = Math.floor(diff / 86400000)
  const h = Math.floor((diff % 86400000) / 3600000)
  const m = Math.floor((diff % 3600000) / 60000)
  const s = Math.floor((diff % 60000) / 1000)

  if (d > 0) return `${d}d ${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
}

const getProgress = (item) => {
  if (!item.presaleTotal || item.presaleTotal === 0) return 0
  return Math.min(100, Math.round((item.presaleSold || 0) / item.presaleTotal * 100))
}

const init = async () => {
  loading.value = true
  params = { page: 1, pageSize: 10 }
  isBottom.value = false
  try {
    const res = await getPresaleGoodList(params)
    if (res.code === 0) {
      list.value = res.data.list || []
      isBottom.value = list.value.length < params.pageSize
    }
  } finally {
    loading.value = false
  }
}

onShow(() => {
  init()
})

const debounce = (func, delay) => {
  let timer
  return function (...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => func.apply(this, args), delay)
  }
}

const lower = async () => {
  if (isBottom.value) return
  params.page += 1
  const res = await getPresaleGoodList(params)
  if (res.code === 0 && res.data.list && res.data.list.length) {
    list.value.push(...res.data.list)
  } else {
    isBottom.value = true
  }
}
const debouncedLower = debounce(lower, 300)

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/index' }) })
}

const goDetail = (item) => {
  uni.navigateTo({ url: `/pages/player/index?id=${item.goodId || item.good_id}` })
}
</script>

<style lang="scss" scoped>
.nf-presale {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-presale-bg {
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
.nf-presale-scroll {
  position: relative;
  z-index: 1;
  height: 100vh;
  padding-top: calc(var(--status-bar-height, 44rpx) + 88rpx);
}
.nf-goods-list {
  padding: 20rpx 24rpx;
}
.nf-goods-card {
  display: flex;
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 20rpx;
}
.nf-goods-img-wrap {
  position: relative;
  width: 240rpx;
  min-height: 240rpx;
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
.nf-presale-badge {
  position: absolute;
  top: 12rpx;
  left: 12rpx;
  background: #e50914;
  border-radius: 8rpx;
  padding: 4rpx 12rpx;
}
.nf-presale-badge-text {
  font-size: 22rpx;
  color: #fff;
  font-weight: bold;
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
.nf-presale-countdown {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin: 8rpx 0;
}
.nf-countdown-label {
  font-size: 22rpx;
  color: #e50914;
}
.nf-countdown-ended {
  color: rgba(255,255,255,0.4);
}
.nf-countdown-time {
  font-size: 24rpx;
  color: #e50914;
  font-weight: bold;
  font-family: monospace;
}
.nf-goods-bottom {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}
.nf-price-row {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}
.nf-price-label {
  font-size: 22rpx;
  color: rgba(255,255,255,0.5);
}
.nf-price {
  font-size: 32rpx;
  font-weight: bold;
  color: #e50914;
}
.nf-presale-progress {
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}
.nf-progress-bar {
  height: 8rpx;
  background: rgba(255,255,255,0.1);
  border-radius: 4rpx;
  overflow: hidden;
}
.nf-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #e50914, #ff6b35);
  border-radius: 4rpx;
}
.nf-progress-text {
  font-size: 20rpx;
  color: rgba(255,255,255,0.5);
}
.nf-presale-footer {
  text-align: center;
  padding: 30rpx 0 60rpx;
}
.nf-presale-footer-text {
  font-size: 24rpx;
  color: rgba(255,255,255,0.3);
}
.nf-presale-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 200rpx;
}
.nf-presale-empty-icon {
  margin-bottom: 20rpx;
}
.nf-presale-empty-text {
  font-size: 28rpx;
  color: rgba(255,255,255,0.4);
}
</style>
