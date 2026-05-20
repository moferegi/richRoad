<template>
	<view class="order-evaluate-container">

		<view class="items-container">
			<view class="item-card" v-for="(item, index) in orderItems" :key="index">
				<view class="item-header">
          <image
              :src="item.good.externalImagePath ? getExternalUrl(item.good.externalImagePath) : getUrl(item.good.imageUrl)"
              class="evaluate_pic_img"
              mode="aspectFill"
          />
					<view class="item-info">
						<text class="item-name">{{ $lt(item.name) || item.name }}</text>
						<text class="item-spec" v-if="item.spec">{{ item.spec }}</text>
						<text class="item-price">{{ cs }}{{ item.price }}</text>
					</view>
				</view>

				<view class="rating-section">
					<text class="section-title">{{ $t('productRating') }}</text>
					<view class="stars">
						<text
							v-for="star in 5"
							:key="star"
							:class="['star', { active: star <= item.rating }]"
							@click="!viewMode && setRating(index, star)"
						>★</text>
					</view>
					<text class="rating-text">{{ getRatingText(item.rating) }}</text>
				</view>

				<view class="comment-section">
					<text class="section-title">{{ $t('reviewContent') }}</text>
					<textarea
						v-model="item.comment"
						class="comment-input"
						:placeholder="viewMode ? '' : $t('reviewPlaceholder')"
						maxlength="200"
						show-confirm-bar="false"
						:disabled="viewMode"
					></textarea>
					<text class="char-count" v-if="!viewMode">{{ item.comment.length }}/200</text>
				</view>

				<view class="image-section">
					<text class="section-title">{{ viewMode ? $t('reviewImages') : $t('uploadImages') + '（' + $t('uploadImagesMax').replace('{n}', '9') + '）' }}</text>
					<view class="image-upload">
						<view class="uploaded-images">
							<view
								v-for="(img, imgIndex) in item.images"
								:key="imgIndex"
								class="image-item"
                @tap="previewImage(index, imgIndex)"
							>
								<image :src="getUrl(img)" mode="aspectFill" class="uploaded-image"></image>
								<text v-if="!viewMode" class="delete-btn" @click="deleteImage(index, imgIndex)">×</text>
							</view>
							<view
								v-if="!viewMode && item.images.length < 9"
								class="add-image-btn"
								@click="chooseImage(index)"
							>
								<text class="add-icon">+</text>
								<text class="add-text">{{ $t('addImage') }}</text>
							</view>
						</view>
					</view>
				</view>
			</view>
		</view>

		<view class="submit-section" v-if="!viewMode">
			<button class="submit-btn" @click="submitEvaluations" :disabled="submitting">
				{{ submitting ? $t('submitting') : $t('submitReview') }}
			</button>
		</view>

		<!-- 加载状态 -->
		<view class="loading" v-if="loading">
			<text>{{ $t('loading') }}</text>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { selfOrder, selfOrderComment } from '@/api/order.js'
import { createComment } from '@/api/comment.js'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl, getExternalUrl } from '@/utils/url.js'
import { baseUrl } from '@/utils/request.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)
const locale = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)

// 响应式数据
const orderID = ref('')
const orderInfo = reactive({})
const orderItems = ref([])
const loading = ref(false)
const submitting = ref(false)
const viewMode = ref(false) // 是否为查看模式
const lastLoadedLocale = ref('')

const updatePageTitle = () => {
	uni.setNavigationBarTitle({
		title: viewMode.value ? $t.value('viewReviewTitle') : $t.value('orderEvaluateTitle')
	})
}

// 页面加载
onLoad((options) => {
	if (options.orderID) {
		orderID.value = options.orderID
		// 检查是否为查看模式
		viewMode.value = options.mode === 'view'
		// 设置页面标题
		updatePageTitle()
		loadOrderData()
		lastLoadedLocale.value = locale.value
	} else {
		uni.showToast({
			title: $t.value('orderIdRequired'),
			icon: 'none'
		})
		setTimeout(() => {
			uni.navigateBack()
		}, 1500)
	}
})

onShow(() => {
	const localeChanged = !!lastLoadedLocale.value && lastLoadedLocale.value !== locale.value
	if (localeChanged) {
		updatePageTitle()
		// 查看模式可安全重拉；编辑模式保留用户草稿，避免切语言后输入丢失。
		if (viewMode.value && orderID.value) {
			loadOrderData()
		}
	}
	lastLoadedLocale.value = locale.value
})

