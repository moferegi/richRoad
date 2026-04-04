<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="触发类型">
          <el-select v-model="searchInfo.triggerType" placeholder="全部" clearable>
            <el-option label="注册奖励" value="register" />
            <el-option label="下级注册奖励" value="sub_register" />
            <el-option label="签到奖励" value="sign_in" />
            <el-option label="下单奖励" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" sortable />
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="触发类型" width="150">
          <template #default="scope">
            <el-tag>{{ triggerTypeMap[scope.row.triggerType] || scope.row.triggerType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="奖励积分" prop="points" width="100" />
        <el-table-column align="left" label="奖励优惠券" width="250" show-overflow-tooltip>
          <template #default="scope">
            <span>{{ formatCouponNames(scope.row.couponIDs) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="要求下级首单" width="120">
          <template #default="scope">{{ scope.row.subRequireOrder ? '是' : '否' }}</template>
        </el-table-column>
        <el-table-column align="left" label="仅奖一次" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.triggerType === 'order'" :type="scope.row.orderOnce ? 'warning' : 'info'">{{ scope.row.orderOnce ? '是' : '否' }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'">{{ scope.row.isEnabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="200">
            <template #default="scope">
            <el-button type="primary" link icon="edit" class="table-button" @click="updateFunc(scope.row)">编辑</el-button>
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
    <el-drawer destroy-on-close size="800" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增':'编辑'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>
          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="触发类型:" prop="triggerType">
              <el-select v-model="formData.triggerType" placeholder="请选择触发类型">
                <el-option label="注册奖励" value="register" />
                <el-option label="下级注册奖励" value="sub_register" />
                <el-option label="签到奖励" value="sign_in" />
                <el-option label="下单奖励" value="order" />
              </el-select>
            </el-form-item>
            <el-form-item label="奖励积分:" prop="points">
              <el-input-number v-model="formData.points" :min="0" />
            </el-form-item>
            <el-form-item label="奖励优惠券:" prop="couponIDs">
              <el-select v-model="selectedCouponIds" multiple filterable placeholder="选择优惠券" style="width: 100%">
                <el-option
                  v-for="c in couponOptions"
                  :key="c.ID"
                  :label="`#${c.ID} ${c.name || ''} (满${c.minSpend/100}减${c.discount/100})`"
                  :value="c.ID"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="要求下级首单付款:" prop="subRequireOrder" v-if="formData.triggerType === 'sub_register'">
              <el-switch v-model="formData.subRequireOrder" />
            </el-form-item>
            <el-form-item label="下单奖励仅发放一次:" prop="orderOnce" v-if="formData.triggerType === 'order'">
              <el-switch v-model="formData.orderOnce" />
            </el-form-item>
            <el-form-item label="启用:" prop="isEnabled">
              <el-switch v-model="formData.isEnabled" />
            </el-form-item>
          </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createMarketingReward,
  deleteMarketingReward,
  updateMarketingReward,
  getMarketingRewardList
} from '@/api/shop/marketingReward'
import { getCouponList } from '@/api/shop/coupon'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, watch, onMounted } from 'vue'

defineOptions({ name: 'MarketingReward' })

const btnLoading = ref(false)
const triggerTypeMap = { register: '注册奖励', sub_register: '下级注册奖励', sign_in: '签到奖励', order: '下单奖励' }

const defaultForm = () => ({
  triggerType: '',
  isEnabled: false,
  points: 0,
  couponIDs: '',
  subRequireOrder: false,
  orderOnce: false,
})

// 优惠券下拉选项
const couponOptions = ref([])
const selectedCouponIds = ref([])

const loadCouponOptions = async () => {
  const res = await getCouponList({ page: 1, pageSize: 999 })
  if (res.code === 0) {
    couponOptions.value = res.data.list || []
  }
}

// couponIDs 字符串 <-> 数组 同步
watch(selectedCouponIds, (ids) => {
  formData.value.couponIDs = ids.join(',')
})

// 根据 couponIDs 反查优惠券名称
const couponMap = ref({})
watch(couponOptions, (list) => {
  const m = {}
  list.forEach(c => { m[c.ID] = c })
  couponMap.value = m
})
const formatCouponNames = (ids) => {
  if (!ids) return '-'
  return ids.split(',').map(id => {
    const c = couponMap.value[Number(id.trim())]
    return c ? `#${c.ID} ${c.name || ''}` : `#${id.trim()}`
  }).join(', ')
}

onMounted(() => {
  loadCouponOptions()
})

const formData = ref(defaultForm())

const rule = reactive({
  triggerType: [{ required: true, message: '请选择触发类型', trigger: 'change' }],
})

const elFormRef = ref()
const elSearchFormRef = ref()

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
  const table = await getMarketingRewardList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

const multipleSelection = ref([])
const handleSelectionChange = (val) => { multipleSelection.value = val }

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(() => deleteFunc(row))
}

const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async() => {
      const IDs = multipleSelection.value.map(item => item.ID)
      if (!IDs.length) { ElMessage({ type: 'warning', message: '请选择要删除的数据' }); return }
      for (const id of IDs) { await deleteMarketingReward({ ID: id }) }
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    })
}

const type = ref('')
const updateFunc = async(row) => {
  type.value = 'update'
  formData.value = { ...row }
  // 回填优惠券多选
  selectedCouponIds.value = row.couponIDs ? row.couponIDs.split(',').map(id => Number(id.trim())).filter(Boolean) : []
  dialogFormVisible.value = true
}
const deleteFunc = async (row) => {
  const res = await deleteMarketingReward({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '删除成功' })
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const dialogFormVisible = ref(false)
const openDialog = () => { type.value = 'create'; selectedCouponIds.value = []; dialogFormVisible.value = true }
const closeDialog = () => { dialogFormVisible.value = false; formData.value = defaultForm(); selectedCouponIds.value = [] }

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    let res
    switch (type.value) {
      case 'create': res = await createMarketingReward(formData.value); break
      case 'update': res = await updateMarketingReward(formData.value); break
      default: res = await createMarketingReward(formData.value); break
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '操作成功' })
      closeDialog()
      getTableData()
    }
  })
}
</script>
