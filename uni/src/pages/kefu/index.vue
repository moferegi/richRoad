<template>
  <view class="nf-kefu">
    <view class="nf-kefu-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#0f172a"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('kefuTitle') }}</text>
        <view class="nf-lang-btn" @tap="showLangPicker = true">
          <text class="nf-lang-label">{{ langLabel }}</text>
        </view>
      </view>
    </view>

    <scroll-view scroll-y :show-scrollbar="false" class="nf-kefu-scroll">
      <!-- 支付上下文提示（从订单页跳入时展示） -->
      <view v-if="hasPaymentContext" class="nf-payment-hint-wrap">
        <view class="nf-payment-hint-card">
          <view class="nf-payment-hint-head">
            <text class="nf-payment-hint-title">{{ $t('kefuPaymentHintTitle') }}</text>
            <view class="nf-payment-copy-btn" @tap.stop="copyPaymentDraft">
              <text class="nf-payment-copy-btn-text">{{ $t('kefuPaymentCopyDraft') }}</text>
            </view>
          </view>

          <view class="nf-payment-meta-row">
            <text class="nf-payment-meta-label">{{ $t('kefuPaymentOrderNo') }}</text>
            <text class="nf-payment-meta-value">{{ prefillOrderID }}</text>
          </view>
          <view class="nf-payment-meta-row">
            <text class="nf-payment-meta-label">{{ $t('kefuPaymentMethod') }}</text>
            <text class="nf-payment-meta-value">{{ expectedPayMethodLabel }}</text>
          </view>
          <view class="nf-payment-meta-row" v-if="prefillSelectedQrName">
            <text class="nf-payment-meta-label">{{ $t('kefuPaymentQrLabel') }}</text>
            <text class="nf-payment-meta-value">{{ prefillSelectedQrName }}</text>
          </view>

          <view class="nf-payment-draft-wrap">
            <text class="nf-payment-meta-label">{{ $t('kefuPaymentDraftLabel') }}</text>
            <text class="nf-payment-draft-text">{{ paymentDraftText }}</text>
          </view>
        </view>
      </view>

      <!-- 外部客服列表（shop/kefu 开关控制） -->
      <view v-if="showExternalList && kefuList.length" class="nf-kefu-list">
        <view
          class="nf-kefu-card"
          v-for="(item, index) in kefuList"
          :key="index"
          @tap="contactKefu(item)"
        >
          <!-- 头像区 -->
          <view class="nf-kefu-avatar-wrap">
            <image v-if="resolveAvatar(item)" class="nf-kefu-avatar" :src="resolveAvatar(item)" mode="aspectFill" />
            <view v-else class="nf-kefu-avatar nf-kefu-avatar-fallback" :style="{ background: avatarColor(displayKefuName(item, index)) }">
              <text class="nf-kefu-avatar-fallback-text">{{ avatarInitial(displayKefuName(item, index)) }}</text>
            </view>
            <view class="nf-kefu-status-dot" :class="'nf-dot-' + normalizeStatus(item.status)"></view>
          </view>
          <!-- 信息区 -->
          <view class="nf-kefu-info">
            <text class="nf-kefu-name">{{ displayKefuName(item, index) }}</text>
            <view class="nf-kefu-status-row">
              <text class="nf-kefu-status-text" :class="'nf-status-' + normalizeStatus(item.status)">
                {{ getStatusText(item.status) }}
              </text>
            </view>
          </view>
          <!-- 联系按钮 -->
          <view class="nf-kefu-action" :class="{ 'nf-action-disabled': normalizeStatus(item.status) === 'offline' }">
            <text class="nf-kefu-action-text">{{ $t('kefuContact') }}</text>
          </view>
        </view>
      </view>

      <!-- 平台内置客服入口（customerService/config 开关控制；若都开，排在最后） -->
      <view v-if="showPlatEntry" :class="['nf-plat-cs-wrap', showExternalList ? 'nf-plat-cs-wrap--tail' : '']">
        <view v-if="showExternalList" class="nf-section-title">{{ $t('kefuPlatformSection') }}</view>
        <view class="nf-plat-cs-card" @tap="enterPlatChat">
          <view class="nf-plat-cs-avatar">
            <uni-icons type="chat-filled" size="40" color="#2563eb" />
          </view>
          <view class="nf-plat-cs-info">
            <text class="nf-plat-cs-name">{{ $t('kefuPlatformName') }}</text>
            <text class="nf-plat-cs-desc">{{ $t('kefuPlatformDesc') }}</text>
          </view>
          <view class="nf-kefu-action">
            <text class="nf-kefu-action-text">{{ $t('kefuConsultNow') }}</text>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="nf-kefu-empty" v-if="showNothing || showExternalEmpty">
        <view class="nf-kefu-empty-icon">
          <uni-icons type="chat" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-kefu-empty-text">{{ $t('kefuEmpty') }}</text>
      </view>
    </scroll-view>

    <lang-switch v-model="showLangPicker" />
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { getKefuList, getCsConfig, getSysConfigByKey } from '@/api/kefu.js'
import { getUrl } from '@/utils/url.js'
import { localText, t as i18nT } from '@/utils/i18n.js'
import { trackVisitorEvent } from '@/utils/visitorEvent.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import langSwitch from '@/components/lang-switch/lang-switch.vue'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const showLangPicker = ref(false)
const langLabel = computed(() => {
  const map = { zh: 'ZH', en: 'EN', mn: 'MN', 'zh-TW': 'TW', th: 'TH', hi: 'HI', id: 'ID', vi: 'VI', ar: 'AR', ja: 'JA', ko: 'KO', ms: 'MS' }
  return map[langStore.locale] || String(langStore.locale || 'zh').slice(0, 2).toUpperCase()
})

