export const DEFAULT_LANG_CURRENCY_MAP = Object.freeze({
  zh: 'CNY',
  en: 'USD',
  mn: 'MNT',
  'zh-TW': 'TWD',
  th: 'THB',
  hi: 'INR',
  id: 'IDR',
  vi: 'VND',
  ar: 'SAR',
  ja: 'JPY',
  ko: 'KRW',
  ms: 'MYR'
})

const FRANKFURTER_API = 'https://api.frankfurter.app/latest'
const OPEN_ER_API = 'https://open.er-api.com/v6/latest'

const normalizeCurrency = (currency) => {
  const value = String(currency || '').trim().toUpperCase()
  return /^[A-Z]{3}$/.test(value) ? value : ''
}

const toFiniteNumber = (value) => {
  const num = Number(value)
  return Number.isFinite(num) ? num : null
}

const collectMissingCurrencies = (targetCurrencies, rates, baseCurrency) => {
  return targetCurrencies.filter((currency) => {
    if (currency === baseCurrency) return false
    return !Number.isFinite(rates[currency])
  })
}

const fetchJson = async (url) => {
  const response = await fetch(url)
  if (!response.ok) {
    throw new Error(`rate request failed: ${response.status}`)
  }
  return response.json()
}

const fetchFrankfurterRates = async (baseCurrency, targetCurrencies) => {
  const query = new URLSearchParams({ from: baseCurrency })
  if (targetCurrencies.length > 0) {
    query.set('to', targetCurrencies.join(','))
  }
  const payload = await fetchJson(`${FRANKFURTER_API}?${query.toString()}`)
  const rates = payload?.rates && typeof payload.rates === 'object' ? payload.rates : {}
  return {
    source: 'frankfurter-ecb',
    date: payload?.date || '',
    rates
  }
}

const fetchOpenErRates = async (baseCurrency) => {
  const payload = await fetchJson(`${OPEN_ER_API}/${baseCurrency}`)
  const rates = payload?.rates && typeof payload.rates === 'object' ? payload.rates : {}
  return {
    source: 'open-er-api',
    date: payload?.time_last_update_utc || '',
    rates
  }
}

export const fetchExchangeRates = async ({
  base = 'CNY',
  currencies = []
} = {}) => {
  const baseCurrency = normalizeCurrency(base) || 'CNY'
  const normalizedCurrencies = Array.from(new Set(
    currencies
      .map((currency) => normalizeCurrency(currency))
      .filter(Boolean)
  ))

  const rates = {
    [baseCurrency]: 1
  }
  const sourceFlags = []
  let fetchedAt = ''

  try {
    const frankfurterResult = await fetchFrankfurterRates(baseCurrency, normalizedCurrencies)
    sourceFlags.push(frankfurterResult.source)
    if (frankfurterResult.date) {
      fetchedAt = String(frankfurterResult.date)
    }
    Object.entries(frankfurterResult.rates || {}).forEach(([currency, value]) => {
      const key = normalizeCurrency(currency)
      const num = toFiniteNumber(value)
      if (key && num && num > 0) {
        rates[key] = num
      }
    })
  } catch (_) {
    // fallback below
  }

  const missingAfterFrankfurter = collectMissingCurrencies(normalizedCurrencies, rates, baseCurrency)

  if (missingAfterFrankfurter.length > 0 || sourceFlags.length === 0) {
    const openErResult = await fetchOpenErRates(baseCurrency)
    sourceFlags.push(openErResult.source)
    if (!fetchedAt && openErResult.date) {
      fetchedAt = String(openErResult.date)
    }
    Object.entries(openErResult.rates || {}).forEach(([currency, value]) => {
      const key = normalizeCurrency(currency)
      const num = toFiniteNumber(value)
      if (key && num && num > 0) {
        rates[key] = num
      }
    })
  }

  return {
    base: baseCurrency,
    source: Array.from(new Set(sourceFlags)).join('+') || 'unknown',
    fetchedAt,
    rates
  }
}

export const calcConvertedFenFromCny = (baseFen, rate, increasePercent = 0) => {
  const fenNum = Number(baseFen)
  const rateNum = Number(rate)
  const percentNum = Number(increasePercent)

  if (!Number.isFinite(fenNum) || !Number.isFinite(rateNum) || rateNum <= 0) {
    return null
  }

  const factor = Number.isFinite(percentNum) ? (100 + percentNum) / 100 : 1
  const convertedMajor = (fenNum / 100) * rateNum * factor
  const convertedFen = Math.round(convertedMajor * 100)

  return convertedFen >= 0 ? convertedFen : null
}

export const normalizePriceI18nMap = (value) => {
  if (!value) return {}

  let payload = value
  if (typeof value === 'string') {
    const trimmed = value.trim()
    if (!trimmed) return {}
    try {
      payload = JSON.parse(trimmed)
    } catch (_) {
      return {}
    }
  }

  if (!payload || typeof payload !== 'object' || Array.isArray(payload)) {
    return {}
  }

  const normalized = {}
  Object.entries(payload).forEach(([lang, rawPrice]) => {
    const priceNum = Number(rawPrice)
    if (Number.isFinite(priceNum) && priceNum >= 0) {
      normalized[String(lang)] = Math.round(priceNum)
    }
  })

  return normalized
}

export const ensurePriceSlots = (priceMap, langCodes, fallbackFen = 0) => {
  const next = { ...normalizePriceI18nMap(priceMap) }
  const fallback = Number.isFinite(Number(fallbackFen)) ? Math.round(Number(fallbackFen)) : 0

  langCodes.forEach((code) => {
    if (!code) return
    if (!Object.prototype.hasOwnProperty.call(next, code)) {
      next[code] = fallback
    }
  })

  return next
}

export const serializePriceI18nMap = (priceMap) => {
  const normalized = normalizePriceI18nMap(priceMap)
  if (!Object.keys(normalized).length) {
    return ''
  }
  return JSON.stringify(normalized)
}

export const formatFen = (fen) => {
  const value = Number(fen)
  if (!Number.isFinite(value)) return '0.00'
  return (value / 100).toFixed(2)
}
