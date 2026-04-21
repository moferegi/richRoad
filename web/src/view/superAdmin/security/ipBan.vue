<template>
  <div class="gva-table-box">
    <!-- 顶部操作栏 -->
    <div class="gva-btn-list mb-4 flex items-center gap-2">
      <el-button type="primary" icon="Plus" @click="openBanDialog()">封禁IP</el-button>
      <el-input
        v-model="searchForm.ip"
        placeholder="搜索IP地址"
        clearable
        style="width: 220px"
        @keyup.enter="fetchList"
        @clear="fetchList"
      />
      <el-select v-model="searchForm.isAuto" placeholder="封禁类型" clearable style="width: 130px" @change="fetchList">
        <el-option label="手动封禁" :value="false" />
        <el-option label="自动封禁" :value="true" />
      </el-select>
      <el-button icon="Search" @click="fetchList">查询</el-button>

      <div class="ml-auto flex items-center gap-2">
        <el-select v-model="statsHours" style="width: 110px" @change="fetchStats">
          <el-option label="近1小时" :value="1" />
          <el-option label="近6小时" :value="6" />
          <el-option label="近24小时" :value="24" />
          <el-option label="近7天" :value="168" />
        </el-select>
        <el-button icon="DataLine" @click="fetchStats">刷新统计</el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <!-- 封禁列表 Tab -->
      <el-tab-pane label="封禁列表" name="banned">
        <el-table
          v-loading="listLoading"
          :data="bannedList"
          border
          stripe
          row-key="ID"
          style="width: 100%"
        >
          <el-table-column prop="ip" label="IP地址" width="160" />
          <el-table-column prop="reason" label="封禁原因" min-width="180" show-overflow-tooltip />
          <el-table-column prop="bannedBy" label="操作人" width="120" />
          <el-table-column label="封禁类型" width="100" align="center">
            <template #default="{ row }">
              <el-tag :type="row.isAuto ? 'warning' : 'danger'" size="small">
                {{ row.isAuto ? '自动' : '手动' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="过期时间" width="180">
            <template #default="{ row }">
              <span v-if="!row.expiredAt" class="text-red-500 font-semibold">永久</span>
              <span v-else>{{ formatTime(row.expiredAt) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="封禁时间" width="180">
            <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right" align="center">
            <template #default="{ row }">
              <el-button type="danger" link size="small" @click="handleUnban(row.ip)">解封</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="gva-pagination mt-4">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :total="total"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next"
            @size-change="fetchList"
            @current-change="fetchList"
          />
        </div>
      </el-tab-pane>

      <!-- 攻击统计 Tab -->
      <el-tab-pane label="攻击统计" name="stats">
        <el-table
          v-loading="statsLoading"
          :data="statsList"
          border
          stripe
          style="width: 100%"
        >
          <el-table-column type="index" label="#" width="60" align="center" />
          <el-table-column prop="ip" label="IP地址" width="155" />
          <el-table-column prop="totalCount" label="总次数" width="90" align="center" sortable>
            <template #default="{ row }">
              <el-tag :type="row.totalCount >= 20 ? 'danger' : row.totalCount >= 10 ? 'warning' : 'info'" size="small">
                {{ row.totalCount }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="loginFail" label="密码错误" width="95" align="center">
            <template #default="{ row }">
              <span v-if="row.loginFail" class="text-orange-500 font-semibold">{{ row.loginFail }}</span>
              <span v-else class="text-gray-300">—</span>
            </template>
          </el-table-column>
          <el-table-column prop="captchaFail" label="验证码错误" width="105" align="center">
            <template #default="{ row }">
              <span v-if="row.captchaFail" class="text-yellow-500 font-semibold">{{ row.captchaFail }}</span>
              <span v-else class="text-gray-300">—</span>
            </template>
          </el-table-column>
          <el-table-column prop="registerLimit" label="注册超限" width="95" align="center">
            <template #default="{ row }">
              <span v-if="row.registerLimit" class="text-purple-500 font-semibold">{{ row.registerLimit }}</span>
              <span v-else class="text-gray-300">—</span>
            </template>
          </el-table-column>
          <el-table-column prop="adminFail" label="后台登录失败" width="120" align="center">
            <template #default="{ row }">
              <span v-if="row.adminFail" class="text-red-500 font-semibold">{{ row.adminFail }}</span>
              <span v-else class="text-gray-300">—</span>
            </template>
          </el-table-column>
          <el-table-column label="最后攻击时间" width="175">
            <template #default="{ row }">{{ formatTime(row.lastTime) }}</template>
          </el-table-column>
          <el-table-column label="封禁状态" width="95" align="center">
            <template #default="{ row }">
              <el-tag :type="row.isBanned ? 'danger' : 'success'" size="small">
                {{ row.isBanned ? '已封禁' : '正常' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right" align="center">
            <template #default="{ row }">
              <el-button v-if="!row.isBanned" type="danger" link size="small" @click="openBanDialog(row.ip)">
                封禁
              </el-button>
              <el-button v-else type="success" link size="small" @click="handleUnban(row.ip)">
                解封
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- 封禁IP对话框 -->
    <el-dialog v-model="banDialogVisible" title="封禁IP" width="480px" @close="resetBanForm">
      <el-form ref="banFormRef" :model="banForm" :rules="banRules" label-width="100px">
        <el-form-item label="IP地址" prop="ip">
          <el-input v-model="banForm.ip" placeholder="输入IP地址，如 192.168.1.1" />
        </el-form-item>
        <el-form-item label="封禁原因" prop="reason">
          <el-input v-model="banForm.reason" placeholder="请输入封禁原因" />
        </el-form-item>
        <el-form-item label="封禁时长">
          <el-select v-model="banForm.duration" style="width: 100%">
            <el-option label="永久封禁" :value="0" />
            <el-option label="1小时" :value="60" />
            <el-option label="6小时" :value="360" />
            <el-option label="24小时" :value="1440" />
            <el-option label="7天" :value="10080" />
            <el-option label="30天" :value="43200" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="banLoading" @click="submitBan">确认封禁</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBannedIPList, getAttackStats, banIP, unbanIP } from '@/api/security'

const activeTab = ref('banned')

// ===== 封禁列表 =====
const listLoading = ref(false)
const bannedList = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchForm = ref({ ip: '', isAuto: null })

const fetchList = async () => {
  listLoading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ip: searchForm.value.ip || undefined,
      isAuto: searchForm.value.isAuto !== null ? searchForm.value.isAuto : undefined
    }
    const res = await getBannedIPList(params)
    if (res.code === 0) {
      bannedList.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } finally {
    listLoading.value = false
  }
}

// ===== 攻击统计 =====
const statsLoading = ref(false)
const statsList = ref([])
const statsHours = ref(24)

const fetchStats = async () => {
  statsLoading.value = true
  try {
    const res = await getAttackStats({ hours: statsHours.value })
    if (res.code === 0) {
      statsList.value = res.data || []
    }
  } finally {
    statsLoading.value = false
  }
}

// ===== 封禁操作 =====
const banDialogVisible = ref(false)
const banLoading = ref(false)
const banFormRef = ref()
const banForm = ref({ ip: '', reason: '手动封禁', duration: 0 })
const banRules = {
  ip: [{ required: true, message: '请输入IP地址', trigger: 'blur' }],
  reason: [{ required: true, message: '请输入封禁原因', trigger: 'blur' }]
}

const openBanDialog = (ip = '') => {
  banForm.value = { ip, reason: '手动封禁', duration: 0 }
  banDialogVisible.value = true
}

const resetBanForm = () => {
  banFormRef.value?.resetFields()
}

const submitBan = async () => {
  await banFormRef.value?.validate()
  banLoading.value = true
  try {
    const res = await banIP(banForm.value)
    if (res.code === 0) {
      ElMessage.success('封禁成功')
      banDialogVisible.value = false
      fetchList()
      fetchStats()
    } else {
      ElMessage.error(res.msg || '封禁失败')
    }
  } finally {
    banLoading.value = false
  }
}

const handleUnban = (ip) => {
  ElMessageBox.confirm(`确认解封 IP：${ip}？`, '解封确认', {
    confirmButtonText: '确认解封',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await unbanIP({ ip })
    if (res.code === 0) {
      ElMessage.success('解封成功')
      fetchList()
      fetchStats()
    } else {
      ElMessage.error(res.msg || '解封失败')
    }
  })
}

// ===== 工具函数 =====
const formatTime = (t) => {
  if (!t) return '-'
  return new Date(t).toLocaleString('zh-CN', { hour12: false })
}

onMounted(() => {
  fetchList()
  fetchStats()
})
</script>