const kefuList = ref([])
const isLoading = ref(true)
const platEnabled = ref(false)
const externalEnabled = ref(true)
const csDefaultAvatarUrl = ref('')

const avatarBgPalette = ['#E50914', '#0EA5E9', '#10B981', '#F59E0B', '#6366F1', '#EC4899', '#14B8A6', '#F97316']

const showExternalList = computed(() => externalEnabled.value)
const showPlatEntry = computed(() => platEnabled.value)
const showNothing = computed(() => !showExternalList.value && !showPlatEntry.value)
const showExternalEmpty = computed(() => showExternalList.value && !kefuList.value.length && !isLoading.value && !showPlatEntry.value)

const KEFU_CHAT_PREFILL_KEY = 'kefu:chat:prefill'
const prefillOrderID = ref('')
const prefillPayMethod = ref('')
const prefillPayMethodLabel = ref('')
const prefillPreferredPayMethod = ref('')
const prefillPreferredPayMethodLabel = ref('')
const prefillPreferredPayMethodCopyText = ref('')
const prefillSelectedQrName = ref('')

const resolveLocalizedText = (value) => {
  const resolved = localText(value, langStore.locale)
  return String(resolved || value || '').trim()
}

const hasPaymentContext = computed(() => !!prefillOrderID.value)
const expectedPayMethodLabel = computed(() => {
  const preferredKey = String(prefillPreferredPayMethod.value || '').trim().toLowerCase()
  const preferredLabel = resolveLocalizedText(prefillPreferredPayMethodLabel.value)
  if (preferredKey && preferredKey !== 'contact' && preferredKey !== 'qrcode') {
    return preferredLabel || getPayMethodLabel(preferredKey)
  }

  const payKey = String(prefillPayMethod.value || '').trim().toLowerCase()
  const payLabel = resolveLocalizedText(prefillPayMethodLabel.value)
  if (payKey === 'contact') {
    return $t.value('paymentPreferredMethodMissing')
  }

  return payLabel || getPayMethodLabel(payKey)
})

const extractLinkHost = (link) => {
  const target = String(link || '').trim()
  if (!target) return ''
  try {
    return new URL(target).host
  } catch (e) {
    return ''
  }
}

