<template>
	<view class="order-evaluate-container">
<!--		<view class="header">
			<text class="title">{{ viewMode ? '查看评价' : '订单评价' }}</text>
		</view>-->

		<view class="items-container">
			<view class="item-card" v-for="(item, index) in orderItems" :key="index">
				<view class="item-header">
          <image
              :src="getUrl(item.good.imageUrl)"
              class="evaluate_pic_img"
              mode="aspectFill"
          />
					<view class="item-info">
						<text class="item-name">{{ item.name }}</text>
						<text class="item-spec" v-if="item.spec">{{ item.spec }}</text>
						<text class="item-price">¥{{ item.price }}</text>
					</view>
				</view>

				<view class="rating-section">
					<text class="section-title">商品评分</text>
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
					<text class="section-title">评价内容</text>
					<textarea
						v-model="item.comment"
						class="comment-input"
						:placeholder="viewMode ? '' : '请输入您的评价...'"
						maxlength="200"
						show-confirm-bar="false"
						:disabled="viewMode"
					></textarea>
					<text class="char-count" v-if="!viewMode">{{ item.comment.length }}/200</text>
				</view>

				<view class="image-section">
					<text class="section-title">{{ viewMode ? '评价图片' : '上传图片（最多9张）' }}</text>
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
								<text class="add-text">添加图片</text>
							</view>
						</view>
					</view>
				</view>
			</view>
		</view>

		<view class="submit-section" v-if="!viewMode">
			<button class="submit-btn" @click="submitEvaluations" :disabled="submitting">
				{{ submitting ? '提交中...' : '提交评价' }}
			</button>
		</view>

		<!-- 加载状态 -->
		<view class="loading" v-if="loading">
			<text>加载中...</text>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { selfOrder, selfOrderComment } from '@/api/order.js'
import { createComment } from '@/api/comment.js'
import { getUrl } from '@/utils/url.js'
import { baseUrl } from '@/utils/request.js'

// 响应式数据
const orderID = ref('')
const orderInfo = reactive({})
const orderItems = ref([])
const loading = ref(false)
const submitting = ref(false)
const viewMode = ref(false) // 是否为查看模式

// 页面加载
onLoad((options) => {
	if (options.orderID) {
		orderID.value = options.orderID
		// 检查是否为查看模式
		viewMode.value = options.mode === 'view'
		// 设置页面标题
		uni.setNavigationBarTitle({
			title: viewMode.value ? '查看评价' : '评价订单'
		})
		loadOrderData()
	} else {
		uni.showToast({
			title: '订单ID不能为空',
			icon: 'none'
		})
		setTimeout(() => {
			uni.navigateBack()
		}, 1500)
	}
})

