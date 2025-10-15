import request from '@/utils/request'

export interface Task {
  id: string
  user_id: string
  name: string
  description: string
  tool_code: string
  config: string
  schedule_type: string
  schedule_value: string
  enabled: boolean
  next_run_time: number | null
  created_at: string
  updated_at: string
}

export interface TaskExecution {
  id: string
  task_id: string
  user_id: string
  status: string
  response_status: number
  response_body: string
  duration_ms: number
  error_message: string
  started_at: number
  completed_at: number
  created_at: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
  request_id?: string
  timestamp: number
}

export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

// 创建任务
export const createTask = async (data: {
  name: string
  description?: string
  tool_code: string
  config: Record<string, any>
  schedule_type: string
  schedule_value: string
}) => {
  // request拦截器提取了data.data到response.data
  const response = await request.post<Task>('/api/v1/tasks', data)
  return response.data
}

// 获取任务列表
export const getTaskList = async (page = 1, pageSize = 20) => {
  // request拦截器提取了data.data到response.data
  const response = await request.get<PaginationResponse<Task>>('/api/v1/tasks', {
    params: { page, page_size: pageSize },
  })
  return response.data
}

// 获取任务详情
export const getTask = async (id: string) => {
  const response = await request.get<Task>(`/api/v1/tasks/${id}`)
  return response.data
}

// 更新任务
export const updateTask = async (
  id: string,
  data: {
    name: string
    description?: string
    tool_code: string
    config: Record<string, any>
    schedule_type: string
    schedule_value: string
  }
) => {
  const response = await request.put<Task>(`/api/v1/tasks/${id}`, data)
  return response.data
}

// 删除任务
export const deleteTask = async (id: string) => {
  const response = await request.delete<null>(`/api/v1/tasks/${id}`)
  return response.data
}

// 启用任务
export const enableTask = async (id: string) => {
  const response = await request.post<null>(`/api/v1/tasks/${id}/enable`)
  return response.data
}

// 禁用任务
export const disableTask = async (id: string) => {
  const response = await request.post<null>(`/api/v1/tasks/${id}/disable`)
  return response.data
}

// 手动触发任务
export const triggerTask = async (id: string) => {
  const response = await request.post<null>(`/api/v1/tasks/${id}/trigger`)
  return response.data
}

// 获取任务执行记录
export const getTaskExecutions = async (taskId: string, page = 1, pageSize = 20) => {
  const response = await request.get<PaginationResponse<TaskExecution>>(
    `/api/v1/tasks/${taskId}/executions`,
    {
      params: { page, page_size: pageSize },
    }
  )
  return response.data
}

// 获取执行记录详情
export const getExecution = async (id: string) => {
  const response = await request.get<TaskExecution>(`/api/v1/executions/${id}`)
  return response.data
}

// 删除执行记录
export const deleteExecution = async (id: string) => {
  const response = await request.delete<null>(`/api/v1/executions/${id}`)
  return response.data
}

// 删除任务的所有执行记录
export const deleteAllExecutions = async (taskId: string) => {
  const response = await request.delete<null>(`/api/v1/tasks/${taskId}/executions`)
  return response.data
}

// 测试任务配置
export interface TestTaskRequest {
  url?: string
  method?: string
  headers?: { key: string; value: string }[]
  params?: { key: string; value: string }[]
  body?: string
  tool_code?: string
  config?: Record<string, any>
}

export interface TestTaskResponse {
  success: boolean
  status_code?: number
  response_body?: string
  duration_ms?: number
  error_message?: string
  output?: Record<string, any>
  message?: string
}

export const testTask = async (data: TestTaskRequest) => {
  const response = await request.post<TestTaskResponse>('/api/v1/tasks/test', data)
  return response.data
}
