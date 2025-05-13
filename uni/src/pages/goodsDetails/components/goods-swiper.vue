<template>
  <view style="padding: 24rpx 0 20rpx;" class="bgc_fff">
    <swiper class="swiper" circular indicator-active-color="#fff" indicator-color="rgba(255,255,255,0.6)"
            :indicator-dots="true" :autoplay="true" :interval="3000" :duration="1000">
      <swiper-item class="swiper-item flexr-jsc flex-aic" v-for="i in list">
        <image
            mode="aspectFill" class="swiper-item-img"
            :src="getUrl(i)">
        </image>
      </swiper-item>
      <!-- <swiper-item class="swiper-item flexr-jsc flex-aic">
        <image mode="aspectFill" class="swiper-item-img"
          src="https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fsafe-img.xhscdn.com%2Fbw1%2Fb103f969-9675-495b-95ff-87ab7e0d515d%3FimageView2%2F2%2Fw%2F1080%2Fformat%2Fjpg&refer=http%3A%2F%2Fsafe-img.xhscdn.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1715314737&t=00a913e3dc1ec77e9be7ae9660a27997">
        </image>
      </swiper-item> -->
    </swiper>
  </view>
</template>

<script setup>
import {getUrl} from "@/utils/url.js"
const props = defineProps({
  list: {
    type: String,
    default: ''
  }
})
</script>

<style lang="scss" scoped>
.swiper {
  height: 560rpx;
  .swiper-item {
    width: 100%;
    height: 100%;
  }
  .swiper-item-img {
    width: 686rpx;
    border-radius: 10rpx;
    height: 100%;
  }
}
</style>
