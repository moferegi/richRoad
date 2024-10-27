<template>
  <view>
    <view class="bgc_fff evaluate_goods_box">
      <view class="flex m_b_24">
        <image class="item_goods_img m_r_16"
               :src="getUrl(SKU.picture)"
               mode=""></image>
        <view class="flex-fitem">
          <view class="m_b_8 flex-aic flexr-jsb">
            <text class="color_333 font_24 text_nowrap"
                  style="max-width: 480rpx;">{{ SKU.name}}</text>
          </view>
          <view class="color_999 font_24 m_b_4 text_nowrap" style="max-width: 560rpx;">
            <view  v-for="(sku, index) in SKU.attrs" :key="index">
              <view class="color_333 m_b_16">
                {{sku.label}}:{{sku.value}}
              </view>
            </view>
          </view>

        </view>
      </view>
      <view class="flex flex-aic m_b_24">
        <text class="color_333 font_30 m_r_16">商品评分</text>
        <uni-rate :disabled="isCheck" size="22" disabled-color="#fe5572" active-color="#fe5572" :value="rating"  @change="onChange"/>
      </view>
      <textarea class="bgc_f9f9f9 evaluate_textarea_box boxs_bb m_b_24" adjust-position :maxlength="200"
                placeholder="请输入您对该商品的评价…" v-model="content" :disabled="isCheck"></textarea>
      <up-upload
      	:fileList="pics"
      	@afterRead="afterRead"
      	@delete="deletePic"
		:disabled="isCheck"
      	name="file"
      	multiple
      	:maxCount="6"
      	></up-upload>
    </view>
    <view class="evaluate_form_nav_view">
      <view v-if="!isCheck" class="evaluate_form_nav_box pos_f bgc_fff flex-aic flexr-jsc">
        <button class="subBtn" @tap="submit">提交</button>
      </view>
    </view>
	<up-toast ref="uToastRef"></up-toast>
  </view>
</template>

<script setup>
import {ref} from 'vue'
import { selfOrderComment} from "@/api/order"
import {createComment} from "@/api/comment"
import {onLoad} from '@dcloudio/uni-app'
import {baseUrl} from "@/utils/request.js"
    import {getUrl} from "@/utils/url.js"
const orderID = ref(0)
const goodID = ref(0)
const SKUID = ref(0)
const content = ref("")
const data = ref({})
const SKU = ref({})
const rating = ref(0)
const isCheck = ref(false)
const pics = ref([])
const uToastRef = ref(null)

const init = async () => {
  const res = await selfOrderComment(orderID.value, goodID.value, SKUID.value)
  if (res.code === 0) {
    data.value = res.data
    SKU.value = res.data.detail.sku
    if (res.data.comment) {
      res.data.comment.content != "" ? isCheck.value = true : isCheck.value = false
      rating.value = res.data.comment.rating
      content.value = res.data.comment.content
    }
  }
}

onLoad(async (options) => {
  orderID.value = Number(options.orderID)
  goodID.value = Number(options.goodID)
  SKUID.value = Number(options.SKUID)
  setTimeout(() => {
    init()
  }, 200)
})

const onChange = (e) => {
  rating.value = e.value
}

// 提交评论 createComment
const submit = async () => {
  const picsArr = pics.value.map(item=>item.url)
  const req = {
    orderID: orderID.value,
    goodID: goodID.value,
    SKUID: SKUID.value,
    rating: rating.value,
    content: content.value,
    pics: picsArr
  }
  const res = await createComment(req)
  if (res.code === 0) {
    uni.showToast({
      title: '评论成功',
      icon: 'success',
      duration: 2000
    })
    uni.redirectTo({
      url: '/pages/order/order'
    })
  }
}


// 删除图片
const deletePic = (event) => {
  pics.value.splice(event.index, 1);
};

// 新增图片
const afterRead = async (event) => {
  // 当设置 mutiple 为 true 时, file 为数组格式，否则为对象格式
  let lists = [].concat(event.file);
  let fileListLen = pics.value.length;
  lists.map((item) => {
    pics.value.push({
      ...item,
      status: 'uploading',
      message: '上传中',
    });
  });
  for (let i = 0; i < lists.length; i++) {
    const result = await uploadFilePromise(lists[i].url);
	if(!result){
		pics.value.splice(fileListLen,1)
		uToastRef.value.show({
			type: 'default',
			message: "图片上传失败",
		})
		return
	}
    let item = pics.value[fileListLen];
    pics.value.splice(fileListLen, 1, {
      ...item,
      status: 'success',
      message: '',
      url: result,
    });
    fileListLen++;
  }
};

const uploadFilePromise = (url) => {
  return new Promise((resolve, reject) => {
    let a = uni.uploadFile({
      url: baseUrl+'/fileUploadAndDownload/upload?noSave=1',
	  header:{
		"x-token":uni.getStorageSync('x-token')  
	  },
      filePath: url,
      name: 'file',
      success: (res) => {
		  const data = JSON.parse(res.data)
		  if(data.code!==0){
			  reject("")
			return
		  }
         resolve(data.data.file.url);
      },
    });
  });
}

</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
}

.item_goods_img {
  width: 108rpx;
  height: 108rpx;
  border-radius: 2rpx;
  overflow: hidden;
}

.evaluate_goods_box {
  padding: 24rpx 24rpx 30rpx;
}

.evaluate_textarea_box {
  width: 100%;
  padding: 12rpx;
}

.evaluate_form_nav_view {
  height: 100rpx;
  width: 100%;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.evaluate_form_nav_box {
  bottom: 0;
  left: 0;
  right: 0;
  height: 100rpx;
  width: 100%;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.subBtn {
  width: 90%;
  color: #fff;
  border-color: #fe5572;
  background-color: #fe5572;
}

</style>
