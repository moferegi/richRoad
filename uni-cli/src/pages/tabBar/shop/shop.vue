<template>
	<view class="shop-container">
		<shop-goods-list></shop-goods-list>
		<view class="recommend-section">
			<view class="section-title">— 猜你喜欢 —</view>
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

	.recommend-section {
		margin-top: 24rpx;
		background-color: #fff;
		border-radius: 24rpx 24rpx 0 0;
		overflow: hidden;
		padding-top: 16rpx;
	}

	.section-title {
		font-size: 32rpx;
		color: #333;
		font-weight: 500;
		text-align: center;
		padding: 24rpx 0;
		position: relative;
	}
</style>
