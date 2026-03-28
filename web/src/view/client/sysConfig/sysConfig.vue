<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="getList">
        <el-form-item label="配置组">
          <el-input v-model="searchInfo.configGroup" placeholder="配置组" clearable />
        </el-form-item>
        <el-form-item label="配置键">
          <el-input v-model="searchInfo.configKey" placeholder="配置键" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getList">查询</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" stripe>
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column prop="configGroup" label="配置组" width="120" />
        <el-table-column prop="configKey" label="配置键" width="200" />
        <el-table-column prop="configName" label="配置名称" width="180" />
        <el-table-column prop="configValue" label="配置值" width="150" />
        <el-table-column prop="remark" label="备注" min-width="200" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openEdit(scope.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="searchInfo.page"
          :page-size="searchInfo.pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑参数" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="配置键">
          <el-input :model-value="editForm.configKey" disabled />
        </el-form-item>
        <el-form-item label="配置名称">
          <el-input :model-value="editForm.configName" disabled />
        </el-form-item>
        <el-form-item label="配置值">
          <el-input v-model="editForm.configValue" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="editForm.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSysConfigList, updateSysConfig } from '@/api/client/sysConfig'
import { ElMessage } from 'element-plus'

const searchInfo = ref({
  page: 1,
  pageSize: 10,
  configGroup: '',
  configKey: '',
})

const tableData = ref([])
const total = ref(0)
const editVisible = ref(false)
const editForm = ref({})

const getList = async () => {
  const res = await getSysConfigList(searchInfo.value)
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

const openEdit = (row) => {
  editForm.value = { ...row }
  editVisible.value = true
}

const handleSave = async () => {
  const res = await updateSysConfig({
    id: editForm.value.ID,
    configValue: editForm.value.configValue,
    remark: editForm.value.remark,
  })
  if (res.code === 0) {
    ElMessage.success('保存成功')
    editVisible.value = false
    getList()
  }
}

onMounted(() => {
  getList()
})
</script>
