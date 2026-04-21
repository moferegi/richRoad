<template>
  <div class="p-6">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xl font-semibold">快捷回复</h2>
      <el-button type="primary" @click="openCreate">新增</el-button>
    </div>
    <el-form inline class="mb-4">
      <el-form-item>
        <el-input v-model="query.title" placeholder="搜索标题" clearable @keydown.enter="loadList" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="list" v-loading="loading" border stripe>
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="content" label="内容" show-overflow-tooltip />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column label="操作" width="140">
        <template #default="{ row }">
          <el-button text size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-popconfirm title="确认删除？" @confirm="handleDelete(row)">
            <template #reference>
              <el-button text size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <div class="flex justify-end mt-4">
      <el-pagination v-model:current-page="query.page" v-model:page-size="query.pageSize" :total="total" layout="total, prev, pager, next" @change="loadList" />
    </div>
    <el-dialog v-model="dialogVisible" :title="editForm.ID ? '编辑快捷回复' : '新增快捷回复'" width="480px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="editForm.title" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="editForm.content" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="editForm.sort" :min="0" />
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
import { ElMessage } from 'element-plus'
import { getQuickReplyList, createQuickReply, updateQuickReply, deleteQuickReply } from '@/api/customerService'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const query = ref({ page: 1, pageSize: 20, title: '' })
const dialogVisible = ref(false)
const editForm = ref({ Title: '', Content: '', Sort: 0 })

async function loadList() {
  loading.value = true
  try {
    const res = await getQuickReplyList(query.value)
    if (res.code === 0) {
      list.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } finally { loading.value = false }
}

function openCreate() {
  editForm.value = { title: '', content: '', sort: 0 }
  dialogVisible.value = true
}
function openEdit(row) {
  editForm.value = { ...row }
  dialogVisible.value = true
}
async function handleSave() {
  if (editForm.value.ID) await updateQuickReply(editForm.value)
  else await createQuickReply(editForm.value)
  ElMessage.success('保存成功')
  dialogVisible.value = false
  loadList()
}
async function handleDelete(row) {
  await deleteQuickReply({ ID: row.ID })
  ElMessage.success('删除成功')
  loadList()
}
onMounted(loadList)
</script>
