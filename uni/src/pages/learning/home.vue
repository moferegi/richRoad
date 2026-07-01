<template>
  <view class="learning-home-container">
    <!-- Header -->
    <view class="header">
      <view class="lang-switch" @click="switchLanguage">
        <text>{{ currentLang === 'zh' ? '中' : 'En' }}</text>
      </view>
      <text class="app-title">{{ sysConfig.appTitle || t('app.title', 'RichRoad English') }}</text>
      <image class="app-logo" :src="sysConfig.logo" mode="aspectFit"></image>
    </view>

    <!-- 打卡区域 -->
    <view class="checkin-board">
      <view class="stats-row">
        <view class="stat-box">
          <text class="stat-num">{{ stats.continuousDays }}</text>
          <text class="stat-label">{{ t('home.continuous_days') }}</text>
        </view>
        <view class="stat-box">
          <text class="stat-num">{{ stats.totalDays }}</text>
          <text class="stat-label">{{ t('home.total_days') }}</text>
        </view>
      </view>
      
      <!-- 今日目标 -->
      <view class="target-card">
        <view class="target-title">
          <text v-if="stats.wordsToday >= sysConfig.dailyTarget">{{ t('home.target_finished') }}</text>
          <text v-else>{{ t('home.daily_target') }} {{ stats.wordsToday }}/{{ sysConfig.dailyTarget }}</text>
        </view>
        <progress :percent="targetPercent" activeColor="#409eff" backgroundColor="#ebeef5" class="target-progress"></progress>
      </view>

      <!-- 本周打卡记录 (小日历) -->
      <view class="week-calendar" @click="openCalendarPopup">
        <view class="week-title">{{ t('home.week_record') }}</view>
        <view class="week-days">
          <view 
            v-for="(day, index) in weekDays" 
            :key="index" 
            class="day-box" 
            :class="{ active: day.checked }">
            {{ day.name }}
          </view>
        </view>
      </view>
      
      <button class="checkin-btn" :disabled="checkedToday" @click="doCheckin">
        {{ checkedToday ? t('home.checked_today', '今日已打卡') : t('home.btn_checkin') }}
      </button>
    </view>

    <view v-if="calendarVisible" class="calendar-mask" @click="closeCalendarPopup">
      <view class="calendar-popup" @click.stop>
        <view class="calendar-popup-head">
          <text class="calendar-popup-title">{{ t('learningHomeCalendarTitle') }}</text>
          <text class="calendar-popup-close" @click="closeCalendarPopup">{{ t('learningHomeCalendarClose') }}</text>
        </view>

        <view class="calendar-popup-stats">
          <view class="calendar-stat-item">
            <text class="calendar-stat-label">{{ t('learningCheckinRecordContinuousDays') }}</text>
            <text class="calendar-stat-value">{{ stats.continuousDays }}</text>
          </view>
          <view class="calendar-stat-item">
            <text class="calendar-stat-label">{{ t('learningCheckinRecordTotalDays') }}</text>
            <text class="calendar-stat-value">{{ stats.totalDays }}</text>
          </view>
        </view>

        <view class="calendar-month-nav">
          <view class="calendar-nav-btn" @click="prevCalendarMonth">‹</view>
          <text class="calendar-month-text">{{ calendarMonthLabel }}</text>
          <view class="calendar-nav-btn" :class="{ disabled: !canGoNextMonth }" @click="nextCalendarMonth">›</view>
        </view>

        <view class="calendar-week-header">
          <text v-for="(label, index) in calendarWeekLabels" :key="`${label}-${index}`" class="calendar-week-label">{{ label }}</text>
        </view>

        <view class="calendar-grid">
          <view
            v-for="cell in calendarCells"
            :key="cell.key"
            class="calendar-cell"
            :class="{ 'calendar-cell-empty': !cell.day }"
            @click="handleCalendarCellTap(cell)"
          >
            <view v-if="cell.day" class="calendar-day" :class="{ 'calendar-day-checked': cell.checked, 'calendar-day-today': cell.isToday }">
              <text class="calendar-day-num">{{ cell.day }}</text>
              <text v-if="cell.checked" class="calendar-day-mark">●</text>
            </view>
          </view>
        </view>

        <view v-if="calendarLoading" class="calendar-loading">{{ t('learningHomeCalendarLoading') }}</view>

        <view class="calendar-popup-footer">
          <text class="calendar-hint">{{ checkedToday ? t('learningCheckinRecordCheckedToday') : t('learningHomeCalendarHint') }}</text>
          <button class="calendar-checkin-btn" size="mini" :disabled="checkedToday" @click="doCheckin">
            {{ checkedToday ? t('home.checked_today') : t('learningHomeCalendarSignToday') }}
          </button>
        </view>
      </view>
    </view>

    <!-- 视频分栏 -->
    <view class="video-section">
      <scroll-view scroll-x class="video-tabs">
        <view 
          v-for="cat in videoCategories" 
          :key="cat.id" 
          class="tab-item"
          :class="{ active: currentCategory === cat.id }"
          @click="currentCategory = cat.id">
          {{ localText(cat.name) }}
        </view>
      </scroll-view>

      <view class="video-list">
        <view class="video-item" v-for="video in videoList" :key="video.id" @click="goDetail(video.id)">
          <image class="video-cover" :src="video.coverUrl || ''" mode="aspectFill"></image>
          <view class="video-info">
            <text class="video-name">{{ localText(video.name) }}</text>
            <view class="video-stats">
              <text>{{ t('video.views') }}: {{ video.viewCount }}</text>
              <text>{{ t('video.users') }}: {{ video.userCount }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useLangStore } from '@/pinia/modules/lang.js'
