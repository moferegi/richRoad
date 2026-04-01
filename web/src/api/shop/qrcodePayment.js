import service from '@/utils/request'

// @Tags QrcodePayment
// @Summary 创建收款码
export const createQrcodePayment = (data) => {
  return service({
    url: '/qrcodePayment/createQrcodePayment',
    method: 'post',
    data
  })
}

// @Tags QrcodePayment
// @Summary 删除收款码
export const deleteQrcodePayment = (params) => {
  return service({
    url: '/qrcodePayment/deleteQrcodePayment',
    method: 'delete',
    params
  })
}

// @Tags QrcodePayment
// @Summary 更新收款码
export const updateQrcodePayment = (data) => {
  return service({
    url: '/qrcodePayment/updateQrcodePayment',
    method: 'put',
    data
  })
}

// @Tags QrcodePayment
// @Summary 分页获取收款码列表
export const getQrcodePaymentList = (params) => {
  return service({
    url: '/qrcodePayment/getQrcodePaymentList',
    method: 'get',
    params
  })
}

// @Tags QrcodePayment
// @Summary 获取启用的收款码（客户端）
export const getEnabledQrcodePayments = () => {
  return service({
    url: '/qrcodePayment/getEnabledQrcodePayments',
    method: 'get'
  })
}
