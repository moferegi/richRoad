# cc-multipleBtn


#### 使用方法 
```使用方法
<!-- colors:设置颜色  remarkList：标签数组 @click：标签点击 -->
<cc-multipleBtn :colors="colors" :remarkList="remarkList" @click="tagsClick"></cc-multipleBtn>
```

#### HTML代码实现部分
```html
<template>
	<view class="cencal_order">
		<view class="remark">
			<textarea maxlength="-1" placeholder="请在此处输入评价" placeholder-class="textarea_p"></textarea>
		</view>

		<!-- colors:设置颜色  remarkList：标签数组 @click：标签点击 -->
		<cc-multipleBtn :colors="colors" :remarkList="remarkList" @click="tagsClick"></cc-multipleBtn>


		<view class="btns" :style="{background: '#fbbd08'}">
			确认提交
		</view>
	</view>
</template>

<script>
	var app = getApp();
	export default {
		data() {
			return {
				colors: '#f37b1d',
				remarkList: [{
					name: '商品品质好'
				}, {
					name: '性价比高'
				}, {
					name: '态度好'
				}, {
					name: '价格合理'
				}, {
					name: '做工不错'
				}, {
					name: '物流时间长'
				}, {
					name: '价格优惠低'
				}, {
					name: '其他原因'
				}],
				data: ""
			};
		},

		components: {},
		props: {},


		/**
		 * 用户点击右上角分享
		 */
		onShareAppMessage: function() {},
		methods: {
			tagsClick(tagArr) {

				console.log("arr = ", JSON.stringify(tagArr));

				uni.showModal({
					title: '选中的标签',
					content: '选中的标签 = ' + JSON.stringify(tagArr)
				})

			}

		}
	};
</script>
<style lang="scss" scoped>
	page {
		background-color: #F5F5FA;
	}

	.cencal_order {
		padding: 20upx 4%;
		background-color: #fff;
	}

	.remark {
		background-color: #F5F5F5;
		border-radius: 10upx;
		height: 28vw;
		padding: 20upx;
		margin-top: 20upx;

		textarea {
			font-size: 26upx;
			color: #797979;
		}
	}

	.textarea_p {
		font-size: 24upx;
		color: #797979;
	}

	

	.btns {
		width: 100%;
		height: 80upx;
		line-height: 80upx;
		font-size: 30upx;
		color: #333;
		text-align: center;
		margin-top: 100upx;
		border-radius: 8upx;
		margin-bottom: 100px;
	}
</style>



```