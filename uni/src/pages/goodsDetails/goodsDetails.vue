<template>
  <view class="nf-goods-detail">
    <!-- 自定义导航栏（悬浮在轮播上方） -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('goodsDetail') || '商品详情' }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <!-- 商品轮播图 -->
    <goods-swiper :list="data.banner"></goods-swiper>

    <!-- 商品基本信息 -->
    <view class="nf-product-info">
      <view class="nf-price-row">
        <text class="nf-price-symbol">¥</text>
        <text class="nf-price-num">{{ data.price && (data.price / 100).toFixed(2) }}</text>
        <view class="nf-sale-tag" v-if="data.saleCount && data.saleCount > 0">
          <text>{{ $t('sold') || '已售' }} {{ data.saleCount }}</text>
        </view>
      </view>
      <text class="nf-product-title">{{ $lt(data.title) }}</text>
      <text class="nf-product-desc">{{ $lt(data.description) }}</text>
      <view class="nf-stock-row">
        <text>{{ $t('stock') || '库存' }}: {{ getTotalInventory(data.skus) }}</text>
      </view>
    </view>

    <!-- 优惠券 -->
    <view class="nf-section-card" @tap="opencoupon">
      <view class="nf-section-row">
        <view class="nf-section-label">
          <text class="nf-label-dot" style="background: #e50914;"></text>
          <text>{{ $t('coupon') || '优惠券' }}</text>
        </view>
        <view class="nf-section-value">
          <text class="nf-coupon-text" v-if="selectedCoupon.couponNum">
            {{ selectedCoupon.minSpend > 0
              ? '满' + selectedCoupon.minSpend/100 + '减' + selectedCoupon.discount/100
              : '无门槛减￥' + selectedCoupon.discount/100 }}
          </text>
          <text class="nf-coupon-hint" v-else>领取优惠券</text>
          <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)"></uni-icons>
        </view>
      </view>
    </view>

    <!-- 促销 + 服务 -->
    <view class="nf-section-card">
      <view class="nf-promo-title">
        <text class="nf-label-dot" style="background: #f59e0b;"></text>
        <text>促销活动</text>
      </view>
      <view class="nf-promo-list">
        <view class="nf-promo-item">🎁 新人首单送20元无门槛代金券</view>
        <view class="nf-promo-item">🔥 订单满50减10</view>
        <view class="nf-promo-item">💰 订单满100减30</view>
        <view class="nf-promo-item">🚚 单笔购买满两件免邮费</view>
      </view>
      <view class="nf-service-row">
        <text class="nf-service-tag">✓ 7天无理由退换货</text>
        <text class="nf-service-tag">✓ 假一赔十</text>
      </view>
    </view>

    <!-- 用户评价 -->
    <view class="nf-section-card">
      <view class="nf-review-header">
        <text class="nf-review-title">评价({{ commentInfo.length || 0 }})</text>
        <view class="nf-review-more" @tap="toEvaluate">
          <text>好评率 100%</text>
          <uni-icons type="right" size="14" color="rgba(255,255,255,0.3)"></uni-icons>
        </view>
      </view>

      <view v-if="hasContent" class="nf-review-content">
        <view class="nf-comment-item" v-if="commentInfo.content && commentInfo.content.length > 0">
          <view class="nf-comment-user">
            <image :src="getUrl(commentInfo.user.avatar)" class="nf-comment-avatar"></image>
            <text class="nf-comment-name">{{ commentInfo.user.nickname }}</text>
          </view>
          <text class="nf-comment-text">{{ commentInfo.content }}</text>
          <text class="nf-comment-time">{{ formatTimeToStr(commentInfo.CreatedAt, 'yyyy-MM-dd HH:mm') }}</text>
        </view>
      </view>
      <view v-else class="nf-no-review">暂无评价</view>
    </view>

    <!-- 图文详情 -->
    <view class="nf-detail-divider">
      <view class="nf-divider-line"></view>
      <text>图文详情</text>
      <view class="nf-divider-line"></view>
    </view>

    <view class="nf-detail-content">
      <goodsDetail :detail="data.detail"></goodsDetail>
      <rich-text style="width: 100%;" />
    </view>

    <!-- SKU选择器 -->
    <goods-sku
      style="z-index: 999;"
      v-if="data.skus"
      ref="goodsSkuRef"
      :isCart="isCart"
      :good="data"
      @toOrder="toOrder"
      :selectedCoupon="selectedCoupon"
    ></goods-sku>

    <!-- 底部导航 -->
    <view class="nf-bottom-nav">
      <view class="nf-nav-icons">
        <view class="nf-nav-icon-item" @tap="goTo()">
          <image class="nf-nav-icon-img" src="./../../static/images/tabBar/home.png"></image>
          <text class="nf-nav-icon-text">{{ $t('home') }}</text>
        </view>
        <view class="nf-nav-icon-item" @tap="goToKefu">
          <uni-icons type="headphones" size="20" color="rgba(255,255,255,0.6)"></uni-icons>
          <text class="nf-nav-icon-text">{{ $t('customerService') }}</text>
        </view>
        <view class="nf-nav-icon-item" @tap="goTo('cart')">
          <image class="nf-nav-icon-img" src="./../../static/images/tabBar/cart.png"></image>
          <text class="nf-nav-icon-text">{{ $t('cart') }}</text>
        </view>
        <view class="nf-nav-icon-item" @tap="addCollect">
          <image class="nf-nav-icon-img" :src="!collectionFlag ? './../../static/collection.png' : './../../static/collect.png'"></image>
          <text class="nf-nav-icon-text">{{ $t('collectText') }}</text>
        </view>
      </view>
      <view class="nf-buy-buttons">
        <view class="nf-add-cart-btn" @tap="addToCart()">{{ $t('addToCart') }}</view>
        <view class="nf-buy-now-btn" @tap="goodsTapPay('pay')">{{ $t('buyNow') }}</view>
      </view>
    </view>

    <!-- 优惠券弹出层 -->
    <view class="nf-mask" v-if="couponshow" @tap="hidecoupon"></view>
    <view class="nf-coupon-popup" :class="{ show: couponshow }">
      <scroll-view class="nf-coupon-scroll" scroll-y>
        <cc-defineCoupon v-if="couponshow" :goodIds="[data.ID]" colors="#e50914" @onReceive="onReceive"></cc-defineCoupon>
      </scroll-view>
    </view>
  </view>
