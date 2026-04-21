<template>
  <div class="cs-workbench flex h-full overflow-hidden">
    <!-- 左侧：会话列表 -->
    <div class="session-list flex-shrink-0 w-72 flex flex-col border-r border-gray-200 bg-white">
      <div class="p-3 border-b border-gray-100">
        <div class="flex items-center justify-between mb-2">
          <span class="font-semibold text-gray-700">会话列表</span>
          <el-tag :type="wsStatus === 'connected' ? 'success' : 'danger'" size="small">
            {{ wsStatus === 'connected' ? '在线' : '离线' }}
          </el-tag>
        </div>
        <el-input v-model="searchKeyword" placeholder="搜索会话..." size="small" clearable prefix-icon="Search" />
      </div>
      <!-- 标签过滤 -->
      <div class="flex gap-1 px-3 py-2 border-b border-gray-100">
        <el-button
          v-for="tab in tabs"
          :key="tab.value"
          :type="activeTab === tab.value ? 'primary' : ''"
          size="small"
          plain
          @click="activeTab = tab.value"
        >{{ tab.label }}</el-button>
      </div>
      <!-- 列表 -->
      <div class="flex-1 overflow-y-auto">
        <div
          v-for="conv in filteredConversations"
          :key="conv.ID"
          class="flex items-start gap-3 px-3 py-3 cursor-pointer hover:bg-gray-50 border-b border-gray-50 transition-colors"
          :class="{ 'bg-blue-50 border-l-4 border-l-blue-500': currentConvId === conv.ID }"
          @click="selectConversation(conv)"
        >
          <el-avatar :size="36" class="flex-shrink-0">
            {{ getUserLabel(conv).charAt(0).toUpperCase() }}
          </el-avatar>
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium text-gray-800 truncate">{{ getUserLabel(conv) }}</span>
              <span class="text-xs text-gray-400 flex-shrink-0 ml-1">{{ formatTime(conv.UpdatedAt) }}</span>
            </div>
            <div class="flex items-center justify-between mt-0.5">
              <p class="text-xs text-gray-500 truncate">{{ conv.lastMsg || '暂无消息' }}</p>
              <el-tag v-if="conv.status === 'pending'" type="warning" size="small" class="flex-shrink-0 ml-1">排队</el-tag>
              <el-tag v-else-if="conv.status === 'active'" type="success" size="small" class="flex-shrink-0 ml-1">进行中</el-tag>
            </div>
          </div>
        </div>
        <el-empty v-if="filteredConversations.length === 0" description="暂无会话" :image-size="60" class="mt-8" />
      </div>
    </div>

    <!-- 中间：聊天区域 -->
    <div class="chat-area flex-1 flex flex-col min-w-0 bg-gray-50">
      <template v-if="currentConvId">
        <!-- 顶部信息栏 -->
        <div class="chat-header flex items-center justify-between px-4 py-2.5 bg-white border-b border-gray-200">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ getUserLabel(currentConv) }}</span>
            <el-tag :type="currentConv?.status === 'active' ? 'success' : 'warning'" size="small">
              {{ currentConv?.status === 'active' ? '服务中' : '排队中' }}
            </el-tag>
          </div>
          <div class="flex gap-2">
            <el-button size="small" @click="handleTransfer">转接</el-button>
            <el-button size="small" type="danger" @click="handleClose">结束会话</el-button>
          </div>
        </div>

        <!-- 消息区 -->
        <div ref="msgContainer" class="flex-1 overflow-y-auto px-4 py-3 space-y-3">
          <div v-if="loadingHistory" class="text-center text-gray-400 text-sm py-4">加载中...</div>
          <template v-else>
            <div
              v-for="msg in messages"
              :key="msg.ID"
              class="flex"
              :class="msg.senderType === 'agent' ? 'justify-end' : 'justify-start'"
            >
              <!-- 系统消息 -->
              <div v-if="msg.senderType === 'system'" class="mx-auto text-xs text-gray-400 bg-gray-100 px-3 py-1 rounded-full">
                {{ msg.content }}
              </div>
              <!-- 普通消息 -->
              <template v-else>
                <el-avatar v-if="msg.senderType !== 'agent'" :size="32" class="flex-shrink-0 mr-2 mt-0.5">U</el-avatar>
                <div class="max-w-[70%]">
                  <div
                    class="px-3 py-2 rounded-2xl text-sm break-words"
                    :class="msg.senderType === 'agent'
                      ? 'bg-blue-500 text-white rounded-tr-sm'
                      : 'bg-white text-gray-800 shadow-sm rounded-tl-sm'"
                  >
                    <span v-if="msg.revoked" class="italic text-opacity-70">[已撤回]</span>
                    <span v-else>{{ msg.content }}</span>
                  </div>
                  <div class="flex items-center gap-2 mt-1" :class="msg.senderType === 'agent' ? 'justify-end' : 'justify-start'">
                    <span class="text-xs text-gray-400">{{ formatTime(msg.CreatedAt) }}</span>
                    <el-button
                      v-if="msg.senderType === 'agent' && !msg.revoked && canRevoke(msg)"
                      text
                      size="small"
                      class="text-xs text-gray-400 p-0"
                      @click="handleRevoke(msg)"
                    >撤回</el-button>
                  </div>
                </div>
                <el-avatar v-if="msg.senderType === 'agent'" :size="32" class="flex-shrink-0 ml-2 mt-0.5">
                  {{ agentName.charAt(0) }}
                </el-avatar>
              </template>
            </div>
          </template>
          <div ref="msgBottom" />
        </div>

        <!-- 快捷回复 -->
        <div v-if="quickReplies.length" class="flex gap-2 px-4 py-2 bg-white border-t border-gray-100 overflow-x-auto flex-shrink-0">
          <el-button
            v-for="qr in quickReplies"
            :key="qr.ID"
            size="small"
            plain
            class="flex-shrink-0"
            @click="inputText = qr.Content"
          >{{ qr.Title }}</el-button>
        </div>

        <!-- 输入区 -->
        <div class="input-area bg-white border-t border-gray-200 p-3">
          <el-input
            v-model="inputText"
            type="textarea"
            :rows="3"
            placeholder="输入消息，Ctrl+Enter 发送..."
            resize="none"
            @keydown.ctrl.enter.prevent="handleSend"
          />
          <div class="flex justify-end mt-2 gap-2">
            <el-button size="small" @click="inputText = ''">清空</el-button>
            <el-button size="small" type="primary" :disabled="!inputText.trim()" @click="handleSend">
              发送 (Ctrl+Enter)
            </el-button>
          </div>
        </div>
      </template>

      <!-- 未选中会话 -->
      <div v-else class="flex-1 flex items-center justify-center">
        <el-empty description="请选择一个会话开始服务" />
      </div>
    </div>

    <!-- 右侧：用户信息 -->
    <div v-if="currentConv" class="user-panel flex-shrink-0 w-60 border-l border-gray-200 bg-white overflow-y-auto">
      <div class="p-4 border-b border-gray-100">
        <h3 class="text-sm font-semibold text-gray-700 mb-3">用户信息</h3>
        <div class="flex flex-col items-center gap-2">
          <el-avatar :size="56">{{ getUserLabel(currentConv).charAt(0).toUpperCase() }}</el-avatar>
          <span class="text-sm font-medium">{{ getUserLabel(currentConv) }}</span>
          <span class="text-xs text-gray-400">UID: {{ currentConv.clientUserId }}</span>
        </div>
      </div>
      <div class="p-4">
        <h3 class="text-sm font-semibold text-gray-700 mb-2">会话信息</h3>
        <div class="space-y-2 text-xs text-gray-600">
          <div class="flex justify-between">
            <span class="text-gray-400">会话ID</span>
            <span>{{ currentConv.ID }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">状态</span>
            <el-tag :type="currentConv.status === 'active' ? 'success' : 'warning'" size="small">
              {{ statusLabel(currentConv.status) }}
            </el-tag>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">创建时间</span>
            <span>{{ formatTime(currentConv.CreatedAt) }}</span>
          </div>
          <div v-if="currentConv.rating" class="flex justify-between">
            <span class="text-gray-400">用户评分</span>
            <el-rate :model-value="currentConv.rating" disabled size="small" />
          </div>
        </div>
        <div class="mt-4">
          <el-button
            size="small"
            type="danger"
            plain
            class="w-full"
            @click="handleAddBlacklist"
          >加入黑名单</el-button>
        </div>
      </div>
    </div>

    <!-- 转接对话框 -->
    <el-dialog v-model="transferVisible" title="转接会话" width="360px">
      <el-select v-model="transferTargetId" placeholder="选择目标坐席" class="w-full">
        <el-option
          v-for="agent in onlineAgents"
          :key="agent.userId"
          :label="`${agent.nickname || agent.sysNickname} (活跃:${agent.activeSessions})`"
          :value="agent.userId"
        />
      </el-select>
      <template #footer>
        <el-button @click="transferVisible = false">取消</el-button>
        <el-button type="primary" :disabled="!transferTargetId" @click="confirmTransfer">确认转接</el-button>
      </template>
    </el-dialog>

    <!-- 黑名单对话框 -->
    <el-dialog v-model="blacklistVisible" title="加入黑名单" width="360px">
      <el-input v-model="blacklistReason" placeholder="请填写原因（可选）" />
      <template #footer>
        <el-button @click="blacklistVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmBlacklist">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import {
  getAgentConversationList,
  getMessageHistory,
  closeConversation,
  transferConversation,
  getAllQuickReplies,
  getAgentList,
  addToBlacklist
} from '@/api/customerService'
import { useUserStore } from '@/pinia/modules/user'

const userStore = useUserStore()

// ===== 状态 =====
const wsStatus = ref('disconnected')
const conversations = ref([])
const currentConvId = ref(null)
const messages = ref([])
const inputText = ref('')
const searchKeyword = ref('')
const activeTab = ref('all')
const quickReplies = ref([])
const agentList = ref([])
const loadingHistory = ref(false)
const msgContainer = ref(null)
const msgBottom = ref(null)
const transferVisible = ref(false)
const transferTargetId = ref(null)
const blacklistVisible = ref(false)
const blacklistReason = ref('')
const agentName = computed(() => userStore.userInfo?.nickName || '坐席')

let ws = null

const tabs = [
  { label: '全部', value: 'all' },
  { label: '进行中', value: 'active' },
  { label: '排队', value: 'pending' },
]

// ===== 计算属性 =====
const currentConv = computed(() => conversations.value.find(c => c.ID === currentConvId.value))

const onlineAgents = computed(() =>
  agentList.value.filter(a => a.onlineStatus && a.userId !== userStore.userInfo?.ID)
)

const filteredConversations = computed(() => {
  let list = conversations.value
  if (activeTab.value !== 'all') {
    list = list.filter(c => c.status === activeTab.value)
  }
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    list = list.filter(c => getUserLabel(c).toLowerCase().includes(kw))
  }
  return list
})

