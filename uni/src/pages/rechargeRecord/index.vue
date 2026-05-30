<template>
  <view class="nf-recharge-record">
    <view class="nf-recharge-record-bg"></view>

    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="nf-navbar-title">{{ $t('tryonRechargeRecord') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body">
      <view class="nf-hero">
        <text class="nf-hero-label">{{ $t('tryonRechargeRecord') }}</text>
        <text class="nf-hero-value">{{ total }}</text>
        <text class="nf-hero-tip">{{ $t('tryonRechargeRecordTip') }}</text>
      </view>

      <view class="nf-empty" v-if="recordList.length === 0 && !loading">
        <view class="nf-empty-icon">🧾</view>
        <text class="nf-empty-text">{{ $t('noTryonRechargeRecord') }}</text>
      </view>

      <view class="nf-list" v-else>
        <view class="nf-card" v-for="item in recordList" :key="item._recordKey">
          <view class="nf-card-head">
            <text class="nf-card-status" :class="`nf-st-${String(item.status || '')}`">{{ item._statusLabel }}</text>
            <text class="nf-card-time">{{ item._displayTime }}</text>
          </view>

          <view class="nf-row">
            <text class="nf-row-label">{{ $t('orderNo') }}</text>
            <text class="nf-row-value">{{ item.outTradeNo || item.OutTradeNo || item.ID }}</text>
          </view>
          <view class="nf-row">
            <text class="nf-row-label">{{ $t('tryonCoins') }}</text>
            <text class="nf-row-value">{{ item.points || item.Points || 0 }}</text>
          </view>
          <view class="nf-row">
            <text class="nf-row-label">{{ $t('payAmount') }}</text>
            <text class="nf-row-value">{{ item._amountText }}</text>
          </view>
          <view class="nf-row">
            <text class="nf-row-label">{{ $t('paymentMethod') }}</text>
            <text class="nf-row-value">{{ item._payMethodLabel }}</text>
          </view>

          <view class="nf-actions">
            <view
              v-if="String(item.status || '') === '0'"
              class="nf-btn nf-btn-ghost"
              @tap="cancelRecord(item)"
            >{{ $t('cancelOrder') }}</view>
            <view
              v-if="String(item.status || '') === '0'"
              class="nf-btn nf-btn-primary"
              @tap="goPay(item)"
            >{{ $t('payNow') }}</view>
            <view
              v-if="String(item.status || '') === '8'"
              class="nf-btn nf-btn-primary"
              @tap="goKefu(item)"
            >{{ $t('contactCustomerService') }}</view>
          </view>
        </view>
      </view>

      <view class="nf-load-more" v-if="recordList.length > 0">
        <text class="nf-load-more-text" v-if="loading">{{ $t('loading') }}</text>
        <text class="nf-load-more-text" v-else-if="noMore">{{ $t('reachedBottom') }}</text>
        <text class="nf-load-more-text" v-else-if="reachedBottom" @tap="loadMore">{{ $t('loadMore') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onReachBottom, onShow } from '@dcloudio/uni-app'
import {
  getMyTryonRechargeOrderList,
  cancelTryonRechargeOrder,
} from '@/api/tryonRechargeOrder.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { resolveLocalizedPriceFen } from '@/utils/price-i18n.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const $t = computed(() => langStore.$t)
const cs = computed(() => appConfigStore.currencySymbol || '¥')
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const page = ref(1)
const pageSize = 10
const total = ref(0)
const recordList = ref([])
const loading = ref(false)
const noMore = ref(false)
const reachedBottom = ref(false)
const lastLoadedLocale = ref('')

const statusLabel = (status) => {
  const map = {
    '0': $t.value('ordersPending'),
    '8': $t.value('ordersPendingConfirm'),
    '1': $t.value('ordersCompleted'),
    '4': $t.value('ordersCancelled'),
  }
  return map[String(status || '')] || String(status || '-')
}

const payMethodLabel = (method) => {
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

const parsePlanSnapshot = (item) => {
  const raw = item?.planSnapshot || item?.PlanSnapshot
  if (!raw) return null
  if (typeof raw === 'object') {
    return raw
  }
  const text = String(raw).trim()
  if (!text) return null
  try {
    const parsed = JSON.parse(text)
    return parsed && typeof parsed === 'object' ? parsed : null
  } catch (e) {
    return null
  }
}

const normalizeSettlementLocale = (value) => {
  const normalized = String(value || '').trim()
  return normalized ? normalized.replace(/_/g, '-') : ''
}

const getOrderSettlementLocale = (item) => {
  const direct = normalizeSettlementLocale(item?.settlementCurrency || item?.SettlementCurrency)
  if (direct) {
    return direct
  }
  const snapshot = parsePlanSnapshot(item)
  const fromSnapshot = normalizeSettlementLocale(snapshot?.settlementCurrency)
  return fromSnapshot || locale.value
}

const getOrderCurrencySymbol = (item) => {
  const direct = String(item?.settlementCurrencySymbol || item?.SettlementCurrencySymbol || '').trim()
  if (direct) {
    return direct
  }
  const snapshot = parsePlanSnapshot(item)
  const fromSnapshot = String(snapshot?.settlementCurrencySymbol || '').trim()
  if (fromSnapshot) {
    return fromSnapshot
  }
  return cs.value
}

const getOrderAmountFen = (item) => {
  const amountFen = Number(item?.amount ?? item?.Amount)
  if (Number.isFinite(amountFen) && amountFen >= 0) {
    return Math.round(amountFen)
  }

  const snapshot = parsePlanSnapshot(item)
  if (snapshot) {
    const resolvedFen = resolveLocalizedPriceFen(snapshot?.price, snapshot?.priceI18n, getOrderSettlementLocale(item))
    const safeResolved = Number(resolvedFen)
    if (Number.isFinite(safeResolved) && safeResolved >= 0) {
      return Math.round(safeResolved)
    }

    const selectedPriceFen = Number(snapshot?.selectedPriceFen)
    if (Number.isFinite(selectedPriceFen) && selectedPriceFen >= 0) {
      return Math.round(selectedPriceFen)
    }
  }

  return 0
}

const amountYuan = (amountFen) => {
  const cents = Number(amountFen || 0)
  return (cents / 100).toFixed(2)
}

const resolveOrderAmountYuan = (item) => {
  return amountYuan(getOrderAmountFen(item))
}

const formatOrderAmount = (item) => {
  return `${getOrderCurrencySymbol(item)}${resolveOrderAmountYuan(item)}`
}

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

// 列表展示字段在加载时统一生成；支付、取消、客服跳转仍读取原始订单字段。
const normalizeRechargeRecord = (item, index) => ({
  ...item,
  _recordKey: item?.ID || item?.id || item?.outTradeNo || item?.OutTradeNo || `recharge-${index}`,
  _statusLabel: statusLabel(item?.status),
  _displayTime: formatTime(item?.CreatedAt || item?.createdAt),
  _amountText: formatOrderAmount(item),
  _payMethodLabel: payMethodLabel(item?.payMethod || item?.PayMethod),
})

// 首屏和分页追加共用同一入口，避免两条数据路径展示字段不一致。
const normalizeRechargeRecordList = (list, offset = 0) => {
  return (Array.isArray(list) ? list : []).map((item, index) => normalizeRechargeRecord(item, offset + index))
}

const fetchList = async (isLoadMore = false) => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getMyTryonRechargeOrderList({ page: page.value, pageSize })
    if (res.code === 0 && res.data) {
      const list = res.data.list || []
      const normalizedList = normalizeRechargeRecordList(list, isLoadMore ? recordList.value.length : 0)
      total.value = Number(res.data.total || 0)
      if (isLoadMore) {
        recordList.value = [...recordList.value, ...normalizedList]
      } else {
        recordList.value = normalizedList
      }
      noMore.value = recordList.value.length >= total.value || list.length < pageSize
    }
  } finally {
    loading.value = false
  }
}

const resetAndLoad = async () => {
  page.value = 1
  noMore.value = false
  reachedBottom.value = false
  recordList.value = []
  await fetchList(false)
}

const loadMore = () => {
  if (loading.value || noMore.value) return
  reachedBottom.value = false
  page.value += 1
  fetchList(true)
}

const goPay = (item) => {
  const orderId = item.ID || item.id
  const orderNo = item.outTradeNo || item.OutTradeNo || orderId
  const payMethod = item.payMethod || item.PayMethod || 'contact'
  const amount = resolveOrderAmountYuan(item)
  const points = Number(item.points || item.Points || 0)
  const closeTime = String(item.closeTime || item.CloseTime || '').trim()
  const closeTimePart = closeTime ? `&closeTime=${encodeURIComponent(closeTime)}` : ''
  const selectedPayMethodLabel = encodeURIComponent(payMethodLabel(payMethod))

  uni.navigateTo({
    url: `/pages/pay/index?orderType=recharge&amount=${encodeURIComponent(amount)}&orderNo=${encodeURIComponent(String(orderNo))}&orderId=${orderId}&payMethod=${encodeURIComponent(payMethod)}&payMethodLabel=${selectedPayMethodLabel}&rechargePoints=${points}${closeTimePart}`
  })
}

const goKefu = (item) => {
  const orderNo = item.outTradeNo || item.OutTradeNo || item.ID || item.id
  const payMethod = item.payMethod || item.PayMethod || ''
  const payMethodText = payMethodLabel(payMethod)
  uni.navigateTo({
    url: `/pages/kefu/index?orderID=${encodeURIComponent(String(orderNo))}&payMethod=${encodeURIComponent(String(payMethod))}&payMethodLabel=${encodeURIComponent(String(payMethodText))}`
  })
}

const cancelRecord = (item) => {
  const id = item.ID || item.id
  uni.showModal({
    title: $t.value('cancelOrderHint'),
    content: $t.value('cancelOrderConfirm'),
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: async (res) => {
      if (!res.confirm) return
      const r = await cancelTryonRechargeOrder(id)
      if (r.code === 0) {
        uni.showToast({ title: $t.value('cancelSuccess'), icon: 'none' })
        resetAndLoad()
      }
    }
  })
}

const goBack = () => {
  uni.navigateBack()
}

onLoad(() => {
  resetAndLoad()
  lastLoadedLocale.value = locale.value
})

onShow(() => {
  const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
  if (localeChanged) {
    resetAndLoad()
    lastLoadedLocale.value = locale.value
    return
  }
  if (recordList.value.length === 0) {
    lastLoadedLocale.value = locale.value
    return
  }
  resetAndLoad()
  lastLoadedLocale.value = locale.value
})

onReachBottom(() => {
  if (loading.value || noMore.value) return
  reachedBottom.value = true
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.nf-recharge-record {
  min-height: 100vh;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 62%), #f4f7fb;
  position: relative;
}

.nf-recharge-record-bg {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 620rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 10% 12%, rgba(59, 130, 246, 0.12), transparent 50%),
    radial-gradient(circle at 92% 20%, rgba(14, 165, 233, 0.1), transparent 46%);
}

.nf-navbar {
  position: sticky;
  top: 0;
  z-index: 99;
  padding: 0 24rpx 12rpx;
  backdrop-filter: blur(16rpx);
  background: rgba(244, 247, 251, 0.72);
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.06);
}

.nf-navbar-status {
  height: var(--status-bar-height, 0px);
}

.nf-navbar-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  border: 1rpx solid rgba(15, 23, 42, 0.1);
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-navbar-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #0f172a;
}

