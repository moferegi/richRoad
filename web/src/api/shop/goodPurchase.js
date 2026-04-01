import service from '@/utils/request'

// @Tags GoodPurchase
// @Summary 创建进货记录
export const createGoodPurchase = (data) => {
  return service({
    url: '/goodPurchase/createGoodPurchase',
    method: 'post',
    data
  })
}

// @Tags GoodPurchase
// @Summary 删除进货记录
export const deleteGoodPurchase = (params) => {
  return service({
    url: '/goodPurchase/deleteGoodPurchase',
    method: 'delete',
    params
  })
}

// @Tags GoodPurchase
// @Summary 更新进货记录
export const updateGoodPurchase = (data) => {
  return service({
    url: '/goodPurchase/updateGoodPurchase',
    method: 'put',
    data
  })
}

// @Tags GoodPurchase
// @Summary 分页获取进货记录列表
export const getGoodPurchaseList = (params) => {
  return service({
    url: '/goodPurchase/getGoodPurchaseList',
    method: 'get',
    params
  })
}

// @Tags GoodPurchase
// @Summary 获取商品进货汇总
export const getGoodPurchaseSummary = (params) => {
  return service({
    url: '/goodPurchase/getGoodPurchaseSummary',
    method: 'get',
    params
  })
}
