import { defineStore } from 'pinia'
import storage from '@/utils/storage'

const LANG_KEY = 'app-lang'
const DEFAULT_LANG = 'zh'

export const useLangStore = defineStore('lang', {
  state: () => ({
    locale: storage.get(LANG_KEY) || DEFAULT_LANG,
    messages: {}
  }),
  actions: {
    setLocale(lang) {
      this.locale = lang
      storage.set(LANG_KEY, lang)
    },
    async initLangs() {
      // 动态导入语言包
      try {
        const modules = import.meta.glob('@/utils/i18n-locales/*.js')
        const langs = {}
        for (const path in modules) {
          const mod = await modules[path]()
          const name = path.split('/').pop().replace('.js', '')
          langs[name] = mod.default || mod
        }
        this.messages = langs
      } catch (e) {
        console.warn('load i18n locales failed', e)
      }
    }
  }
})
