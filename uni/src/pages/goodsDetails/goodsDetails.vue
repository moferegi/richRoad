<template>
  <view class="goods-detail-container">
    <!-- 商品轮播图 -->
    <goods-swiper :list="data.banner"></goods-swiper>

    <!-- 商品基本信息 -->
    <view class="product-info">
      <view class="price-section">
        <view class="current-price">
          <text class="price-symbol">¥</text>
          <text class="price-num">{{ data.price && data.price / 100 }}</text>
        </view>
        <view class="original-price">¥{{ (data.price && data.price / 100 * 1.3).toFixed(1) }}</view>
        <view class="discount-tag">{{ 7 }}折</view>
      </view>
      <view class="title-section">
        <text class="product-title">{{ data.title }}</text>
        <text class="product-desc">{{ data.description }}</text>
      </view>
      <view class="sales-info">
        <text v-if="data.saleCount && data.saleCount>0">销量: {{ data.saleCount }}</text>
        <text class="stock-info">库存: {{ data.stock || 4690 }}</text>
        <text class="view-count">浏览量: {{ data.viewCount || 768 }}</text>
      </view>

      <!-- 优惠信息卡片 -->
      <view class="promotion-card">
        <view class="share-promotion">
<!--          <image class="promo-icon" src="./../../static/promo-icon.png"></image>-->
          <text class="promo-text">该商品分享可领49减10红包</text>
          <view class="share-now">立即分享 ></view>
        </view>
      </view>
    </view>

    <!-- 规格选择 -->
    <view class="specs-section" @tap="goodsSkuRef.showSku()">
      <text class="specs-label">购买类型</text>
      <view class="specs-value">
        <text>XS 白色</text>
        <uni-icons color="#cccccc" size="16" type="right"></uni-icons>
      </view>
    </view>

    <!-- 优惠券 -->
    <view class="coupon-section">
      <text class="coupon-label">优惠券</text>
      <view class="coupon-value" @tap="showCoupons">
        <text>领取优惠券</text>
        <uni-icons color="#cccccc" size="16" type="right"></uni-icons>
      </view>
    </view>

    <!-- 促销活动 -->
    <view class="promotion-section">
      <text class="promotion-label">促销活动</text>
      <view class="promotion-list">
        <view class="promotion-item">新人首单送20元无门槛代金券</view>
        <view class="promotion-item">订单满50减10</view>
        <view class="promotion-item">订单满100减30</view>
        <view class="promotion-item">单笔购买满两件免邮费</view>
      </view>
    </view>

    <!-- 服务保障 -->
    <view class="service-section">
      <text class="service-label">服务</text>
      <view class="service-list">
        <text class="service-item">7天无理由退换货</text>
        <text class="service-dot">·</text>
        <text class="service-item">假一赔十</text>
      </view>
    </view>

    <!-- 用户评价 -->
    <view class="review-section">
      <view class="review-header">
        <text class="review-title">评价({{ commentInfo.length || 0 }})</text>
        <view class="review-rate" @tap="toEvaluate">
          <text class="good-rate">好评率 100%</text>
          <uni-icons color="#cccccc" size="16" type="right"></uni-icons>
        </view>
      </view>

      <!-- 用户评价内容 -->
      <view v-if="hasContent" class="review-content">
        <view class="comment-section section-card">

          <view class="comment-item" v-if="commentInfo.content.length > 0">
            <view class="comment-user">
              <image :src="getUrl(commentInfo.user.avatar)" class="reviewer-avatar"></image>
              <text class="user-name">{{ commentInfo.user.nickname }}</text>
            </view>
            <view class="comment-content">
              <text class="comment-text">{{ commentInfo.content }}</text>
              <view class="comment-images">
                <view class="comment-image" :style="{ backgroundColor: '#f0f0f0' }"></view>
              </view>
              <view class="comment-info">
                <text class="comment-spec">{{ commentInfo.feedbackPics&&commentInfo.feedbackPics[0] }}</text>
                <text class="comment-time">{{ formatTimeToStr(commentInfo.CreatedAt, 'yyyy-MM-dd HH:mm') }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <view v-if="!hasContent" class="no-review">
        暂无评价
      </view>
    </view>

    <!-- 商品详情标题 -->
    <view class="detail-title">
      <view class="title-line"></view>
      <text>图文详情</text>
      <view class="title-line"></view>
    </view>

    <!-- 商品详情内容 -->
    <goodsDetail :detail="data.detail"></goodsDetail>
    <rich-text style="width: 100%;"/>

    <!-- SKU选择器 (隐藏状态) -->
    <goods-sku v-if="data.skus" ref="goodsSkuRef" :isCart="isCart" :good="data" @toOrder="toOrder"></goods-sku>

    <!-- 底部固定导航栏 -->
    <view class="fixed-bottom-nav">
      <view class="nav-action-buttons">
        <view class="nav-button" @tap="goTo()">
          <image class="nav-icon" src="./../../static/images/tabBar/home.png"></image>
          <text class="nav-text">首页</text>
        </view>
        <view class="nav-button" @tap="goTo('cart')">
          <image class="nav-icon" src="./../../static/images/tabBar/cart.png"></image>
          <text class="nav-text">购物车</text>
        </view>
        <view class="nav-button" @tap="addCollect">
          <image class="nav-icon" :src="!collectionFlag? './../../static/collection.png' : './../../static/collect.png'"></image>
          <text class="nav-text">收藏</text>
        </view>
      </view>
      <view class="buy-buttons">
        <view class="add-cart-btn" @tap="addToCart()">加入购物车</view>
        <view class="buy-now-btn" @tap="goodsTapPay('pay')">立即购买</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import goodsSwiper from './components/goods-swiper.vue'
import goodsSku from './components/goods-sku.vue'
import goodsDetail from './components/goods-detail.vue';
import {ref} from "vue";
import {onLoad} from '@dcloudio/uni-app'
import {findGood} from '@/api/product.js'
import {myRouter} from '@/utils/permission';
import {findCollect, createCollect} from '@/api/collect.js'
import {useUserStore} from "@/pinia/modules/user";
import {findComment} from "@/api/comment.js"
import {formatTimeToStr} from "@/utils/date.js"
import evaluateGridImg from '@/pages/evaluate/evaluate-img.vue'
import {getUrl} from "@/utils/url.js"
const data = ref({})
const collectionFlag = ref('')
const hasContent = ref(false)
const goodID = ref(0)
const userStore = useUserStore()
const token = userStore.token || ''
const commentInfo = ref([])
onLoad((options) => {
  // 先获取商品属性
  if (options.id) {
    goodID.value = options.id
    hasContent.value = false
    init()
  }
})

const init = async () => {
  const res = await findGood(goodID.value)
  res.code === 0 ? data.value = res.data.regood : ''
  // 如果登录，就获取当前商品收藏状态 反之则默认未收藏，点击跳转 登录后收藏
  if (token) {
    // 先查看当前商品收藏状态
    const status = await findCollect({
      goodID: goodID.value
    })
    status.code === 0 ? collectionFlag.value = status.data : ''
  }

  const res2 = await findComment(goodID.value)
  if (res2.code === 0 && res2.data.length) {
    commentInfo.value = res2.data[0]
    hasContent.value = true
  }
}

const toEvaluate = () => {
  uni.navigateTo({
    url: `/pages/evaluate/evaluate?goodsID=${goodID.value}`
  })
}

const goodsSkuRef = ref()

const isCart = ref(false)

const addToCart = async () => {
  // 需要校验用户是否登录，未登录跳转至登录页，否则允许加入购物车
  if (!token) {
    uni.showToast({
      title: '请登录后进行操作',
      mask: true,
      icon: 'none'
    });
    uni.redirectTo({
      url: '/pages/user/login'
    })
    return
  }
  //传递参数让子组件知道是加入购物车还是立即购买
  isCart.value = true
  goodsSkuRef.value.showSku()
}

const goTo = (path) => {
  if (path === "cart") {
    uni.switchTab({
      url: '/pages/tabBar/shop/shop'
    })
  } else {
    uni.switchTab({
      url: '/pages/tabBar/index'
    })
  }

}
const toOrder = async () => {
  myRouter(`/pages/orderInfo/orderInfo?skuID=${data.value.skus[0].ID}&goodID=${data.value.skus[0].goodID}`)
}

const goodsTapPay = (pay) => {
  isCart.value = false
  goodsSkuRef.value.showSku()
}

const addCollect = async () => {
  if (token) {
    // 如果已登录并且未收藏 则允许进行收藏操作
    const res = await createCollect({
      goodID: Number(goodID.value)
    })
    if (res.code === 0) {
      collectionFlag.value = !collectionFlag.value
      uni.showToast({
        title: collectionFlag.value ? '已收藏' : '已取消',
        mask: true,
        icon: 'none'
      });
    }
  } else {
    uni.showToast({
      title: '请登录后进行操作',
      mask: true,
      icon: 'none'
    });
    uni.redirectTo({
      url: '/pages/user/login'
    })
  }

}
</script>

<style lang="scss">
page {
  background-color: #f5f5f5;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, 'Helvetica Neue', Helvetica, 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.goods-detail-container {
  padding-bottom: 120rpx;
}

/* 商品基本信息区域 */
.product-info {
  background-color: #fff;
  padding: 20rpx 24rpx;
  margin-bottom: 20rpx;
}

.price-section {
  display: flex;
  align-items: center;
  margin-bottom: 16rpx;

  .current-price {
    color: #e4393c;

    .price-symbol {
      font-size: 32rpx;
      font-weight: bold;
    }

    .price-num {
      font-size: 48rpx;
      font-weight: bold;
    }
  }

  .original-price {
    font-size: 24rpx;
    color: #999;
    text-decoration: line-through;
    margin-left: 16rpx;
    margin-top: 8rpx;
  }

  .discount-tag {
    background-color: #e4393c;
    color: #fff;
    font-size: 20rpx;
    padding: 2rpx 8rpx;
    border-radius: 4rpx;
    margin-left: 12rpx;
  }
}

.title-section {
  margin-bottom: 16rpx;

  .product-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #333;
    line-height: 1.4;
    display: block;
    margin-bottom: 8rpx;
  }

  .product-desc {
    font-size: 26rpx;
    color: #666;
    line-height: 1.4;
    display: block;
  }
}

