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

        <el-form-item label="客户端类型">
          <el-select v-model="searchInfo.clientType" placeholder="全部" clearable>
            <el-option label="uni端" value="uni" />
            <el-option label="web端" value="web" />
            <el-option label="全部" value="all" />
          </el-select>
        </el-form-item>

        <el-form-item label="弹窗类型">
          <el-select v-model="searchInfo.popupType" placeholder="全部" clearable>
            <el-option label="仅图片" value="image" />
            <el-option label="编辑内容" value="content" />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="searchInfo.isEnabled" placeholder="全部" clearable>
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
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
        <el-table-column align="left" label="ID" prop="ID" width="60" />
        <el-table-column align="left" label="日期" width="180" sortable="custom" prop="created_at">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="类型" width="90">
          <template #default="scope">
            <el-tag :type="scope.row.popupType === 'content' ? 'warning' : 'primary'" size="small">
              {{ scope.row.popupType === 'content' ? '编辑内容' : '仅图片' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="标题" width="150">
          <template #default="scope">{{ parseTitle(scope.row.title) }}</template>
        </el-table-column>
        <el-table-column label="图片" width="100">
            <template #default="scope">
              <el-image v-if="scope.row.externalPath || scope.row.image" style="width: 50px; height: 50px" :src="getUrl(scope.row.externalPath || scope.row.image)" fit="cover" :preview-src-list="[getUrl(scope.row.externalPath || scope.row.image)]" preview-teleported />
              <span v-else>-</span>
            </template>
        </el-table-column>
        <el-table-column align="left" label="跳转链接" prop="link" width="150" show-overflow-tooltip />
        <el-table-column align="left" label="展示页面" min-width="200">
          <template #default="scope">
            <template v-if="scope.row.pages">
              <el-tag v-for="p in scope.row.pages.split(',')" :key="p" size="small" class="mr-1 mb-1">{{ p }}</el-tag>
            </template>
            <el-tag v-else-if="scope.row.position" size="small" type="info">{{ scope.row.position }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="客户端" width="80">
          <template #default="scope">{{ clientTypeMap[scope.row.clientType] || scope.row.clientType }}</template>
        </el-table-column>
        <el-table-column align="left" label="生效时间" width="180">
          <template #default="scope">{{ scope.row.startTime ? formatDate(scope.row.startTime) : '不限' }}</template>
        </el-table-column>
        <el-table-column align="left" label="失效时间" width="180">
          <template #default="scope">{{ scope.row.endTime ? formatDate(scope.row.endTime) : '不限' }}</template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="60" />
        <el-table-column align="left" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.isEnabled ? 'success' : 'danger'">{{ scope.row.isEnabled ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
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
            <!-- 弹窗类型 -->
            <el-form-item label="弹窗类型:" prop="popupType">
              <el-radio-group v-model="formData.popupType">
                <el-radio value="image">仅图片</el-radio>
                <el-radio value="content">编辑内容</el-radio>
              </el-radio-group>
            </el-form-item>

            <!-- 标题(多语言) -->
            <el-form-item label="标题(多语言):" prop="title">
              <el-input v-model="formData.title" placeholder="默认标题" class="mb-2" />
              <MultiLangEditor
                :model="titleI18n"
                :languages="enabledLangs"
                title="弹窗标题多语言"
              />
            </el-form-item>

            <!-- 弹窗图片 -->
            <el-form-item label="弹窗图片(上传):" prop="image">
              <SelectImage v-model="formData.image" file-type="image" />
            </el-form-item>

            <el-form-item label="外部图片路径(优先于上传):" prop="externalPath">
              <el-input v-model="formData.externalPath" :clearable="true" placeholder="https://example.com/popup.jpg" />
            </el-form-item>

            <!-- 富文本内容(多语言) - 仅 content 类型 -->
            <el-form-item v-if="formData.popupType === 'content'" label="弹窗内容(多语言):" prop="content">
              <MultiLangEditor
                :model="contentI18n"
                :languages="enabledLangs"
                title="弹窗内容多语言"
                :use-tabs="true"
              >
                <template #editor="{ lang }">
                  <div class="h-[460px]">
                    <RichEdit v-model="contentI18n[lang.code]" />
                  </div>
                </template>
              </MultiLangEditor>
            </el-form-item>

            <!-- 跳转链接 -->
            <el-form-item label="跳转链接:" prop="link">
              <el-input v-model="formData.link" placeholder="点击弹窗跳转的链接" />
            </el-form-item>

            <!-- 展示页面路径(标签输入) -->
            <el-form-item label="展示页面路径:" prop="pages">
              <div style="width:100%;">
                <div class="mb-2">
                  <el-tag v-for="(p, idx) in pagesList" :key="idx" closable class="mr-1 mb-1" @close="removePageTag(idx)">{{ p }}</el-tag>
                </div>
                <div class="flex gap-2">
                  <el-select v-model="newPagePath" filterable allow-create default-first-option placeholder="选择或输入页面路径" style="flex:1;" @keyup.enter="addPageTag">
                    <el-option v-for="p in commonPages" :key="p.value" :label="p.label" :value="p.value" />
                  </el-select>
                  <el-button type="primary" @click="addPageTag">添加</el-button>
                </div>
                <div class="text-xs text-gray-400 mt-1">选择"all"表示所有页面生效。也可手动输入自定义路径。</div>
              </div>
            </el-form-item>

            <!-- 时间范围 -->
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="生效时间:" prop="startTime">
                  <el-date-picker v-model="formData.startTime" type="datetime" placeholder="留空不限" style="width:100%;" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="失效时间:" prop="endTime">
                  <el-date-picker v-model="formData.endTime" type="datetime" placeholder="留空不限" style="width:100%;" />
                </el-form-item>
              </el-col>
            </el-row>

            <!-- 客户端类型 -->
            <el-form-item label="客户端类型:" prop="clientType">
              <el-select v-model="formData.clientType">
                <el-option label="全部" value="all" />
                <el-option label="uni端" value="uni" />
                <el-option label="web端" value="web" />
              </el-select>
            </el-form-item>

            <!-- 开关选项 -->
            <el-row :gutter="16">
              <el-col :span="6">
                <el-form-item label="只弹一次:" prop="onceOnly">
                  <el-switch v-model="formData.onceOnly" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="可关闭:" prop="closeable">
                  <el-switch v-model="formData.closeable" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="排序:" prop="sort">
                  <el-input-number v-model="formData.sort" :min="0" />
                </el-form-item>
              </el-col>
              <el-col :span="6">
                <el-form-item label="启用:" prop="isEnabled">
                  <el-switch v-model="formData.isEnabled" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createPopup,
  deletePopup,
  updatePopup,
  getPopupList
} from '@/api/shop/popup'
import { getUrl } from '@/utils/image'
import { getEnabledLanguages } from '@/api/client/language'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import RichEdit from '@/components/richtext/rich-edit.vue'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted, watch } from 'vue'

defineOptions({ name: 'Popup' })

const btnLoading = ref(false)
const clientTypeMap = { uni: 'uni端', web: 'web端', all: '全部' }

// 常用页面路径选项
const commonPages = [
  { label: '所有页面', value: 'all' },
  { label: '首页', value: '/pages/tabBar/index' },
  { label: '商品详情', value: '/pages/goodsDetails/goodsDetails' },
  { label: '分类页', value: '/pages/tabBar/category' },
  { label: '购物车', value: '/pages/tabBar/cart' },
  { label: '我的', value: '/pages/tabBar/mine' },
  { label: '订单列表', value: '/pages/order/orderList' },
  { label: '预售专区', value: '/pages/presale/presale' },
]

// === 多语言支持 ===
const enabledLangs = ref([])
const titleI18n = ref({})
const contentI18n = ref({})

const loadLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0) {
      enabledLangs.value = res.data || []
    }
  } catch (e) { /* ignore */ }
}

