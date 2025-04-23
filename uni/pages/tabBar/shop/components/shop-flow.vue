<template>
  <view class="flow_box">
      <view class="content">
        <view v-for="(item, index) in data" :key="index" @tap="tapClick(item)">
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
          <up-divider></up-divider>

        </view>
      </view>
  </view>
</template>
<script setup>
import { ref, reactive } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import {getGoodList} from "../../../../api/homePage";
import {getUrl} from "@/utils/url.js"
let data = ref([])
const init = async () => {
  const params = {
    page: 1,
    pageSize: 20
  }
  const res = await getGoodList(params)
  if (res.code === 0) {
    // res.data.list只保留image和title和ID
    data.value = res.data.list
  }
}

onShow(() => {
  init()
})

const tapClick = (item) => {
  console.log(item)
  uni.navigateTo({
    url: `/pages/goodsDetails/goodsDetails?id=${item.ID}`,
  })
}
</script>
<style lang="scss" scoped>
.flow_box {
  padding: 0 32rpx 100rpx;
}
.item {
  padding-top: 16rpx;
  .title {
    max-height: 72rpx;
    line-height: 28rpx;
  }
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
