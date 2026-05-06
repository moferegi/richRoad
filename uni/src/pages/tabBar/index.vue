<template>
  <view class="room-page">
    <view class="room-header">
      <text class="room-title">{{ $t('tryonRoom') }}</text>
      <view class="room-tabs">
        <view class="room-tab active">{{ $t('tryonArea') }}</view>
      </view>
    </view>

    <view class="room-body">
      <view class="room-left" @tap="openUploadDrawer('person')">
        <image v-if="personPreview" class="room-preview" :src="personPreview" mode="aspectFill" />
        <view v-else class="room-upload-empty">
          <uni-icons type="camera" size="26" color="rgba(15,23,42,0.45)" />
          <text class="room-upload-text">{{ $t('uploadModelImage') }}</text>
        </view>
      </view>
      <view class="room-right">
        <view class="cloth-slot" @tap="openUploadDrawer('upper')">
          <image v-if="upperPreview" class="cloth-preview" :src="upperPreview" mode="aspectFill" />
          <text v-else class="cloth-text">{{ $t('uploadUpperImage') }}</text>
        </view>
        <view class="cloth-slot" @tap="openUploadDrawer('lower')">
          <image v-if="lowerPreview" class="cloth-preview" :src="lowerPreview" mode="aspectFill" />
          <text v-else class="cloth-text">{{ $t('uploadLowerImage') }}</text>
        </view>
        <view class="cloth-slot suit-slot" @tap="goClothesPage">
          <text class="cloth-text">{{ $t('suitSet') }}</text>
          <text class="cloth-tip">{{ $t('goClothesPageTip') }}</text>
        </view>
      </view>
    </view>

    <view class="model-card">
      <view class="model-left" @tap="showModelPopup = true">
        <text class="model-label">{{ $t('currentModelLabel') }}</text>
        <view class="model-info">
          <text class="model-value">{{ currentModel.name }}</text>
          <text class="model-cost-hint">{{ $t('pointsCostEach').replace('{cost}', String(currentCost)) }}</text>
          <text class="model-desc-hint" v-if="currentModelDesc">{{ currentModelDesc }}</text>
        </view>
      </view>
      <view class="model-help" @tap="showModelDesc">
        <text>?</text>
      </view>
    </view>

    <view class="generate-btn" @tap="goGenerate">
      <text class="generate-text">{{ $t('generateWithCost').replace('{cost}', String(currentCost)) }}</text>
    </view>

    <view class="footer-tips">
      <text>{{ $t('tryonFootTip') }}</text>
    </view>

    <view class="popup-mask" v-if="showModelPopup" @tap="showModelPopup = false">
      <view class="popup-panel" @tap.stop>
        <view class="popup-title">{{ $t('selectModel') }}</view>
        <scroll-view class="popup-list" scroll-y>
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
              <text class="popup-item-desc" v-if="item.descText">{{ item.descText }}</text>
            </view>
            <uni-icons type="checkmarkempty" size="20" color="#2563eb" v-if="item.key === selectedModelKey" />
          </view>
        </scroll-view>
      </view>
    </view>

    <view class="upload-mask" v-if="showUploadDrawer" @tap="closeUploadDrawer">
      <view class="upload-drawer" @tap.stop>
        <view class="upload-drawer-head">
          <text class="upload-drawer-title">{{ uploadDrawerTitle }}</text>
          <text class="upload-drawer-tip">{{ uploadDrawerTip }}</text>
        </view>

        <view class="upload-drawer-actions">
          <view class="upload-drawer-btn" @tap="chooseFromAlbum">相册上传</view>
          <view class="upload-drawer-btn" @tap="chooseFromCamera">拍照上传</view>
          <view class="upload-drawer-btn ghost" @tap="clearUploadTarget">清空当前</view>
        </view>

        <view class="upload-example-head">
          <text class="upload-example-title">示例图（可一键套用）</text>
          <text class="upload-example-note">建议使用清晰正面、纯背景图片，减少生成失败。</text>
        </view>

        <scroll-view class="upload-example-list" scroll-x>
          <view
            v-for="item in uploadExamples"
            :key="item.url"
            class="upload-example-item"
            @tap="applyUploadExample(item)"
          >
            <image class="upload-example-image" :src="getExamplePreview(item.url)" mode="aspectFill" />
            <text class="upload-example-label">{{ item.label }}</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { getTryonConfig } from '@/api/sysConfig.js'
