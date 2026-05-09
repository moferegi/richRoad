<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="createdAtRange">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="w-[380px]"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>
        <el-form-item label="用户" prop="userId">
          <el-select
            v-model="searchInfo.userId"
            clearable
            filterable
            remote
            reserve-keyword
            placeholder="请选择用户"
            :remote-method="remoteUserSearch"
            :loading="userLoading"
            @clear="() => { searchInfo.userId = undefined }"
          >
            <el-option
              v-for="item in userOptions"
              :key="item.ID"
              :label="`${item.nickname || item.username} (ID: ${item.ID})`"
              :value="item.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="增/减" prop="changeType">
          <el-select v-model="searchInfo.changeType" clearable placeholder="请选择" @clear="() => { searchInfo.changeType = undefined }">
            <el-option label="增加" value="increase" />
            <el-option label="减少" value="decrease" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作类型" prop="operationType">
          <el-input v-model="searchInfo.operationType" clearable placeholder="如 tryon_consume / tryon_beautify_consume" />
        </el-form-item>
        <el-form-item label="原因" prop="reason">
          <el-input v-model="searchInfo.reason" clearable placeholder="搜索原因" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="record-summary">
        <el-tag type="warning">试衣币</el-tag>
        <span>仅展示资产类型为 tryon_point 的流水和统计。</span>
      </div>

      <el-tabs v-model="activeTab" class="panel-tabs" @tab-change="handleTabChange">
        <el-tab-pane label="流水列表" name="list" />
        <el-tab-pane label="统计分析" name="stats" />
      </el-tabs>

      <template v-if="activeTab === 'list'">
        <el-table :data="tableData" row-key="ID" :default-sort="{ prop: 'CreatedAt', order: 'descending' }" @sort-change="sortChange">
          <el-table-column sortable="custom" align="left" label="日期" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
          </el-table-column>
          <el-table-column align="left" label="用户" prop="userId" min-width="160">
            <template #default="scope">
              <span v-if="getUserDisplayName(scope.row.userId)">{{ getUserDisplayName(scope.row.userId) }}</span>
              <span v-else>ID: {{ scope.row.userId }}</span>
            </template>
          </el-table-column>
          <el-table-column align="left" label="变化" prop="pointChange" width="110" sortable="custom">
            <template #default="scope">
              <el-tag :type="scope.row.pointChange > 0 ? 'success' : 'danger'">
                {{ scope.row.pointChange > 0 ? '+' : '' }}{{ scope.row.pointChange }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column align="left" label="当前余额" prop="currentPoints" width="110" />
          <el-table-column align="left" label="操作类型" prop="operationType" width="180" />
          <el-table-column align="left" label="变化原因" prop="reason" min-width="180" show-overflow-tooltip />
          <el-table-column align="left" label="备注" prop="remark" min-width="180" show-overflow-tooltip />
          <el-table-column align="left" label="关联订单" prop="relatedOrderId" width="100" />
        </el-table>
        <div class="gva-pagination">
          <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
          />
        </div>
      </template>

      <template v-else>
        <el-skeleton v-if="statsLoading" :rows="6" animated />
        <template v-else>
          <div class="stats-cards">
            <el-card shadow="hover" class="stats-card"><div class="stats-label">注册赠送</div><div class="stats-value">{{ stats.registerRewardTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">邀请赠送</div><div class="stats-value">{{ stats.inviteRewardTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">充值获得</div><div class="stats-value">{{ stats.rechargeTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">后台增加</div><div class="stats-value">{{ stats.adminIncreaseTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">后台减少</div><div class="stats-value">{{ stats.adminDecreaseTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">总试衣币(获得)</div><div class="stats-value">{{ stats.totalGranted }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">已使用试衣币</div><div class="stats-value">{{ stats.totalUsed }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">模型调用总次数</div><div class="stats-value">{{ stats.modelCallTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">模型总消耗试衣币</div><div class="stats-value">{{ stats.modelCostTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">美肤使用次数</div><div class="stats-value">{{ stats.beautifyUsedTotal }}</div></el-card>
            <el-card shadow="hover" class="stats-card"><div class="stats-label">美肤总消耗试衣币</div><div class="stats-value">{{ stats.beautifyCostTotal }}</div></el-card>
          </div>

          <div class="chart-grid">
            <el-card shadow="never">
              <template #header>
                <div class="chart-title">每个模型调用次数</div>
              </template>
              <div ref="modelCallChartRef" class="chart-container" />
            </el-card>
            <el-card shadow="never">
              <template #header>
                <div class="chart-title">模型消耗试衣币占比</div>
              </template>
              <div ref="modelCostChartRef" class="chart-container" />
            </el-card>
          </div>

          <el-card class="model-stats-wrap" shadow="never">
            <template #header>
              <div class="model-stats-title">模型统计明细</div>
            </template>
            <el-table :data="stats.modelStats" border stripe>
              <el-table-column label="模型" prop="modelKey" min-width="180" />
              <el-table-column label="调用次数" prop="taskCount" width="130" />
              <el-table-column label="开启精修次数" prop="refinerEnabledCount" width="150" />
              <el-table-column label="美肤使用次数" prop="beautifyUsedCount" width="150" />
              <el-table-column label="总消耗试衣币" prop="totalCostPoints" width="150" />
              <el-table-column label="美肤消耗试衣币" prop="beautifyCostPoints" width="150" />
            </el-table>
          </el-card>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'
import { useAppStore } from '@/pinia'
import { getPointRecordList, getPointRecordStats } from '@/api/client/pointRecord'
import { getClientUserList } from '@/api/client/user'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'TryonPointRecord'
})

const appStore = useAppStore()
const activeTab = ref('list')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const statsLoading = ref(false)
const stats = ref({
  registerRewardTotal: 0,
  inviteRewardTotal: 0,
  rechargeTotal: 0,
  adminIncreaseTotal: 0,
  adminDecreaseTotal: 0,
  totalGranted: 0,
  totalUsed: 0,
  beautifyUsedTotal: 0,
  modelCallTotal: 0,
  modelCostTotal: 0,
  beautifyCostTotal: 0,
  modelStats: [],
})
const searchInfo = ref({ assetType: 'tryon_point', sort: 'created_at', order: 'descending' })
const userLoading = ref(false)
const userOptions = ref([])
const userCache = ref(new Map())
const modelCallChartRef = ref(null)
const modelCostChartRef = ref(null)

let modelCallChartInstance = null
let modelCostChartInstance = null

const formatModelKey = (value) => {
  const key = String(value || '').trim()
  return key || 'default'
}

const buildStatsSearchParams = () => {
  return {
    createdAtRange: searchInfo.value.createdAtRange,
    userId: searchInfo.value.userId,
  }
}

const sortChange = ({ prop, order }) => {
  const sortMap = { CreatedAt: 'created_at', pointChange: 'point_change', ID: 'id' }
  searchInfo.value.sort = sortMap[prop] || prop
  searchInfo.value.order = order || 'descending'
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  Promise.all([getTableData(), getStatsData()])
}

const onReset = () => {
  page.value = 1
  searchInfo.value = { assetType: 'tryon_point', sort: 'created_at', order: 'descending' }
  Promise.all([getTableData(), getStatsData()])
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const table = await getPointRecordList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value, assetType: 'tryon_point' })
  if (table.code === 0) {
    tableData.value = table.data.list || []
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
    const userIds = tableData.value.map(item => item.userId).filter(Boolean)
    if (userIds.length > 0) await loadUsersInfo(userIds)
  }
}

const getStatsData = async () => {
  statsLoading.value = true
  try {
    const res = await getPointRecordStats(buildStatsSearchParams())
    if (res.code === 0) {
      const modelStats = Array.isArray(res.data?.modelStats)
        ? res.data.modelStats.map(item => ({
          modelKey: formatModelKey(item.modelKey),
          taskCount: Number(item.taskCount || 0),
          refinerEnabledCount: Number(item.refinerEnabledCount || 0),
          beautifyUsedCount: Number(item.beautifyUsedCount || 0),
          totalCostPoints: Number(item.totalCostPoints || 0),
          beautifyCostPoints: Number(item.beautifyCostPoints || 0),
        }))
        : []
      stats.value = {
        registerRewardTotal: Number(res.data?.registerRewardTotal || 0),
        inviteRewardTotal: Number(res.data?.inviteRewardTotal || 0),
        rechargeTotal: Number(res.data?.rechargeTotal || 0),
        adminIncreaseTotal: Number(res.data?.adminIncreaseTotal || 0),
        adminDecreaseTotal: Number(res.data?.adminDecreaseTotal || 0),
        totalGranted: Number(res.data?.totalGranted || 0),
        totalUsed: Number(res.data?.totalUsed || 0),
        beautifyUsedTotal: Number(res.data?.beautifyUsedTotal || 0),
        modelCallTotal: Number(res.data?.modelCallTotal || 0),
        modelCostTotal: Number(res.data?.modelCostTotal || 0),
        beautifyCostTotal: Number(res.data?.beautifyCostTotal || 0),
        modelStats,
      }
      await nextTick()
      renderCharts()
    }
  } finally {
    statsLoading.value = false
  }
}

const getModelCallChartOption = () => {
  const isDark = appStore.isDark
  const textColor = isDark ? 'rgba(255,255,255,0.7)' : 'rgba(0,0,0,0.7)'
  const xData = stats.value.modelStats.map(item => item.modelKey)
  const yData = stats.value.modelStats.map(item => item.taskCount)
  return {
    backgroundColor: 'transparent',
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: 16, right: 16, bottom: 8, top: 16, containLabel: true },
    xAxis: {
      type: 'category',
      data: xData,
      axisLabel: { color: textColor, rotate: xData.length > 4 ? 18 : 0 },
      axisLine: { lineStyle: { color: isDark ? '#4b5563' : '#dcdfe6' } },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: textColor },
      splitLine: { lineStyle: { color: isDark ? 'rgba(255,255,255,0.08)' : '#ebeef5' } },
    },
    series: [{
      type: 'bar',
      data: yData,
      barMaxWidth: 44,
      itemStyle: {
        color: '#409eff',
        borderRadius: [4, 4, 0, 0],
      },
      label: {
        show: true,
        position: 'top',
        formatter: '{c}',
        color: textColor,
      },
    }],
  }
}

const getModelCostChartOption = () => {
  const isDark = appStore.isDark
  const textColor = isDark ? 'rgba(255,255,255,0.7)' : 'rgba(0,0,0,0.7)'
  return {
    backgroundColor: 'transparent',
    tooltip: { trigger: 'item' },
    legend: {
      orient: 'vertical',
      right: 0,
      top: 'middle',
      textStyle: { color: textColor, fontSize: 12 },
    },
    series: [
      {
        name: '模型扣币',
        type: 'pie',
        radius: ['42%', '68%'],
        center: ['36%', '50%'],
        avoidLabelOverlap: true,
        itemStyle: {
          borderRadius: 6,
          borderColor: isDark ? '#1f2937' : '#fff',
          borderWidth: 2,
        },
        label: { show: false },
        emphasis: {
          label: {
            show: true,
            formatter: '{b}\n{c}币',
            fontSize: 12,
          },
        },
        data: stats.value.modelStats.map(item => ({
          name: item.modelKey,
          value: item.totalCostPoints,
        })),
      },
    ],
  }
}

const renderCharts = () => {
  if (activeTab.value !== 'stats') {
    return
  }

  if (modelCallChartRef.value) {
    if (modelCallChartInstance) {
      modelCallChartInstance.dispose()
      modelCallChartInstance = null
    }
    modelCallChartInstance = echarts.init(modelCallChartRef.value, appStore.isDark ? 'dark' : null)
    modelCallChartInstance.setOption(getModelCallChartOption())
  }

  if (modelCostChartRef.value) {
    if (modelCostChartInstance) {
      modelCostChartInstance.dispose()
      modelCostChartInstance = null
    }
    modelCostChartInstance = echarts.init(modelCostChartRef.value, appStore.isDark ? 'dark' : null)
    modelCostChartInstance.setOption(getModelCostChartOption())
  }
}

const handleTabChange = async (tabName) => {
  if (tabName === 'stats') {
    await getStatsData()
  }
}

const handleResize = () => {
  if (modelCallChartInstance) {
    modelCallChartInstance.resize()
  }
  if (modelCostChartInstance) {
    modelCostChartInstance.resize()
  }
}

const remoteUserSearch = async (query) => {
  if (!query) {
    userOptions.value = Array.from(userCache.value.values())
    return
  }
  const text = String(query).toLowerCase()
  userOptions.value = Array.from(userCache.value.values()).filter(user => {
    return String(user.nickname || '').toLowerCase().includes(text) || String(user.username || '').toLowerCase().includes(text) || String(user.ID).includes(text)
  })
}

const initUserOptions = async () => {
  userLoading.value = true
  const res = await getClientUserList({ page: 1, pageSize: 1000 })
  userLoading.value = false
  if (res.code === 0) {
    userOptions.value = res.data.list || []
    userOptions.value.forEach(user => userCache.value.set(user.ID, user))
  }
}

const loadUsersInfo = async (userIds) => {
  const missing = [...new Set(userIds)].filter(id => id && !userCache.value.has(id))
  if (missing.length === 0) return
  const res = await getClientUserList({ page: 1, pageSize: 100 })
  if (res.code === 0 && res.data.list) {
    res.data.list.forEach(user => userCache.value.set(user.ID, user))
  }
}

const getUserDisplayName = (userId) => {
  const user = userCache.value.get(userId)
  return user ? `${user.nickname || user.username} (ID: ${user.ID})` : ''
}

watch(() => appStore.isDark, async () => {
  await nextTick()
  renderCharts()
})

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  await initUserOptions()
  await Promise.all([getTableData(), getStatsData()])
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (modelCallChartInstance) {
    modelCallChartInstance.dispose()
    modelCallChartInstance = null
  }
  if (modelCostChartInstance) {
    modelCostChartInstance.dispose()
    modelCostChartInstance = null
  }
})
</script>

<style scoped>
.record-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  color: var(--el-text-color-regular);
}

.panel-tabs {
  margin-bottom: 8px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}

.stats-card {
  min-height: 94px;
}

.stats-label {
  font-size: 13px;
  color: var(--el-text-color-secondary);
}

.stats-value {
  margin-top: 8px;
  font-size: 30px;
  line-height: 1;
  font-weight: 600;
  color: var(--el-color-primary);
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.chart-title {
  font-size: 14px;
  font-weight: 600;
}

.chart-container {
  width: 100%;
  height: 320px;
}

.model-stats-wrap {
  margin-top: 16px;
}

.model-stats-title {
  font-size: 14px;
  font-weight: 600;
}

@media (max-width: 1280px) {
  .stats-cards {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .chart-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }
}
</style>
