import request from '@/utils/request'

export interface ToolConfig {
  id: number
  tool_code: string
  tool_name: string
  enabled: boolean
  visible: boolean
  is_deprecated: boolean
  config_schema: string
  description: string
  category: string
  version: string
  author: string
  tags: string
  sort_order: number
  last_sync_at: string
  created_at: string
  updated_at: string
}

export interface ToolConfigDetail extends ToolConfig {
  decrypted_config: Record<string, any>
}

export interface UpdateToolConfigRequest {
  config: Record<string, any>
}

export interface UpdateToolSettingsRequest {
  enabled: boolean
  visible: boolean
  sort_order: number
}

export interface ToolCategory {
  code: string
  name: string
  description: string
  icon: string
}

/**
 * 获取所有工具配置（管理端）
 */
export function getAllToolConfigs() {
  return request<ToolConfig[]>({
    url: '/api/v1/admin/tool-configs',
    method: 'GET',
  })
}

/**
 * 获取工具配置详情（含解密后的配置）
 */
export function getToolConfigDetail(toolCode: string) {
  return request<ToolConfigDetail>({
    url: `/api/v1/admin/tool-configs/${toolCode}`,
    method: 'GET',
  })
}

/**
 * 更新工具配置
 */
export function updateToolConfig(toolCode: string, data: UpdateToolConfigRequest) {
  return request({
    url: `/api/v1/admin/tool-configs/${toolCode}`,
    method: 'PUT',
    data,
  })
}

/**
 * 更新工具设置（启用/禁用、可见性、排序）
 */
export function updateToolSettings(toolCode: string, data: UpdateToolSettingsRequest) {
  return request({
    url: `/api/v1/admin/tool-configs/${toolCode}/settings`,
    method: 'PATCH',
    data,
  })
}

/**
 * 删除工具配置
 */
export function deleteToolConfig(id: number) {
  return request({
    url: `/api/v1/admin/tool-configs/${id}`,
    method: 'DELETE',
  })
}

/**
 * 同步工具定义
 */
export function syncTools() {
  return request({
    url: '/api/v1/admin/tool-configs/sync',
    method: 'POST',
  })
}

/**
 * 获取工具分类列表
 */
export function getToolCategories() {
  return request<ToolCategory[]>({
    url: '/api/v1/tool-categories',
    method: 'GET',
  })
}
