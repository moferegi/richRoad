<template>
  <view class="generate-page">
    <view class="nav">
      <view class="nav-btn" @tap="goBack">
        <uni-icons type="left" size="20" color="#0f172a" />
      </view>
      <text class="nav-title">{{ $t('tryonResultImage') }}</text>
      <view class="nav-btn" @tap="goHistory">
        <text class="nav-link">{{ $t('tryonMyTasks') }}</text>
      </view>
    </view>

    <view class="preview-panel">
      <view class="preview-col">
        <text class="preview-label">{{ $t('tryonSourceImage') }}</text>
        <image class="preview-image" :src="sourcePreview" mode="aspectFill" />
      </view>
      <view class="preview-col">
        <text class="preview-label">{{ $t('tryonTemplateImage') }}</text>
        <image class="preview-image" :src="templatePreview" mode="aspectFill" />
      </view>
    </view>

    <view class="status-card">
      <view class="status-top">
        <text class="status-title">{{ $t('taskStatusTitle') }}</text>
        <text class="status-pill" :class="taskStatusClass">{{ taskStatusText }}</text>
      </view>
      <view class="task-no" v-if="taskNo">{{ $t('tryonTaskNo') }}：{{ taskNo }}</view>
      <view class="progress-wrap" v-if="isGenerating">
        <view class="progress-track">
          <view class="progress-fill" :style="{ width: progress + '%' }"></view>
        </view>
        <text class="progress-text">{{ $t('tryonStatusProcessing') }} {{ progress }}%</text>
      </view>
      <text class="error-msg" v-if="errorMessage">{{ errorMessage }}</text>
    </view>

    <view class="result-card">
      <text class="result-title">{{ $t('tryonResultImage') }}</text>
      <view class="result-box" v-if="resultPreview">
        <image class="result-image" :src="resultPreview" mode="widthFix" @tap="previewResult" />
      </view>
      <view class="result-empty" v-else>
        <text>{{ $t('resultReadyHint') }}</text>
      </view>
    </view>

    <view class="action-row">
      <view class="action-btn secondary" @tap="downloadResult">{{ $t('downloadAction') }}</view>
      <view class="action-btn" @tap="goContinue">{{ $t('continueTryonAction') }}</view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { createTryonTask, findTryonTask } from '@/api/tryonTask.js'
import { getUrl } from '@/utils/url.js'
import {
  appendTryonLocalHistory,
  clearTryonDraft,
  getTryonDraft,
  uploadTryonImage,
} from '@/utils/tryon.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const draft = ref({})
const task = ref(null)
const sourcePreview = ref('')
const templatePreview = ref('')
const resultPreview = ref('')
const progress = ref(5)
const isGenerating = ref(false)
const errorMessage = ref('')

let pollTimer = null
let pollCount = 0

const isTempLocalPath = (value) => {
  const text = String(value || '').trim().toLowerCase()
  if (!text) return false
  return text.startsWith('blob:') || text.startsWith('file:') || text.startsWith('wxfile:') || text.startsWith('content:')
}
const POLL_INTERVAL = 2000
// Plus model tasks can take over 1 minute on provider side.
const POLL_MAX_COUNT = 90

const taskNo = computed(() => task.value?.taskNo || '')
const taskStatus = computed(() => task.value?.status || '')
const taskStatusText = computed(() => {
  if (taskStatus.value === 'success') return $t.value('tryonStatusSuccess')
  if (taskStatus.value === 'failed') return $t.value('tryonStatusFailed')
  if (isGenerating.value) return $t.value('tryonStatusProcessing')
  return $t.value('loading')
})
const taskStatusClass = computed(() => {
  if (taskStatus.value === 'success') return 'success'
  if (taskStatus.value === 'failed') return 'failed'
  if (isGenerating.value) return 'processing'
  return ''
})

const normalizeTask = (res) => {
  if (!res || !res.data) return null
  return res.data.task || res.data.reTryonTask || res.data.tryonTask || null
}

const stopPolling = () => {
  if (pollTimer) {
    clearTimeout(pollTimer)
    pollTimer = null
  }
  isGenerating.value = false
}

