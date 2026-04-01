<template>
  <view class="nf-pay">
    <view class="nf-pay-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff" />
        </view>
        <text class="nf-navbar-title">{{ $t('payTitle') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-pay-content">
      <!-- 订单金额 -->
      <view class="nf-pay-amount-card">
        <text class="nf-pay-amount-label">{{ $t('payAmount') }}</text>
        <text class="nf-pay-amount-value">¥{{ amount }}</text>
        <text class="nf-pay-order-no" v-if="orderNo">{{ $t('orderNo') }}: {{ orderNo }}</text>
      </view>

      <!-- 支付二维码 -->
      <view class="nf-pay-qr-card">
        <text class="nf-pay-qr-title">{{ $t('payQrTitle') }}</text>
        <!-- 多码切换 -->
        <view class="nf-pay-qr-tabs" v-if="qrList.length > 1">
          <view
            v-for="(qr, i) in qrList"
            :key="i"
            class="nf-pay-qr-tab"
            :class="{ active: currentQrIndex === i }"
            @tap="currentQrIndex = i"
          >{{ qr.name }}</view>
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
            <uni-icons type="scan" size="48" color="rgba(255,255,255,0.3)" />
          </view>
        </view>
        <text class="nf-pay-qr-tip">{{ $t('payQrTip') }}</text>
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
        <view class="nf-pay-btn nf-pay-btn-primary" @tap="saveQr">
          <uni-icons type="download" size="18" color="#fff" />
          <text>{{ $t('saveQrCode') }}</text>
        </view>
        <view class="nf-pay-btn nf-pay-btn-kefu" @tap="goKefu">
          <uni-icons type="chatbubble" size="18" color="#e50914" />
          <text>{{ $t('contactCustomerService') }}</text>
        </view>
      </view>

      <!-- 已支付确认 -->
      <view class="nf-pay-confirm" @tap="confirmPaid">
        <text>{{ $t('alreadyPaid') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { request } from '@/utils/request.js'
import { getUrl } from '@/utils/url.js'
import { getEnabledQrcodePayments } from '@/api/qrcodePayment.js'
import { localText } from '@/utils/i18n'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const amount = ref('0.00')
const orderNo = ref('')
const orderId = ref('')

// 多码支持
const qrList = ref([])
const currentQrIndex = ref(0)
const currentQrUrl = computed(() => {
  const qr = qrList.value[currentQrIndex.value]
  if (!qr) return ''
  return getUrl(qr.externalPath || qr.image)
})

// 付款提示文本
const paymentTipText = ref('')
const paymentTipSize = ref(14)
const paymentTipColor = ref('#ffffff')

const steps = computed(() => [
  $t.value('payStep1'),
  $t.value('payStep2'),
  $t.value('payStep3'),
])

onLoad((options) => {
  if (options.amount) amount.value = options.amount
  if (options.orderNo) orderNo.value = options.orderNo
  if (options.orderId) orderId.value = options.orderId
  loadQrCodes()
  loadPaymentTip()
})

const loadQrCodes = async () => {
  try {
    const res = await getEnabledQrcodePayments()
    if (res.code === 0 && res.data && res.data.list) {
      qrList.value = res.data.list
    }
    // 如果多码列表为空，尝试旧的单码配置作为兜底
    if (qrList.value.length === 0) {
      const res2 = await request({ url: '/sysConfig/getConfigByKey', method: 'get', params: { key: 'payment_qr_code' } })
      if (res2.code === 0 && res2.data) {
        qrList.value = [{ name: 'QR', image: res2.data, externalPath: '' }]
      }
    }
  } catch (e) {
    console.error('加载支付二维码失败', e)
  }
}

const loadPaymentTip = async () => {
  try {
    const [tipRes, sizeRes, colorRes] = await Promise.all([
      request({ url: '/sysConfig/getConfigByKey', method: 'get', params: { key: 'payment_tip_text' } }),
      request({ url: '/sysConfig/getConfigByKey', method: 'get', params: { key: 'payment_tip_text_size' } }),
      request({ url: '/sysConfig/getConfigByKey', method: 'get', params: { key: 'payment_tip_text_color' } }),
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
  uni.navigateTo({ url: '/pages/kefu/index' })
}

const confirmPaid = () => {
  uni.showToast({ title: $t.value('paidSuccess'), icon: 'success' })
  setTimeout(() => {
    uni.redirectTo({ url: '/pages/order/order' })
  }, 1500)
}

const goBack = () => {
  uni.navigateBack()
}
</script>

<style lang="scss" scoped>
.nf-pay {
  min-height: 100vh;
  background: #141414;
  position: relative;
}
.nf-pay-bg {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: linear-gradient(180deg, #1a1a2e 0%, #141414 100%);
  z-index: 0;
}
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(20rpx);
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
  background: rgba(255, 255, 255, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-navbar-title {
  font-size: 34rpx;
  font-weight: bold;
  color: #fff;
}
.nf-pay-content {
  position: relative;
  z-index: 1;
  padding: calc(var(--status-bar-height, 44rpx) + 88rpx + 30rpx) 24rpx 60rpx;
}
.nf-pay-amount-card {
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 40rpx;
  text-align: center;
  margin-bottom: 24rpx;
}
.nf-pay-amount-label {
  font-size: 26rpx;
  color: rgba(255,255,255,0.5);
  display: block;
  margin-bottom: 12rpx;
}
.nf-pay-amount-value {
  font-size: 56rpx;
  font-weight: bold;
  color: #e50914;
  display: block;
}
.nf-pay-order-no {
  font-size: 22rpx;
  color: rgba(255,255,255,0.3);
  margin-top: 12rpx;
  display: block;
}
.nf-pay-qr-card {
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 30rpx;
  text-align: center;
  margin-bottom: 24rpx;
}
.nf-pay-qr-title {
  font-size: 28rpx;
  font-weight: bold;
  color: #fff;
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
  background: rgba(255,255,255,0.05);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-pay-qr-tip {
  font-size: 24rpx;
  color: rgba(255,255,255,0.4);
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
  color: rgba(255,255,255,0.5);
  background: rgba(255,255,255,0.08);
  border: 1rpx solid transparent;
}
.nf-pay-qr-tab.active {
  color: #e50914;
  background: rgba(229, 9, 20, 0.15);
  border-color: rgba(229, 9, 20, 0.3);
}
.nf-pay-tip-card {
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 24rpx 30rpx;
  margin-bottom: 24rpx;
  text-align: center;
}
.nf-pay-tip-text {
  line-height: 1.6;
}
.nf-pay-steps {
  background: rgba(255,255,255,0.06);
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 30rpx;
}
.nf-pay-steps-title {
  font-size: 28rpx;
  font-weight: bold;
  color: #fff;
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
  background: #e50914;
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
  color: rgba(255,255,255,0.7);
}
.nf-pay-actions {
  display: flex;
  gap: 20rpx;
  margin-bottom: 20rpx;
}
.nf-pay-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  font-size: 28rpx;
  font-weight: bold;
}
.nf-pay-btn-primary {
  background: #e50914;
  color: #fff;
}
.nf-pay-btn-kefu {
  background: rgba(229, 9, 20, 0.15);
  border: 1rpx solid rgba(229, 9, 20, 0.3);
  color: #e50914;
}
.nf-pay-confirm {
  text-align: center;
  padding: 20rpx;
}
.nf-pay-confirm text {
  font-size: 28rpx;
  color: rgba(255,255,255,0.5);
  text-decoration: underline;
}
</style>
