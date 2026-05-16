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
        <text class="nf-pay-amount-value">{{ displayCs }}{{ displayAmount }}</text>
        <text class="nf-pay-order-no" v-if="orderNo">{{ $t('orderNo') }}: {{ orderNo }}</text>
        <!-- 倒计时 -->
        <view class="nf-pay-countdown" v-if="payCountdown">
          <text class="nf-pay-countdown-label">{{ $t('remainPayTime') }}</text>
          <text class="nf-pay-countdown-time">{{ payCountdown }}</text>
        </view>
      </view>

      <view class="nf-pay-info-card" v-if="purchaseInfoText">
        <text class="nf-pay-info-title">{{ $t('orderDetail') }}</text>
        <text class="nf-pay-info-text">{{ purchaseInfoText }}</text>
      </view>

      <view class="nf-pay-carry-tip" v-if="orderNo">{{ $t('orderDetailKefuCarryTip') }}</view>

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
        <view class="nf-pay-method-hint" v-if="isQrMethod">
          <text>{{ qrFallbackHint }}</text>
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
          <LazyImage
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
        <text class="nf-pay-manual-title">{{ manualCardTitle }}</text>

        <view class="nf-pay-manual-methods" v-if="preferredPayMethods.length">
          <text class="nf-pay-manual-methods-title">{{ $t('preferredPayMethod') }}</text>
          <view class="nf-pay-preferred-grid">
            <view
              v-for="item in preferredPayMethods"
              :key="item.key"
              class="nf-pay-preferred-item"
              :class="{ active: selectedPreferredPayMethod === item.key }"
              @tap="selectPreferredPayMethod(item.key)"
            >
              <LazyImage
                v-if="getPreferredPayMethodImage(item)"
                class="nf-pay-preferred-item-img"
                :src="getPreferredPayMethodImage(item)"
                mode="aspectFit"
              />
              <view v-else class="nf-pay-preferred-item-fallback">
                <text>{{ (item.label || '?').slice(0, 1) }}</text>
              </view>
              <text class="nf-pay-preferred-item-label">{{ item.label }}</text>
            </view>
          </view>
        </view>
        <text class="nf-pay-manual-empty" v-else>{{ $t('paymentPreferredMethodMissing') }}</text>

        <view class="nf-pay-copy-card" v-if="selectedPreferredPayMethodDraftText">
          <view class="nf-pay-copy-head">
            <text class="nf-pay-copy-title">{{ $t('kefuPaymentDraftLabel') }}</text>
            <view class="nf-pay-copy-btn" @tap="copyPreferredPayDraft">
              <text>{{ $t('copy') }}</text>
            </view>
          </view>
          <text class="nf-pay-copy-text">{{ selectedPreferredPayMethodDraftText }}</text>
        </view>

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
import { getPaymentConfig, getUniPreferredPayConfig } from '@/api/sysConfig.js'
import { localText, resolveApiMessage } from '@/utils/i18n'
import { resolveLocalizedPriceFen } from '@/utils/price-i18n.js'
import { selfOrder, updateOrder, updateOrderStatus } from '@/api/order.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'
import {
  selfTryonRechargeOrder,
  updateTryonRechargeOrderPayMethod,
  submitTryonRechargeOrderPayment,
} from '@/api/tryonRechargeOrder.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const orderType = ref('shop')
const orderCurrencySymbol = ref('')
const displayCs = computed(() => {
  const snapshot = String(orderCurrencySymbol.value || '').trim()
  return snapshot || cs.value
})
const $t = computed(() => langStore.$t)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const amount = ref('0.00')
const orderNo = ref('')
const orderId = ref('')
const orderStatus = ref('0')
const rechargePoints = ref(0)
const orderSummaryName = ref('')
const shopOrder = ref(null)
const paymentMethods = ref([])
const preferredPayMethods = ref([])
const selectedPayMethod = ref('qrcode')
const selectedPreferredPayMethod = ref('')
const isRechargeOrder = computed(() => orderType.value === 'recharge')

