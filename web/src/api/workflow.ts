/**
 * 工作流 API
 */

import request from '@/utils/request'
import type {
  Workflow,
  WorkflowNode,
  WorkflowEdge,
  WorkflowExecution,
  WorkflowExecutionDetail,
  CreateWorkflowDto,
  UpdateWorkflowDto,
  ExecuteWorkflowDto,
} from '@/types/workflow'

export interface WorkflowListData {
  items: Workflow[]
  total: number
  page: number
  page_size: number
}

export interface ExecutionListData {
  items: WorkflowExecution[]
  total: number
  page: number
  page_size: number
}

export interface ExecuteWorkflowData {
  execution_id: string
  status: string
  message: string
}

export const workflowApi = {
  /**
   * 获取工作流列表
   */
  list: async (params?: {
    page?: number
    page_size?: number
    keyword?: string
    enabled?: boolean
  }) => {
    const response = await request.get<WorkflowListData>('/api/v1/workflows', { params })
    return response.data
  },

  /**
   * 获取工作流详情
   */
  getById: async (id: string) => {
    const response = await request.get<Workflow>(`/api/v1/workflows/${id}`)
    return response.data
  },

  /**
   * 创建工作流
   */
  create: async (data: CreateWorkflowDto) => {
    const response = await request.post<Workflow>('/api/v1/workflows', data)
    return response.data
  },

  /**
   * 更新工作流
   */
  update: async (id: string, data: UpdateWorkflowDto) => {
    const response = await request.put<Workflow>(`/api/v1/workflows/${id}`, data)
    return response.data
  },

  /**
   * 删除工作流
   */
  delete: async (id: string) => {
    const response = await request.delete<void>(`/api/v1/workflows/${id}`)
    return response.data
  },

  /**
   * 批量删除工作流
   */
  batchDelete: async (ids: string[]) => {
    const response = await request.post<void>('/api/v1/workflows/batch-delete', { ids })
    return response.data
  },

  /**
   * 启用/禁用工作流
   */
  toggleEnabled: async (id: string, enabled: boolean) => {
    const response = await request.patch<Workflow>(`/api/v1/workflows/${id}/toggle`, { enabled })
    return response.data
  },

  /**
   * 执行工作流
   */
  execute: async (id: string, data?: ExecuteWorkflowDto) => {
    const response = await request.post<ExecuteWorkflowData>(
      `/api/v1/workflows/${id}/execute`,
      data
    )
    return response.data
  },

  /**
   * 停止工作流执行
   */
  stopExecution: async (id: string, executionId: string) => {
    const response = await request.post<void>(
      `/api/v1/workflows/${id}/executions/${executionId}/stop`
    )
    return response.data
  },

  /**
   * 获取工作流执行历史
   */
  getExecutions: async (
    id: string,
    params?: {
      page?: number
      page_size?: number
      status?: string
      start_time?: number
      end_time?: number
    }
  ) => {
    const response = await request.get<ExecutionListData>(`/api/v1/workflows/${id}/executions`, {
      params,
    })
    return response.data
  },

  /**
   * 获取工作流执行详情
   */
  getExecutionDetail: async (id: string, executionId: string) => {
    const response = await request.get<WorkflowExecution>(
      `/api/v1/workflows/${id}/executions/${executionId}`
    )
    return response.data
  },

  /**
   * 导出工作流为 JSON
   */
  export: async (id: string) => {
    const response = await request.get<Workflow>(
      `/api/v1/workflows/${id}/export`
    )
    return response.data
  },

  /**
   * 导入工作流 JSON
   */
  import: async (data: { name: string; workflow_json: string }) => {
    const response = await request.post<Workflow>('/api/v1/workflows/import', data)
    return response.data
  },

  /**
   * 复制工作流
   */
  clone: async (id: string, name?: string) => {
    const response = await request.post<Workflow>(`/api/v1/workflows/${id}/clone`, { name })
    return response.data
  },

  /**
   * 获取工作流模板列表
   */
  getTemplates: async () => {
    const response = await request.get<Workflow[]>(
      '/api/v1/workflows/templates'
    )
    return response.data
  },

  /**
   * 从模板创建工作流
   */
  createFromTemplate: async (templateId: string, name: string) => {
    const response = await request.post<Workflow>('/api/v1/workflows/from-template', {
      template_id: templateId,
      name,
    })
    return response.data
  },

  /**
   * 验证工作流配置
   */
  validate: async (data: {
    nodes: WorkflowNode[]
    edges: WorkflowEdge[]
    env_vars?: any[]
  }) => {
    const response = await request.post<{
      valid: boolean
      errors?: string[]
    }>('/api/v1/workflows/validate', data)
    return response.data
  },

  /**
   * 获取工作流统计信息
   */
  getStats: async (id: string) => {
    const response = await request.get<{
      total_executions: number
      success_count: number
      failed_count: number
      avg_duration_ms: number
      last_execution_at?: number
    }>(`/api/v1/workflows/${id}/stats`)
    return response.data
  },

  /**
   * 删除执行记录
   */
  deleteExecution: async (workflowId: string, executionId: string) => {
    const response = await request.delete<void>(
      `/api/v1/workflows/${workflowId}/executions/${executionId}`
    )
    return response.data
  },
}

export default workflowApi