// ===== WebSocket =====
function connectWs() {
  const token = userStore.token
  const wsBase = import.meta.env.VITE_WS_URL || window.location.origin.replace(/^http/, 'ws')
  const apiBase = import.meta.env.VITE_BASE_API || ''
  ws = new WebSocket(`${wsBase}${apiBase}/cs/wsAgent?token=${token}`)

  ws.onopen = () => {
    wsStatus.value = 'connected'
  }

  ws.onclose = () => {
    wsStatus.value = 'disconnected'
    // 5秒后自动重连
    setTimeout(connectWs, 5000)
  }

  ws.onerror = () => {
    wsStatus.value = 'disconnected'
  }

  ws.onmessage = (e) => {
    try {
      const frame = JSON.parse(e.data)
      handleWsFrame(frame)
    } catch {/* ignore */}
  }
}

function handleWsFrame(frame) {
  switch (frame.event) {
    case 'message': {
      const msg = frame.data
      if (msg.conversationId === currentConvId.value) {
        messages.value.push(msg)
        scrollToBottom()
      }
      // 更新会话列表的最后消息
      const conv = conversations.value.find(c => c.ID === msg.conversationId)
      if (conv) conv.lastMsg = msg.content
      break
    }
    case 'new_conv': {
      // 新会话分配过来，刷新列表
      loadConversations()
      break
    }
    case 'revoke': {
      const { messageId } = frame.data || frame
      const idx = messages.value.findIndex(m => m.ID === messageId)
      if (idx !== -1) {
        messages.value[idx].Revoked = true
        messages.value[idx].Content = '[已撤回]'
      }
      break
    }
    case 'closed': {
      const convId = frame.data?.conversationId || frame.conversationId
      const conv = conversations.value.find(c => c.ID === convId)
      if (conv) conv.status = 'closed'
      break
    }
    case 'transfer': {
      loadConversations()
      break
    }
    case 'pong':
      break
  }
}

