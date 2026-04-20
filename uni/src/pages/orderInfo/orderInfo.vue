<template>
  <view class="nf-orderinfo">
    <view class="nf-orderinfo-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('orderConfirm') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 提交中loading -->
    <view v-if="submitting" class="nf-creating">
      <view class="nf-creating-spinner"></view>
      <text class="nf-creating-text">{{ $t('orderCreating') }}</text>
    </view>

    <view class="nf-body" v-else>
      <!-- 收货地址 -->
      <view class="nf-card nf-address-card" @tap="toAddress">
        <view class="nf-address-icon">📍</view>
        <view class="nf-address-info" v-if="hasAddress">
          <text class="nf-address-name">{{ address.name }} {{ address.phone }}</text>
          <text class="nf-address-detail">{{ address.provinceStr }}{{ address.cityStr }}{{ address.areaStr }} {{ address.street }}</text>
        </view>
        <view class="nf-address-info" v-else>
          <text class="nf-address-name">{{ $t('selectAddress') }}</text>
          <text class="nf-address-detail">{{ $t('addAddressHint') }}</text>
        </view>
        <uni-icons type="right" size="16" color="rgba(255,255,255,0.3)"></uni-icons>
      </view>

      <!-- 商品列表 -->
      <view class="nf-card nf-goods-card">
        <view class="nf-goods-item" v-for="(d, i) in goodsList" :key="i">
          <image class="nf-goods-img" :src="d.sku?.externalPicturePath ? getExternalUrl(d.sku.externalPicturePath) : getUrl(d.sku?.picture)" mode="aspectFill"></image>
          <view class="nf-goods-info">
            <text class="nf-goods-name">{{ $lt(d?.sku?.name) || d?.sku?.name }}</text>
            <text class="nf-goods-desc">{{ $lt(d?.good?.description) || d?.good?.description }}</text>
            <text class="nf-goods-specs">{{ formatSpecs(d?.sku?.specs, d?.sku?.attrs) }}</text>
            <view class="nf-goods-bottom">
              <text class="nf-goods-price">{{ cs }}{{ (d.sku?.price || 0) / 100 }}</text>
              <text class="nf-goods-qty">×{{ d.quantity }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 优惠券选择 -->
      <view class="nf-card nf-option-card" @tap="openCouponPopup">
        <view class="nf-option-row">
          <view class="nf-option-label">
            <text class="nf-option-dot" style="background: #e50914;"></text>
            <text>{{ $t('selectCoupon') }}</text>
          </view>
          <view class="nf-option-value">
            <text class="nf-discount-text" v-if="selectedCouponDiscount > 0">-{{ cs }}{{ selectedCouponDiscount / 100 }}</text>
            <text class="nf-discount-hint" v-else>{{ $t('selectCoupon') }}</text>
            <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)"></uni-icons>
          </view>
        </view>
      </view>

      <!-- 价格汇总 -->
      <view class="nf-card nf-price-card">
        <view class="nf-price-row">
          <text class="nf-price-label">{{ $t('productAmount') }}</text>
          <text class="nf-price-val">{{ cs }}{{ (originPrice / 100).toFixed(2) }}</text>
        </view>
        <view class="nf-price-row nf-price-discount" v-if="selectedCouponDiscount > 0">
          <text class="nf-price-label">{{ $t('discountAmount') }}</text>
          <text class="nf-price-val">-{{ cs }}{{ (selectedCouponDiscount / 100).toFixed(2) }}</text>
        </view>
        <view class="nf-price-row nf-price-points">
          <view class="nf-points-left" @tap="togglePoints">
            <view class="nf-points-check">
              <checkbox value="points" :checked="usePoints" color="#e50914" style="transform: scale(0.7); margin-right: 6rpx; pointer-events: none;" />
              <text :style="(!pointsAllowed || userPoints <= 0) ? 'opacity:0.4' : ''">{{ $t('pointsDeduction') }}</text>
            </view>
          </view>
          <text class="nf-points-amount">{{ $t('maxPointsDeduct') }} {{ cs }}{{ maxDeductDisplay }}</text>
        </view>
        <view class="nf-points-detail" v-if="userPoints > 0">
          <text>{{ $t('pointsBalance') }}: {{ userPoints }}</text>
          <text v-if="usePoints && previewPointsUsed > 0" class="nf-points-used">  -{{ previewPointsUsed }} = {{ userPoints - previewPointsUsed }}</text>
        </view>
      </view>

      <view style="height: 140rpx;"></view>
    </view>

    <!-- 底部支付栏 -->
    <view class="nf-footer" v-if="!submitting">
      <view class="nf-footer-info">
        <text class="nf-footer-label">{{ $t('actualPayment') }}</text>
        <text class="nf-footer-price">{{ cs }}{{ (totalPrice / 100).toFixed(2) }}</text>
      </view>
      <view class="nf-footer-btn" @tap="submitOrder">
        <text>{{ $t('submitOrder') }}</text>
      </view>
    </view>

    <!-- Netflix 优惠券弹出层 -->
    <view class="nf-mask" v-if="couponShow" @tap="closeCouponPopup"></view>
    <view class="nf-coupon-popup" :class="{ show: couponShow }">
      <view class="nf-coupon-header">
        <text class="nf-coupon-title">{{ $t('selectCoupon') }}</text>
        <view class="nf-coupon-close" @tap="closeCouponPopup">✕</view>
      </view>
      <scroll-view class="nf-coupon-scroll" scroll-y>
        <view v-if="couponList.length === 0" class="nf-coupon-empty">
          <text>{{ $t('noCoupons') }}</text>
        </view>
        <view v-for="(c, idx) in couponList" :key="idx"
              class="nf-coupon-item" :class="{ 'nf-coupon-selected': selectedCouponNum === c.couponNum }"
              @tap="onCouponTap(c)">
          <view class="nf-coupon-left">
            <text class="nf-coupon-amount">{{ cs }}{{ (c.discount || 0) / 100 }}</text>
            <text class="nf-coupon-cond" v-if="c.minSpend > 0">{{ $t('couponFull').replace('{min}', c.minSpend / 100).replace('{off}', (c.discount || 0) / 100) }}</text>
            <text class="nf-coupon-cond" v-else>{{ $t('couponNoLimit') }}</text>
          </view>
          <view class="nf-coupon-right">
            <text class="nf-coupon-name">{{ $lt(c.name) || c.name }}</text>
            <text class="nf-coupon-exp">{{ $t('couponExpiry') || '' }} {{ formatCouponDate(c.endTime) }}</text>
            <view v-if="!c.couponNum" class="nf-coupon-claim">
              <text>{{ $t('claimCoupon') }}</text>
            </view>
            <view v-else-if="selectedCouponNum === c.couponNum" class="nf-coupon-deselect">
              <text>{{ $t('cancelCoupon') }}</text>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { placeOrder } from '@/api/order.js'
