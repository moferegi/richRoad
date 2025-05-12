<template>
  <view class="container">
    <!-- 商品轮播图 -->
    <goodsSwiper></goodsSwiper>

    <!-- 商品信息区域 -->
    <GoodsInfoComponent :infoData="goodsInfo"></GoodsInfoComponent>

    <!-- 分享和返红包区域 -->
    <view class="share-section">
      <view class="red-packet">
        <text>该商品分享可得49或10红包</text>
        <view class="red-share">
          <text class="detail-btn">立即分享</text>
          <uni-icons type="redo" size="16" color="#C71585"></uni-icons>
        </view>
      </view>
    </view>

    <!-- 规格选择区域 -->
    <goods-choose
        :selected-sku="selectedSku"
        @open-sku="openSkuPopup"
    ></goods-choose>

    <!-- 优惠券区域 -->
    <goodsCoupon></goodsCoupon>

    <!-- 促销活动区域 -->
    <goodsPromotion></goodsPromotion>

    <!-- 服务区域 -->
    <goodsService></goodsService>

    <!-- 评论区域 -->
    <goodsComment :commentData="goodsInfo"></goodsComment>

    <!-- 图文详情 -->
    <goodsDiscribe></goodsDiscribe>

    <!-- 底部操作栏 -->
    <view class="bottom-bar">
      <view class="bottom-left">
        <view class="bottom-btn">
          <uni-icons type="home" size="24" color="#666"></uni-icons>
          <text>首页</text>
        </view>
        <view class="bottom-btn">
          <uni-icons type="cart" size="24" color="#666"></uni-icons>
          <text>购物车</text>
        </view>
        <view class="bottom-btn">
          <uni-icons type="star" size="24" color="#666"></uni-icons>
          <text>收藏</text>
        </view>
      </view>
      <view class="bottom-right">
        <button class="add-cart-btn" @click="openSkuPopup">加入购物车</button>
        <button class="buy-now-btn" @click="buyNow">立即购买</button>
      </view>

    </view>

    <!-- SKU选择组件 -->
    <goods-sku
        v-model="skuShow"
        :sku-data="skuData"
        :action-type="actionType"
        @sku-confirm="handleSkuConfirm"
    ></goods-sku>
  </view>
</template>

<script setup>
import { ref } from 'vue';
import goodsSwiper from './components/goods-swiper.vue'
import GoodsInfoComponent from './components/goods-info.vue'
import goodsChoose from './components/goods-choose.vue'
import goodsCoupon from './components/goods-coupon.vue'
import goodsPromotion from './components/goods-promotion.vue'
import goodsService from './components/goods-service.vue'
import goodsComment from './components/goods-comment.vue'
import goodsDiscribe from './components/goods-describe.vue'
import goodsSku from './components/goods-SKU.vue'


// 商品信息
const goodsInfo = ref({
  id: 1,
  title: '恒源祥2019春季长袖白t恤 新款春装',
  price: '341.6',
  originalPrice: '488',
  discount: '7',
  sales: '108',
  stock: '4660',
  views: '758',
  commentCount: 96,
  goodRate: 100,
  comments: [
    {
      userName: 'Leo yo',
      content: '商品收到了，79元两件，质量不错，试了一下有点宽，但是加了个外套很好看，我很喜欢',
      spec: 'XL 红色',
      time: '2019-04-01 19:21'
    }
  ]
});

// SKU弹窗显示控制
const skuShow = ref(false);
// 操作类型：buy-立即购买，cart-加入购物车
const actionType = ref('');
// 已选择的SKU信息
const selectedSku = ref(null);

// Mock数据 - 实际项目中通过API获取
const skuData = ref({
  skus: [
    {
      id: 1,
      price: 7000,
      stock: 30,
      sku_attrs: {
        '机身颜色': {
          name: '深空黑色',
          img: 'https://example.com/black.jpg',
        },
        '储存容量': '128G',
        '套装': '快充套装'
      }
    },
    // ... 其他SKU数据
  ]
});

// 打开SKU选择弹窗
const openSkuPopup = (type = 'buy') => {
  console.log('打开SKU弹窗，类型:', type);
  actionType.value = type;
  skuShow.value = true;
};

const buyNow = () => {
  // 立即购买逻辑
  uni.navigateTo({
    url: '/pages/order/order-confirm'
  });
};


// 处理SKU确认
const handleSkuConfirm = (sku) => {
  selectedSku.value = sku;

  if (actionType.value === 'cart') {
    // 处理加入购物车逻辑
    console.log('加入购物车', sku);
  } else {
    // 处理立即购买逻辑
    console.log('立即购买', sku);
  }

  skuShow.value = false;
};
</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
  padding-bottom: 200rpx;
}

.container {
  width: 100%;
}

/* 分享和返红包区域 */
.share-section {
  background-color: #fff;
  margin-top: 2rpx;
}

.red-packet {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 28rpx;
  background: #FAEDF6;
  padding: 20rpx 30rpx;
  color: #3B4144;
}

.red-share {
}

.detail-btn {
  color: #C71585;
  padding-right: 12rpx;
}

/* 通用卡片样式 */
.section-card {
  background-color: #fff;
  padding: 30rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 28rpx;
  color: #666;
  margin-right: 20rpx;
  min-width: 80rpx;
}

.section-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 28rpx;
  color: #333;
}

.detail-section {
  background-color: #fff;
  margin-top: 20rpx;
  padding-bottom: 40rpx;
}

.detail-title {
  text-align: center;
  font-size: 32rpx;
  color: #333;
  padding: 30rpx 0;
  position: relative;
}

.detail-title::before,
.detail-title::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 120rpx;
  height: 2rpx;
  background-color: #ddd;
}

.detail-title::before {
  left: 100rpx;
}

.detail-title::after {
  right: 100rpx;
}

.detail-content {
  padding: 0 30rpx;
}

.detail-image {
  width: 100%;
  margin-bottom: 20rpx;
  border-radius: 8rpx;
}

/* 底部操作栏 */
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 100rpx;
  background-color: #fff;
  display: flex;
  border-top: 2rpx solid #eee;
  z-index: 99;
}

.bottom-left {
  display: flex;
  width: 40%;
}

.bottom-btn {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
  color: #666;
}

.bottom-right {
  display: flex;
  width: 60%;
}

.add-cart-btn,
.buy-now-btn {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28rpx;
  border-radius: 0;
}

.add-cart-btn {
  background-color: #ff9800;
  color: #fff;
}

.buy-now-btn {
  background-color: #ff4757;
  color: #fff;
}
</style>

