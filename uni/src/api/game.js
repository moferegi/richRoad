import { request } from '@/utils/request.js'

export const getGameList = () => {
  return request({
    url: '/game/getGameList',
    method: 'GET'
  })
}

export const getGameCategory = (gameID) => {
  return request({
    url: '/game/getGameCategory',
    method: 'GET',
    params: { gameID }
  })
}

export const getDifficultyCategories = (gameID) => {
  return request({
    url: '/game/getDifficultyCategories',
    method: 'GET',
    params: { gameID }
  })
}

export const getDifficultyCategory = (catID) => {
  return request({
    url: '/game/getDifficultyCategory',
    method: 'GET',
    params: { catID }
  })
}

export const getLevelDetail = (levelID) => {
  return request({
    url: '/game/getLevelDetail',
    method: 'GET',
    params: { levelID }
  })
}

export const getLevelList = (categoryID) => {
  return request({
    url: '/game/getLevelList',
    method: 'GET',
    params: { categoryID }
  })
}

export const submitLevelResult = (levelID) => {
  return request({
    url: '/game/submitLevelResult',
    method: 'POST',
    params: { levelID }
  })
}

export const getUserProgress = (categoryID) => {
  return request({
    url: '/game/getUserProgress',
    method: 'GET',
    params: { categoryID }
  })
}

export const getLeaderboard = (type = 'all', limit = 50) => {
  return request({
    url: '/game/getLeaderboard',
    method: 'GET',
    params: { type, limit }
  })
}