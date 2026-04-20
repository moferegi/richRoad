<template>
  <view v-if="false" class="font_28 color_333 goods_item_box bgc_fff">
    {{ $t('noReviewYet') }}
  </view>
  <view v-else class="bgc_fff goods_item_box">
    <view class="flex-aic flexr-jsb m_b_24">
      <text class="font_28 color_333">{{ $t('userReview') }}</text>
      <view class="flex flex-aic" @tap="toEvaluate">
        <text class="font_24 color_fe5572">{{ $t('viewAllReview') }}</text>
        <image class="goods_item_icon" src="http://www.liwanying.top/applate-icon/gengduohuangse.png" mode=""></image>
      </view>
    </view>
  </view>
</template>

<script setup>
import { useLangStore } from '@/pinia/modules/lang.js'
import { computed } from 'vue'
const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const toEvaluate = () => {
  uni.navigateTo({
    url: '/pages/evaluate/evaluate'
  })
}
</script>

<style lang="scss" scoped>
.goods_item_box {
  padding: 20rpx 32rpx;
}
.goods_item_icon {
  width: 12rpx;
  height: 24rpx;
  margin-left: 12rpx;
}
</style>