import { changeOrderPoints } from '@/api/order.js'
import { claimCouponByUser, getAllClaimCoupon } from '@/api/coupon.js'
import { checkNeedPay } from '@/api/base.js'
import { getUserInfo } from '@/api/base.js'
import { findGood } from '@/api/product.js'
import { getDefaultAddress } from '@/api/address.js'
import { getUrl, getExternalUrl } from "@/utils/url.js"
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const goBack = () => { uni.navigateBack() }

/* =================== 数据 =================== */
const submitting = ref(false)
const usePoints = ref(false)
const userPoints = ref(0)
const couponShow = ref(false)
const couponList = ref([])
const selectedCouponNum = ref('')
const selectedCouponDiscount = ref(0)

// 地址
const address = ref({})
const hasAddress = ref(false)

// 商品信息（本地预览）
const goodsList = ref([])         // [{ good, sku, quantity }]
const originPrice = ref(0)        // 原价合计（分）

// 订单参数（从路由传入）
const paramGoodID = ref(0)
const paramSkuID = ref(0)
const paramQuantity = ref(1)
const paramCouponNum = ref('')

/* =================== 价格计算（本地预览） =================== */
// 检查所有商品是否允许积分抵扣
const pointsAllowed = computed(() => {
  if (goodsList.value.length === 0) return false
  return goodsList.value.every(d => d.good?.pointsEnabled === true)
})

// 商品层面允许的最大积分总和
const goodMaxPoints = computed(() => {
  let total = 0
  for (const d of goodsList.value) {
    const maxUse = d.good?.pointsMaxUse || 0
    if (maxUse > 0) {
      total += maxUse * (d.quantity || 1)
    }
  }
  return total  // 0 表示无商品层面限制（仅受订单金额和用户余额限制）
})

const previewPointsUsed = computed(() => {
  if (!usePoints.value || !pointsAllowed.value) return 0
  const afterDiscount = Math.max(0, originPrice.value - selectedCouponDiscount.value)
  let max = afterDiscount
  if (goodMaxPoints.value > 0 && goodMaxPoints.value < max) {
    max = goodMaxPoints.value
  }
  return Math.min(userPoints.value, max)
})

