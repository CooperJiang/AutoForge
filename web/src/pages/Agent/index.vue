<template>
  <div class="h-screen flex overflow-hidden bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800">
    <!-- 左侧：对话列表 -->
    <transition name="slide-fade">
      <ConversationList
        v-show="!sidebarCollapsed"
        :conversations="conversations"
        :current-conversation="currentConversation"
        :loading="loading"
        @select="selectConversation"
        @create="handleCreateConversation"
        @delete="handleDeleteConversation"
        @update="handleUpdateConversation"
      />
    </transition>

    <!-- 中间：对话区域 -->
    <div class="flex-1 flex flex-col min-w-0 relative">
      <!-- 顶部栏 - 玻璃态效果 -->
      <div
        class="h-16 backdrop-blur-xl bg-white/70 dark:bg-slate-900/70 border-b border-slate-200/50 dark:border-slate-700/50 flex items-center justify-between px-6 sticky top-0 z-10"
      >
        <div class="flex items-center gap-4">
          <!-- 折叠按钮 -->
          <button
            @click="sidebarCollapsed = !sidebarCollapsed"
            class="p-2.5 rounded-xl hover:bg-slate-200/50 dark:hover:bg-slate-700/50 transition-all duration-200 hover:scale-105"
            title="切换侧边栏"
          >
            <Menu class="w-5 h-5 text-slate-600 dark:text-slate-300" />
          </button>

          <!-- Logo 和标题 -->
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center shadow-lg">
              <Bot class="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 class="text-lg font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                AI Agent
              </h1>
              <p v-if="currentConversation" class="text-xs text-slate-500 dark:text-slate-400">
                {{ currentConversation.title }}
              </p>
            </div>
          </div>
        </div>

        <!-- 右侧操作按钮 -->
        <div class="flex items-center gap-2">
          <button
            @click="showExamples = !showExamples"
            class="px-4 py-2 rounded-xl bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 hover:bg-blue-100 dark:hover:bg-blue-900/50 transition-all duration-200 text-sm font-medium flex items-center gap-2"
          >
            <Lightbulb class="w-4 h-4" />
            示例
          </button>
          <button
            @click="showSettings = !showSettings"
            class="p-2.5 rounded-xl hover:bg-slate-200/50 dark:hover:bg-slate-700/50 transition-all duration-200 hover:scale-105"
            title="设置"
          >
            <Settings class="w-5 h-5 text-slate-600 dark:text-slate-300" />
          </button>
        </div>
      </div>

      <!-- 示例提示词面板 -->
      <transition name="slide-down">
        <div
          v-if="showExamples"
          class="absolute top-16 left-0 right-0 z-20 backdrop-blur-xl bg-white/90 dark:bg-slate-900/90 border-b border-slate-200/50 dark:border-slate-700/50 p-6 shadow-lg"
        >
          <div class="max-w-4xl mx-auto">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-slate-800 dark:text-slate-200">💡 试试这些示例</h3>
              <button
                @click="showExamples = false"
                class="p-1 rounded-lg hover:bg-slate-200/50 dark:hover:bg-slate-700/50"
              >
                <X class="w-5 h-5 text-slate-500" />
              </button>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
              <button
                v-for="(example, index) in examples"
                :key="index"
                @click="useExample(example)"
                class="text-left p-4 rounded-xl bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-800 dark:to-slate-700 hover:from-blue-50 hover:to-purple-50 dark:hover:from-blue-900/20 dark:hover:to-purple-900/20 border border-slate-200 dark:border-slate-600 hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-200 hover:shadow-md group"
              >
                <div class="flex items-start gap-3">
                  <div class="text-2xl">{{ example.icon }}</div>
                  <div class="flex-1">
                    <h4 class="font-medium text-slate-800 dark:text-slate-200 mb-1 group-hover:text-blue-600 dark:group-hover:text-blue-400">
                      {{ example.title }}
                    </h4>
                    <p class="text-sm text-slate-600 dark:text-slate-400 line-clamp-2">
                      {{ example.prompt }}
                    </p>
                  </div>
                </div>
              </button>
            </div>
          </div>
        </div>
      </transition>

      <!-- 聊天区域 -->
      <ChatArea
        :conversation-id="currentConversation?.id"
        :messages="messages"
        :loading="messagesLoading"
        :is-streaming="isStreaming"
        :current-plan="currentPlan"
        :current-steps="currentSteps"
        @send="handleSendMessage"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Menu, Settings, Bot, Lightbulb, X } from 'lucide-vue-next'
import ConversationList from './components/ConversationList.vue'
import ChatArea from './components/ChatArea.vue'
import { useConversation } from './composables/useConversation'
import { useAgentStream } from './composables/useAgentStream'
import * as agentApi from '@/api/agent'
import type { AgentMessage, SendMessageRequest } from '@/api/agent'
import { message } from '@/utils/message'

const sidebarCollapsed = ref(false)
const messages = ref<AgentMessage[]>([])
const messagesLoading = ref(false)
const showExamples = ref(false)
const showSettings = ref(false)

