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
        <el-table-column align="left" label="国家名称" width="200">
          <template #default="scope">{{ parseI18n(scope.row.countryName) }}</template>
        </el-table-column>
        <el-table-column align="left" label="区号" prop="areaCode" width="100" />
        <el-table-column align="left" label="手机号正则" prop="phoneRegex" width="200" show-overflow-tooltip />
        <el-table-column align="left" label="国旗" width="80">
          <template #default="scope">
            <el-image v-if="scope.row.flagIcon" style="width: 30px; height: 20px" :src="getUrl(scope.row.flagIcon)" fit="cover" />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="80" />
        <el-table-column align="left" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'" size="small">{{ scope.row.isEnabled ? '启用' : '禁用' }}</el-tag>
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
            <el-form-item label="国家名称(多语言):" prop="countryName">
              <MultiLangEditor
                :model="countryNameI18n"
                :languages="enabledLangs"
                title="国家名称多语言"
              />
            </el-form-item>
            <el-form-item label="区号:" prop="areaCode">
              <el-input v-model="formData.areaCode" placeholder="例: +976" />
            </el-form-item>
            <el-form-item label="手机号正则:" prop="phoneRegex">
              <el-input v-model="formData.phoneRegex" placeholder="例: ^[0-9]{8}$" />
            </el-form-item>
            <el-form-item label="国旗图标:" prop="flagIcon">
              <SelectImage
                v-model="formData.flagIcon"
                file-type="image"
                :default-folder="PHONE_AREA_CODE_UPLOAD_FOLDER"
                :fixed-upload-folder="true"
              />
            </el-form-item>
            <el-form-item label="排序:" prop="sort">
              <el-input-number v-model="formData.sort" :min="0" />
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
  createPhoneAreaCode,
  deletePhoneAreaCode,
  updatePhoneAreaCode,
  getPhoneAreaCodeList
} from '@/api/client/phoneAreaCode'
import { getEnabledLanguages } from '@/api/client/language'
import { getUrl } from '@/utils/image'
import SelectImage from '@/components/selectImage/selectImage.vue'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, computed, onMounted } from 'vue'

defineOptions({ name: 'PhoneAreaCode' })

const PHONE_AREA_CODE_UPLOAD_FOLDER = 'cloth-on/web-else/guojiquhao'

const btnLoading = ref(false)

const enabledLangs = ref([])

const DISPLAY_LANG_PRIORITY = ['zh', 'zh-TW', 'en', 'mn', 'th', 'hi', 'id', 'vi', 'ar', 'ja', 'ko', 'ms']

const multilingualLangs = computed(() => {
  return (enabledLangs.value || []).map((lang) => ({
    code: String(lang.code || '').trim(),
    label: `${lang.name || lang.nativeName || lang.code || ''} ${lang.code || ''}`.trim()
  })).filter((lang) => lang.code)
})

const loadEnabledLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0 && res.data) {
      const list = Array.isArray(res.data) ? res.data : (res.data.list || [])
      enabledLangs.value = list.length ? list : [{ code: 'zh', name: '中文' }]
    } else {
      enabledLangs.value = [{ code: 'zh', name: '中文' }]
    }
  } catch (e) {
    enabledLangs.value = [{ code: 'zh', name: '中文' }]
  }
  applyCountryNameToEditor(formData.value.countryName)
}

onMounted(() => {
  loadEnabledLangs()
})

const parseCountryNameObject = (value) => {
  if (!value) {
    return {}
  }
  if (typeof value === 'object') {
    return value
  }
  if (typeof value !== 'string') {
    return {}
  }
  try {
    const parsed = JSON.parse(value)
    return typeof parsed === 'object' && parsed !== null ? parsed : {}
  } catch {
    const normalized = value.trim()
    return normalized ? { zh: normalized } : {}
  }
}

const countryNameI18n = ref({})

const applyCountryNameToEditor = (value) => {
  const parsed = parseCountryNameObject(value)
  const mapped = {}
  Object.entries(parsed).forEach(([rawCode, rawText]) => {
    const code = String(rawCode || '').trim()
    if (!code) return
    mapped[code] = String(rawText ?? '')
  })
  for (const item of multilingualLangs.value) {
    if (!Object.prototype.hasOwnProperty.call(mapped, item.code)) {
      mapped[item.code] = ''
    }
  }
  countryNameI18n.value = mapped
}

const getCountryNamePayload = () => {
  return Object.entries(countryNameI18n.value || {}).reduce((acc, [rawCode, rawText]) => {
    const code = String(rawCode || '').trim()
    if (!code) {
      return acc
    }
    const text = String(rawText ?? '').trim()
    if (text) {
      acc[code] = text
    }
    return acc
  }, {})
}

const parseI18n = (value) => {
  const obj = parseCountryNameObject(value)
  for (const code of DISPLAY_LANG_PRIORITY) {
    const text = String(obj[code] ?? '').trim()
    if (text) {
      return text
    }
  }
  for (const code of multilingualLangs.value.map(item => item.code)) {
    const text = String(obj[code] ?? '').trim()
    if (text) {
      return text
    }
  }
  const firstText = Object.values(obj).find(item => String(item ?? '').trim())
  if (firstText) {
    return String(firstText)
  }
  return typeof value === 'string' ? value : ''
}

const defaultForm = () => ({
  countryName: '',
  areaCode: '',
  phoneRegex: '',
  flagIcon: '',
  sort: 0,
  isEnabled: true,
})

const formData = ref(defaultForm())

const validateCountryName = (_rule, _value, callback) => {
  if (Object.keys(getCountryNamePayload()).length === 0) {
    callback(new Error('请至少填写一种语言的国家名称'))
    return
  }
  callback()
}

const rule = reactive({
  countryName: [{ validator: validateCountryName, trigger: 'change' }],
  areaCode: [{ required: true, message: '请输入区号', trigger: 'blur' }],
})

const elFormRef = ref()

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
  const table = await getPhoneAreaCodeList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
      for (const id of IDs) { await deletePhoneAreaCode({ ID: id }) }
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    })
}

const type = ref('')
const updateFunc = async(row) => {
  type.value = 'update'
  formData.value = { ...row }
  applyCountryNameToEditor(row.countryName)
  dialogFormVisible.value = true
}
const deleteFunc = async (row) => {
  const res = await deletePhoneAreaCode({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '删除成功' })
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const dialogFormVisible = ref(false)
const openDialog = () => {
  type.value = 'create'
  formData.value = defaultForm()
  applyCountryNameToEditor(formData.value.countryName)
  dialogFormVisible.value = true
}
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = defaultForm()
  applyCountryNameToEditor(formData.value.countryName)
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    formData.value.countryName = JSON.stringify(getCountryNamePayload())
    let res
    switch (type.value) {
      case 'create': res = await createPhoneAreaCode(formData.value); break
      case 'update': res = await updatePhoneAreaCode(formData.value); break
      default: res = await createPhoneAreaCode(formData.value); break
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

<style scoped>
.country-name-grid {
  width: 100%;
  display: grid;
  gap: 10px;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
}

.country-name-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.country-name-label {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.country-name-tip {
  margin-top: 8px;
}
</style>
