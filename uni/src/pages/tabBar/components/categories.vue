<template>
  <scroll-view scroll-x class="nf-cate-bar" :scroll-left="scrollLeft" scroll-with-animation>
    <view class="nf-cate-tabs" id="nf-cate-tabs">
      <view
        v-for="item in tabList"
        :key="item._tabKey"
        class="nf-cate-tab nf-cate-tab-nf"
        :class="{ 'nf-cate-tab-active': activeIndex === item._rawIndex }"
        @tap="selectCategory(item, item._rawIndex)"
      >
        <text class="nf-cate-label-nf">{{ item._label || item.title }}</text>
        <view v-if="activeIndex === item._rawIndex" class="nf-cate-indicator"></view>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
import { ref, computed, watch, nextTick, getCurrentInstance } from 'vue'
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
const scrollLeft = ref(0)
const instance = getCurrentInstance()

const tabList = computed(() => {
  const allItem = { ID: 0, title: '', icons: '', _rawIndex: 0, _label: $t.value('all'), _tabKey: 'all' }
  // 分类 tab 只用于展示和选择；稳定 key 避免分类顺序变化时继续用索引复用节点。
  return [allItem, ...props.categoriesData.map((item, index) => ({
    ...item,
    _rawIndex: index + 1,
    _label: $lt.value(item.title),
    _tabKey: `${item?.ID || item?.id || item?.title || 'category'}-${index}`
  }))]
})

const selectCategory = (item, index) => {
  activeIndex.value = index
  emit('update:modelValue', item.ID || 0)
  emit('change', item.ID || 0)
  scrollToCenter(index)
}

// 将选中的tab滚动到中心位置
const scrollToCenter = (index) => {
  nextTick(() => {
    const query = uni.createSelectorQuery().in(instance.proxy)
    // 获取所有tab和scroll容器的宽度信息
    query.select('.nf-cate-bar').boundingClientRect()
    query.selectAll('.nf-cate-tab').boundingClientRect()
    query.exec((res) => {
      if (!res || !res[0] || !res[1] || !res[1][index]) return
      const barWidth = res[0].width
      const tabs = res[1]
      const activeTab = tabs[index]
      // 计算目标tab中心相对于第一个tab左边的偏移
      const tabCenter = activeTab.left - tabs[0].left + activeTab.width / 2
      // 滚动使其居中
      scrollLeft.value = tabCenter - barWidth / 2
    })
  })
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