const parsePlanSnapshot = (rawValue) => {
  if (!rawValue) return null
  if (typeof rawValue === 'object') {
    return rawValue
  }
  const text = String(rawValue || '').trim()
  if (!text) return null
  try {
    const parsed = JSON.parse(text)
    return parsed && typeof parsed === 'object' ? parsed : null
  } catch (e) {
    return null
  }
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

const calculateShopLocalizedTotalFen = (order) => {
  if (!order || !Array.isArray(order.detail)) {
    return null
  }

  const priceLocale = getOrderPriceLocale(order)

  const localizedOrigin = order.detail.reduce((sum, detail) => {
    const priceFen = resolveLocalizedPriceFen(getDetailBasePriceFen(detail), getDetailPriceI18nSnapshot(detail), priceLocale)
    const quantity = Math.max(0, Number(detail?.quantity || 0))
    return sum + Math.max(0, priceFen) * quantity
  }, 0)

  let total = localizedOrigin - Number(order.discount || 0)
  if (order.usePoints) {
    total -= Number(order.pointsUsed || 0)
  }
  return Math.max(0, total)
}

const displayAmount = computed(() => {
  if (isRechargeOrder.value) {
    return amount.value
  }
  const localizedFen = calculateShopLocalizedTotalFen(shopOrder.value)
  if (localizedFen === null) {
    return amount.value
  }
  return (localizedFen / 100).toFixed(2)
})

const purchaseInfoText = computed(() => {
  if (isRechargeOrder.value) {
    if (rechargePoints.value > 0) {
      return `${$t.value('tryonCoins')}: ${rechargePoints.value}`
    }
    return ''
  }
  return orderSummaryName.value
})

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

const getPaymentMethodLabel = (payMethod, name, label) => {
  const localizedName = String(localText(name, langStore.locale) || localText(label, langStore.locale) || '').trim()
  if (localizedName) {
    return localizedName
  }
  const plainLabel = typeof label === 'string' ? label.trim() : ''
  return paymentMethodLabelMap.value[payMethod] || plainLabel || payMethod || $t.value('contactCustomerService')
}

const buildDefaultRawMethods = () => ([
  { key: 'qrcode', label: getPaymentMethodLabel('qrcode'), manual: true, copyText: {} },
  { key: 'contact', label: getPaymentMethodLabel('contact'), manual: true, copyText: {} }
])

const normalizePaymentMethods = (methods) => {
  if (!Array.isArray(methods)) {
    return []
  }
  return methods
    .filter(item => item && typeof item.key === 'string' && item.key)
    .map(item => ({
      key: item.key,
      label: getPaymentMethodLabel(item.key, item.name, item.label),
      manual: Boolean(item.manual),
      name: item.name || {},
      copyText: item.copyText || {}
    }))
}

const buildModePaymentMethods = (methods) => {
  const hasQrcode = methods.some(item => item.key === 'qrcode')
  const hasContact = methods.some(item => item.key === 'contact') || methods.some(item => item.manual && item.key !== 'qrcode')
  const list = []
  if (hasQrcode) {
    list.push({ key: 'qrcode', label: getPaymentMethodLabel('qrcode') })
  }
  if (hasContact) {
    list.push({ key: 'contact', label: getPaymentMethodLabel('contact') })
  }
  return list.length > 0 ? list : [{ key: 'contact', label: getPaymentMethodLabel('contact') }]
}

const normalizePreferredPayMethods = (methods) => {
  if (!Array.isArray(methods)) {
    return []
  }
  return methods
    .filter(item => item && typeof item.key === 'string')
    .map(item => ({
      key: String(item.key || '').trim().toLowerCase(),
      label: getPaymentMethodLabel(String(item.key || '').trim().toLowerCase(), item.name, item.label),
      name: item.name || {},
      image: String(item.image || item.externalPath || '').trim(),
      copyText: item.copyText || {}
    }))
    .filter(item => item.key && item.key !== 'qrcode' && item.key !== 'contact')
}

const parseSkuSpecItems = (payload) => {
  if (Array.isArray(payload)) {
    return payload
  }
  if (typeof payload === 'string' && payload.trim()) {
    try {
      const parsed = JSON.parse(payload)
      return Array.isArray(parsed) ? parsed : []
    } catch {
      return []
    }
  }
  return []
}

const getPreferredPayMethodImage = (method) => {
  if (!method) return ''
  const imagePath = String(method.image || method.externalPath || '').trim()
  if (!imagePath) return ''
  return getUrl(imagePath)
}

const selectPreferredPayMethod = (payMethod) => {
  selectedPreferredPayMethod.value = payMethod
}

const isQrMethod = computed(() => selectedPayMethod.value === 'qrcode')
const selectedPayMethodLabel = computed(() => {
  const found = paymentMethods.value.find(item => item.key === selectedPayMethod.value)
  return found?.label || getPaymentMethodLabel(selectedPayMethod.value)
})
const selectedPreferredPayMethodLabel = computed(() => {
  const found = preferredPayMethods.value.find(item => item.key === selectedPreferredPayMethod.value)
  return found?.label || ''
})
const manualCardTitle = computed(() => {
  if (selectedPreferredPayMethodLabel.value) {
    return selectedPreferredPayMethodLabel.value
  }
  if (selectedPayMethod.value === 'contact') {
    return $t.value('preferredPayMethod')
  }
  return selectedPayMethodLabel.value
})
const selectedPreferredPayMethodCopyText = computed(() => {
  const found = preferredPayMethods.value.find(item => item.key === selectedPreferredPayMethod.value)
  return localText(found?.copyText, langStore.locale) || ''
})
const selectedPreferredPayMethodDraftText = computed(() => {
  if (!selectedPreferredPayMethod.value) {
    return ''
  }
  const orderRef = String(orderNo.value || orderId.value || '').trim() || '-'
  const payLabel = selectedPreferredPayMethodLabel.value || selectedPayMethodLabel.value || '-'
  const template = selectedPreferredPayMethodCopyText.value || $t.value('kefuPaymentDraftTemplatePending')
  return String(template)
    .replace('{orderID}', orderRef)
    .replace('{payMethod}', payLabel)
})
const supportedPreferredMethodsText = computed(() => {
  return preferredPayMethods.value.map(item => item.label).join(' / ')
})
const qrFallbackHint = computed(() => {
  return $t.value('paymentQrFallbackTip').replace('{methods}', supportedPreferredMethodsText.value || $t.value('contactCustomerService'))
})

// 多码支持
const qrList = ref([])
const currentQrIndex = ref(0)
const currentQrLabel = computed(() => {
  const qr = qrList.value[currentQrIndex.value]
  if (!qr) return ''
  return localText(qr.nameI18n, langStore.locale) || qr.name || getPaymentMethodLabel('qrcode')
})
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
  if (selectedPayMethod.value === 'contact' && !selectedPreferredPayMethod.value) {
    return $t.value('paymentPreferredMethodMissing')
  }
  const currentMethod = selectedPreferredPayMethod.value || selectedPayMethod.value
  const currentMethodLabel = selectedPreferredPayMethodLabel.value || selectedPayMethodLabel.value
  return getManualFallbackTip(currentMethod, currentMethodLabel)
})

