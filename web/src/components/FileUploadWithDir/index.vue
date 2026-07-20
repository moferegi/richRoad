<template>
  <div class="file-upload-with-dir">
    <div class="upload-area">
      <el-upload
        :show-file-list="false"
        :http-request="handleUpload"
        :accept="accept"
        :disabled="uploading"
      >
        <el-button :loading="uploading" :type="modelValue ? 'default' : 'primary'" :icon="Upload">
          {{ modelValue ? '重新上传' : '选择文件上传' }}
        </el-button>
      </el-upload>

      <div v-if="modelValue" class="uploaded-file">
        <el-icon class="check-icon"><CircleCheckFilled /></el-icon>
        <span class="file-path" :title="modelValue">{{ modelValue }}</span>
        <el-button link type="danger" size="small" @click="handleClear">清除</el-button>
      </div>
    </div>

    <div class="folder-bar">
      <span class="folder-label">上传目录：</span>
      <el-input
        v-model="currentFolder"
        :disabled="props.fixedUploadFolder"
        size="small"
        class="folder-input"
        placeholder="输入上传目录路径"
        clearable
      />
      <el-button
        v-if="!props.fixedUploadFolder"
        link
        type="primary"
        size="small"
        @click="showFolderDialog = true"
      >
        快捷选择
      </el-button>
    </div>

    <!-- 修改目录弹窗 -->
    <el-dialog v-model="showFolderDialog" title="修改上传目录" width="480px" append-to-body>
      <el-form label-width="80px">
        <el-form-item label="上传目录">
          <el-input v-model="folderInput" placeholder="如：english-learn/pic" clearable />
        </el-form-item>
        <el-form-item label="当前目录">
          <el-tag type="info">{{ currentFolder || '(根目录)' }}</el-tag>
        </el-form-item>
        <el-form-item label="快捷选择">
          <div class="quick-options">
            <el-button
              v-for="opt in quickFolders"
              :key="opt.value"
              size="small"
              :type="currentFolder === opt.value ? 'primary' : 'default'"
              @click="folderInput = opt.value"
            >
              {{ opt.label }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFolderDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmFolder">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Upload, CircleCheckFilled } from '@element-plus/icons-vue'
import { uploadFile } from '@/api/fileUploadAndDownload'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  /** 默认上传目录前缀 */
  defaultFolder: {
    type: String,
    default: ''
  },
  /** 是否固定上传目录（禁止修改） */
  fixedUploadFolder: {
    type: Boolean,
    default: false
  },
  /** 接受的文件类型，如 image/* */
  accept: {
    type: String,
    default: ''
  },
  /** 快捷目录选项 */
  quickFolders: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue'])

const uploading = ref(false)
const currentFolder = ref('')
const showFolderDialog = ref(false)
const folderInput = ref('')

// 初始化目录
watch(() => props.defaultFolder, (val) => {
  if (val && !currentFolder.value) {
    currentFolder.value = val
  }
}, { immediate: true })

// 提取上传返回的URL
const extractURL = (res) => {
  return String(
    res?.data?.file?.url ||
    res?.data?.url ||
    res?.file?.url ||
    res?.url ||
    ''
  ).trim()
}

// 处理上传
const handleUpload = async (options) => {
  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', options.file)
    if (currentFolder.value) {
      formData.append('folder', currentFolder.value)
    }

    const res = await uploadFile(formData)
    if (res.code !== 0) {
      throw new Error(res.msg || '上传失败')
    }

    const url = extractURL(res)
    if (!url) {
      throw new Error('上传成功但未返回URL')
    }

    emit('update:modelValue', url)
    ElMessage.success('上传成功')

    if (typeof options.onSuccess === 'function') {
      options.onSuccess({ url })
    }
  } catch (error) {
    ElMessage.error(error?.message || '上传失败')
    if (typeof options.onError === 'function') {
      options.onError(error)
    }
  } finally {
    uploading.value = false
  }
}

// 清除已上传文件
const handleClear = () => {
  emit('update:modelValue', '')
}

// 确认修改目录
const confirmFolder = () => {
  currentFolder.value = folderInput.value || props.defaultFolder || ''
  showFolderDialog.value = false
}

// 打开发送框时同步值
watch(showFolderDialog, (val) => {
  if (val) {
    folderInput.value = currentFolder.value || props.defaultFolder || ''
  }
})
</script>

<style scoped>
.file-upload-with-dir {
  width: 100%;
}

.upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.uploaded-file {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.check-icon {
  color: #67c23a;
  font-size: 16px;
}

.file-path {
  color: #409eff;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.folder-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  font-size: 13px;
  color: #909399;
}

.folder-label {
  white-space: nowrap;
}

.folder-input {
  width: 220px;
}

.quick-options {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
</style>
