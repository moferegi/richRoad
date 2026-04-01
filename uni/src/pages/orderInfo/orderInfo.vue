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
        <text class="nf-navbar-title">{{ $t('orderDetail') || '订单详情' }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body">
      <!-- 收货地址 -->
      <view class="nf-card nf-address-card" @tap="toAddress">
        <view class="nf-address-icon">📍</view>
        <view class="nf-address-info" v-if="hasAddress">
          <text class="nf-address-name">{{ data.name }} {{ data.phone }}</text>
          <text class="nf-address-detail">{{ data.province }}{{ data.city }}{{ data.area }} {{ data.street }}</text>
        </view>
        <view class="nf-address-info" v-else>
          <text class="nf-address-name">请选择收货地址</text>
          <text class="nf-address-detail">添加收货人信息及地址</text>
        </view>
        <uni-icons type="right" size="16" color="rgba(255,255,255,0.3)"></uni-icons>
      </view>

      <!-- 商品列表 -->
      <view class="nf-card nf-goods-card">
        <view class="nf-goods-item" v-for="(d, i) in data.detail" :key="i">
          <image class="nf-goods-img" :src="getUrl(d.sku.picture)" mode="aspectFill"></image>
          <view class="nf-goods-info">
            <text class="nf-goods-name">{{ d?.sku?.name }}</text>
            <text class="nf-goods-desc">{{ d?.good?.description }}</text>
            <text class="nf-goods-specs">{{ d?.sku?.specs?.map(i => i.value).join(' ') }}</text>
            <view class="nf-goods-bottom">
              <text class="nf-goods-price">¥{{ d.sku.price / 100 }}</text>
              <text class="nf-goods-qty">×{{ d.quantity }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 优惠券选择 -->
      <view class="nf-card nf-option-card" @tap="opencoupon">
        <view class="nf-option-row">
          <view class="nf-option-label">
            <text class="nf-option-dot" style="background: #e50914;"></text>
            <text>优惠券</text>
          </view>
          <view class="nf-option-value">
            <text class="nf-discount-text" v-if="data.discount > 0">-¥{{ data.discount / 100 }}</text>
            <text class="nf-discount-hint" v-else>选择优惠券</text>
            <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)"></uni-icons>
          </view>
        </view>
      </view>

      <!-- 价格汇总 -->
      <view class="nf-card nf-price-card">
        <view class="nf-price-row">
          <text class="nf-price-label">商品金额</text>
          <text class="nf-price-val">¥{{ data.originPrice / 100 }}</text>
        </view>
        <view class="nf-price-row nf-price-discount">
          <text class="nf-price-label">优惠金额</text>
          <text class="nf-price-val">-¥{{ data.discount / 100 }}</text>
        </view>
        <view class="nf-price-row nf-price-points">
          <view class="nf-points-left">
            <checkbox-group @change="onPointsChange">
              <view class="nf-points-check">
                <checkbox value="points" :checked="usePoints" color="#e50914" style="transform: scale(0.7); margin-right: 6rpx;" />
                <text>积分抵扣</text>
              </view>
            </checkbox-group>
          </view>
          <text class="nf-points-amount">可抵扣 ¥{{ Math.min(availablePointsAmount, (data.originPrice - data.discount)) / 100 }}</text>
        </view>
      </view>

      <!-- 订单信息 -->
      <view class="nf-card nf-info-card" v-if="data.ID">
        <view class="nf-info-row">
          <text class="nf-info-label">{{ $t('orderNo') || '订单编号' }}</text>
          <view class="nf-info-value" @tap="copyOrderNo">
            <text>{{ data.ID }}</text>
            <text class="nf-copy-btn">复制</text>
          </view>
        </view>
        <view class="nf-info-row">
          <text class="nf-info-label">下单时间</text>
          <text class="nf-info-value">{{ formatTime(data.CreatedAt) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.payMethod">
          <text class="nf-info-label">支付方式</text>
          <text class="nf-info-value">{{ data.payMethod === 'qrcode' ? ($t('payByQrcode') || 'QR码支付') : ($t('payByContact') || '联系客服') }}</text>
        </view>
        <view class="nf-info-row" v-if="data.paidAt">
          <text class="nf-info-label">付款时间</text>
          <text class="nf-info-value">{{ formatTime(data.paidAt) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.express">
          <text class="nf-info-label">快递单号</text>
          <text class="nf-info-value">{{ data.express }}</text>
        </view>
        <view class="nf-info-row" v-if="data.receivedAt">
          <text class="nf-info-label">收货时间</text>
          <text class="nf-info-value">{{ formatTime(data.receivedAt) }}</text>
        </view>
        <view class="nf-info-row" v-if="data.status === '0' && data.closeTime">
          <text class="nf-info-label">剩余支付时间</text>
          <text class="nf-info-value nf-countdown-val">{{ payCountdown }}</text>
        </view>
      </view>

      <!-- 退款操作 -->
      <view class="nf-card nf-refund-card" v-if="canApplyRefund || isRefunding || isRefunded">
        <view v-if="canApplyRefund" class="nf-refund-btn" @tap="openRefund">申请退款</view>
        <view v-else class="nf-refund-status">{{ isRefunding ? '退款处理中' : '已退款' }}</view>
      </view>

      <view style="height: 140rpx;"></view>
    </view>

    <!-- 底部支付栏 -->
    <view class="nf-footer">
      <view class="nf-footer-info">
        <text class="nf-footer-label">实付款</text>
        <text class="nf-footer-price">¥{{ (totalPrice / 100).toFixed(2) }}</text>
      </view>
      <view class="nf-footer-btn" @tap="tapPay">
        <text>提交订单</text>
      </view>
    </view>

    <!-- 优惠券弹出层 -->
    <view class="nf-mask" v-if="couponshow" @tap="hidecoupon"></view>
    <view class="nf-coupon-popup" :class="{ show: couponshow }">
      <scroll-view class="nf-coupon-scroll" scroll-y>
        <cc-defineCoupon v-if="couponshow" colors="#e50914" @onReceive="onReceive"></cc-defineCoupon>
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
import { ref, nextTick, watch, computed, onUnmounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { selfOrder, changeOrderCoupon } from '@/api/order.js'
import { claimCouponByUser } from '@/api/coupon.js'
import { getPayParams, getOrderById, checkNeedPay } from '@/api/base.js'
import { getUserInfo } from '@/api/base.js'
import { getUrl } from "@/utils/url.js"
import { changeOrderPoints } from '@/api/order.js'
import RefundApplyPopup from '@/components/refund-apply-popup/refund-apply-popup.vue'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const goBack = () => { uni.navigateBack() }

const toAddress = () => {
  uni.navigateTo({ url: `/pages/address/address?ID=${data.value.ID}` })
}

const totalPrice = ref(0)
const hasAddress = ref(false)
const usePoints = ref(false)
const availablePointsAmount = ref(0)
const userPoints = ref(0)
const data = ref({
  detail: [],
  originPrice: 0,
  discount: 0,
})
const orderID = ref("")
const refundVisible = ref(false)
const couponshow = ref(false)

const opencoupon = () => { couponshow.value = true }
const hidecoupon = () => { couponshow.value = false }

const onReceive = async (item, index) => {
  uni.showLoading({ title: item.couponNum == 0 ? '领取中...' : '选择中...', mask: true })
  try {
    if (item.couponNum == 0) {
      const res = await claimCouponByUser({ couponID: item.couponID })
      item.couponNum = res.data
    }
    await changeOrderCoupon({ orderID: orderID.value, couponNum: item.couponNum })
    setTimeout(() => { initSingleOrder(); hidecoupon() }, 500)
  } catch (error) { /* ignore */ }
  finally { uni.hideLoading() }
}

const getUserPoints = async () => {
  try {
    const res = await getUserInfo()
    if (res.code === 0) {
      userPoints.value = res.data.point || 0
      availablePointsAmount.value = userPoints.value
    }
  } catch (error) { console.error('获取用户积分失败:', error) }
}

onLoad((options) => {
  orderID.value = options.orderID
  initSingleOrder()
  getUserPoints()
})

const canApplyRefund = computed(() => ['1', '2', '3', '7'].includes(String(data.value.status || '')))
const isRefunding = computed(() => String(data.value.status || '') === '6')
const isRefunded = computed(() => String(data.value.status || '') === '5')
const openRefund = () => { refundVisible.value = true }

const onPointsChange = async (e) => {
  const isChecked = e.detail.value.includes('points')
  usePoints.value = isChecked
  try {
    const res = await changeOrderPoints({ orderID: orderID.value, usePoints: usePoints.value })
    if (res.code === 0) { await initSingleOrder() }
    else { uni.showToast({ title: '积分抵扣更新失败', icon: 'none' }); usePoints.value = !isChecked }
  } catch (error) {
    uni.showToast({ title: '积分抵扣更新失败', icon: 'none' }); usePoints.value = !isChecked
  }
}

const initSingleOrder = async () => {
  const order = await selfOrder(orderID.value)
  if (order.code === 0) {
    data.value = order.data
    totalPrice.value = order.data.totalPrice
    if (order.data.city) hasAddress.value = true
    startPayCountdown()
  }
}

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const copyOrderNo = () => {
  uni.setClipboardData({
    data: String(data.value.ID),
    success: () => { uni.showToast({ title: $t.value('orderNoCopied') || '已复制', icon: 'none' }) }
  })
}

const payCountdown = ref('')
let payTimer = null

const startPayCountdown = () => {
  if (payTimer) clearInterval(payTimer)
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
  if (remain <= 0) { payCountdown.value = '已超时'; if (payTimer) clearInterval(payTimer); return }
  const m = Math.floor(remain / 60), s = remain % 60
  payCountdown.value = `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

onUnmounted(() => { if (payTimer) clearInterval(payTimer) })

const tapPay = async () => {
  if (!data.value.name || !data.value.phone || !data.value.province || !data.value.city || !data.value.area || !data.value.detail) {
    uni.navigateTo({ url: `/pages/address/address?ID=${data.value.goodID || data.value.ID}` })
    return
  }
  const params = { orderID: Number(orderID.value), openid: uni.getStorageSync('openid') }
  const needPayRes = await checkNeedPay(params)
  if (!needPayRes.data) {
    uni.showToast({ title: $t.value('orderSubmitSuccess') || '订单提交成功！', icon: 'success', duration: 2000 })
    setTimeout(() => { uni.navigateTo({ url: `/pages/order/order?orderID=${orderID.value}` }) }, 2000)
    return
  }
  showPayMethodSelect()
}

const showPayMethodSelect = () => {
  uni.showActionSheet({
    itemList: [$t.value('payByQrcode') || 'QR码支付', $t.value('payByContact') || '联系客服付款'],
    success: (res) => {
      if (res.tapIndex === 0) {
        uni.navigateTo({ url: `/pages/pay/index?amount=${(totalPrice.value / 100).toFixed(2)}&orderNo=${data.value.ID}&orderId=${orderID.value}` })
      } else if (res.tapIndex === 1) {
        uni.setClipboardData({
          data: String(data.value.ID),
          success: () => {
            uni.showToast({ title: ($t.value('orderNoCopied') || '订单号已复制') + '：' + data.value.ID, icon: 'none', duration: 2000 })
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

/* 优惠券弹出层 */
.nf-mask {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.7); z-index: 900;
}
.nf-coupon-popup {
  position: fixed; left: 0; right: 0; bottom: -100vh; z-index: 999;
  background: #1a1a1a; border-radius: 24rpx 24rpx 0 0;
  transition: all 0.3s ease;
  &.show { bottom: 0; }
}
.nf-coupon-scroll { width: 100vw; height: 60vh; padding-top: 16rpx; }
</style>
