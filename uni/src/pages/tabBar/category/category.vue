<template>
    <view>
        <view class="cate-box">
            <scroll-view scroll-y="true" class="left">
                <view class="left-text" :class="activeindex == index ? 'active' : ''" @tap="checkitem(index)"
                    v-for="(item, index) in catelist" :key="index">{{ item.title }}</view>
            </scroll-view>
            <scroll-view scroll-y="true" class="right" scroll-with-animation>
                <template v-for="(ite, i) in catelist" :key="i">
                    <view v-if="i == activeindex">
                        <view class="right-box" v-for="(item, index) in ite.cate1" :key="index">
                            <view class="cate-title">
                                <view class="">{{ item.title }}</view>
                            </view>
                            <view class="goods-box" @tap="clickitem(item1)" v-for="(item1, index1) in item.goods"
                                :key="index1">
                                <view class="img-box ">
                                    <image class="goods-img" :src="item1.img" lazy-load mode=""></image>
                                </view>
                                <view class="goods-info ">
                                    <view class="goods-title">{{ item1.title }}</view>
                                    <view class="goods-desc">{{ item1.desc }}</view>
                                    <view class="goods-price-cart">
                                        <view class="goods-price">￥{{ item1.price }}</view>
                                        <view class="goods-cart">
                                            <uni-icons color="#fff" size="20" type="cart"></uni-icons>
                                        </view>
                                    </view>
                                </view>
                            </view>
                        </view>
                    </view>
                </template>
            </scroll-view>
        </view>
    </view>
</template>
<script setup>


import { ref } from 'vue'
const catelist = ref([
    {
        "title": "女装",
        "cate1": [
            {
                "title": "连衣裙",
                "goods": [
                    {
                        "id": 1,
                        "title": "职业修身连衣裙",
                        "img": "http://images.shopin.net/images/1002/2022/01/10/3922eec2df75be5b901daadb624ca398_0.jpg",
                        "price": "228.00",
                        "desc": "非常好",
                    },
                    {
                        "id": 2,
                        "title": "碎花沙滩裙",
                        "img": "http://images.shopin.net/images/1002/2022/03/10/2433aecb2591d3d3398b4280fb90e28d_0.jpg",
                        "price": "322.00",
                        "desc": "面料舒适柔软",
                    },
                    {
                        "id": 3,
                        "title": "网格职业套裙",
                        "img": "http://images.shopin.net/images/1002/2021/09/15/2be4705eb9ba81faceb0e93d625b543c_0.jpg",
                        "price": "610.00",
                        "desc": "职业白领优选",
                    }

                ],
            },
            {
                "title": "女裤",
                "goods": [
                    {
                        "id": 4,
                        "title": "微喇叭裤",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN018bquaY1H93lgZqBp0_!!2200620780714-0-cib.jpg",
                        "price": "10.00",
                        "desc": "显瘦百搭",
                    },
                    {
                        "id": 5,
                        "title": "抽绳阔腿裤",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01vgImEI2Bn0mL59nLl_!!2200704948382-0-cib.jpg",
                        "price": "360.00",
                        "desc": "亚麻宽松休闲",
                    },
                    {
                        "id": 6,
                        "title": "高腰牛仔女裤",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01cbD8AS2Bae2A9hHpm_!!2215383308355-0-cib.jpg",
                        "price": "98.00",
                        "desc": "时尚性感木耳边",
                    }
                ],
            }
        ],
    },
    {
        "title": "男装",
        "cate1": [
            {
                "title": "上衣",
                "goods": [
                    {
                        "id": 7,
                        "title": "速干冰丝短袖T恤",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN017rTpYL1axsOFfw51u_!!2217500903397-0-cib.jpg",
                        "price": "22.00",
                        "sell_num": 140,
                        "kucun": "488",
                        "desc": "非常好",
                    },
                    {
                        "id": 8,
                        "title": "男式风衣",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01sjDZlQ1BtQ5sHPe7M_!!3076960003-0-cib.jpg",
                        "price": "22.00",
                        "sell_num": 0,
                        "kucun": "900",
                        "desc": "轻熟简约",
                    },
                    {
                        "id": 9,
                        "title": "长袖条纹t裇",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01C3PS4524CdvfAZRyi_!!2298067355-0-cib.jpg",
                        "price": "10.00",
                        "sell_num": 0,
                        "kucun": "600",
                        "desc": "翻领休闲长袖",
                    }
                ],
            },
        ],
    },
    {
        "title": "鞋品",
        "cate1": [
            {
                "title": "男鞋",
                "goods": [
                    {
                        "id": 10,
                        "title": "底减震运动鞋",
                        "img": "https://img.alicdn.com/imgextra/i4/1084736801/O1CN011VBW8K206ubkuAEmN_!!1084736801.jpg",
                        "price": "122.00",
                        "desc": "轻便跑步透气休闲机能软底",
                    },
                    {
                        "id": 11,
                        "title": "时尚加绒男鞋",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN011ITyHb1fqcESmSe8o_!!2212485574058-0-cib.jpg",
                        "price": "228.00",
                        "desc": "百搭休闲保暖",
                    },
                    {
                        "id": 12,
                        "title": "男款厚底小白板",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN016nW9rG1m3KixE0Qm5_!!2207977544898-0-cib.jpg",
                        "price": "190.00",
                        "desc": "休闲百搭运动",
                    }
                ],
            },
            {
                "title": "女鞋",
                "goods": [
                    {
                        "id": 13,
                        "title": "时尚防滑跑鞋",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01ySBj2h2DRbVYLI6r4_!!2218147088606-0-cib.jpg",
                        "price": "22.00",
                        "desc": "网面透气休闲",
                    },
                    {
                        "id": 14,
                        "title": "软底运动鞋",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01dZqwmj1fbxVq9LRrx_!!2211408924026-0-cib.jpg",
                        "price": "22.00",
                        "desc": "休闲跑步鞋",
                    },
                    {
                        "id": 15,
                        "title": "真皮外穿凉鞋",
                        "img": "https://cbu01.alicdn.com/img/ibank/O1CN01URiqcf1gmmTTnlocQ_!!2408174185-0-cib.jpg",
                        "price": "10.00",
                        "desc": "时尚百搭",
                    }
                ],
            },

        ],
    }
])

const getgoods = (item) => {
    uni.showModal({
        title: '商品数据',
        content: '点击商品数据 = ' + JSON.stringify(item)
    })
}

const activeindex = ref(0)

const checkitem = (index) => {
    activeindex.value = index
}

const clickitem = (item) => {
    getgoods(item)
}

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