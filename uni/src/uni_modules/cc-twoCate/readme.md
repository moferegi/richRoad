# cc-twoCate


#### 使用方法 
```使用方法
<!-- colors：选中颜色 current：列表选择序列  左边分类数据 dataList：右边列表数据 @cateClick：左边分类点击 @itemClick：右边条目点击 -->
<cc-twoCate colors="#fa436a" :current="currentTwo" :classList="categoryList" :dataList="dataList"
			:hideShow="true" @cateClick="cateClick" @itemClick="itemClick"></cc-twoCate>
						
```

#### HTML代码实现部分
```html
<template>
	<view class="content">

		<!-- colors：选中颜色 current：列表选择序列  左边分类数据 dataList：右边列表数据 @cateClick：左边分类点击 @itemClick：右边条目点击 -->
		<cc-twoCate colors="#fa436a" :current="currentTwo" :classList="categoryList" :dataList="dataList"
			:hideShow="true" @cateClick="cateClick" @itemClick="itemClick"></cc-twoCate>

	</view>

</template>

<script>
	export default {
		components: {

		},
		data() {
			return {


				//分类列表
				categoryList: [{
					id: 1,
					name: '今日推荐',

				}, {
					id: 2,
					name: '每日特惠',

				}, {
					id: 3,
					name: '进口水果',

				}, {
					id: 4,
					name: '蔬菜豆制品',

				}, {
					id: 5,
					name: '肉禽蛋',

				}, {
					id: 6,
					name: '海鲜水产',

				}, {
					id: 7,
					name: '粮油调味',

				}, {
					id: 8,
					name: '熟食卤味',

				}, {
					id: 9,
					name: '冻品面点',

				}, {
					id: 10,
					name: '乳品烘培',

				}],
				// 选中序列
				currentTwo: 0,

				// 右边列表数据
				dataList: [{
						title: '精美五彩水果',
						goods_id: 201,
						money: '35.90',
						number: 1,
						hmoney: '45.90',
						img: 'https://cdn.pixabay.com/photo/2019/05/14/17/11/fruit-4202929_1280.png',
						youhui: true,
						baoyou: false,
						status: 1, //商品过期状态  0正常  1已失效
						stock: 600,

					},
					{
						title: '智利进口新鲜车厘子',
						goods_id: 202,
						money: '29.9',
						number: 75,
						hmoney: '39.90',
						img: 'https://cdn.pixabay.com/photo/2016/10/30/18/01/apple-1783882_1280.png',
						youhui: false,
						baoyou: true,
						status: 0, //商品过期状态  0正常  1已失效
						stock: 100,
					},
					{
						title: '伊犁冰淇淋',
						type: 3,
						goods_id: 203,
						money: '152.00 ',
						number: 1,
						hmoney: '162.00',
						img: 'https://cdn.pixabay.com/photo/2016/03/23/15/00/ice-cream-1274894_1280.jpg',
						youhui: true,
						baoyou: true,
						status: 0, //商品过期状态  0正常  1已失效
						stock: 200,
					},
					{
						title: '黑美人西瓜5kg',
						type: 6,
						goods_id: 204,
						money: '52.00 ',
						number: 1,
						hmoney: '99.00 ',
						youhui: false,
						baoyou: false,
						stock: 100,
						img: 'https://cdn.pixabay.com/photo/2017/06/02/18/24/watermelon-2367029_1280.jpg',
						status: 0, //商品过期状态  0正常  1已失效
					},
					{
						title: '4杯鸡蛋布甸（双层）',
						type: 4,
						goods_id: 205,
						money: '25.80',
						number: 1,
						hmoney: 35.00,
						img: '/static/images/goods/four.jpg',
						youhui: true,
						baoyou: false,
						stock: 500,
						status: 0, //商品过期状态  0正常  1已失效
					},
					{
						title: '云南草莓夏季草莓新鲜水果3斤礼盒装',
						type: 5,
						goods_id: 206,
						money: '59.90',
						number: 200,
						hmoney: '70.90',
						youhui: true,
						baoyou: true,
						img: '/static/images/goods/five.jpg',
						status: 0, //商品过期状态  0正常  1已失效
						stock: 140,
					}
				],

			};
		},


		methods: {

			cateClick(item) {

				uni.showModal({
					title: "点击分类条目",
					content: '点击分类条目 = ' + JSON.stringify(item)
				})

			},

			itemClick(e) {

				uni.showModal({
					title: "点击右边商品条目",
					content: '点击右边商品条目 = ' + JSON.stringify(e)
				})
			}

		}
	}
</script>

<style lang="scss" scoped>
	.content {
		display: flex;
		flex-direction: column;

	}
</style>



```

