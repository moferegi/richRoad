<template>
  <view>
    <view class="cate-box">
      <scroll-view scroll-y="true" class="left">
        <view class="left-text" :class="activeindex == index ? 'active' : ''" @tap="checkitem(index, item)"
              v-for="(item, index) in catelist" :key="index">{{ item.title }}</view>
      </scroll-view>
        <view class="right-box">
          <view class="cate-title" v-if="currentCategory">
            <view class="">{{ currentCategory.title }}</view>
          </view>
          <view class="goods-box" @tap="clickitem(item)" v-for="(item, index) in goodsList" :key="item.ID">
            <view class="img-box">
              <image class="goods-img" :src="getUrl(item.imageUrl)" mode="aspectFit"></image>
            </view>
            <view class="goods-info">
              <view class="goods-title">{{ item.title }}</view>
              <view class="goods-desc">{{ item.description }}</view>
              <view class="goods-tags" v-if="item.tags && item.tags.length > 0">
                <text class="tag" v-for="tag in item.tags" :key="tag.ID" :style="{color: tag.color}">
                  {{ tag.name }}
                </text>
              </view>
              <view class="goods-price-cart">
                <view class="goods-price">￥{{ item.price/10 }}</view>
                <view class="goods-cart">
                  <uni-icons color="#fff" size="20" type="cart"></uni-icons>
                </view>
              </view>
            </view>
          </view>

          <!-- 加载更多提示 -->
          <view class="load-more" v-if="loading">
            <text class="loading-text">正在加载...</text>
          </view>
          <view class="load-more" v-else-if="!hasMore && goodsList.length > 0">
            <text class="no-more-text">没有更多数据了</text>
          </view>
        </view>
    </view>
  </view>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { getUrl } from '@/utils/url'
import { getCategoryMobile, getGoodList } from '@/api/homePage.js'
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
const scrollTop = ref(0)

// 获取分类数据
const getCategoryData = async () => {
  try {
    const res = await getCategoryMobile()
    if (res && res.data) {
      catelist.value = res.data
      // 默认选择第一个分类
      if (catelist.value.length > 0) {
        currentCategory.value = catelist.value[0]
        await getGoodsData(catelist.value[0].ID, true)
      }
    }
  } catch (error) {
    console.error('获取分类数据失败:', error)
    uni.showToast({
      title: '获取分类失败',
      icon: 'none'
    })
  }
}

// 获取商品数据
const getGoodsData = async (categoryID) => {
  if (loading.value) return

  loading.value = true

  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      categoryID: categoryID
    }

    const res = await getGoodList(params)
    console.log(res);
    if(res.code === 0) {
      goodsList.value = res.data.list || []
    }

  } catch (error) {
    console.error('获取商品数据失败:', error)
    uni.showToast({
      title: '获取商品失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

const getgoods = (item) => {
    uni.showModal({
        title: '商品数据',
        content: '点击商品数据 = ' + JSON.stringify(item)
    })
}


// 切换分类
const checkitem = async (index, item) => {
  if (activeindex.value === index) return

  activeindex.value = index
  currentCategory.value = item

  // 重置分页数据
  currentPage.value = 1
  hasMore.value = true

  // 获取该分类下的商品
  await getGoodsData(item.ID, true)
}

const clickitem = (item) => {
    getgoods(item)
}

// 页面加载
onMounted(() => {
  getCategoryData()
})
</script>
<style lang="scss">
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

    .right {
        flex: 1;
        padding: 20rpx 15rpx;

        .right-box {
            font-size: 28rpx;
            padding-right: 20rpx;

            .cate-title {
                display: flex;
                justify-content: center;
                font-size: 30rpx;
                font-weight: 600;
                text-align: center;
                position: relative;
                margin-bottom: 20rpx;

                view {
                    padding: 0 20rpx;
                    border-left: 3px solid #e11d48;
                    border-right: 3px solid #e11d48;
                    font-weight: bold;
                }
            }

            // .cate-title::after{content: '';width: 80rpx;height: 5rpx;
            // 	background: #e11d48;position: absolute;bottom:-10rpx;

            // }


            .goods-box {
                display: flex;
                padding: 20rpx 0;
                align-items: center;

                .img-box {
                    width: 200rpx;
                    height: 200rpx;
                    display: flex;
                    align-items: center;

                    .goods-img {
                        width: 80%;
                        height: 80%;
                        border-radius: 10rpx;
                    }
                }

                .goods-info {
                    flex: 1;
                    display: flex;
                    flex-direction: column;
                    justify-content: space-between;
                    padding-bottom: 15rpx;
                    border-bottom: 1px solid #f5f4f8;

                    .goods-title {
                        margin-top: -10rpx;
                        font-weight: 500;
                        line-height: 40rpx;
                    }

                    .goods-desc {
                        color: #ef7e05;
                        margin-top: 10rpx;
                    }

                    .goods-price-cart {
                        display: flex;
                        justify-content: space-between;
                        align-items: center;
                        padding-top: 20rpx;

                        .goods-price {
                            color: #e11d48;
                            font-weight: bold;
                        }

                        .goods-cart {
                            height: 50rpx;
                            width: 50rpx;
                            background: #e11d48;
                            display: flex;
                            align-items: center;
                            justify-content: center;
                            border-radius: 50rpx;
                            ;
                        }
                    }
                }
            }
        }
    }
}
</style>