import { t as i18nT, localText as i18nLocalText } from '@/utils/i18n.js'
import { doCheckin as doCheckinApi, getCheckinStats, getLearningCheckinRecordList, getVideoCategoryList, getVideoSeriesList } from '@/api/learning.js'
import { getAppName, getAppLogo } from '@/api/sysConfig.js'

const langStore = useLangStore()
const currentLang = computed(() => langStore.locale || uni.getStorageSync('app-lang') || 'zh')

const fallbackTexts = {
  'home.continuous_days': '连续打卡',
  'home.total_days': '累计打卡',
  'home.target_finished': '今日目标已完成',
  'home.daily_target': '今日目标',
  'home.checked_today': '今日已打卡',
  'home.week_record': '本周打卡记录',
  'home.btn_checkin': '今日打卡',
  'home.checkin_success': '打卡成功',
  learningHomeCalendarTitle: '打卡日历',
  learningHomeCalendarClose: '关闭',
  learningHomeCalendarHint: '点击今日日期可快速签到',
  learningHomeCalendarSignToday: '今日打卡',
  learningHomeCalendarOnlyToday: '仅支持当天签到',
  learningHomeCalendarSigned: '已签到',
  learningHomeCalendarLoading: '签到数据加载中...',
  'video.views': '浏览量',
  'video.users': '观看人数',
}

const t = (key, defaultText = '') => {
  const text = i18nT(key, currentLang.value)
  if (text && text !== key) return text
  return defaultText || fallbackTexts[key] || key
}

// Extract multi-language JSON storage locally
const localText = (jsonStr) => {
  return i18nLocalText(jsonStr, currentLang.value)
}

const sysConfig = ref({ logo: '', dailyTarget: 10, appTitle: '' })
const stats = ref({ continuousDays: 0, totalDays: 0, wordsToday: 0 })
const checkedToday = ref(false)

const targetPercent = computed(() => {
  if (sysConfig.value.dailyTarget === 0) return 0
  return Math.min(100, (stats.value.wordsToday / sysConfig.value.dailyTarget) * 100)
})

const weekDays = ref([])

const videoCategories = ref([])
const currentCategory = ref(0)
const videoList = ref([])