const trackKefuGuideEvent = (action, extra = {}) => {
  const payload = {
    source: 'kefu_index',
    orderNo: prefillOrderID.value || '',
    payMethod: prefillPayMethod.value || '',
    payMethodLabel: prefillPayMethodLabel.value || '',
    preferredPayMethod: prefillPreferredPayMethod.value || '',
    preferredPayMethodLabel: expectedPayMethodLabel.value || '',
    selectedQrName: prefillSelectedQrName.value || '',
    ...extra
  }
  return trackVisitorEvent({
    action,
    label: String(payload.payMethod || ''),
    extra: payload
  })
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

const paymentDraftText = computed(() => {
  if (!prefillOrderID.value) return ''
  const template = resolveLocalizedText(prefillPreferredPayMethodCopyText.value) || $t.value('kefuPaymentDraftTemplatePending')
  let text = template
    .replace('{orderID}', prefillOrderID.value || '-')
    .replace('{payMethod}', expectedPayMethodLabel.value || '-')
  if (prefillSelectedQrName.value) {
    text += `\n${$t.value('kefuPaymentQrSource').replace('{qrName}', prefillSelectedQrName.value)}`
  }
  return text
})

const safeDecode = (val) => {
  const raw = String(val || '').trim()
  if (!raw) return ''
  try {
    return decodeURIComponent(raw)
  } catch (e) {
    return raw
  }
}

const syncChatPrefill = () => {
  if (!paymentDraftText.value) return
  uni.setStorageSync(KEFU_CHAT_PREFILL_KEY, paymentDraftText.value)
}

const copyPaymentDraft = () => {
  if (!paymentDraftText.value) return
  uni.setClipboardData({
    data: paymentDraftText.value,
    success: () => {
      trackKefuGuideEvent('copy_draft')
      uni.showToast({ title: $t.value('kefuPaymentDraftCopied'), icon: 'none' })
    },
    fail: () => {
      trackKefuGuideEvent('copy_draft_fail')
    },
  })
}

const openExternalLink = (link, extra = {}) => {
  trackKefuGuideEvent('open_external_link', {
    linkHost: extractLinkHost(link),
    ...extra
  })
  // #ifdef H5
  window.open(link)
  // #endif
  // #ifndef H5
  uni.navigateTo({
    url: `/pages/webview/index?url=${encodeURIComponent(link)}`
  })
  // #endif
}

const copyDraftThenOpenExternal = (link) => {
  if (!paymentDraftText.value) {
    openExternalLink(link, { openMode: 'open_only' })
    return
  }
  uni.setClipboardData({
    data: paymentDraftText.value,
    success: () => {
      trackKefuGuideEvent('copy_external_draft')
      uni.showToast({ title: $t.value('kefuPaymentDraftCopied'), icon: 'none' })
      setTimeout(() => openExternalLink(link, { openMode: 'copy_then_open' }), 300)
    },
    fail: () => {
      trackKefuGuideEvent('copy_external_draft_fail')
      openExternalLink(link, { openMode: 'open_after_copy_fail' })
    },
  })
}

const STATUS_KEY_TO_I18N = {
  online: 'kefuOnline',
  offline: 'kefuOffline',
  busy: 'kefuBusy',
}
const STATUS_ALIAS_LANGS = ['zh', 'zh-TW', 'en', 'mn', 'th', 'hi', 'id', 'vi', 'ar', 'ja', 'ko', 'ms']

const buildStatusAliasMap = () => {
  const aliasMap = {
    online: 'online',
    offline: 'offline',
    busy: 'busy',
  }

  Object.keys(STATUS_KEY_TO_I18N).forEach((status) => {
    const i18nKey = STATUS_KEY_TO_I18N[status]
    STATUS_ALIAS_LANGS.forEach((lang) => {
      const text = String(i18nT(i18nKey, lang) || '').trim()
      if (!text) return
      aliasMap[text] = status
      aliasMap[text.replace(/\s+/g, '')] = status
    })
  })

  return aliasMap
}

const statusAliasMap = buildStatusAliasMap()

const normalizeStatus = (status) => {
  const rawStatus = String(status || '').trim()
  const normalizedStatus = rawStatus.replace(/\s+/g, '')
  return statusAliasMap[rawStatus] || statusAliasMap[normalizedStatus] || 'offline'
}

const avatarInitial = (name) => {
  const text = String(name || '').trim()
  if (!text) return 'K'
  return text.charAt(0).toUpperCase()
}

const displayKefuName = (item, index = 0) => {
  const localized = String(localText(item?.nameI18n, langStore.locale) || '').trim()
  if (localized) return localized

  const fallback = String(item?.name || '').trim()
  if (fallback) return fallback

  const order = Number(index) + 1
  return `${$t.value('kefuPlatformName') || 'Kefu'} #${order}`
}

const resolveContactId = (item) => {
  return String(item?.contactId || '').trim()
}

const copyContactId = (contactId) => {
  const value = String(contactId || '').trim()
  if (!value) {
    uni.showToast({ title: $t.value('kefuContactIdMissing'), icon: 'none' })
    return
  }
  uni.setClipboardData({
    data: value,
    success: () => {
      uni.showToast({ title: $t.value('copySuccess'), icon: 'none' })
    },
    fail: () => {
      uni.showToast({ title: $t.value('operationFailed'), icon: 'none' })
    },
  })
}

const avatarColor = (seed) => {
  const text = String(seed || 'kefu')
  let hash = 0
  for (let i = 0; i < text.length; i++) {
    hash = (hash * 31 + text.charCodeAt(i)) >>> 0
  }
  return avatarBgPalette[hash % avatarBgPalette.length]
}

const resolveAvatar = (item) => {
  const raw = String(item?.avatar || item?.externalAvatar || csDefaultAvatarUrl.value || '').trim()
  if (!raw) return ''
  if (/^(https?:)?\/\//i.test(raw) || /^data:/i.test(raw)) return raw
  return getUrl(raw)
}

const resolveQrCode = (item) => {
  const raw = String(item?.qrCode || '').trim()
  if (!raw) return ''
  if (/^(https?:)?\/\//i.test(raw) || /^data:/i.test(raw)) return raw
  return getUrl(raw)
}

const openQrCodePreview = (imageUrl, item = {}) => {
  if (!imageUrl) return
  trackKefuGuideEvent('open_qrcode_preview', {
    linkHost: extractLinkHost(item?.link || ''),
    kefuName: String(item?.name || ''),
  })
  uni.previewImage({
    current: imageUrl,
    urls: [imageUrl],
  })
}

const getStatusText = (status) => {
  const key = normalizeStatus(status)
  const map = {
    'online': $t.value('kefuOnline'),
    'offline': $t.value('kefuOffline'),
    'busy': $t.value('kefuBusy')
  }
  return map[key] || status
}

const init = async () => {
  isLoading.value = true
  try {
    const [listRes, cfgRes, switchRes] = await Promise.allSettled([
      getKefuList(),
      getCsConfig(),
      getSysConfigByKey('shop_kefu_enabled')
    ])
    if (listRes.status === 'fulfilled' && listRes.value.code === 0) {
      kefuList.value = Array.isArray(listRes.value.data) ? listRes.value.data : (listRes.value.data?.list || [])
    }
    if (cfgRes.status === 'fulfilled' && cfgRes.value.code === 0) {
      platEnabled.value = !!cfgRes.value.data?.platEnabled
      csDefaultAvatarUrl.value = String(cfgRes.value.data?.defaultAvatarUrl || '')
    }
    if (switchRes.status === 'fulfilled' && switchRes.value.code === 0) {
      externalEnabled.value = String(switchRes.value.data).toLowerCase() === 'true'
    }
  } catch (e) {
    console.error('初始化客服页失败', e)
  } finally {
    isLoading.value = false
  }
}

onShow(() => { init() })

onLoad((options = {}) => {
  prefillOrderID.value = safeDecode(options.orderID)
  prefillPayMethod.value = safeDecode(options.payMethod)
  prefillPayMethodLabel.value = safeDecode(options.payMethodLabel) || getPayMethodLabel(prefillPayMethod.value)
  prefillPreferredPayMethod.value = safeDecode(options.preferredPayMethod)
  prefillPreferredPayMethodLabel.value = safeDecode(options.preferredPayMethodLabel) || getPayMethodLabel(prefillPreferredPayMethod.value)
  prefillPreferredPayMethodCopyText.value = safeDecode(options.preferredPayMethodCopyText)
  prefillSelectedQrName.value = safeDecode(options.selectedQrName)
  syncChatPrefill()
  if (hasPaymentContext.value) {
    trackKefuGuideEvent('kefu_page_view')
  }
})

// 进入平台客服聊天室
const enterPlatChat = () => {
  if (hasPaymentContext.value) {
    trackKefuGuideEvent('navigate_platform_chat')
    syncChatPrefill()
    const orderID = encodeURIComponent(prefillOrderID.value)
    const payMethod = encodeURIComponent(prefillPayMethod.value)
    const payMethodLabel = encodeURIComponent(prefillPayMethodLabel.value)
    uni.navigateTo({ url: `/pages/kefu/chat?orderID=${orderID}&payMethod=${payMethod}&payMethodLabel=${payMethodLabel}` })
    return
  }
  uni.navigateTo({ url: '/pages/kefu/chat' })
}

const contactKefu = (item) => {
  if (normalizeStatus(item.status) === 'offline') {
    uni.showToast({ title: $t.value('kefuOffline'), icon: 'none' })
    return
  }
  const targetLink = String(item?.link || '').trim()
  if (targetLink) {
    if (hasPaymentContext.value && paymentDraftText.value) {
      trackKefuGuideEvent('external_prompt_show', { linkHost: extractLinkHost(targetLink) })
      uni.showModal({
        title: $t.value('kefuPaymentExternalTitle'),
        content: $t.value('kefuPaymentExternalHint'),
        confirmText: $t.value('kefuPaymentExternalCopyOpen'),
        cancelText: $t.value('kefuPaymentExternalOpenOnly'),
        success: (res) => {
          if (res.confirm) {
            trackKefuGuideEvent('external_prompt_confirm', { linkHost: extractLinkHost(targetLink) })
            copyDraftThenOpenExternal(targetLink)
            return
          }
          trackKefuGuideEvent('external_prompt_cancel', { linkHost: extractLinkHost(targetLink) })
          openExternalLink(targetLink, { openMode: 'open_only' })
        },
        fail: () => {
          trackKefuGuideEvent('external_prompt_cancel', { linkHost: extractLinkHost(targetLink), fail: true })
          openExternalLink(targetLink, { openMode: 'open_only' })
        },
      })
      return
    }
    openExternalLink(targetLink, { openMode: 'open_only' })
    return
  }

  const qrCodeUrl = resolveQrCode(item)
  if (qrCodeUrl) {
    if (hasPaymentContext.value && paymentDraftText.value) {
      uni.setClipboardData({
        data: paymentDraftText.value,
        success: () => {
          trackKefuGuideEvent('copy_qrcode_draft')
          uni.showToast({ title: $t.value('kefuPaymentDraftCopied'), icon: 'none' })
          setTimeout(() => openQrCodePreview(qrCodeUrl, item), 300)
        },
        fail: () => {
          trackKefuGuideEvent('copy_qrcode_draft_fail')
          openQrCodePreview(qrCodeUrl, item)
        },
      })
      return
    }
    openQrCodePreview(qrCodeUrl, item)
  } else {
    const contactId = resolveContactId(item)
    if (!contactId) {
      uni.showToast({ title: $t.value('kefuContactIdMissing'), icon: 'none' })
      return
    }

    const modalContent = `${displayKefuName(item)}\n${$t.value('kefuContactIdLabel')}${contactId}`
    uni.showModal({
      title: $t.value('kefuContactIdModalTitle'),
      content: modalContent,
      confirmText: $t.value('kefuCopyId'),
      cancelText: $t.value('cancel'),
      success: (res) => {
        if (res.confirm) {
          copyContactId(contactId)
        }
      },
    })
  }
}

const goBack = () => {
  uni.navigateBack()
}
</script>

<style lang="scss">
page {
  background-color: #f4f7fb;
}

.nf-kefu {
  width: 100%;
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f4f7fb;
  position: relative;
}

.nf-kefu-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 500rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 30% -10%, rgba(37, 99, 235, 0.2) 0%, transparent 58%),
    radial-gradient(ellipse at 75% 10%, rgba(14, 165, 233, 0.15) 0%, transparent 52%);
}

/* ===== 导航栏 ===== */
.nf-navbar {
  background: rgba(244, 247, 251, 0.92);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.2);
  padding: 0 28rpx 16rpx;
  position: sticky;
  top: 0;
  z-index: 99;
}

