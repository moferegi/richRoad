import { defineStore } from 'pinia'
import { ref } from 'vue'

const STORAGE_KEY = 'nf_play_histories'

export const usePlayHistoryStore = defineStore('playHistory', () => {
  // Structure: { [goodID]: { lastEpisodeIndex, episodes: { [epIdx]: { currentTime, duration, episodeName, updatedAt } }, updatedAt } }
  const histories = ref(loadFromStorage())

  function loadFromStorage() {
    try {
      const raw = uni.getStorageSync(STORAGE_KEY)
      return raw ? JSON.parse(raw) : {}
    } catch (e) {
      return {}
    }
  }

  function persist() {
    try {
      uni.setStorageSync(STORAGE_KEY, JSON.stringify(histories.value))
    } catch (e) {}
  }

  /**
   * 保存播放进度（按单集存储）
   */
  function saveProgress(goodID, { episodeIndex, currentTime, duration, episodeName, imageUrl, title }) {
    if (!goodID || currentTime === undefined) return
    const id = String(goodID)
    const epKey = String(episodeIndex ?? 0)
    if (!histories.value[id] || !histories.value[id].episodes) {
      histories.value[id] = { lastEpisodeIndex: episodeIndex ?? 0, episodes: {}, updatedAt: Date.now() }
    }
    histories.value[id].lastEpisodeIndex = episodeIndex ?? 0
    histories.value[id].updatedAt = Date.now()
    if (imageUrl) histories.value[id].imageUrl = imageUrl
    if (title) histories.value[id].title = title
    histories.value[id].episodes[epKey] = {
      currentTime,
      duration: duration || 0,
      episodeName: episodeName || '',
      updatedAt: Date.now()
    }
    persist()
  }

  /**
   * 保存浏览记录（仅记录商品ID、图片、标题，不含播放进度）
   * 用于未登录时本地浏览历史
   */
  function saveBrowse(goodID, { imageUrl, title }) {
    if (!goodID) return
    const id = String(goodID)
    // 如果已存在完整播放记录，只更新时间和基础信息
    if (histories.value[id]) {
      histories.value[id].updatedAt = Date.now()
      if (imageUrl) histories.value[id].imageUrl = imageUrl
      if (title) histories.value[id].title = title
    } else {
      // 创建一个仅有浏览信息的记录
      histories.value[id] = {
        lastEpisodeIndex: 0,
        episodes: {},
        imageUrl: imageUrl || '',
        title: title || '',
        updatedAt: Date.now()
      }
    }
    persist()
  }

  /**
   * 获取最后观看的集的进度（用于“我的”页面和播放按鈕）
   * 返回 null 表示没有有效进度
   */
  function getProgress(goodID) {
    if (!goodID) return null
    const id = String(goodID)
    const rec = histories.value[id]
    if (!rec) return null

    // 将旧格式（没有 episodes）封装层兼容
    if (!rec.episodes) {
      if (!rec.currentTime || rec.currentTime < 5) return null
      if (rec.duration > 0 && rec.duration - rec.currentTime < 10) return null
      return {
        episodeIndex: rec.episodeIndex || 0,
        currentTime: rec.currentTime,
        duration: rec.duration || 0,
        episodeName: rec.episodeName || '',
        updatedAt: rec.updatedAt
      }
    }

    const epKey = String(rec.lastEpisodeIndex ?? 0)
    const epRec = rec.episodes[epKey]
    if (!epRec) return null
    if (epRec.currentTime < 5) return null
    if (epRec.duration > 0 && epRec.duration - epRec.currentTime < 10) return null
    return {
      episodeIndex: Number(epKey),
      currentTime: epRec.currentTime,
      duration: epRec.duration,
      episodeName: epRec.episodeName,
      updatedAt: epRec.updatedAt
    }
  }

  /**
   * 获取指定集的进度（用于点击集数按鈕时恢复）
   */
  function getEpisodeProgress(goodID, episodeIndex) {
    if (!goodID) return null
    const id = String(goodID)
    const epKey = String(episodeIndex ?? 0)
    const rec = histories.value[id]
    if (!rec) return null

    // 屁屠旧格式：如果旧格式且 episodeIndex 匹配
    if (!rec.episodes) {
      if ((rec.episodeIndex || 0) !== Number(episodeIndex)) return null
      if (!rec.currentTime || rec.currentTime < 5) return null
      if (rec.duration > 0 && rec.duration - rec.currentTime < 10) return null
      return { episodeIndex: rec.episodeIndex || 0, currentTime: rec.currentTime, duration: rec.duration || 0 }
    }

    const epRec = rec.episodes[epKey]
    if (!epRec) return null
    if (epRec.currentTime < 5) return null
    if (epRec.duration > 0 && epRec.duration - epRec.currentTime < 10) return null
    return {
      episodeIndex: Number(epKey),
      currentTime: epRec.currentTime,
      duration: epRec.duration,
      episodeName: epRec.episodeName
    }
  }

  /**
   * 获取最近播放列表（本地，按时间排序，最多返20条）
   * 返回 [{ ID, imageUrl, title, updatedAt }, ...]
   */
  function getRecentList(limit = 20) {
    const arr = []
    for (const [id, rec] of Object.entries(histories.value)) {
      if (rec.imageUrl) {
        arr.push({ ID: Number(id) || id, imageUrl: rec.imageUrl, title: rec.title || '', updatedAt: rec.updatedAt || 0 })
      }
    }
    arr.sort((a, b) => b.updatedAt - a.updatedAt)
    return arr.slice(0, limit)
  }

  /**
   * 从本地 storage 重新加载数据（页面 onShow 时调用）
   */
  function reload() {
    histories.value = loadFromStorage()
  }

  /**
   * 清空所有播放历史
   */
  function clearAll() {
    histories.value = {}
    persist()
  }

  return {
    histories,
    saveProgress,
    saveBrowse,
    getProgress,
    getEpisodeProgress,
    getRecentList,
    reload,
    clearAll
  }
})
