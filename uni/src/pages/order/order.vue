<template>
  <view class="nf-order">
    <view class="nf-order-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('myOrders') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 订单状态标签（横排滚动） -->
    <view class="nf-tabs">
      <scroll-view scroll-x class="nf-tabs-scroll" :show-scrollbar="false">
        <view class="nf-tabs-inner">
          <view
            v-for="(item, index) in tabColumns"
            :key="index"
            class="nf-tab"
            :class="{ active: item.id === activeSataus }"
            @tap="tapBtn(item)"
          >
            <text>{{ item.title }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 空状态 -->
    <view class="nf-empty" v-if="orderList.length === 0">
      <view class="nf-empty-icon">📦</view>
      <text class="nf-empty-text">{{ $t('noOrderData') }}</text>
    </view>

    <!-- 订单列表 -->
    <view class="nf-order-list" v-else>
      <view class="nf-order-card" v-for="(item, index) in orderList" :key="index" @tap="goOrderDetail(item)">
        <!-- 订单头部 -->
        <view class="nf-order-header">
          <view class="nf-order-date">
            {{ formatOrderDate(item.CreatedAt) }}
            <text v-if="item.isPresale" class="nf-presale-tag">{{ $t('presale') }}</text>
          </view>
          <view class="nf-order-status" :class="'nf-st-' + item.status">
            {{ getStatusLabel(item.status) }}
            <text v-if="item.status === '0' && countdownMap[item.ID]" class="nf-countdown"> {{ countdownMap[item.ID] }}</text>
          </view>
        </view>

        <!-- 商品列表（多件横滑） -->
        <scroll-view scroll-x class="nf-goods-scroll" :show-scrollbar="false"
          v-if="item.detail && item.detail.length > 1">
          <view class="nf-goods-row">
            <LazyImage
              v-for="(d, di) in item.detail" :key="di"
              :src="d.sku.externalPicturePath ? getExternalUrl(d.sku.externalPicturePath) : getUrl(d.sku.picture)"
              class="nf-goods-thumb"
              mode="aspectFill"
            />
          </view>
        </scroll-view>

        <!-- 单件商品 -->
        <view class="nf-goods-single" v-if="item.detail && item.detail.length === 1">
          <LazyImage :src="item.detail[0].sku.externalPicturePath ? getExternalUrl(item.detail[0].sku.externalPicturePath) : getUrl(item.detail[0].sku.picture)" class="nf-goods-thumb-lg" mode="aspectFill" />
          <view class="nf-goods-single-info">
            <text class="nf-goods-single-name">{{ $lt(item.detail[0].sku.name) || item.detail[0].sku.name }}</text>
            <text class="nf-goods-single-desc">{{ $lt(item.detail[0].good?.description) || $lt(item.detail[0].sku?.description) || item.detail[0].good?.description || item.detail[0].sku?.description }}</text>
          </view>
        </view>

        <!-- 金额统计 -->
        <view class="nf-order-summary">
          <text class="nf-order-count">{{ $t('totalItems').replace('{n}', item.detail ? item.detail.length : 0) }} · {{ $t('paidAmount') }}</text>
          <text class="nf-order-price">{{ getOrderCurrencySymbol(item) }}{{ getLocalizedPaidAmount(item) }}</text>
        </view>

        <!-- 操作按钮 -->
        <view class="nf-order-actions" @tap.stop>
          <!-- 已取消：只显示标签+再次购买 -->
          <template v-if="item.status === '4'">
            <view class="nf-action-tag cancelled">{{ $t('ordersCancelled') }}</view>
            <view class="nf-action nf-action-ghost" @tap.stop="buyAgain(item)">{{ $t('buyAgain') }}</view>
          </template>

          <!-- 待付款 -->
          <template v-if="item.status === '0'">
            <view class="nf-action nf-action-ghost" @tap.stop="cancelOrder(item)">{{ $t('cancelOrder') }}</view>
            <view class="nf-action nf-action-primary" @tap.stop="payOff(item.ID)">{{ $t('payNow') }}</view>
          </template>

          <template v-if="item.status === '8'">
            <view class="nf-action-tag refunding">{{ $t('ordersPendingConfirm') }}</view>
            <view class="nf-action nf-action-primary" @tap.stop="goKefuWithOrder(item)">{{ $t('contactCustomerService') }}</view>
          </template>

          <!-- 待发货/待收货/待评价/已评价 共用退款 -->
          <view v-if="showRefundBtn && canApplyRefund(item)" class="nf-action nf-action-warn" @tap.stop="openRefund(item)">{{ $t('applyRefund') }}</view>
          <view v-else-if="item.status === '6'" class="nf-action-tag refunding">{{ $t('ordersRefunding') }}</view>
          <view v-else-if="item.status === '5'" class="nf-action-tag refunded">{{ $t('ordersRefunded') }}</view>

          <!-- 待收货 -->
          <template v-if="item.status === '2'">
            <view v-if="showLogisticsBtn && item.express" class="nf-action nf-action-ghost" @tap.stop="trackLogistics(item)">{{ $t('viewLogistics') }}</view>
            <view class="nf-action nf-action-primary" @tap.stop="confirm(item)">{{ $t('confirmReceipt') }}</view>
          </template>

          <!-- 评价 -->
          <template v-if="item.status === '3' || item.status === '7'">
            <template v-if="item.detail && item.detail.length > 1">
              <view v-if="hasUncommentedItems(item)" class="nf-action nf-action-primary" @tap.stop="goCommentAll(item)">{{ $t('reviewOrder') }}</view>
              <view v-else-if="hasCommentedItems(item)" class="nf-action nf-action-ghost" @tap.stop="goCommentAll(item)">{{ $t('viewReview') }}</view>
            </template>
            <template v-else-if="item.detail && item.detail.length === 1">
              <view v-if="!item.detail[0].isComment" class="nf-action nf-action-primary" @tap.stop="goComment(item, item.detail[0])">{{ $t('reviewOrder') }}</view>
              <view v-else class="nf-action nf-action-ghost" @tap.stop="goComment(item, item.detail[0])">{{ $t('viewReview') }}</view>
            </template>
          </template>
        </view>
      </view>
    </view>

    <!-- 加载更多提示 -->
    <view class="nf-load-more" v-if="orderList.length > 0">
      <text class="nf-load-text" v-if="isLoading">{{ $t('loading') }}</text>
      <text class="nf-load-text" v-else-if="isBottom">— {{ $t('noMoreData') }} —</text>
      <text class="nf-load-btn" v-else @tap="loadMore">{{ $t('loadMore') }}</text>
    </view>

    <view style="height: 40rpx;"></view>

    <refund-apply-popup
      v-model:visible="refundVisible"
      :order-id="refundOrderId"
      @success="onRefundSuccess"
    />
  </view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import { updateOrderStatus, SelfOrderList } from "../../api/order"
import { getSysConfigByKey } from '@/api/sysConfig.js'
import { getUrl, getExternalUrl } from "@/utils/url.js"
import RefundApplyPopup from '@/components/refund-apply-popup/refund-apply-popup.vue'
import LazyImage from '@/components/lazy-image/lazy-image.vue'
import { resolveLocalizedPriceFen } from '@/utils/price-i18n.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const getOrderCurrencySymbol = (order) => {
  const snapshot = String(order?.settlementCurrencySymbol || '').trim()
  return snapshot || cs.value
}

