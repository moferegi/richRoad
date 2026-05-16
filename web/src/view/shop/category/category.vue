<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <el-form-item
          label="创建日期"
          prop="createdAt"
        >
          <template #label>
            <span>
              创建日期
              <el-tooltip content="搜索范围是开始日期（包含）至结束日期（不包含）">
                <el-icon><QuestionFilled /></el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-date-picker
            v-model="searchInfo.startCreatedAt"
            type="datetime"
            placeholder="开始日期"
            :disabled-date="time=> searchInfo.endCreatedAt ? time.getTime() > searchInfo.endCreatedAt.getTime() : false"
          />
          —
          <el-date-picker
            v-model="searchInfo.endCreatedAt"
            type="datetime"
            placeholder="结束日期"
            :disabled-date="time=> searchInfo.startCreatedAt ? time.getTime() < searchInfo.startCreatedAt.getTime() : false"
          />
        </el-form-item>

        <el-form-item
          label="分类标题"
          prop="title"
        >
          <el-input
            v-model="searchInfo.title"
            placeholder="搜索条件"
          />

        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >查询</el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="openDialog()"
        >新增</el-button>
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        :default-sort="{ prop: 'ID', order: 'descending' }"
        max-height="70vh"
        @selection-change="handleSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          align="left"
          label="ID"
          prop="ID"
          width="140"
          sortable
          ></el-table-column>
        <el-table-column
          align="left"
          label="分类标题"
          prop="title"
          width="120"
        />
        <el-table-column label="图标" width="200">
          <template #default="scope">
            <el-image style="width: 100px; height: 100px" :src="getUrl(scope.row.externalIconPath || scope.row.icons)" fit="cover"/>
          </template>
        </el-table-column>
        <el-table-column align="left" label="外部图标路径" prop="externalIconPath" width="160" show-overflow-tooltip />
        <el-table-column align="center" label="uni显示" width="100">
          <template #default="scope">
            <el-switch
              :model-value="scope.row.showInUni !== false"
              @change="(val) => toggleShowInUni(scope.row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="分类描述"
          prop="desc"
          width="120"
        />
        <el-table-column
          align="left"
          label="操作"
          fixed="right"
          min-width="240"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              class="table-button"
              @click="getDetails(scope.row)"
            >
              <el-icon style="margin-right: 5px"><InfoFilled /></el-icon>
              查看详情
            </el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="updateCategoryFunc(scope.row)"
            >变更</el-button>
            <el-button
              type="primary"
              link
              icon="edit"
              class="table-button"
              @click="openDialog(scope.row.ID)"
              v-if="scope.row.parentId === 0"
            >添加子分类</el-button>
            <el-button
              type="primary"
              link
              icon="delete"
              @click="deleteRow(scope.row)"
            >删除</el-button>
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
    <el-drawer
      v-model="dialogFormVisible"
      size="800"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type==='create'?'添加':'修改' }}</span>
          <div>
            <el-button
              type="primary"
              @click="enterDialog"
            >确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="elFormRef"
        :model="formData"
        label-position="top"
        :rules="rule"
        label-width="80px"
      >
        <el-form-item label="图标(上传):"  prop="icons" >
          <SelectImage
              v-model="formData.icons"
              file-type="image"
              :default-folder="CATEGORY_UPLOAD_FOLDER"
              :fixed-upload-folder="true"
          />
        </el-form-item>
        <el-form-item label="外部图标路径(优先于上传图标):" prop="externalIconPath">
          <el-input v-model="formData.externalIconPath" :clearable="true" placeholder="https://example.com/icon.png" />
        </el-form-item>
        <el-form-item label="分类标题(多语言):" prop="title">
          <el-input v-model="formData.title" :clearable="true" placeholder="默认分类标题" class="mb-2" />
          <MultiLangEditor
            :model="titleI18n"
            :languages="enabledLangs"
            title="分类标题多语言"
          />
        </el-form-item>
        <el-form-item label="分类描述(多语言):" prop="desc">
          <el-input v-model="formData.desc" :clearable="true" placeholder="默认分类描述" class="mb-2" />
          <MultiLangEditor
            :model="descI18n"
            :languages="enabledLangs"
            title="分类描述多语言"
          />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-model="detailShow"
      size="800"
      :before-close="closeDetailShow"
      destroy-on-close
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">查看详情</span>
        </div>
      </template>
      <el-descriptions
        :column="1"
        border
      >
        <el-descriptions-item label="分类标题">
          {{ formData.title }}
        </el-descriptions-item>
        <el-descriptions-item label="分类描述">
          {{ formData.desc }}
        </el-descriptions-item>
        <el-descriptions-item label="外部图标路径">
          {{ formData.externalIconPath || '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createCategory,
  deleteCategory,
  deleteCategoryByIds,
  updateCategory,
  findCategory,
  getCategoryList
} from '@/api/shop/category'
import { getEnabledLanguages } from '@/api/client/language'

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict, ReturnArrImg, onDownloadFile } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import SelectImage from "@/components/selectImage/selectImage.vue";
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'
import {getUrl} from "@/utils/image";

