import service from '@/utils/request'

// @Tags Presale
// @Summary 获取预售商品列表
export const getPresaleGoodList = (params) => {
  return service({
    url: '/presale/getPresaleGoodList',
    method: 'get',
    params
  })
}

// @Tags Presale
// @Summary 检查预售可用性
export const checkPresaleAvailable = (params) => {
  return service({
    url: '/presale/checkAvailable',
    method: 'get',
    params
  })
}

// @Tags Presale
// @Summary 获取预售参与者列表（管理端）
export const getPresaleParticipants = (params) => {
  return service({
    url: '/presale/getPresaleParticipants',
    method: 'get',
    params
  })
}
