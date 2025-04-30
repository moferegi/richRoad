<template>
  <wu-sku v-model="skuShow"
          :data="skus"
          :themeColor="[254, 85, 141]"
          :defaultCover="selectSKU.picture || getUrl(good.imageUrl)"
          :btnConfirmText="isCart? '加入购物车' : '立即购买'"
          notSelectSku="请选择完整的商品信息"
          @skuChange="skuChange"
          @confirm="skuConfirm"></wu-sku>
</template>

<script setup>
import {
  ref
} from "vue";
import {placeOrder} from '@/api/order.js'
import {useUserStore} from "@/pinia/modules/user";
import {addCart} from '@/api/cart.js'
import {getUrl} from "@/utils/url.js"
let props = defineProps({
  isCart: {
    type: Boolean,
    default: false
  },
  good: {
    type: Object,
    default: {}
  },
})

const skuShow = ref(false)

const skus = ref([])

const showSku = (value) => {
  skuShow.value = true
    skus.value = props.good.skus.map(item=>{
        const sku_attrs = {}
        item.specs.forEach(spec=>{
          sku_attrs[spec.label] = spec.value
        })
        return {
          id: item.ID,
          picture: getUrl(item.picture),
          price: item.price,
          stock: item.inventory,
          sku_attrs: sku_attrs
        }
      })
}
const closeSku = () => {
  skuShow.value = false
}
// 创建响应式数据
const userStore = useUserStore()
const token = userStore.token || ''
const selectSKU = ref({})

const skuChange = (sku) =>{
  selectSKU.value = sku
}
const skuConfirm = (select) =>{
  if(props.isCart){
    addToCart(select.num)
  }
  submit(select.num)
}

const submit = async (num) => {
  if (token) {
    skuShow.value = false
    let detail = {
      "detail": [
        {
          "goodID": props.good.ID,
          "skuID": selectSKU.value.id,
          "quantity": num
        }
      ]
    };
    const res = await placeOrder(detail);
    if (res.code === 0) {
      uni.navigateTo({
        url: `/pages/orderInfo/orderInfo?orderID=${res.data.orderID}&type='single'`
      })
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

const addToCart = async () => {
  // 找到当前选中的sku的SKUid和goodID
  const params = {
    goodID: props.good.ID,
    skuID: selectSKU.value.ID,
    num: 1
  }
  const res = await addCart(params)
  if (res.code === 0) {
    uni.showToast({
      title: '添加成功，在购物车等亲!',
      mask: true,
      icon: 'none'
    });
  } else {
    uni.showToast({
      title: '添加失败，请稍后重试',
      mask: true,
      icon: 'none'
    });
  }
}

defineExpose({
  showSku,
  closeSku
})
</script>

<style lang="scss" scoped>
.imgBox{
  width: 100%;
  background: #000;
}
.product-box {
  padding: 20rpx;
  .product-detail {
    padding-left: 20rpx;
    .product-tags {
      min-width: 100rpx;
      padding: 20rpx;
      border-radius: 10rpx;
      background: #eeeeee;
    }
    .tag_box {
      display: inline-block;
      padding: 0 20rpx;
      height: 40rpx;
      line-height: 40rpx;
      border-radius: 4rpx;
      border: 2rpx solid;
      border-color: #c3c3c3;
    }
  }
}
.product-img {
  width: 180rpx;
  height: 200rpx;
  padding: 0 20rpx;
  border-radius: 20rpx;
}
.product-name{
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  text-overflow: ellipsis;
  max-height: 100rpx;
}
.selected {
  /* 选中时的文字样式 */
  border-color: #fe5572 !important;
  background-color: #fe5572 !important;
  color: white !important;
}

.selected-border {
  border-color: #ff7000 !important; /* 或者其他你想要的颜色 */
}

.popup_box {
  padding: 32rpx 32rpx 8rpx 32rpx;
  border-radius: 16rpx 16rpx 0rpx 0rpx;
}

.goods_chart {
  width: 108rpx;
  height: 108rpx;
  border-radius: 6rpx;
  margin-right: 16rpx;
}

.goods_close {
  width: 32rpx;
  height: 32rpx;
}

.confirm {
  width: 100%;
  border-color: #fe5572;
  background-color: #fe5572;
  color: white;
}

.lacking {
  width: 100%;
  border-color: #c0c0c0;
  background-color: #c0c0c0;
  color: white;
}
.active {
  border-color: #fe5572 !important;
  background-color: #fe5572 !important;
  color: white;
}
.disabled{
	color: #c0c0c0;
  background: #eeeeee;
}
.sku-info-box{
	width: 100%;
	.sku-info{
		flex:1;
	}
}
</style>
