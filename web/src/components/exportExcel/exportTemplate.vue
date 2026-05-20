<template>
  <el-button type="primary" icon="download" @click="exportTemplateFunc"
    >下载模板</el-button
  >
</template>

<script setup>
  import { ElMessage } from 'element-plus'
  import { exportTemplate } from '@/api/exportTemplate'

  const EXPORT_TOKEN_HEADER_FALLBACK = 'X-Export-Token'

  const parseDownloadPayload = (rawData) => {
    if (!rawData) {
      return null
    }
    if (typeof rawData === 'string') {
      return {
        legacyUrl: rawData,
        url: rawData,
        token: '',
        tokenHeader: EXPORT_TOKEN_HEADER_FALLBACK
      }
    }
    return {
      url: rawData.url || '',
      token: rawData.token || '',
      tokenHeader: rawData.tokenHeader || EXPORT_TOKEN_HEADER_FALLBACK,
      legacyUrl: rawData.legacyUrl || ''
    }
  }

  const resolveBaseUrl = () => {
    const baseUrl = import.meta.env.VITE_BASE_API
    if (baseUrl === '/') {
      return ''
    }
    return baseUrl || ''
  }

  const joinUrl = (baseUrl, urlPath) => {
    if (/^https?:\/\//i.test(urlPath)) {
      return urlPath
    }
    if (!baseUrl) {
      return urlPath
    }
    if (baseUrl.endsWith('/') && urlPath.startsWith('/')) {
      return `${baseUrl.slice(0, -1)}${urlPath}`
    }
    if (!baseUrl.endsWith('/') && !urlPath.startsWith('/')) {
      return `${baseUrl}/${urlPath}`
    }
    return `${baseUrl}${urlPath}`
  }

  const parseFilenameFromDisposition = (headerValue, fallback) => {
    if (!headerValue) {
      return fallback
    }
    const utf8Match = headerValue.match(/filename\*=UTF-8''([^;]+)/i)
    if (utf8Match && utf8Match[1]) {
      return decodeURIComponent(utf8Match[1])
    }
    const plainMatch = headerValue.match(/filename="?([^";]+)"?/i)
    if (plainMatch && plainMatch[1]) {
      return plainMatch[1]
    }
    return fallback
  }

  const downloadByTokenPayload = async (payload, baseUrl, fallbackName) => {
    const urlPath = payload.url || payload.legacyUrl
    if (!urlPath) {
      throw new Error('导出地址缺失')
    }

    if (!payload.token) {
      window.open(joinUrl(baseUrl, urlPath), '_blank')
      return
    }

    const headers = {
      Authorization: `ExportToken ${payload.token}`
    }

    const resp = await fetch(joinUrl(baseUrl, urlPath), {
      method: 'GET',
      headers
    })

    if (!resp.ok) {
      throw new Error(`下载失败(${resp.status})`)
    }

    const contentType = (resp.headers.get('content-type') || '').toLowerCase()
    if (contentType.includes('application/json')) {
      const jsonBody = await resp.json().catch(() => ({}))
      throw new Error(jsonBody.msg || '下载失败')
    }

    const blob = await resp.blob()
    const disposition = resp.headers.get('content-disposition') || ''
    const fileName = parseFilenameFromDisposition(disposition, fallbackName)
    const objectUrl = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = objectUrl
    link.download = fileName
    link.style.display = 'none'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(objectUrl)
  }

  const props = defineProps({
    templateId: {
      type: String,
      required: true
    }
  })


  const exportTemplateFunc = async () => {
    if (props.templateId === '') {
      ElMessage.error('组件未设置模板ID')
      return
    }
    const baseUrl = resolveBaseUrl()

    const res = await exportTemplate({
      templateID: props.templateId
    })

    if (res.code === 0) {
      ElMessage.success('创建导出任务成功，开始下载')
      const downloadPayload = parseDownloadPayload(res.data)
      if (!downloadPayload) {
        ElMessage.error('导出地址无效')
        return
      }
      try {
        await downloadByTokenPayload(downloadPayload, baseUrl, `template_${Date.now()}.xlsx`)
      } catch (error) {
        ElMessage.error(error?.message || '下载失败')
      }
    }

  }
</script>
