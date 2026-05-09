<template>
  <view class="my-page">
    <view class="bg-glow"></view>

    <view class="profile-card">
      <view class="profile-main" v-if="isLogin" @tap="goProfile">
        <image class="avatar" :src="avatarUrl" mode="aspectFill" />
        <view class="profile-info">
          <text class="name">{{ userName }}</text>
          <text class="phone">{{ phoneText }}</text>
        </view>
        <uni-icons type="right" size="16" color="rgba(15,23,42,0.35)" />
      </view>

      <view class="guest-main" v-else @tap="goLogin">
        <view class="guest-left">
          <view class="guest-icon">?</view>
          <view>
            <text class="guest-title">{{ $t('guestModeTitle') }}</text>
            <text class="guest-sub">{{ $t('guestModeDesc') }}</text>
          </view>
        </view>
        <uni-icons type="right" size="16" color="rgba(15,23,42,0.35)" />
      </view>

      <view class="point-card">
        <view>
          <text class="point-label">{{ $t('tryonCoins') }}</text>
          <text class="point-value">{{ userPoints }}</text>
        </view>
        <view class="recharge-btn" @tap="openRecharge">{{ $t('recharge') }}</view>
      </view>
    </view>

    <view class="section">
      <view class="section-title">{{ $t('commonFeatures') }}</view>
      <view class="quick-grid">
        <view class="quick-item" @tap="goTryonRoom">
          <uni-icons type="image" size="24" color="#2563eb" />
          <text>{{ $t('tryonRoom') }}</text>
        </view>
        <view class="quick-item" @tap="goShoeRoom">
          <uni-icons type="star" size="24" color="#2563eb" />
          <text>{{ $t('shoeRoom') }}</text>
        </view>
        <view class="quick-item" @tap="goHistory">
          <uni-icons type="reload" size="24" color="#2563eb" />
          <text>{{ $t('tryonHistory') }}</text>
        </view>
        <view class="quick-item" @tap="goClothesPage">
          <uni-icons type="shop" size="24" color="#2563eb" />
          <text>{{ $t('clothesPage') }}</text>
        </view>
      </view>
    </view>

    <view class="section">
      <view class="section-title">{{ $t('otherFeatures') }}</view>
      <view class="menu-list">
        <view class="menu-item" @tap="goKefu">
          <text class="menu-text">{{ $t('kefuContact') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goOrder">
          <text class="menu-text">{{ $t('myOrders') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goTryonPointRecord">
          <text class="menu-text">{{ $t('tryonPointRecord') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goTryonRechargeRecord">
          <text class="menu-text">{{ $t('tryonRechargeRecord') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goInvite">
          <text class="menu-text">{{ $t('inviteFriends') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goCollect">
          <text class="menu-text">{{ $t('myCollection') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goMyModel">
          <text class="menu-text">{{ $t('myModels') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goMyCloth">
          <text class="menu-text">{{ $t('myCloset') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="goUseRecord">
          <text class="menu-text">{{ $t('browseHistory') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="openLangSwitch">
          <text class="menu-text">{{ $t('switchLang') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
        <view class="menu-item" @tap="logoutDevice" v-if="isLogin">
          <text class="menu-text danger">{{ $t('logoutDevice') }}</text>
          <uni-icons type="right" size="14" color="rgba(15,23,42,0.35)" />
        </view>
      </view>
    </view>

    <view class="popup-mask" v-if="showRecharge" @tap="showRecharge = false">
      <view class="popup-panel" @tap.stop>
        <view class="popup-title">{{ $t('rechargeTryonCoins') }}</view>
        <view class="popup-list">
          <view class="popup-item" v-for="(item, index) in rechargePlans" :key="item.key || index" @tap="selectRecharge(item)">
            <view>
              <text class="item-title">{{ formatRechargePoints(item) }}{{ formatRechargeCoinLabel(item) }}</text>
              <text class="item-sub">{{ formatRechargePrice(item) }}</text>
            </view>
            <uni-icons type="right" size="16" color="rgba(15,23,42,0.35)" />
          </view>
        </view>
        <view class="popup-close" @tap="showRecharge = false">{{ $t('cancel') }}</view>
      </view>
    </view>

    <lang-switch v-model="showLangPicker" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/pinia/modules/user.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { getTryonRechargePlans } from '@/api/sysConfig.js'
import { getPaymentConfig } from '@/api/sysConfig.js'
import { createTryonRechargeOrder } from '@/api/tryonRechargeOrder.js'
import { localText, resolveApiMessage } from '@/utils/i18n.js'

const userStore = useUserStore()
const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const $t = computed(() => langStore.$t)
const cs = computed(() => appConfigStore.currencySymbol || '¥')

const showRecharge = ref(false)
const showLangPicker = ref(false)
const userInfo = ref({})
const paymentMethods = ref([])

const defaultRechargePlans = [
  { points: { zh: '50', en: '50', mn: '50' }, coinLabel: { zh: '试衣币', en: 'Try-on Coins', mn: 'Туршилтын зоос' }, currencySymbol: { zh: '￥', en: 'CNY ', mn: 'CNY ' }, price: { zh: '9.9', en: '9.9', mn: '9.9' }, currencySuffix: { zh: '元', en: '', mn: '' } },
  { points: { zh: '180', en: '180', mn: '180' }, coinLabel: { zh: '试衣币', en: 'Try-on Coins', mn: 'Туршилтын зоос' }, currencySymbol: { zh: '￥', en: 'CNY ', mn: 'CNY ' }, price: { zh: '29.9', en: '29.9', mn: '29.9' }, currencySuffix: { zh: '元', en: '', mn: '' } },
  { points: { zh: '680', en: '680', mn: '680' }, coinLabel: { zh: '试衣币', en: 'Try-on Coins', mn: 'Туршилтын зоос' }, currencySymbol: { zh: '￥', en: 'CNY ', mn: 'CNY ' }, price: { zh: '99.9', en: '99.9', mn: '99.9' }, currencySuffix: { zh: '元', en: '', mn: '' } },
]
const rechargePlans = ref([...defaultRechargePlans])

const toI18nValue = (value, fallback = '') => {
  if (value && typeof value === 'object') return value
  const text = value === undefined || value === null || value === '' ? fallback : String(value)
  return { zh: text, en: text, mn: text }
}

const resolvePlanText = (value, fallback = '') => {
  return localText(toI18nValue(value, fallback), langStore.locale) || fallback
}

const normalizeRechargePlans = (raw) => {
  if (!Array.isArray(raw)) return [...defaultRechargePlans]
  const list = raw.map((item) => {
    const pointsText = resolvePlanText(item?.points, '0')
    const priceText = resolvePlanText(item?.price, '')
    const points = Number(String(pointsText).replace(/[^0-9.]/g, '') || 0)
    const price = String(priceText ?? '').trim()
    if (!points || !price) return null
    return {
      key: `${points}_${price}`,
      points: toI18nValue(item?.points, String(points)),
      coinLabel: toI18nValue(item?.coinLabel || item?.label, $t.value('tryonCoins')),
      currencySymbol: toI18nValue(item?.currencySymbol, cs.value),
      price: toI18nValue(item?.price, price),
      currencySuffix: toI18nValue(item?.currencySuffix || item?.suffix, ''),
      numericPoints: points,
    }
  }).filter(Boolean)
  return list.length > 0 ? list : [...defaultRechargePlans]
}

const formatRechargePoints = (item) => resolvePlanText(item.points, String(item.numericPoints || ''))

const formatRechargeCoinLabel = (item) => resolvePlanText(item.coinLabel, $t.value('tryonCoins'))

const formatRechargePrice = (item) => {
  const symbol = resolvePlanText(item.currencySymbol, cs.value)
  const price = resolvePlanText(item.price, '')
  const suffix = resolvePlanText(item.currencySuffix, '')
  return `${symbol}${price}${suffix}`
}

const normalizePriceText = (value) => {
  const text = String(value || '').trim()
  const normalized = text.replace(/[^0-9.]/g, '')
  if (!normalized) return ''
  return normalized
}

const parseRechargePointsValue = (item) => {
  const text = String(formatRechargePoints(item) || '').trim()
  const normalized = text.replace(/[^0-9]/g, '')
  return Number(normalized || 0)
}

const parseRechargePriceValue = (item) => {
  const localizedPrice = resolvePlanText(item.price, '')
  return normalizePriceText(localizedPrice)
}

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
  { key: 'contact', label: getPaymentMethodLabel('contact') },
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

const loadPaymentMethods = async () => {
  try {
    const res = await getPaymentConfig()
    const methods = normalizePaymentMethods(res?.data?.methods)
    paymentMethods.value = methods.length > 0 ? methods : buildDefaultPaymentMethods()
  } catch (e) {
    paymentMethods.value = buildDefaultPaymentMethods()
  }
}

const loadRechargePlans = async () => {
  try {
    const res = await getTryonRechargePlans()
    if (res.code !== 0) {
      rechargePlans.value = [...defaultRechargePlans]
      return
    }
    const payload = typeof res.data === 'object' ? res.data?.configValue : res.data
    if (!payload) {
      rechargePlans.value = [...defaultRechargePlans]
      return
    }
    const parsed = JSON.parse(payload)
    rechargePlans.value = normalizeRechargePlans(parsed)
  } catch (e) {
    rechargePlans.value = [...defaultRechargePlans]
  }
}

const isLogin = computed(() => {
  return !!(userStore.token || uni.getStorageSync('x-token'))
})

const userName = computed(() => {
  return userInfo.value.nickname || userInfo.value.username || $t.value('unnamedUser')
})

const avatarUrl = computed(() => {
  return userInfo.value.avatar || 'https://mmbiz.qpic.cn/mmbiz/icTdbqWNOwNRna42FI242Lcia07jQodd2FJGIYQfG0LAJGFxM4FbnQP6yfMxBgJ0F3YRqJCJ1aPAK2dQagdusBZg/0'
})

const userPoints = computed(() => Number((userInfo.value.tryonPoint ?? userInfo.value.point) || 0))

const phoneText = computed(() => {
  const phone = userInfo.value.phone || ''
  if (!phone) return $t.value('unboundPhone')
  const str = String(phone)
  if (str.length < 7) return str
  return `${str.slice(0, 3)}****${str.slice(-4)}`
})

const loadUserInfo = async () => {
  if (!isLogin.value) {
    userInfo.value = {}
    return
  }
  try {
    await userStore.getInfo()
  } catch (e) {
    // ignore
  }
  userInfo.value = uni.getStorageSync('userInfo') || {}
}

const goLogin = () => {
  uni.navigateTo({ url: '/pages/user/login' })
}

const goProfile = () => {
  if (!isLogin.value) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/user/profile' })
}

const openRecharge = () => {
  if (!isLogin.value) {
    uni.showToast({ title: $t.value('pleaseLogin'), icon: 'none' })
    goLogin()
    return
  }
  showRecharge.value = true
}

const selectRecharge = (item) => {
  showRecharge.value = false
  const points = parseRechargePointsValue(item)
  const price = parseRechargePriceValue(item)
  if (!points || !price) {
    uni.showToast({ title: $t.value('apiResponseInvalid'), icon: 'none' })
    return
  }

  const methods = paymentMethods.value.length > 0 ? paymentMethods.value : buildDefaultPaymentMethods()
  uni.showActionSheet({
    itemList: methods.map(m => m.label),
    success: async (sheetRes) => {
      const selected = methods[sheetRes.tapIndex] || methods[0]
      if (!selected?.key) {
        uni.showToast({ title: $t.value('apiResponseInvalid'), icon: 'none' })
        return
      }

      uni.showLoading({ title: $t.value('loading'), mask: true })
      try {
        const res = await createTryonRechargeOrder({
          points,
          price,
          payMethod: selected.key,
        })
        if (res.code !== 0) {
          uni.showToast({ title: resolveApiMessage(res.msg, 'orderCreateFail'), icon: 'none' })
          return
        }

        const order = res?.data?.order || {}
        const orderID = Number(order.ID || order.id || 0)
        if (!orderID) {
          uni.showToast({ title: $t.value('orderCreateFail'), icon: 'none' })
          return
        }
        const orderNo = String(order.outTradeNo || order.OutTradeNo || orderID)

        const amountInCent = Number(order.amount || order.Amount || 0)
        const amount = amountInCent > 0 ? (amountInCent / 100).toFixed(2) : price
        const encodedPayMethod = encodeURIComponent(String(selected.key || 'contact'))
        const encodedPayMethodLabel = encodeURIComponent(String(selected.label || getPaymentMethodLabel(selected.key)))
        const closeTimeRaw = String(order.closeTime || order.CloseTime || '').trim()
        const closeTimePart = closeTimeRaw ? `&closeTime=${encodeURIComponent(closeTimeRaw)}` : ''

        uni.navigateTo({
          url: `/pages/pay/index?orderType=recharge&amount=${encodeURIComponent(amount)}&orderNo=${encodeURIComponent(orderNo)}&orderId=${orderID}&payMethod=${encodedPayMethod}&payMethodLabel=${encodedPayMethodLabel}&rechargePoints=${points}${closeTimePart}`,
        })
      } catch (e) {
        uni.showToast({ title: resolveApiMessage(e?.message, 'orderCreateFail'), icon: 'none' })
      } finally {
        uni.hideLoading()
      }
    },
  })
}

const goTryonRoom = () => {
  uni.switchTab({ url: '/pages/tabBar/index' })
}

const goShoeRoom = () => {
  uni.switchTab({ url: '/pages/tabBar/shop/shop' })
}

const goHistory = () => {
  if (!isLogin.value) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/tryon/history' })
}

const goClothesPage = () => {
  uni.switchTab({ url: '/pages/tabBar/clothes/index' })
}

const goKefu = () => {
  uni.navigateTo({ url: '/pages/kefu/index' })
}

const goOrder = () => {
  if (!isLogin.value) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/order/order' })
}

const goTryonPointRecord = () => {
  if (!isLogin.value) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/integral/integral?assetType=tryon_point&mode=tryon' })
}

const goTryonRechargeRecord = () => {
  if (!isLogin.value) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/rechargeRecord/index' })
}

const goInvite = () => {
  if (!isLogin.value) {
    goLogin()
    return
  }
  uni.navigateTo({ url: '/pages/invite/index' })
}

const goCollect = () => {
  uni.navigateTo({ url: '/pages/collect/collect' })
}

const goMyModel = () => {
  uni.navigateTo({ url: '/pages/myModel/index' })
}

const goMyCloth = () => {
  uni.navigateTo({ url: '/pages/myCloth/index' })
}

const goUseRecord = () => {
  uni.navigateTo({ url: '/pages/browseHistory/index' })
}

const openLangSwitch = async () => {
  await langStore.initLangs()
  showLangPicker.value = true
}

const logoutDevice = () => {
  uni.showModal({
    title: $t.value('pendingOrderTitle'),
    content: $t.value('confirmLogoutDevice'),
    cancelText: $t.value('cancel'),
    confirmText: $t.value('confirm'),
    success: (res) => {
      if (!res.confirm) return
      userStore.loginOut()
      userInfo.value = {}
      uni.showToast({ title: $t.value('logoutSuccess'), icon: 'none' })
    },
  })
}

onShow(() => {
  langStore.initLangs()
  appConfigStore.loadConfig()
  loadRechargePlans()
  loadPaymentMethods()
  loadUserInfo()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.my-page {
  min-height: 100vh;
  color: #0f172a;
  padding: calc(var(--status-bar-height, 0px) + 14rpx) 20rpx 24rpx;
  position: relative;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.bg-glow {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 520rpx;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 35% -10%, rgba(37, 99, 235, 0.2) 0%, transparent 58%),
    radial-gradient(ellipse at 75% 20%, rgba(14, 165, 233, 0.12) 0%, transparent 46%);
}

.profile-card,
.section {
  position: relative;
  z-index: 1;
  border-radius: 18rpx;
  background: rgba(255,255,255,0.93);
  border: 1rpx solid rgba(15,23,42,0.08);
  box-shadow: 0 14rpx 30rpx rgba(15, 23, 42, 0.06);
}

.profile-card {
  padding: 18rpx;
}

.profile-main,
.guest-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.avatar {
  width: 92rpx;
  height: 92rpx;
  border-radius: 50%;
  margin-right: 16rpx;
}

.profile-main {
  gap: 14rpx;
}

.profile-info {
  flex: 1;
  min-width: 0;
}

.name {
  display: block;
  font-size: 28rpx;
  font-weight: 700;
}

.phone {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15,23,42,0.55);
}

.guest-left {
  display: flex;
  align-items: center;
  gap: 14rpx;
}

.guest-icon {
  width: 84rpx;
  height: 84rpx;
  border-radius: 50%;
  background: rgba(219,234,254,0.95);
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 38rpx;
}

.guest-title {
  display: block;
  font-size: 28rpx;
  font-weight: 700;
}

.guest-sub {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15,23,42,0.55);
}

.point-card {
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid rgba(15,23,42,0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.point-label {
  display: block;
  font-size: 22rpx;
  color: rgba(15,23,42,0.55);
}

.point-value {
  display: block;
  margin-top: 6rpx;
  font-size: 36rpx;
  font-weight: 700;
}

.recharge-btn {
  width: 160rpx;
  height: 64rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  font-weight: 600;
  box-shadow: 0 10rpx 24rpx rgba(37, 99, 235, 0.26);
}

.section {
  margin-top: 16rpx;
  padding: 14rpx;
}

.section-title {
  font-size: 26rpx;
  font-weight: 700;
  margin-bottom: 12rpx;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 8rpx;
}

.quick-item {
  height: 132rpx;
  border-radius: 12rpx;
  background: rgba(219,234,254,0.55);
  color: #0f172a;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  font-size: 20rpx;
}

.menu-list {
  background: rgba(248,250,252,0.95);
  border-radius: 12rpx;
  overflow: hidden;
}

.menu-item {
  height: 84rpx;
  padding: 0 16rpx;
  border-bottom: 1rpx solid rgba(15,23,42,0.06);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-text {
  font-size: 24rpx;
}

.menu-text.danger {
  color: #dc2626;
}

.popup-mask {
  position: fixed;
  inset: 0;
  z-index: 99;
  background: rgba(15,23,42,0.36);
  display: flex;
  align-items: flex-end;
}

.popup-panel {
  width: 100%;
  border-top-left-radius: 24rpx;
  border-top-right-radius: 24rpx;
  background: #ffffff;
  padding: 24rpx;
}

.popup-title {
  text-align: center;
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 14rpx;
}

.popup-list {
  max-height: 560rpx;
}

.popup-item {
  margin-bottom: 12rpx;
  height: 86rpx;
  border-radius: 12rpx;
  padding: 0 14rpx;
  background: rgba(248,250,252,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.item-title {
  display: block;
  font-size: 25rpx;
  color: #0f172a;
}

.item-sub {
  display: block;
  margin-top: 6rpx;
  font-size: 21rpx;
  color: rgba(15,23,42,0.55);
}

.popup-close {
  margin-top: 14rpx;
  height: 78rpx;
  border-radius: 12rpx;
  background: rgba(219,234,254,0.65);
  color: #0f172a;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
}
</style>
