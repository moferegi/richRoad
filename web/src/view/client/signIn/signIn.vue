<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="用户ID">
          <el-input v-model="searchInfo.userId" placeholder="请输入用户ID" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchInfo.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="签到日期">
          <el-date-picker
            v-model="searchInfo.signDateRange"
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
      <div class="gva-btn-list">
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="id"
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="id" width="80" />
        <el-table-column align="left" label="用户ID" prop="user_id" width="100" />
        <el-table-column align="left" label="用户名" prop="username" width="140" />
        <el-table-column align="left" label="昵称" prop="nickname" width="140" />
        <el-table-column align="left" label="手机号" prop="phone" width="140" />
        <el-table-column align="left" label="签到日期" prop="sign_date" sortable="custom" width="140">
          <template #default="scope">{{ formatDate(scope.row.sign_date) }}</template>
        </el-table-column>
        <el-table-column align="left" label="连续签到天数" prop="continuous_days" sortable="custom" width="140" />
        <el-table-column align="left" label="总签到天数" prop="total_days" sortable="custom" width="130" />
        <el-table-column align="left" label="操作" fixed="right" min-width="120">
          <template #default="scope">
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
import { getSignInList, deleteSignIn } from '@/api/client/signIn'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref } from 'vue'

defineOptions({
  name: 'SignInManage'
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])
const orderKey = ref('')
const orderDesc = ref(true)

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const handleSortChange = ({ prop, order }) => {
  if (!order) {
    orderKey.value = ''
    orderDesc.value = true
  } else {
    orderKey.value = prop
    orderDesc.value = order === 'descending'
  }
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  orderKey.value = ''
  orderDesc.value = true
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

const formatDate = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
}

const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value,
  }
  if (orderKey.value) {
    params.orderKey = orderKey.value
    params.desc = orderDesc.value
  }
  const res = await getSignInList(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

getTableData()

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    const res = await deleteSignIn({ ID: row.id })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}

const onDelete = () => {
  ElMessageBox.confirm('确定要删除选中项吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    const ids = multipleSelection.value.map(item => item.id)
    for (const id of ids) {
      await deleteSignIn({ ID: id })
    }
    ElMessage.success('删除成功')
    getTableData()
  })
}
</script>
