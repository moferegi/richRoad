<template>
  <view class="nf-tryon">
    <view class="nf-tryon-bg"></view>

    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('tryonTitle') }}</text>
        <view class="nf-navbar-back"></view>
      </view>
    </view>

    <view class="nf-body">
      <view class="nf-card nf-config-card">
        <view class="nf-config-line">
          <text class="nf-config-label">{{ $t('tryonCostEach') }}</text>
          <text class="nf-config-value">{{ tryonConfig.tryon_cost_points }} {{ $t('pointsUnit') }}</text>
        </view>
        <view class="nf-config-line">
          <text class="nf-config-label">{{ $t('tryonRefundRule') }}</text>
          <text class="nf-config-value">{{ tryonConfig.tryon_fail_refund_percent }}%</text>
        </view>
        <view class="nf-config-line">
          <text class="nf-config-label">{{ $t('availablePoints') }}</text>
          <text class="nf-config-value">{{ userPoints }}</text>
        </view>
      </view>

      <view class="nf-card">
        <text class="nf-card-title">{{ $t('tryonScene') }}</text>
        <view class="nf-scene-row">
          <view class="nf-scene-btn" :class="{ active: sceneType === 'clothes' }" @tap="sceneType = 'clothes'">
            <text class="nf-scene-text">{{ $t('tryonClothes') }}</text>
          </view>
          <view class="nf-scene-btn" :class="{ active: sceneType === 'shoes' }" @tap="sceneType = 'shoes'">
            <text class="nf-scene-text">{{ $t('tryonShoes') }}</text>
          </view>
        </view>
      </view>

      <view class="nf-card">
        <text class="nf-card-title">{{ $t('tryonSourceImage') }}</text>
        <view class="nf-upload-box" @tap="pickImage('source')">
          <image v-if="sourcePreview" class="nf-upload-image" :src="sourcePreview" mode="aspectFill"></image>
          <view v-else class="nf-upload-empty">
            <uni-icons type="camera" size="24" color="rgba(255,255,255,0.65)"></uni-icons>
            <text class="nf-upload-text">{{ $t('tryonChooseSource') }}</text>
          </view>
        </view>
      </view>

      <view class="nf-card">
        <text class="nf-card-title">{{ $t('tryonTemplateImage') }}</text>
        <view class="nf-upload-box" @tap="pickImage('template')">
          <image v-if="templatePreview" class="nf-upload-image" :src="templatePreview" mode="aspectFill"></image>
          <view v-else class="nf-upload-empty">
            <uni-icons type="camera" size="24" color="rgba(255,255,255,0.65)"></uni-icons>
            <text class="nf-upload-text">{{ $t('tryonChooseTemplate') }}</text>
          </view>
        </view>
      </view>

      <view class="nf-submit" :class="{ disabled: isSubmitting }" @tap="submitTryon">
        <text class="nf-submit-text">{{ isSubmitting ? $t('submitting') : $t('tryonStart') }}</text>
      </view>

      <view class="nf-card" v-if="currentTask">
        <view class="nf-result-header">
          <text class="nf-card-title">{{ $t('tryonResultImage') }}</text>
          <view class="nf-status" :class="statusClass(currentTask.status)">
            <text class="nf-status-text">{{ statusLabel(currentTask.status) }}</text>
          </view>
        </view>
        <view class="nf-result-meta">
          <text class="nf-meta-text">{{ $t('tryonTaskNo') }}: {{ currentTask.taskNo || '-' }}</text>
        </view>
        <view class="nf-polling" v-if="isPolling && currentTask.status === 'processing'">
          <text class="nf-polling-text">{{ $t('tryonPolling') }}</text>
        </view>
        <view class="nf-result-images">
          <view class="nf-mini-image" v-if="currentTask.sourceImage" @tap="previewImage(currentTask.sourceImage)">
            <image class="nf-mini-image-img" :src="getUrl(currentTask.sourceImage)" mode="aspectFill"></image>
          </view>
          <view class="nf-mini-image" v-if="currentTask.templateImage" @tap="previewImage(currentTask.templateImage)">
            <image class="nf-mini-image-img" :src="getUrl(currentTask.templateImage)" mode="aspectFill"></image>
          </view>
          <view class="nf-mini-image" v-if="currentTask.resultImage" @tap="previewImage(currentTask.resultImage)">
            <image class="nf-mini-image-img" :src="getUrl(currentTask.resultImage)" mode="aspectFill"></image>
          </view>
        </view>
        <view class="nf-fail" v-if="currentTask.status === 'failed' && currentTask.errorMessage">
          <text class="nf-fail-text">{{ $t('tryonFailReason') }}: {{ currentTask.errorMessage }}</text>
        </view>
        <view class="nf-retry" v-if="currentTask.status === 'failed'" @tap="retryTryon">
          <text class="nf-retry-text">{{ $t('tryonRetry') }}</text>
        </view>
      </view>

      <view class="nf-card">
        <view class="nf-result-header">
          <text class="nf-card-title">{{ $t('tryonMyTasks') }}</text>
          <view class="nf-refresh" @tap="loadTaskList">
            <uni-icons type="reload" size="16" color="rgba(255,255,255,0.75)"></uni-icons>
          </view>
        </view>
        <view v-if="myTaskList.length === 0" class="nf-empty">
          <text class="nf-empty-text">{{ $t('tryonNoTask') }}</text>
        </view>
        <view v-else>
          <view
            class="nf-task-item"
            v-for="item in myTaskList"
            :key="item.ID"
            @tap="pickTask(item)"
          >
            <view class="nf-task-left">
              <text class="nf-task-no">{{ item.taskNo || '-' }}</text>
              <text class="nf-task-time">{{ formatTime(item.CreatedAt) }}</text>
            </view>
            <view class="nf-status" :class="statusClass(item.status)">
              <text class="nf-status-text">{{ statusLabel(item.status) }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow, onHide, onUnload } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useUserStore } from '@/pinia/modules/user.js'