const parseI18nJson = (jsonStr) => {
  if (!jsonStr) return {}
  try { return JSON.parse(jsonStr) } catch { return {} }
}

const serializeI18nJson = (obj) => {
  const filtered = {}
  for (const [k, v] of Object.entries(obj)) {
    if (v) filtered[k] = v
  }
  return Object.keys(filtered).length ? JSON.stringify(filtered) : ''
}

const parseTitle = (title) => {
  if (!title) return ''
  try {
    const obj = JSON.parse(title)
    return obj.zh || obj.en || Object.values(obj)[0] || title
  } catch { return title }
}

onMounted(() => { loadLangs() })

// === 页面路径标签管理 ===
const pagesList = ref([])
const newPagePath = ref('')

const addPageTag = () => {
  const v = (newPagePath.value || '').trim()
  if (!v) return
  if (!pagesList.value.includes(v)) {
    pagesList.value.push(v)
  }
  newPagePath.value = ''
}

const removePageTag = (idx) => {
  pagesList.value.splice(idx, 1)
}

// 同步 pagesList <-> formData.pages
watch(pagesList, (val) => {
  formData.value.pages = val.join(',')
}, { deep: true })

const defaultForm = () => ({
  title: '',
  image: '',
  externalPath: '',
  link: '',
  content: '',
  popupType: 'image',
  pages: '',
  startTime: null,
  endTime: null,
  clientType: 'all',
  onceOnly: false,
  sort: 0,
  closeable: true,
  position: '',
  isEnabled: true,
})

const formData = ref(defaultForm())

const rule = reactive({})

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

const onReset = () => { searchInfo.value = {}; getTableData() }

// 排序
const sortChange = ({ prop, order }) => {
  if (order) {
    searchInfo.value.orderBy = prop
    searchInfo.value.orderDir = order === 'ascending' ? 'asc' : 'desc'
  } else {
    searchInfo.value.orderBy = ''
    searchInfo.value.orderDir = ''
  }
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
  const table = await getPopupList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
      for (const id of IDs) { await deletePopup({ ID: id }) }
      ElMessage({ type: 'success', message: '删除成功' })
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    })
}

const type = ref('')

const updateFunc = async(row) => {
  type.value = 'update'
  formData.value = { ...row }
  // 解析多语言字段
  titleI18n.value = parseI18nJson(row.title)
  contentI18n.value = parseI18nJson(row.content)
  // 解析页面路径
  pagesList.value = row.pages ? row.pages.split(',').map(s => s.trim()).filter(Boolean) : []
  dialogFormVisible.value = true
}

const deleteFunc = async (row) => {
  const res = await deletePopup({ ID: row.ID })
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
  titleI18n.value = {}
  contentI18n.value = {}
  pagesList.value = []
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = defaultForm()
  titleI18n.value = {}
  contentI18n.value = {}
  pagesList.value = []
}

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
    // 序列化多语言字段
    if (enabledLangs.value.length) {
      formData.value.title = serializeI18nJson(titleI18n.value) || formData.value.title
      if (formData.value.popupType === 'content') {
        formData.value.content = serializeI18nJson(contentI18n.value) || formData.value.content
      }
    }
    // pages 从 pagesList 同步（watch 已处理，这里确保）
    formData.value.pages = pagesList.value.join(',')

    let res
    switch (type.value) {
      case 'create': res = await createPopup(formData.value); break
      case 'update': res = await updatePopup(formData.value); break
      default: res = await createPopup(formData.value); break
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
