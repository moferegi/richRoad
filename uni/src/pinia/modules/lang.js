import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import messages, { t } from '@/utils/i18n.js'

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
      { index: 1, key: 'tabCollect' },
      { index: 2, key: 'tabMy' },
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

  return { locale, setLocale, updateTabBar, $t }
})
