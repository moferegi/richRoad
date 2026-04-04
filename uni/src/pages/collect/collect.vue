<template>
  <view class="nf-collect">
    <view class="nf-collect-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view style="width: 64rpx;"></view>
        <text class="nf-navbar-title">{{ $t('myCollection') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 商品列表 -->
    <scroll-view
      scroll-y
      :show-scrollbar="false"
      class="nf-collect-scroll"
      @scrolltolower="debouncedLower"
    >
      <view v-if="collectList.length" class="nf-goods-list">
        <view
          class="nf-goods-card"
          v-for="(item, index) in collectList"
          :key="index"
          @tap="goTo(item)"
        >
          <view class="nf-goods-img-wrap">
            <image class="nf-goods-img" :src="getUrl(item.imageUrl)" mode="aspectFill" />
            <view class="nf-goods-img-overlay"></view>
            <view class="nf-goods-badge" v-if="item.discount && item.discount < 10">
              <text class="nf-goods-badge-text">{{ getDiscountText(item.discount) }}</text>
            </view>
          </view>
          <view class="nf-goods-info">
            <text class="nf-goods-title">{{ $lt(item.title) }}</text>
            <view class="nf-goods-tags">
              <text class="nf-tag nf-tag-red">{{ $t('selfOperated') }}</text>
              <text class="nf-tag nf-tag-blue">{{ $t('qualityAssured') }}</text>
              <text class="nf-tag nf-tag-green">{{ $t('freeShipping') }}</text>
            </view>
            <view class="nf-goods-bottom">
              <view class="nf-price-row">
                <text class="nf-price">{{ cs }}{{ formatPrice(item.price) }}</text>
              </view>
              <view class="nf-goods-meta">
                <text class="nf-rating-stars">★</text>
                <text class="nf-rating-score">{{ item.rating || '5.0' }}</text>
                <text class="nf-sold">{{ $t('sold') }} {{ item.saleNum || 0 }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <!-- 底部加载状态 -->
      <view class="nf-collect-footer" v-if="collectList.length > 0">
        <text class="nf-collect-footer-text">{{ isBottom ? $t('reachedBottom') : $t('loading') }}</text>
      </view>

      <!-- 空状态 -->
      <view class="nf-collect-empty" v-if="!collectList.length">
        <view class="nf-collect-empty-icon">
          <uni-icons type="heart" size="48" color="rgba(229,9,20,0.4)" />
        </view>
        <text class="nf-collect-empty-text">{{ $t('noCollectData') }}</text>
      </view>
    </scroll-view>

    <!-- Netflix风格登录弹窗 -->
    <view v-if="showLoginModal" class="nf-modal-mask" @tap.stop>
      <view class="nf-modal-box">
        <text class="nf-modal-title">{{ $t('loginFirst') }}</text>
        <view class="nf-modal-btns">
          <view class="nf-modal-btn nf-modal-cancel" @tap="onModalCancel">
            <text>{{ $t('cancel') }}</text>
          </view>
          <view class="nf-modal-btn nf-modal-confirm" @tap="onModalConfirm">
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
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const userStore = useUserStore()
const token = userStore.token || ''
const collectList = ref([])
const collectionFlag = ref('')
const showLoginModal = ref(false)

const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  const cents = parseInt(priceInCents)
  if (isNaN(cents)) return '0.00'
  return (cents / 100).toFixed(2)
}

const init = async () => {
  params = { page: 1, pageSize: 10 }
  isBottom.value = false
  const res = await getCollectList(params)
  if (res.code === 0) {
    const list = res.data.list || []
    collectList.value = list
    isBottom.value = list.length < params.pageSize
  }
}

onShow(() => {
  const tk = uni.getStorageSync('x-token')
  if (!tk) {
    showLoginModal.value = true
    return
  }
  showLoginModal.value = false
  init()
})

const onModalCancel = () => {
  showLoginModal.value = false
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.reLaunch({ url: '/pages/tabBar/index' })
  }
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
  if (res.code === 0 && res.data.list && res.data.list.length) {
    collectList.value.push(...res.data.list)
  } else {
    isBottom.value = true
  }
}
const debouncedLower = debounce(lower, 300)

const cancelCollect = async (ID) => {
  if (token) {
    const status = await findCollect({ goodID: ID })
    if (status.code === 0) collectionFlag.value = status.data
    const res = await createCollect({ goodID: Number(ID) })
    if (res.code === 0) {
      collectionFlag.value = !collectionFlag.value
      init()
      uni.showToast({
        title: collectionFlag.value ? $t.value('collected') : $t.value('uncollected'),
        mask: true,
        icon: 'none'
      })
    }
  } else {
    uni.showToast({ title: $t.value('loginFirst'), mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
  }
}

const goTo = (item) => {
  uni.navigateTo({ url: `/pages/player/index?id=${item.ID}` })
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
page { background: #000; }

/* Netflix风格登录弹窗 */
.nf-modal-mask {
  position: fixed;
  left: 0; top: 0; right: 0; bottom: 0;
  z-index: 9999;
  background: rgba(0,0,0,0.75);
  display: flex;
  align-items: center;
  justify-content: center;
}
.nf-modal-box {
  width: 560rpx;
  background: #1a1a1a;
  border-radius: 24rpx;
  padding: 60rpx 48rpx 48rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 8rpx 40rpx rgba(229,9,20,0.18);
}
.nf-modal-title {
  color: #fff;
  font-size: 34rpx;
  font-weight: 600;
  text-align: center;
  margin-bottom: 56rpx;
  letter-spacing: 1rpx;
}
.nf-modal-btns {
  display: flex;
  width: 100%;
  gap: 24rpx;
}
.nf-modal-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30rpx;
  font-weight: 500;
  letter-spacing: 1rpx;
}
.nf-modal-cancel {
  background: #333;
  color: #ccc;
}
.nf-modal-cancel:active { background: #444; }
.nf-modal-confirm {
  background: #e50914;
  color: #fff;
}
.nf-modal-confirm:active { background: #b20710; }

.nf-collect {
  min-height: 100vh;
  background: #000;
  display: flex;
  flex-direction: column;
}

.nf-collect-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 400rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 50% 0%, rgba(229, 9, 20, 0.12) 0%, transparent 60%);
}

/* ===== 导航栏 ===== */
.nf-navbar {
  position: fixed;
  top: 0; left: 0; right: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}

.nf-navbar-status {
  height: var(--status-bar-height, 44px);
}

.nf-navbar-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 2rpx;
}

/* ===== 滚动区域 ===== */
.nf-collect-scroll {
  flex: 1;
  margin-top: calc(var(--status-bar-height, 44px) + 88rpx);
  min-height: calc(100vh - var(--status-bar-height, 44px) - 88rpx);
}

/* ===== 商品列表 ===== */
.nf-goods-list {
  padding: 20rpx 24rpx;
}

.nf-goods-card {
  display: flex;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: transform 0.3s;

  &:active { transform: scale(0.98); }
}

.nf-goods-img-wrap {
  position: relative;
  width: 220rpx;
  height: 220rpx;
  flex-shrink: 0;
}

.nf-goods-img {
  width: 100%;
  height: 100%;
}

.nf-goods-img-overlay {
  position: absolute;
  top: 0; right: 0; bottom: 0;
  width: 40%;
  background: linear-gradient(to left, rgba(0,0,0,0.3), transparent);
  pointer-events: none;
}

.nf-goods-badge {
  position: absolute;
  top: 12rpx;
  left: 12rpx;
  background: linear-gradient(135deg, #e50914, #b20710);
  padding: 4rpx 14rpx;
  border-radius: 10rpx;
  box-shadow: 0 4rpx 12rpx rgba(229, 9, 20, 0.4);
}

.nf-goods-badge-text {
  font-size: 20rpx;
  color: #fff;
  font-weight: 700;
}

.nf-goods-info {
  flex: 1;
  padding: 16rpx 20rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
}

.nf-goods-title {
  font-size: 28rpx;
  color: #fff;
  font-weight: 600;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.nf-goods-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6rpx;
  margin: 8rpx 0;
}

.nf-tag {
  font-size: 18rpx;
  padding: 2rpx 10rpx;
  border-radius: 6rpx;
  border: 1rpx solid;
  font-weight: 500;
}

.nf-tag-red {
  color: #e50914;
  border-color: rgba(229, 9, 20, 0.3);
  background: rgba(229, 9, 20, 0.1);
}

.nf-tag-blue {
  color: #4a9eff;
  border-color: rgba(74, 158, 255, 0.3);
  background: rgba(74, 158, 255, 0.1);
}

.nf-tag-green {
  color: #4caf50;
  border-color: rgba(76, 175, 80, 0.3);
  background: rgba(76, 175, 80, 0.1);
}

.nf-goods-bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nf-price-row {
  display: flex;
  align-items: baseline;
}

.nf-price {
  font-size: 34rpx;
  color: #e50914;
  font-weight: 800;
}

.nf-goods-meta {
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.nf-rating-stars {
  font-size: 20rpx;
  color: #ffd700;
}

.nf-rating-score {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.6);
  font-weight: 600;
  margin-right: 8rpx;
}

.nf-sold {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.35);
}

/* ===== 底部状态 ===== */
.nf-collect-footer {
  padding: 40rpx;
  text-align: center;
}

.nf-collect-footer-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.25);
  letter-spacing: 2rpx;
}

/* ===== 空状态 ===== */
.nf-collect-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.nf-collect-empty-icon {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.08);
  border: 1rpx solid rgba(229, 9, 20, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
}

.nf-collect-empty-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.3);
  letter-spacing: 2rpx;
}
</style>
