<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="类型">
          <el-select v-model="searchInfo.type" placeholder="全部" clearable style="width:120px">
            <el-option label="规格(spec)" value="spec" />
            <el-option label="属性(attr)" value="attr" />
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
        max-height="70vh"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="70" sortable />
        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="类型" prop="type" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'spec' ? 'primary' : 'success'">
              {{ scope.row.type === 'spec' ? '规格' : '属性' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="名称(默认)" prop="label" width="160" />
        <el-table-column align="left" label="名称(多语言)" width="260">
          <template #default="scope">{{ formatI18n(scope.row.labelI18n) }}</template>
        </el-table-column>
        <el-table-column align="left" label="值(默认)" prop="value" width="160" />
        <el-table-column align="left" label="值(多语言)" width="260">
          <template #default="scope">{{ formatI18n(scope.row.valueI18n) }}</template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="80" sortable />
        <el-table-column align="left" label="操作" fixed="right" min-width="160">
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

    <!-- 新增/编辑抽屉 -->
    <el-drawer destroy-on-close size="520px" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ dialogType === 'create' ? '新增' : '编辑' }}规格/属性</span>
          <div>
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rules" label-width="80px">
        <el-form-item label="类型:" prop="type">
          <el-select v-model="formData.type" placeholder="请选择类型" style="width:100%">
            <el-option label="规格(spec)" value="spec" />
            <el-option label="属性(attr)" value="attr" />
          </el-select>
        </el-form-item>
        <el-form-item label="名称(默认):" prop="label">
          <el-input v-model="formData.label" placeholder="请输入默认名称，如：颜色" />
        </el-form-item>
        <el-form-item label="名称(多语言):">
          <MultiLangEditor
            :model="labelI18nObj"
            :languages="enabledLangs"
            title="名称多语言"
          />
        </el-form-item>
        <el-form-item label="值(默认):" prop="value">
          <el-input v-model="formData.value" placeholder="请输入默认值，如：红色" />
        </el-form-item>
        <el-form-item label="值(多语言):">
          <MultiLangEditor
            :model="valueI18nObj"
            :languages="enabledLangs"
            title="值多语言"
          />
        </el-form-item>
        <el-form-item label="排序:" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" :max="9999" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createSkuSpec, deleteSkuSpec, deleteSkuSpecByIds, updateSkuSpec, findSkuSpec, getSkuSpecList
} from '@/api/shop/skuSpec'
import { getEnabledLanguages } from '@/api/client/language'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'

defineOptions({ name: 'SkuSpec' })

const btnLoading = ref(false)

// === i18n ===
const enabledLangs = ref([])
const labelI18nObj = reactive({})
const valueI18nObj = reactive({})

const loadEnabledLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0 && res.data) {
      enabledLangs.value = Array.isArray(res.data) ? res.data : (res.data.list || [])
    }
  } catch(e) { /* ignore */ }
}

const parseI18nToObj = (jsonStr, target) => {
  Object.keys(target).forEach(k => delete target[k])
  try {
    const parsed = typeof jsonStr === 'string' ? JSON.parse(jsonStr || '{}') : (jsonStr || {})
    Object.assign(target, parsed)
  } catch { /* ignore */ }
}

const serializeI18nObj = (obj) => {
  const result = {}
  Object.entries(obj || {}).forEach(([rawCode, rawText]) => {
    const code = String(rawCode || '').trim()
    if (!code) return
    const text = String(rawText ?? '').trim()
    if (!text) return
    result[code] = text
  })
  return JSON.stringify(result)
}

const formatI18n = (jsonStr) => {
  try {
    const obj = typeof jsonStr === 'string' ? JSON.parse(jsonStr || '{}') : (jsonStr || {})
    const parts = Object.entries(obj).filter(([,v]) => v).map(([k,v]) => `${k}:${v}`)
    return parts.length ? parts.join(' | ') : '-'
  } catch { return '-' }
}

onMounted(() => { loadEnabledLangs() })

// === 表单 ===
const formData = ref({
  type: 'spec',
  label: '',
  labelI18n: '',
  value: '',
  valueI18n: '',
  sort: 0,
})

const rules = reactive({
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  label: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  value: [{ required: true, message: '请输入值', trigger: 'blur' }],
})

// === 表格 ===
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const elFormRef = ref()
const elSearchFormRef = ref()

const sortChange = ({ prop, order }) => {
  searchInfo.value.sort = prop
  searchInfo.value.order = order
  getTableData()
}

const onReset = () => { searchInfo.value = {}; getTableData() }
const onSubmit = () => { page.value = 1; getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }

const getTableData = async () => {
  const res = await getSkuSpecList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}
getTableData()

// === 多选 ===
const multipleSelection = ref([])
const handleSelectionChange = (val) => { multipleSelection.value = val }

// === 删除 ===
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const res = await deleteSkuSpec({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        if (tableData.value.length === 1 && page.value > 1) page.value--
        getTableData()
      }
    })
}

const onDelete = () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const IDs = multipleSelection.value.map(item => item.ID)
      if (!IDs.length) { ElMessage({ type: 'warning', message: '请选择要删除的数据' }); return }
      const res = await deleteSkuSpecByIds({ IDs })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        if (tableData.value.length === IDs.length && page.value > 1) page.value--
        getTableData()
      }
    })
}

// === 新增/编辑 ===
const dialogType = ref('create')
const dialogFormVisible = ref(false)

const openDialog = () => {
  dialogType.value = 'create'
  formData.value = { type: 'spec', label: '', labelI18n: '', value: '', valueI18n: '', sort: 0 }
  parseI18nToObj('{}', labelI18nObj)
  parseI18nToObj('{}', valueI18nObj)
  dialogFormVisible.value = true
}

const updateFunc = async (row) => {
  const res = await findSkuSpec({ ID: row.ID })
  if (res.code === 0) {
    dialogType.value = 'update'
    formData.value = res.data
    parseI18nToObj(res.data.labelI18n, labelI18nObj)
    parseI18nToObj(res.data.valueI18n, valueI18nObj)
    dialogFormVisible.value = true
  }
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { type: 'spec', label: '', labelI18n: '', value: '', valueI18n: '', sort: 0 }
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) { btnLoading.value = false; return }
    formData.value.labelI18n = serializeI18nObj(labelI18nObj)
    formData.value.valueI18n = serializeI18nObj(valueI18nObj)
    let res
    if (dialogType.value === 'create') {
      res = await createSkuSpec(formData.value)
    } else {
      res = await updateSkuSpec(formData.value)
    }
    btnLoading.value = false
    if (res.code === 0) {
      ElMessage({ type: 'success', message: dialogType.value === 'create' ? '新增成功' : '编辑成功' })
      closeDialog()
      getTableData()
    }
  })
}
</script>

<style scoped>
</style>