// 加载订单数据
const loadOrderData = async () => {
	loading.value = true
	try {
		// 首先获取基本订单信息
		const orderRes = await selfOrder(orderID.value)
		if (orderRes.code !== 0 || !orderRes.data) {
			throw new Error('order data load failed')
		}

		Object.assign(orderInfo, orderRes.data)

		// 如果是查看模式，需要为每个商品获取评价信息
		if (viewMode.value && orderInfo.detail && orderInfo.detail.length > 0) {
			await loadOrderItemsWithComments()
		} else {
			// 评价模式，直接初始化商品列表
			initOrderItems()
		}
	} catch (error) {
		console.error('加载订单数据失败:', error)
		uni.showToast({
			title: $t.value('loadFail'),
			icon: 'none'
		})
		setTimeout(() => {
			uni.navigateBack()
		}, 1500)
	} finally {
		loading.value = false
	}
}

// 加载订单商品及其评价信息（查看模式）
const loadOrderItemsWithComments = async () => {
	if (!orderInfo.detail || orderInfo.detail.length === 0) {
		uni.showToast({
			title: $t.value('noItemsInOrder'),
			icon: 'none'
		})
		setTimeout(() => {
			uni.navigateBack()
		}, 1500)
		return
	}

	try {
		const itemsWithComments = await Promise.all(
			orderInfo.detail.map(async (item) => {
				try {
					// 为每个商品获取评价信息
					const commentRes = await selfOrderComment(orderID.value, item.goodID, item.skuID)

					if (commentRes.code === 0 && commentRes.data) {
						const commentData = commentRes.data
						return {
							...item,
							name: commentData.detail?.good?.name || item.sku?.name || item.name || '',
							image: commentData.detail?.good?.imageUrl || item.sku?.picture || item.image,
							spec: formatSpec(commentData.detail?.sku?.attrs || item.sku?.attrs),
							price: (commentData.detail?.price || item.price || 0).toFixed(2),
							rating: commentData.comment?.rating || 5,
							comment: commentData.comment?.content || '',
							images: commentData.comment?.pics || []
						}
					} else {
						// 如果获取评价失败，使用默认数据
						return {
							...item,
						name: item.sku?.name || item.name || '',
						image: item.sku?.picture || item.image,
						spec: formatSpec(item.sku?.attrs),
						price: (item.price || 0).toFixed(2),
						rating: 5,
						comment: '',
						images: []
					}
				}
			} catch (error) {
					console.error(`get comment error ${item.goodID}-${item.skuID}:`, error)
					// 出错时返回默认数据
					return {
						...item,
						name: item.sku?.name || item.name || '',
						image: item.sku?.picture || item.image,
						spec: formatSpec(item.sku?.attrs),
						price: (item.price || 0).toFixed(2),
						rating: 5,
						comment: '',
						images: []
					}
				}
			})
		)

		orderItems.value = itemsWithComments
	} catch (error) {
		console.error('加载商品评价信息失败:', error)
		// 如果批量加载失败，回退到基本初始化
		initOrderItems()
	}
}

// 初始化订单商品（评价模式）
const initOrderItems = () => {
	if (orderInfo.detail && orderInfo.detail.length > 0) {
		orderItems.value = orderInfo.detail.map(item => ({
			...item,
			name: item.sku?.name || item.name || '',
			image: item.sku?.picture || item.image,
			spec: formatSpec(item.sku?.attrs),
			price: (item.price || 0).toFixed(2),
			rating: 5, // 默认5星
			comment: '',
			images: []
		}))
	} else {
		uni.showToast({
			title: $t.value('noItemsInOrder'),
			icon: 'none'
		})
		setTimeout(() => {
			uni.navigateBack()
		}, 1500)
	}
}

// 格式化规格
const formatSpec = (attrs) => {
	if (!attrs || !Array.isArray(attrs)) return ''
	return attrs.map(attr => `${$lt.value(attr.label)}：${$lt.value(attr.value)}`).join('；')
}

// 设置评分
const setRating = (itemIndex, rating) => {
	orderItems.value[itemIndex].rating = rating
}

// 获取评分文本
const getRatingText = (rating) => {
	const texts = ['', $t.value('ratingVeryBad'), $t.value('ratingBad'), $t.value('ratingOk'), $t.value('ratingGood'), $t.value('ratingExcellent')]
	return texts[rating] || ''
}

