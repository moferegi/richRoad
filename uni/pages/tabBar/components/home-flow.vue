<template>
  <view class="flow_box">
        <view class="content">
          <view v-for="(item, index) in flowData" :key="index" @tap="goTo(item)">
            <view class="grid-item">
              <view class="gird-item-card">
                <image :src="getUrl(item.imageUrl)" class="img"></image>
                <p class="title"> {{ item.title }}</p>
                <view class="flexr-jsb flex-aic desc">
                  <p class="price"><span class="icons">￥</span>{{ item.price && item.price/100 }}</p>
                  <p class="color_b5b5b5 font_24">已售：{{item.saleNum}}</p>
                </view>
              </view>
            </view>
          </view>
      </view>
  </view>
</template>
<script setup>
import {ref } from 'vue';
import {getUrl} from "@/utils/url.js"
import { getGoodList } from '/api/homePage.js'
// props接收父组件传递过来的数据
const props = defineProps({
  flowData: {
    type: Array,
    default: []
  },
})

	const goTo = (item) => {
	  uni.navigateTo({
	    url: `/pages/goodsDetails/goodsDetails?id=${item.ID}`,
	  })
	}

</script>
<style lang="scss" scoped>
.flow_box {
  padding: 0 20rpx;
}

.item {
  padding-top: 16rpx;
}

.content {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  text-align: center;
  margin-top: 30rpx;
  gap: 20rpx;
  .grid-item {
    border: 2rpx solid whitesmoke;
    border-radius: 10rpx;
    .gird-item-card {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 20rpx;
      .img {
        width: 270rpx;
        height: 320rpx;
        border-radius: 20rpx;
      }

      .title {
        width: 100%;
        font-size: 30rpx;
        max-hight: 30rpx;
        overflow: hidden;				//溢出内容隐藏
        text-overflow: ellipsis;		//文本溢出部分用省略号表示
        display: -webkit-box;			//特别显示模式
        -webkit-line-clamp: 1;			//行数
        line-clamp: 1;
        -webkit-box-orient: vertical;	//盒子中内容竖直排列
        margin: 20rpx 0;
      }
      .desc{
        width: 100%;
        .price {
          color: #FE5572;
          font-size: 32rpx;
          .icons {
            font-size: 20rpx !important;
          }
        }
      }
    }
  }
}
</style>
