<template>
  <view class="test" style="padding: 3rem 1.5rem;">
    <wu-sku
        v-model="show"
        :data="skuData.skus"
        :themeColor="[226, 35, 26]"
        :defaultCover="defaultCover"
        :btnConfirmText="confirmButtonText"
        notSelectSku="请选择完整的商品信息"
        @skuChange="handleSkuChange"
        @confirm="handleConfirm"
    ></wu-sku>
  </view>
</template>

<script setup>
import {ref, computed, watch} from 'vue';

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  skuData: {
    type: Object,
    required: true
  },
  actionType: {
    type: String,
    default: 'buy'
  }
});

const emit = defineEmits(['update:modelValue', 'sku-confirm']);

// 创建一个计算属性来双向绑定 show
const show = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
});

// 默认封面图
const defaultCover = 'https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/default.jpeg';

// 确认按钮文本
const confirmButtonText = computed(() => {
  return props.actionType === 'cart' ? '加入购物车' : '立即购买';
});

// SKU变化处理
const handleSkuChange = (sku) => {
  console.log('SKU变化', sku);
};

// 确认选择
const handleConfirm = (sku) => {
  emit('sku-confirm', sku);
};
</script>
