<template>
    <z-paging ref="paging" v-model="goodsList" @query="queryList" auto-show-system-loading>
        <view class="goods-list">
            <view class="goods-item" v-for="(item, index) in goodsList" :key="index" @click="handleGoodsClick(item)">
                <image class="goods-image" :src="item.image" mode="aspectFill" />
                <view class="goods-info">
                    <text class="goods-name">{{ item.name }}</text>
                    <view class="merchant-tags">
                        <text v-if="item.isSelfOperated" class="merchant-tag self-operated">{{ $t('selfOperated') }}</text>
                        <text v-if="item.hasQualityAssurance" class="merchant-tag quality-assured">{{ $t('qualityAssured') }}</text>
                        <text v-if="item.isPlusDelivery" class="merchant-tag plus-delivery">{{ $t('freeShipping') }}</text>
                    </view>
                    <view class="price-container">
                        <text class="discount-price">{{ cs }}{{ item.discountPrice }}</text>
                        <text class="original-price">{{ cs }}{{ item.originalPrice }}</text>
                        <text class="discount-tag">{{ getDiscountText(item.discount) }}</text>
                    </view>
                    <view class="shop-info">
                        <view class="shop-left">
                            <image class="shop-avatar" :src="item.shop.avatar" mode="aspectFill" />
                            <text class="shop-name">{{ item.shop.name }}</text>
                        </view>
                        <view class="enter-shop" @click.stop="enterShop(item.shop)">{{ $t('goShopping') }}</view>
                    </view>
                </view>
            </view>
        </view>
    </z-paging>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'