</template>

<script setup>
import goodsSwiper from './components/goods-swiper.vue'
import goodsSku from './components/goods-sku.vue'
import goodsDetail from './components/goods-detail.vue'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { findGood } from '@/api/product.js'
import { myRouter } from '@/utils/permission'
import { findCollect, createCollect } from '@/api/collect.js'
import { claimCouponByUser } from '@/api/coupon.js'
import { useUserStore } from '@/pinia/modules/user'
import { findComment } from '@/api/comment.js'
import { formatTimeToStr } from '@/utils/date.js'
import { getUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $lt = computed(() => langStore.$lt)
const $t = computed(() => langStore.$t)

const goBack = () => { uni.navigateBack() }

const data = ref({})
const collectionFlag = ref('')
const hasContent = ref(false)
const goodID = ref(0)
const userStore = useUserStore()
const token = userStore.token || ''
const commentInfo = ref([])

onLoad((options) => {
  if (options.id) {
    goodID.value = options.id
    hasContent.value = false
    init()
  }
})

const init = async () => {
  const res = await findGood(goodID.value)
  if (res.code === 0) data.value = res.data.regood
  if (token) {
    const status = await findCollect({ goodID: goodID.value })
    if (status.code === 0) collectionFlag.value = status.data
  }
  const res2 = await findComment({ ID: goodID.value })
  if (res2.code === 0 && res2.data.length) {
    commentInfo.value = res2.data[0]
    hasContent.value = true
  }
}

const toEvaluate = () => {
  uni.navigateTo({ url: `/pages/evaluate/evaluate?goodsID=${goodID.value}` })
}

const goodsSkuRef = ref()
const isCart = ref(false)

const addToCart = () => {
  if (!token) {
    uni.showToast({ title: '请登录后进行操作', mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
    return
  }
  isCart.value = true
  goodsSkuRef.value.showSku()
}

const goTo = (path) => {
  if (path === 'cart') uni.switchTab({ url: '/pages/tabBar/shop/shop' })
  else uni.switchTab({ url: '/pages/tabBar/index' })
}

const goToKefu = () => { uni.navigateTo({ url: '/pages/kefu/index' }) }

const toOrder = () => {
  myRouter(`/pages/orderInfo/orderInfo?skuID=${data.value.skus[0].ID}&goodID=${data.value.skus[0].goodID}`)
}

const goodsTapPay = () => {
  isCart.value = false
  goodsSkuRef.value.showSku()
}

const addCollect = async () => {
  if (token) {
    const res = await createCollect({ goodID: Number(goodID.value) })
    if (res.code === 0) {
      collectionFlag.value = !collectionFlag.value
      uni.showToast({ title: collectionFlag.value ? '已收藏' : '已取消', mask: true, icon: 'none' })
    }
  } else {
    uni.showToast({ title: '请登录后进行操作', mask: true, icon: 'none' })
    uni.redirectTo({ url: '/pages/user/login' })
  }
}

const selectedCoupon = ref({})
const couponshow = ref(false)

const getTotalInventory = (skus) => {
  if (!skus || skus.length === 0) return 0
  return skus.reduce((total, sku) => total + sku.inventory, 0)
}

const opencoupon = () => { couponshow.value = true }
const hidecoupon = () => { couponshow.value = false }

const onReceive = async (item) => {
  uni.showLoading({ title: item.couponNum == 0 ? '领取中...' : '选择中...', mask: true })
  try {
    if (item.couponNum == 0) {
      const res = await claimCouponByUser({ couponID: item.couponID })
      item.couponNum = res.data
      uni.showToast({ title: '领取成功', icon: 'success', duration: 1500 })
    }
    selectedCoupon.value = item
    setTimeout(() => { hidecoupon() }, 500)
  } catch (error) {
    uni.showToast({ title: '操作失败，请重试', icon: 'none' })
  } finally { uni.hideLoading() }
}
</script>

<style lang="scss">
page { background-color: #000; }

.nf-goods-detail { min-height: 100vh; background: #000; padding-bottom: 120rpx; }

/* 导航栏 */
.nf-navbar {
  position: fixed; top: 0; left: 0; right: 0; z-index: 100;
  background: rgba(0, 0, 0, 0.6); backdrop-filter: blur(24px);
  padding: 0 28rpx 16rpx;
}
.nf-navbar-status { height: var(--status-bar-height, 0px); }
.nf-navbar-content { display: flex; align-items: center; justify-content: space-between; height: 88rpx; }
.nf-navbar-back {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: rgba(255, 255, 255, 0.12); border: 1rpx solid rgba(255, 255, 255, 0.15);
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-title { font-size: 34rpx; font-weight: 700; color: #fff; letter-spacing: 2rpx; }

/* 商品信息 */
.nf-product-info {
  padding: 28rpx 24rpx;
  background: rgba(255, 255, 255, 0.04);
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);
}
.nf-price-row {
  display: flex; align-items: baseline; gap: 4rpx; margin-bottom: 16rpx;
}
.nf-price-symbol { font-size: 30rpx; font-weight: 700; color: #e50914; }
.nf-price-num { font-size: 52rpx; font-weight: 800; color: #e50914; line-height: 1; }
.nf-sale-tag {
  margin-left: 16rpx; font-size: 20rpx; color: rgba(255, 255, 255, 0.4);
  padding: 4rpx 14rpx; background: rgba(255, 255, 255, 0.06); border-radius: 8rpx;
}
.nf-product-title { display: block; font-size: 32rpx; font-weight: 700; color: #fff; line-height: 1.4; margin-bottom: 8rpx; }
.nf-product-desc { display: block; font-size: 26rpx; color: rgba(255, 255, 255, 0.5); line-height: 1.4; margin-bottom: 12rpx; }
.nf-stock-row { font-size: 24rpx; color: rgba(255, 255, 255, 0.3); }

/* 通用区块卡片 */
.nf-section-card {
  margin: 20rpx 24rpx; padding: 24rpx 28rpx;
  background: rgba(255, 255, 255, 0.04); border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx; backdrop-filter: blur(8px);
}
.nf-section-row { display: flex; justify-content: space-between; align-items: center; }
.nf-section-label { display: flex; align-items: center; gap: 12rpx; font-size: 28rpx; color: #fff; }
.nf-label-dot { width: 12rpx; height: 12rpx; border-radius: 4rpx; }
.nf-section-value { display: flex; align-items: center; gap: 8rpx; }
.nf-coupon-text { font-size: 26rpx; color: #e50914; font-weight: 600; }
.nf-coupon-hint { font-size: 26rpx; color: rgba(255, 255, 255, 0.4); }

/* 促销 */
.nf-promo-title { display: flex; align-items: center; gap: 12rpx; font-size: 28rpx; color: #fff; margin-bottom: 16rpx; }
.nf-promo-list { margin-bottom: 16rpx; }
.nf-promo-item {
  font-size: 24rpx; color: rgba(255, 255, 255, 0.6); padding: 6rpx 0;
}
.nf-service-row {
  display: flex; gap: 20rpx; padding-top: 16rpx;
  border-top: 1rpx solid rgba(255, 255, 255, 0.04);
}
.nf-service-tag { font-size: 24rpx; color: rgba(255, 255, 255, 0.5); }

/* 评价 */
.nf-review-header {
  display: flex; justify-content: space-between; align-items: center;
  padding-bottom: 16rpx; border-bottom: 1rpx solid rgba(255, 255, 255, 0.04); margin-bottom: 16rpx;
}
.nf-review-title { font-size: 28rpx; font-weight: 700; color: #fff; }
.nf-review-more { display: flex; align-items: center; gap: 8rpx; font-size: 24rpx; color: rgba(255, 255, 255, 0.4); }
.nf-comment-item { }
.nf-comment-user { display: flex; align-items: center; gap: 12rpx; margin-bottom: 12rpx; }
.nf-comment-avatar { width: 48rpx; height: 48rpx; border-radius: 50%; border: 1rpx solid rgba(255, 255, 255, 0.1); }
.nf-comment-name { font-size: 26rpx; color: rgba(255, 255, 255, 0.7); }
.nf-comment-text { font-size: 26rpx; color: rgba(255, 255, 255, 0.6); line-height: 1.5; display: block; margin-bottom: 8rpx; }
.nf-comment-time { font-size: 22rpx; color: rgba(255, 255, 255, 0.3); }
.nf-no-review { text-align: center; padding: 30rpx 0; color: rgba(255, 255, 255, 0.3); font-size: 26rpx; }

/* 图文详情分割 */
.nf-detail-divider {
  display: flex; align-items: center; justify-content: center;
  gap: 20rpx; padding: 30rpx 0;
  text { font-size: 26rpx; color: rgba(255, 255, 255, 0.4); }
}
.nf-divider-line {
  height: 1rpx; width: 100rpx;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.12), transparent);
}
.nf-detail-content { padding: 0 24rpx; }

/* 底部导航 */
.nf-bottom-nav {
  position: fixed; left: 0; right: 0; bottom: 0; z-index: 98;
  display: flex; align-items: center; height: 110rpx;
  background: rgba(0, 0, 0, 0.95); backdrop-filter: blur(24px);
  border-top: 1rpx solid rgba(255, 255, 255, 0.06);
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}
.nf-nav-icons { display: flex; flex: 1; }
.nf-nav-icon-item {
  display: flex; flex-direction: column; align-items: center; justify-content: center; width: 96rpx;
}
.nf-nav-icon-img { width: 40rpx; height: 40rpx; margin-bottom: 4rpx; opacity: 0.6; }
.nf-nav-icon-text { font-size: 20rpx; color: rgba(255, 255, 255, 0.5); }
.nf-buy-buttons { display: flex; height: 100%; }
.nf-add-cart-btn, .nf-buy-now-btn {
  padding: 0 36rpx; height: 100%; display: flex; align-items: center; justify-content: center;
  font-size: 28rpx; font-weight: 600; color: #fff;
}
.nf-add-cart-btn { background: rgba(255, 149, 0, 0.9); }
.nf-buy-now-btn { background: #e50914; }

/* 优惠券弹出层 */
.nf-mask {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0, 0, 0, 0.7); z-index: 900;
}
.nf-coupon-popup {
  position: fixed; left: 0; right: 0; bottom: -100vh; z-index: 999;
  background: #1a1a1a; border-radius: 24rpx 24rpx 0 0;
  transition: all 0.3s ease;
  &.show { bottom: 0; }
}
.nf-coupon-scroll { width: 100vw; height: 60vh; padding-top: 16rpx; }
</style>


