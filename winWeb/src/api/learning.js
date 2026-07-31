import { request } from '@/utils/request.js'

const withIncludeI18n = (params = {}) => ({
  ...params,
  includeI18n: params?.includeI18n ?? 1
})

// ========== 打卡/积分 ==========
export const getCheckinStats = () => {
  return request({
    url: '/englishLearning/checkin/getStats',
    method: 'get'
  })
}

export const getLearningCheckinRecordList = (params = {}) => {
  return request({
    url: '/englishLearning/checkin/getCheckinRecordList',
    method: 'get',
    params
  })
}

export const doCheckin = () => {
  return request({
    url: '/englishLearning/checkin/do',
    method: 'post'
  })
}

export const exchangeTime = (points) => {
  return request({
    url: '/englishLearning/checkin/exchange',
    method: 'post',
    data: { points }
  })
}

export const getLearningPointRecordList = (params = {}) => {
  return request({
    url: '/englishLearning/checkin/getPointRecordList',
    method: 'get',
    params
  })
}

export const getLearningFreeTimeRecordList = (params = {}) => {
  return request({
    url: '/englishLearning/asset/getFreeTimeRecordList',
    method: 'get',
    params
  })
}

// ========== 内容分类 ==========
export const getCategoryList = (params = {}) => {
  return request({
    url: '/englishLearning/content/getCategoryList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getChapterList = (params = {}) => {
  return request({
    url: '/englishLearning/content/getChapterList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getVideoCategoryList = (params = {}) => {
  return request({
    url: '/englishLearning/content/getVideoCategoryList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

// ========== 视频系列/集数 ==========
export const getVideoSeriesList = (params = {}) => {
  return request({
    url: '/englishLearning/content/getVideoSeriesList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const findVideoSeries = (id) => {
  return request({
    url: '/englishLearning/content/findVideoSeries',
    method: 'get',
    params: withIncludeI18n({ ID: id })
  })
}

export const findVideoEpisode = (id) => {
  return request({
    url: '/englishLearning/content/findVideoEpisode',
    method: 'get',
    params: withIncludeI18n({ ID: id })
  })
}

export const getVideoEpisodeList = (params = {}) => {
  return request({
    url: '/englishLearning/content/getVideoEpisodeList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getVideoEpisodeListByTag = (params = {}) => {
  return request({
    url: '/englishLearning/content/getVideoEpisodeListByTag',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getVideoTagList = (params = {}) => {
  return request({
    url: '/videoTag/getVideoTagPublic',
    method: 'get',
    params
  })
}

// ========== 单词 ==========
export const getWordList = (params = {}) => {
  return request({
    url: '/englishLearning/word/getWordList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const findWord = (id) => {
  return request({
    url: '/englishLearning/word/findWord',
    method: 'get',
    params: withIncludeI18n({ ID: id })
  })
}

export const reportWordError = (data) => {
  return request({
    url: '/englishLearning/word/reportError',
    method: 'post',
    data
  })
}

export const getWordErrorLogList = (params = {}) => {
  return request({
    url: '/englishLearning/word/getErrorLogList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const deleteWordErrorLog = (wordId) => {
  return request({
    url: '/englishLearning/word/deleteErrorLog',
    method: 'delete',
    data: { wordId }
  })
}

// ========== 句子/字幕 ==========
export const getSentenceList = (episodeId) => {
  return request({
    url: '/englishLearning/video/getSentenceList',
    method: 'get',
    params: withIncludeI18n({ episodeId })
  })
}

// ========== 用户数据/进度 ==========
export const saveWordProgress = (data) => {
  return request({
    url: '/englishLearning/userData/saveWordProgress',
    method: 'post',
    data
  })
}

export const getWordProgress = () => {
  return request({
    url: '/englishLearning/userData/getWordProgress',
    method: 'get'
  })
}

export const getWatchProgress = (episodeId) => {
  return request({
    url: '/englishLearning/userData/getWatchProgress',
    method: 'get',
    params: { episodeId }
  })
}

export const getSeriesWatchProgressList = (seriesId) => {
  return request({
    url: '/englishLearning/userData/getSeriesWatchProgressList',
    method: 'get',
    params: { seriesId }
  })
}

// ========== 观看历史 ==========
export const getWatchHistoryList = (params = {}) => {
  return request({
    url: '/englishLearning/userData/getWatchHistoryList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

// ========== 收藏 ==========
export const collect = (targetType, targetId) => {
  return request({
    url: '/englishLearning/userData/collect',
    method: 'post',
    data: { targetType, targetId }
  })
}

export const uncollect = (targetType, targetId) => {
  return request({
    url: '/englishLearning/userData/uncollect',
    method: 'delete',
    data: { targetType, targetId }
  })
}

export const getCollectionList = (params = {}) => {
  return request({
    url: '/englishLearning/userData/getCollectionList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getCollectionDetailList = (params = {}) => {
  return request({
    url: '/englishLearning/userData/getCollectionDetailList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

// ========== 资产 ==========
export const getAsset = () => {
  return request({
    url: '/englishLearning/userData/getAsset',
    method: 'get'
  })
}

export const heartbeat = (data) => {
  return request({
    url: '/englishLearning/asset/heartbeat',
    method: 'post',
    data
  })
}

// ========== 日记相关 ==========
export const getDiaryCategoryList = (params = {}) => {
  return request({
    url: '/englishLearning/diary/getDiaryCategoryList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getDiaryTagPublic = (params = {}) => {
  return request({
    url: '/diaryTag/getDiaryTagPublic',
    method: 'get',
    params
  })
}

export const getDiaryList = (params = {}) => {
  return request({
    url: '/englishLearning/diary/getDiaryList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const getDiaryListByTag = (params = {}) => {
  return request({
    url: '/englishLearning/diary/getDiaryListByTag',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

export const findDiary = (id) => {
  return request({
    url: '/englishLearning/diary/findDiary',
    method: 'get',
    params: withIncludeI18n({ ID: id })
  })
}

export const getDiarySentenceList = (diaryId) => {
  return request({
    url: '/englishLearning/diary/getSentenceList',
    method: 'get',
    params: withIncludeI18n({ diaryId })
  })
}
