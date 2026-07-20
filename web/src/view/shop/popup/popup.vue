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
              <el-tag v-for="p in getDisplayPageTags(scope.row.pages)" :key="p.value" size="small" class="mr-1 mb-1">{{ p.text }}</el-tag>
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
              <FileUploadWithDir
                v-model="formData.image"
                :default-folder="POPUP_UPLOAD_FOLDER"
                :fixed-upload-folder="false"
                accept="image/*"
              />
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
                    <RichEdit
                      v-model="contentI18n[lang.code]"
                      :upload-folder="POPUP_DETAIL_UPLOAD_FOLDER"
                    />
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
                <div class="flex gap-2 mb-2">
                  <el-input v-model="newPageName" placeholder="名称(可选)，如：首页" style="width: 220px;" />
                  <el-select v-model="newPagePath" filterable allow-create default-first-option placeholder="选择或输入页面路径" style="flex:1;" @keyup.enter="addPageEntry">
                    <el-option v-for="p in pageOptions" :key="p.value" :label="`${p.label} (${p.value})`" :value="p.value">
                      <div class="flex justify-between items-center">
                        <span>{{ p.label }}</span>
                        <span class="text-xs text-gray-400 ml-3">{{ p.value }}</span>
                      </div>
                    </el-option>
                  </el-select>
                  <el-button type="primary" @click="addPageEntry">添加</el-button>
                </div>

                <el-table v-if="pageEntries.length" :data="pageEntries" size="small" border class="mb-2">
                  <el-table-column label="名称" min-width="150">
                    <template #default="scope">
                      <template v-if="editingPageIndex === scope.$index">
                        <el-input v-model="editingPageName" placeholder="名称(可选)" />
                      </template>
                      <template v-else>
                        {{ scope.row.name || '-' }}
                      </template>
                    </template>
                  </el-table-column>
                  <el-table-column label="路径" min-width="320">
                    <template #default="scope">
                      <template v-if="editingPageIndex === scope.$index">
                        <el-select v-model="editingPagePath" filterable allow-create default-first-option placeholder="选择或输入页面路径" style="width:100%;">
                          <el-option v-for="p in pageOptions" :key="p.value" :label="`${p.label} (${p.value})`" :value="p.value" />
                        </el-select>
                      </template>
                      <template v-else>
                        {{ scope.row.path }}
                      </template>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" width="220" fixed="right">
                    <template #default="scope">
                      <template v-if="editingPageIndex === scope.$index">
                        <el-button link type="primary" @click="saveEditPageEntry(scope.$index)">保存</el-button>
                        <el-button link @click="cancelEditPageEntry">取消</el-button>
                      </template>
                      <template v-else>
                        <el-button link type="primary" @click="startEditPageEntry(scope.$index)">修改</el-button>
                        <el-button link type="danger" @click="removePageEntry(scope.$index)">删除</el-button>
                      </template>
                    </template>
                  </el-table-column>
                </el-table>
                <div v-else class="text-xs text-gray-400 mb-2">暂无已配置页面路径</div>

                <div class="mb-2">
                  <el-tag v-for="(entry, idx) in pageEntries" :key="`${entry.path}_${idx}`" class="mr-1 mb-1" type="info" effect="plain">{{ formatPageEntryTag(entry) }}</el-tag>
                </div>
                <div class="text-xs text-gray-400 mt-1">选择"all"表示不按页面过滤，客户端下所有页面生效。也可手动输入自定义路径。</div>
                <div class="text-xs text-gray-500 mt-2">all 当前覆盖（系统路由 + 已配置路径）：</div>
                <div class="mt-1" v-if="allPageScopeList.length">
                  <el-tag v-for="p in allPageScopeList" :key="p" size="small" effect="plain" class="mr-1 mb-1">{{ p }}</el-tag>
                </div>
                <div class="text-xs text-gray-400 mt-1" v-else>暂无已收录路径</div>
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
  getPopupList,
  getPopupPagePathOptions
} from '@/api/shop/popup'
import { getUrl } from '@/utils/image'
import { getEnabledLanguages } from '@/api/client/language'
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import FileUploadWithDir from '@/components/FileUploadWithDir/index.vue'
import RichEdit from '@/components/richtext/rich-edit.vue'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted, watch, computed } from 'vue'

defineOptions({ name: 'Popup' })

const POPUP_UPLOAD_FOLDER = 'cloth-on/web-else/popup'
const POPUP_DETAIL_UPLOAD_FOLDER = 'cloth-on/web-else/popup/detail'

const btnLoading = ref(false)
const clientTypeMap = { uni: 'uni端', web: 'web端', all: '全部' }

