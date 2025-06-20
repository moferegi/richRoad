<template>
  <view class="evaluate-container">
<!--    <view class="header">
      <text class="title">{{ isCheck ? '查看评价' : '商品评价' }}</text>
    </view>-->

    <view class="item-card">
      <view class="item-header">
        <image
            :src="getUrl(SKU.picture)"
            class="evaluate_pic_img"
            mode="aspectFill"
        />
        <view class="item-info">
          <text class="item-name">{{ SKU.name }}</text>
          <view v-for="(sku, index) in SKU.attrs" :key="index" class="item-spec">
            {{ sku.label }}:{{ sku.value }}
          </view>
        </view>
      </view>

      <view class="rating-section">
        <text class="section-title">商品评分</text>
        <view class="stars">
          <text
              v-for="star in 5"
              :key="star"
              :class="['star', { active: star <= rating }]"
              @click="!isCheck && setRating(star)"
          >★</text>
        </view>
        <text class="rating-text">{{ getRatingText(rating) }}</text>
      </view>

      <view class="comment-section">
        <text class="section-title">评价内容</text>
        <textarea
            v-model="content"
            class="comment-input"
            :placeholder="isCheck ? '' : '请输入您对该商品的评价…'"
            maxlength="200"
            show-confirm-bar="false"
            :disabled="isCheck"
        ></textarea>
        <text class="char-count" v-if="!isCheck">{{ content.length }}/200</text>
      </view>

      <view class="image-section">
        <text class="section-title">{{ isCheck ? '评价图片' : '上传图片（最多9张）' }}</text>
        <view class="image-upload">
          <view class="uploaded-images">
            <view
                v-for="(pic, index) in pics"
                :key="index"
                class="image-item"
                @tap="previewImage(index)"
            >
              <image
                  :src="getUrl(pic.url || pic.tempFilePath)"
                  mode="aspectFill"
                  class="uploaded-image"
                  @error="onImageError(index)"
              />
              <text v-if="!isCheck" class="delete-btn" @click.stop="deletePic(index)">×</text>
            </view>
            <view
                v-if="!isCheck && pics.length < 9"
                class="add-image-btn"
                @click="chooseImage"
            >
              <text class="add-icon">+</text>
              <text class="add-text">添加图片</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <view class="submit-section" v-if="!isCheck">
      <button class="submit-btn" @click="submit" :disabled="isSubmitting">
        {{ isSubmitting ? '提交中...' : '提交评价' }}
      </button>
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