const getDetailBasePriceFen = (detail) => {
  const detailPrice = Number(detail?.price)
  if (Number.isFinite(detailPrice) && detailPrice >= 0) {
    return detailPrice
  }
  return 0
}

const getDetailPriceI18nSnapshot = (detail) => {
  return detail?.priceI18n
}

const getOrderPriceLocale = (order) => {
  const snapshot = String(order?.settlementCurrency || '').trim()
  if (!snapshot) {
    return locale.value
  }
  return snapshot.replace(/_/g, '-')
}

const getLocalizedPaidAmount = (order) => {
  const details = Array.isArray(order?.detail) ? order.detail : []
  if (!details.length) {
    return (Number(order?.totalPrice || 0) / 100).toFixed(2)
  }

  const priceLocale = getOrderPriceLocale(order)

  const localizedOrigin = details.reduce((sum, detail) => {
    const priceFen = resolveLocalizedPriceFen(getDetailBasePriceFen(detail), getDetailPriceI18nSnapshot(detail), priceLocale)
    const quantity = Math.max(0, Number(detail?.quantity || 0))
    return sum + Math.max(0, priceFen) * quantity
  }, 0)

  let localizedTotal = localizedOrigin - Number(order?.discount || 0)
  if (order?.usePoints) {
    localizedTotal -= Number(order?.pointsUsed || 0)
  }
  return (Math.max(0, localizedTotal) / 100).toFixed(2)
}

