import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getCurrencySymbol, getAppName, getAppLogo } from '@/api/sysConfig.js'

export const useAppConfigStore = defineStore('appConfig', () => {
  const currencySymbol = ref(uni.getStorageSync('currency_symbol') || '¥')
  const appName = ref(uni.getStorageSync('app_name') || 'RichRoad')
  const appLogo = ref(uni.getStorageSync('app_logo') || '')
  const loaded = ref(false)

  const loadConfig = async () => {
    if (loaded.value) return
    try {
      const [symRes, nameRes, logoRes] = await Promise.all([
        getCurrencySymbol(),
        getAppName(),
        getAppLogo(),
      ])
      if (symRes.code === 0 && symRes.data) {
        const sym = typeof symRes.data === 'object' ? symRes.data.configValue : symRes.data
        if (sym) { currencySymbol.value = sym; uni.setStorageSync('currency_symbol', sym) }
      }
      if (nameRes.code === 0 && nameRes.data) {
        const name = typeof nameRes.data === 'object' ? nameRes.data.configValue : nameRes.data
        if (name) { appName.value = name; uni.setStorageSync('app_name', name) }
      }
      if (logoRes.code === 0 && logoRes.data) {
        const logo = typeof logoRes.data === 'object' ? logoRes.data.configValue : logoRes.data
        if (logo) { appLogo.value = logo; uni.setStorageSync('app_logo', logo) }
      }
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
    appLogo,
    loaded,
    loadConfig,
    formatPrice,
  }
})
