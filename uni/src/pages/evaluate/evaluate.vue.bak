<template>
	<view class="big_box">
		<view style="height: 88rpx;width: 100%;">
			<view class="tabs_box pos_f bgc_fff">
				<view class="font_28 words color_333">
					<text @tap="choose('all')" class="font_bold color_fe5572">{{`全部（${commentInfo.length}）`}}</text>
					<text @tap="choose('pics')" class="font_bold color_fe5572">{{`图文（${num}）`}}</text>
				</view>
			</view>
		</view>
		<view class="content-box">
			<view class="bgc_fff mt_20 cards" v-for="(item, key) in commentInfo" :key="key">
				<view class="flex m_b_24">
					<image class="item_head m_r_16" :src="getUrl(item.user.avatar)"></image>
					<view class="flex-fitem">
						<view class="flex-aic flexr-jsb">
							<text class="color_333 font_28">{{item.user.nickname}}</text>
							<text class="color_999 font_24">{{formatTimeToStr(item.CreatedAt, 'yyyy-MM-dd')}}</text>
						</view>
						<uni-rate size="16" :readonly="true" active-color="#fe5572" :value="item.rating" />
					</view>
				</view>
				<view style="padding: 0 16rpx 0 80rpx;" class="boxs_bb">
					<view class="flex-aic flexr-jsb color_999 font_24 bgc_f8f8f8 evaluate_num boxs_bb">
						<text>已购</text>
					</view>
					<view class="text_pre_wrap m_b_24 m_t_24 color_333 font_28">{{item.content}}</view>
          <view class="pics_grid">
            <view
                v-for="(pic, index) in item.pics"
                :key="index"
                class="pic_item"
                @tap="previewImage(pic, index, item.pics)"
            >
              <image
                  :src="getUrl(pic)"
                  class="evaluate_pic_img"
                  mode="aspectFill"
                  @error="onImageError(index)"
              />
            </view>
          </view>
					<view v-if="item.shopReply !== ''" class="shop-reply m_t_24  color_999 font_24 bgc_f8f8f8 evaluate_num boxs_bb">
						<p >商家回复：{{item.shopReply}}</p>
					</view>
				</view>
			</view>
		</view>
	</view>

</template>

<script setup>
	import { ref } from "vue";
	import { findComment } from "@/api/comment.js"
	import evaluateGridImg from './evaluate-img.vue'
	import { formatTimeToStr } from "@/utils/date.js"
	import {
		onLoad,
	} from '@dcloudio/uni-app'
    import {getUrl} from "@/utils/url.js"
	const commentInfo = ref([])
	const num = ref(0)
	let ID = ""
	const findFunc = async(params) => {
		num.value = 0
		const res = await findComment({ID:params})
		if (res.code === 0) {
			commentInfo.value = res.data
			commentInfo.value.map((i) =>{
				if(i.feedbackPics) {
					++num.value
				}
			})
		}
	}
	onLoad((options) => {
		ID = options.goodsID
		setTimeout(async () => {
			// 获取评论
			findFunc(ID)
		}, 500)
	})

	const choose = (params) => {
		const temp = commentInfo.value;
		const filteredItems = temp.filter(item => item.feedbackPics !== null);
		if (params === 'all') {
			findFunc(ID)
		} else {
			commentInfo.value = filteredItems;
		}
}

  const pics = ref([])

  // 新增图片预览方法
  const previewImage = (currentPic, index, allPics) => {
    const urls = allPics.map(pic => getUrl(pic));

    uni.previewImage({
      current: getUrl(currentPic), // 当前图片
      urls: urls, // 所有图片
      fail: (err) => {
        console.error('图片预览失败:', err)
        uni.showToast({
          title: '图片加载失败',
          icon: 'error'
        })
      }
    })
  }

  // 图片加载错误处理
  const onImageError = (index) => {
    const failedPic = pics.value[index];
    console.error(`图片加载失败:`, {
      originalUrl: failedPic.url,
      processedUrl: getUrl(failedPic.url),
      index: index
    });
  }

</script>

<style lang="scss" scoped>
	.big_box{
    height: 100vh;
		background-color: #f5f5f5;
		.words{
			text-align: center;
			font-size: 32rpx;
			.line_active {
				width: 100%;
				background-color: #fe5572;
			}
		}
	}
	.item_head {
		width: 64rpx;
		height: 64rpx;
		border-radius: 50%;
	}

	.evaluate_num {
		margin-top: 16rpx;
		padding: 6rpx 16rpx;
	}
	.content-box{
		margin: 40rpx;
		.shop-reply{
			display: flex;
			justify-content: space-between;
		}
		.cards{
			border-radius: 20rpx;
			padding: 20rpx;
		}
	}
		page {
			background-color: #F8F8F8;
		}

		.tabs_box {
			/* #ifdef H5 */
			top: var(--window-top);
			/* #endif */
			/* #ifndef H5 */
			top: 0;
			/* #endif */
			z-index: 1;
			left: 0;
			right: 0;
			padding: 24rpx 0;
		}

		.item_box {
			padding: 24rpx 32rpx;
		}
		.mt_20{
			margin-top: 40rpx;
		}

  .pics_grid {
    display: grid;
    grid-template-columns: repeat(3, 100rpx); /* 改为100rpx匹配pic_item */
    gap: 8rpx;
    justify-content: flex-start;
    margin-bottom: 16rpx; /* 添加底部间距 */
  }

  .pic_item {
    width: 100rpx;
    height: 100rpx;
    border-radius: 8rpx;
    overflow: hidden;
    background-color: #f5f5f5;
  }

  .evaluate_pic_img {
    width: 100%;
    height: 100%;
    border-radius: 8rpx;
  }
</style>
