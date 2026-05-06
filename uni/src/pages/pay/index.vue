<template>
  <view class="nf-pay">
    <view class="nf-pay-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a" />
        </view>
        <text class="nf-navbar-title">{{ $t('payTitle') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-pay-content">
      <!-- 订单金额 -->
      <view class="nf-pay-amount-card">
        <text class="nf-pay-amount-label">{{ $t('payAmount') }}</text>
        <text class="nf-pay-amount-value">{{ cs }}{{ amount }}</text>
        <text class="nf-pay-order-no" v-if="orderNo">{{ $t('orderNo') }}: {{ orderNo }}</text>
        <!-- 倒计时 -->
        <view class="nf-pay-countdown" v-if="payCountdown">
          <text class="nf-pay-countdown-label">{{ $t('remainPayTime') }}</text>
          <text class="nf-pay-countdown-time">{{ payCountdown }}</text>
        </view>
      </view>

      <!-- 支付方式 -->
      <view class="nf-pay-method-card" v-if="paymentMethods.length">
        <view class="nf-pay-method-tabs">
          <view
            v-for="item in paymentMethods"
            :key="item.key"
            class="nf-pay-method-tab"
            :class="{ active: selectedPayMethod === item.key }"
            @tap="selectPayMethod(item.key)"
          >{{ item.label }}</view>
        </view>
      </view>

      <!-- 支付二维码 -->
      <view class="nf-pay-qr-card" v-if="isQrMethod">
        <text class="nf-pay-qr-title">{{ $t('payQrTitle') }}</text>
        <!-- 多码切换 -->
        <view class="nf-pay-qr-tabs" v-if="qrList.length > 1">
          <view
            v-for="(qr, i) in qrList"
            :key="i"
            class="nf-pay-qr-tab"
            :class="{ active: currentQrIndex === i }"
            @tap="currentQrIndex = i"
          >{{ localText(qr.nameI18n) || qr.name }}</view>
        </view>
        <view class="nf-pay-qr-wrap">
          <image
            v-if="currentQrUrl"
            class="nf-pay-qr-img"
            :src="currentQrUrl"
            mode="aspectFit"
            @tap="previewQr"
          />
          <view v-else class="nf-pay-qr-placeholder">
            <uni-icons type="scan" size="48" color="rgba(148,163,184,0.7)" />
          </view>
        </view>
        <text class="nf-pay-qr-tip">{{ $t('payQrTip') }}</text>
      </view>

      <!-- 人工支付提示 -->
      <view class="nf-pay-manual-card" v-else>
        <text class="nf-pay-manual-title">{{ selectedPayMethodLabel }}</text>
        <text class="nf-pay-manual-tip">{{ manualFallbackTip }}</text>
        <text class="nf-pay-manual-proof">{{ $t('paymentManualProofHint') }}</text>
      </view>

      <!-- 付款提示文本 -->
      <view class="nf-pay-tip-card" v-if="paymentTipText">
        <text class="nf-pay-tip-text" :style="{ fontSize: paymentTipSize + 'px', color: paymentTipColor }">{{ paymentTipText }}</text>
      </view>

      <!-- 支付步骤 -->
      <view class="nf-pay-steps">
        <text class="nf-pay-steps-title">{{ $t('payStepsTitle') }}</text>
        <view class="nf-pay-step" v-for="(step, i) in steps" :key="i">
          <view class="nf-pay-step-num">
            <text>{{ i + 1 }}</text>
          </view>
          <text class="nf-pay-step-text">{{ step }}</text>
        </view>
      </view>

      <!-- 操作按钮 -->
      <view class="nf-pay-actions">
        <view class="nf-pay-btn nf-pay-btn-primary" v-if="isQrMethod" @tap="saveQr">
          <uni-icons type="download" size="18" color="#fff" />
          <text>{{ $t('saveQrCode') }}</text>
        </view>
        <view class="nf-pay-btn nf-pay-btn-kefu" @tap="goKefu">
          <uni-icons type="chatbubble" size="18" color="#2563eb" />
          <text>{{ $t('contactCustomerService') }}</text>
        </view>
      </view>

      <!-- 底部操作 -->
      <view class="nf-pay-bottom-actions">
        <view class="nf-pay-bottom-btn nf-pay-bottom-home" @tap="goHome">
          <text>{{ $t('goHome') }}</text>
        </view>
        <view class="nf-pay-bottom-btn nf-pay-bottom-paid" @tap="confirmPaid">
          <text>{{ $t('alreadyPaid') }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { onLoad, onBackPress } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { request } from '@/utils/request.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { getEnabledQrcodePayments } from '@/api/qrcodePayment.js'
import { getPaymentConfig } from '@/api/sysConfig.js'
import { localText } from '@/utils/i18n'
import { selfOrder } from '@/api/order.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)

const amount = ref('0.00')
const orderNo = ref('')
const orderId = ref('')
const paymentMethods = ref([])
const selectedPayMethod = ref('qrcode')

const paymentMethodLabelMap = computed(() => ({
  qrcode: $t.value('payByQrcode'),
  contact: $t.value('payByContact'),
  wechat: $t.value('payMethodWechat'),
  alipay: $t.value('payMethodAlipay'),
  bank_card_cn: $t.value('payMethodBankCn'),
  bank_card_us: $t.value('payMethodBankUs'),
  bank_card_mn: $t.value('payMethodBankMn'),
  paypal: $t.value('payMethodPaypal'),
}))

const getPaymentMethodLabel = (payMethod, label) => {
  return paymentMethodLabelMap.value[payMethod] || label || payMethod || $t.value('contactCustomerService')
}

const buildDefaultPaymentMethods = () => ([
  { key: 'qrcode', label: getPaymentMethodLabel('qrcode') },
  { key: 'contact', label: getPaymentMethodLabel('contact') }
])

const normalizePaymentMethods = (methods) => {
  if (!Array.isArray(methods)) {
    return []
  }
  return methods
    .filter(item => item && typeof item.key === 'string' && item.key)
    .map(item => ({
      key: item.key,
      label: getPaymentMethodLabel(item.key, item.label)
    }))
}

const isQrMethod = computed(() => selectedPayMethod.value === 'qrcode')
const selectedPayMethodLabel = computed(() => {
  const found = paymentMethods.value.find(item => item.key === selectedPayMethod.value)
  return found?.label || getPaymentMethodLabel(selectedPayMethod.value)
})

// 多码支持
const qrList = ref([])
const currentQrIndex = ref(0)
const currentQrUrl = computed(() => {
  const qr = qrList.value[currentQrIndex.value]
  if (!qr) return ''
  return qr.externalPath ? getExternalUrl(qr.externalPath) : getUrl(qr.image)
})

// 付款提示文本
const paymentTipText = ref('')
const paymentTipSize = ref(14)
const paymentTipColor = ref('#334155')

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
  const fallbackLabel = payMethodLabel || getPaymentMethodLabel(payMethod)
  return $t.value(tipKey).replace('{}', fallbackLabel)
}

const manualFallbackTip = computed(() => {
  return getManualFallbackTip(selectedPayMethod.value, selectedPayMethodLabel.value)
})

const steps = computed(() => {
  if (isQrMethod.value) {
    return [$t.value('payStep1'), $t.value('payStep2'), $t.value('payStep3')]
  }
  return [manualFallbackTip.value, $t.value('paymentManualProofHint'), $t.value('contactCustomerService')]
})

const selectPayMethod = (payMethod) => {
  selectedPayMethod.value = payMethod
  if (payMethod === 'qrcode' && qrList.value.length === 0) {
    loadQrCodes()
  }
}

const loadPaymentMethods = async (preferredMethod, preferredLabel) => {
  let methods = []
  try {
    const res = await getPaymentConfig()
    if (res.code === 0 && res.data) {
      methods = normalizePaymentMethods(res.data.methods)
    }
  } catch (e) {
    methods = []
  }

  if (methods.length === 0) {
    methods = buildDefaultPaymentMethods()
  }

  if (preferredMethod && typeof preferredMethod === 'string') {
    const existing = methods.find(item => item.key === preferredMethod)
    if (existing) {
      if (preferredLabel) {
        existing.label = preferredLabel
      }
    } else {
      methods.unshift({
        key: preferredMethod,
        label: preferredLabel || getPaymentMethodLabel(preferredMethod)
      })
    }
  }

  paymentMethods.value = methods
  const preferred = methods.find(item => item.key === preferredMethod)?.key
  if (preferred) {
    selectedPayMethod.value = preferred
  } else if (methods.some(item => item.key === 'qrcode')) {
    selectedPayMethod.value = 'qrcode'
  } else {
    selectedPayMethod.value = methods[0]?.key || 'contact'
  }

  if (selectedPayMethod.value === 'qrcode') {
    loadQrCodes()
  }
}

onLoad((options) => {
  if (options.amount) amount.value = options.amount
  if (options.orderNo) orderNo.value = options.orderNo
  if (options.orderId) orderId.value = options.orderId
  if (options.closeTime) closeTime.value = decodeURIComponent(options.closeTime)

  const preferredMethod = typeof options.payMethod === 'string' ? decodeURIComponent(options.payMethod) : ''
  const preferredLabel = typeof options.payMethodLabel === 'string' ? decodeURIComponent(options.payMethodLabel) : ''
  loadPaymentMethods(preferredMethod, preferredLabel)
  loadPaymentTip()
  loadOrderCloseTime()
})

// ========== 倒计时 ==========
const closeTime = ref('')
const payCountdown = ref('')
let payTimer = null

const loadOrderCloseTime = async () => {
  // 如果没有传入closeTime，从订单接口获取
  if (!closeTime.value && orderId.value) {
    try {
      const res = await selfOrder(orderId.value)
      if (res.code === 0 && res.data && res.data.closeTime) {
        closeTime.value = res.data.closeTime
      }
    } catch (e) {}
  }
  if (closeTime.value) startPayCountdown()
}

const startPayCountdown = () => {
  if (payTimer) clearInterval(payTimer)
  updatePayCountdown()
  payTimer = setInterval(updatePayCountdown, 1000)
}

const updatePayCountdown = () => {
  if (!closeTime.value) { payCountdown.value = ''; return }
  const remain = Math.max(0, Math.floor((new Date(closeTime.value).getTime() - Date.now()) / 1000))
  if (remain <= 0) {
    payCountdown.value = ''
    if (payTimer) clearInterval(payTimer)
    uni.showToast({ title: $t.value('payTimeout'), icon: 'none' })
    setTimeout(() => {
      goToOrders()
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

// ========== 离开页面（不取消订单） ==========
const goToOrders = (isTab = false) => {
  if (isTab) { uni.switchTab({ url: '/pages/tabBar/index' }) } else { uni.redirectTo({ url: '/pages/order/order' }) }
}

// ========== 返回拦截 ==========
onBackPress(() => {
  goToOrders()
  return true // 阻止默认返回，直接跳转订单列表
})

const loadQrCodes = async () => {
  try {
    const res = await getEnabledQrcodePayments()
    if (res.code === 0 && res.data) {
      // API返回的data可能是数组（直接是列表）或对象（含list字段）
      const list = Array.isArray(res.data) ? res.data : (res.data.list || [])
      qrList.value = list
      currentQrIndex.value = 0
    }
    // 如果多码列表为空，尝试旧的单码配置作为兜底
    if (qrList.value.length === 0) {
      const res2 = await request({ url: '/sysConfig/getSysConfigByKey', method: 'get', params: { configKey: 'payment_qr_code' } })
      if (res2.code === 0 && res2.data) {
        qrList.value = [{ name: $t.value('payByQrcode'), image: res2.data, externalPath: '' }]
        currentQrIndex.value = 0
      }
    }
  } catch (e) {
    console.error('Failed to load payment QR codes', e)
  }
}

const loadPaymentTip = async () => {
  try {
    const [tipRes, sizeRes, colorRes] = await Promise.all([
      request({ url: '/sysConfig/getSysConfigByKey', method: 'get', params: { configKey: 'payment_tip_text' } }),
      request({ url: '/sysConfig/getSysConfigByKey', method: 'get', params: { configKey: 'payment_tip_text_size' } }),
      request({ url: '/sysConfig/getSysConfigByKey', method: 'get', params: { configKey: 'payment_tip_text_color' } }),
    ])
    if (tipRes.code === 0 && tipRes.data) {
      paymentTipText.value = localText(tipRes.data)
    }
    if (sizeRes.code === 0 && sizeRes.data) {
      paymentTipSize.value = Number(sizeRes.data) || 14
    }
    if (colorRes.code === 0 && colorRes.data) {
      paymentTipColor.value = colorRes.data
    }
  } catch (e) { /* ignore */ }
}

const previewQr = () => {
  if (currentQrUrl.value) {
    uni.previewImage({ urls: [currentQrUrl.value] })
  }
}

const saveQr = () => {
  if (!currentQrUrl.value) {
    uni.showToast({ title: $t.value('noQrCode'), icon: 'none' })
    return
  }
  uni.downloadFile({
    url: currentQrUrl.value,
    success: (downloadRes) => {
      if (downloadRes.statusCode === 200) {
        uni.saveImageToPhotosAlbum({
          filePath: downloadRes.tempFilePath,
          success: () => {
            uni.showToast({ title: $t.value('saveSuccess'), icon: 'success' })
          },
          fail: () => {
            uni.showToast({ title: $t.value('saveFailed'), icon: 'none' })
          }
        })
      }
    }
  })
}

const goKefu = () => {
  const orderID = encodeURIComponent(String(orderNo.value || orderId.value || ''))
  const payMethod = encodeURIComponent(String(selectedPayMethod.value || ''))
  const payMethodLabel = encodeURIComponent(String(selectedPayMethodLabel.value || ''))
  uni.navigateTo({ url: `/pages/kefu/index?orderID=${orderID}&payMethod=${payMethod}&payMethodLabel=${payMethodLabel}` })
}

const confirmPaid = () => {
  if (!isQrMethod.value) {
    if (selectedPayMethod.value === 'contact') {
      goKefu()
      return
    }

    uni.showModal({
      title: selectedPayMethodLabel.value || $t.value('paymentManualFallbackTitle'),
      content: `${manualFallbackTip.value}\n\n${$t.value('paymentManualProofHint')}`,
      confirmText: $t.value('paymentManualFallbackContact'),
      cancelText: $t.value('paymentManualFallbackLater'),
      success: (res) => {
        if (res.confirm) {
          goKefu()
          return
        }
        uni.showToast({ title: $t.value('paymentManualSavedOrder'), icon: 'none' })
        setTimeout(() => {
          if (orderId.value) {
            uni.redirectTo({ url: `/pages/orderDetail/orderDetail?orderID=${orderId.value}` })
            return
          }
          uni.redirectTo({ url: '/pages/order/order' })
        }, 1200)
      }
    })
    return
  }

  uni.showToast({ title: $t.value('paidSuccess'), icon: 'success' })
  setTimeout(() => {
    uni.redirectTo({ url: '/pages/order/order' })
  }, 1500)
}

const goHome = () => {
  goToOrders(true)
}

const goBack = () => {
  goToOrders()
}
</script>

<style lang="scss" scoped>
.nf-pay {
  --nf-text-secondary: rgba(15, 23, 42, 0.62);
  --nf-text-tertiary: rgba(15, 23, 42, 0.52);
  --nf-text-weak: rgba(15, 23, 42, 0.45);
  min-height: 100vh;
  background: #f4f7fb;
  position: relative;
}
.nf-pay-bg {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background:
    radial-gradient(120% 80% at 100% -10%, rgba(14, 165, 233, 0.2) 0%, transparent 60%),
    radial-gradient(120% 80% at 0% 0%, rgba(37, 99, 235, 0.2) 0%, transparent 60%);
  z-index: 0;
}
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(244, 247, 251, 0.9);
  backdrop-filter: blur(20rpx);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.24);
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
  border-radius: 50%;
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.3);
  box-shadow: 0 10rpx 24rpx rgba(15, 23, 42, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  &:active {
    transform: scale(0.96);
    box-shadow: 0 6rpx 16rpx rgba(15, 23, 42, 0.1);
  }
}
.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
}
.nf-pay-content {
  position: relative;
  z-index: 1;
  padding: calc(var(--status-bar-height, 44rpx) + 88rpx + 30rpx) 24rpx 72rpx;
}
.nf-pay-amount-card {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 36rpx;
  text-align: center;
  margin-bottom: 20rpx;
}
.nf-pay-method-card {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
}
.nf-pay-method-tabs {
  display: flex;
  gap: 12rpx;
  flex-wrap: wrap;
}
.nf-pay-method-tab {
  padding: 10rpx 22rpx;
  border-radius: 28rpx;
  font-size: 24rpx;
  color: var(--nf-text-secondary);
  border: 1rpx solid rgba(148, 163, 184, 0.3);
  background: #f8fafc;
  transition: transform 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
  &:active { transform: scale(0.97); }
}
.nf-pay-method-tab.active {
  color: #2563eb;
  background: rgba(37, 99, 235, 0.1);
  border-color: rgba(37, 99, 235, 0.4);
}
.nf-pay-amount-label {
  font-size: 26rpx;
  color: var(--nf-text-tertiary);
  display: block;
  margin-bottom: 12rpx;
}
.nf-pay-amount-value {
  font-size: 56rpx;
  font-weight: bold;
  color: #2563eb;
  display: block;
}
.nf-pay-order-no {
  font-size: 22rpx;
  color: var(--nf-text-weak);
  margin-top: 12rpx;
  display: block;
}
.nf-pay-countdown {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  margin-top: 20rpx;
  padding: 16rpx 24rpx;
  background: rgba(37, 99, 235, 0.08);
  border: 1rpx solid rgba(37, 99, 235, 0.2);
  border-radius: 12rpx;
}
.nf-pay-countdown-label {
  font-size: 24rpx;
  color: var(--nf-text-tertiary);
}
.nf-pay-countdown-time {
  font-size: 32rpx;
  font-weight: 700;
  color: #2563eb;
  font-variant-numeric: tabular-nums;
}
.nf-pay-qr-card {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 30rpx;
  text-align: center;
  margin-bottom: 20rpx;
}
.nf-pay-qr-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 20rpx;
  display: block;
}
.nf-pay-qr-wrap {
  display: flex;
  justify-content: center;
  margin-bottom: 16rpx;
}
.nf-pay-qr-img {
  width: 400rpx;
  height: 400rpx;
  border-radius: 12rpx;
  background: #fff;
}
.nf-pay-qr-placeholder {
  width: 400rpx;
  height: 400rpx;
  border-radius: 12rpx;
  background: #f8fafc;
  border: 1rpx dashed rgba(148, 163, 184, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-pay-qr-tip {
  font-size: 24rpx;
  color: var(--nf-text-weak);
}
.nf-pay-manual-card {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}
.nf-pay-manual-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
  display: block;
  margin-bottom: 18rpx;
}
.nf-pay-manual-tip {
  font-size: 26rpx;
  line-height: 1.7;
  color: var(--nf-text-secondary);
  display: block;
}
.nf-pay-manual-proof {
  margin-top: 16rpx;
  font-size: 24rpx;
  line-height: 1.7;
  color: rgba(37, 99, 235, 0.9);
  display: block;
}
.nf-pay-qr-tabs {
  display: flex;
  justify-content: center;
  gap: 16rpx;
  margin-bottom: 20rpx;
  flex-wrap: wrap;
}
.nf-pay-qr-tab {
  padding: 8rpx 24rpx;
  border-radius: 24rpx;
  font-size: 24rpx;
  color: var(--nf-text-tertiary);
  background: #f8fafc;
  border: 1rpx solid rgba(148, 163, 184, 0.3);
  transition: transform 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
  &:active { transform: scale(0.97); }
}
.nf-pay-qr-tab.active {
  color: #2563eb;
  background: rgba(37, 99, 235, 0.1);
  border-color: rgba(37, 99, 235, 0.4);
}
.nf-pay-tip-card {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 24rpx 30rpx;
  margin-bottom: 20rpx;
  text-align: center;
}
.nf-pay-tip-text {
  line-height: 1.6;
}
.nf-pay-steps {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
}
.nf-pay-steps-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 20rpx;
  display: block;
}
.nf-pay-step {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 16rpx;
}
.nf-pay-step:last-child {
  margin-bottom: 0;
}
.nf-pay-step-num {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
  color: #fff;
  font-weight: bold;
  flex-shrink: 0;
}
.nf-pay-step-text {
  font-size: 26rpx;
  color: var(--nf-text-secondary);
}
.nf-pay-actions {
  display: flex;
  gap: 20rpx;
  margin-bottom: 24rpx;
}
.nf-pay-btn {
  flex: 1;
  height: 92rpx;
  border-radius: 14rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  font-size: 28rpx;
  font-weight: 700;
  transition: transform 0.2s ease, filter 0.2s ease;
  &:active {
    transform: translateY(1rpx);
    filter: brightness(0.96);
  }
}
.nf-pay-btn-primary {
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
  color: #fff;
}
.nf-pay-btn-kefu {
  background: rgba(37, 99, 235, 0.1);
  border: 1rpx solid rgba(37, 99, 235, 0.3);
  color: #2563eb;
}
.nf-pay-confirm {
  text-align: center;
  padding: 20rpx;
}
.nf-pay-confirm text {
  font-size: 28rpx;
  color: var(--nf-text-weak);
  text-decoration: underline;
}
.nf-pay-bottom-actions {
  display: flex;
  gap: 24rpx;
  margin-top: 10rpx;
}
.nf-pay-bottom-btn {
  flex: 1;
  height: 84rpx;
  border-radius: 42rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 700;
  transition: transform 0.2s ease, filter 0.2s ease;
  &:active {
    transform: translateY(1rpx);
    filter: brightness(0.97);
  }
}
.nf-pay-bottom-home {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.28);
  color: var(--nf-text-secondary);
}
.nf-pay-bottom-paid {
  background: rgba(37, 99, 235, 0.12);
  border: 1rpx solid rgba(37, 99, 235, 0.32);
  color: #2563eb;
}
</style>
