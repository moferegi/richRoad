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
    <goodsChoose></goodsChoose>

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
        <button class="add-cart-btn">加入购物车</button>
        <button class="buy-now-btn">立即购买</button>
      </view>
    </view>
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
</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
  padding-bottom: 100px;
}

.container {
  width: 100%;
}

/* 分享和返红包区域 */
.share-section {
  background-color: #fff;
  margin-top: 1px;
}

.red-packet {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  background: #FAEDF6;
  padding: 10px 15px;
  color: #3B4144;
}
.red-share{
}

.detail-btn {
  color: #C71585;
  padding-right: 6px;
}

/* 通用卡片样式 */
.section-card {
  background-color: #fff;
  padding: 15px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  font-size: 14px;
  color: #666;
  margin-right: 10px;
  min-width: 40px;
}

.section-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #333;
}
.detail-section {
  background-color: #fff;
  margin-top: 10px;
  padding-bottom: 20px;
}

.detail-title {
  text-align: center;
  font-size: 16px;
  color: #333;
  padding: 15px 0;
  position: relative;
}

.detail-title::before, .detail-title::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 60px;
  height: 1px;
  background-color: #ddd;
}

.detail-title::before {
  left: 50px;
}

.detail-title::after {
  right: 50px;
}

.detail-content {
  padding: 0 15px;
}

.detail-image {
  width: 100%;
  margin-bottom: 10px;
  border-radius: 4px;
}

/* 底部操作栏 */
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 50px;
  background-color: #fff;
  display: flex;
  border-top: 1px solid #eee;
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
  font-size: 10px;
  color: #666;
}

.bottom-right {
  display: flex;
  width: 60%;
}

.add-cart-btn, .buy-now-btn {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
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
