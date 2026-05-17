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

const normalizeRawConfig = (value) => {
  if (value === undefined || value === null) return ''
  if (typeof value === 'string') return value
  if (typeof value === 'object') {
    try {
      return JSON.stringify(value)
    } catch {
      return ''
    }
  }
  return String(value)
}

export const useAppConfigStore = defineStore('appConfig', () => {
  const currencySymbolRaw = ref(uni.getStorageSync('currency_symbol_raw') || uni.getStorageSync('currency_symbol') || '¥')
  const currencySymbol = ref(localText(currencySymbolRaw.value, getActiveLocale()) || normalizeString(currencySymbolRaw.value) || '¥')
  const appNameRaw = ref(uni.getStorageSync('app_name_raw') || uni.getStorageSync('app_name') || 'RichRoad')
  const appName = ref(uni.getStorageSync('app_name') || 'RichRoad')
  const appLogo = ref(uni.getStorageSync('app_logo') || '')
  const loaded = ref(false)

  const refreshLocalizedConfig = () => {
    const localizedName = localText(appNameRaw.value, getActiveLocale()) || normalizeString(appNameRaw.value) || 'RichRoad'
    appName.value = localizedName
    uni.setStorageSync('app_name', localizedName)

    const localizedCurrency = localText(currencySymbolRaw.value, getActiveLocale()) || normalizeString(currencySymbolRaw.value) || '¥'
    currencySymbol.value = localizedCurrency
    uni.setStorageSync('currency_symbol', localizedCurrency)
  }

  const loadConfig = async () => {
    if (loaded.value) {
      refreshLocalizedConfig()
      return
    }

    try {
      const [symRes, nameRes, logoRes] = await Promise.all([
        getCurrencySymbol({ includeI18n: true }),
        getAppName({ includeI18n: true }),
        getAppLogo(),
      ])

      if (symRes.code === 0 && symRes.data) {
        const rawSymbol = normalizeRawConfig(resolveConfigValue(symRes.data))
        if (rawSymbol) {
          currencySymbolRaw.value = rawSymbol
          uni.setStorageSync('currency_symbol_raw', rawSymbol)
        }
      }

      if (nameRes.code === 0 && nameRes.data) {
        const rawNameValue = normalizeRawConfig(resolveConfigValue(nameRes.data))
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
    currencySymbolRaw,
    appName,
    appNameRaw,
    appLogo,
    loaded,
    loadConfig,
    refreshLocalizedConfig,
    formatPrice,
  }
})
