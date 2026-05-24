
import { t } from '@/utils/i18n.js'
import { baseUrl } from '@/utils/request.js'
/*
 *上传文件到阿里云oss
 *@param - filePath :文件的本地资源路径
  @param - path :上传oss哪个路径下
  postfix : 文件后缀
 *@param - successc:成功回调
 *@param - failc:失败回调
 */
export const  uploadFile = (filePath, fix,successc, failc) => {
  fix = fix.replace('.','')
  if (!filePath || filePath.length < 9) {
    uni.showToast({
      title: t('uploadFail'),
      icon: 'none'
    })
    return;
  }
  if(!fix){
	  fix = ''
  }
  return uni.uploadFile({
    url: `${baseUrl}/fileUploadAndDownload/upload?noSave=1`,
    filePath: filePath, //要上传文件资源的路径
    name: 'file', //必须填file
    header: {
      'x-token': uni.getStorageSync('x-token') || '',
    },
    success: function(res) {
      if (res.statusCode !== 200) {
        failc(new Error('上传错误:' + JSON.stringify(res)))
        return;
      }
      try {
        const data = JSON.parse(res.data)
        const fileURL = data && data.code === 0 && data.data && data.data.file && data.data.file.url
        if (!fileURL) {
          failc(new Error(data && data.msg ? data.msg : t('uploadFail')))
          return
        }
        successc(fileURL)
      } catch (e) {
        failc(new Error(t('uploadFail')))
      }
    },
    fail: function(err) {
      failc(err);
    },
  })
}
