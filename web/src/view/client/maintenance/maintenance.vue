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
        </el-form-item>

        <el-form-item label="弹窗开关:">
          <el-switch v-model="formData.maintenance_popup_enabled" active-value="true" inactive-value="false" />
        </el-form-item>

        <el-form-item label="弹窗标题 (JSON多语言):">
          <el-input v-model="formData.maintenance_popup_title" type="textarea" :rows="3" placeholder='{"mn":"...","zh":"系统维护中","en":"System Maintenance"}' />
        </el-form-item>

        <el-form-item label="弹窗内容 (JSON多语言):">
          <el-input v-model="formData.maintenance_popup_content" type="textarea" :rows="4" placeholder='{"mn":"...","zh":"系统正在维护，请稍后再试","en":"..."}' />
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
import { ElMessage } from 'element-plus'
import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({
  name: 'MaintenanceSetting'
})

const saving = ref(false)

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
  }
}

const saveAll = async () => {
  saving.value = true
  try {
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
  loadConfigs()
})
</script>
