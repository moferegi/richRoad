<template>
  <!-- 浮动按钮 + 聊天窗口 -->
  <div class="cs-widget">
    <!-- 聊天窗口 -->
    <transition name="slide-up">
      <div v-if="open" class="chat-window">
        <!-- 标题栏 -->
        <div class="chat-header">
          <span>在线客服</span>
          <el-icon class="cursor-pointer" @click="open = false"><Close /></el-icon>
        </div>

        <!-- 状态提示 -->
        <div v-if="convStatus === 'pending'" class="status-bar pending">
          排队等待中，您前面还有 {{ queuePosition }} 人...
        </div>
        <div v-else-if="convStatus === 'active'" class="status-bar active">
          已连接客服，请直接发送消息
        </div>

        <!-- 消息区 -->
        <div ref="msgBox" class="msg-box">
          <div
            v-for="msg in messages"
            :key="msg.ID || msg.clientMsgId"
            class="msg-row"
            :class="msg.SenderType === 'user' ? 'mine' : ''"
          >
            <div v-if="msg.SenderType === 'system'" class="sys-msg">{{ msg.Content }}</div>
            <template v-else>
              <el-avatar v-if="msg.SenderType !== 'user'" :size="28" class="avatar">客</el-avatar>
              <div
                class="bubble"
                :class="msg.SenderType === 'user' ? 'bubble-mine' : 'bubble-other'"
              >
                <span v-if="msg.Revoked" class="italic opacity-60">[已撤回]</span>
                <span v-else>{{ msg.Content }}</span>
              </div>
              <el-avatar v-if="msg.SenderType === 'user'" :size="28" class="avatar">我</el-avatar>
            </template>
          </div>
          <div ref="msgBottom" />
        </div>

        <!-- 评价区（会话结束后显示） -->
        <div v-if="convStatus === 'closed' && !rated" class="rate-bar">
          <span class="text-sm text-gray-600">请为本次服务评分：</span>
          <el-rate v-model="rating" @change="submitRating" />
        </div>

        <!-- 输入区 -->
        <div class="input-area">
          <el-input
            v-model="inputText"
            placeholder="输入消息..."
            :disabled="convStatus === 'closed'"
            @keydown.enter.prevent="sendMsg"
          >
            <template #append>
              <el-button :disabled="convStatus === 'closed'" @click="sendMsg">发送</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </transition>

    <!-- 浮动按钮 -->
    <div class="float-btn" @click="handleOpen">
      <el-badge :is-dot="unread > 0">
        <el-icon :size="28" color="#fff"><Service /></el-icon>
      </el-badge>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Close, Service } from '@element-plus/icons-vue'
import { getOrCreateConversation, getMessageHistory, rateConversation } from '@/api/customerService'
import { useUserStore } from '@/pinia/modules/user'

const userStore = useUserStore()

const open = ref(false)
const messages = ref([])
const inputText = ref('')
const convId = ref(null)
const convStatus = ref('idle') // idle | pending | active | closed
const queuePosition = ref(0)
const unread = ref(0)
const rating = ref(0)
const rated = ref(false)
const msgBox = ref(null)
const msgBottom = ref(null)

let ws = null
let clientMsgCounter = 0

async function handleOpen() {
  open.value = true
  unread.value = 0
  if (convId.value) return
  await initConversation()
}

async function initConversation() {
  try {
    const res = await getOrCreateConversation()
    if (res.code !== 0) return
    const conv = res.data
    convId.value = conv.ID
    convStatus.value = conv.Status
    // 加载历史消息
    const histRes = await getMessageHistory({ conversationId: conv.ID, page: 1, pageSize: 30 })
    if (histRes.code === 0) {
      messages.value = (histRes.data?.list || []).reverse()
    }
    scrollBottom()
    connectWs()
  } catch {
    ElMessage.error('连接客服失败，请稍后重试')
  }
}

function connectWs() {
  if (ws) return
  const token = userStore.token
  const wsBase = import.meta.env.VITE_WS_URL || window.location.origin.replace(/^http/, 'ws')
  const apiBase = import.meta.env.VITE_BASE_API || ''
  ws = new WebSocket(`${wsBase}${apiBase}/cs/ws?token=${token}`)

  ws.onopen = () => {}

  ws.onclose = () => {
    ws = null
    if (convStatus.value !== 'closed') {
      setTimeout(connectWs, 5000)
    }
  }

  ws.onmessage = (e) => {
    try {
      const frame = JSON.parse(e.data)
      handleFrame(frame)
    } catch {/* ignore */}
  }
}

