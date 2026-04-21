<template>
  <view class="nf-chat">
    <view class="nf-chat-bg"></view>

    <!-- 自定义导航栏 -->
    <view class="nf-navbar">
      <view class="nf-navbar-status"></view>
      <view class="nf-navbar-content">
        <view class="nf-navbar-back" @tap="handleBack">
          <uni-icons type="left" size="20" color="#fff"></uni-icons>
        </view>
        <view class="nf-navbar-center">
          <text class="nf-navbar-title">在线客服</text>
          <view class="nf-conn-dot" :class="wsStatus === 'connected' ? 'nf-dot-on' : 'nf-dot-off'"></view>
        </view>
        <view style="width:64rpx;"></view>
      </view>
    </view>

    <!-- 排队提示 -->
    <view v-if="queuePos > 0 && convStatus === 'pending'" class="nf-queue-bar">
      <uni-icons type="clock" size="14" color="#E50914"></uni-icons>
      <text class="nf-queue-text">前方还有 {{ queuePos }} 人排队，请耐心等待…</text>
    </view>

    <!-- 消息列表 -->
    <scroll-view
      scroll-y
      :scroll-top="scrollTop"
      :scroll-with-animation="true"
      :show-scrollbar="false"
      class="nf-msg-list"
      @scrolltolower="() => {}"
    >
      <view class="nf-msg-padding">
        <!-- 系统消息：会话开始 -->
        <view v-if="convStatus !== 'closed'" class="nf-sys-msg">
          <text>{{ convStatus === 'pending' ? '正在为您分配客服，请稍候…' : '已为您接入人工客服，有任何问题请直接发送' }}</text>
        </view>
        <view v-if="convStatus === 'closed'" class="nf-sys-msg">
          <text>本次会话已结束</text>
        </view>

        <!-- 消息气泡 -->
        <view
          v-for="msg in messages"
          :key="msg.ID || msg.tempId"
          class="nf-msg-row"
          :class="msg.senderType === 'user' ? 'nf-row-right' : 'nf-row-left'"
        >
          <!-- 客服头像 -->
          <view v-if="msg.senderType !== 'user'" class="nf-avatar-wrap">
            <image class="nf-avatar" src="/static/cs-avatar.png" mode="aspectFill" />
          </view>

          <view class="nf-bubble-wrap">
            <view
              class="nf-bubble"
              :class="msg.senderType === 'user' ? 'nf-bubble-user' : 'nf-bubble-agent'"
            >
              <text class="nf-bubble-text" selectable>{{ msg.content }}</text>
            </view>
            <!-- 撤回标记 -->
            <text v-if="msg.revoked" class="nf-revoke-hint">消息已撤回</text>
            <!-- 发送状态 -->
            <view v-if="msg.senderType === 'user'" class="nf-msg-status">
              <uni-icons v-if="msg.status === 'sending'" type="spinner-cycle" size="12" color="#888"></uni-icons>
              <uni-icons v-else-if="msg.status === 'failed'" type="warn" size="12" color="#E50914"></uni-icons>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>

    <!-- 底部输入区 -->
    <view class="nf-input-bar" v-if="convStatus !== 'closed'">
      <input
        v-model="inputText"
        class="nf-input"
        type="text"
        placeholder="有什么可以帮您…"
        :placeholder-style="'color:rgba(255,255,255,0.35)'"
        :disabled="convStatus === 'pending' || wsStatus !== 'connected'"
        @confirm="sendMessage"
        confirm-type="send"
        :maxlength="500"
      />
      <view
        class="nf-send-btn"
        :class="canSend ? 'nf-send-active' : 'nf-send-disabled'"
        @tap="sendMessage"
      >
        <uni-icons type="paperplane-filled" size="20" color="#fff"></uni-icons>
      </view>
    </view>

    <!-- 评分弹层 -->
    <view v-if="showRating" class="nf-rating-overlay" @tap.self="() => {}">
      <view class="nf-rating-card">
        <text class="nf-rating-title">请为本次服务评分</text>
        <view class="nf-stars">
          <view
            v-for="n in 5"
            :key="n"
            class="nf-star"
            @tap="ratingScore = n"
          >
            <uni-icons
              :type="n <= ratingScore ? 'star-filled' : 'star'"
              size="32"
              :color="n <= ratingScore ? '#E50914' : '#666'"
            ></uni-icons>
          </view>
        </view>
        <view class="nf-rating-btns">
          <view class="nf-rating-skip" @tap="closeRating">跳过</view>
          <view class="nf-rating-submit" @tap="submitRating">提交评分</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { onLoad, onUnload } from '@dcloudio/uni-app'
import { rateConversation } from '@/api/kefu.js'

// -------- WS 地址 --------
let wsBase = 'wss://back.mnmovie.icu'
// #ifdef H5
if (process.env.NODE_ENV === 'development') wsBase = 'ws://localhost:8888'
// #endif
// #ifndef H5
if (process.env.NODE_ENV === 'development') wsBase = 'ws://localhost:8888'
// #endif

