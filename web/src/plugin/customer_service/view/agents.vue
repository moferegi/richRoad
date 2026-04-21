<template>
  <div class="p-6">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xl font-semibold">坐席管理</h2>
      <el-button type="primary" @click="openCreate">新增坐席</el-button>
    </div>

    <el-table :data="list" v-loading="loading" border stripe>
      <el-table-column prop="userId" label="用户ID" width="90" />
      <el-table-column label="昵称">
        <template #default="{ row }">{{ row.nickname || row.sysNickname }}</template>
      </el-table-column>
      <el-table-column label="在线状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.onlineStatus ? 'success' : 'info'" size="small">
            {{ row.onlineStatus ? '在线' : '离线' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="activeSessions" label="当前会话" width="100" />
      <el-table-column prop="maxSessions" label="最大会话" width="100" />
      <el-table-column label="启用状态" width="100">
        <template #default="{ row }">
          <el-switch :model-value="row.isEnabled" @change="(v) => toggleEnabled(row, v)" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{ row }">
          <el-button text size="small" type="primary" @click="openEdit(row)">编辑</el-button>
          <el-popconfirm title="确认删除该坐席？" @confirm="handleDelete(row)">
            <template #reference>
              <el-button text size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="editForm.ID ? '编辑坐席' : '新增坐席'" width="400px">
      <el-form :model="editForm" label-width="90px">
        <el-form-item v-if="!editForm.ID" label="用户ID">
          <el-input-number v-model="editForm.userID" :min="1" />
        </el-form-item>
        <el-form-item label="显示昵称">
          <el-input v-model="editForm.nickname" placeholder="可选，留空使用系统昵称" />
        </el-form-item>
        <el-form-item label="最大会话">
          <el-input-number v-model="editForm.maxSessions" :min="1" :max="20" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="editForm.isEnabled" />
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
import { getAgentList, createAgent, updateAgent, deleteAgent } from '@/api/customerService'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editForm = ref({ userID: null, nickname: '', maxSessions: 5, isEnabled: true })

async function loadList() {
  loading.value = true
  try {
    const res = await getAgentList()
    if (res.code === 0) list.value = res.data?.list || []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editForm.value = { userID: null, nickname: '', maxSessions: 5, isEnabled: true }
  dialogVisible.value = true
}

function openEdit(row) {
  editForm.value = { ID: row.ID, nickname: row.nickname || '', maxSessions: row.maxSessions, isEnabled: row.isEnabled }
  dialogVisible.value = true
}

async function handleSave() {
  try {
    let res
    if (editForm.value.ID) {
      res = await updateAgent(editForm.value)
    } else {
      res = await createAgent({ userId: editForm.value.userID, nickname: editForm.value.nickname, maxSessions: editForm.value.maxSessions, isEnabled: editForm.value.isEnabled })
    }
    if (res?.code === 0) {
      ElMessage.success('保存成功')
      dialogVisible.value = false
      loadList()
    }
  } catch (e) {
    console.error('保存坐席失败', e)
  }
}

async function handleDelete(row) {
  await deleteAgent({ ID: row.ID })
  ElMessage.success('删除成功')
  loadList()
}

async function toggleEnabled(row, val) {
  await updateAgent({ ID: row.ID, isEnabled: val, maxSessions: row.maxSessions })
  row.isEnabled = val
}

onMounted(loadList)
</script>
