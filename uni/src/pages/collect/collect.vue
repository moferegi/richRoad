<template>
	<view v-if="collectList.length">

    <view class="goods-list">
      <view class="goods-item" v-for="(item, index) in collectList" :key="index" @tap="goTo(item)">
        <image class="goods-image" :src="getUrl(item.imageUrl)" mode="aspectFill" />
        <view class="goods-info">
          <text class="goods-name">{{ item.title }}</text>
          <view class="merchant-tags">
            <text class="merchant-tag self-operated">自营</text>
            <text class="merchant-tag quality-assured">放心购</text>
            <text class="merchant-tag plus-delivery">Plus免邮</text>
          </view>
          <view class="price-container">
            <text class="discount-price">¥{{ item.price/100 }}</text>
            <text class="discount-tag">{{ getDiscountText(item.discount) }}</text>
          </view>
          <view class="goods-extra">
            <view class="rating">
              <text class="rating-score">{{ item.rating }}</text>
              <text class="rating-stars">★★★★★</text>
              <text class="rating-count">({{ item.ratingCount || 0 }})</text>
            </view>
            <view class="sales">
              <text>已售 {{ item.saleNum }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
	</view>
  <view v-if="!collectList.length" class="empty-state">
      <view class="empty-image-container">
        <image class="empty-image-placeholder" src="./../../static/emptyStatus.jpg"></image>
      </view>
      <view class="empty-text">暂无订单数据</view>
  </view>
</template>

<script setup>
	import { ref } from 'vue'
  import { onShow } from '@dcloudio/uni-app'
  import { getCollectList } from '@/api/collect'
  import {useUserStore} from "@/pinia/modules/user";
  import {findCollect, createCollect} from '@/api/collect.js'
  	import {getUrl} from "@/utils/url.js"
	const options = ref([{
		text: '删除',
		style: {
			backgroundColor: '#F56C6C'
		}
	}])

  const userStore = useUserStore()
  const token = userStore.token || ''
  const collectList = ref([])
  const collectionFlag = ref('')
  const init = async () => {
    const pageInfo = {
      page: 1,
      pageSize: 10
    }
    const res = await getCollectList(pageInfo)
    if(res.code === 0) {
      collectList.value = res.data.list
    }
  }

  onShow(() => {
    init()
  })
  // 防抖函数
  const debounce = (func, delay) => {
    let debounceTimer;
    return function(...args) {
      if (debounceTimer) clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        func.apply(this, args);
      }, delay);
    };
  };

  let params = {
    page: 1,
    pageSize: 10
  }
  const isBottom = ref(false)
  const lower = async (e) => {
    if(isBottom.value) {
      uni.showToast({
        title: '没有更多数据了',
        icon: 'none'
      })
      return
    } else {
      // 滑动到底了，然后每次给page+1 调接口继续加载下一页 如果接口已经没有数据了，给出提示并且不允许再次加载
      params.page += 1
      const res = await getCollectList(params)
      if (res.code === 0 && res.data.list.length) {
        collectList.value.push(...res.data.list)
        isBottom.value = false
      } else {
        isBottom.value = true
      }
    }
  }
  // 防抖包装的 lower 方法
  const debouncedLower = debounce(lower, 300);


  const cancelCollect = async (ID) => {
    if (token) {
      // 先查看当前商品收藏状态
      const status = await findCollect({
        goodID: ID
      })
      status.code === 0 ? collectionFlag.value = status.data : ''
      // 如果已登录并且未收藏 则允许进行收藏操作
      const res = await createCollect({
        goodID: Number(ID)
      })
      if (res.code === 0) {
        collectionFlag.value = !collectionFlag.value
        init()
        uni.showToast({
          title: collectionFlag.value ? '已收藏' : '已取消收藏',
          mask: true,
          icon: 'none'
        });
      }
    } else {
      uni.showToast({
        title: '请登录后进行操作',
        mask: true,
        icon: 'none'
      });
      uni.redirectTo({
        url: '/pages/user/login'
      })
    }
  }

  const goTo = (item) => {
    uni.navigateTo({
      url: `/pages/goodsDetails/goodsDetails?id=${item.ID}`,
    })
  }

  const getDiscountText = (discount) => {
    if (discount >= 9.5) return '小降'
    if (discount >= 9.0) return '优惠'
    if (discount >= 8.0) return '特惠'
    if (discount >= 7.0) return '好价'
    if (discount >= 6.0) return '低价'
    if (discount >= 5.0) return '特价'
    return '折扣'
  }