import { baseUrl } from '@/utils/request.js'
import { getUrl } from '@/utils/url.js'
import { getTryonConfig } from '@/api/sysConfig.js'
import { createTryonTask, findTryonTask, getMyTryonTaskList } from '@/api/tryonTask.js'

const langStore = useLangStore()
const userStore = useUserStore()
const $t = computed(() => langStore.$t)

const isSubmitting = ref(false)
const sceneType = ref('clothes')
const sourceLocalPath = ref('')
const templateLocalPath = ref('')
const sourceRemoteUrl = ref('')
const templateRemoteUrl = ref('')
const currentTask = ref(null)
const myTaskList = ref([])
const userPoints = ref(0)
const isPolling = ref(false)
const pendingRequestID = ref('')

const TRYON_POLL_INTERVAL = 2000
// Plus model tasks can take over 1 minute on provider side.
const TRYON_POLL_MAX_TIMES = 90
let pollTimer = null
let pollTimes = 0

const isTempLocalPath = (value) => {
  const text = String(value || '').trim().toLowerCase()
  if (!text) return false
  return text.startsWith('blob:') || text.startsWith('file:') || text.startsWith('wxfile:') || text.startsWith('content:')
}

const tryonConfig = ref({
  tryon_guest_init_points: '3',
  tryon_register_reward_points: '8',
  tryon_cost_points: '1',
  tryon_fail_refund_percent: '100'
})

const sourcePreview = computed(() => {
  if (sourceRemoteUrl.value) return getUrl(sourceRemoteUrl.value)
  return sourceLocalPath.value
})

const templatePreview = computed(() => {
  if (templateRemoteUrl.value) return getUrl(templateRemoteUrl.value)
  return templateLocalPath.value
})

onShow(async () => {
  if (!ensureLogin()) return
  await Promise.all([loadTryonConfig(), loadTaskList(), refreshUserPoints()])
  if (currentTask.value && currentTask.value.status === 'processing' && currentTask.value.ID) {
    startPolling(currentTask.value.ID)
  }
})

onHide(() => {
  stopPolling()
})

onUnload(() => {
  stopPolling()
})

const ensureLogin = () => {
  const token = userStore.token || uni.getStorageSync('x-token')
  if (token) return true
  uni.showToast({ title: $t.value('tryonNeedLogin'), icon: 'none' })
  setTimeout(() => {
    uni.navigateTo({ url: '/pages/user/login' })
  }, 300)
  return false
}

const refreshUserPoints = async () => {
  try {
    await userStore.getInfo()
  } catch (e) {
    // ignore
  }
  const info = uni.getStorageSync('userInfo') || {}
  userPoints.value = Number(info.point || 0)
}

const goBack = () => {
  uni.navigateBack({
    delta: 1,
    fail: () => {
      uni.switchTab({ url: '/pages/tabBar/my/index' })
    }
  })
}