.nf-navbar-status {
  width: 100%;
  height: var(--status-bar-height, 0px);
}

/* #ifdef MP-WEIXIN */
.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}
/* #endif */

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
  background: #ffffff;
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  box-shadow: 0 8rpx 22rpx rgba(15, 23, 42, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    background: rgba(241, 245, 249, 0.98);
    transform: scale(0.93);
  }
}

.nf-lang-btn {
  width: 64rpx;
  height: 64rpx;
  border-radius: 18rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.24);
  background: rgba(255, 255, 255, 0.86);
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    transform: scale(0.93);
  }
}

.nf-lang-label {
  font-size: 22rpx;
  font-weight: 800;
  color: #7c2d12;
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: 2rpx;
}

/* ===== 滚动区域 ===== */
.nf-kefu-scroll {
  flex: 1;
  width: 100%;
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
    width: 0;
    height: 0;
  }
}

/* ===== 客服列表 ===== */
.nf-kefu-list {
  padding: 24rpx 28rpx;
}

.nf-payment-hint-wrap {
  padding: 24rpx 28rpx 0;
}

.nf-payment-hint-card {
  background: rgba(219, 234, 254, 0.75);
  border: 1rpx solid rgba(37, 99, 235, 0.2);
  border-radius: 20rpx;
  padding: 20rpx;
}

