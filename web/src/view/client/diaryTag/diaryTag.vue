<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
      <el-form-item label="创建日期" prop="CreatedAt">
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
      
        <template v-if="showAllQuery">
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery=true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery=false" v-else>收起</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button  type="primary" icon="plus" @click="openDialog()">新增</el-button>
            <el-button  icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
        @sort-change="sortChange"
        >
        <el-table-column type="selection" width="55" />
        
        <el-table-column sortable align="left" label="日期" prop="CreatedAt" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        
            <el-table-column sortable align="left" label="标签名" prop="name" width="120" />

            <el-table-column align="left" label="颜色" prop="color" width="120">
              <template #default="scope">
                <div v-if="scope.row.color" style="display:flex;align-items:center;gap:8px;">
                  <span :style="{ display:'inline-block', width:'20px', height:'20px', borderRadius:'4px', backgroundColor: scope.row.color, border: '1px solid #dcdfe6' }" />
                  <span>{{ scope.row.color }}</span>
                </div>
                <span v-else>-</span>
              </template>
            </el-table-column>

        <el-table-column align="left" label="操作" fixed="right" :min-width="appStore.operateMinWith">
            <template #default="scope">
            <el-button  type="primary" link class="table-button" @click="getDetails(scope.row)"><el-icon style="margin-right: 5px"><InfoFilled /></el-icon>查看</el-button>
            <el-button  type="primary" link icon="edit" class="table-button" @click="updateDiaryTagFunc(scope.row)">编辑</el-button>
            <el-button   type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
       <template #header>
              <div class="flex justify-between items-center">
                <span class="text-lg">{{type==='create'?'新增日记标签':'编辑日记标签'}}</span>
                <div>
                  <el-button :loading="btnLoading" type="primary" @click="enterDialog">确 定</el-button>
                  <el-button @click="closeDialog">取 消</el-button>
                </div>
              </div>
            </template>

          <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule" label-width="80px">
            <el-form-item label="标签名:" prop="name">
    <el-input v-model="formData.name" :clearable="false" placeholder="请输入标签名" />
</el-form-item>
            <el-form-item label="标签名(多语言):">
              <MultiLangEditor
                :model="nameI18nObj"
                :languages="enabledLangs"
                title="标签名多语言"
              />
            </el-form-item>
            <el-form-item label="描述:" prop="description">
    <el-input v-model="formData.description" :clearable="false" placeholder="请输入描述" />
</el-form-item>
            <el-form-item label="颜色:" prop="color">
            <el-color-picker v-model="formData.color" />
</el-form-item>
          </el-form>
    </el-drawer>

    <el-drawer destroy-on-close :size="appStore.drawerSize" v-model="detailShow" :show-close="true" :before-close="closeDetailShow" title="查看日记标签">
            <el-descriptions :column="1" border>
                    <el-descriptions-item label="标签名">
    {{ detailFrom.name }}
</el-descriptions-item>
                    <el-descriptions-item label="描述">
    {{ detailFrom.description }}
</el-descriptions-item>
            </el-descriptions>
        </el-drawer>

  </div>
</template>

<script setup>
import {
  createDiaryTag,
  deleteDiaryTag,
  deleteDiaryTagByIds,
  updateDiaryTag,
  findDiaryTag,
  getDiaryTagList
} from '@/api/client/diaryTag'
import { getEnabledLanguages } from '@/api/client/language'

import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import { useAppStore } from "@/pinia"
import MultiLangEditor from '@/components/multilingual/multi-lang-editor.vue'

defineOptions({
    name: 'DiaryTag'
})

const btnLoading = ref(false)
const appStore = useAppStore()

const showAllQuery = ref(false)

const formData = ref({
            name: '',
            nameI18n: '',
            description: '',
            color: '',
        })

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
  Object.entries(nameI18nObj).forEach(([rawCode, rawText]) => {
    const code = String(rawCode || '').trim()
    if (!code) return
    const text = String(rawText ?? '').trim()
    if (!text) return
    obj[code] = text
  })
  return JSON.stringify(obj)
}

onMounted(() => {
  loadEnabledLangs()
})

const rule = reactive({
               name : [{
                   required: true,
                   message: '请输入标签名',
                   trigger: ['input','blur'],
               },
               {
                   whitespace: true,
                   message: '不能只输入空格',
                   trigger: ['input', 'blur'],
              }
              ],
})

const searchRule = reactive({
  CreatedAt: [
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

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const sortChange = ({ prop, order }) => {
  const sortMap = {
    CreatedAt:"created_at",
    ID:"id",
            name: 'name',
  }

  let sort = sortMap[prop]
  if(!sort){
   sort = prop.replace(/[A-Z]/g, match => `_${match.toLowerCase()}`)
  }

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

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async() => {
  const table = await getDiaryTagList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

const setOptions = async () =>{}
setOptions()

const multipleSelection = ref([])
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteDiaryTagFunc(row)
        })
    }

const onDelete = async() => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
      const ids = multipleSelection.value.map(item => item.ID)
      const res = await deleteDiaryTagByIds({ ids: ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        getTableData()
      }
  })
}

const deleteDiaryTagFunc = async(row) => {
    const res = await deleteDiaryTag({ id: row.ID })
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

const type = ref('')
const dialogFormVisible = ref(false)

const updateDiaryTagFunc = async(row) => {
    const res = await findDiaryTag({ id: row.ID })
    if (res.code !== 0) return
    type.value = 'update'
    formData.value = { ...res.data }
    parseNameI18n(res.data?.nameI18n || '')
    dialogFormVisible.value = true
}

const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        name: '',
        nameI18n: '',
        description: '',
        color: '',
    }
    Object.keys(nameI18nObj).forEach(k => delete nameI18nObj[k])
}

const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

const enterDialog = async() => {
    btnLoading.value = true
    elFormRef.value?.validate(async(valid) => {
        if (!valid) {
            btnLoading.value = false
            return
        }
        formData.value.nameI18n = serializeNameI18n()
        let res
        if (type.value === 'create') {
            res = await createDiaryTag(formData.value)
        } else {
            res = await updateDiaryTag(formData.value)
        }
        if (res.code === 0) {
            ElMessage({
                type: 'success',
                message: type.value === 'create' ? '创建成功' : '更新成功'
            })
            closeDialog()
            getTableData()
        }
        btnLoading.value = false
    })
}

const detailFrom = ref({})
const detailShow = ref(false)

const getDetails = async(row) => {
    const res = await findDiaryTag({ id: row.ID })
    if (res.code === 0) {
        detailFrom.value = res.data
        detailShow.value = true
    }
}

const closeDetailShow = () => {
    detailShow.value = false
    detailFrom.value = {}
}
</script>

<style scoped>
</style>