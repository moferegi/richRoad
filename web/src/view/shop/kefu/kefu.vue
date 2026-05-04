<template>
  <div class="shop-kefu-page">
    <el-card class="mb-4" shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold text-gray-700">客服展示设置</span>
          <el-tag type="info" size="small">双开关并行</el-tag>
        </div>
      </template>

      <div class="grid gap-4 md:grid-cols-2">
        <div class="setting-item">
          <div>
            <div class="setting-title">外部客服展示（shop/kefu）</div>
            <div class="setting-desc">控制用户端是否展示外部客服列表</div>
          </div>
          <el-switch
            :model-value="externalEnabled"
            :loading="savingExternal || loading"
            active-text="开"
            inactive-text="关"
            @change="handleToggleExternal"
          />
        </div>

        <div class="setting-item">
          <div>
            <div class="setting-title">内置客服展示（customerService/config）</div>
            <div class="setting-desc">当前状态只读，可前往客服系统配置页调整</div>
          </div>
          <div class="flex items-center gap-2">
            <el-tag :type="platEnabled ? 'success' : 'info'" size="small">
              {{ platEnabled ? '已开启' : '已关闭' }}
            </el-tag>
            <el-button link type="primary" @click="openCsConfig">去设置</el-button>
          </div>
        </div>
      </div>

      <el-alert
        class="mt-4"
        type="success"
        :closable="false"
        title="展示规则：两个开关互不冲突；若都开启，用户端会先展示外部客服，内置客服入口排在最后。"
      />
    </el-card>

    <KefuService />
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'
import KefuService from '@/view/shop/kefuService/kefuService.vue'
import { getCsConfig } from '@/api/customerService'
import { getSysConfigList, updateSysConfig } from '@/api/client/sysConfig'

defineOptions({
  name: 'KefuRedirect'
})

const loading = ref(false)
const savingExternal = ref(false)
const platEnabled = ref(false)
const externalEnabled = ref(true)
const externalConfigId = ref(0)
const externalRemark = ref('')

async function loadShopKefuSwitch() {
  const res = await getSysConfigList({ page: 1, pageSize: 20, configKey: 'shop_kefu_enabled' })
  if (res.code !== 0) return
  const list = res.data?.list || []
  const row = list.find(item => item.configKey === 'shop_kefu_enabled')
  if (!row) {
    ElMessage.warning('未找到 shop_kefu_enabled 配置，请重启后端初始化默认配置')
    return
  }
  externalConfigId.value = row.ID
  externalRemark.value = row.remark || ''
  externalEnabled.value = row.configValue === 'true'
}

async function loadPlatConfig() {
  const res = await getCsConfig()
  if (res.code === 0) {
    platEnabled.value = !!res.data?.platEnabled
  }
}

async function loadConfig() {
  loading.value = true
  try {
    await Promise.all([loadShopKefuSwitch(), loadPlatConfig()])
  } finally {
    loading.value = false
  }
}

async function handleToggleExternal(val) {
  if (!externalConfigId.value) {
    externalEnabled.value = !val
    ElMessage.warning('配置项不存在，无法更新，请重启后端后重试')
    return
  }
  savingExternal.value = true
  try {
    const res = await updateSysConfig({
      id: externalConfigId.value,
      configValue: val ? 'true' : 'false',
      remark: externalRemark.value,
    })
    if (res.code !== 0) {
      externalEnabled.value = !val
      return
    }
    externalEnabled.value = val
    ElMessage.success('外部客服展示开关已更新')
  } finally {
    savingExternal.value = false
  }
}

function openCsConfig() {
  router.push({ name: 'csConfig' })
}

onMounted(loadConfig)
</script>

<style scoped>
.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border: 1px solid var(--el-border-color-light);
  border-radius: 12px;
  padding: 14px 16px;
  background: linear-gradient(180deg, #ffffff, #f8fafc);
}

.setting-title {
  color: var(--el-text-color-primary);
  font-weight: 600;
  font-size: 14px;
  line-height: 20px;
}

.setting-desc {
  color: var(--el-text-color-secondary);
  font-size: 12px;
  line-height: 18px;
}
</style>