const steps = computed(() => {
  if (isQrMethod.value) {
    return [$t.value('payStep1'), $t.value('payStep2'), $t.value('payStep3')]
  }
  return [manualFallbackTip.value, $t.value('paymentManualProofHint'), $t.value('contactCustomerService')]
})

const syncOrderPayMethod = async (payMethod) => {
  if (!orderId.value || !payMethod) return
  try {
    if (isRechargeOrder.value) {
      await updateTryonRechargeOrderPayMethod({ id: Number(orderId.value), payMethod })
      return
    }
    await updateOrder({ ID: Number(orderId.value), payMethod })
  } catch (e) {
    uni.showToast({ title: resolveApiMessage(e?.message, 'operationFailed'), icon: 'none' })
  }
}

const selectPayMethod = async (payMethod) => {
  selectedPayMethod.value = payMethod
  if (payMethod === 'contact' && !selectedPreferredPayMethod.value && preferredPayMethods.value.length > 0) {
    selectPreferredPayMethod(preferredPayMethods.value[0].key)
  }
  if (payMethod === 'qrcode' && qrList.value.length === 0) {
    loadQrCodes()
  }
  await syncOrderPayMethod(payMethod)
}

const loadPaymentMethods = async (preferredMethod, preferredLabel) => {
  const normalizedPreferredMethod = String(preferredMethod || '').trim().toLowerCase()
  let allMethods = []
  let preferredMethods = []

  const [paymentConfigResult, preferredConfigResult] = await Promise.allSettled([
    getPaymentConfig(),
    getUniPreferredPayConfig()
  ])

  if (paymentConfigResult.status === 'fulfilled') {
    const paymentRes = paymentConfigResult.value
    if (paymentRes.code === 0 && paymentRes.data) {
      allMethods = normalizePaymentMethods(paymentRes.data.methods)
    }
  }

  if (preferredConfigResult.status === 'fulfilled') {
    const preferredRes = preferredConfigResult.value
    if (preferredRes.code === 0 && preferredRes.data) {
      const preferredList = Array.isArray(preferredRes.data)
        ? preferredRes.data
        : (preferredRes.data.methods || preferredRes.data.list || [])
      const normalizedPreferred = normalizePreferredPayMethods(preferredList)
      if (normalizedPreferred.length > 0) {
        preferredMethods = normalizedPreferred
      }
    }
  }

  if (allMethods.length === 0) {
    allMethods = buildDefaultRawMethods()
  }

  preferredPayMethods.value = preferredMethods

  if (normalizedPreferredMethod && normalizedPreferredMethod !== 'qrcode' && normalizedPreferredMethod !== 'contact') {
    const existing = preferredPayMethods.value.find(item => item.key === normalizedPreferredMethod)
    if (!existing) {
      preferredPayMethods.value.unshift({
        key: normalizedPreferredMethod,
        label: preferredLabel || getPaymentMethodLabel(normalizedPreferredMethod),
        name: {},
        image: '',
        copyText: {}
      })
    }
  }

  paymentMethods.value = buildModePaymentMethods(allMethods)

  if (normalizedPreferredMethod === 'qrcode' && paymentMethods.value.some(item => item.key === 'qrcode')) {
    selectedPayMethod.value = 'qrcode'
  } else if (normalizedPreferredMethod === 'contact' && paymentMethods.value.some(item => item.key === 'contact')) {
    selectedPayMethod.value = 'contact'
  } else if (normalizedPreferredMethod && normalizedPreferredMethod !== 'qrcode' && normalizedPreferredMethod !== 'contact') {
    selectedPayMethod.value = 'contact'
    if (preferredPayMethods.value.some(item => item.key === normalizedPreferredMethod)) {
      selectPreferredPayMethod(normalizedPreferredMethod)
    }
  } else if (paymentMethods.value.some(item => item.key === 'qrcode')) {
    selectedPayMethod.value = 'qrcode'
  } else {
    selectedPayMethod.value = paymentMethods.value[0]?.key || 'contact'
  }

  if (selectedPayMethod.value === 'contact' && !selectedPreferredPayMethod.value && preferredPayMethods.value.length > 0) {
    selectPreferredPayMethod(preferredPayMethods.value[0].key)
  }

  if (selectedPayMethod.value === 'qrcode') {
    loadQrCodes()
  }
}

