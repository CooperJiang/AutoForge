import request from '@/utils/request'

export interface Tool {
  id?: string
  code: string
  name: string
  description: string
  category: string
  version: string
  author: string
  icon: string
  config_schema: string
  ai_callable: boolean
  enabled: boolean
  tags: string
  created_at?: string
  updated_at?: string
}

// 获取工具列表
export const getToolList = async () => {
  const response = await request.get<Tool[]>('/api/v1/tools')
  return response.data
}

// 获取工具详情
export const getToolDetail = async (code: string) => {
  const response = await request.get<Tool>(`/api/v1/tools/${code}`)
  return response.data
}
