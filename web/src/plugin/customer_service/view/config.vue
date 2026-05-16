<template>
  <div class="p-4">
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold text-gray-700">客服系统配置</span>
          <el-tag type="info" size="small">实时生效</el-tag>
        </div>
      </template>

      <el-form label-width="180px" class="max-w-xl">
        <el-form-item label="启用平台内置客服">
          <el-switch v-model="form.platEnabled" />
          <span class="ml-3 text-xs text-gray-400">关闭后用户无法创建平台客服会话</span>
        </el-form-item>

        <el-form-item label="会话空闲自动关闭（分钟）">
          <el-input-number
            v-model="form.autoCloseMinutes"
            :min="0"
            :max="1440"
            :step="5"
            controls-position="right"
          />
          <span class="ml-3 text-xs text-gray-400">0 表示不自动关闭</span>
        </el-form-item>

        <el-form-item label="图片允许扩展名">
          <el-input v-model="form.uploadAllowExt" placeholder="jpg,jpeg,png,webp,gif" />
          <span class="ml-3 text-xs text-gray-400">英文逗号分隔，支持填写不带点扩展名</span>
        </el-form-item>

        <el-form-item label="单图最大体积（MB）">
          <el-input-number
            v-model="form.uploadMaxSizeMB"
            :min="1"
            :max="50"
            :step="1"
            controls-position="right"
          />
          <span class="ml-3 text-xs text-gray-400">客户端和服务端都会按该配置校验</span>
        </el-form-item>

        <el-form-item label="客服默认头像">
          <SelectImage
            v-model="form.defaultAvatarUrl"
            file-type="image"
            :default-folder="CUSTOMER_SERVICE_AVATAR_UPLOAD_FOLDER"
            :fixed-upload-folder="true"
          />
          <span class="ml-3 text-xs text-gray-400">留空时客户端会使用首字母随机底色头像</span>
        </el-form-item>

        <el-form-item>
          <el-button :loading="saving" type="primary" @click="handleSave">保存配置</el-button>
          <el-button :loading="loading" @click="loadConfig">刷新</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getCsConfig, updateCsConfig } from '@/api/customerService'
import SelectImage from '@/components/selectImage/selectImage.vue'

const CUSTOMER_SERVICE_AVATAR_UPLOAD_FOLDER = 'cloth-on/web-else'

const loading = ref(false)
const saving = ref(false)
const form = reactive({
  platEnabled: false,
  autoCloseMinutes: 30,
  uploadMaxSizeMB: 5,
  uploadAllowExt: 'jpg,jpeg,png,webp,gif',
  defaultAvatarUrl: '',
})

async function loadConfig() {
  loading.value = true
  try {
    const res = await getCsConfig()
    if (res.code !== 0) return
    form.platEnabled = !!res.data?.platEnabled
    form.autoCloseMinutes = Number(res.data?.autoCloseMinutes ?? 30)
    form.uploadMaxSizeMB = Number(res.data?.uploadMaxSizeMB ?? 5)
    form.uploadAllowExt = String(res.data?.uploadAllowExt || 'jpg,jpeg,png,webp,gif')
    form.defaultAvatarUrl = String(res.data?.defaultAvatarUrl || '')
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const res = await updateCsConfig({
      platEnabled: form.platEnabled,
      autoCloseMinutes: Number(form.autoCloseMinutes || 0),
      uploadMaxSizeMB: Number(form.uploadMaxSizeMB || 5),
      uploadAllowExt: String(form.uploadAllowExt || '').trim(),
      defaultAvatarUrl: String(form.defaultAvatarUrl || '').trim(),
    })
    if (res.code === 0) {
      ElMessage.success('保存成功')
    }
  } finally {
    saving.value = false
  }
}

onMounted(loadConfig)
</script>
