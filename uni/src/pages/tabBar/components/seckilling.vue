<template>
  <view class="flash-sale">
    <view class="flash-header">
      <text class="flash-title">近期热销</text>
    </view>

    <view class="product-list">
      <view class="product-item" v-for="(item, index) in productData" :key="index">
          <image class="product-image" :src="getUrl(item.imageUrl)" mode="aspectFill"></image>
        <view class="desc">
          <text class="product-title">{{ item.title }}</text>
          <text class="product-price">¥ {{ item.price }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import {getUrl} from "@/utils/url";

const props = defineProps({
  productData: {
    type: Array,
    default: () => [] // 修改为函数返回空数组
  }
})

const init = () => {
  // 初始化逻辑，如果需要的话
  console.log(props.productData);
}
setTimeout(() => {
  init();
}, 1000);
</script>

<style scoped lang="scss">
// 限时秒杀区域
.flash-sale {
  margin: 24rpx;
  background-color: #fff;
  border-radius: 24rpx;
  padding: 24rpx;
}

.flash-header {
  display: flex;
  align-items: center;
  margin-bottom: 24rpx;
}

.flash-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #FF4500;
  margin-right: auto;
}

.flash-countdown {
  font-size: 24rpx;
  color: #666;
  margin-right: 8rpx;
}

.product-list {
  display: flex;
  overflow-x: auto;
  padding-bottom: 16rpx;
  &::-webkit-scrollbar {
    display: none;
  }
  .desc{
    padding: 12rpx;
    display: flex;
    align-items: baseline;
    justify-content: space-between;
  }
}

.product-item {
  width: 240rpx;
  margin-right: 24rpx;
  flex-shrink: 0;
}

.product-image {
  width: 240rpx;
  height: 300rpx;
  background-color: #f5f5f5;
  border-radius: 16rpx;
  margin-bottom: 12rpx;
  box-shadow: 0 6rpx 6rpx rgba(0, 0, 0, 0.1);

}

.product-title {
  font-size: 28rpx;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 8rpx;
}

.product-price {
  font-size: 24rpx;
  color: #FF4500;
  font-weight: bold;
}
</style>

