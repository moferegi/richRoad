<template>
  <view class="shop-container">
    <shop-goods-list></shop-goods-list>
    <gva-divider text="猜你喜欢"></gva-divider>
    <noPaginGridGoodList style="margin-top: 40rpx;" v-if="flowData" :goodsList="flowData"></noPaginGridGoodList>
    <view style="height: 80rpx;"></view>
  </view>
</template>

<script setup>
import noPaginGridGoodList from '@/components/good-list/no-pagin-grid-good-list.vue'
import shopGoodsList from './components/shop-goods-list.vue'
import {ref} from "vue";
import {getGoodList} from '@/api/homePage.js'

const flowData = ref([])
const params = ref({
  page: 1,
  pageSize: 10,
  categoryID: 0
})
const init = async () => {
  const res = await getGoodList(params.value)
  if (res.code === 0 && res.data.list.length) {
    flowData.value.push(...res.data.list)
  }
}
init()


</script>

<style lang="scss">
page {
  background-color: #f7f7f7;
  min-height: 100vh;
}

.shop-container {
  padding: 20rpx;
  min-height: 100vh;
}

.recommend-section {
  margin-top: 32rpx;
  background: linear-gradient(135deg, #ffffff 0%, #fafafa 100%);
  border-radius: 32rpx 32rpx 0 0;
  overflow: hidden;
  padding: 24rpx;
  box-shadow: 0 8rpx 32rpx rgba(255, 76, 125, 0.15);
  backdrop-filter: blur(10rpx);
  border: 1rpx solid rgba(255, 255, 255, 0.2);
}

.section-title {
  font-size: 36rpx;
  color: #333;
  font-weight: 600;
  text-align: center;
  padding: 32rpx 0;
  position: relative;
  background: linear-gradient(135deg, #ff4c7d, #ff6b9d);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;

  &::before {
    content: '';
    position: absolute;
    left: 50%;
    bottom: 16rpx;
    transform: translateX(-50%);
    width: 60rpx;
    height: 4rpx;
    background: linear-gradient(90deg, #ff4c7d, #ff6b9d);
    border-radius: 2rpx;
  }
}
</style>
