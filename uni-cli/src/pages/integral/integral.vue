<template>
	<view>
		<view class="integral_title_box color_333 font_36">
			可用积分：100
		</view>
		<view class="bgc_fff integral_ul">
			<template v-for="(item,index) in 3">
				<view class="flex-aic flexr-jsb integral_item" :class="index == 0 ? '' : 'b_t_2'">
					<image class="integral_item_img m_r_24" src="http://www.liwanying.top/applate-icon/integral-icon.png" mode=""></image>
					<view class="flex-fitem">
						<view class="color_333 font_28 font_500 m_b_8">积分商城兑换</view>
						<view class="color_666 font_24">2020.10.11 13:40:12</view>
					</view>
					<view class="font_32 color_de1f26">-20</view>
				</view>
			</template>
		</view>
	</view>
</template>

<script setup>
	
</script>

<style lang="scss">
page {
	background-color: #f8f8f8;
}
.integral_ul {
	width: 686rpx;
	margin: 0 auto;
	border-radius: 20rpx;
}
.integral_title_box {
	padding: 24rpx 32rpx;
}
.integral_item {
	padding: 32rpx;
}
.integral_item_img {
	width: 80rpx;
	height: 80rpx;
}
</style>