const activeSataus = ref("")
const showRefundBtn = ref(true)
const showLogisticsBtn = ref(true)

// Tab定义，使用i18n
const tabColumns = computed(() => [
  { title: $t.value('viewAll'), id: '' },
  { title: $t.value('ordersPending'), id: '0' },
  { title: $t.value('ordersPendingConfirm'), id: '8' },
  { title: $t.value('ordersShipping'), id: '1' },
  { title: $t.value('ordersReceiving'), id: '2' },
  { title: $t.value('ordersToReview'), id: '3' },
  { title: $t.value('ordersRefunding'), id: '6' },
  { title: $t.value('ordersRefunded'), id: '5' },
  { title: $t.value('ordersCancelled'), id: '4' },
  { title: $t.value('ordersReviewed'), id: '7' },
])

const getStatusLabel = (status) => {
  const tab = tabColumns.value.find(t => t.id === status)
  return tab ? tab.title : ''
}

const getPayMethodLabel = (method) => {
  const map = {
    qrcode: $t.value('payByQrcode'),
    contact: $t.value('payByContact'),
    wechat: $t.value('payMethodWechat'),
    alipay: $t.value('payMethodAlipay'),
    bank_card_cn: $t.value('payMethodBankCn'),
    bank_card_us: $t.value('payMethodBankUs'),
    bank_card_mn: $t.value('payMethodBankMn'),
    paypal: $t.value('payMethodPaypal'),
  }
  return map[method] || method || '-'
}

const orderList = ref([])
const refundVisible = ref(false)
const refundOrderId = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const isLoading = ref(false)
const isBottom = ref(false)
const lastLoadedLocale = ref('')

const loadConfig = async () => {
  try {
    const [refundRes, logisticsRes] = await Promise.all([
      getSysConfigByKey('order_refund_enabled').catch(() => null),
      getSysConfigByKey('order_logistics_enabled').catch(() => null),
    ])
    if (refundRes?.code === 0 && refundRes.data !== undefined) {
      const val = typeof refundRes.data === 'object' ? refundRes.data.configValue : refundRes.data
      showRefundBtn.value = val !== 'false' && val !== '0'
    }
    if (logisticsRes?.code === 0 && logisticsRes.data !== undefined) {
      const val = typeof logisticsRes.data === 'object' ? logisticsRes.data.configValue : logisticsRes.data
      showLogisticsBtn.value = val !== 'false' && val !== '0'
    }
  } catch (e) { /* 配置获取失败时默认显示 */ }
}

const init = async (params) => {
  activeSataus.value = params || ''
  page.value = 1
  isBottom.value = false
  isLoading.value = true
  try {
    const res = await SelfOrderList({ status: activeSataus.value, page: page.value, pageSize: pageSize.value })
    if (res.code === 0) {
      orderList.value = res.data.list || []
      total.value = res.data.total || 0
      if (orderList.value.length >= total.value) isBottom.value = true
    }
  } finally {
    isLoading.value = false
  }
  startCountdown()
}

const loadMore = async () => {
  if (isBottom.value || isLoading.value) return
  isLoading.value = true
  page.value++
  try {
    const res = await SelfOrderList({ status: activeSataus.value, page: page.value, pageSize: pageSize.value })
    if (res.code === 0 && res.data.list && res.data.list.length > 0) {
      orderList.value = [...orderList.value, ...res.data.list]
      if (orderList.value.length >= (res.data.total || 0)) isBottom.value = true
    } else {
      isBottom.value = true
    }
  } finally {
    isLoading.value = false
  }
  startCountdown()
}

