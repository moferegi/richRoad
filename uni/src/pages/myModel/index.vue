<template>
  <view class="model-page">
    <view class="nav">
      <view class="nav-btn" @tap="goBack">
        <uni-icons type="left" size="20" color="#0f172a" />
      </view>
      <text class="nav-title">{{ $t('myModels') }}</text>
      <view class="nav-btn" @tap="addModel">
        <uni-icons type="plusempty" size="22" color="#0f172a" />
      </view>
    </view>

    <view class="tips">{{ $t('myModelTips') }}</view>

    <scroll-view class="list-wrap" scroll-y>
      <view class="grid">
        <view class="card" v-for="item in modelList" :key="item.id">
          <image class="card-image" :src="getUrl(item.url)" mode="aspectFill" @tap="preview(item)" />
          <view class="card-foot">
            <text class="card-name">{{ item.name || $t('unnamedModel') }}</text>
            <view class="card-actions">
              <view class="mini-btn use" @tap.stop="useForRoom(item, 'tryon')">{{ $t('tryonRoom') }}</view>
              <view class="mini-btn use" @tap.stop="useForRoom(item, 'shoe')">{{ $t('shoeRoom') }}</view>
            </view>
            <view class="card-actions card-actions-secondary">
              <view class="mini-btn" @tap.stop="renameModel(item)">{{ $t('rename') }}</view>
              <view class="mini-btn danger" @tap.stop="removeModel(item)">{{ $t('delete') }}</view>
            </view>
          </view>
        </view>
      </view>

      <view class="empty" v-if="modelList.length === 0">
        <text>{{ $t('myModelEmpty') }}</text>
      </view>
      <view style="height: 40rpx"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getUrl } from '@/utils/url.js'
import { setSelectedTryonModel, uploadTryonImage } from '@/utils/tryon.js'
import {
  createTryonModel,
  updateTryonModel,
  deleteTryonModel,
  getMyTryonModelList,
} from '@/api/tryonTask.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const modelList = ref([])

const isTempLocalPath = (value) => {
  const text = String(value || '').trim().toLowerCase()
  if (!text) return false
  return text.startsWith('blob:') || text.startsWith('file:') || text.startsWith('wxfile:') || text.startsWith('content:')
}

const normalizeModelItem = (item) => {
  const id = item?.ID || item?.id || ''
  return {
    id: String(id),
    name: item?.name || '',
    url: item?.image || item?.url || '',
    createdAt: item?.CreatedAt || item?.createdAt || '',
  }
}

const loadModels = async () => {
  const res = await getMyTryonModelList()
  if (res.code !== 0) {
    return
  }
  const list = Array.isArray(res?.data?.list) ? res.data.list : []
  modelList.value = list.map(normalizeModelItem).filter(item => item.id)
}

const addModel = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const path = res.tempFilePaths && res.tempFilePaths[0]
      if (!path) return

      uni.showLoading({ title: $t.value('uploading') || '上传中', mask: true })
      let uploadedUrl = ''
      try {
        uploadedUrl = await uploadTryonImage(path, 'tryon/model/source')
      } catch (e) {
        uni.showToast({ title: e.message || $t.value('uploadFail'), icon: 'none' })
        return
      } finally {
        uni.hideLoading()
      }

      const createRes = await createTryonModel({
        name: `${$t.value('modelDefaultPrefix')}${modelList.value.length + 1}`,
        image: uploadedUrl,
      })
      if (createRes.code !== 0) {
        return
      }

      await loadModels()
      uni.showToast({ title: $t.value('createSuccess') || '创建成功', icon: 'none' })
    },
  })
}

const renameModel = (item) => {
  uni.showModal({
    title: $t.value('rename'),
    editable: true,
    placeholderText: $t.value('modelNamePlaceholder'),
    content: item.name || '',
    success: async (res) => {
      if (!res.confirm) return
      const value = (res.content || '').trim()
      if (!value) return

      const renameRes = await updateTryonModel({
        ID: Number(item.id),
        name: value,
      })
      if (renameRes.code !== 0) {
        return
      }

      item.name = value
      uni.showToast({ title: $t.value('updateSuccess') || '更新成功', icon: 'none' })
    },
  })
}

const removeModel = (item) => {
  uni.showModal({
    title: $t.value('pendingOrderTitle'),
    content: $t.value('confirmDeleteModel'),
    success: async (res) => {
      if (!res.confirm) return

      const deleteRes = await deleteTryonModel({ ID: Number(item.id) })
      if (deleteRes.code !== 0) {
        return
      }

      modelList.value = modelList.value.filter(v => v.id !== item.id)
      uni.showToast({ title: $t.value('deleteSuccess') || '删除成功', icon: 'none' })
    },
  })
}

const preview = (item) => {
  if (!item?.url) return
  uni.previewImage({ urls: [getUrl(item.url)] })
}

const useForRoom = (item, roomType) => {
  if (!item?.url) {
    uni.showToast({ title: $t.value('modelImageMissing'), icon: 'none' })
    return
  }

  const value = String(item.url || '').trim()
  const isRemote = !!value && !isTempLocalPath(value)
  setSelectedTryonModel({
    roomType,
    modelID: String(item.id || ''),
    modelName: item.name,
    localPath: isRemote ? '' : value,
    remoteUrl: isRemote ? getUrl(value) : '',
  })

  uni.switchTab({
    url: roomType === 'shoe' ? '/pages/tabBar/shop/shop' : '/pages/tabBar/index',
  })
}

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

onShow(() => {
  loadModels()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.model-page {
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

.tips {
  margin-top: 10rpx;
  margin-bottom: 12rpx;
  font-size: 22rpx;
  color: rgba(15,23,42,0.62);
}

.list-wrap {
  height: calc(100vh - var(--status-bar-height, 0px) - 120rpx);
}

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12rpx;
}

.card {
  background: rgba(255,255,255,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  border-radius: 12rpx;
  overflow: hidden;
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.card-image {
  width: 100%;
  height: 260rpx;
}

.card-foot {
  padding: 10rpx;
}

.card-name {
  font-size: 22rpx;
}

.card-actions {
  margin-top: 10rpx;
  display: flex;
  gap: 8rpx;
}

.card-actions-secondary {
  margin-top: 8rpx;
}

.mini-btn {
  flex: 1;
  height: 50rpx;
  border-radius: 999rpx;
  background: rgba(248,250,252,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  color: #0f172a;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
}

.mini-btn.use {
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  border: none;
  color: #fff;
}

.mini-btn.danger {
  background: rgba(220, 38, 38, 0.12);
  border-color: rgba(220, 38, 38, 0.24);
  color: #b91c1c;
}

.empty {
  margin-top: 140rpx;
  text-align: center;
  color: rgba(15,23,42,0.55);
}
</style>
