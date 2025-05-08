<template>
  <view class="test" style="padding: 3rem 1.5rem;">
    <wu-sku v-model="skuShow" :data="skus" :themeColor="[226, 35, 26]"
            defaultCover="https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/default.jpeg"
            btnConfirmText="购买" notSelectSku="请选择完整的商品信息" @skuChange="skuChange" @confirm="skuConfirm"></wu-sku>
    <!-- 打开sku组件 -->
    <button @click="skuShow = true">打开sku组件</button>
  </view>
</template>

<script setup>
import {
  ref,
  reactive,
  inject,
  toRefs
} from "vue";
import {} from '@dcloudio/uni-app';

let skuShow = ref(false);
// sku列表
let skus = ref(// 以下代码可以直接复制使用
    [
      {
        id: 1,
        price: 7000,
        stock: 30,
        sku_attrs: {
          '机身颜色': {
            name: '深空黑色',
            img: 'https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/heise.jpeg',
          },
          '储存容量': '128G',
          '套装': '快充套装'
        }
      },
      {
        id: 2,
        price: 8500,
        stock: 10,
        sku_attrs: {
          '机身颜色': {
            name: '暗紫色',
            img: 'https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/anyezi.jpeg',
          },
          '储存容量': '256G',
          '套装': '快充套装'
        }
      },
      {
        id: 3,
        price: 9500,
        stock: 0,
        sku_attrs: {
          '机身颜色': {
            name: '暗紫色',
            img: 'https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/anyezi.jpeg',
          },
          '储存容量': '512G',
          '套装': 'AirPods套装'
        }
      },
      {
        id: 4,
        price: 9200,
        stock: 60,
        sku_attrs: {
          '机身颜色': {
            name: '银色',
            img: 'https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/baise.jpeg',
          },
          '储存容量': '512G',
          '套装': 'AirPods套装'
        }
      },
      {
        id: 5,
        price: 9200,
        stock: 80,
        sku_attrs: {
          '机身颜色': {
            name: '金色',
            img: 'https://mp-16389c52-85e7-413f-b566-3da08e2fe054.cdn.bspapp.com/img/sku-test/jinse.jpeg',
          },
          '储存容量': '512G',
          '套装': 'AirPods3套装'
        }
      }
    ])

// sku发生改变事件, 选择完整的sku则回返回 否则sku的值为{}
let skuChange = (sku) => {
  console.log(sku);
}
// sku确认事件
let skuConfirm = (e) => {
  console.log(e);
}
</script>
