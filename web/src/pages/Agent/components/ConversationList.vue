<template>
  <div class="w-80 backdrop-blur-xl bg-white/70 dark:bg-slate-900/70 border-r border-slate-200/50 dark:border-slate-700/50 flex flex-col h-full shadow-xl">
    <!-- 头部 -->
    <div class="p-4 border-b border-slate-200/50 dark:border-slate-700/50">
      <button
        @click="$emit('create')"
        class="w-full flex items-center justify-center gap-2 px-4 py-3 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl hover:from-blue-600 hover:to-purple-700 transition-all duration-200 font-medium shadow-lg hover:shadow-xl hover:scale-[1.02] active:scale-[0.98]"
      >
        <Plus class="w-5 h-5" />
        新建对话
      </button>
    </div>

    <!-- 对话列表 -->
    <div class="flex-1 overflow-y-auto custom-scrollbar">
      <div v-if="loading" class="p-4 text-center text-slate-500 dark:text-slate-400">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-slate-300 border-t-blue-500"></div>
      </div>

      <div v-else-if="conversations.length === 0" class="p-8 text-center">
        <div class="text-4xl mb-3">💬</div>
        <p class="text-slate-500 dark:text-slate-400 text-sm">
          还没有对话<br />点击上方按钮开始
        </p>
      </div>

      <div v-else class="p-3 space-y-2">
        <div
          v-for="conversation in conversations"
          :key="conversation.id"
          @click="$emit('select', conversation)"
          class="group relative p-3 rounded-xl cursor-pointer transition-all duration-200"
          :class="
            currentConversation?.id === conversation.id
              ? 'bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/30 dark:to-purple-900/30 border-l-4 border-blue-500 shadow-md'
              : 'hover:bg-slate-100/50 dark:hover:bg-slate-800/50 hover:shadow-sm'
          "
        >
          <!-- 标题 -->
          <div class="flex items-start justify-between gap-2">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <MessageSquare class="w-4 h-4 text-blue-500 flex-shrink-0" />
                <h3 class="text-sm font-medium text-slate-800 dark:text-slate-200 truncate">
                  {{ conversation.title }}
                </h3>
              </div>
              <p class="text-xs text-slate-500 dark:text-slate-400 flex items-center gap-1">
                <Clock class="w-3 h-3" />
                {{ formatTime(conversation.updated_at) }}
              </p>
            </div>

            <!-- 操作按钮 -->
            <div
              class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <button
                @click.stop="handleEdit(conversation)"
                class="p-1.5 rounded-lg hover:bg-blue-100 dark:hover:bg-blue-900/50 transition-colors"
                title="编辑"
              >
                <Edit2 class="w-3.5 h-3.5 text-slate-600 dark:text-slate-400" />
              </button>
              <button
                @click.stop="handleDelete(conversation)"
                class="p-1.5 rounded-lg hover:bg-red-100 dark:hover:bg-red-900/50 transition-colors"
                title="删除"
              >
                <Trash2 class="w-3.5 h-3.5 text-red-500" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Plus, Edit2, Trash2, MessageSquare, Clock } from 'lucide-vue-next'
import type { AgentConversation } from '@/api/agent'
import { message } from '@/utils/message'

interface Props {
  conversations: AgentConversation[]
  currentConversation: AgentConversation | null
  loading: boolean
}

defineProps<Props>()

const emit = defineEmits<{
  select: [conversation: AgentConversation]
  create: []
  delete: [id: string]
  update: [id: string, title: string]
}>()

// 格式化时间
const formatTime = (timestamp: number) => {
  const now = Date.now() / 1000
  const diff = now - timestamp

  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)} 分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)} 小时前`
  if (diff < 604800) return `${Math.floor(diff / 86400)} 天前`

  const date = new Date(timestamp * 1000)
  return `${date.getMonth() + 1}/${date.getDate()}`
}

// 编辑对话标题
const handleEdit = (conversation: AgentConversation) => {
  const newTitle = prompt('修改对话标题', conversation.title)
  if (newTitle && newTitle.trim()) {
    emit('update', conversation.id, newTitle.trim())
  }
}

// 删除对话
const handleDelete = (conversation: AgentConversation) => {
  if (confirm(`确定要删除对话 "${conversation.title}" 吗？`)) {
    emit('delete', conversation.id)
  }
}
</script>

<style scoped>
/* 自定义滚动条 */
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(148, 163, 184, 0.3);
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(148, 163, 184, 0.5);
}
</style>



