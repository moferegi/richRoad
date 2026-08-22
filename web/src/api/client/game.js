import service from '@/utils/request'

// ==================== 游戏大分类 ====================

// @Tags GameAdmin
// @Summary 创建游戏大分类
export const createCategory = (data) => {
  return service({
    url: '/game/admin/createCategory',
    method: 'post',
    data
  })
}

// @Tags GameAdmin
// @Summary 更新游戏大分类
export const updateCategory = (data) => {
  return service({
    url: '/game/admin/updateCategory',
    method: 'put',
    data
  })
}

// @Tags GameAdmin
// @Summary 删除游戏大分类
export const deleteCategory = (params) => {
  return service({
    url: '/game/admin/deleteCategory',
    method: 'delete',
    params
  })
}

// @Tags GameAdmin
// @Summary 获取游戏大分类列表
export const getCategoryList = (params) => {
  return service({
    url: '/game/admin/getCategoryList',
    method: 'get',
    params
  })
}

// ==================== 难度分类 ====================

// @Tags GameAdmin
// @Summary 创建难度分类
export const createDifficultyCategory = (data) => {
  return service({
    url: '/game/admin/createDifficultyCategory',
    method: 'post',
    data
  })
}

// @Tags GameAdmin
// @Summary 更新难度分类
export const updateDifficultyCategory = (data) => {
  return service({
    url: '/game/admin/updateDifficultyCategory',
    method: 'put',
    data
  })
}

// @Tags GameAdmin
// @Summary 删除难度分类
export const deleteDifficultyCategory = (params) => {
  return service({
    url: '/game/admin/deleteDifficultyCategory',
    method: 'delete',
    params
  })
}

// @Tags GameAdmin
// @Summary 获取难度分类列表
export const getDifficultyCategoryList = (params) => {
  return service({
    url: '/game/admin/getDifficultyCategoryList',
    method: 'get',
    params
  })
}

// ==================== 关卡 ====================

// @Tags GameAdmin
// @Summary 创建关卡
export const createLevel = (data) => {
  return service({
    url: '/game/admin/createLevel',
    method: 'post',
    data
  })
}

// @Tags GameAdmin
// @Summary 更新关卡
export const updateLevel = (data) => {
  return service({
    url: '/game/admin/updateLevel',
    method: 'put',
    data
  })
}

// @Tags GameAdmin
// @Summary 删除关卡
export const deleteLevel = (params) => {
  return service({
    url: '/game/admin/deleteLevel',
    method: 'delete',
    params
  })
}

// @Tags GameAdmin
// @Summary 获取关卡列表
export const getLevelList = (params) => {
  return service({
    url: '/game/admin/getLevelList',
    method: 'get',
    params
  })
}

// ==================== 排行榜 ====================

// @Tags GameAdmin
// @Summary 获取排行榜
export const getLeaderboard = (params) => {
  return service({
    url: '/game/admin/getLeaderboard',
    method: 'get',
    params
  })
}

// @Tags GameAdmin
// @Summary 查看用户进度
export const getUserProgress = (params) => {
  return service({
    url: '/game/admin/getUserProgress',
    method: 'get',
    params
  })
}

// @Tags GameAdmin
// @Summary 设置用户进度
export const setUserProgress = (data) => {
  return service({
    url: '/game/admin/setUserProgress',
    method: 'post',
    data
  })
}

// ==================== 密码推理关卡 ====================

// @Tags GameAdmin
// @Summary 创建密码关卡
export const createPwdLevel = (data) => {
  return service({
    url: '/game/admin/createPwdLevel',
    method: 'post',
    data
  })
}

// @Tags GameAdmin
// @Summary 更新密码关卡
export const updatePwdLevel = (data) => {
  return service({
    url: '/game/admin/updatePwdLevel',
    method: 'put',
    data
  })
}

// @Tags GameAdmin
// @Summary 删除密码关卡
export const deletePwdLevel = (params) => {
  return service({
    url: '/game/admin/deletePwdLevel',
    method: 'delete',
    params
  })
}

// @Tags GameAdmin
// @Summary 获取密码关卡列表
export const getPwdLevelList = (params) => {
  return service({
    url: '/game/admin/getPwdLevelList',
    method: 'get',
    params
  })
}