.sales-info {
  display: flex;
  font-size: 24rpx;
  color: #999;
  margin-bottom: 16rpx;

  text {
    margin-right: 24rpx;
  }
}

/* 分享优惠卡片 */
.promotion-card {
  background-color: #fff8f8;
  border-radius: 8rpx;
  padding: 16rpx;

  .share-promotion {
    display: flex;
    align-items: center;

    .promo-icon {
      width: 32rpx;
      height: 32rpx;
      margin-right: 8rpx;
    }

    .promo-text {
      flex: 1;
      font-size: 24rpx;
      color: #e4393c;
    }

    .share-now {
      font-size: 24rpx;
      color: #e4393c;
    }
  }
}

/* 规格选择 */
.specs-section, .coupon-section, .promotion-section, .service-section {
  background-color: #fff;
  padding: 24rpx;
  display: flex;
  align-items: center;
  border-bottom: 1rpx solid #f0f0f0;
}

.specs-label, .coupon-label, .promotion-label, .service-label {
  width: 140rpx;
  font-size: 28rpx;
  color: #666;
}

.specs-value, .coupon-value {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 28rpx;
  color: #333;
}

/* 促销活动 */
.promotion-section {
  flex-direction: column;
  align-items: flex-start;
}

.promotion-list {
  width: 100%;
  margin-top: 16rpx;
}

