import { baseUrl } from '@/utils/request.js'
import { localText, resolveApiMessage, t } from '@/utils/i18n.js'

const TRYON_DRAFT_KEY = 'tryon:draft:v1'
const TRYON_HISTORY_KEY = 'tryon:history:v1'
const TRYON_SELECTED_CLOTHES_KEY = 'tryon:selectedClothes:v1'
const TRYON_SELECTED_MODEL_KEY = 'tryon:selectedModel:v1'

const getUploadFailMessage = () => t('uploadFail')

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

export const setTryonLocalHistory = (list) => {
  if (!Array.isArray(list)) {
    uni.setStorageSync(TRYON_HISTORY_KEY, [])
    return
  }
  uni.setStorageSync(TRYON_HISTORY_KEY, list)
}

export const appendTryonLocalHistory = (item) => {
  if (!item || typeof item !== 'object') return
  const oldList = getTryonLocalHistory()
  const taskNo = item.taskNo || item.requestID || ''
  const filtered = oldList.filter(v => (v.taskNo || v.requestID || '') !== taskNo)
  const merged = [{ ...item, savedAt: Date.now() }, ...filtered].slice(0, 100)
  uni.setStorageSync(TRYON_HISTORY_KEY, merged)
}

export const removeTryonLocalHistoryByTaskKey = (taskKey) => {
  const key = String(taskKey || '').trim()
  if (!key) return
  const oldList = getTryonLocalHistory()
  const filtered = oldList.filter(v => String(v.taskNo || v.requestID || '').trim() !== key)
  setTryonLocalHistory(filtered)
}

// Keep uni upload behavior aligned with backend default OSS strategy.
export const getTryonUploadUrl = () => `${baseUrl}/fileUploadAndDownload/upload`

const normalizeUploadFolder = (folder) => {
  const raw = String(folder || '').trim().replace(/\\/g, '/')
  if (!raw) return ''

  const parts = raw
    .split('/')
    .map(v => String(v || '').trim())
    .filter(v => v && v !== '.' && v !== '..')

  return parts.join('/')
}

export const getTryonUploadFolder = ({ sceneType = 'clothes', operationType = 'tryon', role = 'source' } = {}) => {
  const operation = String(operationType || 'tryon').trim().toLowerCase()
  let scene = String(sceneType || 'clothes').trim().toLowerCase()
  let pathRole = String(role || 'source').trim().toLowerCase()

  if (operation === 'takeoff') {
    scene = 'takeoff'
  }

  if (!['clothes', 'shoes', 'takeoff'].includes(scene)) {
    scene = 'clothes'
  }

  if (!pathRole) {
    pathRole = 'source'
  }

  return normalizeUploadFolder(`tryon/${scene}/${pathRole}`)
}

