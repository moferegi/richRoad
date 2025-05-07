<template>
  <view class="category-section">
    <view class="category-item" v-for="(item, index) in categoriesData" :key="index">
      <view class="category-icon" :style="{ backgroundColor: item.color }"></view>
      <text class="category-title">{{ item.title }}</text>
    </view>
  </view>

  <!-- 促销横幅 -->
  <view class="promo-banner">
    <view class="promo-left">
      <view class="promo-content">
        <text class="promo-text">生鲜果蔬 5折起</text>
        <text class="promo-subtitle">端午乐三天 立刻折上折</text>
      </view>
      <view class="promo-image-container">
        <view class="promo-image-left"></view>
        <view class="promo-image-right"></view>
      </view>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  categoriesData: {
    type: Array,
    default: () => [] // 修改为函数返回空数组
  }
})
</script>

<style scoped lang="scss">
// 分类导航
.category-section {
  display: flex;
  justify-content: space-around;
  padding: 18px;
  background-color: #fff;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 20%;
}

.category-icon {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  margin-bottom: 8px;
}

.category-title {
  font-size: 12px;
  color: #333;
}

</style>
