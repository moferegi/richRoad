<template>
  <view class="goods-info">
    <view class="price-section">
      <text class="price">¥{{ infoData.price }}</text>
      <text class="original-price">¥{{ infoData.originalPrice }}</text>
      <text class="discount">{{ infoData.discount }}折</text>
    </view>
    <view class="title-section">
      <text class="goods-title">{{ infoData.title }}</text>
    </view>
    <view class="sales-info">
      <text class="sales">销量: {{ infoData.sales }}</text>
      <text class="stock">库存: {{ infoData.stock }}</text>
      <text class="views">浏览量: {{ infoData.views }}</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';
const props = defineProps({
  infoData: {
    type: Object, // 修改为Object类型，因为传入的是一个对象
    default: () => ({}) // 默认值应该是空对象，而不是空数组
  }
})
</script>
<style scoped lang="scss">
/* 商品信息区域 */
.goods-info {
  background-color: #fff;
  padding: 15px;
}

.price-section {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.price {
  font-size: 24px;
  color: #ff4757;
  font-weight: bold;
}

.original-price {
  font-size: 14px;
  color: #999;
  text-decoration: line-through;
  margin-left: 10px;
}

.discount {
  font-size: 12px;
  color: #fff;
  background-color: #ff4757;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 10px;
}

.title-section {
  margin-bottom: 10px;
}

.goods-title {
  font-size: 16px;
  color: #333;
  line-height: 1.4;
}

.sales-info {
  display: flex;
  font-size: 12px;
  color: #999;
}

.sales, .stock {
  margin-right: 15px;
}

</style>
