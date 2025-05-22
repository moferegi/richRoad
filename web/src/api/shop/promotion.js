import service from '@/utils/request'
// @Tags Promotion
// @Summary 创建促销信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Promotion true "创建促销信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /Promo/createPromotion [post]
export const createPromotion = (data) => {
  return service({
    url: '/Promo/createPromotion',
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
// @Router /Promo/deletePromotion [delete]
export const deletePromotion = (params) => {
  return service({
    url: '/Promo/deletePromotion',
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
// @Router /Promo/deletePromotion [delete]
export const deletePromotionByIds = (params) => {
  return service({
    url: '/Promo/deletePromotionByIds',
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
// @Router /Promo/updatePromotion [put]
export const updatePromotion = (data) => {
  return service({
    url: '/Promo/updatePromotion',
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
// @Router /Promo/findPromotion [get]
export const findPromotion = (params) => {
  return service({
    url: '/Promo/findPromotion',
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
// @Router /Promo/getPromotionList [get]
export const getPromotionList = (params) => {
  return service({
    url: '/Promo/getPromotionList',
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
// @Router /Promo/getPromotionPublic [get]
export const getPromotionPublic = () => {
  return service({
    url: '/Promo/getPromotionPublic',
    method: 'get',
  })
}
