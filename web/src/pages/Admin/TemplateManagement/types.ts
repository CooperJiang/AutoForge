/**
 * TemplateManagement 类型定义
 */

export interface WorkflowTemplate {
  id: string
  name: string
  description: string
  category: string
  cover_image?: string
  status: 'draft' | 'published'
  featured: boolean
  install_count: number
  view_count: number
  created_at: string
  updated_at: string
}

export interface SelectOption {
  label: string
  value: string
}

