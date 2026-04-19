<template>
  <div>
    <!-- 分组标签切换 -->
    <div class="gva-search-box">
      <el-radio-group v-model="activeGroup" @change="handleGroupChange">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button v-for="g in groupList" :key="g.value" :label="g.value">{{ g.label }}</el-radio-button>
      </el-radio-group>
    </div>

    <div class="gva-table-box">
      <el-table :data="tableData" stripe>
        <el-table-column prop="configGroup" label="配置组" width="120">
          <template #default="scope">
            <el-tag size="small" :type="groupTagType(scope.row.configGroup)">{{ groupLabel(scope.row.configGroup) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="configName" label="配置名称" width="180" />
        <el-table-column prop="configKey" label="配置键" width="220">
          <template #default="scope">
            <span class="config-key">{{ scope.row.configKey }}</span>
          </template>
        </el-table-column>
        <el-table-column label="配置值" min-width="250">
          <template #default="scope">
            <!-- 布尔类型展示开关 -->
            <template v-if="isBooleanConfig(scope.row)">
              <el-switch
                :model-value="scope.row.configValue === 'true'"
                @change="(val) => quickToggle(scope.row, val)"
                active-text="开"
                inactive-text="关"
              />
            </template>
            <!-- 颜色类型 -->
            <template v-else-if="isColorConfig(scope.row)">
              <div class="color-preview">
                <span class="color-dot" :style="{ background: scope.row.configValue }" />
                <span>{{ scope.row.configValue }}</span>
              </div>
            </template>
            <!-- 长文本截断 -->
            <template v-else-if="scope.row.configValue && scope.row.configValue.length > 60">
              <el-tooltip :content="scope.row.configValue" placement="top">
                <span>{{ scope.row.configValue.substring(0, 60) }}...</span>
              </el-tooltip>
            </template>
            <template v-else>
              <span>{{ scope.row.configValue || '-' }}</span>
            </template>
          </template>
        </el-table-column>
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
    <el-dialog v-model="editVisible" title="编辑参数" width="600px">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="配置键">
          <el-input :model-value="editForm.configKey" disabled />
        </el-form-item>
        <el-form-item label="配置名称">
          <el-input :model-value="editForm.configName" disabled />
        </el-form-item>
        <el-form-item label="配置值">
          <!-- 布尔类型用开关 -->
          <template v-if="isBooleanConfig(editForm)">
            <el-switch
              v-model="editBoolValue"
              active-text="开"
              inactive-text="关"
            />
          </template>
          <!-- 颜色类型用颜色选择器 -->
          <template v-else-if="isColorConfig(editForm)">
            <el-color-picker v-model="editForm.configValue" show-alpha />
            <el-input v-model="editForm.configValue" class="ml-2" style="width: 200px" />
          </template>
          <!-- JSON多语言类型 -->
          <template v-else-if="isJsonConfig(editForm)">
            <div class="json-editor">
              <div v-for="(val, lang) in editJsonValue" :key="lang" class="json-row">
                <el-tag size="small" class="mr-2">{{ lang }}</el-tag>
                <el-input v-model="editJsonValue[lang]" style="flex:1" />
              </div>
              <div class="json-row mt-2">
                <el-input v-model="newJsonLang" placeholder="语言代码" style="width: 100px" class="mr-2" />
                <el-input v-model="newJsonVal" placeholder="值" style="flex:1" class="mr-2" />
                <el-button size="small" @click="addJsonLang">添加</el-button>
              </div>
            </div>
          </template>
          <!-- 数字类型 -->
          <template v-else-if="isNumberConfig(editForm)">
            <el-input-number v-model="editNumberValue" :min="0" />
          </template>
          <!-- 普通文本 -->
          <template v-else>
            <el-input v-model="editForm.configValue" type="textarea" :rows="3" />
          </template>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="editForm.remark" type="textarea" :rows="2" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { getSysConfigList, updateSysConfig } from '@/api/client/sysConfig'
import { ElMessage } from 'element-plus'

// 配置组定义
const groupList = [
  { value: 'system', label: '系统设置', type: '' },
  { value: 'maintenance', label: '维护模式', type: 'danger' },
  { value: 'security', label: '安全设置', type: 'danger' },
  { value: 'auth', label: '认证设置', type: 'warning' },
  { value: 'payment', label: '支付设置', type: 'success' },
  { value: 'points', label: '积分设置', type: 'info' },
  { value: 'order', label: '订单设置', type: '' },
  { value: 'display', label: '显示设置', type: 'warning' },
  { value: 'announcement', label: '公告设置', type: 'info' },
  { value: 'invite', label: '邀请设置', type: 'success' },
]

const groupTagType = (group) => {
  const found = groupList.find(g => g.value === group)
  return found ? found.type : 'info'
}

const groupLabel = (group) => {
  const found = groupList.find(g => g.value === group)
  return found ? found.label : group
}

// 布尔键列表
const booleanKeys = [
  'maintenance_enabled', 'maintenance_popup_enabled', 'maintenance_home_btn_enabled',
  'phone_login_enabled', 'username_login_enabled', 'password_change_enabled',
  'sign_in_enabled', 'announcement_enabled'
]
const isBooleanConfig = (row) => booleanKeys.includes(row.configKey)

// 颜色键列表
const colorKeys = ['announcement_text_color', 'payment_tip_text_color']
const isColorConfig = (row) => colorKeys.includes(row.configKey)

// JSON键列表（多语言配置）
const jsonKeys = [
  'payment_tip_text', 'maintenance_popup_title', 'maintenance_popup_content',
  'announcement_content', 'maintenance_message',
  'username_regex_tip', 'password_regex_tip'
]
const isJsonConfig = (row) => jsonKeys.includes(row.configKey)

// 数字键列表
const numberKeys = [
  'captcha_expiry_seconds', 'captcha_rate_limit', 'register_ip_limit',
  'login_fail_max', 'login_fail_wait_seconds', 'points_exchange_rate',
  'order_close_minutes', 'presale_home_count', 'announcement_speed',
  'invite_reward_points', 'payment_tip_text_size'
]
const isNumberConfig = (row) => numberKeys.includes(row.configKey)

const activeGroup = ref('')
const searchInfo = ref({
  page: 1,
  pageSize: 50,
  configGroup: '',
  configKey: '',
})

const tableData = ref([])
const total = ref(0)
const editVisible = ref(false)
const editForm = ref({})
const editBoolValue = ref(false)
const editNumberValue = ref(0)
const editJsonValue = reactive({})
const newJsonLang = ref('')
const newJsonVal = ref('')

const handleGroupChange = (val) => {
  searchInfo.value.configGroup = val
  searchInfo.value.page = 1
  getList()
}

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
  if (isBooleanConfig(row)) {
    editBoolValue.value = row.configValue === 'true'
  }
  if (isNumberConfig(row)) {
    editNumberValue.value = parseInt(row.configValue) || 0
  }
  if (isJsonConfig(row)) {
    try {
      const parsed = JSON.parse(row.configValue || '{}')
      Object.keys(editJsonValue).forEach(k => delete editJsonValue[k])
      Object.assign(editJsonValue, parsed)
    } catch {
      Object.keys(editJsonValue).forEach(k => delete editJsonValue[k])
    }
  }
  newJsonLang.value = ''
  newJsonVal.value = ''
  editVisible.value = true
}

const addJsonLang = () => {
  if (newJsonLang.value) {
    editJsonValue[newJsonLang.value] = newJsonVal.value
    newJsonLang.value = ''
    newJsonVal.value = ''
  }
}

const quickToggle = async (row, val) => {
  const res = await updateSysConfig({
    id: row.ID,
    configValue: val ? 'true' : 'false',
    remark: row.remark,
  })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    getList()
  }
}

const handleSave = async () => {
  let configValue = editForm.value.configValue
  if (isBooleanConfig(editForm.value)) {
    configValue = editBoolValue.value ? 'true' : 'false'
  } else if (isNumberConfig(editForm.value)) {
    configValue = String(editNumberValue.value)
  } else if (isJsonConfig(editForm.value)) {
    configValue = JSON.stringify(editJsonValue)
  }
  const res = await updateSysConfig({
    id: editForm.value.ID,
    configValue: configValue,
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

<style scoped>
.config-key {
  font-family: monospace;
  font-size: 12px;
  color: #909399;
}
.color-preview {
  display: flex;
  align-items: center;
  gap: 8px;
}
.color-dot {
  display: inline-block;
  width: 16px;
  height: 16px;
  border-radius: 3px;
  border: 1px solid #dcdfe6;
}
.json-editor {
  width: 100%;
}
.json-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}
</style>
