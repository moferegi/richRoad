<template>
  <view class="spec-section section-card" @click="openPopup">
    <text class="section-title">规格</text>
    <view class="section-content">
      <text>{{ selectedText }}</text>
      <uni-icons type="right" size="16" color="#999"></uni-icons>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue';
const props = defineProps({
  selectedSku: {
    type: Object,
    default: null
  }
});

// 修正这里的emit定义，使用正确的事件名称
const emit = defineEmits(['open-sku']);

// 计算显示文本
const selectedText = computed(() => {
  if (!props.selectedSku) return '请选择规格';

  const attrs = props.selectedSku.sku_attrs;
  return Object.entries(attrs)
      .map(([key, value]) => `${key}:${typeof value === 'object' ? value.name : value}`)
      .join(' ');
});

// 打开SKU选择弹窗
const openPopup = () => {
  emit('open-sku');
};
</script>

<style scoped lang="scss">
.spec-section{
  border-bottom: 1px solid #F5F5F5;
}
</style>
