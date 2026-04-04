<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="createdAt">
      <template #label>
        <span>
          创建日期
          <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </span>
      </template>
      <el-date-picker v-model="searchInfo.startCreatedAt" type="datetime" placeholder="开始日期" :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="datetime" placeholder="结束日期" :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"></el-date-picker>
      </el-form-item>

        <el-form-item label="商品ID">
          <el-input v-model.number="searchInfo.goodID" placeholder="请输入商品ID" clearable />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <!-- 进货汇总 -->
    <div v-if="summaryData" class="gva-search-box" style="padding: 12px 20px; margin-bottom: 0;">
      <el-descriptions title="商品进货汇总" :column="3" border size="small">
        <el-descriptions-item label="总进货数量">{{ summaryData.totalQty || 0 }}</el-descriptions-item>
        <el-descriptions-item label="总金额(元)">{{ ((summaryData.totalCost || 0) / 100).toFixed(2) }}</el-descriptions-item>
        <el-descriptions-item label="商品ID">{{ searchInfo.goodID }}</el-descriptions-item>
      </el-descriptions>
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
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />
        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="商品" prop="goodID" width="120">
          <template #default="scope">
            {{ scope.row.good?.name || scope.row.goodID }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="SKU" prop="skuID" width="100">
          <template #default="scope">
            {{ scope.row.sku?.name || scope.row.skuID || '-' }}
          </template>
        </el-table-column>
        <el-table-column sortable align="left" label="数量" prop="quantity" width="100" />
        <el-table-column sortable align="left" label="单价(元)" prop="unitCost" width="120">
          <template #default="scope">{{ (scope.row.unitCost / 100).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column sortable align="left" label="总金额(元)" prop="totalCost" width="120">
          <template #default="scope">{{ (scope.row.totalCost / 100).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column align="left" label="供应商" prop="supplier" width="150" />
        <el-table-column align="left" label="进货日期" prop="purchaseDate" width="180">
          <template #default="scope">{{ scope.row.purchaseDate ? formatDate(scope.row.purchaseDate) : '-' }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
            <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
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
            <el-form-item label="商品ID:" prop="goodID">
              <el-input v-model.number="formData.goodID" placeholder="请输入商品ID" />
            </el-form-item>
            <el-form-item label="SKU ID:" prop="skuID">
              <el-select v-model="formData.skuID" placeholder="输入商品ID后自动加载SKU" clearable filterable style="width:100%;">
                <el-option v-for="s in skuOptions" :key="s.id" :label="s.label" :value="s.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="进货数量:" prop="quantity">
              <el-input-number v-model="formData.quantity" :min="1" placeholder="请输入数量" />
            </el-form-item>
            <el-form-item label="单价(分):" prop="unitCost">
              <el-input-number v-model="formData.unitCost" :min="0" placeholder="请输入单价(分)" />
            </el-form-item>
            <el-form-item label="总金额(分):" prop="totalCost">
              <el-input-number v-model="formData.totalCost" :min="0" placeholder="请输入总金额(分)" />
            </el-form-item>
            <el-form-item label="供应商:" prop="supplier">
              <el-input v-model="formData.supplier" placeholder="请输入供应商" />
            </el-form-item>
            <el-form-item label="进货日期:" prop="purchaseDate">
              <el-date-picker v-model="formData.purchaseDate" type="datetime" placeholder="选择进货日期" />
            </el-form-item>
            <el-form-item label="备注:" prop="remark">
              <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="请输入备注" />
            </el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看详情">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="商品">{{ detailFrom.good?.name || detailFrom.goodID }}</el-descriptions-item>
                    <el-descriptions-item label="SKU">{{ detailFrom.sku?.name || detailFrom.skuID || '-' }}</el-descriptions-item>
                    <el-descriptions-item label="进货数量">{{ detailFrom.quantity }}</el-descriptions-item>
                    <el-descriptions-item label="单价(元)">{{ (detailFrom.unitCost / 100).toFixed(2) }}</el-descriptions-item>
                    <el-descriptions-item label="总金额(元)">{{ (detailFrom.totalCost / 100).toFixed(2) }}</el-descriptions-item>
                    <el-descriptions-item label="供应商">{{ detailFrom.supplier }}</el-descriptions-item>
                    <el-descriptions-item label="进货日期">{{ detailFrom.purchaseDate ? formatDate(detailFrom.purchaseDate) : '-' }}</el-descriptions-item>
                    <el-descriptions-item label="备注">{{ detailFrom.remark }}</el-descriptions-item>
            </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createGoodPurchase,
  deleteGoodPurchase,
  updateGoodPurchase,
  getGoodPurchaseList,
  getGoodPurchaseSummary
} from '@/api/shop/goodPurchase'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getSkuList } from '@/api/shop/sku'

