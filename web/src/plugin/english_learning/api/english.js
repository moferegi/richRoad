import service from '@/utils/request'

// 创建英语单词
export const createEnglishWord = (data) => {
  return service({
    url: '/englishLearning/word/create',
    method: 'post',
    data
  })
}

// 更新英语单词
export const updateEnglishWord = (data) => {
  return service({
    url: '/englishLearning/word/update',
    method: 'put',
    data
  })
}

// 删除英语单词
export const deleteEnglishWord = (params) => {
  return service({
    url: '/englishLearning/word/delete',
    method: 'delete',
    params
  })
}

// 重生成单词发音
export const regenerateWordAudio = (data) => {
  return service({
    url: '/englishLearning/word/regenerateAudio',
    method: 'post',
    data
  })
}

// 预检TTS服务连通性
export const preflightWordTTS = (data) => {
  return service({
    url: '/englishLearning/word/preflightTTS',
    method: 'post',
    data
  })
}

// 获取单词详情
export const findEnglishWord = (params) => {
  return service({
    url: '/englishLearning/word/findWord',
    method: 'get',
    params
  })
}

// 分页查询单词
export const getEnglishWordList = (params) => {
  return service({
    url: '/englishLearning/word/getWordList',
    method: 'get',
    params
  })
}

// 分类管理
export const createCategory = (data) => {
  return service({
    url: '/englishLearning/content/createCategory',
    method: 'post',
    data
  })
}

export const updateCategory = (data) => {
  return service({
    url: '/englishLearning/content/updateCategory',
    method: 'put',
    data
  })
}

export const deleteCategory = (params) => {
  return service({
    url: '/englishLearning/content/deleteCategory',
    method: 'delete',
    params
  })
}

export const findCategory = (params) => {
  return service({
    url: '/englishLearning/content/findCategory',
    method: 'get',
    params
  })
}

export const getCategoryList = (params) => {
  return service({
    url: '/englishLearning/content/getCategoryList',
    method: 'get',
    params
  })
}

// 章节管理
export const createChapter = (data) => {
  return service({
    url: '/englishLearning/content/createChapter',
    method: 'post',
    data
  })
}

export const updateChapter = (data) => {
  return service({
    url: '/englishLearning/content/updateChapter',
    method: 'put',
    data
  })
}

export const deleteChapter = (params) => {
  return service({
    url: '/englishLearning/content/deleteChapter',
    method: 'delete',
    params
  })
}

export const getChapterList = (params) => {
  return service({
    url: '/englishLearning/content/getChapterList',
    method: 'get',
    params
  })
}

// 视频分类管理
export const createVideoCategory = (data) => {
  return service({
    url: '/englishLearning/content/createVideoCategory',
    method: 'post',
    data
  })
}

export const updateVideoCategory = (data) => {
  return service({
    url: '/englishLearning/content/updateVideoCategory',
    method: 'put',
    data
  })
}

export const deleteVideoCategory = (params) => {
  return service({
    url: '/englishLearning/content/deleteVideoCategory',
    method: 'delete',
    params
  })
}

export const getVideoCategoryList = (params) => {
  return service({
    url: '/englishLearning/content/getVideoCategoryList',
    method: 'get',
    params
  })
}

// 视频剧集管理
export const createVideoSeries = (data) => {
  return service({
    url: '/englishLearning/content/createVideoSeries',
    method: 'post',
    data
  })
}

export const updateVideoSeries = (data) => {
  return service({
    url: '/englishLearning/content/updateVideoSeries',
    method: 'put',
    data
  })
}

export const deleteVideoSeries = (params) => {
  return service({
    url: '/englishLearning/content/deleteVideoSeries',
    method: 'delete',
    params
  })
}

export const getVideoSeriesList = (params) => {
  return service({
    url: '/englishLearning/content/getVideoSeriesList',
    method: 'get',
    params
  })
}

// 视频单集管理
export const createVideoEpisode = (data) => {
  return service({
    url: '/englishLearning/content/createVideoEpisode',
    method: 'post',
    data
  })
}

export const updateVideoEpisode = (data) => {
  return service({
    url: '/englishLearning/content/updateVideoEpisode',
    method: 'put',
    data
  })
}

export const deleteVideoEpisode = (params) => {
  return service({
    url: '/englishLearning/content/deleteVideoEpisode',
    method: 'delete',
    params
  })
}

export const getVideoEpisodeList = (params) => {
  return service({
    url: '/englishLearning/content/getVideoEpisodeList',
    method: 'get',
    params
  })
}

// 提交字幕进行高亮解析并入库
export const parseSubtitle = (data) => {
  return service({
    url: '/englishLearning/video/parseSubtitle',
    method: 'post',
    data
  })
}

// 用户资源授权管理
export const grantEntitlement = (data) => {
  return service({
    url: '/englishLearning/asset/grantEntitlement',
    method: 'post',
    data
  })
}

export const revokeEntitlement = (data) => {
  return service({
    url: '/englishLearning/asset/revokeEntitlement',
    method: 'delete',
    data
  })
}

export const getEntitlementList = (params) => {
  return service({
    url: '/englishLearning/asset/getEntitlementList',
    method: 'get',
    params
  })
}
