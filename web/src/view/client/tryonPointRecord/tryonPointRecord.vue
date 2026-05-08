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
          <el-input v-model="searchInfo.operationType" clearable placeholder="如 tryon_consume" />
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
        <span>仅展示资产类型为 tryon_point 的流水，积分流水不会混入。</span>
      </div>
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
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getPointRecordList } from '@/api/client/pointRecord'
import { getClientUserList } from '@/api/client/user'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'TryonPointRecord'
})

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const tableData = ref([])
const searchInfo = ref({ assetType: 'tryon_point', sort: 'created_at', order: 'descending' })
const userLoading = ref(false)
const userOptions = ref([])
const userCache = ref(new Map())

const sortChange = ({ prop, order }) => {
  const sortMap = { CreatedAt: 'created_at', pointChange: 'point_change', ID: 'id' }
  searchInfo.value.sort = sortMap[prop] || prop
  searchInfo.value.order = order || 'descending'
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  page.value = 1
  searchInfo.value = { assetType: 'tryon_point', sort: 'created_at', order: 'descending' }
  getTableData()
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

initUserOptions()
getTableData()
</script>

<style scoped>
.record-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  color: #606266;
}
</style>
