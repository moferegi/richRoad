<template>
	<view class="sunui-uploader-bd">
		<view class="sunui-uploader-files">
			<block v-for="(item, index) in images" :key="index">
				<view class="sunui-uploader-file" :class="[item.progress < 100 ? 'sunui-uploader-file-status' : '']"
					@click="previewImage(index)">
					<!-- step1.这里修改服务器返回字段 ！！！ -->
					<!-- <view class="sunui-uploader-file"  @click="previewImage(index)"> -->
					<image class="sunui-uploader-img" :style="upload_img_wh" :src="item.path" mode="aspectFill" />
					<view class="sunui-img-removeicon right" @click.stop="removeImage(index)" v-show="upimg_move"><text
							class="close_box">×</text></view>
					<view class="sunui-loader-filecontent" v-if="item.progress < 100">{{ item.progress }}%</view>
				</view>
			</block>
			<image v-if="images.length <= upload_count - 1" class="sunui-uploader-inputbox" @click="chooseImage"
				:style="upload_img_wh" src="http://www.liwanying.top/applate-icon/upload-icon.png" mode=""></image>
		</view>
	</view>
</template>
<script setup>
	import {
		uploadFile
	} from '@/js_sdk/lyz-aloss-uploader/lyz-aloss-uploader.js'
	import configData from '@/common/config.js'
	import {
		getOssImgName
	} from '@/common/util.js'
	import {
		ref
	} from 'vue'
	let images = ref([])
	const emit = defineEmits(['change'])
	let props = defineProps({
		// 上传样式宽高
		upload_img_wh: {
			type: String,
			default: 'width:216rpx;height:216rpx;'
		},
		// 上传数量
		upload_count: {
			type: [Number, String],
			default: 3
		},
		// 是否显示删除
		upimg_move: {
			type: Boolean,
			default: true
		},
		// 服务器预览图片
		upimg_preview: {
			type: Array,
			default: () => {
				return [];
			}
		},
		// 服务器返回预览(看服务器卡顿情况设定)
		upimg_delaytime: {
			type: [Number, String],
			default: 300
		}
	})
	const init = (imgs) => {
		imgs.forEach(img => {
			images.value.push({
				path: img.src,
				url: img.url,
				uploadTask: null,
				progress: 100,
				uploadInfo: null,
				status: 6
			})
		})
		toEmit()
	}
	const toEmit = () => {
		let list = []
		let upSuccessCount = 0
		let upLoadingCount = 0
		images.value.forEach(img => {
			if (img.status == 2 || img.status == 6) {
				upSuccessCount++
				list.push(configData.uploadImageUrl + img.fileKey)
			}
			if (img.status == 1) {
				upLoadingCount++
			}
		})
		let data = {
			details: images.value,
			list:list,
			upLoadingCount: upLoadingCount,
			upSuccessCount: upSuccessCount,
			imgCount: images.value.length
		}
		emit('change', data);
	}
	/* ************************上传方法***************************************** */
	const upImage = () => {
		let promises = []
		images.value.forEach(img => {
			if (img.status == 0) {
				img.status = 1
				let promise = new Promise((resolve, reject) => {
					let uploadTask = uploadFile(img.path, 'png', data => {
						toEmit()
						resolve(data)
					}, reject)
					img.uploadTask = uploadTask
					uploadTask.onProgressUpdate(async (res) => {
						img.progress = res.progress
						img.uploadInfo = res
					});
				}).then(res => {
					img.progress = 100
					img.status = 2
					img.originalFileName = getOssImgName() + '.png'
					img.fileKey = res
				}).catch(e => {
					img.status = 99
				})
				promises.push(promise)
			}
		})

		uni.showLoading({
			title: `正在上传...`
		});
		// console.log("promises===>",promises)
		Promise.all(promises)
			.then(res => {
				uni.hideLoading()
				toEmit()
			}).catch(function(e) {
				uni.hideLoading()
			})
	}
	const chooseImage = () => {
		uni.chooseImage({
			count: 1,
			sizeType: ['compressed', 'original'],
			sourceType: ['camera'],
			success: function(res) {
				for (let i = 0, len = res.tempFiles.length; i < len; i++) {
					images.value.push({
						path: res.tempFiles[i].path,
						fileSize: res.tempFiles[i].size,
						fileEtag: '',
						originalFileName: '',
						fileKey: '',
						uploadTask: null,
						progress: 0,
						uploadInfo: null,
						status: 0
					})
				}
				upImage()
			},
			fail: function(err) {
				console.log(err);
			}
		});
	}
	const previewImage = (index) => {
		// console.log(index)
		let list = []
		images.value.forEach(img => {
			list.push(img.path)
		})
		// console.log(list)
		uni.previewImage({
			urls: list,
			current: index,
		})
	}
	const removeImage = (idx) => {
		images.value[idx].uploadTask.abort();
		images.value.splice(idx, 1);
		toEmit();
	}
</script>
<style lang="scss" scoped>
	.sunui-uploader-img {
		display: block;
		border-radius: 6rpx;
	}

	.sunui-uploader-input {
		position: absolute;
		z-index: 1;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		opacity: 0;
	}

	.sunui-uploader-inputbox {
		position: relative;
		box-sizing: border-box;
		background-color: #FFFFFF;
	}

	.sunui-img-removeicon {
		position: absolute;
		color: #fff;
		width: 40upx;
		height: 40upx;
		background-color: rgba(0, 0, 0, 0.7);
		border-bottom-left-radius: 24rpx;
		border-top-right-radius: 6rpx;

		&.right {
			top: 0;
			right: 0;
		}

		.close_box {
			position: absolute;
			font-size: 30rpx;
			top: -0rpx;
			line-height: 30rpx;
			right: 10rpx;
		}
	}

	.sunui-uploader-file {
		position: relative;
		margin-right: 24rpx;
		// margin-bottom: 16rpx;
	}

	.sunui-uploader-file-status:before {
		content: ' ';
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		background-color: rgba(0, 0, 0, 0.5);
	}

	.sunui-loader-filecontent {
		position: absolute;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		color: #fff;
		z-index: 9;
	}

	.sunui-uploader-bd {
		// padding: 26rpx;
		margin: 0;
	}

	.sunui-uploader-files {
		display: flex;
		flex-wrap: wrap;
		// justify-content: space-around;
	}

	.sunui-uploader-file:nth-child(3n + 0) {
		margin-right: 0;
	}

	.sunui-uploader-inputbox>view {
		text-align: center;
	}

	.sunui-uploader-file-status:after {
		content: ' ';
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		background-color: rgba(0, 0, 0, 0.5);
	}

	.sunui-uploader-hover {
		box-shadow: 0 0 0 #e5e5e5;
		background: #e5e5e5;
	}
</style>
