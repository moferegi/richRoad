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

    <div class="gva-search-box stats-wrap">
      <el-row :gutter="16">
        <el-col :xs="24" :sm="12" :md="6">
          <el-card
            class="stats-card stats-card-clickable"
            :class="{ active: !searchInfo.status }"
            shadow="hover"
            v-loading="statsLoading"
            @click="handleStatusCardClick('')"
          >
            <el-statistic title="任务总数" :value="stats.total" />
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card
            class="stats-card stats-card-clickable"
            :class="{ active: searchInfo.status === 'processing' }"
            shadow="hover"
            v-loading="statsLoading"
            @click="handleStatusCardClick('processing')"
          >
            <el-statistic title="处理中" :value="stats.processing" />
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card
            class="stats-card stats-card-clickable"
            :class="{ active: searchInfo.status === 'success' }"
            shadow="hover"
            v-loading="statsLoading"
            @click="handleStatusCardClick('success')"
          >
            <el-statistic title="成功" :value="stats.success" />
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6">
          <el-card
            class="stats-card stats-card-clickable"
            :class="{ active: searchInfo.status === 'failed' }"
            shadow="hover"
            v-loading="statsLoading"
            @click="handleStatusCardClick('failed')"
          >
            <el-statistic title="失败" :value="stats.failed" />
          </el-card>
        </el-col>
      </el-row>
    </div>

    <div class="gva-search-box trend-wrap" v-loading="chartLoading">
      <div class="trend-header">
        <span class="trend-title">任务趋势</span>
        <el-radio-group v-model="chartDays" size="small" @change="fetchTrendData">
          <el-radio-button :value="7">近7天</el-radio-button>
          <el-radio-button :value="14">近14天</el-radio-button>
          <el-radio-button :value="30">近30天</el-radio-button>
        </el-radio-group>
      </div>
      <div ref="chartRef" class="trend-chart" />
    </div>

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

    <div class="gva-search-box model-search-box">
      <div class="model-title">我的模特管理</div>
      <el-form :inline="true" :model="modelSearchInfo" @keyup.enter="onModelSubmit">
        <el-form-item label="用户ID">
          <el-input v-model.number="modelSearchInfo.userID" clearable placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="模特名称">
          <el-input v-model="modelSearchInfo.name" clearable placeholder="请输入模特名称" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onModelSubmit">查询</el-button>
          <el-button icon="refresh" @click="onModelReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="danger"
          icon="delete"
          :disabled="!modelMultipleSelection.length"
          @click="handleModelBatchDelete"
        >
          批量删除模特
        </el-button>
      </div>

      <el-table
        :data="modelTableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        v-loading="modelLoading"
        @selection-change="handleModelSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="用户ID" prop="userID" width="90" />
        <el-table-column align="left" label="模特名称" prop="name" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="模特图" width="90">
          <template #default="scope">
            <el-image
              v-if="scope.row.image"
              class="task-image"
              :src="getUrl(scope.row.image)"
              :preview-src-list="[getUrl(scope.row.image)]"
              preview-teleported
              fit="cover"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="170">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" width="170">
          <template #default="scope">{{ formatDate(scope.row.UpdatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="danger" link @click="handleModelDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="modelPage"
          :page-size="modelPageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="modelTotal"
          @current-change="handleModelCurrentChange"
          @size-change="handleModelSizeChange"
        />
      </div>
    </div>

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
import { ref, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import { useAppStore } from '@/pinia'
import {
  getTryonTaskList,
  getTryonTaskStats,
  getTryonTaskTrend,
  deleteTryonTask,
  deleteTryonTaskByIds,
  getTryonModelList,
  deleteTryonModel,
  deleteTryonModelByIds,
} from '@/api/client/tryonTask'
import { formatDate } from '@/utils/format'
import { getUrl } from '@/utils/image'

defineOptions({
  name: 'TryonTaskManage'
})

const appStore = useAppStore()

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
})
const chartDays = ref(7)
const chartLoading = ref(false)
const chartRef = ref(null)
const chartData = ref([])
const modelPage = ref(1)
const modelTotal = ref(0)
const modelPageSize = ref(10)
const modelLoading = ref(false)
const modelTableData = ref([])
const modelMultipleSelection = ref([])
const searchInfo = ref({
  startCreatedAt: undefined,
  endCreatedAt: undefined,
  userID: undefined,
  taskNo: '',
  requestID: '',
  sceneType: '',
  status: ''
})
const modelSearchInfo = ref({
  userID: undefined,
  name: '',
})