const loadTryonConfig = async () => {
  const res = await getTryonConfig()
  if (res.code === 0 && res.data) {
    tryonConfig.value = {
      ...tryonConfig.value,
      ...res.data
    }
  }
}

const loadTaskList = async () => {
  const res = await getMyTryonTaskList({
    page: 1,
    pageSize: 10
  })
  if (res.code === 0 && res.data) {
    myTaskList.value = res.data.list || []
    if (currentTask.value && currentTask.value.ID) {
      const latest = myTaskList.value.find(item => item.ID === currentTask.value.ID)
      if (latest) {
        currentTask.value = latest
      }
    }
    if (!currentTask.value && myTaskList.value.length > 0) {
      currentTask.value = myTaskList.value[0]
    }
  }
}

const pickTask = (task) => {
  currentTask.value = task
  if (task && task.sourceImage) {
    sourceRemoteUrl.value = task.sourceImage
    sourceLocalPath.value = ''
  }
  if (task && task.templateImage) {
    templateRemoteUrl.value = task.templateImage
    templateLocalPath.value = ''
  }
  sceneType.value = (task && task.sceneType) || 'clothes'
  if (task && task.status === 'processing' && task.ID) {
    startPolling(task.ID)
    return
  }
  stopPolling()
}

const pickImage = (kind) => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const path = (res.tempFilePaths && res.tempFilePaths[0]) || ''
      if (!path) return
      if (kind === 'source') {
        sourceLocalPath.value = path
        sourceRemoteUrl.value = ''
      } else {
        templateLocalPath.value = path
        templateRemoteUrl.value = ''
      }
      pendingRequestID.value = ''
    }
  })
}

const uploadSingleImage = (tempFilePath) => {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: baseUrl + '/fileUploadAndDownload/upload',
      header: { 'x-token': uni.getStorageSync('x-token') },
      filePath: tempFilePath,
      name: 'file',
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          if (data.code !== 0 || !data.data || !data.data.file || !data.data.file.url) {
            reject(new Error(data.msg || $t.value('uploadFail')))
            return
          }
          resolve(data.data.file.url)
        } catch (e) {
          reject(new Error($t.value('uploadFail')))
        }
      },
      fail: (err) => reject(new Error(err.errMsg || $t.value('uploadFail')))
    })
  })
}

const createActionRequestID = () => {
  return `tryon_${Date.now()}_${Math.random().toString(16).slice(2, 10)}`
}

const normalizeTaskFromResponse = (res) => {
  if (!res || !res.data) return null
  return res.data.task || res.data.reTryonTask || res.data.tryonTask || null
}

const updateTaskInList = (task) => {
  if (!task || !task.ID) return
  const index = myTaskList.value.findIndex(item => item.ID === task.ID)
  if (index < 0) return
  myTaskList.value.splice(index, 1, {
    ...myTaskList.value[index],
    ...task
  })
}

const stopPolling = () => {
  if (pollTimer) {
    clearTimeout(pollTimer)
    pollTimer = null
  }
  isPolling.value = false
}

const startPolling = (taskID) => {
  if (!taskID) return
  stopPolling()
  isPolling.value = true
  pollTimes = 0

  const poll = async () => {
    pollTimes += 1
    try {
      const res = await findTryonTask({ ID: taskID })
      if (res.code === 0) {
        const task = normalizeTaskFromResponse(res)
        if (task) {
          currentTask.value = task
          updateTaskInList(task)
          if (task.status !== 'processing') {
            stopPolling()
            await Promise.all([loadTaskList(), refreshUserPoints()])
            return
          }
        }
      }
    } catch (e) {
      // ignore single poll failure and continue polling
    }

    if (pollTimes >= TRYON_POLL_MAX_TIMES) {
      stopPolling()
      uni.showToast({ title: $t.value('tryonPollingTimeout'), icon: 'none' })
      return
    }

    pollTimer = setTimeout(poll, TRYON_POLL_INTERVAL)
  }

  pollTimer = setTimeout(poll, TRYON_POLL_INTERVAL)
}