import { getUrl } from '@/utils/url.js'
import { localText } from '@/utils/i18n.js'
import {
  createTryonRequestId,
  saveTryonDraft,
  getSelectedClothes,
  clearSelectedClothes,
  getSelectedTryonModel,
  clearSelectedTryonModel,
  parseTryonModels,
} from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const showModelPopup = ref(false)
const showUploadDrawer = ref(false)
const uploadTarget = ref('')
const selectedModelKey = ref('')
const modelList = ref([])
const tryonConfig = ref({
  tryon_cost_points: '1',
  tryon_models: '',
})

const personLocal = ref('')
const personRemote = ref('')
const upperLocal = ref('')
const upperRemote = ref('')
const lowerLocal = ref('')
const lowerRemote = ref('')

const personPreview = computed(() => personRemote.value ? getUrl(personRemote.value) : personLocal.value)
const upperPreview = computed(() => upperRemote.value ? getUrl(upperRemote.value) : upperLocal.value)
const lowerPreview = computed(() => lowerRemote.value ? getUrl(lowerRemote.value) : lowerLocal.value)

const activeSceneType = computed(() => 'clothes')

const currentModel = computed(() => {
  const selected = modelList.value.find(v => v.key === selectedModelKey.value)
  return selected || modelList.value[0] || { key: 'aitryon', name: 'aitryon', cost: Number(tryonConfig.value.tryon_cost_points || 1), desc: {} }
})

const currentCost = computed(() => Number(currentModel.value.cost || 1))
const currentModelDesc = computed(() => currentModel.value.descText || localText(currentModel.value.desc, langStore.locale) || '')

const uploadExamplesMap = {
  person: [
    { label: '模特示例A', url: '/uploads/file/f46124e7c57a2d3fbbe5ec45cc66201b_20260403092613.PNG' },
    { label: '模特示例B', url: '/uploads/file/01fd5aeddc82f4f584fc9b4927e4ce2b_20260403092813.JPEG' },
  ],
  upper: [
    { label: '上衣示例A', url: '/uploads/file/233e54e186d06ded7c4cf2b3c68b97b7_20260403094625.JPEG' },
    { label: '上衣示例B', url: '/uploads/file/6290204ab51e001c89b4063a0821afee_20260403094625.JPEG' },
  ],
  lower: [
    { label: '下装示例A', url: '/uploads/file/2b6a999cefac3fe3d282fc8bf96ae32b_20260325140212.jpg' },
    { label: '下装示例B', url: '/uploads/file/44e2f6cc243a32f0854caa9019fe99b3_20260325140114.jpg' },
  ],
}

const uploadSceneKey = computed(() => {
  return uploadTarget.value || 'person'
})

const uploadDrawerTitle = computed(() => {
  const titleMap = {
    person: '上传模特图',
    upper: '上传上装图',
    lower: '上传下装图',
  }
  return titleMap[uploadSceneKey.value] || '上传图片'
})

const uploadDrawerTip = computed(() => {
  const tipMap = {
    person: '上传目录：tryon/clothes/source',
    upper: '上传目录：tryon/clothes/template/upper',
    lower: '上传目录：tryon/clothes/template/lower',
  }
  return tipMap[uploadSceneKey.value] || '上传目录：tryon/clothes/source'
})

const uploadExamples = computed(() => uploadExamplesMap[uploadSceneKey.value] || [])

const rebuildModelList = () => {
  modelList.value = parseTryonModels(
    tryonConfig.value.tryon_models,
    activeSceneType.value,
    Number(tryonConfig.value.tryon_cost_points || 1),
    langStore.locale
  )
  if (!modelList.value.find(item => item.key === selectedModelKey.value)) {
    selectedModelKey.value = modelList.value[0]?.key || ''
  }
}

const assignImageToTarget = (target, value, isRemote = false) => {
  const localPath = isRemote ? '' : value
  const remotePath = isRemote ? value : ''

  if (target === 'person') {
    personLocal.value = localPath
    personRemote.value = remotePath
  }
  if (target === 'upper') {
    upperLocal.value = localPath
    upperRemote.value = remotePath
  }
  if (target === 'lower') {
    lowerLocal.value = localPath
    lowerRemote.value = remotePath
  }
}

