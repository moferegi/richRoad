<template>
  <view>
    <view class="cate-box">
      <scroll-view class="left" scroll-y="true">
        <view v-for="(item, index) in catelist" :key="index" :class="activeindex == index ? 'active' : ''"
              class="left-text" @tap="checkitem(index, item)">{{ item.title }}
        </view>
      </scroll-view>
      <scroll-view :lower-threshold="lowerThresholdInPx" class="right-box" scroll-y="true"
                   @scrolltolower="handleScrollToLower">
        <!--        <view class="cate-title" v-if="currentCategory">
                  <view class="">GVA商城</view>
                </view>-->
        <view v-for="(item, index) in goodsList" :key="item.ID" class="goods-box" @tap="goto(item)">
          <image :src="getUrl(item.imageUrl)" class="goods-img" mode="aspectFit"></image>
          <view class="goods-info">
            <view class="goods-title">{{ item.title }}</view>
            <view class="goods-desc">{{ item.description }}</view>
            <view v-if="item.tags && item.tags.length > 0" class="merchant-tags">
              <text
                  v-for="tag in item.tags"
                  :key="tag.ID"
                  :style="{
                color: tag.color,
                background: `${tag.color}1A`,
                border: `1px solid ${tag.color}33`
            }"
                  class="merchant-tag"
              >
                {{ tag.name }}
              </text>
            </view>
            <view class="goods-price-cart">
              <view class="goods-price">￥{{ item.price / 100 }}</view>
              <view class="goods-cart" >
                <uni-icons color="#fff" size="20" type="search"></uni-icons>
              </view>
            </view>
          </view>
        </view>

        <view v-if="loading && goodsList.length === 0 && currentCategory" class="load-more">
          <text class="loading-text">正在加载...</text>
        </view>
        <view v-if="loading && goodsList.length > 0" class="load-more">
          <text class="loading-text">正在加载更多...</text>
        </view>
        <view v-if="!hasMore && goodsList.length > 0" class="load-more">
          <text class="no-more-text">没有更多数据了</text>
        </view>
        <view v-if="!hasMore && goodsList.length === 0 && !loading && currentCategory" class="load-more">
          <text class="no-more-text">暂无商品数据</text>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup>
import {ref, onMounted} from 'vue'
import {getUrl} from '@/utils/url'
import {getCategoryMobile, getGoodList} from '@/api/homePage.js'

// 分类数据
const catelist = ref([])
const activeindex = ref(0)
const currentCategory = ref(null)

// 商品数据
const goodsList = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const hasMore = ref(true)
const loading = ref(false)

// 新增：用于存储转换后的 lower-threshold 的 px 值
const lowerThresholdInPx = ref(50) // 给一个默认的px值，在onMounted中会被实际rpx转换值覆盖

// 获取分类数据
const getCategoryData = async () => {
  try {
    const res = await getCategoryMobile()
    if (res && res.data) {
      catelist.value = res.data
      if (catelist.value.length > 0) {
        currentCategory.value = catelist.value[0]
        await loadInitialGoods(currentCategory.value.ID)
      } else {
        goodsList.value = [];
        hasMore.value = false;
        currentCategory.value = null;
      }
    } else {
      goodsList.value = [];
      hasMore.value = false;
      currentCategory.value = null;
      uni.showToast({title: '获取分类数据失败', icon: 'none'});
    }
  } catch (error) {
    console.error('获取分类数据失败:', error)
    goodsList.value = [];
    hasMore.value = false;
    currentCategory.value = null;
    uni.showToast({
      title: '获取分类失败',
      icon: 'none'
    })
  }
}

// 加载初始商品数据
const loadInitialGoods = async (categoryID) => {
  goodsList.value = []
  currentPage.value = 1
  hasMore.value = true
  await fetchMoreGoods(categoryID)
}