const countdownMap = ref({})
let countdownTimer = null

const updateCountdowns = () => {
  const map = {}
  let hasExpired = false
  orderList.value.forEach(item => {
    if (item.status === '0' && item.closeTime) {
      const remain = Math.max(0, Math.floor((new Date(item.closeTime).getTime() - Date.now()) / 1000))
      if (remain > 0) {
        const h = Math.floor(remain / 3600)
        const m = Math.floor((remain % 3600) / 60)
        const s = remain % 60
        map[item.ID] = h > 0
          ? `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
          : `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
      } else {
        hasExpired = true
      }
    }
  })
  countdownMap.value = map
  // 有订单倒计时到期，刷新列表
  if (hasExpired) {
    setTimeout(() => { init(activeSataus.value) }, 1000)
  }
}

const startCountdown = () => {
  if (countdownTimer) clearInterval(countdownTimer)
  updateCountdowns()
  countdownTimer = setInterval(updateCountdowns, 1000)
}

onUnmounted(() => { if (countdownTimer) clearInterval(countdownTimer) })

onLoad((options) => {
  loadConfig()
  init(options.status)
  lastLoadedLocale.value = locale.value
})

let isFirstShow = true
onShow(() => {
  if (isFirstShow) {
    isFirstShow = false
    lastLoadedLocale.value = locale.value
    return
  }

  const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
  if (localeChanged) {
    loadConfig()
  }

  init(activeSataus.value)
  lastLoadedLocale.value = locale.value
})

onReachBottom(() => {
  loadMore()
})

const formatOrderDate = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const cancelOrder = (item) => {
  uni.showModal({
    title: $t.value('cancelOrderHint'), content: $t.value('cancelOrderConfirm'), confirmColor: "#e50914", cancelText: $t.value('cancel'), confirmText: $t.value('confirm'),
    success: async (res) => {
      if (res.confirm) {
        const r = await updateOrderStatus({ ID: item.ID, status: "4" })
        if (r.code === 0) { uni.showToast({ title: $t.value('cancelSuccess'), icon: "none" }); init() }
      }
    }
  })
}

const payOff = (ID) => { uni.navigateTo({ url: `/pages/orderDetail/orderDetail?orderID=${ID}` }) }
const trackLogistics = (item) => { uni.navigateTo({ url: `/pages/logistics/logistics?express=${item.express}` }) }
const canApplyRefund = (item) => ['1', '2', '3', '7'].includes(item.status)
const openRefund = (item) => { refundOrderId.value = item.ID; refundVisible.value = true }
const onRefundSuccess = () => { init(activeSataus.value) }

const goOrderDetail = (item) => {
  uni.navigateTo({ url: `/pages/orderDetail/orderDetail?orderID=${item.ID}` })
}