// 选择图片
const chooseImage = (itemIndex) => {
	uni.chooseImage({
		count: 9 - orderItems.value[itemIndex].images.length,
		sizeType: ['compressed'],
		sourceType: ['album', 'camera'],
		success: (res) => {
			orderItems.value[itemIndex].images.push(...res.tempFilePaths)
		},
		fail: (error) => {
			console.error('选择图片失败:', error)
		}
	})
}

// 删除图片
const deleteImage = (itemIndex, imageIndex) => {
	orderItems.value[itemIndex].images.splice(imageIndex, 1)
}

// 图片预览方法
const previewImage = (itemIndex, imageIndex = 0) => {
  const currentItem = orderItems.value[itemIndex]
  if (!currentItem || !currentItem.images || currentItem.images.length === 0) {
    uni.showToast({
      title: $t.value('noImagePreview'),
      icon: 'none'
    })
    return
  }

  // 获取当前商品的所有图片URL
  const urls = currentItem.images.map(img => {
    // 如果是本地临时文件路径，直接使用
    if (img.startsWith('blob:') || img.startsWith('file://') || img.startsWith('/')) {
      return img
    }
    // 如果是网络图片，使用getUrl处理
    return getUrl(img)
  })

  uni.previewImage({
    current: imageIndex, // 当前预览的图片索引
    urls: urls, // 图片URL数组
    longPressActions: {
      itemList: [$t.value('saveImage')],
      success: function (res) {
        if (res.tapIndex === 0) {
          // 保存图片到相册
          uni.saveImageToPhotosAlbum({
            filePath: urls[res.index],
            success: () => {
              uni.showToast({
                title: $t.value('saveSuccess'),
                icon: 'success'
              })
            },
            fail: () => {
              uni.showToast({
                title: $t.value('saveFail'),
                icon: 'error'
              })
            }
          })
        }
      }
    },
    fail: (err) => {
      console.error('图片预览失败:', err)
      uni.showToast({
        title: $t.value('imagePreviewFail'),
        icon: 'error'
      })
    }
  })
}

// 表单验证
const validateForm = () => {
	for (let i = 0; i < orderItems.value.length; i++) {
		const item = orderItems.value[i]
		if (!item.rating || item.rating < 1) {
			uni.showToast({
				title: $t.value('selectRating'),
				icon: 'none'
			})
			return false
		}
		if (!item.comment.trim()) {
			uni.showToast({
				title: $t.value('reviewPlaceholder'),
				icon: 'none'
			})
			return false
		}
	}
	return true
}

// 单个图片上传
const uploadSingleImage = (tempFilePath, _index, _itemIndex) => {
	return new Promise((resolve, reject) => {
		uni.uploadFile({
			url: baseUrl + '/fileUploadAndDownload/upload',
			header: {
				"x-token": uni.getStorageSync('x-token')
			},
			filePath: tempFilePath,
			name: 'file',
			success: (res) => {
				try {
					const data = JSON.parse(res.data)
					if (data.code !== 0) {
						reject(new Error(resolveApiMessage(data.msg, 'uploadFail')))
						return
					}
					resolve(data.data.file.url)
				} catch (parseError) {
					reject(new Error($t.value('uploadFail')))
				}
			},
			fail: () => {
				reject(new Error($t.value('uploadFail')))
			}
		})
	})
}

// 批量上传单个商品的图片
const uploadItemImages = async (item, itemIndex) => {
	if (!item.images || item.images.length === 0) return []

	const uploadPromises = item.images.map((imagePath, index) =>
		uploadSingleImage(imagePath, index, itemIndex)
	)

	try {
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

		// 如果有失败的图片，给用户提示
		if (failedIndexes.length > 0) {
			const failedCount = failedIndexes.length
			const successCount = successUrls.length

			if (successCount === 0) {
				throw new Error($t.value('imageUploadAllFail'))
			} else {
				uni.showToast({
					title: $t.value('nImagesUploadFail').replace('{n}', failedCount),
					icon: 'none',
					duration: 2000
				})
			}
		}

		return successUrls

	} catch (error) {
		console.error(`商品${itemIndex + 1}图片批量上传失败:`, error)
		throw error
	}
}

