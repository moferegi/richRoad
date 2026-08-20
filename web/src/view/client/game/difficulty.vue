<template>
  <div>
    <div class="gva-search-box">
      <div class="lang-switch-row">
        <span>展示语言：</span>
        <el-select v-model="displayLang" style="width: 140px">
          <el-option v-for="lang in displayLanguageOptions" :key="lang.value" :label="lang.label" :value="lang.value" />
        </el-select>
      </div>
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" style="margin-top: 10px">
        <el-form-item label="游戏">
          <el-select v-model="searchInfo.gameID" placeholder="选择游戏" clearable @change="onSearch">
            <el-option v-for="item in gameList" :key="item.ID" :label="formatI18nText(item.name)" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSearch">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" :disabled="!searchInfo.gameID" @click="openDialog()">新增难度分类</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" border>
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="名称" min-width="180">
          <template #default="scope">
            {{ formatI18nText(scope.row.name) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="80" />
        <el-table-column align="left" label="状态" prop="status" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openDialog(scope.row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="700px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="名称">
          <MultiLangEditor
            :model="form.nameI18n"
            title="难度分类多语言名称"
            input-type="input"
            :rows="2"
            :use-tabs="true"
          />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.statusBool" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCategoryList } from '@/api/client/game'
import { getDifficultyCategoryList, createDifficultyCategory, updateDifficultyCategory, deleteDifficultyCategory } from '@/api/client/game'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'

defineOptions({ name: 'GameDifficulty' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const gameList = ref([])
const searchInfo = reactive({ gameID: '' })

const displayLang = ref('zh')
const displayLanguageOptions = [
  { value: 'zh', label: '中文' },
  { value: 'en', label: 'English' },
  { value: 'mn', label: 'Монгол' },
  { value: 'th', label: 'ไทย' },
  { value: 'hi', label: 'हिन्दी' },
  { value: 'id', label: 'Bahasa' },
  { value: 'vi', label: 'Tiếng Việt' },
  { value: 'ar', label: 'العربية' },
  { value: 'ja', label: '日本語' },
  { value: 'ko', label: '한국어' },
  { value: 'ms', label: 'Melayu' }
]

const normalizeI18nObject = (raw) => {
  if (!raw) return { zh: '' }
  if (typeof raw === 'object') return Object.keys(raw).length > 0 ? { ...raw } : { zh: '' }
  const text = String(raw).trim()
  if (!text) return { zh: '' }
  try {
    const parsed = JSON.parse(text)
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      return Object.keys(parsed).length > 0 ? { ...parsed } : { zh: '' }
    }
    return { zh: text }
  } catch { return { zh: text } }
}

const stringifyI18nObject = (i18nObject) => {
  const source = i18nObject && typeof i18nObject === 'object' ? i18nObject : {}
  const cleaned = {}
  for (const [key, value] of Object.entries(source)) {
    const lang = String(key || '').trim()
    if (!lang) continue
    const text = String(value ?? '').trim()
    if (text) cleaned[lang] = text
  }
  return JSON.stringify(Object.keys(cleaned).length > 0 ? cleaned : { zh: '' })
}

const formatI18nText = (raw) => {
  if (!raw) return ''
  try {
    const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
    if (typeof parsed === 'object' && parsed !== null) {
      const lang = String(displayLang.value || 'zh')
      return parsed[lang] || parsed.zh || parsed.en || Object.values(parsed)[0] || ''
    }
    return String(raw)
  } catch { return String(raw) }
}

const dialogVisible = ref(false)
const dialogTitle = ref('新增难度分类')
const formRef = ref(null)
const form = reactive({
  ID: 0,
  nameI18n: { zh: '' },
  sort: 0,
  statusBool: true
})
const rules = {}

const loadGameList = async () => {
  const res = await getCategoryList({ page: 1, pageSize: 100 })
  if (res.code === 0) gameList.value = res.data.list
}

const getTableData = async () => {
  if (!searchInfo.gameID) { tableData.value = []; total.value = 0; return }
  const params = { gameID: searchInfo.gameID, page: page.value, pageSize: pageSize.value }
  const res = await getDifficultyCategoryList(params)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const onSearch = () => { page.value = 1; getTableData() }
const onReset = () => { searchInfo.gameID = ''; page.value = 1; getTableData() }

const handleCurrentChange = (val) => { page.value = val; getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; page.value = 1; getTableData() }

const openDialog = (row) => {
  if (row) {
    dialogTitle.value = '编辑难度分类'
    form.ID = row.ID
    form.sort = row.sort || 0
    form.statusBool = row.status === 1
    form.nameI18n = normalizeI18nObject(row.name)
  } else {
    dialogTitle.value = '新增难度分类'
    form.ID = 0
    form.nameI18n = { zh: '' }
    form.sort = 0
    form.statusBool = true
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  const data = {
    gameID: Number(searchInfo.gameID),
    name: stringifyI18nObject(form.nameI18n),
    sort: form.sort,
    status: form.statusBool ? 1 : 0
  }
  if (form.ID) {
    data.ID = form.ID
    const res = await updateDifficultyCategory(data)
    if (res.code === 0) { ElMessage.success('更新成功'); dialogVisible.value = false; getTableData() }
  } else {
    const res = await createDifficultyCategory(data)
    if (res.code === 0) { ElMessage.success('创建成功'); dialogVisible.value = false; getTableData() }
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该难度分类吗？删除后关联的关卡也会被删除。', '提示', { type: 'warning' })
  const res = await deleteDifficultyCategory({ ID: row.ID })
  if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
}

loadGameList()
getTableData()
</script>

<style scoped>
.lang-switch-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 14px;
  color: #606266;
}
</style>