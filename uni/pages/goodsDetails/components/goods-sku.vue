<template>
    <uni-popup ref="popupRef" type="bottom">
      <view class="popup_box bgc_fff boxs_bb">
        <view class="flex b_b_2 p_b_24">
          <view class="font_40 flex-fitem flex flex-aife">
            <!-- <text>{{ currentPrice ? currentPrice : list[0] && list[0].price / 100 }}</text> -->
            <view class="flex sku-info-box">
              <image class="product-img" :src="getUrl(selectSKU?selectSKU.picture: good.imageUrl)" alt="" />
              <view class="font_40 flex flexc-jsb sku-info">
                <p class="font_28 product-name">{{ selectSKU?selectSKU.name : good.title }} </p>
                <p class="color_ff0003"><span class="font_28">¥ </span>{{ selectSKU?selectSKU.price/100 : good.price/100}}</p>
                <up-number-box v-if="selectSKU"
                               v-model="currentQuantity"
                               :min="1" :max="selectSKU?selectSKU.inventory:0"
                               button-size="20">
                </up-number-box>
              </view>
              </view>
          </view>
          <image class="goods_close" mode=""
                 src="../../../static/images/close-icons.png"
                 @tap="popupRef.close()">
          </image>
        </view>
        <view class="product-box">
          <view class="p_b_24">请选择：</view>
          <view class="product-detail" v-for="(specs, key) in specsMap" :key="index" >
            <view class="color_333 font_28 m_b_16 ">
              {{ key }}
            </view>
          <view class="flex flex-fww" style="gap:20rpx;">
            <view v-for="spec in specs" class="color_333 font_28 m_b_16"  @click="selectSpec(key,spec)" >
              <view class="product-tags tac" :class="{
              active:selectMap[key] === spec.text,
              disabled:spec.disabled
            }">{{ spec.text }} </view>
            </view>
          </view>
          </view>

          <view class="product-inventory b_b_2 p_b_24 p_t_24 flexr-jsb flex-aic boxs_bb" style="height: 88rpx;margin-bottom: 32rpx;">
            <text class="font_28 color_333 m_r_8">数量</text>
            <text class="flex-fitem font_24 color_999">库存 {{ selectSKU?selectSKU.inventory : 0}}</text>
          </view>
        </view>
        <view class="tac">
          <button v-if="selectSKU&&selectSKU.inventory >= 1 && !isCart" class="confirm" @tap="submit()">确认</button>
          <button v-if="!selectSKU && !isCart" class="confirm" disabled>确认</button>
          <button v-if="selectSKU&&selectSKU.inventory >= 1 && isCart" class="confirm" @tap="addToCart()">加入购物车</button>
          <button v-if="!selectSKU && isCart" class="confirm" disabled>加入购物车</button>
        </view>
      </view>
    </uni-popup>
</template>

<script setup>
import {
  ref,
  watch
} from "vue";
import {placeOrder} from '@/api/order.js'
import {useUserStore} from "@/pinia/modules/user";
import {addCart} from '@/api/cart.js'
import {getUrl} from "@/utils/url.js"
const emit = defineEmits(['toOrder'])
let popupRef = ref()
let props = defineProps({
  isCart: {
    type: Boolean,
    default: false
  },
  good: {
    type: Object,
    default: {}
  }
})
const showSku = (val) => {
  popupRef.value.open()
}
// 创建响应式数据
const userStore = useUserStore()
const token = userStore.token || ''
const currentQuantity = ref(1);
const specsMap = ref({})
const selectMap = ref({})
let selectSKU = ref(undefined)


const getSkuInfo = (skus,key) =>{
	const info = {}
	skus.forEach((sku)=>{
		sku.specs.forEach(s=>{
			if(s.label === key){
				info[s.value] = true
			}
		})
	})
	return Object.keys(info).map(key=>{
		return {text:key,disabled:false}
	})
}

const selectSpec = (key,spec) =>{
	const hasSpecs = []
	if(spec.disabled){
		return
	}
	for(let k in specsMap.value){
		if(k === key){
			continue
		}
		specsMap.value[k].forEach(item=>{
			item.disabled = true
		})
	}


	 props.good.skus.forEach(sku=>{
		 let hasKey = false
		 sku.specs.forEach(s=>{
			 if (s.label === key && s.value === spec.text){
				 hasKey = true
			 }
		 })
		 if(hasKey){
			 hasSpecs.push(sku.specs)
		 }
	 })

	 hasSpecs.forEach(sp=>{
		 sp.forEach(s=>{
			 specsMap.value[s.label].forEach(item=>{
				 if(key === s.label){
					 return
				 }
			 	if (s.value === item.text){
			 		item.disabled = false
			 	}
			 })
		 })
	 })

	selectMap.value[key] = spec.text
}

const init = () => {
	const specs = props.good.specs
	const skus = props.good.skus
	specs.forEach((item)=>{
		specsMap.value[item.name] = getSkuInfo(skus,item.name)
		selectMap.value[item.name] = ""
	})
}
init()

watch(()=>selectMap,(newValue)=>{
	const sku = props.good.skus.find(sku=>{
		let allSame = true
		sku.specs.forEach(item=>{
			if (newValue.value[item.label] != item.value){
				allSame = false
			}
		})
		return allSame
	})
	selectSKU = sku
},{deep:true})

const submit = async () => {
  if (token) {
    popupRef.value.close()
    let detail = {
      "detail": [
        {
          "goodID": selectSKU.goodID,
          "skuID": selectSKU.ID,
          "quantity": currentQuantity.value
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
    goodID: selectSKU.goodID,
    skuID: selectSKU.ID
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
  showSku
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
