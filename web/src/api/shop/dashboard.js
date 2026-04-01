import service from '@/utils/request'

// @Tags Dashboard
// @Summary 获取数据看板概览
export const getDashboardOverview = (params) => {
  return service({
    url: '/dashboard/getOverview',
    method: 'get',
    params
  })
}
