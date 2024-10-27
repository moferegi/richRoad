<template>
	<view class="bgc_fff m_b_24">
		<view class="flex m_b_24">
			<image class="item_head m_r_16" src="https://inews.gtimg.com/newsapp_bt/0/14670295003/641" mode=""></image>
			<view class="flex-fitem">
				<view class="flex-aic flexr-jsb">
					<text class="color_333 font_28">匿名买家</text>
					<text class="color_999 font_24">2023-02-07 12:25:36</text>
				</view>
				<uni-rate size="16" :readonly="true" active-color="#fe5572" :value="2.5" />
			</view>
		</view>
		<view style="padding: 0 16rpx 0 80rpx;" class="boxs_bb">
			<view class="text_pre_wrap m_b_24 color_333 font_28">{{`发货快，包装完好，味道很好闻跟之前一样
发货快，包装完好，味道很好闻跟之前一样`}}</view>
			<evaluate-grid-img></evaluate-grid-img>
			<view class="flex-aic flexr-jsb color_999 font_24 bgc_f8f8f8 evaluate_num boxs_bb">
				<text>已购柔情黄/50ml</text>
				<text>x1</text>
			</view>
		</view>
	</view>
</template>

<script setup>
	import evaluateGridImg from '@/components/evaluate-grid-img/evaluate-grid-img.vue'
</script>

<style lang="scss" scoped>
.item_head {
	width: 64rpx;
	height: 64rpx;
	border-radius: 50%;
}
.evaluate_num {
	margin-top: 16rpx;
	padding: 6rpx 16rpx;
}
</style>
