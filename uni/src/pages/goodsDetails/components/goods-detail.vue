<template>
  <view class="nf-intro-box" v-if="detail">
    <rich-text :nodes="detail"></rich-text>
  </view>
</template>

<script setup>
const props = defineProps({
  detail: {
    type: String,
    default: ''
  }
})
</script>

<style lang="scss" scoped>
.nf-intro-box {
  width: 100%;
  padding: 20rpx 0;
  box-sizing: border-box;
  color: rgba(255, 255, 255, 0.7);
  font-size: 28rpx;
  line-height: 1.6;
}
</style>
