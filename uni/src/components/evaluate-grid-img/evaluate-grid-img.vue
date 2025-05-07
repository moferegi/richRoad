<template>
	<view>
		<template v-for="(item,index) in imgList">
			<image @tap="previewImage(item,index)" class="evaluate_img m_b_16 m_r_16" :src="item" mode="aspectFill"></image>
		</template>
	</view>
</template>

<script setup>
	import {ref} from 'vue'
	let props = defineProps({
		imgList: {
			type: Array,
			default: () => {
				return ['https://img0.baidu.com/it/u=1323189122,456399369&fm=253&fmt=auto&app=138&f=JPEG?w=500&h=500',
				'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fsafe-img.xhscdn.com%2Fbw1%2F22bf8610-6415-4d90-9781-613dbdfed967%3FimageView2%2F2%2Fw%2F1080%2Fformat%2Fjpg&refer=http%3A%2F%2Fsafe-img.xhscdn.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1678954474&t=0f8e9ad31db360bbd367edbe8bf83be2',
				'https://img0.baidu.com/it/u=181967424,3158939913&fm=253&fmt=auto&app=120&f=JPEG?w=800&h=801',
				'https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fsafe-img.xhscdn.com%2Fbw1%2F79ecc656-3654-4ba9-9eef-96c4fba37c7e%3FimageView2%2F2%2Fw%2F1080%2Fformat%2Fjpg&refer=http%3A%2F%2Fsafe-img.xhscdn.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1678954473&t=686ed8dc3b222fa9e52dab9c587dbfde']
			}
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
