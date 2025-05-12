import service from '@/utils/request'
// @Tags Coupon
// @Summary 创建优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Coupon true "创建优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /Cou/createCoupon [post]
export const createCoupon = (data) => {
  return service({
    url: '/Cou/createCoupon',
    method: 'post',
    data
  })
}

// @Tags Coupon
// @Summary 删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Coupon true "删除优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /Cou/deleteCoupon [delete]
export const deleteCoupon = (params) => {
  return service({
    url: '/Cou/deleteCoupon',
    method: 'delete',
    params
  })
}

// @Tags Coupon
// @Summary 批量删除优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /Cou/deleteCoupon [delete]
export const deleteCouponByIds = (params) => {
  return service({
    url: '/Cou/deleteCouponByIds',
    method: 'delete',
    params
  })
}

// @Tags Coupon
// @Summary 更新优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Coupon true "更新优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /Cou/updateCoupon [put]
export const updateCoupon = (data) => {
  return service({
    url: '/Cou/updateCoupon',
    method: 'put',
    data
  })
}

// @Tags Coupon
// @Summary 用id查询优惠券
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Coupon true "用id查询优惠券"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /Cou/findCoupon [get]
export const findCoupon = (params) => {
  return service({
    url: '/Cou/findCoupon',
    method: 'get',
    params
  })
}

// @Tags Coupon
// @Summary 分页获取优惠券列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取优惠券列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /Cou/getCouponList [get]
export const getCouponList = (params) => {
  return service({
    url: '/Cou/getCouponList',
    method: 'get',
    params
  })
}
// @Tags Coupon
// @Summary 获取数据源
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /Cou/findCouponDataSource [get]
export const getCouponDataSource = () => {
  return service({
    url: '/Cou/getCouponDataSource',
    method: 'get',
  })
}

// @Tags Coupon
// @Summary 不需要鉴权的优惠券接口
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.CouponSearch true "分页获取优惠券列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /Cou/getCouponPublic [get]
export const getCouponPublic = () => {
  return service({
    url: '/Cou/getCouponPublic',
    method: 'get',
  })
}
