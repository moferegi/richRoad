import service from '@/utils/request'

// @Tags PhoneAreaCode
// @Summary 创建国际区号
export const createPhoneAreaCode = (data) => {
  return service({
    url: '/phoneAreaCode/createPhoneAreaCode',
    method: 'post',
    data
  })
}

// @Tags PhoneAreaCode
// @Summary 删除国际区号
export const deletePhoneAreaCode = (params) => {
  return service({
    url: '/phoneAreaCode/deletePhoneAreaCode',
    method: 'delete',
    params
  })
}

// @Tags PhoneAreaCode
// @Summary 更新国际区号
export const updatePhoneAreaCode = (data) => {
  return service({
    url: '/phoneAreaCode/updatePhoneAreaCode',
    method: 'put',
    data
  })
}

// @Tags PhoneAreaCode
// @Summary 获取区号列表（管理端）
export const getPhoneAreaCodeList = (params) => {
  return service({
    url: '/phoneAreaCode/getPhoneAreaCodeList',
    method: 'get',
    params
  })
}

// @Tags PhoneAreaCode
// @Summary 获取启用的区号列表（客户端）
export const getEnabledPhoneAreaCodes = () => {
  return service({
    url: '/phoneAreaCode/getEnabledPhoneAreaCodes',
    method: 'get'
  })
}