.promotion-item {
  font-size: 26rpx;
  color: #e4393c;
  margin-bottom: 8rpx;
  position: relative;
  padding-left: 16rpx;

  &:before {
    content: '';
    position: absolute;
    left: 0;
    top: 12rpx;
    width: 8rpx;
    height: 8rpx;
    background-color: #e4393c;
    border-radius: 50%;
  }
}

/* 服务保障 */
.service-section {
  margin-bottom: 20rpx;
}

.service-list {
  flex: 1;
  display: flex;
  align-items: center;
}

.service-item {
  font-size: 26rpx;
  color: #666;
}

.service-dot {
  margin: 0 8rpx;
  color: #ccc;
}

/* 用户评价 */
.review-section {
  background-color: #fff;
  padding: 24rpx;
  margin-bottom: 20rpx;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid #f0f0f0;
  margin-bottom: 20rpx;
}

.review-title {
  font-size: 28rpx;
  font-weight: bold;
  color: #333;
}

.review-rate {
  display: flex;
  align-items: center;

  .good-rate {
    font-size: 24rpx;
    color: #666;
    margin-right: 8rpx;
  }
}

.reviewer-info {
  display: flex;
  align-items: center;
  margin-bottom: 16rpx;
}

.reviewer-avatar {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  margin-right: 12rpx;
}

.reviewer-name {
  font-size: 26rpx;
  color: #333;
  margin-right: 12rpx;
}

.review-time {
  font-size: 24rpx;
  color: #999;
}

.review-text {
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
  margin-bottom: 16rpx;
}

.purchase-info {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 16rpx;
}

.no-review {
  text-align: center;
  padding: 30rpx 0;
  color: #999;
  font-size: 28rpx;
}

/* 商品详情标题 */
.detail-title {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx 0;

  text {
    margin: 0 16rpx;
    font-size: 28rpx;
    color: #666;
  }

  .title-line {
    height: 1rpx;
    background-color: #e0e0e0;
    flex: 1;
    max-width: 200rpx;
  }
}

/* 底部导航 */
.fixed-bottom-nav {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  height: 100rpx;
  background-color: #fff;
  display: flex;
  align-items: center;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
  z-index: 999;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.nav-action-buttons {
  display: flex;
  flex: 1;
}

.nav-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100rpx;

  .nav-icon {
    width: 44rpx;
    height: 44rpx;
    margin-bottom: 4rpx;
  }

  .nav-text {
    font-size: 20rpx;
    color: #666;
  }
}

.comment-section {
  flex-direction: column;
  align-items: flex-start;
}

.comment-header {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.comment-rate {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #666;
}

.comment-item {
  width: 100%;
}

.comment-user {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.user-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  margin-right: 8px;
}

.user-name {
  font-size: 14px;
  color: #333;
}

.comment-content {
  padding-left: 38px;
}

.comment-text {
  font-size: 14px;
  color: #333;
  line-height: 1.4;
  margin-bottom: 10px;
}

.comment-images {
  display: flex;
  margin-bottom: 10px;
}

.comment-image {
  width: 80px;
  height: 80px;
  margin-right: 8px;
  border-radius: 4px;
}

.comment-info {
  display: flex;
  font-size: 12px;
  color: #999;
}

.comment-spec {
  margin-right: 15px;
}

.buy-buttons {
  display: flex;
  align-items: center;
  height: 100%;

  .add-cart-btn, .buy-now-btn {
    padding: 0 40rpx;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28rpx;
    font-weight: 500;
    color: #fff;
  }

  .add-cart-btn {
    background-color: #ff9500;
  }

  .buy-now-btn {
    background-color: #ff5000;
  }
}
</style>