// 提交评价
const submitEvaluations = async () => {
	if (!validateForm()) {
		return
	}

	submitting.value = true

	try {
		// 1. 先批量上传所有商品的图片
		const itemsWithUploadedImages = []

		for (let i = 0; i < orderItems.value.length; i++) {
			const item = orderItems.value[i]

			// 如果有图片需要上传
			let uploadedPics = []
			if (item.images && item.images.length > 0) {
				uni.showLoading({
					title: $t.value('uploadingImages'),
					mask: true
				})

				try {
					uploadedPics = await uploadItemImages(item, i)
				} catch (error) {
					uni.hideLoading()
					uni.showToast({
						title: resolveApiMessage(error?.message, 'uploadFail'),
						icon: 'none',
						duration: 2000
					})
					return
				}

				uni.hideLoading()
			}

			itemsWithUploadedImages.push({
				...item,
				uploadedPics
			})
		}

		// 2. 提交所有评价数据
		uni.showLoading({
			title: $t.value('submitting'),
			mask: true
		})

		for (const item of itemsWithUploadedImages) {
			const commentData = {
				orderID: parseInt(orderID.value),
				goodID: item.goodID,
				SKUID: item.skuID,
				rating: item.rating,
				content: item.comment.trim(),
				pics: item.uploadedPics
			}

			await createComment(commentData)
		}

		uni.hideLoading()
		uni.showToast({
			title: $t.value('reviewSuccess'),
			icon: 'success'
		})

		setTimeout(() => {
			uni.navigateBack()
		}, 1500)

	} catch (error) {
		console.error('提交评价失败:', error)
		uni.hideLoading()
		uni.showToast({
			title: $t.value('submitFail'),
			icon: 'none'
		})
	} finally {
		submitting.value = false
	}
}
</script>

<style lang="scss" scoped>
.order-evaluate-container {
	min-height: 100vh;
	background: #000;
	padding-bottom: 120rpx;
}

.header {
	background: #1a1a1a;
	padding: 40rpx 30rpx 30rpx;
	color: #fff;

	.title {
		font-size: 36rpx;
		font-weight: bold;
		text-align: center;
	}
}

.order-info {
	background: rgba(255, 255, 255, 0.04);
	border: 1rpx solid rgba(255, 255, 255, 0.06);
	margin: 20rpx;
	padding: 25rpx 30rpx;
	border-radius: 16rpx;

	.order-number {
		font-size: 28rpx;
		color: rgba(255, 255, 255, 0.6);
		font-weight: 500;
	}
}

.items-container {
	padding: 20rpx;
}

.item-card {
	background: rgba(255, 255, 255, 0.04);
	border: 1rpx solid rgba(255, 255, 255, 0.06);
	border-radius: 20rpx;
	margin-bottom: 20rpx;
	overflow: hidden;
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	transition: transform 0.2s ease;

	&:hover {
		transform: translateY(-2rpx);
	}
}

.item-header {
	display: flex;
	padding: 30rpx;
	border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);

	.evaluate_pic_img {
		width: 120rpx;
		height: 120rpx;
		border-radius: 12rpx;
		margin-right: 20rpx;
		border: 1rpx solid rgba(255, 255, 255, 0.08);
	}

	.item-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: space-between;

		.item-name {
			font-size: 30rpx;
			font-weight: 500;
			color: rgba(255, 255, 255, 0.9);
			line-height: 1.4;
			margin-bottom: 8rpx;
		}

		.item-spec {
			font-size: 24rpx;
			color: rgba(255, 255, 255, 0.4);
			margin-bottom: 8rpx;
		}

		.item-price {
			font-size: 32rpx;
			font-weight: bold;
			color: #e50914;
		}
	}
}

.rating-section {
	padding: 30rpx;
	border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);

	.section-title {
		font-size: 28rpx;
		font-weight: 500;
		color: rgba(255, 255, 255, 0.9);
		margin-bottom: 15rpx;
	}

	.stars {
		display: flex;
		align-items: center;
		margin-bottom: 10rpx;

		.star {
			font-size: 40rpx;
			color: rgba(255, 255, 255, 0.15);
			margin-right: 8rpx;
			transition: color 0.2s ease;
			cursor: pointer;

			&.active {
				color: #ffd700;
				text-shadow: 0 0 10rpx rgba(255, 215, 0, 0.5);
			}

			&:hover {
				transform: scale(1.1);
			}
		}
	}

	.rating-text {
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.5);
		font-style: italic;
	}
}

