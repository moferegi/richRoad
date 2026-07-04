import service from '@/utils/request'
// @Tags FileUploadAndDownload
// @Summary 分页文件列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body modelInterface.PageInfo true "分页获取文件户列表"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /fileUploadAndDownload/getFileList [post]
export const getFileList = (data) => {
  return service({
    url: '/fileUploadAndDownload/getFileList',
    method: 'post',
    data
  })
}

// @Tags FileUploadAndDownload
// @Summary 删除文件
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body dbModel.FileUploadAndDownload true "传入文件里面id即可"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"返回成功"}"
// @Router /fileUploadAndDownload/deleteFile [post]
export const deleteFile = (data) => {
  return service({
    url: '/fileUploadAndDownload/deleteFile',
    method: 'post',
    data
  })
}

/**
 * 编辑文件名或者备注
 * @param data
 * @returns {*}
 */
export const editFileName = (data) => {
  return service({
    url: '/fileUploadAndDownload/editFileName',
    method: 'post',
    data
  })
}

/**
 * 导入URL
 * @param data
 * @returns {*}
 */
export const importURL = (data) => {
  return service({
    url: '/fileUploadAndDownload/importURL',
    method: 'post',
    data
  })
}

/**
 * 生成扫码上传一次性票据
 * @param {Object} data
 * @param {number} data.classId 分类ID
 * @param {string} data.folder 上传目录
 * @param {string} data.uploadType 上传类型
 * @param {string} data.uploadPosition 上传位置
 * @returns {Promise} {ticket: string, expiresAt: number}
 */
export const createScanUploadTicket = (data) => {
  return service({
    url: '/fileUploadAndDownload/createScanUploadTicket',
    method: 'post',
    data
  })
}


// 上传文件 暂时用于头像上传
export const uploadFile = (data) => {
  return service({
    url: "/fileUploadAndDownload/upload",
    method: "post",
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data,
  });
};

/**
 * 生成防盗链签名URL
 * @param {Object} data
 * @param {string} data.filePath 文件相对路径
 * @returns {Promise} {url: string} 签名后的完整URL
 */
export const signURL = (data) => {
  return service({
    url: '/fileUploadAndDownload/signURL',
    method: 'post',
    data
  })
}

/**
 * 获取防盗链配置信息
 * @returns {Promise} {enabled: boolean, cdnDomain: string}
 */
export const getHotlinkConfig = () => {
  return service({
    url: '/fileUploadAndDownload/hotlinkConfig',
    method: 'get'
  })
}

/**
 * 列举OSS存储中的文件夹
 * @returns {Promise} {folders: string[]}
 */
export const getOSSFolders = () => {
  return service({
    url: '/fileUploadAndDownload/listFolders',
    method: 'get'
  })
}