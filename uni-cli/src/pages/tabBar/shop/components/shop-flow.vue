<template>
  <view class="flow_box">
    <view v-if="loading" class="loading-container">
      <view class="loading-spinner"></view>
      <text class="loading-text">加载中...</text>
    </view>
    <view v-else-if="data.length === 0" class="empty-container">
      <image src="/static/empty-cart.png" class="empty-image"></image>
      <text class="empty-text">暂无商品</text>
      <button class="empty-button" @tap="init">刷新</button>
    </view>
    <view v-else class="content">
      <view v-for="(item, index) in data" :key="index" @tap="tapClick(item)" class="grid-item">
        <view class="gird-item-card">
          <image :src="getUrl(item.imageUrl)" class="img" mode="aspectFill"></image>
          <view class="info-container">
            <text class="title">{{ item.title }}</text>
            <view class="flexr-jsb flex-aic desc">
              <text class="price"><text class="icons">￥</text>{{ item.price && item.price/100 }}</text>
              <text class="sale-num">已售：{{item.saleNum}}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>
<script setup>
import { ref, reactive } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import {getGoodList} from "@/api/homePage";
import {getUrl} from "@/utils/url.js"

let data = ref([])
const loading = ref(false)

const init = async () => {
  loading.value = true
  const params = {
    page: 1,
    pageSize: 20
  }
  try {
    const res = await getGoodList(params)
    if (res.code === 0) {
      data.value = res.data.list
    }
  } catch (error) {
    console.error('获取商品列表失败:', error)
  } finally {
    loading.value = false
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
  padding: 0 24rpx 100rpx;
  min-height: 40vh;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300rpx;

  .loading-spinner {
    width: 60rpx;
    height: 60rpx;
    border: 6rpx solid #f3f3f3;
    border-top: 6rpx solid #FE5572;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  .loading-text {
    margin-top: 20rpx;
    color: #666;
    font-size: 28rpx;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
}

.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;

  .empty-image {
    width: 200rpx;
    height: 200rpx;
    margin-bottom: 30rpx;
  }

  .empty-text {
    font-size: 30rpx;
    color: #999;
    margin-bottom: 30rpx;
  }

  .empty-button {
    background-color: #FE5572;
    color: white;
    font-size: 28rpx;
    padding: 10rpx 60rpx;
    border-radius: 40rpx;
  }
}

.content {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24rpx;
  margin-top: 20rpx;

  .grid-item {
    position: relative;
    border-radius: 16rpx;
    overflow: hidden;
    transition: transform 0.2s;
    box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);

    &:active {
      transform: scale(0.98);
    }

    .gird-item-card {
      background-color: #fff;
      border-radius: 16rpx;
      overflow: hidden;

      .img {
        width: 100%;
        height: 340rpx;
        border-radius: 16rpx 16rpx 0 0;
        transition: transform 0.3s;

        &:hover {
          transform: scale(1.05);
        }
      }

      .info-container {
        padding: 16rpx 20rpx;
      }

      .title {
        width: 100%;
        font-size: 28rpx;
        line-height: 40rpx;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        height: 80rpx;
        text-align: left;
        color: #333;
        margin-bottom: 10rpx;
      }

      .desc {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;

        .price {
          color: #FE5572;
          font-size: 34rpx;
          font-weight: bold;
          text-align: left;

          .icons {
            font-size: 24rpx;
          }
        }

        .sale-num {
          color: #999;
          font-size: 24rpx;
        }
      }
    }
  }
}
</style>
