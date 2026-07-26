<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="用户">
          <el-select
            v-model="searchInfo.userId"
            filterable
            remote
            reserve-keyword
            clearable
            placeholder="搜索用户(用户名/昵称/手机号)"
            :remote-method="searchUsers"
            :loading="userLoading"
            @clear="onReset"
          >
            <el-option
              v-for="item in userOptions"
              :key="item.id"
              :label="`${item.username}${item.nickName ? '(' + item.nickName + ')' : ''}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="更新时间">
          <el-date-picker
            v-model="searchInfo.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
      >
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="用户" min-width="150">
          <template #default="scope">
            <span>{{ scope.row.userName || scope.row.nickName || scope.row.username || ('用户#' + scope.row.userId) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="单集名称" prop="episodeName" min-width="180" show-overflow-tooltip />
        <el-table-column align="left" label="观看进度" width="130">
          <template #default="scope">
            <el-tag v-if="scope.row.progressSecs > 0" type="warning" size="small">
              {{ formatProgress(scope.row.progressSecs) }}
            </el-tag>
            <span v-else class="text-gray">未开始</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.UpdatedAt) }}</template>
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
import { getAdminWatchHistoryList, getAdminUserList } from '@/api/client/englishLearningAdmin'
import { formatDate } from '@/utils/format'
import { ref } from 'vue'

defineOptions({
  name: 'EnglishWatchHistory'
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const elSearchFormRef = ref()
const userOptions = ref([])
const userLoading = ref(false)

const formatProgress = (secs) => {
  if (!secs || secs <= 0) return '0:00'
  const m = Math.floor(secs / 60)
  const s = Math.floor(secs % 60)
  return `${m}:${String(s).padStart(2, '0')}`
}

const searchUsers = async (keyword) => {
  if (!keyword) {
    userOptions.value = []
    return
  }
  userLoading.value = true
  try {
    const res = await getAdminUserList({ keyword })
    if (res.code === 0) {
      userOptions.value = res.data || []
    }
  } finally {
    userLoading.value = false
  }
}

const onReset = () => {
  searchInfo.value = {}
  userOptions.value = []
  getTableData()
}

const onSubmit = () => {
  page.value = 1
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
  const params = {
    page: page.value,
    pageSize: pageSize.value
  }
  if (searchInfo.value.userId) {
    params.userId = searchInfo.value.userId
  }
  if (searchInfo.value.dateRange && searchInfo.value.dateRange.length === 2) {
    params.startDate = searchInfo.value.dateRange[0]
    params.endDate = searchInfo.value.dateRange[1]
  }
  const res = await getAdminWatchHistoryList(params)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

getTableData()
</script>

<style scoped>
.text-gray {
  color: #909399;
  font-size: 12px;
}
</style>
