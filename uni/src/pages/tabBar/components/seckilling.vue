<template>
  <view class="flash-sale">
    <view class="flash-header">
      <text class="flash-title">限时秒杀</text>
      <text class="flash-countdown">8点场 07:13:55</text>
      <uni-icons type="right" size="14" color="#999"></uni-icons>
    </view>

    <view class="product-list">
      <view class="product-item" v-for="(item, index) in productData" :key="index">
        <view class="product-image"></view>
        <text class="product-title">{{ item.title }}</text>
        <text class="product-price">¥ {{ item.price }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  productData: {
    type: Array,
    default: () => [] // 修改为函数返回空数组
  }
})
</script>

<style scoped lang="scss">
// 限时秒杀区域
.flash-sale {
  margin: 12px;
  background-color: #fff;
  border-radius: 12px;
  padding: 12px;
}

.flash-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.flash-title {
  font-size: 16px;
  font-weight: bold;
  color: #FF4500; // 猩红色，像原型图
  margin-right: auto;
}

.flash-countdown {
  font-size: 12px;
  color: #666;
  margin-right: 4px;
}

.product-list {
  display: flex;
  overflow-x: auto;
  padding-bottom: 8px;
  &::-webkit-scrollbar {
    display: none;
  }
}

.product-item {
  width: 120px;
  margin-right: 12px;
  flex-shrink: 0;
}

.product-image {
  width: 120px;
  height: 150px;
  background-color: #f5f5f5;
  border-radius: 8px;
  margin-bottom: 6px;
}

.product-title {
  font-size: 12px;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 4px;
}

.product-price {
  font-size: 14px;
  color: #FF4500;
  font-weight: bold;
}
</style>