const chooseLocalImage = (target, sourceType) => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType,
    success: (res) => {
      const path = res.tempFilePaths && res.tempFilePaths[0]
      if (!path) return
      assignImageToTarget(target, path, false)
      showUploadDrawer.value = false
    },
  })
}

const openUploadDrawer = (target) => {
  uploadTarget.value = target
  showUploadDrawer.value = true
}

const closeUploadDrawer = () => {
  showUploadDrawer.value = false
}

const chooseFromAlbum = () => {
  if (!uploadTarget.value) return
  chooseLocalImage(uploadTarget.value, ['album'])
}

const chooseFromCamera = () => {
  if (!uploadTarget.value) return
  chooseLocalImage(uploadTarget.value, ['camera'])
}

const clearUploadTarget = () => {
  if (!uploadTarget.value) return
  assignImageToTarget(uploadTarget.value, '', false)
  showUploadDrawer.value = false
}

const applyUploadExample = (item) => {
  if (!uploadTarget.value || !item?.url) return
  assignImageToTarget(uploadTarget.value, item.url, true)
  showUploadDrawer.value = false
  uni.showToast({ title: '已套用示例图', icon: 'none' })
}

const getExamplePreview = (url) => {
  const value = String(url || '').trim()
  if (!value) return ''
  if (value.startsWith('http://') || value.startsWith('https://')) {
    return value
  }
  return getUrl(value)
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
  rebuildModelList()
}

const applySelectedClothes = () => {
  const selected = getSelectedClothes()
  if (!selected || typeof selected !== 'object') return
  if (selected.upperImage) {
    upperRemote.value = selected.upperImage
    upperLocal.value = ''
  }
  if (selected.lowerImage) {
    lowerRemote.value = selected.lowerImage
    lowerLocal.value = ''
  }
  clearSelectedClothes()
}

const applySelectedModel = () => {
  const selectedModel = getSelectedTryonModel()
  if (!selectedModel || typeof selectedModel !== 'object') return
  if (selectedModel.roomType && selectedModel.roomType !== 'tryon') return

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
    uni.showToast({ title: $t.value('uploadModelFirst'), icon: 'none' })
    return
  }

  let templateRemoteUrl = upperRemote.value || lowerRemote.value
  let templateLocalPath = upperLocal.value || lowerLocal.value

  if (!templateRemoteUrl && !templateLocalPath) {
    uni.showToast({ title: $t.value('uploadClothesFirst'), icon: 'none' })
    return
  }

  const templateSegment = upperRemote.value || upperLocal.value ? 'upper' : 'lower'
  const sourceUploadFolder = 'tryon/clothes/source'
  const templateUploadFolder = `tryon/clothes/template/${templateSegment}`

  saveTryonDraft({
    roomType: 'tryon',
    sceneType: 'clothes',
    operationType: 'tryon',
    sourceLocalPath: personLocal.value,
    sourceRemoteUrl: personRemote.value,
    sourceUploadFolder,
    templateLocalPath,
    templateRemoteUrl,
    templateUploadFolder,
    modelKey: currentModel.value.key,
    modelName: currentModel.value.name,
    modelCost: currentCost.value,
    requestID: createTryonRequestId(),
  })

  uni.navigateTo({ url: '/pages/tryon/generate' })
}

watch(
  [() => tryonConfig.value.tryon_models, () => tryonConfig.value.tryon_cost_points, () => langStore.locale],
  () => {
    rebuildModelList()
  }
)

onShow(() => {
  loadConfig()
  applySelectedModel()
  applySelectedClothes()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.room-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 24rpx) 24rpx 24rpx;
  color: #0f172a;
  background: radial-gradient(120% 80% at 100% -10%, #dbeafe 0%, transparent 60%), #f4f7fb;
}

.room-header {
  margin-bottom: 20rpx;
}

.room-title {
  font-size: 38rpx;
  font-weight: 700;
}

.room-tabs {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
}

.room-tab {
  flex: 1;
  text-align: center;
  padding: 16rpx 0;
  border-radius: 999rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.12);
  color: rgba(15, 23, 42, 0.55);
  background: rgba(255, 255, 255, 0.75);
}

