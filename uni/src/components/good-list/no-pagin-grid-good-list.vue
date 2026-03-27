<template>
  <view class="goods-list">
    <view class="goods-row"  v-if="props.goodsList.length">
      <view v-for="(item, index) in props.goodsList" :key="index" class="goods-item" @tap="handleGoodsClick(item)">
        <image :src="getUrl(item.imageUrl)" class="goods-image" mode="aspectFill"></image>
        <view class="goods-info">
          <text class="goods-name">{{ $lt(item.title) }}</text>
          <view class="price-container">
            <text class="discount-price">¥{{ item.price / 100 }}</text>
            <text class="original-price">¥{{ item.originalPrice }}</text>
            <text class="discount-tag">{{ getDiscountText(item.discount) }}</text>
          </view>
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
          <view class="goods-extra">
            <view class="rating">
              <text class="rating-score">{{ item.rating }}</text>
              <text class="rating-stars">{{ getRatingStars(item.rating) }}</text>
              <text class="rating-count">({{ item.ratingCount || 0 }})</text>
            </view>
            <view class="sales">
              <text>售出 {{ item.saleNum }} 件</text>
            </view>
          </view>
        </view>

      </view>
    </view>
    <view v-if="!props.goodsList.length" class="empty-state">
      <view class="empty-image-container">
        <image class="empty-image-placeholder" src="./../../static/emptyStatus.jpg"></image>
      </view>
      <view class="empty-text">暂无订单数据</view>
    </view>
  </view>
</template>

<script setup>
import {ref, computed, onMounted, onUnmounted} from 'vue'
import {getUrl} from "@/utils/url";
import { useLangStore } from '@/pinia/modules/lang.js'

const langStore = useLangStore()
const $lt = computed(() => langStore.$lt)

const props = defineProps({
  goodsList: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  noMore: {
    type: Boolean,
    default: false
  }
})
const emit = defineEmits(['loadMore'])
// 节流标识
const isThrottling = ref(false)

const getDiscountText = (discount) => {
  if (discount >= 9.5) return '小降'
  if (discount >= 9.0) return '优惠'
  if (discount >= 8.0) return '特惠'
  if (discount >= 7.0) return '好价'
  if (discount >= 6.0) return '低价'
  if (discount >= 5.0) return '特价'
  return '折扣'
}

// 根据评分生成星星
const getRatingStars = (rating) => {
  // 如果评分为0，返回1颗星
  if (rating === 0) return '★☆☆☆☆';

  // 计算实心星星数量（最大5颗）
  const fullStars = Math.min(Math.floor(rating), 5);
  // 计算空心星星数量
  const emptyStars = 5 - fullStars;

  // 返回对应数量的星星
  return '★'.repeat(fullStars) + '☆'.repeat(emptyStars);
}


// 跳转到商品详情页
const handleGoodsClick = (item) => {
  uni.navigateTo({
    url: '/pages/player/index?id=' + item.ID
  })
}

// 滚动事件处理
const handleScroll = () => {
  // 如果正在加载、没有更多数据、或者正在节流，则不处理
  if (props.loading || props.noMore || isThrottling.value) {
    return
  }

  // H5环境下直接使用window对象
  // #ifdef H5
  const scrollTop = document.documentElement.scrollTop || document.body.scrollTop
  const scrollHeight = document.documentElement.scrollHeight || document.body.scrollHeight
  const windowHeight = window.innerHeight

  console.log('H5滚动检测:', { scrollTop, scrollHeight, windowHeight })

  // 距离底部还有100px时开始加载
  if (scrollTop + windowHeight >= scrollHeight - 100) {
    console.log('触发加载更多')
    // 开启节流
    isThrottling.value = true

    // 触发加载更多
    emit('loadMore')

    // 500ms后关闭节流
    setTimeout(() => {
      isThrottling.value = false
    }, 500)
  }
  // #endif

  // 小程序等其他平台使用uni API
  // #ifndef H5
  uni.createSelectorQuery().selectViewport().scrollOffset((res) => {
    const scrollTop = res.scrollTop
    const scrollHeight = res.scrollHeight
    const windowHeight = uni.getSystemInfoSync().windowHeight

    console.log('小程序滚动检测:', { scrollTop, scrollHeight, windowHeight })

    // 距离底部还有100px时开始加载
    if (scrollTop + windowHeight >= scrollHeight - 100) {
      console.log('触发加载更多')
      // 开启节流
      isThrottling.value = true

      // 触发加载更多
      emit('loadMore')

      // 500ms后关闭节流
      setTimeout(() => {
        isThrottling.value = false
      }, 500)
    }
  }).exec()
  // #endif
}


// 监听页面滚动
// 监听页面滚动
onMounted(() => {
  // H5环境下监听window滚动
  // #ifdef H5
  window.addEventListener('scroll', handleScroll, { passive: true })
  console.log('H5滚动监听已添加')
  // #endif

  // 其他平台使用uni API
  // #ifndef H5
  handleScroll()
  console.log('uni滚动监听已添加')
  // #endif
})

// 清理监听
onUnmounted(() => {
  // #ifdef H5
  window.removeEventListener('scroll', handleScroll)
  console.log('H5滚动监听已移除')
  // #endif

  // #ifndef H5
  uni.offPageScroll(handleScroll)
  console.log('uni滚动监听已移除')
  // #endif
})

</script>

<style lang="scss" scoped>
/* 页面背景色 */


/* 商品列表容器 */
.goods-list {
  padding: 16rpx 12rpx;

  /* 商品行布局 */
  .goods-row {
    display: flex;
    flex-wrap: wrap;
    margin: 0 -8rpx;
  }
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


/* 商品卡片 */
.goods-item {
  width: calc(50% - 16rpx);
  background: #ffffff;
  border-radius: 12rpx;
  margin: 8rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;

  &:active {
    transform: scale(0.98);
  }

  .goods-image {
    width: 100%;
    height: 320rpx;
  }
}

/* 商品信息区域 */
.goods-info {
  padding: 12rpx 12rpx 8rpx;
  display: flex;
  flex-direction: column;
  min-height: 220rpx;
  position: relative;

  .goods-name {
    font-size: 28rpx;
    color: #333;
    margin-bottom: 8rpx;
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    word-break: break-all;
  }
}

/* 价格区域 */
.price-container {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  margin: 4rpx 0;

  .discount-price {
    font-size: 32rpx;
    color: #ff4444;
    font-weight: bold;
    margin-right: 6rpx;
  }

  .original-price {
    font-size: 22rpx;
    color: #999;
    text-decoration: line-through;
    margin-right: 6rpx;
  }

  .discount-tag {
    font-size: 20rpx;
    color: #fff;
    background: #ff4444;
    padding: 2rpx 8rpx;
    border-radius: 4rpx;
  }
}

/* 商品额外信息 */
.goods-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 20rpx;
  color: #666;
  margin: 6rpx 0;

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
      font-size: 28rpx;
    }
  }

  .sales {
    color: #999;
    font-size: 28rpx;
  }
}

/* 商品标签 */
.goods-tags {
  display: flex;
  flex-wrap: wrap;
  margin-top: 6rpx;

  .tag {
    font-size: 18rpx;
    color: #666;
    background: #f7f7f7;
    padding: 0 6rpx;
    border-radius: 2rpx;
    line-height: 24rpx;
    margin-right: 4rpx;
    margin-bottom: 4rpx;
  }
}

/* 商家标签 */
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
</style>