.nf-payment-hint-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.nf-payment-hint-title {
  font-size: 26rpx;
  color: #0f172a;
  font-weight: 600;
}

.nf-payment-copy-btn {
  padding: 8rpx 16rpx;
  border-radius: 999rpx;
  background: rgba(37, 99, 235, 0.14);
}

.nf-payment-copy-btn-text {
  color: #1d4ed8;
  font-size: 22rpx;
}

.nf-payment-meta-row {
  display: flex;
  align-items: center;
  gap: 8rpx;
  margin-bottom: 8rpx;
}

.nf-payment-meta-label {
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.nf-payment-meta-value {
  font-size: 22rpx;
  color: #0f172a;
}

.nf-payment-draft-wrap {
  margin-top: 10rpx;
  padding-top: 10rpx;
  border-top: 1rpx solid rgba(15, 23, 42, 0.08);
}

.nf-payment-draft-text {
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.72);
  line-height: 1.45;
}

.nf-section-title {
  color: rgba(15, 23, 42, 0.6);
  font-size: 24rpx;
  letter-spacing: 2rpx;
  margin-bottom: 16rpx;
}

/* ===== 平台专属客服卡片 ===== */
.nf-plat-cs-wrap {
  padding: 40rpx 28rpx;
}

.nf-plat-cs-wrap--tail {
  padding-top: 8rpx;
}

