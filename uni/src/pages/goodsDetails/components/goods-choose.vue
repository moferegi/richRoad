<template>
  <view class="spec-section section-card" @click="openPopup">
    <text class="section-title">规格</text>
    <view class="section-content">
      <text>请选择规格</text>
      <uni-icons type="right" size="16" color="#999"></uni-icons>
    </view>
  </view>
</template>

<script setup>
const openPopup = () => {
  // 这里可以调用父组件的方法来打开弹窗
  // 例如: emit('openPopup')
  console.log('打开规格选择弹窗');
}
</script>

<style scoped lang="scss">
.spec-section{
  border-bottom: 1px solid #F5F5F5;
}
</style>
