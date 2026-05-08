<template>
  <view class="nf-evaluate">
    <view class="nf-evaluate-bg"></view>

    <!-- custom navbar -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ isCheck ? $t('viewReviewTitle') : $t('productReview') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body">
      <!-- product info card -->
      <view class="nf-item-card">
        <image :src="SKU.externalPicturePath ? getExternalUrl(SKU.externalPicturePath) : getUrl(SKU.picture)" class="nf-item-img" mode="aspectFill" />
        <view class="nf-item-info">
          <text class="nf-item-name">{{ $lt(SKU.name) || SKU.name }}</text>
          <view v-for="(sku, index) in SKU.attrs" :key="index" class="nf-item-spec">
            <text>{{ $lt(sku.label) || sku.label }}：{{ $lt(sku.value) || sku.value }}</text>
          </view>
        </view>
      </view>

      <!-- rating -->
      <view class="nf-card">
        <text class="nf-card-title">{{ $t('productRating') }}</text>
        <view class="nf-stars">
          <text
            v-for="star in 5"
            :key="star"
            :class="['nf-star', { active: star <= rating }]"
            @click="!isCheck && setRating(star)"
          ></text>
        </view>
        <text class="nf-rating-text">{{ getRatingText(rating) }}</text>
      </view>

      <!-- review content -->
      <view class="nf-card">
        <text class="nf-card-title">{{ $t('reviewContent') }}</text>
        <textarea
          v-model="content"
          class="nf-textarea"
          :placeholder="isCheck ? '' : $t('reviewPlaceholder')"
          maxlength="200"
          :disabled="isCheck"
        ></textarea>
        <text class="nf-char-count" v-if="!isCheck">{{ content.length }}/200</text>
      </view>

      <!-- image upload -->
      <view class="nf-card" v-if="picEnabled || (isCheck && pics.length > 0)">
        <text class="nf-card-title">{{ isCheck ? $t('reviewImages') : `${$t('uploadImages')}（${$t('uploadImagesMax').replace('{n}', '9')}）` }}</text>
        <view class="nf-pics-grid">
          <view
            v-for="(pic, index) in pics"
            :key="index"
            class="nf-pic-item"
            @tap="previewImage(index)"
          >
            <image
              :src="getUrl(pic.url || pic.tempFilePath)"
              mode="aspectFill"
              class="nf-pic-img"
            />
            <text v-if="!isCheck" class="nf-pic-del" @click.stop="deletePic(index)"></text>
          </view>
          <view
            v-if="!isCheck && pics.length < 9"
            class="nf-pic-add"
            @click="chooseImage"
          >
            <text class="nf-pic-add-icon">+</text>
            <text class="nf-pic-add-text">{{ $t('addImage') }}</text>
          </view>
        </view>
      </view>

      <!-- shop reply (view mode) -->
      <view class="nf-card" v-if="isCheck && shopReply">
        <text class="nf-card-title nf-reply-title">{{ $t('shopReplyLabel') }}</text>
        <view class="nf-reply-box">
          <text class="nf-reply-text">{{ shopReply }}</text>
          <text class="nf-reply-time" v-if="shopReplyAt">{{ formatTime(shopReplyAt) }}</text>
        </view>
      </view>
    </view>

    <!-- submit button -->
    <view class="nf-bottom-bar" v-if="!isCheck">
      <view class="nf-submit-btn" :class="{ disabled: isSubmitting }" @click="submit">
        <text>{{ isSubmitting ? $t('submitting') : $t('submitReview') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { selfOrderComment } from '@/api/order'
import { createComment } from '@/api/comment'
import { getSysConfigByKey } from '@/api/sysConfig.js'
import { onLoad } from '@dcloudio/uni-app'
import { baseUrl } from '@/utils/request.js'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const orderID = ref(0)
const goodID = ref(0)
const SKUID = ref(0)
const content = ref('')
const data = ref({})
const SKU = ref({})
const rating = ref(0)
const isCheck = ref(false)
const pics = ref([])
const isSubmitting = ref(false)
const shopReply = ref('')
const shopReplyAt = ref('')
const picEnabled = ref(true)

const goBack = () => {
  uni.navigateBack({ delta: 1, fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const setRating = (star) => {
  if (!isCheck.value) rating.value = star
}

const getRatingText = (r) => {
  const keys = ['', 'ratingVeryBad', 'ratingBad', 'ratingOk', 'ratingGood', 'ratingExcellent']
  return $t.value(keys[r] || '')
}

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

const loadPicConfig = async () => {
  try {
    const res = await getSysConfigByKey('review_pic_enabled')
    if (res.code === 0 && res.data && res.data.configValue !== undefined) {
      picEnabled.value = res.data.configValue === 'true' || res.data.configValue === '1'
    }
  } catch (e) { /* default enabled */ }
}

const init = async () => {
  const res = await selfOrderComment(orderID.value, goodID.value, SKUID.value)
  if (res.code === 0) {
    data.value = res.data
    SKU.value = res.data.detail.sku
    if (res.data.comment) {
      isCheck.value = res.data.comment.content !== ''
      rating.value = res.data.comment.rating
      content.value = res.data.comment.content
      shopReply.value = res.data.comment.shopReply || ''
      shopReplyAt.value = res.data.comment.shopReplyAt || ''
      if (res.data.comment.pics && Array.isArray(res.data.comment.pics) && res.data.comment.pics.length > 0) {
        pics.value = res.data.comment.pics.map(pic => ({
          url: typeof pic === 'string' ? pic : pic.url,
          isUploaded: true
        }))
      } else {
        pics.value = []
      }
    }
  }
}

onLoad(async (options) => {
  orderID.value = Number(options.orderID)
  goodID.value = Number(options.goodID)
  SKUID.value = Number(options.SKUID)
  loadPicConfig()
  setTimeout(() => init(), 200)
})

const chooseImage = () => {
  if (pics.value.length >= 9) {
    uni.showToast({ title: $t.value('imageMaxReached'), icon: 'none' })
    return
  }
  uni.chooseImage({
    count: 9 - pics.value.length,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const newPics = res.tempFilePaths.map(path => ({ url: path, tempFilePath: path }))
      pics.value.push(...newPics)
    }
  })
}

const deletePic = (index) => { pics.value.splice(index, 1) }

const previewImage = (index) => {
  uni.previewImage({
    current: index,
    urls: pics.value.map(pic => getUrl(pic.url || pic.tempFilePath))
  })
}

const uploadAllImages = async () => {
  if (pics.value.length === 0) return []
  const results = await Promise.allSettled(
    pics.value.map((pic, i) => pic.isUploaded ? Promise.resolve(pic.url) : uploadSingleImage(pic.tempFilePath, i))
  )
  const urls = []
  results.forEach((r) => { if (r.status === 'fulfilled') urls.push(r.value) })
  if (urls.length === 0 && pics.value.length > 0) throw new Error($t.value('uploadFail'))
  return urls
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
          const d = JSON.parse(res.data)
          if (d.code !== 0) return reject(new Error(d.msg))
          resolve(d.data.file.url)
        } catch { reject(new Error($t.value('uploadFail'))) }
      },
      fail: (err) => reject(new Error(err.errMsg))
    })
  })
}

