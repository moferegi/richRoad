import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import messages, { t, localText } from '@/utils/i18n.js'

export const useLangStore = defineStore('lang', () => {
  const locale = ref(uni.getStorageSync('app-lang') || 'zh')

  const setLocale = (lang) => {
    locale.value = lang
    uni.setStorageSync('app-lang', lang)
    updateTabBar(lang)
  }

  const updateTabBar = (lang) => {
    const tabs = [
      { index: 0, key: 'tabHome' },
      { index: 1, key: 'tabCart' },
      { index: 2, key: 'tabCollect' },
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

  return { locale, setLocale, updateTabBar, $t, $lt }
})