.comment-section {
	padding: 30rpx;
	border-bottom: 1rpx solid rgba(255, 255, 255, 0.06);

	.section-title {
		font-size: 28rpx;
		font-weight: 500;
		color: rgba(255, 255, 255, 0.9);
		margin-bottom: 15rpx;
	}

	.comment-input {
		width: 100%;
		min-height: 120rpx;
		padding: 20rpx;
		border: 1rpx solid rgba(255, 255, 255, 0.1);
		border-radius: 12rpx;
		font-size: 28rpx;
		line-height: 1.5;
		background-color: rgba(255, 255, 255, 0.04);
		color: #fff;
		box-sizing: border-box;
		transition: border-color 0.2s ease;

		&:focus {
			border-color: #e50914;
			background-color: rgba(255, 255, 255, 0.06);
		}
	}

	.char-count {
		display: block;
		text-align: right;
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.35);
		margin-top: 10rpx;
	}
}

.section-title {
	font-size: 28rpx;
	font-weight: 500;
	color: rgba(255, 255, 255, 0.9);
	margin-bottom: 15rpx;
	padding: 0 30rpx;
}

.comment-input {
	width: calc(100% - 60rpx);
	min-height: 120rpx;
	padding: 20rpx;
	margin: 0 30rpx;
	border: 1rpx solid rgba(255, 255, 255, 0.1);
	border-radius: 12rpx;
	font-size: 28rpx;
	line-height: 1.5;
	background-color: rgba(255, 255, 255, 0.04);
	color: #fff;
	box-sizing: border-box;
	transition: border-color 0.2s ease;

	&:focus {
		border-color: #e50914;
		background-color: rgba(255, 255, 255, 0.06);
	}
}

.char-count {
	display: block;
	text-align: right;
	font-size: 24rpx;
	color: rgba(255, 255, 255, 0.35);
	margin-top: 10rpx;
	padding: 0 30rpx;
}

.image-section {
	padding: 30rpx;

	.section-title {
		font-size: 28rpx;
		font-weight: 500;
		color: rgba(255, 255, 255, 0.9);
		margin-bottom: 15rpx;
		padding: 0;
	}

	.image-upload {
		.uploaded-images {
			display: flex;
			flex-wrap: wrap;
			gap: 15rpx;
			margin-top: 30rpx;

			.image-item {
				position: relative;
				width: 120rpx;
				height: 120rpx;
				border-radius: 12rpx;
				overflow: hidden;
				border: 1rpx solid rgba(255, 255, 255, 0.1);

				.uploaded-image {
					width: 100%;
					height: 100%;
				}

				.delete-btn {
					position: absolute;
					top: -5rpx;
					right: -5rpx;
					width: 30rpx;
					height: 30rpx;
					background: #e50914;
					color: white;
					border-radius: 50%;
					font-size: 20rpx;
					display: flex;
					align-items: center;
					justify-content: center;
					line-height: 1;
					cursor: pointer;
				}
			}

			.add-image-btn {
				width: 120rpx;
				height: 120rpx;
				border: 2rpx dashed rgba(255, 255, 255, 0.2);
				border-radius: 12rpx;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				background: rgba(255, 255, 255, 0.04);
				cursor: pointer;
				transition: all 0.2s ease;

				&:hover {
					border-color: #e50914;
					background: rgba(229, 9, 20, 0.06);
				}

				.add-icon {
					font-size: 40rpx;
					color: rgba(255, 255, 255, 0.4);
					margin-bottom: 5rpx;
				}

				.add-text {
					font-size: 20rpx;
					color: rgba(255, 255, 255, 0.4);
				}
			}
		}
	}
}

.submit-section {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #1a1a1a;
	padding: 20rpx 30rpx;
	border-top: 1rpx solid rgba(255, 255, 255, 0.06);

	.submit-btn {
		width: 100%;
		height: 80rpx;
		background: #e50914;
		color: white;
		border: none;
		border-radius: 40rpx;
		font-size: 32rpx;
		font-weight: 500;
		transition: all 0.2s ease;

		&:active {
			opacity: 0.85;
		}

		&:disabled {
			opacity: 0.5;
		}
	}
}

.loading {
	position: fixed;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	background: rgba(26, 26, 26, 0.9);
	color: white;
	padding: 20rpx 40rpx;
	border-radius: 10rpx;
	font-size: 28rpx;
	z-index: 9999;
	backdrop-filter: blur(10px);
}

/* 动画效果 */
@keyframes fadeIn {
	from {
		opacity: 0;
		transform: translateY(20rpx);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.item-card {
	animation: fadeIn 0.3s ease-out;
}
</style>
