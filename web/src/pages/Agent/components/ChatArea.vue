<template>
  <div class="flex-1 flex flex-col min-h-0">
    <!-- 消息列表 -->
    <div ref="messagesContainer" class="flex-1 overflow-y-auto p-6">
      <div v-if="!conversationId" class="h-full flex items-center justify-center">
        <div class="text-center">
          <div class="text-6xl mb-4">🤖</div>
          <h2 class="text-2xl font-bold text-text-primary mb-2">AI Agent</h2>
          <p class="text-text-secondary">开始新对话，让 AI 帮你完成任务</p>
        </div>
      </div>

      <div v-else-if="loading" class="text-center text-text-tertiary py-8">加载消息中...</div>

      <div v-else class="max-w-4xl mx-auto space-y-6">
        <MessageItem
          v-for="message in messages"
          :key="message.id"
          :message="message"
          :is-streaming="isStreaming && message.role === 'agent' && (message.status === 'pending' || message.status === 'running')"
          :current-plan="isStreaming && message.role === 'agent' && message.id.startsWith('temp-') ? currentPlan : null"
          :current-steps="isStreaming && message.role === 'agent' && message.id.startsWith('temp-') ? currentSteps : []"
        />
      </div>
    </div>

    <!-- 输入区域 -->
    <InputArea @send="handleSend" :disabled="isStreaming" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import MessageItem from './MessageItem.vue'
import InputArea from './InputArea.vue'
import type { AgentMessage, AgentPlan, AgentStep, SendMessageRequest } from '@/api/agent'

interface Props {
  conversationId?: string
  messages: AgentMessage[]
  loading: boolean
  isStreaming: boolean
  currentPlan: AgentPlan | null
  currentSteps: AgentStep[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  send: [data: SendMessageRequest]
}>()

const messagesContainer = ref<HTMLElement>()

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 监听消息变化，自动滚动
watch(
  () => props.messages.length,
  () => {
    scrollToBottom()
  }
)

// 监听流式步骤变化，自动滚动
watch(
  () => props.currentSteps.length,
  () => {
    scrollToBottom()
  }
)

// 发送消息
const handleSend = (data: SendMessageRequest) => {
  emit('send', data)
  scrollToBottom()
}
</script>



