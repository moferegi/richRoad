<script setup>

</script>

<template>
  <view class="promotion-section section-card">
    <text class="section-title">促销</text>
    <view class="promotion-list">
      <view class="promotion-item">
        <text class="promotion-tag">新人</text>
        <text class="promotion-text">新人首单减20元无门槛代金券</text>
      </view>
      <view class="promotion-item">
        <text class="promotion-tag">满减</text>
        <text class="promotion-text">订单满50减10</text>
      </view>
      <view class="promotion-item">
        <text class="promotion-tag">满减</text>
        <text class="promotion-text">订单满100减30</text>
      </view>
      <view class="promotion-item">
        <text class="promotion-tag">包邮</text>
        <text class="promotion-text">单笔满两件包邮</text>
      </view>
    </view>
  </view>
</template>

<style scoped lang="scss">
/* 促销活动区域 */
.promotion-section {
  flex-direction: column;
  align-items: flex-start;
  border-bottom: 1px solid #F5F5F5;
}

.promotion-list {
  width: 100%;
  margin-top: 10px;
}

.promotion-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.promotion-tag {
  font-size: 12px;
  color: #fff;
  background-color: #ff4757;
  padding: 1px 4px;
  border-radius: 2px;
  margin-right: 8px;
}

.promotion-text {
  font-size: 12px;
  color: #333;
}
</style>
