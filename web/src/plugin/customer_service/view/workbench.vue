<template>
  <div class="cs-workbench flex min-h-0 overflow-y-auto rounded-2xl border border-slate-200 shadow-sm">
    <!-- 左侧：会话列表 -->
    <div class="session-list cs-panel flex-shrink-0 w-80 min-h-0 flex flex-col border-r border-slate-200 bg-white/90">
      <div class="p-3 border-b border-gray-100">
        <div class="flex items-center justify-between mb-2">
          <span class="font-semibold text-gray-700">会话列表</span>
          <el-tag :type="wsStatus === 'connected' ? 'success' : 'danger'" size="small">
            {{ wsStatus === 'connected' ? '在线' : '离线' }}
          </el-tag>
        </div>
        <el-input v-model="searchKeyword" placeholder="搜索会话..." size="small" clearable prefix-icon="Search" />
        <div class="mt-3 grid grid-cols-2 gap-2">
          <div class="rounded-lg border border-amber-100 bg-amber-50 px-2 py-1.5">
            <div class="text-xs text-amber-600">排队中</div>
            <div class="text-base font-semibold text-amber-700">{{ pendingCount }}</div>
          </div>
          <div class="rounded-lg border border-emerald-100 bg-emerald-50 px-2 py-1.5">
            <div class="text-xs text-emerald-600">服务中</div>
            <div class="text-base font-semibold text-emerald-700">{{ activeCount }}</div>
          </div>
        </div>
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
          class="session-item flex items-start gap-3 px-3 py-3 cursor-pointer hover:bg-gray-50 border-b border-gray-50 transition-colors"
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
              <div class="flex items-center gap-1 flex-shrink-0 ml-1">
                <el-tag v-if="conv.status === 'pending'" type="warning" size="small">
                  排队{{ conv.queuePosition > 0 ? ' #' + conv.queuePosition : '' }}
                </el-tag>
                <el-tag v-else-if="conv.status === 'active'" type="success" size="small">进行中</el-tag>
                <el-badge
                  v-if="(conv.unreadCount || 0) > 0"
                  :value="conv.unreadCount"
                  :max="99"
                  type="danger"
                />
              </div>
            </div>
            <div v-if="conv.status === 'pending'" class="mt-1">
              <el-button size="small" type="warning" plain @click.stop="handleAcceptConversation(conv)">立即接入</el-button>
            </div>
          </div>
        </div>
        <el-empty v-if="filteredConversations.length === 0" description="暂无会话" :image-size="60" class="mt-8" />
      </div>
    </div>

    <!-- 中间：聊天区域 -->
    <div class="chat-area cs-panel flex-1 min-h-0 flex flex-col min-w-0 overflow-hidden bg-slate-50/70">
      <template v-if="currentConvId">
        <!-- 顶部信息栏 -->
        <div class="chat-header chat-header-panel sticky top-0 z-20 flex items-center justify-between px-4 py-3 bg-white/95 border-b border-slate-200">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ getUserLabel(currentConv) }}</span>
            <el-tag :type="currentConv?.status === 'active' ? 'success' : 'warning'" size="small">
              {{ currentConv?.status === 'active' ? '服务中' : '排队中' }}
            </el-tag>
            <el-tag v-if="currentConv?.status === 'pending' && currentConv?.queuePosition > 0" type="warning" effect="plain" size="small">
              队列位置 #{{ currentConv.queuePosition }}
            </el-tag>
          </div>
          <div class="flex gap-2">
            <el-button
              v-if="currentConv?.status === 'pending'"
              size="small"
              type="warning"
              :loading="accepting"
              @click="handleAcceptConversation()"
            >接入会话</el-button>
            <el-button size="small" :disabled="currentConv?.status !== 'active'" @click="handleTransfer">转接</el-button>
            <el-button size="small" type="danger" :disabled="currentConv?.status !== 'active'" @click="handleClose">结束会话</el-button>
          </div>
        </div>

        <div v-if="currentConv?.status === 'pending'" class="mx-4 mt-3 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-700">
          当前会话仍在排队。你可以点击上方“接入会话”，或直接回复首条消息自动接入。
        </div>

        <!-- 消息区 -->
        <div ref="msgContainer" class="flex-1 overflow-y-auto px-4 py-3 pb-4 space-y-3">
          <div v-if="loadingHistory" class="text-center text-gray-400 text-sm py-4">加载中...</div>
          <template v-else>
            <div class="text-center text-xs text-gray-400 py-1">
              <el-button
                v-if="historyHasMore"
                size="small"
                text
                :loading="historyLoadingMore"
                @click="loadMoreHistory"
              >加载更多</el-button>
              <span v-else-if="messages.length > 0">已加载完</span>
            </div>
            <div
              v-for="msg in messages"
              :key="msg.ID || msg.clientMsgId || msg.CreatedAt"
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
                    <img
                      v-else-if="msg.msgType === 'image'"
                      :src="msg.content"
                      class="max-w-full rounded-lg cursor-pointer"
                      style="max-height:240px"
                      @click="previewImage(msg.content)"
                    />
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
            @click="inputText = qr.content"
          >{{ qr.title }}</el-button>
        </div>

        <!-- 输入区 -->
        <div class="input-area sticky bottom-0 z-20 bg-white/95 border-t border-gray-200 p-3">
          <el-input
            ref="inputRef"
            v-model="inputText"
            type="textarea"
            :rows="3"
            :placeholder="currentConv?.status === 'pending' ? '输入首条回复后将自动接入，Ctrl+Enter 发送...' : '输入消息，Ctrl+Enter 发送...'"
            resize="none"
            @keydown.ctrl.enter.prevent="handleSend"
          />
          <div class="flex justify-between mt-2 gap-2">
            <div class="flex gap-1">
              <el-popover placement="top-start" :width="320" trigger="click">
                <template #reference>
                  <el-button size="small" title="插入表情">😊</el-button>
                </template>
                <div class="emoji-grid">
                  <button
                    v-for="emoji in emojiList"
                    :key="emoji"
                    type="button"
                    class="emoji-item"
                    @click="insertEmoji(emoji)"
                  >{{ emoji }}</button>
                </div>
              </el-popover>
              <el-button size="small" title="发送图片" @click="triggerImageUpload">
                <el-icon><Picture /></el-icon>
              </el-button>
              <input ref="imageInput" type="file" accept="image/*" class="hidden" @change="handleImageUpload" />
            </div>
            <div class="flex gap-2">
              <el-button size="small" @click="inputText = ''">清空</el-button>
              <el-button size="small" type="primary" :disabled="!inputText.trim() || currentConv?.status === 'closed'" @click="handleSend">
                发送 (Ctrl+Enter)
              </el-button>
            </div>
          </div>
        </div>
      </template>

      <!-- 未选中会话 -->
      <div v-else class="flex-1 flex items-center justify-center">
        <el-empty description="请选择一个会话开始服务" />
      </div>
    </div>

    <!-- 右侧：用户信息 -->
    <div v-if="currentConv" class="user-panel cs-panel flex-shrink-0 w-60 min-h-0 border-l border-gray-200 bg-white overflow-y-auto">
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
import { Picture } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import {
  getAgentConversationList,
  acceptConversation,
  getCsConfig,
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
const historyLoadingMore = ref(false)
const historyPage = ref(1)
const historyTotal = ref(0)
const inputRef = ref(null)
const msgContainer = ref(null)
const msgBottom = ref(null)
const transferVisible = ref(false)
const transferTargetId = ref(null)
const blacklistVisible = ref(false)
const blacklistReason = ref('')
const accepting = ref(false)
const agentName = computed(() => userStore.userInfo?.nickName || '坐席')
const uploadMaxSizeMB = ref(5)
const uploadAllowMimeTypes = ref(['image/jpeg', 'image/png', 'image/webp', 'image/gif'])
const uploadAllowExtLabel = ref('JPG/JPEG/PNG/WEBP/GIF')

let ws = null
let pingTimer = null
let pollTimer = null
const imageInput = ref(null)
const wsAllowQueryTokenFallback = String(import.meta.env.VITE_WS_ALLOW_QUERY_TOKEN || '').toLowerCase() === 'true'

const tabs = [
  { label: '全部', value: 'all' },
  { label: '进行中', value: 'active' },
  { label: '排队', value: 'pending' },
  { label: '已退出', value: 'closed' },
]

const emojiList = [
  '😀', '😄', '😁', '😆', '😊', '😉', '😍', '😘', '🤝', '🙌',
  '👏', '👍', '💪', '🎉', '✨', '🔥', '❤️', '🧡', '💚', '💙',
  '🤖', '🤔', '😎', '🥳', '😭', '😮', '😴', '😡', '🙏', '🌈'
]

const uploadMimeByExt = {
  jpg: 'image/jpeg',
  jpeg: 'image/jpeg',
  png: 'image/png',
  webp: 'image/webp',
  gif: 'image/gif',
}

// ===== 计算属性 =====
const currentConv = computed(() => conversations.value.find(c => c.ID === currentConvId.value))

const onlineAgents = computed(() =>
  agentList.value.filter(a => a.onlineStatus === 'online' && a.userId !== userStore.userInfo?.ID)
)

const pendingCount = computed(() => conversations.value.filter(c => c.status === 'pending').length)
const activeCount = computed(() => conversations.value.filter(c => c.status === 'active').length)
const historyHasMore = computed(() => messages.value.length < historyTotal.value)

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
function buildWsUrl(path) {
  const baseApi = import.meta.env.VITE_BASE_API || ''
  if (/^https?:\/\//.test(baseApi)) {
    // 生产：VITE_BASE_API 是绝对 URL，如 https://clothapi.235235.vip
    return baseApi.replace(/^http/, 'ws') + path
  }
  // 开发：VITE_BASE_API 是相对路径，如 /api，通过 vite proxy 转发
  const proto = window.location.protocol === 'https:' ? 'wss://' : 'ws://'
  return proto + window.location.host + baseApi + path
}

function connectWs() {
  const token = userStore.token
  const wsURL = buildWsUrl('/cs/wsAgent')
  try {
    ws = token ? new WebSocket(wsURL, ['bearer', token]) : new WebSocket(wsURL)
  } catch (err) {
    if (!wsAllowQueryTokenFallback) {
      wsStatus.value = 'disconnected'
      console.error('[CS WS] protocol handshake failed and query-token fallback is disabled', err)
      return
    }
    ws = new WebSocket(buildWsUrl(`/cs/wsAgent?token=${encodeURIComponent(token || '')}`))
  }

  ws.onopen = () => {
    wsStatus.value = 'connected'
    // 建连成功立即拉取一次会话，避免错过早期事件导致排队列表滞后
    loadConversations()
    // 应用层心跳：每30秒发一次ping，防止代理/nginx超时断开
    pingTimer = setInterval(() => {
      sendWs({ event: 'ping' })
    }, 30000)
  }

  ws.onclose = () => {
    wsStatus.value = 'disconnected'
    clearInterval(pingTimer)
    pingTimer = null
    // 5秒后自动重连，重连后刷新会话列表
    setTimeout(() => {
      connectWs()
      loadConversations()
    }, 5000)
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
    case 'pending_conv': {
      // 有用户处于排队中（分配失败或无坐席），刷新会话列表
      loadConversations()
      break
    }
    case 'convList': {
      // 坐席连接后服务端推送当前进行中会话，直接刷新完整列表
      loadConversations()
      break
    }
    case 'message': {
      const msg = frame.data
      if (!msg) break
      const conv = conversations.value.find(c => c.ID === msg.conversationId)
      if (conv) {
        conv.lastMsg = summarizeMessage(msg)
        conv.UpdatedAt = msg.CreatedAt || conv.UpdatedAt
      }
      if (msg.conversationId === currentConvId.value) {
        // 坐席自己发的消息：替换乐观更新的临时消息
        const tempIdx = messages.value.findIndex(m => m._tempId && m.clientMsgId === msg.clientMsgId)
        if (tempIdx !== -1) {
          messages.value[tempIdx] = msg
        } else if (msg.senderType !== 'agent') {
          // 用户发来的消息直接追加
          messages.value.push(msg)
        }
        if (msg.senderType === 'user') {
          markCurrentConversationRead()
        }
        scrollToBottom()
      } else if (msg.senderType === 'user' && conv) {
        conv.unreadCount = Number(conv.unreadCount || 0) + 1
      }
      // 用户发来的消息：标签不聚焦时发送浏览器通知
      if (msg.senderType === 'user' && document.visibilityState !== 'visible') {
        if ('Notification' in window && Notification.permission === 'granted') {
          const conv = conversations.value.find(c => c.ID === msg.conversationId)
          const label = conv ? getUserLabel(conv) : '用户'
          new Notification(`新消息来自 ${label}`, {
            body: msg.msgType === 'image' ? '[图片]' : msg.content,
            icon: '/favicon.ico',
          })
        }
      }
      break
    }
    case 'new_conv': {
      // 新会话分配给当前坐席（可能是 pending → active）
      // 先乐观更新本地列表，再刷新完整列表以确保数据同步
      const newConv = frame.data
      if (newConv && newConv.ID) {
        const idx = conversations.value.findIndex(c => c.ID === newConv.ID)
        if (idx !== -1) {
          conversations.value[idx] = { ...conversations.value[idx], ...newConv }
        } else {
          conversations.value.unshift(newConv)
        }
      }
      // 无论如何都刷新，保证与服务端状态完全一致
      loadConversations()
      break
    }
    case 'revoke': {
      const { messageId } = frame.data || {}
      const idx = messages.value.findIndex(m => m.ID === messageId)
      if (idx !== -1) {
        messages.value[idx] = { ...messages.value[idx], revoked: true, content: '[已撤回]' }
      }
      break
    }
    case 'closed': {
      const convId = frame.data?.conversationId
      const conv = conversations.value.find(c => c.ID === convId)
      if (conv) conv.status = 'closed'
      if (convId === currentConvId.value) {
        messages.value.push({ senderType: 'system', content: '会话已结束', CreatedAt: new Date().toISOString() })
      }
      break
    }
    case 'transfer': {
      loadConversations()
      break
    }
    case 'error': {
      const msg = frame.data?.message || '操作失败'
      ElMessage.warning(msg)
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
    const params = { page: 1, pageSize: 50 }
    if (activeTab.value !== 'all') {
      params.status = activeTab.value
    }
    const res = await getAgentConversationList(params)
    if (res.code === 0) {
      conversations.value = res.data?.list || []
      if (currentConvId.value && !conversations.value.some(c => c.ID === currentConvId.value)) {
        currentConvId.value = null
        messages.value = []
        historyPage.value = 1
        historyTotal.value = 0
      }
    }
  } catch { /* ignore */ }
}

async function selectConversation(conv) {
  currentConvId.value = conv.ID
  historyPage.value = 1
  historyTotal.value = 0
  messages.value = []
  loadingHistory.value = true
  try {
    await fetchHistory(conv.ID, 1, false)
    markCurrentConversationRead()
  } finally {
    loadingHistory.value = false
  }
}

async function fetchHistory(convID, page, prepend) {
  const res = await getMessageHistory({ conversationId: convID, page, pageSize: 10 })
  if (res.code !== 0) return
  const list = (res.data?.list || []).reverse()
  historyTotal.value = Number(res.data?.total || 0)
  if (prepend) {
    const oldHeight = msgContainer.value?.scrollHeight || 0
    messages.value = [...list, ...messages.value]
    await nextTick()
    if (msgContainer.value) {
      const newHeight = msgContainer.value.scrollHeight || 0
      msgContainer.value.scrollTop = newHeight - oldHeight
    }
  } else {
    messages.value = list
    await nextTick()
    scrollToBottom()
  }
  historyPage.value = page
}

async function loadMoreHistory() {
  if (!currentConvId.value || !historyHasMore.value || historyLoadingMore.value) return
  historyLoadingMore.value = true
  try {
    await fetchHistory(currentConvId.value, historyPage.value + 1, true)
  } finally {
    historyLoadingMore.value = false
  }
}

function markCurrentConversationRead() {
  if (!currentConvId.value) return
  const conv = conversations.value.find(c => c.ID === currentConvId.value)
  if (conv) conv.unreadCount = 0
  sendWs({ event: 'read', messageId: currentConvId.value })
}

function summarizeMessage(msg) {
  if (!msg) return '暂无消息'
  if (msg.revoked) return '[已撤回]'
  if (msg.msgType === 'image') return '[图片]'
  if (msg.msgType === 'emoji') return '[表情]'
  return msg.content || '暂无消息'
}

function handleSend() {
  const content = inputText.value.trim()
  if (!content || !currentConvId.value) return
  if (currentConv.value?.status === 'closed') return
  if (!ws || ws.readyState !== WebSocket.OPEN) {
    ElMessage.warning('连接已断开，正在重连...')
    return
  }
  const clientMsgId = `agent_${Date.now()}`
  // 乐观更新：立即显示在聊天框
  const tempMsg = {
    _tempId: true,
    clientMsgId,
    senderType: 'agent',
    content,
    conversationId: currentConvId.value,
    CreatedAt: new Date().toISOString(),
    revoked: false,
  }
  messages.value.push(tempMsg)
  scrollToBottom()
  inputText.value = ''
  sendWs({
    event: 'send_message',
    msgType: 'text',
    content,
    messageId: currentConvId.value, // 后端用 messageId 字段作为 conversationId
    clientMsgId,
  })
  const conv = conversations.value.find(c => c.ID === currentConvId.value)
  if (conv) {
    conv.lastMsg = content
    conv.UpdatedAt = new Date().toISOString()
  }
}

async function handleAcceptConversation(targetConv) {
  const conv = targetConv || currentConv.value
  if (!conv?.ID || conv.status !== 'pending') return
  accepting.value = true
  try {
    const res = await acceptConversation({ conversationId: conv.ID })
    if (res.code === 0) {
      ElMessage.success('接入成功')
      const idx = conversations.value.findIndex(c => c.ID === conv.ID)
      if (idx !== -1) {
        conversations.value[idx] = {
          ...conversations.value[idx],
          ...(res.data || {}),
          status: 'active',
          queuePosition: 0,
        }
      }
      if (currentConvId.value === conv.ID) {
        messages.value.push({ senderType: 'system', content: '你已接入该会话', CreatedAt: new Date().toISOString() })
      }
      await loadConversations()
    }
  } finally {
    accepting.value = false
  }
}

function handleRevoke(msg) {
  sendWs({ event: 'revoke', messageId: msg.ID })
}

function canRevoke(msg) {
  return dayjs().diff(dayjs(msg.CreatedAt), 'second') < 120
}

function previewImage(url) {
  window.open(url, '_blank')
}

function triggerImageUpload() {
  imageInput.value?.click()
}

async function handleImageUpload(e) {
  const file = e.target.files?.[0]
  if (!file || !currentConvId.value) return
  e.target.value = ''

  const maxImageSize = uploadMaxSizeMB.value * 1024 * 1024
  if (!uploadAllowMimeTypes.value.includes(file.type)) {
    ElMessage.error(`仅支持 ${uploadAllowExtLabel.value} 图片`)
    return
  }
  if (file.size > maxImageSize) {
    ElMessage.error(`图片大小不能超过 ${uploadMaxSizeMB.value}MB`)
    return
  }

  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await fetch(
      (import.meta.env.VITE_BASE_API || '') + '/cs/agent/message/upload',
      {
        method: 'POST',
        headers: { 'x-token': userStore.token },
        body: formData,
      }
    )
    const json = await res.json()
    if (json.code !== 0) { ElMessage.error(json.msg || '图片上传失败'); return }
    const imageUrl = json.data.file.url
    const clientMsgId = `agent_img_${Date.now()}`
    const tempMsg = {
      _tempId: true, clientMsgId,
      senderType: 'agent', msgType: 'image',
      content: imageUrl,
      conversationId: currentConvId.value,
      CreatedAt: new Date().toISOString(), revoked: false,
    }
    messages.value.push(tempMsg)
    scrollToBottom()
    sendWs({ event: 'send_message', msgType: 'image', content: imageUrl, messageId: currentConvId.value, clientMsgId })
  } catch {
    ElMessage.error('图片上传失败')
  }
}

async function handleClose() {
  if (!currentConvId.value || currentConv.value?.status !== 'active') return
  await ElMessageBox.confirm('确认结束当前会话？', '提示', { type: 'warning' })
  await closeConversation({ conversationId: currentConvId.value })
  const conv = currentConv.value
  if (conv) conv.status = 'closed'
  ElMessage.success('已结束会话')
}

async function handleTransfer() {
  if (currentConv.value?.status !== 'active') return
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
  if (res.code === 0) agentList.value = Array.isArray(res.data) ? res.data : (res.data?.list || [])
}

async function loadQuickReplies() {
  const res = await getAllQuickReplies()
  if (res.code === 0) quickReplies.value = res.data || []
}

async function loadUploadConfig() {
  try {
    const res = await getCsConfig()
    if (res.code !== 0) return
    uploadMaxSizeMB.value = Math.max(1, Number(res.data?.uploadMaxSizeMB || 5))

    const extList = String(res.data?.uploadAllowExt || '')
      .split(',')
      .map(ext => ext.trim().toLowerCase().replace(/^\./, ''))
      .filter(Boolean)
    const uniqueExtList = extList.length ? Array.from(new Set(extList)) : ['jpg', 'jpeg', 'png', 'webp', 'gif']
    uploadAllowExtLabel.value = uniqueExtList.map(ext => ext.toUpperCase()).join('/')

    const mimeList = uniqueExtList
      .map(ext => uploadMimeByExt[ext])
      .filter(Boolean)
    uploadAllowMimeTypes.value = mimeList.length ? Array.from(new Set(mimeList)) : ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  } catch {
    // 使用默认上传策略
  }
}

function getUserLabel(conv) {
  if (!conv) return '未知用户'
  return conv.clientNickname || `用户${conv.clientUserId}`
}

function statusLabel(status) {
  const map = { pending: '排队中', active: '服务中', closed: '已结束' }
  return map[status] || status
}

function insertEmoji(emoji) {
  if (currentConv.value?.status === 'closed') {
    ElMessage.warning('会话已结束，无法发送消息')
    return
  }
  inputText.value = `${inputText.value || ''}${emoji}`
  nextTick(() => {
    inputRef.value?.focus?.()
  })
}

function formatTime(t) {
  if (!t) return ''
  const d = dayjs(t)
  const isToday = d.format('YYYY-MM-DD') === dayjs().format('YYYY-MM-DD')
  if (isToday) return d.format('HH:mm')
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

watch(activeTab, () => {
  currentConvId.value = null
  messages.value = []
  historyPage.value = 1
  historyTotal.value = 0
  loadConversations()
})

// ===== 生命周期 =====
onMounted(async () => {
  await Promise.all([loadConversations(), loadQuickReplies(), loadUploadConfig()])
  connectWs()
  // 安全网：每10秒轮询一次，防止错过WS事件导致排队状态不可见
  pollTimer = setInterval(loadConversations, 10000)
  // 请求浏览器通知权限
  if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission()
  }
})

onUnmounted(() => {
  clearInterval(pingTimer)
  clearInterval(pollTimer)
  if (ws) {
    ws.onclose = null
    ws.close()
  }
})
</script>

<style scoped>
.cs-workbench {
  height: calc(100vh - 84px);
  overflow-y: auto;
  overflow-x: hidden;
  background:
    radial-gradient(circle at 12% 10%, rgba(59, 130, 246, 0.14), transparent 36%),
    radial-gradient(circle at 87% 80%, rgba(14, 165, 233, 0.1), transparent 34%),
    linear-gradient(180deg, #f8fafc, #f1f5f9);
}

.cs-panel {
  backdrop-filter: blur(8px);
}

@media (max-height: 860px) {
  .cs-workbench {
    height: auto;
    min-height: calc(100vh - 84px);
  }
}

.session-item {
  position: relative;
}

.session-item::after {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: transparent;
  transition: background-color 0.2s ease;
}

.session-item:hover::after {
  background: rgba(37, 99, 235, 0.35);
}

.chat-header-panel {
  backdrop-filter: blur(12px);
}

.input-area {
  backdrop-filter: blur(12px);
  box-shadow: 0 -8px 24px rgba(15, 23, 42, 0.06);
}

.emoji-grid {
  display: grid;
  grid-template-columns: repeat(10, minmax(0, 1fr));
  gap: 6px;
}

.emoji-item {
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #fff;
  font-size: 18px;
  line-height: 1;
  height: 34px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.emoji-item:hover {
  transform: translateY(-1px);
  background: #f8fafc;
  border-color: #cbd5e1;
}

:deep(.input-area .el-textarea__inner) {
  border-radius: 12px;
  border-color: #dbe3ef;
  background: #f8fafc;
}

:deep(.input-area .el-textarea__inner:focus) {
  border-color: #60a5fa;
}
</style>
