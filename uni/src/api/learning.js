import { request } from '@/utils/request.js'

const withIncludeI18n = (params = {}) => ({
  ...params,
  includeI18n: params?.includeI18n ?? 1,
})

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

export const getSentenceList = (episodeId) => {
  return request({
    url: '/englishLearning/video/getSentenceList',
    method: 'get',
    params: withIncludeI18n({ episodeId })
  })
}

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

export const getWatchHistoryList = (params = {}) => {
  return request({
    url: '/englishLearning/userData/getWatchHistoryList',
    method: 'get',
    params: withIncludeI18n(params)
  })
}

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