const submit = async () => {
  if (isSubmitting.value) return
  try {
    isSubmitting.value = true
    if (rating.value === 0) {
      uni.showToast({ title: $t.value('selectRating'), icon: 'none' })
      return
    }
    let uploadedPics = []
    if (pics.value.length > 0) {
      uni.showLoading({ title: $t.value('uploadingImages'), mask: true })
      try {
        uploadedPics = await uploadAllImages()
      } catch (error) {
        uni.hideLoading()
        uni.showToast({ title: resolveApiMessage(error?.message, 'uploadFail'), icon: 'none' })
        return
      }
      uni.hideLoading()
    }
    uni.showLoading({ title: $t.value('submitting'), mask: true })
    const req = {
      orderID: orderID.value,
      goodID: goodID.value,
      SKUID: SKUID.value,
      rating: rating.value,
      content: content.value,
      pics: uploadedPics
    }
    const res = await createComment(req)
    uni.hideLoading()
    if (res.code === 0) {
      uni.showToast({ title: $t.value('reviewSuccess'), icon: 'success' })
      setTimeout(() => uni.redirectTo({ url: '/pages/order/order' }), 500)
    } else {
      uni.showToast({ title: resolveApiMessage(res.msg, 'reviewFail'), icon: 'none' })
    }
  } catch (error) {
    uni.hideLoading()
    uni.showToast({ title: $t.value('submitFail'), icon: 'none' })
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style lang="scss" scoped>
page { background: #000; }
.nf-evaluate {
  min-height: 100vh;
  background: #000;
  color: #fff;
  padding-bottom: 140rpx;
}
.nf-evaluate-bg {
  position: fixed; top: 0; left: 0; right: 0;
  height: 500rpx;
  background: radial-gradient(ellipse at 50% 0%, rgba(229, 9, 20, 0.15) 0%, transparent 70%);
  pointer-events: none; z-index: 0;
}
.nf-navbar {
  position: sticky; top: 0; z-index: 100;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(20rpx);
}
.nf-navbar-status { height: var(--status-bar-height); }
.nf-navbar-content {
  display: flex; align-items: center; justify-content: space-between;
  height: 88rpx; padding: 0 24rpx;
}
.nf-navbar-back {
  width: 64rpx; height: 64rpx;
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-title { font-size: 34rpx; font-weight: 600; color: #fff; }
.nf-body { position: relative; z-index: 1; padding: 24rpx; }

.nf-item-card {
  display: flex; padding: 24rpx;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 16rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.08);
  margin-bottom: 24rpx;
}
.nf-item-img {
  width: 140rpx; height: 140rpx;
  border-radius: 12rpx; margin-right: 24rpx; flex-shrink: 0;
}
.nf-item-info { flex: 1; display: flex; flex-direction: column; justify-content: center; }
.nf-item-name { font-size: 28rpx; color: #fff; font-weight: 500; line-height: 1.4; margin-bottom: 10rpx; }
.nf-item-spec { font-size: 24rpx; color: rgba(255, 255, 255, 0.5); margin-bottom: 4rpx; }

.nf-card {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 16rpx;
  border: 1rpx solid rgba(255, 255, 255, 0.08);
  padding: 28rpx; margin-bottom: 24rpx;
}
.nf-card-title {
  font-size: 28rpx; font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 20rpx; display: block;
}

.nf-stars { display: flex; align-items: center; margin-bottom: 12rpx; }
.nf-star {
  font-size: 48rpx; color: rgba(255, 255, 255, 0.15);
  margin-right: 12rpx; transition: all 0.2s;
  &.active { color: #e50914; text-shadow: 0 0 12rpx rgba(229, 9, 20, 0.5); }
}
.nf-rating-text { font-size: 24rpx; color: rgba(255, 255, 255, 0.5); }

.nf-textarea {
  width: 100%; min-height: 160rpx; padding: 20rpx;
  border-radius: 12rpx;
  background: rgba(255, 255, 255, 0.05);
  border: 1rpx solid rgba(255, 255, 255, 0.1);
  color: #fff; font-size: 28rpx; line-height: 1.6; box-sizing: border-box;
  &:focus { border-color: #e50914; }
}
.nf-char-count {
  display: block; text-align: right;
  font-size: 22rpx; color: rgba(255, 255, 255, 0.3); margin-top: 8rpx;
}

.nf-pics-grid { display: flex; flex-wrap: wrap; gap: 16rpx; }
.nf-pic-item {
  position: relative; width: 160rpx; height: 160rpx;
  border-radius: 12rpx; overflow: hidden;
}
.nf-pic-img { width: 100%; height: 100%; }
.nf-pic-del {
  position: absolute; top: 4rpx; right: 4rpx;
  width: 36rpx; height: 36rpx;
  background: rgba(229, 9, 20, 0.9); color: #fff;
  border-radius: 50%; font-size: 24rpx;
  display: flex; align-items: center; justify-content: center; line-height: 1;
}
.nf-pic-add {
  width: 160rpx; height: 160rpx;
  border: 2rpx dashed rgba(255, 255, 255, 0.2);
  border-radius: 12rpx;
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  background: rgba(255, 255, 255, 0.03);
}
.nf-pic-add-icon { font-size: 48rpx; color: rgba(255, 255, 255, 0.3); }
.nf-pic-add-text { font-size: 20rpx; color: rgba(255, 255, 255, 0.3); margin-top: 4rpx; }

.nf-reply-title { color: #e50914; }
.nf-reply-box {
  background: rgba(229, 9, 20, 0.08);
  border-radius: 12rpx; padding: 20rpx;
  border-left: 4rpx solid #e50914;
}
.nf-reply-text { font-size: 28rpx; color: rgba(255, 255, 255, 0.8); line-height: 1.6; }
.nf-reply-time { display: block; font-size: 22rpx; color: rgba(255, 255, 255, 0.35); margin-top: 12rpx; }

.nf-bottom-bar {
  position: fixed; bottom: 0; left: 0; right: 0;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: rgba(0, 0, 0, 0.9);
  backdrop-filter: blur(20rpx);
  border-top: 1rpx solid rgba(255, 255, 255, 0.06);
  z-index: 100;
}
.nf-submit-btn {
  height: 88rpx; background: #e50914;
  border-radius: 44rpx;
  display: flex; align-items: center; justify-content: center;
  font-size: 32rpx; font-weight: 600; color: #fff;
  letter-spacing: 2rpx;
  box-shadow: 0 4rpx 20rpx rgba(229, 9, 20, 0.4);
  &.disabled { opacity: 0.5; }
}
</style>
