<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="用户ID">
          <el-input v-model.number="searchInfo.userID" clearable placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="任务编号">
          <el-input v-model="searchInfo.taskNo" clearable placeholder="请输入任务编号" />
        </el-form-item>
        <el-form-item label="请求ID">
          <el-input v-model="searchInfo.requestID" clearable placeholder="请输入请求ID" />
        </el-form-item>
        <el-form-item label="创建日期">
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="开始日期"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="结束日期"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"
          />
        </el-form-item>
        <el-form-item label="场景">
          <el-select v-model="searchInfo.sceneType" clearable placeholder="请选择场景">
            <el-option label="试衣" value="clothes" />
            <el-option label="试鞋" value="shoes" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择状态">
            <el-option label="处理中" value="processing" />
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-tabs v-model="activeTab" class="panel-tabs" @tab-change="handleTabChange">
      <el-tab-pane label="任务列表" name="list">
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button type="primary" icon="download" :loading="exporting" @click="handleExportCsv">
              导出CSV
            </el-button>
            <el-button
              type="danger"
              icon="delete"
              :disabled="!multipleSelection.length"
              @click="handleBatchDelete"
            >
              批量删除
            </el-button>
          </div>

          <el-table
            :data="tableData"
            style="width: 100%"
            tooltip-effect="dark"
            row-key="ID"
            :row-class-name="tableRowClassName"
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="50" />
            <el-table-column align="left" label="ID" prop="ID" width="80" />
            <el-table-column align="left" label="用户ID" prop="userID" width="90" />
            <el-table-column align="left" label="任务编号" min-width="220" show-overflow-tooltip>
              <template #default="scope">
                <div class="copy-cell">
                  <span class="copy-text">{{ scope.row.taskNo || '-' }}</span>
                  <el-button
                    v-if="scope.row.taskNo"
                    type="primary"
                    link
                    size="small"
                    @click="copyText(scope.row.taskNo, '任务编号')"
                  >
                    复制
                  </el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column align="left" label="请求ID" min-width="220" show-overflow-tooltip>
              <template #default="scope">
                <div class="copy-cell">
                  <span class="copy-text">{{ scope.row.requestID || '-' }}</span>
                  <el-button
                    v-if="scope.row.requestID"
                    type="primary"
                    link
                    size="small"
                    @click="copyText(scope.row.requestID, '请求ID')"
                  >
                    复制
                  </el-button>
                </div>
              </template>
            </el-table-column>
            <el-table-column align="left" label="场景" width="100">
              <template #default="scope">
                {{ sceneLabel(scope.row.sceneType) }}
              </template>
            </el-table-column>
            <el-table-column align="left" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="statusTagType(scope.row.status)">
                  {{ statusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column align="left" label="模型" min-width="140" show-overflow-tooltip>
              <template #default="scope">
                <span>{{ formatModelKey(scope.row.provider) }}</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="精修" width="110">
              <template #default="scope">
                <el-tag size="small" :type="refinerTagType(scope.row)">{{ refinerLabel(scope.row) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column align="left" label="扣币" prop="costPoints" width="80" />
            <el-table-column align="left" label="退币" prop="refundPoints" width="80" />
            <el-table-column align="left" label="原图" width="90">
              <template #default="scope">
                <el-image
                  v-if="scope.row.sourceImage"
                  class="task-image"
                  :src="getUrl(scope.row.sourceImage)"
                  :preview-src-list="[getUrl(scope.row.sourceImage)]"
                  preview-teleported
                  fit="cover"
                />
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="模板图" width="90">
              <template #default="scope">
                <el-image
                  v-if="scope.row.templateImage"
                  class="task-image"
                  :src="getUrl(scope.row.templateImage)"
                  :preview-src-list="[getUrl(scope.row.templateImage)]"
                  preview-teleported
                  fit="cover"
                />
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="结果图" width="90">
              <template #default="scope">
                <el-image
                  v-if="scope.row.resultImage"
                  class="task-image"
                  :src="getUrl(scope.row.resultImage)"
                  :preview-src-list="[getUrl(scope.row.resultImage)]"
                  preview-teleported
                  fit="cover"
                />
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="失败原因" min-width="200" show-overflow-tooltip>
              <template #default="scope">
                <span :class="{ 'error-message': scope.row.status === 'failed' && scope.row.errorMessage }">
                  {{ scope.row.errorMessage || '-' }}
                </span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="创建时间" width="170">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column align="left" label="完成时间" width="170">
              <template #default="scope">{{ formatDate(scope.row.completedAt) || '-' }}</template>
            </el-table-column>
            <el-table-column align="left" label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button type="primary" link @click="openDetail(scope.row)">详情</el-button>
                <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
              </template>
            </el-table-column>
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
        </div>
      </el-tab-pane>

      <el-tab-pane label="统计分析" name="stats">
        <div class="gva-search-box stats-wrap" v-loading="statsLoading">
          <el-row :gutter="16">
            <el-col :xs="24" :sm="12" :md="4">
              <el-card class="stats-card stats-card-clickable" :class="{ active: !searchInfo.status }" shadow="hover" @click="handleStatusCardClick('')">
                <el-statistic title="任务总数" :value="stats.total" />
              </el-card>
            </el-col>
            <el-col :xs="24" :sm="12" :md="4">
              <el-card class="stats-card stats-card-clickable" :class="{ active: searchInfo.status === 'processing' }" shadow="hover" @click="handleStatusCardClick('processing')">
                <el-statistic title="处理中" :value="stats.processing" />
              </el-card>
            </el-col>
            <el-col :xs="24" :sm="12" :md="4">
              <el-card class="stats-card stats-card-clickable" :class="{ active: searchInfo.status === 'success' }" shadow="hover" @click="handleStatusCardClick('success')">
                <el-statistic title="成功" :value="stats.success" />
              </el-card>
            </el-col>
            <el-col :xs="24" :sm="12" :md="4">
              <el-card class="stats-card stats-card-clickable" :class="{ active: searchInfo.status === 'failed' }" shadow="hover" @click="handleStatusCardClick('failed')">
                <el-statistic title="失败" :value="stats.failed" />
              </el-card>
            </el-col>
            <el-col :xs="24" :sm="12" :md="4">
              <el-card class="stats-card" shadow="hover">
                <el-statistic title="精修任务数" :value="stats.refinerEnabledCount" />
              </el-card>
            </el-col>
            <el-col :xs="24" :sm="12" :md="4">
              <el-card class="stats-card" shadow="hover">
                <el-statistic title="总扣币" :value="stats.totalCostPoints" />
              </el-card>
            </el-col>
          </el-row>
        </div>

        <div class="chart-grid">
          <div class="gva-search-box trend-wrap" v-loading="chartLoading">
            <div class="trend-header">
              <span class="trend-title">任务趋势</span>
              <el-radio-group v-model="chartDays" size="small" @change="fetchTrendData">
                <el-radio-button :value="7">近7天</el-radio-button>
                <el-radio-button :value="14">近14天</el-radio-button>
                <el-radio-button :value="30">近30天</el-radio-button>
              </el-radio-group>
            </div>
            <div ref="trendChartRef" class="trend-chart" />
          </div>

          <div class="gva-search-box trend-wrap" v-loading="statsLoading">
            <div class="trend-header">
              <span class="trend-title">模型调用占比</span>
              <span class="trend-subtitle">总调用 {{ stats.modelCallTotal }}</span>
            </div>
            <div ref="modelChartRef" class="trend-chart" />
          </div>
        </div>

        <div class="gva-table-box model-stats-wrap">
          <div class="model-stats-title">模型调用明细</div>
          <el-table :data="stats.modelStats || []" row-key="modelKey" empty-text="暂无模型统计数据">
            <el-table-column align="left" label="模型" prop="modelKey" min-width="180" show-overflow-tooltip>
              <template #default="scope">{{ formatModelKey(scope.row.modelKey) }}</template>
            </el-table-column>
            <el-table-column align="left" label="调用次数" prop="taskCount" width="110" />
            <el-table-column align="left" label="成功" prop="successCount" width="90" />
            <el-table-column align="left" label="失败" prop="failedCount" width="90" />
            <el-table-column align="left" label="处理中" prop="processingCount" width="90" />
            <el-table-column align="left" label="精修次数" prop="refinerEnabledCount" width="100" />
            <el-table-column align="left" label="总扣币" prop="totalCostPoints" width="110" />
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-drawer v-model="detailVisible" title="试衣任务详情" size="640px" destroy-on-close>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="ID">{{ detailRow.ID || '-' }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ detailRow.userID || '-' }}</el-descriptions-item>
        <el-descriptions-item label="任务编号">
          <div class="copy-cell">
            <span class="copy-text">{{ detailRow.taskNo || '-' }}</span>
            <el-button v-if="detailRow.taskNo" type="primary" link size="small" @click="copyText(detailRow.taskNo, '任务编号')">
              复制
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="请求ID">
          <div class="copy-cell">
            <span class="copy-text">{{ detailRow.requestID || '-' }}</span>
            <el-button
              v-if="detailRow.requestID"
              type="primary"
              link
              size="small"
              @click="copyText(detailRow.requestID, '请求ID')"
            >
              复制
            </el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="场景">{{ sceneLabel(detailRow.sceneType) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTagType(detailRow.status)">{{ statusLabel(detailRow.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="模型">{{ formatModelKey(detailRow.provider) }}</el-descriptions-item>
        <el-descriptions-item label="精修">
          <el-tag size="small" :type="refinerTagType(detailRow)">{{ refinerLabel(detailRow) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="扣币">{{ detailRow.costPoints ?? '-' }}</el-descriptions-item>
        <el-descriptions-item label="退币">{{ detailRow.refundPoints ?? '-' }}</el-descriptions-item>
        <el-descriptions-item label="失败原因">{{ detailRow.errorMessage || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ formatDate(detailRow.CreatedAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ formatDate(detailRow.completedAt) || '-' }}</el-descriptions-item>
        <el-descriptions-item label="原图">
          <el-image
            v-if="detailRow.sourceImage"
            class="detail-image"
            :src="getUrl(detailRow.sourceImage)"
            :preview-src-list="[getUrl(detailRow.sourceImage)]"
            preview-teleported
            fit="cover"
          />
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="模板图">
          <el-image
            v-if="detailRow.templateImage"
            class="detail-image"
            :src="getUrl(detailRow.templateImage)"
            :preview-src-list="[getUrl(detailRow.templateImage)]"
            preview-teleported
            fit="cover"
          />
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="结果图">
          <el-image
            v-if="detailRow.resultImage"
            class="detail-image"
            :src="getUrl(detailRow.resultImage)"
            :preview-src-list="[getUrl(detailRow.resultImage)]"
            preview-teleported
            fit="cover"
          />
          <span v-else>-</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import { useAppStore } from '@/pinia'
import {
  getTryonTaskList,
  getTryonTaskStats,
  getTryonTaskTrend,
  deleteTryonTask,
  deleteTryonTaskByIds,
} from '@/api/client/tryonTask'
import { formatDate } from '@/utils/format'
import { getUrl } from '@/utils/image'

defineOptions({
  name: 'TryonTaskManage'
})

const appStore = useAppStore()

const activeTab = ref('list')
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const multipleSelection = ref([])
const exporting = ref(false)
const statsLoading = ref(false)
const detailVisible = ref(false)
const detailRow = ref({})
const stats = ref({
  total: 0,
  processing: 0,
  success: 0,
  failed: 0,
  totalCostPoints: 0,
  refinerEnabledCount: 0,
  modelCallTotal: 0,
  modelStats: [],
})
const chartDays = ref(7)
const chartLoading = ref(false)
const trendChartRef = ref(null)
const modelChartRef = ref(null)
const chartData = ref([])
const searchInfo = ref({
  startCreatedAt: undefined,
  endCreatedAt: undefined,
  userID: undefined,
  taskNo: '',
  requestID: '',
  sceneType: '',
  status: ''
})

const MAX_EXPORT_ROWS = 5000
let trendChartInstance = null
let modelChartInstance = null

const statusLabel = (status) => {
  if (status === 'processing') return '处理中'
  if (status === 'success') return '成功'
  if (status === 'failed') return '失败'
  return status || '-'
}

const statusTagType = (status) => {
  if (status === 'processing') return 'warning'
  if (status === 'success') return 'success'
  if (status === 'failed') return 'danger'
  return 'info'
}

const sceneLabel = (sceneType) => {
  if (sceneType === 'clothes') return '试衣'
  if (sceneType === 'shoes') return '试鞋'
  return sceneType || '-'
}

const formatModelKey = (modelKey) => {
  const key = String(modelKey || '').trim()
  return key || 'default'
}

const isAliyunProvider = (provider) => {
  const value = String(provider || '').trim().toLowerCase()
  if (!value) return false
  return value.includes('aliyun') || value.includes('dashscope') || value.includes('aitryon')
}

const refinerLabel = (row) => {
  if (!isAliyunProvider(row?.provider)) {
    return '不支持'
  }
  return row?.enableRefiner ? '已开启' : '未开启'
}

const refinerTagType = (row) => {
  if (!isAliyunProvider(row?.provider)) {
    return 'info'
  }
  return row?.enableRefiner ? 'success' : 'warning'
}

const handleStatusCardClick = (status) => {
  if (searchInfo.value.status === status) {
    return
  }
  searchInfo.value.status = status
  activeTab.value = 'list'
  page.value = 1
  getTableData()
}

const copyText = async (text, label = '内容') => {
  if (!text) {
    ElMessage.warning(`${label}为空`)
    return
  }
  const content = String(text)
  if (navigator.clipboard && window.isSecureContext) {
    try {
      await navigator.clipboard.writeText(content)
      ElMessage.success(`${label}已复制`)
      return
    } catch (e) {
      // continue to fallback copy
    }
  }
  const textarea = document.createElement('textarea')
  textarea.value = content
  textarea.style.position = 'fixed'
  textarea.style.top = '-9999px'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
  ElMessage.success(`${label}已复制`)
}

const openDetail = (row) => {
  detailRow.value = { ...row }
  detailVisible.value = true
}

const tableRowClassName = ({ row }) => {
  if (row.status !== 'processing') {
    return 'tryon-row-non-processing'
  }
  return ''
}

const handleSelectionChange = (rows) => {
  multipleSelection.value = rows || []
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除任务【${row.taskNo || row.ID}】吗？`, '提示', { type: 'warning' })
    const res = await deleteTryonTask({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value -= 1
      }
      await getTableData()
    }
  } catch (e) {
    // 用户取消删除时不做提示
  }
}

const handleBatchDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请先选择要删除的数据')
    return
  }

  try {
    await ElMessageBox.confirm(`确定批量删除选中的${multipleSelection.value.length}条任务吗？`, '提示', { type: 'warning' })
    const ids = multipleSelection.value.map(item => item.ID).filter(Boolean)
    const res = await deleteTryonTaskByIds({ 'IDs[]': ids })
    if (res.code === 0) {
      ElMessage.success('批量删除成功')
      if (ids.length >= tableData.value.length && page.value > 1) {
        page.value -= 1
      }
      multipleSelection.value = []
      await getTableData()
    }
  } catch (e) {
    // 用户取消删除时不做提示
  }
}

const escapeCsvValue = (value) => {
  if (value === null || value === undefined) return ''
  const text = String(value).replace(/"/g, '""')
  return /[",\n]/.test(text) ? `"${text}"` : text
}

const buildCsvContent = (list) => {
  const headers = [
    'ID',
    '用户ID',
    '任务编号',
    '请求ID',
    '场景',
    '状态',
    '模型',
    '精修',
    '扣币',
    '退币',
    '原图',
    '模板图',
    '结果图',
    '失败原因',
    '创建时间',
    '完成时间',
  ]
  const rows = list.map((item) => {
    return [
      item.ID,
      item.userID,
      item.taskNo,
      item.requestID,
      sceneLabel(item.sceneType),
      statusLabel(item.status),
      formatModelKey(item.provider),
      refinerLabel(item),
      item.costPoints,
      item.refundPoints,
      item.sourceImage,
      item.templateImage,
      item.resultImage,
      item.errorMessage,
      formatDate(item.CreatedAt) || '',
      formatDate(item.completedAt) || '',
    ]
  })
  const lines = [headers, ...rows].map((row) => row.map(escapeCsvValue).join(','))
  return `\uFEFF${lines.join('\n')}`
}

const downloadCsv = (content, fileName) => {
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

const getExportFileName = () => {
  const now = new Date()
  const pad = (num) => String(num).padStart(2, '0')
  return `tryon_task_${now.getFullYear()}${pad(now.getMonth() + 1)}${pad(now.getDate())}_${pad(now.getHours())}${pad(now.getMinutes())}${pad(now.getSeconds())}.csv`
}

const buildSearchParams = () => {
  return {
    ...searchInfo.value,
  }
}

const getTrendChartOption = () => {
  const isDark = appStore.isDark
  const textColor = isDark ? 'rgba(255,255,255,0.70)' : 'rgba(0,0,0,0.70)'
  const subtextColor = isDark ? 'rgba(255,255,255,0.45)' : 'rgba(0,0,0,0.45)'
  const borderColor = isDark ? 'rgba(255,255,255,0.08)' : 'rgba(0,0,0,0.06)'

  const xAxisData = chartData.value.map(item => item.date?.slice(5) || '')
  const totalData = chartData.value.map(item => item.total || 0)
  const processingData = chartData.value.map(item => item.processing || 0)
  const successData = chartData.value.map(item => item.success || 0)
  const failedData = chartData.value.map(item => item.failed || 0)

  return {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
    },
    legend: {
      data: ['总数', '处理中', '成功', '失败'],
      top: 0,
      textStyle: { color: textColor, fontSize: 12 },
    },
    grid: { left: 45, right: 20, top: 40, bottom: 20, containLabel: true },
    xAxis: {
      type: 'category',
      data: xAxisData,
      boundaryGap: false,
      axisLabel: { color: subtextColor, fontSize: 11 },
      axisLine: { show: false },
      axisTick: { show: false },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: subtextColor, fontSize: 11 },
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: borderColor, type: 'dashed' } },
    },
    series: [
      { name: '总数', type: 'line', smooth: 0.4, data: totalData, lineStyle: { width: 2.5, color: '#409EFF' }, itemStyle: { color: '#409EFF' } },
      { name: '处理中', type: 'line', smooth: 0.4, data: processingData, lineStyle: { width: 2, color: '#E6A23C' }, itemStyle: { color: '#E6A23C' } },
      { name: '成功', type: 'line', smooth: 0.4, data: successData, lineStyle: { width: 2, color: '#67C23A' }, itemStyle: { color: '#67C23A' } },
      { name: '失败', type: 'line', smooth: 0.4, data: failedData, lineStyle: { width: 2, color: '#F56C6C' }, itemStyle: { color: '#F56C6C' } },
    ],
    animationDuration: 500,
  }
}

const renderTrendChart = () => {
  if (!trendChartRef.value) {
    return
  }
  if (trendChartInstance) {
    trendChartInstance.dispose()
    trendChartInstance = null
  }
  trendChartInstance = echarts.init(trendChartRef.value, appStore.isDark ? 'dark' : null)
  trendChartInstance.setOption(getTrendChartOption())
}

const getModelChartOption = () => {
  const isDark = appStore.isDark
  const textColor = isDark ? 'rgba(255,255,255,0.7)' : 'rgba(0,0,0,0.7)'
  const list = Array.isArray(stats.value.modelStats) ? stats.value.modelStats : []
  const pieData = list.map((item) => ({
    name: formatModelKey(item.modelKey),
    value: Number(item.taskCount || 0),
  }))

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
        name: '模型调用',
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
            formatter: '{b}\n{c}次',
            fontSize: 12,
          },
        },
        data: pieData,
      },
    ],
  }
}

const renderModelChart = () => {
  if (!modelChartRef.value) {
    return
  }
  if (modelChartInstance) {
    modelChartInstance.dispose()
    modelChartInstance = null
  }
  modelChartInstance = echarts.init(modelChartRef.value, appStore.isDark ? 'dark' : null)
  modelChartInstance.setOption(getModelChartOption())
}

const fetchTrendData = async () => {
  chartLoading.value = true
  try {
    const params = {
      ...buildSearchParams(),
      status: '',
      trendDays: chartDays.value,
    }
    const res = await getTryonTaskTrend(params)
    if (res.code === 0) {
      chartData.value = res.data?.list || []
      await nextTick()
      renderTrendChart()
    }
  } finally {
    chartLoading.value = false
  }
}

const handleResize = () => {
  if (trendChartInstance) {
    trendChartInstance.resize()
  }
  if (modelChartInstance) {
    modelChartInstance.resize()
  }
}

const handleExportCsv = async () => {
  if (exporting.value) {
    return
  }
  exporting.value = true
  try {
    const currentTotal = Number(total.value || 0)
    if (currentTotal === 0) {
      ElMessage.warning('暂无可导出数据')
      return
    }

    const exportSize = Math.min(currentTotal, MAX_EXPORT_ROWS)
    const params = {
      page: 1,
      pageSize: exportSize,
      ...buildSearchParams(),
    }
    const res = await getTryonTaskList(params)
    if (res.code !== 0) {
      ElMessage.error(res.msg || '导出失败')
      return
    }

    const exportList = res.data?.list || []
    if (exportList.length === 0) {
      ElMessage.warning('暂无可导出数据')
      return
    }
    if (currentTotal > MAX_EXPORT_ROWS) {
      ElMessage.warning(`导出数量超过${MAX_EXPORT_ROWS}条，已截取前${MAX_EXPORT_ROWS}条`)
    }

    const content = buildCsvContent(exportList)
    downloadCsv(content, getExportFileName())
    ElMessage.success(`已导出${exportList.length}条记录`)
  } catch (e) {
    ElMessage.error(e?.message || '导出失败')
  } finally {
    exporting.value = false
  }
}

const getStatsData = async () => {
  statsLoading.value = true
  try {
    const params = {
      ...buildSearchParams(),
      status: '',
    }
    const res = await getTryonTaskStats(params)
    if (res.code === 0) {
      const modelStats = Array.isArray(res.data?.modelStats) ? res.data.modelStats : []
      stats.value = {
        total: Number(res.data?.total || 0),
        processing: Number(res.data?.processing || 0),
        success: Number(res.data?.success || 0),
        failed: Number(res.data?.failed || 0),
        totalCostPoints: Number(res.data?.totalCostPoints || 0),
        refinerEnabledCount: Number(res.data?.refinerEnabledCount || 0),
        modelCallTotal: Number(res.data?.modelCallTotal || 0),
        modelStats,
      }
      await nextTick()
      renderModelChart()
    }
  } finally {
    statsLoading.value = false
  }
}

const refreshStatsTabData = async () => {
  await Promise.all([getStatsData(), fetchTrendData()])
}

const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...buildSearchParams(),
  }
  const res = await getTryonTaskList(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = res.data.page || page.value
    pageSize.value = res.data.pageSize || pageSize.value
  }
}

const onSubmit = () => {
  page.value = 1
  Promise.all([getTableData(), refreshStatsTabData()])
}

const onReset = () => {
  searchInfo.value = {
    startCreatedAt: undefined,
    endCreatedAt: undefined,
    userID: undefined,
    taskNo: '',
    requestID: '',
    sceneType: '',
    status: ''
  }
  page.value = 1
  pageSize.value = 10
  Promise.all([getTableData(), refreshStatsTabData()])
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  getTableData()
}

const handleTabChange = (tabName) => {
  if (tabName === 'stats') {
    refreshStatsTabData()
  }
}

watch(() => appStore.isDark, async () => {
  await nextTick()
  renderTrendChart()
  renderModelChart()
})

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  await Promise.all([getTableData(), refreshStatsTabData()])
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (trendChartInstance) {
    trendChartInstance.dispose()
    trendChartInstance = null
  }
  if (modelChartInstance) {
    modelChartInstance.dispose()
    modelChartInstance = null
  }
})
</script>

<style scoped>
.panel-tabs {
  margin-bottom: 8px;
}

.stats-wrap {
  margin-bottom: 16px;
}

.stats-card {
  margin-bottom: 8px;
}

.stats-card-clickable {
  cursor: pointer;
}

.stats-card-clickable.active {
  border: 1px solid var(--el-color-primary);
  box-shadow: 0 0 0 1px color-mix(in srgb, var(--el-color-primary) 20%, transparent);
}

.trend-wrap {
  margin-bottom: 16px;
}

.chart-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.trend-subtitle {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.model-stats-wrap {
  margin-top: 8px;
}

.model-stats-title {
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 600;
}

.trend-header {
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.trend-title {
  font-size: 14px;
  font-weight: 600;
}

.trend-chart {
  width: 100%;
  height: 280px;
}

.copy-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.copy-text {
  display: inline-block;
  max-width: 170px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-image {
  width: 52px;
  height: 52px;
  border-radius: 6px;
  border: 1px solid var(--el-border-color-light);
}

.detail-image {
  width: 180px;
  height: 180px;
  border-radius: 8px;
  border: 1px solid var(--el-border-color-light);
}

.error-message {
  color: var(--el-color-danger);
  font-weight: 500;
}

@media (max-width: 1280px) {
  .chart-grid {
    grid-template-columns: 1fr;
  }
}

:deep(.el-table .tryon-row-non-processing > td.el-table__cell) {
  background: rgba(230, 162, 60, 0.12);
}
</style>
