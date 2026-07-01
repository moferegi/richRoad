import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import messages, { t, localText } from '@/utils/i18n.js'
import { getEnabledLanguages } from '@/api/language.js'
import { useAppConfigStore } from '@/pinia/modules/appConfig.js'

// 国旗映射
const flagMap = {
  zh: '🇨🇳', 'zh-TW': '🇹🇼', en: '🇬🇧', mn: '🇲🇳',
  th: '🇹🇭', hi: '🇮🇳', id: '🇮🇩', vi: '🇻🇳', ar: '🇸🇦', ja: '🇯🇵', ko: '🇰🇷', ms: '🇲🇾',
  ru: '🇷🇺', fr: '🇫🇷', de: '🇩🇪', es: '🇪🇸', pt: '🇧🇷',
}

const LANG_MANUAL_SELECTED_KEY = 'app-lang-manual-selected'
const LANG_PICKER_PROMPTED_KEY = 'app-lang-picker-prompted'

export const useLangStore = defineStore('lang', () => {
  const locale = ref(uni.getStorageSync('app-lang') || '')
  const enabledLangs = ref([]) // 从后端加载的启用语言列表
  const loaded = ref(false)
  const manuallySelected = ref(uni.getStorageSync(LANG_MANUAL_SELECTED_KEY) === '1')
  const pickerPrompted = ref(uni.getStorageSync(LANG_PICKER_PROMPTED_KEY) === '1')
  const tabBarSyncTimers = ref([])

  const markLanguageManualSelected = () => {
    manuallySelected.value = true
    uni.setStorageSync(LANG_MANUAL_SELECTED_KEY, '1')
  }

  const markLanguagePickerPrompted = () => {
    pickerPrompted.value = true
    uni.setStorageSync(LANG_PICKER_PROMPTED_KEY, '1')
  }

  const shouldAutoShowLanguagePicker = () => {
    return !manuallySelected.value && !pickerPrompted.value
  }

  const getEffectiveLocale = (lang) => {
    return lang || locale.value || uni.getStorageSync('app-lang') || 'mn'
  }

  const clearTabBarSyncTimers = () => {
    if (!Array.isArray(tabBarSyncTimers.value) || !tabBarSyncTimers.value.length) {
      return
    }
    tabBarSyncTimers.value.forEach((timer) => {
      clearTimeout(timer)
    })
    tabBarSyncTimers.value = []
  }

  // 初始化：从后端获取启用的语言和默认语言
  const initLangs = async () => {
    if (loaded.value) return
    try {
      const res = await getEnabledLanguages()
      if (res.code === 0 && res.data) {
        const list = Array.isArray(res.data) ? res.data : (res.data.list || [])
        enabledLangs.value = list.map(item => ({
          value: item.code,
          label: item.nativeName || item.name,
          native: item.name,
          flag: flagMap[item.code] || '🏳️',
          isDefault: item.isDefault,
        }))
        // 如果本地没有存储过语言，使用后端默认语言
        if (!uni.getStorageSync('app-lang')) {
          const defaultLang = list.find(l => l.isDefault)
          if (defaultLang) {
            locale.value = defaultLang.code
            uni.setStorageSync('app-lang', defaultLang.code)
          } else {
            locale.value = 'mn'
          }
        }
      }
    } catch (e) {
      console.error('initLangs error', e)
    }
    // 如果仍然没有locale值（API失败等），回退到mn
    if (!locale.value) {
      locale.value = 'mn'
    }
    // 无论语言来源于本地还是后端，初始化后都同步一次 tabBar 文案
    updateTabBar(locale.value)
    loaded.value = true
  }

  const setLocale = (lang, options = {}) => {
    const { manual = true } = options
    if (!lang) return
    locale.value = lang
    uni.setStorageSync('app-lang', lang)
    if (manual) {
      markLanguageManualSelected()
    }
    updateTabBar(lang)
    clearTabBarSyncTimers()
    ;[90, 260, 520].forEach((delay) => {
      const timer = setTimeout(() => {
        updateTabBar(lang)
      }, delay)
      tabBarSyncTimers.value.push(timer)
    })
    useAppConfigStore().refreshLocalizedConfig()
  }

  const updateTabBar = (lang, retry = 0) => {
    const currentLang = getEffectiveLocale(lang)
    const tabFallbackTextMap = {
      zh: ['首页', '跟打', '我的'],
      en: ['Home', 'Typing', 'Me'],
      mn: ['Нүүр', 'Бичих', 'Миний'],
    }
    const fallbackTexts = tabFallbackTextMap[currentLang] || tabFallbackTextMap.en
    const tabs = [
      { index: 0, key: 'englishTabHome', fallbackIndex: 0 },
      { index: 1, key: 'englishTabTyping', fallbackIndex: 1 },
      { index: 2, key: 'englishTabMy', fallbackIndex: 2 },
    ]

    let finished = 0
    let failed = 0
    const tryNext = () => {
      if (finished < tabs.length) return
      if (failed > 0 && retry < 8) {
        setTimeout(() => updateTabBar(currentLang, retry + 1), 120 * (retry + 1))
      }
    }

    tabs.forEach(item => {
      uni.setTabBarItem({
        index: item.index,
        text: (() => {
          const text = t(item.key, currentLang)
          if (text && text !== item.key) return text
          return fallbackTexts[item.fallbackIndex]
        })(),
        success: () => {
          finished += 1
          tryNext()
        },
        fail: () => {
          failed += 1
          finished += 1
          tryNext()
        }
      })
    })
  }

  const $t = computed(() => {
    return (key) => t(key, locale.value)
  })

  // 解析多语言数据库字段（响应式，语言切换自动更新）
  const $lt = computed(() => {
    return (value) => localText(value, locale.value)
  })

  return {
    locale,
    enabledLangs,
    loaded,
    manuallySelected,
    pickerPrompted,
    initLangs,
    setLocale,
    updateTabBar,
    shouldAutoShowLanguagePicker,
    markLanguagePickerPrompted,
    markLanguageManualSelected,
    $t,
    $lt,
  }
})
