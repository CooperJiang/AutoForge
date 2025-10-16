import request from '@/utils/request'

export interface OutputFieldDef {
  type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  label: string
  children?: Record<string, OutputFieldDef>
}

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
  output_fields_schema?: Record<string, OutputFieldDef>
  created_at?: string
  updated_at?: string
}

export const getToolList = async () => {
  const response = await request.get<Tool[]>('/api/v1/tools')
  return response.data
}

export const getToolDetail = async (code: string) => {
  const response = await request.get<Tool>(`/api/v1/tools/${code}`)
  return response.data
}

export const describeToolOutput = async (code: string, config: Record<string, any>) => {
  const response = await request.post<Record<string, OutputFieldDef>>(
    `/api/v1/tools/${code}/describe-output`,
    config
  )
  return response.data
}