.nf-body {
  position: relative;
  z-index: 1;
  padding: 24rpx;
}

.nf-hero {
  text-align: center;
  padding: 46rpx 30rpx;
  margin-bottom: 24rpx;
  border-radius: 20rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.14);
  background: linear-gradient(135deg, rgba(219, 234, 254, 0.92), rgba(239, 246, 255, 0.98));
  box-shadow: 0 16rpx 32rpx rgba(59, 130, 246, 0.12);
}

.nf-hero-label {
  display: block;
  margin-bottom: 10rpx;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.64);
}

.nf-hero-value {
  display: block;
  margin-bottom: 8rpx;
  font-size: 76rpx;
  line-height: 1;
  font-weight: 800;
  color: #0f172a;
}

.nf-hero-tip {
  display: block;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.55);
}

.nf-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 90rpx 0;
}

.nf-empty-icon {
  font-size: 90rpx;
  margin-bottom: 14rpx;
}

.nf-empty-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.5);
}

.nf-card {
  margin-bottom: 14rpx;
  padding: 20rpx 22rpx;
  border-radius: 16rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.06);
}

.nf-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14rpx;
}

.nf-card-status {
  font-size: 24rpx;
  font-weight: 700;
}

.nf-st-0 {
  color: #dc2626;
}

.nf-st-8 {
  color: #d97706;
}

.nf-st-1 {
  color: #16a34a;
}

.nf-st-4 {
  color: rgba(15, 23, 42, 0.42);
}

.nf-card-time {
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.5);
}

.nf-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.nf-row:last-child {
  margin-bottom: 0;
}

.nf-row-label {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.54);
}

.nf-row-value {
  font-size: 24rpx;
  color: #0f172a;
}

.nf-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12rpx;
  margin-top: 14rpx;
}

.nf-btn {
  min-width: 132rpx;
  height: 56rpx;
  border-radius: 28rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  font-weight: 600;
}

.nf-btn-ghost {
  background: rgba(241, 245, 249, 0.95);
  border: 1rpx solid rgba(148, 163, 184, 0.24);
  color: rgba(15, 23, 42, 0.72);
}

.nf-btn-primary {
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
}

.nf-load-more {
  text-align: center;
  padding: 22rpx 0 60rpx;
}

.nf-load-more-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.5);
}
</style>
