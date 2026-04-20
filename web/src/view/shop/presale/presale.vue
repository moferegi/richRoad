<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="预售开关">
          <el-select v-model="searchInfo.presaleEnabled" placeholder="全部" clearable>
            <el-option label="已开启" :value="true" />
            <el-option label="未开启" :value="false" />
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
        max-height="70vh"
        @sort-change="sortChange"
        >
        <el-table-column align="left" label="日期" width="180" sortable="custom" prop="created_at">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="商品ID" prop="ID" width="80" />
        <el-table-column align="left" label="商品名称" prop="name" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="是否预售" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.isPresale ? 'warning' : 'info'">{{ scope.row.isPresale ? '预售' : '否' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="预售开关" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.isPresale" :type="scope.row.presaleEnabled ? 'success' : 'danger'">{{ scope.row.presaleEnabled ? '开启' : '关闭' }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="预售数量" width="120">
          <template #default="scope">
            <span v-if="scope.row.isPresale">{{ scope.row.presaleSold || 0 }} / {{ scope.row.presaleQty || 0 }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="预售开始" width="180">
          <template #default="scope">{{ scope.row.presaleStart ? formatDate(scope.row.presaleStart) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="预售结束" width="180">
          <template #default="scope">{{ scope.row.presaleEnd ? formatDate(scope.row.presaleEnd) : '-' }}</template>
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

    <el-drawer destroy-on-close size="800" v-model="participantShow" :show-close="true" :title="`预售参与者 - ${currentGoodName}`">
      <el-table :data="participants" style="width: 100%">
        <el-table-column label="订单ID" prop="ID" width="80" />
        <el-table-column label="用户ID" prop="userID" width="80" />
        <el-table-column label="订单金额(元)" width="120">
          <template #default="scope">¥{{ (scope.row.totalPrice / 100).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="订单状态" width="100">
          <template #default="scope">
            <el-tag :type="orderStatusType(scope.row.status)">{{ orderStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="购买商品" min-width="200">
          <template #default="scope">
            <div v-for="d in (scope.row.detail || [])" :key="d.ID" class="participant-detail-item">
              {{ d.good?.name || '-' }} × {{ d.quantity }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="下单时间" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination" style="margin-top: 10px;">
        <el-pagination
          layout="total, prev, pager, next"
          :current-page="participantPage"
          :page-size="participantPageSize"
          :total="participantTotal"
          @current-change="handleParticipantPageChange"
        />
      </div>
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

// 排序
const sortChange = ({ prop, order }) => {
  if (order) {
    searchInfo.value.orderBy = prop
    searchInfo.value.orderDir = order === 'ascending' ? 'asc' : 'desc'
  } else {
    searchInfo.value.orderBy = ''
    searchInfo.value.orderDir = ''
  }
  getTableData()
}

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
const participantTotal = ref(0)
const participantPage = ref(1)
const participantPageSize = ref(10)
const currentGoodID = ref(0)
const currentGoodName = ref('')

const viewParticipants = async (row) => {
  currentGoodID.value = row.ID
  currentGoodName.value = row.name || `商品#${row.ID}`
  participantPage.value = 1
  await loadParticipants()
  participantShow.value = true
}

const loadParticipants = async () => {
  const res = await getPresaleParticipants({ goodID: currentGoodID.value, page: participantPage.value, pageSize: participantPageSize.value })
  if (res.code === 0) {
    participants.value = res.data.list || []
    participantTotal.value = res.data.total
  }
}

const handleParticipantPageChange = (val) => {
  participantPage.value = val
  loadParticipants()
}

const orderStatusMap = { '0': '待付款', '1': '待发货', '2': '待收货', '3': '已完成', '4': '已关闭', '5': '已退款', '6': '退款中' }
const orderStatusTypeMap = { '0': 'warning', '1': 'primary', '2': 'primary', '3': 'success', '4': 'info', '5': 'info', '6': 'danger' }
const orderStatusText = (s) => orderStatusMap[s] || s
const orderStatusType = (s) => orderStatusTypeMap[s] || 'info'
</script>
