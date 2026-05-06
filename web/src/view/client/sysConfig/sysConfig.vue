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
            <!-- 密钥类脱敏展示 -->
            <template v-else-if="isSecretConfig(scope.row)">
              <span>{{ maskSecretValue(scope.row.configValue) }}</span>
            </template>
            <!-- 试衣模型可视化摘要 -->
            <template v-else-if="isTryonModelsConfig(scope.row)">
              <div class="tryon-model-summary">
                <el-tag size="small" type="warning">{{ tryonModelsCount(scope.row.configValue) }} 个模型</el-tag>
                <span class="tryon-model-summary-text">启用 {{ enabledTryonModelsCount(scope.row.configValue) }} 个</span>
              </div>
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
          <!-- 试衣模型可视化编辑 -->
          <template v-else-if="isTryonModelsConfig(editForm)">
            <div class="tryon-model-editor">
              <div class="tryon-model-toolbar">
                <el-button type="primary" plain size="small" @click="addTryonModel">新增模型</el-button>
              </div>

              <div v-if="tryonModels.length === 0" class="tryon-model-empty">
                暂无模型，点击“新增模型”开始配置
              </div>

              <el-collapse v-else>
                <el-collapse-item
                  v-for="(model, index) in tryonModels"
                  :key="model.__uid"
                  :name="model.__uid"
                >
                  <template #title>
                    <div class="tryon-model-title">
                      <span>{{ model.name.zh || model.name.en || model.key || ('模型' + (index + 1)) }}</span>
                      <el-tag size="small" :type="model.enabled ? 'success' : 'info'">
                        {{ model.enabled ? '启用' : '关闭' }}
                      </el-tag>
                    </div>
                  </template>

                  <div class="tryon-model-panel">
                    <div class="tryon-model-grid">
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">模型键 key</span>
                        <el-input v-model="model.key" placeholder="如 aliyun_aitryon" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">模型类型 model</span>
                        <el-input v-model="model.model" placeholder="如 aitryon / aitryon-plus" />
                      </div>

                      <div class="tryon-model-field">
                        <span class="tryon-model-label">提供商 provider</span>
                        <el-input v-model="model.provider" placeholder="如 aliyun" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">运行模式 mode</span>
                        <el-select v-model="model.mode" style="width: 100%">
                          <el-option label="prod" value="prod" />
                          <el-option label="mock_success" value="mock_success" />
                        </el-select>
                      </div>

                      <div class="tryon-model-field">
                        <span class="tryon-model-label">单次消耗 cost</span>
                        <el-input-number v-model="model.cost" :min="0" :step="1" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">是否启用</span>
                        <el-switch v-model="model.enabled" active-text="开" inactive-text="关" />
                      </div>

                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">场景 scenes</span>
                        <el-checkbox-group v-model="model.scenes">
                          <el-checkbox value="clothes">试衣 clothes</el-checkbox>
                          <el-checkbox value="shoes">试鞋 shoes</el-checkbox>
                          <el-checkbox value="takeoff">取衣 takeoff</el-checkbox>
                        </el-checkbox-group>
                      </div>

                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">模型地址 url</span>
                        <el-input v-model="model.url" placeholder="如 https://dashscope.aliyuncs.com/api/v1/services/..." />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">查询地址 taskQueryUrl</span>
                        <el-input v-model="model.taskQueryUrl" placeholder="如 https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">模型 token</span>
                        <el-input v-model="model.token" type="password" show-password placeholder="留空则回退使用 tryon_provider_token" />
                      </div>

                      <div class="tryon-model-field">
                        <span class="tryon-model-label">分辨率 resolution</span>
                        <el-input-number v-model="model.resolution" :min="-1" :step="1" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">人脸修复 restoreFace</span>
                        <el-switch v-model="model.restoreFace" active-text="开" inactive-text="关" />
                      </div>

                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">取衣分割 clothesType</span>
                        <el-checkbox-group v-model="model.clothesType">
                          <el-checkbox value="upper">upper</el-checkbox>
                          <el-checkbox value="lower">lower</el-checkbox>
                        </el-checkbox-group>
                      </div>

                      <div class="tryon-model-subtitle">名称多语言 name</div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">中文 zh</span>
                        <el-input v-model="model.name.zh" placeholder="中文名称" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">英文 en</span>
                        <el-input v-model="model.name.en" placeholder="English name" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">蒙文 mn</span>
                        <el-input v-model="model.name.mn" placeholder="Монгол нэр" />
                      </div>

                      <div class="tryon-model-subtitle">说明多语言 desc</div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">中文 zh</span>
                        <el-input v-model="model.desc.zh" type="textarea" :rows="2" placeholder="中文说明" />
                      </div>
                      <div class="tryon-model-field">
                        <span class="tryon-model-label">英文 en</span>
                        <el-input v-model="model.desc.en" type="textarea" :rows="2" placeholder="English description" />
                      </div>
                      <div class="tryon-model-field full">
                        <span class="tryon-model-label">蒙文 mn</span>
                        <el-input v-model="model.desc.mn" type="textarea" :rows="2" placeholder="Монгол тайлбар" />
                      </div>
                    </div>

                    <div class="tryon-model-actions">
                      <el-button size="small" @click="cloneTryonModel(index)">复制</el-button>
                      <el-button size="small" type="danger" plain @click="removeTryonModel(index)">删除</el-button>
                    </div>
                  </div>
                </el-collapse-item>
              </el-collapse>
            </div>
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
          <!-- 密钥类型 -->
          <template v-else-if="isSecretConfig(editForm)">
            <el-input v-model="editForm.configValue" type="password" show-password />
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
  { value: 'tryon', label: '试衣设置', type: 'warning' },
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
  'sign_in_enabled', 'announcement_enabled',
  'payment_auto_enabled', 'payment_manual_qrcode_enabled', 'payment_manual_contact_enabled',
  'payment_wechat_enabled', 'payment_alipay_enabled', 'payment_bank_cn_enabled',
  'payment_bank_us_enabled', 'payment_bank_mn_enabled', 'payment_paypal_enabled'
]
const isBooleanConfig = (row) => booleanKeys.includes(row.configKey)

