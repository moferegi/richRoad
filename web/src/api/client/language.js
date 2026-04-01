import service from '@/utils/request'

// @Tags Language
// @Summary 创建语言
export const createLanguage = (data) => {
  return service({
    url: '/language/createLanguage',
    method: 'post',
    data
  })
}

// @Tags Language
// @Summary 删除语言
export const deleteLanguage = (params) => {
  return service({
    url: '/language/deleteLanguage',
    method: 'delete',
    params
  })
}

// @Tags Language
// @Summary 更新语言
export const updateLanguage = (data) => {
  return service({
    url: '/language/updateLanguage',
    method: 'put',
    data
  })
}

// @Tags Language
// @Summary 获取语言列表（管理端）
export const getLanguageList = (params) => {
  return service({
    url: '/language/getLanguageList',
    method: 'get',
    params
  })
}

// @Tags Language
// @Summary 获取启用的语言列表（客户端）
export const getEnabledLanguages = () => {
  return service({
    url: '/language/getEnabledLanguages',
    method: 'get'
  })
}
