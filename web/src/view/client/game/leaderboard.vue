<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="榜单类型">
          <el-select v-model="searchInfo.type" placeholder="选择榜单类型" @change="onSearch">
            <el-option label="累计闯关榜" value="all" />
            <el-option label="每日闯关榜" value="daily" />
          </el-select>
        </el-form-item>
        <el-form-item label="显示数量">
          <el-input-number v-model="searchInfo.limit" :min="1" :max="200" @change="onSearch" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="userID" border>
        <el-table-column align="left" label="排名" width="80">
          <template #default="scope">
            <el-tag :type="scope.$index < 3 ? 'warning' : 'info'">{{ scope.$index + 1 }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="用户ID" prop="userID" width="100" />
        <el-table-column align="left" label="用户名" prop="username" min-width="150" />
        <el-table-column align="left" label="昵称" prop="nickname" min-width="150" />
        <el-table-column align="left" label="闯关数" prop="total" width="120">
          <template #default="scope">
            <el-tag type="success">{{ scope.row.total }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { getLeaderboard } from '@/api/client/game'

defineOptions({ name: 'GameLeaderboard' })

const tableData = ref([])
const searchInfo = reactive({ type: 'all', limit: 50 })

const getTableData = async () => {
  const res = await getLeaderboard({ type: searchInfo.type, limit: searchInfo.limit })
  if (res.code === 0) {
    tableData.value = res.data || []
  }
}

const onSearch = () => { getTableData() }

getTableData()
</script>