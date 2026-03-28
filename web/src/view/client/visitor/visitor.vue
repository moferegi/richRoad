<template>
  <div>
    <!-- 今日实时统计卡片 -->
    <div class="gva-search-box" style="padding: 20px;">
      <el-row :gutter="20">
        <el-col :span="4">
          <el-statistic title="今日 PV" :value="todayStats.pv || 0">
            <template #prefix>
              <el-icon><View /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="今日 UV" :value="todayStats.uv || 0">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="已登录用户" :value="todayStats.registeredUser || 0">
            <template #prefix>
              <el-icon><UserFilled /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-statistic title="新访客" :value="todayStats.newVisitor || 0">
            <template #prefix>
              <el-icon><Plus /></el-icon>
            </template>
          </el-statistic>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="handleAggregate">手动聚合今日数据</el-button>
        </el-col>
      </el-row>
    </div>

    <!-- Tab 切换 -->
    <el-tabs v-model="activeTab" class="gva-table-box" style="padding: 0 20px;">
      <!-- 访客日志 Tab -->
      <el-tab-pane label="访客日志" name="log">
        <div class="gva-search-box">
          <el-form :inline="true" :model="logSearch" @keyup.enter="getLogList">
            <el-form-item label="日期范围">
              <el-date-picker v-model="logSearch.startDate" type="date" value-format="YYYY-MM-DD" placeholder="开始日期" />
              —
              <el-date-picker v-model="logSearch.endDate" type="date" value-format="YYYY-MM-DD" placeholder="结束日期" />
            </el-form-item>
            <el-form-item label="访客ID">
              <el-input v-model="logSearch.visitorId" placeholder="访客指纹ID" clearable />
            </el-form-item>
            <el-form-item label="IP">
              <el-input v-model="logSearch.ip" placeholder="IP地址" clearable />
            </el-form-item>
            <el-form-item label="平台">
              <el-select v-model="logSearch.platform" placeholder="全部" clearable>
                <el-option label="H5" value="h5" />
                <el-option label="iOS" value="ios" />
                <el-option label="Android" value="android" />
                <el-option label="微信小程序" value="wxmp" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="getLogList">查询</el-button>
              <el-button icon="refresh" @click="resetLogSearch">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-table :data="logList" style="width: 100%" tooltip-effect="dark">
          <el-table-column align="left" label="ID" prop="ID" width="80" />
          <el-table-column align="left" label="访客ID" prop="visitorId" width="200" show-overflow-tooltip />
          <el-table-column align="left" label="用户ID" prop="userId" width="80">
            <template #default="scope">
              {{ scope.row.userId || '未登录' }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="IP" prop="ip" width="140" />
          <el-table-column align="left" label="归属地" prop="location" width="160" show-overflow-tooltip />
          <el-table-column align="left" label="平台" prop="platform" width="100" />
          <el-table-column align="left" label="页面路径" prop="pagePath" width="200" show-overflow-tooltip />
          <el-table-column align="left" label="语言" prop="language" width="80" />
          <el-table-column align="left" label="屏幕" width="120">
            <template #default="scope">
              {{ scope.row.screenWidth }}×{{ scope.row.screenHeight }}
            </template>
          </el-table-column>
          <el-table-column align="left" label="UA" prop="userAgent" min-width="200" show-overflow-tooltip />
          <el-table-column align="left" label="访问时间" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
          </el-table-column>
        </el-table>
        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="logSearch.page"
            :page-size="logSearch.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="logTotal"
            @current-change="(val) => { logSearch.page = val; getLogList() }"
            @size-change="(val) => { logSearch.pageSize = val; getLogList() }"
          />
        </div>
      </el-tab-pane>

      <!-- 汇总统计 Tab -->
      <el-tab-pane label="汇总统计" name="summary">
        <!-- 趋势图表 -->
        <div class="gva-search-box" style="padding: 20px 24px;">
          <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
            <span style="font-weight: 600; font-size: 15px; letter-spacing: 0.5px;">访客趋势</span>
            <el-radio-group v-model="chartDays" size="small" @change="fetchChartData">
              <el-radio-button :value="7">近7天</el-radio-button>
              <el-radio-button :value="14">近14天</el-radio-button>
              <el-radio-button :value="30">近30天</el-radio-button>
            </el-radio-group>
          </div>
          <div ref="chartRef" style="width: 100%; height: 340px;"></div>
        </div>

        <!-- 汇总表格 -->
        <div class="gva-search-box">
          <el-form :inline="true" :model="summarySearch" @keyup.enter="getSummaryList">
            <el-form-item label="日期范围">
              <el-date-picker v-model="summarySearch.startDate" type="date" value-format="YYYY-MM-DD" placeholder="开始日期" />
              —
              <el-date-picker v-model="summarySearch.endDate" type="date" value-format="YYYY-MM-DD" placeholder="结束日期" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="getSummaryList">查询</el-button>
              <el-button icon="refresh" @click="resetSummarySearch">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <el-table :data="summaryList" style="width: 100%" tooltip-effect="dark">
          <el-table-column align="left" label="日期" prop="date" width="150" />
          <el-table-column align="left" label="PV (页面浏览)" prop="pv" width="150" />
          <el-table-column align="left" label="UV (独立访客)" prop="uv" width="150" />
          <el-table-column align="left" label="新访客" prop="newVisitor" width="150" />
          <el-table-column align="left" label="已登录用户" prop="registeredUser" min-width="150" />
        </el-table>
        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="summarySearch.page"
            :page-size="summarySearch.pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="summaryTotal"
            @current-change="(val) => { summarySearch.page = val; getSummaryList() }"
            @size-change="(val) => { summarySearch.pageSize = val; getSummaryList() }"
          />
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getVisitorLogList, getVisitorSummaryList, getTodayStats, aggregateDailySummary } from '@/api/client/visitor'
import { formatDate } from '@/utils/format'
import * as echarts from 'echarts'
import { useAppStore } from '@/pinia'

