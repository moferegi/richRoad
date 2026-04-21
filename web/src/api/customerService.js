import service from '@/utils/request'

// ===== 系统配置接口 =====

/**
 * 获取客服系统配置（公开）
 */
export const getCsConfig = () => {
  return service({ url: '/cs/config/get', method: 'get' })
}

/**
 * 更新客服系统配置（管理员）
 * @param {Object} data - { platEnabled: boolean }
 */
export const updateCsConfig = (data) => {
  return service({ url: '/cs/admin/config/update', method: 'put', data })
}

// ===== 会话接口 =====

/**
 * 发起或获取当前活跃会话（客户端用户调用）
 */
export const getOrCreateConversation = () => {
  return service({ url: '/cs/conversation/getOrCreate', method: 'post' })
}

/**
 * 评价会话
 * @param {Object} data - { conversationId, rating }
 */
export const rateConversation = (data) => {
  return service({ url: '/cs/conversation/rate', method: 'post', data })
}

/**
 * 获取所有会话（管理员）
 * @param {Object} params - { page, pageSize, status, agentUserID, clientUserID }
 */
export const getConversationList = (params) => {
  return service({ url: '/cs/admin/conversation/list', method: 'get', params })
}

/**
 * 获取坐席自己的会话列表
 * @param {Object} params - { page, pageSize, status }
 */
export const getAgentConversationList = (params) => {
  return service({ url: '/cs/agent/conversation/list', method: 'get', params })
}

/**
 * 关闭会话（坐席）
 * @param {Object} data - { conversationId }
 */
export const closeConversation = (data) => {
  return service({ url: '/cs/agent/conversation/close', method: 'post', data })
}

/**
 * 转接会话
 * @param {Object} data - { conversationId, targetAgentUserID }
 */
export const transferConversation = (data) => {
  return service({ url: '/cs/agent/conversation/transfer', method: 'post', data })
}

/**
 * 手动分配会话（管理员）
 * @param {Object} data - { conversationId, agentUserID }
 */
export const assignConversation = (data) => {
  return service({ url: '/cs/admin/conversation/assign', method: 'post', data })
}

// ===== 消息接口 =====

/**
 * 获取消息历史
 * @param {Object} params - { conversationId, page, pageSize }
 */
export const getMessageHistory = (params) => {
  return service({ url: '/cs/message/history', method: 'get', params })
}

/**
 * 发送消息（坐席 REST 接口，WebSocket 优先）
 * @param {Object} data - { conversationId, msgType, content }
 */
export const sendMessage = (data) => {
  return service({ url: '/cs/agent/message/send', method: 'post', data })
}

/**
 * 撤回消息
 * @param {Object} data - { messageId }
 */
export const revokeMessage = (data) => {
  return service({ url: '/cs/message/revoke', method: 'post', data })
}

// ===== 坐席接口 =====

/**
 * 获取坐席列表（管理员）
 */
export const getAgentList = () => {
  return service({ url: '/cs/admin/agent/list', method: 'get' })
}

/**
 * 创建坐席
 * @param {Object} data - { userID, maxSessions, isEnabled, nickname }
 */
export const createAgent = (data) => {
  return service({ url: '/cs/admin/agent/create', method: 'post', data })
}

/**
 * 更新坐席信息
 * @param {Object} data - { ID, maxSessions, isEnabled, nickname }
 */
export const updateAgent = (data) => {
  return service({ url: '/cs/admin/agent/update', method: 'put', data })
}

/**
 * 删除坐席
 * @param {Object} params - { ID }
 */
export const deleteAgent = (params) => {
  return service({ url: '/cs/admin/agent/delete', method: 'delete', params })
}

// ===== 快捷回复接口 =====

/**
 * 获取快捷回复列表（分页+搜索，管理员/坐席）
 * @param {Object} params - { page, pageSize, title }
 */
export const getQuickReplyList = (params) => {
  return service({ url: '/cs/agent/quickReply/list', method: 'get', params })
}

/**
 * 获取全部快捷回复（供工作台使用）
 */
export const getAllQuickReplies = () => {
  return service({ url: '/cs/quickReply/all', method: 'get' })
}

/**
 * 创建快捷回复
 * @param {Object} data - { title, content, sort }
 */
export const createQuickReply = (data) => {
  return service({ url: '/cs/admin/quickReply/create', method: 'post', data })
}

/**
 * 更新快捷回复
 * @param {Object} data - { ID, title, content, sort }
 */
export const updateQuickReply = (data) => {
  return service({ url: '/cs/admin/quickReply/update', method: 'put', data })
}

/**
 * 删除快捷回复
 * @param {Object} params - { ID }
 */
export const deleteQuickReply = (params) => {
  return service({ url: '/cs/admin/quickReply/delete', method: 'delete', params })
}

// ===== 黑名单接口 =====

/**
 * 获取黑名单列表
 * @param {Object} params - { page, pageSize }
 */
export const getBlacklist = (params) => {
  return service({ url: '/cs/admin/blacklist/list', method: 'get', params })
}

/**
 * 加入黑名单
 * @param {Object} data - { clientUserID, reason }
 */
export const addToBlacklist = (data) => {
  return service({ url: '/cs/admin/blacklist/add', method: 'post', data })
}

/**
 * 移出黑名单
 * @param {Object} params - { ID }
 */
export const removeFromBlacklist = (params) => {
  return service({ url: '/cs/admin/blacklist/remove', method: 'delete', params })
}
