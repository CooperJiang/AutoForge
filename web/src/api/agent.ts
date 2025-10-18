import request from '@/utils/request'

export interface AgentConversation {
  id: string
  user_id: string
  title: string
  created_at: number
  updated_at: number
}

export interface AgentFile {
  path: string
  filename: string
  size: number
  mime_type: string
}

export interface AgentAction {
  type: string
  tool?: string
  args?: Record<string, any>
}

export interface AgentStep {
  step: number
  action?: AgentAction
  observation?: string
  tool_output?: Record<string, any>
  elapsed_ms: number
  timestamp: string
}

export interface AgentTrace {
  steps: AgentStep[]
  final_answer: string
  finish_reason: string
  used_tools?: Record<string, any>
  total_ms: number
}

export interface AgentPlanStep {
  step: number
  description: string
  tool?: string
  status: 'pending' | 'running' | 'completed' | 'skipped' | 'failed'
}

export interface AgentPlan {
  steps: AgentPlanStep[]
  total_steps: number
  created_at: string
  generated_by: string
}

export interface TokenUsage {
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
}

export interface AgentMessage {
  id: string
  conversation_id: string
  role: 'user' | 'agent' | 'system'
  content: string
  files?: AgentFile[]
  trace?: AgentTrace
  plan?: AgentPlan
  config?: Record<string, any>
  token_usage?: TokenUsage
  status: 'pending' | 'running' | 'completed' | 'failed'
  error?: string
  created_at: number
}

export interface CreateConversationRequest {
  title: string
}

export interface SendMessageRequest {
  message: string
  config?: {
    model?: string
    mode?: 'direct' | 'plan'
    max_steps?: number
    temperature?: number
    allowed_tools?: string[]
  }
}

export interface SendMessageResponse {
  user_message: AgentMessage
  agent_message: AgentMessage
}

// 创建对话
export const createConversation = async (data: CreateConversationRequest) => {
  const response = await request.post<AgentConversation>('/agent/conversations', data)
  return response.data
}

// 获取对话列表
export const getConversations = async (params?: { page?: number; page_size?: number }) => {
  const response = await request.get<{
    list: AgentConversation[]
    total: number
    page: number
    page_size: number
  }>('/agent/conversations', { params })
  return response.data
}

// 获取对话详情
export const getConversationById = async (id: string) => {
  const response = await request.get<AgentConversation>(`/agent/conversations/${id}`)
  return response.data
}

// 更新对话
export const updateConversation = async (id: string, data: { title: string }) => {
  const response = await request.put(`/agent/conversations/${id}`, data)
  return response.data
}

// 删除对话
export const deleteConversation = async (id: string) => {
  const response = await request.delete(`/agent/conversations/${id}`)
  return response.data
}

// 获取消息列表
export const getMessages = async (conversationId: string) => {
  const response = await request.get<AgentMessage[]>(`/agent/conversations/${conversationId}/messages`)
  return response.data
}

// 发送消息（普通响应）
export const sendMessage = async (conversationId: string, data: SendMessageRequest) => {
  const response = await request.post<SendMessageResponse>(`/agent/conversations/${conversationId}/messages`, data)
  return response.data
}

// 发送消息（流式响应）- 返回 EventSource
export const sendMessageStream = (conversationId: string, data: SendMessageRequest) => {
  // 注意：这个函数已废弃，请使用 useAgentStream composable
  // 保留此函数仅为兼容性
  const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
  const url = `${baseURL}/agent/conversations/${conversationId}/messages`

  const eventSource = new EventSource(url, {
    withCredentials: true,
  })

  return eventSource
}



