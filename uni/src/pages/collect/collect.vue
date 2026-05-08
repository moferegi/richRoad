<template>
  <view class="collect-page">
    <view class="collect-bg"></view>

    <view class="collect-navbar">
      <view class="collect-navbar-status"></view>
      <view class="collect-navbar-content">
        <view class="nav-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#1e293b" />
        </view>
        <text class="collect-navbar-title">{{ $t('myCollection') }}</text>
        <view class="nav-space"></view>
      </view>
    </view>

    <scroll-view
      scroll-y
      :show-scrollbar="false"
      class="collect-scroll"
      @scrolltolower="debouncedLower"
    >
      <!-- <view class="collect-summary" v-if="collectList.length">
        <text class="collect-summary-title">{{ $t('myCollectionMenu') }}</text>
        <text class="collect-summary-count">{{ collectList.length }}</text>
      </view> -->

      <view class="collect-grid" v-if="collectList.length">
        <view
          class="collect-card"
          v-for="(item, index) in collectList"
          :key="index"
          @tap="goTo(item)"
        >
          <view class="card-image-wrap">
            <image class="card-image" :src="getUrl(item.imageUrl)" mode="aspectFill" />
            <view class="card-discount" v-if="item.discount && item.discount < 10">
              <text>{{ getDiscountText(item.discount) }}</text>
            </view>
            <view class="card-heart" @tap.stop="cancelCollect(item.ID)">
              <uni-icons type="heart-filled" size="16" color="#ffffff" />
            </view>
          </view>
          <view class="card-info">
            <text class="card-title">{{ $lt(item.title) }}</text>
            <text class="card-price">{{ cs }}{{ formatPrice(item.price) }}</text>
            <view class="card-meta">
              <text class="card-sold">{{ $t('sold') }} {{ item.saleNum || 0 }}</text>
              <uni-icons type="right" size="14" color="rgba(30,41,59,0.36)" />
            </view>
          </view>
        </view>
      </view>

      <view class="collect-footer" v-if="collectList.length > 0">
        <text class="collect-footer-text">{{ isBottom ? $t('reachedBottom') : $t('loading') }}</text>
      </view>

      <view class="collect-empty" v-if="!collectList.length">
        <view class="collect-empty-icon">
          <uni-icons type="heart" size="50" color="rgba(37,99,235,0.45)" />
        </view>
        <text class="collect-empty-text">{{ $t('noCollectData') }}</text>
      </view>
    </scroll-view>

    <view v-if="showLoginModal" class="login-mask" @tap.stop>
      <view class="login-modal">
        <text class="login-title">{{ $t('loginFirst') }}</text>
        <view class="login-actions">
          <view class="login-btn login-btn-cancel" @tap="onModalCancel">
            <text>{{ $t('cancel') }}</text>
          </view>
          <view class="login-btn login-btn-confirm" @tap="onModalConfirm">
            <text>{{ $t('goLogin') }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getCollectList, findCollect, createCollect } from '@/api/collect'
import { useUserStore } from '@/pinia/modules/user'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol || '¥')
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const userStore = useUserStore()
const collectList = ref([])
const collectionFlag = ref(false)
const showLoginModal = ref(false)

const getToken = () => {
  return userStore.token || uni.getStorageSync('x-token') || ''
}

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  const amount = Number(priceInCents)
  if (Number.isNaN(amount)) return '0.00'
  if (amount >= 1000) return (amount / 100).toFixed(2)
  return amount.toFixed(2)
}

const init = async () => {
  params = { page: 1, pageSize: 10 }
  isBottom.value = false
  const res = await getCollectList(params)
  if (res.code === 0) {
    const list = Array.isArray(res?.data?.list) ? res.data.list : []
    collectList.value = list
    isBottom.value = list.length < params.pageSize
  }
}

onShow(() => {
  const tk = getToken()
  if (!tk) {
    showLoginModal.value = true
    return
  }
  showLoginModal.value = false
  init()
})

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack({ delta: 1 })
    return
  }
  uni.switchTab({ url: '/pages/tabBar/my/index' })
}

const onModalCancel = () => {
  showLoginModal.value = false
  goBack()
}

const onModalConfirm = () => {
  showLoginModal.value = false
  uni.redirectTo({ url: '/pages/user/login' })
}

const debounce = (func, delay) => {
  let timer
  return function (...args) {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => func.apply(this, args), delay)
  }
}

let params = { page: 1, pageSize: 10 }
const isBottom = ref(false)

const lower = async () => {
  if (isBottom.value) return
  params.page += 1
  const res = await getCollectList(params)
  if (res.code !== 0) {
    isBottom.value = true
    return
  }

  const list = Array.isArray(res?.data?.list) ? res.data.list : []
  if (!list.length) {
    isBottom.value = true
    return
  }

  collectList.value.push(...list)
  if (list.length < params.pageSize) {
    isBottom.value = true
  }
}

const debouncedLower = debounce(lower, 300)

