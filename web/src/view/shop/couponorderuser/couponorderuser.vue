<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="CreatedAt">
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon>
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期"
            :disabled-date="time => searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
          —
          <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期"
            :disabled-date="time => searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
        </el-form-item>

        <el-form-item label="优惠券ID" prop="couponID">
          <el-select v-model="searchInfo.couponID" filterable placeholder="请选择优惠券ID" :clearable="false">
            <el-option v-for="(item, key) in dataSource.couponID" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="订单ID" prop="orderID">
          <el-select v-model="searchInfo.orderID" filterable placeholder="请选择订单ID" :clearable="false">
            <el-option v-for="(item, key) in dataSource.orderID" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="使用状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width:120px">
            <el-option label="未使用" :value="false" />
            <el-option label="已使用" :value="true" />
          </el-select>
        </el-form-item>

        <el-form-item label="用户ID" prop="userID">
          <el-input-number v-model="searchInfo.userID" :min="0" controls-position="right" placeholder="用户ID" clearable style="width:140px" />
        </el-form-item>

        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true"
            v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length"
          @click="onDelete">删除</el-button>

      </div>
      <el-table ref="multipleTable" style="width: 100%" tooltip-effect="dark" :data="tableData" row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />

        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="券编号" prop="couponNum" width="180" show-overflow-tooltip />

        <el-table-column align="left" label="优惠券" prop="couponID" width="140">
          <template #default="scope">
            <span>{{ filterDataSource(dataSource.couponID, scope.row.couponID) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="订单ID" prop="orderID" width="120">
          <template #default="scope">
            <span>{{ filterDataSource(dataSource.orderID, scope.row.orderID) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="用户ID" prop="userID" width="100" />
        <el-table-column align="left" label="领取时间" prop="claimedAt" width="180">
          <template #default="scope">{{ scope.row.claimedAt ? formatDate(scope.row.claimedAt) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="使用时间" prop="usedAt" width="180">
          <template #default="scope">{{ scope.row.usedAt ? formatDate(scope.row.usedAt) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="状态" prop="status" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status ? 'success' : 'info'">{{ scope.row.status ? '已使用' : '未使用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon
                style="margin-right: 5px">
                <InfoFilled />
              </el-icon>查看</el-button>
            <el-button type="primary" link icon="edit" class="table-button"
              @click="updateCouponOrderUserFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper" :current-page="page" :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]" :total="total" @current-change="handleCurrentChange"
          @size-change="handleSizeChange" />
      </div>
    </div>
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false"
      :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增' : '编辑' }}</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="优惠券:" prop="couponID">
          <el-select v-model="formData.couponID" placeholder="请选择优惠券" filterable style="width:100%"
            :clearable="false">
            <el-option v-for="(item, key) in dataSource.couponID" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="订单ID:" prop="orderID">
          <el-select v-model="formData.orderID" placeholder="请选择订单ID" filterable style="width:100%" :clearable="false">
            <el-option v-for="(item, key) in dataSource.orderID" :key="key" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="券编号:" prop="couponNum">
          <el-input v-model="formData.couponNum" placeholder="自动生成" disabled />
        </el-form-item>
        <el-form-item label="用户ID:" prop="userID">
          <el-input-number v-model="formData.userID" :min="0" controls-position="right" style="width:100%" />
        </el-form-item>
        <el-form-item label="领取时间:" prop="claimedAt">
          <el-date-picker v-model="formData.claimedAt" type="datetime" placeholder="领取时间" style="width:100%" />
        </el-form-item>
        <el-form-item label="使用时间:" prop="usedAt">
          <el-date-picker v-model="formData.usedAt" type="datetime" placeholder="使用时间" style="width:100%" />
        </el-form-item>
        <el-form-item label="已使用:" prop="status">
          <el-switch v-model="formData.status" active-color="#13ce66" inactive-color="#ff4949" active-text="是"
            inactive-text="否" clearable></el-switch>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true"
      :before-close="closeDetailShow" title="查看">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="券编号">{{ detailFrom.couponNum }}</el-descriptions-item>
        <el-descriptions-item label="优惠券">
          {{ filterDataSource(dataSource.couponID, detailFrom.couponID) }}
        </el-descriptions-item>
        <el-descriptions-item label="订单ID">
          {{ filterDataSource(dataSource.orderID, detailFrom.orderID) }}
        </el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ detailFrom.userID }}</el-descriptions-item>
        <el-descriptions-item label="领取时间">{{ detailFrom.claimedAt ? formatDate(detailFrom.claimedAt) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="使用时间">{{ detailFrom.usedAt ? formatDate(detailFrom.usedAt) : '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailFrom.status ? 'success' : 'info'">{{ detailFrom.status ? '已使用' : '未使用' }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>

  </div>
</template>

<script setup>
import {
  getCouponOrderUserDataSource,
  createCouponOrderUser,
  deleteCouponOrderUser,
  deleteCouponOrderUserByIds,
  updateCouponOrderUser,
  findCouponOrderUser,
  getCouponOrderUserList
} from '@/api/shop/couponorderuser'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, filterDataSource, returnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { useAppStore } from "@/pinia"




defineOptions({
  name: 'CouponOrderUser'
})

// 提交按钮loading
const btnLoading = ref(false)
const appStore = useAppStore()

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  couponID: undefined,
  orderID: undefined,
  couponNum: '',
  userID: undefined,
  claimedAt: null,
  usedAt: null,
  status: false,
})
const dataSource = ref([])
const getDataSourceFunc = async () => {
  const res = await getCouponOrderUserDataSource()
  if (res.code === 0) {
    dataSource.value = res.data
  }
}
getDataSourceFunc()



// 验证规则
const rule = reactive({
  couponID: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
  status: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  },
  ],
})

const searchRule = reactive({
  CreatedAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 排序
const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt: 'created_at',
    ID: 'id',
  }
  let sort = sortMap[prop]
  if (!sort) {
    sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }
  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    if (searchInfo.value.status === "") {
      searchInfo.value.status = null
    }
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getCouponOrderUserList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    deleteCouponOrderUserFunc(row)
  })
}

// 多选删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = []
    if (multipleSelection.value.length === 0) {
      ElMessage({
        type: 'warning',
        message: '请选择要删除的数据'
      })
      return
    }
    multipleSelection.value &&
      multipleSelection.value.map(item => {
        IDs.push(item.ID)
      })
    const res = await deleteCouponOrderUserByIds({ IDs })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '删除成功'
      })
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const updateCouponOrderUserFunc = async (row) => {
  const res = await findCouponOrderUser({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}


// 删除行
const deleteCouponOrderUserFunc = async (row) => {
  const res = await deleteCouponOrderUser({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '删除成功'
    })
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    couponID: undefined,
    orderID: undefined,
    couponNum: '',
    userID: undefined,
    claimedAt: null,
    usedAt: null,
    status: false,
  }
}
// 弹窗确定
const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    let res
    switch (type.value) {
      case 'create':
        res = await createCouponOrderUser(formData.value)
        break
      case 'update':
        res = await updateCouponOrderUser(formData.value)
        break
      default:
        res = await createCouponOrderUser(formData.value)
        break
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
      closeDialog()
      getTableData()
    }
  })
}

const detailFrom = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findCouponOrderUser({ ID: row.ID })
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}


</script>

<style></style>
