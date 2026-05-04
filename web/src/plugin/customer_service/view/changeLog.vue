<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true">
        <el-form-item label="记录范围">
          <el-select v-model="scope" class="!w-56" @change="handleScopeChange">
            <el-option label="全部配置变更" value="all" />
            <el-option label="客服配置" value="cs" />
            <el-option label="系统参数" value="sys" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态码">
          <el-input v-model="statusFilter" clearable placeholder="如: 200 / 500" @keyup.enter="onSubmit" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button :loading="loading" @click="refresh">刷新</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <span class="text-xs text-gray-500">仅展示配置类写操作（PUT）</span>
      </div>

      <el-table :data="tableData" v-loading="loading" row-key="ID" style="width: 100%">
        <el-table-column label="时间" min-width="170">
          <template #default="scopeRow">
            {{ formatDate(scopeRow.row.CreatedAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作人" min-width="180" show-overflow-tooltip>
          <template #default="scopeRow">
            {{ formatOperator(scopeRow.row) }}
          </template>
        </el-table-column>

        <el-table-column label="状态码" prop="status" width="100">
          <template #default="scopeRow">
            <el-tag :type="statusTagType(scopeRow.row.status)">
              {{ scopeRow.row.status }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="请求方法" prop="method" width="100" />
        <el-table-column label="请求路径" prop="path" min-width="300" show-overflow-tooltip />

        <el-table-column label="请求体" width="100">
          <template #default="scopeRow">
            <el-popover v-if="scopeRow.row.body" placement="left-start" :width="480">
              <pre class="log-pre">{{ fmtBody(scopeRow.row.body) }}</pre>
              <template #reference>
                <el-button type="primary" link>查看</el-button>
              </template>
            </el-popover>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>

        <el-table-column label="响应体" width="100">
          <template #default="scopeRow">
            <el-popover v-if="scopeRow.row.resp" placement="left-start" :width="480">
              <pre class="log-pre">{{ fmtBody(scopeRow.row.resp) }}</pre>
              <template #reference>
                <el-button type="primary" link>查看</el-button>
              </template>
            </el-popover>
            <span v-else class="text-gray-400">-</span>
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
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/format'
import { getSysOperationRecordList } from '@/api/sysOperationRecord'

defineOptions({
  name: 'CsChangeLog'
})

const scope = ref('all')
const statusFilter = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const tableData = ref([])

const scopePathMap = {
  // 这里用不带前导斜杠的模糊关键字，能够同时命中:
  // /cs/admin/config/update 和 /sysConfig/updateSysConfig
  all: 'config/update',
  cs: '/cs/admin/config/update',
  sys: '/sysConfig/updateSysConfig'
}

function formatOperator(row) {
  const userName = row?.user?.userName || ''
  const nickName = row?.user?.nickName || ''
  if (!userName && !nickName) {
    return '-'
  }
  if (!nickName) {
    return userName
  }
  return `${userName}(${nickName})`
}

function statusTagType(status) {
  if (status >= 500) return 'danger'
  if (status >= 400) return 'warning'
  return 'success'
}

function fmtBody(value) {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  if (typeof value === 'object') {
    return JSON.stringify(value, null, 2)
  }
  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch (_) {
    return String(value)
  }
}

async function getTableData() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      method: 'PUT',
      path: scopePathMap[scope.value] || scopePathMap.all
    }
    const status = Number(statusFilter.value)
    if (statusFilter.value !== '' && !Number.isNaN(status) && status > 0) {
      params.status = status
    }

    const res = await getSysOperationRecordList(params)
    if (res.code !== 0) {
      ElMessage.error(res.msg || '获取失败')
      return
    }

    tableData.value = res.data?.list || []
    total.value = Number(res.data?.total || 0)
    page.value = Number(res.data?.page || page.value)
    pageSize.value = Number(res.data?.pageSize || pageSize.value)
  } catch (_) {
    ElMessage.error('获取失败')
  } finally {
    loading.value = false
  }
}

function onSubmit() {
  page.value = 1
  getTableData()
}

function onReset() {
  scope.value = 'all'
  statusFilter.value = ''
  page.value = 1
  pageSize.value = 10
  getTableData()
}

function handleScopeChange() {
  page.value = 1
  getTableData()
}

function handleCurrentChange(current) {
  page.value = current
  getTableData()
}

function handleSizeChange(size) {
  pageSize.value = size
  page.value = 1
  getTableData()
}

function refresh() {
  getTableData()
}

getTableData()
</script>

<style scoped>
.log-pre {
  margin: 0;
  max-height: 480px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  line-height: 1.5;
}
</style>
