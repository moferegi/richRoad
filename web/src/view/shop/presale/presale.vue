<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="预售状态">
          <el-select v-model="searchInfo.isPresale" placeholder="全部" clearable>
            <el-option label="预售中" :value="true" />
            <el-option label="非预售" :value="false" />
          </el-select>
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
        <el-table-column align="left" label="商品ID" prop="ID" width="80" />
        <el-table-column align="left" label="商品名称" prop="name" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="是否预售" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.isPresale ? 'warning' : 'info'">{{ scope.row.isPresale ? '预售' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="预售开始" width="180">
          <template #default="scope">{{ scope.row.presaleStartTime ? formatDate(scope.row.presaleStartTime) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="预售结束" width="180">
          <template #default="scope">{{ scope.row.presaleEndTime ? formatDate(scope.row.presaleEndTime) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="预计发货" width="180">
          <template #default="scope">{{ scope.row.presaleShipDate ? formatDate(scope.row.presaleShipDate) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="150">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="viewParticipants(scope.row)" v-if="scope.row.isPresale">
              <el-icon style="margin-right: 5px"><User /></el-icon>参与者
            </el-button>
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

    <el-drawer destroy-on-close size="600" v-model="participantShow" :show-close="true" title="预售参与者">
      <el-table :data="participants" style="width: 100%">
        <el-table-column label="用户ID" prop="userID" width="80" />
        <el-table-column label="订单号" prop="orderNo" width="200" show-overflow-tooltip />
        <el-table-column label="数量" prop="quantity" width="80" />
        <el-table-column label="下单时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
      </el-table>
    </el-drawer>
  </div>
</template>

<script setup>
import { getPresaleGoodList, getPresaleParticipants } from '@/api/shop/presale'
import { formatDate } from '@/utils/format'
import { ref } from 'vue'

defineOptions({ name: 'Presale' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const onReset = () => { searchInfo.value = {}; getTableData() }
const onSubmit = () => { page.value = 1; getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }

const getTableData = async() => {
  const table = await getPresaleGoodList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

const participantShow = ref(false)
const participants = ref([])

const viewParticipants = async (row) => {
  const res = await getPresaleParticipants({ goodID: row.ID })
  if (res.code === 0) {
    participants.value = res.data.list || []
    participantShow.value = true
  }
}
</script>
