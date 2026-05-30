<template>
	<view>
		<template v-for="item in imageViews" :key="item._imageKey">
			<image @tap="previewImage(item._rawIndex)" class="evaluate_img m_b_16 m_r_16" :src="item._imageUrl" mode="aspectFill"></image>
		</template>
	</view>
</template>

<script setup>
	import { computed } from 'vue'
	import {getUrl} from "@/utils/url.js"
	let props = defineProps({
		imgList: {
			type: Array,
			default: () => []
		}
	})
	// 评价图片只用于展示和预览；提前归一化 URL，模板和 previewImage 共用同一组地址。
	const imageViews = computed(() => (props.imgList || []).map((item, index) => {
		const imageUrl = getUrl(item)
		return {
			_rawIndex: index,
			_imageKey: `${imageUrl || item || 'evaluate-img'}-${index}`,
			_imageUrl: imageUrl,
		}
	}))

	const previewImage = (index) => {
		uni.previewImage({
			current:index,
			urls:imageViews.value.map(item => item._imageUrl)
		})
	}
</script>

<style lang="scss" scoped>
.evaluate_img {
	width: 168rpx;
	height: 168rpx;
	border-radius: 2rpx;
	overflow: hidden;
}
</style>