const totalPrice = computed(() => {
  let price = originPrice.value - selectedCouponDiscount.value
  if (usePoints.value) price -= previewPointsUsed.value
  return Math.max(0, price)
})

const maxDeductDisplay = computed(() => {
  if (!pointsAllowed.value) return '0.00'
  const afterDiscount = Math.max(0, originPrice.value - selectedCouponDiscount.value)
  let max = afterDiscount
  if (goodMaxPoints.value > 0 && goodMaxPoints.value < max) {
    max = goodMaxPoints.value
  }
  return (Math.min(userPoints.value, max) / 100).toFixed(2)
})

const formatSpecs = (specs, attrs) => {
  const arr = [...(Array.isArray(specs) ? specs : []), ...(Array.isArray(attrs) ? attrs : [])]
  if (!arr.length) return ''
  return arr.map(s => {
    const label = $lt.value(s.labelI18n || s.nameI18n) || $lt.value(s.label || s.name) || s.label || s.name || ''
    const value = $lt.value(s.valueI18n) || $lt.value(s.value) || s.value || ''
    return label ? `${label}: ${value}` : value
  }).join('  ')
}

/* =================== 地址 =================== */
const loadAddress = async () => {
  try {
    const res = await getDefaultAddress()
    if (res.code === 0 && res.data) {
      address.value = res.data
      hasAddress.value = !!(res.data.name && res.data.phone)
    }
  } catch (e) { /* ignore */ }
}

const toAddress = () => {
  uni.navigateTo({ url: '/pages/address/address?from=orderInfo' })
}

/* =================== 优惠券 =================== */
const openCouponPopup = async () => {
  couponShow.value = true
  await loadCoupons()
}
const closeCouponPopup = () => { couponShow.value = false }

const loadCoupons = async () => {
  try {
    const res = await getAllClaimCoupon({ goodIds: goodsList.value.map(g => g.good.ID) })
    if (res.code === 0 && res.data) {
      // 只显示可用的优惠券：未使用、未过期、商品可用
      couponList.value = (Array.isArray(res.data) ? res.data : res.data.list || []).filter(c => {
        if (c.used || c.expired || c.canUse === 0) return false
        return true
      })
    }
  } catch (e) { /* ignore */ }
}

const onCouponTap = async (item) => {
  // 取消选择
  if (selectedCouponNum.value === item.couponNum) {
    selectedCouponNum.value = ''
    selectedCouponDiscount.value = 0
    closeCouponPopup()
    return
  }
  // 需要先领取
  if (!item.couponNum) {
    uni.showLoading({ title: $t.value('claiming') || '...', mask: true })
    try {
      const res = await claimCouponByUser({ couponID: item.couponID })
      item.couponNum = res.data
    } catch (e) { uni.hideLoading(); return }
    uni.hideLoading()
    uni.showToast({ title: $t.value('claimSuccess') || '✓', icon: 'success' })
    return
  }
  // 选择优惠券（本地设置）
  selectedCouponNum.value = item.couponNum
  selectedCouponDiscount.value = item.discount || 0
  closeCouponPopup()
}

const formatCouponDate = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
}

/* =================== 积分 =================== */
const loadUserPoints = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) {
      userPoints.value = res.data.point || 0
    }
  } catch (e) { /* ignore */ }
}

const togglePoints = () => {
  if (!pointsAllowed.value) {
    uni.showToast({ title: $t.value('pointsNotAllowed'), icon: 'none' })
    return
  }
  if (userPoints.value <= 0) {
    uni.showToast({ title: $t.value('noPointsAvailable'), icon: 'none' })
    return
  }
  usePoints.value = !usePoints.value
}

/* =================== 初始化 =================== */
onLoad(async (options) => {
  if (options.goodID && options.skuID) {
    paramGoodID.value = Number(options.goodID)
    paramSkuID.value = Number(options.skuID)
    paramQuantity.value = Number(options.quantity) || 1
    paramCouponNum.value = (options.couponNum && options.couponNum !== '0') ? options.couponNum : ''
    selectedCouponNum.value = paramCouponNum.value
    await loadGoodInfo()
  }
  loadAddress()
  loadUserPoints()
})

