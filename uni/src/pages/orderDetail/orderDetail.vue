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
        <text class="nf-navbar-title">{{ $t('orderDetail') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body" v-if="data.ID">
      <!-- 订单状态 -->
      <view class="nf-card nf-status-card">
        <text class="nf-status-text" :class="'nf-st-' + data.status">{{ getStatusLabel(data.status) }}</text>
        <text class="nf-status-countdown" v-if="data.status === '0' && payCountdown">{{ $t('remainPayTime') }}: {{ payCountdown }}</text>
      </view>

      <!-- 收货地址 -->
      <view class="nf-card nf-address-card" @tap="toAddress" :style="data.status !== '0' ? 'opacity: 0.6;' : ''">
        <view class="nf-address-icon">📍</view>
        <view class="nf-address-info" v-if="hasAddress">
          <text class="nf-address-name">{{ data.name }} {{ data.phone }}</text>
          <text class="nf-address-detail">{{ data.province }}{{ data.city }}{{ data.area }} {{ data.street }}</text>
        </view>
        <view class="nf-address-info" v-else>
          <text class="nf-address-name">{{ $t('selectAddress') }}</text>
          <text class="nf-address-detail">{{ $t('addAddressHint') }}</text>
        </view>
        <uni-icons v-if="data.status === '0'" type="right" size="16" color="rgba(255,255,255,0.3)"></uni-icons>
      </view>

      <!-- 商品列表 -->
      <view class="nf-card nf-goods-card">
        <view class="nf-goods-item" v-for="(d, i) in data.detail" :key="i" @tap="goGoodsDetail(d)">
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

      <!-- 价格汇总 -->
      <view class="nf-card nf-price-card">
        <view class="nf-price-row">
          <text class="nf-price-label">{{ $t('productAmount') }}</text>
          <text class="nf-price-val">{{ cs }}{{ ((data.originPrice || 0) / 100).toFixed(2) }}</text>
        </view>
        <view class="nf-price-row nf-price-discount" v-if="data.discount > 0">
          <text class="nf-price-label">{{ $t('discountAmount') }}</text>
          <text class="nf-price-val">-{{ cs }}{{ (data.discount / 100).toFixed(2) }}</text>
        </view>
        <view class="nf-price-row" v-if="data.usePoints && data.pointsUsed > 0">
          <text class="nf-price-label">{{ $t('pointsDeduction') }}</text>
          <text class="nf-price-val nf-points-used">-{{ data.pointsUsed }} {{ $t('pointsUnit') }}</text>
        </view>
        <view class="nf-price-row nf-price-total">
          <text class="nf-price-label">{{ $t('actualPayment') }}</text>
          <text class="nf-price-val nf-price-total-val">{{ cs }}{{ ((data.totalPrice || 0) / 100).toFixed(2) }}</text>
        </view>
      </view>

      <!-- 订单信息 -->
      <view class="nf-card nf-info-card">
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
          <text class="nf-info-value">{{ getPayMethodLabel(data.payMethod) }}</text>
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
      </view>

      <!-- 退款操作 -->
      <view class="nf-card nf-refund-card" v-if="canApplyRefund || isRefunding || isRefunded">
        <view v-if="canApplyRefund" class="nf-refund-btn" @tap="openRefund">{{ $t('applyRefund') }}</view>
        <view v-else class="nf-refund-status">{{ isRefunding ? $t('refundProcessing') : $t('refunded') }}</view>
      </view>

      <view style="height: 140rpx;"></view>
    </view>

    <!-- 底部操作栏 -->
    <view class="nf-footer" v-if="data.ID">
      <!-- 待付款 -->
      <template v-if="data.status === '0'">
        <view class="nf-footer-info">
          <text class="nf-footer-label">{{ $t('actualPayment') }}</text>
          <text class="nf-footer-price">{{ cs }}{{ ((data.totalPrice || 0) / 100).toFixed(2) }}</text>
        </view>
        <view class="nf-footer-btn nf-footer-btn-ghost" @tap="cancelOrder">
          <text>{{ $t('cancelOrder') }}</text>
        </view>
        <view class="nf-footer-btn" @tap="payNow">
          <text>{{ $t('payNow') }}</text>
        </view>
      </template>
      <!-- 待收货 -->
      <template v-else-if="data.status === '2'">
        <view class="nf-footer-info"></view>
        <view class="nf-footer-btn nf-footer-btn-ghost" v-if="showLogisticsBtn && data.express" @tap="trackLogistics">
          <text>{{ $t('viewLogistics') }}</text>
        </view>
        <view class="nf-footer-btn" @tap="confirmReceipt">
          <text>{{ $t('confirmReceipt') }}</text>
        </view>
      </template>
      <!-- 待评价 -->
      <template v-else-if="data.status === '3'">
        <view class="nf-footer-info"></view>
        <view class="nf-footer-btn" @tap="goEvaluate">
          <text>{{ $t('reviewOrder') }}</text>
        </view>
      </template>
      <!-- 已取消 -->
      <template v-else-if="data.status === '4'">
        <view class="nf-footer-info"></view>
        <view class="nf-footer-btn" @tap="buyAgain">
          <text>{{ $t('buyAgain') }}</text>
        </view>
      </template>
      <!-- 其他状态不显示操作 -->
      <template v-else>
        <view class="nf-footer-info">
          <text class="nf-footer-label">{{ $t('actualPayment') }}</text>
          <text class="nf-footer-price">{{ cs }}{{ ((data.totalPrice || 0) / 100).toFixed(2) }}</text>
        </view>
      </template>
    </view>

    <refund-apply-popup
      v-model:visible="refundVisible"
      :order-id="orderID"
      @success="loadOrder"
    />
  </view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { selfOrder, updateOrderStatus, updateOrder } from '@/api/order.js'
import { checkNeedPay } from '@/api/base.js'
import { getSysConfigByKey } from '@/api/sysConfig.js'
import { getPaymentConfig } from '@/api/sysConfig.js'
import { getUrl, getExternalUrl } from "@/utils/url.js"
import { trackVisitorEvent } from '@/utils/visitorEvent.js'
import RefundApplyPopup from '@/components/refund-apply-popup/refund-apply-popup.vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const orderID = ref('')
const data = ref({})
const hasAddress = ref(false)
const pendingAddress = ref(null)
const refundVisible = ref(false)
const showRefundBtn = ref(true)
const showLogisticsBtn = ref(true)
const paymentMethods = ref([])

const goBack = () => { uni.navigateBack() }

const getStatusLabel = (status) => {
  const map = {
    '0': $t.value('ordersPending'),
    '1': $t.value('ordersShipping'),
    '2': $t.value('ordersReceiving'),
    '3': $t.value('ordersToReview'),
    '4': $t.value('ordersCancelled'),
    '5': $t.value('ordersRefunded'),
    '6': $t.value('ordersRefunding'),
    '7': $t.value('ordersReviewed'),
  }
  return map[status] || ''
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

const buildDefaultPaymentMethods = () => ([
  { key: 'qrcode', label: $t.value('payByQrcode') },
  { key: 'contact', label: $t.value('payByContact') },
])

const getManualFallbackTip = (payMethod, payMethodLabel) => {
  const tipKeyMap = {
    wechat: 'paymentManualTipWechat',
    alipay: 'paymentManualTipAlipay',
    bank_card_cn: 'paymentManualTipBankCn',
    bank_card_us: 'paymentManualTipBankUs',
    bank_card_mn: 'paymentManualTipBankMn',
    paypal: 'paymentManualTipPaypal',
  }
  const tipKey = tipKeyMap[payMethod] || 'paymentManualTipDefault'
  const fallbackLabel = payMethodLabel || getPayMethodLabel(payMethod)
  return $t.value(tipKey).replace('{}', fallbackLabel)
}

const buildKefuUrl = (orderNo, payMethod, payMethodLabel) => {
  const orderID = encodeURIComponent(String(orderNo || ''))
  const method = encodeURIComponent(String(payMethod || ''))
  const methodLabel = encodeURIComponent(String(payMethodLabel || ''))
  return `/pages/kefu/index?orderID=${orderID}&payMethod=${method}&payMethodLabel=${methodLabel}`
}

const trackKefuGuideEvent = (action, extra = {}) => {
  const payload = { source: 'order_detail', ...extra }
  return trackVisitorEvent({
    action,
    label: String(payload.payMethod || ''),
    extra: payload
  })
}

const routeToKefuByMethod = (orderNo, payMethod, payMethodLabel) => {
  const url = buildKefuUrl(orderNo, payMethod, payMethodLabel)
  uni.setClipboardData({
    data: String(orderNo),
    success: () => {
      trackKefuGuideEvent('copy_order_no', { orderNo, payMethod, payMethodLabel })
      uni.showToast({ title: `${$t.value('orderNoCopied')}: ${orderNo}`, icon: 'none', duration: 2000 })
      setTimeout(() => {
        trackKefuGuideEvent('navigate_kefu', {
          orderNo,
          payMethod,
          payMethodLabel,
          navigateMode: 'navigateTo'
        })
        uni.navigateTo({ url })
      }, 1500)
    },
    fail: () => {
      trackKefuGuideEvent('copy_order_no_fail', { orderNo, payMethod, payMethodLabel })
      trackKefuGuideEvent('navigate_kefu', {
        orderNo,
        payMethod,
        payMethodLabel,
        navigateMode: 'navigateTo'
      })
      uni.navigateTo({ url })
    },
  })
}

const confirmManualFallback = (payMethod, payMethodLabel) => {
  trackKefuGuideEvent('manual_fallback_modal_show', { payMethod, payMethodLabel, orderNo: data.value.ID })
  return new Promise((resolve) => {
    uni.showModal({
      title: payMethodLabel || $t.value('paymentManualFallbackTitle'),
      content: `${getManualFallbackTip(payMethod, payMethodLabel)}\n\n${$t.value('paymentManualProofHint')}`,
      confirmText: $t.value('paymentManualFallbackContact'),
      cancelText: $t.value('paymentManualFallbackLater'),
      success: (res) => {
        trackKefuGuideEvent(res.confirm ? 'manual_fallback_confirm' : 'manual_fallback_cancel', {
          payMethod,
          payMethodLabel,
          orderNo: data.value.ID
        })
        resolve(!!res.confirm)
      },
      fail: () => {
        trackKefuGuideEvent('manual_fallback_cancel', { payMethod, payMethodLabel, orderNo: data.value.ID, fail: true })
        resolve(false)
      },
    })
  })
}

const loadPaymentMethods = async () => {
  try {
    const res = await getPaymentConfig()
    const methods = Array.isArray(res?.data?.methods)
      ? res.data.methods
        .filter(m => m && m.key && m.enabled !== false)
        .map(m => ({
          key: m.key,
          label: m.label || getPayMethodLabel(m.key)
        }))
      : []
    paymentMethods.value = methods.length > 0 ? methods : buildDefaultPaymentMethods()
  } catch (e) {
    paymentMethods.value = buildDefaultPaymentMethods()
  }
}

const formatSpecs = (specs, attrs) => {
  const arr = [...(Array.isArray(specs) ? specs : []), ...(Array.isArray(attrs) ? attrs : [])]
  if (!arr.length) return ''
  return arr.map(s => {
    const label = $lt.value(s.labelI18n || s.nameI18n) || $lt.value(s.label || s.name) || s.label || s.name || ''
    const value = $lt.value(s.valueI18n) || $lt.value(s.value) || s.value || ''
    return label ? `${label}: ${value}` : value
  }).join('  ')
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

/* =================== 初始化 =================== */
onLoad(async (options) => {
  if (options.orderID) {
    orderID.value = options.orderID
    await loadOrder()
    loadConfig()
    loadPaymentMethods()
  }
})

let isFirstShow = true
onShow(() => {
  if (isFirstShow) { isFirstShow = false; return }
  // 检查是否有从地址页选回的地址
  const addr = uni.getStorageSync('selectedAddress')
  if (addr) {
    uni.removeStorageSync('selectedAddress')
    pendingAddress.value = addr
    data.value.name = addr.name
    data.value.phone = addr.phone
    data.value.province = addr.provinceStr || addr.province
    data.value.city = addr.cityStr || addr.city
    data.value.area = addr.areaStr || addr.area
    data.value.street = addr.street
    hasAddress.value = true
    return
  }
  if (orderID.value) loadOrder()
})

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
  } catch (e) { /* ignore */ }
}

const loadOrder = async () => {
  const res = await selfOrder(orderID.value)
  if (res.code === 0) {
    data.value = res.data
    hasAddress.value = !!(res.data.city || res.data.name)
    startPayCountdown()
  }
}

/* =================== 倒计时 =================== */
const payCountdown = ref('')
let payTimer = null
let payExpiredHandled = false

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
    payCountdown.value = ''
    if (payTimer) clearInterval(payTimer)
    // 只在首次到期时刷新一次，避免循环轮询
    if (!payExpiredHandled) {
      payExpiredHandled = true
      uni.showToast({ title: $t.value('payTimeout'), icon: 'none' })
      // 先尝试从后端刷新，如果后端尚未自动取消，本地也强制展示已取消
      setTimeout(async () => {
        await loadOrder()
        if (data.value.status === '0') {
          data.value.status = '4'
        }
      }, 1500)
    }
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

/* =================== 权限计算 =================== */
const canApplyRefund = computed(() => showRefundBtn.value && ['1', '2', '3', '7'].includes(String(data.value.status || '')))
const isRefunding = computed(() => String(data.value.status || '') === '6')
const isRefunded = computed(() => String(data.value.status || '') === '5')
const openRefund = () => { refundVisible.value = true }

/* =================== 地址 =================== */
const toAddress = () => {
  if (data.value.status !== '0') return
  if (data.value.ID) {
    uni.navigateTo({ url: `/pages/address/address?ID=${data.value.ID}&from=orderDetail` })
  }
}

/* =================== 操作 =================== */
const cancelOrder = () => {
  uni.showModal({
    title: $t.value('cancelOrderHint'),
    content: $t.value('cancelOrderConfirm'),
    confirmColor: '#e50914',
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: async (res) => {
      if (res.confirm) {
        const r = await updateOrderStatus({ ID: data.value.ID, status: '4' })
        if (r.code === 0) {
          uni.showToast({ title: $t.value('cancelSuccess'), icon: 'none' })
          await loadOrder()
        }
      }
    }
  })
}

const payNow = async () => {
  if (!hasAddress.value) {
    uni.showToast({ title: $t.value('addressMissing'), icon: 'none' })
    setTimeout(() => { toAddress() }, 500)
    return
  }
  // 如果有待提交的地址变更，先提交
  if (pendingAddress.value) {
    const addr = pendingAddress.value
    await updateOrder({
      ID: Number(orderID.value),
      name: addr.name,
      phone: addr.phone,
      province: addr.provinceStr || addr.province,
      city: addr.cityStr || addr.city,
      area: addr.areaStr || addr.area,
      Street: addr.street,
      active: addr.active,
    })
    pendingAddress.value = null
  }
  const needPayRes = await checkNeedPay({ orderID: Number(orderID.value), openid: uni.getStorageSync('openid') })
  if (!needPayRes.data) {
    uni.showToast({ title: $t.value('orderSubmitSuccess'), icon: 'success', duration: 2000 })
    setTimeout(() => { loadOrder() }, 2000)
    return
  }
  if (!paymentMethods.value.length) {
    paymentMethods.value = buildDefaultPaymentMethods()
  }
  uni.showActionSheet({
    itemList: paymentMethods.value.map(item => item.label),
    success: async (res) => {
      const payMethod = paymentMethods.value[res.tapIndex]?.key || 'contact'
      const payMethodLabel = paymentMethods.value[res.tapIndex]?.label || getPayMethodLabel(payMethod)
      // 同步支付方式到后端
      await updateOrder({ ID: Number(orderID.value), payMethod })
      if (payMethod === 'qrcode') {
        uni.navigateTo({ url: `/pages/pay/index?amount=${((data.value.totalPrice || 0) / 100).toFixed(2)}&orderNo=${data.value.ID}&orderId=${orderID.value}` })
      } else {
        trackKefuGuideEvent('guide_entry', {
          orderNo: data.value.ID,
          payMethod,
          payMethodLabel
        })

        if (payMethod !== 'contact') {
          const shouldContactNow = await confirmManualFallback(payMethod, payMethodLabel)
          if (!shouldContactNow) {
            uni.showToast({ title: $t.value('paymentManualSavedMethod'), icon: 'none' })
            return
          }
        }

        routeToKefuByMethod(data.value.ID, payMethod, payMethodLabel)
      }
    }
  })
}

const confirmReceipt = async () => {
  const res = await updateOrderStatus({ ID: data.value.ID, status: '3' })
  if (res.code === 0) {
    uni.showToast({ title: $t.value('confirmReceiptSuccess'), icon: 'none' })
    await loadOrder()
  }
}

const trackLogistics = () => {
  uni.navigateTo({ url: `/pages/logistics/logistics?express=${data.value.express}` })
}

const goEvaluate = () => {
  if (data.value.detail && data.value.detail.length > 1) {
    uni.navigateTo({ url: `/pages/evaluate/orderEvaluate?orderID=${orderID.value}` })
  } else if (data.value.detail && data.value.detail.length === 1) {
    const d = data.value.detail[0]
    uni.navigateTo({ url: `/pages/evaluate/addEvaluate?orderID=${orderID.value}&goodID=${d.goodID}&SKUID=${d.skuID}` })
  }
}

const buyAgain = () => {
  if (data.value.detail && data.value.detail.length > 0) {
    uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${data.value.detail[0].goodID}` })
  }
}

const goGoodsDetail = (d) => {
  uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${d.goodID}` })
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

/* 状态卡片 */
.nf-status-card {
  padding: 28rpx; display: flex; flex-direction: column; align-items: center; gap: 12rpx;
}
.nf-status-text { font-size: 36rpx; font-weight: 700; }
.nf-status-countdown { font-size: 26rpx; color: #e50914; }
.nf-st-0 { color: #e50914; }
.nf-st-1 { color: #22c55e; }
.nf-st-2 { color: #3b82f6; }
.nf-st-3 { color: #f59e0b; }
.nf-st-4 { color: rgba(255,255,255,0.3); }
.nf-st-5 { color: rgba(255,255,255,0.4); }
.nf-st-6 { color: #f59e0b; }
.nf-st-7 { color: #22c55e; }

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
.nf-price-total { border-top: 1rpx solid rgba(255,255,255,0.08); margin-top: 8rpx; padding-top: 20rpx; }
.nf-price-total-val { font-size: 32rpx; font-weight: 700; color: #e50914; }
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

/* 底部操作栏 */
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
  min-width: 180rpx; height: 100%; background: #e50914;
  display: flex; justify-content: center; align-items: center;
  font-size: 28rpx; font-weight: 700; color: #fff; padding: 0 24rpx;
  &:active { background: #b30710; }
}
.nf-footer-btn-ghost {
  background: transparent; border-left: 1rpx solid rgba(255,255,255,0.06);
  color: rgba(255,255,255,0.6);
  &:active { background: rgba(255,255,255,0.04); }
}
</style>