const calendarVisible = ref(false)
const calendarYear = ref(new Date().getFullYear())
const calendarMonth = ref(new Date().getMonth() + 1)
const calendarLoading = ref(false)

const checkinDateKeys = ref(new Set())
const checkinRecordPage = ref(1)
const checkinRecordTotal = ref(null)
const checkinRecordPageSize = 100
const checkinRecordLoading = ref(false)
const oldestCheckinDateTs = ref(null)

const getEntityId = (item) => Number(item?.id || item?.ID || 0)

const resetCheckinDateCache = () => {
  checkinDateKeys.value = new Set()
  checkinRecordPage.value = 1
  checkinRecordTotal.value = null
  oldestCheckinDateTs.value = null
}

const formatDateYMD = (date) => {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const parseDateKey = (value) => {
  if (!value) return ''
  const raw = String(value).trim()
  const matched = raw.match(/^(\d{4})-(\d{2})-(\d{2})/)
  if (matched) {
    return `${matched[1]}-${matched[2]}-${matched[3]}`
  }

  const date = new Date(raw)
  if (Number.isNaN(date.getTime())) return ''
  return formatDateYMD(date)
}

const dateKeyToTimestamp = (dateKey) => {
  if (!dateKey) return null
  const date = new Date(`${dateKey}T00:00:00`)
  if (Number.isNaN(date.getTime())) return null
  return date.getTime()
}

const calendarWeekLabels = computed(() => [
  t('weekDayMonShort', 'M'),
  t('weekDayTueShort', 'T'),
  t('weekDayWedShort', 'W'),
  t('weekDayThuShort', 'T'),
  t('weekDayFriShort', 'F'),
  t('weekDaySatShort', 'S'),
  t('weekDaySunShort', 'S'),
])

const calendarMonthLabel = computed(() => `${calendarYear.value}.${String(calendarMonth.value).padStart(2, '0')}`)

const canGoNextMonth = computed(() => {
  const now = new Date()
  const currentValue = calendarYear.value * 100 + calendarMonth.value
  const latestValue = now.getFullYear() * 100 + (now.getMonth() + 1)
  return currentValue < latestValue
})

const getCurrentWeekDateKeys = () => {
  const today = new Date()
  const dayIndex = (today.getDay() + 6) % 7
  const monday = new Date(today)
  monday.setDate(today.getDate() - dayIndex)

  return Array.from({ length: 7 }).map((_, index) => {
    const day = new Date(monday)
    day.setDate(monday.getDate() + index)
    return formatDateYMD(day)
  })
}

const updateWeekDays = (checkedDateKeys = new Set()) => {
  const currentWeekDateKeys = getCurrentWeekDateKeys()
  weekDays.value = calendarWeekLabels.value.map((name, index) => {
    const dateKey = currentWeekDateKeys[index]
    return {
      name,
      checked: checkedDateKeys.has(dateKey)
    }
  })
}

const appendCheckinDateKeys = (records = []) => {
  const nextSet = new Set(checkinDateKeys.value)
  let oldestTs = oldestCheckinDateTs.value

  records.forEach((item) => {
    const key = parseDateKey(item?.checkinDate || item?.CheckinDate)
    if (!key) return
    nextSet.add(key)
    const ts = dateKeyToTimestamp(key)
    if (ts === null) return
    if (oldestTs === null || ts < oldestTs) {
      oldestTs = ts
    }
  })

  checkinDateKeys.value = nextSet
  oldestCheckinDateTs.value = oldestTs
}

const hasLoadedAllCheckinRecords = () => {
  if (checkinRecordTotal.value === null) return false
  return (checkinRecordPage.value - 1) * checkinRecordPageSize >= checkinRecordTotal.value
}

const fetchNextCheckinRecordPage = async () => {
  if (checkinRecordLoading.value || hasLoadedAllCheckinRecords()) {
    return false
  }

  checkinRecordLoading.value = true
  try {
    const res = await getLearningCheckinRecordList({ page: checkinRecordPage.value, pageSize: checkinRecordPageSize })
    if (res.code !== 0 || !res.data) {
      return false
    }

    const data = res.data || {}
    checkinRecordTotal.value = Number(data.total || 0)
    const list = Array.isArray(data.list) ? data.list : []

    appendCheckinDateKeys(list)
    checkinRecordPage.value += 1
    return list.length > 0
  } catch (error) {
    return false
  } finally {
    checkinRecordLoading.value = false
  }
}

const ensureCheckinDateCoverage = async (year, month) => {
  const targetStartTs = new Date(year, month - 1, 1, 0, 0, 0, 0).getTime()
  let guard = 0

  while (guard < 24) {
    const hasCoverage = oldestCheckinDateTs.value !== null && oldestCheckinDateTs.value <= targetStartTs
    if (hasCoverage || hasLoadedAllCheckinRecords()) {
      break
    }

    const loaded = await fetchNextCheckinRecordPage()
    if (!loaded) {
      break
    }
    guard += 1
  }
}

const refreshWeekCheckins = () => {
  updateWeekDays(checkinDateKeys.value)
}

const refreshCalendarMonth = async () => {
  calendarLoading.value = true
  await ensureCheckinDateCoverage(calendarYear.value, calendarMonth.value)
  refreshWeekCheckins()
  calendarLoading.value = false
}

const calendarCells = computed(() => {
  const year = calendarYear.value
  const month = calendarMonth.value
  const firstDayIndex = (new Date(year, month - 1, 1).getDay() + 6) % 7
  const daysInMonth = new Date(year, month, 0).getDate()
  const totalCells = Math.ceil((firstDayIndex + daysInMonth) / 7) * 7
  const todayKey = formatDateYMD(new Date())

  return Array.from({ length: totalCells }).map((_, index) => {
    const day = index - firstDayIndex + 1
    if (day < 1 || day > daysInMonth) {
      return {
        key: `empty-${year}-${month}-${index}`,
        day: 0,
        dateKey: '',
        checked: false,
        isToday: false,
        isFuture: false,
      }
    }

    const dateKey = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    return {
      key: dateKey,
      day,
      dateKey,
      checked: checkinDateKeys.value.has(dateKey),
      isToday: dateKey === todayKey,
      isFuture: dateKey > todayKey,
    }
  })
})

const prevCalendarMonth = async () => {
  if (calendarMonth.value === 1) {
    calendarYear.value -= 1
    calendarMonth.value = 12
  } else {
    calendarMonth.value -= 1
  }
  await refreshCalendarMonth()
}

const nextCalendarMonth = async () => {
  if (!canGoNextMonth.value) return
  if (calendarMonth.value === 12) {
    calendarYear.value += 1
    calendarMonth.value = 1
  } else {
    calendarMonth.value += 1
  }
  await refreshCalendarMonth()
}

const handleCalendarCellTap = async (cell) => {
  if (!cell?.day) return

  if (cell.isToday && !checkedToday.value) {
    await doCheckin()
    return
  }

  if (cell.checked) {
    uni.showToast({ title: t('learningHomeCalendarSigned'), icon: 'none' })
    return
  }

  if (cell.isFuture || !cell.isToday) {
    uni.showToast({ title: t('learningHomeCalendarOnlyToday'), icon: 'none' })
  }
}

const switchLanguage = () => {
  const next = currentLang.value === 'zh' ? 'en' : 'zh'
  langStore.setLocale(next, { manual: true })
}

const openCalendarPopup = async () => {
  const now = new Date()
  calendarYear.value = now.getFullYear()
  calendarMonth.value = now.getMonth() + 1
  calendarVisible.value = true
  await refreshCalendarMonth()
}

const closeCalendarPopup = () => {
  calendarVisible.value = false
}

const loadWeekCheckins = async () => {
  const now = new Date()
  await ensureCheckinDateCoverage(now.getFullYear(), now.getMonth() + 1)
  refreshWeekCheckins()
}

const loadCheckin = async () => {
  const res = await getCheckinStats()
  if (res.code !== 0 || !res.data) {
    await loadWeekCheckins()
    return
  }
  stats.value = {
    continuousDays: Number(res.data.continuousDays || 0),
    totalDays: Number(res.data.totalDays || 0),
    wordsToday: Number(res.data.wordsToday || 0)
  }
  checkedToday.value = !!res.data.checkedToday
  if (res.data.dailyTarget && Number(res.data.dailyTarget) > 0) {
    sysConfig.value = { ...sysConfig.value, dailyTarget: Number(res.data.dailyTarget) }
  }
  await loadWeekCheckins()
}

const loadSysConfig = async () => {
  try {
    const [nameRes, logoRes] = await Promise.all([getAppName(), getAppLogo()])
    if (nameRes.code === 0 && nameRes.data?.configValue) {
      sysConfig.value = { ...sysConfig.value, appTitle: String(nameRes.data.configValue) }
    }
    if (logoRes.code === 0 && logoRes.data?.configValue) {
      sysConfig.value = { ...sysConfig.value, logo: String(logoRes.data.configValue) }
    }
  } catch (e) {
    // 使用默认值继续
  }
}

const loadVideoCategories = async () => {
  const res = await getVideoCategoryList({ page: 1, pageSize: 50 })
  if (res.code !== 0 || !res.data) {
    return
  }
  const list = Array.isArray(res.data.list) ? res.data.list : []
  videoCategories.value = list.map((item) => ({
    ...item,
    id: getEntityId(item)
  })).filter((item) => item.id > 0)
  if (videoCategories.value.length > 0 && !currentCategory.value) {
    currentCategory.value = videoCategories.value[0].id
  }
}

const loadVideoList = async () => {
  if (!currentCategory.value) {
    videoList.value = []
    return
  }
  const res = await getVideoSeriesList({ page: 1, pageSize: 20, categoryId: currentCategory.value })
  if (res.code !== 0 || !res.data) {
    return
  }
  const list = Array.isArray(res.data.list) ? res.data.list : []
  videoList.value = list.map((item) => ({
    ...item,
    id: getEntityId(item),
    viewCount: Number(item?.viewCount || 0),
    userCount: Number(item?.userCount || 0)
  }))
}

const doCheckin = async () => {
  if (checkedToday.value) {
    return
  }
  const res = await doCheckinApi()
  if (res.code === 0) {
    uni.showToast({ title: t('home.checkin_success'), icon: 'success' })
    resetCheckinDateCache()
    await loadCheckin()
    if (calendarVisible.value) {
      await refreshCalendarMonth()
    }
  }
}

const goDetail = (seriesId) => {
  uni.navigateTo({ url: `/pages/learning/video-detail?seriesId=${seriesId}` })
}

watch(currentCategory, () => {
  loadVideoList()
})

onMounted(async () => {
  updateWeekDays()
  resetCheckinDateCache()
  await Promise.all([loadCheckin(), loadVideoCategories(), loadSysConfig()])
  await loadVideoList()
})

onShow(() => {
  resetCheckinDateCache()
  loadCheckin()
})
</script>

<style scoped>
.learning-home-container { padding: 20rpx; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 30rpx; }
.checkin-board { background: #fff; border-radius: 20rpx; padding: 30rpx; margin-bottom: 30rpx; box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.05); }
.stats-row { display: flex; justify-content: space-around; margin-bottom: 20rpx; }
.stat-box { display: flex; flex-direction: column; align-items: center; }
.stat-num { font-size: 40rpx; font-weight: bold; color: #333; }
.stat-label { font-size: 24rpx; color: #888; margin-top: 10rpx; }
.target-card { margin: 30rpx 0; }
.target-title { font-size: 28rpx; margin-bottom: 10rpx; }
.week-calendar { margin-bottom: 30rpx; }
.week-days { display: flex; justify-content: space-between; margin-top: 15rpx; }
.day-box { width: 50rpx; height: 50rpx; line-height: 50rpx; text-align: center; border-radius: 8rpx; background: #eee; font-size: 24rpx; }
.day-box.active { background: #409eff; color: #fff; }
.checkin-btn { background: #409eff; color: #fff; border-radius: 40rpx; }
.calendar-mask {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(15, 23, 42, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx;
  box-sizing: border-box;
}
.calendar-popup {
  width: 100%;
  max-width: 680rpx;
  border-radius: 20rpx;
  background: #fff;
  padding: 24rpx;
  box-sizing: border-box;
}
.calendar-popup-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.calendar-popup-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #111827;
}
.calendar-popup-close {
  font-size: 24rpx;
  color: #6b7280;
}
.calendar-popup-stats {
  margin-top: 16rpx;
  display: flex;
  gap: 16rpx;
}
.calendar-stat-item {
  flex: 1;
  background: #f8fafc;
  border-radius: 12rpx;
  padding: 16rpx;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}
.calendar-stat-label {
  font-size: 22rpx;
  color: #6b7280;
}
.calendar-stat-value {
  font-size: 30rpx;
  color: #111827;
  font-weight: 700;
}
.calendar-month-nav {
  margin-top: 16rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.calendar-nav-btn {
  width: 56rpx;
  height: 56rpx;
  border-radius: 28rpx;
  background: #f3f4f6;
  color: #111827;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 34rpx;
  line-height: 1;
}
.calendar-nav-btn.disabled {
  color: #9ca3af;
}
.calendar-month-text {
  font-size: 30rpx;
  color: #111827;
  font-weight: 600;
}
.calendar-week-header {
  margin-top: 14rpx;
  display: flex;
}
.calendar-week-label {
  width: 14.2857%;
  text-align: center;
  font-size: 22rpx;
  color: #6b7280;
}
.calendar-grid {
  margin-top: 10rpx;
  display: flex;
  flex-wrap: wrap;
}
.calendar-cell {
  width: 14.2857%;
  padding: 8rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.calendar-day {
  width: 68rpx;
  height: 68rpx;
  border-radius: 34rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4rpx;
  background: #f8fafc;
}
.calendar-day-num {
  font-size: 24rpx;
  color: #4b5563;
}
.calendar-day-mark {
  font-size: 14rpx;
  color: #2563eb;
  line-height: 1;
}
.calendar-day-checked {
  background: #dbeafe;
}
.calendar-day-checked .calendar-day-num {
  color: #1d4ed8;
  font-weight: 600;
}
.calendar-day-today {
  border: 1rpx solid #2563eb;
}
.calendar-loading {
  text-align: center;
  color: #9ca3af;
  font-size: 22rpx;
  padding: 8rpx 0;
}
.calendar-popup-footer {
  margin-top: 14rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12rpx;
}
.calendar-hint {
  flex: 1;
  font-size: 22rpx;
  color: #6b7280;
}
.calendar-checkin-btn {
  background: #2563eb;
  color: #fff;
  font-size: 24rpx;
  padding: 0 22rpx;
}
/* Video stuff styling */
.video-tabs { white-space: nowrap; margin-bottom: 20rpx; }
.tab-item { display: inline-block; padding: 15rpx 30rpx; font-size: 28rpx; color: #555; }
.tab-item.active { color: #409eff; font-weight: bold; border-bottom: 4rpx solid #409eff; }
.video-item { display: flex; margin-bottom: 20rpx; background: #fff; padding: 15rpx; border-radius: 12rpx; }
.video-cover { width: 200rpx; height: 120rpx; border-radius: 8rpx; margin-right: 20rpx; }
.video-info { flex: 1; display: flex; flex-direction: column; justify-content: space-between; }
.video-name { font-size: 28rpx; font-weight: bold; }
.video-stats { font-size: 22rpx; color: #888; display: flex; justify-content: space-between; }
</style>
