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
  padding: 36rpx;
  background-color: #fff;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 20%;
}

.category-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  margin-bottom: 16rpx;
}

.category-title {
  font-size: 24rpx;
  color: #333;
}

// 胶囊形状的促销横幅
.promo-banner {
  padding: 24rpx;
  background-color: #fff;
}

.promo-left {
  background: linear-gradient(90deg, #FF8C69, #FFD700);
  border-radius: 80rpx;
  height: 140rpx;
  padding: 0 40rpx;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.promo-content {
  display: flex;
  flex-direction: column;
}

.promo-text {
  font-size: 36rpx;
  font-weight: bold;
  color: white;
}

.promo-subtitle {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4rpx;
}

.promo-tag {
  position: absolute;
  left: 240rpx;
  top: 36rpx;
  background-color: #ff4a4a;
  color: white;
  font-size: 20rpx;
  padding: 0 12rpx;
  border-radius: 20rpx;
}
</style>

