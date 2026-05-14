<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="表名/用途">
          <el-input v-model="searchInfo.keyword" clearable placeholder="输入表名或用途关键词" />
        </el-form-item>
        <el-form-item>
          <el-button :loading="loading" icon="search" type="primary" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-row :gutter="16" class="status-row">
        <el-col :lg="12" :md="24" :sm="24" :xs="24">
          <el-card shadow="never">
            <template #header>
              <div class="status-header">
                <span>数据库运行状态</span>
                <el-tag :type="databaseStatus.healthy ? 'success' : 'danger'">
                  {{ databaseStatus.healthy ? '正常' : '异常' }}
                </el-tag>
              </div>
            </template>
            <div class="status-content">
              <div class="status-row-item">数据库：{{ databaseStatus.databaseName || '-' }}</div>
              <div class="status-row-item">类型：{{ databaseStatus.dbType || '-' }}</div>
              <div class="status-row-item">总大小：{{ databaseStatus.totalSizeHuman || '-' }}</div>
              <div class="status-row-item">总记录数：{{ databaseStatus.totalRows || 0 }}</div>
              <div class="status-row-item">原因：{{ databaseStatus.reason || '-' }}</div>
              <div class="status-row-item">连接数：{{ databaseStatus.openConnections || 0 }}/{{ databaseStatus.maxOpenConnections || 0 }}</div>
              <div v-if="databaseStatus.warning" class="status-warning">{{ databaseStatus.warning }}</div>
              <div class="status-guides" v-if="(databaseStatus.manualFixGuide || []).length">
                <div class="guide-title">手动修复教程：</div>
                <div
                  v-for="(step, idx) in databaseStatus.manualFixGuide"
                  :key="`db-step-${idx}`"
                  class="guide-item"
                >
                  {{ idx + 1 }}. {{ step }}
                </div>
              </div>
              <div class="status-action">
                <el-button :loading="autoFixingTarget === 'database'" type="warning" @click="handleAutoFix('database')">
                  自动修复数据库连接
                </el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :lg="12" :md="24" :sm="24" :xs="24">
          <el-card shadow="never">
            <template #header>
              <div class="status-header">
                <span>Redis运行状态</span>
                <el-tag :type="redisStatus.healthy ? 'success' : 'danger'">
                  {{ redisStatus.healthy ? '正常' : '异常' }}
                </el-tag>
              </div>
            </template>
            <div class="status-content">
              <div class="status-row-item">是否启用：{{ redisStatus.enabled ? '是' : '否' }}</div>
              <div class="status-row-item">原因：{{ redisStatus.reason || '-' }}</div>
              <div class="status-row-item">版本：{{ redisStatus.version || '-' }}</div>
              <div class="status-row-item">内存：{{ redisStatus.usedMemory || '-' }}</div>
              <div class="status-guides" v-if="(redisStatus.manualFixGuide || []).length">
                <div class="guide-title">手动修复教程：</div>
                <div
                  v-for="(step, idx) in redisStatus.manualFixGuide"
                  :key="`redis-step-${idx}`"
                  class="guide-item"
                >
                  {{ idx + 1 }}. {{ step }}
                </div>
              </div>
              <div class="status-action" v-if="redisStatus.enabled">
                <el-button :loading="autoFixingTarget === 'redis'" type="warning" @click="handleAutoFix('redis')">
                  自动修复Redis连接
                </el-button>
              </div>
              <div class="keyspace-box" v-if="redisKeyspaces.length">
                <div class="guide-title">Redis Keyspace：</div>
                <el-table :data="redisKeyspaces" size="small" style="width: 100%">
                  <el-table-column label="DB" prop="db" width="80" />
                  <el-table-column label="Key数" prop="keys" width="100" />
                  <el-table-column label="过期Key" prop="expires" width="100" />
                  <el-table-column label="平均TTL(ms)" prop="avgTTL" min-width="120" />
                </el-table>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-alert
        type="warning"
        show-icon
        :closable="false"
        title="危险操作提醒：按日期范围真删除会直接物理删除数据库记录；若开启联动删除，会同时删除对应本地/OSS文件。"
      />

      <el-table
        :data="tableData"
        :loading="loading"
        row-key="tableName"
        style="margin-top: 12px; width: 100%"
      >
        <el-table-column align="left" label="表名" min-width="200" prop="tableName" />
        <el-table-column align="left" label="总大小" min-width="120" prop="sizeHuman" />
        <el-table-column align="left" label="记录数" min-width="100" prop="rowCount" />
        <el-table-column align="left" label="作用" min-width="180" prop="purpose" show-overflow-tooltip />
        <el-table-column align="left" label="系统判定" min-width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.systemRequired" type="danger">系统必须</el-tag>
            <el-tag v-else-if="scope.row.likelyUnused" type="warning">疑似无用</el-tag>
            <el-tag v-else type="info">业务数据</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="可清理时间字段" min-width="160">
          <template #default="scope">
            <span v-if="(scope.row.dateColumns || []).length">{{ scope.row.dateColumns.join(', ') }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="文件关联字段" min-width="180" show-overflow-tooltip>
          <template #default="scope">
            <span v-if="(scope.row.fileColumns || []).length">{{ scope.row.fileColumns.join(', ') }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="异常原因/保护说明" min-width="220" show-overflow-tooltip>
          <template #default="scope">
            {{ scope.row.deleteGuardReason || '未发现异常' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="160">
          <template #default="scope">
            <el-button
              :disabled="!scope.row.canDeleteByDate"
              link
              type="danger"
              @click="openDeleteDialog(scope.row)"
            >
              按日期真删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="deleteDialogVisible" destroy-on-close title="按日期范围真删除" width="640px">
      <el-form label-width="130px">
        <el-form-item label="目标表">
          <el-input v-model="deleteForm.tableName" disabled />
        </el-form-item>
        <el-form-item label="时间字段">
          <el-select v-model="deleteForm.timeColumn" class="w-full" placeholder="请选择时间字段">
            <el-option
              v-for="item in deleteForm.dateColumns"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="deleteForm.dateRange"
            class="w-full"
            end-placeholder="结束时间"
            start-placeholder="开始时间"
            type="datetimerange"
            unlink-panels
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="联动删文件">
          <el-switch v-model="deleteForm.deleteRelatedFiles" />
          <span class="inline-tip">开启后会删除匹配记录内解析出的本地/OSS文件。</span>
        </el-form-item>
        <el-form-item label="单次最大删除数">
          <el-input-number v-model="deleteForm.maxDeleteRows" :max="50000" :min="1" />
        </el-form-item>
        <el-form-item label="二次确认文本">
          <el-input v-model="deleteForm.confirmText" :placeholder="`请输入 DELETE ${deleteForm.tableName}`" />
        </el-form-item>
      </el-form>

      <el-alert
        type="error"
        show-icon
        :closable="false"
        title="此操作不可恢复，请确认日期范围和表名无误。"
      />

      <template #footer>
        <div>
          <el-button @click="deleteDialogVisible = false">取消</el-button>
          <el-button :loading="deleteSubmitting" type="danger" @click="submitDelete">确认真删除</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { autoFixDBInspector, deleteRecordsByRange, getDBInspectorOverview } from '@/api/system/dbInspector'

const loading = ref(false)
const deleteSubmitting = ref(false)
const autoFixingTarget = ref('')

const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const tableData = ref([])

const searchInfo = reactive({
  keyword: ''
})

const databaseStatus = ref({})
const redisStatus = ref({})

const redisKeyspaces = computed(() => redisStatus.value.keyspaces || [])

const deleteDialogVisible = ref(false)
const deleteForm = reactive({
  tableName: '',
  dateColumns: [],
  timeColumn: '',
  dateRange: [],
  deleteRelatedFiles: true,
  confirmText: '',
  maxDeleteRows: 5000
})

const resolveErrorMessage = (err, fallback = '操作失败') => {
  return err?.response?.data?.msg || err?.message || fallback
}

const fetchOverview = async () => {
  loading.value = true
  try {
    const res = await getDBInspectorOverview({
      page: page.value,
      pageSize: pageSize.value,
      keyword: searchInfo.keyword
    })
    if (res.code === 0) {
      databaseStatus.value = res.data.database || {}
      redisStatus.value = res.data.redis || {}
      tableData.value = res.data.tableList || []
      total.value = res.data.total || 0
      page.value = res.data.page || page.value
      pageSize.value = res.data.pageSize || pageSize.value
    } else {
      ElMessage.error(res.msg || '获取数据库巡检数据失败')
    }
  } catch (err) {
    ElMessage.error(resolveErrorMessage(err, '获取数据库巡检数据失败'))
  } finally {
    loading.value = false
  }
}

const onSubmit = () => {
  page.value = 1
  fetchOverview()
}

const onReset = () => {
  searchInfo.keyword = ''
  page.value = 1
  pageSize.value = 20
  fetchOverview()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  fetchOverview()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchOverview()
}

const handleAutoFix = async (target) => {
  try {
    await ElMessageBox.confirm(
      `将尝试自动修复${target === 'database' ? '数据库' : 'Redis'}连接，是否继续？`,
      '自动修复确认',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
  } catch {
    return
  }

  autoFixingTarget.value = target
  try {
    const res = await autoFixDBInspector({ target })
    if (res.code === 0) {
      ElMessage.success(res.msg || '自动修复执行成功')
    } else {
      ElMessage.error(res.msg || '自动修复执行失败')
    }
    await fetchOverview()
  } catch (err) {
    ElMessage.error(resolveErrorMessage(err, '自动修复请求失败'))
  } finally {
    autoFixingTarget.value = ''
  }
}

const openDeleteDialog = (row) => {
  deleteForm.tableName = row.tableName
  deleteForm.dateColumns = row.dateColumns || []
  deleteForm.timeColumn = (row.dateColumns && row.dateColumns[0]) || ''
  deleteForm.dateRange = []
  deleteForm.deleteRelatedFiles = true
  deleteForm.confirmText = ''
  deleteForm.maxDeleteRows = 5000
  deleteDialogVisible.value = true
}

const submitDelete = async () => {
  if (!deleteForm.timeColumn) {
    ElMessage.warning('请选择时间字段')
    return
  }
  if (!deleteForm.dateRange || deleteForm.dateRange.length !== 2) {
    ElMessage.warning('请选择完整的日期范围')
    return
  }
  const expected = `DELETE ${deleteForm.tableName}`
  if (deleteForm.confirmText !== expected) {
    ElMessage.warning(`确认文本不正确，请输入 ${expected}`)
    return
  }

  try {
    await ElMessageBox.confirm(
      `将对 ${deleteForm.tableName} 执行不可恢复的真删除，是否继续？`,
      '高危操作确认',
      {
        confirmButtonText: '继续删除',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
  } catch {
    return
  }

  deleteSubmitting.value = true
  try {
    const res = await deleteRecordsByRange({
      tableName: deleteForm.tableName,
      timeColumn: deleteForm.timeColumn,
      startTime: deleteForm.dateRange[0],
      endTime: deleteForm.dateRange[1],
      confirmText: deleteForm.confirmText,
      deleteRelatedFiles: deleteForm.deleteRelatedFiles,
      maxDeleteRows: deleteForm.maxDeleteRows
    })

    if (res.code === 0) {
      const data = res.data || {}
      let msg = `已删除 ${data.deletedRows || 0} 条记录`
      if (deleteForm.deleteRelatedFiles) {
        msg += `，文件删除成功 ${data.deletedFileCount || 0} 个`
        if ((data.failedFileCount || 0) > 0) {
          msg += `，失败 ${data.failedFileCount || 0} 个`
        }
      }
      ElMessage.success(msg)
      if ((data.fileDeleteWarning || []).length > 0) {
        ElMessage.warning(data.fileDeleteWarning[0])
      }
      deleteDialogVisible.value = false
      await fetchOverview()
    } else {
      ElMessage.error(res.msg || '真删除失败')
    }
  } catch (err) {
    ElMessage.error(resolveErrorMessage(err, '真删除请求失败'))
  } finally {
    deleteSubmitting.value = false
  }
}

fetchOverview()
</script>

<style scoped>
.status-row {
  margin-bottom: 12px;
}

.status-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.status-row-item {
  color: var(--el-text-color-regular);
}

.status-warning {
  color: var(--el-color-warning);
}

.status-guides {
  padding: 8px;
  background: var(--el-fill-color-lighter);
  border-radius: 6px;
}

.guide-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.guide-item {
  color: var(--el-text-color-secondary);
  line-height: 1.7;
}

.status-action {
  margin-top: 6px;
}

.keyspace-box {
  margin-top: 10px;
}

.inline-tip {
  margin-left: 10px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.w-full {
  width: 100%;
}
</style>
