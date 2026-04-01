<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="名称">
          <el-input v-model="searchInfo.name" placeholder="域名名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList">查询</el-button>
          <el-button type="primary" @click="openCreate">新增域名</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" stripe>
        <el-table-column prop="name" label="名称" width="150" />
        <el-table-column prop="domain" label="域名地址" min-width="300">
          <template #default="scope">
            <el-link type="primary" :href="scope.row.domain" target="_blank">{{ scope.row.domain }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="默认" width="80" align="center">
          <template #default="scope">
            <el-tag v-if="scope.row.isDefault" type="success" size="small">默认</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用" width="80" align="center">
          <template #default="scope">
            <el-switch
              :model-value="scope.row.isEnabled"
              @change="(val) => handleToggle(scope.row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column prop="remark" label="备注" width="200" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openEdit(scope.row)">编辑</el-button>
            <el-button v-if="!scope.row.isDefault" type="warning" link @click="handleSetDefault(scope.row)">设为默认</el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="searchInfo.page"
          :page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 编辑/创建弹窗 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑域名' : '新增域名'" width="500px">
      <el-form :model="form" label-width="100px" :rules="rules" ref="formRef">
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如：阿里云CDN" />
        </el-form-item>
        <el-form-item label="域名地址" prop="domain">
          <el-input v-model="form.domain" placeholder="如：https://cdn.example.com" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="设为默认">
          <el-switch v-model="form.isDefault" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.isEnabled" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  getExternalLinkDomainList,
  createExternalLinkDomain,
  updateExternalLinkDomain,
  deleteExternalLinkDomain,
  setDefaultDomain
} from '@/api/client/externalLinkDomain'
import { ElMessage, ElMessageBox } from 'element-plus'

const searchInfo = ref({
  page: 1,
  pageSize: 10,
  name: '',
})

const tableData = ref([])
const total = ref(0)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const defaultForm = {
  name: '',
  domain: '',
  isDefault: false,
  isEnabled: true,
  sort: 0,
  remark: '',
}

const form = ref({ ...defaultForm })

const rules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  domain: [{ required: true, message: '请输入域名地址', trigger: 'blur' }],
}

const getList = async () => {
  const res = await getExternalLinkDomainList(searchInfo.value)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
  }
}

const handleCurrentChange = (val) => {
  searchInfo.value.page = val
  getList()
}

const handleSizeChange = (val) => {
  searchInfo.value.pageSize = val
  searchInfo.value.page = 1
  getList()
}

const openCreate = () => {
  isEdit.value = false
  form.value = { ...defaultForm }
  dialogVisible.value = true
}

const openEdit = (row) => {
  isEdit.value = true
  form.value = { ...row }
  dialogVisible.value = true
}

const handleSave = async () => {
  const fn = isEdit.value ? updateExternalLinkDomain : createExternalLinkDomain
  const res = await fn(form.value)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    dialogVisible.value = false
    getList()
  }
}

const handleDelete = async (row) => {
  await ElMessageBox.confirm('确定删除该域名？', '提示', { type: 'warning' })
  const res = await deleteExternalLinkDomain({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    getList()
  }
}

const handleSetDefault = async (row) => {
  const res = await setDefaultDomain({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getList()
  }
}

const handleToggle = async (row, val) => {
  const res = await updateExternalLinkDomain({
    ...row,
    isEnabled: val,
  })
  if (res.code === 0) {
    ElMessage.success('更新成功')
    getList()
  }
}

onMounted(() => {
  getList()
})
</script>