.nf-plat-cs-card {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.98);
  border: 1rpx solid rgba(37, 99, 235, 0.16);
  border-radius: 24rpx;
  padding: 40rpx 32rpx;
  box-shadow: 0 14rpx 28rpx rgba(37, 99, 235, 0.1);
  transition: transform 0.3s, background 0.3s;

  &:active {
    transform: scale(0.985);
    background: rgba(239, 246, 255, 0.98);
  }
}

.nf-plat-cs-avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background: rgba(37, 99, 235, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-right: 28rpx;
}

.nf-plat-cs-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.nf-plat-cs-name {
  font-size: 34rpx;
  font-weight: 600;
  color: #0f172a;
}

.nf-plat-cs-desc {
  font-size: 26rpx;
  color: rgba(15, 23, 42, 0.58);
}

.nf-kefu-card {
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.98);
  border: 1rpx solid rgba(148, 163, 184, 0.2);
  border-radius: 24rpx;
  padding: 32rpx 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 12rpx 24rpx rgba(15, 23, 42, 0.07);
  transition: transform 0.3s, background 0.3s;

  &:active {
    transform: scale(0.985);
    background: #f8fbff;
  }
}

/* ===== 头像 ===== */
.nf-kefu-avatar-wrap {
  position: relative;
  flex-shrink: 0;
  margin-right: 24rpx;
}