const finishTask = (latestTask) => {
  task.value = latestTask
  isGenerating.value = false
  progress.value = 100
  if (latestTask?.resultImage) {
    resultPreview.value = getUrl(latestTask.resultImage)
  }
  if (latestTask?.status === 'failed') {
    errorMessage.value = latestTask.errorMessage || $t.value('operationFailed')
  }
  appendTryonLocalHistory({
    ID: latestTask?.ID,
    taskNo: latestTask?.taskNo,
    requestID: draft.value.requestID,
    roomType: draft.value.roomType,
    sceneType: latestTask?.sceneType || draft.value.sceneType,
    status: latestTask?.status,
    sourceImage: latestTask?.sourceImage || draft.value.sourceRemoteUrl,
    templateImage: latestTask?.templateImage || draft.value.templateRemoteUrl,
    resultImage: latestTask?.resultImage || '',
    costPoints: latestTask?.costPoints || draft.value.modelCost || 0,
    errorMessage: latestTask?.errorMessage || '',
    CreatedAt: latestTask?.CreatedAt,
    CompletedAt: latestTask?.completedAt,
  })
  clearTryonDraft()
}

const startPolling = (taskID) => {
  if (!taskID) {
    isGenerating.value = false
    return
  }
  pollCount = 0
  isGenerating.value = true

  const loop = async () => {
    pollCount += 1
    progress.value = Math.min(95, 10 + pollCount * 3)
    try {
      const res = await findTryonTask({ ID: taskID })
      if (res.code === 0) {
        const latestTask = normalizeTask(res)
        if (latestTask) {
          task.value = latestTask
          if (latestTask.resultImage) {
            resultPreview.value = getUrl(latestTask.resultImage)
          }
          if (latestTask.status && latestTask.status !== 'processing') {
            stopPolling()
            finishTask(latestTask)
            return
          }
        }
      }
    } catch (e) {
      // 忽略单次轮询错误，继续重试
    }

    if (pollCount >= POLL_MAX_COUNT) {
      stopPolling()
      errorMessage.value = $t.value('tryonPollingTimeout')
      return
    }

    pollTimer = setTimeout(loop, POLL_INTERVAL)
  }

  pollTimer = setTimeout(loop, POLL_INTERVAL)
}

const ensureRemoteImage = async (remoteUrl, localPath) => {
  const remote = String(remoteUrl || '').trim()
  if (remote && !isTempLocalPath(remote)) return remote
  if (!localPath) return ''
  return uploadTryonImage(localPath)
}

const createTask = async () => {
  const token = uni.getStorageSync('x-token') || ''
  if (!token) {
    uni.showToast({ title: $t.value('tryonNeedLogin'), icon: 'none' })
    uni.navigateTo({ url: '/pages/user/login' })
    return
  }

  if (!draft.value.sourceLocalPath && !draft.value.sourceRemoteUrl) {
    uni.showToast({ title: $t.value('tryonPickSourceFirst'), icon: 'none' })
    return
  }
  if (!draft.value.templateLocalPath && !draft.value.templateRemoteUrl) {
    uni.showToast({ title: $t.value('tryonPickTemplateFirst'), icon: 'none' })
    return
  }

  uni.showLoading({ title: $t.value('submitting'), mask: true })
  try {
    const sourceImage = await ensureRemoteImage(draft.value.sourceRemoteUrl, draft.value.sourceLocalPath)
    const templateImage = await ensureRemoteImage(draft.value.templateRemoteUrl, draft.value.templateLocalPath)

    if (!sourceImage) {
      throw new Error($t.value('tryonPickSourceFirst'))
    }
    if (!templateImage) {
      throw new Error($t.value('tryonPickTemplateFirst'))
    }

    const res = await createTryonTask({
      requestID: draft.value.requestID,
      sceneType: draft.value.sceneType || 'clothes',
      sourceImage,
      templateImage,
      modelKey: draft.value.modelKey || '',
    })

    const taskData = normalizeTask(res)
    if (!taskData) {
      errorMessage.value = res.msg || $t.value('operationFailed')
      return
    }

    task.value = taskData
    sourcePreview.value = getUrl(taskData.sourceImage || sourceImage)
    templatePreview.value = getUrl(taskData.templateImage || templateImage)

    if (taskData.status === 'processing') {
      progress.value = 12
      startPolling(taskData.ID)
      return
    }

    finishTask(taskData)
  } catch (e) {
    errorMessage.value = e.message || $t.value('operationFailed')
  } finally {
    uni.hideLoading()
  }
}

const goBack = () => {
  uni.navigateBack({ delta: 1 })
}

const goHistory = () => {
  uni.navigateTo({ url: '/pages/tryon/history' })
}

const previewResult = () => {
  if (!resultPreview.value) return
  uni.previewImage({ urls: [resultPreview.value] })
}

