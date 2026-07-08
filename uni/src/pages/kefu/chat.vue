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
          <text class="nf-navbar-title">{{ $t('kefuTitle') }}</text>
          <view class="nf-conn-dot" :class="wsStatus === 'connected' ? 'nf-dot-on' : 'nf-dot-off'"></view>
        </view>
        <view style="width:64rpx;"></view>
      </view>
    </view>

    <!-- 排队提示 -->
    <view v-if="queuePos > 0 && convStatus === 'pending'" class="nf-queue-bar">
      <uni-icons type="clock" size="14" color="#E50914"></uni-icons>
      <text class="nf-queue-text">{{ $t('kefuQueueHint').replace('{}', String(queuePos)) }}</text>
    </view>

    <!-- 消息列表 -->
    <scroll-view
      scroll-y
      :scroll-top="scrollTop"
      :scroll-with-animation="true"
      :show-scrollbar="false"
      class="nf-msg-list"
      @scrolltoupper="loadMoreHistory"
    >
      <view v-if="historyLoading" class="nf-history-loading">
        <text class="nf-history-loading-text">{{ $t('loading') }}</text>
      </view>
      <view v-else-if="noMoreHistory" class="nf-history-loading">
        <text class="nf-history-loading-text">{{ $t('kefuNoMore') }}</text>
      </view>
      <view class="nf-msg-padding">
        <!-- 系统消息：会话开始 -->
        <view v-if="convStatus !== 'closed'" class="nf-sys-msg">
          <text>{{ convStatus === 'pending' ? $t('kefuPendingHint') : $t('kefuAssignedHint') }}</text>
        </view>
        <view v-if="convStatus === 'closed'" class="nf-sys-msg">
          <text>{{ $t('kefuClosedHint') }}</text>
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
            <image v-if="resolveAgentAvatarUrl()" class="nf-avatar" :src="resolveAgentAvatarUrl()" mode="aspectFill" />
            <view v-else class="nf-avatar nf-avatar-fallback" :style="{ background: avatarColor(msg.senderNickname || $t('kefuTitle')) }">
              <text class="nf-avatar-fallback-text">{{ avatarInitial(msg.senderNickname || $t('kefuTitle')) }}</text>
            </view>
          </view>

          <view class="nf-bubble-wrap">
            <view
              class="nf-bubble"
              :class="msg.senderType === 'user' ? 'nf-bubble-user' : 'nf-bubble-agent'"
            >
              <!-- 图片消息 -->
              <image
                v-if="msg.msgType === 'image'"
                :src="msg.content"
                class="nf-bubble-image"
                mode="aspectFit"
                @tap="previewImg(msg.content)"
              />
              <!-- 文本/表情消息 -->
              <text v-else class="nf-bubble-text" selectable>{{ msg.content }}</text>
            </view>
            <!-- 撤回标记 -->
            <text v-if="msg.revoked" class="nf-revoke-hint">{{ $t('kefuRevoked') }}</text>
            <!-- 发送状态 -->
            <view v-if="msg.senderType === 'user'" class="nf-msg-status">
              <uni-icons v-if="msg.status === 'sending'" type="spinner-cycle" size="12" color="#888"></uni-icons>
              <uni-icons v-else-if="msg.status === 'failed'" type="warn" size="12" color="#E50914"></uni-icons>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>

    <!-- 表情面板 -->
    <view v-if="showEmojiPanel" class="nf-emoji-panel">
      <scroll-view scroll-x class="nf-emoji-scroll">
        <view class="nf-emoji-row">
          <view
            v-for="emoji in emojiList"
            :key="emoji"
            class="nf-emoji-item"
            @tap="insertEmoji(emoji)"
          >
            <text class="nf-emoji-char">{{ emoji }}</text>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- 底部输入区 -->
    <view class="nf-input-bar" v-if="convStatus !== 'closed'">
      <!-- 表情按鈕 -->
      <view class="nf-tool-btn" @tap="toggleEmojiPanel">
        <text class="nf-tool-icon">{{ showEmojiPanel ? '✕' : '😊' }}</text>
      </view>
      <!-- 图片按鈕 -->
      <view class="nf-tool-btn" @tap="sendImage">
        <uni-icons type="image" size="22" color="rgba(255,255,255,0.6)"></uni-icons>
      </view>
      <input
        v-model="inputText"
        class="nf-input"
        type="text"
        :placeholder="$t('kefuInputPlaceholder')"
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
        <text class="nf-rating-title">{{ $t('kefuRatingTitle') }}</text>
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
          <view class="nf-rating-skip" @tap="closeRating">{{ $t('kefuSkip') }}</view>
          <view class="nf-rating-submit" @tap="submitRating">{{ $t('kefuSubmitRating') }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad, onPullDownRefresh, onUnload } from '@dcloudio/uni-app'