onLoad((options) => {
  if (options.amount) amount.value = decodeURIComponent(options.amount)
  if (options.orderNo) orderNo.value = decodeURIComponent(options.orderNo)
  if (options.orderId) orderId.value = decodeURIComponent(options.orderId)
  if (options.orderType) {
    const normalizedType = String(options.orderType || '').trim().toLowerCase()
    orderType.value = normalizedType === 'recharge' ? 'recharge' : 'shop'
  }
  if (options.rechargePoints) {
    rechargePoints.value = Number(decodeURIComponent(options.rechargePoints) || 0)
  }
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
  if (orderId.value) {
    try {
      if (isRechargeOrder.value) {
        shopOrder.value = null
        const rechargeRes = await selfTryonRechargeOrder(orderId.value)
        const order = rechargeRes?.data?.order || rechargeRes?.data || {}
        orderStatus.value = String(order.status || order.Status || orderStatus.value || '0')
        orderCurrencySymbol.value = String(order.settlementCurrencySymbol || order.SettlementCurrencySymbol || '').trim()
        const latestCloseTime = String(order.closeTime || order.CloseTime || '').trim()
        if (!closeTime.value && latestCloseTime) {
          closeTime.value = latestCloseTime
        }
        if (!orderNo.value) {
          orderNo.value = String(order.outTradeNo || order.OutTradeNo || order.ID || order.id || '')
        }
        const points = Number(order.points || order.Points || 0)
        if (points > 0) {
          rechargePoints.value = points
        }

        const orderAmountCents = Number(order.amount || order.Amount)
        if (Number.isFinite(orderAmountCents) && orderAmountCents >= 0) {
          amount.value = (orderAmountCents / 100).toFixed(2)
        } else {
          const snapshot = parsePlanSnapshot(order.planSnapshot || order.PlanSnapshot)
          if (snapshot) {
            if (!orderCurrencySymbol.value) {
              orderCurrencySymbol.value = String(snapshot.settlementCurrencySymbol || '').trim()
            }
            const settlementLocale = String(order.settlementCurrency || order.SettlementCurrency || snapshot.settlementCurrency || locale.value || 'zh')
              .replace(/_/g, '-')
            const snapshotPriceFen = resolveLocalizedPriceFen(snapshot.price, snapshot.priceI18n, settlementLocale)
            const selectedPriceFen = Number(snapshot.selectedPriceFen)
            const amountFen = Number.isFinite(Number(snapshotPriceFen)) && Number(snapshotPriceFen) >= 0
              ? Math.round(Number(snapshotPriceFen))
              : (Number.isFinite(selectedPriceFen) && selectedPriceFen >= 0 ? Math.round(selectedPriceFen) : 0)
            amount.value = (amountFen / 100).toFixed(2)
          }
        }
      } else {
        const res = await selfOrder(orderId.value)
        const order = res?.data || {}
        shopOrder.value = order
        orderStatus.value = String(order.status || order.Status || orderStatus.value || '0')
        orderCurrencySymbol.value = String(order.settlementCurrencySymbol || order.SettlementCurrencySymbol || '').trim()
        if (!closeTime.value && order.closeTime) {
          closeTime.value = order.closeTime
        }
        if (!orderNo.value) {
          orderNo.value = String(order.outTradeNo || order.OutTradeNo || order.ID || '')
        }
        const detail = Array.isArray(order.detail) ? order.detail : []
        if (detail.length > 0) {
          orderSummaryName.value = detail.map(item => {
            const sku = item?.sku || {}
            const good = item?.good || {}
            const name = localText(sku?.nameI18n || sku?.name || good?.titleI18n || good?.title || good?.nameI18n || good?.name, langStore.locale) || ''
            const specs = [...parseSkuSpecItems(sku?.specs), ...parseSkuSpecItems(sku?.attrs)]
              .map(spec => {
                const label = localText(spec?.labelI18n || spec?.nameI18n || spec?.label || spec?.name, langStore.locale)
                const value = localText(spec?.valueI18n || spec?.value, langStore.locale)
                if (!label && !value) return ''
                return label ? `${label}:${value}` : value
              })
              .filter(Boolean)
              .join(', ')
            const quantity = Number(item?.quantity || 1)
            const quantityText = quantity > 1 ? ` x${quantity}` : ''
            return `${name}${quantityText}${specs ? ` (${specs})` : ''}`
          }).filter(Boolean).join(' ; ')
        }
      }
    } catch (e) {}
  }
  if (closeTime.value && orderStatus.value === '0') startPayCountdown()
}

const startPayCountdown = () => {
  if (payTimer) clearInterval(payTimer)
  updatePayCountdown()
  payTimer = setInterval(updatePayCountdown, 1000)
}

const updatePayCountdown = () => {
  if (orderStatus.value !== '0') {
    payCountdown.value = ''
    if (payTimer) clearInterval(payTimer)
    return
  }
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
  if (isRechargeOrder.value) {
    uni.switchTab({ url: '/pages/tabBar/my/index' })
    return
  }
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

const copyPreferredPayDraft = () => {
  if (!selectedPreferredPayMethodDraftText.value) {
    uni.showToast({ title: $t.value('selectPreferredPayMethodFirst'), icon: 'none' })
    return
  }
  uni.setClipboardData({
    data: selectedPreferredPayMethodDraftText.value,
    success: () => {
      uni.showToast({ title: $t.value('copySuccess'), icon: 'none' })
    }
  })
}

const goKefu = () => {
  if (!isQrMethod.value && preferredPayMethods.value.length === 0) {
    uni.showToast({ title: $t.value('paymentPreferredMethodMissing'), icon: 'none' })
    return
  }
  if (!isQrMethod.value && !selectedPreferredPayMethod.value) {
    uni.showToast({ title: $t.value('selectPreferredPayMethodFirst'), icon: 'none' })
    return
  }
  const preferredMethod = selectedPreferredPayMethod.value || ''
  const preferredMethodLabel = selectedPreferredPayMethodLabel.value || ''
  const preferredMethodCopyText = selectedPreferredPayMethodCopyText.value || ''
  const forwardPayMethodLabel = selectedPayMethod.value === 'contact'
    ? (preferredMethodLabel || '')
    : selectedPayMethodLabel.value
  const qrName = isQrMethod.value ? currentQrLabel.value : ''
  const orderID = encodeURIComponent(String(orderNo.value || orderId.value || ''))
  const payMethod = encodeURIComponent(String(selectedPayMethod.value || ''))
  const payMethodLabel = encodeURIComponent(String(forwardPayMethodLabel || ''))
  const preferredPayMethod = encodeURIComponent(String(preferredMethod))
  const preferredPayMethodLabel = encodeURIComponent(String(preferredMethodLabel))
  const preferredPayMethodCopyText = encodeURIComponent(String(preferredMethodCopyText))
  const selectedQrName = encodeURIComponent(String(qrName || ''))
  uni.navigateTo({
    url: `/pages/kefu/index?orderID=${orderID}&payMethod=${payMethod}&payMethodLabel=${payMethodLabel}&preferredPayMethod=${preferredPayMethod}&preferredPayMethodLabel=${preferredPayMethodLabel}&preferredPayMethodCopyText=${preferredPayMethodCopyText}&selectedQrName=${selectedQrName}`
  })
}

const confirmPaid = async () => {
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
          if (isRechargeOrder.value) {
            goToOrders()
            return
          }
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

  if (!orderId.value) {
    uni.showToast({ title: $t.value('operationFailed'), icon: 'none' })
    return
  }

  uni.showLoading({ title: $t.value('checking'), mask: true })
  try {
    if (isRechargeOrder.value) {
      await submitTryonRechargeOrderPayment({ id: Number(orderId.value) })
    } else {
      await updateOrderStatus({ ID: Number(orderId.value), status: '8' })
    }
    orderStatus.value = '8'
    uni.showToast({ title: $t.value('paidSuccess'), icon: 'success' })
    setTimeout(() => {
      goToOrders()
    }, 1500)
  } catch (e) {
    uni.showToast({ title: resolveApiMessage(e?.message, 'operationFailed'), icon: 'none' })
  } finally {
    uni.hideLoading()
  }
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
.nf-pay-info-card {
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 12rpx 28rpx rgba(15, 23, 42, 0.08);
  border-radius: 20rpx;
  padding: 20rpx 24rpx;
  margin-bottom: 20rpx;
}
.nf-pay-info-title {
  display: block;
  font-size: 24rpx;
  color: var(--nf-text-tertiary);
}
.nf-pay-info-text {
  display: block;
  margin-top: 10rpx;
  font-size: 28rpx;
  color: #0f172a;
  font-weight: 600;
}
.nf-pay-carry-tip {
  margin: 0 0 20rpx;
  padding: 14rpx 18rpx;
  border-radius: 12rpx;
  background: rgba(219, 234, 254, 0.72);
  border: 1rpx solid rgba(37, 99, 235, 0.24);
  font-size: 22rpx;
  color: rgba(30, 64, 175, 0.95);
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
.nf-pay-method-hint {
  margin-top: 14rpx;
  font-size: 22rpx;
  color: rgba(30, 64, 175, 0.9);
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
.nf-pay-manual-methods {
  margin-bottom: 16rpx;
}
.nf-pay-manual-methods-title {
  display: block;
  margin-bottom: 10rpx;
  font-size: 24rpx;
  color: var(--nf-text-tertiary);
}
.nf-pay-preferred-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
}
.nf-pay-preferred-item {
  display: flex;
  align-items: center;
  gap: 12rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(148, 163, 184, 0.3);
  background: #f8fafc;
  padding: 14rpx 16rpx;
  transition: transform 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
  &:active {
    transform: scale(0.98);
  }
}
.nf-pay-preferred-item.active {
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(37, 99, 235, 0.08);
}
.nf-pay-preferred-item-img {
  width: 48rpx;
  height: 48rpx;
  border-radius: 10rpx;
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.28);
}
.nf-pay-preferred-item-fallback {
  width: 48rpx;
  height: 48rpx;
  border-radius: 10rpx;
  background: rgba(37, 99, 235, 0.12);
  color: #1d4ed8;
  font-size: 24rpx;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-pay-preferred-item-label {
  flex: 1;
  min-width: 0;
  font-size: 24rpx;
  color: #0f172a;
  line-height: 1.4;
}
.nf-pay-copy-card {
  margin-top: 16rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.24);
  border-radius: 14rpx;
  background: rgba(219, 234, 254, 0.5);
  padding: 16rpx;
}
.nf-pay-copy-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10rpx;
}
.nf-pay-copy-title {
  font-size: 22rpx;
  color: rgba(30, 64, 175, 0.95);
}
.nf-pay-copy-btn {
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(37, 99, 235, 0.35);
  color: #1d4ed8;
  font-size: 22rpx;
  background: #ffffff;
}
.nf-pay-copy-text {
  display: block;
  font-size: 23rpx;
  line-height: 1.6;
  color: #1e293b;
  word-break: break-all;
}
.nf-pay-manual-empty {
  display: block;
  margin-bottom: 10rpx;
  font-size: 24rpx;
  color: rgba(220, 38, 38, 0.85);
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
