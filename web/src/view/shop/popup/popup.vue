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

        <el-form-item label="展示位置">
          <el-select v-model="searchInfo.position" placeholder="全部" clearable>
            <el-option label="首页" value="home" />
            <el-option label="商品页" value="goods" />
            <el-option label="全部页面" value="all" />
          </el-select>
        </el-form-item>

        <el-form-item label="客户端类型">
          <el-select v-model="searchInfo.clientType" placeholder="全部" clearable>
            <el-option label="uni端" value="uni" />
            <el-option label="web端" value="web" />
            <el-option label="全部" value="all" />
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
        @selection-change="handleSelectionChange"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="标题" prop="title" width="150" />
        <el-table-column label="图片" width="100">
            <template #default="scope">
              <el-image v-if="scope.row.image" style="width: 50px; height: 50px" :src="getUrl(scope.row.image)" fit="cover" :preview-src-list="[getUrl(scope.row.image)]" preview-teleported />
              <span v-else>-</span>
            </template>
        </el-table-column>
        <el-table-column align="left" label="跳转链接" prop="link" width="150" show-overflow-tooltip />
        <el-table-column align="left" label="位置" width="80">
          <template #default="scope">
            <el-tag size="small">{{ positionMap[scope.row.position] || scope.row.position }}</el-tag>
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
            <el-form-item label="标题(JSON多语言):" prop="title">
              <el-input v-model="formData.title" placeholder='例: {"zh":"标题","en":"Title","mn":"Гарчиг"}' />
            </el-form-item>
            <el-form-item label="弹窗图片:" prop="image">
              <SelectImage v-model="formData.image" file-type="image" />
            </el-form-item>
            <el-form-item label="跳转链接:" prop="link">
              <el-input v-model="formData.link" placeholder="点击弹窗跳转的链接" />
            </el-form-item>
            <el-form-item label="生效时间:" prop="startTime">
              <el-date-picker v-model="formData.startTime" type="datetime" placeholder="选择生效时间（留空不限）" />
            </el-form-item>
            <el-form-item label="失效时间:" prop="endTime">
              <el-date-picker v-model="formData.endTime" type="datetime" placeholder="选择失效时间（留空不限）" />
            </el-form-item>
            <el-form-item label="客户端类型:" prop="clientType">
              <el-select v-model="formData.clientType">
                <el-option label="全部" value="all" />
                <el-option label="uni端" value="uni" />
                <el-option label="web端" value="web" />
              </el-select>
            </el-form-item>
            <el-form-item label="展示位置:" prop="position">
              <el-select v-model="formData.position">
                <el-option label="首页" value="home" />
                <el-option label="商品页" value="goods" />
                <el-option label="全部页面" value="all" />
              </el-select>
            </el-form-item>
            <el-form-item label="只弹一次:" prop="onceOnly">
              <el-switch v-model="formData.onceOnly" />
            </el-form-item>
            <el-form-item label="可关闭:" prop="closeable">
              <el-switch v-model="formData.closeable" />
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
  createPopup,
  deletePopup,
  updatePopup,
  getPopupList
} from '@/api/shop/popup'
import { getUrl } from '@/utils/image'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({ name: 'Popup' })

const btnLoading = ref(false)

const positionMap = { home: '首页', goods: '商品页', all: '全部' }
const clientTypeMap = { uni: 'uni端', web: 'web端', all: '全部' }

const defaultForm = () => ({
  title: '',
  image: '',
  link: '',
  startTime: null,
  endTime: null,
  clientType: 'all',
  onceOnly: false,
  sort: 0,
  closeable: true,
  position: 'home',
  isEnabled: true,
})

const formData = ref(defaultForm())

const rule = reactive({
  title: [{ required: true, message: '请输入标题', trigger: ['input','blur'] }],
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

const onReset = () => { searchInfo.value = {}; getTableData() }
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
const updateFunc = async(row) => { type.value = 'update'; formData.value = { ...row }; dialogFormVisible.value = true }
const deleteFunc = async (row) => {
  const res = await deletePopup({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '删除成功' })
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const dialogFormVisible = ref(false)
const openDialog = () => { type.value = 'create'; dialogFormVisible.value = true }
const closeDialog = () => { dialogFormVisible.value = false; formData.value = defaultForm() }

const enterDialog = async () => {
  btnLoading.value = true
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return btnLoading.value = false
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