// 颜色键列表
const colorKeys = ['announcement_text_color', 'payment_tip_text_color']
const isColorConfig = (row) => colorKeys.includes(row.configKey)

const isTryonModelsConfig = (row) => row?.configKey === 'tryon_models'

// 密钥键列表
const secretKeys = ['tryon_provider_token']
const isSecretConfig = (row) => secretKeys.includes(row.configKey)

const maskSecretValue = (value) => {
  if (!value) {
    return '-'
  }
  if (value.length <= 8) {
    return '*'.repeat(value.length)
  }
  return `${value.slice(0, 4)}${'*'.repeat(Math.max(value.length - 8, 4))}${value.slice(-4)}`
}

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
  'invite_reward_points', 'payment_tip_text_size',
  'tryon_guest_init_points', 'tryon_register_reward_points', 'tryon_cost_points'
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
const tryonModels = ref([])

const createTryonModelUid = () => `tryon_model_${Date.now()}_${Math.random().toString(16).slice(2, 8)}`

const toBool = (value, fallback = true) => {
  if (typeof value === 'boolean') {
    return value
  }
  if (value === undefined || value === null || value === '') {
    return fallback
  }
  const text = String(value).trim().toLowerCase()
  if (['true', '1', 'yes', 'on'].includes(text)) {
    return true
  }
  if (['false', '0', 'no', 'off'].includes(text)) {
    return false
  }
  return fallback
}

const toInt = (value, fallback = 0) => {
  const numberValue = Number.parseInt(value, 10)
  return Number.isFinite(numberValue) ? numberValue : fallback
}

const toStringArray = (value, fallback = []) => {
  if (Array.isArray(value)) {
    return Array.from(new Set(value.map(item => String(item).trim()).filter(Boolean)))
  }
  return [...fallback]
}

const normalizeI18nObject = (value) => {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    return {
      zh: String(value.zh || ''),
      en: String(value.en || ''),
      mn: String(value.mn || ''),
    }
  }
  return {
    zh: '',
    en: '',
    mn: '',
  }
}

const createDefaultTryonModel = () => ({
  __uid: createTryonModelUid(),
  key: '',
  enabled: true,
  scenes: ['clothes'],
  model: 'aitryon',
  name: { zh: '', en: '', mn: '' },
  desc: { zh: '', en: '', mn: '' },
  cost: 1,
  provider: 'aliyun',
  mode: 'prod',
  url: '',
  taskQueryUrl: '',
  token: '',
  resolution: -1,
  restoreFace: true,
  clothesType: ['upper'],
})

const normalizeTryonModel = (item = {}, index = 0) => {
  const defaultModel = createDefaultTryonModel()
  return {
    __uid: createTryonModelUid(),
    key: String(item.key || item.modelKey || ''),
    enabled: toBool(item.enabled, true),
    scenes: toStringArray(item.scenes, [String(item.sceneType || '').trim() || 'clothes']),
    model: String(item.model || defaultModel.model),
    name: normalizeI18nObject(item.name),
    desc: normalizeI18nObject(item.desc),
    cost: Math.max(0, toInt(item.cost, 1)),
    provider: String(item.provider || defaultModel.provider),
    mode: String(item.mode || defaultModel.mode),
    url: String(item.url || item.providerUrl || ''),
    taskQueryUrl: String(item.taskQueryUrl || ''),
    token: String(item.token || item.providerToken || ''),
    resolution: toInt(item.resolution, -1),
    restoreFace: toBool(item.restoreFace, true),
    clothesType: toStringArray(item.clothesType, ['upper']),
  }
}

