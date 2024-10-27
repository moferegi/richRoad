import service from '@/utils/request'

// @Tags Address
// @Summary 创建用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Address true "创建用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /address/createAddress [post]
export const createAddress = (data) => {
  return service({
    url: '/address/createAddress',
    method: 'post',
    data
  })
}

// @Tags Address
// @Summary 删除用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Address true "删除用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /address/deleteAddress [delete]
export const deleteAddress = (params) => {
  return service({
    url: '/address/deleteAddress',
    method: 'delete',
    params
  })
}

// @Tags Address
// @Summary 批量删除用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /address/deleteAddress [delete]
export const deleteAddressByIds = (params) => {
  return service({
    url: '/address/deleteAddressByIds',
    method: 'delete',
    params
  })
}

// @Tags Address
// @Summary 更新用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Address true "更新用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /address/updateAddress [put]
export const updateAddress = (data) => {
  return service({
    url: '/address/updateAddress',
    method: 'put',
    data
  })
}

// @Tags Address
// @Summary 用id查询用户地址
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Address true "用id查询用户地址"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /address/findAddress [get]
export const findAddress = (params) => {
  return service({
    url: '/address/findAddress',
    method: 'get',
    params
  })
}

// @Tags Address
// @Summary 分页获取用户地址列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户地址列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /address/getAddressList [get]
export const getAddressList = (params) => {
  return service({
    url: '/address/getAddressList',
    method: 'get',
    params
  })
}
// @Tags Address
// @Summary 获取数据源
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /address/findAddressDataSource [get]
export const getAddressDataSource = () => {
  return service({
    url: '/address/getAddressDataSource',
    method: 'get',
  })
}
