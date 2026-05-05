import { baseUrl } from '@/utils/request.js'
import { localText } from '@/utils/i18n.js'

const TRYON_DRAFT_KEY = 'tryon:draft:v1'
const TRYON_HISTORY_KEY = 'tryon:history:v1'
const TRYON_SELECTED_CLOTHES_KEY = 'tryon:selectedClothes:v1'
const TRYON_SELECTED_MODEL_KEY = 'tryon:selectedModel:v1'

export const createTryonRequestId = () => `uni_${Date.now()}_${Math.random().toString(16).slice(2, 10)}`

export const saveTryonDraft = (draft) => {
  uni.setStorageSync(TRYON_DRAFT_KEY, draft || {})
}

export const getTryonDraft = () => {
  return uni.getStorageSync(TRYON_DRAFT_KEY) || {}
}

export const clearTryonDraft = () => {
  uni.removeStorageSync(TRYON_DRAFT_KEY)
}

export const setSelectedClothes = (payload) => {
  uni.setStorageSync(TRYON_SELECTED_CLOTHES_KEY, payload || {})
}

export const getSelectedClothes = () => {
  return uni.getStorageSync(TRYON_SELECTED_CLOTHES_KEY) || {}
}

export const clearSelectedClothes = () => {
  uni.removeStorageSync(TRYON_SELECTED_CLOTHES_KEY)
}

export const setSelectedTryonModel = (payload) => {
  uni.setStorageSync(TRYON_SELECTED_MODEL_KEY, payload || {})
}

export const getSelectedTryonModel = () => {
  return uni.getStorageSync(TRYON_SELECTED_MODEL_KEY) || {}
}

export const clearSelectedTryonModel = () => {
  uni.removeStorageSync(TRYON_SELECTED_MODEL_KEY)
}

export const getTryonLocalHistory = () => {
  const data = uni.getStorageSync(TRYON_HISTORY_KEY)
  if (Array.isArray(data)) return data
  return []
}

export const appendTryonLocalHistory = (item) => {
  if (!item || typeof item !== 'object') return
  const oldList = getTryonLocalHistory()
  const taskNo = item.taskNo || item.requestID || ''
  const filtered = oldList.filter(v => (v.taskNo || v.requestID || '') !== taskNo)
  const merged = [{ ...item, savedAt: Date.now() }, ...filtered].slice(0, 100)
  uni.setStorageSync(TRYON_HISTORY_KEY, merged)
}

// Keep uni upload behavior aligned with backend default OSS strategy.
export const getTryonUploadUrl = () => `${baseUrl}/fileUploadAndDownload/upload`

export const uploadTryonImage = (tempFilePath) => {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: getTryonUploadUrl(),
      filePath: tempFilePath,
      name: 'file',
      header: {
        'x-token': uni.getStorageSync('x-token') || '',
      },
      success: (res) => {
        try {
          const data = JSON.parse(res.data)
          if (data.code === 0 && data.data && data.data.file && data.data.file.url) {
            resolve(data.data.file.url)
            return
          }
          reject(new Error(data.msg || '上传失败'))
        } catch (e) {
          reject(new Error('上传失败'))
        }
      },
      fail: (err) => reject(new Error(err.errMsg || '上传失败')),
    })
  })
}

const toLocalizedText = (value, lang) => {
  if (value === null || value === undefined) return ''
  return localText(value, lang)
}

const isModelEnabled = (item) => {
  if (!item || typeof item !== 'object') return false
  if (item.enabled === undefined || item.enabled === null) return true
  if (typeof item.enabled === 'boolean') return item.enabled
  const text = String(item.enabled).trim().toLowerCase()
  return text === '1' || text === 'true' || text === 'yes' || text === 'on'
}

const matchSceneType = (item, sceneType) => {
  if (!sceneType) return true
  const scene = String(sceneType).trim().toLowerCase()

  if (Array.isArray(item?.scenes) && item.scenes.length > 0) {
    return item.scenes.some(s => String(s).trim().toLowerCase() === scene)
  }

  const sceneTypeVal = item?.sceneType || item?.scene || item?.roomType
  if (!sceneTypeVal) return true
  return String(sceneTypeVal).trim().toLowerCase() === scene
}

export const parseTryonModels = (modelsRaw, sceneType, fallbackCost = 1, lang = '') => {
  let list = []
  if (Array.isArray(modelsRaw)) {
    list = modelsRaw
  } else if (typeof modelsRaw === 'string' && modelsRaw.trim()) {
    try {
      const parsed = JSON.parse(modelsRaw)
      if (Array.isArray(parsed)) {
        list = parsed
      }
    } catch (e) {
      list = []
    }
  }

  const normalized = list
    .filter(item => item && typeof item === 'object')
    .filter(item => isModelEnabled(item) && matchSceneType(item, sceneType))
    .map((item, index) => {
      const key = String(item.key || item.modelKey || `model_${index + 1}`)

      const rawName = item.nameI18n || item.name || item.titleI18n || item.title || key
      const rawDesc = item.descI18n || item.desc || item.descriptionI18n || item.description || ''

      return {
        key,
        name: toLocalizedText(rawName, lang) || key,
        cost: Number(item.cost || fallbackCost || 1),
        desc: rawDesc,
        descText: toLocalizedText(rawDesc, lang),
        provider: item.provider || '',
        mode: item.mode || '',
        url: item.url || item.providerUrl || '',
        token: item.token || item.providerToken || '',
      }
    })

  if (normalized.length > 0) {
    return normalized
  }

  return [{
    key: 'aitryon',
    name: 'aitryon',
    cost: Number(fallbackCost || 1),
    desc: {
      zh: '默认AI试衣模型，适合常规试衣场景',
      en: 'Default AI try-on model for common try-on scenes',
    },
  }]
}
