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
							>
								<image :src="img" mode="aspectFill" class="uploaded-image"></image>
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
import { selfOrder } from '@/api/order.js'
import { createComment } from '@/api/comment.js'
import { getUrl } from '@/utils/url.js'

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
		const res = await selfOrder(orderID.value)
		if (res.code === 0 && res.data) {
			Object.assign(orderInfo, res.data)
			initOrderItems()
		} else {
			throw new Error('获取订单数据失败')
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

// 初始化订单商品
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
		count: 5 - orderItems.value[itemIndex].images.length,
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

// 提交评价
const submitEvaluations = async () => {
	if (!validateForm()) {
		return
	}

	submitting.value = true

	try {
		// 批量提交评价
		for (const item of orderItems.value) {
			const commentData = {
				orderID: parseInt(orderID.value),
				goodID: item.goodID,
				SKUID: item.skuID,
				rating: item.rating,
				content: item.comment.trim(),
				pics: item.images
			}

			await createComment(commentData)
		}

		uni.showToast({
			title: '评价提交成功',
			icon: 'success'
		})

		setTimeout(() => {
			uni.navigateBack()
		}, 1500)

	} catch (error) {
		console.error('提交评价失败:', error)
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
