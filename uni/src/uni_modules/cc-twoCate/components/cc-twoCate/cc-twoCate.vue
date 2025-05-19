<template>
	<view class="towstage" @touchmove.stop.prevent="moveHandle">
		<scroll-view :class="['left',showHeader == false?'left_h5':'left_app']" :scroll-y="true" scroll-with-animation
			:scroll-top="scrollTop">
			<view v-for="(item, index) in classList" :key="index"
				:class="'classTitle ' + (current == index ? 'actives':'')" @tap="setCurrent(item, index)"
				:style="{color:current == index ?colors:''}">
				<text>{{item.name}}</text>
				<view class="classactive" v-if="current == index" :style="'background-color:' + colors">
				</view>
			</view>
		</scroll-view>
		<scroll-view :class="['right',showHeader == false?'right_h5':'right_app']" :scroll-y="true"
			:scroll-with-animation="true">


			<view class="goods_list" v-if="dataList.length !== 0">
				<view v-for="(item, index) in dataList"  @tap="jumpDetail(item)" :key="index" class="goods_item">
					<image :src="item.img"></image>
					<view class="goods_right">
						<view class="goods_name">{{item.title}}</view>
						<view class="numbers">
							<text v-if="item.youhui == true">优惠券</text>
							<text v-if="item.baoyou == true">包邮</text>
						</view>
						<view class="price">
							<text class="money">￥{{item.money}}</text>
							<text class="hx_money">￥{{item.hmoney}}</text>

						</view>
					</view>
				</view>
				<view class="nodata">—— 到底啦 ——</view>
			</view>
			<view class="zanwu" v-if="dataList.length == 0">
				<text class="iconfont icon-none" style="font-size: 32upx"></text>
				<text style="font-size: 24upx">暂无商品</text>
			</view>
		</scroll-view>

	</view>
</template>

<script>
	import Vue from 'vue'
	
	Vue.mixin({
		methods: {
			setData: function(obj, callback) {
				let that = this;
				const handleData = (tepData, tepKey, afterKey) => {
					tepKey = tepKey.split('.');
					tepKey.forEach(item => {
						if (tepData[item] === null || tepData[item] === undefined) {
							let reg = /^[0-9]+$/;
							tepData[item] = reg.test(afterKey) ? [] : {};
							tepData = tepData[item];
						} else {
							tepData = tepData[item];
						}
					});
					return tepData;
				};
				const isFn = function(value) {
					return typeof value == 'function' || false;
				};
				Object.keys(obj).forEach(function(key) {
					let val = obj[key];
					key = key.replace(/\]/g, '').replace(/\[/g, '.');
					let front, after;
					let index_after = key.lastIndexOf('.');
					if (index_after != -1) {
						after = key.slice(index_after + 1);
						front = handleData(that, key.slice(0, index_after), after);
					} else {
						after = key;
						front = that;
					}
					if (front.$data && front.$data[after] === undefined) {
						Object.defineProperty(front, after, {
							get() {
								return front.$data[after];
							},
							set(newValue) {
								front.$data[after] = newValue;
								that.$forceUpdate();
							},
							enumerable: true,
							configurable: true
						});
						front[after] = val;
					} else {
						that.$set(front, after, val);
					}
				});
				isFn(callback) && this.$nextTick(callback);
			}
		}
	});
	
	export default {
		data() {
			return {
				scrollTop: 0,


			};
		},


		props: {
			colors: {
				type: String
			},

			classList: {
				type: Array
			},
			dataList: {
				type: Array
			},
			current: {
				type: Number,
				default: 0
			},
			showHeader: { //用来控制在APP下样式
				type: Boolean,
				default: false
			},

		},
		methods: {
			moveHandle() {},
			setCurrent(item, index) {
				//设置tab
				let scrollTop = index * 20;
				if (index <= 2) {
					scrollTop = 0;
				}
				this.setData({
					current: index,
					scrollTop: scrollTop,
				});
				this.$emit('cateClick', item)
			},



			jumpDetail(e) {
				
				this.$emit("itemClick",e);
			}

		}
	};
</script>
<style scoped lang="scss">
	@import "./twoCate.scss";
</style>