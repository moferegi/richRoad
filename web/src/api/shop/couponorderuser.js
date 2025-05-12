import service from '@/utils/request'
// @Tags CouponOrderUser
// @Summary 创建优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.CouponOrderUser true "创建优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /cou/createCouponOrderUser [post]
export const createCouponOrderUser = (data) => {
  return service({
    url: '/cou/createCouponOrderUser',
    method: 'post',
    data
  })
}

// @Tags CouponOrderUser
// @Summary 删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.CouponOrderUser true "删除优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cou/deleteCouponOrderUser [delete]
export const deleteCouponOrderUser = (params) => {
  return service({
    url: '/cou/deleteCouponOrderUser',
    method: 'delete',
    params
  })
}

// @Tags CouponOrderUser
// @Summary 批量删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cou/deleteCouponOrderUser [delete]
export const deleteCouponOrderUserByIds = (params) => {
  return service({
    url: '/cou/deleteCouponOrderUserByIds',
    method: 'delete',
    params
  })
}

// @Tags CouponOrderUser
// @Summary 更新优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.CouponOrderUser true "更新优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /cou/updateCouponOrderUser [put]
export const updateCouponOrderUser = (data) => {
  return service({
    url: '/cou/updateCouponOrderUser',
    method: 'put',
    data
  })
}

// @Tags CouponOrderUser
// @Summary 用id查询优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.CouponOrderUser true "用id查询优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /cou/findCouponOrderUser [get]
export const findCouponOrderUser = (params) => {
  return service({
    url: '/cou/findCouponOrderUser',
    method: 'get',
    params
  })
}

// @Tags CouponOrderUser
// @Summary 分页获取优惠券列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取优惠券列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cou/getCouponOrderUserList [get]
export const getCouponOrderUserList = (params) => {
  return service({
    url: '/cou/getCouponOrderUserList',
    method: 'get',
    params
  })
}
// @Tags CouponOrderUser
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /cou/findCouponOrderUserDataSource [get]
export const getCouponOrderUserDataSource = () => {
  return service({
    url: '/cou/getCouponOrderUserDataSource',
    method: 'get',
  })
}

// @Tags CouponOrderUser
// @Summary 不需要鉴权的优惠券接口
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.CouponOrderUserSearch true "分页获取优惠券列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /cou/getCouponOrderUserPublic [get]
export const getCouponOrderUserPublic = () => {
  return service({
    url: '/cou/getCouponOrderUserPublic',
    method: 'get',
  })
}
