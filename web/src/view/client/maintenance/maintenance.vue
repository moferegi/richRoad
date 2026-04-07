<template>
  <div>
    <div class="gva-table-box" style="padding: 20px;">
      <h3 style="margin-bottom: 20px;">系统维护模式设置</h3>
      <el-form :model="formData" label-position="right" label-width="180px" style="max-width: 700px;">
        <el-form-item label="维护模式开关:">
          <el-switch v-model="formData.maintenance_enabled" active-value="true" inactive-value="false" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">开启后客户端用户将无法访问业务功能</span>
        </el-form-item>

        <el-form-item label="维护背景图:">
          <SelectImage v-model="formData.maintenance_bg_image" file-type="image" />
          <el-input v-model="formData.maintenance_bg_image" :clearable="true" placeholder="也可直接输入外部图片URL" style="margin-top: 8px;" />
        </el-form-item>

        <el-form-item label="弹窗开关:">
          <el-switch v-model="formData.maintenance_popup_enabled" active-value="true" inactive-value="false" />
        </el-form-item>

        <el-form-item label="弹窗标题(多语言):">
          <el-tabs v-if="enabledLangs.length" type="border-card" style="width:100%;">
            <el-tab-pane v-for="lang in enabledLangs" :key="lang.code" :label="lang.name">
              <el-input v-model="popupTitleI18n[lang.code]" :placeholder="`${lang.name} 弹窗标题`" />
            </el-tab-pane>
          </el-tabs>
          <el-input v-else v-model="formData.maintenance_popup_title" type="textarea" :rows="3" placeholder='{"zh":"系统维护中","en":"System Maintenance"}' />
        </el-form-item>

        <el-form-item label="弹窗内容(多语言):">
          <el-tabs v-if="enabledLangs.length" type="border-card" style="width:100%;">
            <el-tab-pane v-for="lang in enabledLangs" :key="lang.code" :label="lang.name">
              <el-input v-model="popupContentI18n[lang.code]" type="textarea" :rows="3" :placeholder="`${lang.name} 弹窗内容`" />
            </el-tab-pane>
          </el-tabs>
          <el-input v-else v-model="formData.maintenance_popup_content" type="textarea" :rows="4" placeholder='{"zh":"系统正在维护，请稍后再试","en":"..."}' />
        </el-form-item>

        <el-form-item label="进入首页按钮:">
          <el-switch v-model="formData.maintenance_home_btn_enabled" active-value="true" inactive-value="false" />
          <span style="margin-left: 10px; color: #909399; font-size: 12px;">是否允许用户在维护期间进入首页</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="saving" @click="saveAll">保存设置</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getSysConfigList, updateSysConfig } from '@/api/client/sysConfig'
import { getEnabledLanguages } from '@/api/client/language'
import { ElMessage } from 'element-plus'
import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({
  name: 'MaintenanceSetting'
})

const saving = ref(false)

// === 多语言支持 ===
const enabledLangs = ref([])
const popupTitleI18n = ref({})
const popupContentI18n = ref({})

const loadLangs = async () => {
  try {
    const res = await getEnabledLanguages()
    if (res.code === 0) {
      enabledLangs.value = res.data || []
    }
  } catch (e) { /* ignore */ }
}

const parseI18nJson = (jsonStr) => {
  if (!jsonStr) return {}
  try { return JSON.parse(jsonStr) } catch { return {} }
}

const serializeI18nJson = (obj) => {
  const filtered = {}
  for (const [k, v] of Object.entries(obj)) {
    if (v) filtered[k] = v
  }
  return Object.keys(filtered).length ? JSON.stringify(filtered) : ''
}

const configKeys = [
  'maintenance_enabled',
  'maintenance_bg_image',
  'maintenance_popup_enabled',
  'maintenance_popup_title',
  'maintenance_popup_content',
  'maintenance_home_btn_enabled',
]

const formData = ref({
  maintenance_enabled: 'false',
  maintenance_bg_image: '',
  maintenance_popup_enabled: 'false',
  maintenance_popup_title: '',
  maintenance_popup_content: '',
  maintenance_home_btn_enabled: 'false',
})

// ID映射，用于更新时找到对应记录
const configIdMap = ref({})

const loadConfigs = async () => {
  const res = await getSysConfigList({ page: 1, pageSize: 100, configGroup: 'maintenance' })
  if (res.code === 0 && res.data.list) {
    res.data.list.forEach(item => {
      if (configKeys.includes(item.configKey)) {
        formData.value[item.configKey] = item.configValue
        configIdMap.value[item.configKey] = item.ID
      }
    })
    // 解析多语言字段
    popupTitleI18n.value = parseI18nJson(formData.value.maintenance_popup_title)
    popupContentI18n.value = parseI18nJson(formData.value.maintenance_popup_content)
  }
}

const saveAll = async () => {
  saving.value = true
  try {
    // 序列化多语言字段
    if (enabledLangs.value.length) {
      formData.value.maintenance_popup_title = serializeI18nJson(popupTitleI18n.value)
      formData.value.maintenance_popup_content = serializeI18nJson(popupContentI18n.value)
    }
    for (const key of configKeys) {
      const id = configIdMap.value[key]
      if (id) {
        await updateSysConfig({ ID: id, configValue: formData.value[key] })
      }
    }
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadLangs()
  loadConfigs()
})
</script>
