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

/**
 * 后台调整客户端用户试衣币
 * @param {Object} data 调整参数
 * @param {number} data.userID 用户ID
 * @param {'increase'|'decrease'} data.changeType 调整类型
 * @param {number} data.amount 调整数量
 * @param {string} data.reason 调整原因
 * @param {string} data.remark 备注
 * @returns {Promise} 调整结果
 */
export const adjustTryonPoint = (data) => {
  return service({
    url: '/clientUser/adjustTryonPoint',
    method: 'post',
    data
  })
}
