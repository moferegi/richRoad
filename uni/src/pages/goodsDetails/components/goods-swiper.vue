<script setup>

</script>

<template>
  <swiper class="goods-swiper" circular indicator-dots autoplay :interval="3000" :duration="500" indicator-active-color="#ff4757">
    <swiper-item v-for="(item, index) in 3" :key="index">
      <view class="swiper-item" :style="{ backgroundColor: 'pink' }"></view>
    </swiper-item>
  </swiper>
</template>

<style scoped lang="scss">
/* 商品轮播图 */
.goods-swiper {
  width: 100%;
  height: 832rpx; // 原 26rem
}

.swiper-item {
  width: 100%;
  height: 100%;
}
</style>
