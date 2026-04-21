<template>
  <div class="p-6">
    <div class="mb-4 flex items-center justify-between">
      <h2 class="text-xl font-semibold">黑名单管理</h2>
      <el-button type="danger" @click="openAdd">手动加入</el-button>
    </div>
    <el-table :data="list" v-loading="loading" border stripe>
      <el-table-column prop="ID" label="ID" width="70" />
      <el-table-column prop="clientUserId" label="用户ID" width="100" />
      <el-table-column prop="reason" label="原因" show-overflow-tooltip />
      <el-table-column prop="createdBy" label="操作人 ID" width="100" />
      <el-table-column prop="CreatedAt" label="时间" :formatter="(r) => fmtTime(r.CreatedAt)" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-popconfirm title="确认移出黑名单？" @confirm="handleRemove(row)">
            <template #reference>
              <el-button text size="small" type="primary">移出</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <div class="flex justify-end mt-4">
      <el-pagination v-model:current-page="query.page" v-model:page-size="query.pageSize" :total="total" layout="total, prev, pager, next" @change="loadList" />
    </div>
    <el-dialog v-model="dialogVisible" title="加入黑名单" width="380px">
      <el-form label-width="80px">
        <el-form-item label="用户ID">
          <el-input-number v-model="form.clientUserID" :min="1" />
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="form.reason" placeholder="可选" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleAdd">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { getBlacklist, addToBlacklist, removeFromBlacklist } from '@/api/customerService'

const list = ref([])
const total = ref(0)
const loading = ref(false)
const query = ref({ page: 1, pageSize: 20 })
const dialogVisible = ref(false)
const form = ref({ clientUserID: null, reason: '' })

async function loadList() {
  loading.value = true
  try {
    const res = await getBlacklist(query.value)
    if (res.code === 0) {
      list.value = res.data?.list || []
      total.value = res.data?.total || 0
    }
  } finally { loading.value = false }
}

function openAdd() {
  form.value = { clientUserID: null, reason: '' }
  dialogVisible.value = true
}
async function handleAdd() {
  await addToBlacklist({ clientUserId: form.value.clientUserID, reason: form.value.reason })
  ElMessage.success('操作成功')
  dialogVisible.value = false
  loadList()
}
async function handleRemove(row) {
  await removeFromBlacklist({ ID: row.ID })
  ElMessage.success('已移出黑名单')
  loadList()
}
const fmtTime = (t) => t ? dayjs(t).format('YYYY-MM-DD HH:mm') : '-'

onMounted(loadList)
</script>
