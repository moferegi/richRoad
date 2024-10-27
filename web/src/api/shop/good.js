import service from '@/utils/request'

// @Tags Good
// @Summary 创建商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Good true "创建商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /good/createGood [post]
export const createGood = (data) => {
  return service({
    url: '/good/createGood',
    method: 'post',
    data
  })
}

// @Tags Good
// @Summary 删除商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Good true "删除商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /good/deleteGood [delete]
export const deleteGood = (params) => {
  return service({
    url: '/good/deleteGood',
    method: 'delete',
    params
  })
}

// @Tags Good
// @Summary 批量删除商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /good/deleteGood [delete]
export const deleteGoodByIds = (params) => {
  return service({
    url: '/good/deleteGoodByIds',
    method: 'delete',
    params
  })
}

// @Tags Good
// @Summary 更新商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Good true "更新商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /good/updateGood [put]
export const updateGood = (data) => {
  return service({
    url: '/good/updateGood',
    method: 'put',
    data
  })
}

// @Tags Good
// @Summary 用id查询商品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Good true "用id查询商品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /good/findGood [get]
export const findGood = (params) => {
  return service({
    url: '/good/findGood',
    method: 'get',
    params
  })
}

// @Tags Good
// @Summary 分页获取商品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取商品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /good/getGoodList [get]
export const getGoodList = (params) => {
  return service({
    url: '/good/getGoodList',
    method: 'get',
    params
  })
}
