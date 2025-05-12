<script setup>

</script>

<template>
  <view class="coupon-section section-card">
    <text class="section-title">优惠券</text>
    <view class="section-content">
      <view class="coupon-tag">领取优惠券</view>
      <uni-icons type="right" size="16" color="#999"></uni-icons>
    </view>
  </view>
</template>

<style scoped lang="scss">
/* 优惠券区域 */
.coupon-section {
  border-bottom: 1px solid #F5F5F5; // 保留 1px 边框线
}

.coupon-tag {
  font-size: 24rpx;
  color: #ff4757;
  border: 1px solid #ff4757; // 同样保留为细线描边
  padding: 4rpx 12rpx;
  border-radius: 8rpx;
}


</style>
