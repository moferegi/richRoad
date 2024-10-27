<template>
  <view>
    <goods-swiper :list="data.banner"></goods-swiper>
    <view class="goods_title boxs_bb m_b_24 bgc_fff">
      <view class="m_b_8 color_333 font_32 font_bold">{{ data.title }}</view>
      <view class="m_b_16 color_999 font_28">{{ data.description }}</view>
      <view v-if="data.saleCount&&data.saleCount>0" class="m_b_16 color_999 font_28">{{ data.saleCount }}</view>
      <view class="flex p_b_24 flex-aife">
        <view class="price_color">
          <text class="font_36">¥</text>
          <text class="font_64 m_r_24">{{ data.price && data.price / 100 }}</text>
        </view>
      </view>
    </view>
    <view class="bgc_fff font_28 color_333 m_b_24">
      <view class="goods_item_box boxs_bb" v-for="(i, k) in data.attrs" :key="k">
          <text class="goods_item_cla color_999">{{ i.name }}</text>
        <text>{{ i.value }}</text>
      </view>
      <view class="goods_item_box boxs_bb">
        <text class="goods_item_cla color_999">运费</text>
        <text>免运费
          <text class="color_eee">｜</text>
          48小时内发货
        </text>
      </view>
    </view>

    <view class="bgc_fff goods_item_box">
      <view class="flex-aic flexr-jsb m_b_24">
        <text class="font_28 color_333">用户评价</text>
        <view class="flex flex-aic" @tap="toEvaluate">
          <text class="font_24 color_fe5572">查看全部</text>
          <uni-icons color="#ff7000" size="16" type="right"></uni-icons>
        </view>
      </view>
      <view v-if="hasContent">
        <view class="flex m_b_24 ">
          <image :src="getUrl(commentInfo.user.avatar)" class="item_head m_r_16"></image>
          <view class="flex-fitem">
            <view class="flex-aic flexr-jsb">
              <text class="color_333 font_28">{{ commentInfo.user.nickname }}</text>
              <text class="color_999 font_24">{{ formatTimeToStr(commentInfo.CreatedAt, 'yyyy-MM-dd') }}</text>
            </view>
            <uni-rate :readonly="true" :value="commentInfo.rating" active-color="#fe5572" size="16"/>
          </view>
        </view>
        <view class="boxs_bb" style="padding: 0 16rpx 0 80rpx;">
          <view class="text_pre_wrap m_b_24 color_333 font_28">{{ commentInfo.content }}</view>
          <evaluate-grid-img :imgList="commentInfo.feedbackPics&&commentInfo.feedbackPics"></evaluate-grid-img>

        </view>
      </view>
      <view v-if="!hasContent">
        暂无评价
      </view>
    </view>


    <view class="bb_title font_28 color_999 flex-aic flexr-jsc">
      <view class="bgc_ddd"></view>
      <text>宝贝详情</text>
    </view>
    <goodsDetail :detail="data.detail"></goodsDetail>
    <rich-text style="width: 100%;"/>
    <goods-sku v-if="data.skus" ref="goodsSkuRef" :isCart="isCart" :good="data" @toOrder="toOrder"></goods-sku>
    <view class="goods_nav_view">
      <view class="pos_f g_nav_box bgc_fff flex flex-aic">
        <view class="flex-fitem flex-aic flexr-jsc font_20 color_333" @tap="goTo">
          <image class="g_nav_icon m_r_8" mode="" src="./../../static/images/tabBar/index.png">
          </image>
          <!-- <text>首页</text> -->
        </view>
        <view class="flex-fitem flex-aic flexr-jsc font_20 color_333" @tap="addCollect">
          <image :src="!collectionFlag? './../../static/collection.png' : './../../static/collect.png'"
                 class="g_nav_icon m_r_8" mode="">
          </image>
          <!-- <text>收藏</text> -->
        </view>
        <view class="flex-fitem flex-aic flexr-jsc font_20 color_333" @tap="goTo('cart')">
          <image class="g_nav_icon m_r_8" mode="" src="./../../static/images/tabBar/cart.png">
          </image>
          <!-- <text>购物车</text> -->
        </view>
        <view class="btns_box color_fff font_32 flex flex-aic">
          <view class="bgc_ffa259 tac" @tap="addToCart()">加入购物车</view>
          <view class="bgc_ff7000 tac" @tap="goodsTapPay('pay')">立即购买</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import goodsSwiper from './components/goods-swiper.vue'
