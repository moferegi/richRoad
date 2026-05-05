<template>
  <view class="history-page">
    <view class="nav">
      <view class="nav-btn" @tap="goBack">
        <uni-icons type="left" size="20" color="#0f172a" />
      </view>
      <text class="nav-title">{{ $t('browseHistory') }}</text>
      <view class="nav-btn" @tap="loadHistory">
        <uni-icons type="reload" size="18" color="#0f172a" />
      </view>
    </view>

    <scroll-view class="filter-row" scroll-x :show-scrollbar="false">
      <view class="filter-chip" :class="{ active: currentFilter === 'all' }" @tap="currentFilter = 'all'">{{ $t('all') }}</view>
      <view class="filter-chip" :class="{ active: currentFilter === 'tryon' }" @tap="currentFilter = 'tryon'">{{ $t('tryonClothes') }}</view>
      <view class="filter-chip" :class="{ active: currentFilter === 'takeoff' }" @tap="currentFilter = 'takeoff'">{{ $t('takeoffAction') }}</view>
      <view class="filter-chip" :class="{ active: currentFilter === 'shoe' }" @tap="currentFilter = 'shoe'">{{ $t('tryonShoes') }}</view>
    </scroll-view>

    <scroll-view class="list-wrap" scroll-y>
      <view class="card" v-for="item in filteredList" :key="item._key" @tap="preview(item)">
        <view class="thumb-wrap">
          <image class="thumb" :src="item.resultImage ? getUrl(item.resultImage) : getUrl(item.templateImage || item.sourceImage)" mode="aspectFill" />
        </view>
        <view class="card-main">
          <view class="line">
            <text class="type">{{ typeLabel(item) }}</text>
            <text class="status" :class="statusClass(item.status)">{{ statusLabel(item.status) }}</text>
          </view>
          <text class="no">{{ $t('tryonTaskNo') }}：{{ item.taskNo || '-' }}</text>
          <text class="time">{{ formatTime(item.CreatedAt || item.savedAt) }}</text>
          <text class="error" v-if="item.status === 'failed' && item.errorMessage">{{ item.errorMessage }}</text>
        </view>
      </view>

      <view class="empty" v-if="filteredList.length === 0">
        <text>{{ $t('tryonNoTask') }}</text>
      </view>
      <view style="height: 30rpx"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getMyTryonTaskList } from '@/api/tryonTask.js'
import { getUrl } from '@/utils/url.js'
import { getTryonLocalHistory } from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const currentFilter = ref('all')
const allList = ref([])

const normalizeServerType = (item) => {
  if (item.sceneType === 'shoes') return 'shoe'
  return 'tryon'
}

const recordType = (item) => {
  if (item.roomType === 'shoe' || item.sceneType === 'shoes') return 'shoe'
  if (item.operationType === 'takeoff') return 'takeoff'
  return normalizeServerType(item)
}

const filteredList = computed(() => {
  if (currentFilter.value === 'all') return allList.value
  return allList.value.filter(item => recordType(item) === currentFilter.value)
})

const typeLabel = (item) => {
  const type = recordType(item)
  if (type === 'shoe') return $t.value('tryonShoes')
  if (type === 'takeoff') return $t.value('takeoffAction')
  return $t.value('tryonClothes')
}

const statusLabel = (status) => {
  if (status === 'success') return $t.value('tryonStatusSuccess')
  if (status === 'failed') return $t.value('tryonStatusFailed')
  return $t.value('tryonStatusProcessing')
}

const statusClass = (status) => {
  if (status === 'success') return 'success'
  if (status === 'failed') return 'failed'
  return 'processing'
}

const formatTime = (time) => {
  if (!time) return '-'
  const date = new Date(time)
  if (Number.isNaN(date.getTime())) return '-'
  const Y = date.getFullYear()
  const M = String(date.getMonth() + 1).padStart(2, '0')
  const D = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const m = String(date.getMinutes()).padStart(2, '0')
  return `${Y}-${M}-${D} ${h}:${m}`
}

const mergeHistory = (serverList, localList) => {
  const map = new Map()

  ;(localList || []).forEach((item, index) => {
    const key = item.taskNo || item.requestID || `local_${index}`
    map.set(key, {
      ...item,
      _key: key,
    })
  })

  ;(serverList || []).forEach((item) => {
    const key = item.taskNo || item.requestID || String(item.ID || Math.random())
    const old = map.get(key)
    map.set(key, {
      ...old,
      ...item,
      _key: key,
    })
  })

  return Array.from(map.values()).sort((a, b) => {
    const ta = new Date(a.CreatedAt || a.savedAt || 0).getTime()
    const tb = new Date(b.CreatedAt || b.savedAt || 0).getTime()
    return tb - ta
  })
}

const loadHistory = async () => {
  const token = uni.getStorageSync('x-token') || ''
  const localList = getTryonLocalHistory()

  if (!token) {
    allList.value = mergeHistory([], localList)
    return
  }

  const res = await getMyTryonTaskList({
    page: 1,
    pageSize: 50,
  })

  if (res.code === 0 && res.data) {
    allList.value = mergeHistory(res.data.list || [], localList)
    return
  }

  allList.value = mergeHistory([], localList)
}

const preview = (item) => {
  const url = item.resultImage || item.templateImage || item.sourceImage
  if (!url) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  uni.previewImage({ urls: [getUrl(url)] })
}

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

onShow(() => {
  loadHistory()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.history-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 14rpx) 20rpx 20rpx;
  color: #0f172a;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.nav {
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-btn {
  width: 90rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-title {
  font-size: 30rpx;
  font-weight: 700;
}

.filter-row {
  margin-top: 12rpx;
  display: flex;
  gap: 10rpx;
  overflow-x: auto;
  white-space: nowrap;
}

.filter-chip {
  min-width: 110rpx;
  height: 58rpx;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15,23,42,0.12);
  background: rgba(255,255,255,0.9);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
  color: rgba(15,23,42,0.62);
}

.filter-chip.active {
  color: #0f172a;
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(219, 234, 254, 0.82);
}

.list-wrap {
  margin-top: 14rpx;
  height: calc(100vh - var(--status-bar-height, 0px) - 170rpx);
}

.card {
  padding: 12rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(15,23,42,0.08);
  background: rgba(255,255,255,0.95);
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
  margin-bottom: 12rpx;
  display: flex;
  gap: 12rpx;
}

.thumb-wrap {
  width: 156rpx;
  height: 156rpx;
  border-radius: 10rpx;
  overflow: hidden;
}

.thumb {
  width: 100%;
  height: 100%;
}

.card-main {
  flex: 1;
  min-width: 0;
}

.line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8rpx;
}

.type {
  font-size: 24rpx;
  font-weight: 600;
}

.status {
  padding: 4rpx 12rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
}

.status.processing {
  background: rgba(255,205,105,0.2);
}

.status.success {
  background: rgba(86,199,102,0.22);
}

.status.failed {
  background: rgba(255,107,107,0.22);
}

.no,
.time,
.error {
  margin-top: 7rpx;
  display: block;
  font-size: 21rpx;
  color: rgba(15,23,42,0.62);
}

.error {
  color: #ffaaaa;
}

.empty {
  margin-top: 120rpx;
  text-align: center;
  color: rgba(15,23,42,0.55);
}
</style>
