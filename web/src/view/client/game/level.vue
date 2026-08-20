<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="游戏">
          <el-select v-model="searchInfo.gameID" placeholder="选择游戏" clearable @change="onGameChange">
            <el-option v-for="item in gameList" :key="item.ID" :label="localText(item.name)" :value="item.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="难度分类">
          <el-select v-model="searchInfo.categoryID" placeholder="选择难度分类" clearable @change="onSearch">
            <el-option v-for="item in diffList" :key="item.ID" :label="localText(item.name)" :value="item.ID" />
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
        <el-table-column align="left" label="数字" prop="numbers" width="200" />
        <el-table-column align="left" label="目标结果" prop="targetResult" width="100" />
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="关卡序号" prop="levelNumber">
          <el-input-number v-model="form.levelNumber" :min="1" />
        </el-form-item>
        <el-form-item label="4个数字" prop="numbers">
          <el-input v-model="form.numbers" placeholder="逗号分隔，如: 3,8,3,8" />
        </el-form-item>
        <el-form-item label="目标结果">
          <el-input-number v-model="form.targetResult" :min="1" :max="999" />
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
import { getLevelList, createLevel, updateLevel, deleteLevel } from '@/api/client/game'

defineOptions({ name: 'GameLevel' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const gameList = ref([])
const diffList = ref([])
const searchInfo = reactive({ gameID: '', categoryID: '' })

const dialogVisible = ref(false)
const dialogTitle = ref('新增关卡')
const formRef = ref(null)
const form = reactive({
  ID: 0,
  levelNumber: 1,
  numbers: '',
  targetResult: 24,
  sort: 0
})
const rules = {
  levelNumber: [{ required: true, message: '请输入关卡序号', trigger: 'blur' }],
  numbers: [{ required: true, message: '请输入4个数字', trigger: 'blur' }]
}

const localText = (value) => {
  if (!value) return ''
  if (typeof value === 'object') return value.zh || value.en || ''
  try {
    const obj = JSON.parse(value)
    return obj.zh || obj.en || ''
  } catch {
    return value
  }
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
  const res = await getLevelList(params)
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
    form.numbers = row.numbers
    form.targetResult = row.targetResult
    form.sort = row.sort
  } else {
    dialogTitle.value = '新增关卡'
    form.ID = 0
    form.levelNumber = 1
    form.numbers = ''
    form.targetResult = 24
    form.sort = 0
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return
  const data = {
    categoryID: Number(searchInfo.categoryID),
    levelNumber: form.levelNumber,
    numbers: form.numbers,
    targetResult: form.targetResult,
    sort: form.sort
  }
  if (form.ID) {
    data.ID = form.ID
    const res = await updateLevel(data)
    if (res.code === 0) { ElMessage.success('更新成功'); dialogVisible.value = false; getTableData() }
  } else {
    const res = await createLevel(data)
    if (res.code === 0) { ElMessage.success('创建成功'); dialogVisible.value = false; getTableData() }
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定要删除该关卡吗？', '提示', { type: 'warning' })
  const res = await deleteLevel({ ID: row.ID })
  if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
}

loadGameList()
getTableData()
</script>