import { getCsConfig, rateConversation, getMessageHistory } from '@/api/kefu.js'
import { useLangStore } from '@/pinia/modules/lang.js'
import { resolveApiMessage } from '@/utils/i18n.js'
import { getUrl } from '@/utils/url.js'

// -------- i18n --------
const langStore = useLangStore()
const $t = computed(() => langStore.$t)

// -------- WS 地址 --------
let wsBase = 'wss://enapi.235235.vip'
// #ifdef H5
if (process.env.NODE_ENV === 'development') wsBase = 'ws://localhost:8888'
// #endif
// #ifndef H5
if (process.env.NODE_ENV === 'development') wsBase = 'ws://localhost:8888'
// #endif

// HTTP 上传基地址（ws -> http）
const httpBase = wsBase.replace(/^wss?/, m => m === 'wss' ? 'https' : 'http')

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

// -------- 历史消息分页 --------
const historyPage = ref(1)
const historyLoading = ref(false)
const noMoreHistory = ref(false)
const shouldCloseOnUnload = ref(false)
const closingBySelf = ref(false)

const agentAvatarUrl = ref('')
const uploadMaxSizeMB = ref(5)
const uploadAllowExt = ref(['jpg', 'jpeg', 'png', 'webp', 'gif'])
const KEFU_CHAT_PREFILL_KEY = 'kefu:chat:prefill'

const canSend = computed(() =>
  inputText.value.trim().length > 0 &&
  convStatus.value === 'active' &&
  wsStatus.value === 'connected'
)

// -------- 表情 --------
const showEmojiPanel = ref(false)
const emojiList = [
  '😄','😂','🥹','😍','😭','😡','🥰','🤩','😎','😅',
  '🙂','🤔','😇','🥳','🙄','🙏','👍','👎','❤️','💚',
  '🔥','✨','🎉','💯','🙌','👏','💛','💫','🌟','💪'
]

const avatarBgPalette = ['#E50914', '#0EA5E9', '#10B981', '#F59E0B', '#6366F1', '#EC4899', '#14B8A6', '#F97316']

function avatarInitial(name) {
  const text = String(name || '').trim()
  if (!text) return 'K'
  return text.charAt(0).toUpperCase()
}

function avatarColor(seed) {
  const text = String(seed || 'kefu')
  let hash = 0
  for (let i = 0; i < text.length; i++) {
    hash = (hash * 31 + text.charCodeAt(i)) >>> 0
  }
  return avatarBgPalette[hash % avatarBgPalette.length]
}