// 在init方法中修改图片数据处理
const init = async () => {
  const res = await selfOrderComment(orderID.value, goodID.value, SKUID.value)
  if (res.code === 0) {
    data.value = res.data
    SKU.value = res.data.detail.sku

    if (res.data.comment) {
      // 判断是否为查看模式
      isCheck.value = res.data.comment.content !== ""
      rating.value = res.data.comment.rating
      content.value = res.data.comment.content

      // 处理评价图片数据 - 直接使用原始数据，让getUrl处理URL转换
      if (res.data.comment.pics && Array.isArray(res.data.comment.pics) && res.data.comment.pics.length > 0) {
        pics.value = res.data.comment.pics.map(pic => ({
          url: typeof pic === 'string' ? pic : pic.url, // 保持原始URL，交给getUrl处理
          isUploaded: true // 标记为已上传的图片
        }))
      } else {
        pics.value = [] // 没有图片时清空数组
      }
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

// 设置评分
const setRating = (star) => {
  if (!isCheck.value) {
    rating.value = star
  }
}

// 获取评分文本
const getRatingText = (rating) => {
  const texts = ['', '很差', '一般', '满意', '很好', '非常好']
  return texts[rating] || ''
}

// 选择图片
const chooseImage = () => {
  if (pics.value.length >= 9) {
    uni.showToast({
      title: '最多只能上传9张图片',
      icon: 'none'
    })
    return
  }

  uni.chooseImage({
    count: 9 - pics.value.length,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const newPics = res.tempFilePaths.map(path => ({
        url: path,
        tempFilePath: path
      }))
      pics.value.push(...newPics)
    },
    fail: (err) => {
      console.error('选择图片失败:', err)
    }
  })
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

// 新增图片预览方法
const previewImage = (index) => {
  const urls = pics.value.map(pic => getUrl(pic.url))
  uni.previewImage({
    current: index,
    urls: urls,
    fail: (err) => {
      console.error('图片预览失败:', err)
      uni.showToast({
        title: '图片加载失败',
        icon: 'error'
      })
    }
  })
}

// 图片加载错误处理
const onImageError = (index) => {
  const failedPic = pics.value[index];
  console.error(`图片加载失败:`, {
    originalUrl: failedPic.url,
    processedUrl: getUrl(failedPic.url),
    index: index
  });
}

</script>

<style lang="scss">
page {
  background-color: #f8f8f8;
}

.evaluate-container {
  min-height: 100vh;
  background: #f7f7f7;
  padding-bottom: 120rpx;
}

.header {
  background: white;
  padding: 30rpx;
  text-align: center;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;

  .title {
    font-size: 36rpx;
    font-weight: 600;
    color: #333;
  }
}

.item-card {
  background: white;
  margin: 20rpx;
  border-radius: 20rpx;
  overflow: hidden;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
  animation: fadeIn 0.3s ease-out;
}

.item-header {
  display: flex;
  padding: 30rpx;
  border-bottom: 1rpx solid #f0f0f0;

  .evaluate_pic_img {
    width: 120rpx;
    height: 120rpx;
    border-radius: 12rpx;
    margin-right: 20rpx;
    border: 1rpx solid #e0e0e0;
  }

  .item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;

    .item-name {
      font-size: 30rpx;
      font-weight: 500;
      color: #333;
      margin-bottom: 10rpx;
      line-height: 1.4;
    }

    .item-spec {
      font-size: 24rpx;
      color: #666;
      margin-bottom: 5rpx;
    }
  }
}

.rating-section {
  padding: 30rpx;
  border-bottom: 1rpx solid #f0f0f0;

  .section-title {
    font-size: 28rpx;
    font-weight: 500;
    color: #333;
    margin-bottom: 15rpx;
  }

  .stars {
    display: flex;
    align-items: center;
    margin-bottom: 10rpx;

    .star {
      font-size: 40rpx;
      color: #ddd;
      margin-right: 8rpx;
      cursor: pointer;
      transition: all 0.2s ease;

      &.active {
        color: #ffd700;
      }

      &:hover {
        transform: scale(1.1);
      }
    }
  }

  .rating-text {
    font-size: 24rpx;
    color: #666;
    font-style: italic;
  }
}

.comment-section {
  padding: 30rpx;
  border-bottom: 1rpx solid #f0f0f0;

  .section-title {
    font-size: 28rpx;
    font-weight: 500;
    color: #333;
    margin-bottom: 15rpx;
  }

  .comment-input {
    width: 100%;
    min-height: 120rpx;
    padding: 20rpx;
    border: 1rpx solid #e0e0e0;
    border-radius: 12rpx;
    font-size: 28rpx;
    line-height: 1.5;
    background-color: #fafafa;
    box-sizing: border-box;
    transition: border-color 0.2s ease;

    &:focus {
      border-color: #667eea;
      background-color: white;
    }
  }

  .char-count {
    display: block;
    text-align: right;
    font-size: 24rpx;
    color: #999;
    margin-top: 10rpx;
  }
}

.image-section {
  padding: 30rpx;

  .section-title {
    font-size: 28rpx;
    font-weight: 500;
    color: #333;
    margin-bottom: 15rpx;
  }

  .image-upload {
    .uploaded-images {
      display: flex;
      flex-wrap: wrap;
      gap: 15rpx;

      .image-item {
        position: relative;
        width: 120rpx;
        height: 120rpx;
        border-radius: 12rpx;
        overflow: hidden;
        border: 1rpx solid #e0e0e0;

        .uploaded-image {
          width: 100%;
          height: 100%;
        }

        .delete-btn {
          position: absolute;
          top: -5rpx;
          right: -5rpx;
          width: 30rpx;
          height: 30rpx;
          background: #ff4757;
          color: white;
          border-radius: 50%;
          font-size: 20rpx;
          display: flex;
          align-items: center;
          justify-content: center;
          line-height: 1;
          cursor: pointer;
          box-shadow: 0 2rpx 8rpx rgba(255, 71, 87, 0.3);
        }
      }

      .add-image-btn {
        width: 120rpx;
        height: 120rpx;
        border: 2rpx dashed #ccc;
        border-radius: 12rpx;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        background: #fafafa;
        cursor: pointer;
        transition: all 0.2s ease;

        &:hover {
          border-color: #667eea;
          background: #f0f2ff;
        }

        .add-icon {
          font-size: 40rpx;
          color: #999;
          margin-bottom: 5rpx;
        }

        .add-text {
          font-size: 20rpx;
          color: #999;
        }
      }
    }
  }
}

.submit-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  padding: 20rpx 30rpx;
  border-top: 1rpx solid #e0e0e0;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.1);

  .submit-btn {
    width: 100%;
    height: 80rpx;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 40rpx;
    font-size: 32rpx;
    font-weight: 500;
    box-shadow: 0 4rpx 15rpx rgba(102, 126, 234, 0.4);
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-2rpx);
      box-shadow: 0 6rpx 20rpx rgba(102, 126, 234, 0.5);
    }

    &:disabled {
      opacity: 0.6;
      transform: none;
      box-shadow: 0 4rpx 15rpx rgba(102, 126, 234, 0.2);
    }
  }
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.item-card {
  animation: fadeIn 0.3s ease-out;
}
</style>
