<script setup>

</script>

<template>
  <view class="detail-section">
    <view class="detail-title">图文详情</view>
    <view class="detail-content">
      <view class="detail-image" :style="{ backgroundColor: '#f5f5f5', height: '200px' }"></view>
      <view class="detail-image" :style="{ backgroundColor: '#f0f0f0', height: '300px' }"></view>
    </view>
  </view>
</template>

<style scoped lang="scss">
/* 图文详情 */
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
</style>
