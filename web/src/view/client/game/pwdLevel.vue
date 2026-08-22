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
          <el-select v-model="searchInfo.gameID" placeholder="选择游戏" clearable @change="onGameChange">
            <el-option v-for="item in gameList" :key="item.ID" :label="formatI18nText(item.name)" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="难度分类">
          <el-select v-model="searchInfo.categoryID" placeholder="选择难度分类" clearable @change="onSearch">
            <el-option v-for="item in diffList" :key="item.ID" :label="formatI18nText(item.name)" :value="item.ID" />
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
        <el-button type="primary" :disabled="!searchInfo.categoryID" @click="openDialog()">新增关卡</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" border>
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="关卡序号" prop="levelNumber" width="100" />
        <el-table-column align="left" label="答案" prop="answer" width="100" />
        <el-table-column align="left" label="提示数字" prop="hintDigits" width="180" />
        <el-table-column align="left" label="文字提示" min-width="300">
          <template #default="scope">
            <span>{{ formatHintTextsDisplay(scope.row) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="80" />
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="关卡序号" prop="levelNumber">
          <el-input-number v-model="form.levelNumber" :min="1" />
        </el-form-item>
        <el-form-item label="答案(4位)" prop="answer">
          <el-input v-model="form.answer" placeholder="4位数字，如: 1234" maxlength="4" />
        </el-form-item>
        <el-divider />
        <el-form-item v-for="(hint, idx) in form.hints" :key="'hint' + idx" :label="'提示' + (idx + 1)">
          <div style="display: flex; align-items: flex-start; gap: 12px; width: 100%">
            <div style="flex-shrink: 0; width: 100px">
              <div style="font-size: 12px; color: #909399; margin-bottom: 4px">提示数字</div>
              <el-input v-model="hint.digit" placeholder="数字" maxlength="4" />
            </div>
            <div style="flex: 1; min-width: 0">
              <MultiLangEditor
                :model="hint.text"
                :title="'提示' + (idx + 1) + '多语言文案'"
                input-type="input"
                :use-tabs="true"
              />
            </div>
          </div>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
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
import { getCategoryList, getDifficultyCategoryList } from '@/api/client/game'
import { getPwdLevelList, createPwdLevel, updatePwdLevel, deletePwdLevel } from '@/api/client/game'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'

defineOptions({ name: 'GamePwdLevel' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const gameList = ref([])
const diffList = ref([])
const searchInfo = reactive({ gameID: '', categoryID: '' })

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

const formatHintTextsDisplay = (row) => {
  if (!row) return ''
  const digits = (row.hintDigits || '').split(',').map(d => d.trim())
  try {
    const arr = typeof row.hintTexts === 'string' ? JSON.parse(row.hintTexts) : row.hintTexts
    if (!Array.isArray(arr)) return ''
    const lang = String(displayLang.value || 'zh')
    return arr.map((o, i) => {
      const text = (o && typeof o === 'object') ? (o[lang] || o.zh || o.en || Object.values(o)[0] || '') : ''
      return `[${digits[i] || '?'}] ${text}`
    }).join(' | ')
  } catch { return String(row.hintTexts) }
}

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

const emptyHint = () => ({ digit: '', text: { zh: '' } })

const buildHintDigits = () => {
  return form.hints.map(h => (h.digit || '').trim()).join(',')
}

const buildHintTexts = () => {
  const arr = form.hints.map(h => {
    const cleaned = {}
    const source = h.text && typeof h.text === 'object' ? h.text : {}
    for (const [key, value] of Object.entries(source)) {
      const lang = String(key || '').trim()
      if (!lang) continue
      const text = String(value ?? '').trim()
      if (text) cleaned[lang] = text
    }
    return Object.keys(cleaned).length > 0 ? cleaned : { zh: '' }
  })
  return JSON.stringify(arr)
}

const parseHints = (digitsStr, textsStr) => {
  const digits = (digitsStr || '').split(',').map(d => d.trim())
  let texts = []
  try {
    const arr = JSON.parse(textsStr || '[]')
    if (Array.isArray(arr)) texts = arr
  } catch { /* ignore */ }
  form.hints = []
  for (let i = 0; i < 4; i++) {
    form.hints.push({
      digit: digits[i] || '',
      text: normalizeI18nObject(texts[i])
    })
  }
}

const dialogVisible = ref(false)
const dialogTitle = ref('新增关卡')
const formRef = ref(null)
const form = reactive({
  ID: 0,
  levelNumber: 1,
  answer: '',
  hints: [emptyHint(), emptyHint(), emptyHint(), emptyHint()],
  sort: 0
})
const rules = {
  levelNumber: [{ required: true, message: '请输入关卡序号', trigger: 'blur' }],
  answer: [{ required: true, message: '请输入4位答案', trigger: 'blur' }]
}

const loadGameList = async () => {
  const res = await getCategoryList({ page: 1, pageSize: 100 })
  if (res.code === 0) gameList.value = res.data.list
}

const onGameChange = async (val) => {
  searchInfo.categoryID = ''
  diffList.value = []
  if (!val) { getTableData(); return }
  const res = await getDifficultyCategoryList({ gameID: val, page: 1, pageSize: 100 })
  if (res.code === 0) diffList.value = res.data.list
  getTableData()
}

const getTableData = async () => {
  if (!searchInfo.categoryID) { tableData.value = []; total.value = 0; return }
  const params = { categoryID: searchInfo.categoryID, page: page.value, pageSize: pageSize.value }
  const res = await getPwdLevelList(params)
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const onSearch = () => { page.value = 1; getTableData() }
const onReset = () => { searchInfo.gameID = ''; searchInfo.categoryID = ''; diffList.value = []; page.value = 1; getTableData() }

const handleCurrentChange = (val) => { page.value = val; getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; page.value = 1; getTableData() }

const openDialog = (row) => {
  if (row) {
    dialogTitle.value = '编辑关卡'
    form.ID = row.ID
    form.levelNumber = row.levelNumber
    form.answer = row.answer
    form.sort = row.sort
    parseHints(row.hintDigits, row.hintTexts)
  } else {
    dialogTitle.value = '新增关卡'
    form.ID = 0
    form.levelNumber = 1
    form.answer = ''
    form.sort = 0
    form.hints = [emptyHint(), emptyHint(), emptyHint(), emptyHint()]
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  const data = {
    categoryID: Number(searchInfo.categoryID),
    levelNumber: form.levelNumber,
    answer: form.answer,
    hintDigits: buildHintDigits(),
    hintTexts: buildHintTexts(),
    sort: form.sort
  }
  if (form.ID) {
    data.ID = form.ID
    const res = await updatePwdLevel(data)
    if (res.code === 0) { ElMessage.success('更新成功'); dialogVisible.value = false; getTableData() }
  } else {
    const res = await createPwdLevel(data)
    if (res.code === 0) { ElMessage.success('创建成功'); dialogVisible.value = false; getTableData() }
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该关卡吗？', '提示', { type: 'warning' })
  const res = await deletePwdLevel({ ID: row.ID })
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