defineOptions({
  name: 'Category'
})

const CATEGORY_UPLOAD_FOLDER = 'cloth-on/web-else/fenlei'

// === 多语言支持 ===
const enabledLangs = ref([])
const titleI18n = ref({})
const descI18n = ref({})

const loadLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0) enabledLangs.value = res.data || []
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

onMounted(() => { loadLangs() })

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
  parentID: 0,
  title: '',
  desc: '',
  icons: '',
  externalIconPath: '',
  showInUni: true
})

// 验证规则
const rule = reactive({
})

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
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

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
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
const getTableData = async() => {
  const table = await getCategoryList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// 切换uni显示开关
const toggleShowInUni = async (row, val) => {
  const res = await updateCategory({ ...row, showInUni: val })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getTableData()
  }
}

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async() => {
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
    deleteCategoryFunc(row)
  })
}

// 多选删除
const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
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
    const res = await deleteCategoryByIds({ IDs })
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
const updateCategoryFunc = async(row) => {
  const res = await findCategory({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data.recategory
    titleI18n.value = parseI18nJson(formData.value.title)
    descI18n.value = parseI18nJson(formData.value.desc)
    dialogFormVisible.value = true
  }
}

// 删除行
const deleteCategoryFunc = async(row) => {
  const res = await deleteCategory({ ID: row.ID })
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

// 查看详情控制标记
const detailShow = ref(false)

// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}

// 打开详情
const getDetails = async(row) => {
  // 打开弹窗
  const res = await findCategory({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data.recategory
    openDetailShow()
  }
}

// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  formData.value = {
    parentID: 0,
    title: '',
    desc: '',
    icons: '',
    externalIconPath: '',
    showInUni: true
  }
}

// 打开弹窗
const openDialog = (id) => {
  if (id) {
    formData.value.parentID = id
  }
  titleI18n.value = {}
  descI18n.value = {}
  type.value = 'create'
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  titleI18n.value = {}
  descI18n.value = {}
  formData.value = {
    parentID: 0,
    title: '',
    desc: '',
    icons: '',
    externalIconPath: '',
    showInUni: true
  }
}
// 弹窗确定
const enterDialog = async() => {
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return
    // 序列化多语言字段
    if (enabledLangs.value.length) {
      formData.value.title = serializeI18nJson(titleI18n.value) || formData.value.title
      formData.value.desc = serializeI18nJson(descI18n.value) || formData.value.desc
    }
    let res
    switch (type.value) {
      case 'create':
        res = await createCategory(formData.value)
        break
      case 'update':
        res = await updateCategory(formData.value)
        break
      default:
        res = await createCategory(formData.value)
        break
    }
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

</script>

<style>

</style>
