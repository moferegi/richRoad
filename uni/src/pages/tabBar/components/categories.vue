<template>
  <scroll-view scroll-x class="nf-cate-bar">
    <view class="nf-cate-tabs">
      <view
        v-for="(item, index) in tabList"
        :key="index"
        class="nf-cate-tab nf-cate-tab-nf"
        :class="{ 'nf-cate-tab-active': activeIndex === index }"
        @tap="selectCategory(item, index)"
      >
        <text class="nf-cate-label-nf">{{ item._label || item.title }}</text>
        <view v-if="activeIndex === index" class="nf-cate-indicator"></view>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { getUrl } from '@/utils/url'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const props = defineProps({
  categoriesData: {
    type: Array,
    default: () => []
  },
  modelValue: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const activeIndex = ref(0)

const tabList = computed(() => {
  const allItem = { ID: 0, title: '', icons: '', _label: $t.value('all') }
  return [allItem, ...props.categoriesData.map(item => ({ ...item, _label: $lt.value(item.title) }))]
})

const selectCategory = (item, index) => {
  activeIndex.value = index
  emit('update:modelValue', item.ID || 0)
  emit('change', item.ID || 0)
}

watch(() => props.modelValue, (val) => {
  const idx = tabList.value.findIndex(t => (t.ID || 0) === val)
  if (idx >= 0) activeIndex.value = idx
})
</script>

<style scoped lang="scss">
.nf-cate-bar {
  width: 100%;
  white-space: nowrap;
  background: #000;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.04);
}

.nf-cate-tabs {
  display: flex;
  padding: 16rpx 12rpx 8rpx;
  gap: 6rpx;
}

.nf-cate-tab {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  min-width: 110rpx;
  height: 72rpx;
  padding: 0 16rpx;
  border-radius: 20rpx;
  position: relative;
  transition: all 0.25s ease;
  border: 1rpx solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.05);
}

.nf-cate-tab:active {
  transform: scale(0.95);
}

.nf-cate-tab-active {
  background: linear-gradient(135deg, rgba(229, 9, 20, 0.95), rgba(224, 0, 19, 0.75));
  border-color: #e50914;
  transform: scale(1.03);
}

.nf-cate-tab-active .nf-cate-label-nf {
  color: #fff;
  font-weight: 800;
}


.nf-cate-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  margin-bottom: 10rpx;
  border: 2rpx solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  transition: all 0.35s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.nf-cate-icon-all {
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.06);
}

.nf-cate-icon-all-text {
  font-size: 36rpx;
  color: rgba(229, 9, 20, 0.7);
  font-weight: 300;
}

.nf-cate-label-nf {
  font-size: 28rpx;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180rpx;
  letter-spacing: 1.2rpx;
  text-transform: capitalize;
  transition: color 0.2s ease;
}

.nf-cate-tab-active .nf-cate-label-nf {
  color: #fff;
}


.nf-cate-indicator {
  position: absolute;
  bottom: 4rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 32rpx;
  height: 4rpx;
  border-radius: 2rpx;
  background: linear-gradient(90deg, #e50914, #ff4d58);
  box-shadow: 0 0 10rpx rgba(229, 9, 20, 0.5);
  animation: nf-indicator-in 0.3s ease-out;
}

@keyframes nf-indicator-in {
  0% {
    width: 0;
    opacity: 0;
  }
  100% {
    width: 32rpx;
    opacity: 1;
  }
}
</style>