const downloadResult = async () => {
  if (!resultPreview.value) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  try {
    uni.showLoading({ title: $t.value('loading'), mask: true })
    const downloadRes = await uni.downloadFile({ url: resultPreview.value })
    const filePath = downloadRes?.tempFilePath
    if (!filePath) throw new Error($t.value('downloadFailed'))
    await uni.saveImageToPhotosAlbum({ filePath })
    uni.showToast({ title: $t.value('savedToAlbum'), icon: 'none' })
  } catch (e) {
    uni.showToast({ title: $t.value('saveToAlbumFailed'), icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

const goContinue = () => {
  if (draft.value.roomType === 'shoe') {
    uni.switchTab({ url: '/pages/tabBar/shop/shop' })
    return
  }
  uni.switchTab({ url: '/pages/tabBar/index' })
}

onLoad(() => {
  draft.value = getTryonDraft()
  sourcePreview.value = draft.value.sourceRemoteUrl ? getUrl(draft.value.sourceRemoteUrl) : draft.value.sourceLocalPath
  templatePreview.value = draft.value.templateRemoteUrl ? getUrl(draft.value.templateRemoteUrl) : draft.value.templateLocalPath
  createTask()
})

onUnload(() => {
  stopPolling()
})
</script>

<style lang="scss">
page {
  background: #f4f7fb;
}

.generate-page {
  min-height: 100vh;
  padding: calc(var(--status-bar-height, 0px) + 14rpx) 20rpx 24rpx;
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
  min-width: 90rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-title {
  font-size: 30rpx;
  font-weight: 700;
}

.nav-link {
  font-size: 24rpx;
  color: rgba(15,23,42,0.62);
}

.preview-panel {
  margin-top: 16rpx;
  display: flex;
  gap: 12rpx;
}

.preview-col {
  flex: 1;
  background: rgba(255,255,255,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  border-radius: 14rpx;
  padding: 12rpx;
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.preview-label {
  font-size: 22rpx;
  color: rgba(15,23,42,0.58);
}

.preview-image {
  margin-top: 8rpx;
  width: 100%;
  height: 280rpx;
  border-radius: 10rpx;
}

.status-card,
.result-card {
  margin-top: 16rpx;
  background: rgba(255,255,255,0.95);
  border: 1rpx solid rgba(15,23,42,0.08);
  border-radius: 14rpx;
  padding: 16rpx;
  box-shadow: 0 12rpx 26rpx rgba(15, 23, 42, 0.06);
}

.status-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-title,
.result-title {
  font-size: 26rpx;
  font-weight: 600;
}

.status-pill {
  padding: 6rpx 14rpx;
  border-radius: 999rpx;
  font-size: 20rpx;
  background: rgba(15,23,42,0.08);
}

.status-pill.processing {
  background: rgba(255,205,105,0.24);
}

.status-pill.success {
  background: rgba(86,199,102,0.26);
}

.status-pill.failed {
  background: rgba(255,107,107,0.24);
}

.task-no {
  margin-top: 8rpx;
  font-size: 22rpx;
  color: rgba(15,23,42,0.62);
}

.progress-wrap {
  margin-top: 14rpx;
}

.progress-track {
  width: 100%;
  height: 12rpx;
  border-radius: 999rpx;
  background: rgba(15,23,42,0.08);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
}

.progress-text {
  margin-top: 8rpx;
  color: rgba(15,23,42,0.62);
  font-size: 22rpx;
}

.error-msg {
  margin-top: 10rpx;
  color: #ff9c9c;
  font-size: 22rpx;
}

.result-box {
  margin-top: 12rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.result-image {
  width: 100%;
}

.result-empty {
  margin-top: 12rpx;
  height: 220rpx;
  border-radius: 12rpx;
  border: 1rpx dashed rgba(15,23,42,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(15,23,42,0.55);
  font-size: 22rpx;
}

.action-row {
  margin-top: 18rpx;
  display: flex;
  gap: 12rpx;
}

.action-btn {
  flex: 1;
  height: 84rpx;
  border-radius: 999rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26rpx;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(90deg, #2563eb, #0ea5e9);
  box-shadow: 0 12rpx 26rpx rgba(37, 99, 235, 0.28);
}

.action-btn.secondary {
  background: rgba(255,255,255,0.92);
  color: #0f172a;
  border: 1rpx solid rgba(15,23,42,0.1);
  box-shadow: none;
}
</style>
