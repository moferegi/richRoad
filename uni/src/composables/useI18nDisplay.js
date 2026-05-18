import { computed } from 'vue'
import { useLangStore } from '@/pinia/modules/lang.js'
import { resolveI18nDisplayText } from '@/utils/i18n.js'

export function useI18nDisplay(localeRef) {
  const langStore = useLangStore()

  const locale = computed(() => {
    const preferredLocale = localeRef && typeof localeRef === 'object' ? localeRef.value : ''
    return preferredLocale || langStore.locale || uni.getStorageSync('app-lang') || 'zh'
  })

  const resolveDisplayText = (candidate, fallback = '') => {
    return resolveI18nDisplayText(candidate, fallback, locale.value)
  }

  return {
    locale,
    resolveDisplayText,
  }
}
