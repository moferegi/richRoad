<template>
        <view class="goods-list">
            <view class="goods-item" v-for="(item, index) in goodsList" :key="index" @click="handleGoodsClick(item)">
                <image class="goods-image" :src="item.image" mode="aspectFill" />
                <view class="goods-info">
                    <text class="goods-name">{{ item.name }}</text>
                    <view class="merchant-tags">
                        <text v-if="item.isSelfOperated" class="merchant-tag self-operated">自营</text>
                        <text v-if="item.hasQualityAssurance" class="merchant-tag quality-assured">放心购</text>
                        <text v-if="item.isPlusDelivery" class="merchant-tag plus-delivery">Plus免邮</text>
                    </view>
                    <view class="price-container">
                        <text class="discount-price">¥{{ item.discountPrice }}</text>
                        <text class="original-price">¥{{ item.originalPrice }}</text>
                        <text class="discount-tag">{{ getDiscountText(item.discount) }}</text>
                    </view>
                    <view class="goods-extra">
                            <view class="rating">
                                <text class="rating-score">{{ item.rating }}</text>
                                <text class="rating-stars">★★★★★</text>
                                <text class="rating-count">({{ item.ratingCount }})</text>
                            </view>
                            <view class="sales">
                                <text>月销 {{ item.monthSales }}</text>
                            </view>
                        </view>
                </view>
            </view>
        </view>
</template>

<script>
export default {
    name: 'NoPaginrowGoodList',
    
    data() {
        return {
            goodsList:[
                {
                    name: '2023新款时尚运动鞋男女同款透气网面跑步鞋减震耐磨休闲运动鞋',
                    image: 'https://picsum.photos/300/300?random=1',
                    originalPrice: 399,
                    discountPrice: 299,
                    discount: 7.5,
                    rating: 4.8,
                    ratingCount: 2531,
                    monthSales: 1688,
                    tags: ['正品保证', '极速发货', '七天退换'],
                    isSelfOperated: true,
                    hasQualityAssurance: true,
                    isPlusDelivery: true,
                    shop: {
                        name: '运动户外专营店',
                        avatar: 'https://picsum.photos/64/64?random=1',
                        rating: 4.8,
                        isOfficial: true
                    }
                },
                {
                    name: '苹果14代 iPhone14 Pro Max 5G手机 暗紫色 256GB全网通',
                    image: 'https://picsum.photos/300/300?random=2',
                    originalPrice: 9999,
                    discountPrice: 9299,
                    discount: 9.3,
                    rating: 4.9,
                    ratingCount: 12420,
                    monthSales: 3562,
                    tags: ['官方授权', '全国联保', '24期免息'],
                    isSelfOperated: true,
                    hasQualityAssurance: true,
                    isPlusDelivery: true,
                    shop: {
                        name: 'Apple官方旗舰店',
                        avatar: 'https://picsum.photos/64/64?random=2',
                        rating: 4.9,
                        isOfficial: true
                    }
                },
                {
                    name: '科沃斯扫地机器人T10 Pro家用扫拖一体机智能全自动吸尘器',
                    image: 'https://picsum.photos/300/300?random=3',
                    originalPrice: 4999,
                    discountPrice: 3999,
                    discount: 8.0,
                    rating: 4.7,
                    ratingCount: 8526,
                    monthSales: 2156,
                    tags: ['新品上市', '智能家电', '以旧换新'],
                    isSelfOperated: false,
                    hasQualityAssurance: true,
                    isPlusDelivery: true,
                    shop: {
                        name: '科沃斯官方店',
                        avatar: 'https://picsum.photos/64/64?random=3',
                        rating: 4.8,
                        isOfficial: true
                    }
                },
                {
                    name: '华为智慧屏V65 2022款65英寸4K超高清智能电视机',
                    image: 'https://picsum.photos/300/300?random=4',
                    originalPrice: 5999,
                    discountPrice: 4999,
                    discount: 8.3,
                    rating: 4.8,
                    ratingCount: 3654,
                    monthSales: 986,
                    tags: ['超清画质', '智能语音', '大屏影音'],
                    isSelfOperated: true,
                    hasQualityAssurance: true,
                    isPlusDelivery: false,
                    shop: {
                        name: '华为官方旗舰店',
                        avatar: 'https://picsum.photos/64/64?random=4',
                        rating: 4.9,
                        isOfficial: true
                    }
                },
                {
                    name: '蒙牛特仑苏纯牛奶250ml*12盒整箱装',
                    image: 'https://picsum.photos/300/300?random=5',
                    originalPrice: 69.9,
                    discountPrice: 59.9,
                    discount: 8.5,
                    rating: 4.6,
                    ratingCount: 15689,
                    monthSales: 8562,
                    tags: ['新鲜直达', '冷链配送', '营养早餐'],
                    isSelfOperated: true,
                    hasQualityAssurance: true,
                    isPlusDelivery: true,
                    shop: {
                        name: '蒙牛乳业官方店',
                        avatar: 'https://picsum.photos/64/64?random=5',
                        rating: 4.7,
                        isOfficial: true
                    }
                },
                {
                    name: '荣耀手环7 NFC版 血氧心率监测智能运动手环',
                    image: 'https://picsum.photos/300/300?random=6',
                    originalPrice: 299,
                    discountPrice: 249,
                    discount: 8.3,
                    rating: 4.5,
                    ratingCount: 6523,
                    monthSales: 2365,
                    tags: ['续航持久', '健康监测', '运动计步'],
                    isSelfOperated: false,
                    hasQualityAssurance: true,
                    isPlusDelivery: true,
                    shop: {
                        name: '荣耀智能专卖店',
                        avatar: 'https://picsum.photos/64/64?random=6',
                        rating: 4.6,
                        isOfficial: false
                    }
                }
            ]
        }
    },

    methods: {
        getDiscountText(discount) {
            if (discount >= 9.5) return '小降'
            if (discount >= 9.0) return '优惠'
            if (discount >= 8.0) return '特惠'
            if (discount >= 7.0) return '好价'
            if (discount >= 6.0) return '低价'
            if (discount >= 5.0) return '特价'
            return '折扣'
        },

        enterShop(shop) {
            uni.showToast({
                title: '进入店铺：' + shop.name,
                icon: 'none'
            })
        },

        async queryList(pageNo, pageSize) {
            try {
                await new Promise(resolve => setTimeout(resolve, 1500))
                const start = (pageNo - 1) * pageSize
                const totalItems = this.mockData.length
                const pageData = []
                
                if (pageNo <= 3) {
                    for (let i = 0; i < pageSize && start + i < totalItems; i++) {
                        const index = start + i
                        pageData.push({
                            ...this.mockData[index % totalItems],
                            name: `${this.mockData[index % totalItems].name} - ${pageNo}-${i + 1}`,
                            image: `${this.mockData[index % totalItems].image}&page=${pageNo}&item=${i}`
                        })
                    }
                    this.$refs.paging.complete(pageData)
                } else {
                    this.$refs.paging.complete([])
                }
            } catch (error) {
                this.$refs.paging.complete([])
                uni.showToast({
                    title: '加载失败',
                    icon: 'none'
                })
            }
        },

        handleGoodsClick(item) {
            uni.showToast({
                title: '点击商品：' + item.name,
                icon: 'none'
            })
        }
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