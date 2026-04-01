import service from '@/utils/request'

// @Tags Popup
// @Summary 创建弹窗
export const createPopup = (data) => {
  return service({
    url: '/popup/createPopup',
    method: 'post',
    data
  })
}

// @Tags Popup
// @Summary 删除弹窗
export const deletePopup = (params) => {
  return service({
    url: '/popup/deletePopup',
    method: 'delete',
    params
  })
}

// @Tags Popup
// @Summary 更新弹窗
export const updatePopup = (data) => {
  return service({
    url: '/popup/updatePopup',
    method: 'put',
    data
  })
}

// @Tags Popup
// @Summary 分页获取弹窗列表
export const getPopupList = (params) => {
  return service({
    url: '/popup/getPopupList',
    method: 'get',
    params
  })
}

// @Tags Popup
// @Summary 获取当前生效弹窗（客户端）
export const getActivePopups = (params) => {
  return service({
    url: '/popup/getActivePopups',
    method: 'get',
    params
  })
}