const submitTryon = async (options = {}) => {
  const forceNewRequestID = !!options.forceNewRequestID
  if (isSubmitting.value) return
  if (!ensureLogin()) return

  if (!sourceLocalPath.value && !sourceRemoteUrl.value) {
    uni.showToast({ title: $t.value('tryonPickSourceFirst'), icon: 'none' })
    return
  }
  if (!templateLocalPath.value && !templateRemoteUrl.value) {
    uni.showToast({ title: $t.value('tryonPickTemplateFirst'), icon: 'none' })
    return
  }

  try {
    isSubmitting.value = true
    stopPolling()
    uni.showLoading({ title: $t.value('loading'), mask: true })

    let sourceImage = sourceRemoteUrl.value
    if (isTempLocalPath(sourceImage)) {
      sourceImage = ''
    }
    if (!sourceImage && sourceLocalPath.value) {
      sourceImage = await uploadSingleImage(sourceLocalPath.value)
      sourceRemoteUrl.value = sourceImage
    }
    if (!sourceImage) {
      throw new Error($t.value('tryonPickSourceFirst'))
    }

    let templateImage = templateRemoteUrl.value
    if (isTempLocalPath(templateImage)) {
      templateImage = ''
    }
    if (!templateImage && templateLocalPath.value) {
      templateImage = await uploadSingleImage(templateLocalPath.value)
      templateRemoteUrl.value = templateImage
    }
    if (!templateImage) {
      throw new Error($t.value('tryonPickTemplateFirst'))
    }

    let requestID = pendingRequestID.value
    if (!requestID || forceNewRequestID) {
      requestID = createActionRequestID()
      pendingRequestID.value = requestID
    }

    const res = await createTryonTask({
      sceneType: sceneType.value,
      sourceImage,
      templateImage,
      requestID
    })

    const taskData = normalizeTaskFromResponse(res)
    if (!taskData) {
      pendingRequestID.value = ''
    }

    if (taskData) {
      pendingRequestID.value = ''
      currentTask.value = taskData
      updateTaskInList(taskData)
      if (taskData.status === 'processing' && taskData.ID) {
        startPolling(taskData.ID)
      } else {
        stopPolling()
      }
      await Promise.all([loadTaskList(), refreshUserPoints()])
    }

    if (res.code === 0 && taskData && res.data.reused) {
      uni.showToast({ title: $t.value('tryonTaskReused'), icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: e.message || $t.value('operationFailed'), icon: 'none' })
  } finally {
    uni.hideLoading()
    isSubmitting.value = false
  }
}

const retryTryon = async () => {
  if (!currentTask.value) return
  if (currentTask.value.sourceImage) {
    sourceRemoteUrl.value = currentTask.value.sourceImage
    sourceLocalPath.value = ''
  }
  if (currentTask.value.templateImage) {
    templateRemoteUrl.value = currentTask.value.templateImage
    templateLocalPath.value = ''
  }
  sceneType.value = currentTask.value.sceneType || sceneType.value
  pendingRequestID.value = ''
  await submitTryon({ forceNewRequestID: true })
}

const previewImage = (url) => {
  if (!url) {
    uni.showToast({ title: $t.value('noImagePreview'), icon: 'none' })
    return
  }
  uni.previewImage({
    urls: [getUrl(url)]
  })
}

const statusLabel = (status) => {
  if (status === 'success') return $t.value('tryonStatusSuccess')
  if (status === 'failed') return $t.value('tryonStatusFailed')
  return $t.value('tryonStatusProcessing')
}

const statusClass = (status) => {
  if (status === 'success') return 'is-success'
  if (status === 'failed') return 'is-failed'
  return 'is-processing'
}

const formatTime = (time) => {
  if (!time) return '-'
  const d = new Date(time)
  if (Number.isNaN(d.getTime())) return '-'
  const Y = d.getFullYear()
  const M = String(d.getMonth() + 1).padStart(2, '0')
  const D = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const m = String(d.getMinutes()).padStart(2, '0')
  return `${Y}-${M}-${D} ${h}:${m}`
}
</script>

<style lang="scss" scoped>
page {
  background: #000;
}

.nf-tryon {
  min-height: 100vh;
  background: #000;
  color: #fff;
}

.nf-tryon-bg {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 560rpx;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 20% 0%, rgba(229, 9, 20, 0.18) 0%, transparent 60%),
    radial-gradient(ellipse at 80% 12%, rgba(229, 9, 20, 0.1) 0%, transparent 55%);
}

.nf-navbar {
  position: sticky;
  top: 0;
  z-index: 20;
  background: rgba(0, 0, 0, 0.82);
  backdrop-filter: blur(20rpx);
}

.nf-navbar-status {
  height: var(--status-bar-height);
}

.nf-navbar-content {
  height: 88rpx;
  padding: 0 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-navbar-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
}

.nf-body {
  position: relative;
  z-index: 1;
  padding: 24rpx;
  padding-bottom: 80rpx;
}

