import { localText } from '@/utils/i18n.js'

const toFiniteNumber = (value) => {
  const num = Number(value)
  return Number.isFinite(num) ? num : null
}

export const resolveLocalizedPriceFen = (basePriceFen, priceI18n, locale) => {
  const localized = localText(priceI18n, locale)
  const localizedNum = toFiniteNumber(localized)
  if (localizedNum !== null && localizedNum >= 0) {
    return localizedNum
  }

  const baseNum = toFiniteNumber(basePriceFen)
  if (baseNum !== null && baseNum >= 0) {
    return baseNum
  }

  return 0
}

export const formatLocalizedPrice = (basePriceFen, priceI18n, locale) => {
  const fen = resolveLocalizedPriceFen(basePriceFen, priceI18n, locale)
  return (fen / 100).toFixed(2)
}
