<template>
	<view class="coupon_body">
		<view v-for="(item, index) in couponList" :key="index" class="coupon_box" :class="{'disabled-coupon': item.canUse === 0}">
			<view class="left">
				<view class="left_top">
					<text class="hui">券</text>
					<text class="coupon_name">{{item.name}}</text>
					<text class="hui_name">{{item.minSpend > 0 ? '满'+item.minSpend+'减'+item.discount : '无门槛优惠券'}}</text>
				</view>
				<view class="left_bottom">
					<text>有效日期：{{item.startTime}} - {{item.endTime}}</text>
				</view>
				<image src="./ylq.png" v-if="item.status == 1" class="ylq"></image>
			</view>
			<view class="right">
				<view class="money">￥{{item.discount}}</view>
				<text>{{item.minSpend > 0 ? '满'+item.minSpend+'可用' : '无门槛'}}</text>
			</view>

			<view class="bottom">
				<view class="unavailable-tip" v-if="item.canUse === 0">
					此商品不支持使用该优惠券
				</view>
				<view class="receiveBtn" :style="{background: item.canUse === 0 ? '#cccccc' : (item.couponNum == 0 ? colors:'#fbbd08')}"
					@tap="item.canUse !== 0 && onreceive(item, index)">{{item.couponNum == 0 ? '领取':'立即使用'}}</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { defineProps, defineEmits, ref } from 'vue';
import { getAllClaimCoupon } from "@/api/coupon";

// 定义props
const props = defineProps({
	colors: {
		type: String,
		default: ''
	},
	goodIds: {
		type: Array,
		default: []
	}
});

// 定义事件
const emit = defineEmits(['onReceive']);

// 定义数据
const couponList = ref([]);


const getCouponList = async () => {
	const res = await getAllClaimCoupon({goodIds:props.goodIds});
	if (res.code == 0) {
		couponList.value = res.data;
	}
};

getCouponList()

// 定义方法
const onreceive = (item, index) => {
	emit('onReceive', item, index);
};
</script>

<style lang="scss" scoped>
	.coupon_box {
		margin: 20upx;
		padding: 20upx;
		box-shadow: 0upx 0upx 10upx #ddd;
		position: relative;
		border-radius: 10upx;
		padding-bottom: 10upx;
		overflow: hidden;

	}

	.coupon_box .left {
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		border-bottom: 1upx solid #eee;
		padding-bottom: 20upx;
		position: relative;
	}

	.coupon_box .left .ylq {
		width: 60upx;
		height: 45upx;
		position: absolute;
		top: 0;
		right: 140upx;
	}

	.coupon_box .left .hui {
		width: 40upx;
		height: 40upx;
		font-size: 22upx;
		color: #fff;
		background-color: rgba(255, 84, 110, .8);
		border-radius: 8upx;
		line-height: 40upx;
		text-align: center;
		display: inline-block;
		transform: translateY(-5upx);
	}

	.coupon_box .left .left_top {
		width: 80vw;
		display: block;
		font-size: 26upx;
		font-weight: bold;

	}

	.left_top .coupon_name {
		line-height: 60upx;
		height: 60upx;
		margin-left: 20upx;
		display: inline-block;
		font-size: 30upx;
		font-weight: bold;
		color: #333;
	}

	.left_top .hui_name {
		line-height: 60upx;
		height: 60upx;
		margin-left: 20upx;
		display: inline-block;
		font-size: 28upx;
		color: #666;
	}

	.left_bottom {
		font-size: 24upx;
		font-weight: 500;
		color: #333;
		height: 60upx;
		line-height: 60upx;
	}

	.coupon_box .right {
		position: absolute;
		right: 20upx;
		top: 25upx;
		text-align: center;
	}

	.coupon_box .right .money {
		font-size: 45upx;
		margin-bottom: 10upx;
	}

	.coupon_box .right text {
		font-size: 24upx;
		color: #999;
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
	}

	.receiveBtn {
		position: absolute;
		left: calc(100vw - 86px);

		width: 58px;
		height: 24px;
		line-height: 24px;

		border-radius: 4px;
		background-color: #fa436a;
		color: white !important;

		font-size: 12px;
		text-align: center;
		margin-top: 2px;

	}
	
	.disabled-coupon {
		opacity: 0.6;
	}
	
	.disabled-coupon .receiveBtn {
		cursor: not-allowed;
	}

	// 在样式部分新增样式规则
	.unavailable-tip {
	  color: #e4393c;
	}
</style>