function sendWs(payload) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(payload))
  }
}

// ===== 方法 =====
async function loadConversations() {
  try {
    const res = await getAgentConversationList({ page: 1, pageSize: 50 })
    if (res.code === 0) {
      conversations.value = res.data?.list || []
    }
  } catch { /* ignore */ }
}

async function selectConversation(conv) {
  currentConvId.value = conv.ID
  messages.value = []
  loadingHistory.value = true
  try {
    const res = await getMessageHistory({ conversationId: conv.ID, page: 1, pageSize: 50 })
    if (res.code === 0) {
      messages.value = (res.data?.list || []).reverse()
    }
  } finally {
    loadingHistory.value = false
    await nextTick()
    scrollToBottom()
  }
}

function handleSend() {
  const content = inputText.value.trim()
  if (!content || !currentConvId.value) return
  sendWs({
    event: 'send_message',
    msgType: 'text',
    content,
    messageId: currentConvId.value, // agent 用此字段表示 conversationId
    clientMsgId: `agent_${Date.now()}`,
  })
  inputText.value = ''
}

function handleRevoke(msg) {
  sendWs({ event: 'revoke', messageId: msg.ID })
}

function canRevoke(msg) {
  return dayjs().diff(dayjs(msg.CreatedAt), 'second') < 120
}

