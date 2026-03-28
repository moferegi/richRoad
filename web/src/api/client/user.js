import service from '@/utils/request'

// @Tags ClientUser
// @Summary 创建客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ClientUser true "创建客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /clientUser/createClientUser [post]
export const createClientUser = (data) => {
  return service({
    url: '/clientUser/createClientUser',
    method: 'post',
    data
  })
}

// @Tags ClientUser
// @Summary 删除客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ClientUser true "删除客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /clientUser/deleteClientUser [delete]
export const deleteClientUser = (params) => {
  return service({
    url: '/clientUser/deleteClientUser',
    method: 'delete',
    params
  })
}

// @Tags ClientUser
// @Summary 批量删除客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /clientUser/deleteClientUser [delete]
export const deleteClientUserByIds = (params) => {
  return service({
    url: '/clientUser/deleteClientUserByIds',
    method: 'delete',
    params
  })
}

// @Tags ClientUser
// @Summary 更新客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ClientUser true "更新客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /clientUser/updateClientUser [put]
export const updateClientUser = (data) => {
  return service({
    url: '/clientUser/updateClientUser',
    method: 'put',
    data
  })
}

// @Tags ClientUser
// @Summary 用id查询客户端用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ClientUser true "用id查询客户端用户"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /clientUser/findClientUser [get]
export const findClientUser = (params) => {
  return service({
    url: '/clientUser/findClientUser',
    method: 'get',
    params
  })
}

// @Tags ClientUser
// @Summary 分页获取客户端用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取客户端用户列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /clientUser/getClientUserList [get]
export const getClientUserList = (params) => {
  return service({
    url: '/clientUser/getClientUserList',
    method: 'get',
    params
  })
}

/**
 * 获取用户的下级列表
 * @param {Object} params { userID, page, pageSize }
 * @returns {Promise}
 */
export const getSubordinates = (params) => {
  return service({
    url: '/clientUser/getSubordinates',
    method: 'get',
    params
  })
}
