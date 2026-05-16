<template>
  <view v-if="visible" class="ias-mask" @tap="handleCancel" @touchmove.stop.prevent>
    <view class="ias-panel" @tap.stop>
      <view v-if="title" class="ias-title">{{ title }}</view>
      <view class="ias-list">
        <view
          v-for="(item, index) in normalizedItems"
          :key="`item-${index}`"
          class="ias-item"
          @tap="handleSelect(index)"
        >
          <text class="ias-item-text">{{ item }}</text>
        </view>
      </view>
      <view class="ias-gap"></view>
      <view class="ias-cancel" @tap="handleCancel">
        <text class="ias-cancel-text">{{ cancelText }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '',
  },
  items: {
    type: Array,
    default: () => [],
  },
  cancelText: {
    type: String,
    default: 'Cancel',
  },
})

const emit = defineEmits(['update:visible', 'select', 'cancel'])

const normalizedItems = computed(() => {
  return (Array.isArray(props.items) ? props.items : []).map((item) => {
    if (typeof item === 'string') {
      return item
    }
    if (item && typeof item === 'object' && Object.prototype.hasOwnProperty.call(item, 'label')) {
      return String(item.label || '')
    }
    return String(item || '')
  })
})

const closeSheet = () => {
  emit('update:visible', false)
}

const handleCancel = () => {
  closeSheet()
  emit('cancel')
}

const handleSelect = (index) => {
  closeSheet()
  emit('select', { index })
}
</script>

<style scoped>
.ias-mask {
  position: fixed;
  inset: 0;
  z-index: 1200;
  background: rgba(0, 0, 0, 0.28);
  display: flex;
  align-items: flex-end;
}

.ias-panel {
  width: 100%;
  padding: 0 20rpx calc(env(safe-area-inset-bottom, 0px) + 20rpx);
  box-sizing: border-box;
}

.ias-title {
  padding: 24rpx 20rpx;
  text-align: center;
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.58);
  background: #ffffff;
  border-radius: 22rpx 22rpx 0 0;
}

.ias-list {
  background: #ffffff;
  border-radius: 22rpx;
  overflow: hidden;
}

.ias-title + .ias-list {
  border-radius: 0 0 22rpx 22rpx;
}

.ias-item {
  min-height: 96rpx;
  padding: 0 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1rpx solid rgba(15, 23, 42, 0.08);
}

.ias-item:last-child {
  border-bottom: none;
}

.ias-item-text {
  font-size: 30rpx;
  color: #0f172a;
  text-align: center;
}

.ias-gap {
  height: 14rpx;
}

.ias-cancel {
  min-height: 96rpx;
  background: #ffffff;
  border-radius: 22rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ias-cancel-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #0f172a;
}
</style>
