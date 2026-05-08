<template>
  <view class="goods-list" ref="goodsListRef">
    <view class="goods-grid">
      <view v-for="(item, index) in props.goodsList" :key="index" class="goods-item" @tap="handleGoodsClick(item)">
        <LazyImage :src="item.externalImagePath ? getExternalUrl(item.externalImagePath) : getUrl(item.imageUrl)" class="goods-image" mode="aspectFill"></LazyImage>
        <view class="goods-info">
          <text class="goods-name">{{ $lt(item.title) }}</text>
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
              {{ $lt(tag.nameI18n || tag.name) }}
            </text>
          </view>
          <view class="price-container">
            <text class="discount-price">{{ cs }}{{ (item.price / 100).toFixed(2) }}</text>
            <text class="original-price">{{ cs }}{{ item.originalPrice }}</text>
          </view>
          <!-- 优惠券/积分标记 -->
          <view class="goods-badges" v-if="item.pointsEnabled || hasCoupon(item)">
            <text class="goods-badge goods-badge-coupon" v-if="hasCoupon(item)">{{ $t('canUseCoupon') }}</text>
            <text class="goods-badge goods-badge-points" v-if="item.pointsEnabled">{{ $t('canUsePoints') }}</text>
          </view>
          <view class="goods-extra">
            <view class="rating">
              <text class="rating-stars">{{ getRatingStars(item.rating) }}</text>
            </view>
            <view class="sales">
              <text>{{ $t('sold') }} {{ item.saleNum }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- 滚动检测触发区域 -->
    <view
      class="scroll-trigger"
      ref="scrollTriggerRef"
      v-if="!isLastPage && goodsList.length > 0"
    >
    </view>
  </view>
</template>

<script setup>
import {ref, computed, onMounted, onUnmounted, nextTick, watch} from 'vue'
import {getUrl, getExternalUrl} from "@/utils/url.js"
import { onReachBottom } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import LazyImage from '@/components/lazy-image/lazy-image.vue'

const langStore = useLangStore()
const appConfigStore = useAppConfigStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const $lt = computed(() => langStore.$lt)

const goodsList = ref([])
const props = defineProps({
  goodsList: {
    type: Array,
    default: () => []
  },
  currentPage: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    default: 0
  },
  // 距离底部多少rpx时触发加载
  loadOffset: {
    type: Number,
    default: 400
  }
})
const getDiscountText = (discount) => {
  if (discount >= 9.5) return $t.value('discountSmall')
  if (discount >= 9.0) return $t.value('discountNormal')
  if (discount >= 8.0) return $t.value('discountGood')
  if (discount >= 7.0) return $t.value('discountGreat')
  if (discount >= 6.0) return $t.value('discountLow')
  if (discount >= 5.0) return $t.value('discountSpecial')
  return $t.value('discountDefault')
}
const queryList = async (pageNo, pageSize) => {

// 加载数据方法

}
// 向父组件发送事件
const emit = defineEmits(['load-more'])

// 组件引用
const goodsListRef = ref(null)

// 计算是否为最后一页
const isLastPage = computed(() => {
  if (props.total % props.pageSize === 0) {
    return props.currentPage >= props.total / props.pageSize
  }
  return props.currentPage >= Math.floor(props.total / props.pageSize) + 1
})

// 页面滚动到底部触发加载
onReachBottom(() => {
  if (!isLastPage.value) {
    emit('load-more')
  }
})

/*const emit = defineEmits(['load-more'])
// 组件引用
const goodsListRef = ref(null)
const scrollTriggerRef = ref(null)

// 滚动检测相关
const isLoadingMore = ref(false)
const intersectionObserver = ref(null)

// 计算是否为最后一页
const isLastPage = computed(() => {
  if (props.total === 0) return false
  return props.currentPage * props.pageSize >= props.total
})*/
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

// 判断商品是否有可用优惠券
const hasCoupon = (item) => {
  return item.couponAvailable || (item.coupons && item.coupons.length > 0)
}

// 跳转到商品详情页
const handleGoodsClick = (item) => {
  uni.navigateTo({
    url: '/pages/goodsDetails/goodsDetails?id=' + item.ID
  })
}