.nf-card {
  margin-bottom: 24rpx;
  padding: 24rpx;
  border-radius: 20rpx;
  background: rgba(255, 255, 255, 0.05);
  border: 1rpx solid rgba(255, 255, 255, 0.08);
}

.nf-config-card {
  padding-top: 18rpx;
  padding-bottom: 18rpx;
}

.nf-config-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10rpx 0;
}

.nf-config-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.66);
}

.nf-config-value {
  font-size: 26rpx;
  color: #fff;
  font-weight: 700;
}

.nf-card-title {
  display: block;
  margin-bottom: 16rpx;
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
}

.nf-scene-row {
  display: flex;
  gap: 16rpx;
}

.nf-scene-btn {
  flex: 1;
  height: 78rpx;
  border-radius: 14rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.03);
}

.nf-scene-btn.active {
  border-color: rgba(229, 9, 20, 0.8);
  background: rgba(229, 9, 20, 0.22);
}

.nf-scene-text {
  font-size: 26rpx;
  color: #fff;
  font-weight: 600;
}

.nf-upload-box {
  width: 100%;
  height: 360rpx;
  border-radius: 16rpx;
  overflow: hidden;
  border: 1rpx dashed rgba(255, 255, 255, 0.25);
  background: rgba(255, 255, 255, 0.02);
}

.nf-upload-image {
  width: 100%;
  height: 100%;
}

.nf-upload-empty {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
}

.nf-upload-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.75);
}

.nf-submit {
  height: 88rpx;
  border-radius: 999rpx;
  margin: 8rpx 0 28rpx;
  background: linear-gradient(90deg, #f20d17, #a6060d);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-submit.disabled {
  opacity: 0.6;
}

.nf-submit-text {
  color: #fff;
  font-size: 30rpx;
  font-weight: 700;
  letter-spacing: 2rpx;
}

.nf-result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.nf-status {
  padding: 6rpx 16rpx;
  border-radius: 999rpx;
  border: 1rpx solid transparent;
}

.nf-status.is-processing {
  border-color: rgba(255, 205, 105, 0.7);
  background: rgba(255, 205, 105, 0.12);
}

.nf-status.is-success {
  border-color: rgba(101, 229, 123, 0.7);
  background: rgba(101, 229, 123, 0.12);
}

.nf-status.is-failed {
  border-color: rgba(255, 107, 107, 0.7);
  background: rgba(255, 107, 107, 0.14);
}

.nf-status-text {
  font-size: 22rpx;
  color: #fff;
}

.nf-result-meta {
  margin: 10rpx 0 14rpx;
}

.nf-polling {
  margin: 0 0 12rpx;
  padding: 10rpx 12rpx;
  border-radius: 10rpx;
  border: 1rpx solid rgba(255, 205, 105, 0.35);
  background: rgba(255, 205, 105, 0.12);
}

.nf-polling-text {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.9);
}

.nf-meta-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.65);
}

.nf-result-images {
  display: flex;
  gap: 14rpx;
}

.nf-mini-image {
  width: 30%;
  aspect-ratio: 1/1;
  border-radius: 12rpx;
  overflow: hidden;
  border: 1rpx solid rgba(255, 255, 255, 0.1);
}

.nf-mini-image-img {
  width: 100%;
  height: 100%;
}

.nf-fail {
  margin-top: 12rpx;
  padding: 12rpx;
  border-radius: 10rpx;
  background: rgba(255, 107, 107, 0.12);
  border: 1rpx solid rgba(255, 107, 107, 0.25);
}

.nf-fail-text {
  font-size: 23rpx;
  color: rgba(255, 255, 255, 0.9);
}

.nf-retry {
  margin-top: 12rpx;
  height: 72rpx;
  border-radius: 12rpx;
  border: 1rpx solid rgba(229, 9, 20, 0.45);
  background: rgba(229, 9, 20, 0.18);
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-retry-text {
  font-size: 24rpx;
  color: #fff;
  font-weight: 700;
}

.nf-refresh {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.08);
}

.nf-empty {
  padding: 16rpx 0 8rpx;
}

.nf-empty-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.55);
}

.nf-task-item {
  padding: 14rpx 0;
  border-bottom: 1rpx solid rgba(255, 255, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}

.nf-task-item:last-child {
  border-bottom: none;
}

.nf-task-left {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.nf-task-no {
  font-size: 24rpx;
  color: #fff;
  font-weight: 600;
}

.nf-task-time {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.55);
}
</style>
