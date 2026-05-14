import zh from './zh.js'
import en from './en.js'
import mn from './mn.js'
import zhTW from './zh-TW.js'
import th from './th.js'
import hi from './hi.js'
import id from './id.js'

import viOverrides from './vi.js'
import arOverrides from './ar.js'
import jaOverrides from './ja.js'
import koOverrides from './ko.js'
import msOverrides from './ms.js'

const extendLocaleFromEnglish = (overrides = {}) => ({
  ...en,
  ...overrides,
})

const extendLocaleFromZh = (locale = {}) => ({
  ...zh,
  ...locale,
})

const messages = {
  zh,
  en,
  // 历史词库中部分语言存在缺键，这里以中文为基线补齐，保持原有显示兜底行为。
  mn: extendLocaleFromZh(mn),
  'zh-TW': extendLocaleFromZh(zhTW),
  th: extendLocaleFromZh(th),
  hi: extendLocaleFromZh(hi),
  id: extendLocaleFromZh(id),
  vi: extendLocaleFromEnglish(viOverrides),
  ar: extendLocaleFromEnglish(arOverrides),
  ja: extendLocaleFromEnglish(jaOverrides),
  ko: extendLocaleFromEnglish(koOverrides),
  ms: extendLocaleFromEnglish(msOverrides),
}

export default messages