// 每次页面显示时刷新地址（用户可能从地址页返回 / 选择了指定地址）
onShow(() => {
  const selected = uni.getStorageSync('selectedAddress')
  if (selected) {
    uni.removeStorageSync('selectedAddress')
    address.value = selected
    hasAddress.value = !!(selected.name && selected.phone)
  } else {
    loadAddress()
  }
})

const loadGoodInfo = async () => {
  try {
    const res = await findGood(paramGoodID.value)
    if (res.code === 0 && res.data) {
      const good = res.data.regood || res.data.good || res.data
      const skus = good.skus || good.SKUS || []
      const matchedSku = skus.find(s => s.ID === paramSkuID.value)
      if (matchedSku) {
        goodsList.value = [{
          good: good,
          sku: matchedSku,
          quantity: paramQuantity.value
        }]
        originPrice.value = matchedSku.price * paramQuantity.value
        // 如果有预选优惠券，找到其折扣额
        if (selectedCouponNum.value) {
          await loadCoupons()
          const found = couponList.value.find(c => c.couponNum === selectedCouponNum.value)
          if (found) selectedCouponDiscount.value = found.discount || 0
        }
      }
    }
  } catch (e) {
    uni.showToast({ title: $t.value('loadFail') || 'Load failed', icon: 'none' })
  }
}

/* =================== 提交订单 =================== */
const submitOrder = async () => {
  if (!hasAddress.value) {
    uni.showToast({ title: $t.value('addressMissing'), icon: 'none' })
    setTimeout(() => { uni.navigateTo({ url: '/pages/address/address' }) }, 500)
    return
  }
  if (goodsList.value.length === 0) return

  // 先选择支付方式，再创建订单
  showPayMethodSelect()
}

const showPayMethodSelect = () => {
  uni.showActionSheet({
    itemList: [$t.value('payByQrcode'), $t.value('payByContact')],
    success: async (sheetRes) => {
      await doCreateOrder(sheetRes.tapIndex)
    },
    fail: () => {
      // 用户取消选择，不做任何操作
    }
  })
}