const cancelCollect = async (ID) => {
  const token = getToken()
  if (!token) {
    uni.showToast({ title: $t.value('loginFirst'), mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }

  const status = await findCollect({ goodID: ID })
  if (status.code === 0) {
    collectionFlag.value = !!status.data
  }

  const res = await createCollect({ goodID: Number(ID) })
  if (res.code === 0) {
    collectionFlag.value = !collectionFlag.value
    await init()
    uni.showToast({
      title: collectionFlag.value ? $t.value('collectSuccess') : $t.value('uncollectSuccess'),
      mask: true,
      icon: 'none'
    })
  }
}

const goTo = (item) => {
  const goodID = Number(item?.ID || 0)
  if (!goodID) return
  uni.navigateTo({ url: `/pages/goodsDetails/goodsDetails?id=${goodID}` })
}

const getDiscountText = (discount) => {
  const t = $t.value
  if (discount >= 9.5) return t('priceDrop')
  if (discount >= 9.0) return t('discount')
  if (discount >= 8.0) return t('specialOffer')
  if (discount >= 7.0) return t('goodPrice')
  if (discount >= 6.0) return t('lowPrice')
  if (discount >= 5.0) return t('bargain')
  return t('saleTag')
}
</script>

<style lang="scss" scoped>
page {
  background: #f4f7fb;
}

.collect-page {
  min-height: 100vh;
  background: #f4f7fb;
}

.collect-bg {
  position: fixed;
  inset: 0 0 auto 0;
  height: 360rpx;
  pointer-events: none;
  background:
    radial-gradient(120% 90% at 100% -10%, rgba(14, 165, 233, 0.18) 0%, transparent 60%),
    radial-gradient(120% 90% at 0% 0%, rgba(37, 99, 235, 0.18) 0%, transparent 55%);
}

.collect-navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(244, 247, 251, 0.9);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  border-bottom: 1rpx solid rgba(148, 163, 184, 0.2);
}

.collect-navbar-status {
  height: var(--status-bar-height, 44px);
}

.collect-navbar-content {
  height: 88rpx;
  padding: 0 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-back,
.nav-space {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-back {
  border-radius: 50%;
  background: #ffffff;
  box-shadow: 0 8rpx 20rpx rgba(15, 23, 42, 0.08);
}

.collect-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: 1rpx;
}

.collect-scroll {
  min-height: calc(100vh - var(--status-bar-height, 44px) - 88rpx);
  margin-top: calc(var(--status-bar-height, 44px) + 88rpx);
}

.collect-summary {
  margin: 24rpx 24rpx 10rpx;
  padding: 24rpx 28rpx;
  border-radius: 20rpx;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
  color: #ffffff;
  box-shadow: 0 16rpx 28rpx rgba(37, 99, 235, 0.24);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.collect-summary-title {
  font-size: 28rpx;
  font-weight: 600;
}

.collect-summary-count {
  min-width: 54rpx;
  height: 54rpx;
  border-radius: 27rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  font-weight: 700;
  background: rgba(255, 255, 255, 0.2);
}

.collect-grid {
  padding: 16rpx 24rpx 10rpx;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.collect-card {
  width: 338rpx;
  margin-bottom: 18rpx;
  overflow: hidden;
  border-radius: 18rpx;
  background: #ffffff;
  box-shadow: 0 12rpx 24rpx rgba(15, 23, 42, 0.08);
}

.card-image-wrap {
  position: relative;
  width: 100%;
  height: 338rpx;
}

.card-image {
  width: 100%;
  height: 100%;
}

.card-discount {
  position: absolute;
  left: 12rpx;
  top: 12rpx;
  border-radius: 8rpx;
  padding: 4rpx 10rpx;
  background: rgba(239, 68, 68, 0.92);
  color: #ffffff;
  font-size: 20rpx;
  font-weight: 600;
}

.card-heart {
  position: absolute;
  right: 12rpx;
  top: 12rpx;
  width: 48rpx;
  height: 48rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(15, 23, 42, 0.4);
}

.card-info {
  padding: 16rpx 14rpx 14rpx;
}

.card-title {
  font-size: 25rpx;
  line-height: 1.45;
  color: #0f172a;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  min-height: 72rpx;
}

.card-price {
  margin-top: 10rpx;
  display: block;
  font-size: 32rpx;
  color: #ef4444;
  font-weight: 700;
}

.card-meta {
  margin-top: 8rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-sold {
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.55);
}

.collect-footer {
  padding: 20rpx 0 32rpx;
  text-align: center;
}

.collect-footer-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.42);
}

.collect-empty {
  min-height: 68vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.collect-empty-icon {
  width: 130rpx;
  height: 130rpx;
  border-radius: 65rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(37, 99, 235, 0.08);
  border: 1rpx solid rgba(37, 99, 235, 0.2);
}

.collect-empty-text {
  margin-top: 18rpx;
  font-size: 28rpx;
  color: rgba(15, 23, 42, 0.52);
}

.login-mask {
  position: fixed;
  inset: 0;
  z-index: 999;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-modal {
  width: 560rpx;
  padding: 48rpx 38rpx 34rpx;
  border-radius: 20rpx;
  background: #ffffff;
  box-shadow: 0 20rpx 40rpx rgba(15, 23, 42, 0.2);
}

.login-title {
  display: block;
  text-align: center;
  color: #0f172a;
  font-size: 32rpx;
  font-weight: 700;
  margin-bottom: 42rpx;
}

.login-actions {
  display: flex;
  gap: 16rpx;
}

.login-btn {
  flex: 1;
  height: 78rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  font-weight: 600;
}

.login-btn-cancel {
  color: #475569;
  background: #f1f5f9;
}

.login-btn-confirm {
  color: #ffffff;
  background: linear-gradient(135deg, #2563eb, #0ea5e9);
}
</style>
