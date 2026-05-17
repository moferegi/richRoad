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
const EXCHANGE_RATE_CACHE_KEY = '__rr_exchange_rate_cache_v1__'
const DEFAULT_CACHE_TTL_MS = 30 * 60 * 1000

let exchangeRateCacheStore = null

const normalizeCurrency = (currency) => {
  const value = String(currency || '').trim().toUpperCase()
  return /^[A-Z]{3}$/.test(value) ? value : ''
}

const toFiniteNumber = (value) => {
  const num = Number(value)
  return Number.isFinite(num) ? num : null
}

const hasLocalStorage = () => {
  return typeof window !== 'undefined' && typeof window.localStorage !== 'undefined'
}

const normalizeRatesMap = (rawRates, baseCurrency) => {
  const rates = {
    [baseCurrency]: 1
  }

  if (!rawRates || typeof rawRates !== 'object') {
    return rates
  }

  Object.entries(rawRates).forEach(([currency, value]) => {
    const key = normalizeCurrency(currency)
    const num = toFiniteNumber(value)
    if (!key || !num || num <= 0) {
      return
    }
    rates[key] = num
  })

  rates[baseCurrency] = 1
  return rates
}

const normalizeCacheItem = (rawItem, baseCurrency) => {
  const item = rawItem && typeof rawItem === 'object' ? rawItem : {}
  const source = String(item.source || '').trim()
  const fetchedAt = String(item.fetchedAt || '').trim()
  const updatedAt = String(item.updatedAt || '').trim()
  const rates = normalizeRatesMap(item.rates, baseCurrency)
  return {
    source,
    fetchedAt,
    updatedAt,
    rates
  }
}

const readExchangeRateCacheStore = () => {
  if (!hasLocalStorage()) {
    return {}
  }

  try {
    const raw = window.localStorage.getItem(EXCHANGE_RATE_CACHE_KEY)
    if (!raw) {
      return {}
    }
    const parsed = JSON.parse(raw)
    if (!parsed || typeof parsed !== 'object') {
      return {}
    }

    const nextStore = {}
    Object.entries(parsed).forEach(([base, item]) => {
      const normalizedBase = normalizeCurrency(base)
      if (!normalizedBase) {
        return
      }
      nextStore[normalizedBase] = normalizeCacheItem(item, normalizedBase)
    })
    return nextStore
  } catch (_) {
    return {}
  }
}

const ensureExchangeRateCacheStore = () => {
  if (exchangeRateCacheStore) {
    return exchangeRateCacheStore
  }
  exchangeRateCacheStore = readExchangeRateCacheStore()
  return exchangeRateCacheStore
}

const persistExchangeRateCacheStore = () => {
  if (!hasLocalStorage()) {
    return
  }
  try {
    window.localStorage.setItem(EXCHANGE_RATE_CACHE_KEY, JSON.stringify(ensureExchangeRateCacheStore()))
  } catch (_) {
    // ignore localStorage write errors
  }
}

const getCachedBaseItem = (baseCurrency) => {
  const store = ensureExchangeRateCacheStore()
  const rawItem = store[baseCurrency]
  if (!rawItem) {
    return null
  }
  return normalizeCacheItem(rawItem, baseCurrency)
}

const saveCachedBaseItem = (baseCurrency, item) => {
  const store = ensureExchangeRateCacheStore()
  const normalizedItem = normalizeCacheItem(item, baseCurrency)
  store[baseCurrency] = normalizedItem
  persistExchangeRateCacheStore()
  return normalizedItem
}

const hasRatesForCurrencies = (cacheItem, currencies, baseCurrency) => {
  if (!cacheItem || !cacheItem.rates || typeof cacheItem.rates !== 'object') {
    return false
  }
  return currencies.every((currency) => {
    if (currency === baseCurrency) {
      return true
    }
    return Number.isFinite(Number(cacheItem.rates[currency])) && Number(cacheItem.rates[currency]) > 0
  })
}

const isCacheFresh = (cacheItem, cacheTtlMs) => {
  if (!cacheItem) {
    return false
  }
  const ttl = Number.isFinite(Number(cacheTtlMs)) ? Number(cacheTtlMs) : DEFAULT_CACHE_TTL_MS
  if (ttl <= 0) {
    return false
  }
  const updatedAt = String(cacheItem.updatedAt || '').trim()
  if (!updatedAt) {
    return false
  }
  const timestamp = Date.parse(updatedAt)
  if (!Number.isFinite(timestamp)) {
    return false
  }
  return (Date.now() - timestamp) <= ttl
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

export const getExchangeRateSnapshot = ({
  base = 'CNY'
} = {}) => {
  const baseCurrency = normalizeCurrency(base) || 'CNY'
  const cacheItem = getCachedBaseItem(baseCurrency)
  if (!cacheItem) {
    return null
  }
  return {
    base: baseCurrency,
    source: cacheItem.source || '',
    fetchedAt: cacheItem.fetchedAt || '',
    updatedAt: cacheItem.updatedAt || '',
    rates: normalizeRatesMap(cacheItem.rates, baseCurrency)
  }
}

export const clearExchangeRateCache = (base) => {
  const store = ensureExchangeRateCacheStore()
  const normalizedBase = normalizeCurrency(base)
  if (normalizedBase) {
    delete store[normalizedBase]
  } else {
    Object.keys(store).forEach((key) => {
      delete store[key]
    })
  }
  persistExchangeRateCacheStore()
}

export const fetchExchangeRates = async ({
  base = 'CNY',
  currencies = [],
  forceRefresh = false,
  cacheTtlMs = DEFAULT_CACHE_TTL_MS,
  useCache = true
} = {}) => {
  const baseCurrency = normalizeCurrency(base) || 'CNY'
  const normalizedCurrencies = Array.from(new Set(
    currencies
      .map((currency) => normalizeCurrency(currency))
      .filter(Boolean)
  ))

  const cachedItem = useCache ? getCachedBaseItem(baseCurrency) : null
  if (!forceRefresh && useCache && cachedItem && isCacheFresh(cachedItem, cacheTtlMs) && hasRatesForCurrencies(cachedItem, normalizedCurrencies, baseCurrency)) {
    return {
      base: baseCurrency,
      source: cachedItem.source ? `cache:${cachedItem.source}` : 'cache',
      fetchedAt: cachedItem.fetchedAt || '',
      rates: normalizeRatesMap(cachedItem.rates, baseCurrency)
    }
  }

  const rates = normalizeRatesMap(cachedItem?.rates, baseCurrency)
  const sourceFlags = []
  let fetchedAt = cachedItem?.fetchedAt || ''

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

  const networkSource = Array.from(new Set(sourceFlags)).join('+') || 'unknown'
  const cachedResult = saveCachedBaseItem(baseCurrency, {
    source: networkSource,
    fetchedAt,
    updatedAt: new Date().toISOString(),
    rates
  })

  return {
    base: baseCurrency,
    source: cachedResult.source || networkSource,
    fetchedAt: cachedResult.fetchedAt || fetchedAt,
    rates: normalizeRatesMap(cachedResult.rates, baseCurrency)
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
