<template>
	<view>
		<shop-goods-list></shop-goods-list>
		<view class="shop_list bgc_fff">
			<view class="color_333 font_32" style="padding: 24rpx 32rpx;">猜你喜欢</view>
			<shop-flow></shop-flow>
		</view>
	</view>
</template>

<script setup>
	import shopFlow from './components/shop-flow.vue'
	import shopGoodsList from './components/shop-goods-list.vue'
</script>

<style lang="scss">
	page {
		background-color: #F8F8F8;
	}

	.shop_list {
		border-radius: 24rpx;
	}
</style>