const parseTryonModelsValue = (rawValue) => {
  try {
    const parsed = JSON.parse(rawValue || '[]')
    if (!Array.isArray(parsed)) {
      return []
    }
    return parsed.map((item, index) => normalizeTryonModel(item, index))
  } catch {
    return []
  }
}

const buildTryonModelsPayload = () => {
  return tryonModels.value.map((item) => ({
    key: String(item.key || '').trim(),
    enabled: !!item.enabled,
    scenes: toStringArray(item.scenes, ['clothes']),
    model: String(item.model || '').trim(),
    name: normalizeI18nObject(item.name),
    desc: normalizeI18nObject(item.desc),
    cost: Math.max(0, toInt(item.cost, 0)),
    provider: String(item.provider || '').trim(),
    mode: String(item.mode || '').trim(),
    url: String(item.url || '').trim(),
    taskQueryUrl: String(item.taskQueryUrl || '').trim(),
    token: String(item.token || '').trim(),
    resolution: toInt(item.resolution, -1),
    restoreFace: !!item.restoreFace,
    clothesType: toStringArray(item.clothesType, []),
  }))
}

const validateTryonModels = () => {
  for (let i = 0; i < tryonModels.value.length; i++) {
    const item = tryonModels.value[i]
    const modelIndex = i + 1
    if (!String(item.key || '').trim()) {
      ElMessage.warning(`第 ${modelIndex} 个模型缺少 key`)
      return false
    }
    if (!Array.isArray(item.scenes) || item.scenes.length === 0) {
      ElMessage.warning(`第 ${modelIndex} 个模型至少要选择一个场景`)
      return false
    }
  }
  return true
}

const tryonModelsCount = (rawValue) => parseTryonModelsValue(rawValue).length

const enabledTryonModelsCount = (rawValue) => {
  return parseTryonModelsValue(rawValue).filter(item => item.enabled).length
}

const addTryonModel = () => {
  tryonModels.value.push(createDefaultTryonModel())
}

const cloneTryonModel = (index) => {
  const item = tryonModels.value[index]
  if (!item) {
    return
  }
  const cloned = normalizeTryonModel(item)
  if (cloned.key) {
    cloned.key = `${cloned.key}_copy`
  }
  tryonModels.value.splice(index + 1, 0, cloned)
}

const removeTryonModel = (index) => {
  tryonModels.value.splice(index, 1)
}

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
    return
  }
  tableData.value = []
  total.value = 0
  ElMessage.error(res.msg || '获取系统参数失败')
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
  tryonModels.value = []
  if (isBooleanConfig(row)) {
    editBoolValue.value = row.configValue === 'true'
  }
  if (isNumberConfig(row)) {
    editNumberValue.value = parseInt(row.configValue) || 0
  }
  if (isTryonModelsConfig(row)) {
    tryonModels.value = parseTryonModelsValue(row.configValue)
    const rawValue = String(row.configValue || '').trim()
    if (rawValue && rawValue !== '[]' && tryonModels.value.length === 0) {
      ElMessage.warning('当前试衣模型配置格式异常，已按空列表打开，请确认后保存')
    }
  } else if (isJsonConfig(row)) {
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
  } else if (isTryonModelsConfig(editForm.value)) {
    if (!validateTryonModels()) {
      return
    }
    configValue = JSON.stringify(buildTryonModelsPayload())
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

.tryon-model-summary {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tryon-model-summary-text {
  color: #606266;
}

.tryon-model-editor {
  width: 100%;
}

.tryon-model-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 10px;
}

.tryon-model-empty {
  color: #909399;
  font-size: 13px;
  padding: 4px 0 8px;
}

.tryon-model-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.tryon-model-panel {
  padding: 6px 4px;
}

.tryon-model-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px 12px;
}

.tryon-model-field {
  min-width: 0;
}

.tryon-model-field.full {
  grid-column: 1 / -1;
}

.tryon-model-label {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-bottom: 6px;
}

.tryon-model-subtitle {
  grid-column: 1 / -1;
  font-size: 12px;
  color: #606266;
  font-weight: 600;
  margin-top: 2px;
}

.tryon-model-actions {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 900px) {
  .tryon-model-grid {
    grid-template-columns: 1fr;
  }
}
</style>
