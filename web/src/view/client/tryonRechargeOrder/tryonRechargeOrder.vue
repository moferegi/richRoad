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
        <el-form-item label="用户ID" prop="userID">
          <el-input v-model.number="searchInfo.userID" clearable placeholder="请输入用户ID" style="width: 140px" />
        </el-form-item>
        <el-form-item label="订单状态" prop="status">
          <el-select v-model="searchInfo.status" clearable placeholder="全部" style="width: 150px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="订单号" prop="outTradeNo">
          <el-input v-model="searchInfo.outTradeNo" clearable placeholder="请输入订单号" style="width: 180px" />
        </el-form-item>
        <el-form-item label="支付方式" prop="payMethod">
          <el-input v-model="searchInfo.payMethod" clearable placeholder="如 qrcode/contact" style="width: 160px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID" :default-sort="{ prop: 'ID', order: 'descending' }" @sort-change="sortChange">
        <el-table-column align="left" label="ID" prop="ID" width="80" sortable="custom" />
        <el-table-column align="left" label="创建时间" prop="CreatedAt" width="180" sortable="custom">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="订单号" prop="outTradeNo" min-width="180" />
        <el-table-column align="left" label="用户ID" prop="userID" width="100" />
        <el-table-column align="left" label="试衣币" prop="points" width="110" />
        <el-table-column align="left" label="金额(元)" prop="amount" width="120">
          <template #default="scope">{{ amountYuan(scope.row.amount) }}</template>
        </el-table-column>
        <el-table-column align="left" label="支付方式" prop="payMethod" width="140">
          <template #default="scope">{{ payMethodLabel(scope.row.payMethod) }}</template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="status" width="130">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.status)">{{ statusLabel(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="确认时间" prop="paidAt" width="180">
          <template #default="scope">{{ scope.row.paidAt ? formatDate(scope.row.paidAt) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" min-width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openDetail(scope.row)">查看详情</el-button>
            <el-button
              v-if="scope.row.status === '0' || scope.row.status === '8'"
              type="success"
              link
              @click="confirmPayment(scope.row)"
            >确认到账</el-button>
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

    <el-drawer v-model="detailVisible" size="680" :show-close="false">
      <template #header>
        <div class="drawer-header">
          <span class="text-lg">试衣币充值订单详情</span>
          <el-button @click="detailVisible = false">关 闭</el-button>
        </div>
      </template>

      <div v-if="detailData">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="订单ID">{{ detailData.ID }}</el-descriptions-item>
          <el-descriptions-item label="订单号">{{ detailData.outTradeNo || '-' }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ detailData.userID }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTagType(detailData.status)">{{ statusLabel(detailData.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="试衣币">{{ detailData.points }}</el-descriptions-item>
          <el-descriptions-item label="金额(元)">{{ amountYuan(detailData.amount) }}</el-descriptions-item>
          <el-descriptions-item label="支付方式">{{ payMethodLabel(detailData.payMethod) }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(detailData.CreatedAt) }}</el-descriptions-item>
          <el-descriptions-item label="关闭时间">{{ detailData.closeTime ? formatDate(detailData.closeTime) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="确认时间">{{ detailData.paidAt ? formatDate(detailData.paidAt) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="取消时间">{{ detailData.cancelledAt ? formatDate(detailData.cancelledAt) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="备注">{{ detailData.remark || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'
import {
  getTryonRechargeOrderList,
  findTryonRechargeOrder,
  confirmTryonRechargeOrderPayment,
} from '@/api/client/tryonRechargeOrder'

defineOptions({
  name: 'TryonRechargeOrder'
})

const statusOptions = [
  { label: '待支付', value: '0' },
  { label: '待后台确认', value: '8' },
  { label: '已支付', value: '1' },
  { label: '已取消', value: '4' },
]

const detailVisible = ref(false)
const detailData = ref(null)

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  createdAtRange: [],
  status: '',
  userID: undefined,
  outTradeNo: '',
  payMethod: '',
  sort: 'id',
  order: 'descending',
})

const statusLabel = (status) => {
  const hit = statusOptions.find(item => item.value === String(status))
  return hit?.label || String(status || '-')
}

const statusTagType = (status) => {
  const map = {
    '0': 'info',
    '8': 'warning',
    '1': 'success',
    '4': 'danger',
  }
  return map[String(status)] || 'info'
}

const amountYuan = (amount) => {
  const cents = Number(amount || 0)
  return (cents / 100).toFixed(2)
}

const payMethodLabel = (method) => {
  const map = {
    qrcode: '扫码支付',
    contact: '联系客服',
    wechat: '微信支付',
    alipay: '支付宝',
    bank_card_cn: '银行卡(国内)',
    bank_card_us: 'Bank Card (US)',
    bank_card_mn: 'Bank Card (MN)',
    paypal: 'PayPal',
  }
  return map[method] || method || '-'
}

const sortChange = ({ prop, order }) => {
  const sortMap = {
    ID: 'id',
    CreatedAt: 'created_at',
  }
  searchInfo.value.sort = sortMap[prop] || 'id'
  searchInfo.value.order = order || 'descending'
  getTableData()
}

const buildSearchParams = () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    status: searchInfo.value.status || undefined,
    userID: searchInfo.value.userID || undefined,
    outTradeNo: searchInfo.value.outTradeNo || undefined,
    payMethod: searchInfo.value.payMethod || undefined,
  }

  const range = searchInfo.value.createdAtRange
  if (Array.isArray(range) && range.length === 2 && range[0] && range[1]) {
    params.startCreatedAt = range[0]
    params.endCreatedAt = range[1]
  }

  return params
}

const getTableData = async () => {
  const res = await getTryonRechargeOrderList(buildSearchParams())
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = res.data.page || page.value
    pageSize.value = res.data.pageSize || pageSize.value
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {
    createdAtRange: [],
    status: '',
    userID: undefined,
    outTradeNo: '',
    payMethod: '',
    sort: 'id',
    order: 'descending',
  }
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

const openDetail = async (row) => {
  const res = await findTryonRechargeOrder({ ID: row.ID })
  if (res.code === 0) {
    detailData.value = res.data.order || null
    detailVisible.value = true
  }
}

const confirmPayment = (row) => {
  ElMessageBox.prompt('请输入备注（可选）', '确认到账', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPlaceholder: '到账备注'
  }).then(async ({ value }) => {
    const res = await confirmTryonRechargeOrderPayment({ ID: row.ID, remark: value || '' })
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '确认到账成功' })
      if (detailVisible.value && detailData.value?.ID === row.ID) {
        await openDetail(row)
      }
      getTableData()
    }
  })
}

getTableData()
</script>

<style scoped>
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