// 加载订单数据
const loadOrderData = async () => {
	loading.value = true
	try {
		// 首先获取基本订单信息
		const orderRes = await selfOrder(orderID.value)
		if (orderRes.code !== 0 || !orderRes.data) {
			throw new Error('获取订单数据失败')
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
			title: '加载订单失败',
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
			title: '订单中没有商品',
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
							name: commentData.detail?.good?.name || item.sku?.name || item.name || '商品名称',
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
							name: item.sku?.name || item.name || '商品名称',
							image: item.sku?.picture || item.image,
							spec: formatSpec(item.sku?.attrs),
							price: (item.price || 0).toFixed(2),
							rating: 5,
							comment: '',
							images: []
						}
					}
				} catch (error) {
					console.error(`获取商品 ${item.goodID}-${item.skuID} 评价失败:`, error)
					// 出错时返回默认数据
					return {
						...item,
						name: item.sku?.name || item.name || '商品名称',
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
			name: item.sku?.name || item.name || '商品名称',
			image: item.sku?.picture || item.image,
			spec: formatSpec(item.sku?.attrs),
			price: (item.price || 0).toFixed(2),
			rating: 5, // 默认5星
			comment: '',
			images: []
		}))
	} else {
		uni.showToast({
			title: '订单中没有商品',
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
	return attrs.map(attr => `${attr.label}：${attr.value}`).join('；')
}

// 设置评分
const setRating = (itemIndex, rating) => {
	orderItems.value[itemIndex].rating = rating
}

// 获取评分文本
const getRatingText = (rating) => {
	const texts = ['', '很差', '较差', '一般', '满意', '非常满意']
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
      title: '暂无图片可预览',
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
      itemList: ['保存图片'],
      success: function (res) {
        if (res.tapIndex === 0) {
          // 保存图片到相册
          uni.saveImageToPhotosAlbum({
            filePath: urls[res.index],
            success: () => {
              uni.showToast({
                title: '保存成功',
                icon: 'success'
              })
            },
            fail: () => {
              uni.showToast({
                title: '保存失败',
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
        title: '图片预览失败',
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
				title: `请为第${i + 1}个商品评分`,
				icon: 'none'
			})
			return false
		}
		if (!item.comment.trim()) {
			uni.showToast({
				title: `请为第${i + 1}个商品填写评价内容`,
				icon: 'none'
			})
			return false
		}
	}
	return true
}

// 单个图片上传
const uploadSingleImage = (tempFilePath, index, itemIndex) => {
	return new Promise((resolve, reject) => {
		uni.uploadFile({
			url: baseUrl + '/fileUploadAndDownload/upload?noSave=1',
			header: {
				"x-token": uni.getStorageSync('x-token')
			},
			filePath: tempFilePath,
			name: 'file',
			success: (res) => {
				try {
					const data = JSON.parse(res.data)
					if (data.code !== 0) {
						reject(new Error(data.msg || `商品${itemIndex + 1}的图片${index + 1}上传失败`))
						return
					}
					resolve(data.data.file.url)
				} catch (parseError) {
					reject(new Error(`商品${itemIndex + 1}的图片${index + 1}响应数据解析失败`))
				}
			},
			fail: (error) => {
				reject(new Error(`商品${itemIndex + 1}的图片${index + 1}网络请求失败: ${error.errMsg}`))
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
				console.error(`商品${itemIndex + 1}的图片${index + 1}上传失败:`, result.reason)
			}
		})

		// 如果有失败的图片，给用户提示
		if (failedIndexes.length > 0) {
			const failedCount = failedIndexes.length
			const successCount = successUrls.length

			if (successCount === 0) {
				throw new Error(`商品${itemIndex + 1}的所有图片上传失败，请检查网络后重试`)
			} else {
				uni.showToast({
					title: `商品${itemIndex + 1}有${failedCount}张图片上传失败`,
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
					title: `上传商品${i + 1}的图片...`,
					mask: true
				})

				try {
					uploadedPics = await uploadItemImages(item, i)
					console.log(`商品${i + 1}图片上传成功:`, uploadedPics)
				} catch (error) {
					uni.hideLoading()
					uni.showToast({
						title: error.message || `商品${i + 1}图片上传失败`,
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
			title: '提交评价中...',
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
			title: '评价提交成功',
			icon: 'success'
		})

		setTimeout(() => {
			uni.navigateBack()
		}, 1500)

	} catch (error) {
		console.error('提交评价失败:', error)
		uni.hideLoading()
		uni.showToast({
			title: '提交失败，请重试',
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
	background: #f7f7f7;
	padding-bottom: 120rpx;
}

.header {
	background: linear-gradient(135deg, #667eea 0%, #ff4c7d 100%);
	padding: 40rpx 30rpx 30rpx;
	color: white;
	box-shadow: 0 4rpx 20rpx rgba(102, 126, 234, 0.3);

	.title {
		font-size: 36rpx;
		font-weight: bold;
		text-align: center;
	}
}

.order-info {
	background: white;
	margin: 20rpx;
	padding: 25rpx 30rpx;
	border-radius: 16rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.08);

	.order-number {
		font-size: 28rpx;
		color: #666;
		font-weight: 500;
	}
}

.items-container {
	padding: 20rpx;
}

.item-card {
	background: white;
	border-radius: 20rpx;
	margin-bottom: 20rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
	transition: transform 0.2s ease;

	&:hover {
		transform: translateY(-2rpx);
	}
}

.item-header {
	display: flex;
	padding: 30rpx;
	border-bottom: 1rpx solid #f0f0f0;

	.evaluate_pic_img {
		width: 120rpx;
		height: 120rpx;
		border-radius: 12rpx;
		margin-right: 20rpx;
		border: 1rpx solid #f0f0f0;
	}

	.item-info {
		flex: 1;
		display: flex;
		flex-direction: column;
		justify-content: space-between;

		.item-name {
			font-size: 30rpx;
			font-weight: 500;
			color: #333;
			line-height: 1.4;
			margin-bottom: 8rpx;
		}

		.item-spec {
			font-size: 24rpx;
			color: #999;
			margin-bottom: 8rpx;
		}

		.item-price {
			font-size: 32rpx;
			font-weight: bold;
			color: #ff6b6b;
		}
	}
}

.rating-section {
	padding: 30rpx;
	border-bottom: 1rpx solid #f0f0f0;

	.section-title {
		font-size: 28rpx;
		font-weight: 500;
		color: #333;
		margin-bottom: 15rpx;
	}

	.stars {
		display: flex;
		align-items: center;
		margin-bottom: 10rpx;

		.star {
			font-size: 40rpx;
			color: #ddd;
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
		color: #666;
		font-style: italic;
	}
}

.comment-section {
	padding: 30rpx;
	border-bottom: 1rpx solid #f0f0f0;

	.section-title {
		font-size: 28rpx;
		font-weight: 500;
		color: #333;
		margin-bottom: 15rpx;
	}

	.comment-input {
		width: 100%;
		min-height: 120rpx;
		padding: 20rpx;
		border: 1rpx solid #e0e0e0;
		border-radius: 12rpx;
		font-size: 28rpx;
		line-height: 1.5;
		background-color: #fafafa;
		box-sizing: border-box;
		transition: border-color 0.2s ease;

		&:focus {
			border-color: #667eea;
			background-color: white;
		}
	}

	.char-count {
		display: block;
		text-align: right;
		font-size: 24rpx;
		color: #999;
		margin-top: 10rpx;
	}
}

.image-section {
	padding: 30rpx;

	.section-title {
		font-size: 28rpx;
		font-weight: 500;
		color: #333;
		margin-bottom: 15rpx;
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
				border: 1rpx solid #e0e0e0;

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
					background: #ff4757;
					color: white;
					border-radius: 50%;
					font-size: 20rpx;
					display: flex;
					align-items: center;
					justify-content: center;
					line-height: 1;
					cursor: pointer;
					box-shadow: 0 2rpx 8rpx rgba(255, 71, 87, 0.3);
				}
			}

			.add-image-btn {
				width: 120rpx;
				height: 120rpx;
				border: 2rpx dashed #ccc;
				border-radius: 12rpx;
				display: flex;
				flex-direction: column;
				align-items: center;
				justify-content: center;
				background: #fafafa;
				cursor: pointer;
				transition: all 0.2s ease;

				&:hover {
					border-color: #667eea;
					background: #f0f2ff;
				}

				.add-icon {
					font-size: 40rpx;
					color: #999;
					margin-bottom: 5rpx;
				}

				.add-text {
					font-size: 20rpx;
					color: #999;
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
	background: white;
	padding: 20rpx 30rpx;
	border-top: 1rpx solid #e0e0e0;
	box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.1);

	.submit-btn {
		width: 100%;
		height: 80rpx;
		background: linear-gradient(135deg, #667eea 0%, #ff4c7d 100%);
		color: white;
		border: none;
		border-radius: 40rpx;
		font-size: 32rpx;
		font-weight: 500;
		box-shadow: 0 4rpx 15rpx rgba(102, 126, 234, 0.4);
		transition: all 0.2s ease;

		&:hover {
			transform: translateY(-2rpx);
			box-shadow: 0 6rpx 20rpx rgba(102, 126, 234, 0.5);
		}

		&:disabled {
			opacity: 0.6;
			transform: none;
			box-shadow: 0 4rpx 15rpx rgba(102, 126, 234, 0.2);
		}
	}
}

.loading {
	position: fixed;
	top: 50%;
	left: 50%;
	transform: translate(-50%, -50%);
	background: rgba(0, 0, 0, 0.7);
	color: white;
	padding: 20rpx 40rpx;
	border-radius: 10rpx;
	font-size: 28rpx;
	z-index: 9999;
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