const MAX_EXPORT_ROWS = 5000
let chartInstance = null

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

const handleStatusCardClick = (status) => {
  if (searchInfo.value.status === status) {
    return
  }
  searchInfo.value.status = status
  page.value = 1
  getTableData(true)
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
      await getTableData(true)
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
      await getTableData(true)
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

const buildModelSearchParams = () => {
  return {
    ...modelSearchInfo.value,
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
  if (!chartRef.value) {
    return
  }
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  chartInstance = echarts.init(chartRef.value, appStore.isDark ? 'dark' : null)
  chartInstance.setOption(getTrendChartOption())
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
      renderTrendChart()
    }
  } finally {
    chartLoading.value = false
  }
}

const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
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
      stats.value = {
        total: Number(res.data?.total || 0),
        processing: Number(res.data?.processing || 0),
        success: Number(res.data?.success || 0),
        failed: Number(res.data?.failed || 0),
      }
    }
  } finally {
    statsLoading.value = false
  }
}

const getTableData = async (refreshStats = false) => {
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
  if (refreshStats) {
    await Promise.all([getStatsData(), fetchTrendData()])
  }
}

const getModelTableData = async () => {
  modelLoading.value = true
  try {
    const params = {
      page: modelPage.value,
      pageSize: modelPageSize.value,
      ...buildModelSearchParams(),
    }
    const res = await getTryonModelList(params)
    if (res.code === 0) {
      modelTableData.value = res.data?.list || []
      modelTotal.value = res.data?.total || 0
      modelPage.value = res.data?.page || modelPage.value
      modelPageSize.value = res.data?.pageSize || modelPageSize.value
    }
  } finally {
    modelLoading.value = false
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData(true)
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
  getTableData(true)
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

const onModelSubmit = () => {
  modelPage.value = 1
  getModelTableData()
}

const onModelReset = () => {
  modelSearchInfo.value = {
    userID: undefined,
    name: '',
  }
  modelPage.value = 1
  modelPageSize.value = 10
  getModelTableData()
}

const handleModelCurrentChange = (val) => {
  modelPage.value = val
  getModelTableData()
}

const handleModelSizeChange = (val) => {
  modelPageSize.value = val
  modelPage.value = 1
  getModelTableData()
}

const handleModelSelectionChange = (rows) => {
  modelMultipleSelection.value = rows || []
}

const handleModelDelete = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除模特【${row.name || row.ID}】吗？`, '提示', { type: 'warning' })
    const res = await deleteTryonModel({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (modelTableData.value.length === 1 && modelPage.value > 1) {
        modelPage.value -= 1
      }
      await getModelTableData()
    }
  } catch (e) {
    // 用户取消删除时不做提示
  }
}

const handleModelBatchDelete = async () => {
  if (!modelMultipleSelection.value.length) {
    ElMessage.warning('请先选择要删除的数据')
    return
  }

  try {
    await ElMessageBox.confirm(`确定批量删除选中的${modelMultipleSelection.value.length}个模特吗？`, '提示', { type: 'warning' })
    const ids = modelMultipleSelection.value.map(item => item.ID).filter(Boolean)
    const res = await deleteTryonModelByIds({ 'IDs[]': ids })
    if (res.code === 0) {
      ElMessage.success('批量删除成功')
      if (ids.length >= modelTableData.value.length && modelPage.value > 1) {
        modelPage.value -= 1
      }
      modelMultipleSelection.value = []
      await getModelTableData()
    }
  } catch (e) {
    // 用户取消删除时不做提示
  }
}

window.addEventListener('resize', handleResize)

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})

getTableData(true)
getModelTableData()
</script>

<style scoped>
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

.model-search-box {
  margin-top: 16px;
}

.model-title {
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

:deep(.el-table .tryon-row-non-processing > td.el-table__cell) {
  background: rgba(230, 162, 60, 0.12);
}
</style>