defineOptions({ name: 'GoodPurchase' })

const route = useRoute()

const btnLoading = ref(false)

const formData = ref({
  goodID: null,
  skuID: null,
  quantity: 1,
  unitCost: 0,
  totalCost: 0,
  supplier: '',
  purchaseDate: null,
  remark: '',
})

const rule = reactive({
  goodID: [{ required: true, message: '请输入商品ID', trigger: 'blur' }],
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }],
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && searchInfo.value.startCreatedAt.getTime() >= searchInfo.value.endCreatedAt.getTime()) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// === 进货汇总 ===
const summaryData = ref(null)
const loadSummary = async (goodID) => {
  if (!goodID) { summaryData.value = null; return }
  try {
    const res = await getGoodPurchaseSummary({ goodId: goodID })
    if (res.code === 0) {
      summaryData.value = res.data
    }
  } catch (e) { summaryData.value = null }
}

watch(() => searchInfo.value.goodID, (val) => { loadSummary(val) })

const sortChange = ({ prop, order }) => {
  let sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  searchInfo.value.sort = sort
  searchInfo.value.order = order
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    getTableData()
  })
}

const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }

const getTableData = async() => {
  const table = await getGoodPurchaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}
getTableData()

// 从商品管理页跳转过来时，自动打开新增表单并预填商品ID
onMounted(() => {
  const qGoodID = route.query.goodID
  if (qGoodID) {
    openDialog(Number(qGoodID))
  }
})

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
      // 逐个删除
      for (const id of IDs) {
        await deleteGoodPurchase({ ID: id })
      }
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    })
}

const type = ref('')

const updateFunc = async(row) => {
  type.value = 'update'
  formData.value = { ...row }
  dialogFormVisible.value = true
}

const deleteFunc = async (row) => {
  const res = await deleteGoodPurchase({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '删除成功' })
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const dialogFormVisible = ref(false)

// === SKU 下拉列表 ===
const skuOptions = ref([])
const loadSkuOptions = async (goodID) => {
  skuOptions.value = []
  if (!goodID) return
  try {
    const res = await getSkuList({ page: 1, pageSize: 100, goodID })
    if (res.code === 0 && res.data.list) {
      skuOptions.value = res.data.list.map(sku => {
        const specStr = formatSkuSpecs(sku.attrs)
        return { id: sku.ID, label: `${sku.ID} - ${sku.name || ''} ${specStr}`.trim() }
      })
    }
  } catch (e) { /* ignore */ }
}

const formatSkuSpecs = (attrs) => {
  if (!attrs) return ''
  try {
    const arr = typeof attrs === 'string' ? JSON.parse(attrs) : attrs
    if (!Array.isArray(arr)) return ''
    return arr.map(a => `${a.name || a.key}:${a.value}`).join(', ')
  } catch { return '' }
}

watch(() => formData.value.goodID, (val) => {
  if (val) loadSkuOptions(val)
  else skuOptions.value = []
})

const openDialog = (preGoodID) => {
  type.value = 'create'
  if (preGoodID) {
    formData.value.goodID = preGoodID
  }
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { goodID: null, skuID: null, quantity: 1, unitCost: 0, totalCost: 0, supplier: '', purchaseDate: null, remark: '' }
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    let res
    switch (type.value) {
      case 'create': res = await createGoodPurchase(formData.value); break
      case 'update': res = await updateGoodPurchase(formData.value); break
      default: res = await createGoodPurchase(formData.value); break
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({ type: 'success', message: '操作成功' })
      closeDialog()
      getTableData()
    }
  })
}

const detailFrom = ref({})
const detailShow = ref(false)

const getDetails = async (row) => {
  detailFrom.value = row
  detailShow.value = true
}

const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}
</script>
