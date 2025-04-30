<template>
	<view class="waterfalls-flow">
		<view v-for="(item,index) in data.column" :key="index" class="waterfalls-flow-column"
			:id="`waterfalls_flow_column_${index+1}`" :style="{'width':w,'margin-left':index==0?0:m}">
			<view class="column-value" v-for="(item2,index2) in data[`column_${index+1}_values`]" :key="index2" @tap.stop="wapperClick(item2)">
				<view class="inner" v-if="data.seat==1">
					<!-- #ifdef MP-WEIXIN -->
					<slot :name="`slot${item2.index}`"></slot>
					<!-- #endif -->
					<!-- #ifndef MP-WEIXIN -->
					<slot v-bind="item2"></slot>
					<!-- #endif -->
				</view>
				<image :class="['img',{'img-error':!item2[data.imageKey]}]" :src="item2[data.imageKey]" mode="widthFix"
					@load="imgLoad(item2)" @error="imgError(item2)"></image>
				<view class="inner" v-if="data.seat==2">
					<!-- #ifdef MP-WEIXIN -->
					<slot :name="`slot${item2.index}`"></slot>
					<!-- #endif -->
					<!-- #ifndef MP-WEIXIN -->
					<slot name="card" :data="item2"></slot>
					<!-- #endif -->
				</view>
			</view>
		</view>
	</view>
	<view style="width: 100%;" :style="{
		height: minHeightResNum +'px'
	}">

	</view>
</template>
<script setup>
	import {
		ref,
		reactive,
		watch,
		computed,
		getCurrentInstance,
		onMounted
	} from 'vue';
	const emit = defineEmits(['tapClick'])
	const _this = getCurrentInstance();
	const data = reactive({});
	let minHeightResNum = ref(100)
	const props_data = defineProps({
		value: Array,
		column: { // 列的数量 
			type: [String, Number],
			default: 2
		},
		columnSpace: { // 列之间的间距 百分比
			type: [String, Number],
			default: 2
		},
		imageKey: { // 图片key
			type: [String],
			default: 'image'
		},
		seat: { // 文本的位置，1图片之上 2图片之下
			type: [String, Number],
			default: 2
		}
	});
	// 数据赋值
	data.list = props_data.value ? props_data.value : [];
	data.column = props_data.column < 2 ? 2 : props_data.column;
	data.columnSpace = props_data.columnSpace <= 5 ? props_data.columnSpace : 5;
	data.imageKey = props_data.imageKey;
	data.seat = props_data.seat;
	// 计算列宽
	const w = computed(() => {
		const column_rate = `${100 / data.column - (+data.columnSpace)}%`;
		return column_rate;
	})
	// 计算margin
	const m = computed(() => {
		const column_margin =
			`${(100-(100 / data.column - (+data.columnSpace)).toFixed(5)*data.column)/(data.column-1)}%`;
		return column_margin;
	})
	// 每列的数据初始化
	for (let i = 1; i <= data.column; i++) {
		data[`column_${i}_values`] = [];
	}
	// 获取最小值的对象
	const getMin = (a, s) => {
		let m = a[0][s];
		let mo = a[0];
		for (var i = a.length - 1; i >= 0; i--) {
			if (a[i][s] < m) {
				m = a[i][s];
			}
		}
		mo = a.filter(i => i[s] == m);
		setTimeout(() => {
				getMaxColumnHeight()
		}, 1000)
		return mo[0];
	}
	// 计算每列的高度
	function getMinColumnHeight() {
		return new Promise(resolve => {
			const heightArr = [];
			for (let i = 1; i <= data.column; i++) {
				const query = uni.createSelectorQuery().in(_this);
				query.select(`#waterfalls_flow_column_${i}`).boundingClientRect(data => {
					heightArr.push({
						column: i,
						height:  data ? data.height : 0
					});
				}).exec(() => {
					if (data.column <= heightArr.length) {
						resolve(getMin(heightArr, 'height'));
					}
				});
			}
		})
	};
	//延时计算列的最高高度
	const getMaxColumnHeight = () => {
		const heightArr = [];
		for (let i = 1; i <= data.column; i++) {
			const query = uni.createSelectorQuery().in(_this);
			query.select(`#waterfalls_flow_column_${i}`).boundingClientRect(res => {
				heightArr.push({
					column: i,
					height: res ? res.height : 0
				});
			}).exec(() => {
				if (data.column <= heightArr.length) {
					minHeightResNum.value = heightArr[0].height > heightArr[1].height ? heightArr[0].height : heightArr[1].height
				}
			});
		}
	};
	async function initValue(i) {
		if (i >= data.list.length) return false;
		const minHeightRes = await getMinColumnHeight();
		data[`column_${minHeightRes.column}_values`].push({
			...data.list[i],
			index: i
		});
		
	}
	onMounted(() => {
		initValue(0);		
	})
	// 图片加载完成
	function imgLoad(item) {
		const i = item.index;
		initValue(i + 1);
	}
	// 图片加载失败
	function imgError(item) {
		const i = item.index;
		initValue(i + 1);
		item[data.imageKey] = null;
	}
	// 单项点击事件
	const wapperClick = (item) => {
		emit('tapClick', JSON.parse(JSON.stringify(item)))
	}
	// 监听数据的变化
	watch(() => props_data.value, (newValue, oldValue) => {
		const oldLength = oldValue ? oldValue.length : 0;
		data.list = newValue;
		if (oldLength > 0) initValue(oldLength);
	}, {
		immediate: true
	});
</script>
<style lang="scss" scoped>
	.waterfalls-flow {
		&-column {
			float: left;
		}
	}

	.column-value {
		width: 100%;
		font-size: 0;
		border: 2rpx solid #efefef;
		border-radius: 12rpx;
		padding: 20rpx;
		box-sizing: border-box;
		margin-bottom: 20rpx;

		.inner {
			font-size: 30rpx;
		}

		.img {
			width: 100%;
			border-radius: 12rpx;

			&-error {
				background: #f2f2f2 url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAC4AAAAiAQMAAAAatXkPAAAABlBMVEUAAADMzMzIT8AyAAAAAXRSTlMAQObYZgAAAIZJREFUCNdlzjEKwkAUBNAfEGyCuYBkLyLuxRYW2SKlV1JSeA2tUiZg4YrLjv9PGsHqNTPMSAQuyAJgRDHSyvBPwtZoSJXakeJI9iuRLGDygdl6V0yKDtyMAeMPZySj8yfD+UapvRPj2JOwkyAooSV5IwdDjPdCPspe8LyTl9IKJvDETKKRv6vnlUasgg0fAAAAAElFTkSuQmCC) no-repeat center center;
			}
		}
	}
</style>
