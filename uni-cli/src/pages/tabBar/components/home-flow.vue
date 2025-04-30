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
import { getGoodList } from '@/api/homePage.js'
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
  margin-top: 20rpx;
  gap: 24rpx;
  .grid-item {
    background-color: #ffffff;
    border-radius: 16rpx;
    box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.05);
    overflow: hidden;
    transition: transform 0.3s ease;
    
    &:active {
      transform: scale(0.98);
    }
    
    .gird-item-card {
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 20rpx;
      .img {
        width: 100%;
        height: 340rpx;
        border-radius: 0;
        object-fit: cover;
      }

      .title {
        width: 90%;
        font-size: 28rpx;
        max-height: 80rpx;
        overflow: hidden;				//溢出内容隐藏
        text-overflow: ellipsis;		//文本溢出部分用省略号表示
        display: -webkit-box;			//特别显示模式
        -webkit-line-clamp: 2;			//行数
        line-clamp: 2;
        -webkit-box-orient: vertical;	//盒子中内容竖直排列
        margin: 16rpx 0;
        color: #333;
        text-align: left;
        padding: 0 16rpx;
      }
      .desc{
        width: 92%;
        padding: 0 16rpx 16rpx;
        .price {
          color: #FF6A6A;
          font-size: 34rpx;
          font-weight: bold;
          .icons {
            font-size: 22rpx !important;
          }
        }
      }
    }
  }
}
</style>
