import { ref, computed } from 'vue'
import * as agentApi from '@/api/agent'
import type { AgentConversation } from '@/api/agent'
import { message } from '@/utils/message'

export function useConversation() {
  const conversations = ref<AgentConversation[]>([])
  const currentConversation = ref<AgentConversation | null>(null)
  const loading = ref(false)
  const total = ref(0)

  // 加载对话列表
  const loadConversations = async (page = 1, pageSize = 50) => {
    loading.value = true
    try {
      const res = await agentApi.getConversations({ page, page_size: pageSize })
      conversations.value = res.list
      total.value = res.total
    } catch (error: any) {
      message.error(error.response?.data?.message || '加载对话列表失败')
    } finally {
      loading.value = false
    }
  }

  // 创建新对话
  const createConversation = async (title: string) => {
    try {
      const conversation = await agentApi.createConversation({ title })
      conversations.value.unshift(conversation)
      currentConversation.value = conversation
      message.success('创建对话成功')
      return conversation
    } catch (error: any) {
      message.error(error.response?.data?.message || '创建对话失败')
      throw error
    }
  }

  // 选择对话
  const selectConversation = (conversation: AgentConversation) => {
    currentConversation.value = conversation
  }

  // 更新对话标题
  const updateConversationTitle = async (id: string, title: string) => {
    try {
      await agentApi.updateConversation(id, { title })
      const index = conversations.value.findIndex((c) => c.id === id)
      if (index !== -1) {
        conversations.value[index].title = title
      }
      if (currentConversation.value?.id === id) {
        currentConversation.value.title = title
      }
      message.success('更新成功')
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败')
      throw error
    }
  }

  // 删除对话
  const deleteConversation = async (id: string) => {
    try {
      await agentApi.deleteConversation(id)
      conversations.value = conversations.value.filter((c) => c.id !== id)
      if (currentConversation.value?.id === id) {
        currentConversation.value = conversations.value[0] || null
      }
      message.success('删除成功')
    } catch (error: any) {
      message.error(error.response?.data?.message || '删除失败')
      throw error
    }
  }

  const hasConversations = computed(() => conversations.value.length > 0)

  return {
    conversations,
    currentConversation,
    loading,
    total,
    hasConversations,
    loadConversations,
    createConversation,
    selectConversation,
    updateConversationTitle,
    deleteConversation,
  }
}



