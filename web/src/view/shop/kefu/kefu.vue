<template>
  <div class="p-6">
    <!-- 平台内置客服开关 -->
    <el-card class="mb-6 shadow-sm">
      <div class="flex items-center justify-between">
        <div>
          <div class="flex items-center gap-3 mb-1">
            <el-icon :size="22" color="#409eff"><Service /></el-icon>
            <span class="text-base font-semibold text-gray-800">平台内置客服系统</span>
            <el-tag v-if="platEnabled" type="success" size="small" effect="plain">已启用</el-tag>
            <el-tag v-else type="info" size="small" effect="plain">未启用</el-tag>
          </div>
          <p class="text-sm text-gray-500 ml-8">
            开启后，用户在 App 点击"联系客服"将直接进入平台内置的实时聊天；
            关闭则展示下方外部客服列表。
          </p>
        </div>
        <el-switch
          v-model="platEnabled"
          :loading="configLoading"
          size="large"
          active-text="内置客服"
          inactive-text="外部客服"
          @change="onPlatToggle"
        />
      </div>

      <transition name="el-fade-in">
        <div v-if="platEnabled" class="mt-4 pt-4 border-t border-gray-100 grid grid-cols-3 gap-3">
          <el-button
            v-for="link in platformLinks"
            :key="link.label"
            :icon="link.icon"
            plain
            class="!h-12 text-sm"
            @click="$router.push(link.to)"
          >{{ link.label }}</el-button>
        </div>
      </transition>
    </el-card>

    <!-- 外部客服管理 -->
    <el-card class="shadow-sm" :class="platEnabled ? 'opacity-50 pointer-events-none' : ''">
      <template #header>
        <div class="flex items-center gap-2">
          <span class="font-semibold">外部客服列表</span>
          <el-tooltip content="关闭内置客服后，App 将展示此列表中的客服账号">
            <el-icon color="#aaa"><QuestionFilled /></el-icon>
          </el-tooltip>
          <el-tag v-if="platEnabled" type="warning" size="small" class="ml-auto">内置客服已开启，此列表暂不生效</el-tag>
        </div>
      </template>
      <KefuService />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, markRaw } from 'vue'
import { ElMessage } from 'element-plus'
import { Service, QuestionFilled, ChatLineRound, User, Lightning, Clock, Document } from '@element-plus/icons-vue'
import KefuService from '@/view/shop/kefuService/kefuService.vue'
import { getCsConfig, updateCsConfig } from '@/api/customerService'

defineOptions({ name: 'KefuConfig' })

const platEnabled = ref(false)
const configLoading = ref(false)

const platformLinks = [
  { label: '坐席工作台', icon: markRaw(ChatLineRound), to: '/layout/superAdmin/csWorkbench' },
  { label: '坐席管理', icon: markRaw(User), to: '/layout/superAdmin/csAgents' },
  { label: '快捷回复', icon: markRaw(Lightning), to: '/layout/superAdmin/csQuickReplies' },
  { label: '历史会话', icon: markRaw(Clock), to: '/layout/superAdmin/csHistory' },
  { label: '黑名单', icon: markRaw(Document), to: '/layout/superAdmin/csBlacklist' },
]

async function loadConfig() {
  configLoading.value = true
  try {
    const res = await getCsConfig()
    if (res.code === 0) platEnabled.value = !!res.data?.platEnabled
  } finally {
    configLoading.value = false
  }
}

async function onPlatToggle(val) {
  const action = val ? '启用' : '关闭'
  configLoading.value = true
  try {
    const res = await updateCsConfig({ platEnabled: val })
    if (res && res.code === 0) {
      ElMessage.success(`已${action}内置客服系统`)
    } else {
      // 接口失败时回滚开关
      platEnabled.value = !val
    }
  } catch {
    platEnabled.value = !val
  } finally {
    configLoading.value = false
  }
}

onMounted(loadConfig)
</script>