function resolveAgentAvatarUrl() {
  const raw = String(agentAvatarUrl.value || '').trim()
  if (!raw) return ''
  if (/^(https?:)?\/\//i.test(raw) || /^data:/i.test(raw)) return raw
  return getUrl(raw)
}

function parseUploadAllowExt(raw) {
  const list = String(raw || '')
    .split(',')
    .map(ext => ext.trim().toLowerCase().replace(/^\./, ''))
    .filter(Boolean)
  return list.length ? Array.from(new Set(list)) : ['jpg', 'jpeg', 'png', 'webp', 'gif']
}

function safeDecode(val) {
  const raw = String(val || '').trim()
  if (!raw) return ''
  try {
    return decodeURIComponent(raw)
  } catch (e) {
    return raw
  }
}

function buildPaymentDraft(orderID, payMethodLabel) {
  return $t.value('kefuPaymentDraftTemplate')
    .replace('{orderID}', orderID || '-')
    .replace('{payMethod}', payMethodLabel || '-')
}

function resolvePaymentDraft(options = {}) {
  const orderID = safeDecode(options.orderID)
  const payMethodLabel = safeDecode(options.payMethodLabel)
  const queryPrefill = safeDecode(options.prefill)

  if (queryPrefill) {
    return queryPrefill
  }
  if (orderID) {
    return buildPaymentDraft(orderID, payMethodLabel)
  }

  const cachePrefill = String(uni.getStorageSync(KEFU_CHAT_PREFILL_KEY) || '').trim()
  return cachePrefill
}

function applyPaymentDraft(options = {}) {
  const draft = resolvePaymentDraft(options)
  if (draft) {
    inputText.value = draft
  }
  uni.removeStorageSync(KEFU_CHAT_PREFILL_KEY)
}

async function loadChatConfig() {
  try {
    const res = await getCsConfig()
    if (res.code !== 0) return
    agentAvatarUrl.value = String(res.data?.defaultAvatarUrl || '')
    uploadMaxSizeMB.value = Number(res.data?.uploadMaxSizeMB || 5)
    uploadAllowExt.value = parseUploadAllowExt(res.data?.uploadAllowExt)
  } catch (e) {
    // 使用默认值继续流程
  }
}

function toggleEmojiPanel() {
  showEmojiPanel.value = !showEmojiPanel.value
}

function insertEmoji(emoji) {
  inputText.value += emoji
}

// -------- 下拉加载历史消息 --------
async function loadMoreHistory() {
  if (historyLoading.value || noMoreHistory.value || !convId.value) return
  historyLoading.value = true
  const nextPage = historyPage.value + 1
  try {
    const res = await getMessageHistory({ conversationId: convId.value, page: nextPage, pageSize: 10 })
    if (res.code !== 0) {
      noMoreHistory.value = true
      return
    }
    const list = res.data?.list || []
    if (list.length === 0) {
      noMoreHistory.value = true
    } else {
      const older = [...list].reverse()
      messages.value = [...older, ...messages.value]
      historyPage.value = nextPage
    }
  } catch (e) {
    console.error('[CS] loadMoreHistory error', e)
  } finally {
    historyLoading.value = false
  }
}

// -------- WebSocket 连接 --------
function connectWs() {
  const token = uni.getStorageSync('x-token')
  if (!token) {
    uni.showToast({ title: $t.value('loginPrompt'), icon: 'none' })
    return
  }
  wsStatus.value = 'connecting'
  try {
    socketTask = uni.connectSocket({
      url: `${wsBase}/cs/ws`,
      protocols: ['bearer', token],
      header: { Authorization: `Bearer ${token}` },
      complete: () => {}
    })
  } catch {
    socketTask = uni.connectSocket({
      url: `${wsBase}/cs/ws?token=${encodeURIComponent(token)}`,
      complete: () => {}
    })
  }

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
// 服务端帧格式: { event: string, data: object }
function handleWsMessage(frame) {
  const { event, data: payload } = frame
  switch (event) {
    case 'session': {
      // 连接后服务端发送当前会话状态
      convId.value = payload?.ID || 0
      if (payload?.status === 'active') {
        convStatus.value = 'active'
        queuePos.value = 0
      } else if (payload?.status === 'pending') {
        convStatus.value = 'pending'
        queuePos.value = payload?.queuePosition || 0
      }
      break
    }
    case 'history': {
      // 服务端发送最近消息历史
      if (Array.isArray(payload)) {
        messages.value = payload.map(m => ({ ...m, status: 'sent' }))
        historyPage.value = 1
        noMoreHistory.value = false
        scrollToBottom()
      }
      break
    }
    case 'message': {
      // payload 就是 CsMessage 对象
      const m = payload
      if (!m) break
      if (m.senderType === 'user') {
        // 找到对应的临时消息（用 clientMsgId 匹配），替换为服务端确认消息
        const idx = messages.value.findIndex(x => x.clientMsgId && x.clientMsgId === m.clientMsgId)
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
      // payload: { conversationId, agentUserId }
      convStatus.value = 'active'
      if (payload?.conversationId) convId.value = payload.conversationId
      queuePos.value = 0
      appendSysHint($t.value('kefuAssignedHint'))
      scrollToBottom()
      break
    }
    case 'queue': {
      // payload: { conversationId, position }
      queuePos.value = payload?.position || 0
      break
    }
    case 'revoke': {
      // payload: { messageId, conversationId }
      const idx = messages.value.findIndex(x => x.ID === payload?.messageId)
      if (idx !== -1) messages.value[idx].revoked = true
      break
    }
    case 'closed': {
      convStatus.value = 'closed'
      if (!closingBySelf.value) {
        showRating.value = true
      }
      break
    }
    case 'transfer': {
      appendSysHint($t.value('kefuTransferred'))
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

// -------- 图片预览 --------
function previewImg(url) {
  uni.previewImage({ current: url, urls: [url] })
}

// -------- 发送图片 --------
function sendImage() {
  if (convStatus.value !== 'active' || wsStatus.value !== 'connected') {
    uni.showToast({ title: $t.value('kefuWaitAgent'), icon: 'none' })
    return
  }
  const maxSize = Math.max(1, Number(uploadMaxSizeMB.value || 5)) * 1024 * 1024
  const allowExt = uploadAllowExt.value
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempFilePath = res.tempFilePaths[0]
      const size = Array.isArray(res.tempFiles) && res.tempFiles[0] ? res.tempFiles[0].size || 0 : 0
      const lower = String(tempFilePath || '').toLowerCase()
      const ext = lower.includes('.') ? lower.substring(lower.lastIndexOf('.') + 1) : ''
      if (!allowExt.includes(ext)) {
        const extHint = allowExt.map(item => item.toUpperCase()).join('/')
        uni.showToast({ title: $t.value('kefuUploadTypeInvalid').replace('{}', extHint), icon: 'none' })
        return
      }
      if (size > maxSize) {
        uni.showToast({ title: $t.value('kefuUploadSizeInvalid').replace('{}', String(uploadMaxSizeMB.value || 5)), icon: 'none' })
        return
      }
      uploadAndSendImage(tempFilePath)
    }
  })
}

function uploadAndSendImage(tempFilePath) {
  const clientMsgId = `u-img-${Date.now()}`
  const tempId = `tmp-${++tempIdCounter}`
  const tempMsg = {
    tempId,
    clientMsgId,
    senderType: 'user',
    msgType: 'image',
    content: tempFilePath, // 展示临时本地路径
    status: 'sending',
    revoked: false
  }
  messages.value.push(tempMsg)
  scrollToBottom()

  uni.uploadFile({
    url: httpBase + '/cs/message/upload',
    header: { 'x-token': uni.getStorageSync('x-token') },
    filePath: tempFilePath,
    name: 'file',
    success: (res) => {
      try {
        const data = JSON.parse(res.data)
        if (data.code !== 0) {
          const idx = messages.value.findIndex(m => m.tempId === tempId)
          if (idx !== -1) messages.value[idx].status = 'failed'
          uni.showToast({ title: resolveApiMessage(data.msg, 'uploadFail'), icon: 'none' })
          return
        }
        const imageUrl = data.data.file.url
        // 替换临时消息为真实 URL
        const idx = messages.value.findIndex(m => m.tempId === tempId)
        if (idx !== -1) {
          messages.value[idx].content = imageUrl
          messages.value[idx].status = 'sent'
        }
        // 通过 WS 发送图片消息
        socketTask?.send({
          data: JSON.stringify({
            event: 'send_message',
            msgType: 'image',
            content: imageUrl,
            clientMsgId
          }),
          fail: () => {
            const i = messages.value.findIndex(m => m.tempId === tempId)
            if (i !== -1) messages.value[i].status = 'failed'
          }
        })
      } catch (e) {
        const idx = messages.value.findIndex(m => m.tempId === tempId)
        if (idx !== -1) messages.value[idx].status = 'failed'
        uni.showToast({ title: $t.value('uploadFail'), icon: 'none' })
      }
    },
    fail: (err) => {
      const idx = messages.value.findIndex(m => m.tempId === tempId)
      if (idx !== -1) messages.value[idx].status = 'failed'
      uni.showToast({ title: resolveApiMessage(err?.errMsg, 'uploadFail'), icon: 'none' })
    }
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
let bottomSeq = 0
function scrollToBottom() {
  clearTimeout(scrollTimer)
  scrollTimer = setTimeout(() => {
    bottomSeq += 1
    scrollTop.value = bottomSeq * 100000
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
  uni.showToast({ title: $t.value('kefuThanksRating'), icon: 'success' })
}

// -------- 返回 --------
function handleBack() {
  if (convStatus.value === 'active') {
    uni.showModal({
      title: '',
      content: $t.value('kefuLeaveConfirm'),
      confirmText: $t.value('confirm'),
      cancelText: $t.value('cancel'),
      success: (res) => {
        if (res.confirm) {
          shouldCloseOnUnload.value = true
          uni.navigateBack()
        }
      }
    })
  } else {
    shouldCloseOnUnload.value = true
    uni.navigateBack()
  }
}

// -------- 生命周期 --------
onLoad((options = {}) => {
  applyPaymentDraft(options)
  loadChatConfig().finally(() => {
    connectWs()
  })
})

onPullDownRefresh(() => {
  uni.stopPullDownRefresh()
})

onUnload(() => {
  if (socketTask && wsStatus.value === 'connected' && convId.value && convStatus.value !== 'closed' && shouldCloseOnUnload.value) {
    closingBySelf.value = true
    socketTask.send({
      data: JSON.stringify({ event: 'close' }),
      fail: () => {},
    })
  }
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

/* 历史消息加载提示 */
.nf-history-loading {
  text-align: center;
  padding: 16rpx 0;
}
.nf-history-loading-text {
  font-size: 22rpx;
  color: rgba(255,255,255,0.3);
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

.nf-avatar-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
}

.nf-avatar-fallback-text {
  font-size: 28rpx;
  font-weight: 700;
  color: #fff;
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

/* ===== 工具栏按钮 ===== */
.nf-tool-btn {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(255,255,255,0.06);
  flex-shrink: 0;
}

.nf-tool-icon {
  font-size: 36rpx;
  line-height: 1;
}

/* ===== 表情面板 ===== */
.nf-emoji-panel {
  background: rgba(20,20,20,0.98);
  border-top: 1rpx solid rgba(255,255,255,0.08);
  padding: 16rpx 12rpx;
}

.nf-emoji-scroll {
  width: 100%;
  white-space: nowrap;
}

.nf-emoji-row {
  display: flex;
  gap: 4rpx;
}

.nf-emoji-item {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12rpx;
  flex-shrink: 0;
}

.nf-emoji-item:active {
  background: rgba(255,255,255,0.1);
}

.nf-emoji-char {
  font-size: 44rpx;
  line-height: 1;
}

/* ===== 图片气泡 ===== */
.nf-bubble-image {
  width: 360rpx;
  max-height: 480rpx;
  border-radius: 16rpx;
  display: block;
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
