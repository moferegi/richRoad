<template>
  <view class="nf-goods-detail" :class="{ 'lock-scroll': isPageLocked }">
    <!-- 页面滚动锁定：弹窗打开时禁止body滚动 -->
    <page-meta :page-style="isPageLocked ? 'overflow: hidden;' : ''" />
    <!-- 自定义导航栏（悬浮在轮播上方） -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('goodsDetail') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 商品轮播图 -->
    <goods-swiper class="nf-top-swiper" :list="data.banner"></goods-swiper>

    <!-- 商品基本信息 -->
    <view class="nf-product-info">
      <view class="nf-price-row">
        <text class="nf-price-symbol">{{ cs }}</text>
        <text class="nf-price-num">{{ goodsDisplayPrice }}</text>
        <view class="nf-sale-tag" v-if="data.saleCount && data.saleCount > 0">
          <text>{{ $t('sold') }} {{ data.saleCount }}</text>
        </view>
        <!-- 预售标志 -->
        <view class="nf-presale-badge" v-if="data.isPresale">
          <text v-if="presaleCountdownType === 'ended'">{{ $t('presaleEnded') }}</text>
          <text v-else>{{ $t('presaleBadge') }}</text>
        </view>
      </view>
      <text class="nf-product-title">{{ $lt(data.title) }}</text>
      <text class="nf-product-desc">{{ $lt(data.description) }}</text>
      <view class="nf-meta-row">
        <text>{{ $t('stock') }}: {{ getTotalInventory(data.skus) }}</text>
        <text v-if="data.view_num" class="nf-meta-sep">{{ $t('viewCount') }}: {{ data.view_num }}</text>
      </view>
      <!-- 预售时间 -->
      <view class="nf-presale-time" v-if="data.isPresale && (data.presaleStart || data.presaleEnd)">
        <text class="nf-label-dot" style="background: #f59e0b;"></text>
        <text>{{ $t('presaleTime') }}：{{ formatDate(data.presaleStart) }} ~ {{ formatDate(data.presaleEnd) }}</text>
      </view>
      <!-- 预售倒计时 + 进度条 -->
      <view class="nf-presale-countdown-section" v-if="data.isPresale">
        <view class="nf-presale-countdown-row">
          <text class="nf-countdown-label" v-if="presaleCountdownType === 'start'">{{ $t('presaleStartsIn') }}</text>
          <text class="nf-countdown-label" v-else-if="presaleCountdownType === 'end'">{{ $t('presaleEndsIn') }}</text>
          <text class="nf-countdown-label nf-countdown-ended" v-else>{{ $t('presaleEnded') }}</text>
          <text class="nf-countdown-time" v-if="presaleCountdownType !== 'ended'">{{ presaleCountdownText }}</text>
        </view>
        <view class="nf-presale-progress" v-if="data.presaleQty > 0">
          <view class="nf-progress-bar">
            <view class="nf-progress-fill" :style="{ width: presaleProgress + '%' }"></view>
          </view>
          <text class="nf-progress-text">{{ $t('sold') }} {{ data.presaleSold || 0 }}/{{ data.presaleQty }}</text>
        </view>
      </view>
    </view>

    <!-- 商品属性 -->
    <view class="nf-section-card" v-if="parsedAttrs.length > 0">
      <view class="nf-section-label" style="margin-bottom: 16rpx;">
        <text class="nf-label-dot" style="background: #6366f1;"></text>
        <text>{{ $t('goodAttrs') }}</text>
      </view>
      <view class="nf-attrs-grid">
        <view class="nf-attr-item" v-for="(attr, idx) in parsedAttrs" :key="idx">
          <text class="nf-attr-label">{{ $lt(attr.labelI18n) || attr.label }}</text>
          <text class="nf-attr-value">{{ $lt(attr.valueI18n) || attr.value }}</text>
        </view>
      </view>
    </view>

    <!-- 优惠券 -->
    <view class="nf-section-card" @tap="opencoupon">
      <view class="nf-section-row">
        <view class="nf-section-label">
          <text class="nf-label-dot" style="background: #2563eb;"></text>
          <text>{{ $t('coupon') }}</text>
        </view>
        <view class="nf-section-value">
          <text class="nf-coupon-text" v-if="selectedCoupon.couponNum">
            {{ selectedCoupon.minSpend > 0
              ? $t('couponFull').replace('{min}', selectedCoupon.minSpend/100).replace('{off}', selectedCoupon.discount/100)
              : $t('couponNoMin').replace('{off}', selectedCoupon.discount/100) }}
          </text>
          <text class="nf-coupon-hint" v-else>{{ $t('claimCoupon') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.3)"></uni-icons>
        </view>
      </view>
    </view>

    <!-- 图文详情 -->
    <view class="nf-detail-divider">
      <view class="nf-divider-line"></view>
      <text>{{ $t('graphicDetail') }}</text>
      <view class="nf-divider-line"></view>
    </view>

    <view class="nf-detail-content">
      <goodsDetail :detail="localDetail"></goodsDetail>
    </view>

    <!-- SKU选择器 -->
    <goods-sku
      v-if="data.skus"
      ref="goodsSkuRef"
      :isCart="isCart"
      :good="data"
      :selectedCoupon="selectedCoupon"
      @skuVisibleChange="onSkuVisibleChange"
    ></goods-sku>

    <!-- 底部导航 -->
    <view class="nf-bottom-nav">
      <view class="nf-nav-icons">
        <view class="nf-nav-icon-item" @tap="goTo()">
          <LazyImage class="nf-nav-icon-img" src="./../../static/images/tabBar/home.png"></LazyImage>
          <text class="nf-nav-icon-text">{{ $t('home') }}</text>
        </view>
        <view class="nf-nav-icon-item" @tap="goToKefu">
          <LazyImage class="nf-nav-icon-img" src="./../../static/images/tabBar/service.svg"></LazyImage>
          <text class="nf-nav-icon-text">{{ $t('customerService') }}</text>
        </view>
        <view class="nf-nav-icon-item" @tap="goTo('cart')">
          <LazyImage class="nf-nav-icon-img" src="./../../static/images/tabBar/cart.png"></LazyImage>
          <text class="nf-nav-icon-text">{{ $t('cart') }}</text>
        </view>
        <view class="nf-nav-icon-item" @tap="addCollect">
          <LazyImage class="nf-nav-icon-img" :src="!collectionFlag ? './../../static/collection.png' : './../../static/collect.png'"></LazyImage>
          <text class="nf-nav-icon-text">{{ $t('collectText') }}</text>
        </view>
      </view>
      <view class="nf-buy-buttons">
        <view class="nf-add-cart-btn" @tap="addToCart()">{{ $t('addToCart') }}</view>
        <view class="nf-buy-now-btn" @tap="goodsTapPay('pay')">{{ $t('buyNow') }}</view>
      </view>
    </view>

    <!-- 优惠券弹出层 Netflix风格 -->
    <view class="nf-mask" v-if="couponshow" @tap="hidecoupon" @touchmove.stop.prevent></view>
    <view class="nf-coupon-popup" :class="{ show: couponshow }">
      <view class="nf-coupon-popup-header">
        <text class="nf-coupon-popup-title">{{ $t('claimCoupon') }}</text>
        <view class="nf-coupon-popup-close" @tap="hidecoupon">
          <uni-icons type="close" size="18" color="rgba(15,23,42,0.68)"></uni-icons>
        </view>
      </view>
      <scroll-view class="nf-coupon-scroll" scroll-y>
        <view v-if="couponList.length === 0" class="nf-coupon-empty">
          <text>{{ $t('noCoupons') }}</text>
        </view>
        <view v-for="(item, index) in couponList" :key="index"
          class="nf-coupon-card" :class="{ 'nf-coupon-disabled': item.canUse === 0 }">
          <view class="nf-coupon-left">
            <text class="nf-coupon-amount">{{ cs }}{{ item.discount / 100 }}</text>
            <text class="nf-coupon-condition">{{ item.minSpend > 0 ? $t('couponFull').replace('{min}', item.minSpend/100).replace('{off}', item.discount/100) : $t('couponNoLimit') }}</text>
          </view>
          <view class="nf-coupon-right">
            <text class="nf-coupon-name">{{ $lt(item.name) }}</text>
            <text class="nf-coupon-date">{{ $t('couponExpiry') }}: {{ item.startTime }} ~ {{ item.endTime }}</text>
            <text class="nf-coupon-unavail" v-if="item.canUse === 0">{{ $t('couponUnavailable') }}</text>
          </view>
          <view class="nf-coupon-action">
            <view v-if="item.status === 1" class="nf-coupon-btn nf-coupon-btn-used">{{ $t('couponInUse') }}</view>
            <view v-else-if="item.canUse === 0" class="nf-coupon-btn nf-coupon-btn-disabled"></view>
            <view v-else class="nf-coupon-btn" :class="item.couponNum ? 'nf-coupon-btn-use' : 'nf-coupon-btn-claim'"
              @tap="onReceive(item, index)">
              {{ item.couponNum ? $t('couponUse') : $t('couponClaim') }}
            </view>
          </view>
        </view>
      </scroll-view>
    </view>

  </view>
</template>

<script setup>
import goodsSwiper from './components/goods-swiper.vue'
import goodsSku from './components/goods-sku.vue'
import goodsDetail from './components/goods-detail.vue'
import { ref, computed, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { findGood } from '@/api/product.js'
import { SelfOrderList } from '@/api/order.js'
import { myRouter } from '@/utils/permission'
import { findCollect, createCollect } from '@/api/collect.js'
import { claimCouponByUser, getAllClaimCoupon } from '@/api/coupon.js'
import { useUserStore } from '@/pinia/modules/user'
import { usePlayHistoryStore } from '@/pinia/modules/playHistory.js'
import { getUrl } from '@/utils/url.js'
import { localText } from '@/utils/i18n.js'
import { formatLocalizedPrice } from '@/utils/price-i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'
const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const playHistoryStore = usePlayHistoryStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $lt = computed(() => langStore.$lt)
const $t = computed(() => langStore.$t)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const goBack = () => { uni.navigateBack() }

const data = ref({})
const goodsDisplayPrice = computed(() => {
  return formatLocalizedPrice(data.value?.price, data.value?.priceI18n, locale.value)
})
const collectionFlag = ref('')
const goodID = ref(0)
const userStore = useUserStore()
const token = userStore.token || ''

onLoad((options) => {
  if (options.id) {
    goodID.value = options.id
    init()
  }
})

const init = async () => {
  const res = await findGood(goodID.value)
  if (res.code === 0) {
    data.value = res.data.regood
    // 保存本地浏览记录（无论是否登录）
    playHistoryStore.saveBrowse(goodID.value, {
      imageUrl: data.value.imageUrl || (data.value.banner && data.value.banner[0]) || '',
      title: data.value.title || ''
    })
    if (data.value.isPresale) startCountdown()
    // 预售弹窗开关开启时自动弹出
    if (data.value.isPresale && data.value.presalePopupEnabled) {
      const title = localText(data.value.presalePopupTitle) || $t.value('presale')
      const content = localText(data.value.presalePopupContent) || ''
      uni.showModal({ title, content, showCancel: false, confirmText: $t.value('confirm') })
    }
  }
  if (token) {
    const status = await findCollect({ goodID: goodID.value })
    if (status.code === 0) collectionFlag.value = status.data
  }
}

// 解析商品属性（JSON数组）
const parsedAttrs = computed(() => {
  if (!data.value.attrs) return []
  try {
    const arr = typeof data.value.attrs === 'string' ? JSON.parse(data.value.attrs) : data.value.attrs
    if (!Array.isArray(arr)) return []
    const tryParseJson = (s) => {
      if (!s || typeof s !== 'string' || s.charAt(0) !== '{') return null
      try { return JSON.parse(s) } catch { return null }
    }
    return arr.map(item => {
      // name 可能是序列化的i18n JSON字符串，nameI18n 可能被删除
      let labelI18n = item.nameI18n || item.labelI18n
      let label = item.name || item.label || ''
      if (!labelI18n && typeof label === 'string') {
        const parsed = tryParseJson(label)
        if (parsed) { labelI18n = parsed; label = parsed['zh'] || Object.values(parsed)[0] || label }
      }
      let valueI18n = item.valueI18n
      let value = item.value || ''
      if (!valueI18n && typeof value === 'string') {
        const parsed = tryParseJson(value)
        if (parsed) { valueI18n = parsed; value = parsed['zh'] || Object.values(parsed)[0] || value }
      }
      return { label, value, labelI18n, valueI18n }
    })
  } catch { return [] }
})

// 解析商品详情（多语言JSON字段 → HTML字符串）
const localDetail = computed(() => {
  return localText(data.value.detail, locale.value) || ''
})

const formatDate = (d) => {
  if (!d) return ''
  const dt = new Date(d)
  return `${dt.getFullYear()}-${String(dt.getMonth()+1).padStart(2,'0')}-${String(dt.getDate()).padStart(2,'0')}`
}

// ========== 预售倒计时 ==========
const countdownTick = ref(0)
let countdownTimer = null

const presaleCountdownType = computed(() => {
  void countdownTick.value
  if (!data.value.isPresale) return 'ended'
  const now = Date.now()
  const start = new Date(data.value.presaleStart).getTime()
  const end = new Date(data.value.presaleEnd).getTime()
  if (now < start) return 'start'
  if (now < end) return 'end'
  return 'ended'
})

const presaleCountdownText = computed(() => {
  void countdownTick.value
  if (!data.value.isPresale) return ''
  const now = Date.now()
  const type = presaleCountdownType.value
  let target
  if (type === 'start') target = new Date(data.value.presaleStart).getTime()
  else if (type === 'end') target = new Date(data.value.presaleEnd).getTime()
  else return ''
  const diff = Math.max(0, target - now)
  const dd = Math.floor(diff / 86400000)
  const hh = Math.floor((diff % 86400000) / 3600000)
  const mm = Math.floor((diff % 3600000) / 60000)
  const ss = Math.floor((diff % 60000) / 1000)
  const dayUnit = $t.value('countdownDay')
  if (dd > 0) return `${dd}${dayUnit} ${String(hh).padStart(2,'0')}:${String(mm).padStart(2,'0')}:${String(ss).padStart(2,'0')}`
  return `${String(hh).padStart(2,'0')}:${String(mm).padStart(2,'0')}:${String(ss).padStart(2,'0')}`
})

const presaleProgress = computed(() => {
  if (!data.value.presaleQty || data.value.presaleQty === 0) return 0
  return Math.min(100, Math.round((data.value.presaleSold || 0) / data.value.presaleQty * 100))
})

const startCountdown = () => {
  if (countdownTimer) clearInterval(countdownTimer)
  countdownTimer = setInterval(() => { countdownTick.value++ }, 1000)
}

onUnmounted(() => {
  if (countdownTimer) clearInterval(countdownTimer)
})

const goodsSkuRef = ref()
const isCart = ref(false)
const skuVisible = ref(false)
const onSkuVisibleChange = (visible) => { skuVisible.value = visible }

const addToCart = () => {
  if (!token) {
    uni.showToast({ title: $t.value('pleaseLogin'), mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  // 预售已结束，不能加购
  if (data.value.isPresale && presaleCountdownType.value === 'ended') {
    uni.showToast({ title: $t.value('presaleEndedToast'), icon: 'none' })
    return
  }
  isCart.value = true
  goodsSkuRef.value.showSku()
}

const goTo = (path) => {
  if (path === 'cart') uni.navigateTo({ url: '/pages/cart/index' })
  else uni.switchTab({ url: '/pages/tabBar/index' })
}

const goToKefu = () => { uni.navigateTo({ url: '/pages/kefu/index' }) }

const toOrder = () => {
  myRouter(`/pages/orderInfo/orderInfo?skuID=${data.value.skus[0].ID}&goodID=${data.value.skus[0].goodID}`)
}

const goodsTapPay = async () => {
  // 预售已结束，不能购买
  if (data.value.isPresale && presaleCountdownType.value === 'ended') {
    uni.showToast({ title: $t.value('presaleEndedToast'), icon: 'none' })
    return
  }
  // 检查是否有该商品的待付款订单
  try {
    const token = uni.getStorageSync('x-token')
    if (token) {
      const res = await SelfOrderList({ status: '0', page: 1, pageSize: 50 })
      if (res.code === 0 && res.data?.list) {
        const pending = res.data.list.find(order =>
          order.detail && order.detail.some(d => String(d.goodID) === String(goodID.value))
        )
        if (pending) {
          uni.showModal({
            title: $t.value('pendingOrderTitle'),
            content: $t.value('pendingOrderExist'),
            confirmText: $t.value('goToPay'),
            cancelText: $t.value('continueBuy'),
            success: (modalRes) => {
              if (modalRes.confirm) {
                uni.navigateTo({ url: `/pages/orderDetail/orderDetail?orderID=${pending.ID}` })
                return
              }
              isCart.value = false
              goodsSkuRef.value.showSku()
            }
          })
          return
        }
      }
    }
  } catch (e) { /* ignore */ }
  isCart.value = false
  goodsSkuRef.value.showSku()
}

const addCollect = async () => {
  if (token) {
    const res = await createCollect({ goodID: Number(goodID.value) })
    if (res.code === 0) {
      collectionFlag.value = !collectionFlag.value
      uni.showToast({ title: collectionFlag.value ? $t.value('collectSuccess') : $t.value('uncollectSuccess'), mask: true, icon: 'none' })
    }
  } else {
    uni.showToast({ title: $t.value('pleaseLogin'), mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
  }
}

const selectedCoupon = ref({})
const couponshow = ref(false)
const couponList = ref([])
const isPageLocked = computed(() => skuVisible.value || couponshow.value)

const getTotalInventory = (skus) => {
  if (!skus || skus.length === 0) return 0
  return skus.reduce((total, sku) => total + sku.inventory, 0)
}

const opencoupon = async () => {
  couponshow.value = true
  const res = await getAllClaimCoupon({ goodIds: [data.value.ID] })
  if (res.code === 0) {
    // 只显示可用的优惠券：未使用、未过期、商品可用
    couponList.value = (res.data || []).filter(c => {
      if (c.used || c.expired || c.canUse === 0) return false
      return true
    })
  }
}
const hidecoupon = () => {
  couponshow.value = false
}

const onReceive = async (item) => {
  uni.showLoading({ title: item.couponNum == 0 ? $t.value('claiming') : $t.value('selecting'), mask: true })
  try {
    if (!item.couponNum) {
      const res = await claimCouponByUser({ couponID: item.couponID })
      item.couponNum = res.data
      uni.showToast({ title: $t.value('claimSuccess'), icon: 'success', duration: 1500 })
    }
    selectedCoupon.value = item
    setTimeout(() => { hidecoupon() }, 500)
  } catch (error) {
    uni.showToast({ title: $t.value('operationFailed'), icon: 'none' })
  } finally { uni.hideLoading() }
}
</script>

<style lang="scss">
page { background-color: #f4f7fb; }

.nf-goods-detail {
  --nf-text-secondary: rgba(15, 23, 42, 0.62);
  --nf-text-tertiary: rgba(15, 23, 42, 0.52);
  --nf-text-weak: rgba(15, 23, 42, 0.45);
  min-height: 100vh;
  background: #f4f7fb;
  padding-bottom: calc(220rpx + constant(safe-area-inset-bottom));
  padding-bottom: calc(220rpx + env(safe-area-inset-bottom));
}
.nf-goods-detail.lock-scroll { height: 100vh; overflow: hidden; }

.nf-top-swiper {
  display: block;
  margin-top: calc(var(--status-bar-height, 0px) + 88rpx);
}

/* 导航栏 */
.nf-navbar {
  position: fixed; top: 0; left: 0; right: 0; z-index: 100;
  background: rgba(244, 247, 251, 0.9); backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.24);
  padding: 0 28rpx 16rpx;
}
.nf-navbar-status { height: var(--status-bar-height, 0px); }
.nf-navbar-content { display: flex; align-items: center; justify-content: space-between; height: 88rpx; }
.nf-navbar-back {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: #ffffff; border: 1rpx solid rgba(148, 163, 184, 0.3);
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.08);
  display: flex; align-items: center; justify-content: center;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  &:active {
    transform: scale(0.96);
    box-shadow: 0 6rpx 16rpx rgba(15, 23, 42, 0.1);
  }
}
.nf-navbar-title { font-size: 34rpx; font-weight: 700; color: #0f172a; letter-spacing: 1rpx; }

/* 商品信息 */
.nf-product-info {
  margin: 16rpx 24rpx 0;
  padding: 28rpx 24rpx;
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  border-radius: 20rpx;
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
}
.nf-price-row {
  display: flex; align-items: baseline; gap: 4rpx; margin-bottom: 16rpx; flex-wrap: wrap;
}
.nf-price-symbol { font-size: 30rpx; font-weight: 700; color: #2563eb; }
.nf-price-num { font-size: 52rpx; font-weight: 800; color: #2563eb; line-height: 1; }
.nf-sale-tag {
  margin-left: 16rpx; font-size: 20rpx; color: var(--nf-text-tertiary);
  padding: 4rpx 14rpx; background: rgba(226, 232, 240, 0.7); border-radius: 8rpx;
}
.nf-presale-badge {
  margin-left: 12rpx; font-size: 20rpx; color: #fff; font-weight: 600;
  padding: 4rpx 16rpx; background: linear-gradient(135deg, #f59e0b, #ef4444); border-radius: 8rpx;
}
.nf-product-title { display: block; font-size: 32rpx; font-weight: 700; color: #0f172a; line-height: 1.4; margin-bottom: 8rpx; }
.nf-product-desc { display: block; font-size: 26rpx; color: var(--nf-text-secondary); line-height: 1.45; margin-bottom: 12rpx; }
.nf-meta-row {
  font-size: 24rpx; color: var(--nf-text-tertiary); display: flex; gap: 24rpx;
}
.nf-meta-sep { }
.nf-presale-time {
  display: flex; align-items: center; gap: 10rpx; margin-top: 16rpx; padding: 16rpx 20rpx;
  background: rgba(245, 158, 11, 0.12); border: 1rpx solid rgba(245, 158, 11, 0.2);
  border-radius: 12rpx; font-size: 24rpx; color: #f59e0b;
}
/* 预售倒计时+进度条 */
.nf-presale-countdown-section {
  margin-top: 16rpx; padding: 20rpx 24rpx;
  background: rgba(245, 158, 11, 0.1); border: 1rpx solid rgba(245, 158, 11, 0.18);
  border-radius: 12rpx;
}
.nf-presale-countdown-row {
  display: flex; align-items: center; gap: 12rpx; margin-bottom: 16rpx;
}
.nf-countdown-label { font-size: 24rpx; color: #f59e0b; font-weight: 600; }
.nf-countdown-ended { color: rgba(15, 23, 42, 0.45); }
.nf-countdown-time {
  font-size: 28rpx; color: #0f172a; font-weight: 700;
  padding: 4rpx 16rpx; background: rgba(37, 99, 235, 0.12); border-radius: 8rpx;
  font-variant-numeric: tabular-nums;
}
.nf-presale-progress { margin-top: 4rpx; }
.nf-progress-bar {
  height: 12rpx; background: rgba(148, 163, 184, 0.28); border-radius: 6rpx; overflow: hidden;
}
.nf-progress-fill {
  height: 100%; background: linear-gradient(90deg, #f59e0b, #ef4444); border-radius: 6rpx;
  transition: width 0.3s;
}
.nf-progress-text {
  display: block; margin-top: 8rpx; font-size: 22rpx; color: var(--nf-text-tertiary); text-align: right;
}

/* 属性网格 */
.nf-attrs-grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.nf-attr-item {
  display: flex; flex-direction: column; gap: 4rpx;
  min-width: 200rpx; padding: 12rpx 16rpx;
  background: #f8fafc; border-radius: 10rpx;
  border: 1rpx solid rgba(148, 163, 184, 0.24);
}
.nf-attr-label { font-size: 22rpx; color: var(--nf-text-tertiary); }
.nf-attr-value { font-size: 26rpx; color: #0f172a; }

/* 通用区块卡片 */
.nf-section-card {
  margin: 20rpx 24rpx; padding: 24rpx 28rpx;
  background: #ffffff; border: 1rpx solid rgba(148, 163, 184, 0.2);
  border-radius: 20rpx; box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  &:active {
    transform: translateY(1rpx);
    box-shadow: 0 8rpx 20rpx rgba(15, 23, 42, 0.1);
  }
}
.nf-section-row { display: flex; justify-content: space-between; align-items: center; }
.nf-section-label { display: flex; align-items: center; gap: 12rpx; font-size: 28rpx; color: #0f172a; }
.nf-label-dot { width: 12rpx; height: 12rpx; border-radius: 4rpx; }
.nf-section-value { display: flex; align-items: center; gap: 8rpx; }
.nf-coupon-text { font-size: 26rpx; color: #2563eb; font-weight: 600; }
.nf-coupon-hint { font-size: 26rpx; color: var(--nf-text-weak); }

/* 图文详情分割 */
.nf-detail-divider {
  display: flex; align-items: center; justify-content: center;
  gap: 20rpx; padding: 30rpx 0;
  text { font-size: 26rpx; color: var(--nf-text-weak); }
}
.nf-divider-line {
  height: 1rpx; width: 100rpx;
  background: linear-gradient(90deg, transparent, rgba(15, 23, 42, 0.16), transparent);
}
.nf-detail-content { padding: 0 24rpx; }

/* 底部导航 */
.nf-bottom-nav {
  position: fixed; left: 0; right: 0; bottom: 0; z-index: 98;
  display: flex; align-items: center; height: 110rpx;
  background: rgba(255, 255, 255, 0.95); backdrop-filter: blur(24px);
  border-top: 1rpx solid rgba(148, 163, 184, 0.24);
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}
.nf-nav-icons { display: flex; flex: 1; }
.nf-nav-icon-item {
  display: flex; flex-direction: column; align-items: center; justify-content: center; width: 96rpx;
  transition: transform 0.2s ease;
  &:active { transform: scale(0.95); }
}
.nf-nav-icon-img { width: 40rpx; height: 40rpx; margin-bottom: 4rpx; opacity: 0.72; }
.nf-nav-icon-text { font-size: 20rpx; color: var(--nf-text-tertiary); }
.nf-buy-buttons { display: flex; height: 100%; }
.nf-add-cart-btn, .nf-buy-now-btn {
  padding: 0 36rpx; height: 100%; display: flex; align-items: center; justify-content: center;
  font-size: 28rpx; font-weight: 600; color: #fff;
  transition: filter 0.2s ease, transform 0.2s ease;
  &:active {
    filter: brightness(0.95);
    transform: translateY(1rpx);
  }
}
.nf-add-cart-btn { background: linear-gradient(135deg, #f59e0b, #fb923c); }
.nf-buy-now-btn { background: linear-gradient(135deg, #2563eb, #0ea5e9); }

/* 优惠券弹出层 Netflix风格 */
.nf-mask {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(2, 6, 23, 0.52); z-index: 900;
}
.nf-coupon-popup {
  position: fixed; left: 0; right: 0; bottom: -100vh; z-index: 999;
  background: #f8fafc; border-radius: 28rpx 28rpx 0 0;
  transition: all 0.3s ease;
  &.show { bottom: 0; }
}
.nf-coupon-popup-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 28rpx 32rpx 16rpx; border-bottom: 1rpx solid rgba(148, 163, 184, 0.18);
}
.nf-coupon-popup-title { font-size: 32rpx; font-weight: 700; color: #0f172a; }
.nf-coupon-popup-close {
  width: 56rpx; height: 56rpx; display: flex; align-items: center; justify-content: center;
  background: rgba(148, 163, 184, 0.16); border-radius: 50%;
}
.nf-coupon-scroll { width: 100vw; height: 55vh; padding: 16rpx 0; }
.nf-coupon-empty {
  display: flex; align-items: center; justify-content: center; height: 200rpx;
  color: rgba(71, 85, 105, 0.75); font-size: 28rpx;
}

/* 优惠券卡片 */
.nf-coupon-card {
  display: flex; align-items: center; margin: 16rpx 24rpx; padding: 24rpx;
  background: #ffffff; border: 1rpx solid rgba(148, 163, 184, 0.2);
  border-radius: 16rpx; position: relative; overflow: hidden;
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  &:active {
    transform: scale(0.996);
    box-shadow: 0 8rpx 20rpx rgba(15, 23, 42, 0.1);
  }
}
.nf-coupon-card::before {
  content: ''; position: absolute; left: 0; top: 0; bottom: 0; width: 8rpx;
  background: linear-gradient(180deg, #2563eb, #0ea5e9);
}
.nf-coupon-disabled { opacity: 0.45; }
.nf-coupon-left {
  min-width: 160rpx; text-align: center; padding-right: 20rpx;
  border-right: 1rpx dashed rgba(148, 163, 184, 0.35);
}
.nf-coupon-amount { display: block; font-size: 44rpx; font-weight: 800; color: #2563eb; line-height: 1.2; }
.nf-coupon-condition { display: block; font-size: 20rpx; color: var(--nf-text-tertiary); margin-top: 4rpx; }
.nf-coupon-right { flex: 1; padding: 0 20rpx; }
.nf-coupon-name { display: block; font-size: 26rpx; color: #0f172a; font-weight: 600; margin-bottom: 6rpx; }
.nf-coupon-date { display: block; font-size: 20rpx; color: var(--nf-text-tertiary); }
.nf-coupon-unavail { display: block; font-size: 20rpx; color: #2563eb; margin-top: 4rpx; }
.nf-coupon-action { min-width: 120rpx; display: flex; align-items: center; justify-content: center; }
.nf-coupon-btn {
  padding: 10rpx 24rpx; border-radius: 8rpx; font-size: 24rpx; font-weight: 600; text-align: center;
  transition: filter 0.2s ease, transform 0.2s ease;
  &:active {
    filter: brightness(0.95);
    transform: translateY(1rpx);
  }
}
.nf-coupon-btn-claim { background: linear-gradient(135deg, #2563eb, #0ea5e9); color: #fff; }
.nf-coupon-btn-use { background: linear-gradient(135deg, #f59e0b, #fb923c); color: #fff; }
.nf-coupon-btn-used { background: rgba(148, 163, 184, 0.2); color: rgba(15, 23, 42, 0.35); }
.nf-coupon-btn-disabled { display: none; }
</style>


