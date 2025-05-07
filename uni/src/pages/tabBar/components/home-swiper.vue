<template>
	<view class="swiper-container">
		<uni-swiper-dot :info="props.list" :current="current" field="content" :mode="mode">
			<swiper class="swiper-box" @change="change">
				<swiper-item v-for="(item ,index) in props.list" :key="index">
					<view class="swiper-item">
						<image class="image" :src="getUrl(item.src)" mode="aspectFill"></image>
					</view>
				</swiper-item>
			</swiper>
		</uni-swiper-dot>
	</view>
</template>

<script setup>
	import {getUrl} from "@/utils/url.js"
	import {
		ref
	} from 'vue'
	const props = defineProps({
		list: {
			type: Array,
			default: []
		}
	})

	const current = ref(0)
	const mode = ref('round')

	const change = (e) => {
		current.value = e.detail.current;
	}
</script>

<style scoped lang="scss">
	.swiper-container {
		height: 400rpx;
		width: 100%;
		padding: 0 20rpx;

		.swiper-box {
			width: 100%;
			height: 400rpx;
			border-radius: 16rpx;
			overflow: hidden;
			box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);

			.image {
				width: 100%;
				height: 400rpx;
				transition: transform 0.3s ease;
				
				&:hover {
					transform: scale(1.02);
				}
			}
		}
	}
</style>