// 常用页面路径选项
const commonPages = [
  { label: '所有页面', value: 'all' },
  { label: '首页', value: '/pages/tabBar/index' },
  { label: '衣服商品页', value: '/pages/tabBar/clothes/index' },
  { label: '我的', value: '/pages/tabBar/my/index' },
  { label: '商品详情', value: '/pages/goodsDetails/goodsDetails' },
  { label: '购物车', value: '/pages/cart/index' },
  { label: '订单列表', value: '/pages/order/order' },
  { label: '预售专区', value: '/pages/presale/list' },
]

const CUSTOM_PAGE_LABEL = '自定义路径'
const RUNTIME_PAGE_LABEL = '可用页面'
const customPageOptions = ref([])
const runtimePageOptions = ref([])

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

onMounted(() => {
  loadLangs()
  loadRuntimePageOptions('all')
})

// === 页面路径标签管理 ===
const PAGE_NAME_PATH_SEPARATOR = '|'

const pageEntries = ref([])
const newPageName = ref('')
const newPagePath = ref('')
const editingPageIndex = ref(-1)
const editingPageName = ref('')
const editingPagePath = ref('')

const normalizePagePathInput = (rawPath) => {
  const raw = String(rawPath || '').trim()
  if (!raw) return ''

  if (raw.toLowerCase() === 'all' || raw === '所有页面') return 'all'

  const byLabel = commonPages.find(item => item.label === raw)
  let normalized = byLabel ? byLabel.value : raw

  if (String(normalized).toLowerCase() === 'all') return 'all'
  normalized = String(normalized).replace(/\\/g, '/')
  if (!normalized.startsWith('/')) {
    normalized = '/' + normalized
  }
  normalized = normalized.replace(/\/+/g, '/')
  if (normalized.length > 1) {
    normalized = normalized.replace(/\/+$/, '')
  }
  return normalized
}

const dedupePageEntries = (entries = []) => {
  const result = []
  const seen = new Set()

  entries.forEach(item => {
    const path = normalizePagePathInput(item?.path)
    if (!path || seen.has(path)) return
    seen.add(path)
    result.push({
      name: String(item?.name || '').trim(),
      path,
    })
  })

  return result
}

const parsePageEntryToken = (token) => {
  const raw = String(token || '').trim()
  if (!raw) return null

  let name = ''
  let pathValue = raw
  if (raw.includes(PAGE_NAME_PATH_SEPARATOR)) {
    const [tokenName, tokenPath] = raw.split(PAGE_NAME_PATH_SEPARATOR)
    name = String(tokenName || '').trim()
    pathValue = String(tokenPath || '').trim()
  }

  const path = normalizePagePathInput(pathValue)
  if (!path) return null
  return { name, path }
}

const parsePageEntriesFromRaw = (rawPages) => {
  const entries = String(rawPages || '')
    .split(',')
    .map(parsePageEntryToken)
    .filter(Boolean)
  return dedupePageEntries(entries)
}

const serializePageEntry = (entry) => {
  const path = normalizePagePathInput(entry?.path)
  if (!path) return ''

  const name = String(entry?.name || '').trim()
  if (!name || path === 'all') {
    return path
  }

  return `${name}${PAGE_NAME_PATH_SEPARATOR}${path}`
}

const normalizedPathList = (rawPages) => {
  return parsePageEntriesFromRaw(rawPages).map(item => item.path)
}

const pageOptions = computed(() => {
  const merged = []
  const seen = new Set()

  const pushOption = (label, value) => {
    const normalizedValue = normalizePagePathInput(value)
    if (!normalizedValue || seen.has(normalizedValue)) return
    seen.add(normalizedValue)
    merged.push({
      label: label || (normalizedValue === 'all' ? '所有页面' : CUSTOM_PAGE_LABEL),
      value: normalizedValue,
    })
  }

  commonPages.forEach(item => pushOption(item.label, item.value))
  runtimePageOptions.value.forEach(item => pushOption(item.label, item.value))
  customPageOptions.value.forEach(item => pushOption(item.label, item.value))

  return merged
})

const allPageScopeList = computed(() => {
  return pageOptions.value
    .map(item => item.value)
    .filter(value => value !== 'all')
})

const ensureCustomPageOption = (path) => {
  const normalized = normalizePagePathInput(path)
  if (!normalized || normalized === 'all') return

  const inCommon = commonPages.some(item => normalizePagePathInput(item.value) === normalized)
  if (inCommon) return

  const inCustom = customPageOptions.value.some(item => normalizePagePathInput(item.value) === normalized)
  if (inCustom) return

  customPageOptions.value.push({
    label: CUSTOM_PAGE_LABEL,
    value: normalized,
  })
}

const resolvePageLabel = (path) => {
  const normalized = normalizePagePathInput(path)
  if (!normalized) return ''
  if (normalized === 'all') return '所有页面'
  const found = pageOptions.value.find(item => item.value === normalized)
  if (!found || found.label === RUNTIME_PAGE_LABEL) return ''
  return found.label
}

