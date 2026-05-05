<template>
  <view class="shoe-page">
    <view class="shoe-header">
      <text class="shoe-title">{{ $t('shoeRoom') }}</text>
    </view>

    <view class="shoe-body">
      <view class="shoe-left" @tap="pickImage('person')">
        <image v-if="personPreview" class="shoe-preview" :src="personPreview" mode="aspectFill" />
        <view v-else class="upload-empty">
          <uni-icons type="camera" size="26" color="rgba(15,23,42,0.45)" />
          <text class="upload-text">{{ $t('uploadPersonImage') }}</text>
        </view>
      </view>
      <view class="shoe-right">
        <view class="slot" @tap="pickImage('shoe')">
          <image v-if="shoePreview" class="slot-preview" :src="shoePreview" mode="aspectFill" />
          <text v-else class="slot-text">{{ $t('uploadShoeImage') }}</text>
        </view>
        <view class="slot suit-slot" @tap="goClothesPage">
          <text class="slot-text">{{ $t('suitSet') }}</text>
          <text class="slot-tip">{{ $t('goClothesPageTip') }}</text>
        </view>
      </view>
    </view>

    <view class="model-card">
      <view class="model-left" @tap="showModelPopup = true">
        <text class="model-label">{{ $t('currentModelLabel') }}</text>
        <text class="model-value">{{ currentModel.name }}</text>
      </view>
      <view class="model-help" @tap="showModelDesc">
        <text>?</text>
      </view>
    </view>

    <view class="action-btn" @tap="goGenerate">
      <text class="action-text">{{ $t('tryOnShoesWithCost').replace('{cost}', String(currentCost)) }}</text>
    </view>
    <view class="tips">{{ $t('shoeTryonTip') }}</view>

    <view class="popup-mask" v-if="showModelPopup" @tap="showModelPopup = false">
      <view class="popup-panel" @tap.stop>
        <view class="popup-title">{{ $t('selectModel') }}</view>
        <scroll-view scroll-y class="popup-list">
          <view
            v-for="item in modelList"
            :key="item.key"
            class="popup-item"
            :class="{ active: item.key === selectedModelKey }"
            @tap="selectModel(item.key)"
          >
            <view>
              <text class="popup-item-name">{{ item.name }}</text>
              <text class="popup-item-cost">{{ $t('pointsCostEach').replace('{cost}', String(item.cost)) }}</text>
            </view>
            <uni-icons type="checkmarkempty" size="20" color="#2563eb" v-if="item.key === selectedModelKey" />
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getTryonConfig } from '@/api/sysConfig.js'
import { getUrl } from '@/utils/url.js'
import { localText } from '@/utils/i18n.js'
import {
  createTryonRequestId,
  saveTryonDraft,
  getSelectedTryonModel,
  clearSelectedTryonModel,
  parseTryonModels,
} from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const showModelPopup = ref(false)
const selectedModelKey = ref('')
const modelList = ref([])
const tryonConfig = ref({
  tryon_cost_points: '1',
  tryon_models: '',
})

const personLocal = ref('')
const personRemote = ref('')
const shoeLocal = ref('')
const shoeRemote = ref('')

const personPreview = computed(() => personRemote.value ? getUrl(personRemote.value) : personLocal.value)
const shoePreview = computed(() => shoeRemote.value ? getUrl(shoeRemote.value) : shoeLocal.value)

const currentModel = computed(() => {
  const selected = modelList.value.find(v => v.key === selectedModelKey.value)
  return selected || modelList.value[0] || { key: 'aitryon', name: 'aitryon', cost: Number(tryonConfig.value.tryon_cost_points || 1), desc: {} }
})
const currentCost = computed(() => Number(currentModel.value.cost || 1))

const pickImage = (target) => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const path = res.tempFilePaths && res.tempFilePaths[0]
      if (!path) return
      if (target === 'person') {
        personLocal.value = path
        personRemote.value = ''
      }
      if (target === 'shoe') {
        shoeLocal.value = path
        shoeRemote.value = ''
      }
    },
  })
}

const selectModel = (key) => {
  selectedModelKey.value = key
  showModelPopup.value = false
}

const showModelDesc = () => {
  const model = currentModel.value
  const desc = model.descText || localText(model.desc, langStore.locale) || $t.value('noDescription')

  uni.showModal({
    title: model.name,
    content: desc,
    showCancel: false,
  })
}

const loadConfig = async () => {
  const res = await getTryonConfig()
  if (res.code === 0 && res.data) {
    tryonConfig.value = { ...tryonConfig.value, ...res.data }
  }

  modelList.value = parseTryonModels(
    tryonConfig.value.tryon_models,
    'shoes',
    Number(tryonConfig.value.tryon_cost_points || 1),
    langStore.locale
  )
  if (!modelList.value.find(item => item.key === selectedModelKey.value)) {
    selectedModelKey.value = modelList.value[0]?.key || ''
  }
}

