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
          <el-tag :type="row.onlineStatus === 'online' ? 'success' : 'info'" size="small">
            {{ row.onlineStatus === 'online' ? '在线' : '离线' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="当前会话" width="180">
        <template #default="{ row }">
          <div class="flex items-center gap-2">
            <span class="text-xs text-gray-600">{{ row.activeSessions || 0 }}/{{ row.maxSessions || 0 }}</span>
            <el-progress
              :percentage="Math.min(100, Math.round(((row.activeSessions || 0) / Math.max(1, row.maxSessions || 1)) * 100))"
              :stroke-width="8"
              :show-text="false"
              style="width: 90px"
              :status="(row.activeSessions || 0) >= (row.maxSessions || 1) ? 'exception' : ''"
            />
          </div>
        </template>
      </el-table-column>
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
        <el-form-item v-if="!editForm.ID" label="选择用户">
          <el-select
            v-model="editForm.userID"
            filterable
            remote
            reserve-keyword
            placeholder="搜索用户昵称/ID"
            :remote-method="searchUsers"
            :loading="userSearchLoading"
            style="width:100%"
          >
            <el-option
              v-for="u in userOptions"
              :key="u.ID"
              :label="`${u.nickName} (ID:${u.ID})`"
              :value="u.ID"
            />
          </el-select>
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
import { getUserList } from '@/api/user'

const list = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editForm = ref({ userID: null, nickname: '', maxSessions: 5, isEnabled: true })

// 用户选择器
const userOptions = ref([])
const userSearchLoading = ref(false)

async function searchUsers(query) {
  userSearchLoading.value = true
  try {
    const res = await getUserList({ page: 1, pageSize: 20, nickName: query || '' })
    userOptions.value = res?.data?.list || []
  } finally {
    userSearchLoading.value = false
  }
}

async function loadList() {
  loading.value = true
  try {
    const res = await getAgentList()
    if (res.code === 0) list.value = Array.isArray(res.data) ? res.data : (res.data?.list || [])
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editForm.value = { userID: null, nickname: '', maxSessions: 5, isEnabled: true }
  userOptions.value = []
  // 默认加载一批用户
  searchUsers('')
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
      res = await createAgent({ userId: editForm.value.userID, nickname: editForm.value.nickname, maxSessions: editForm.value.maxSessions || 5, isEnabled: editForm.value.isEnabled })
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
  const res = await updateAgent({ ID: row.ID, isEnabled: val, maxSessions: row.maxSessions || 5 })
  if (res?.code === 0) {
    row.isEnabled = val
  }
  // 失败时 request 拦截器已弹出错误提示，UI 保持不变（el-switch 绑定的是 :model-value 不会自动回滚）
}

onMounted(loadList)
</script>