const buyAgain = (item) => {
  if (item.detail && item.detail.length > 0) {
    const d = item.detail[0]
    uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${d.goodID}` })
  }
}

const goKefuWithOrder = (item) => {
  const orderNo = item.outTradeNo || item.OutTradeNo || item.ID
  const payMethod = item.payMethod || item.PayMethod || ''
  const payMethodLabel = getPayMethodLabel(payMethod)
  uni.navigateTo({
    url: `/pages/kefu/index?orderID=${encodeURIComponent(String(orderNo))}&payMethod=${encodeURIComponent(String(payMethod))}&payMethodLabel=${encodeURIComponent(String(payMethodLabel || payMethod || '-'))}`
  })
}

const confirm = async (item) => {
  const res = await updateOrderStatus({ ID: item.ID, status: '3' })
  if (res.code === 0) { uni.showToast({ title: $t.value('confirmReceiptSuccess'), icon: "none" }); init() }
}

const goComment = (order, detail) => {
  uni.navigateTo({ url: `/pages/evaluate/addEvaluate?orderID=${order.ID}&goodID=${detail.goodID}&SKUID=${detail.skuID}` })
}
const hasUncommentedItems = (order) => order.detail && order.detail.some(i => !i.isComment)
const hasCommentedItems = (order) => order.detail && order.detail.some(i => i.isComment)

const goCommentAll = (order) => {
  if (order.detail && order.detail.length > 1) {
    uni.navigateTo({ url: hasUncommentedItems(order)
      ? `/pages/evaluate/orderEvaluate?orderID=${order.ID}`
      : `/pages/evaluate/orderEvaluate?orderID=${order.ID}&mode=view`
    })
  } else if (order.detail && order.detail.length === 1) {
    const d = order.detail[0]
    uni.navigateTo({ url: `/pages/evaluate/addEvaluate?orderID=${order.ID}&goodID=${d.goodID}&SKUID=${d.skuID}` })
  }
}

const tapBtn = async (item) => {
  init(item.id)
}

const goBack = () => { uni.navigateBack() }
</script>

<style lang="scss">
page { background-color: #f4f7fb; }

.nf-order { min-height: 100vh; background: #f4f7fb; position: relative; }

.nf-order-bg {
  position: fixed; top: 0; left: 0; right: 0; height: 500rpx; z-index: 0; pointer-events: none;
  background: radial-gradient(ellipse at 30% -10%, rgba(37, 99, 235, 0.18) 0%, transparent 58%),
    radial-gradient(ellipse at 75% 10%, rgba(14, 165, 233, 0.14) 0%, transparent 52%);
}

.nf-navbar {
  background: rgba(244, 247, 251, 0.92); backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.2);
  padding: 0 28rpx 16rpx; position: sticky; top: 0; z-index: 99;
}
.nf-navbar-status { height: var(--status-bar-height, 0px); }
.nf-navbar-content { display: flex; align-items: center; justify-content: space-between; height: 88rpx; }
.nf-navbar-back {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: #ffffff; border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 8rpx 22rpx rgba(15, 23, 42, 0.08);
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-title { font-size: 34rpx; font-weight: 700; color: #0f172a; letter-spacing: 2rpx; }

.nf-tabs {
  position: sticky; top: 0; z-index: 98;
  background: rgba(244, 247, 251, 0.95); backdrop-filter: blur(16px);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.2);
}
.nf-tabs-scroll { white-space: nowrap; }
.nf-tabs-inner { display: inline-flex; padding: 0 16rpx; }
.nf-tab {
  display: inline-flex; padding: 20rpx 28rpx;
  font-size: 26rpx; color: rgba(15, 23, 42, 0.52); font-weight: 500;
  position: relative; white-space: nowrap; flex-shrink: 0;
  &.active {
    color: #0f172a; font-weight: 700;
    &::after {
      content: ''; position: absolute; bottom: 0;
      left: 28rpx; right: 28rpx; height: 4rpx;
      background: #2563eb; border-radius: 2rpx;
    }
  }
}

.nf-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 60vh;
}
.nf-empty-icon { font-size: 120rpx; margin-bottom: 20rpx; opacity: 0.5; }
.nf-empty-text { font-size: 28rpx; color: rgba(15, 23, 42, 0.46); }

.nf-order-list { padding: 20rpx 24rpx; position: relative; z-index: 1; }

.nf-order-card {
  background: rgba(255, 255, 255, 0.98);
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  border-radius: 20rpx; margin-bottom: 20rpx;
  overflow: hidden;
  box-shadow: 0 12rpx 24rpx rgba(15, 23, 42, 0.07);
}

.nf-order-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 24rpx 28rpx; border-bottom: 1rpx solid rgba(148, 163, 184, 0.15);
}
.nf-order-date {
  font-size: 24rpx; color: rgba(15, 23, 42, 0.55);
  display: flex; align-items: center; gap: 10rpx;
}
.nf-presale-tag {
  font-size: 20rpx; color: #fff;
  background: linear-gradient(135deg, #0ea5e9, #2563eb);
  padding: 2rpx 12rpx; border-radius: 6rpx;
}
.nf-order-status { font-size: 24rpx; font-weight: 600; }
.nf-st-0 { color: #ef4444; }
.nf-st-8 { color: #d97706; }
.nf-st-1 { color: #22c55e; }
.nf-st-2 { color: #3b82f6; }
.nf-st-3 { color: #f59e0b; }
.nf-st-4 { color: rgba(15, 23, 42, 0.35); }
.nf-st-5 { color: rgba(15, 23, 42, 0.4); }
.nf-st-6 { color: #f59e0b; }
.nf-st-7 { color: #22c55e; }
.nf-countdown { font-size: 22rpx; color: #ef4444; margin-left: 8rpx; }

.nf-goods-scroll { padding: 20rpx 28rpx; }
.nf-goods-row { display: inline-flex; gap: 12rpx; }
.nf-goods-thumb {
  width: 140rpx; height: 140rpx; border-radius: 12rpx;
  border: 1rpx solid rgba(148, 163, 184, 0.2); flex-shrink: 0;
}
.nf-goods-single { display: flex; padding: 20rpx 28rpx; gap: 20rpx; }
.nf-goods-thumb-lg {
  width: 160rpx; height: 160rpx; border-radius: 12rpx;
  border: 1rpx solid rgba(148, 163, 184, 0.2); flex-shrink: 0;
}
.nf-goods-single-info { flex: 1; display: flex; flex-direction: column; justify-content: center; }
.nf-goods-single-name {
  font-size: 28rpx; font-weight: 600; color: #0f172a; margin-bottom: 8rpx;
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}
.nf-goods-single-desc {
  font-size: 24rpx; color: rgba(15, 23, 42, 0.45);
  overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
}

.nf-order-summary {
  display: flex; justify-content: flex-end; align-items: center;
  padding: 16rpx 28rpx; border-top: 1rpx solid rgba(148, 163, 184, 0.15); gap: 8rpx;
}
.nf-order-count { font-size: 24rpx; color: rgba(15, 23, 42, 0.5); }
.nf-order-price { font-size: 30rpx; font-weight: 700; color: #ef4444; }

.nf-order-actions {
  display: flex; justify-content: flex-end; align-items: center;
  padding: 16rpx 28rpx 20rpx; gap: 12rpx; flex-wrap: wrap;
}
.nf-action {
  height: 60rpx; line-height: 58rpx; padding: 0 28rpx;
  border-radius: 30rpx; font-size: 24rpx; font-weight: 600;
  text-align: center; transition: all 0.2s;
  &:active { transform: scale(0.96); }
}
.nf-action-primary { background: linear-gradient(90deg, #2563eb, #0ea5e9); color: #fff; }
.nf-action-ghost {
  background: rgba(241, 245, 249, 0.95);
  border: 1rpx solid rgba(148, 163, 184, 0.24); color: rgba(15, 23, 42, 0.72);
}
.nf-action-warn {
  background: rgba(245, 158, 11, 0.15);
  border: 1rpx solid rgba(245, 158, 11, 0.3); color: #f59e0b;
}
.nf-action-tag {
  font-size: 22rpx; padding: 6rpx 20rpx; border-radius: 20rpx;
  &.cancelled { color: rgba(15, 23, 42, 0.35); background: rgba(241, 245, 249, 0.95); }
  &.refunding { color: #f59e0b; background: rgba(245, 158, 11, 0.1); }
  &.refunded { color: rgba(15, 23, 42, 0.42); background: rgba(241, 245, 249, 0.95); }
}

.nf-load-more {
  display: flex; justify-content: center; align-items: center;
  padding: 32rpx 0 16rpx; min-height: 80rpx;
}
.nf-load-text { font-size: 24rpx; color: rgba(15, 23, 42, 0.42); }
.nf-load-btn {
  font-size: 26rpx; color: rgba(15, 23, 42, 0.72);
  padding: 16rpx 48rpx; border-radius: 30rpx;
  background: rgba(255, 255, 255, 0.96);
  border: 1rpx solid rgba(148, 163, 184, 0.25);
}

</style>
