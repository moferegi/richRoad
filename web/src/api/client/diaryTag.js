import service from '@/utils/request'

export const createDiaryTag = (data) => {
  return service({
    url: '/diaryTag/createDiaryTag',
    method: 'post',
    data
  })
}

export const deleteDiaryTag = (data) => {
  return service({
    url: '/diaryTag/deleteDiaryTag',
    method: 'delete',
    data
  })
}

export const deleteDiaryTagByIds = (data) => {
  return service({
    url: '/diaryTag/deleteDiaryTagByIds',
    method: 'delete',
    data
  })
}

export const updateDiaryTag = (data) => {
  return service({
    url: '/diaryTag/updateDiaryTag',
    method: 'put',
    data
  })
}

export const findDiaryTag = (params) => {
  return service({
    url: '/diaryTag/findDiaryTag',
    method: 'get',
    params
  })
}

export const getDiaryTagList = (params) => {
  return service({
    url: '/diaryTag/getDiaryTagList',
    method: 'get',
    params
  })
}

export const getDiaryTagPublic = (params) => {
  return service({
    url: '/diaryTag/getDiaryTagPublic',
    method: 'get',
    params
  })
}