// 示例提示词
const examples = ref([
  {
    icon: '🌐',
    title: '获取网页数据并分析',
    prompt: '请帮我访问 https://httpbin.org/json 获取数据，然后提取其中的 slideshow.title 字段，并用格式化工具输出结果。',
  },
  {
    icon: '📧',
    title: '发送邮件通知',
    prompt: '请给 example@example.com 发送一封邮件，主题是"测试邮件"，内容是"这是一封由 AI Agent 自动发送的测试邮件"。',
  },
  {
    icon: '🔍',
    title: '数据转换与处理',
    prompt: '请访问 https://api.github.com/users/github 获取 GitHub 用户信息，提取 name、bio 和 public_repos 字段，然后格式化输出。',
  },
  {
    icon: '⏰',
    title: '健康检查',
    prompt: '请检查 https://httpbin.org 的健康状态，如果正常就发送一个格式化的成功消息。',
  },
  {
    icon: '🎨',
    title: '多步骤任务',
    prompt: '请执行以下任务：1. 访问 https://httpbin.org/uuid 获取一个 UUID；2. 将这个 UUID 格式化输出；3. 检查 httpbin.org 的健康状态。',
  },
  {
    icon: '🔗',
    title: '链式 API 调用',
    prompt: '请先访问 https://httpbin.org/json 获取数据，然后将获取到的 JSON 数据发送到 https://httpbin.org/post 进行 POST 请求，最后格式化输出响应结果。',
  },
])

// 使用示例
const useExample = (example: any) => {
  showExamples.value = false
  // 触发发送消息
  handleSendMessage({ message: example.prompt })
}

// 对话管理
const {
  conversations,
  currentConversation,
  loading,
  loadConversations,
  createConversation,
  selectConversation: selectConv,
  updateConversationTitle,
  deleteConversation,
} = useConversation()

// 流式响应
const { isStreaming, currentPlan, currentSteps, startStream } = useAgentStream()

// 初始化
onMounted(async () => {
  await loadConversations()

  // 如果有对话，加载第一个
  if (conversations.value.length > 0) {
    selectConversation(conversations.value[0])
  }
})

// 监听当前对话变化，加载消息
watch(
  () => currentConversation.value?.id,
  async (newId) => {
    if (newId) {
      await loadMessages(newId)
    } else {
      messages.value = []
    }
  }
)

// 加载消息列表
const loadMessages = async (conversationId: string) => {
  messagesLoading.value = true
  try {
    const loadedMessages = await agentApi.getMessages(conversationId)
    // 确保消息按创建时间升序排列（旧消息在前）
    messages.value = loadedMessages.sort((a, b) => a.created_at - b.created_at)
  } catch (error: any) {
    message.error(error.response?.data?.message || '加载消息失败')
  } finally {
    messagesLoading.value = false
  }
}

// 选择对话
const selectConversation = (conversation: any) => {
  selectConv(conversation)
}

// 创建对话
const handleCreateConversation = async () => {
  const title = `新对话 ${new Date().toLocaleTimeString()}`
  await createConversation(title)
}

// 删除对话
const handleDeleteConversation = async (id: string) => {
  await deleteConversation(id)
}

// 更新对话
const handleUpdateConversation = async (id: string, title: string) => {
  await updateConversationTitle(id, title)
}

// 发送消息
const handleSendMessage = async (data: SendMessageRequest) => {
  if (!currentConversation.value) {
    // 如果没有当前对话，先创建一个
    const title = data.message.substring(0, 30) + (data.message.length > 30 ? '...' : '')
    await createConversation(title)
  }

  if (!currentConversation.value) {
    message.error('无法创建对话')
    return
  }

  const conversationId = currentConversation.value.id

  // 添加用户消息到列表（临时显示）
  const tempUserMessageId = `temp-user-${Date.now()}`
  const nowTimestamp = Math.floor(Date.now() / 1000) // 转换为秒级时间戳，与后端一致
  const userMessage: AgentMessage = {
    id: tempUserMessageId,
    conversation_id: conversationId,
    role: 'user',
    content: data.message,
    status: 'completed',
    created_at: nowTimestamp,
  }
  messages.value.push(userMessage)

  // 添加一个占位的 Agent 消息
  const tempAgentMessageId = `temp-agent-${Date.now()}`
  const placeholderAgentMessage: AgentMessage = {
    id: tempAgentMessageId,
    conversation_id: conversationId,
    role: 'agent',
    content: '',
    status: 'running',
    created_at: nowTimestamp + 1, // 确保 Agent 消息在用户消息之后
  }
  messages.value.push(placeholderAgentMessage)

  // 启动流式接收
  startStream(conversationId, data.message, data.config || {}, async (agentMessage) => {
    // 流式完成后，重新加载消息列表以获取完整数据（包括真实的 ID）
    await loadMessages(conversationId)
  })
}
</script>

<style scoped>
/* 滑动淡入动画 */
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(-20px);
  opacity: 0;
}

/* 下滑动画 */
.slide-down-enter-active {
  transition: all 0.3s ease-out;
}

.slide-down-leave-active {
  transition: all 0.2s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(-20px);
  opacity: 0;
}

/* 文本截断 */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>



