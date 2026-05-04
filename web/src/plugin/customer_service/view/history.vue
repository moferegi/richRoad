<template>
  <div class="p-6">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xl font-semibold">历史会话</h2>
    </div>

    <!-- 筛选 -->
    <el-form inline class="mb-4">
      <el-form-item label="状态">
        <el-select v-model="query.status" clearable placeholder="全部" style="width:120px">
          <el-option label="排队中" value="pending" />
          <el-option label="进行中" value="active" />
          <el-option label="已结束" value="closed" />
        </el-select>
      </el-form-item>
      <el-form-item label="用户ID">
        <el-input v-model.number="query.clientUserId" clearable placeholder="输入客户端用户ID" style="width: 180px" />
      </el-form-item>
      <el-form-item label="时间范围">
        <el-date-picker
          v-model="dateRange"
          type="datetimerange"
          value-format="YYYY-MM-DDTHH:mm:ssZ"
          range-separator="至"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          style="width: 340px"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" border stripe>
      <el-table-column prop="ID" label="ID" width="70" />
      <el-table-column prop="clientUserId" label="用户ID" width="90" />
      <el-table-column prop="agentUserId" label="坐席ID" width="90" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="评分" width="130">
        <template #default="{ row }">
          <el-rate v-if="row.rating" :model-value="row.rating" disabled size="small" />
          <span v-else class="text-gray-400 text-sm">未评价</span>
        </template>
      </el-table-column>
      <el-table-column prop="CreatedAt" label="创建时间" :formatter="(r) => fmtTime(r.CreatedAt)" />
      <el-table-column prop="UpdatedAt" label="更新时间" :formatter="(r) => fmtTime(r.UpdatedAt)" />
    </el-table>

    <div class="flex justify-end mt-4">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @change="loadList"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import { getConversationList } from '@/api/customerService'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const dateRange = ref([])
const query = ref({
  page: 1,
  pageSize: 20,
  status: '',
  clientUserId: undefined,
  startTime: undefined,
  endTime: undefined,
})

async function loadList() {
  loading.value = true
  try {
    query.value.startTime = dateRange.value?.[0] || undefined
    query.value.endTime = dateRange.value?.[1] || undefined
    const res = await getConversationList(query.value)
    if (res.code === 0) {
      list.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.value = {
    page: 1,
    pageSize: 20,
    status: '',
    clientUserId: undefined,
    startTime: undefined,
    endTime: undefined,
  }
  dateRange.value = []
  loadList()
}

function statusLabel(s) {
  return { pending: '排队中', active: '进行中', closed: '已结束' }[s] || s
}
function statusType(s) {
  return { pending: 'warning', active: 'success', closed: 'info' }[s] || ''
}
function fmtTime(t) {
  return t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'
}

onMounted(loadList)
</script>
