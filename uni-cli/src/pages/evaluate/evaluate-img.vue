<template>
	<view>
		<template v-for="(item,index) in imgList">
			<image @tap="previewImage(item,index)" class="evaluate_img m_b_16 m_r_16" :src="getUrl(item)" mode="aspectFill"></image>
		</template>
	</view>
</template>

<script setup>
	import {getUrl} from "@/utils/url.js"
	let props = defineProps({
		imgList: {
			type: Array,
			default: []
		}
	})
	const previewImage = (item,index) => {
		uni.previewImage({
			current:index,
			urls:props.imgList
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