const doCreateOrder = async (payMethodIndex) => {
  submitting.value = true
  const payMethodMap = { 0: 'qrcode', 1: 'contact' }
  try {
    // 1. 创建订单（携带地址和支付方式）
    const orderData = {
      couponNum: selectedCouponNum.value,
      payMethod: payMethodMap[payMethodIndex] || '',
      name: address.value.name || '',
      phone: address.value.phone || '',
      province: address.value.provinceStr || '',
      city: address.value.cityStr || '',
      area: address.value.areaStr || '',
      street: address.value.street || '',
      detail: goodsList.value.map(d => ({
        goodID: d.good.ID,
        skuID: d.sku.ID,
        quantity: d.quantity,
      }))
    }
    const res = await placeOrder(orderData)
    if (res.code !== 0 || !res.data?.orderID) {
      uni.showToast({ title: res.msg || $t.value('orderCreateFail'), icon: 'none' })
      submitting.value = false
      return
    }
    const newOrderID = String(res.data.orderID)

    // 2. 如果使用积分，调用积分抵扣
    if (usePoints.value) {
      try {
        await changeOrderPoints({ orderID: newOrderID, usePoints: true })
      } catch (e) { /* 积分失败不阻断流程 */ }
    }

    // 3. 检查是否需要支付
    const needPayRes = await checkNeedPay({ orderID: Number(newOrderID), openid: uni.getStorageSync('openid') })
    if (!needPayRes.data) {
      // 0元购，无需支付
      uni.showToast({ title: $t.value('orderSubmitSuccess'), icon: 'success', duration: 2000 })
      setTimeout(() => { uni.redirectTo({ url: '/pages/order/order' }) }, 2000)
      return
    }

    // 4. 按选择的支付方式跳转（统一用 redirectTo 离开本页，防止返回重复提交）
    if (payMethodIndex === 0) {
      // 二维码支付
      uni.redirectTo({ url: `/pages/pay/index?amount=${(totalPrice.value / 100).toFixed(2)}&orderNo=${newOrderID}&orderId=${newOrderID}` })
    } else if (payMethodIndex === 1) {
      // 联系客服支付：redirectTo 替换本页，防止返回重复提交
      uni.setClipboardData({
        data: String(newOrderID),
        success: () => {
          uni.showToast({ title: $t.value('orderNoCopied') + '：' + newOrderID, icon: 'none', duration: 2000 })
          setTimeout(() => { uni.redirectTo({ url: '/pages/kefu/index' }) }, 1500)
        }
      })
    }
  } catch (e) {
    uni.showToast({ title: $t.value('orderCreateFail'), icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss">
page { background-color: #000; }

.nf-orderinfo { min-height: 100vh; background: #000; position: relative; }

.nf-orderinfo-bg {
  position: fixed; top: 0; left: 0; right: 0; height: 500rpx; z-index: 0; pointer-events: none;
  background: radial-gradient(ellipse at 50% 0%, rgba(229, 9, 20, 0.10) 0%, transparent 60%);
}

.nf-navbar {
  background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
  padding: 0 28rpx 16rpx; position: sticky; top: 0; z-index: 99;
}
.nf-navbar-status { height: var(--status-bar-height, 0px); }
.nf-navbar-content { display: flex; align-items: center; justify-content: space-between; height: 88rpx; }
.nf-navbar-back {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: rgba(255, 255, 255, 0.06); border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-title { font-size: 34rpx; font-weight: 700; color: #fff; letter-spacing: 2rpx; }

/* 创建中loading */
.nf-creating {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding-top: 200rpx; gap: 24rpx;
}
.nf-creating-spinner {
  width: 64rpx; height: 64rpx; border: 4rpx solid rgba(255,255,255,0.1);
  border-top-color: #e50914; border-radius: 50%; animation: nf-spin 0.8s linear infinite;
}
@keyframes nf-spin { to { transform: rotate(360deg); } }
.nf-creating-text { font-size: 28rpx; color: rgba(255,255,255,0.5); }

.nf-body { padding: 20rpx 24rpx; position: relative; z-index: 1; }

.nf-card {
  background: rgba(255, 255, 255, 0.04); border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx; margin-bottom: 20rpx; backdrop-filter: blur(8px); overflow: hidden;
}

/* 地址卡片 */
.nf-address-card {
  display: flex; align-items: center; padding: 28rpx; gap: 16rpx;
}
.nf-address-icon { font-size: 40rpx; }
.nf-address-info { flex: 1; }
.nf-address-name { font-size: 28rpx; font-weight: 600; color: #fff; display: block; margin-bottom: 6rpx; }
.nf-address-detail { font-size: 24rpx; color: rgba(255, 255, 255, 0.4); display: block; }

/* 商品卡片 */
.nf-goods-card { padding: 24rpx; }
.nf-goods-item {
  display: flex; gap: 20rpx; padding: 12rpx 0;
  & + .nf-goods-item { border-top: 1rpx solid rgba(255, 255, 255, 0.04); }
}
.nf-goods-img {
  width: 168rpx; height: 168rpx; border-radius: 12rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.06); flex-shrink: 0;
}
.nf-goods-info { flex: 1; display: flex; flex-direction: column; justify-content: space-between; }
.nf-goods-name { font-size: 24rpx; color: rgba(255, 255, 255, 0.5); margin-bottom: 4rpx; }
.nf-goods-desc { font-size: 28rpx; color: #fff; font-weight: 500; margin-bottom: 6rpx; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.nf-goods-specs { font-size: 22rpx; color: rgba(255, 255, 255, 0.3); margin-bottom: 8rpx; }
.nf-goods-bottom { display: flex; justify-content: space-between; align-items: center; }
.nf-goods-price { font-size: 28rpx; font-weight: 700; color: #e50914; }
.nf-goods-qty { font-size: 24rpx; color: rgba(255, 255, 255, 0.4); }

/* 选项卡片 */
.nf-option-card { padding: 24rpx 28rpx; }
.nf-option-row { display: flex; justify-content: space-between; align-items: center; }
.nf-option-label { display: flex; align-items: center; gap: 12rpx; font-size: 28rpx; color: #fff; }
.nf-option-dot { width: 16rpx; height: 16rpx; border-radius: 4rpx; }
.nf-option-value { display: flex; align-items: center; gap: 8rpx; }
.nf-discount-text { font-size: 28rpx; color: #e50914; font-weight: 600; }
.nf-discount-hint { font-size: 26rpx; color: rgba(255, 255, 255, 0.4); }

/* 价格卡片 */
.nf-price-card { padding: 20rpx 28rpx; }
.nf-price-row {
  display: flex; justify-content: space-between; align-items: center;
  padding: 14rpx 0; border-bottom: 1rpx solid rgba(255, 255, 255, 0.04);
  &:last-child { border-bottom: none; }
}
.nf-price-label { font-size: 26rpx; color: rgba(255, 255, 255, 0.5); }
.nf-price-val { font-size: 26rpx; color: #fff; }
.nf-price-discount .nf-price-val { color: #e50914; }
.nf-points-left { display: flex; align-items: center; }
.nf-points-check { display: flex; align-items: center; font-size: 26rpx; color: rgba(255, 255, 255, 0.6); }
.nf-points-amount { font-size: 24rpx; color: rgba(255, 255, 255, 0.4); }
.nf-points-detail {
  padding: 8rpx 0 4rpx; font-size: 22rpx; color: rgba(255,255,255,0.3);
  display: flex; align-items: center;
}
.nf-points-used { color: #e50914; }

/* 底部支付栏 */
.nf-footer {
  position: fixed; bottom: 0; left: 0; right: 0; z-index: 99;
  display: flex; align-items: center; height: 110rpx;
  background: rgba(0, 0, 0, 0.95); backdrop-filter: blur(24px);
  border-top: 1rpx solid rgba(255, 255, 255, 0.06);
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}
.nf-footer-info { flex: 1; padding-left: 32rpx; display: flex; align-items: baseline; gap: 8rpx; }
.nf-footer-label { font-size: 26rpx; color: rgba(255, 255, 255, 0.5); }
.nf-footer-price { font-size: 38rpx; font-weight: 700; color: #e50914; }
.nf-footer-btn {
  width: 240rpx; height: 100%; background: #e50914;
  display: flex; justify-content: center; align-items: center;
  font-size: 30rpx; font-weight: 700; color: #fff;
  &:active { background: #b30710; }
}

/* Netflix 优惠券弹出层 */
.nf-mask {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.7); z-index: 900;
}
.nf-coupon-popup {
  position: fixed; left: 0; right: 0; bottom: -100vh; z-index: 999;
  background: #141414; border-radius: 24rpx 24rpx 0 0;
  transition: all 0.3s ease;
  &.show { bottom: 0; }
}
.nf-coupon-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 28rpx 32rpx 16rpx; border-bottom: 1rpx solid rgba(255,255,255,0.06);
}
.nf-coupon-title { font-size: 32rpx; font-weight: 700; color: #fff; }
.nf-coupon-close {
  width: 56rpx; height: 56rpx; border-radius: 50%;
  background: rgba(255,255,255,0.08); display: flex;
  align-items: center; justify-content: center;
  font-size: 28rpx; color: rgba(255,255,255,0.6);
}
.nf-coupon-scroll { width: 100%; height: 55vh; padding: 16rpx 24rpx; box-sizing: border-box; }
.nf-coupon-empty {
  display: flex; align-items: center; justify-content: center;
  height: 200rpx; color: rgba(255,255,255,0.3); font-size: 28rpx;
}
.nf-coupon-item {
  display: flex; margin-bottom: 20rpx; border-radius: 16rpx; overflow: hidden;
  border: 1rpx solid rgba(255,255,255,0.06); background: rgba(255,255,255,0.04);
  transition: all 0.2s;
  &.nf-coupon-selected { border-color: rgba(229,9,20,0.5); background: rgba(229,9,20,0.08); }
}
.nf-coupon-left {
  width: 200rpx; display: flex; flex-direction: column;
  align-items: center; justify-content: center; padding: 24rpx 16rpx;
  background: rgba(229,9,20,0.12); flex-shrink: 0;
}
.nf-coupon-amount { font-size: 40rpx; font-weight: 800; color: #e50914; }
.nf-coupon-cond { font-size: 20rpx; color: rgba(255,255,255,0.4); margin-top: 4rpx; }
.nf-coupon-right {
  flex: 1; padding: 20rpx 24rpx; display: flex; flex-direction: column; justify-content: center;
}
.nf-coupon-name { font-size: 28rpx; color: #fff; font-weight: 500; margin-bottom: 8rpx; }
.nf-coupon-exp { font-size: 22rpx; color: rgba(255,255,255,0.3); }
.nf-coupon-deselect {
  margin-top: 8rpx; display: inline-flex; align-self: flex-start;
  padding: 4rpx 16rpx; border-radius: 8rpx; font-size: 22rpx;
  color: #e50914; border: 1rpx solid rgba(229,9,20,0.4); background: rgba(229,9,20,0.06);
}
.nf-coupon-claim {
  margin-top: 8rpx; display: inline-flex; align-self: flex-start;
  padding: 4rpx 16rpx; border-radius: 8rpx; font-size: 22rpx;
  color: #e5a609; border: 1rpx solid rgba(229,166,9,0.4); background: rgba(229,166,9,0.06);
}
</style>