/*// 初始化滚动检测
const initScrollDetection = () => {
  // 使用 IntersectionObserver 检测触发区域是否进入可视区域
  intersectionObserver.value = uni.createIntersectionObserver()

  intersectionObserver.value
    .relativeToViewport({ bottom: props.loadOffset })
    .observe('.scroll-trigger', (res) => {
      if (res.intersectionRatio > 0 && !isLoadingMore.value && !props.loading && !isLastPage.value) {
        handleAutoLoadMore()
      }
    })
}

// 自动加载更多
const handleAutoLoadMore = async () => {
  if (isLoadingMore.value || props.loading || isLastPage.value) {
    return
  }

  isLoadingMore.value = true

  try {
    // 发送加载更多事件给父组件
    emit('load-more')
  } finally {
    // 延迟重置状态，避免重复触发
    setTimeout(() => {
      isLoadingMore.value = false
    }, 500)
  }
}

// 监听商品列表变化，重新初始化滚动检测
watch(() => props.goodsList.length, async (newLength, oldLength) => {
  if (newLength > oldLength) {
    // 新数据加载完成后，重新初始化滚动检测
    await nextTick()
    if (intersectionObserver.value) {
      intersectionObserver.value.disconnect()
      initScrollDetection()
    }
  }
})

// 组件挂载时初始化
onMounted(() => {
  nextTick(() => {
    initScrollDetection()
  })
})

// 组件卸载时清理
onUnmounted(() => {
  if (intersectionObserver.value) {
    intersectionObserver.value.disconnect()
  }
})*/

</script>

<style lang="scss" scoped>
.goods-list {
  padding: 12rpx 20rpx;
  background: #000;
}

.goods-grid {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
}

.goods-item {
  width: 49%;
  box-sizing: border-box;
  margin-bottom: 16rpx;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.06);
  border-radius: 20rpx;
  overflow: hidden;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: transform 0.3s;

  &:active {
    transform: scale(0.97);
  }

  .goods-image {
    width: 100%;
    height: 340rpx;
  }
}

.goods-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16rpx 18rpx;
  min-width: 0;

  .goods-name {
    font-size: 26rpx;
    color: #fff;
    line-height: 1.35;
    margin-bottom: 8rpx;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 2;
    font-weight: 600;
    letter-spacing: 0.5rpx;
  }
}

.merchant-tags {
  display: flex;
  flex-wrap: wrap;
  margin: 4rpx 0;

  .merchant-tag {
    font-size: 18rpx;
    padding: 2rpx 10rpx;
    border-radius: 16rpx;
    height: 28rpx;
    line-height: 22rpx;
    margin-right: 6rpx;
    margin-bottom: 6rpx;
    font-weight: 500;
  }
}

.price-container {
  display: flex;
  align-items: baseline;
  margin: 8rpx 0 4rpx;
  flex-wrap: wrap;
  gap: 6rpx;

  .discount-price {
    font-size: 32rpx;
    color: #e50914;
    font-weight: 800;
  }

  .original-price {
    font-size: 22rpx;
    color: rgba(255, 255, 255, 0.3);
    text-decoration: line-through;
  }
}

.goods-extra {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.45);
  margin-top: 6rpx;

  .rating {
    display: flex;
    align-items: center;

    .rating-stars {
      color: #ffd700;
      font-size: 18rpx;
    }
  }

  .sales {
    color: rgba(255, 255, 255, 0.35);
    font-size: 20rpx;
  }
}

.goods-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 8rpx;
  margin: 6rpx 0 2rpx;

  .goods-badge {
    font-size: 18rpx;
    padding: 2rpx 10rpx;
    border-radius: 16rpx;
    height: 28rpx;
    line-height: 24rpx;
    font-weight: 500;
  }

  .goods-badge-coupon {
    color: #e50914;
    background: rgba(229, 9, 20, 0.12);
    border: 1rpx solid rgba(229, 9, 20, 0.3);
  }

  .goods-badge-points {
    color: #ffd700;
    background: rgba(255, 215, 0, 0.12);
    border: 1rpx solid rgba(255, 215, 0, 0.3);
  }
}

.scroll-trigger {
  height: 1rpx;
  width: 100%;
}
</style>
