<template>
  <view v-if="visible" class="refund-popup-wrapper">
    <view class="refund-mask" @tap="close"></view>
    <view class="refund-popup">
      <view class="refund-title">申请退款</view>
      <textarea
        class="refund-textarea"
        v-model="reason"
        placeholder="请填写退款原因"
        maxlength="200"
      />
      <view class="refund-images">
        <view class="image-item" v-for="(pic, index) in pics" :key="index">
          <image class="image-thumb" :src="pic.url" mode="aspectFill" @tap="previewImage(index)" />
          <view class="image-delete" @tap="removeImage(index)">×</view>
        </view>
        <view v-if="pics.length < maxImages" class="image-add" @tap="chooseImage">
          <text>+</text>
          <text class="image-add-text">上传</text>
        </view>
      </view>
      <view class="refund-actions">
        <view class="action-btn cancel" @tap="close">取消</view>
        <view class="action-btn submit" @tap="submitRefund">提交申请</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch } from 'vue'
import { baseUrl } from '@/utils/request.js'
import { applyRefund } from '@/api/order.js'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  orderId: {
    type: [String, Number],
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'success'])

const reason = ref('')
const pics = ref([])
const maxImages = 3
const submitting = ref(false)

const resetForm = () => {
  reason.value = ''
  pics.value = []
  submitting.value = false
}

watch(
  () => props.visible,
  (val) => {
    if (!val) resetForm()
  }
)

const close = () => {
  emit('update:visible', false)
}

const chooseImage = () => {
  if (pics.value.length >= maxImages) {
    uni.showToast({
      title: `最多上传${maxImages}张图片`,
      icon: 'none'
    })
    return
  }
  uni.chooseImage({
    count: maxImages - pics.value.length,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const newPics = res.tempFilePaths.map(path => ({
        url: path,
        tempFilePath: path
      }))
      pics.value.push(...newPics)
    }
  })
}

const removeImage = (index) => {
  pics.value.splice(index, 1)
}

const previewImage = (index) => {
  const urls = pics.value.map(pic => pic.url)
  uni.previewImage({
    current: index,
    urls
  })
}

const uploadSingleImage = (tempFilePath, index) => {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: baseUrl + '/fileUploadAndDownload/upload?noSave=1',
      header: {
        'x-token': uni.getStorageSync('x-token')
      },
      filePath: tempFilePath,
      name: 'file',
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          if (data.code !== 0) {
            reject(new Error(data.msg || `图片${index + 1}上传失败`))
            return
          }
          resolve(data.data.file.url)
        } catch (parseError) {
          reject(new Error(`图片${index + 1}响应解析失败`))
        }
      },
      fail: (error) => {
        reject(new Error(`图片${index + 1}上传失败: ${error.errMsg}`))
      }
    })
  })
}

const uploadAllImages = async () => {
  if (!pics.value.length) return []
  const uploadPromises = pics.value.map((pic, index) =>
    uploadSingleImage(pic.tempFilePath, index)
  )
  const results = await Promise.allSettled(uploadPromises)
  const successUrls = []
  const failedIndexes = []
  results.forEach((result, index) => {
    if (result.status === 'fulfilled') {
      successUrls.push(result.value)
    } else {
      failedIndexes.push(index)
    }
  })
  if (failedIndexes.length && successUrls.length === 0) {
    throw new Error('图片上传失败，请稍后重试')
  }
  if (failedIndexes.length && successUrls.length > 0) {
    uni.showToast({
      title: `${failedIndexes.length}张图片上传失败`,
      icon: 'none'
    })
  }
  return successUrls
}

const submitRefund = async () => {
  if (submitting.value) return
  if (!props.orderId) {
    uni.showToast({
      title: '订单信息异常',
      icon: 'none'
    })
    return
  }
  if (!reason.value.trim()) {
    uni.showToast({
      title: '请填写退款原因',
      icon: 'none'
    })
    return
  }
  submitting.value = true
  let uploadedPics = []
  try {
    if (pics.value.length > 0) {
      uni.showLoading({
        title: '上传图片中...',
        mask: true
      })
      uploadedPics = await uploadAllImages()
      uni.hideLoading()
    }

    uni.showLoading({
      title: '提交中...',
      mask: true
    })
    const res = await applyRefund({
      orderID: Number(props.orderId),
      reason: reason.value.trim(),
      images: uploadedPics
    })
    uni.hideLoading()
    if (res.code === 0) {
      uni.showToast({
        title: '申请成功',
        icon: 'none'
      })
      emit('success')
      close()
    }
  } catch (error) {
    uni.hideLoading()
    uni.showToast({
      title: error.message || '提交失败',
      icon: 'none'
    })
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss">
.refund-popup-wrapper {
  position: fixed;
  inset: 0;
  z-index: 999;
}

.refund-mask {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
}

.refund-popup {
  position: absolute;
  left: 50%;
  top: 50%;
  width: 640rpx;
  transform: translate(-50%, -50%);
  background: #fff;
  border-radius: 20rpx;
  padding: 28rpx;
}

.refund-title {
  font-size: 32rpx;
  font-weight: 600;
  text-align: center;
  margin-bottom: 20rpx;
}

.refund-textarea {
  width: 100%;
  min-height: 160rpx;
  padding: 16rpx;
  border: 1rpx solid #e6e6e6;
  border-radius: 12rpx;
  font-size: 26rpx;
  box-sizing: border-box;
  margin-bottom: 20rpx;
}

.refund-images {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-bottom: 20rpx;
}

.image-item {
  position: relative;
  width: 140rpx;
  height: 140rpx;
}

.image-thumb {
  width: 100%;
  height: 100%;
  border-radius: 12rpx;
  background: #f5f5f5;
}

.image-delete {
  position: absolute;
  right: -10rpx;
  top: -10rpx;
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  background: #ff4d4f;
  color: #fff;
  font-size: 24rpx;
  text-align: center;
  line-height: 36rpx;
}

.image-add {
  width: 140rpx;
  height: 140rpx;
  border: 1rpx dashed #ccc;
  border-radius: 12rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 36rpx;
}

.image-add-text {
  margin-top: 8rpx;
  font-size: 22rpx;
}

.refund-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16rpx;
}

.action-btn {
  padding: 12rpx 28rpx;
  border-radius: 28rpx;
  font-size: 26rpx;
}

.action-btn.cancel {
  border: 1rpx solid #ddd;
  color: #666;
}

.action-btn.submit {
  background: #fe5572;
  color: #fff;
}
</style>