// -------- 响应式状态 --------
const messages = ref([])
const inputText = ref('')
const scrollTop = ref(0)
const wsStatus = ref('disconnected') // disconnected | connecting | connected
const convStatus = ref('pending')    // pending | active | closed
const queuePos = ref(0)
const showRating = ref(false)
const ratingScore = ref(5)
let socketTask = null
let convId = ref(0)
let tempIdCounter = 0

const canSend = computed(() =>
  inputText.value.trim().length > 0 &&
  convStatus.value === 'active' &&
  wsStatus.value === 'connected'
)

// -------- WebSocket 连接 --------
function connectWs() {
  const token = uni.getStorageSync('x-token')
  if (!token) {
    uni.showToast({ title: '请先登录', icon: 'none' })
    return
  }
  wsStatus.value = 'connecting'
  socketTask = uni.connectSocket({
    url: `${wsBase}/cs/ws?token=${token}`,
    complete: () => {}
  })

  socketTask.onOpen(() => {
    wsStatus.value = 'connected'
  })

  socketTask.onMessage((evt) => {
    try {
      const frame = JSON.parse(evt.data)
      handleWsMessage(frame)
    } catch (e) {
      console.error('[CS WS] parse error', e)
    }
  })

  socketTask.onError((err) => {
    wsStatus.value = 'disconnected'
    console.error('[CS WS] error', err)
  })

  socketTask.onClose(() => {
    wsStatus.value = 'disconnected'
  })
}

// -------- 处理服务器消息 --------
function handleWsMessage(data) {
  const { type } = data
  switch (type) {
    case 'message': {
      // 普通聊天消息
      const m = data.message
      if (!m) break
      // 如果是自己发的，不重复添加（服务端 echo 回来的，找到临时消息更新状态）
      if (m.senderType === 'user') {
        const idx = messages.value.findIndex(x => x.tempId && x.tempId === data.tempId)
        if (idx !== -1) {
          messages.value[idx] = { ...m, status: 'sent' }
        } else {
          messages.value.push({ ...m, status: 'sent' })
        }
      } else {
        messages.value.push({ ...m, status: 'sent' })
      }
      scrollToBottom()
      break
    }
    case 'assigned': {
      convStatus.value = 'active'
      convId.value = data.conversationId || 0
      queuePos.value = 0
      appendSysHint(`客服 ${data.agentName || '在线客服'} 已接入，请问有什么需要帮助？`)
      scrollToBottom()
      break
    }
    case 'queue': {
      queuePos.value = data.position || 0
      break
    }
    case 'revoke': {
      const idx = messages.value.findIndex(x => x.id === data.messageId)
      if (idx !== -1) messages.value[idx].revoked = true
      break
    }
    case 'closed': {
      convStatus.value = 'closed'
      showRating.value = true
      break
    }
    default:
      break
  }
}

function appendSysHint(text) {
  messages.value.push({
    tempId: `sys-${Date.now()}`,
    senderType: 'system',
    content: text,
    status: 'sent'
  })
}

// -------- 发送消息 --------
function sendMessage() {
  const text = inputText.value.trim()
  if (!text || !canSend.value) return

  const tempId = `tmp-${++tempIdCounter}`
  const clientMsgId = `u-${Date.now()}-${tempIdCounter}`
  const tempMsg = {
    tempId,
    clientMsgId,
    senderType: 'user',
    content: text,
    status: 'sending',
    revoked: false
  }
  messages.value.push(tempMsg)
  inputText.value = ''
  scrollToBottom()

  socketTask?.send({
    data: JSON.stringify({
      event: 'send_message',
      msgType: 'text',
      content: text,
      clientMsgId
    }),
    fail: () => {
      const idx = messages.value.findIndex(x => x.tempId === tempId)
      if (idx !== -1) messages.value[idx].status = 'failed'
    }
  })
}

// -------- 滚动到底部 --------
let scrollTimer = null
function scrollToBottom() {
  clearTimeout(scrollTimer)
  scrollTimer = setTimeout(() => {
    scrollTop.value = 99999
  }, 80)
}

// -------- 评分 --------
function closeRating() {
  showRating.value = false
}

async function submitRating() {
  if (convId.value) {
    try {
      await rateConversation({ conversationId: convId.value, rating: ratingScore.value })
    } catch (e) {
      // 静默处理：评分失败不阻塞流程
    }
  }
  showRating.value = false
  uni.showToast({ title: '感谢您的评价', icon: 'success' })
}

// -------- 返回 --------
function handleBack() {
  if (convStatus.value === 'active') {
    uni.showModal({
      title: '提示',
      content: '会话进行中，确认离开？离开后会话将继续保留。',
      success: (res) => {
        if (res.confirm) uni.navigateBack()
      }
    })
  } else {
    uni.navigateBack()
  }
}

// -------- 生命周期 --------
onLoad(() => {
  connectWs()
})

onUnload(() => {
  socketTask?.close({})
  clearTimeout(scrollTimer)
})
</script>

<style lang="scss">
page {
  background-color: #000;
}

.nf-chat {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #000;
  position: relative;
}

