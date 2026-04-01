import service from '@/utils/request'

// @Tags MarketingReward
// @Summary 创建营销奖励规则
export const createMarketingReward = (data) => {
  return service({
    url: '/marketingReward/createMarketingReward',
    method: 'post',
    data
  })
}

// @Tags MarketingReward
// @Summary 删除营销奖励规则
export const deleteMarketingReward = (params) => {
  return service({
    url: '/marketingReward/deleteMarketingReward',
    method: 'delete',
    params
  })
}

// @Tags MarketingReward
// @Summary 更新营销奖励规则
export const updateMarketingReward = (data) => {
  return service({
    url: '/marketingReward/updateMarketingReward',
    method: 'put',
    data
  })
}

// @Tags MarketingReward
// @Summary 获取营销奖励规则列表
export const getMarketingRewardList = (params) => {
  return service({
    url: '/marketingReward/getMarketingRewardList',
    method: 'get',
    params
  })
}

// @Tags MarketingReward
// @Summary 根据类型获取奖励规则
export const getMarketingRewardByType = (params) => {
  return service({
    url: '/marketingReward/getMarketingRewardByType',
    method: 'get',
    params
  })
}
