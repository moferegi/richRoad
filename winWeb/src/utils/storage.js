/**
 * 本地存储封装
 */
const storage = {
  get(key) {
    try {
      const val = localStorage.getItem(key)
      if (val === null) return null
      try {
        return JSON.parse(val)
      } catch {
        return val
      }
    } catch {
      return null
    }
  },
  set(key, value) {
    try {
      const val = typeof value === 'object' ? JSON.stringify(value) : value
      localStorage.setItem(key, val)
    } catch (e) {
      console.warn('storage set error', e)
    }
  },
  remove(key) {
    try {
      localStorage.removeItem(key)
    } catch (e) {
      console.warn('storage remove error', e)
    }
  },
  clear() {
    try {
      localStorage.clear()
    } catch (e) {
      console.warn('storage clear error', e)
    }
  }
}

export default storage