defineOptions({
  name: 'Visitor'
})

const appStore = useAppStore()
const activeTab = ref('log')

// ===== 今日统计 =====
const todayStats = ref({})
const fetchTodayStats = async () => {
  const res = await getTodayStats()
  if (res.code === 0) {
    todayStats.value = res.data
  }
}

// ===== 访客日志 =====
const logSearch = ref({
  page: 1,
  pageSize: 10,
  startDate: '',
  endDate: '',
  visitorId: '',
  ip: '',
  platform: ''
})
const logList = ref([])
const logTotal = ref(0)

const getLogList = async () => {
  const res = await getVisitorLogList(logSearch.value)
  if (res.code === 0) {
    logList.value = res.data.list || []
    logTotal.value = res.data.total
  }
}

const resetLogSearch = () => {
  logSearch.value = { page: 1, pageSize: 10, startDate: '', endDate: '', visitorId: '', ip: '', platform: '' }
  getLogList()
}

// ===== 汇总统计 =====
const summarySearch = ref({
  page: 1,
  pageSize: 10,
  startDate: '',
  endDate: ''
})
const summaryList = ref([])
const summaryTotal = ref(0)

const getSummaryList = async () => {
  const res = await getVisitorSummaryList(summarySearch.value)
  if (res.code === 0) {
    summaryList.value = res.data.list || []
    summaryTotal.value = res.data.total
  }
}

const resetSummarySearch = () => {
  summarySearch.value = { page: 1, pageSize: 10, startDate: '', endDate: '' }
  getSummaryList()
}

// ===== 图表 =====
const chartDays = ref(7)
const chartData = ref([])
const chartRef = ref(null)
let chartInstance = null

const fetchChartData = async () => {
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - chartDays.value + 1)
  const startDate = start.toISOString().slice(0, 10)
  const endDate = end.toISOString().slice(0, 10)
  const res = await getVisitorSummaryList({ page: 1, pageSize: chartDays.value, startDate, endDate })
  if (res.code === 0) {
    chartData.value = (res.data.list || []).sort((a, b) => a.date.localeCompare(b.date))
    updateChart()
  }
}

