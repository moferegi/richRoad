<template>
  <view class="ad-watch-btn" v-if="hasAds" @click="openAdPlayer">
    <text class="ad-btn-text">{{ adText }}</text>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getAdByPosition } from '@/api/ad.js'

const props = defineProps({
  positionKey: { type: String, required: true },
  text: { type: String, default: '看广告' }
})

const emit = defineEmits(['open'])
const hasAds = ref(false)
const adText = ref(props.text)
const adList = ref([])

const fetchAds = async () => {
  try {
    const res = await getAdByPosition(props.positionKey)
    if (res.code === 0 && res.data && res.data.length > 0) {
      adList.value = res.data
      hasAds.value = true
    } else {
      hasAds.value = false
    }
  } catch (e) {
    hasAds.value = false
  }
}

const openAdPlayer = () => {
  if (adList.value.length > 0) {
    emit('open', adList.value)
  }
}

onMounted(() => {
  fetchAds()
})
</script>

<style scoped>
.ad-watch-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 10rpx 24rpx;
  background: linear-gradient(135deg, #ff9500, #ff5e3a);
  border-radius: 30rpx;
  box-shadow: 0 4rpx 12rpx rgba(255, 94, 58, 0.4);
}
.ad-btn-text {
  font-size: 24rpx;
  color: #fff;
  font-weight: 600;
  white-space: nowrap;
}
</style>