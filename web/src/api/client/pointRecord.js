import service from '@/utils/request'
// @Tags PointRecord
// @Summary 创建积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PointRecord true "创建积分记录管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /cpr/createPointRecord [post]
export const createPointRecord = (data) => {
  return service({
    url: '/cpr/createPointRecord',
    method: 'post',
    data
  })
}

// @Tags PointRecord
// @Summary 删除积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PointRecord true "删除积分记录管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cpr/deletePointRecord [delete]
export const deletePointRecord = (params) => {
  return service({
    url: '/cpr/deletePointRecord',
    method: 'delete',
    params
  })
}

// @Tags PointRecord
// @Summary 批量删除积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除积分记录管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cpr/deletePointRecord [delete]
export const deletePointRecordByIds = (params) => {
  return service({
    url: '/cpr/deletePointRecordByIds',
    method: 'delete',
    params
  })
}

// @Tags PointRecord
// @Summary 更新积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.PointRecord true "更新积分记录管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cpr/updatePointRecord [put]
export const updatePointRecord = (data) => {
  return service({
    url: '/cpr/updatePointRecord',
    method: 'put',
    data
  })
}

// @Tags PointRecord
// @Summary 用id查询积分记录管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.PointRecord true "用id查询积分记录管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /cpr/findPointRecord [get]
export const findPointRecord = (params) => {
  return service({
    url: '/cpr/findPointRecord',
    method: 'get',
    params
  })
}

// @Tags PointRecord
// @Summary 分页获取积分记录管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取积分记录管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cpr/getPointRecordList [get]
export const getPointRecordList = (params) => {
  return service({
    url: '/cpr/getPointRecordList',
    method: 'get',
    params
  })
}

// @Tags PointRecord
// @Summary 获取试衣币统计（管理端）
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query clientReq.TryonPointStatsSearch true "查询试衣币统计"
// @Success 200 {string} string "{\"success\":true,\"data\":{},\"msg\":\"获取成功\"}"
// @Router /cpr/getPointRecordStats [get]
export const getPointRecordStats = (params) => {
  return service({
    url: '/cpr/getPointRecordStats',
    method: 'get',
    params
  })
}

// @Tags PointRecord
// @Summary 不需要鉴权的积分记录管理接口
// @Accept application/json
// @Produce application/json
// @Param data query clientReq.PointRecordSearch true "分页获取积分记录管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cpr/getPointRecordPublic [get]
export const getPointRecordPublic = () => {
  return service({
    url: '/cpr/getPointRecordPublic',
    method: 'get',
  })
}
