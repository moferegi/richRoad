import { defineStore } from 'pinia'
import { ref } from 'vue'
import { localText } from '@/utils/i18n.js'
import { getCurrencySymbol, getAppName, getAppLogo } from '@/api/sysConfig.js'

const resolveConfigValue = (data) => {
  return typeof data === 'object' ? data?.configValue : data
}

const getActiveLocale = () => uni.getStorageSync('app-lang') || 'mn'

const normalizeString = (value) => {
  if (value === undefined || value === null) return ''
  return String(value)
}

export const useAppConfigStore = defineStore('appConfig', () => {
  const currencySymbol = ref(uni.getStorageSync('currency_symbol') || '¥')
  const appNameRaw = ref(uni.getStorageSync('app_name_raw') || uni.getStorageSync('app_name') || 'RichRoad')
  const appName = ref(uni.getStorageSync('app_name') || 'RichRoad')
  const appLogo = ref(uni.getStorageSync('app_logo') || '')
  const loaded = ref(false)

  const refreshLocalizedConfig = () => {
    const localized = localText(appNameRaw.value, getActiveLocale()) || normalizeString(appNameRaw.value) || 'RichRoad'
    appName.value = localized
    uni.setStorageSync('app_name', localized)
  }

  const loadConfig = async () => {
    if (loaded.value) {
      refreshLocalizedConfig()
      return
    }

    try {
      const [symRes, nameRes, logoRes] = await Promise.all([
        getCurrencySymbol(),
        getAppName(),
        getAppLogo(),
      ])

      if (symRes.code === 0 && symRes.data) {
        const sym = resolveConfigValue(symRes.data)
        if (sym) { currencySymbol.value = sym; uni.setStorageSync('currency_symbol', sym) }
      }

      if (nameRes.code === 0 && nameRes.data) {
        const rawNameValue = resolveConfigValue(nameRes.data)
        const rawText = normalizeString(rawNameValue).trim()
        if (rawText) {
          appNameRaw.value = rawText
          uni.setStorageSync('app_name_raw', rawText)
        }
      }

      if (logoRes.code === 0 && logoRes.data) {
        const logo = resolveConfigValue(logoRes.data)
        if (logo) { appLogo.value = logo; uni.setStorageSync('app_logo', logo) }
      }

      refreshLocalizedConfig()
      loaded.value = true
    } catch (e) {
      console.error('加载应用配置失败', e)
    }
  }

  /**
   * 格式化金额（分→元，带货币符号）
   * @param {number} cents 金额（分）
   * @returns {string} 格式化后的金额字符串
   */
  const formatPrice = (cents) => {
    if (cents === null || cents === undefined) return currencySymbol.value + '0.00'
    return currencySymbol.value + (Number(cents) / 100).toFixed(2)
  }

  return {
    currencySymbol,
    appName,
    appNameRaw,
    appLogo,
    loaded,
    loadConfig,
    refreshLocalizedConfig,
    formatPrice,
  }
})