import { useLangStore } from '@/pinia/modules/lang.js'
const appConfigStore = useAppConfigStore()
const langStore = useLangStore()
const cs = computed(() => appConfigStore.currencySymbol)
const $t = computed(() => langStore.$t)
const goodsList = ref([])
const mockData = [
    {
        name: '2023 Breathable running sneakers for men and women',
        image: 'https://picsum.photos/300/300?random=1',
        originalPrice: 399,
        discountPrice: 299,
        discount: 7.5,
        discountEmoji: '🔥',
        rating: 4.8,
        ratingCount: 2531,
        monthSales: 1688,
        tags: ['Authentic', 'Fast Shipping', '7-Day Return'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: true,
        shop: {
            name: 'Sport & Outdoor Store',
            avatar: 'https://picsum.photos/64/64?random=1',
            rating: 4.8,
            isOfficial: true
        }
    },
    {
        name: 'Large waterproof canvas backpack for school and travel',
        image: 'https://picsum.photos/300/300?random=2',
        originalPrice: 199,
        discountPrice: 139,
        discount: 7.0,
        discountEmoji: '⚡',
        rating: 4.6,
        ratingCount: 1234,
        monthSales: 966,
        tags: ['Brand Choice', 'Free Shipping'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: false,
        shop: {
            name: 'Fashion Bags Flagship',
            avatar: 'https://picsum.photos/64/64?random=2',
            rating: 4.7,
            isOfficial: true
        }
    },
    {
        name: 'Smart watch with sports tracking and heart monitoring',
        image: 'https://picsum.photos/300/300?random=3',
        originalPrice: 899,
        discountPrice: 699,
        discount: 7.8,
        discountEmoji: '💥',
        rating: 4.7,
        ratingCount: 1876,
        monthSales: 1245,
        tags: ['Smart Watch', 'Waterproof'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: true,
        shop: {
            name: 'Smart Devices Flagship',
            avatar: 'https://picsum.photos/64/64?random=3',
            rating: 4.7,
            isOfficial: true
        }
    },
    {
        name: 'True wireless earbuds with active noise cancellation',
        image: 'https://picsum.photos/300/300?random=4',
        originalPrice: 299,
        discountPrice: 199,
        discount: 6.6,
        discountEmoji: '🎉',
        rating: 4.5,
        ratingCount: 1023,
        monthSales: 789,
        tags: ['Wireless Earbuds', 'Noise Canceling'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: false,
        shop: {
            name: 'Audio Devices Flagship',
            avatar: 'https://picsum.photos/64/64?random=4',
            rating: 4.5,
            isOfficial: true
        }
    },
    {
        name: 'Health smart band with heart rate and activity tracking',
        image: 'https://picsum.photos/300/300?random=5',
        originalPrice: 199,
        discountPrice: 149,
        discount: 7.5,
        discountEmoji: '🎯',
        rating: 4.3,
        ratingCount: 852,
        monthSales: 654,
        tags: ['Smart Band', 'Health Tracking'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: true,
        shop: {
            name: 'Health Monitoring Flagship',
            avatar: 'https://picsum.photos/64/64?random=5',
            rating: 4.3,
            isOfficial: true
        }
    },
    {
        name: 'Portable bluetooth speaker with waterproof design',
        image: 'https://picsum.photos/300/300?random=6',
        originalPrice: 299,
        discountPrice: 239,
        discount: 8.0,
        discountEmoji: '⚡',
        rating: 4.9,
        ratingCount: 3000,
        monthSales: 2000,
        tags: ['Bluetooth Speaker', 'Wireless'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: true,
        shop: {
            name: 'Smart Home Flagship',
            avatar: 'https://picsum.photos/64/64?random=6',
            rating: 4.9,
            isOfficial: true
        }
    },
    {
        name: 'Business laptop backpack with anti-theft design',
        image: 'https://picsum.photos/300/300?random=7',
        originalPrice: 259,
        discountPrice: 189,
        discount: 7.3,
        discountEmoji: '💫',
        rating: 4.2,
        ratingCount: 1500,
        monthSales: 1000,
        tags: ['Laptop Backpack', 'Waterproof'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: false,
        shop: {
            name: 'Computer Accessories Flagship',
            avatar: 'https://picsum.photos/64/64?random=7',
            rating: 4.2,
            isOfficial: true
        }
    },
    {
        name: 'Mechanical keyboard with RGB backlight and dual mode',
        image: 'https://picsum.photos/300/300?random=8',
        originalPrice: 499,
        discountPrice: 399,
        discount: 8.0,
        discountEmoji: '🌟',
        rating: 4.7,
        ratingCount: 2200,
        monthSales: 1500,
        tags: ['Mechanical Keyboard', 'Backlight'],
        isSelfOperated: true,
        hasQualityAssurance: true,
        isPlusDelivery: true,
        shop: {
            name: 'Electronics Accessories Flagship',
            avatar: 'https://picsum.photos/64/64?random=8',
            rating: 4.7,
            isOfficial: true
        }
    }
]

const paging = ref(null)
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
        // 模拟网络延迟1.5秒
        await new Promise(resolve => setTimeout(resolve, 1500))

        // 计算当前页数据的起始位置
        const start = (pageNo - 1) * pageSize
        const totalItems = mockData.length
        const pageData = []

        // 生成当前页数据
        for (let i = 0; i < pageSize; i++) {
            const index = (start + i) % totalItems
            if (pageNo <= 3) { // 只模拟3页数据
                pageData.push({
                    ...mockData[index],
                    // 添加页码标记，方便区分不同页的数据
                    name: `${mockData[index].name} - ${pageNo}-${i + 1}`,
                    // 确保不同页的图片URL不重复
                    image: `${mockData[index].image}&page=${pageNo}&item=${i}`
                })
            }
        }

        // 处理数据加载完成的情况
        if (pageNo >= 3 || pageData.length === 0) {
            paging.value.complete(pageData)
        } else {
            paging.value.complete(pageData)
        }

}

</script>

<style scoped lang="scss">
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
        font-size: 20rpx;
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
    font-size: 20rpx;
    color: #666;
    margin-top: 8rpx;

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