function handleFrame(frame) {
  switch (frame.event) {
    case 'message': {
      const msg = frame.data
      messages.value.push(msg)
      if (!open.value) unread.value++
      scrollBottom()
      break
    }
    case 'revoke': {
      const idx = messages.value.findIndex(m => m.ID === (frame.data?.messageId || frame.messageId))
      if (idx !== -1) {
        messages.value[idx].Revoked = true
        messages.value[idx].Content = '[已撤回]'
      }
      break
    }
    case 'assigned': {
      convStatus.value = 'active'
      messages.value.push({ SenderType: 'system', Content: '客服已接入，请开始咨询', ID: `sys_${Date.now()}` })
      scrollBottom()
      break
    }
    case 'closed': {
      convStatus.value = 'closed'
      messages.value.push({ SenderType: 'system', Content: '会话已结束，感谢您的咨询', ID: `sys_${Date.now()}` })
      scrollBottom()
      break
    }
    case 'queue': {
      const pos = frame.data?.position || frame.position || 0
      queuePosition.value = pos
      convStatus.value = 'pending'
      break
    }
    case 'pong':
      break
  }
}

function sendMsg() {
  const content = inputText.value.trim()
  if (!content || !ws || ws.readyState !== WebSocket.OPEN) return
  clientMsgCounter++
  const clientMsgId = `c_${Date.now()}_${clientMsgCounter}`
  const payload = { event: 'send_message', msgType: 'text', content, clientMsgId }
  ws.send(JSON.stringify(payload))
  // 乐观渲染
  messages.value.push({
    SenderType: 'user',
    Content: content,
    clientMsgId,
    CreatedAt: new Date().toISOString(),
  })
  inputText.value = ''
  scrollBottom()
}

async function submitRating(val) {
  if (!convId.value) return
  await rateConversation({ conversationId: convId.value, rating: val })
  rated.value = true
  ElMessage.success('感谢您的评价！')
}

function scrollBottom() {
  nextTick(() => {
    msgBottom.value?.scrollIntoView({ behavior: 'smooth' })
  })
}

onUnmounted(() => {
  if (ws) {
    ws.onclose = null
    ws.close()
  }
})
</script>

<style scoped>
.cs-widget {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 999;
}

.float-btn {
  width: 54px;
  height: 54px;
  border-radius: 50%;
  background: linear-gradient(135deg, #409eff, #36a3ff);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(64, 158, 255, 0.45);
  transition: transform 0.2s;
}
.float-btn:hover {
  transform: scale(1.08);
}

.chat-window {
  position: absolute;
  bottom: 66px;
  right: 0;
  width: 340px;
  height: 500px;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.chat-header {
  background: linear-gradient(135deg, #409eff, #36a3ff);
  color: #fff;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-weight: 600;
  flex-shrink: 0;
}

.status-bar {
  padding: 6px 12px;
  font-size: 12px;
  text-align: center;
  flex-shrink: 0;
}
.status-bar.pending {
  background: #fdf6ec;
  color: #e6a23c;
}
.status-bar.active {
  background: #f0f9eb;
  color: #67c23a;
}

.msg-box {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  background: #f5f7fa;
}

.msg-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}
.msg-row.mine {
  flex-direction: row-reverse;
}

.sys-msg {
  width: 100%;
  text-align: center;
  font-size: 11px;
  color: #999;
  background: #e8e8e8;
  border-radius: 10px;
  padding: 2px 10px;
  align-self: center;
}

.bubble {
  max-width: 72%;
  padding: 8px 12px;
  border-radius: 16px;
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
}
.bubble-mine {
  background: #409eff;
  color: #fff;
  border-bottom-right-radius: 4px;
}
.bubble-other {
  background: #fff;
  color: #333;
  border-bottom-left-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.rate-bar {
  padding: 8px 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  border-top: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.input-area {
  padding: 10px 12px;
  border-top: 1px solid #f0f0f0;
  flex-shrink: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}
</style>
