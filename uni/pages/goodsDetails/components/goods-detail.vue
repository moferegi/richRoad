<template>
		<view class="intro_box">
			<rich-text :nodes="detail"></rich-text>
		</view>
</template>

<script setup>
	import {ref} from 'vue'
	import { onLoad, onShow } from '@dcloudio/uni-app'
	
	const baseUrl = ref('http://localhost:8080/api')
	const props = defineProps({
		detail: {
			type: String,
			default: ''
		}
	})
	
	const data = ref({})
	
</script>

<style lang="scss" scoped>
	.intro_box{
		width: 100%;
		height: 100%;
		padding: 20rpx 32rpx;
		background-color: #fff;
		box-sizing: border-box;
		.intro_pic{
			width: 100%;
		}
		.desc{
			padding: 40rpx 16rpx;
			text-align: center;
			font-size: 36rpx;
		}
		
	}
</style>