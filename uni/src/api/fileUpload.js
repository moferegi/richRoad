import { request } from '@/utils/request.js'

/**
 * 生成防盗链签名URL（付费视频使用）
 * @param {string} filePath 文件完整URL或相对路径
 * @returns {Promise} {url: string} 带签名参数的时效URL
 */
export const signURL = (filePath) => {
  return request({
    url: '/fileUploadAndDownload/signURL',
    method: 'post',
    data: { filePath }
  })
}