// 获取更多商品数据
const fetchMoreGoods = async (categoryID) => {
  if (loading.value || !hasMore.value) {
    return
  }
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      categoryID: categoryID
    }
    const res = await getGoodList(params)
    if (res && res.code === 0 && res.data) {
      const newItems = res.data.list || []
      if (newItems.length > 0) {
        goodsList.value.push(...newItems)
        currentPage.value++
        if (newItems.length < pageSize.value) {
          hasMore.value = false;
        }
      } else {
        hasMore.value = false
      }
    } else {
      hasMore.value = false
      console.error('获取商品数据API响应异常:', res)
    }
  } catch (error) {
    console.error('获取商品数据失败:', error)
    hasMore.value = false
    uni.showToast({
      title: '获取商品列表失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

// 切换分类
const checkitem = async (index, item) => {
  if (activeindex.value === index && goodsList.value.length > 0) {
    return;
  }
  activeindex.value = index
  currentCategory.value = item
  await loadInitialGoods(item.ID)
}

// 滚动到底部时触发加载更多
const handleScrollToLower = () => {
  if (currentCategory.value && currentCategory.value.ID) {
    fetchMoreGoods(currentCategory.value.ID)
  }
}

// 点击商品
const clickitem = (item) => {
  getgoods(item)
}

const getgoods = (item) => {
  uni.showModal({
    title: '商品数据',
    content: '点击商品数据 = ' + JSON.stringify(item)
  })
}

const goto = (item) => {
  uni.navigateTo({
    url: '/pages/goodsDetails/goodsDetails?id=' + item.ID
  })
}
// 页面加载时执行
onMounted(() => {
  // 将 200rpx 转换为 px 值并设置
  // uni.upx2px 在编译到不同平台时，会基于基准宽度（通常是750rpx）进行转换
  lowerThresholdInPx.value = uni.upx2px(200);
  // console.log(`200rpx is approximately ${lowerThresholdInPx.value}px on this device.`); // 用于调试

  getCategoryData()
})
</script>

<style lang="scss">
// SCSS 样式与上一版完全相同，此处省略以减少重复
.cate-box {
  display: flex;
  height: 100vh;
  /* #ifdef H5 */
  height: calc(100vh - var(--window-top) - var(--window-bottom));
  /* #endif */

  .left {
    width: 200rpx;
    background: #f5f4f8;
    font-size: 28rpx;
    height: 100%;
    box-sizing: border-box;

    .left-text {
      text-align: center;
      padding: 20rpx 10rpx;
      border-bottom: 1px solid #fff;
      position: relative;
    }

    .active {
      padding: 24rpx 10rpx;
      background: #fff;
      font-size: 30rpx;
      font-weight: 600;
    }

    .active.left-text::before {
      content: '';
      height: 50rpx;
      width: 6rpx;
      background: #e11d48;
      position: absolute;
      left: 0;
      border-radius: 0 0rpx 0rpx 0;
    }
  }

  .merchant-tags {
    display: flex;
    flex-wrap: wrap;
    margin: 4rpx 0;

    .merchant-tag {
      font-size: 18rpx;
      padding: 0 6rpx;
      border-radius: 4rpx;
      height: 26rpx;
      line-height: 26rpx;
      margin-right: 4rpx;
      margin-bottom: 4rpx;

      &.self-operated {
        color: #ff6b6b;
        background: rgba(255, 107, 107, 0.1);
        border: 1px solid rgba(255, 107, 107, 0.2);
      }

      &.quality-assured {
        color: #2196f3;
        background: rgba(33, 150, 243, 0.1);
        border: 1px solid rgba(33, 150, 243, 0.2);
      }

      &.plus-delivery {
        color: #4caf50;
        background: rgba(76, 175, 80, 0.1);
        border: 1px solid rgba(76, 175, 80, 0.2);
      }
    }
  }

  .right-box {
    flex: 1;
    height: 100%;
    box-sizing: border-box;

    .cate-title {
      display: flex;
      justify-content: center;
      font-size: 30rpx;
      font-weight: 600;
      text-align: center;
      position: sticky;
      top: 0;
      z-index: 10;
      padding: 20rpx 15rpx 10rpx 15rpx;
      margin-bottom: 10rpx;

      view {
        padding: 0 20rpx;
        border-left: 3px solid #e11d48;
        border-right: 3px solid #e11d48;
        font-weight: bold;
      }
    }

    .goods-box {
      display: flex;
      padding: 40rpx 15rpx;
      align-items: center;


      .goods-img {
        width: 220rpx;
        height: 220rpx;
        border-radius: 8rpx;
        margin-right: 16rpx;
      }

      .goods-info {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        padding-bottom: 15rpx;
        border-bottom: 1px solid #f5f4f8;
        min-height: 180rpx;
        box-sizing: border-box;


        .goods-title {
          margin-top: 0;
          font-weight: 500;
          line-height: 1.4;
          font-size: 28rpx;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          -webkit-box-orient: vertical;
          margin-bottom: 8rpx;
        }

        .goods-desc {
          color: #999;
          margin-top: 0;
          font-size: 24rpx;
          line-height: 1.3;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          margin-bottom: 8rpx;
        }

        .goods-tags {
          margin-bottom: 8rpx;

          .tag {
            margin-right: 8rpx;
            padding: 2rpx 8rpx;
            font-size: 20rpx;
            border-radius: 4rpx;
            border: 1px solid;
            display: inline-block;
            line-height: 1.2;
          }
        }


        .goods-price-cart {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding-top: 10rpx;

          .goods-price {
            color: #e11d48;
            font-weight: bold;
            font-size: 32rpx;
          }

          .goods-cart {
            height: 50rpx;
            width: 50rpx;
            background: #e11d48;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 50%;
            margin-right: 20rpx;
          }
        }
      }
    }

    .load-more {
      text-align: center;
      padding: 20rpx 0;
      color: #999;
      font-size: 24rpx;
    }
  }
}
</style>
