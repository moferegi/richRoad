<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
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
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="名称" prop="name" width="150" />
        <el-table-column label="收款码图片" width="120">
            <template #default="scope">
              <el-image style="width: 60px; height: 60px" :src="getUrl(scope.row.image)" fit="cover" :preview-src-list="[getUrl(scope.row.image)]" preview-teleported />
            </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="80" />
        <el-table-column align="left" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'">{{ scope.row.isEnabled ? '启用' : '禁用' }}</el-tag>
          </template>
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
            <el-form-item label="名称:" prop="name">
              <el-input v-model="formData.name" placeholder="请输入名称" />
            </el-form-item>
            <el-form-item label="名称(多语言):">
              <div style="width:100%">
                <div v-for="lang in enabledLangs" :key="lang.code" style="display:flex;align-items:center;margin-bottom:8px;">
                  <el-tag size="small" style="margin-right:8px;min-width:50px;text-align:center;">{{ lang.code }}</el-tag>
                  <el-input v-model="nameI18nObj[lang.code]" :placeholder="lang.name" style="flex:1" />
                </div>
                <div v-if="!enabledLangs.length" style="color:#999;font-size:12px;">请先在语言管理中启用语言</div>
              </div>
            </el-form-item>
            <el-form-item label="收款码图片:" prop="image">
              <SelectImage v-model="formData.image" file-type="image" />
            </el-form-item>
            <el-form-item label="排序:" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" />
            </el-form-item>
            <el-form-item label="是否启用:" prop="isEnabled">
              <el-switch v-model="formData.isEnabled" />
            </el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="800" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看详情">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="名称">{{ detailFrom.name }}</el-descriptions-item>
                    <el-descriptions-item label="收款码图片">
                      <el-image style="width: 100px; height: 100px" :src="getUrl(detailFrom.image)" fit="cover" />
                    </el-descriptions-item>
                    <el-descriptions-item label="排序">{{ detailFrom.sort }}</el-descriptions-item>
                    <el-descriptions-item label="状态">{{ detailFrom.isEnabled ? '启用' : '禁用' }}</el-descriptions-item>
            </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createQrcodePayment,
  deleteQrcodePayment,
  updateQrcodePayment,
  getQrcodePaymentList
} from '@/api/shop/qrcodePayment'
import { getEnabledLanguages } from '@/api/client/language'
import { getUrl } from '@/utils/image'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({ name: 'QrcodePayment' })

const btnLoading = ref(false)

const formData = ref({
  name: '',
  nameI18n: '',
  image: '',
  sort: 0,
  isEnabled: true,
})

// 多语言编辑
const nameI18nObj = reactive({})
const enabledLangs = ref([])

const loadEnabledLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0 && res.data) {
      enabledLangs.value = Array.isArray(res.data) ? res.data : (res.data.list || [])
    }
  } catch(e) {}
}

const parseNameI18n = (jsonStr) => {
  Object.keys(nameI18nObj).forEach(k => delete nameI18nObj[k])
  try {
    const parsed = JSON.parse(jsonStr || '{}')
    Object.assign(nameI18nObj, parsed)
  } catch { /* ignore */ }
}

const serializeNameI18n = () => {
  const obj = {}
  for (const lang of enabledLangs.value) {
    if (nameI18nObj[lang.code]) {
      obj[lang.code] = nameI18nObj[lang.code]
    }
  }
  return JSON.stringify(obj)
}

onMounted(() => {
  loadEnabledLangs()
})

const rule = reactive({
  name: [{ required: true, message: '请输入名称', trigger: ['input','blur'] }],
  image: [{ required: true, message: '请上传收款码图片', trigger: 'blur' }],
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
  const table = await getQrcodePaymentList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
      for (const id of IDs) { await deleteQrcodePayment({ ID: id }) }
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    })
}

const type = ref('')
const updateFunc = async(row) => { type.value = 'update'; formData.value = { ...row }; parseNameI18n(row.nameI18n); dialogFormVisible.value = true }
const deleteFunc = async (row) => {
  const res = await deleteQrcodePayment({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '删除成功' })
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const dialogFormVisible = ref(false)
const openDialog = () => { type.value = 'create'; parseNameI18n('{}'); dialogFormVisible.value = true }
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = { name: '', nameI18n: '', image: '', sort: 0, isEnabled: true }
  parseNameI18n('{}')
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    formData.value.nameI18n = serializeNameI18n()
    let res
    switch (type.value) {
      case 'create': res = await createQrcodePayment(formData.value); break
      case 'update': res = await updateQrcodePayment(formData.value); break
      default: res = await createQrcodePayment(formData.value); break
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
const getDetails = async (row) => { detailFrom.value = row; detailShow.value = true }
const closeDetailShow = () => { detailShow.value = false; detailFrom.value = {} }
</script>