const formatPageTag = (path, name = '') => {
  const normalized = normalizePagePathInput(path)
  if (!normalized) return ''
  if (normalized === 'all') return '所有页面 (all)'

  const displayName = String(name || '').trim() || resolvePageLabel(normalized)
  const label = displayName
  if (!label || label === CUSTOM_PAGE_LABEL || label === normalized) {
    return normalized
  }
  return `${label} (${normalized})`
}

const formatPageEntryTag = (entry) => {
  return formatPageTag(entry?.path, entry?.name)
}

const getDisplayPageTags = (rawPages) => {
  const entries = parsePageEntriesFromRaw(rawPages)
  entries.forEach(entry => ensureCustomPageOption(entry.path))
  return entries.map(entry => ({ value: `${entry.path}_${entry.name || ''}`, text: formatPageTag(entry.path, entry.name) }))
}

const extractPagePathList = (data) => {
  if (Array.isArray(data)) return data
  if (Array.isArray(data?.list)) return data.list
  return []
}

const loadRuntimePageOptions = async (clientType = 'all') => {
  try {
    const res = await getPopupPagePathOptions({ clientType })
    if (res.code !== 0) return

    const options = dedupePageEntries(extractPagePathList(res.data).map(path => ({ path })))
      .map(item => item.path)
      .filter(path => path !== 'all')
      .map(path => ({ label: RUNTIME_PAGE_LABEL, value: path }))

    runtimePageOptions.value = options
  } catch (e) {
    runtimePageOptions.value = []
  }
}

const addPageEntry = () => {
  const normalized = normalizePagePathInput(newPagePath.value)
  if (!normalized) return

  const name = String(newPageName.value || '').trim()
  const existedIndex = pageEntries.value.findIndex(item => item.path === normalized)
  if (existedIndex >= 0) {
    pageEntries.value[existedIndex] = {
      ...pageEntries.value[existedIndex],
      name: name || pageEntries.value[existedIndex].name,
    }
    ElMessage.warning('该路径已存在，已更新名称')
  } else {
    pageEntries.value.push({ name, path: normalized })
  }

  ensureCustomPageOption(normalized)
  newPageName.value = ''
  newPagePath.value = ''
}

const removePageEntry = (idx) => {
  pageEntries.value.splice(idx, 1)
  if (editingPageIndex.value === idx) {
    cancelEditPageEntry()
  }
}

const startEditPageEntry = (idx) => {
  const current = pageEntries.value[idx]
  if (!current) return
  editingPageIndex.value = idx
  editingPageName.value = current.name || ''
  editingPagePath.value = current.path || ''
}

const cancelEditPageEntry = () => {
  editingPageIndex.value = -1
  editingPageName.value = ''
  editingPagePath.value = ''
}

const saveEditPageEntry = (idx) => {
  if (editingPageIndex.value !== idx) return

  const normalized = normalizePagePathInput(editingPagePath.value)
  if (!normalized) {
    ElMessage.warning('请填写正确的页面路径')
    return
  }

  const duplicateIndex = pageEntries.value.findIndex((item, itemIndex) => item.path === normalized && itemIndex !== idx)
  if (duplicateIndex >= 0) {
    ElMessage.warning('该路径已存在，请勿重复添加')
    return
  }

  pageEntries.value[idx] = {
    name: String(editingPageName.value || '').trim(),
    path: normalized,
  }
  ensureCustomPageOption(normalized)
  cancelEditPageEntry()
}

// 同步 pageEntries <-> formData.pages
watch(pageEntries, (val) => {
  const normalizedEntries = dedupePageEntries(val)
  formData.value.pages = normalizedEntries.map(item => serializePageEntry(item)).filter(Boolean).join(',')
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
    tableData.value.forEach(row => {
      normalizedPathList(row.pages).forEach(path => ensureCustomPageOption(path))
    })
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
  pageEntries.value = parsePageEntriesFromRaw(row.pages)
  pageEntries.value.forEach(entry => ensureCustomPageOption(entry.path))
  cancelEditPageEntry()
  newPageName.value = ''
  newPagePath.value = ''
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
  pageEntries.value = []
  cancelEditPageEntry()
  newPageName.value = ''
  newPagePath.value = ''
  dialogFormVisible.value = true
  loadRuntimePageOptions(formData.value.clientType || 'all')
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = defaultForm()
  titleI18n.value = {}
  contentI18n.value = {}
  pageEntries.value = []
  cancelEditPageEntry()
  newPageName.value = ''
  newPagePath.value = ''
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
    // pages 从 pageEntries 同步（watch 已处理，这里确保）
    const normalizedEntries = dedupePageEntries(pageEntries.value)
    normalizedEntries.forEach(entry => ensureCustomPageOption(entry.path))
    formData.value.pages = normalizedEntries
      .map(item => serializePageEntry(item))
      .filter(Boolean)
      .join(',')

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

watch(() => formData.value.clientType, (val) => {
  if (!dialogFormVisible.value) return
  loadRuntimePageOptions(val || 'all')
})
</script>
