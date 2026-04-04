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
        <text class="nf-navbar-title">{{ orderReady ? ($t('orderDetail')) : ($t('orderConfirm')) }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 创建中loading -->
    <view v-if="creating" class="nf-creating">
      <view class="nf-creating-spinner"></view>
      <text class="nf-creating-text">{{ $t('orderCreating') }}</text>
    </view>

    <view class="nf-body" v-else>
      <!-- 收货地址 -->
      <view class="nf-card nf-address-card" @tap="toAddress">
        <view class="nf-address-icon">📍</view>
        <view class="nf-address-info" v-if="hasAddress">
          <text class="nf-address-name">{{ data.name }} {{ data.phone }}</text>
          <text class="nf-address-detail">{{ data.province }}{{ data.city }}{{ data.area }} {{ data.street }}</text>
        </view>
        <view class="nf-address-info" v-else>
          <text class="nf-address-name">{{ $t('selectAddress') }}</text>
          <text class="nf-address-detail">{{ $t('addAddressHint') }}</text>
        </view>
        <uni-icons type="right" size="16" color="rgba(255,255,255,0.3)"></uni-icons>
      </view>

      <!-- 商品列表 -->
      <view class="nf-card nf-goods-card">
        <view class="nf-goods-item" v-for="(d, i) in data.detail" :key="i">
          <image class="nf-goods-img" :src="getUrl(d.sku?.picture)" mode="aspectFill"></image>
          <view class="nf-goods-info">
            <text class="nf-goods-name">{{ $lt(d?.sku?.name) || d?.sku?.name }}</text>
            <text class="nf-goods-desc">{{ $lt(d?.good?.description) || d?.good?.description }}</text>
            <text class="nf-goods-specs">{{ formatSpecs(d?.sku?.specs) }}</text>
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
            <text class="nf-discount-text" v-if="data.discount > 0">-{{ cs }}{{ data.discount / 100 }}</text>
            <text class="nf-discount-hint" v-else>{{ $t('selectCoupon') }}</text>
            <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)"></uni-icons>
          </view>
        </view>
      </view>

      <!-- 价格汇总 -->
      <view class="nf-card nf-price-card">
        <view class="nf-price-row">
          <text class="nf-price-label">{{ $t('productAmount') }}</text>
          <text class="nf-price-val">{{ cs }}{{ (data.originPrice || 0) / 100 }}</text>
        </view>
        <view class="nf-price-row nf-price-discount" v-if="data.discount > 0">
          <text class="nf-price-label">{{ $t('discountAmount') }}</text>
          <text class="nf-price-val">-{{ cs }}{{ data.discount / 100 }}</text>
        </view>
        <view class="nf-price-row nf-price-points">
          <view class="nf-points-left">
            <checkbox-group @change="onPointsChange">
              <view class="nf-points-check">
                <checkbox value="points" :checked="usePoints" color="#e50914" style="transform: scale(0.7); margin-right: 6rpx;" />
                <text>{{ $t('pointsDeduction') }}</text>
              </view>
            </checkbox-group>
          </view>
          <text class="nf-points-amount">{{ $t('maxPointsDeduct') }} {{ cs }}{{ maxDeductDisplay }}</text>
        </view>
        <view class="nf-points-detail" v-if="userPoints > 0">
          <text>{{ $t('pointsBalance') }}: {{ userPoints }}</text>
          <text v-if="usePoints && data.pointsUsed > 0" class="nf-points-used">  -{{ data.pointsUsed }} = {{ userPoints - (data.pointsUsed || 0) }}</text>
        </view>
      </view>

      <!-- 订单信息 (仅已创建的订单显示) -->
      <view class="nf-card nf-info-card" v-if="orderReady && data.ID">
        <view class="nf-info-row">
          <text class="nf-info-label">{{ $t('orderNo') }}</text>
          <view class="nf-info-value" @tap="copyOrderNo">
            <text>{{ data.ID }}</text>
            <text class="nf-copy-btn">{{ $t('copy') }}</text>
          </view>
        </view>
        <view class="nf-info-row">
          <text class="nf-info-label">{{ $t('orderTime') }}</text>
          <text class="nf-info-value">{{ formatTime(data.CreatedAt) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.payMethod">
          <text class="nf-info-label">{{ $t('paymentMethod') }}</text>
          <text class="nf-info-value">{{ data.payMethod === 'qrcode' ? ($t('payByQrcode')) : ($t('payByContact')) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.paidAt">
          <text class="nf-info-label">{{ $t('paymentTime') }}</text>
          <text class="nf-info-value">{{ formatTime(data.paidAt) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.express">
          <text class="nf-info-label">{{ $t('trackingNo') }}</text>
          <text class="nf-info-value">{{ data.express }}</text>
        </view>
        <view class="nf-info-row" v-if="data.receivedAt">
          <text class="nf-info-label">{{ $t('deliveryTime') }}</text>
          <text class="nf-info-value">{{ formatTime(data.receivedAt) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.status === '0' && data.closeTime">
          <text class="nf-info-label">{{ $t('remainPayTime') }}</text>
          <text class="nf-info-value nf-countdown-val">{{ payCountdown }}</text>
        </view>
      </view>

      <!-- 退款操作 -->
      <view class="nf-card nf-refund-card" v-if="canApplyRefund || isRefunding || isRefunded">
        <view v-if="canApplyRefund" class="nf-refund-btn" @tap="openRefund">{{ $t('applyRefund') }}</view>
        <view v-else class="nf-refund-status">{{ isRefunding ? $t('refundProcessing') : $t('refunded') }}</view>
      </view>

      <view style="height: 140rpx;"></view>
    </view>

    <!-- 底部支付栏 -->
    <view class="nf-footer" v-if="!creating">
      <view class="nf-footer-info">
        <text class="nf-footer-label">{{ $t('actualPayment') }}</text>
        <text class="nf-footer-price">{{ cs }}{{ (totalPrice / 100).toFixed(2) }}</text>
      </view>
      <view class="nf-footer-btn" @tap="tapPay">
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
              class="nf-coupon-item" :class="{ 'nf-coupon-selected': data.couponNum && data.couponNum === c.couponNum }"
              @tap="onCouponTap(c)">
          <view class="nf-coupon-left">
            <text class="nf-coupon-amount">{{ cs }}{{ (c.discount || 0) / 100 }}</text>
            <text class="nf-coupon-cond" v-if="c.minAmount > 0">{{ $t('couponNoMin') || '' }} {{ cs }}{{ c.minAmount / 100 }}</text>
            <text class="nf-coupon-cond" v-else>{{ $t('couponNoMin') || '' }}</text>
          </view>
          <view class="nf-coupon-right">
            <text class="nf-coupon-name">{{ $lt(c.name) || c.name }}</text>
            <text class="nf-coupon-exp">{{ $t('couponExpiry') || '' }} {{ formatCouponDate(c.endTime) }}</text>
            <view v-if="data.couponNum && data.couponNum === c.couponNum" class="nf-coupon-deselect">
              <text>{{ $t('cancelCoupon') }}</text>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>

    <refund-apply-popup
      v-model:visible="refundVisible"
      :order-id="orderID"
      @success="initSingleOrder"
    />
  </view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { selfOrder, changeOrderCoupon, placeOrder } from '@/api/order.js'
import { claimCouponByUser, getAllClaimCoupon } from '@/api/coupon.js'
import { checkNeedPay } from '@/api/base.js'
import { getUserInfo } from '@/api/base.js'
import { getUrl } from "@/utils/url.js"
import { changeOrderPoints } from '@/api/order.js'
import RefundApplyPopup from '@/components/refund-apply-popup/refund-apply-popup.vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const goBack = () => { uni.navigateBack() }

const toAddress = () => {
  uni.navigateTo({ url: `/pages/address/address?ID=${data.value.ID}` })
}

const totalPrice = ref(0)
const hasAddress = ref(false)
const usePoints = ref(false)
const availablePointsAmount = ref(0)
const userPoints = ref(0)
const creating = ref(false)
const orderReady = ref(false)
const data = ref({
  detail: [],
  originPrice: 0,
  discount: 0,
})
const orderID = ref("")
const refundVisible = ref(false)
const couponShow = ref(false)
const couponList = ref([])

const maxDeductDisplay = computed(() => {
  const afterDiscount = (data.value.originPrice || 0) - (data.value.discount || 0)
  return (Math.min(availablePointsAmount.value, Math.max(0, afterDiscount)) / 100).toFixed(2)
})

const formatSpecs = (specs) => {
  if (!specs || !Array.isArray(specs)) return ''
  return specs.map(s => {
    const label = $lt.value(s.label) || s.label || ''
    const value = $lt.value(s.value) || s.value || ''
    return label ? `${label}: ${value}` : value
  }).join('  ')
}

/* =================== 优惠券 =================== */
const openCouponPopup = async () => {
  couponShow.value = true
  await loadCoupons()
}
const closeCouponPopup = () => { couponShow.value = false }

const loadCoupons = async () => {
  try {
    const res = await getAllClaimCoupon({ page: 1, pageSize: 50 })
    if (res.code === 0 && res.data?.list) {
      couponList.value = res.data.list.filter(c => {
        if (c.couponNum && c.used) return false
        return true
      })
    }
  } catch (e) { /* ignore */ }
}

const onCouponTap = async (item) => {
  // 取消选择：如果已选中的是同一个，移除优惠券
  if (data.value.couponNum && data.value.couponNum === item.couponNum) {
    uni.showLoading({ mask: true })
    try {
      await changeOrderCoupon({ orderID: orderID.value, couponNum: '' })
      await initSingleOrder()
      closeCouponPopup()
    } catch (e) { /* ignore */ }
    finally { uni.hideLoading() }
    return
  }
  // 选择优惠券
  uni.showLoading({ title: item.couponNum ? ($t.value('selecting') || '...') : ($t.value('claiming') || '...'), mask: true })
  try {
    if (!item.couponNum) {
      const res = await claimCouponByUser({ couponID: item.couponID })
      item.couponNum = res.data
    }
    await changeOrderCoupon({ orderID: orderID.value, couponNum: item.couponNum })
    await initSingleOrder()
    closeCouponPopup()
  } catch (e) { /* ignore */ }
  finally { uni.hideLoading() }
}

const formatCouponDate = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
}

/* =================== 积分 =================== */
const getUserPoints = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) {
      userPoints.value = res.data.point || 0
      availablePointsAmount.value = userPoints.value
    }
  } catch (e) { /* ignore */ }
}

const onPointsChange = async (e) => {
  const isChecked = e.detail.value.includes('points')
  usePoints.value = isChecked
  try {
    const res = await changeOrderPoints({ orderID: orderID.value, usePoints: usePoints.value })
    if (res.code === 0) { await initSingleOrder() }
    else { uni.showToast({ title: $t.value('pointsFail'), icon: 'none' }); usePoints.value = !isChecked }
  } catch (e) {
    uni.showToast({ title: $t.value('pointsFail'), icon: 'none' }); usePoints.value = !isChecked
  }
}

/* =================== 初始化 =================== */
onLoad(async (options) => {
  getUserPoints()

  if (options.orderID) {
    // 模式1: 查看/操作已有订单
    orderID.value = options.orderID
    orderReady.value = true
    await initSingleOrder()
  } else if (options.goodID && options.skuID) {
    // 模式2: 从SKU弹窗跳过来，需先创建订单
    creating.value = true
    try {
      const orderData = {
        couponNum: (options.couponNum && options.couponNum !== '0') ? options.couponNum : '',
        detail: [{
          goodID: Number(options.goodID),
          skuID: Number(options.skuID),
          quantity: Number(options.quantity) || 1,
        }]
      }
      const res = await placeOrder(orderData)
      if (res.code === 0 && res.data?.orderID) {
        orderID.value = String(res.data.orderID)
        orderReady.value = true
        await initSingleOrder()
      } else {
        uni.showToast({ title: res.msg || $t.value('orderCreateFail'), icon: 'none' })
        setTimeout(() => uni.navigateBack(), 1500)
      }
    } catch (e) {
      uni.showToast({ title: $t.value('orderCreateFail'), icon: 'none' })
      setTimeout(() => uni.navigateBack(), 1500)
    } finally {
      creating.value = false
    }
  }
})

const canApplyRefund = computed(() => ['1', '2', '3', '7'].includes(String(data.value.status || '')))
const isRefunding = computed(() => String(data.value.status || '') === '6')
const isRefunded = computed(() => String(data.value.status || '') === '5')
const openRefund = () => { refundVisible.value = true }

const initSingleOrder = async () => {
  const order = await selfOrder(orderID.value)
  if (order.code === 0) {
    data.value = order.data
    totalPrice.value = order.data.totalPrice
    usePoints.value = !!order.data.usePoints
    hasAddress.value = !!order.data.city
    startPayCountdown()
  }
}

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const copyOrderNo = () => {
  uni.setClipboardData({
    data: String(data.value.ID),
    success: () => { uni.showToast({ title: $t.value('orderNoCopied'), icon: 'none' }) }
  })
}

/* =================== 倒计时 =================== */
const payCountdown = ref('')
let payTimer = null

const startPayCountdown = () => {
  if (payTimer) clearInterval(payTimer)
  if (data.value.status !== '0' || !data.value.closeTime) return
  updatePayCountdown()
  payTimer = setInterval(updatePayCountdown, 1000)
}

const updatePayCountdown = () => {
  if (data.value.status !== '0' || !data.value.closeTime) {
    payCountdown.value = ''
    if (payTimer) clearInterval(payTimer)
    return
  }
  const remain = Math.max(0, Math.floor((new Date(data.value.closeTime).getTime() - Date.now()) / 1000))
  if (remain <= 0) {
    payCountdown.value = $t.value('payTimeout')
    if (payTimer) clearInterval(payTimer)
    // 超时自动跳转订单列表
    uni.showToast({ title: $t.value('payTimeout'), icon: 'none' })
    setTimeout(() => {
      uni.redirectTo({ url: '/pages/order/order' })
    }, 1500)
    return
  }
  const h = Math.floor(remain / 3600)
  const m = Math.floor((remain % 3600) / 60)
  const s = remain % 60
  payCountdown.value = h > 0
    ? `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
    : `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

onUnmounted(() => { if (payTimer) clearInterval(payTimer) })

/* =================== 支付 =================== */
const tapPay = async () => {
  if (!data.value.name || !data.value.phone || !data.value.province) {
    uni.showToast({ title: $t.value('addressMissing'), icon: 'none' })
    setTimeout(() => {
      uni.navigateTo({ url: `/pages/address/address?ID=${data.value.ID}` })
    }, 500)
    return
  }
  const params = { orderID: Number(orderID.value), openid: uni.getStorageSync('openid') }
  const needPayRes = await checkNeedPay(params)
  if (!needPayRes.data) {
    uni.showToast({ title: $t.value('orderSubmitSuccess'), icon: 'success', duration: 2000 })
    setTimeout(() => { uni.navigateTo({ url: `/pages/order/order?orderID=${orderID.value}` }) }, 2000)
    return
  }
  showPayMethodSelect()
}

const showPayMethodSelect = () => {
  uni.showActionSheet({
    itemList: [$t.value('payByQrcode'), $t.value('payByContact')],
    success: (res) => {
      if (res.tapIndex === 0) {
        uni.navigateTo({ url: `/pages/pay/index?amount=${(totalPrice.value / 100).toFixed(2)}&orderNo=${data.value.ID}&orderId=${orderID.value}` })
      } else if (res.tapIndex === 1) {
        uni.setClipboardData({
          data: String(data.value.ID),
          success: () => {
            uni.showToast({ title: $t.value('orderNoCopied') + '：' + data.value.ID, icon: 'none', duration: 2000 })
            setTimeout(() => { uni.navigateTo({ url: '/pages/kefu/index' }) }, 1500)
          }
        })
      }
    }
  })
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

/* 订单信息卡片 */
.nf-info-card { padding: 24rpx 28rpx; }
.nf-info-row {
  display: flex; justify-content: space-between; align-items: center; padding: 10rpx 0;
}
.nf-info-label { font-size: 24rpx; color: rgba(255, 255, 255, 0.4); min-width: 140rpx; }
.nf-info-value { font-size: 24rpx; color: rgba(255, 255, 255, 0.7); text-align: right; flex: 1; }
.nf-copy-btn {
  font-size: 22rpx; color: #e50914; margin-left: 12rpx;
  padding: 4rpx 12rpx; border: 1rpx solid rgba(229, 9, 20, 0.3);
  border-radius: 8rpx; background: rgba(229, 9, 20, 0.1);
}
.nf-countdown-val { color: #e50914; font-weight: 700; }

/* 退款 */
.nf-refund-card { padding: 20rpx 28rpx; display: flex; justify-content: flex-end; }
.nf-refund-btn {
  padding: 12rpx 28rpx; border-radius: 30rpx; font-size: 26rpx; font-weight: 600;
  background: rgba(245, 158, 11, 0.15); border: 1rpx solid rgba(245, 158, 11, 0.3); color: #f59e0b;
}
.nf-refund-status {
  padding: 12rpx 28rpx; border-radius: 30rpx; font-size: 26rpx;
  background: rgba(255, 255, 255, 0.04); border: 1rpx solid rgba(255, 255, 255, 0.1); color: rgba(255, 255, 255, 0.4);
}

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
</style>
