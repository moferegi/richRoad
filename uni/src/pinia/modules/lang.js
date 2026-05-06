import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import messages, { t, localText } from '@/utils/i18n.js'
import { getEnabledLanguages } from '@/api/language.js'

// 国旗映射
const flagMap = {
  zh: '🇨🇳', 'zh-TW': '🇹🇼', en: '🇬🇧', mn: '🇲🇳',
  th: '🇹🇭', hi: '🇮🇳', id: '🇮🇩', ja: '🇯🇵', ko: '🇰🇷',
  ru: '🇷🇺', fr: '🇫🇷', de: '🇩🇪', es: '🇪🇸', pt: '🇧🇷',
}

export const useLangStore = defineStore('lang', () => {
  const locale = ref(uni.getStorageSync('app-lang') || '')
  const enabledLangs = ref([]) // 从后端加载的启用语言列表
  const loaded = ref(false)

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

  const setLocale = (lang) => {
    locale.value = lang
    uni.setStorageSync('app-lang', lang)
    updateTabBar(lang)
  }

  const updateTabBar = (lang) => {
    const tabs = [
      { index: 0, key: 'tryonRoom' },
      { index: 1, key: 'shoeRoom' },
      { index: 2, key: 'clothesPage' },
      { index: 3, key: 'tabMy' },
    ]
    tabs.forEach(item => {
      uni.setTabBarItem({
        index: item.index,
        text: t(item.key, lang)
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

  return { locale, enabledLangs, loaded, initLangs, setLocale, updateTabBar, $t, $lt }
})
