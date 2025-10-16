/**
 * 模板市场 API
 */

import request from '@/utils/request'

export interface TemplateBasicInfo {
  id: string
  name: string
  description: string
  category: string
  cover_image: string
  icon: string
  install_count: number
  view_count: number
  is_official: boolean
  is_featured: boolean
  author_name: string
  required_tools: string[]
  created_at: string
}

export interface TemplateData {
  nodes: any[]
  edges: any[]
  env_vars?: any[]
}

export interface TemplateDetail {
  id: string
  name: string
  description: string
  category: string
  cover_image: string
  icon: string
  install_count: number
  view_count: number
  is_official: boolean
  is_featured: boolean
  author_name: string
  required_tools: string[]
  usage_guide: string
  template_data: TemplateData
  status: string
  created_at: string
  updated_at: string
}

export interface TemplateListData {
  items: TemplateBasicInfo[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface CreateTemplateDto {
  name: string
  description: string
  category: string
  workflow_id: string
  cover_image?: string
  icon?: string
  required_tools?: string[]
  usage_guide?: string
  is_featured?: boolean
}

export interface UpdateTemplateDto {
  name?: string
  description?: string
  category?: string
  cover_image?: string
  icon?: string
  usage_guide?: string
  is_featured?: boolean
  status?: string
}

export interface InstallTemplateDto {
  template_id: string
  workflow_name?: string
  env_vars?: Record<string, string>
}

export interface InstallTemplateResult {
  workflow_id: string
  workflow_name: string
  message: string
}

export const templateApi = {
  /**
   * 获取模板列表
   */
  list: async (params?: {
    page?: number
    page_size?: number
    category?: string
    is_featured?: boolean
    search?: string
  }) => {
    const response = await request.get<TemplateListData>('/api/v1/templates', { params })
    return response.data
  },

  /**
   * 获取模板详情
   */
  getById: async (id: string) => {
    const response = await request.get<TemplateDetail>(`/api/v1/templates/${id}`)
    return response.data
  },

  /**
   * 创建模板（管理员）
   */
  create: async (data: CreateTemplateDto) => {
    const response = await request.post<TemplateDetail>('/api/v1/templates', data)
    return response.data
  },

  /**
   * 更新模板（管理员）
   */
  update: async (id: string, data: UpdateTemplateDto) => {
    const response = await request.put<TemplateDetail>(`/api/v1/templates/${id}`, data)
    return response.data
  },

  /**
   * 删除模板（管理员）
   */
  delete: async (id: string) => {
    const response = await request.delete<void>(`/api/v1/templates/${id}`)
    return response.data
  },

  /**
   * 安装模板
   */
  install: async (data: InstallTemplateDto) => {
    const response = await request.post<InstallTemplateResult>('/api/v1/templates/install', data)
    return response.data
  },

  /**
   * 获取安装历史
   */
  getInstallHistory: async () => {
    const response = await request.get<any[]>('/api/v1/templates/installs')
    return response.data
  },
}