.room-tab.active {
  color: #0f172a;
  border-color: rgba(37, 99, 235, 0.45);
  background: linear-gradient(90deg, rgba(191, 219, 254, 0.95), rgba(219, 234, 254, 0.95));
}

.room-body {
  display: flex;
  gap: 16rpx;
}

.room-left {
  width: 70%;
  height: 660rpx;
  border-radius: 20rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 12rpx 36rpx rgba(15, 23, 42, 0.06);
}

.room-right {
  width: 30%;
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.cloth-slot {
  flex: 1;
  border-radius: 16rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  background: rgba(255, 255, 255, 0.92);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 10rpx 26rpx rgba(15, 23, 42, 0.05);
}

.suit-slot {
  border-style: solid;
  border-color: rgba(37, 99, 235, 0.45);
}

.cloth-tip {
  margin-top: 8rpx;
  font-size: 20rpx;
  color: rgba(15, 23, 42, 0.55);
}

.cloth-text {
  color: rgba(15, 23, 42, 0.82);
  font-size: 24rpx;
}

.room-preview,
.takeoff-preview,
.cloth-preview {
  width: 100%;
  height: 100%;
}

.takeoff-body {
  height: 660rpx;
  border-radius: 20rpx;
  border: 1rpx dashed rgba(15, 23, 42, 0.2);
  overflow: hidden;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 12rpx 36rpx rgba(15, 23, 42, 0.06);
}

.room-upload-empty {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 12rpx;
}

.room-upload-text {
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
  align-items: flex-start;
  gap: 8rpx;
  flex: 1;
}

.model-info {
  min-width: 0;
}

.model-label {
  color: rgba(15, 23, 42, 0.55);
  font-size: 24rpx;
  margin-top: 2rpx;
}

.model-value {
  display: block;
  color: #0f172a;
  font-size: 26rpx;
  font-weight: 600;
}

.model-cost-hint {
  display: block;
  margin-top: 4rpx;
  font-size: 21rpx;
  color: rgba(37, 99, 235, 0.9);
}

.model-desc-hint {
  display: -webkit-box;
  margin-top: 4rpx;
  font-size: 21rpx;
  color: rgba(15, 23, 42, 0.6);
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
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

.generate-btn {
  margin-top: 20rpx;
  height: 88rpx;
  border-radius: 999rpx;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 14rpx 30rpx rgba(37, 99, 235, 0.32);
}

.generate-text {
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
}

.footer-tips {
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
  display: block;
  color: rgba(15, 23, 42, 0.55);
  font-size: 22rpx;
}

.popup-item-desc {
  display: -webkit-box;
  margin-top: 6rpx;
  max-width: 460rpx;
  color: rgba(15, 23, 42, 0.56);
  font-size: 20rpx;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.upload-mask {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.34);
  display: flex;
  justify-content: center;
  align-items: flex-end;
  z-index: 1000;
}

.upload-drawer {
  width: 100%;
  background: #ffffff;
  border-top-left-radius: 28rpx;
  border-top-right-radius: 28rpx;
  padding: 24rpx;
}

.upload-drawer-head {
  margin-bottom: 16rpx;
}

.upload-drawer-title {
  display: block;
  font-size: 30rpx;
  font-weight: 700;
  color: #0f172a;
}

.upload-drawer-tip {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.56);
}

.upload-drawer-actions {
  display: flex;
  gap: 12rpx;
  margin-bottom: 16rpx;
}

.upload-drawer-btn {
  flex: 1;
  text-align: center;
  border-radius: 999rpx;
  padding: 14rpx 10rpx;
  font-size: 24rpx;
  color: #fff;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.upload-drawer-btn.ghost {
  color: rgba(15, 23, 42, 0.76);
  background: rgba(15, 23, 42, 0.08);
}

.upload-example-head {
  margin-bottom: 12rpx;
}

.upload-example-title {
  display: block;
  font-size: 26rpx;
  font-weight: 600;
  color: #0f172a;
}

.upload-example-note {
  display: block;
  margin-top: 6rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.52);
}

.upload-example-list {
  white-space: nowrap;
}

.upload-example-item {
  width: 200rpx;
  display: inline-flex;
  flex-direction: column;
  margin-right: 12rpx;
}

.upload-example-image {
  width: 200rpx;
  height: 200rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(15, 23, 42, 0.1);
}

.upload-example-label {
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15, 23, 42, 0.74);
}
</style>
