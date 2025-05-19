<template>
	<view>
	<view v-for="(item, index) in couponList" :key="index" class="coupon_box">
		<view class="left">
			<view class="left_top">
				<text class="hui" :style="'background:' + colors">惠</text>
				<text class="hui_name">{{item.name}}</text>
			</view>
			<view class="left_bottom">
				<text> <slot>有效期：</slot>{{item.dates}}</text>
			</view>
		</view>
		<image src="./ysy.png" class="ysy" v-if="item.status == 1"></image>
		<view class="right" :style="'background:' + (item.status == 0 ? colors:'')">
			<view class="money">￥{{item.sub}}</view>
			<text>满{{item.money}}可用</text>
			<text class="shiyong" :style="{'color': colors}" v-if="item.status == 0" @click="itemClick(item)">去使用</text>
		</view>
	</view>
	</view>
</template>


<script setup>

    const colors = '#e54d42';
    const couponList = [{
						name: '满105减5',
						dates: '2023-07-09 2023-08-02',
						status: 0,
						money: 105,
						sub: 5
					},
					{
						name: '满200减10',
						dates: '2023-07-19 2023-08-22',
						status: 0,
						money: 200,
						sub: 10
					}, {
						name: '满100减10',
						dates: '2023-05-09 2023-06-02',
						status: 1,
						money: 100,
						sub: 10
					},
					{
						name: '满400减20',
						dates: '2023-04-09 2023-05-08',
						status: 1,
						money: 400,
						sub: 20
					}]

    const itemClick = (item) => {
    	uni.showModal({
    		title: '点击优惠券条目',
    		content: '点击优惠券条目 ='+ JSON.stringify(item)
    	})
    }
</script>

<style scoped lang="scss">
	.coupon_box {
		display: flex;
		margin: 28upx 28upx;
		box-shadow: 0upx 0upx 10upx #ddd;
		position: relative;
		border-radius: 10upx;
		overflow: hidden;
	}
	
	.coupon_box .left {
		width: 70%;
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		padding: 20upx;
		float: left;
	}
	
	.coupon_box .left .hui {
		width: 40upx;
		height: 40upx;
		font-size: 20upx;
		color: #fff;
		background-color: #EC1818;
		border-radius: 8upx;
		line-height: 40upx;
		text-align: center;
		display: inline-block;
		transform: translateY(-5upx);
	}
	
	.coupon_box .left .left_top {
		width: 60vw;
		display: block;
		font-size: 26upx;
		font-weight: bold;
	
	}
	
	.left_top .hui_name {
		line-height: 60upx;
		height: 60upx;
		margin-left: 20upx;
		display: inline-block;
	}
	
	.left_bottom {
		font-size: 24upx;
		font-weight: bold;
		color: #333;
		height: 60upx;
		line-height: 60upx;
	}
	
	.coupon_box .right {
		text-align: center;
		height: 160upx;
		width: 180upx;
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		align-content: center;
		background-color: #A8A8A8;
		float: right;
	}
	
	.coupon_box .right .shiyong {
		height: 40upx;
		line-height: 40upx;
		background-color: #fff;
		border-radius: 20upx;
		padding: 0 20upx;
		color: #A8A8A8;
		
	}
	
	.coupon_box .right .money {
		font-size: 45upx;
		color: #fff;
	}
	
	.coupon_box .right text {
		font-size: 24upx;
		color: #fff;
		height: 40upx;
		line-height: 34upx;
	}
	
	.coupon_box .bottom {
		height: 60upx;
		line-height: 60upx;
		display: flex;
		align-content: flex-start;
		font-size: 24upx;
		margin-top: 10upx;
	}
	
	.coupon_box .bottom view {
		margin-right: 20upx;
		color: #888;
		font-weight: bold;
	}
	
	.ysy {
		width: 80upx;
		height: 80upx;
		position: absolute;
		top: 20upx;
		right: 200upx;
	}
</style>