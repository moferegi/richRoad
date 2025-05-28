<template>
  <view class="bgc_fff">
    <view class=" evaluate_goods_box">
      <view class="flex m_b_24">
        <image :src="getUrl(SKU.picture)"
               class="item_goods_img m_r_16"
               mode=""></image>
        <view class="flex-fitem">
          <view class="m_b_8 flex-aic flexr-jsb">
            <text class="color_333 font_24 text_nowrap"
                  style="max-width: 480rpx;">{{ SKU.name }}
            </text>
          </view>
          <view class="color_999 font_24 m_b_4 text_nowrap" style="max-width: 560rpx;">
            <view v-for="(sku, index) in SKU.attrs" :key="index">
              <view class="color_333 m_b_16">
                {{ sku.label }}:{{ sku.value }}
              </view>
            </view>
          </view>

        </view>
      </view>
      <view class="flex flex-aic m_b_24">
        <text class="color_333 font_30 m_r_16">商品评分</text>
        <uni-rate :disabled="isCheck" :value="rating" active-color="#fe5572" disabled-color="#fe5572" size="22"
                  @change="onChange"/>
      </view>
      <textarea v-model="content" :disabled="isCheck" :maxlength="200"
                adjust-position class="bgc_f9f9f9 evaluate_textarea_box boxs_bb m_b_24" placeholder="请输入您对该商品的评价…"></textarea>
      <uni-file-picker
          v-model="pics"
          :auto-upload="false"
          :disabled="isCheck"
          :source-type="['album', 'camera']"
          limit="9"
          title="最多选择9张图片"
          @delete="deletePic"
          @select="afterRead"
      >
      </uni-file-picker>
      <!--      <uni-file-picker
                v-model="pics"
                file-mediatype="image"
                mode="grid"
                file-extname="png,jpg"
                :limit="1"
                @success="success"
                @fail="fail"
                @select="select"
            />-->
    </view>
    <view class="evaluate_form_nav_view">
      <view v-if="!isCheck" class="evaluate_form_nav_box pos_f m_b_20  flex-aic flexr-jsc">
        <button
            :disabled="isSubmitting"
            class="subBtn"
            @tap="submit"
        >
          {{ isSubmitting ? '提交中...' : '提交' }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import {ref} from 'vue'
import {selfOrderComment} from "@/api/order"
import {createComment} from "@/api/comment"
import {onLoad} from '@dcloudio/uni-app'
import {baseUrl} from "@/utils/request.js"
import {getUrl} from "@/utils/url.js"
import UniFilePicker from "@/uni_modules/uni-file-picker/components/uni-file-picker/uni-file-picker.vue";

const orderID = ref(0)
const goodID = ref(0)
const SKUID = ref(0)
const content = ref("")
const data = ref({})
const SKU = ref({})
const rating = ref(0)
const isCheck = ref(false)


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

// 修改提交方法
/*const submit = async () => {

  const req = {
    orderID: orderID.value,
    goodID: goodID.value,
    SKUID: SKUID.value,
    rating: rating.value,
    content: content.value,
    pics: currentUploadedUrls // 使用已上传的 URLs
  }

  console.log('提交数据:', req);

  // 取消注释进行实际提交
  /!*
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
  *!/
};*/

const submit = async () => {
  if (isSubmitting.value) return; // 防止重复提交

  try {
    isSubmitting.value = true;

    // 验证必填项
    if (rating.value === 0) {
      uni.showToast({title: "请选择商品评分", icon: "none"})
      return;
    }

    // 1. 先批量上传图片
    let uploadedPics = [];
    if (pics.value.length > 0) {
      uni.showLoading({
        title: `上传图片中...`,
        mask: true
      });

      try {
        uploadedPics = await uploadAllImages();
        console.log('图片上传成功:', uploadedPics);
      } catch (error) {
        uni.hideLoading();
        uni.showToast({
          type: 'default',
          title: error.message || '图片上传失败',
          duration: 2000
        });
        return;
      }

      uni.hideLoading();
    }

    // 2. 提交评论数据
    uni.showLoading({
      title: '提交评论中...',
      mask: true
    });

    const req = {
      orderID: orderID.value,
      goodID: goodID.value,
      SKUID: SKUID.value,
      rating: rating.value,
      content: content.value,
      pics: uploadedPics
    };

    console.log('提交评论数据:', req);

    const res = await createComment(req);
    uni.hideLoading();

    if (res.code === 0) {
      uni.showToast({title: "评论成功", icon: "success"})

      setTimeout(() => {
        uni.redirectTo({
          url: '/pages/order/order'
        });
      }, 100);
    } else {
      uni.showToast({title: res.msg || '评论提交失败', icon: "success"})
    }

  } catch (error) {
    uni.hideLoading();
    console.error('提交失败:', error);
    uni.showToast({
      type: 'default',
      title: '提交失败，请重试',
      icon: "error"
    });
  } finally {
    isSubmitting.value = false;
  }
};

const pics = ref([])
const isSubmitting = ref(false) // 提交状态

// 删除图片 - 纯前端操作
const deletePic = (event) => {
  pics.value.splice(event.index, 1);
};

/*// 新增图片
const afterRead = async (event) => {
  console.log(event);

  // 当设置 mutiple 为 true 时, file 为数组格式，否则为对象格式
  let lists = event.tempFilePaths.map((path, index) => ({
    tempFilePath: path, // 保存原始临时路径
    url: path, // uni-file-picker 需要的 url 字段
    ...event.tempFiles[index]
  }));

  console.log(lists);

  let fileListLen = pics.value.length;
  lists.map((item) => {
    pics.value.push({
      ...item,
      status: 'uploading',
      message: '上传中',
    });
  });

  for (let i = 0; i < lists.length; i++) {
    try {
      const result = await uploadFilePromise(lists[i].tempFilePath);
      // 上传成功，更新图片状态，但保留 tempFilePath
      let item = pics.value[fileListLen];
      pics.value.splice(fileListLen, 1, {
        ...item,
        status: 'success',
        message: '',
        url: result, // 这是上传后的网络 URL
        uploadedUrl: result, // 明确标识上传后的 URL
      });
      fileListLen++;
    } catch (error) {
      // 上传失败，移除该图片
      pics.value.splice(fileListLen, 1);
      uToastRef.value.show({
        type: 'default',
        message: "图片上传失败",
      });
      console.error('图片上传失败:', error);
    }
  }
};*/
// 新增图片 - 仅做预览处理，不上传
const afterRead = (event) => {
  console.log('选择图片:', event);

  // 将新选择的图片添加到预览列表
  const newPics = event.tempFilePaths.map((path, index) => ({
    url: path, // uni-file-picker 预览需要
    tempFilePath: path, // 保存原始路径用于后续上传
    ...event.tempFiles[index]
  }));

  // 添加到现有列表
  pics.value.push(...newPics);

  console.log('当前图片列表:', pics.value);
};
// 批量上传图片
const uploadAllImages = async () => {
  if (pics.value.length === 0) return [];

  const uploadPromises = pics.value.map((pic, index) =>
      uploadSingleImage(pic.tempFilePath, index)
  );

  try {
    const results = await Promise.allSettled(uploadPromises);
    const successUrls = [];
    const failedIndexes = [];

    results.forEach((result, index) => {
      if (result.status === 'fulfilled') {
        successUrls.push(result.value);
      } else {
        failedIndexes.push(index);
        console.error(`图片${index + 1}上传失败:`, result.reason);
      }
    });

    // 如果有失败的图片，给用户提示
    if (failedIndexes.length > 0) {
      const failedCount = failedIndexes.length;
      const successCount = successUrls.length;

      if (successCount === 0) {
        throw new Error('所有图片上传失败，请检查网络后重试');
      } else {
        uni.showToast({
          type: 'default',
          title: `${failedCount}张图片上传失败，已提交${successCount}张图片`,
          icon: 'error'
        });
      }
    }

    return successUrls;

  } catch (error) {
    console.error('批量上传失败:', error);
    throw error;
  }
};

// 单个图片上传
const uploadSingleImage = (tempFilePath, index) => {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: baseUrl + '/fileUploadAndDownload/upload?noSave=1',
      header: {
        "x-token": uni.getStorageSync('x-token')
      },
      filePath: tempFilePath,
      name: 'file',
      success: (res) => {
        try {
          const data = JSON.parse(res.data);
          if (data.code !== 0) {
            reject(new Error(data.msg || `图片${index + 1}上传失败`));
            return;
          }
          resolve(data.data.file.url);
        } catch (parseError) {
          reject(new Error(`图片${index + 1}响应数据解析失败`));
        }
      },
      fail: (error) => {
        reject(new Error(`图片${index + 1}网络请求失败: ${error.errMsg}`));
      }
    });
  });
};


const uploadFilePromise = (url) => {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: baseUrl + '/fileUploadAndDownload/upload?noSave=1',
      header: {
        "x-token": uni.getStorageSync('x-token')
      },
      filePath: url,
      name: 'file',
      success: (res) => {
        try {
          const data = JSON.parse(res.data);
          if (data.code !== 0) {
            reject(new Error(data.msg || '上传失败'));
            return;
          }
          resolve(data.data.file.url);
        } catch (parseError) {
          reject(new Error('响应数据解析失败'));
        }
      },
      fail: (error) => {
        reject(new Error('网络请求失败: ' + error.errMsg));
      }
    });
  });
};

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
