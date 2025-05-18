<template>
  <swiper
      class="goods-swiper"
      circular
      :indicator-dots="true"
      :autoplay="true"
      :interval="3000"
      :duration="500"
      indicator-active-color="#ff4757">
    <swiper-item v-for="(i, index) in list" :key="index" class="swiper-item">
      <image
          mode="aspectFill"
          class="swiper-item-img"
          :src="getUrl(i)">
      </image>
    </swiper-item>
  </swiper>
</template>

<script setup>
import {getUrl} from "@/utils/url.js"

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  }
})
</script>

<style lang="scss" scoped>
.goods-swiper {
  width: 100%;
  height: 26rem;
}

.swiper-item {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.swiper-item-img {
  width: 100%;
  height: 100%;
  object-fit: cover; /* 确保图片覆盖整个容器同时保持比例 */
}
</style>