const getChartOption = () => {
  const isDark = appStore.isDark
  const textColor = isDark ? 'rgba(255,255,255,0.70)' : 'rgba(0,0,0,0.70)'
  const subtextColor = isDark ? 'rgba(255,255,255,0.45)' : 'rgba(0,0,0,0.45)'
  const borderColor = isDark ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)'
  const bgColor = isDark ? '#141414' : '#ffffff'
  const dates = chartData.value.map(item => item.date.slice(5)) // MM-DD 更简洁
  const pvData = chartData.value.map(item => item.pv)
  const uvData = chartData.value.map(item => item.uv)
  const newVisitorData = chartData.value.map(item => item.newVisitor)
  const registeredData = chartData.value.map(item => item.registeredUser)

  return {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: isDark ? 'rgba(30,30,30,0.95)' : 'rgba(255,255,255,0.95)',
      borderColor: isDark ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.1)',
      textStyle: { color: isDark ? '#fff' : '#333', fontSize: 13 },
      padding: [12, 16],
      extraCssText: 'border-radius: 8px; box-shadow: 0 4px 20px rgba(0,0,0,0.15);',
      axisPointer: {
        type: 'cross',
        crossStyle: { color: subtextColor },
        lineStyle: { color: borderColor, type: 'dashed' }
      }
    },
    legend: {
      data: ['PV', 'UV', '新访客', '已登录用户'],
      top: 0,
      right: 0,
      itemWidth: 16,
      itemHeight: 8,
      itemGap: 20,
      textStyle: { color: textColor, fontSize: 12 },
      icon: 'roundRect'
    },
    grid: { left: 55, right: 20, top: 45, bottom: 35, containLabel: false },
    xAxis: {
      type: 'category',
      data: dates,
      boundaryGap: false,
      axisLabel: { color: subtextColor, fontSize: 11, margin: 12 },
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: subtextColor, fontSize: 11, margin: 12 },
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: borderColor, type: 'dashed' } },
      splitNumber: 4
    },
    series: [
      {
        name: 'PV', type: 'line', smooth: 0.4, data: pvData,
        symbol: 'circle', symbolSize: 6, showSymbol: false,
        emphasis: { scale: true, itemStyle: { borderWidth: 2, borderColor: '#fff' } },
        lineStyle: { width: 2.5, color: '#409EFF' },
        itemStyle: { color: '#409EFF' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64,158,255,0.20)' },
            { offset: 0.7, color: 'rgba(64,158,255,0.05)' },
            { offset: 1, color: 'rgba(64,158,255,0)' }
          ])
        }
      },
      {
        name: 'UV', type: 'line', smooth: 0.4, data: uvData,
        symbol: 'circle', symbolSize: 6, showSymbol: false,
        emphasis: { scale: true, itemStyle: { borderWidth: 2, borderColor: '#fff' } },
        lineStyle: { width: 2.5, color: '#67C23A' },
        itemStyle: { color: '#67C23A' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(103,194,58,0.15)' },
            { offset: 1, color: 'rgba(103,194,58,0)' }
          ])
        }
      },
      {
        name: '新访客', type: 'line', smooth: 0.4, data: newVisitorData,
        symbol: 'diamond', symbolSize: 6, showSymbol: false,
        emphasis: { scale: true, itemStyle: { borderWidth: 2, borderColor: '#fff' } },
        lineStyle: { width: 2, color: '#E6A23C', type: 'dashed' },
        itemStyle: { color: '#E6A23C' }
      },
      {
        name: '已登录用户', type: 'bar', data: registeredData,
        barMaxWidth: 24, barMinHeight: 2,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(245,108,108,0.85)' },
            { offset: 1, color: 'rgba(245,108,108,0.35)' }
          ]),
          borderRadius: [4, 4, 0, 0]
        },
        emphasis: {
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(245,108,108,1)' },
              { offset: 1, color: 'rgba(245,108,108,0.6)' }
            ])
          }
        }
      }
    ],
    animationDuration: 800,
    animationEasing: 'cubicInOut'
  }
}

const updateChart = () => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  if (chartRef.value) {
    chartInstance = echarts.init(chartRef.value, appStore.isDark ? 'dark' : null)
    chartInstance.setOption(getChartOption())
  }
}

const handleResize = () => {
  chartInstance && chartInstance.resize()
}

// 监听 tab 切换，在汇总 tab 激活后初始化/刷新图表
watch(activeTab, (val) => {
  if (val === 'summary') {
    nextTick(() => {
      updateChart()
    })
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})

// ===== 手动聚合 =====
const handleAggregate = async () => {
  const today = new Date().toISOString().slice(0, 10)
  const res = await aggregateDailySummary(today)
  if (res.code === 0) {
    ElMessage.success('聚合成功')
    getSummaryList()
    fetchTodayStats()
    fetchChartData()
  }
}

onMounted(() => {
  fetchTodayStats()
  getLogList()
  getSummaryList()
  fetchChartData()
  window.addEventListener('resize', handleResize)
})
</script>