async function handleClose() {
  await ElMessageBox.confirm('确认结束当前会话？', '提示', { type: 'warning' })
  await closeConversation({ conversationId: currentConvId.value })
  const conv = currentConv.value
  if (conv) conv.Status = 'closed'
  ElMessage.success('已结束会话')
}

async function handleTransfer() {
  await loadAgentList()
  transferTargetId.value = null
  transferVisible.value = true
}

async function confirmTransfer() {
  await transferConversation({ conversationId: currentConvId.value, targetAgentUserID: transferTargetId.value })
  transferVisible.value = false
  ElMessage.success('转接成功')
  await loadConversations()
}

function handleAddBlacklist() {
  blacklistReason.value = ''
  blacklistVisible.value = true
}

async function confirmBlacklist() {
  await addToBlacklist({
    clientUserId: currentConv.value?.clientUserId,
    reason: blacklistReason.value,
  })
  blacklistVisible.value = false
  ElMessage.success('已加入黑名单')
}

async function loadAgentList() {
  const res = await getAgentList()
  if (res.code === 0) agentList.value = res.data?.list || []
}

async function loadQuickReplies() {
  const res = await getAllQuickReplies()
  if (res.code === 0) quickReplies.value = res.data || []
}

function getUserLabel(conv) {
  if (!conv) return '未知用户'
  return conv.clientNickname || `用户${conv.clientUserId}`
}

function statusLabel(status) {
  const map = { pending: '排队中', active: '服务中', closed: '已结束' }
  return map[status] || status
}

function formatTime(t) {
  if (!t) return ''
  const d = dayjs(t)
  if (d.isToday()) return d.format('HH:mm')
  return d.format('MM-DD HH:mm')
}

function scrollToBottom() {
  nextTick(() => {
    msgBottom.value?.scrollIntoView({ behavior: 'smooth' })
  })
}

watch(currentConvId, () => {
  scrollToBottom()
})

// ===== 生命周期 =====
onMounted(async () => {
  await Promise.all([loadConversations(), loadQuickReplies()])
  connectWs()
})

onUnmounted(() => {
  if (ws) {
    ws.onclose = null
    ws.close()
  }
})
</script>

<style scoped>
.cs-workbench {
  height: calc(100vh - 84px);
}
</style>