</script>

<style lang="scss" scoped>
	.collect_nav_view {
		width: 100%;
		height: 124rpx;
	}

  /* 空状态样式 */
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: calc(100vh - 88rpx);
    /* #ifdef H5 */
    height: calc(100vh - 88rpx - var(--window-top));
    /* #endif */
  }

  .empty-image-container {
    margin-bottom: 30rpx;
  }

  .empty-image-placeholder {
    width: 280rpx;
    height: 280rpx;
    border-radius: 8rpx;
  }

  .empty-text {
    font-size: 28rpx;
    color: #999;
  }

  .scroll-Y {
    height: 100vh;
  }
	.collect_nav_box {
		width: 100%;
		height: 124rpx;
		z-index: 1;
		/* #ifdef H5 */
		top: var(--window-top);
		/* #endif */
		/* #ifndef H5 */
		top: 0;
		/* #endif */
		left: 0;
		right: 0;
	}

	.collect_search_box {
		width: 100%;
		padding: 24rpx 32rpx;
	}

	.search_icon {
		top: 42rpx;
		left: 64rpx;
		width: 40rpx;
		height: 40rpx;
	}

	.search_input {
		width: 100%;
		display: block;
		height: 76rpx;
		border-radius: 76rpx;
		padding-left: 88rpx;
	}
  .collect_goods {
    padding: 24rpx 32rpx 0;
  }
  .collect_goods_img {
    width: 200rpx;
    height: 172rpx;
    border-radius: 12rpx;
    overflow: hidden;
  }
	.non-collect{
    height: 100vh;
  }

  .desc{
    width: 100%;
    margin-top: 20rpx;
  }

  page {
    background: #f5f7fa;
  }

  .goods-list {
    padding: 12rpx;
    background: #f5f7fa;
  }

  .goods-item {
    display: flex;
    background: #ffffff;
    margin-bottom: 12rpx;
    border-radius: 12rpx;
    padding: 12rpx;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);

    &:active {
      transform: scale(0.98);
    }

    .goods-image {
      width: 240rpx;
      height: 240rpx;
      border-radius: 8rpx;
      margin-right: 16rpx;
    }
  }

  .goods-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;

    .goods-name {
      font-size: 28rpx;
      color: #333;
      line-height: 1.4;
      margin-bottom: 8rpx;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 2;
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

  .price-container {
    display: flex;
    align-items: center;
    margin: 8rpx 0;

    .discount-price {
      font-size: 32rpx;
      color: #ff4444;
      font-weight: bold;
      margin-right: 8rpx;
    }

    .original-price {
      font-size: 22rpx;
      color: #999;
      text-decoration: line-through;
      margin-right: 8rpx;
    }

    .discount-tag {
      font-size: 24rpx;
      color: #fff;
      background: #ff4444;
      padding: 2rpx 8rpx;
      border-radius: 4rpx;
    }
  }

  .goods-extra {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 24rpx;
    color: #666;
    margin-top: 8px;
    .rating {
      display: flex;
      align-items: center;

      .rating-score {
        color: #ff4444;
        font-weight: bold;
        margin-right: 2rpx;
      }

      .rating-stars {
        color: #ffd700;
        font-size: 18rpx;
        margin-right: 2rpx;
      }

      .rating-count {
        color: #999;
      }
    }

    .sales {
      color: #999;
    }
  }

  .goods-tags {
    display: flex;
    flex-wrap: wrap;
    margin: 4rpx 0;

    .tag {
      font-size: 18rpx;
      color: #666;
      background: #f7f7f7;
      padding: 0 6rpx;
      border-radius: 2rpx;
      margin-right: 4rpx;
      margin-bottom: 4rpx;
    }
  }
</style>