.nf-chat-bg {
  position: fixed;
  top: 0; left: 0; right: 0;
  height: 500rpx;
  z-index: 0;
  pointer-events: none;
  background:
    radial-gradient(ellipse at 30% 0%, rgba(229, 9, 20, 0.1) 0%, transparent 60%),
    radial-gradient(ellipse at 70% 10%, rgba(229, 9, 20, 0.07) 0%, transparent 50%);
}

/* ===== 导航栏 ===== */
.nf-navbar {
  background: rgba(0, 0, 0, 0.9);
  backdrop-filter: blur(20px);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2rpx 16rpx rgba(229,9,20,0.08);
}

.nf-navbar-status {
  height: var(--status-bar-height);
}

.nf-navbar-content {
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24rpx;
}

.nf-navbar-back {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255,255,255,0.06);
}

.nf-navbar-center {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.nf-navbar-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #fff;
  letter-spacing: 1rpx;
}

.nf-conn-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
}
.nf-dot-on { background: #4caf50; box-shadow: 0 0 8rpx #4caf50; }
.nf-dot-off { background: #888; }

/* ===== 排队提示 ===== */
.nf-queue-bar {
  display: flex;
  align-items: center;
  gap: 10rpx;
  background: rgba(229,9,20,0.08);
  border-bottom: 1rpx solid rgba(229,9,20,0.2);
  padding: 16rpx 32rpx;
}

.nf-queue-text {
  font-size: 24rpx;
  color: rgba(255,255,255,0.7);
}

/* ===== 消息列表 ===== */
.nf-msg-list {
  flex: 1;
  overflow: hidden;
}

.nf-msg-padding {
  padding: 24rpx 24rpx 32rpx;
}

/* 系统提示 */
.nf-sys-msg {
  text-align: center;
  margin: 20rpx 0 28rpx;
  font-size: 22rpx;
  color: rgba(255,255,255,0.35);
  padding: 8rpx 24rpx;
  background: rgba(255,255,255,0.04);
  border-radius: 24rpx;
  display: inline-block;
  width: 100%;
  box-sizing: border-box;
}

/* 消息行 */
.nf-msg-row {
  display: flex;
  align-items: flex-end;
  margin-bottom: 24rpx;
  gap: 16rpx;
}

.nf-row-right {
  flex-direction: row-reverse;
}

.nf-avatar-wrap {
  flex-shrink: 0;
}

.nf-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  background: #222;
  border: 2rpx solid rgba(229,9,20,0.3);
}

/* 气泡 */
.nf-bubble-wrap {
  max-width: 70%;
  display: flex;
  flex-direction: column;
  gap: 6rpx;
}

.nf-bubble {
  padding: 20rpx 28rpx;
  border-radius: 24rpx;
  word-break: break-all;
}

.nf-bubble-user {
  background: #E50914;
  border-bottom-right-radius: 6rpx;
  box-shadow: 0 4rpx 20rpx rgba(229,9,20,0.35);
}

.nf-bubble-agent {
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.12);
  border-bottom-left-radius: 6rpx;
  backdrop-filter: blur(8px);
}

.nf-bubble-text {
  font-size: 28rpx;
  color: #fff;
  line-height: 1.6;
}

.nf-revoke-hint {
  font-size: 22rpx;
  color: rgba(255,255,255,0.3);
  font-style: italic;
}

.nf-msg-status {
  display: flex;
  justify-content: flex-end;
}

/* ===== 输入区 ===== */
.nf-input-bar {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: rgba(15,15,15,0.96);
  border-top: 1rpx solid rgba(255,255,255,0.08);
  backdrop-filter: blur(20px);
}

.nf-input {
  flex: 1;
  height: 80rpx;
  background: rgba(255,255,255,0.08);
  border: 1rpx solid rgba(255,255,255,0.12);
  border-radius: 40rpx;
  padding: 0 32rpx;
  color: #fff;
  font-size: 28rpx;
}

.nf-send-btn {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nf-send-active {
  background: #E50914;
  box-shadow: 0 4rpx 16rpx rgba(229,9,20,0.5);
}

.nf-send-disabled {
  background: rgba(229,9,20,0.25);
}

/* ===== 评分弹层 ===== */
.nf-rating-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.75);
  z-index: 999;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 60rpx;
}

.nf-rating-card {
  width: 90%;
  background: #1a1a1a;
  border: 1rpx solid rgba(255,255,255,0.12);
  border-radius: 32rpx;
  padding: 48rpx 40rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 36rpx;
}

.nf-rating-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #fff;
}

.nf-stars {
  display: flex;
  gap: 20rpx;
}

.nf-star {
  padding: 8rpx;
}

.nf-rating-btns {
  width: 100%;
  display: flex;
  gap: 20rpx;
}

.nf-rating-skip {
  flex: 1;
  height: 84rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 42rpx;
  background: rgba(255,255,255,0.08);
  color: rgba(255,255,255,0.6);
  font-size: 28rpx;
}

.nf-rating-submit {
  flex: 2;
  height: 84rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 42rpx;
  background: #E50914;
  color: #fff;
  font-size: 28rpx;
  font-weight: 600;
  box-shadow: 0 4rpx 20rpx rgba(229,9,20,0.4);
}
</style>