import goodsSku from './components/goods-sku.vue'
import goodsDetail from './components/goods-detail.vue';
import {ref} from "vue";
import {onLoad} from '@dcloudio/uni-app'
import {findGood} from '@/api/product.js'
import {myRouter} from '../../utils/permission';
import {findCollect, createCollect} from '@/api/collect.js'
import {useUserStore} from "@/pinia/modules/user";
import {findComment} from "@/api/comment.js"
import {formatTimeToStr} from "@/utils/date.js"
import evaluateGridImg from '/pages/evaluate/evaluate-img.vue'
import {getUrl} from "@/utils/url.js"
const data = ref({})
const collectionFlag = ref('')
const hasContent = ref(false)
const goodID = ref(0)
const userStore = useUserStore()
const token = userStore.token || ''
const commentInfo = ref([])
onLoad((options) => {
  // 先获取商品属性
  if (options.id) {
    goodID.value = options.id
    hasContent.value = false
    init()
  }
})

const init = async () => {
  const res = await findGood(goodID.value)
  res.code === 0 ? data.value = res.data.regood : ''
  // 如果登录，就获取当前商品收藏状态 反之则默认未收藏，点击跳转 登录后收藏
  if (token) {
    // 先查看当前商品收藏状态
    const status = await findCollect({
      goodID: goodID.value
    })
    status.code === 0 ? collectionFlag.value = status.data : ''
  }

  const res2 = await findComment(goodID.value)
  if (res2.code === 0 && res2.data.length) {
    commentInfo.value = res2.data[0]
    hasContent.value = true
  }
}

const toEvaluate = () => {
  uni.navigateTo({
    url: `/pages/evaluate/evaluate?goodsID=${goodID.value}`
  })
}

const goodsSkuRef = ref()

const isCart = ref(false)

const addToCart = async () => {
  // 需要校验用户是否登录，未登录跳转至登录页，否则允许加入购物车
  if (!token) {
    uni.showToast({
      title: '请登录后进行操作',
      mask: true,
      icon: 'none'
    });
    uni.redirectTo({
      url: '/pages/user/login'
    })
    return
  }


  //传递参数让子组件知道是加入购物车还是立即购买
  isCart.value = true
  goodsSkuRef.value.showSku()
}

const goTo = (path) => {
  if (path === "cart") {
    uni.switchTab({
      url: '/pages/tabBar/shop/shop'
    })
  } else {
    uni.switchTab({
      url: '/pages/tabBar/index'
    })
  }

}
const toOrder = async () => {
  myRouter(`/pages/orderInfo/orderInfo?skuID=${data.value.skus[0].ID}&goodID=${data.value.skus[0].goodID}`)
}

const goodsTapPay = () => {
  isCart.value = false
  goodsSkuRef.value.showSku()
}

const addCollect = async () => {
  if (token) {
    // 如果已登录并且未收藏 则允许进行收藏操作
    const res = await createCollect({
      goodID: Number(goodID.value)
    })
    if (res.code === 0) {
      collectionFlag.value = !collectionFlag.value
      uni.showToast({
        title: collectionFlag.value ? '已收藏' : '已取消',
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
</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
}

.goods_title {
  padding: 0 32rpx;
}

.price_color {
  color: #FAAD14;
  background: linear-gradient(360deg, #FF3636 0%, #FF6600 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.goods_item_box {
  padding: 20rpx 32rpx;
}

.goods_item_cla {
  margin-right: 16rpx;
}

.goods_item_icon {
  width: 12rpx;
  height: 24rpx;
}

.bb_title {
  padding: 24rpx 0;

  text {
    margin: 0 16rpx;
  }

  view {
    width: 100rpx;
    height: 2rpx;
  }
}

.goods_nav_view {
  height: 100rpx;
  width: 100%;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.g_nav_box {
  padding: 0 32rpx;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100rpx;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
  box-shadow: 0rpx -2rpx 6rpx 0rpx rgba(215, 215, 215, 0.5);
}

.g_nav_icon {
  width: 60rpx;
  height: 60rpx;
}

.btns_box {
  border-radius: 6rpx;
  overflow: hidden;

  view {
    width: 200rpx;
    height: 80rpx;
    line-height: 80rpx;
  }
}

.item_head {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
}

.evaluate_num {
  margin-top: 16rpx;
  padding: 6rpx 16rpx;
}

</style>
