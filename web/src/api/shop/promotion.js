import service from '@/utils/request'
// @Tags Promotion
// @Summary 创建促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Promotion true "创建促销信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /promo/createPromotion [post]
export const createPromotion = (data) => {
  return service({
    url: '/promo/createPromotion',
    method: 'post',
    data
  })
}

// @Tags Promotion
// @Summary 删除促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Promotion true "删除促销信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /promo/deletePromotion [delete]
export const deletePromotion = (params) => {
  return service({
    url: '/promo/deletePromotion',
    method: 'delete',
    params
  })
}

// @Tags Promotion
// @Summary 批量删除促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除促销信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /promo/deletePromotion [delete]
export const deletePromotionByIds = (params) => {
  return service({
    url: '/promo/deletePromotionByIds',
    method: 'delete',
    params
  })
}

// @Tags Promotion
// @Summary 更新促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Promotion true "更新促销信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /promo/updatePromotion [put]
export const updatePromotion = (data) => {
  return service({
    url: '/promo/updatePromotion',
    method: 'put',
    data
  })
}

// @Tags Promotion
// @Summary 用id查询促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Promotion true "用id查询促销信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /promo/findPromotion [get]
export const findPromotion = (params) => {
  return service({
    url: '/promo/findPromotion',
    method: 'get',
    params
  })
}

// @Tags Promotion
// @Summary 分页获取促销信息列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取促销信息列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /promo/getPromotionList [get]
export const getPromotionList = (params) => {
  return service({
    url: '/promo/getPromotionList',
    method: 'get',
    params
  })
}

// @Tags Promotion
// @Summary 不需要鉴权的促销信息接口
// @Accept application/json
// @Produce application/json
// @Param data query shopReq.PromotionSearch true "分页获取促销信息列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /promo/getPromotionPublic [get]
export const getPromotionPublic = () => {
  return service({
    url: '/promo/getPromotionPublic',
    method: 'get',
  })
}
