<template>
  <view class="nf-signin">
    <view class="nf-signin-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="goBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <text class="nf-navbar-title">{{ $t('signInTitle') }}</text>
        <view style="width: 64rpx;"></view>
      </view>
    </view>

    <view class="nf-body">
      <!-- 签到卡片 -->
      <view class="nf-signin-card">
        <view class="nf-signin-header">
          <text class="nf-signin-title">{{ $t('signInTitle') }}</text>
          <text class="nf-signin-sub">{{ $t('signInConsecutive').replace('{}', status.continuousDays || 0) }}</text>
        </view>
        <view class="nf-signin-btn-wrap">
          <view class="nf-signin-btn" :class="{ 'nf-signin-done': status.signed }" @tap="onSignIn">
            <text class="nf-signin-btn-text">{{ status.signed ? $t('signedToday') : $t('signInNow') }}</text>
          </view>
        </view>
      </view>

      <!-- 签到日历 -->
      <view class="nf-calendar-section">
        <view class="nf-section-title">
          <view class="nf-section-line"></view>
          <text>{{ $t('signInCalendar') }}</text>
          <view class="nf-section-line"></view>
        </view>
        <view class="nf-calendar"
          @touchstart="onCalendarTouchStart"
          @touchend="onCalendarTouchEnd"
        >
          <view class="nf-calendar-header">
            <view class="nf-calendar-nav" @tap="prevMonth">
              <uni-icons type="left" size="16" color="rgba(255,255,255,0.5)"></uni-icons>
            </view>
            <text class="nf-calendar-month">{{ calendarYear }}.{{ String(calendarMonth).padStart(2, '0') }}</text>
            <view class="nf-calendar-nav" @tap="nextMonth">
              <uni-icons type="right" size="16" color="rgba(255,255,255,0.5)"></uni-icons>
            </view>
          </view>
          <view class="nf-calendar-weekdays">
            <text class="nf-weekday" v-for="d in weekDays" :key="d">{{ d }}</text>
          </view>
          <view class="nf-calendar-grid">
            <view class="nf-calendar-cell" v-for="(cell, i) in calendarCells" :key="i">
              <template v-if="cell.day">
                <view class="nf-day" :class="{
                  'nf-day-signed': cell.signed,
                  'nf-day-today': cell.isToday
                }">
                  <text class="nf-day-num">{{ cell.day }}</text>
                  <text class="nf-day-dot" v-if="cell.signed">✓</text>
                </view>
              </template>
            </view>
          </view>
        </view>
      </view>

      <!-- 签到记录 -->
      <view class="nf-section-title">
        <view class="nf-section-line"></view>
        <text>{{ $t('signInRecord') }}</text>
        <view class="nf-section-line"></view>
      </view>

      <view class="nf-empty" v-if="records.length === 0 && !loading">
        <text class="nf-empty-text">{{ $t('noSignInRecord') }}</text>
      </view>

      <view class="nf-record-list" v-else>
        <view class="nf-record-card" v-for="(item, index) in records" :key="index">
          <view class="nf-record-icon">✅</view>
          <view class="nf-record-info">
            <text class="nf-record-date">{{ formatDate(item.signDate) }}</text>
          </view>
          <text class="nf-record-tag">{{ $t('signedDay') }}</text>
        </view>
      </view>

      <view class="nf-load-more" v-if="records.length > 0">
        <text class="nf-load-more-text" v-if="loading">...</text>
        <text class="nf-load-more-text" v-else-if="noMore">{{ $t('reachedBottom') }}</text>
        <text class="nf-load-more-text" v-else @tap="loadMoreRecords">{{ $t('loadMore') }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { doSignIn, getSignInStatus, getSignInRecords } from '@/api/signIn.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { onShow } from '@dcloudio/uni-app'

const langStore = useLangStore()
const $t = computed(() => langStore.$t)

const goBack = () => {
  uni.navigateBack({ fail: () => uni.switchTab({ url: '/pages/tabBar/my/index' }) })
}

const status = ref({
  signed: false,
  continuousDays: 0,
})

const records = ref([])
const loading = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20

// Calendar
const now = new Date()
const calendarYear = ref(now.getFullYear())
const calendarMonth = ref(now.getMonth() + 1)
const signedDatesSet = ref(new Set())

// Calendar swipe
const touchStartX = ref(0)

const onCalendarTouchStart = (e) => {
  touchStartX.value = e.touches[0].clientX
}

const onCalendarTouchEnd = (e) => {
  const deltaX = e.changedTouches[0].clientX - touchStartX.value
  if (Math.abs(deltaX) < 50) return // ignore small swipes
  if (deltaX < 0) {
    nextMonth()
  } else {
    prevMonth()
  }
}

const prevMonth = () => {
  if (calendarMonth.value === 1) {
    calendarMonth.value = 12
    calendarYear.value--
  } else {
    calendarMonth.value--
  }
  loadCalendarSignedDates()
}

const nextMonth = () => {
  const now = new Date()
  // Don't allow going beyond current month
  if (calendarYear.value === now.getFullYear() && calendarMonth.value === now.getMonth() + 1) return
  if (calendarMonth.value === 12) {
    calendarMonth.value = 1
    calendarYear.value++
  } else {
    calendarMonth.value++
  }
  loadCalendarSignedDates()
}

// Load signed dates for the displayed calendar month
const loadCalendarSignedDates = async () => {
  try {
    // Load enough records to cover the month — use a large pageSize
    const res = await getSignInRecords({ page: 1, pageSize: 100 })
    if (res.code === 0 && res.data && res.data.list) {
      const dateSet = new Set()
      res.data.list.forEach(item => {
        if (item.signDate) {
          const d = new Date(item.signDate)
          const ds = `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
          dateSet.add(ds)
        }
      })
      signedDatesSet.value = dateSet
    }
  } catch (e) { console.error(e) }
}

const weekDays = computed(() => [
  $t.value('weekDaySunShort'),
  $t.value('weekDayMonShort'),
  $t.value('weekDayTueShort'),
  $t.value('weekDayWedShort'),
  $t.value('weekDayThuShort'),
  $t.value('weekDayFriShort'),
  $t.value('weekDaySatShort'),
])

const calendarCells = computed(() => {
  const y = calendarYear.value
  const m = calendarMonth.value
  const firstDay = new Date(y, m - 1, 1).getDay()
  const daysInMonth = new Date(y, m, 0).getDate()
  const todayStr = `${now.getFullYear()}-${String(now.getMonth()+1).padStart(2,'0')}-${String(now.getDate()).padStart(2,'0')}`
  const cells = []
  for (let i = 0; i < firstDay; i++) cells.push({ day: 0 })
  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = `${y}-${String(m).padStart(2,'0')}-${String(d).padStart(2,'0')}`
    cells.push({
      day: d,
      signed: signedDatesSet.value.has(dateStr),
      isToday: dateStr === todayStr,
    })
  }
  return cells
})

const loadStatus = async () => {
  const res = await getSignInStatus()
  if (res.code === 0 && res.data) {
    status.value = {
      signed: res.data.signed || false,
      continuousDays: res.data.continuousDays || 0,
    }
  }
}

const loadRecords = async (isLoadMore = false) => {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getSignInRecords({ page: page.value, pageSize })
    if (res.code === 0 && res.data && res.data.list) {
      const list = res.data.list
      if (isLoadMore) {
        records.value = [...records.value, ...list]
      } else {
        records.value = list
      }
      // Build signed dates set for calendar
      list.forEach(item => {
        if (item.signDate) {
          const d = new Date(item.signDate)
          const ds = `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`
          signedDatesSet.value.add(ds)
        }
      })
      if (list.length < pageSize || records.value.length >= (res.data.total || Infinity)) {
        noMore.value = true
      }
    }
  } catch (e) { console.error(e) }
  loading.value = false
}

const loadMoreRecords = () => {
  if (noMore.value || loading.value) return
  page.value++
  loadRecords(true)
}

const onSignIn = async () => {
  if (status.value.signed) {
    uni.showToast({ title: $t.value('alreadySigned'), icon: 'none' })
    return
  }
  const res = await doSignIn()
  if (res.code === 0) {
    uni.showToast({ title: $t.value('signInSuccess'), icon: 'success' })
    loadStatus()
    // Reload records to reflect new sign-in
    page.value = 1
    noMore.value = false
    loadRecords()
  }
}

const formatDate = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
}

onShow(() => {
  loadStatus()
  page.value = 1
  noMore.value = false
  signedDatesSet.value = new Set()
  loadRecords()
})
</script>

<style scoped>
.nf-signin { min-height: 100vh; background: #000; position: relative; }
.nf-signin-bg {
  position: fixed; top: 0; left: 0; right: 0; height: 600rpx; z-index: 0; pointer-events: none;
  background: radial-gradient(ellipse at 50% 0%, rgba(229,9,20,0.12) 0%, transparent 60%);
}

.nf-navbar {
  background: rgba(0,0,0,0.85); backdrop-filter: blur(24px);
  border-bottom: 1rpx solid rgba(255,255,255,0.06);
  padding: 0 28rpx 16rpx; position: sticky; top: 0; z-index: 99;
}
.nf-navbar-status { height: var(--status-bar-height, 0px); }
.nf-navbar-content { display: flex; align-items: center; justify-content: space-between; height: 88rpx; }
.nf-navbar-back {
  width: 64rpx; height: 64rpx; border-radius: 50%;
  background: rgba(255,255,255,0.06); border: 1rpx solid rgba(255,255,255,0.1);
  display: flex; align-items: center; justify-content: center;
}
.nf-navbar-back:active { background: rgba(255,255,255,0.12); transform: scale(0.93); }
.nf-navbar-title { font-size: 34rpx; font-weight: 700; color: #fff; letter-spacing: 2rpx; }

.nf-body { padding: 24rpx; position: relative; z-index: 1; }

.nf-signin-card {
  background: rgba(255,255,255,0.04); border: 1rpx solid rgba(255,255,255,0.06);
  border-radius: 24rpx; padding: 60rpx 40rpx; text-align: center; margin-bottom: 40rpx;
}
.nf-signin-header {}
.nf-signin-title { font-size: 44rpx; font-weight: bold; color: #fff; display: block; }
.nf-signin-sub { font-size: 26rpx; color: rgba(255,255,255,0.5); margin-top: 10rpx; display: block; }
.nf-signin-btn-wrap { margin-top: 48rpx; }
.nf-signin-btn {
  display: inline-flex; align-items: center; justify-content: center;
  width: 300rpx; height: 100rpx; border-radius: 50rpx;
  background: linear-gradient(135deg, #e50914, #b20710);
}
.nf-signin-btn.nf-signin-done { background: rgba(255,255,255,0.1); }
.nf-signin-btn-text { font-size: 32rpx; font-weight: bold; color: #fff; }

/* Calendar */
.nf-calendar-section { margin-bottom: 40rpx; }
.nf-calendar {
  background: rgba(255,255,255,0.04); border: 1rpx solid rgba(255,255,255,0.06);
  border-radius: 20rpx; padding: 24rpx; margin-top: 16rpx;
}
.nf-calendar-header { display: flex; align-items: center; justify-content: center; gap: 24rpx; margin-bottom: 16rpx; }
.nf-calendar-nav {
  width: 48rpx; height: 48rpx; border-radius: 50%;
  background: rgba(255,255,255,0.06); display: flex; align-items: center; justify-content: center;
}
.nf-calendar-nav:active { background: rgba(255,255,255,0.12); }
.nf-calendar-month { font-size: 28rpx; color: rgba(255,255,255,0.6); font-weight: 600; }
.nf-calendar-weekdays { display: flex; margin-bottom: 8rpx; }
.nf-weekday { flex: 1; text-align: center; font-size: 22rpx; color: rgba(255,255,255,0.3); }
.nf-calendar-grid { display: flex; flex-wrap: wrap; }
.nf-calendar-cell { width: 14.2857%; display: flex; align-items: center; justify-content: center; padding: 8rpx 0; }
.nf-day {
  width: 64rpx; height: 64rpx; border-radius: 50%; display: flex; flex-direction: column;
  align-items: center; justify-content: center; position: relative;
}
.nf-day-num { font-size: 24rpx; color: rgba(255,255,255,0.5); }
.nf-day-dot { font-size: 16rpx; position: absolute; bottom: 2rpx; }
.nf-day-signed { background: rgba(229,9,20,0.15); }
.nf-day-signed .nf-day-num { color: #e50914; font-weight: 600; }
.nf-day-signed .nf-day-dot { color: #e50914; }
.nf-day-today { border: 1rpx solid rgba(229,9,20,0.4); }

/* Section title */
.nf-section-title {
  display: flex; align-items: center; justify-content: center;
  gap: 20rpx; margin-bottom: 24rpx;
}
.nf-section-title text { font-size: 26rpx; color: rgba(255,255,255,0.5); }
.nf-section-line {
  height: 1rpx; width: 80rpx;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.15), transparent);
}

/* Records */
.nf-empty { display: flex; flex-direction: column; align-items: center; padding: 60rpx 0; }
.nf-empty-text { font-size: 26rpx; color: rgba(255,255,255,0.3); }
.nf-record-list {}
.nf-record-card {
  display: flex; align-items: center; gap: 20rpx; padding: 24rpx 28rpx;
  background: rgba(255,255,255,0.04); border: 1rpx solid rgba(255,255,255,0.06);
  border-radius: 16rpx; margin-bottom: 12rpx;
}
.nf-record-icon { font-size: 28rpx; flex-shrink: 0; }
.nf-record-info { flex: 1; }
.nf-record-date { font-size: 28rpx; color: #fff; }
.nf-record-tag { font-size: 22rpx; color: #22c55e; background: rgba(34,197,94,0.1); padding: 4rpx 16rpx; border-radius: 8rpx; }
.nf-load-more { text-align: center; padding: 24rpx 0 60rpx; }
.nf-load-more-text { font-size: 24rpx; color: rgba(255,255,255,0.3); }
</style>
