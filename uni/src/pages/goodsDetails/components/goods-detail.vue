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
  color: #111827;
  font-size: 28rpx;
  line-height: 1.6;
}

.nf-intro-box :deep(p),
.nf-intro-box :deep(div),
.nf-intro-box :deep(span),
.nf-intro-box :deep(li),
.nf-intro-box :deep(td),
.nf-intro-box :deep(th) {
  color: inherit;
}
</style>