const applySelectedModel = () => {
  const selectedModel = getSelectedTryonModel()
  if (selectedModel.roomType && selectedModel.roomType !== 'shoe') return

  if (selectedModel.remoteUrl) {
    personRemote.value = selectedModel.remoteUrl
    personLocal.value = ''
    clearSelectedTryonModel()
    return
  }

  if (selectedModel.localPath) {
    personLocal.value = selectedModel.localPath
    personRemote.value = ''
    clearSelectedTryonModel()
  }
}

const goClothesPage = () => {
  uni.switchTab({ url: '/pages/tabBar/clothes/index' })
}

const goGenerate = () => {
  if (!personLocal.value && !personRemote.value) {
    uni.showToast({ title: $t.value('uploadPersonFirst'), icon: 'none' })
    return
  }
  if (!shoeLocal.value && !shoeRemote.value) {
    uni.showToast({ title: $t.value('uploadShoeFirst'), icon: 'none' })
    return
  }

  saveTryonDraft({
    roomType: 'shoe',
    sceneType: 'shoes',
    operationType: 'tryon',
    sourceLocalPath: personLocal.value,
    sourceRemoteUrl: personRemote.value,
    templateLocalPath: shoeLocal.value,
    templateRemoteUrl: shoeRemote.value,
    modelKey: currentModel.value.key,
    modelName: currentModel.value.name,
    modelCost: currentCost.value,
    requestID: createTryonRequestId(),
  })

  uni.navigateTo({ url: '/pages/tryon/generate' })
}

onShow(() => {
  loadConfig()
  applySelectedModel()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.shoe-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 24rpx) 24rpx 24rpx;
  color: #0f172a;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.shoe-title {
  font-size: 38rpx;
  font-weight: 700;
  margin-bottom: 20rpx;
}

.shoe-body {
  display: flex;
  gap: 16rpx;
}

.shoe-left {
  width: 70%;
  height: 660rpx;
  border-radius: 20rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 12rpx 34rpx rgba(15, 23, 42, 0.06);
}

.shoe-right {
  width: 30%;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.slot {
  flex: 1;
  border-radius: 16rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  overflow: hidden;
}

.suit-slot {
  border-style: solid;
  border-color: rgba(37, 99, 235, 0.45);
}

.slot-tip {
  margin-top: 8rpx;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.55);
}

.slot-text {
  color: rgba(15, 23, 42, 0.82);
  font-size: 24rpx;
}

.shoe-preview,
.slot-preview {
  width: 100%;
  height: 100%;
}

.upload-empty {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 12rpx;
}

.upload-text {
  font-size: 24rpx;
  color: rgba(15, 23, 42, 0.6);
}

.model-card {
  margin-top: 18rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 16rpx;
  padding: 20rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 10rpx 28rpx rgba(15, 23, 42, 0.05);
}

.model-left {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.model-label {
  color: rgba(15, 23, 42, 0.55);
  font-size: 24rpx;
}

.model-value {
  color: #0f172a;
  font-size: 26rpx;
  font-weight: 600;
}

.model-help {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  border: 1rpx solid rgba(15, 23, 42, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #0f172a;
}

.action-btn {
  margin-top: 20rpx;
  height: 88rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 14rpx 30rpx rgba(37, 99, 235, 0.32);
}

.action-text {
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
}

.tips {
  margin-top: 14rpx;
  color: rgba(15, 23, 42, 0.55);
  font-size: 22rpx;
}

.popup-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.34);
  display: flex;
  justify-content: center;
  align-items: flex-end;
  z-index: 999;
}

.popup-panel {
  width: 100%;
  border-top-left-radius: 24rpx;
  border-top-right-radius: 24rpx;
  background: #ffffff;
  padding: 24rpx;
}

.popup-title {
  font-size: 30rpx;
  font-weight: 700;
  margin-bottom: 16rpx;
  color: #0f172a;
}

.popup-list {
  max-height: 520rpx;
}

.popup-item {
  border: 1rpx solid rgba(15, 23, 42, 0.08);
  border-radius: 14rpx;
  padding: 18rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
  background: #fff;
}

.popup-item.active {
  border-color: rgba(37, 99, 235, 0.45);
  background: rgba(219, 234, 254, 0.75);
}

.popup-item-name {
  display: block;
  font-size: 26rpx;
  margin-bottom: 6rpx;
  color: #0f172a;
}

.popup-item-cost {
  color: rgba(15, 23, 42, 0.55);
  font-size: 22rpx;
}
</style>