export const uploadTryonImage = (tempFilePath, folder = '', uploadType = '', uploadPosition = 'tryon') => {
  return new Promise((resolve, reject) => {
    const formData = {}
    const uploadFolder = normalizeUploadFolder(folder)
    if (uploadFolder) {
      formData.folder = uploadFolder
    }
    const normalizedUploadType = String(uploadType || '').trim()
    if (normalizedUploadType) {
      formData.uploadType = normalizedUploadType
    }
    const normalizedUploadPosition = String(uploadPosition || '').trim().toLowerCase()
    if (normalizedUploadPosition) {
      formData.uploadPosition = normalizedUploadPosition
    }

    uni.uploadFile({
      url: getTryonUploadUrl(),
      filePath: tempFilePath,
      name: 'file',
      formData,
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
          reject(new Error(resolveApiMessage(data.msg, 'uploadFail')))
        } catch (e) {
          reject(new Error(getUploadFailMessage()))
        }
      },
      fail: (err) => reject(new Error(resolveApiMessage(err?.errMsg, 'uploadFail'))),
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

const toBool = (value, fallback = false) => {
  if (typeof value === 'boolean') return value
  if (value === undefined || value === null || value === '') return fallback
  const text = String(value).trim().toLowerCase()
  if (['1', 'true', 'yes', 'on'].includes(text)) return true
  if (['0', 'false', 'no', 'off'].includes(text)) return false
  return fallback
}

const normalizeModelUsage = (value) => {
  const usage = String(value || '').trim().toLowerCase()
  return usage === 'beautify' ? 'beautify' : 'tryon'
}

const isBeautifyModelUsage = (item = {}) => normalizeModelUsage(item?.modelUsage) === 'beautify'

const isTryonModelUsage = (item = {}) => !isBeautifyModelUsage(item)

const isAliyunTryonModel = (item = {}) => {
  if (!isTryonModelUsage(item)) return false
  const provider = String(item.provider || '').trim().toLowerCase()
  const key = String(item.key || item.modelKey || '').trim().toLowerCase()
  const model = String(item.model || '').trim().toLowerCase()
  if (provider.includes('aliyun') || provider.includes('dashscope')) return true
  if (key.includes('aliyun') || key.includes('aitryon')) return true
  return model.startsWith('aitryon')
}

const inferSupportsRefiner = (item = {}) => {
  if (!isAliyunTryonModel(item)) return false
  const model = String(item.model || '').trim().toLowerCase()
  const key = String(item.key || item.modelKey || '').trim().toLowerCase()
  return model === 'aitryon' || model === 'aitryon-plus' || key.includes('aliyun_aitryon') || key.includes('aliyun_aitryon_plus')
}

const inferSupportsBeautify = (item = {}) => {
  if (isBeautifyModelUsage(item)) return true
  if (!isAliyunTryonModel(item)) return false
  if (item.beautifyModel || item.beautifyUrl || item.beautifyToken || item.beautifyAccessKeyId) {
    return true
  }
  return false
}

const normalizeRefinerGender = (value, fallback = 'woman') => {
  const gender = String(value || '').trim().toLowerCase()
  if (gender === 'woman' || gender === 'man') return gender
  return fallback
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

const parseTryonModelsRaw = (modelsRaw) => {
  if (Array.isArray(modelsRaw)) {
    return modelsRaw
  }
  if (typeof modelsRaw === 'string' && modelsRaw.trim()) {
    try {
      const parsed = JSON.parse(modelsRaw)
      if (Array.isArray(parsed)) {
        return parsed
      }
    } catch (e) {
      return []
    }
  }
  return []
}

const normalizeTryonModelItem = (item = {}, index = 0, fallbackCost = 1, lang = '') => {
  const key = String(item.key || item.modelKey || `model_${index + 1}`)
  const rawName = item.nameI18n || item.name || item.titleI18n || item.title || key
  const rawDesc = item.descI18n || item.desc || item.descriptionI18n || item.description || ''
  const rawRefinerDesc = item.refinerDescI18n || item.refinerDesc || ''
  const rawBeautifyDesc = item.beautifyDescI18n || item.beautifyDesc || ''
  const supportsRefiner = toBool(item.supportsRefiner, inferSupportsRefiner(item))
  const supportsBeautify = toBool(item.supportsBeautify, inferSupportsBeautify(item))
  const refinerExtraCost = Math.max(0, Number(item.refinerExtraCost || item.refinerExtraPoints || 0))
  const beautifyExtraCost = Math.max(0, Number(item.beautifyExtraCost || item.beautifyExtraPoints || 0))

  return {
    key,
    modelUsage: 'tryon',
    beautifyModelKey: String(item.beautifyModelKey || ''),
    name: toLocalizedText(rawName, lang) || key,
    cost: Number(item.cost || fallbackCost || 1),
    desc: rawDesc,
    descText: toLocalizedText(rawDesc, lang),
    provider: item.provider || '',
    mode: item.mode || '',
    url: item.url || item.providerUrl || '',
    token: item.token || item.providerToken || '',
    supportsRefiner,
    refinerExtraCost,
    refinerModel: String(item.refinerModel || 'aitryon-refiner'),
    refinerGender: normalizeRefinerGender(item.refinerGender, 'woman'),
    refinerDesc: rawRefinerDesc,
    refinerDescText: toLocalizedText(rawRefinerDesc, lang),
    supportsBeautify,
    beautifyModel: String(item.beautifyModel || 'RetouchSkin'),
    beautifyExtraCost,
    beautifyRetouchDegree: Number(item.beautifyRetouchDegree || 70),
    beautifyWhiteningDegree: Number(item.beautifyWhiteningDegree || 30),
    beautifyDesc: rawBeautifyDesc,
    beautifyDescText: toLocalizedText(rawBeautifyDesc, lang),
    freeQuotaTotal: Math.max(0, Number(item.freeQuotaTotal || 0)),
  }
}

export const parseTryonModels = (modelsRaw, sceneType, fallbackCost = 1, lang = '') => {
  const list = parseTryonModelsRaw(modelsRaw)

  const normalized = list
    .filter(item => item && typeof item === 'object')
    .filter(item => isTryonModelUsage(item))
    .filter(item => isModelEnabled(item) && matchSceneType(item, sceneType))
    .map((item, index) => normalizeTryonModelItem(item, index, fallbackCost, lang))

  if (normalized.length > 0) {
    return normalized
  }

  return [{
    key: 'aitryon',
    modelUsage: 'tryon',
    beautifyModelKey: '',
    name: 'aitryon',
    cost: Number(fallbackCost || 1),
    supportsRefiner: true,
    refinerExtraCost: 1,
    refinerModel: 'aitryon-refiner',
    refinerGender: 'woman',
    supportsBeautify: false,
    beautifyModel: 'RetouchSkin',
    beautifyExtraCost: 0,
    beautifyRetouchDegree: 70,
    beautifyWhiteningDegree: 30,
    desc: {
      en: 'Default AI try-on model for common try-on scenes',
      zh: 'Default AI try-on model for common try-on scenes',
      'zh-TW': 'Default AI try-on model for common try-on scenes',
      mn: 'Default AI try-on model for common try-on scenes',
      th: 'Default AI try-on model for common try-on scenes',
      hi: 'Default AI try-on model for common try-on scenes',
      id: 'Default AI try-on model for common try-on scenes',
    },
  }]
}

const normalizeBeautifyModelItem = (item = {}, index = 0, fallbackCost = 0, lang = '') => {
  const key = String(item.key || item.modelKey || `beautify_${index + 1}`)
  const rawName = item.nameI18n || item.name || item.titleI18n || item.title || key
  const rawDesc = item.beautifyDescI18n || item.beautifyDesc || item.descI18n || item.desc || ''
  const beautifyModel = String(item.beautifyModel || item.model || 'RetouchSkin')
  const beautifyExtraCost = Math.max(0, Number(item.beautifyExtraCost || item.beautifyExtraPoints || item.cost || fallbackCost || 0))

  return {
    key,
    modelUsage: 'beautify',
    name: toLocalizedText(rawName, lang) || key,
    supportsBeautify: true,
    beautifyModelKey: key,
    beautifyModel,
    beautifyExtraCost,
    beautifyRetouchDegree: Number(item.beautifyRetouchDegree || 70),
    beautifyWhiteningDegree: Number(item.beautifyWhiteningDegree || 30),
    beautifyDesc: rawDesc,
    beautifyDescText: toLocalizedText(rawDesc, lang),
    provider: item.provider || '',
    mode: item.mode || '',
    url: item.beautifyUrl || item.url || item.providerUrl || '',
    token: item.beautifyToken || item.token || item.providerToken || '',
  }
}

export const parseTryonBeautifyModels = (modelsRaw, sceneType, fallbackCost = 0, lang = '') => {
  const list = parseTryonModelsRaw(modelsRaw)
  const normalizedList = list
    .filter(item => item && typeof item === 'object')
    .filter(item => isModelEnabled(item) && matchSceneType(item, sceneType))

  return normalizedList
    .filter(item => isBeautifyModelUsage(item))
    .map((item, index) => normalizeBeautifyModelItem(item, index, fallbackCost, lang))
}
