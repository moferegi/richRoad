<template>
  <view class="flash-sale">
    <view class="flash-header">
      <image class="top-sell" src="@/static/top.png" mode="aspectFill"></image>
      <text class="flash-title">近期热销</text>
    </view>

    <view class="product-list">
      <view class="product-item" v-for="(item, index) in productData" :key="index" @click="goto(item)">
          <image class="product-image" :src="getUrl(item.imageUrl)" mode="aspectFill"></image>
        <view class="desc">
          <text class="product-title">{{ item.title }}</text>
          <text class="product-price">¥ {{ formatPrice(item.price) }}</text>
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

// 处理分转元的价格格式化，避免浮点数精度问题
const formatPrice = (priceInCents) => {
  if (!priceInCents && priceInCents !== 0) return '0.00'
  
  // 确保输入是数字
  const cents = parseInt(priceInCents)
  if (isNaN(cents)) return '0.00'
  
  // 使用整数运算避免精度问题
  const yuan = Math.floor(cents / 100)
  const remainingCents = cents % 100
  
  // 格式化为两位小数
  return `${yuan}.${remainingCents.toString().padStart(2, '0')}`
}

const goto = (item) => {
  uni.navigateTo({
    url: '/pages/goodsDetails/goodsDetails?id=' + item.ID
  })
}
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
  background-color: #1a1a1a;
  border-radius: 24rpx;
  padding: 24rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.3);
}

.flash-header {
  display: flex;
  align-items: center;
  margin-bottom: 24rpx;
  .top-sell{
    width: 40rpx;
    height: 40rpx;
    margin-right: 20rpx;
  }
}

.flash-title {
  font-size: 32rpx;
  font-weight: bold;
  color: #e50914;
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
    flex-direction: column;
    gap: 8rpx;
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
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

.product-price {
  font-size: 24rpx;
  color: #e50914;
  font-weight: bold;
  white-space: nowrap;
  flex-shrink: 0;
}
</style>

