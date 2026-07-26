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
        <el-form-item label="签到日期">
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
            <span>{{ scope.row.userName || scope.row.nickName || scope.row.username || ('用户' + scope.row.userId) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="签到日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.checkinDate) }}</template>
        </el-table-column>
        <el-table-column align="left" label="获得积分" width="120">
          <template #default="scope">
            <el-tag type="success">+{{ scope.row.pointAward }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
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
import { getAdminCheckinRecordList, getAdminUserList } from '@/api/client/englishLearningAdmin'
import { formatDate } from '@/utils/format'
import { ref } from 'vue'

defineOptions({
  name: 'EnglishCheckinRecord'
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const elSearchFormRef = ref()
const userOptions = ref([])
const userLoading = ref(false)

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
  const res = await getAdminCheckinRecordList(params)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

getTableData()
</script>