.nf-kefu-avatar {
  width: 108rpx;
  height: 108rpx;
  border-radius: 50%;
  border: 3rpx solid rgba(37, 99, 235, 0.28);
  box-shadow: 0 10rpx 18rpx rgba(37, 99, 235, 0.12);
}

.nf-kefu-avatar-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-kefu-avatar-fallback-text {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
}

.nf-kefu-status-dot {
  position: absolute;
  bottom: 4rpx;
  right: 4rpx;
  width: 24rpx;
  height: 24rpx;
  border-radius: 50%;
  border: 3rpx solid #ffffff;
}

.nf-dot-online {
  background: #22c55e;
  box-shadow: 0 0 8rpx rgba(34, 197, 94, 0.6);
}

.nf-dot-offline {
  background: #6b7280;
}

.nf-dot-busy {
  background: #f59e0b;
  box-shadow: 0 0 8rpx rgba(245, 158, 11, 0.6);
}

/* ===== 信息区 ===== */
.nf-kefu-info {
  flex: 1;
  min-width: 0;
}

.nf-kefu-name {
  font-size: 32rpx;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: 1rpx;
  margin-bottom: 8rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.nf-kefu-status-row {
  display: flex;
  align-items: center;
}

.nf-kefu-status-text {
  font-size: 24rpx;
  font-weight: 500;
}

.nf-status-online {
  color: #22c55e;
}

.nf-status-offline {
  color: #94a3b8;
}

.nf-status-busy {
  color: #f59e0b;
}

/* ===== 按钮 ===== */
.nf-kefu-action {
  flex-shrink: 0;
  margin-left: 16rpx;
  padding: 14rpx 32rpx;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
  border-radius: 40rpx;
  box-shadow: 0 8rpx 18rpx rgba(37, 99, 235, 0.3);
  transition: all 0.3s;

  &:active {
    transform: scale(0.95);
    box-shadow: 0 2rpx 8rpx rgba(37, 99, 235, 0.45);
  }
}

.nf-action-disabled {
  background: rgba(148, 163, 184, 0.18);
  box-shadow: none;

  .nf-kefu-action-text {
    color: rgba(15, 23, 42, 0.35);
  }
}

.nf-kefu-action-text {
  font-size: 26rpx;
  font-weight: 600;
  color: #fff;
  letter-spacing: 1rpx;
}

/* ===== 空状态 ===== */
.nf-kefu-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.nf-kefu-empty-icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(37, 99, 235, 0.1);
  border: 1rpx solid rgba(37, 99, 235, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
}

.nf-kefu-empty-text {
  font-size: 28rpx;
  color: rgba(15, 23, 42, 0.48);
  letter-spacing: 2rpx;
}

.nf-kefu {
  background: linear-gradient(180deg, #f8f4e7 0%, #f4efe1 100%);
}

.nf-navbar {
  background: rgba(255, 253, 248, 0.88);
  border-bottom: 1rpx solid rgba(146, 64, 14, 0.12);
}

.nf-navbar-back {
  border-radius: 18rpx;
  border: 1rpx solid rgba(20, 184, 166, 0.22);
  background: #fffdf8;
}

.nf-navbar-title {
  color: #7c2d12;
  font-weight: 800;
}

.nf-payment-hint-card,
.nf-plat-cs-card,
.nf-kefu-card {
  background: rgba(255, 253, 248, 0.96);
  border: 1rpx solid rgba(20, 184, 166, 0.18);
  box-shadow: 0 12rpx 24rpx rgba(120, 53, 15, 0.1);
}

.nf-payment-copy-btn,
.nf-plat-cs-avatar {
  background: rgba(15, 118, 110, 0.14);
}

.nf-payment-copy-btn-text {
  color: #0f766e;
}

.nf-kefu-name,
.nf-plat-cs-name,
.nf-payment-hint-title {
  color: #1e293b;
}

.nf-kefu-action {
  background: linear-gradient(120deg, #0f766e 0%, #f97316 100%);
  box-shadow: 0 8rpx 18rpx rgba(15, 118, 110, 0.28);
}

.nf-kefu-empty-icon {
  background: rgba(15, 118, 110, 0.1);
  border: 1rpx solid rgba(15, 118, 110, 0.24);
}
</style>
