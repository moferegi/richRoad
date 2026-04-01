import service from '@/utils/request'

// @Tags SignIn
// @Summary 获取签到记录列表（管理端）
export const getSignInList = (params) => {
  return service({
    url: '/signIn/getSignInList',
    method: 'get',
    params
  })
}

// @Tags SignIn
// @Summary 删除签到记录
export const deleteSignIn = (params) => {
  return service({
    url: '/signIn/deleteSignIn',
    method: 'delete',
    params
  })
}
