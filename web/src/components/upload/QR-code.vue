<template>
  <div>
    <el-button type="primary" icon="iphone" @click="createQrCode"> 扫码上传</el-button>
  </div>

  <el-dialog v-model="dialogVisible" title="扫码上传" width="320px" :show-close="false" append-to-body :close-on-click-modal="false"
             draggable
  >
    <div class="m-2">
      <vue-qr :logoSrc="logoSrc"
              :size="291"
              :margin="0"
              :autoColor="true"
              :dotScale="1"
              :text="codeUrl"
              colorDark="green"
              colorLight="white"
              ref="qrcode"
      />
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="onFinished">完成上传</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import logoSrc from '@/assets/logo.png'
import { createScanUploadTicket } from '@/api/fileUploadAndDownload'
import { ElMessage } from 'element-plus'
import vueQr from 'vue-qr/src/packages/vue-qr.vue'
import { ref } from 'vue'

defineOptions({
  name: 'QRCodeUpload'
})

const emit = defineEmits(['on-success'])

const props = defineProps({
  classId: {
    type: Number,
    default: 0
  },
  folder: {
    type: String,
    default: ''
  }
})

const dialogVisible = ref(false)
const codeUrl = ref('')

const createQrCode = async () => {
  try {
    const local = window.location
    const { data } = await createScanUploadTicket({
      classId: props.classId,
      folder: props.folder
    })
    const ticket = data?.ticket
    if (!ticket) {
      throw new Error('empty ticket')
    }
    codeUrl.value = `${local.protocol}//${local.host}/#/scanUpload?ticket=${encodeURIComponent(ticket)}&t=${Date.now()}`
    dialogVisible.value = true
  } catch {
    ElMessage.error('生成扫码上传凭证失败，请稍后重试')
  }
}

const onFinished = () => {
  dialogVisible.value = false
  codeUrl.value = ''
  emit('on-success', '')
}
</script>
